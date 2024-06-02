package entity

import "time"

type Window struct {
	StartsAt time.Time
	EndsAt   time.Time
}

func (w Window) IsClosed() bool {
	return !w.EndsAt.IsZero()
}
