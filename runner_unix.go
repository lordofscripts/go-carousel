//go:build unix

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Unix-specific code for Runner.
 *-----------------------------------------------------------------*/
package carousel

import (
	"fmt"
	"log"
	"os"
	"strings"

	"lordofscripts/carousel/app"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Unix specific
	EXT_LSUSB    = "/usr/bin/lsusb"    // @note from JSON config
	EXT_LSBLK    = "/usr/bin/lsblk"    // @note idem
	EXT_LOGINCTL = "/usr/bin/loginctl" // @note idem
)

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

func init() {
	app.AssertOrDie(!FileExists(EXT_LSUSB), "Missing lsusb", 12)
	app.AssertOrDie(!FileExists(EXT_LSBLK), "Missing lsblk", 13)
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetMountPoint(volumeLabel string) (string, error) {
	outStr, err := ExecuteProgram(EXT_LSBLK, "-P", "-o", "name,label,mountpoint") // or use -P instead of -l
	if err != nil {
		log.Println("error IsDeviceOnline", err)
		return "", err
	}

	// with -P: NAME="sda7" LABEL="" MOUNTPOINT="/home"
	// with -l: sda7      /home
	lines := strings.Split(outStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, volumeLabel) { // LABEL @todo improve
			mpoint := line[strings.Index(line, "MOUNTPOINT=")+12 : len(line)-1]
			return mpoint, nil
		}
	}

	return "", fmt.Errorf("mountpoint of '%s' not found", volumeLabel)
}

/**
 * Verify that the USB device by VendorID:ProductID is connected
 */
func IsDeviceOnline(vendorId, productId string) (bool, error) {
	devId := vendorId + ":" + productId
	outStr, err := ExecuteProgram(EXT_LSUSB, "-d", devId)
	if err != nil {
		log.Println("error IsDeviceOnline", err)
		return false, err
	}

	if outStr == "" {
		return false, nil
	}

	return true, nil
}

// @audit need a more reliable version but so far it works well.
func IsCronJob() (bool, error) {
	return os.Getenv("HOME") == "", nil
}
