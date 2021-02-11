package network

import (
	"bufio"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type server struct {
	listener  *net.TCPListener
	delimiter byte
	handler   func([]byte) []byte
	conns     []*net.TCPConn
	done      chan struct{}
}

func NewServer(addr string, delimiter byte) (*server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create tcp listener(address: %s)", addr)
	}
	l := listener.(*net.TCPListener)
	return &server{
		l,
		delimiter,
		nil,
		make([]*net.TCPConn, 0, 128),
		make(chan struct{}),
	}, nil
}

func (s *server) RegisterHandler(handler func([]byte) []byte) {
	s.handler = handler
}

func (s *server) Run() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			fmt.Println(errors.Wrap(err, "failed to accept tcp connection"))
			return

		}
		if err := conn.SetKeepAlive(true); err != nil {
			fmt.Println(errors.Wrap(err, "failed to set keep alive"))
			continue
		}
		s.conns = append(s.conns, conn)
		go func(conn *net.TCPConn) {
			reader := bufio.NewReader(conn)
			for {
				b, err := reader.ReadBytes(s.delimiter)
				if err != nil {
					fmt.Println(errors.Wrap(err, "server failed to read tcp connection"))
					return
				}
				resp := s.handler(b)
				resp = append(resp, s.delimiter)
				if _, err := conn.Write(resp); err != nil {
					fmt.Println(errors.Wrap(err, "server failed to write tcp connection"))
					return
				}
			}
		}(conn)
	}
}

func (s *server) Stop() {
	s.listener.Close()
	for _, conn := range s.conns {
		conn.Close()
	}
}
