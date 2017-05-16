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

func Wrap(conn net.Conn, onClose InfoCallback) (Conn, error) {
	if onClose == nil {
		onClose = func(bytesWritten int, info *tcpinfo.Info, bbrInfo *tcpinfo.BBRInfo, err error) {
		}
	}
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

	tconn, err := tcp.NewConn(tcpConn)
	if err != nil {
		return nil, fmt.Errorf("Unable to wrap TCP conn: %v", err)
	}
	return &bbrconn{Conn: conn, tconn: tconn, onClose: onClose}, nil
}

func (conn *bbrconn) BytesWritten() int {
	return int(atomic.LoadUint64(&conn.bytesWritten))
}

func (conn *bbrconn) TCPInfo() (*tcpinfo.Info, error) {
	var o tcpinfo.Info
	b := make([]byte, o.Size())
	i, err := conn.tconn.Option(o.Level(), o.Name(), b)
	if err != nil {
		return nil, err
	}
	return i.(*tcpinfo.Info), nil
}

func (conn *bbrconn) BBRInfo() (*tcpinfo.BBRInfo, error) {
	var bo tcpinfo.BBRInfo
	var o tcpinfo.CCInfo
	b := make([]byte, bo.Size())
	i, err := conn.tconn.Option(o.Level(), o.Name(), b)
	if err != nil {
		return nil, err
	}
	ai, err := tcpinfo.ParseCCAlgorithmInfo("bbr", i.(*tcpinfo.CCInfo).Raw)
	if err != nil {
		return nil, err
	}
	return ai.(*tcpinfo.BBRInfo), nil
}

func (conn *bbrconn) Close() error {
	bytesWritten := conn.BytesWritten()
	info, err1 := conn.TCPInfo()
	bbrInfo, err2 := conn.BBRInfo()
	err := err1
	if err == nil {
		err = err2
	}
	conn.onClose(bytesWritten, info, bbrInfo, err)
	return conn.Conn.Close()
}
