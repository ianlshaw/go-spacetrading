package main

import (
	"fmt"
	"time"
)

func ApplyRoleTrader(ship Ship) time.Time {
	fmt.Println("[INFO] " + ship.Symbol + " ApplyRoleTrader")
	return ThreeHoursFromNow()
}