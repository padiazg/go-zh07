package zh07

import (
	"bufio"
	"bytes"
	"testing"
)

var sampleInitiative = []byte{
	0x42,       // Start byte 1
	0x4D,       // Start byte 2
	0x00, 0x1C, // Frame length
	0x00, 0x54, // Data  1 | Reserved
	0x00, 0x6E, // Data  2 | Reserved
	0x00, 0x7C, // Data  3 | Reserved
	0x00, 0x54, // Data  4 | PM1.0 concentration (ug/m3)
	0x00, 0x6E, // Data  5 | PM2.5 concentration (ug/m3)
	0x00, 0x7C, // Data  6 | PM10 concentration (ug/m3)
	0x00, 0x00, // Data  7 | Reserved
	0x00, 0x00, // Data  8 | Reserved
	0x00, 0x00, // Data  9 | Reserved
	0x00, 0x00, // Data 10 | Reserved
	0x00, 0x00, // Data 11 | Reserved
	0x00, 0x00, // Data 12 | Reserved
	0x00, 0x00, // Data 13 | Reserved
	0x03, 0x27, // Checksum
} // sampleInitiative ...

var badChecksum = []byte{
	0x42,       // Start byte 1
	0x4D,       // Start byte 2
	0x00, 0x1C, // Frame length
	0x00, 0x54, // Data  1 | Reserved
	0x00, 0x6E, // Data  2 | Reserved
	0x00, 0x7C, // Data  3 | Reserved
	0x00, 0x54, // Data  4 | PM1.0 concentration (ug/m3)
	0x00, 0x6E, // Data  5 | PM2.5 concentration (ug/m3)
	0x00, 0x7C, // Data  6 | PM10 concentration (ug/m3)
	0x00, 0x00, // Data  7 | Reserved
	0x00, 0x00, // Data  8 | Reserved
	0x00, 0x00, // Data  9 | Reserved
	0x00, 0x00, // Data 10 | Reserved
	0x00, 0x00, // Data 11 | Reserved
	0x00, 0x00, // Data 12 | Reserved
	0x00, 0x00, // Data 13 | Reserved
	0x03, 0x28, // Checksum
} // badChecksum ...

const (
	fameworkLength = 0x001C // 28
	checksum       = 0x0327
)

func TestGetFrameLengthInitiative(t *testing.T) {
	var z *ZH07i = &ZH07i{data: sampleInitiative}
	if fl := z.getFrameLength(); fl != fameworkLength {
		t.Errorf("TestGetFrameLengthInitiative, got %d, expected %d", fl, fameworkLength)
	}
}

func TestGetChecksumInitiative(t *testing.T) {
	var z *ZH07i = &ZH07i{data: sampleInitiative}
	if cs := z.getChecksum(); cs != checksum {
		t.Errorf("TestGetChecksumInitiative, got %d, expected %d", cs, checksum)
	}
}

func TestCalculateChecksumInitiative(t *testing.T) {
	var z *ZH07i = &ZH07i{data: sampleInitiative}
	if cs := z.CalculateChecksum(); cs != z.getChecksum() {
		t.Errorf("TestCalculateChecksumInitiative, got %d, expected %d", cs, checksum)
	}
}

func TestIsReadingValidInitiative(t *testing.T) {
	var z *ZH07i = &ZH07i{data: sampleInitiative}
	if !z.IsReadingValid() {
		t.Errorf("TestIsReadingValidInitiative, got false, expected true")
	}
}

func TestIsReadingInvalidInitiative(t *testing.T) {
	var z *ZH07i = &ZH07i{data: badChecksum}
	if z.IsReadingValid() {
		t.Errorf("TestIsReadingInvalidInitiative, got true, expected false")
	}
}

func readWriterFromBytes(b []byte) *bufio.ReadWriter {
	return bufio.NewReadWriter(bufio.NewReader(bytes.NewReader(b)), nil)
}

func TestReadInitiative(t *testing.T) {
	var z *ZH07i
	var e error
	var r0 *Reading

	rw := bufio.NewReadWriter(bufio.NewReader(bytes.NewReader(sampleInitiative)), nil)

	// check successful read

	z = &ZH07i{}
	r0, _ = z.Read(rw)

	if r0.MassPM1 != 0x54 {
		t.Errorf("TestReadInitiative PM1.0, got %d, expected %d", r0.MassPM1, 0x54)
	}

	if r0.MassPM25 != 0x6E {
		t.Errorf("TestReadInitiative PM2.5, got %d, expected %d", r0.MassPM1, 0x6E)
	}

	if r0.MassPM10 != 0x7C {
		t.Errorf("TestReadInitiative PM10, got %d, expected %d", r0.MassPM10, 0x7C)
	}

	r0, e = (&ZH07i{}).Read(readWriterFromBytes([]byte{0x00}))
	if r0 != nil || e != nil {
		t.Errorf("TestReadInitiative, expected to skip reading")
	}

	r0, e = (&ZH07i{}).Read(readWriterFromBytes([]byte{0x42, 0x00}))
	if r0 != nil || e != nil {
		t.Errorf("TestReadInitiative, expected to skip reading")
	}

	r0, e = (&ZH07i{}).Read(readWriterFromBytes([]byte{0x42, 0x4d, 0x00, 0x00}))
	if r0 != nil || e != nil {
		t.Errorf("TestReadInitiative, expected to skip reading")
	}

	r0, e = (&ZH07i{}).Read(readWriterFromBytes([]byte{0x42, 0x4d, 0x00, 0x1C, 0x00}))
	if r0 != nil || e == nil {
		t.Errorf("TestReadInitiative, expected error")
	}
}
