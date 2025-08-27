//go:build unix

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   workspace.json
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Unix-specific settings.
 *-----------------------------------------------------------------*/
package carousel

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
var DefaultUserOptions = Options{Notify: true, AssumeSession: FLAVOR_GNOME}

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
