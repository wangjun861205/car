package main

import (
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
	driver1, err := driver.NewDriver(rpio.Pin(5), rpio.Pin(6), 0, 1000)
	if err != nil {
		panic(err)
	}
	driver2, err := driver.NewDriver(rpio.Pin(24), rpio.Pin(23), 1, 1000)
	for i := 1; i <= 10; i++ {
		driver1.Forward(uint64(i) * 100)
		driver2.Forward(uint64(i) * 100)
		fmt.Println(driver1.Status(), driver2.Status())
		time.Sleep(5 * time.Second)
	}
	driver1.Brake()
	driver2.Brake()
	fmt.Println(driver1.Status(), driver2.Status())
	for i := 1; i <= 10; i++ {
		driver1.Backward(uint64(i) * 100)
		driver2.Backward(uint64(i) * 100)
		fmt.Println(driver1.Status(), driver2.Status())
		time.Sleep(5 * time.Second)
	}
	driver1.Glide()
	driver2.Glide()
	fmt.Println(driver1.Status(), driver2.Status())
	driver1.Close()
	driver2.Close()
}
