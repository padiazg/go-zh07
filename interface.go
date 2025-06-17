package zh07

// SensorInterface defines the common interface for ZH07 sensors.
type SensorInterface interface {
	// Init initializes the sensor and sets the communication mode
	Init() error
	// CalculateChecksum computes the checksum for data validation
	CalculateChecksum() int
	// IsReadingValid checks if the received data has a valid checksum
	IsReadingValid() bool
	// Read returns a sensor reading or an error
	Read() (*Reading, error)
}
