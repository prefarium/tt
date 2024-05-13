package entity

import "time"

type Window struct {
	StartsAt time.Time
	EndsAt   time.Time
}

func (w *Window) IsCovered(from, to time.Time) bool {
	return (from.Before(w.StartsAt) || from.Equal(w.StartsAt)) && to.After(w.EndsAt)
}
