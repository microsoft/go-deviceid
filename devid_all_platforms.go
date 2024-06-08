// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package deviceid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// generateDeviceID generates values in the format of:
// `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
// Where 'x' is any legal lowercased hex digit.
func generateDeviceID() (string, error) {
	randBytes := make([]byte, 4+2+2+2+6)

	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	return formatGUID(randBytes), nil
}

// formatGUID takes 16 bytes and formats it into a lowercased GUID (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
// NOTE: there's no error checking here, I just split this out from generateDeviceID so we can unit test it.
func formatGUID(randBytes []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		randBytes[0:4],
		randBytes[4:6],
		randBytes[6:8],
		randBytes[8:10],
		randBytes[10:])
}

// readWriteDeviceIDFile reads a deviceid from a file in dir + "/deviceid" and returns it.
// If the file doesn't exist it creates the file with a newly generated device id and returns the new deviceid.
//   - dir is the folder we'll write our deviceid file to.
//   - scrubValue is a free form description of what this directory represents, like the environment variable we used
//     to form the path. When errors occur that have the 'dir' in them it'll be replaced by this value.
func readWriteDeviceIDFile(dir string, scrubValue string) (string, error) {
	err := os.MkdirAll(dir, 0700)

	if err != nil {
		return "", scrubPathError(dir, scrubValue, err)
	}

	filePath := path.Join(dir, "deviceid")

	contents, err := os.ReadFile(filePath)

	if os.IsNotExist(err) {
		deviceID, err := generateDeviceID()

		if err != nil {
			return "", newError(err)
		}

		if err := os.WriteFile(filePath, []byte(deviceID), 0600); err != nil {
			// this error might have some user information in it (via the home folder)
			// so let's scrub that out and make it something non-identifying and generic.
			return "", scrubPathError(dir, scrubValue, err)
		}

		return deviceID, nil
	} else if err != nil {
		return "", scrubPathError(dir, scrubValue, err)
	}

	return string(contents), nil
}

func newError(err error) error {
	// we're purposefully removing any type/chain information here - there's no type
	// guarantees for errors returned from this library.
	return errors.New(err.Error())
}

func scrubPathError(folder string, replace string, err error) error {
	oldMessage := err.Error()
	newMessage := strings.Replace(oldMessage, folder, replace, 1)

	return newError(errors.New(newMessage))
}
