/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   workspace.json
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package carousel

import (
	_ "embed"

	"github.com/gen2brain/beeep"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	EXT_NOTIFY = "/usr/bin/notify-send"
)

//go:embed docs/assets/goCarousel.png
var defaultIconData []byte

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
func init() {
	beeep.AppName = NAME
}

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

func NotifyDesktopExt(body, iconPath string) error {
	_, err := ExecuteProgram(EXT_NOTIFY, // @note never returns...
		"-e",
		"--app-name='gnomeBackgroundChange",
		"-t", "1500",
		"-i", iconPath,
		"Gnome Background Change",
		"<i>"+body+"</i>")
	return err
}

func NotifyDesktop(body, iconPath string) error {
	body = "<i>" + body + "</i>"
	err := beeep.Notify(NAME, body, iconPath)
	return err
}

func NotifyAlert(body, iconPath string) error {
	err := beeep.Alert(NAME, body, iconPath)
	return err
}

func NotifySound() error {
	err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	return err
}
