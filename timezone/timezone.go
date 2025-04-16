package timezone

import "time"

func LoadLocation(name string) (*time.Location, error) {
	if n, ok := windowsZones[name]; ok {
		name = n
	}

	return time.LoadLocation(name)
}
