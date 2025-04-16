package calendar

import (
	"time"
)

type Event struct {
	Start time.Time
	End   time.Time
}

func (e Event) IsValid() bool {
	return !e.Start.IsZero() && !e.End.IsZero() && e.End.After(e.Start)
}

type Source interface {
	Events() ([]Event, error)
}

type Calendar struct {
	events []Event
}

func (c *Calendar) Events() []Event {
	return c.events
}

func (c *Calendar) Sync(source Source) error {
	events, err := source.Events()
	if err != nil {
		return err
	}

	c.events = events
	return nil
}

func New() *Calendar {
	return &Calendar{}
}
