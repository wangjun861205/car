package car

import (
	"buxiong/car/model"
	"encoding/json"

	"github.com/pkg/errors"
)

type car struct {
	ctl    controller
	server server
}

func (c *car) handleInstruction(inst model.CarInstruction) error {
	switch inst.Action {
	case model.ActionAccelerate:
		c.ctl.Accelerate()
	case model.ActionBrake:
		c.ctl.Brake()
	case model.ActionTurnLeft:
		c.ctl.TurnLeft()
	case model.ActionTurnRight:
		c.ctl.TurnRight()
	case model.ActionStop:
		c.ctl.Stop()
		c.server.Stop()
	default:
		return errors.Errorf("unknown instruction: %s", inst)
	}
	return nil
}

func (c *car) registerHandler() {
	c.server.RegisterHandler(func(b []byte) []byte {
		var inst model.CarInstruction
		if err := json.Unmarshal(b, &inst); err != nil {
			err = errors.Wrap(err, "failed to unmarshal instruction")
			stat := c.ctl.Status()
			resp, _ := json.Marshal(model.CarResponse{
				LeftBase:        stat.Left.Base,
				RightBase:       stat.Right.Base,
				LeftStep:        stat.Left.Step,
				RightStep:       stat.Right.Step,
				LeftProportion:  stat.Left.Proportion,
				RightProportion: stat.Right.Proportion,
				Error:           err,
			})
			return resp
		}
		err := c.handleInstruction(inst)
		stat := c.ctl.Status()
		resp, _ := json.Marshal(model.CarResponse{
			LeftBase:        stat.Left.Base,
			RightBase:       stat.Right.Base,
			LeftStep:        stat.Left.Step,
			RightStep:       stat.Right.Step,
			LeftProportion:  stat.Left.Proportion,
			RightProportion: stat.Right.Proportion,
			Error:           err,
		})
		return resp
	})
}

func NewCar(ctl controller, server server) *car {
	car := &car{
		ctl,
		server,
	}
	car.registerHandler()
	return car
}

func (c *car) Run() {
	c.ctl.Run()
	c.server.Run()
}
