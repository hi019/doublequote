//go:build tools
// +build tools

// From https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package dq

import (
	_ "github.com/google/wire/cmd/wire"
)
