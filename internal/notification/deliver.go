package notification

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Configuration is what a sender needs from a notification configuration. The
// store decrypts SigningSecret before it reaches this package.
type Configuration struct {
	ChannelType    string
	DestinationURL string
	SigningSecret  *string
}

// DeliveryError classifies one failed delivery. Retryable is false when the
// destination rejected the request (a 4xx other than 429) and true for 429,
// 5xx and network failures. RetryAfter is non-zero only when the destination
// sent a usable Retry-After header.
type DeliveryError struct {
	Message    string
	Retryable  bool
	RetryAfter time.Duration
	Err        error
}

func (e *DeliveryError) Error() string { return e.Message }
func (e *DeliveryError) Unwrap() error { return e.Err }

// Deliverer posts rendered payloads.
type Deliverer struct {
	client       *http.Client
	blockPrivate bool
}

// NewDeliverer builds the client upstream's RestTemplateFactory uses: 10s to
// connect, 30s overall, and no redirects. Redirects stay disabled because the
// SSRF check only sees the configured host; following a 3xx from an allowed
// host could reach an internal address after validation passed.
func NewDeliverer(blockPrivate bool) *Deliverer {
	return &Deliverer{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			},
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		blockPrivate: blockPrivate,
	}
}

// Deliver sends payload to the configured destination.
func (d *Deliverer) Deliver(ctx context.Context, cfg Configuration, payload string) error {
	label := channelLabel(cfg.ChannelType)
	if err := validateDestination(label, cfg.DestinationURL, d.blockPrivate); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.DestinationURL, strings.NewReader(payload))
	if err != nil {
		return &DeliveryError{Message: label + " destination URL is invalid", Err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.ChannelType == ChannelWebhook && cfg.SigningSecret != nil && notBlank(*cfg.SigningSecret) {
		req.Header.Set("X-Terrakube-Signature", Sign(*cfg.SigningSecret, payload))
	}

	resp, err := d.client.Do(req)
	if err != nil {
		// Only the cause, never err.Error(): a Slack webhook URL carries its
		// token in the path and must not appear in logs or errors.
		return &DeliveryError{
			Message:   "Failed to deliver " + label + " notification: " + networkCause(err),
			Retryable: true,
			Err:       err,
		}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Incoming Webhooks report payload/config errors as HTTP 200 with a
		// plain-text body; only the literal "ok" is success.
		if cfg.ChannelType == ChannelSlack && strings.TrimSpace(string(body)) != "ok" {
			return &DeliveryError{
				Message: label + " endpoint rejected the payload: " + strings.TrimSpace(string(body)),
			}
		}
		return nil
	}

	deliveryErr := &DeliveryError{Message: fmt.Sprintf("%s endpoint returned status %d", label, resp.StatusCode)}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		deliveryErr.Retryable = true
		deliveryErr.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"))
	}
	return deliveryErr
}

func networkCause(err error) string {
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Err != nil {
		return urlErr.Err.Error()
	}
	return err.Error()
}

func parseRetryAfter(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		if seconds < 0 {
			return 0
		}
		return time.Duration(seconds) * time.Second
	}
	target, err := http.ParseTime(value)
	if err != nil {
		return 0
	}
	if delay := time.Until(target); delay > 0 {
		return delay
	}
	return 0
}

func channelLabel(channelType string) string {
	switch channelType {
	case ChannelSlack:
		return "Slack"
	case ChannelTeams:
		return "Teams"
	case ChannelWebhook:
		return "Webhook"
	default:
		return channelType
	}
}
