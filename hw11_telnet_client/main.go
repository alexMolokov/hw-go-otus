package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout connection in seconds")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("Please enter host and port. Example --timeout=5s 127.0.0.1 80")
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	tcpClient := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	err := tcpClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tcpClient.Close(); err != nil {
			log.Println("Can't close connection to tcp server")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go shutdown(cancel)
	go write(tcpClient, cancel)
	go read(tcpClient, cancel)

	<-ctx.Done()
}

func shutdown(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	<-signals
	fmt.Fprint(os.Stderr, "Graceful shutdown")
	cancel()
}

func write(tcpClient TelnetClient, cancel context.CancelFunc) {
	if err := tcpClient.Send(); err != nil {
		log.Println(err)
	}
	cancel()
}

func read(tcpClient TelnetClient, cancel context.CancelFunc) {
	if err := tcpClient.Receive(); err != nil {
		log.Println(err)
	}
	cancel()
}
