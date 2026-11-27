package notification

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testDeliverer() *Deliverer { return NewDeliverer(false) }

func deliveryError(t *testing.T, err error) *DeliveryError {
	t.Helper()
	var delivered *DeliveryError
	if !errors.As(err, &delivered) {
		t.Fatalf("error = %v, want *DeliveryError", err)
	}
	return delivered
}

func TestDeliverPostsPayload(t *testing.T) {
	for _, channel := range []string{ChannelSlack, ChannelTeams, ChannelWebhook} {
		t.Run(channel, func(t *testing.T) {
			var body, contentType, method string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				raw, _ := io.ReadAll(r.Body)
				body, contentType, method = string(raw), r.Header.Get("Content-Type"), r.Method
				if channel == ChannelSlack {
					_, _ = w.Write([]byte("ok"))
				}
			}))
			defer server.Close()

			payload := `{"a":1}`
			err := testDeliverer().Deliver(context.Background(), Configuration{
				ChannelType: channel, DestinationURL: server.URL + "/hook",
			}, payload)
			if err != nil {
				t.Fatalf("Deliver() error = %v", err)
			}
			if method != http.MethodPost {
				t.Fatalf("method = %q, want POST", method)
			}
			if contentType != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", contentType)
			}
			if body != payload {
				t.Fatalf("body = %q, want %q", body, payload)
			}
		})
	}
}

func TestDeliverSlackRejectsA200ErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("no_text"))
	}))
	defer server.Close()

	err := testDeliverer().Deliver(context.Background(), Configuration{
		ChannelType: ChannelSlack, DestinationURL: server.URL + "/services/X/Y/Z",
	}, `{}`)
	delivered := deliveryError(t, err)
	if delivered.Retryable {
		t.Fatal("a rejected payload is retryable, want terminal")
	}
	if !strings.Contains(delivered.Error(), "no_text") {
		t.Fatalf("error = %q, want the Slack body", delivered.Error())
	}
}

func TestDeliverStatusClassification(t *testing.T) {
	tests := []struct {
		status     int
		retryable  bool
		retryAfter time.Duration
	}{
		{http.StatusBadRequest, false, 0},
		{http.StatusUnauthorized, false, 0},
		{http.StatusForbidden, false, 0},
		{http.StatusNotFound, false, 0},
		{http.StatusGone, false, 0},
		{http.StatusTooManyRequests, true, 5 * time.Second},
		{http.StatusInternalServerError, true, 0},
		{http.StatusBadGateway, true, 0},
		{http.StatusServiceUnavailable, true, 0},
	}
	for _, test := range tests {
		t.Run(strconv.Itoa(test.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.status == http.StatusTooManyRequests {
					w.Header().Set("Retry-After", "5")
				}
				w.WriteHeader(test.status)
			}))
			defer server.Close()

			err := testDeliverer().Deliver(context.Background(), Configuration{
				ChannelType: ChannelWebhook, DestinationURL: server.URL + "/hook",
			}, `{}`)
			delivered := deliveryError(t, err)
			if delivered.Retryable != test.retryable {
				t.Fatalf("status %d retryable = %v, want %v", test.status, delivered.Retryable, test.retryable)
			}
			if delivered.RetryAfter != test.retryAfter {
				t.Fatalf("status %d retryAfter = %v, want %v", test.status, delivered.RetryAfter, test.retryAfter)
			}
		})
	}
}

func TestDeliverParsesHTTPDateRetryAfter(t *testing.T) {
	target := time.Now().UTC().Add(90 * time.Second).Truncate(time.Second)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", target.Format(http.TimeFormat))
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	err := testDeliverer().Deliver(context.Background(), Configuration{
		ChannelType: ChannelTeams, DestinationURL: server.URL + "/hook",
	}, `{}`)
	delivered := deliveryError(t, err)
	if !delivered.Retryable || delivered.RetryAfter < 60*time.Second || delivered.RetryAfter > 90*time.Second {
		t.Fatalf("retryAfter = %v, want about 90s", delivered.RetryAfter)
	}
}

func TestDeliverIgnoresMalformedRetryAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "not-a-valid-value")
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	err := testDeliverer().Deliver(context.Background(), Configuration{
		ChannelType: ChannelWebhook, DestinationURL: server.URL + "/hook",
	}, `{}`)
	delivered := deliveryError(t, err)
	if !delivered.Retryable || delivered.RetryAfter != 0 {
		t.Fatalf("retryable = %v, retryAfter = %v, want retryable without a hint",
			delivered.Retryable, delivered.RetryAfter)
	}
}

func TestWebhookSignature(t *testing.T) {
	payload := `{"a":1}`
	for _, secret := range []*string{nil, stringPtr(""), stringPtr("my-secret")} {
		name := "unsigned"
		if secret != nil && *secret != "" {
			name = "signed"
		}
		t.Run(name, func(t *testing.T) {
			var signature string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				signature = r.Header.Get("X-Terrakube-Signature")
			}))
			defer server.Close()

			err := testDeliverer().Deliver(context.Background(), Configuration{
				ChannelType: ChannelWebhook, DestinationURL: server.URL + "/hook", SigningSecret: secret,
			}, payload)
			if err != nil {
				t.Fatalf("Deliver() error = %v", err)
			}
			if secret == nil || *secret == "" {
				if signature != "" {
					t.Fatalf("signature = %q, want no header", signature)
				}
				return
			}
			if signature != Sign(*secret, payload) {
				t.Fatalf("signature = %q, want %q", signature, Sign(*secret, payload))
			}
			if signature == Sign("wrong-secret", payload) {
				t.Fatal("signature verifies with the wrong secret")
			}
		})
	}
}

func TestDeliverDoesNotFollowRedirects(t *testing.T) {
	var targetHits atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetHits.Add(1)
	}))
	defer target.Close()
	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/internal", http.StatusFound)
	}))
	defer redirector.Close()

	err := testDeliverer().Deliver(context.Background(), Configuration{
		ChannelType: ChannelWebhook, DestinationURL: redirector.URL + "/hook",
	}, `{}`)
	if delivered := deliveryError(t, err); delivered.Retryable {
		t.Fatal("a redirect is retryable, want terminal")
	}
	if targetHits.Load() != 0 {
		t.Fatal("Deliver() followed a redirect")
	}
}

func TestDeliverNetworkErrorIsRetryableAndHidesTheURL(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := listener.Addr().String()
	_ = listener.Close()

	destination := "http://" + address + "/hook?token=supersecret"
	err = testDeliverer().Deliver(context.Background(), Configuration{
		ChannelType: ChannelWebhook, DestinationURL: destination,
	}, `{}`)
	delivered := deliveryError(t, err)
	if !delivered.Retryable {
		t.Fatal("a network error is terminal, want retryable")
	}
	if strings.Contains(delivered.Error(), "supersecret") {
		t.Fatalf("error leaks the destination URL: %q", delivered.Error())
	}
}

func TestValidateDestinationBlocksPrivateAndReservedAddresses(t *testing.T) {
	blocked := []string{
		"http://127.0.0.1/hook",
		"http://localhost/hook",
		"http://[::1]/hook",
		"http://10.0.0.5/hook",
		"http://172.16.4.4/hook",
		"http://192.168.1.1/hook",
		"http://169.254.169.254/latest/meta-data/",
		"http://100.64.0.1/hook",
		"http://0.0.0.0/hook",
		"http://[fd00::1]/hook",
	}
	for _, destination := range blocked {
		t.Run(destination, func(t *testing.T) {
			err := validateDestination("Webhook", destination, true)
			if delivered := deliveryError(t, err); delivered.Retryable {
				t.Fatalf("blocked destination is retryable: %v", delivered)
			}
		})
	}
}

func TestValidateDestinationAllowsPublicAndPermissiveDestinations(t *testing.T) {
	if err := validateDestination("Webhook", "http://8.8.8.8/hook", true); err != nil {
		t.Fatalf("public destination rejected: %v", err)
	}
	if err := validateDestination("Webhook", "http://127.0.0.1:8080/hook", false); err != nil {
		t.Fatalf("permissive loopback rejected: %v", err)
	}
}

func TestValidateDestinationRejectsBadURLs(t *testing.T) {
	for _, destination := range []string{"file:///etc/passwd", "not a url", "", "http:///hook"} {
		if err := validateDestination("Webhook", destination, false); err == nil {
			t.Fatalf("validateDestination(%q) = nil, want an error", destination)
		}
	}
}

func stringPtr(value string) *string { return &value }
