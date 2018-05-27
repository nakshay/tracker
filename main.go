package main

import (
	"fmt"

	"github.com/sqweek/dialog"
)

const (
	high = iota
	low
)

func main() {
	//notify("Break", "Short break", "", high)
	ok := dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
	fmt.Println(ok)
}
