package nettool

import (
	"io"
	"net"
)

type Conn struct {
	net.Conn
	r io.ReadCloser
	w io.WriteCloser
}

func NewConn(conn net.Conn, r io.ReadCloser, w io.WriteCloser) *net.Conn {
	return &Conn{
		Conn: conn,
		r:    r,
		w:    w,
	}
}

func (c *Conn) Read(p []byte) (n int, err error) {
	if c.r != nil {
		return c.r.Read(p)
	}
	return c.Conn.Read(p)
}

func (c *Conn) Write(p []byte) (n int, err error) {
	if c.w != nil {
		return c.w.Write(p)
	}
	return c.Conn.Write(p)
}

func (c *Conn) Close() error {
	if c.r != nil {
		c.r.Close()
	}

	if c.w != nil {
		c.w.Close()
	}

	return c.Conn.Close()
}


