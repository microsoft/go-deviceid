//go:build windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package deviceid

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

// Get returns the device id for the current system.
func Get() (string, error) {
	return readWriteDeviceIDRegistry(devToolsSubPath)
}

const devToolsSubPath = `SOFTWARE\Microsoft\DeveloperTools`
const deviceIDValueName = "deviceid"

// readWriteDeviceIDRegistry reads a deviceid from a registry key value in subKeyPath + "\deviceid" and returns it.
// If the value doesn't exist, it creates the registry key value with a newly generated device id
// and returns the new deviceid.
func readWriteDeviceIDRegistry(subKeyPath string) (string, error) {
	// 1.1 Windows
	// * The value is cached in the 64-bit Windows Registry under HKeyCurrentUser\SOFTWARE\Microsoft\DeveloperTools.
	// * The key should be named 'deviceid' and should be of type REG_SZ (String value).
	// * The value should be stored in plain text and in the format specified in Section 1.
	key, _, err := registry.CreateKey(registry.CURRENT_USER, subKeyPath, registry.READ|registry.WRITE)

	if err != nil {
		return "", newError(err)
	}

	defer key.Close()

	value, _, err := key.GetStringValue(deviceIDValueName)

	if err == nil {
		return value, nil
	}

	if !errors.Is(err, registry.ErrNotExist) {
		return "", newError(err)
	}

	devID, err := generateDeviceID()

	if err != nil {
		return "", newError(err)
	}

	if err := key.SetStringValue(deviceIDValueName, devID); err != nil {
		return "", newError(err)
	}

	return devID, nil
}
