package bbrconn

import (
	"io"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
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
		for i := 0; i < 100000; i++ {
			n, _, err := bconn.Info()
			if err != nil {
				t.Fatal(err)
			}
			sent += n
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

func BenchmarkInfo(b *testing.B) {
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
		_, _, err := bc.Info()
		if err != nil {
			b.Fatal(err)
		}
	}
}
