/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * External command runner.
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

	"github.com/gen2brain/beeep"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
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

func FileExists(filename string) bool { //@audit deprecate in favor of app.FileExists
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
		if err == nil {
			beeep.Notify("Beware", "Wallpapers locked", defaultIconData)
		}
	} else {
		beeep.Alert("Beware", "Could not lock carousel", defaultIconData)
		log.Print(err)
	}

	if fdOut != nil {
		fdOut.Close()
	}
	return err
}

/**
 * Remove the lock, thus allowing changing the wallpaper
 */
func UnlockCarousel(settings *Settings) error {
	var err error
	if err = os.Remove(getLockFile(settings)); err != nil {
		beeep.Alert("Beware", "Could not unlock carousel", defaultIconData)
		log.Print(err)
	} else {
		beeep.Notify("Just to let you know", "Wallpapers unlocked", defaultIconData)
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
