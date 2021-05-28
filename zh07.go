package zh07

import (
	"bufio"
	"fmt"
	"log"
	"time"
)

type CommunicationMode int

func (v CommunicationMode) String() string {
	var r string
	switch v {
	case ModeInitiative:
		r = "InitiativeUpload"
	case ModeQA:
		r = "QuestionAndAnswer"
	}
	return r
} // CommunicationMode.String() ...

type Reading struct {
	MassPM1  int // Mass Concentration PM1.0   [g/m3]
	MassPM25 int // Mass Concentration PM2.5   [g/m3]
	MassPM10 int // Mass Concentration PM10    [g/m3]
}

const (
	ModeInitiative CommunicationMode = iota
	ModeQA
)

type SensorInterface interface {
	CalculateChecksum() int
	IsReadingValid() bool
	Read(*bufio.ReadWriter) (*Reading, error)
}

type ZH07 struct {
	mode       CommunicationMode
	zh07       SensorInterface
	readWriter *bufio.ReadWriter
}

func NewZH07(mode CommunicationMode, rw *bufio.ReadWriter) (*ZH07, error) {
	z := &ZH07{
		mode:       mode,
		readWriter: rw,
	} // z ...

	if e := z.setMode(); e != nil {
		return nil, e
	}

	switch mode {
	case ModeInitiative:
		z.zh07 = &ZH07i{}
	case ModeQA:
		z.zh07 = &ZH07q{}
	}

	return z, nil
} // NewZH07 ...

func (z *ZH07) Read() (*Reading, error) {
	return z.zh07.Read(z.readWriter)
} // ZH07.Read ...

func (z *ZH07) setMode() error {
	var setMode []byte = []byte{0xFF, 0x01, 0x78, 0x40, 0x00, 0x00, 0x00, 0x00, 0x47}

	if z.mode == ModeQA {
		setMode[3] = 0x41
		setMode[8] = 0x46
	}
	// log.Printf("setMode | data: %s", utils.ToHex(setMode))

	_, e := z.readWriter.Write(setMode) // send set mode command
	if e != nil {
		return e
	}
	z.readWriter.Writer.Flush() // flush write buffer
	time.Sleep(2 * time.Second) // wait command to be executed
	log.Printf("Set mode to %s", z.mode)

	return nil
} // ZH07.setMode ...

func SendCommand(rw *bufio.ReadWriter, c *[]byte) ([]byte, error) {
	// log.Printf("SendCommand | c => %v", ToHex(*c))
	if _, e := rw.Write(*c); e != nil { // write commant to tty
		return nil, e
	}

	rw.Writer.Flush()                   // flush write buffer
	time.Sleep(1500 * time.Millisecond) // wait for the response
	r := make([]byte, 9)                // buffer to receive response
	if _, e := rw.Read(r); e != nil {   // read response from tty
		return nil, e
	}
	// log.Printf("SendCommand | r => %v", ToHex(r))
	return r, nil
} // SendCommand ...

func CalculateChecksum(d *[]byte) int {
	var tempq byte
	for _, v := range (*d)[1 : len(*d)-1] {
		tempq += v
	}
	return int((^tempq) + 1)
} // CalculateChecksum ...

// ByteToInt converts 2 bytes to Int
func ByteToInt(data []byte) int {
	return int(data[1]) + (int(data[0]) << 8)
} // byteToInt ...

func ToHex(data []byte) string {
	var result = ""

	for _, c := range data {
		result = fmt.Sprintf("%s%#02x ", result, c)
	}

	return result
} // ToHex ...

func CommunicationModeFromString(s string) (*CommunicationMode, error) {
	var m CommunicationMode
	switch s {
	case "initiative":
		m = ModeInitiative
	case "qa":
		m = ModeQA
	default:
		return nil, fmt.Errorf("cannot covert %s to CommunicationMode", s)
	}
	return &m, nil
} // CommunicationModeFromString ...
