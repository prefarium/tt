package entity

import "time"

type Week [7]Day

func (w Week) Total() time.Duration {
	var total time.Duration
	for _, d := range w {
		total += time.Duration(d)
	}
	return total
}
