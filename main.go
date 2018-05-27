package main

import (
	"fmt"
	"log"
	"os"
	"time"

	dialog "github.com/gen2brain/dlgs"
)

var totalOverdue time.Duration

func main() {

	const (
		high = iota
		low
	)

	notify("Tracker", "Work tracking started", "", low)

	shortBreak := time.Second * 5
	longBreak := time.Second * 10

	// Create channel for timers
	short := time.After(shortBreak)
	long := time.After(longBreak)

	for {
		select {
		case <-short:

			// re-initialize timer
			short = time.After(shortBreak)

			notify("Tracker", "Time to take a short break", "", low)

			// ask to resume work after short break expired
			time.AfterFunc(shortBreak, resumeWork)

		case <-long:
			notify("Tracker", "Time to take a long break", "", low)

			// ask to resume work after long break expired
			time.AfterFunc(longBreak, resumeWork)

			if confirm() {
				os.Exit(1)
			} else {
				// re-initialize timer
				long = time.After(longBreak)
			}
		default:
			fmt.Println("default case")
		}
	}

}

func confirm() bool {
	fmt.Println("function called")
	//result := dialog.Message("%s", "Do you want to exit?").Title("Are you sure?").YesNo()
	result, err := dialog.Question("Confirm", "Do you want to exit", true)
	if err != nil {
		log.Fatal("error occuered", err)
	}
	return result
}

func resumeWork() {
	current := time.Now()
	_, err := dialog.Warning("Break over", "Do you want to resume work")
	if err != nil {
		log.Fatal("error occuered", err)
	}

	totalOverdue = totalOverdue + time.Since(current)

	fmt.Println(totalOverdue)

}
