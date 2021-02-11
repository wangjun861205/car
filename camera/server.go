package camera

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"
)

type displayServer struct {
	server  server
	doneIn  chan struct{}
	doneOut chan struct{}
}

func NewDisplayServer(server server) *displayServer {
	return &displayServer{
		server:  server,
		doneIn:  make(chan struct{}),
		doneOut: make(chan struct{}),
	}
}

func (s *displayServer) Run() {
	defer s.server.Stop()
	defer close(s.doneOut)
	window := gocv.NewWindow("from network")
	defer window.Close()
	s.server.RegisterHandler(func(b []byte) []byte {
		var img Image
		if err := json.Unmarshal(b, &img); err != nil {
			log.Println(errors.Wrap(err, "failed to unmarshal image"))
			return nil
		}
		mat, err := gocv.NewMatFromBytes(img.Rows, img.Cols, img.Type, img.Data)
		if err != nil {
			log.Println(errors.Wrap(err, "failed to create mat"))
			return nil
		}
		defer mat.Close()
		window.IMShow(mat)
		return nil
	})
	s.server.Run()
}

func (s *displayServer) Close() {
	close(s.doneIn)
}
