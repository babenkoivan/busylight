package light

import (
	"fmt"
	"github.com/babenkoivan/busylight/status"
)

const (
	Green = iota
	Yellow
	Red
)

type Color int

type Provider interface {
	TurnOn() error
	TurnOff() error
	ChangeColor(color Color) error
}

type Controller struct {
	provider Provider
}

func (c Controller) ProcessStatusTransition(trans status.Transition) error {
	if trans.To == status.Inactive {
		return c.provider.TurnOff()
	}

	if trans.From == status.Inactive {
		if err := c.provider.TurnOn(); err != nil {
			return err
		}
	}

	var color Color
	switch trans.To {
	case status.Active:
		color = Green
	case status.Focused:
		color = Yellow
	case status.Busy:
		color = Red
	default:
		return fmt.Errorf("unknown status %d", trans.To)
	}

	return c.provider.ChangeColor(color)
}

func NewController(provider Provider) Controller {
	return Controller{provider: provider}
}
