package main

import (
	"buxiong/car/controller"
	"buxiong/car/driver"
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()
	ld, err := driver.NewDriver(rpio.Pin(24), rpio.Pin(23), 1, 1000)
	if err != nil {
		panic(err)
	}
	rd, err := driver.NewDriver(rpio.Pin(5), rpio.Pin(6), 0, 1000)
	if err != nil {
		panic(err)
	}
	ctl := controller.NewController(
		1000,
		400,
		100,
		ld,
		rd,
	)
	// for i := 0; i < 7; i++ {
	// 	ctl.Forward()
	// 	fmt.Println(ctl.Status())
	// 	time.Sleep(time.Second * 5)
	// }
	// for i := 0; i < 7; i++ {
	// 	ctl.Backward()
	// 	fmt.Println(ctl.Status())
	// 	time.Sleep(time.Second * 5)
	// }
	// for i := 0; i < 7; i++ {
	// 	ctl.Backward()
	// 	fmt.Println(ctl.Status())
	// 	time.Sleep(time.Second * 5)
	// }
	// for i := 0; i < 7; i++ {
	// 	ctl.Forward()
	// 	fmt.Println(ctl.Status())
	// 	time.Sleep(time.Second * 5)
	// }
	for i := 0; i < 7; i++ {
		ctl.TurnLeft()
		fmt.Println(ctl.Status())
		time.Sleep(time.Second * 5)
	}
	for i := 0; i < 7; i++ {
		ctl.TurnRight()
		fmt.Println(ctl.Status())
		time.Sleep(time.Second * 5)
	}
	ctl.Close()
}
