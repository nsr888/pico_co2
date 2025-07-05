package types

// Readings represents sensor data
type Readings struct {
	CO2         uint16  `json:"eco2"`
	CO2String   string  `json:"eco2_human"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Description string  `json:"description,omitempty"`
	IsValid     bool    `json:"is_valid"` // Indicates if the readings are valid
}
