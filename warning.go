/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							 GoÜBRU
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Go Bulk Rename & Rename Workflow Utility
 *-----------------------------------------------------------------*/
package carousel

import "fmt"

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

// These are numeric error codes
const (
	WarnEmpty WarningCode = iota
	WarnAuthorizationDenied
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type WarningCode int

/**
 * Although the Warning object fulfills the Errors interface, it is
 * NOT an error, it is simply a warning that could be shown to the
 * user but does NOT warrant termination of the application NOR
 * termination of the current task.
 */
type Warning struct {
	warnum   WarningCode
	message  string
	location string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewWarning(num WarningCode, err error) *Warning {
	return &Warning{num, err.Error(), ""}
}

func NewWarningMsg(num WarningCode, message string) *Warning {
	return &Warning{num, message, ""}
}

func NewWarningWith(num WarningCode, message string, err error) *Warning {
	return &Warning{num, fmt.Sprintf("%s: %s", message, err.Error()), ""}
}

func NewWarningf(num WarningCode, format string, v ...any) *Warning {
	return &Warning{num, fmt.Sprintf(format, v...), ""}
}

func IsWarning(err error) bool {
	_, ok := err.(*Warning)
	return ok
}

/* ----------------------------------------------------------------
 *					P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Implement Error interface
 */
func (e *Warning) Error() string {
	var locStr string = ""
	if len(e.location) > 0 {
		locStr = fmt.Sprintf(" @{%s} ", e.location)
	}
	return fmt.Sprintf("(Warning) #W%03d%s%s", e.warnum, locStr, e.message)
}

/**
 * Set location of error.
 * err := NewAesGcmError(14, "Bad thing happened").At("module")
 */
func (e *Warning) At(location string) *Warning {
	e.location = location
	return e
}

func (e *Warning) Pretty() string {
	return fmt.Sprintf("\t(Warning)\n\tNumber: #W%03d\n\tMessage: %s\n\tLocation: %s\n", e.warnum, e.message, e.location)
}
