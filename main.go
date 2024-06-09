package main

import (
	"fmt"
	"os"
)

var (
	csvPath  string
	offset   string
	location string
)

func main() {
	var err error

	if len(os.Args) == 1 {
		err = trackTime()
	} else if os.Args[1] == "week" {
		err = calculateWeek()
	} else {
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		fmt.Println(err)
	}
}
