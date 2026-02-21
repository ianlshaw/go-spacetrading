package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
)

// TODO

var ShuttleMarketplaceGraph dijkstra.Graph = make(dijkstra.Graph)

func DecideTraderAction(ship_ptr *Ship, world *WorldState) ShipAction {

	ship := *ship_ptr

	fmt.Println("[INFO] " + ship.Symbol + " DecideTraderAction")

	derived_trade_routes := DeriveTradeRoutes(world)

	//fmt.Printf("[DEBUG] %d derived trade routes\n", len(derived_trade_routes))

	//for _, route := range derived_trade_routes {
	//	fmt.Println(route)
	//}

	most_profitable_derived_trade_route := MostProfitableDerivedTradeRoute(derived_trade_routes)

	fmt.Printf("[DEBUG] %s Trade route: Buy %s at %s sell at %s pr %.2f\n",
	ship.Symbol,
	most_profitable_derived_trade_route.Good,
	most_profitable_derived_trade_route.From,
	most_profitable_derived_trade_route.To,
	most_profitable_derived_trade_route.ProfitabilityRating)

	//most_profitable_trade_route := MostProfitableTradeRoute(trade_routes)
	current_waypoint := *world.Waypoints[ship.Nav.WaypointSymbol]

	// Always update market data for a market we're at if it needs it.
	if IsShipDocked(ship){
		if world.IsMarketStale(ship.Nav.WaypointSymbol) {
			return ShipAction{
				Type: ActionUpdateMarketData,
				ShipSymbol: ship.Symbol,
				WaypointSymbol: ship.Nav.WaypointSymbol,
			}
		}
	}

	// Always refuel if fuel is not full
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

	// This can get stuck when the trader is holding cargo which becomes unprofitable after a trade.
	// Further trades of that cargo will not be worth it and the trader will wait until it becomes profitable
	if most_profitable_derived_trade_route.ProfitabilityRating <= 0 {
		fmt.Println("[INFO] Most profitable trade route is not profitable enough. Doing nothing...")
		return ShipAction{
			Type: ActionWait,
			NotBefore: ThreeMinutesFromNow(),
		}
	}

	var sell_market_symbol string
	var trade_good_in_cargo string

	// large trades with low trade volume can cause the derived route to cease to exist. eventually causing a nil pointer from here.
	if !IsShipCargoEmpty(ship){
		trade_good_in_cargo = ship.Cargo.Inventory[0].Symbol
		fmt.Printf("[INFO] trade good in cargo: %s\n", trade_good_in_cargo)
		trade_routes_with_cargo := DerivedTradeRoutesWithTradeGood(derived_trade_routes, trade_good_in_cargo)
		fmt.Printf("[INFO] %d trade routes with carried cargo %s\n", len(trade_routes_with_cargo), trade_good_in_cargo)
		if (len(trade_routes_with_cargo)) == 0 {
			fmt.Printf("[WARN] No trade route exists for held cargo %s\n", trade_good_in_cargo)

			// There are three options here
			// 1) Jettison the remaining cargo (wasteful)
			// 2) Find a different market which will take the remaining trade goods (complex) (current)
			// 3) Wait until the current market will take the remaining trade goods (may get stuck)

			sell_market, _, success := BestMarketToSellGood(world, trade_good_in_cargo)
			if !success {
				fmt.Println("[ERROR] No markets will accept %s\n", trade_good_in_cargo)
				return ShipAction{
					Type: ActionWait,
					NotBefore: ThreeMinutesFromNow(),
				}
			}
			//fmt.Printf("[DEBUG] Backup market found for %s\n", trade_good_in_cargo)
			sell_market_symbol = sell_market.Symbol
		
		} else {
			most_profitable_derived_trade_route = MostProfitableDerivedTradeRoute(trade_routes_with_cargo)
			sell_market_symbol = most_profitable_derived_trade_route.To
		}

	}

	if IsShipAlreadyAtWaypoint(ship, sell_market_symbol) {
		if !IsShipCargoEmpty(ship) {
			if !IsShipDocked(ship) {
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			} else {
				// at sell wp, not empty, docked.
				trade_good_cargo_count := CountTradeGoodCargo(ship, trade_good_in_cargo)
				units := trade_good_cargo_count
				
				sell_market := world.Markets[sell_market_symbol].Market

				success, sell_market_trade_good := TradeGoodFromMarket(trade_good_in_cargo, sell_market)
				if !success {
					fmt.Println("[ERROR] DecideTraderAction TradeGoodFromMarket failed")
				}
				trade_volume := sell_market_trade_good.TradeVolume
				if trade_volume < trade_good_cargo_count {
					units = trade_volume
				}
				return ShipAction{
					Type: ActionSellCargo,
					ShipSymbol: ship.Symbol,
					TradeGoodSymbol: trade_good_in_cargo,
					Units: units,
				}
			}
		}
	}

	//fmt.Println(most_profitable_derived_trade_route)
	buy_market := world.Markets[most_profitable_derived_trade_route.From].Market
	success, buy_market_trade_good := TradeGoodFromMarket(most_profitable_derived_trade_route.Good, buy_market)
	if !success {
		fmt.Println("[ERROR] DecideTraderAction TradeGoodFromMarket failed")
	}
	//fmt.Println(buy_market_trade_good)
	max_affordable_units := HowManyTradeGoodCanIAfford(agent, buy_market_trade_good)

	//if IsShipCargoEmpty(ship) {
	if !IsShipCargoFull(ship) && max_affordable_units > 0 {
		//fmt.Println("[DEBUG] Cargo hold empty")
		if IsShipAlreadyAtWaypoint(ship, most_profitable_derived_trade_route.From) {
			if !IsShipDocked(ship) {
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			} else {
				space_in_cargo_hold := ship.Cargo.Capacity - ship.Cargo.Units
				units := space_in_cargo_hold
				if buy_market_trade_good.TradeVolume < space_in_cargo_hold {
					units = buy_market_trade_good.TradeVolume
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
					TradeGoodSymbol: most_profitable_derived_trade_route.Good,
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
		
		most_profitable_trade_route_buy_marketplace_waypoint := *world.Waypoints[most_profitable_derived_trade_route.From]
		path, _, err := CalculateShortestPathBetweenTwoWaypoints(ShuttleMarketplaceGraph, current_waypoint, most_profitable_trade_route_buy_marketplace_waypoint)

		if err != nil {
			fmt.Println(err)
			fmt.Println("[ERROR] cannot path")
			return ShipAction{
				Type: ActionWait,
				NotBefore: FifteenMinutesFromNow(),
			}
		}
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

	most_profitable_trade_route_sell_marketplace_waypoint := *world.Waypoints[sell_market_symbol]
	path, _, err := CalculateShortestPathBetweenTwoWaypoints(ShuttleMarketplaceGraph, current_waypoint, most_profitable_trade_route_sell_marketplace_waypoint)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] cannot path")
		return ShipAction{
			Type: ActionWait,
			NotBefore: FifteenMinutesFromNow(),
		}
	}
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