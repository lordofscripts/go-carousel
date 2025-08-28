/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Desktop notifications.
 *-----------------------------------------------------------------*/
package carousel

import (
	_ "embed"
	"os"

	"github.com/gen2brain/beeep"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const ()

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

func NotifyDesktop(body, iconPath string) error {
	var icon any = nil
	body = "<i>" + body + "</i>"

	if iconPath != "" {
		icon = getIconFile(iconPath)
	}

	err := beeep.Notify(NAME, body, icon)
	return err
}

func NotifyAlert(body, iconPath string) error {
	iconData := getIconFile(iconPath)
	err := beeep.Alert(NAME, body, iconData)
	return err
}

func NotifySound() error {
	err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	return err
}

func getIconFile(filename string) []byte {
	var data []byte
	data, err := os.ReadFile(filename)
	if err != nil {
		data = nil
	}

	return data
}
