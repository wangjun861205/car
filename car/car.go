package car

import (
	"bufio"
	"buxiong/car/model"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"
)

// Car Car
type Car struct {
	addr string
	ctl  controller
}

// NewCar NewCar
func NewCar(ctl controller, addr string) *Car {
	return &Car{
		addr,
		ctl,
	}
}

// Run Run
func (c *Car) Run() {
	go c.ctl.Run()
	listener, err := net.Listen("tcp", c.addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer listener.Close()
	l := listener.(*net.TCPListener)
	var wg sync.WaitGroup
	closeSignal := make(chan os.Signal)
	signal.Notify(closeSignal, os.Interrupt, os.Kill)
	done := make(chan interface{})
	go func() {
		s := <-closeSignal
		log.Printf("got close signal: %v", s)
		close(done)
	}()
	connChan := make(chan *net.TCPConn)
	errChan := make(chan error)
	go func() {
		for {
			conn, err := l.AcceptTCP()
			if err != nil {
				errChan <- err
				return
			}
			connChan <- conn
		}
	}()
OUTER:
	for {
		select {
		case <-done:
			break OUTER
		case err := <-errChan:
			log.Println(err)
			close(done)
			break OUTER
		case conn := <-connChan:
			if err := conn.SetKeepAlive(true); err != nil {
				fmt.Println(errors.Wrap(err, "failed to set keep alive"))
				close(done)
				break OUTER
			}
			wg.Add(1)
			go func(conn *net.TCPConn) {
				reader := bufio.NewReader(conn)
				readChan := make(chan []byte)
				errChan := make(chan error)
				go func() {
					for {
						b, err := reader.ReadBytes('\n')
						if err != nil {
							errChan <- err
							return
						}
						readChan <- b
					}
				}()
				for {
					select {
					case <-done:
						conn.Close()
						wg.Done()
						return
					case err := <-errChan:
						log.Println(errors.Wrap(err, "server failed to read tcp connection"))
						conn.Close()
						wg.Done()
						return
					case b := <-readChan:
						var req model.Request
						if err := json.Unmarshal(b, &req); err != nil {
							resp := model.Response{
								Error: errors.Errorf("invalid request"),
							}
							b, _ := json.Marshal(resp)
							if _, err := conn.Write(append(b, '\n')); err != nil {
								log.Println(errors.Wrap(err, "failed to write to connection"))
								conn.Close()
								wg.Done()
								return
							}
							continue
						}
						switch req.Action {
						case model.Forward:
							c.ctl.Forward()
						case model.Backward:
							c.ctl.Backward()
						case model.TurnLeft:
							c.ctl.TurnLeft()
						case model.TurnRight:
							c.ctl.TurnRight()
						case model.Stop:
							c.ctl.Stop()
						}
						stat := c.ctl.Status()
						b, _ = json.Marshal(stat)
						if _, err := conn.Write(append(b, '\n')); err != nil {
							fmt.Println(errors.Wrap(err, "server failed to write tcp connection"))
							conn.Close()
							wg.Done()
							return
						}
					}
				}
			}(conn)
		}
	}
	wg.Wait()
	c.ctl.Close()
	log.Println("exited")
}
