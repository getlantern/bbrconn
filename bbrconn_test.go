package bbrconn

import (
	"net"
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	conn, err := net.Dial("tcp", "www.google.com:443")
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	bc, err := Wrap(conn)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bc.Info()
		if err != nil {
			b.Fatal(err)
		}
	}
}
