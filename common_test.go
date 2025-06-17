package zh07

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type checkFn func(t *testing.T, r *Reading, e error)

var (
	check = func(fns ...checkFn) []checkFn { return fns }

	hasError = func(want bool) checkFn {
		return func(t *testing.T, _ *Reading, err error) {
			t.Helper()
			if want {
				assert.NotNil(t, err, "hasError: error expected, none produced")
			} else {
				assert.Nil(t, err, "hasError = [+%v], no error expected")
			}
		}
	}

	pm = func(pm25, pm10, pm1 int) checkFn {
		return func(t *testing.T, r *Reading, err error) {
			assert.Equalf(t, pm1, r.PM1, "Expected PM1.0=%d, got %d", pm1, r.PM1)
			assert.Equalf(t, pm25, r.PM25, "Expected PM2.5=%d, got %d", pm25, r.PM25)
			assert.Equalf(t, pm10, r.PM10, "Expected PM10=%d, got %d", pm10, r.PM10)
		}
	}

	isNil = func(t *testing.T, r *Reading, err error) {
		assert.Empty(t, r)
	}
)

func dummyCommandResponder(t *testing.T, r0 *bufio.ReadWriter, c, r []byte) {
	b := make([]byte, 9) // buffer to receive response
	for {
		if _, e := r0.Read(b); e != nil { // read response from tty
			// command received?
			if b[0] == c[0] && b[1] == c[1] && b[8] == c[8] {
				// write response
				_, e := r0.Write(r)
				if e != nil { // write commant to tty
					t.Error(e)
				}
				r0.Writer.Flush()
				break
			}
		}
		time.Sleep(sleepAfterWrite)
	}
}

func Test_writeAndRead(t *testing.T) {
	var (
		command  = []byte{0xFF, 0x86, 0x00, 0x47, 0x00, 0xC7, 0x03, 0x0F, 0x5A}
		response = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
		b        = &bytes.Buffer{}
		// buffer to simulate tty
		rw = bufio.NewReadWriter(
			bufio.NewReader(b),
			bufio.NewWriter(b),
		)
	)

	go dummyCommandResponder(t, rw, command, response)

	res, e0 := writeAndRead(rw, command)
	if e0 != nil {
		t.Errorf("Test_writeAndRead | Sending command: %v", e0)
	}

	for i, v := range res {
		if response[i] != v {
			t.Errorf("Test_writeAndRead | At response index %d. Got %X, expected %X", i, v, response[i])
		}
	}
}

func Test_calculateChecksum(t *testing.T) {
	var data *[]byte = &[]byte{0xFF, 0x86, 0x00, 0x47, 0x00, 0xC7, 0x03, 0x0F, 0x5A}

	if c := calculateChecksum(data); c != 0x5A {
		t.Errorf("Test_calculateChecksum. Got: %d, expexted: %d", c, 0x5A)
	}
}

func Test_toHex(t *testing.T) {
	if v := toHex([]byte{0x00, 0x01, 0x02, 0x0a, 0xff}); v != "0x00 0x01 0x02 0x0a 0xff " {
		t.Errorf("Test_toHex, mistmatch: [%s], expected 0x00 0x01 0x02 0x0a 0xff", v)
	}
}
