//go:build !release

package web

import (
	"net/http"
)

// HTTP auto generated
var HTTP http.FileSystem = http.Dir("./webgui")
