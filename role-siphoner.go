package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
	"time"
)

var SiphonerMarketplaceGraph dijkstra.Graph = make(dijkstra.Graph)

func ApplyRoleSiphoner(ship Ship,
						all_waypoints_in_system []Waypoint, 
						all_markets_in_system []Market, 
						contract Contract) time.Time {
	fmt.Println("[DEBUG] ApplyRoleSiphoner")

	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)
	explosive_gas_waypoints := WaypointsWithTrait(all_waypoints_in_system, "EXPLOSIVE_GASES")
	closest_explosive_gas_waypoint := ClosestWaypointFromSliceToWaypoint(explosive_gas_waypoints, current_waypoint)
	closest_market_to_closest_explosive_gas_waypoint := ClosestMarketToWaypoint(closest_explosive_gas_waypoint, all_waypoints_in_system, all_markets_in_system)

	target_trade_good := contract.Terms.Deliver[0].TradeSymbol	
	contract_delivery_destination_symbol := contract.Terms.Deliver[0].DestinationSymbol
	contract_delivery_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, contract_delivery_destination_symbol)

	if IsShipAlreadyAtWaypoint(ship, closest_explosive_gas_waypoint.Symbol) {
		fmt.Println("Siphoner already at closest explosive gas waypoint")
		siphon_result := SiphonResources(ship.Symbol)
		siphon := siphon_result.Siphon
		yield := siphon.Yield
		cooldown := siphon_result.Cooldown
		expiration := cooldown.Expiration
		cargo := siphon_result.Cargo

		if target_trade_good != "" {
			if yield.Symbol != target_trade_good {
				JettisonCargo(ship, yield.Symbol, yield.Units)
			}
			
		}

		// We need extra logic to stage from the closest market to the explosive gas waypoint, since that waypoint wont exist in the graph of markets

		if cargo.Units == cargo.Capacity {
			// head to contract delivery waypoint
			path, _ := CalculateShortestPathBetweenTwoWaypoints(SiphonerMarketplaceGraph, current_waypoint, contract_delivery_waypoint)
			expiration := FollowPath(ship, path)
			return expiration
		}

		expiration_timestamp := StringToTimestamp(expiration)
		return expiration_timestamp
		// ship is not at closest_explosive_gas_waypoint
	} else {
		if IsShipAlreadyAtWaypoint(ship, closest_market_to_closest_explosive_gas_waypoint.Symbol) {

		}
		path, _ := CalculateShortestPathBetweenTwoWaypoints(SiphonerMarketplaceGraph, current_waypoint, closest_explosive_gas_waypoint)
		expiration := FollowPath(ship, path)
		return expiration
	}
	// Go to closest gas deposit

	// siphon until full

	if IsShipCargoFull(ship) {
		fmt.Println("ship cargo full")
		// if contract, goto contract dropoff destination
		// otherwise, goto most profitable marketplace
	}

	// go to contract delivery location

	// give goods to contract

	fmt.Println("[ERROR] ApplyRoleSiphoner returned default")
	return time.Now()
}
