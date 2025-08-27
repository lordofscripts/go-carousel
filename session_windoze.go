//go:build windows

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MS-Windows Session Manager Built-in.
 * Status: Works
 *-----------------------------------------------------------------*/
package carousel

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FLAVOR_WINDOWS = "Windows"
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724947.aspx
const (
	spiGetDeskWallpaper = 0x0073
	spiSetDeskWallpaper = 0x0014

	uiParam = 0x0000

	spifUpdateINIFile = 0x01
	spifSendChange    = 0x02
)

const (
	Center Mode = iota
	Crop
	Fit
	Span
	Stretch
	Tile
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724947.aspx
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/
var _ ISessionManager = (*WindowsSession)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type Mode int

type WindowsSession struct{}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newWindozeHandler() *WindowsSession {
	return &WindowsSession{}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Sets the Window Manager Flavor when a single handler can
 * handle several types of flavors. For example, our WindozeHandler
 * could in principle deal with Windows7/10/11.
 * NOTE: No Windows flavors supported  yet
 */
func (s *WindowsSession) WithFlavor(schema string) {} // @note NoOp

/**
 * Get the current Color Theme (Light, Dark) by querying the
 * current session manager
 *
 * @returns (string) scheme name, i.e. "prefer-dark"
 * @returns (error) error if unable to determine
 */
func (s *WindowsSession) QueryColorScheme() (string, error) {
	return "", fmt.Errorf("not implemented exception")
}

/**
 * After determining the preferred/current color scheme, attempt
 * to set the wallpaper.
 *
 * @param (string) full path to wallpaper file
 * @returns (error) error if unable to set wallpaper
 */
func (s *WindowsSession) SetWallpaperAuto(filename string) error {
	return SetFromFile(filename)
}

func (s *WindowsSession) SetWallpaperDark(filename string) error {
	return SetFromFile(filename)
}

func (s *WindowsSession) SetWallpaperLight(filename string) error {
	return SetFromFile(filename)
}

func (s *WindowsSession) String() string {
	return FLAVOR_WINDOWS
}

// Get returns the current wallpaper.
func Get() (string, error) {
	// the maximum length of a windows path is 256 utf16 characters
	var filename [256]uint16
	systemParametersInfo.Call(
		uintptr(spiGetDeskWallpaper),
		uintptr(cap(filename)),
		// the memory address of the first byte of the array
		uintptr(unsafe.Pointer(&filename[0])),
		uintptr(0),
	)
	return strings.Trim(string(utf16.Decode(filename[:])), "\x00"), nil
}

// SetFromFile sets the wallpaper for the current user.
func SetFromFile(filename string) error {
	filenameUTF16, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}

	systemParametersInfo.Call(
		uintptr(spiSetDeskWallpaper),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(filenameUTF16)),
		uintptr(spifUpdateINIFile|spifSendChange),
	)
	return nil
}

// SetMode sets the wallpaper mode.
func SetMode(mode Mode) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, "Control Panel\\Desktop", registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	var tile string
	if mode == Tile {
		tile = "1"
	} else {
		tile = "0"
	}
	err = key.SetStringValue("TileWallpaper", tile)
	if err != nil {
		return err
	}

	var style string
	switch mode {
	case Center, Tile:
		style = "0"
	case Fit:
		style = "6"
	case Span:
		style = "22"
	case Stretch:
		style = "2"
	case Crop:
		style = "10"
	default:
		panic("invalid wallpaper mode")
	}
	err = key.SetStringValue("WallpaperStyle", style)
	if err != nil {
		return err
	}

	// updates wallpaper
	path, err := Get()
	if err != nil {
		return err
	}

	return SetFromFile(path)
}

func getCacheDir() (string, error) {
	return os.TempDir(), nil
}
