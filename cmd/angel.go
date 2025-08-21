/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCarousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"context"
	"fmt"
	"log"
	"lordofscripts/carousel"
	"time"

	"github.com/adhocore/gronx/pkg/tasker"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	DAEMON_VERBOSE          bool = false
	DAEMON_CONCURRENT_TASKS bool = false
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func CarouselTasker(settings *carousel.Settings, runUntilMinutes int) {
	log.Println("Angel battering wings in wallpaper heaven...")
	carousel.ExecuteCommand(settings.AngelOptions.FirstAction, settings)
	log.Print("Executed Angel.FirstAction")

	taskr := tasker.New(tasker.Option{
		Verbose: DAEMON_VERBOSE,
		// optional: defaults to local
		//Tz:      "Asia/Bangkok",
		// optional: defaults to stderr log stream
		//Out:     "/full/path/to/output-file",
	})

	// run task without overlap, set concurrent flag to false:
	concurrent := DAEMON_CONCURRENT_TASKS

	for jid, job := range settings.Schedules {
		taskr.Task(job.CronTab, func(ctx context.Context) (int, error) {
			taskr.Log.Printf("running Job #%d %s", jid+1, job.Title)

			err := carousel.Execute(job.Command, job.Argument, settings)
			return 0, err
		}, concurrent)
	}

	// optionally if you want tasker to stop after 2 hour, pass the duration with Until():
	//taskr.Until(2 * time.Hour)
	taskr.Until(time.Duration(runUntilMinutes) * time.Minute)

	// finally run the tasker, it ticks sharply on every minute and runs all the tasks due on that time!
	// it exits gracefully when ctrl+c is received making sure pending tasks are completed.
	taskr.Run()

	carousel.ExecuteCommand(settings.AngelOptions.LastAction, settings)
	log.Print("Executed Angel.LastAction")

	fmt.Println("Angels says goodbye...")
}

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/
