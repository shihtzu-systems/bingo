package main

import (
	"fmt"
	"github.com/shihtzu-systems/bingo/cmd"
	"os"
	"os/signal"
)

func main() {
	// setup signal catching
	signals := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	signal.Notify(signals)
	// method invoked upon seeing signal
	go func() {
		s := <-signals
		fmt.Printf("%s received", s.String())
		switch s.String() {
		case "window size changes":
			fmt.Printf("ok")

		case "trace/breakpoint trap":
			fallthrough
		case "interrupt":
			fallthrough
		default:
			fmt.Printf("stopping!")
		}
	}()
	// run
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
