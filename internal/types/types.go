package types

import "time"

// Readings represents sensor data
type Readings struct {
	CO2         uint16    `json:"eco2"`
	Temperature float32   `json:"temperature"`
	Humidity    float32   `json:"humidity"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}
