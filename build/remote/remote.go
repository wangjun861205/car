package main

import (
	"buxiong/car/keyboard"
	"buxiong/car/model"
	"buxiong/car/network"
	"buxiong/car/remote"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	dailAddr := os.Getenv(model.DailAddr)
	client, err := network.NewClient(dailAddr, '\n')
	if err != nil {
		panic(err)
	}
	keyboardEvent := os.Getenv(model.KeyboardEvent)
	keyboard, err := keyboard.NewKeyboardReader(keyboardEvent)
	if err != nil {
		panic(err)
	}
	remote := remote.NewRemote(client, keyboard)
	remote.Run()
}
