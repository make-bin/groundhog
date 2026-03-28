// Package web provides embedded frontend static assets.
// Build frontend first with: make frontend-build
package web

import "embed"

// FS contains the embedded frontend dist files.
//
//go:embed all:app/dist
var FS embed.FS
