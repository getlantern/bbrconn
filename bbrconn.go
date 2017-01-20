// Package bbrconn provides a wrapper around net.Conn that exposes BBR
// congestion control information. This works only on Linux with Kernel 4.9.0 or
// newer installed.
package bbrconn

import (
	"net"

	"github.com/getlantern/tcpinfo"
	"github.com/mikioh/tcp"
)

type Conn interface {
	net.Conn
	Info() (*tcpinfo.BBRInfo, error)
}

type bbrconn struct {
	net.Conn
	tconn *tcp.Conn
}
