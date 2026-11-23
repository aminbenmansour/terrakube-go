package clilogin

import (
	"fmt"
	"html"
)

// ConsentPage renders the page that follows the Dex login, where the user
// chooses the token lifetime and whether to authorize the CLI. It is plain
// HTML on purpose: the contract is the form it posts, not the styling.
func ConsentPage(email string, defaultDays, maxDays int, problem string) string {
	notice := ""
	if problem != "" {
		notice = `<p class="error">` + html.EscapeString(problem) + `</p>`
	}
	return `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>Authorize Terraform CLI</title></head>
<body>
<h1>Authorize Terraform CLI</h1>
<p>Signed in as <strong>` + html.EscapeString(email) + `</strong>.</p>
` + notice + `
<form method="post" action="/oauth/consent">
  <p>
    <label>Token name
      <input name="name" maxlength="64" placeholder="terraform login">
    </label>
  </p>
  <p>
    <label>Valid for (days)
      <input name="days" type="number" min="1" max="` + fmt.Sprint(maxDays) + `" value="` + fmt.Sprint(defaultDays) + `" required>
    </label>
  </p>
  <p>
    <button name="decision" value="approve">Authorize</button>
    <button name="decision" value="deny">Deny</button>
  </p>
</form>
</body></html>`
}

// ErrorPage reports a broker failure to the browser.
func ErrorPage(message string) string {
	return `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>Terraform login failed</title></head>
<body><h1>Terraform login failed</h1><p>` + html.EscapeString(message) + `</p></body></html>`
}
