package status

import "encoding/json"

type CO2Index uint8

const (
	ExcellentCO2 CO2Index = iota
	GoodCO2
	FairCO2
	PoorCO2
	BadCO2
	UnknownCO2
)

var CO2IndexStrings = [...]string{
	"Excellent",
	"Good",
	"Fair",
	"Poor",
	"Bad",
	"Unknown CO2",
}

func ToCO2Index(co2 uint16) CO2Index {
	switch {
	case co2 < 600:
		return ExcellentCO2
	case co2 < 800:
		return GoodCO2
	case co2 < 1000:
		return FairCO2
	case co2 < 1500:
		return PoorCO2
	case co2 >= 1500:
		return BadCO2
	default:
		return UnknownCO2
	}
}

func (c CO2Index) String() string {
	if c < ExcellentCO2 || c > UnknownCO2 {
		return "Unknown CO2"
	}
	return CO2IndexStrings[c]
}

func (c CO2Index) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
