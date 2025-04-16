package ics

import (
	"bufio"
	"fmt"
	"github.com/babenkoivan/busylight/calendar"
	"github.com/babenkoivan/busylight/timezone"
	"net/http"
	"strings"
	"time"
)

const (
	sepComponent = ":"
	sepToken     = ";"
	sepKeyValue  = "="

	entryBeginEvent = "BEGIN:VEVENT"
	entryEndEvent   = "END:VEVENT"

	tokenStartDate = "DTSTART"
	tokenEndDate   = "DTEND"
	tokenTimezone  = "TZID"

	timeFormatDefault = "20060102T150405"
	timeFormatUTC     = "20060102T150405Z"
)

type ICS struct {
	url string
}

func (s ICS) Events() ([]calendar.Event, error) {
	resp, err := http.Get(s.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	events := make([]calendar.Event, 0)
	var event calendar.Event

	for scanner.Scan() {
		entry := scanner.Text()

		switch {
		case entry == entryBeginEvent:
			event = calendar.Event{}
		case strings.HasPrefix(entry, tokenStartDate):
			if event.Start, err = s.parseTime(entry); err != nil {
				continue
			}
		case strings.HasPrefix(entry, tokenEndDate):
			if event.End, err = s.parseTime(entry); err != nil {
				continue
			}
		case entry == entryEndEvent && event.IsValid():
			events = append(events, event)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s ICS) parseTime(entry string) (time.Time, error) {
	components := strings.SplitN(entry, sepComponent, 2)
	if len(components) < 2 {
		return time.Time{}, fmt.Errorf("invalid entry: %s", entry)
	}

	if strings.HasSuffix(components[1], "Z") {
		return time.ParseInLocation(timeFormatUTC, components[1], time.UTC)
	}

	tokens := strings.Split(components[0], sepToken)
	if len(tokens) > 1 && strings.HasPrefix(tokens[1], tokenTimezone) {
		parts := strings.SplitN(tokens[1], sepKeyValue, 2)
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("invalid token: %s", tokens[1])
		}

		location, err := timezone.LoadLocation(parts[1])
		if err != nil {
			return time.Time{}, err
		}

		return time.ParseInLocation(timeFormatDefault, components[1], location)
	}

	return time.Parse(timeFormatDefault, components[1])
}

func New(url string) ICS {
	return ICS{url: url}
}
