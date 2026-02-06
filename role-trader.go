package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
)

// TODO
// Allow for purchasing multiple times in the case the trade volume is low but we have sufficient credits.
// UpdateMarketData after performing ActionPurchaseCargo since it may alter prices and therefore trade_routes

var ShuttleMarketplaceGraph dijkstra.Graph = make(dijkstra.Graph)

func DecideTraderAction(ship Ship, all_waypoints_in_system []Waypoint) ShipAction {

	fmt.Println("[INFO] " + ship.Symbol + " DecideTraderAction")

	most_profitable_trade_route := MostProfitableTradeRoute(trade_routes)
	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)

	fmt.Print("[DEBUG] " + ship.Symbol + " Trade route: Buy ")
	fmt.Print(most_profitable_trade_route.TradeGoodSymbol)
	fmt.Print(" at ")
	fmt.Print(most_profitable_trade_route.BuyMarketplaceWaypointSymbol)
	fmt.Print(" for ")
	fmt.Print(most_profitable_trade_route.BuyMarketTradeGood.PurchasePrice)
	fmt.Print(" sell at ")
	fmt.Print(most_profitable_trade_route.SellMarketplaceWaypointSymbol)
	fmt.Print(" for ")
	fmt.Print(most_profitable_trade_route.SellMarketTradeGood.SellPrice)
	fmt.Print(" ppu ")
	fmt.Print(most_profitable_trade_route.ProfitPerUnit)
	fmt.Print(" pr ")
	fmt.Printf("%.2f", most_profitable_trade_route.ProfitabilityRating)
	fmt.Print( " distance ")
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

	max_affordable_units := HowManyTradeGoodCanIAfford(agent, most_profitable_trade_route.BuyMarketTradeGood)

	//if IsShipCargoEmpty(ship) {
	if !IsShipCargoFull(ship) && max_affordable_units > 0 {
		//fmt.Println("[DEBUG] Cargo hold empty")
		if IsShipAlreadyAtWaypoint(ship, most_profitable_trade_route.BuyMarketplaceWaypointSymbol) {
			if !IsShipDocked(ship) {
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
				space_in_cargo_hold := ship.Cargo.Capacity - ship.Cargo.Units
				units := space_in_cargo_hold
				if most_profitable_trade_route.BuyMarketTradeGood.TradeVolume < space_in_cargo_hold {
					units = most_profitable_trade_route.BuyMarketTradeGood.TradeVolume
				}
				if max_affordable_units < units {
					units = max_affordable_units
				}
				if units == 0 {
					return ShipAction{
						Type: ActionOrbit,
						ShipSymbol: ship.Symbol,
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