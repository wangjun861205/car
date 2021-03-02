package main

import (
	"buxiong/car/balancer"
	"buxiong/car/driver"
	"fmt"
	"math"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func compare(ld, rd *driver.Driver, lb, rb *balancer.Balancer, lduty, rduty uint64) (ls, rs uint64) {
	ld.Brake()
	rd.Brake()
	lb.Reset()
	lb.Reset()
	ld.Forward(lduty)
	rd.Forward(rduty)
	time.Sleep(5 * time.Second)
	ls, rs = lb.GetSpeed(), rb.GetSpeed()
	ld.Brake()
	rd.Brake()
	lb.Reset()
	rb.Reset()
	return
}

func main() {
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	lb := balancer.NewBalancer(17, 27)
	rb := balancer.NewBalancer(22, 26)
	go lb.Run()
	go rb.Run()
	ld, err := driver.NewDriver(rpio.Pin(24), rpio.Pin(23), 1, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	rd, err := driver.NewDriver(rpio.Pin(5), rpio.Pin(6), 0, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 4; i < 10; i++ {
		left, right := uint64(i*100), uint64(i*100)
		var leftSpeed, rightSpeed uint64
		for {
			ld.Forward(left)
			rd.Forward(right)
			time.Sleep(time.Millisecond * 100)
			lb.Reset()
			rb.Reset()
			time.Sleep(time.Second)
			ls, rs := lb.GetSpeed(), rb.GetSpeed()
			diff := math.Abs(float64(ls-rs) / float64(ls+rs) * 1000)
			if diff <= 5 {
				leftSpeed, rightSpeed = ls, rs
				break
			}
			switch {
			case ls < rs:
				left += 2
				right -= 2
				continue
			case ls > rs:
				left -= 2
				right += 2
				continue
			}
		}
		fmt.Printf("level: %d, left pwm: %d, right pwm: %d, left speed: %d, right speed: %d, left rate: %2f, right rate: %2f\n", i*100, left, right, leftSpeed, rightSpeed, float64(left-150)/float64(leftSpeed), float64(right-150)/float64(rightSpeed))

	}
	ld.Close()
	rd.Close()
	lb.Close()
	rb.Close()

}
