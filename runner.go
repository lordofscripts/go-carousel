/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   workspace.json
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package carousel

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	EXT_LSUSB    = "/usr/bin/lsusb"    // @note from JSON config
	EXT_LSBLK    = "/usr/bin/lsblk"    // @note idem
	EXT_LOGINCTL = "/usr/bin/loginctl" // @note idem

	STOP_FILE string = "%W/.nochange"
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

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

func ExecuteProgram(programPath string, args ...string) (string, error) {
	cmd := exec.Command(programPath, args...)

	// Create a buffer to capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	err := cmd.Run()
	if err != nil && cmd.ProcessState.ExitCode() == 2 {
		log.Printf("Error: %d %s", cmd.ProcessState.ExitCode(), err)
		return "", err
	}

	// Print the output
	return out.String(), nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

func CalculateMD5(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err) //@audit custom error
	}
	defer file.Close()

	// Create a new MD5 hash
	hash := md5.New()

	// Copy the file contents to the hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error calculating MD5 hash: %w", err) //@audit custom error
	}

	// Calculate the checksum in hexadecimal format
	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

/**
 * Apply the lock to prevent changing the wallpaper via this application
 */
func LockCarousel(settings *Settings) error {
	var err error
	var fdOut *os.File
	if fdOut, err = os.Create(getLockFile(settings)); err == nil {
		_, err = fdOut.WriteString(time.Now().String())
	} else {
		log.Print(err)
	}

	return err
}

/**
 * Remove the lock, thus allowing changing the wallpaper
 */
func UnlockCarousel(settings *Settings) error {
	var err error
	if err = os.Remove(getLockFile(settings)); err != nil {
		log.Print(err)
	}

	return err
}

/**
 * If the lock file exists we are locked from changing the wallpaper
 * via this application.
 */
func IsLocked(settings *Settings) bool {
	_, err := os.Stat(getLockFile(settings))
	return !errors.Is(err, os.ErrNotExist)
}

func ExecuteCommand(cmd ScheduleAction, settings *Settings) error {
	return Execute(cmd.Command, cmd.Argument, settings)
}

/**
 * Execute the application sub-command
 */
func Execute(command Action, argument string, settings *Settings) error {
	var err error = nil
	wm := NewWallpaperMgr(settings)
	if err = wm.Init(); err != nil {
		return err
	}

	switch command {
	case ActIdentify:
		fmt.Printf("%s for %s\n", NAME, wm.Identify())

	case ActDefaultWallpaper:
		wm.SetWallpaperAuto(settings.DefaultWallpaper)

	case ActAnyWallpaper:
		err = wm.SetAnyWallpaper()

	case ActLockCarousel:
		err = LockCarousel(settings)

	case ActUnlockCarousel:
		err = UnlockCarousel(settings)

	case ActChosenFile:
		wm.SetWallpaperAuto(argument)

	case ActChosenCategory:
		wm.SetWallpaperFromCategory(argument)

	case ActChosenCarousel:
		wm.SetWallpaperFromCarousel(argument)

	case ActStatus:
		if IsLocked(settings) {
			fmt.Println("Carousel is Locked")
		} else {
			fmt.Println("Carousel is NOT locked")
		}

	case ActNone:

	default:
		fmt.Println("unknown command ", command)
	}

	return err
}

/**
 * Get the name of the lock file by untokenizing %W with Settings.DefaultDir
 */
func getLockFile(settings *Settings) string {
	stopFilename := STOP_FILE
	stopFilename = strings.Replace(stopFilename, "%W", settings.DefaultDir, 1)
	return stopFilename
}

// @audit this is not working....
func IsCronJob() (bool, error) {
	// For OS running systemd: loginctl show-session "$(</proc/self/sessionid)" | sed -n 's/^Service=//p'
	// that should return crond if running from CRON, else systemd-user

	/*
		exec.Command(EXT_LOGINCTL, "show-session", )
		parentProcess, err := exec.Command("ps", "-o", "comm=", fmt.Sprint(os.Getppid())).Output()
		if err != nil {
			log.Println("isCronJob Error:", err)
			return false, err
		}

		if strings.EqualFold(strings.TrimSpace(string(parentProcess)), "cron") {
			return true, nil
		} else {
			return false, nil
		}
	*/
	return os.Getenv("HOME") == "", nil
}
