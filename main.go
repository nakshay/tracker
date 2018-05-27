package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	dialog "github.com/gen2brain/dlgs"
)

var totalOverdue time.Duration
var timeBuffer time.Duration
var shortBreakHour time.Duration
var shortBreakAllowed time.Duration
var longBreakHour time.Duration
var longBreakAllowed time.Duration

func main() {

	const (
		high = iota
		low
	)

	totalShortBreaks := 0
	totalLongBreaks := 0

	// buffer to resume work
	flag.DurationVar(&timeBuffer, "buffer", time.Second*10, "To provide additional buffer")

	// After shortBreakHour short break will start
	flag.DurationVar(&shortBreakHour, "shortBreakHour", time.Second*30, "To provide additional buffer")
	// Short break allowed for shortBreakAllowed duration
	flag.DurationVar(&shortBreakAllowed, "shortBreakAllowed", time.Second*10, "To provide additional buffer")

	// After shortBreakHour long break will start
	flag.DurationVar(&longBreakHour, "longBreakHour", time.Second*60, "To provide additional buffer")
	// Long break allowed for longBreakAllowed duration
	flag.DurationVar(&longBreakAllowed, "longBreakAllowed", time.Second*20, "To provide additional buffer")

	flag.Parse()

	notify("Tracker", "Work tracking started", "", low)

	fmt.Println("Time tracking started")

	// Create channel for timers
	short := time.After(shortBreakHour)
	long := time.After(longBreakHour)

	for {
		select {
		case <-long:

			// skip one short break when long break hits
			go skipShortBreak(short)

			notify("Tracker", "Time to take a long break", "", low)

			if confirm() {
				fmt.Println("\nYour total time Overdue ", totalOverdue)
				fmt.Printf("Total short breaks teaken %d\n", totalShortBreaks)
				fmt.Printf("Total short breaks teaken %d\n", totalLongBreaks)
				os.Exit(1)
			} else {
				// re-initialize timer
				totalLongBreaks++
				long = time.After(longBreakHour)
			}
		case <-short:

			notify("Tracker", "Time to take a short break", "", low)
			totalShortBreaks++
			// re-initialize timer
			short = time.After(shortBreakHour)

			// ask to resume work after short break expired
			time.AfterFunc(shortBreakAllowed, resumeWork)

		default:

			for _, char := range `-\|/` {
				fmt.Printf("\r%c ", char)
				time.Sleep(time.Millisecond * 100)
			}

		}
	}

}

func confirm() bool {

	result, err := dialog.Question("Press No for exit", "Do you want to resume work or exit?", false)
	if err != nil {
		log.Fatal("error occuered", err)
	}

	if !result {
		return true
	} else {
		// ask to resume work after long break expired
		time.AfterFunc(longBreakAllowed, resumeWork)
		return false
	}
}

func resumeWork() {

	start := time.Now()
	_, err := dialog.Warning("Break over", "Do you want to resume work")
	if err != nil {
		log.Fatal("error occuered", err)
	}

	overdue := time.Since(start)
	if overdue > timeBuffer {
		totalOverdue += overdue - timeBuffer
	} else {
		totalOverdue += 0 * time.Second
	}

	fmt.Println("\nYour time Overdue till now ", totalOverdue)
	fmt.Println("Time tracking Resumed")

}

func skipShortBreak(ch <-chan time.Time) {
	<-ch
}
