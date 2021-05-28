package zh07

import (
	"bufio"
	"fmt"
	"log"
)

type ZH07q struct {
	data []byte
}

// CalculateChecksum calculates the checksum from the payload
func (z *ZH07q) CalculateChecksum() int {
	return CalculateChecksum(&z.data)
} // ZH07q.CalculateChecksum ...

func (z *ZH07q) getChecksum() int {
	return int(z.data[8])
} // ZH07q.getChecksum ...

func (z *ZH07q) IsReadingValid() bool {
	return z.CalculateChecksum() == z.getChecksum()
} // ZH07q.IsReadingValid ...

func (z *ZH07q) Read(rw *bufio.ReadWriter) (*Reading, error) {
	command := []byte{0xFF, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
	var e error
	if z.data, e = SendCommand(rw, &command); e != nil {
		return nil, e
	}

	if !z.IsReadingValid() {
		log.Printf("Read | z.data => %v\n", ToHex(z.data))
		log.Printf("Read | Checksum: %X Calculated:%X \n", z.getChecksum(), z.CalculateChecksum())
		return nil, fmt.Errorf("checksum mistmatch")
	}

	r := Reading{
		MassPM25: ByteToInt(z.data[2:4]),
		MassPM10: ByteToInt(z.data[4:6]),
		MassPM1:  ByteToInt(z.data[6:8]),
	}

	return &r, nil
} // ZH07q.Read ...
