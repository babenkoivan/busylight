package calendar

import (
	"time"
)

const maxEventDuration = 10 * time.Hour

type Event struct {
	Start time.Time
	End   time.Time
}

func (e Event) IsValid() bool {
	return !e.Start.IsZero() && !e.End.IsZero() && e.End.After(e.Start)
}

func (e Event) Duration() time.Duration {
	return e.End.Sub(e.Start)
}

type Source interface {
	Today() ([]Event, error)
}

type Calendar struct {
	events []Event
}

func (c *Calendar) Events() []Event {
	return c.events
}

func (c *Calendar) Sync(source Source) error {
	todayEvents, err := source.Today()
	if err != nil {
		return err
	}

	validEvents := make([]Event, 0, len(todayEvents))
	for _, e := range todayEvents {
		if e.IsValid() && e.Duration() <= maxEventDuration {
			validEvents = append(validEvents, e)
		}
	}
	c.events = validEvents
	return nil
}

func New() *Calendar {
	return &Calendar{}
}
