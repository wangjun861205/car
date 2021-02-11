package server

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
	addr := os.Getenv("DISPLAY_SERVER_LISTEN_ADDR")
	server, err := network.NewServer(addr, '\n')
	if err != nil {
		panic(err)
	}
	display := camera.NewDisplayServer(server)
	display.Run()
}
