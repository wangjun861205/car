package driver

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/fsnotify.v0"
)

const basePath = "/sys/class/pwm/pwmchip0"

func exportPath(pin uint8) string {
	return fmt.Sprintf(path.Join(basePath, "export"))
}

func unexportPath(pin uint8) string {
	return fmt.Sprintf(path.Join(basePath, "unexport"))
}

func periodPath(pin uint8) string {
	return fmt.Sprintf(path.Join(basePath, fmt.Sprintf("pwm%d/period", pin)))
}

func dutyCyclePath(pin uint8) string {
	return fmt.Sprintf(path.Join(basePath, fmt.Sprintf("pwm%d/duty_cycle", pin)))
}

func enablePath(pin uint8) string {
	return fmt.Sprintf(path.Join(basePath, fmt.Sprintf("pwm%d/enable", pin)))
}

func isExported(num uint8) (bool, error) {
	stat, err := os.Stat(path.Join(basePath, fmt.Sprintf("pwm%d", num)))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, "isExported() failed")
	}
	if !stat.IsDir() {
		return false, errors.Wrap(errors.Errorf("pwm%d is not a directory", num), "isExported() failed")
	}
	return true, nil
}

// PWM pwm controller
type PWM struct {
	num       uint8
	export    *os.File
	unexport  *os.File
	period    *os.File
	dutyCycle *os.File
	enable    *os.File
}

// NewPWM NewPWM
func NewPWM(num uint8) (*PWM, error) {
	export, err := os.OpenFile(exportPath(num), os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0770)
	if err != nil {
		return nil, errors.Wrap(err, "NewPWM() failed")
	}
	exported, err := isExported(num)
	if err != nil {
		export.Close()
		return nil, errors.Wrap(err, "NewPWM() failed")
	}
	if !exported {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil, errors.Wrap(err, "NewPWM() failed")
		}
		defer watcher.Close()
		if err := watcher.Watch(basePath); err != nil {
			return nil, errors.Wrap(err, "NewPWM() failed")
		}
		if _, err := export.WriteString(fmt.Sprintf("%d\n", num)); err != nil {
			return nil, errors.Wrap(err, "NewPWM() failed")
		}
		timer := time.NewTimer(time.Second)
		var periodReady, dutyCycleReady, enableReady bool
	OUTER:
		for {
			select {
			case event := <-watcher.Event:
				if event.IsModify() {
					switch event.Name {
					case path.Join(basePath, fmt.Sprintf("pwm%d", num)):
						if err := watcher.Watch(path.Join(basePath, fmt.Sprintf("pwm%d", num))); err != nil {
							return nil, errors.Wrap(err, "NewPWM() failed")
						}
					case periodPath(num):
						periodReady = true
						if periodReady && dutyCycleReady && enableReady {
							break OUTER
						}
					case dutyCyclePath(num):
						dutyCycleReady = true
						if periodReady && dutyCycleReady && enableReady {
							break OUTER
						}
					case enablePath(num):
						enableReady = true
						if periodReady && dutyCycleReady && enableReady {
							break OUTER
						}
					}
				}
			case <-timer.C:
				return nil, errors.Wrap(errors.Errorf("timeout(num : %d)", num), "NewPWM() failed")
			}
		}
	}
	unexport, err := os.OpenFile(unexportPath(num), os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0770)
	if err != nil {
		return nil, errors.Wrap(err, "NewPWM() failed")
	}
	period, err := os.OpenFile(periodPath(num), os.O_RDWR|os.O_SYNC|os.O_TRUNC, 0770)
	if err != nil {
		return nil, errors.Wrap(err, "NewPWM() failed")
	}
	dutyCycle, err := os.OpenFile(dutyCyclePath(num), os.O_RDWR|os.O_SYNC|os.O_TRUNC, 0770)
	if err != nil {
		return nil, errors.Wrap(err, "NewPWM() failed")
	}
	enable, err := os.OpenFile(enablePath(num), os.O_RDWR|os.O_SYNC|os.O_TRUNC, 0770)
	if err != nil {
		return nil, errors.Wrap(err, "NewPWM() failed")
	}
	return &PWM{
		num:       num,
		export:    export,
		unexport:  unexport,
		period:    period,
		dutyCycle: dutyCycle,
		enable:    enable,
	}, nil
}

// SetPeriod SetPeriod
func (p *PWM) SetPeriod(period uint64) error {
	if _, err := p.period.WriteString(fmt.Sprintf("%d\n", period)); err != nil {
		return errors.Wrap(err, "SetPeriod() failed")
	}
	if _, err := p.period.Seek(0, 0); err != nil {
		return errors.Wrap(err, "SetPeriod() failed")
	}
	return nil
}

// GetPeriod GetPeriod
func (p *PWM) GetPeriod() (uint64, error) {
	buf := bufio.NewReader(p.period)
	s, err := buf.ReadString('\n')
	if err != nil {
		return 0, errors.Wrap(err, "GetPeriod() failed")
	}
	u, err := strconv.ParseUint(strings.Trim(s, "\n"), 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "GetPeriod() failed")
	}
	return u, nil
}

// SetDutyCycle SetDutyCycle
func (p *PWM) SetDutyCycle(dutyCycle uint64) error {
	if _, err := p.dutyCycle.WriteString(fmt.Sprintf("%d\n", dutyCycle)); err != nil {
		return errors.Wrap(err, "SetDutyCycle() failed")
	}
	if _, err := p.dutyCycle.Seek(0, 0); err != nil {
		return errors.Wrap(err, "SetDutyCycle() failed")
	}
	return nil
}

// GetDutyCycle GetDutyCycle
func (p *PWM) GetDutyCycle() (uint64, error) {
	buf := bufio.NewReader(p.dutyCycle)
	s, err := buf.ReadString('\n')
	if err != nil {
		return 0, errors.Wrap(err, "GetDutyCycle() failed")
	}
	u, err := strconv.ParseUint(strings.Trim(s, "\n"), 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "GetDutyCycle() failed")
	}
	return u, nil
}

// Enable Enable
func (p *PWM) Enable() error {
	if _, err := p.enable.WriteString("1\n"); err != nil {
		return errors.Wrap(err, "Enable() failed")
	}
	if _, err := p.enable.Seek(0, 0); err != nil {
		return errors.Wrap(err, "Enable() failed")
	}
	return nil
}

// Disable Disable
func (p *PWM) Disable() error {
	if _, err := p.enable.WriteString("0\n"); err != nil {
		return errors.Wrap(err, "Disable() failed")
	}
	if _, err := p.enable.Seek(0, 0); err != nil {
		return errors.Wrap(err, "Disable() failed")
	}
	return nil
}

// IsEnabled IsEnabled
func (p *PWM) IsEnabled() (bool, error) {
	buf := bufio.NewReader(p.enable)
	s, err := buf.ReadString('\n')
	if err != nil {
		return false, errors.Wrap(err, "IsEnabled() failed")
	}
	s = strings.Trim(s, "\n")
	return s == "1", nil
}

// Close Close
func (p *PWM) Close() error {
	if err := p.Disable(); err != nil {
		return err
	}
	if err := p.enable.Close(); err != nil {
		return err
	}
	if err := p.dutyCycle.Close(); err != nil {
		return err
	}
	if err := p.period.Close(); err != nil {
		return err
	}
	if _, err := p.unexport.WriteString(fmt.Sprintf("%d", p.num)); err != nil {
		return err
	}
	if err := p.export.Close(); err != nil {
		return err
	}
	if err := p.unexport.Close(); err != nil {
		return err
	}
	return nil
}
