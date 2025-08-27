/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   workspace.json
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"lordofscripts/carousel"
	"lordofscripts/carousel/app"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/adhocore/gronx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
const (
	CONFIG_GROUP  string = "coralys"
	SETTINGS_FILE string = "goCarousel.json"
)

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

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

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

func getConfigPath(checkExists bool) string {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		app.DieWithError(err, 20)
	}

	cfgDir = path.Join(cfgDir, CONFIG_GROUP)
	if checkExists && !app.DirExists(cfgDir) {
		app.DieWithError(carousel.NewAppErrorMsg(carousel.ErrNoConfigurationDir, "config dir").At("main"), 21)
	}

	return cfgDir
}

func getConfigFilename() string {
	filename := path.Join(getConfigPath(true), SETTINGS_FILE)
	return filename
}

func getSettings(filename string) (*carousel.Settings, error) {
	var settings carousel.Settings
	if data, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	} else {
		return &settings, nil
	}
}

/**
 * Create a default Wallpaper Manager JSON configuration file that could
 * be custom-tailored by the end-user.
 */
func CreateDefault(location string) error {
	const (
		KEY_DEVICE_1 = "MMC_Card"
		KEY_DEVICE_2 = "USB-Maxell"

		PRIVATE_CAROUSEL = "Private"
		PUBLIC_CAROUSEL  = "Public"

		PRIVATE_CATEGORY_A = "Anime"
		PRIVATE_CATEGORY_B = "Gothic"
	)
	home, _ := os.UserHomeDir()

	angelEnter := carousel.ScheduleAction{
		Command:  carousel.ActChosenCarousel,
		Argument: "Public",
	}
	angelLeave := carousel.ScheduleAction{
		Command:  carousel.ActDefaultWallpaper,
		Argument: "",
	}

	var def *carousel.Settings = &carousel.Settings{
		DefaultDir:       path.Join(home, "Pictures", "Wallpapers"),
		DefaultWallpaper: "/usr/share/desktop-base/emerald-theme/wallpaper/gnome-background.xml",
		UserOptions:      carousel.DefaultUserOptions,
		Categories: map[string]*carousel.Category{
			"Aviation":         carousel.NewCategory(home + "/Pictures/Wallpapers/Aviation"),
			PRIVATE_CATEGORY_A: carousel.NewCategoryWithProtection(home+"/Pictures/Wallpapers/Anime", KEY_DEVICE_1),
			"Misc":             carousel.NewCategory(home + "/Pictures/Wallpapers/Misc"),
			"Nature":           carousel.NewCategory(home + "/Pictures/Wallpapers/Nature"),
			PRIVATE_CATEGORY_B: carousel.NewCategoryWithProtection(home+"/Pictures/Goth", KEY_DEVICE_1),
		},
		Carousels: map[string]carousel.CategoryCollection{
			PUBLIC_CAROUSEL:  carousel.NewCategoryCollection("Aviation", "Misc", "Nature"),
			PRIVATE_CAROUSEL: carousel.NewCategoryCollection("Anime", "Gothic"),
		},
		KeyDevices: map[string]string{
			KEY_DEVICE_2: "058f:6387 MAXELL_BLUE 5844bef71c16299cc5d73334153544be",
			KEY_DEVICE_1: "058f:6335 E0FD-1813 5844bef71c16299cc5d73334153544be",
		},
		AngelOptions: carousel.AngelOpts{
			FirstAction: angelEnter,
			LastAction:  angelLeave,
		},
		Schedules: []carousel.Schedule{
			*carousel.NewSchedule("Public Carousel workhour", "*/10 * * * 0,6", carousel.ActChosenCarousel, PUBLIC_CAROUSEL),
			*carousel.NewSchedule("Public Carousel afterhour", "*/10 15-23 * * 1-5", carousel.ActChosenCarousel, PUBLIC_CAROUSEL),
			*carousel.NewSchedule("Revert to default", "56 11 * * 1-5", carousel.ActDefaultWallpaper, ""),
			*carousel.NewSchedule("Lock Prior", "59 11 * * 1-5", carousel.ActLockCarousel, ""),
			*carousel.NewSchedule("Unlock After", "00 13 * * 1-5", carousel.ActUnlockCarousel, ""),
			*carousel.NewSchedule("Anime Time", "*/10 8-12 * * 1-5", carousel.ActChosenCategory, PRIVATE_CATEGORY_A),
			*carousel.NewSchedule("Fun Time", "*/10 13-14 * * 1-5", carousel.ActChosenCarousel, PRIVATE_CAROUSEL),
			*carousel.NewSchedule("The Woods", "* 16 * * 1-5", carousel.ActChosenFile, home+"/Pictures/Wallpapers/Nature/enchantedwood-fhd.jpg"),
		},
	}

	var err error
	var data []byte

	if data, err = json.MarshalIndent(def, "", "  "); err == nil {
		var fdOut *os.File
		if fdOut, err = os.Create(location); err == nil {
			_, err = fdOut.Write(data)
		}
	}

	return err
}

