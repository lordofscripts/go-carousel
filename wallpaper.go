/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   workspace.json
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package carousel

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ENV_SESSION = "GDMSESSION" // to determine whether it is Gnome, XFCE, etc.

	DEFAULT_ICON_FILE = ".category_icon.png"
	DEFAULT_AUTH_FILE = "goCarousel.png"
	DEFAULT_ICON      = "/home/lordofscripts/Pictures/Wallpapers/.category_icon.png" // 100x100 @audit
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/**
 * Every session manager (Gnome, XFCE, Cinammon) has its own handler
 */
type ISessionManager interface {
	/**
	 * Sets the Window Manager Flavor when a single handler can
	 * handle several types of flavors. For example, our GnomeHandler
	 * deals with the standard Gnome as well as Cinnamon.
	 */
	WithFlavor(string)

	/**
	* Get the current Color Theme (Light, Dark) by querying the
	* current session manager
	*
	* @returns (string) scheme name, i.e. "prefer-dark"
	* @returns (error) error if unable to determine
	 */
	QueryColorScheme() (string, error)

	/**
	* After determining the preferred/current color scheme, attempt
	* to set the wallpaper.
	*
	* @param (string) full path to wallpaper file
	* @returns (error) error if unable to set wallpaper
	 */
	SetWallpaperAuto(string) error

	SetWallpaperDark(string) error
	SetWallpaperLight(string) error

	String() string
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type WallpaperManager struct {
	settings       *Settings
	sessionHandler ISessionManager
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) New wallpaper manager instance. After this the Init()
 * method must be called because we need to determine the current
 * session manager in order to know how to set wallpapers.
 */
func NewWallpaperMgr(settings *Settings) *WallpaperManager {
	return &WallpaperManager{settings: settings, sessionHandler: nil}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Determine the session manager in use in order to know which
 * handler to use.
 */
func (w *WallpaperManager) Init() error {
	if name, err := w.getSessionManager(); err == nil {
		switch name {
		case FLAVOR_GNOME:
			w.sessionHandler = newGnomeHandler()

		case FLAVOR_CINNAMON:
			w.sessionHandler = newGnomeHandler()
			w.sessionHandler.WithFlavor(FLAVOR_CINNAMON)

		case FLAVOR_LXDE:
			w.sessionHandler = newLxdeHandler()

		case "xfce":
			w.sessionHandler = newXfceHandler()

		default:
			return NewAppErrorMsg(ErrUnknownSessionManager, name)
		}
	} else {
		return err
	}

	return nil // go-static-check issue
}

/**
 * Set the wallpaper but auto-determine whether it is chosen is Light|Dark
 */
func (w *WallpaperManager) SetWallpaperAuto(filename string) error {
	return w.sessionHandler.SetWallpaperAuto(filename)
}

/**
 * Set a Dark-themed wallpaper
 */
func (w *WallpaperManager) SetWallpaperDark(filename string) error {
	return w.sessionHandler.SetWallpaperDark(filename)
}

/**
 * Set a Light-themed wallpaper
 */
func (w *WallpaperManager) SetWallpaperLight(filename string) error {
	return w.sessionHandler.SetWallpaperLight(filename)
}

/**
 * Set a random wallpaper from the default non-categorized wallpaper
 * directory.
 */
func (w *WallpaperManager) SetAnyWallpaper() error {
	wallpaper, err := w.pickRandomFileIn(w.settings.DefaultDir)
	if err == nil {
		err = w.SetWallpaperAuto(wallpaper)
	}
	// Debian 13: /usr/share/images/desktop-base/default
	return err
}

/**
 * Set a random wallpaper from the chosen category.
 */
func (w *WallpaperManager) SetWallpaperFromCategory(chosenCategory string) error {
	if category, exists := w.settings.Categories[chosenCategory]; exists {
		preAuthorized := true
		if category.Protected {
			preAuthorized = w.authorize(category.KeyName)
			if !preAuthorized {
				log.Printf("authorization denied on %s", category.KeyName)
				if w.settings.UserOptions.Notify {
					NotifySound()
					NotifyAlert("Authorization denied", DEFAULT_ICON)
				}

				return NewWarningMsg(WarnAuthorizationDenied, "Authorization Denied")
			}
		}

		// Pick a random wallpaper from the chosen category
		if randomWallpaper, err := w.pickRandomFileIn(category.Directory); err != nil {
			return err
		} else {
			err := w.SetWallpaperAuto(randomWallpaper)
			if err == nil {
				iconFile := w.getIcon(category.Directory)
				message := fmt.Sprintf("Background from %s", chosenCategory)
				if w.settings.UserOptions.Notify {
					NotifyDesktop(message, iconFile)
				} else {
					log.Print(message)
				}
			}
		}
	}

	return NewAppErrorf(ErrUnknownCategory, "category named '%s' does not exist", chosenCategory)
}

/**
 * If the named carousel exists in the configuration, retrieve the categories
 * it is allowed to serve. Pick a random category from that list and then
 * delegate the Category work to @see SetWallpaperFromCategory()
 */
func (w *WallpaperManager) SetWallpaperFromCarousel(chosenCarousel string) error {
	if categories, exists := w.settings.Carousels[chosenCarousel]; exists {
		maxItems := len(categories)
		chosenIndex := w.getRandom(maxItems)
		categoryName := categories[chosenIndex]

		return w.SetWallpaperFromCategory(categoryName)
	}

	return NewAppErrorf(ErrUnknownCarousel, "carousel named '%s' does not exist", chosenCarousel).At("carousel")
}

func (w *WallpaperManager) Identify() string {
	return w.sessionHandler.String()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (w *WallpaperManager) authorize(deviceName string) bool {
	if deviceName == "" {
		return true
	}

	// (a) there is a matching device entry in the configuration file
	if deviceComposite, exists := w.settings.KeyDevices[deviceName]; exists {
		const KEY_PART_COUNT = 3
		const KEY_PART_IDENT = 0
		const KEY_PART_LABEL = 1
		const KEY_PART_HASH = 2
		// vendorId:productId volumeLabel MD5
		deviceParts := strings.Split(deviceComposite, " ")
		if len(deviceParts) != KEY_PART_COUNT {
			log.Printf("bad device-spec '%s' must be 'vendorId:productId volumeLabel'", deviceComposite)
			return false
		}

		// (b) the referenced device is plugged in
		// @todo what to do with vendorId:productId ?
		mountPoint, err := GetMountPoint(deviceParts[KEY_PART_LABEL])
		if err == nil {
			log.Printf("device %s mounted on %s", deviceParts[KEY_PART_LABEL], mountPoint)

			// (c) Auth File must be in place
			authFilename := path.Join(mountPoint, DEFAULT_AUTH_FILE)
			if FileExists(authFilename) {
				if md5, err := CalculateMD5(authFilename); err == nil {
					// (d) Its MD5 must match the one in the key
					if strings.EqualFold(md5, deviceParts[KEY_PART_HASH]) {
						return true
					}
				}
			} else {
				log.Printf("key object for %s not found", deviceParts[KEY_PART_LABEL])
			}
		} else {
			log.Printf("key carrier '%s' not plugged in %s", deviceParts[KEY_PART_IDENT], err)
		}
	} else {
		log.Printf("missing key %s", deviceName)
	}

	return false
}

/**
 * generate a true random integer between 0..N-1
 */
func (w *WallpaperManager) getRandom(upperLimit int) int64 {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(upperLimit)))
	if err != nil {
		log.Println("Error:", err)
		return -1
	}

	return int64(randomInt.Int64())
}

/**
 * Select a random image file from the selected directory
 */
func (w *WallpaperManager) pickRandomFileIn(dir string) (string, error) {
	// Read the directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return "", err
	}

	// Filter out non-regular files (like directories)
	var regularFiles []string
	for _, file := range files {
		if !file.IsDir() {
			if file.Name() == DEFAULT_ICON_FILE {
				continue
			}

			ext := strings.ToLower(path.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".svg" { // @todo globalize
				regularFiles = append(regularFiles, file.Name())
			}
		}
	}

	// Check if there are any files to choose from
	if len(regularFiles) == 0 {
		log.Println("No files found in the directory.")
		return "", NewAppErrorf(ErrNoQualifyingWallpaper, "no qualifying wallpaper files").At("carousel")
	}

	// Pick a random file therein
	randomIndex := w.getRandom(len(regularFiles))
	randomFile := regularFiles[randomIndex]

	return path.Join(dir, randomFile), nil
}

func (w *WallpaperManager) getIcon(dir string) string {
	filename := path.Join(dir, DEFAULT_ICON_FILE)
	_, err := os.Stat(filename)
	//return !errors.Is(err, os.ErrNotExist)
	if err != nil {
		filename = ""
	}

	return filename
}

func (w *WallpaperManager) getSessionManager() (string, error) {
	value, isSet := os.LookupEnv(ENV_SESSION)
	if isSet {
		return value, nil
	}

	return value, NewAppErrorMsg(ErrUnknownSessionManager, "couldn't determine Session Manager")
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
