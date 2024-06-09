package main

import (
	"fmt"
	"time"
)

func calculateWeek() error {
	monday, err := thisMonday()
	if err != nil {
		return err
	}

	var wk week
	for i := range wk {
		from := monday.AddDate(0, 0, i)
		to := monday.AddDate(0, 0, i+1)

		if windows, err := read(from, to); err != nil {
			return err
		} else {
			wk[i] = windows
		}
	}

	fmt.Println(wk)
	return nil
}

func thisMonday() (time.Time, error) {
	parsedLocation, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, fmt.Errorf("loading location: %w", err)
	}

	today := beginningOfDay(time.Now().In(parsedLocation))
	return today.AddDate(0, 0, -dayNumber(today)), nil
}

func dayNumber(t time.Time) int {
	wd := t.Weekday()
	switch wd {
	case time.Sunday:
		return 6
	default:
		return int(wd) - 1
	}
}
