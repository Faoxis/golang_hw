package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("connect to nonexistent server", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("1s")
		require.NoError(t, err)

		client := NewTelnetClient("127.0.0.1:9999", timeout, io.NopCloser(in), out)
		err = client.Connect()
		require.Error(t, err, "expected connection error")
	})

	t.Run("send without connect", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:9999", 1*time.Second, io.NopCloser(in), out)

		in.WriteString("test message\n")
		err := client.Send()
		require.Error(t, err, "expected error when sending without connection")
	})

	t.Run("receive without connect", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:9999", 1*time.Second, io.NopCloser(in), out)

		err := client.Receive()
		require.Error(t, err, "expected error when receiving without connection")
	})

	t.Run("close without connect", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:9999", 1*time.Second, io.NopCloser(in), out)

		err := client.Close()
		require.Error(t, err)
	})
}
