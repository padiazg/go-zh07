package zh07

import (
	"bufio"
	"fmt"
	"time"
)

type Config struct {
	RW *bufio.ReadWriter
}

type Reading struct {
	PM1  int // Mass Concentration PM1.0   [g/m3]
	PM25 int // Mass Concentration PM2.5   [g/m3]
	PM10 int // Mass Concentration PM10    [g/m3]
}

var (
	commandSetInitiativeUploadMode = []byte{
		0xFF,
		0x01,
		0x78,
		0x40, // first byte for initiative mode
		0x00,
		0x00,
		0x00,
		0x00,
		0x47, // second byte for initiative mode
	}

	commandSetQAMode = []byte{
		0xFF,
		0x01,
		0x78,
		0x41, // first byte for q&a mode
		0x00,
		0x00,
		0x00,
		0x00,
		0x46, // second byte for q&a mode
	}

	commandQuery = []byte{ // q&a mode - query the sensor
		0xFF,
		0x01,
		0x86,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x79,
	}

	commandDormantEnter = []byte{ // enter dormant mode
		0xFF,
		0x01,
		0xA7,
		0x01,
		0x00,
		0x00,
		0x00,
		0x00,
		0x57,
	}

	commandDormantQuit = []byte{ // quit dormant mode
		0xFF,
		0x01,
		0xA7,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x58,
	}

	sleepAfterWrite = 250 * time.Millisecond
)

func calculateChecksum(d *[]byte) int {
	var tempq byte
	for _, v := range (*d)[1 : len(*d)-1] {
		tempq += v
	}
	return int((^tempq) + 1)
}

// byteToInt converts 2 bytes to Int
func byteToInt(data []byte) int {
	return int(data[1]) + (int(data[0]) << 8)
}

// toHex formats hex values to string
func toHex(data []byte) string {
	var result = ""

	for _, c := range data {
		result = fmt.Sprintf("%s%#02x ", result, c)
	}

	return result
}

// writeAndRead writes a command to the sensor and returns the response
func writeAndRead(rw *bufio.ReadWriter, c []byte) ([]byte, error) {
	write(rw, c)
	time.Sleep(sleepAfterWrite) // wait for the response

	r := make([]byte, 9)              // buffer to receive response
	if _, e := rw.Read(r); e != nil { // read response from tty
		return nil, e
	}

	return r, nil
}

// write writes to the sensor
func write(rw *bufio.ReadWriter, c []byte) error {
	if _, e := rw.Write(c); e != nil { // send command
		return e
	}

	return rw.Writer.Flush() // flush write buffer
}
