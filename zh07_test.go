package zh07

// TestCalculateChecksum ...

//  // TestToHex ...

// func TestcommandSetInitiativeUploadModeReadSuccessful(t *testing.T) {
// 	b := new(bytes.Buffer)
// 	rw := bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b))

// 	if _, e := NewZH07(commandSetInitiativeUploadMode, rw); e != nil {
// 		t.Error(e)
// 	}

// 	// 	w.Flush()

// 	buf := make([]byte, 9)
// 	n, err := rw.Reader.Read(buf)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if n != 9 {
// 		t.Errorf("TestSetcommandSetInitiativeUploadMode. Got %d, expected 9", n)
// 	}
// } // TestcommandSetInitiativeUploadModeReadSuccessful ...

// func TestcommandSetQAModeReadSuccessful(t *testing.T) {
// 	var (
// 		command  []byte = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
// 		response []byte = []byte{0xff, 0x86, 0x00, 0x85, 0x00, 0x96, 0x00, 0x65, 0xfa}
// 	)

// 	// buffer to simulate tty
// 	b := new(bytes.Buffer)
// 	rw := bufio.NewReadWriter(
// 		bufio.NewReader(b),
// 		bufio.NewWriter(b),
// 	)

// 	z, e := NewZH07(commandSetQAMode, rw)
// 	if e != nil {
// 		t.Errorf("TestSetcommandSetQAMode | Creating instance: %v", e)
// 	}

// 	// go function to capture command and send dummy response
// 	go dummyCommandResponder(rw, &command, &response, t) // go func ...

// 	if r0, e := z.Read(); e != nil {
// 		t.Errorf("TestSetcommandSetQAMode | Reading: %v", e)
// 	} else {
// 		var v int

// 		v = byteToInt(response[2:4])
// 		if r0.PM25 != v {
// 			t.Errorf("TestSetcommandSetQAMode | PM 2.5, Got %d, expected %d", r0.PM25, v)
// 		}

// 		v = byteToInt(response[4:6])
// 		if r0.PM10 != v {
// 			t.Errorf("TestSetcommandSetQAMode | PM 10, Got %d, expected %d", r0.PM10, v)
// 		}

// 		v = byteToInt(response[6:8])
// 		if r0.PM1 != v {
// 			t.Errorf("TestSetcommandSetQAMode | PM 1.0, Got %d, expected %d", r0.PM1, v)
// 		}

// 	}
// } // TestcommandSetQAModeReadSuccessful ...

// func TestcommandSetQAModeReadChecksumError(t *testing.T) {
// 	var (
// 		command  []byte = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
// 		response []byte = []byte{0xff, 0x86, 0x00, 0x85, 0x00, 0x96, 0x00, 0x65, 0xfb}
// 	)

// 	// buffer to simulate tty
// 	b := new(bytes.Buffer)
// 	rw := bufio.NewReadWriter(
// 		bufio.NewReader(b),
// 		bufio.NewWriter(b),
// 	)

// 	z, e := NewZH07(commandSetQAMode, rw)
// 	if e != nil {
// 		t.Errorf("TestSetcommandSetQAMode | Creating instance: %v", e)
// 	}

// 	// go function to capture command and send dummy response
// 	go dummyCommandResponder(rw, &command, &response, t) // go func ...

// 	if _, e := z.Read(); e == nil {
// 		t.Errorf("TestSetcommandSetQAMode | Got nil, expected error")
// 	}
// } // TestcommandSetQAModeReadChecksumError ...
