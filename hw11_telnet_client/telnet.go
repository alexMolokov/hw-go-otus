package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &tClient{
		address: address,
		timeOut: timeout,
		in:      in,
		out:     out,
	}
}

type tClient struct {
	address      string
	timeOut      time.Duration
	in           io.ReadCloser
	out          io.Writer
	conn         net.Conn
	isConnClosed bool
}

func (t *tClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeOut)
	if err != nil {
		return fmt.Errorf("can't connect to server address %s", t.address)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.address)
	t.conn = conn
	return nil
}

func (t *tClient) Close() error {
	if t.isConnClosed {
		return nil
	}
	t.isConnClosed = true
	return t.conn.Close()
}

func (t *tClient) Send() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err == nil {
			fmt.Fprintf(os.Stderr, "...Read from stdin %s", str)
			_, err := t.conn.Write([]byte(str))
			if err == nil {
				continue
			}

			if t.isConnClosed {
				return nil
			}

			return fmt.Errorf("can't send to server %s", str)
		}

		if errors.Is(err, io.EOF) {
			fmt.Fprint(os.Stderr, "...Send EOF\n")
			return nil
		}

		return err
	}
}

func (t *tClient) Receive() error {
	reader := bufio.NewReader(t.conn)
	for {
		str, err := reader.ReadString('\n')
		if err == nil {
			fmt.Fprintf(os.Stderr, "...Read from conn %s", str)
			_, err := t.out.Write([]byte(str))
			if err == nil {
				continue
			}

			return fmt.Errorf("can't receive %s", str)
		}

		if t.isConnClosed {
			return nil
		}

		if errors.Is(err, io.EOF) {
			fmt.Fprint(os.Stderr, "...Receive EOF\n")
			return nil
		}

		return err
	}
}
