package zh07

import (
	"testing"
)

var sampleChecksumQA []byte = []byte{0xFF, 0x86, 0x00, 0x47, 0x00, 0xC7, 0x03, 0x0F, 0x5A}
var sampleReadingQA []byte = []byte{0xFF, 0x86, 0x00, 0x85, 0x00, 0x96, 0x00, 0x65, 0xFA}

func TestGetChecksumQA(t *testing.T) {
	var z *ZH07q = &ZH07q{data: sampleChecksumQA}
	if cs := z.getChecksum(); cs != 0x5A {
		t.Errorf("TestGetChecksumQA, got %d, expected %d", cs, checksum)
	}
} // TestGetChecksumQA ...

func TestCalculateChecksumQA(t *testing.T) {
	var z *ZH07q = &ZH07q{data: sampleChecksumQA}
	if cs := z.CalculateChecksum(); cs != z.getChecksum() {
		t.Errorf("TestCalculateChecksumInitiative, got %d, expected %d", cs, checksum)
	}
} // TestCalculateChecksumQA ...

func TestIsReadingValidQA(t *testing.T) {
	var z *ZH07q = &ZH07q{data: sampleReadingQA}
	if !z.IsReadingValid() {
		t.Errorf("TestIsReadingValidInitiative, got false, expected true")
	}
} // TestIsReadingValidQA ...

// func TestReadQA(t *testing.T) {
// 	var z *ZH07q = &ZH07q{}
// 	var e error
// 	var r0 *Reading

// 	b := new(bytes.Buffer)
// 	w := bufio.NewWriter(b)
// 	r := bufio.NewReader(b)
// 	rw := bufio.NewReadWriter(r, w)

// 	r0, e = z.Read(rw)
// 	if e != nil {
// 		t.Error(e)
// 	}
// 	t.Error(r0)
// }
