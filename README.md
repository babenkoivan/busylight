# Busy Light

This lightweight library enables a busy light integration for your smart home, helping you stay aware of your schedule at a glance.

The light syncs with your calendar and changes color based on upcoming or ongoing events:

🟡 Yellow – event starting in 10 minutes

🔴 Red – event currently in progress

The light is only active when there is an upcoming or ongoing event. Events that take more than 10 hours are ignored.

## Current Support

* Calendar source: ICS format
* Smart bulb provider: Govee (e.g. `H6006`)

Additional calendar formats and smart home devices may be supported in the future.

## Usage

Refer to the application code listed below for implementation details.

```go
package main

import (
	"context"
	"github.com/babenkoivan/busylight/calendar"
	"github.com/babenkoivan/busylight/calendar/ics"
	"github.com/babenkoivan/busylight/light"
	"github.com/babenkoivan/busylight/light/govee"
	"github.com/babenkoivan/busylight/scheduler"
	"github.com/babenkoivan/busylight/status"
	"log/slog"
	"time"
)

func main() {
	ctx := context.Background()

	// Initialize an ICS calendar source
	src := ics.New("https://my.calendar.io/public.ics")

	// Create a calendar instance and sync it periodically with the source
	cal := calendar.New()
	go scheduler.Every(ctx, 5*time.Minute, func() {
		if err := cal.Sync(src); err != nil {
			slog.Error(err.Error())
		}
	})

	// Create a status instance and a tracker
	stat := status.New()
	tracker := status.NewTracker()

	// Sync status with the calendar and track changes every minute
	go scheduler.Every(ctx, time.Minute, func() {
		stat.Sync(cal)
		tracker.Record(*stat)
	})

	// Initialize a light controller with Govee provider
	provider := govee.New("my_api_key", "H6006", "64:09:C5:32:37:36:2D:13")
	lightCtl := light.NewController(provider)

	// React to status transitions and control the light accordingly
	for trans := range tracker.C {
		if err := lightCtl.ProcessStatusTransition(trans); err != nil {
			slog.Error(err.Error())
		}
	}
}
```