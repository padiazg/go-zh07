package zh07

type SensorInterface interface {
	Init() error
	CalculateChecksum() int
	IsReadingValid() bool
	Read() (*Reading, error)
	getChecksum() int
}
