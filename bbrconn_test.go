package bbrconn

import (
	"io"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	iters = 100000
)

func TestConcurrency(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if !assert.NoError(t, err) {
		return
	}
	defer l.Close()

	go func() {
		conn, acceptErr := l.Accept()
		if !assert.NoError(t, acceptErr) {
			return
		}
		go io.Copy(ioutil.Discard, conn)
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	if !assert.NoError(t, err) {
		return
	}
	bconn, err := Wrap(conn, nil)
	if !assert.NoError(t, err) {
		conn.Close()
		return
	}
	defer bconn.Close()

	result := make(chan int)
	go func() {
		var sent int
		for i := 0; i < iters; i++ {
			n := bconn.BytesWritten()
			sent += n

			if i == iters-1 {
				// Last iteration
				info, err := bconn.TCPInfo()
				if err != nil {
					t.Fatal(err)
				}
				_, err = bconn.BBRInfo()
				if err != nil {
					t.Fatal(err)
				}
				assert.True(t, info.SenderMSS > 0)
				assert.True(t, info.Sys.SegsOut > 0)
				assert.True(t, info.Sys.SegsOut > info.Sys.TotalRetransSegs)
			}
		}
		result <- sent
	}()

	for {
		_, err := bconn.Write([]byte("hello there"))
		if !assert.NoError(t, err) {
			return
		}
		select {
		case sent := <-result:
			assert.True(t, sent > 0)
			return
		default:
			// keep writing
		}
	}
}

func BenchmarkTCPInfo(b *testing.B) {
	conn, err := net.Dial("tcp", "www.google.com:443")
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	bc, err := Wrap(conn, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bc.TCPInfo()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBBRInfo(b *testing.B) {
	conn, err := net.Dial("tcp", "www.google.com:443")
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	bc, err := Wrap(conn, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bc.BBRInfo()
		if err != nil {
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
