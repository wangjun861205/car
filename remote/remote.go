package remote

import (
	"buxiong/car/model"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

type remote struct {
	client   client
	keyboard keyboardReader
}

func NewRemote(client client, keyboard keyboardReader) *remote {
	return &remote{
		client,
		keyboard,
	}
}

func (r *remote) registerHandler() {
	r.client.RegisterHandler(func(b []byte) {
		var resp model.Response
		if err := json.Unmarshal(b, &resp); err != nil {
			log.Println(errors.Wrap(err, "failed to unmarshal response"))
			return
		}
		fmt.Printf("Left Target(Direction): %f(%s), Right Target(Direction): %f(%s)\n", resp.LeftTarget, resp.LeftDirection, resp.RightTarget, resp.RightDirection)
	})
}

func (r *remote) Run() {
	r.registerHandler()
	go r.client.Run()
	go r.keyboard.Run()
	for {
		select {
		case event := <-r.keyboard.Out():
			if event.KeyPress() {
				switch event.KeyString() {
				case "Up":
					b, _ := json.Marshal(model.Request{Action: model.Forward})
					r.client.Write(b)
				case "Down":
					b, _ := json.Marshal(model.Request{Action: model.Backward})
					r.client.Write(b)
				case "Left":
					b, _ := json.Marshal(model.Request{Action: model.TurnLeft})
					r.client.Write(b)
				case "Right":
					b, _ := json.Marshal(model.Request{Action: model.TurnRight})
					r.client.Write(b)
				case "S":
					b, _ := json.Marshal(model.Request{Action: model.Stop})
					r.client.Write(b)
				default:
					log.Println(errors.Errorf("unknown key map(key: %s)", event.KeyString()))
				}
			}
		}
	}
}
