package main

import (
	"fmt"
	"time"

	"github.com/namsral/flag"
)

const (
	defaultInterval = 120 * time.Second
	defaultMode     = "qa"
	defaultPort     = "/dev/serial0"
)

type Config struct {
	Mode     string
	Port     string
	Interval time.Duration
} // Config ...

// Init reads configuration
func (c *Config) Init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "", "Path to config file")

	var (
		tick = flags.Duration("interval", defaultInterval, "Reading interval, defaults to 2m")
		mode = flags.String("mode", defaultMode, "Reading mode <qa|initiative>, defaults to `qa`")
		port = flags.String("port", defaultPort, "TTY Port, defaults to `/dev/serial0`")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.Interval = *tick
	c.Port = *port
	c.Mode = *mode

	return nil
}

func (c *Config) Show() {
	fmt.Printf("Using:\n")
	fmt.Printf("Mode: %s\n", c.Mode)
	fmt.Printf("Port: %s\n", c.Port)
	fmt.Printf("Interval: %s\n", c.Interval)
}
