package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		signalChan  = make(chan os.Signal, 1)
		c           = &Config{}
	)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	fmt.Printf("ZHx sensors reading demo")
	if err := c.Init(os.Args); err != nil {
		log.Fatalf("reading config: %+v", err)
		return
	}

	c.Show()

	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			fmt.Printf("Got SIGINT/SIGTERM, exiting.\n")
			cancel()
			os.Exit(1)
		case syscall.SIGHUP:
			fmt.Printf("Got SIGHUP, reloading configuration.\n")
			c.Init(os.Args)
		}
	case <-ctx.Done():
		fmt.Printf("Done\n")
		os.Exit(1)
	}
}

func run(ctx context.Context, c *Config) {

}
