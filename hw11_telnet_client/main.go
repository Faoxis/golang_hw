package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := fmt.Sprintf("%s:%s", host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		printAndExit(err)
	}

	go func() {
		err := client.Send()
		if err != nil {
			printAndExit(err)
			return
		}
	}()
	go func() {
		err := client.Receive()
		if err != nil {
			printAndExit(err)
			return
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	err = client.Close()
	if err != nil {
		printAndExit(err)
	}
}

func printAndExit(err error) {
	fmt.Println("ERROR: ", err)
	os.Exit(1)
}
