// Package bbrconn provides a wrapper around net.Conn that exposes BBR
// congestion control information. This works only on Linux with Kernel 4.9.0 or
// newer installed.
package bbrconn

import (
	"net"
	"sync/atomic"

	"github.com/getlantern/tcpinfo"
	"github.com/mikioh/tcp"
)

type InfoCallback func(bytesWritten int, info *tcpinfo.BBRInfo, err error)

type Conn interface {
	net.Conn
	Info() (bytesWritten int, info *tcpinfo.BBRInfo, err error)
}

type bbrconn struct {
	net.Conn
	tconn        *tcp.Conn
	bytesWritten uint64
	onClose      InfoCallback
}

func (c *bbrconn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	if n > 0 {
		atomic.AddUint64(&c.bytesWritten, uint64(n))
	}
	return n, err
}

// Wrapped implements the interface netx.Wrapped
func (c *bbrconn) Wrapped() net.Conn {
	return c.Conn
}
