package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
)

// TODO
// Allow for purchasing multiple times in the case the trade volume is low but we have sufficient credits.

var ShuttleMarketplaceGraph dijkstra.Graph = make(dijkstra.Graph)

func DecideTraderAction(ship Ship, all_waypoints_in_system []Waypoint) ShipAction {

	fmt.Println("[INFO] " + ship.Symbol + " ApplyRoleTrader")

	most_profitable_trade_route := MostProfitableTradeRoute(trade_routes)
	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)

	fmt.Println(most_profitable_trade_route.TradeGoodSymbol)
	fmt.Println(most_profitable_trade_route.BuyMarketplaceWaypointSymbol)
	fmt.Println(most_profitable_trade_route.BuyMarketTradeGood.PurchasePrice)
	fmt.Println(most_profitable_trade_route.SellMarketplaceWaypointSymbol)
	fmt.Println(most_profitable_trade_route.SellMarketTradeGood.SellPrice)
	fmt.Println(most_profitable_trade_route.ProfitPerUnit)
	fmt.Println(most_profitable_trade_route.ProfitabilityRating)
	fmt.Println(most_profitable_trade_route.Distance)

	if most_profitable_trade_route.ProfitabilityRating < 1 {
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

	// This is to account for most_profitable_trade_route changing after a purchase. We still want to sell our current cargo.
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
				if World.IsMarketStale(ship.Nav.WaypointSymbol) {
					return ShipAction{
						Type: ActionUpdateMarketData,
						ShipSymbol: ship.Symbol,
						WaypointSymbol: ship.Nav.WaypointSymbol,
					}
				}
				// at sell wp, not empty, docked.
				trade_good_cargo_count := CountTradeGoodCargo(ship, most_profitable_trade_route.TradeGoodSymbol)
				units := trade_good_cargo_count
				trade_volume := most_profitable_trade_route.SellMarketTradeGood.TradeVolume
				if trade_volume < trade_good_cargo_count {
					units = trade_volume
				}
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
				if World.IsMarketStale(ship.Nav.WaypointSymbol) {
					return ShipAction{
						Type: ActionUpdateMarketData,
						ShipSymbol: ship.Symbol,
						WaypointSymbol: ship.Nav.WaypointSymbol,
					}
				}
				units := space_in_cargo_hold
				if most_profitable_trade_route.BuyMarketTradeGood.TradeVolume < space_in_cargo_hold {
					units = most_profitable_trade_route.BuyMarketTradeGood.TradeVolume
				}
				max_affordable_units := HowManyTradeGoodCanIAfford(agent, most_profitable_trade_route.BuyMarketTradeGood)
				if max_affordable_units < units {
					units = max_affordable_units
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