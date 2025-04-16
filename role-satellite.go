package main

import "fmt"

func ApplyRoleSatellite(ship Ship, markets_to_cover map[string]string, trade_routes []TradeRoute) {

	fmt.Println("[INFO] ApplyRoleSatellite " + ship.Symbol)

	if !SatelliteToMarketAssignmentComplete(markets_to_cover) {
		fmt.Println("[INFO] Satellites have not yet been assigned to markets. Waiting...")
		return
	}

	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[INFO] IN_TRANSIT TO " + ship.Nav.Route.Destination.Symbol)
		fmt.Println("[INFO] Arrival " + ship.Nav.Route.Arrival)
		fmt.Println()
		return
	}

	// find my assignment waypoint
	// am i there?
	// dock and get market
	// if no orbit and navigate there

	var assigned_market_waypoint = markets_to_cover[ship.Symbol]

	// find the name of this satellite as a value in the markets_to_cover map, return the key of that value as assigned_market_waypoint
	for market_symbol, assigned_satellite := range markets_to_cover {

		//print("ship.Symbol == ")
		//println(ship.Symbol)
		//print("assigned_satellite == ")
		//println(assigned_satellite)

		if assigned_satellite == ship.Symbol {
			//println("[DEBUG] satellite market assignment found")
			assigned_market_waypoint = market_symbol
			break
		}

	}

	//println("[DEBUG] markets_to_cover:")

	//for k, v := range markets_to_cover {
	//	println(k)
	//	println(v)
	//}

	if IsShipAlreadyAtWaypoint(ship, assigned_market_waypoint) {
		fmt.Println("[INFO] Already at assigned market waypoint")
		if !IsShipDocked(ship) {
			DockShip(ship.Symbol)
		}
		fmt.Println("[INFO] Updating trade data...")
		UpdateTradeRoutesIncludingThisWaypoint(assigned_market_waypoint, trade_routes)
	} else {
		fmt.Println("[INFO] Not at assigned market waypoint, heading there now")
		if IsShipDocked(ship) {
			OrbitShip(ship.Symbol)
		}
		fmt.Println("[DEBUG] assigned_market_waypoint: " + assigned_market_waypoint)
		NavigateShip(ship.Symbol, assigned_market_waypoint)
	}

}
