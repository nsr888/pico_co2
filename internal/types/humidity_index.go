package types

// HumidityComfortIndex returns an index from -3 to 4 indicating human comfort
// levels based on humidity readings. See the following source for details:
// https://www.myqualitycomfort.com/general/best-humidity-level-for-home/
func (r *Readings) HumidityComfortIndex() int16 {
	humidity := r.Humidity

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
