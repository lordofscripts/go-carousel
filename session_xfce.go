/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * XFCE Session Manager Built-in. This is known to work with
 * XFCE v4.
 *-----------------------------------------------------------------*/
package carousel

import (
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	EXT_XFQUERY = "/usr/bin/xfconf-query" // @todo get from JSON config
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/
var _ ISessionManager = (*XfceSession)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type XfceSession struct{}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newXfceHandler() *XfceSession {
	return &XfceSession{}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
* Sets the Window Manager Flavor when a single handler can
* handle several types of flavors. For example, our GnomeHandler
* deals with the standard Gnome as well as Cinnamon.
 */
func (s *XfceSession) WithFlavor(schema string) {} // @note NoOp

/**
 * Get the current Color Theme (Light, Dark) by querying the
 * current session manager
 *
 * @returns (string) scheme name, i.e. "prefer-dark"
 * @returns (error) error if unable to determine
 */
func (s *XfceSession) QueryColorScheme() (string, error) {
	var outStr string
	var err error
	outStr, err = ExecuteProgram(EXT_XFQUERY,
		"--channel", // -c
		"xsettings",
		"--property", // -p
		"/Net/ThemeName",
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
func (s *XfceSession) SetWallpaperAuto(filename string) error {
	// gsettings get org.gnome.desktop.interface color-scheme
	var colorScheme string
	var err error
	if colorScheme, err = s.QueryColorScheme(); err == nil {
		if strings.Contains(colorScheme, "dark") { // "Adwaita-dark"
			err = s.SetWallpaperDark(filename)
		} else { // "Adwaita", "Xfce"
			err = s.SetWallpaperLight(filename)
		}
	}

	return err
}

func (s *XfceSession) SetWallpaperDark(filename string) error {
	// xfconf-query -c xfce4-desktop -p /backdrop/screen0/monitor0/image-path --set VALUE
	_, err := ExecuteProgram(EXT_XFQUERY,
		"--channel", // -c
		"xfce4-desktop",
		"--property", // -p
		"/backdrop/screen0/monitor0/image-path",
		"--set",
		filename,
	)

	return err
}

func (s *XfceSession) SetWallpaperLight(filename string) error {
	// xfconf-query -c xfce4-desktop -p /backdrop/screen0/monitor0/image-path --set VALUE
	_, err := ExecuteProgram(EXT_XFQUERY,
		"--channel", // -c
		"xfce4-desktop",
		"--property", // -p
		"/backdrop/screen0/monitor0/image-path",
		"--set",
		filename,
	)

	return err
}

func (s *XfceSession) String() string {
	return "Xfce4"
}
