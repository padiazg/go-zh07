package main

import (
	"context"

	"github.com/padiazg/go-zh07"
)

type Monitor struct {
	config *Config
}

func NewMonitor(config *Config) *Monitor {
	return &Monitor{config: config}
}

func (m *Monitor) Init() (chan zh07.Reading, func()) {
	readingChan := make(chan zh07.Reading, 1)
	return readingChan, func() {
		close(readingChan)
	}
}

func (m *Monitor) Run(ctx context.Context) {

}
