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
	fmt.Println("[DEBUG] ApplyRoleSiphoner " + ship.Symbol)

	ship = GetShip(ship.Symbol)

	//
	for _, item := range ship.Cargo.Inventory {
		fmt.Print("[DEBUG] Inventory: ")
		fmt.Print(item.Name)
		fmt.Print(" ")
		fmt.Print(item.Units)
		fmt.Println()
	}
	//

	// ad-hoc dump inventory because buy flow is incorrect
	//for _, cargo_trade_good := range ship.Cargo.Inventory {
	//	units := CountTradeGoodCargo(ship, cargo_trade_good.Symbol)
	//	JettisonCargo(ship, cargo_trade_good.Symbol, units)
	//}

	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)
	gas_giant_waypoints := ListWaypointInSystemByType(base_system_symbol, "GAS_GIANT")
	closest_gas_giant_waypoint := ClosestWaypointFromSliceToWaypoint(gas_giant_waypoints, current_waypoint)
	closest_market_to_closest_gas_giant := ClosestMarketToWaypoint(closest_gas_giant_waypoint, all_waypoints_in_system, all_markets_in_system)
	closest_market_to_closest_gas_giant_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, closest_market_to_closest_gas_giant.Symbol)
	target_trade_good := contract.Terms.Deliver[0].TradeSymbol	
	contract_delivery_destination_symbol := contract.Terms.Deliver[0].DestinationSymbol
	contract_delivery_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, contract_delivery_destination_symbol)
	
	// is cargo full? -> go to closest market, followed by delivery waypoint

	if IsShipCargoFull(ship) {
		if IsShipAlreadyAtWaypoint(ship, contract_delivery_destination_symbol) {
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}
			units := CountTradeGoodCargo(ship, target_trade_good)
			DeliverCargoToContract(contract.ID, ship.Symbol, target_trade_good, units)
			if CanContractBeCompleted(contract) {
				FulfillContract(contract.ID)
				NegotiateContract(ship.Symbol)
			}
			path, _ := CalculateShortestPathBetweenTwoWaypoints(MarketplaceGraph, current_waypoint, closest_market_to_closest_gas_giant_waypoint)
			arrival_time := FollowPath(ship, path)
			return arrival_time
		}
		if IsShipAlreadyAtWaypoint(ship, closest_market_to_closest_gas_giant_waypoint.Symbol) {
			path, _ := CalculateShortestPathBetweenTwoWaypoints(MarketplaceGraph, current_waypoint, contract_delivery_waypoint)
			arrival_time := FollowPath(ship, path)
			return arrival_time
		}
		if IsShipAlreadyAtWaypoint(ship, closest_gas_giant_waypoint.Symbol) {
			_, arrival_time := NavigateShip(ship.Symbol, closest_market_to_closest_gas_giant_waypoint.Symbol)
			return arrival_time
		}
	}

	// cargo not full
	if IsShipAlreadyAtWaypoint(ship, closest_gas_giant_waypoint.Symbol) {
		fmt.Println("[DEBUG] ship is already at closest_gas_giant_waypoint")
		// TODO
		siphon_result := SiphonResources(ship.Symbol)
		siphon := siphon_result.Siphon
		yield := siphon.Yield
		cooldown := siphon_result.Cooldown
		expiration := cooldown.Expiration
		cargo := siphon_result.Cargo

		if target_trade_good != "" && yield.Symbol != target_trade_good && yield.Units != 0 {
			cargo = JettisonCargo(ship, yield.Symbol, yield.Units)
		}

		if cargo.Units == cargo.Capacity {
		// head to nearest market, to get back onto the market graph
			_, arrival_time := NavigateShip(ship.Symbol, closest_market_to_closest_gas_giant_waypoint.Symbol)
			return arrival_time
		}
		expiration_timestamp := StringToTimestamp(expiration)
		fmt.Print("[DEBUG] On cooldown after siphoning until ")
		fmt.Println(expiration)
		return expiration_timestamp
	}

	if IsShipAlreadyAtWaypoint(ship, closest_market_to_closest_gas_giant_waypoint.Symbol) {
		_, arrival_time := NavigateShip(ship.Symbol, closest_gas_giant_waypoint.Symbol)
		return arrival_time
	}

	if IsShipAlreadyAtWaypoint(ship, contract_delivery_destination_symbol) {
		// TODO Dump excess cargo here
		path, _ := CalculateShortestPathBetweenTwoWaypoints(MarketplaceGraph, current_waypoint, closest_market_to_closest_gas_giant_waypoint)
		arrival_time := FollowPath(ship, path)
		return arrival_time
	}

	path, _ := CalculateShortestPathBetweenTwoWaypoints(MarketplaceGraph, current_waypoint, closest_market_to_closest_gas_giant_waypoint)
	arrival_time := FollowPath(ship, path)
	return arrival_time

	fmt.Println("[ERROR] ApplyRoleSiphoner returning time.Now() because uncaught branch")
	return time.Now()
}
