package ics

import (
	"github.com/apognu/gocal"
	"github.com/babenkoivan/busylight/calendar"
	"net/http"
	"time"
)

type ICS struct {
	url string
}

func (s ICS) Today() ([]calendar.Event, error) {
	resp, err := http.Get(s.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	cal := gocal.NewParser(resp.Body)
	cal.Start, cal.End = &start, &end
	err = cal.Parse()
	if err != nil {
		return nil, err
	}

	events := make([]calendar.Event, 0, len(cal.Events))
	for _, e := range cal.Events {
		events = append(events, calendar.Event{
			Start: *e.Start,
			End:   *e.End,
		})
	}

	return events, nil
}

func New(url string) ICS {
	return ICS{url: url}
}
