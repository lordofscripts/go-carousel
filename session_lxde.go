//go:build unix

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * LXDE Session Manager Built-in. This is known to work with
 * pcmanfm v1.4.
 * Status: Works
 *-----------------------------------------------------------------*/
package carousel

import (
	"os"
	"path"
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FLAVOR_LXDE = "LXDE"
	EXT_PCMANFM = "/usr/bin/pcmanfm" // @todo get from JSON config
	EXT_GREP    = "/usr/bin/grep"
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/
var _ ISessionManager = (*LxdeSession)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type LxdeSession struct{}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newLxdeHandler() *LxdeSession {
	return &LxdeSession{}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
* Sets the Window Manager Flavor when a single handler can
* handle several types of flavors. For example, our GnomeHandler
* deals with the standard Gnome as well as Cinnamon.
 */
func (s *LxdeSession) WithFlavor(schema string) {} // @note NoOp

/**
 * Get the current Color Theme (Light, Dark) by querying the
 * current session manager
 *
 * @returns (string) scheme name, i.e. "prefer-dark"
 * @returns (error) error if unable to determine
 */
func (s *LxdeSession) QueryColorScheme() (string, error) {
	var outStr string
	var err error

	home, _ := os.UserHomeDir()
	configFile := path.Join(home, ".config/gtk-3.0/settings.ini")

	// Grep: 0=found 1=Not found 2=Error
	outStr, err = ExecuteProgram(EXT_GREP,
		"-F",
		"gtk-theme-name",
		configFile,
	)

	if err != nil {
		return "", err
	}

	return outStr, nil
}

/**
 * After determining the preferred/current color scheme, attempt
 * to set the wallpaper.
 *
 * @param (string) full path to wallpaper file
 * @returns (error) error if unable to set wallpaper
 */
func (s *LxdeSession) SetWallpaperAuto(filename string) error {
	// gsettings get org.gnome.desktop.interface color-scheme
	var colorScheme string
	var err error
	if colorScheme, err = s.QueryColorScheme(); err == nil {
		if strings.Contains(strings.ToLower(colorScheme), "dark") { // "Adwaita-dark"
			err = s.SetWallpaperDark(filename)
		} else { // "Adwaita", "Xfce"
			err = s.SetWallpaperLight(filename)
		}
	}

	return err
}

func (s *LxdeSession) SetWallpaperDark(filename string) error {
	// pcmanfm --set-wallpaper=FILE
	// pcmanfm -w FILE
	_, err := ExecuteProgram(EXT_PCMANFM,
		"-wallpaper-mode=crop",
		"-w",
		filename,
	)

	return err
}

func (s *LxdeSession) SetWallpaperLight(filename string) error {
	_, err := ExecuteProgram(EXT_PCMANFM,
		"--wallpaper-mode=crop",
		"-w",
		filename,
	)

	return err
}

func (s *LxdeSession) String() string {
	return FLAVOR_LXDE
}
