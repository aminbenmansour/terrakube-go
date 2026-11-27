package vcs

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// validateEndpoint refuses to fetch an internal address on an admin's behalf.
// The API URL of a VCS connection is admin-controlled, but a careless or
// compromised admin could point it at cloud metadata or a loopback service.
// Private ranges (RFC 1918 and IPv6 unique-local) stay allowed because
// self-hosted GitLab and GitHub Enterprise commonly live there; loopback,
// link-local (where cloud metadata lives), CGNAT and other reserved ranges are
// refused unless AllowLocal is set for a development stack.
//
// shortcut: resolve-then-connect, not a pinned-IP dial, so a same-millisecond
// DNS-rebinding race is not closed; pin the IP in the dialer if that becomes a
// requirement.
func validateEndpoint(rawURL string, allowLocal bool) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return failed(http.StatusUnprocessableEntity, "The VCS API URL is not a valid URL")
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
	default:
		return failed(http.StatusUnprocessableEntity, "The VCS API URL must use http or https")
	}
	if allowLocal {
		return nil
	}

	addresses, err := net.LookupIP(parsed.Hostname())
	if err != nil || len(addresses) == 0 {
		return failed(http.StatusBadGateway, "The VCS API host could not be resolved")
	}
	for _, address := range addresses {
		if blockedAddress(address) {
			return failed(http.StatusUnprocessableEntity,
				"The VCS API URL resolves to a loopback, link-local or reserved address")
		}
	}
	return nil
}

func blockedAddress(address net.IP) bool {
	if address.IsLoopback() || address.IsLinkLocalUnicast() || address.IsLinkLocalMulticast() ||
		address.IsUnspecified() || address.IsMulticast() {
		return true
	}
	// AWS and GCP publish their IPv6 metadata endpoint in unique-local space,
	// which is otherwise allowed for a self-hosted provider.
	if address.Equal(net.ParseIP("fd00:ec2::254")) {
		return true
	}
	// 100.64.0.0/10 is carrier-grade NAT shared space; no net.IP helper covers it.
	if v4 := address.To4(); v4 != nil && v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 {
		return true
	}
	return false
}
