package types

// CalculateComfortIndex computes a 7-bar comfort index (-3 to +3) using temperature (°C) and humidity (%)
func (r *Readings) CalculateComfortIndex() int16 {
	var adjustedTemp float32
	T := r.Temperature
	RH := r.Humidity

	adjustedTemp = T

	if T >= 27.0 {
		adjustedTemp = HeatIndex(T, RH)
	}

	// Map to 7-bar scale using adjusted temperature thresholds
	// Thresholds over 27°C are based on heat index calculations
	// Other thresholds are based on general comfort levels
	switch {
	case adjustedTemp >= 54.0:
		return 4 // Extreme danger: heat cramps and heat exhaustion are likely; heat stroke is imminent with continued activity.
	case adjustedTemp >= 41.0:
		return 3 // Danger: heat cramps and heat exhaustion are likely; heat stroke is probable with continued activity.
	case adjustedTemp >= 32.0:
		return 2 // Extreme caution: heat cramps and heat exhaustion are possible. Continuing activity could result in heat stroke.
	case adjustedTemp >= 27.0:
		return 1 // Caution: fatigue is possible with prolonged exposure and activity. Continuing activity could result in heat cramps.
	case adjustedTemp >= 20.0:
		return 0 // Neutral (comfortable)
	case adjustedTemp >= 16.0:
		return -1 // Slightly cool
	default:
		return -2 // Cool to cold
	}
}
