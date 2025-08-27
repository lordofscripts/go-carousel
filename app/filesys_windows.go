/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Windoze-specific code for FileSys.
 *-----------------------------------------------------------------*/
package app

import (
	"os"
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetUserTempDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return `C:\`
	}
	return dir
}
