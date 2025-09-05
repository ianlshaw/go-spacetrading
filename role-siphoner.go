package main

import "fmt"

func ApplyRoleSiphoner(ship Ship, all_waypoints_in_system []Waypoint) {
	return
	fmt.Println(ship.Fuel)
	AddWaypointToGraph(SystemGraph, all_waypoints_in_system[0])

	fmt.Println("ApplyRoleSiphoner")
	explosive_gas_waypoints := WaypointsWithTrait(all_waypoints_in_system, "EXPLOSIVE_GASES")
	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)
	closest_explosive_gas_waypoint := ClosestWaypointFromSliceToWaypoint(explosive_gas_waypoints, current_waypoint)
	if IsShipCargoFull(ship) {
		fmt.Println("ship cargo full")
		// if contract, goto contract dropoff destination
		// otherwise, goto most profitable marketplace
	}
	if IsShipAlreadyAtWaypoint(ship, closest_explosive_gas_waypoint.Symbol) {
		fmt.Println("Siphoner already at closest explosive gas waypoint")
		SiphonResources(ship.Symbol)
	} else {
		if IsShipDocked(ship) {
			OrbitShip(ship.Symbol)
		}
		CalculateShortestPathBetweenTwoWaypoints(SystemGraph, current_waypoint, closest_explosive_gas_waypoint)
		NavigateShip(ship.Symbol, closest_explosive_gas_waypoint.Symbol)
	}
	// Go to closest gas deposit

	// siphon until full

	// go to contract delivery location

	// give goods to contract
}
