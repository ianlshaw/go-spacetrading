package main

import (
	"fmt"
)

func DecideSatelliteAction(
	ship Ship,
	trade_routes []TradeRoute,
) ShipAction {
	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + "DecideSatelliteAction")
	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[INFO] IN_TRANSIT TO " + ship.Nav.Route.Destination.Symbol)
		fmt.Println("[INFO] Arrival " + ship.Nav.Route.Arrival)
		fmt.Println("[ERROR] ApplyRoleSatellite SHIP IN TRANSIT THIS SHOULD NOT HAPPEN")
		return ShipAction{
			Type: ActionWait,
			ShipSymbol: ship.Symbol,
			NotBefore: StringToTimestamp(ship.Nav.Route.Arrival),
		}
	}

	trade_route_missing_data_success, trade_route := TradeRouteMissingData(trade_routes)
	if !trade_route_missing_data_success {
		fmt.Println("[WARN] " + ship.Symbol + " All trade routes have data. Sleeping for 3 hours")
		return ShipAction{
			Type: ActionWait,
			ShipSymbol: ship.Symbol,
			NotBefore: ThreeHoursFromNow(),
		}
	}

	if trade_route.BuyMarketTradeGood.PurchasePrice == 0 {
		if IsShipAlreadyAtWaypoint(ship, trade_route.BuyMarketplaceWaypointSymbol){
			if IsShipDocked(ship){
				return ShipAction{
					Type: ActionUpdateMarketData,
					ShipSymbol: ship.Symbol,
					WaypointSymbol: trade_route.BuyMarketplaceWaypointSymbol,
					TradeRoutes: trade_routes,
				}
			} else {
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			}
		} else if IsShipDocked(ship){
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

		} else if IsShipDocked(ship) {
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
	
	fmt.Println(" uncaught branch, returning 15 minute delay")
	return ShipAction{
		Type: ActionWait,
		ShipSymbol: ship.Symbol,
		NotBefore: FifteenMinutesFromNow(),
	}
}