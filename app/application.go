/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Application-related functions.
 *-----------------------------------------------------------------*/
package app

import (
	"fmt"
	"log"
	"os"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	UC_RED_EXCLAMATION = rune(0x2757) // Dingbats
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * Death of an application by outputting a good-bye and setting
 * the OS exit code.
 */
func Die(message string, exitCode int) {
	log.Printf("die: (%d) %s", exitCode, message)
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
	os.Exit(exitCode)
}

func DieWithError(err error, exitCode int) {
	log.Printf("die: (%d) %s", exitCode, err)
	fmt.Printf("\n\t💀 x 💀 x 💀\n\t(%T)\n\t%s\n\tExit code: %d\n", err, err, exitCode)
	os.Exit(exitCode)
}

func Assert(condition bool, warnMessage string) {
	if condition {
		fmt.Printf("\n\t%c Assertion Failed:\n\t%s\n", UC_RED_EXCLAMATION, warnMessage)
	}
}

func AssertOrDie(condition bool, deathMessage string, exitCode int) {
	if condition {
		fmt.Printf("\n\t%c Assertion Failed:", UC_RED_EXCLAMATION)
		Die(deathMessage, exitCode)
	}
}
