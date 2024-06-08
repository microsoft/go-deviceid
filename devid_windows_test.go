//go:build windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package deviceid

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

func TestGet_Windows(t *testing.T) {
	now := time.Now().UnixNano()
	root := fmt.Sprintf(`DevIDTestsGo\%d\`, now)
	subPath := fmt.Sprintf(root + devToolsSubPath)

	defer func() {
		require.True(t, strings.HasPrefix(root, `DevIDTestsGo\`))

		// `registry` doesn't export RegDeleteTree or SHDeleteKey but we don't
		// actually need it for production.
		cmd := exec.Command("reg.exe", "delete", "HKCU\\"+root, "/f")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		require.NoError(t, err)
	}()

	var devID string

	t.Run("FromZero", func(t *testing.T) {
		// when we have to create the key and the value
		tmpDevID, err := readWriteDeviceIDRegistry(subPath)
		require.NoError(t, err)
		requireValidGUID(t, tmpDevID)
		requireRegID(t, subPath, tmpDevID)

		devID = tmpDevID
	})

	t.Run("ValueAlreadyExists", func(t *testing.T) {
		// value already exists, should return the same one.
		cachedDevID, err := readWriteDeviceIDRegistry(subPath)
		require.NoError(t, err)
		require.Equal(t, devID, cachedDevID)
	})

	// slight variation - precreate the key, but not the value.
	t.Run("KeyExistsButNoValue", func(t *testing.T) {
		err := registry.DeleteKey(registry.CURRENT_USER, subPath)
		require.NoError(t, err)

		key, existing, err := registry.CreateKey(registry.CURRENT_USER, subPath, registry.ALL_ACCESS)
		require.NoError(t, err)
		require.False(t, existing)

		err = key.Close()
		require.NoError(t, err)

		newDevID, err := readWriteDeviceIDRegistry(subPath)
		require.NoError(t, err)
		requireValidGUID(t, newDevID)
	})
}

func requireRegID(t *testing.T, subPath string, expectedID string) {
	// manual check that it's there.
	key, err := registry.OpenKey(registry.CURRENT_USER, subPath, registry.ALL_ACCESS)
	require.NoError(t, err)

	defer key.Close()

	actualDevID, _, err := key.GetStringValue("deviceid")
	require.NoError(t, err)

	require.Equal(t, actualDevID, expectedID)
}
