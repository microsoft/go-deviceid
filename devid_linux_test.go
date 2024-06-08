//go:build linux

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

// 1.2 Linux
// * The folder path will be <RootPath>/Microsoft/DeveloperTools where <RootPath> is $XDG_CACHE_HOME if it is set and not empty, else use $HOME/.cache.
// * The file will be called 'deviceid'.
// * The value should be stored in plain text, UTF-8, and in the format specified in Section 1.

func TestGet_Linux(t *testing.T) {
	t.Run("XDG_CACHE_HOME", func(t *testing.T) {
		xdgDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "xdg")
		homeDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "home")

		defer func() {
			err := os.RemoveAll(xdgDir)
			require.NoError(t, err)

			err = os.RemoveAll(homeDir)
			require.NoError(t, err)
		}()

		// xdg takes precedence
		t.Setenv("XDG_CACHE_HOME", xdgDir)
		t.Setenv("HOME", homeDir)

		deviceID, err := Get()
		require.NoError(t, err)
		requireValidGUID(t, deviceID)

		// validate it went to the right spot
		bytes, err := os.ReadFile(path.Join(xdgDir, "Microsoft/DeveloperTools", "deviceid"))
		require.NoError(t, err)
		require.Equal(t, deviceID, string(bytes))

		// and nothing ended up in HOME
		_, err = os.Stat(homeDir)
		require.True(t, os.IsNotExist(err), "pseudo-HOME folder isn't created because we never wrote to it")
	})

	t.Run("HOME", func(t *testing.T) {
		homeDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "home")

		t.Setenv("XDG_CACHE_HOME", "") // when empty we default to HOME
		t.Setenv("HOME", homeDir)

		deviceID, err := Get()
		require.NoError(t, err)
		requireValidGUID(t, deviceID)

		// validate it went to the right spot
		bytes, err := os.ReadFile(path.Join(homeDir, ".cache", "Microsoft/DeveloperTools", "deviceid"))
		require.NoError(t, err)
		require.Equal(t, deviceID, string(bytes))
	})

	t.Run("NeitherAreSet", func(t *testing.T) {
		// I can't imagine this happening, but we also don't want to be unpredictable.
		t.Setenv("XDG_CACHE_HOME", "")
		t.Setenv("HOME", "")

		deviceID, err := Get()
		require.Empty(t, deviceID)
		require.EqualError(t, err, "neither XDG_CACHE_HOME or HOME are set")
	})
}