func CronTask(settings *carousel.Settings, tellNext bool) error {
	if len(settings.Schedules) > 0 {
		gron := gronx.New()

		exprs := make([]string, len(settings.Schedules))
		for idx, sched := range settings.Schedules {
			exprs[idx] = sched.CronTab
		}

		const TIMESTAMP_LAYOUT = "2006-01-02 15:04:05 -0700 MST"
		var jobSlice carousel.JobInfoSlice = make(carousel.JobInfoSlice, 0)
		anyTaskDue := false
		for idx, job := range settings.Schedules {
			if gron.IsValid(job.CronTab) {
				due, err := gron.IsDue(job.CronTab, time.Now())
				if err != nil {
					log.Printf("job #%d '%s' due error: %s", idx+1, job.Title, err)
				} else if due {
					if err := carousel.Execute(job.Command, job.Argument, settings); err != nil {
						log.Printf("job #%d '%s' exec error: %s", idx+1, job.Title, err)
						return err
					} else {
						log.Printf("Success running %s", job.Title)
						anyTaskDue = true
					}
				}

				if !due && tellNext {
					allowCurrent := true // include current time
					nextTime, err := gronx.NextTick(job.CronTab, allowCurrent)
					if err == nil {
						jobSlice = append(jobSlice, carousel.JobInfo{
							Id:        uint(idx + 1),
							TimeStamp: nextTime.Format(TIMESTAMP_LAYOUT),
							Title:     job.Title})
						//fmt.Printf("\tjob #%02d next due on %s\n", idx+1, nextTime)
					}
				}
			}
		}

		if tellNext {
			sort.Sort(carousel.JobInfoSlice(jobSlice))
			fmt.Println("\tTasks Next Due...")
			fmt.Println("\t" + strings.Repeat("-", 39))
			fmt.Printf("\t%10s %s\n", "Now is:", time.Now().Format(TIMESTAMP_LAYOUT))
			for _, ji := range jobSlice {
				fmt.Printf("\tjob #%02d on %s %s\n", ji.Id, ji.TimeStamp, ji.Title)
			}
		}

		if !anyTaskDue {
			log.Print("No Carousel tasks due")
		}
	}

	return nil
}

func Version() {
	carousel.Copyright(carousel.CO1, true)
	carousel.BuyMeCoffee("lostinwriting")
}

func Help() {
	carousel.Copyright(carousel.CO1, true)

	const NAME = "\tgoCarousel "
	fmt.Println("Usage:\t\t\t(Environment)")
	fmt.Println(NAME, "-init")
	fmt.Println(NAME, "-verify")
	fmt.Println(NAME, "-ident")
	fmt.Println("\t\t\t(Control)")
	fmt.Println(NAME, "-lock|-unlock|-status")
	fmt.Println("\t\t\t(Panic Mode)")
	fmt.Println(NAME, "-default")
	fmt.Println(NAME, "-any")
	fmt.Println("\t\t\t(Wallpapers)")
	fmt.Println(NAME, "-C|-category CATEGORY")
	fmt.Println(NAME, "-G|-carousel NAME")
	fmt.Println(NAME, "-F /path/to/wallpaper.jpg")
	fmt.Println("\t\t\t(Scheduling)")
	fmt.Println(NAME, "-task [-next]")
	fmt.Println(NAME, "-daemon MINUTES")
	//flag.PrintDefaults()

	carousel.BuyMeCoffee("lostinwriting")
}

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

