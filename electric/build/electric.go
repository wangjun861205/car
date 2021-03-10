package main

import (
	"buxiong/car/electric"
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
)

func main() {
	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer chip.Close()
	mainSwitchLine, err := chip.RequestLine(16, gpiod.AsOutput())
	if err != nil {
		fmt.Println(err)
		return
	}
	headLightLine, err := chip.RequestLine(25, gpiod.AsOutput())
	if err != nil {
		fmt.Println(err)
		return
	}
	ctl := electric.NewController(mainSwitchLine, headLightLine)
	go ctl.Run()
	defer ctl.Close()
	fmt.Println("turn on main switch...")
	time.Sleep(time.Second)
	if err := ctl.TurnOnMainSwitch(); err != nil {
		fmt.Println(err)
		return
	}
	stat, err := ctl.GetMainSwitchStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("main switch status: ", stat)
	time.Sleep(time.Second)
	fmt.Println("turn on headlight")
	if err := ctl.ToggleHeadLight(); err != nil {
		fmt.Println(err)
		return
	}
	stat, err = ctl.GetHeadLightStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("headlight status: ", stat)
	time.Sleep(time.Second)
	fmt.Println("turn off headlight")
	if err := ctl.ToggleHeadLight(); err != nil {
		fmt.Println(err)
		return
	}
	stat, err = ctl.GetHeadLightStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("headlight status: ", stat)
	time.Sleep(time.Second)
	fmt.Println("turn off main switch")
	if err := ctl.TurnOffMainSwitch(); err != nil {
		fmt.Println(err)
		return
	}
	stat, err = ctl.GetMainSwitchStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("headlight status: ", stat)
	time.Sleep(time.Second)
}
