package remote

import "github.com/MarinX/keylogger"

type client interface {
	Run()
	Stop()
	RegisterHandler(func(b []byte))
	Write(b []byte) error
}

type keyboardReader interface {
	Run()
	Stop()
	Out() chan keylogger.InputEvent
}
