package main

import (
	"fmt"
	"time"
)

func ApplyRoleTransportOre(ship Ship) time.Time {
	fmt.Println("[INFO] " + ship.Symbol + " ApplyRoleTransportOre")
	return ThreeHoursFromNow()
}
