package zh07

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"
)

var _ SensorInterface = (*ZH07i)(nil)

// ZH07i implements the SensorInterface for initiative upload mode.
// In this mode, the sensor continuously broadcasts readings.
type ZH07i struct {
	data  []byte
	rw    *bufio.ReadWriter
	write func(rw *bufio.ReadWriter, c []byte) error
}

// NewZH07i creates a new ZH07i sensor instance for initiative upload mode.
func NewZH07i(config *Config) *ZH07i {
	if config == nil {
		config = &Config{}
	}

	if config.RW == nil {
		config.RW = bufio.NewReadWriter(bufio.NewReader(bytes.NewReader([]byte{})), nil)
	}

	return &ZH07i{
		data:  make([]byte, 32),
		rw:    config.RW,
		write: write,
	}
}

// Init initializes the sensor for initiative upload mode.
func (z *ZH07i) Init() error {
	if err := z.write(z.rw, commandSetInitiativeUploadMode); err != nil {
		return err
	}
	time.Sleep(sleepAfterWrite) // wait command to be executed

	return nil
}

// CalculateChecksum calculates the checksum from the payload.
// The checksum is calculated by adding all the first 30 bytes of the data received;
// the last 2 bytes are the checksum.
func (z *ZH07i) CalculateChecksum() int {
	var r0 int
	for _, v := range z.data[:30] {
		r0 += int(v)
	}
	return r0
}

// IsReadingValid checks if the calculated checksum matches the payload checksum.
func (z *ZH07i) IsReadingValid() bool {
	return z.CalculateChecksum() == byteToInt(z.data[30:32])
}

// Read reads particulate matter data from the sensor in initiative upload mode.
func (z *ZH07i) Read() (*Reading, error) {
	var (
		b0   = make([]byte, 1)
		b1   = make([]byte, 3)
		data = make([]byte, 28)
		err  error
	)

	// read byte by byte until we find the 1st start character (0x42)
	if _, err = z.rw.Read(b0); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSensorCommunication, err)
	}
	if b0[0] != 0x42 {
		return nil, nil
	}
	// then we read the next 3 bytes, which should be:
	//   2nd character start (0x4d)
	//   frame length high bits
	//   frame length low bits
	if _, err = z.rw.Read(b1); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSensorCommunication, err)
	}
	if b1[0] != 0x4d {
		return nil, nil
	}

	// frame length must be 0x00 0x1C => 28
	if byteToInt(b1[1:3]) != 28 {
		return nil, nil
	}

	// if everything matches so far, we read the remaining data (28 bytes)
	if _, err = io.ReadFull(z.rw, data); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSensorCommunication, err)
	}

	// let's concatenate all the bytes read into a single slice
	z.data = append(append(b0, b1...), data...)

	r := Reading{
		PM1:  byteToInt(z.data[10:12]),
		PM25: byteToInt(z.data[12:14]),
		PM10: byteToInt(z.data[14:16]),
	}

	return &r, nil
}

// getChecksum recovers the checksum received in the payload.
func (z *ZH07i) getChecksum() int {
	return byteToInt(z.data[30:])
}
