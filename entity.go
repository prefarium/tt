package main

import (
	"fmt"
	"strings"
	"time"
)

var weekDays = map[int]string{
	0: "Mon",
	1: "Tue",
	2: "Wed",
	3: "Thu",
	4: "Fri",
	5: "Sat",
	6: "San",
}

type window [2]time.Time
type day []window
type week [7]day

func (w window) String() string {
	return prettyDuration(w.duration())
}

func (w window) duration() time.Duration {
	return w[1].Sub(w[0])
}

func (d day) String() string {
	return prettyDuration(d.duration())
}

func (d day) duration() time.Duration {
	var total time.Duration
	for _, w := range d {
		total += w.duration()
	}
	return total
}

func (w week) String() string {
	b := strings.Builder{}

	for i, d := range w {
		prettyDay := fmt.Sprintf("%s %s\n", weekDays[i], d)
		b.WriteString(prettyDay)
	}

	daysSummary := fmt.Sprintf("---\nSum %s", w.duration())
	b.WriteString(daysSummary)
	return b.String()
}

func (w week) duration() time.Duration {
	var total time.Duration
	for _, d := range w {
		total += d.duration()
	}
	return total
}

func prettyDuration(d time.Duration) string {
	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute
	d -= minutes * time.Minute

	seconds := d / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
