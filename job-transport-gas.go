package main

import (
	"fmt"
	"time"
)

func ApplyRoleTransportGas(ship Ship) time.Time {
	fmt.Println("[INFO] " + ship.Symbol + " ApplyRoleTransportGas")
	return ThreeHoursFromNow()
}
