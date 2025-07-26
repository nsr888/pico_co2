package types

type ENSRawReadings struct {
	CO2                 uint16  `json:"eco2"`
	TVOC                uint16  `json:"tvoc"`
	AQI                 uint8   `json:"aqi"`
	Temperature         float32 `json:"temperature"`
	Humidity            float32 `json:"humidity"`
	DataValidityWarning string  `json:"data_validity_warning,omitempty"`
}
