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

type InfoCallback func(bytesWritten int, info *tcpinfo.Info, bbrInfo *tcpinfo.BBRInfo, err error)

type Conn interface {
	net.Conn

	// BytesWritten returns the number of bytes written to this connection
	BytesWritten() int

	// TCPInfo returns TCP connection info from the kernel
	TCPInfo() (*tcpinfo.Info, error)

	// BBRInfo returns BBR congestion avoidance info from the kernel
	BBRInfo() (*tcpinfo.BBRInfo, error)
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
