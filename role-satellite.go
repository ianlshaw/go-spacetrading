package main

import (
	"fmt"
)


func DecideSatelliteAction(
	ship Ship,
	
) ShipAction {
	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + "DecideSatelliteAction")
	
	//if IsShipInTransit(ship){
	//	fmt.Println("[INFO] IN_TRANSIT TO " + ship.Nav.Route.Destination.Symbol)
	//	fmt.Println("[INFO] Arrival " + ship.Nav.Route.Arrival)
	//	fmt.Println("[ERROR] ApplyRoleSatellite SHIP IN TRANSIT THIS SHOULD NOT HAPPEN")
	//	return ShipAction{
	//		Type: ActionWait,
	//		ShipSymbol: ship.Symbol,
	//		NotBefore: StringToTimestamp(ship.Nav.Route.Arrival),
	//	}
	//}

	// This doesn't really work since we're not setting profitability rating on a per trade route basis
	trade_route_missing_data_success, trade_route := TradeRouteMissingData()
	if !trade_route_missing_data_success {
		fmt.Println("[WARN] " + ship.Symbol + " All trade routes have data. Sleeping for 3 hours")
		return ShipAction{
			Type: ActionWait,
			ShipSymbol: ship.Symbol,
			NotBefore: ThreeHoursFromNow(),
		}
	}

	// buy marketplace
	if trade_route.BuyMarketTradeGood.PurchasePrice == 0 {
		if IsShipAlreadyAtWaypoint(ship, trade_route.BuyMarketplaceWaypointSymbol) {
			if IsShipDocked(ship){
				return ShipAction{
					Type: ActionUpdateMarketData,
					ShipSymbol: ship.Symbol,
					WaypointSymbol: trade_route.BuyMarketplaceWaypointSymbol,
				}
			} else {
				// ship is at buy marketplace but not docked
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			}
		} else {
			if IsShipDocked(ship){
				return ShipAction{
					Type: ActionOrbit,
					ShipSymbol: ship.Symbol,
				}
			} else {
				return ShipAction{
					Type: ActionNavigate,
					ShipSymbol: ship.Symbol,
					WaypointSymbol: trade_route.BuyMarketplaceWaypointSymbol,
				}
			}
		}
	}

	// sell marketplace
	if trade_route.SellMarketTradeGood.SellPrice == 0 {
		if IsShipAlreadyAtWaypoint(ship, trade_route.SellMarketplaceWaypointSymbol){
			if IsShipDocked(ship){
				return ShipAction{
					Type: ActionUpdateMarketData,
					ShipSymbol: ship.Symbol,
					WaypointSymbol: trade_route.SellMarketplaceWaypointSymbol,
				}
			} else {
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			}

		} else {
			if IsShipDocked(ship) {
				return ShipAction{
					Type: ActionOrbit,
					ShipSymbol: ship.Symbol,
				}
			} else {
				return ShipAction{
					Type: ActionNavigate,
					ShipSymbol: ship.Symbol,
					WaypointSymbol: trade_route.SellMarketplaceWaypointSymbol,
				}
			}
		}
	}
	
	fmt.Println(" uncaught branch, returning 15 minute delay")
	return ShipAction{
		Type: ActionWait,
		ShipSymbol: ship.Symbol,
		NotBefore: FifteenMinutesFromNow(),
	}
}