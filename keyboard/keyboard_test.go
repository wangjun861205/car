package keyboard

import (
	"fmt"
	"testing"
)

func TestKeyboardReader(t *testing.T) {
	reader, err := NewKeyboardReader()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Stop()
	reader.Run()
	for e := range reader.out {
		fmt.Println(e.KeyString())
	}
}
