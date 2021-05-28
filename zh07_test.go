package zh07

import (
	"bufio"
	"bytes"
	"testing"
	"time"
)

func TestCalculateChecksum(t *testing.T) {
	var data *[]byte = &[]byte{0xFF, 0x86, 0x00, 0x47, 0x00, 0xC7, 0x03, 0x0F, 0x5A}

	if c := CalculateChecksum(data); c != 0x5A {
		t.Errorf("TestCalculateChecksum. Got: %d, expexted: %d", c, 0x5A)
	}
} // TestCalculateChecksum ...

func TestToHex(t *testing.T) {
	if ToHex([]byte{0x00, 0x01, 0x02, 0x0a, 0xff}) != "0x00 0x01 0x02 0x0a 0xff " {
		t.Errorf("Mistmatch: [%s]", ToHex([]byte{0x00, 0x01, 0x02, 0x0a, 0xff}))
	}
} // TestToHex ...

func TestCommunicationModeFromString(t *testing.T) {
	if c, e := CommunicationModeFromString("qa"); e != nil {
		t.Errorf("TestCommunicationModeFromString | Got %s, expected %s", c, ModeQA)
	}

	if c, e := CommunicationModeFromString("initiative"); e != nil {
		t.Errorf("TestCommunicationModeFromString | Got %s, expected %s", c, ModeInitiative)
	}

	if _, e := CommunicationModeFromString("wrong"); e == nil {
		t.Errorf("TestCommunicationModeFromString | Got nil, expected error")
	}
} // TestCommunicationModeFromString ...

func dummyCommandResponder(r0 *bufio.ReadWriter, c *[]byte, r *[]byte, t *testing.T) {
	b := make([]byte, 9) // buffer to receive response
	for {
		if _, e := r0.Read(b); e != nil { // read response from tty
			// command received?
			if b[0] == (*c)[0] && b[1] == (*c)[1] && b[8] == (*c)[8] {
				// write response
				_, e := r0.Write(*r)
				if e != nil { // write commant to tty
					t.Error(e)
				}
				r0.Writer.Flush()
				break
			} // if (*r)[0] == (*c)[0] &&  ...
		} // if _, e := r0.Read(r); e != nil ....
		time.Sleep(250 * time.Millisecond)
	} // for ...
} // DummyCommandResponder ...

func TestSendCommand(t *testing.T) {
	var (
		command  []byte = []byte{0xFF, 0x86, 0x00, 0x47, 0x00, 0xC7, 0x03, 0x0F, 0x5A}
		response []byte = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	)

	// buffer to simulate tty
	b := new(bytes.Buffer)
	rw := bufio.NewReadWriter(
		bufio.NewReader(b),
		bufio.NewWriter(b),
	)

	go dummyCommandResponder(rw, &command, &response, t)

	res, e0 := SendCommand(rw, &command)
	if e0 != nil {
		t.Errorf("TestSendCommand | Sending command: %v", e0)
	}

	for i, v := range res {
		if response[i] != v {
			t.Errorf("TestSendCommand | At response index %d. Got %X, expected %X", i, v, response[i])
		}
	} // for i, v := range res ...
} // TestSendCommand ...

func TestModeInitiativeReadSuccessful(t *testing.T) {
	b := new(bytes.Buffer)
	rw := bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b))

	if _, e := NewZH07(ModeInitiative, rw); e != nil {
		t.Error(e)
	}

	// 	w.Flush()

	buf := make([]byte, 9)
	n, err := rw.Reader.Read(buf)
	if err != nil {
		t.Error(err)
	}
	if n != 9 {
		t.Errorf("TestSetModeInitiative. Got %d, expected 9", n)
	}
} // TestModeInitiativeReadSuccessful ...

func TestModeQaReadSuccessful(t *testing.T) {
	var (
		command  []byte = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
		response []byte = []byte{0xff, 0x86, 0x00, 0x85, 0x00, 0x96, 0x00, 0x65, 0xfa}
	)

	// buffer to simulate tty
	b := new(bytes.Buffer)
	rw := bufio.NewReadWriter(
		bufio.NewReader(b),
		bufio.NewWriter(b),
	)

	z, e := NewZH07(ModeQA, rw)
	if e != nil {
		t.Errorf("TestSetModeQA | Creating instance: %v", e)
	}

	// go function to capture command and send dummy response
	go dummyCommandResponder(rw, &command, &response, t) // go func ...

	if r0, e := z.Read(); e != nil {
		t.Errorf("TestSetModeQA | Reading: %v", e)
	} else {
		var v int

		v = ByteToInt(response[2:4])
		if r0.MassPM25 != v {
			t.Errorf("TestSetModeQA | PM 2.5, Got %d, expected %d", r0.MassPM25, v)
		}

		v = ByteToInt(response[4:6])
		if r0.MassPM10 != v {
			t.Errorf("TestSetModeQA | PM 10, Got %d, expected %d", r0.MassPM10, v)
		}

		v = ByteToInt(response[6:8])
		if r0.MassPM1 != v {
			t.Errorf("TestSetModeQA | PM 1.0, Got %d, expected %d", r0.MassPM1, v)
		}

	}
} // TestModeQaReadSuccessful ...

func TestModeQaReadChecksumError(t *testing.T) {
	var (
		command  []byte = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
		response []byte = []byte{0xff, 0x86, 0x00, 0x85, 0x00, 0x96, 0x00, 0x65, 0xfb}
	)

	// buffer to simulate tty
	b := new(bytes.Buffer)
	rw := bufio.NewReadWriter(
		bufio.NewReader(b),
		bufio.NewWriter(b),
	)

	z, e := NewZH07(ModeQA, rw)
	if e != nil {
		t.Errorf("TestSetModeQA | Creating instance: %v", e)
	}

	// go function to capture command and send dummy response
	go dummyCommandResponder(rw, &command, &response, t) // go func ...

	if _, e := z.Read(); e == nil {
		t.Errorf("TestSetModeQA | Got nil, expected error")
	}
} // TestModeQaReadChecksumError ...
