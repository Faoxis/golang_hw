package main

import (
	"context"
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
	context    context.Context
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.

	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		context: context.WithoutCancel(context.Background()),
	}
}

func (tc *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.connection = conn
	return nil
}

func (tc *telnetClient) Close() error {
	if tc.connection == nil {
		return errors.New("client not connected")
	}
	err := tc.connection.Close()
	tc.context.Done()
	if err != nil {
		return err
	}
	return nil
}

func (tc *telnetClient) Send() error {
	if tc.connection == nil {
		return errors.New("client not connected")
	}
	return readWriteAsync(&tc.context, tc.connection, tc.in)
}

func (tc *telnetClient) Receive() error {
	if tc.connection == nil {
		return errors.New("client not connected")
	}
	return readWriteAsync(&tc.context, tc.out, tc.connection)
}

func readWriteAsync(ctx *context.Context, writer io.Writer, reader io.Reader) error {
	errorCh := make(chan error, 1)
	defer close(errorCh)
	go func() {
		_, err := io.Copy(writer, reader)
		errorCh <- err
	}()

	select {
	case <-(*ctx).Done():
		return (*ctx).Err()
	case err := <-errorCh:
		return err
	}
}
