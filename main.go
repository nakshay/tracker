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

func main() {

	const (
		high = iota
		low
	)

	var timeBuffer time.Duration
	var shortBreakHour time.Duration
	var shortBreakAllowed time.Duration
	var longBreakHour time.Duration
	var longBreakAllowed time.Duration

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

			notify("Tracker", "Time to take a long break", "", low)

			// ask to resume work after long break expired
			time.AfterFunc(longBreakAllowed, resumeWork)

			if confirm() {
				fmt.Println("\nYour total time Overdue ", totalOverdue)
				os.Exit(1)
			} else {
				// re-initialize timer
				long = time.After(longBreakHour)
			}
		case <-short:

			// re-initialize timer
			short = time.After(shortBreakHour)

			notify("Tracker", "Time to take a short break", "", low)

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

func countShortBreak() {

}

func countLongBreak() {

}

func confirm() bool {

	result, err := dialog.Question("Confirm", "Do you want to exit", true)
	if err != nil {
		log.Fatal("error occuered", err)
	}
	return result
}

func resumeWork() {
	start := time.Now()
	_, err := dialog.Warning("Break over", "Do you want to resume work")
	if err != nil {
		log.Fatal("error occuered", err)
	}

	Overdue := time.Since(start)
	if Overdue > timeBuffer {
		totalOverdue += Overdue - timeBuffer
	} else {
		totalOverdue += 0 * time.Second
	}

	fmt.Println("\nYour time Overdue till now ", totalOverdue)
	fmt.Println("Time tracking Resumed")

}
