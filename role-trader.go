package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
)

// TODO
// Account for most_profitable_trade_route changing while we have cargo on board.
// Add an UpdateMarketData action when trader is at a market

var ShuttleMarketplaceGraph dijkstra.Graph = make(dijkstra.Graph)

func DecideTraderAction(ship Ship, all_waypoints_in_system []Waypoint) ShipAction {

	fmt.Println("[INFO] " + ship.Symbol + " ApplyRoleTrader")

	most_profitable_trade_route := MostProfitableTradeRoute(trade_routes)
	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)

	if most_profitable_trade_route.ProfitabilityRating < 0 {
		fmt.Println("[INFO] Most profitable trade route is not profitable enough. Doing nothing...")
		return ShipAction{
			Type: ActionWait,
			NotBefore: FifteenMinutesFromNow(),
		}
	}

	if !IsFuelFull(ship) {
		if IsShipDocked(ship){
			return ShipAction{
				Type: ActionRefuel,
				ShipSymbol: ship.Symbol,
			}
		} else {
			return ShipAction{
				Type: ActionDock,
				ShipSymbol: ship.Symbol,
			}
		}
	}

	if !IsShipCargoEmpty(ship){
		trade_good_in_cargo := ship.Cargo.Inventory[0].Symbol
		trade_routes_with_cargo := TradeRoutesWithTradeGood(trade_routes, trade_good_in_cargo)
		most_profitable_trade_route = MostProfitableTradeRoute(trade_routes_with_cargo)
	}

	space_in_cargo_hold := ship.Cargo.Capacity - ship.Cargo.Units

	if IsShipAlreadyAtWaypoint(ship, most_profitable_trade_route.SellMarketplaceWaypointSymbol) {
		if !IsShipCargoEmpty(ship) {
			if !IsShipDocked(ship) {
				// maybe this never triggers since the refuel logic above would always dock
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			} else {
				// at sell wp, not empty, docked.

				// calculate units
				trade_good_cargo_count := CountTradeGoodCargo(ship, most_profitable_trade_route.TradeGoodSymbol)
				units := trade_good_cargo_count
				trade_volume := most_profitable_trade_route.SellMarketTradeGood.TradeVolume
				if trade_volume < trade_good_cargo_count {
					units = trade_volume
				}
				// keep low trade volume in mind
				return ShipAction{
					Type: ActionSellCargo,
					ShipSymbol: ship.Symbol,
					TradeGoodSymbol: most_profitable_trade_route.TradeGoodSymbol,
					Units: units,
				}
			}
		}
	}
	
	
	if IsShipCargoEmpty(ship) {
		fmt.Println("[INFO] Cargo hold empty")
		if IsShipAlreadyAtWaypoint(ship, most_profitable_trade_route.BuyMarketplaceWaypointSymbol) {
			fmt.Println("[DEBUG] Already at waypoint")
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			} else {
				units := space_in_cargo_hold
				if most_profitable_trade_route.BuyMarketTradeGood.TradeVolume < space_in_cargo_hold {
					units = most_profitable_trade_route.BuyMarketTradeGood.TradeVolume
				}

				if World.IsMarketStale(ship.Nav.WaypointSymbol) {
					return ShipAction{
						Type: ActionUpdateMarketData,
						ShipSymbol: ship.Symbol,
						WaypointSymbol: ship.Nav.WaypointSymbol,
					}
				}
				return ShipAction{
					Type: ActionPurchaseCargo,
					ShipSymbol: ship.Symbol,
					TradeGoodSymbol: most_profitable_trade_route.BuyMarketTradeGood.Symbol,
					Units: units,
				}
			}
		}
		if IsShipDocked(ship){
			return ShipAction{
				Type: ActionOrbit,
				ShipSymbol: ship.Symbol,
			}
		}

		
		most_profitable_trade_route_buy_marketplace_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, most_profitable_trade_route.BuyMarketplaceWaypointSymbol)
		path, _ := CalculateShortestPathBetweenTwoWaypoints(ShuttleMarketplaceGraph, current_waypoint, most_profitable_trade_route_buy_marketplace_waypoint)
		// DEBUG
		fmt.Println(path)
		return ShipAction{
			Type: ActionFollowPath,
			ShipSymbol: ship.Symbol,
			Path: path,
		}
		
	}
	if IsShipDocked(ship) {
		return ShipAction{
			Type: ActionOrbit,
			ShipSymbol: ship.Symbol,
		}
	}

	most_profitable_trade_route_sell_marketplace_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, most_profitable_trade_route.SellMarketplaceWaypointSymbol)
	path, _ := CalculateShortestPathBetweenTwoWaypoints(ShuttleMarketplaceGraph, current_waypoint, most_profitable_trade_route_sell_marketplace_waypoint)
	// DEBUG
	fmt.Println(path)
	return ShipAction{
		Type: ActionFollowPath,
		ShipSymbol: ship.Symbol,
		Path: path,
	}

	fmt.Println("[ERROR] DecideTraderAction unhandled branch")
	return ShipAction{
		Type: ActionWait,
		NotBefore: FifteenMinutesFromNow(),
	}
}