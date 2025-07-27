package status

import "encoding/json"

// https://en.wikipedia.org/wiki/Heat_index#Formula
func HeatIndex(tempC, rh float32) float32 {
	if tempC < 27 {
		return tempC
	}

	T := tempC
	R := rh

	T2 := tempC * tempC
	R2 := rh * rh

	// coefficients for °C
	const (
		c1 float32 = -8.78469475556
		c2 float32 = 1.61139411
		c3 float32 = 2.33854883889
		c4 float32 = -0.14611605
		c5 float32 = -0.012308094
		c6 float32 = -0.0164248277778
		c7 float32 = 0.002211732
		c8 float32 = 0.00072546
		c9 float32 = -0.000003582
	)

	resultC := c1 + c2*T + c3*R + c4*T*R + c5*T2 + c6*R2 + c7*T2*R + c8*T*R2 + c9*T2*R2

	return resultC
}

type HeatIndexStatus uint8

const (
	NoHeat HeatIndexStatus = iota
	Caution
	ExtremeCaution
	Danger
	ExtremeDanger
	UnknownHeatIndex
)

var HeatIndexStatusStrings = [...]string{
	"No heat",
	"Caution",
	"Extreme caution",
	"Danger",
	"Extreme danger",
	"Unknown Heat Index",
}

func ToHeatIndexStatus(heatIndex float32) HeatIndexStatus {
	switch {
	case heatIndex < 27:
		return NoHeat
	case heatIndex < 32:
		return Caution
	case heatIndex < 41:
		return ExtremeCaution
	case heatIndex < 54:
		return Danger
	default:
		return ExtremeDanger
	}
}

func (h HeatIndexStatus) String() string {
	if h < NoHeat || h > UnknownHeatIndex {
		return "Unknown Heat Index"
	}
	return HeatIndexStatusStrings[h]
}

func (h HeatIndexStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}
