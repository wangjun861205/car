package network

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/pkg/errors"
)

type client struct {
	conn      *net.TCPConn
	delimiter byte
	handler   func(b []byte)
}

func NewClient(addr string, delimiter byte) (*client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dail to server")
	}
	c := conn.(*net.TCPConn)
	if err := c.SetKeepAlive(true); err != nil {
		return nil, errors.Wrap(err, "failed to set keep alive")
	}
	return &client{
		conn:      conn.(*net.TCPConn),
		delimiter: delimiter,
		handler:   nil,
	}, nil
}

func (c *client) RegisterHandler(handler func(b []byte)) {
	c.handler = handler
}

func (c *client) Run() {
	reader := bufio.NewReader(c.conn)
	for {
		b, err := reader.ReadBytes(c.delimiter)
		if err != nil {
			if err != io.EOF {
				fmt.Println(errors.Wrap(err, "failed to read connection"))
				continue
			} else {
				c.conn.Close()
				return
			}
		}
		c.handler(b)
	}
}

func (c *client) Write(b []byte) error {
	b = append(b, c.delimiter)
	if _, err := c.conn.Write(b); err != nil {
		return errors.Wrap(err, "failed to write to connection")
	}
	return nil
}

func (c *client) Stop() {
	c.conn.Close()
}
