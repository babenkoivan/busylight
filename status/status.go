package status

import (
	"github.com/babenkoivan/busylight/calendar"
	"time"
)

const (
	_ = Status(iota)
	Idle
	Focused
	Busy

	focusTime = 10 * time.Minute
)

type Status int

func (s *Status) Sync(cal *calendar.Calendar) {
	now := time.Now()
	focusCutoff := time.Now().Add(focusTime)
	status := Idle

	for _, e := range cal.Events() {
		switch {
		case e.Start.Before(now) && e.End.After(now):
			status = Busy
			break
		case e.Start.After(now) && e.Start.Before(focusCutoff):
			status = Focused
		}
	}

	*s = status
}

func New() *Status {
	status := Idle
	return &status
}
