//go:build !release

// Package web - assets_dev.go provides development-time web asset loading from disk.
package web

import (
	"net/http"
)

// HTTP auto generated.
var HTTP http.FileSystem = http.Dir("./webgui")
