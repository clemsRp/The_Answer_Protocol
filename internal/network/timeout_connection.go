package network

import (
	"net"
	"time"
)

type TimeoutConn struct {
	net.Conn
	timeout time.Duration
}

// override Read method of net.Conn interface
// wrapper for timeout read
func (c *TimeoutConn) Read(b []byte) (int, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (s *Server) createTimeoutConn(conn net.Conn) *TimeoutConn {
	return &TimeoutConn{
		Conn:    conn,
		timeout: time.Duration(s.timeoutSeconds) * time.Second,
	}
}
