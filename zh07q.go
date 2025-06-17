package zh07

import (
	"bufio"
	"bytes"
	"fmt"
	"time"
)

var _ SensorInterface = (*ZH07q)(nil)

// ZH07q implements the SensorInterface for question and answer mode.
// In this mode, readings are requested on demand.
type ZH07q struct {
	data         []byte
	rw           *bufio.ReadWriter
	writeAndRead func(rw *bufio.ReadWriter, c []byte) ([]byte, error)
	write        func(rw *bufio.ReadWriter, c []byte) error
}

// NewZH07q creates a new ZH07q sensor instance for question and answer mode.
func NewZH07q(config *Config) *ZH07q {
	if config == nil {
		config = &Config{}
	}

	if config.RW == nil {
		config.RW = bufio.NewReadWriter(bufio.NewReader(bytes.NewReader([]byte{})), nil)
	}

	return &ZH07q{
		rw:           config.RW,
		writeAndRead: writeAndRead,
		write:        write,
	}
}

// Init initializes the sensor for question and answer mode.
func (z *ZH07q) Init() error {
	if err := z.write(z.rw, commandSetQAMode); err != nil {
		return err
	}
	time.Sleep(sleepAfterWrite) // wait command to be executed

	return nil
}

// CalculateChecksum calculates the checksum from the payload.
func (z *ZH07q) CalculateChecksum() int {
	return calculateChecksum(&z.data)
}

// IsReadingValid checks if the calculated checksum matches the payload checksum.
func (z *ZH07q) IsReadingValid() bool {
	return z.CalculateChecksum() == z.getChecksum()
}

// Read sends a query command and reads particulate matter data from the sensor.
func (z *ZH07q) Read() (*Reading, error) {
	var e error
	if z.data, e = z.writeAndRead(z.rw, commandQuery); e != nil {
		return nil, e
	}

	if !z.IsReadingValid() {
		return nil, fmt.Errorf("%w: received=%X, calculated=%X", ErrChecksumMismatch, z.getChecksum(), z.CalculateChecksum())
	}

	r := Reading{
		PM1:  byteToInt(z.data[6:8]),
		PM25: byteToInt(z.data[2:4]),
		PM10: byteToInt(z.data[4:6]),
	}

	return &r, nil
}

// getChecksum recovers the checksum from the payload.
func (z *ZH07q) getChecksum() int {
	return int(z.data[8])
}
