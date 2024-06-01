package entity

import (
	"fmt"
	"time"
)

type Day time.Duration

func (d Day) String() string {
	dur := time.Duration(d)

	hours := dur / time.Hour
	dur -= hours * time.Hour

	minutes := dur / time.Minute
	dur -= minutes * time.Minute

	seconds := dur / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
