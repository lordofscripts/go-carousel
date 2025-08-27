//go:build unix

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Unix-specific Wallpaper Manager logic.
 *-----------------------------------------------------------------*/
package carousel

import (
	"fmt"
	"log"
	"os"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Determine the session manager in use in order to know which
 * handler to use.
 */
func (w *WallpaperManager) Init() error {
	if name, err := w.getSessionManager(); err == nil {
		switch name {
		case FLAVOR_GNOME:
			w.sessionHandler = newGnomeHandler()

		case FLAVOR_CINNAMON:
			w.sessionHandler = newGnomeHandler()
			w.sessionHandler.WithFlavor(FLAVOR_CINNAMON)

		case FLAVOR_LXDE:
			w.sessionHandler = newLxdeHandler()

		case FLAVOR_XFCE4:
			w.sessionHandler = newXfceHandler()

		default:
			return NewAppErrorMsg(ErrUnknownSessionManager, name)
		}
	} else {
		return err
	}

	//if cronned, err := IsCronJob(); err == nil && cronned {
	w.exportSessionBusAddress()
	//}

	return nil // go-static-check issue
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Find out which Desktop Session Manager is being used (Gnome, XFCE4, Cinnamon, LXDE).
 * This digs out the information for both interactive and CRON execution.
 */
func (w *WallpaperManager) getSessionManager() (string, error) {
	// (a) goCarousel interactive (from tty) we have an environment variable
	value, isSet := os.LookupEnv(ENV_SESSION)
	if isSet {
		return value, nil
	}

	// (b) goCarousel from CRON, no environment. We must query processes and
	//	   look for signs of a session.
	guess, err := w.cronDetermineSession()
	if err == nil {
		return guess, nil
	}

	// (c) goCarousel from CRON. Process query failed, use assumption
	//		set in the configuration file.
	if IsSupportedSession(w.settings.UserOptions.AssumeSession) {
		return w.settings.UserOptions.AssumeSession, nil
	}

	return value, NewAppErrorMsg(ErrUnknownSessionManager, "couldn't determine Session Manager")
}

func (w *WallpaperManager) exportSessionBusAddress() error {
	userId := os.Getuid()
	dBusSessionAddress := fmt.Sprintf("unix:path=/run/user/%d/bus", userId)
	err := os.Setenv("DBUS_SESSION_BUS_ADDRESS", dBusSessionAddress)
	if err != nil {
		log.Print("Unable to set DBUS_SESSION_BUS_ADDRESS")
	}
	xdgRuntimeDir := fmt.Sprintf("/run/user/%d", userId)
	err = os.Setenv("XDG_RUNTIME_DIR", xdgRuntimeDir)
	if err != nil {
		log.Print("Unable to set XDG_RUNTIME_DIR")
	}
	err = os.Setenv("DISPLAY", ":0")
	if err != nil {
		log.Print("Unable to set DISPLAY")
	}
	return err
}

func (w *WallpaperManager) cronDetermineSession() (string, error) {
	// @todo In CRON we should query and guess which one it is
	// ps -U userIDnum | grep SESSION
	// gnome-session, startkde/plasma/kwin, xfce4-session, lxsession, cinnamon-session

	return "", NewAppErrorMsg(ErrUnknownSessionManager, "process query failed")
}
