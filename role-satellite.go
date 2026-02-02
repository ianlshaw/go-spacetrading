package main

import (
	"fmt"
	"time"
)

//func ApplyRoleSatellite(ship Ship, markets_to_cover map[string]string, trade_routes []TradeRoute) time.Time {
func ApplyRoleSatellite(ship Ship, trade_routes []TradeRoute) time.Time {

	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + "ApplyRoleSatellite")
	// This logic prevents trading from begining until enough probes have been purchased to cover all markets
	//if !SatelliteToMarketAssignmentComplete(markets_to_cover) {
	//	fmt.Println("[INFO] Satellites have not yet been assigned to markets. Waiting...")
	//	return time.Now().Add(5 * time.Minute)
	//}

	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[INFO] IN_TRANSIT TO " + ship.Nav.Route.Destination.Symbol)
		fmt.Println("[INFO] Arrival " + ship.Nav.Route.Arrival)
		fmt.Println("[ERROR] ApplyRoleSatellite SHIP IN TRANSIT THIS SHOULD NOT HAPPEN")
		return StringToTimestamp(ship.Nav.Route.Arrival)
	}

	trade_route_missing_data_success, trade_route := TradeRouteMissingData(trade_routes)
	if !trade_route_missing_data_success {
		fmt.Println("[WARN] " + ship.Symbol + " All trade routes have data. Sleeping for 3 hours")
		return ThreeHoursFromNow()
	}

	// given a trade route
	// go to buy waypoint
	// get market details and update trade_routes with new data
	// go to sell waypoint and update trade_routes with new data
	// I should only be given trade routes which lack pricing data

	if IsShipDocked(ship){
		OrbitShip(ship.Symbol)
	}

	if trade_route.BuyMarketTradeGood.PurchasePrice == 0 {
		if IsShipAlreadyAtWaypoint(ship, trade_route.BuyMarketplaceWaypointSymbol){
			DockShip(ship.Symbol)
			UpdateTradeRoutesIncludingThisWaypoint(trade_route.BuyMarketplaceWaypointSymbol, trade_routes)
			return time.Now()
		}
		_, arrival_time := NavigateShip(ship.Symbol, trade_route.BuyMarketplaceWaypointSymbol)
		return arrival_time
	}

	if trade_route.SellMarketTradeGood.SellPrice == 0 {
		if IsShipAlreadyAtWaypoint(ship, trade_route.SellMarketplaceWaypointSymbol){
			DockShip(ship.Symbol)
			UpdateTradeRoutesIncludingThisWaypoint(trade_route.SellMarketplaceWaypointSymbol, trade_routes)
			return time.Now()
		}
		_, arrival_time := NavigateShip(ship.Symbol, trade_route.SellMarketplaceWaypointSymbol)
		return arrival_time
	}

	fmt.Println("[ERROR] Satellite was given a trade route which already has pricing data.")
	return time.Now()


	// find my assignment waypoint
	// am i there?
	// dock and get market
	// if no orbit and navigate there

	//var assigned_market_waypoint = markets_to_cover[ship.Symbol]

	// find the name of this satellite as a value in the markets_to_cover map, return the key of that value as assigned_market_waypoint
	//for market_symbol, assigned_satellite := range markets_to_cover {

		//print("ship.Symbol == ")
		//println(ship.Symbol)
		//print("assigned_satellite == ")
		//println(assigned_satellite)

	//	if assigned_satellite == ship.Symbol {
			//println("[DEBUG] satellite market assignment found")
	//		assigned_market_waypoint = market_symbol
	//		break
	//	}

	//}

	//println("[DEBUG] markets_to_cover:")

	//for k, v := range markets_to_cover {
	//	println(k)
	//	println(v)
	//}

	//if IsShipAlreadyAtWaypoint(ship, assigned_market_waypoint) {
	//	fmt.Println("[INFO] Already at assigned market waypoint")
	//	if !IsShipDocked(ship) {
	//		DockShip(ship.Symbol)
	//	}
	//	fmt.Println("[INFO] Updating trade data...")
	//	UpdateTradeRoutesIncludingThisWaypoint(assigned_market_waypoint, trade_routes)
	//} else {
	//	fmt.Println("[INFO] Not at assigned market waypoint, heading there now")
	//	if IsShipDocked(ship) {
	//		OrbitShip(ship.Symbol)
	//	}
	//	fmt.Println("[DEBUG] assigned_market_waypoint: " + assigned_market_waypoint)
	//	NavigateShip(ship.Symbol, assigned_market_waypoint)
	//}
	//fmt.Print("[ERROR] " + ship.Symbol)
	//fmt.Println(" uncaught branch, returning default 5 minute delay")
	//return time.Now().Add(5 * time.Minute)
}