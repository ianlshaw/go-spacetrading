package main

import "math"

func IsShipAlreadyAtWaypoint(ship_to_test Ship, waypoint_symbol string) bool {
	return (ship_to_test.Nav.WaypointSymbol == waypoint_symbol && ship_to_test.Nav.Status != "IN_TRANSIT")
}

func IsShipDocked(ship Ship) bool {
	return ship.Nav.Status == "DOCKED"
}

func IsShipCargoEmpty(ship Ship) bool {
	return ship.Cargo.Units == 0
}

func IsASatelliteDockedAtMarketplace(list_ships_result []Ship, waypoint_symbol string) (answer bool) {
	for _, ship := range list_ships_result {
		if ship.Registration.Role == "SATELLITE" {
			if ship.Nav.WaypointSymbol == waypoint_symbol {
				if ship.Nav.Status == "DOCKED" {
					return true
				}
			}
		}
	}
	return false
}

func DistanceBetweenTwoWaypoints(waypoint1 Waypoint, waypoint2 Waypoint) int {
	return DistanceBetweenTwoCoordinates(waypoint1.X, waypoint1.Y, waypoint2.X, waypoint2.Y)
}

func IsWaypointWithinDistanceOfWaypoint(waypoint1 Waypoint, waypoint2 Waypoint, distance int) bool {
	return DistanceBetweenTwoWaypoints(waypoint1, waypoint2) < distance
}

func IsWaypointWithinDistanceOfTwoWaypoints(waypoint_to_test Waypoint, origin_waypoint Waypoint, destination_waypoint Waypoint, max_distance int) bool {
	return IsWaypointWithinDistanceOfWaypoint(waypoint_to_test, origin_waypoint, max_distance) && IsWaypointWithinDistanceOfWaypoint(waypoint_to_test, destination_waypoint, max_distance)
}

func DistanceBetweenTwoCoordinates(waypoint1X int64, waypoint1Y int64, waypoint2X int64, waypoint2Y int64) (resultant_distance int) {
	//fmt.Println("[DEBUG] distance_between_two_coordinates")
	XIntermediate := waypoint1X - waypoint2X
	YIntermediate := waypoint1Y - waypoint2Y
	XSquared := XIntermediate * XIntermediate
	YSquared := YIntermediate * YIntermediate
	XPlusY := XSquared + YSquared
	XPlusYFloat := float64(XPlusY)

	distance_float := math.Sqrt(XPlusYFloat)
	resultant_distance = int(distance_float)
	return resultant_distance
}

func GetWaypointCoordinate(waypoint Waypoint) (waypointX int64, waypointY int64) {
	return waypoint.X, waypoint.Y
}

func CountTradeGoodCargo(ship Ship, trade_good_symbol string) int64 {
	inventory := ship.Cargo.Inventory
	for _, trade_good := range inventory {
		if trade_good_symbol == trade_good.Symbol {
			return trade_good.Units
		}
	}
	return 0
}

func HowManySatellitesDoIOwn() int {
	satellite_count := 0
	list_ships := ListShips()
	for _, ship := range list_ships {
		if ship.Frame.Name == "SATELLITE" {
			satellite_count++
		}
	}
	return satellite_count
}
