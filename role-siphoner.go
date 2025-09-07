package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
)

var SiphonerMarketplaceGraph dijkstra.Graph = make(dijkstra.Graph)

func ApplyRoleSiphoner(ship Ship, all_waypoints_in_system []Waypoint, all_markets_in_system []Market) {
	fmt.Println("[DEBUG] ApplyRoleSiphoner")



	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)
	explosive_gas_waypoints := WaypointsWithTrait(all_waypoints_in_system, "EXPLOSIVE_GASES")


	for _,  explosive_gas_waypoint := range explosive_gas_waypoints {
		for _, market := range all_markets_in_system {
			market_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, market.Symbol)
			distance := DistanceBetweenTwoWaypoints(explosive_gas_waypoint, market_waypoint)
			if distance <= 80 {
				fmt.Println(market_waypoint.Symbol)
				fmt.Println(explosive_gas_waypoint.Symbol)
			}
		}
	}


	return


	closest_explosive_gas_waypoint := ClosestWaypointFromSliceToWaypoint(explosive_gas_waypoints, current_waypoint)


	// find closest market to to closest explosive gas waypoint

	ClosestMarketToWaypoint(closest_explosive_gas_waypoint, all_waypoints_in_system, all_markets_in_system)
	//
	//for _, waypoint := range explosive_gas_waypoints {
	//	path, _ := CalculateShortestPathBetweenTwoWaypoints(SystemGraph, current_waypoint, waypoint)
	//	fmt.Println(path)
	//}
	//

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
		path, _ := CalculateShortestPathBetweenTwoWaypoints(SiphonerMarketplaceGraph, current_waypoint, closest_explosive_gas_waypoint)
		fmt.Println(path)
		FollowPath(ship, path)
	}
	// Go to closest gas deposit

	// siphon until full

	// go to contract delivery location

	// give goods to contract
}
