// +build linux

package bbrconn

import (
	"fmt"
	"net"
	"reflect"
	"sync/atomic"

	"github.com/getlantern/netx"
	"github.com/getlantern/tcpinfo"
	"github.com/mikioh/tcp"
)

func Wrap(conn net.Conn) (Conn, error) {
	var tcpConn net.Conn
	netx.WalkWrapped(conn, func(candidate net.Conn) bool {
		switch t := candidate.(type) {
		case *net.TCPConn:
			tcpConn = t
			return false
		}
		return true
	})
	if tcpConn == nil {
		return nil, fmt.Errorf("Could not find a net.TCPConn from connection of type %v", reflect.TypeOf(conn))
	}

	tconn, err := tcp.NewConn(conn)
	if err != nil {
		return nil, err
	}
	return &bbrconn{Conn: conn, tconn: tconn}, nil
}

func (conn *bbrconn) Info() (int, *tcpinfo.BBRInfo, error) {
	var o tcpinfo.CCInfo
	b := make([]byte, 16)
	i, err := conn.tconn.Option(o.Level(), o.Name(), b)
	if err != nil {
		return 0, nil, err
	}
	ai, err := tcpinfo.ParseCCAlgorithmInfo("bbr", i.(*tcpinfo.CCInfo).Raw)
	if err != nil {
		return 0, nil, err
	}
	return int(atomic.LoadUint64(&conn.bytesWritten)), ai.(*tcpinfo.BBRInfo), nil
}
