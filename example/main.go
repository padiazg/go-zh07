package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/namsral/flag"
	"github.com/padiazg/go-zh07"

	"github.com/tarm/serial"
)

const (
	defaultTick = 120 * time.Second
	defaultMode = "qa"
	defaultPort = "/dev/serial0"
	defaultWait = 1500 * time.Millisecond
)

type Config struct {
	Tick time.Duration
	Mode zh07.CommunicationMode
	Port string
	Wait time.Duration
} // Config ...

// Init reads configuration
func (c *Config) Init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "", "Path to config file")

	var (
		tick = flags.Duration("tick", defaultTick, "Ticking interval")
		mode = flags.String("mode", defaultMode, "Reading mode")
		port = flags.String("port", defaultPort, "TTY Port")
		wait = flags.Duration("wait", defaultWait, "Send command wait to request response")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.Tick = *tick
	c.Port = *port
	c.Wait = *wait

	m, e := zh07.CommunicationModeFromString(*mode)
	if e != nil {
		return e
	}
	c.Mode = *m

	return nil
} // Config.Init

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	c := &Config{}

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}() // defer func...

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Printf("Got SIGINT/SIGTERM, exiting.")
					cancel()
					os.Exit(1)
				case syscall.SIGHUP:
					log.Printf("Got SIGHUP, reloading configuration.")
					c.Init(os.Args)
				} // switch ...
			case <-ctx.Done():
				log.Printf("Done")
				os.Exit(1)
			} // select ...
		} // for ...
	}() // go func ...

	if err := run(ctx, c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
} // main ...

// Run monitoring main loop
func run(ctx context.Context, c *Config, out io.Writer) error {
	c.Init(os.Args)
	log.SetOutput(os.Stdout)

	var (
		e      error
		z      *zh07.ZH07
		r      *zh07.Reading
		ticker *time.Ticker = time.NewTicker(c.Tick)
	)

	// open TTY port
	s, err := serial.OpenPort(&serial.Config{
		Name:     c.Port,
		Baud:     9600,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating sensor instance")

	// we wrap the tty port with a bufio.ReadWriter
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// create a sensor instance
	z, e = zh07.NewZH07(c.Mode, rw)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e)
		log.Fatal(e)
	}

	log.Printf("Starting ticker. Triggering every %s", c.Tick)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if r, e = z.Read(); e != nil {
				fmt.Printf("Reading from tty: %v\n", e)
				continue
			}
			fmt.Printf("Reading:\nPM 1.0: %d\nPM 2.5: %d\nPM 10 : %d\n\n", r.MassPM1, r.MassPM25, r.MassPM10)
		} // select
	} // for ...
} // run ...
