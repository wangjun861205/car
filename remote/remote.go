package remote

import (
	"buxiong/car/model"
	"encoding/json"
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
		var resp model.CarResponse
		if err := json.Unmarshal(b, &resp); err != nil {
			log.Println(errors.Wrap(err, "failed to unmarshal response"))
			return
		}
		log.Println(resp)
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
					b, _ := json.Marshal(model.CarInstruction{Action: model.ActionAccelerate})
					r.client.Write(b)
				case "Down":
					b, _ := json.Marshal(model.CarInstruction{Action: model.ActionBrake})
					r.client.Write(b)
				case "Left":
					b, _ := json.Marshal(model.CarInstruction{Action: model.ActionTurnLeft})
					r.client.Write(b)
				case "Right":
					b, _ := json.Marshal(model.CarInstruction{Action: model.ActionTurnRight})
					r.client.Write(b)
				case "Q":
					b, _ := json.Marshal(model.CarInstruction{Action: model.ActionStop})
					r.client.Write(b)

				default:
					log.Println(errors.Errorf("unknown key map(key: %s)", event.KeyString()))
				}
			}
		}
	}
}