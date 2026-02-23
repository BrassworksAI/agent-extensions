package embedded

import (
	"embed"
	"io/fs"
)

// Content embeds tools.yaml and repository/ for distribution in the binary.
// Run `go generate ./...` or `sh scripts/embed.sh` before building to update.
//
//go:generate sh ../../scripts/embed.sh
//go:embed content
var Content embed.FS

// FS returns the embedded filesystem rooted at "content"
func FS() (fs.FS, error) {
	return fs.Sub(Content, "content")
}
