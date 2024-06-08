//go:build darwin

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package deviceid

import (
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGet_darwin(t *testing.T) {
	t.Run("HOME", func(t *testing.T) {
		homeDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "home")

		t.Setenv("HOME", homeDir)

		deviceID, err := Get()
		require.NoError(t, err)
		requireValidGUID(t, deviceID)

		// validate it went to the right spot
		bytes, err := os.ReadFile(path.Join(homeDir, "Library/Application Support/Microsoft/DeveloperTools", "deviceid"))
		require.NoError(t, err)
		require.Equal(t, deviceID, string(bytes))
	})

	t.Run("HOME is not set", func(t *testing.T) {
		t.Setenv("HOME", "")
		deviceID, err := Get()
		require.Empty(t, deviceID)
		require.EqualError(t, err, "environment variable HOME is not set")
	})
}
