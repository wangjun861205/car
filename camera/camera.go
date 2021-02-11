package camera

import (
	"encoding/json"
	"log"

	_ "github.com/hybridgroup/mjpeg"
	"github.com/pkg/errors"
	"gocv.io/x/gocv"
)

type camera struct {
	capture *gocv.VideoCapture
	client  client
	doneIn  chan struct{}
	doneOut chan struct{}
}

func NewCamera(devid int, client client) (*camera, error) {
	webcam, err := gocv.OpenVideoCapture(devid)
	if err != nil {
		return nil, errors.Wrap(err, "create camera error")
	}
	return &camera{
		webcam,
		client,
		make(chan struct{}),
		make(chan struct{}),
	}, nil
}

func (c *camera) Run() {
	defer c.capture.Close()
	defer c.client.Stop()
	defer close(c.doneOut)
	c.client.RegisterHandler(func(b []byte) {
		log.Println(string(b))
	})
	c.client.Run()
	img := gocv.NewMat()
	for {
		select {
		case <-c.doneIn:
			log.Println("close camera")
			return
		default:
			if ok := c.capture.Read(&img); ok {
				log.Println("device closed")
				return
			}
			if img.Empty() {
				continue
			}
			rows := img.Rows()
			cols := img.Cols()
			typ := img.Type()
			b, _ := json.Marshal(Image{
				Rows: rows,
				Cols: cols,
				Type: typ,
				Data: img.ToBytes(),
			})
			if err := c.client.Write(b); err != nil {
				log.Println(errors.Wrap(err, "failed to write"))
				return
			}
		}
	}
}

func (c *camera) Stop() {
	close(c.doneIn)
}
