// Package migrations embeds the SQL migrations so every binary can run them
// without needing the source tree at runtime.
package migrations

import "embed"

// Files contains every SQL migration in this directory.
//
//go:embed *.sql
var Files embed.FS
