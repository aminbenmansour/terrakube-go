package tfe

import (
	"encoding/json"
	"net/http"

	"github.com/aminbenmansour/terrakube-go/internal/clilogin"
)

// Discovery returns the /.well-known/terraform.json document. The state.v2
// entry intentionally points at /remote/state/v2/ to match upstream even
// though the CLI state backend is mounted at /tfstate/v1.
func Discovery(deps Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		var login map[string]any
		if deps.TerraformLogin {
			// Upstream advertises the broker's own client and the whole
			// 10000-10010 loopback range, not the Dex client and two ports.
			login = map[string]any{
				"client":      clilogin.ClientID,
				"grant_types": []string{"authz_code"},
				"authz":       deps.PublicURL + "/oauth/authorize",
				"token":       deps.PublicURL + "/oauth/token",
				"ports":       []int{clilogin.PortLow, clilogin.PortHigh},
			}
		} else {
			login = map[string]any{
				"client": deps.DexClientID,
				"grant_types": []string{
					"authz_code", "openid", "profile", "email", "offline_access", "groups",
				},
				"authz": deps.DexIssuerURL + "/auth?scope=openid+profile+email+offline_access+groups",
				"token": deps.DexIssuerURL + "/token",
				"ports": []int{10000, 10001},
			}
		}

		body := map[string]any{
			"login.v1": login,
			"state.v2": "/remote/state/v2/",
			"tfe.v2":   "/remote/tfe/v2/",
			"tfe.v2.1": "/remote/tfe/v2/",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(body)
	}
}
