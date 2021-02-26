package keyboard

import (
	"fmt"
	"testing"
	"time"

	"github.com/MarinX/keylogger"
)

func TestKeyboardReader(t *testing.T) {
	reader, err := NewKeyboardReader("dev/input/event5")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Stop()
	reader.Run()
	for e := range reader.out {
		fmt.Println(e.KeyString())
	}
}

func TestKeyLogger(t *testing.T) {
	logger, err := keylogger.New("/dev/input/event5")
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	go func() {
		for event := range logger.Read() {
			fmt.Println(event.KeyString())
		}
	}()
	time.Sleep(10 * time.Second)
}
