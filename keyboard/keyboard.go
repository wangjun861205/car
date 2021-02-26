package keyboard

import (
	"github.com/MarinX/keylogger"
	"github.com/pkg/errors"
)

type keyboardReader struct {
	logger  *keylogger.KeyLogger
	doneIn  chan struct{}
	doneOut chan struct{}
	out     chan keylogger.InputEvent
}

// NewKeyboardReader NewKeyboardReader
func NewKeyboardReader(inputEvent string) (*keyboardReader, error) {
	logger, err := keylogger.New(inputEvent)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to init keyboard(event: %s)", inputEvent)
	}
	return &keyboardReader{
		logger,
		make(chan struct{}),
		make(chan struct{}),
		make(chan keylogger.InputEvent),
	}, nil
}

func (k *keyboardReader) Run() {
	events := k.logger.Read()
	for {
		select {
		case <-k.doneIn:
			k.logger.Close()
			k.doneOut <- struct{}{}
		case e := <-events:
			k.out <- e
		}
	}
}

func (k *keyboardReader) Stop() {
	k.doneIn <- struct{}{}
	<-k.doneOut
}

func (k *keyboardReader) Out() chan keylogger.InputEvent {
	return k.out
}
