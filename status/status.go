package status

import (
	"github.com/babenkoivan/busylight/calendar"
	"time"
)

const (
	Inactive = Status(iota)
	Active
	Focused
	Busy

	activeHourFrom   = 9
	inactiveHourFrom = 18

	focusTime = 10 * time.Minute
)

type Status int

func (s *Status) Sync(cal *calendar.Calendar) {
	now := time.Now()

	if now.Hour() >= inactiveHourFrom || now.Hour() < activeHourFrom {
		*s = Inactive
		return
	}

	focusCutoff := time.Now().Add(focusTime)
	status := Active

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
	status := Inactive
	return &status
}
