package electric

import (
	"strings"

	"github.com/pkg/errors"
)

type commandType string

const (
	turnOnMainSwitch    commandType = "TURN_ON_MAIN_SWITCH"
	turnOffMainSwitch   commandType = "TURN_OFF_MAIN_SWITCH"
	troggleHeadLight    commandType = "TROGGLE_HEAD_LIGHT"
	getMainSwitchStatus commandType = "GET_MAIN_SWITCH_STATUS"
	getHeadLightStatus  commandType = "GET_HEAD_LIGHT_STATUS"
	shutdown            commandType = "SHUTDOWN"
)

type command struct {
	typ commandType
	val int
	err chan error
}

// Controller Controller
type Controller struct {
	mainSwitch IOPinner
	headLight  IOPinner
	commands   chan *command
	done       chan interface{}
}

// NewController NewController
func NewController(mainSwitch, headLight IOPinner) *Controller {
	return &Controller{
		mainSwitch: mainSwitch,
		headLight:  headLight,
		commands:   make(chan *command),
		done:       make(chan interface{}),
	}
}

func (c *Controller) turnOnMainSwitch() error {
	return c.mainSwitch.SetValue(1)
}

func (c *Controller) turnOffMainSwitch() error {
	return c.mainSwitch.SetValue(0)
}

func (c *Controller) troggleHeadLight() error {
	stat, err := c.headLight.Value()
	if err != nil {
		return err
	}
	return c.headLight.SetValue(stat ^ 1)
}

func (c *Controller) getMainSwitchStatus() (int, error) {
	return c.mainSwitch.Value()
}

func (c *Controller) getHeadLightStatus() (int, error) {
	return c.headLight.Value()
}

// Run Run
func (c *Controller) Run() {
	for {
		cmd := <-c.commands
		switch cmd.typ {
		case turnOnMainSwitch:
			err := c.turnOnMainSwitch()
			if err != nil {
				err = errors.Wrap(err, "failed to turn on main switch")
			}
			cmd.err <- err
			close(cmd.err)
		case turnOffMainSwitch:
			err := c.turnOffMainSwitch()
			if err != nil {
				err = errors.Wrap(err, "failed to turn off main switch")
			}
			cmd.err <- err
			close(cmd.err)
		case troggleHeadLight:
			err := c.troggleHeadLight()
			if err != nil {
				err = errors.Wrap(err, "failed to troggle headlight")
			}
			cmd.err <- err
			close(cmd.err)
		case getMainSwitchStatus:
			var err error
			cmd.val, err = c.getMainSwitchStatus()
			if err != nil {
				err = errors.Wrap(err, "failed to get main switch status")
			}
			cmd.err <- err
			close(cmd.err)
		case getHeadLightStatus:
			var err error
			cmd.val, err = c.getHeadLightStatus()
			if err != nil {
				err = errors.Wrap(err, "failed to get headlight status")
			}
			cmd.err <- err
			close(cmd.err)
		case shutdown:
			errs := make([]error, 0, 10)
			errMainSwitch := c.mainSwitch.Close()
			if errMainSwitch != nil {
				errs = append(errs, errors.Wrap(errMainSwitch, "failed to close main switch pin"))
			}
			errHeadLight := c.headLight.Close()
			if errHeadLight != nil {
				errs = append(errs, errors.Wrap(errHeadLight, "failed to close head light pin"))
			}
			if len(errs) != 0 {
				errStrs := make([]string, 0, len(errs))
				for _, e := range errs {
					errStrs = append(errStrs, e.Error())
				}
				cmd.err <- errors.Errorf("failed to close electric controller: \n%s", strings.Join(errStrs, "\n"))
			}
			close(cmd.err)
			return
		default:
			cmd.err <- errors.New("unknown electric controller command")
			close(cmd.err)
		}
	}
}

// TurnOnMainSwitch TurnOnMainSwitch
func (c *Controller) TurnOnMainSwitch() error {
	cmd := command{
		typ: turnOnMainSwitch,
		err: make(chan error),
	}
	c.commands <- &cmd
	return <-cmd.err
}

// TurnOffMainSwitch TurnOffMainSwitch
func (c *Controller) TurnOffMainSwitch() error {
	cmd := command{
		typ: turnOffMainSwitch,
		err: make(chan error),
	}
	c.commands <- &cmd
	return <-cmd.err
}

// ToggleHeadLight ToggleHeadLight
func (c *Controller) ToggleHeadLight() error {
	cmd := command{
		typ: troggleHeadLight,
		err: make(chan error),
	}
	c.commands <- &cmd
	return <-cmd.err
}

// GetMainSwitchStatus GetMainSwitchStatus
func (c *Controller) GetMainSwitchStatus() (stat int, err error) {
	cmd := command{
		typ: getMainSwitchStatus,
		err: make(chan error),
	}
	c.commands <- &cmd
	err = <-cmd.err
	stat = cmd.val
	return
}

// GetHeadLightStatus GetHeadLightStatus
func (c *Controller) GetHeadLightStatus() (stat int, err error) {
	cmd := command{
		typ: getHeadLightStatus,
		err: make(chan error),
	}
	c.commands <- &cmd
	err = <-cmd.err
	stat = cmd.val
	return
}

// Close Close
func (c *Controller) Close() error {
	cmd := command{
		typ: shutdown,
		err: make(chan error),
	}
	c.commands <- &cmd
	return <-cmd.err
}
