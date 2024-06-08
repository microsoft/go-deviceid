//go:build linux

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package deviceid

import (
	"fmt"
	"os"
	"path"
)

// Get returns the device id for the current system.
func Get() (string, error) {
	// 1.2 Linux
	// * The folder path will be <RootPath>/Microsoft/DeveloperTools where <RootPath> is $XDG_CACHE_HOME if it is set and not empty, else use $HOME/.cache.
	// * The file will be called 'deviceid'.
	// * The value should be stored in plain text, UTF-8, and in the format specified in Section 1.
	xdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	home := os.Getenv("HOME")

	const devToolsSubPath = `Microsoft/DeveloperTools`

	switch {
	case xdgCacheHome != "":
		dir := path.Join(xdgCacheHome, devToolsSubPath)
		return readWriteDeviceIDFile(dir, "XDG_CACHE_HOME")
	case home != "":
		dir := path.Join(home, ".cache", devToolsSubPath)
		return readWriteDeviceIDFile(dir, "HOME")
	default:
		return "", fmt.Errorf("neither XDG_CACHE_HOME or HOME are set")
	}
}
