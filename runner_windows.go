//go:build windows

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * External command runner.
 *-----------------------------------------------------------------*/
package carousel

import (
	"fmt"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

var (
	// setupapi.dll handle
	setupapi = syscall.NewLazyDLL("setupapi.dll")

	// procedures
	procSetupDiGetClassDevs              = setupapi.NewProc("SetupDiGetClassDevsW")
	procSetupDiEnumDeviceInfo            = setupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiGetDeviceRegistryProperty = setupapi.NewProc("SetupDiGetDeviceRegistryPropertyW")
	procSetupDiDestroyDeviceInfoList     = setupapi.NewProc("SetupDiDestroyDeviceInfoList")
)

// GUID for USB device class (from devguid.h)
var guidUsbClass = windows.GUID{
	Data1: 0x36fc9e60,
	Data2: 0xc465,
	Data3: 0x11cf,
	Data4: [8]byte{0x80, 0x56, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00},
}

const (
	DIGCF_PRESENT       = 0x00000002
	SPDRP_HARDWAREID    = 0x00000001
	ERROR_NO_MORE_ITEMS = 259
)

type SP_DEVINFO_DATA struct {
	CbSize    uint32
	ClassGuid windows.GUID
	DevInst   uint32
	Reserved  uintptr
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetMountPoint(volumeLabel string) (string, error) {
	mountPoint, err := GetRootPathByVolumeLabel(volumeLabel)
	return mountPoint, err
}

/**
 * Verify that the USB device by VendorID:ProductID is connected
 */
func IsDeviceOnline(vendorId, productId string) (bool, error) {
	isConnected, err := isDeviceConnected(vendorId, productId)
	if err != nil {
		return false, err
	}

	return isConnected, nil
}

func IsCronJob() (bool, error) {
	return false, nil
}

// GetRootPathByVolumeLabel finds the root path for a given volume label on Windows.
func GetRootPathByVolumeLabel(label string) (string, error) {
	// Step 1: Get a list of logical drive strings.
	buf := make([]uint16, 256)
	n, err := windows.GetLogicalDriveStrings(uint32(len(buf)), &buf[0])
	if err != nil {
		return "", fmt.Errorf("GetLogicalDriveStrings: %w", err)
	}

	// Step 2: Split the strings into individual drive letters.
	drives := strings.Split(string(utf16.Decode(buf[:n])), "\x00")

	// Step 3: Iterate through each drive letter.
	for _, drive := range drives {
		if len(drive) == 0 {
			continue
		}

		// Prepare buffers for volume information.
		volumeNameBuf := make([]uint16, 256)
		// No serial number or flags needed for this task.
		var serialNumber, maxComponentLen, fileSystemFlags uint32
		fileSystemNameBuf := make([]uint16, 256)

		drivePtr, err := syscall.UTF16PtrFromString(drive)
		if err != nil {
			return "", err
		}

		// Get volume information for the current drive.
		err = windows.GetVolumeInformation(
			drivePtr,
			&volumeNameBuf[0],
			uint32(len(volumeNameBuf)),
			&serialNumber,
			&maxComponentLen,
			&fileSystemFlags,
			&fileSystemNameBuf[0],
			uint32(len(fileSystemNameBuf)),
		)
		if err != nil {
			// Skip drives we can't access, like some network drives.
			continue
		}

		// Convert the volume name to a Go string and compare it.
		volumeName := windows.UTF16ToString(volumeNameBuf)
		if strings.EqualFold(volumeName, label) {
			return drive, nil
		}
	}

	return "", fmt.Errorf("volume with label '%s' not found", label)
}

// IsDeviceConnected checks for a connected USB device by its Vendor ID and Product ID.
func isDeviceConnected(vendorID, productID string) (bool, error) {
	// 1. Get a handle to all present devices of the specified class.
	r0, _, err := procSetupDiGetClassDevs.Call(
		uintptr(unsafe.Pointer(&guidUsbClass)),
		0,
		0,
		DIGCF_PRESENT,
	)
	devInfoSet := windows.Handle(r0)
	if devInfoSet == windows.InvalidHandle {
		return false, fmt.Errorf("SetupDiGetClassDevs: %w", err)
	}
	defer procSetupDiDestroyDeviceInfoList.Call(uintptr(devInfoSet))

	// 2. Iterate through each device in the set.
	var devInfo SP_DEVINFO_DATA
	devInfo.CbSize = uint32(unsafe.Sizeof(devInfo))
	for i := uint32(0); ; i++ {
		r1, _, err := procSetupDiEnumDeviceInfo.Call(
			uintptr(devInfoSet),
			uintptr(i),
			uintptr(unsafe.Pointer(&devInfo)),
		)
		if r1 == 0 {
			if err.(syscall.Errno) == ERROR_NO_MORE_ITEMS {
				break
			}
			return false, fmt.Errorf("SetupDiEnumDeviceInfo: %w", err)
		}

		// 3. Get the hardware ID.
		var regType uint32
		var reqSize uint32
		buf := make([]uint16, 2048)
		r2, _, err := procSetupDiGetDeviceRegistryProperty.Call(
			uintptr(devInfoSet),
			uintptr(unsafe.Pointer(&devInfo)),
			SPDRP_HARDWAREID,
			uintptr(unsafe.Pointer(&regType)),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(len(buf)*2),
			uintptr(unsafe.Pointer(&reqSize)),
		)
		if r2 == 0 {
			// Skip to next device if we can't get the hardware ID.
			continue
		}

		// 4. Parse the hardware ID string and check for VID/PID.
		hardwareID := syscall.UTF16ToString(buf)
		if strings.Contains(strings.ToUpper(hardwareID), "VID_"+strings.ToUpper(vendorID)) &&
			strings.Contains(strings.ToUpper(hardwareID), "PID_"+strings.ToUpper(productID)) {
			return true, nil
		}
	}

	return false, nil
}

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

/*
func main() {
	// Replace "MYUSB" with the volume label of your device.
	label := "MYUSB"
	rootPath, err := GetRootPathByVolumeLabel(label)
	if err != nil {
		fmt.Printf("Error finding device: %v\n", err)
		return
	}

	fmt.Printf("Found root path for device with label '%s': %s\n", label, rootPath)
}

func main() {
	// Replace with the VID and PID of your device.
	// Example: A common keyboard often uses VID 046D and PID C077.
	vendorID := "046D"
	productID := "C077"

	connected, err := IsDeviceConnected(vendorID, productID)
	if err != nil {
		fmt.Printf("Error checking device: %v\n", err)
		return
	}

	if connected {
		fmt.Printf("Device with VID %s and PID %s is connected.\n", vendorID, productID)
	} else {
		fmt.Printf("Device with VID %s and PID %s is not connected.\n", vendorID, productID)
	}
}
*/
