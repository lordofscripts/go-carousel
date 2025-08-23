/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   workspace.json
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package carousel

import (
	"log"

	"github.com/adhocore/gronx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
const (
	CATEGORY_ICON_FILE = ".category_icon.png"
	NOTIFIER           = "/usr/bin/notify-send"
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type Settings struct {
	DefaultDir       string                        `json:"default_dir"`
	DefaultWallpaper string                        `json:"default_wallpaper"`
	UserOptions      Options                       `json:"options"`
	Categories       map[string]*Category          `json:"categories"`
	Carousels        map[string]CategoryCollection `json:"carousels"`
	KeyDevices       map[string]string             `json:"key_devices"`
	AngelOptions     AngelOpts                     `json:"angel"`
	Schedules        []Schedule                    `json:"schedules"`
}

type Options struct {
	Notify        bool   `json:"notify"`
	AssumeSession string `json:"assume_session"`
}

type Category struct {
	Protected bool   `json:"protected"`
	KeyName   string `json:"key_name,omitempty"`
	Directory string `json:"directory"`
}

type Schedule struct {
	Title    string `json:"title"`
	Command  Action `json:"action"` // random-in-cat, specific-file,
	Argument string `json:"argument"`
	CronTab  string `json:"cron_tab"`
}

type AngelOpts struct {
	FirstAction ScheduleAction `json:"first_action"`
	LastAction  ScheduleAction `json:"last_action"`
}

type ScheduleAction struct {
	Command  Action `json:"action"` // random-in-cat, specific-file,
	Argument string `json:"argument"`
}

type CategoryCollection = []string

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/
func NewCategory(dir string) *Category {
	return &Category{false, "", dir}
}

func NewCategoryWithProtection(dir string, keyName string) *Category {
	return &Category{true, keyName, dir}
}

func NewCategoryCollection(categories ...string) CategoryCollection {
	all := make([]string, len(categories))
	copy(all, categories)
	return all
}

func NewSchedule(title, cron string, action Action, arg string) *Schedule {
	if !gronx.IsValid(cron) {
		log.Println("ERR-Crontab ", cron)
		return nil
	}

	return &Schedule{title, action, arg, cron}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * DESCR
 * @params a (type):
 * @returns
 */

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * Whether the named session manager is one we can handle
 * @audit this must be updated as new DMs are supported
 */
func IsSupportedSession(sessionManager string) bool {
	if sessionManager == FLAVOR_GNOME ||
		sessionManager == FLAVOR_CINNAMON ||
		sessionManager == FLAVOR_LXDE ||
		sessionManager == FLAVOR_XFCE4 {
		return true
	}

	return false
}
