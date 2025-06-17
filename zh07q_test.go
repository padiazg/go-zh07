package zh07

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

var (
	sampleQAPayload = []byte{
		0xFF,       // starting
		0x86,       // command
		0x00, 0x85, // pm2.5
		0x00, 0x96, // pm10
		0x00, 0x65, // pm1.0
		0xFA, // checksum
	}
	sampleQABadChecksum = []byte{
		0xFF,       // starting
		0x86,       // command
		0x00, 0x85, // pm2.5
		0x00, 0x96, // pm10
		0x00, 0x65, // pm1.0
		0xFB, // checksum
	}
)

func TestZH07q_Init(t *testing.T) {
	tests := []struct {
		name    string
		before  func(z *ZH07q)
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
		{
			name: "fail",
			before: func(z *ZH07q) {
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
				z = NewZH07q(&Config{
					RW: bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b)),
				})
			)

			if tt.before != nil {
				tt.before(z)
			}

			go dummyCommandResponder(t, z.rw, commandQuery, []byte{0x00})

			if err := z.Init(); (err != nil) != tt.wantErr {
				t.Errorf("ZH07q.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZH07q_CalculateChecksum(t *testing.T) {
	var z *ZH07q = &ZH07q{data: sampleQAPayload}
	if cs := z.CalculateChecksum(); cs != z.getChecksum() {
		t.Errorf("TestCalculateChecksumInitiative, got %d, expected %d", cs, checksum)
	}
}

func TestZH07q_IsReadingInvalid(t *testing.T) {
	var (
		tests = []struct {
			name string
			data []byte
			want bool
		}{
			{
				name: "success",
				data: sampleQAPayload,
				want: true,
			},
			{
				name: "fail",
				data: sampleQABadChecksum,
				want: false,
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			z := &ZH07q{data: tt.data}
			if got := z.IsReadingValid(); got != tt.want {
				t.Errorf("ZH07q.IsReadingValid() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestZH07q_Read(t *testing.T) {
	var (
		tests = []struct {
			name     string
			response []byte
			checks   []checkFn
			before   func(z *ZH07q)
		}{
			{
				name:     "success",
				response: sampleQAPayload,
				checks: check(
					hasError(false),
					pm(0x85, 0x96, 0x65),
				),
			},
			{
				name:     "fail-checksum-mismatch",
				response: sampleQABadChecksum,
				checks: check(
					hasError(true),
				),
			},
			{
				name: "fail-sendcommand",
				before: func(z *ZH07q) {
					z.writeAndRead = func(_ *bufio.ReadWriter, _ []byte) ([]byte, error) {
						return nil, fmt.Errorf("test error from sendCommand")
					}
				},
				checks: check(
					hasError(true),
				),
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				b = &bytes.Buffer{}
				z = NewZH07q(&Config{
					RW: bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b)),
				})
				got *Reading
				err error
			)

			if tt.before != nil {
				tt.before(z)
			}

			go dummyCommandResponder(t, z.rw, commandQuery, tt.response)

			got, err = z.Read()
			for _, c := range tt.checks {
				c(t, got, err)
			}
		})
	}
}

func TestZH07q_getChecksum(t *testing.T) {
	var z *ZH07q = &ZH07q{data: sampleQAPayload}
	if cs := z.getChecksum(); cs != 0xFA {
		t.Errorf("TestGetChecksumQA, got %d, expected %d", cs, checksum)
	}
}
