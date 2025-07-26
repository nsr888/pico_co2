package status

// CalculateComfortIndex computes a 7-bar comfort index (-3 to +3) using temperature (°C) and humidity (%)
func CalculateComfortIndex(T, RH float32) int16 {
	var adjustedTemp float32

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

// ComfortStatus returns a human-readable comfort status based on sensor readings.
func ComfortStatus(co2 uint16, aqi uint8, heatIndex float32, humidity float32, temperature float32) string {
	switch {
	case co2 < 1000 && aqi >= 3:
		return "Poor Air"
	case co2 >= 1000 || aqi >= 3:
		return "High CO2"
	case heatIndex >= 54:
		return "Danger heat"
	case heatIndex >= 41:
		return "Extreme heat"
	case heatIndex >= 32:
		return "Very heat"
	case heatIndex >= 27:
		return "Heat"
	case humidity > 65:
		return "High humidity"
	case humidity < 35:
		return "Dry"
	case temperature < 18:
		return "Cold"
	case co2 < 800 &&
		aqi <= 2 &&
		temperature >= 18 && temperature <= 25 &&
		humidity >= 35 && humidity <= 60:
		return "Comfort"
	default:
		return "Normal"
	}
}

func TempComfortIndex(temp float32) int16 {
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

// HumidityComfortIndex returns an index from -3 to 4 indicating human comfort
// levels based on humidity readings. See the following source for details:
// https://www.myqualitycomfort.com/general/best-humidity-level-for-home/
func HumidityComfortIndex(humidity float32) int16 {
	switch {
	case humidity > 80:
		return 4 // Extremely high humidity, risk of mold and discomfort
	case humidity > 70:
		return 3 // Uncomfortable humidity
	case humidity > 60:
		return 2 // Very humid
	case humidity > 50:
		return 1 // Slightly humid
	case humidity > 40:
		return 0 // Comfortable range
	case humidity > 30:
		return -1 // Slightly dry
	default:
		return -2 // Very dry, risk of respiratory issues
	}
}
