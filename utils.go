package client

import "net"

// 连接
func dial(addr string) (*net.TCPConn, error) {
	rAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp4", nil, rAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
