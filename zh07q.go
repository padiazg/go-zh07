package zh07

import (
	"bufio"
	"bytes"
	"fmt"
	"time"
)

var _ SensorInterface = (*ZH07q)(nil)

type ZH07q struct {
	data         []byte
	rw           *bufio.ReadWriter
	writeAndRead func(rw *bufio.ReadWriter, c []byte) ([]byte, error)
	write        func(rw *bufio.ReadWriter, c []byte) error
}

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

func (z *ZH07q) Init() error {
	if err := z.write(z.rw, commandSetQAMode); err != nil {
		return err
	}
	time.Sleep(sleepAfterWrite) // wait command to be executed

	return nil
}

// CalculateChecksum calculates the checksum from the payload
func (z *ZH07q) CalculateChecksum() int {
	return calculateChecksum(&z.data)
}

func (z *ZH07q) IsReadingValid() bool {
	return z.CalculateChecksum() == z.getChecksum()
}

func (z *ZH07q) Read() (*Reading, error) {
	var e error
	if z.data, e = z.writeAndRead(z.rw, commandQuery); e != nil {
		return nil, e
	}

	if !z.IsReadingValid() {
		return nil, fmt.Errorf("checksum mistmatch=%X, calculated %X", z.getChecksum(), z.CalculateChecksum())
	}

	r := Reading{
		PM25: byteToInt(z.data[2:4]),
		PM10: byteToInt(z.data[4:6]),
		PM1:  byteToInt(z.data[6:8]),
	}

	return &r, nil
}

func (z *ZH07q) getChecksum() int {
	return int(z.data[8])
}