func main() {
	fmt.Println("GO-GnomeChangeBackground")

	// ============= CLI FLAGS ===============
	var actInit, actHelp, actVersion, actAnyGlobal, actLock, actUnlock, actStatus, actDefault, actVerify, actWhoAmI bool
	var actTask, optNextTime bool
	var actDaemon int
	var group, category, filename string

	flag.BoolVar(&actHelp, "help", false, "Cry for help!")
	flag.BoolVar(&actVersion, "version", false, "Show version")
	flag.BoolVar(&actInit, "init", false, "Create default configuration")
	flag.BoolVar(&actVerify, "verify", false, "Verify CRON")
	flag.BoolVar(&actLock, "lock", false, "Prevent change")
	flag.BoolVar(&actUnlock, "unlock", false, "Allow changes")
	flag.BoolVar(&actStatus, "status", false, "Check Lock/Unlock state")
	flag.BoolVar(&actDefault, "default", false, "Set default wallpaper")
	flag.BoolVar(&actAnyGlobal, "any", false, "Select from Default wallpapers")
	flag.BoolVar(&actTask, "task", false, "Run any Wallpaper task that is due")
	flag.BoolVar(&optNextTime, "next", false, "Show when task is next due (only with -task)")
	flag.IntVar(&actDaemon, "daemon", -1, "Run as a dumb daemon for N minutes")
	flag.BoolVar(&actWhoAmI, "ident", false, "Identify and exit")
	flag.StringVar(&category, "C", "", "Select from this category")
	flag.StringVar(&category, "category", "", "Select from this category")
	flag.StringVar(&filename, "F", "", "Select this wallpaper")
	flag.StringVar(&group, "G", "", "Select this caroussel group")
	flag.StringVar(&group, "carousel", "", "Select this caroussel group")
	flag.Parse()

	// ============= CLI PROCESS ===============
	if actVersion {
		Version()
		os.Exit(0)
	}

	if actHelp {
		Help()
		os.Exit(0)
	}

	var err error
	var isCron bool
	var logFile *os.File = nil

	isCron, err = carousel.IsCronJob() // @audit IsCronJob is broken
	if err == nil && isCron {
		//if true {
		LOG_FILENAME := path.Join(app.GetUserTempDir(), "goCarousel.log")
		logFile, err = os.OpenFile(LOG_FILENAME, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Print(err)
		}

		log.SetOutput(logFile)
		log.Printf("goCarousel Started from CRON: %t", isCron)
	}
	defer func() {
		if logFile != nil {
			log.Println("... Closing ...")
			logFile.Close()
		}
	}()

	if actInit {
		// Ensure our configuration directory exists
		configPath := getConfigPath(false)
		if !app.DirExists(configPath) {
			if err := os.MkdirAll(configPath, 0644); err != nil {
				app.Die("could not create configuration path", 1)
			}
			fmt.Printf("Configuration Path: %s\n", configPath)
		}

		if err = CreateDefault(getConfigFilename()); err != nil {
			app.DieWithError(err, 1)
		}
		os.Exit(0)
	}

	var settings *carousel.Settings
	if settings, err = getSettings(getConfigFilename()); err != nil {
		app.DieWithError(err, 2)
	}

	if actLock {
		if err = carousel.LockCarousel(settings); err != nil {
			app.DieWithError(err, 3)
		}
		os.Exit(0)
	}

	if actUnlock {
		if err = carousel.UnlockCarousel(settings); err != nil {
			app.DieWithError(err, 4)
		}
		os.Exit(0)
	}

	if actStatus {
		if carousel.IsLocked(settings) {
			fmt.Println("Carousel is Locked")
			os.Exit(125)
		} else {
			fmt.Println("Carousel is NOT locked")
			os.Exit(0)
		}
	}

	if actVerify {
		wm := carousel.NewWallpaperMgr(settings)
		if err = wm.Init(); err == nil {
			fmt.Println("Window manager: ", wm.Identify())
		}
		fmt.Printf("Configuration: %s\n", getConfigFilename())
		fmt.Println("Verifying Cron Jobs...")
		var cumulative bool = true
		for idx, crontab := range settings.Schedules {
			ok := gronx.IsValid(crontab.CronTab)
			cumulative = cumulative && ok
			fmt.Printf("\t#%2d %t %s\n", idx+1, ok, crontab.Title)
		}
		if !cumulative {
			app.Die("Some Cron entries are invalid", 5)
		}
		os.Exit(0)
	}

	if actDaemon > -1 {
		CarouselTasker(settings, actDaemon)
		os.Exit(0)
	}

	if carousel.IsLocked(settings) {
		os.Exit(124)
	}

	if actTask {
		// (@) No command-line arguments? Execute Cron
		err = CronTask(settings, optNextTime)
		if err != nil {
			app.DieWithError(err, 5)
		}
	}

	// (@) Execute Wallpaper actions
	var action carousel.Action = carousel.ActNone
	var argument string = ""

	if actAnyGlobal {
		action = carousel.ActAnyWallpaper
	}
	if group != "" {
		action = carousel.ActChosenCarousel
		argument = group
	}
	if category != "" {
		action = carousel.ActChosenCategory
		argument = category
	}
	if filename != "" {
		action = carousel.ActChosenFile
		argument = filename
	}
	if actDefault {
		action = carousel.ActDefaultWallpaper
		argument = settings.DefaultWallpaper
	}
	if actWhoAmI {
		action = carousel.ActIdentify
	}

	err = carousel.Execute(action, argument, settings)
	log.Printf("exec %s %s returns %v", action, argument, err)
	if err != nil {
		app.DieWithError(err, 6)
	}
	log.Print("...done")
}
