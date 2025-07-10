package types

// Readings represents sensor data
type Readings struct {
	CO2             uint16  `json:"eco2"`
	AQI             uint8   `json:"aqi"`
	CO2Status       string  `json:"eco2_human"`
	Temperature     float32 `json:"temperature"`
	Humidity        float32 `json:"humidity"`
	ValidityError   string  `json:"validity_notes,omitempty"`
	HeatIndex       float32 `json:"heat_index"`
}

// ComfortStatus returns a human-readable comfort status based on sensor readings.
// Priority order:
//  1. If there is a validity error, return that error.
//  2. If CO2 ≥ 1500 ppm or AQI ≥ 4:                 return "Poor air"
//  3. If CO2 ≥ 1000 ppm or AQI ≥ 3:                 return "Ventilate"
//  4. If HeatIndex ≥ 27 °C:                         return "Heat"
//  5. If Humidity > 65%:                            return "High humidity"
//  6. If Humidity < 35%:                            return "Low humidity"
//  7. If Temperature < 18 °C:                       return "Cold"
//  8. If Temperature > 30 °C:                       return "Hot"
//  9. If CO2 < 800 ppm && AQI ≤ 2 &&                 // all conditions within comfort zone
//     Temperature ≥ 18 °C && Temperature ≤ 25 °C &&
//     Humidity ≥ 35% && Humidity ≤ 60%:              return "Comfort"
//  10. Otherwise (e.g. CO2 800–999 ppm, AQI 1–2,      // промежуточная зона
//     temp 25–30 °C, humid 35–65%):                  return "Normal"
func (r *Readings) ComfortStatus() string {
	// TODO: Implement validation logic for sensor readings
	// 1. Validate the sensor reading
	// if err := r.Validate(); err != nil {
	//     return fmt.Sprintf("Error: %v", err)
	// }

	// 2. Very poor air quality
	if r.CO2 >= 1500 || r.AQI >= 4 {
		return "Poor air quality"
	}

	// 3. Need to ventilate
	if r.CO2 >= 1000 || r.AQI >= 3 {
		return "Ventilate the room"
	}

	// 4. Heat discomfort
	if r.HeatIndex >= 27 {
		return "Heat discomfort"
	}

	// 5. High humidity (muggy)
	if r.Humidity > 65 {
		return "High humidity"
	}

	// 6. Low humidity (dry)
	if r.Humidity < 35 {
		return "Low humidity"
	}

	// 7. Cold temperature
	if r.Temperature < 18 {
		return "Cold temperature"
	}

	// 8. Very hot temperature
	if r.Temperature > 30 {
		return "Hot temperature"
	}

	// 9. All parameters within defined comfort ranges
	if r.CO2 < 800 &&
		r.AQI <= 2 &&
		r.Temperature >= 18 && r.Temperature <= 25 &&
		r.Humidity >= 35 && r.Humidity <= 60 {
		return "Comfort zone"
	}

	// 10. Intermediate zone (neither comfortable nor critical)
	return "Normal conditions"
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

// CO2Rating returns a CO2 rating based on the CO2 value from 0 to 4.
func (r *Readings) CO2Rating() int16 {
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
