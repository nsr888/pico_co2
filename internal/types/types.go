package types

// Readings represents sensor data
type Readings struct {
	CO2                 uint16  `json:"eco2"`
	CO2Status           string  `json:"eco2_human"`
	Temperature         float32 `json:"temperature"`
	Humidity            float32 `json:"humidity"`
	ValidityWaitMinutes int64   `json:"validity_minutes"`
	ValidityError       string  `json:"validity_notes,omitempty"`
	HeatIndex           float32 `json:"heat_index"`
	HeatIndexStatus     string  `json:"heat_index_human"`
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
func (r *Readings) HeatIndexRating() int {
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

// CO2Rating returns a CO2 rating based on the CO2 value from 0 to 4.
func (r *Readings) CO2Rating() int {
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
