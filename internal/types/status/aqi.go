package status

import "encoding/json"

type AQIIndex uint8

const (
	Excellent AQIIndex = iota
	Good
	Moderate
	Poor
	Unhealthy
	UnknownAQI
)

var AQIIndexStrings = [...]string{
	"Excellent",
	"Good",
	"Moderate",
	"Poor",
	"Unhealthy",
	"Unknown AQI",
}

func ToAQIIndex(aqi uint8) AQIIndex {
	switch {
	case aqi == 0:
		return Excellent
	case aqi == 1:
		return Good
	case aqi == 2:
		return Moderate
	case aqi == 3:
		return Poor
	case aqi == 4:
		return Unhealthy
	default:
		return UnknownAQI
	}
}

func (a AQIIndex) String() string {
	if a < Excellent || a > UnknownAQI {
		return "Unknown AQI"
	}
	return AQIIndexStrings[a]
}

func (a AQIIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}
