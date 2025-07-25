package types

// Readings represents sensor data
type Readings struct {
	CO2           uint16  `json:"eco2"`
	AQI           uint8   `json:"aqi"`
	CO2Status     string  `json:"eco2_human"`
	Temperature   float32 `json:"temperature"`
	Humidity      float32 `json:"humidity"`
	ValidityError string  `json:"validity_notes,omitempty"`
	HeatIndex     float32 `json:"heat_index"`
}

func (r *Readings) AQIStatus() string {
	if r == nil {
		return "No data"
	}

	switch r.AQI {
	case 1:
		return "Excellent"
	case 2:
		return "Good"
	case 3:
		return "Moderate"
	case 4:
		return "Poor"
	case 5:
		return "Unhealthy"
	default:
		return "Unknown AQI"
	}
}

// ComfortStatus returns a human-readable comfort status based on sensor readings.
func (r *Readings) ComfortStatus() string {
	switch {
	case r.CO2 < 1000 && r.AQI >= 3:
		return "Poor Air"
	case r.CO2 >= 1000 || r.AQI >= 3:
		return "High CO2"
	case r.HeatIndex >= 54:
		return "Danger heat"
	case r.HeatIndex >= 41:
		return "Extreme heat"
	case r.HeatIndex >= 32:
		return "Very heat"
	case r.HeatIndex >= 27:
		return "Heat"
	case r.Humidity > 65:
		return "High humidity"
	case r.Humidity < 35:
		return "Dry"
	case r.Temperature < 18:
		return "Cold"
	case r.CO2 < 800 &&
		r.AQI <= 2 &&
		r.Temperature >= 18 && r.Temperature <= 25 &&
		r.Humidity >= 35 && r.Humidity <= 60:
		return "Comfort"
	default:
		return "Normal"
	}
}

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

// HeatIndexStatus returns a human-readable status based on the heat index
// value.
func HeatIndexStatus(heatIndex float32) string {
	if heatIndex < 27 {
		return "No heat"
	} else if heatIndex < 32 {
		return "Caution"
	} else if heatIndex < 41 {
		return "Extreme caution"
	} else if heatIndex < 54 {
		return "Danger"
	} else {
		return "Extreme danger"
	}
}

// HeatIndexRating returns the heat index level as an integer from 0 to 4.
func (r *Readings) HeatIndexRating() int16 {
	if r == nil {
		return 0
	}

	switch {
	case r.HeatIndex < 27:
		return 0 // No heat
	case r.HeatIndex < 32:
		return 1 // Caution
	case r.HeatIndex < 41:
		return 2 // Extreme caution
	case r.HeatIndex < 54:
		return 3 // Danger
	default:
		return 4 // Extreme danger
	}
}

// CO2Index returns a CO2 rating based on the CO2 value from 0 to 4.
func (r *Readings) CO2Index() int16 {
	if r == nil {
		return 0
	}

	switch {
	case r.CO2 < 600:
		return 0 // Excellent
	case r.CO2 < 800:
		return 1 // Good
	case r.CO2 < 1000:
		return 2 // Fair
	case r.CO2 < 1500:
		return 3 // Poor
	default:
		return 4 // Bad
	}
}

func (r *Readings) TempComfortIndex() int16 {
	temp := r.Temperature

	switch {
	// Comfort
	case temp >= 22 && temp <= 24:
		return 0

	// Normal
	case temp >= 20 && temp < 22:
		return -1
	case temp > 24 && temp <= 26:
		return 1

	// Caution
	case temp >= 18 && temp < 20:
		return -2
	case temp > 26 && temp <= 28:
		return 2

	// Uncomfortable
	case temp >= 16 && temp < 18:
		return -3
	case temp > 28 && temp <= 32:
		return 3

	// Extreme
	case temp < 16:
		return -3
	case temp > 32:
		return 3

	default:
		return 0
	}
}
