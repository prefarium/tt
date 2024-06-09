package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func trackTime() error {
	trackFrom := time.Now()

	today, err := todayWindows(trackFrom)
	if err != nil {
		return err
	}

	lastWindow := len(today) - 1

	ticker := time.NewTicker(time.Second)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			today[lastWindow][1] = time.Now()
			fmt.Printf("\r%s", today)
		case <-exit:
			ticker.Stop()
			today[lastWindow][1] = time.Now()
			if err := write(today[lastWindow]); err != nil {
				return fmt.Errorf("failed to finish tracking last %s: %w", today[lastWindow], err)
			}
			return nil
		}
	}
}

func todayWindows(t time.Time) (day, error) {
	from, to, err := todayRange(t)
	if err != nil {
		return day{}, err
	}

	windows, err := read(from, to)
	if err != nil {
		return day{}, err
	}

	return append(windows, window{t, t}), nil
}

func todayRange(t time.Time) (time.Time, time.Time, error) {
	offsetDuration, err := time.ParseDuration(offset)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parsing offset duration: %w", err)
	}

	withOffset := t.Add(-offsetDuration)
	from := beginningOfDay(withOffset).Add(offsetDuration)
	to := nextDay(withOffset).Add(offsetDuration)

	return from, to, nil
}
