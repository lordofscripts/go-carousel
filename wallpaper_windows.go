//go:build windows

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Windows-specific Wallpaper manager logic.
 *-----------------------------------------------------------------*/
package carousel

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
	w.sessionHandler = newWindozeHandler()

	return nil
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/
