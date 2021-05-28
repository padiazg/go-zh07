package zh07

import (
	"bufio"
	"fmt"
	"io"
	// "golang.org/x/tools/cmd/guru/serial"
)

// ZH07i data read in initiative upload mode
type ZH07i struct {
	data []byte
}

// GetFrameLength returns the data frame length
func (z *ZH07i) getFrameLength() int {
	return ByteToInt((z.data)[2:5])
} // ZH07i.getFrameLength ...

// GetChecksum recovers the checksun received in the payload
func (z *ZH07i) getChecksum() int {
	return ByteToInt(z.data[30:])
} // ZH07i.getChecksum ...

// CalculateChecksum calculates the checksum from the payload
func (z *ZH07i) CalculateChecksum() int {
	// we calculate adding all the first 30 bytes of the data received,
	// the last 2 bytes are the checkcum
	var r0 int
	for _, v := range z.data[:30] {
		r0 += int(v)
	}
	return r0
} // ZH07i.CalculateChecksum ...

func (z *ZH07i) IsReadingValid() bool {
	return z.CalculateChecksum() == ByteToInt(z.data[30:32])
} // ZH07i.IsReadingValid ...

func (z *ZH07i) Read(rw *bufio.ReadWriter) (*Reading, error) {
	// read byte by byte until we find the 1st start character 0x42
	b0, _ := rw.ReadByte()
	if b0 != 0x42 {
		fmt.Printf("First character != 0x42 [%v]\n", b0)
		return nil, nil
	}
	// fmt.Printf("Start character 1 (0x42) detected\n")

	// then we read the next 3 bytes, which should be:
	//   2nd character start 0x4d
	//   frame length high bits
	//   frame length low bits
	b1, _ := rw.Peek(3)
	if b1[0] != 0x4d {
		return nil, nil
	}
	// fmt.Printf("Start character 2 (0x4d) detected\n")

	// frame length must be 0x001C => 28
	if ByteToInt(b1[1:3]) != 28 {
		return nil, nil
	}

	// if everything matches till now, we read the remaining data (31 bytes)
	// as peek doesnt moves the index, we'll be reading those 3 bytes again:
	// 2nd start character + frame len high + frame len low + data frame
	data := make([]byte, 31)
	if _, err := io.ReadFull(rw, data); err != nil {
		return nil, err
	}
	// fmt.Printf("Full data frame found\n")

	// here we prepend the first byte of the start character so we
	// can have all payload in one single slice
	z.data = append([]byte{b0}, data...)

	r := Reading{
		MassPM1:  ByteToInt(z.data[10:12]),
		MassPM25: ByteToInt(z.data[12:14]),
		MassPM10: ByteToInt(z.data[14:16]),
	}

	return &r, nil
} // ZH07i.Read ...
