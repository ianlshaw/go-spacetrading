package main

import (
	"fmt"
	"time"
)

func ApplyRoleTrader(ship Ship, all_waypoints_in_system []Waypoint, all_markets_in_system []Market, command_ship_waypoint_symbol string) time.Time {
	fmt.Println("[INFO] " + ship.Symbol + " ApplyRoleTrader")

	if IsShipAlreadyAtWaypoint(ship, command_ship_waypoint_symbol) {
		TransferCargo(ship.Symbol, "TVRJ-TEST-3-1", "FUEL", 4)
		return time.Now().Add(3 * time.Hour)
	}

	if IsShipCargoEmpty(ship) {
		PurchaseCargo(ship.Symbol, "FUEL", 40)
	}
	OrbitShip(ship.Symbol)
	current_waypoint :=	WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)
	command_ship_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, command_ship_waypoint_symbol)
	closest_market_to_command_ship := ClosestMarketToWaypoint(command_ship_waypoint, all_waypoints_in_system, all_markets_in_system)
	closest_market_waypoint_to_command_ship := WaypointFromWaypointSymbol(all_waypoints_in_system, closest_market_to_command_ship.Symbol)
	

	if IsShipAlreadyAtWaypoint(ship, closest_market_to_command_ship.Symbol) {
		_, arrival_time := NavigateShip(ship.Symbol, command_ship_waypoint_symbol)
		return arrival_time
	}

	path, _ := CalculateShortestPathBetweenTwoWaypoints(MarketplaceGraph, current_waypoint, closest_market_waypoint_to_command_ship)
	arrival_time := FollowPath(ship, path)
	return arrival_time
}