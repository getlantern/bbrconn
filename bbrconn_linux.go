// +build linux

package bbrconn

import (
	"net"

	"github.com/getlantern/tcpinfo"
	"github.com/mikioh/tcp"
)

func Wrap(conn net.Conn) (Conn, error) {
	tconn, err := tcp.NewConn(conn)
	if err != nil {
		return nil, err
	}
	return &bbrconn{conn, tconn}, nil
}

func (conn *bbrconn) Info() (*tcpinfo.BBRInfo, error) {
	var o tcpinfo.CCInfo
	b := make([]byte, 16)
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
