package network

import (
	"fmt"
	"testing"
	"time"
)

func TestNetwork(t *testing.T) {
	server, err := NewServer(":9999", '\n')
	if err != nil {
		t.Fatal(err)
	}
	server.Run()
	client, err := NewClient("localhost:9999", '\n')
	client.RegisterHandler(
		func(b []byte) {
			fmt.Printf("client received: %s", b)
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	client.Run()
	client.Write([]byte("test"))
	time.Sleep(time.Second)
	server.Stop()
}
