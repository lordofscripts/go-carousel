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
	"encoding/json"
	"fmt"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
const (
	ActNone Action = iota
	ActDefaultWallpaper
	ActAnyWallpaper
	ActLockCarousel
	ActUnlockCarousel
	ActChosenFile
	ActChosenCategory
	ActChosenCarousel
	ActStatus
	ActIdentify
)

/* ----------------------------------------------------------------
 *						L o c a l s
 *-----------------------------------------------------------------*/

var toString = map[Action]string{
	ActNone:             "ActNone",
	ActDefaultWallpaper: "ActDefaultWallpaper",
	ActAnyWallpaper:     "ActAnyWallpaper",
	ActLockCarousel:     "ActLockCarousel",
	ActUnlockCarousel:   "ActUnlockCarousel",
	ActChosenFile:       "ActChosenFile",
	ActChosenCategory:   "ActChosenCategory",
	ActChosenCarousel:   "ActChosenCarousel",
	ActStatus:           "ActStatus",
	ActIdentify:         "ActIdentify",
}

var toID = map[string]Action{
	"ActNone":             ActNone,
	"ActDefaultWallpaper": ActDefaultWallpaper,
	"ActAnyWallpaper":     ActAnyWallpaper,
	"ActLockCarousel":     ActLockCarousel,
	"ActUnlockCarousel":   ActUnlockCarousel,
	"ActChosenFile":       ActChosenFile,
	"ActChosenCategory":   ActChosenCategory,
	"ActChosenCarousel":   ActChosenCarousel,
	"ActStatus":           ActStatus,
	"ActIdentify":         ActIdentify,
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type Action uint8

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (s Action) String() string {
	return toString[s]
}

// MarshalJSON marshals the enum as a quoted json string
func (s Action) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Action) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toID[j]
	return nil
}

func (c Action) Parse(v string) (Action, error) {
	if v, ok := toID[v]; !ok {
		return ActDefaultWallpaper, fmt.Errorf("invalid enum '%s'", v)
	} else {
		return v, nil
	}
}
