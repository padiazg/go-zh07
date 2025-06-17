package zh07

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

const (
	checksum = 0x0327
)

var (
	sampleInitiativePayload = []byte{
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
	}

	sampleInitiativeBadChecksum = []byte{
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

)

func TestZH07i_Init(t *testing.T) {
	tests := []struct {
		name    string
		before  func(z *ZH07i)
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
		{
			name: "fail",
			before: func(z *ZH07i) {
				z.write = func(_ *bufio.ReadWriter, _ []byte) error {
					return fmt.Errorf("test error from write")
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				b = &bytes.Buffer{}
				z = NewZH07i(&Config{
					RW: bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b)),
				})
			)

			if tt.before != nil {
				tt.before(z)
			}

			go dummyCommandResponder(t, z.rw, commandSetInitiativeUploadMode, []byte{0x01, 0x02, 0x03})

			if err := z.Init(); (err != nil) != tt.wantErr {
				t.Errorf("ZH07q.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZH07i_CalculateChecksum(t *testing.T) {
	var z *ZH07i = &ZH07i{data: sampleInitiativePayload}
	if cs := z.CalculateChecksum(); cs != z.getChecksum() {
		t.Errorf("TestCalculateChecksumInitiative, got %d, expected %d", cs, checksum)
	}
}

func TestZH07i_IsReadingInvalid(t *testing.T) {
	var (
		tests = []struct {
			name string
			data []byte
			want bool
		}{
			{
				name: "success",
				data: sampleInitiativePayload,
				want: true,
			},
			{
				name: "fail",
				data: sampleInitiativeBadChecksum,
				want: false,
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			z := &ZH07i{data: tt.data}
			if got := z.IsReadingValid(); got != tt.want {
				t.Errorf("ZH07i.IsReadingValid() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestZH07i_Read(t *testing.T) {
	var (
		tests = []struct {
			name   string
			data   []byte
			checks []checkFn
		}{
			{
				name: "success",
				data: sampleInitiativePayload,
				checks: check(
					hasError(false),
					pm(0x6E, 0x7C, 0x54),
				),
			},
			{
				name: "skip-reading-empty-buffer",
				data: []byte{0x00},
				checks: check(
					hasError(false),
					isNil,
				),
			},
			{
				name: "skip-reading-incomplete-header-bytes",
				data: []byte{0x42, 0x00},
				checks: check(
					hasError(false),
					isNil,
				),
			},
			{
				name: "skip-reading-missing-frame-length",
				data: []byte{0x42, 0x4d, 0x00, 0x00},
				checks: check(
					hasError(false),
					isNil,
				),
			},
			{
				name: "fail-unexpected-eof",
				data: []byte{0x42, 0x4d, 0x00, 0x1C, 0x00},
				checks: check(
					hasError(true),
				),
			},
		}
	)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				z = NewZH07i(&Config{
					RW: bufio.NewReadWriter(bufio.NewReader(bytes.NewReader(tt.data)), nil),
				})
				got *Reading
				err error
			)

			got, err = z.Read()
			for _, c := range tt.checks {
				c(t, got, err)
			}
		})
	}
}

func TestZH07i_getChecksum(t *testing.T) {
	var z *ZH07i = &ZH07i{data: sampleInitiativePayload}
	if cs := z.getChecksum(); cs != checksum {
		t.Errorf("TestGetChecksumInitiative, got %d, expected %d", cs, checksum)
	}
}
