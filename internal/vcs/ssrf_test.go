package vcs

import (
	"errors"
	"net"
	"net/http"
	"testing"
)

func TestValidateEndpointRefusesInternalAddresses(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		allow    bool
		wantCode int
	}{
		{"loopback v4", "http://127.0.0.1:8080", false, http.StatusUnprocessableEntity},
		{"loopback v6", "http://[::1]:8080", false, http.StatusUnprocessableEntity},
		{"metadata v4", "http://169.254.169.254/latest/meta-data", false, http.StatusUnprocessableEntity},
		{"metadata v6", "http://[fd00:ec2::254]/latest/meta-data", false, http.StatusUnprocessableEntity},
		{"unspecified", "http://0.0.0.0:8080", false, http.StatusUnprocessableEntity},
		{"carrier-grade NAT", "http://100.64.0.10", false, http.StatusUnprocessableEntity},
		{"wrong scheme", "ftp://10.0.0.1/repos", false, http.StatusUnprocessableEntity},
		{"not a URL", "not-a-url", false, http.StatusUnprocessableEntity},
		{"self-hosted private", "http://10.1.2.3/api/v4", false, 0},
		{"home private", "http://192.168.1.10", false, 0},
		{"unique-local v6", "http://[fd12::1]/api/v4", false, 0},
		{"loopback allowed", "http://127.0.0.1:8080", true, 0},
		{"metadata allowed by opt-in", "http://169.254.169.254", true, 0},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateEndpoint(testCase.url, testCase.allow)
			if testCase.wantCode == 0 {
				if err != nil {
					t.Fatalf("error = %v, want nil", err)
				}
				return
			}
			var providerErr *Error
			if !errors.As(err, &providerErr) || providerErr.Status != testCase.wantCode {
				t.Fatalf("error = %v, want status %d", err, testCase.wantCode)
			}
		})
	}
}

func TestBlockedAddressRanges(t *testing.T) {
	for _, raw := range []string{"127.0.0.1", "::1", "169.254.169.254", "fd00:ec2::254", "0.0.0.0", "100.100.100.200", "224.0.0.1"} {
		if address := net.ParseIP(raw); !blockedAddress(address) {
			t.Errorf("%s is not blocked", raw)
		}
	}
	for _, raw := range []string{"10.0.0.1", "172.16.0.1", "192.168.0.1", "fd12::1", "93.184.216.34"} {
		if address := net.ParseIP(raw); blockedAddress(address) {
			t.Errorf("%s is blocked", raw)
		}
	}
}

func TestProviderForRefusesALoopbackAPIURL(t *testing.T) {
	_, err := providerFor(Connection{Type: TypeGitHub, APIURL: "http://127.0.0.1:9", AccessToken: "token"})
	var providerErr *Error
	if !errors.As(err, &providerErr) || providerErr.Status != http.StatusUnprocessableEntity {
		t.Fatalf("error = %v, want a 422 provider error", err)
	}
}
