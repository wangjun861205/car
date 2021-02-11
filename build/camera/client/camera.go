package main

import (
	"buxiong/car/camera"
	"buxiong/car/network"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	addr := os.Getenv("DISPLAY_SERVER_ADDR")
	client, err := network.NewClient(addr, '\n')
	if err != nil {
		panic(err)
	}
	camera, err := camera.NewCamera(0, client)
	if err != nil {
		panic(err)
	}
	camera.Run()

}
