/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Gnome Session Manager Built-in. This is known to work with
 * Gnome 43..48 and Cinnamon.
 * Status: Works
 *-----------------------------------------------------------------*/
package carousel

import (
	"fmt"
	"log"
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FLAVOR_GNOME    = "gnome"
	FLAVOR_CINNAMON = "cinnamon"

	EXT_GSETTINGS = "/usr/bin/gsettings" // @todo get from JSON config

	orgGnomeBackground    = "org.gnome.desktop.background"
	orgCinnamonBackground = "org.cinnamon.desktop.background"
	orgGnomeScheme        = "org.gnome.desktop.interface"
	orgCinnamonScheme     = "org.cinnamon.desktop.interface"
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/
var _ ISessionManager = (*GnomeSession)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type GnomeSession struct {
	schemaBackground string
	schemaInterface  string
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newGnomeHandler() *GnomeSession {
	return &GnomeSession{
		schemaBackground: orgGnomeBackground,
		schemaInterface:  orgGnomeScheme,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
* Sets the Window Manager Flavor when a single handler can
* handle several types of flavors. For example, our GnomeHandler
* deals with the standard Gnome as well as Cinnamon.
 */
func (s *GnomeSession) WithFlavor(schema string) {
	switch schema {
	case FLAVOR_GNOME:
		s.schemaBackground = orgGnomeBackground
		s.schemaInterface = orgGnomeScheme

	case FLAVOR_CINNAMON:
		s.schemaBackground = orgCinnamonBackground
		s.schemaInterface = orgCinnamonScheme

	default:
		log.Printf("unknown GnomeSession flavor '%s'", schema)
	}
}

/**
 * Get the current Color Theme (Light, Dark) by querying the
 * current session manager
 *
 * @returns (string) scheme name, i.e. "prefer-dark"
 * @returns (error) error if unable to determine
 */
func (s *GnomeSession) QueryColorScheme() (string, error) {
	var outStr string
	var err error
	var key string = "color-scheme"

	if s.schemaInterface == orgCinnamonScheme {
		key = "gtk-theme"
	}

	outStr, err = ExecuteProgram(EXT_GSETTINGS,
		"get",
		s.schemaInterface,
		key,
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
func (s *GnomeSession) SetWallpaperAuto(filename string) error {
	// gsettings get org.gnome.desktop.interface color-scheme
	var colorScheme string
	var err error
	if colorScheme, err = s.QueryColorScheme(); err == nil {
		if strings.Contains(colorScheme, "prefer-dark") {
			err = s.SetWallpaperDark(filename)
		} else {
			err = s.SetWallpaperLight(filename)
		}
	}

	return err
}

func (s *GnomeSession) SetWallpaperDark(filename string) error {
	// gsettings set org.gnome.desktop.background picture-uri-dark file://$1
	_, err := ExecuteProgram(EXT_GSETTINGS,
		"set",
		s.schemaBackground,
		"picture-uri-dark",
		fmt.Sprintf("file://%s", filename),
	)

	return err
}

func (s *GnomeSession) SetWallpaperLight(filename string) error {
	// gsettings set org.gnome.desktop.background picture-uri file://$1
	_, err := ExecuteProgram(EXT_GSETTINGS,
		"set",
		s.schemaBackground,
		"picture-uri",
		fmt.Sprintf("file://%s", filename),
	)
	return err
}

func (s *GnomeSession) String() string {
	var identity string = FLAVOR_GNOME
	if s.schemaBackground == orgCinnamonBackground {
		identity = FLAVOR_CINNAMON
	}

	return identity
}
