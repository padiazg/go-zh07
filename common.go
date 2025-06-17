// Package zh07 provides a driver for Winsen ZH06 and ZH07 laser dust sensors.
//
// The ZH06 and ZH07 are laser dust sensor modules used to check air quality by
// measuring particulate matter concentrations (PM1.0, PM2.5, and PM10).
//
// This package supports two communication modes:
//   - Initiative upload mode: The sensor continuously broadcasts readings
//   - Question and answer mode: Readings are requested on demand
//
// Example usage:
//
//	// Create a sensor instance for Q&A mode
//	sensor := zh07.NewZH07q(&zh07.Config{RW: rw})
//	if err := sensor.Init(); err != nil {
//		log.Fatal(err)
//	}
//
//	// Read sensor data
//	reading, err := sensor.Read()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("PM1.0: %d, PM2.5: %d, PM10: %d\n", reading.PM1, reading.PM25, reading.PM10)
package zh07

import (
	"bufio"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrChecksumMismatch is returned when the calculated checksum doesn't match the received checksum
	ErrChecksumMismatch = errors.New("checksum mismatch")
	// ErrInvalidFrame is returned when the received data frame is invalid
	ErrInvalidFrame = errors.New("invalid data frame")
	// ErrSensorCommunication is returned when communication with the sensor fails
	ErrSensorCommunication = errors.New("sensor communication failed")
)

// Config holds configuration options for sensor instances.
type Config struct {
	// RW is the ReadWriter interface for communicating with the sensor
	RW *bufio.ReadWriter
}

// Reading represents a sensor reading with particulate matter concentrations.
type Reading struct {
	PM1  int // Mass Concentration PM1.0 [μg/m³]
	PM25 int // Mass Concentration PM2.5 [μg/m³]
	PM10 int // Mass Concentration PM10 [μg/m³]
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

// calculateChecksum computes the checksum for sensor data validation.
func calculateChecksum(d *[]byte) int {
	var tempq byte
	for _, v := range (*d)[1 : len(*d)-1] {
		tempq += v
	}
	return int((^tempq) + 1)
}

// byteToInt converts 2 bytes to int in big-endian format.
func byteToInt(data []byte) int {
	return int(data[1]) + (int(data[0]) << 8)
}

// toHex formats byte slice as hexadecimal string for debugging.
func toHex(data []byte) string {
	var result = ""

	for _, c := range data {
		result = fmt.Sprintf("%s%#02x ", result, c)
	}

	return result
}

// writeAndRead writes a command to the sensor and returns the response.
func writeAndRead(rw *bufio.ReadWriter, c []byte) ([]byte, error) {
	write(rw, c)
	time.Sleep(sleepAfterWrite) // wait for the response

	r := make([]byte, 9)              // buffer to receive response
	if _, e := rw.Read(r); e != nil { // read response from tty
		return nil, e
	}

	return r, nil
}

// write sends a command to the sensor.
func write(rw *bufio.ReadWriter, c []byte) error {
	if _, e := rw.Write(c); e != nil { // send command
		return e
	}

	return rw.Writer.Flush() // flush write buffer
}
