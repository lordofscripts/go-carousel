/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							 GoÜBRU
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Go Bulk Rename & Rename Workflow Utility
 *-----------------------------------------------------------------*/
package carousel

import (
	"fmt"
	"runtime"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ErrMissingTarget ApplicationErrorCode = iota // @note update String()
	ErrNoConfigurationDir
	ErrNoQualifyingWallpaper
	ErrUnknownCarousel
	ErrUnknownCategory
	ErrUnknownSessionManager
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type ApplicationErrorCode int

type AppErrorBase[T ~int] struct {
	errnum   T
	message  string
	location string
}

type BackgroundChangerError struct {
	AppErrorBase[ApplicationErrorCode]
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewAppError(num ApplicationErrorCode, err error) *BackgroundChangerError {
	return &BackgroundChangerError{AppErrorBase[ApplicationErrorCode]{num, err.Error(), getCallerFlex(2)}}
}

func NewAppErrorMsg(num ApplicationErrorCode, message string) *BackgroundChangerError {
	return &BackgroundChangerError{AppErrorBase[ApplicationErrorCode]{num, message, getCallerFlex(2)}}
}

func NewAppErrorWith(num ApplicationErrorCode, message string, err error) *BackgroundChangerError {
	return &BackgroundChangerError{AppErrorBase[ApplicationErrorCode]{num, fmt.Sprintf("%s: %s", message, err.Error()), getCallerFlex(2)}}
}

func NewAppErrorf(num ApplicationErrorCode, format string, v ...any) *BackgroundChangerError {
	return &BackgroundChangerError{AppErrorBase[ApplicationErrorCode]{num, fmt.Sprintf(format, v...), getCallerFlex(2)}}
}

/*
func newAppError[T int](num T, err error) *AppErrorBase[T] {
	return &AppErrorBase[T]{num, err.Error(), getCallerFlex(2)}
}

func newAppErrorMsg[T int](num T, message string) *AppErrorBase[T] {
	return &AppErrorBase[T]{num, message, getCallerFlex(2)}
}

func newAppErrorWith[T int](num T, message string, err error) *AppErrorBase[T] {
	return &AppErrorBase[T]{num, fmt.Sprintf("%s: %s", message, err.Error()), getCallerFlex(2)}
}

func newAppErrorf[T int](num T, format string, v ...any) *AppErrorBase[T] {
	return &AppErrorBase[T]{num, fmt.Sprintf(format, v...), getCallerFlex(2)}
}
*/
/* ----------------------------------------------------------------
 *					P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Implement Error interface
 */
func (e *AppErrorBase[T]) Error() string {
	t := fmt.Sprintf("%T", e.errnum)
	if p := strings.Index(t, "."); p != -1 {
		t = t[p:]
	}
	return fmt.Sprintf("(%s) #E%03d %s\n\t%s %s", t, e.errnum, e.errnum, e.location, e.message)
}

/**
 * Set location of error.
 * err := NewAesGcmError(14, "Bad thing happened").At("module")
 */
func (e *AppErrorBase[T]) At(location string) *AppErrorBase[T] {
	e.location = location
	return e
}

func (e *AppErrorBase[T]) Pretty() string {
	return fmt.Sprintf("\t(Error)\n\tNumber: #E%03d\n\tMessage: %s\n\tLocation: %s\n", e.errnum, e.message, e.location)
}

func (n ApplicationErrorCode) String() string {
	toString := map[ApplicationErrorCode]string{
		ErrMissingTarget:         "ErrMissingTarget",
		ErrNoConfigurationDir:    "ErrNoConfigurationDir",
		ErrNoQualifyingWallpaper: "ErrNoQualifyingWallpaper",
		ErrUnknownCarousel:       "ErrUnknownCarousel",
		ErrUnknownCategory:       "ErrUnknownCategory",
		ErrUnknownSessionManager: "ErrUnknownSessionManager",
	}
	return toString[n]
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// @todo fix this output format
// @lordofscripts/bulkrename/flow.(*Workflow).recurseDirectory()#139
func getCallerFlex(stackIdx int) string {
	// 3 works when test.Announce() is used
	// 2 works when ShowCaseOK/Failed called directly from Test*()
	pc, _, line, ok := runtime.Caller(stackIdx) // PC,file,line,ok
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		info := details.Name()
		lastSlash := strings.LastIndexByte(info, '/')
		if lastSlash < 0 {
			lastSlash = 0
		}
		lastDot := strings.LastIndexByte(info[lastSlash:], '.') + lastSlash

		//fmt.Printf("INFO %s\n", info)
		//fileS := filepath.Base(filename)
		packageS := info[:lastDot] // in tests it returns 'command-line-*'
		if strings.HasPrefix(packageS, "command-line") {
			packageS = "main"
		}
		funcS := info[lastDot+1:]
		return fmt.Sprintf("@%s.%s()#%d\n", packageS, funcS, line)
		//fmt.Printf("called from %s\n", details.Name())
		//return filename + "+" + string(line) + " " + packageS + ":" + funcS
	}
	return ""
}
