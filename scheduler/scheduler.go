package scheduler

import (
	"context"
	"log/slog"
	"time"
)

func Every(ctx context.Context, tick time.Duration, fn func()) {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Error(ctx.Err().Error())
			return
		case <-ticker.C:
			fn()
		}
	}
}
