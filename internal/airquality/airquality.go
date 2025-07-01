package airquality

// AirQualitySensor defines the standard interface for a sensor module that
// provides environmental readings.
type AirQualitySensor interface {
	Configure() error
	Read() (*Readings, error)
	Temperature() float32
	Humidity() float32
	CO2() uint16
}

// Readings represents airquality data
type Readings struct {
	CO2               uint16      `json:"co2"`
	CO2Interpretation string      `json:"co2_interpretation"`
	Temperature       float32     `json:"temperature"`
	Humidity          float32     `json:"humidity"`
	Quality           DataQuality `json:"quality"`
}

func (r Readings) Interpretation() string {
	if r.Quality.Status != StatusOK {
		return r.Quality.Status.String()
	}

	return r.CO2Interpretation
}

type QualityStatus string

const (
	StatusOK        QualityStatus = "ok"
	StatusWarmUp    QualityStatus = "warming_up"
	StatusStartUp   QualityStatus = "start_up"
	StatusNotValid  QualityStatus = "not_valid"
)

func (qs QualityStatus) String() string {
	switch qs {
	case StatusOK:
		return "OK"
	case StatusWarmUp:
		return "Warming up"
	case StatusStartUp:
		return "Start up"
	case StatusNotValid:
		return "Not Valid"
	default:
		return "Unknown"
	}
}

type DataQuality struct {
	Status      QualityStatus `json:"status"`
	Description string        `json:"description"`
}
