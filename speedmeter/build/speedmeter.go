package main

import (
	"buxiong/car/driver"
	"buxiong/car/speedmeter"
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// func compare(ld, rd *driver.Driver, b *balancer.Balancer, lduty, rduty uint64) (ls, rs float64) {
// 	ld.Brake()
// 	rd.Brake()
// 	b.Reset()
// 	ld.Forward(lduty)
// 	rd.Forward(rduty)
// 	time.Sleep(5 * time.Second)
// 	ls, rs = b.LeftSpeed(), b.RightSpeed()
// 	ld.Brake()
// 	rd.Brake()
// 	b.Reset()
// 	return
// }

func pid(power, target, speed float64) float64 {
	diff := target - speed
	return power + 0.1*diff
}

func main() {
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	lsm, err := speedmeter.NewSpeedMeter(17, 27)
	if err != nil {
		fmt.Println(err)
		return
	}
	go lsm.Run()
	defer lsm.Close()
	rsm, err := speedmeter.NewSpeedMeter(26, 22)
	if err != nil {
		fmt.Println(err)
		return
	}
	go rsm.Run()
	defer rsm.Close()
	ld, err := driver.NewDriver(rpio.Pin(24), rpio.Pin(23), 1, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ld.Close()
	rd, err := driver.NewDriver(rpio.Pin(5), rpio.Pin(6), 0, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rd.Close()
	timer := time.NewTimer(time.Second * 10)
	target := 150.0
OUTER:
	for {
		select {
		case <-timer.C:
			ld.Brake()
			rd.Brake()
			break OUTER
		default:
			ls, rs := lsm.Measure(), rsm.Measure()
			lcp, rcp := float64(ld.GetDuty()), float64(rd.GetDuty())
			lp, rp := pid(lcp, target, ls), pid(rcp, target, rs)
			ld.SetDuty(int64(lp))
			rd.SetDuty(int64(rp))
			fmt.Printf("left speed(power): %.4f(%.4f), right speed(power): %.4f(%.4f), left diff: %.4f, right diff: %.4f\n", ls, lp, rs, rp, target-ls, target-rs)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// func main() {
// 	reader := bufio.NewReader(os.Stdin)
// 	var kp, ki, kd, target float64
// 	for {
// 		fmt.Printf("please enter command: ")
// 		cmd, err := reader.ReadString('\n')
// 		if err != nil {
// 			fmt.Println(errors.Wrap(err, "failed to read command"))
// 			return
// 		}
// 		switch strings.Trim(cmd, "\n") {
// 		case "config":
// 		CONFIG_KP:
// 			for {
// 				fmt.Printf("please enter Kp: ")
// 				ps, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read Kp"))
// 					return
// 				}
// 				kp, err = strconv.ParseFloat(strings.Trim(ps, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid Kp: %s", ps))
// 					continue CONFIG_KP
// 				}
// 				break
// 			}
// 		CONFIG_KI:
// 			for {
// 				fmt.Printf("please enter Ki: ")
// 				is, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read Ki"))
// 					return
// 				}
// 				ki, err = strconv.ParseFloat(strings.Trim(is, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid Ki: %s", is))
// 					continue CONFIG_KI
// 				}
// 				break
// 			}
// 		CONFIG_KD:
// 			for {
// 				fmt.Printf("please enter Kd: ")
// 				ds, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read Kd"))
// 					return
// 				}
// 				kd, err = strconv.ParseFloat(strings.Trim(ds, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid Kd: %s", ds))
// 					continue CONFIG_KD
// 				}
// 				break
// 			}
// 		CONFIG_TARGET:
// 			for {
// 				fmt.Printf("please enter target: ")
// 				ts, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read target"))
// 					return
// 				}
// 				target, err = strconv.ParseFloat(strings.Trim(ts, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid target: %s", ts))
// 					continue CONFIG_TARGET
// 				}
// 				break
// 			}
// 		case "start":
// 			var speed float64
// 			for j := 0; j < 50; j++ {
// 				d := target - speed
// 				power := speed*2 + 200 + kp*d + ki + kd*d*d
// 				if power < 0 {
// 					speed = (power + 200) / 2
// 				} else {
// 					speed = (power - 200) / 2
// 				}
// 				fmt.Printf("target: %f, power: %f, speed: %f, diff: %f\n", target, power, speed, target-speed)
// 			}
// 		case "kp":
// 			for {
// 				fmt.Printf("please enter Kp: ")
// 				ps, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read Kp"))
// 					return
// 				}
// 				kp, err = strconv.ParseFloat(strings.Trim(ps, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid Kp: %s", ps))
// 					continue
// 				}
// 				break
// 			}
// 		case "ki":
// 			for {
// 				fmt.Printf("please enter Ki: ")
// 				is, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read Ki"))
// 					return
// 				}
// 				ki, err = strconv.ParseFloat(strings.Trim(is, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid Ki: %s", is))
// 					continue
// 				}
// 				break
// 			}
// 		case "kd":
// 			for {
// 				fmt.Printf("please enter Kd: ")
// 				ds, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read Kd"))
// 					return
// 				}
// 				kd, err = strconv.ParseFloat(strings.Trim(ds, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid Kd: %s", ds))
// 					continue
// 				}
// 				break
// 			}
// 		case "target":
// 			for {
// 				fmt.Printf("please enter target: ")
// 				ts, err := reader.ReadString('\n')
// 				if err != nil {
// 					fmt.Println(errors.Wrap(err, "failed to read target"))
// 					return
// 				}
// 				target, err = strconv.ParseFloat(strings.Trim(ts, "\n"), 64)
// 				if err != nil {
// 					fmt.Println(errors.Wrapf(err, "invalid target: %s", ts))
// 					continue
// 				}
// 				break
// 			}
// 		case "show":
// 			fmt.Printf("target: %f, Kp: %f, Ki: %f, Kd: %f\n", target, kp, ki, kd)
// 		case "exit":
// 			return
// 		}
// 	}
// }
