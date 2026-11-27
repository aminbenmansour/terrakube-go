package notification

import (
	"net"
	"net/url"
	"strings"
)

// validateDestination refuses to let a notification configuration be used as
// an SSRF pivot: destinationUrl is set by an admin, but a careless or
// compromised admin account could point it at cloud metadata or internal
// infrastructure. The hostname is resolved and private/reserved, loopback and
// link-local answers are rejected. It runs on every send so a DNS record that
// later starts resolving internally is caught.
//
// shortcut: this is resolve-then-connect, not a pinned-IP dial, so a
// same-millisecond DNS-rebinding race is not closed; pin the IP in the dialer
// if that becomes a requirement.
func validateDestination(label, rawURL string, blockPrivate bool) error {
	u, err := url.Parse(rawURL)
	if err != nil || u.Hostname() == "" {
		return blockedDelivery(label, "destination URL is not a valid URI")
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return blockedDelivery(label, "destination URL must use http or https")
	}
	if !blockPrivate {
		return nil
	}
	addresses, err := net.LookupIP(u.Hostname())
	if err != nil || len(addresses) == 0 {
		return blockedDelivery(label, "destination host could not be resolved")
	}
	for _, address := range addresses {
		if disallowedDestination(address) {
			return blockedDelivery(label, "destination resolves to a private/reserved address")
		}
	}
	return nil
}

func disallowedDestination(address net.IP) bool {
	if address.IsLoopback() || address.IsLinkLocalUnicast() || address.IsLinkLocalMulticast() ||
		address.IsPrivate() || address.IsUnspecified() || address.IsMulticast() {
		return true
	}
	// 100.64.0.0/10 is carrier-grade NAT shared space; no net.IP helper covers it.
	if v4 := address.To4(); v4 != nil && v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 {
		return true
	}
	return false
}

func blockedDelivery(label, reason string) *DeliveryError {
	return &DeliveryError{Message: label + " delivery blocked: " + reason}
}
