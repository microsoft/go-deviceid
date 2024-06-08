//go:build darwin

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package deviceid

import (
	"errors"
	"os"
	"path"
)

// Get returns the device id for the current system.
func Get() (string, error) {
	// 1.3 Mac
	// * The folder path will be $HOME\Library\Application Support\Microsoft\DeveloperTools where $HOME is the user's home directory.
	// * The file will be called 'deviceid'
	// * The value should be stored in plain text, UTF-8, and in the format specified in Section 1.
	home := os.Getenv("HOME")

	if home == "" {
		return "", errors.New("environment variable HOME is not set")
	}

	return readWriteDeviceIDFile(
		path.Join(home, `Library/Application Support/Microsoft/DeveloperTools`),
		"HOME",
	)
}
