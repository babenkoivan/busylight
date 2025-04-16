package govee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/babenkoivan/busylight/light"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	url = "https://openapi.api.govee.com/router/api/v1/device/control"

	timeoutRequest = 10 * time.Second

	capabilityTypeOnOff = "devices.capabilities.on_off"
	capabilityTypeColor = "devices.capabilities.color_setting"

	capabilityInstanceOnOff = "powerSwitch"
	capabilityInstanceColor = "colorRgb"
)

var colorMap = map[light.Color]int{
	light.Green:  65280,
	light.Yellow: 16776960,
	light.Red:    16711680,
}

type Request struct {
	RequestId uuid.UUID `json:"requestId"`
	Payload   Payload   `json:"payload"`
}

type Payload struct {
	SKU        string     `json:"sku"`
	Device     string     `json:"device"`
	Capability Capability `json:"capability"`
}

type Capability struct {
	Type     string `json:"type"`
	Instance string `json:"instance"`
	Value    int    `json:"value"`
}

type Govee struct {
	apiKey string
	sku    string
	device string
	client *http.Client
}

func (g Govee) TurnOn() error {
	return g.controlDevice(Capability{
		Type:     capabilityTypeOnOff,
		Instance: capabilityInstanceOnOff,
		Value:    1,
	})
}

func (g Govee) TurnOff() error {
	return g.controlDevice(Capability{
		Type:     capabilityTypeOnOff,
		Instance: capabilityInstanceOnOff,
		Value:    0,
	})
}

func (g Govee) ChangeColor(color light.Color) error {
	mappedColor, ok := colorMap[color]
	if !ok {
		return fmt.Errorf("color %d not recognized", color)
	}

	return g.controlDevice(Capability{
		Type:     capabilityTypeColor,
		Instance: capabilityInstanceColor,
		Value:    mappedColor,
	})
}

func (g Govee) controlDevice(capability Capability) error {
	data := Request{
		RequestId: uuid.New(),
		Payload: Payload{
			SKU:        g.sku,
			Device:     g.device,
			Capability: capability,
		},
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Govee-API-Key", g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	_, err = g.client.Do(req)
	return err
}

func New(apiKey, sku, device string) Govee {
	return Govee{
		apiKey: apiKey,
		sku:    sku,
		device: device,
		client: &http.Client{Timeout: timeoutRequest},
	}
}
