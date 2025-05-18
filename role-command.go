package main

import (
	"fmt"
	"math"
)

func ApplyRoleCommand(ship Ship, markets_to_cover map[string]string, trade_routes []TradeRoute, callsign string) {

	fmt.Println("[INFO] " + ship.Symbol)

	//fmt.Println("[DEBUG] ApplyRoleCommand")

	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[DEBUG] IN_TRANSIT TO " + ship.Nav.Route.Destination.Symbol)
		fmt.Println("[DEBUG] Arrival " + ship.Nav.Route.Arrival)
		return
	}

	ship_list := ListShips()

	// TESTING
	//PrintTradeRoutes(ship_list, trade_routes)
	//fmt.Println(len(markets_to_cover))
	//os.Exit(1)
	return
	// TESTING

	number_of_satellites := CountShipsByFrame(ship_list, "SATELLITE")
	number_of_markets_to_cover := len(markets_to_cover)

	if number_of_satellites >= number_of_markets_to_cover {
		fmt.Println("[INFO] Enough satellites")
		if !SatelliteToMarketAssignmentComplete(markets_to_cover) {
			AssignSatellitesToMarkets(markets_to_cover)
		}
	} else {
		fmt.Println("[INFO] Not enough satellites, postponing market assignments...")
	}

	// we have enough satellites
	//fmt.Println("[INFO] We have enough satellites, boss. It's time to start trading!")

	if !SatelliteToMarketAssignmentComplete(markets_to_cover) {
		AssignSatellitesToMarkets(markets_to_cover)
	}

	if MarketScanComplete(trade_routes) {
		PopulateTradeRoutesProfitPerUnit(trade_routes)
		WriteTradeRoutesToFile(trade_routes, callsign)
	}

	PrintTradeRoutes(ship_list, trade_routes)

	most_profitable_trade_route := MostProfitableTradeRoute(trade_routes)

	if IsWaypointWithinDistanceOfWaypoint(most_profitable_trade_route.BuyWaypoint, most_profitable_trade_route.SellWaypoint, int(ship.Frame.FuelCapacity)) {
		//fmt.Println("[DEBUG] IsWaypointWithinDistanceOfWaypoint true")
	} else {
		//fmt.Println("[DEBUG] IsWaypointWithinDistanceOfWaypoint false")
	}

	if most_profitable_trade_route.ProfitabilityRating < 0 {
		fmt.Println("[INFO] Most profitable trade route is not profitable enough. Doing nothing...")
		return
	}

	if IsShipCargoEmpty(ship) {
		fmt.Println("[INFO] Cargo hold empty")
		if IsShipAlreadyAtWaypoint(ship, most_profitable_trade_route.BuyMarketplaceWaypointSymbol) {
			fmt.Println("[DEBUG] Already at waypoint")
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}

			// BUY STUFF

			maximum_affordable_units := HowManyTradeGoodCanIAfford(GetAgent(), most_profitable_trade_route.BuyMarketTradeGood)
			fmt.Print("[DEBUG] maximum_affordable_units = ")
			fmt.Println(maximum_affordable_units)

			units_to_purchase := maximum_affordable_units
			fmt.Println("[DEBUG] units_to_purchase = ")
			fmt.Println(units_to_purchase)

			space_in_cargo_hold := ship.Cargo.Capacity - ship.Cargo.Units
			fmt.Print("[DEBUG] space_in_cargo_hold = ")
			fmt.Println(space_in_cargo_hold)

			if space_in_cargo_hold < units_to_purchase {
				units_to_purchase = space_in_cargo_hold
				fmt.Print("[DEBUG] units_to_purchase = ")
				fmt.Println(space_in_cargo_hold)
			}

			buy_market_trade_volume := most_profitable_trade_route.BuyMarketTradeGood.TradeVolume

			fmt.Print("[DEBUG] buy_market_trade_volume = ")
			fmt.Println(buy_market_trade_volume)

			if space_in_cargo_hold > buy_market_trade_volume {
				fmt.Println("[DEBUG] space_in_cargo_hold > buy_market_trade_volume")

				number_of_purchases_required := float64(space_in_cargo_hold) / float64(buy_market_trade_volume)
				rounded_number_of_purchases_required := math.Ceil(number_of_purchases_required)
				fmt.Println("[DEBUG] rounded_number_of_purchases_required = ")
				fmt.Println(rounded_number_of_purchases_required)

				units_to_purchase = buy_market_trade_volume
				fmt.Println("[DEBUG] units_to_purchase = ")
				fmt.Println(units_to_purchase)

				for i := 0; float64(i) < rounded_number_of_purchases_required; i++ {
					fmt.Println("[DEBUG] units_to_purchase = ")
					fmt.Println(units_to_purchase)

					buy_cargo_result := PurchaseCargo(ship.Symbol, most_profitable_trade_route.TradeGoodSymbol, units_to_purchase)

					space_in_cargo_hold = buy_cargo_result.Cargo.Capacity - buy_cargo_result.Cargo.Units

					fmt.Println("[DEBUG] space_in_cargo_hold = ")
					fmt.Println(space_in_cargo_hold)

					if space_in_cargo_hold < units_to_purchase {
						units_to_purchase = space_in_cargo_hold
					}

				}
			} else {
				PurchaseCargo(ship.Symbol, most_profitable_trade_route.TradeGoodSymbol, units_to_purchase)
			}

			fmt.Print("[DEBUG] units_to_purchase = ")
			fmt.Print(units_to_purchase)
			fmt.Println()
			RefuelShip(ship.Symbol)
			OrbitShip(ship.Symbol)
			fmt.Println("[INFO] " + ship.Symbol + " Heading to SellMarketplaceWaypointSymbol")
			NavigateShip(ship.Symbol, most_profitable_trade_route.SellMarketplaceWaypointSymbol)
			return
		}

		if MarketScanComplete(trade_routes) {
			fmt.Println("[INFO] Heading to buy marketplace")
			if IsShipDocked(ship) {
				RefuelShip(ship.Symbol)
				OrbitShip(ship.Symbol)
			}
			NavigateShip(ship.Symbol, most_profitable_trade_route.BuyMarketplaceWaypointSymbol)
			return
		}
	} else {
		fmt.Println("[INFO] Cargo not empty, we have " + ship.Cargo.Inventory[0].Symbol)

		if !MarketScanComplete(trade_routes) {
			fmt.Println("[DEBUG] wait for market data")
			return
		}

		first_item_in_inventory := ship.Cargo.Inventory[0]

		trade_routes_with_inventory_good := TradeRoutesWithTradeGood(trade_routes, first_item_in_inventory.Symbol)

		fmt.Print("[DEBUG] first_item_in_inventory.Symbol ")
		fmt.Println(first_item_in_inventory.Symbol)

		fmt.Print("[DEBUG] trade_routes_with_inventory_good length ")
		fmt.Println(len(trade_routes_with_inventory_good))

		if !MarketScanComplete(trade_routes) {
			println("[DEBUG] market data incomplete, returning to avoid running MostProfitableTradeRoute")
			return
		}

		most_profitable_trade_route_with_inventory_good := MostProfitableTradeRoute(trade_routes_with_inventory_good)

		fmt.Println("[DEBUG] most_profitable_trade_route_with_inventory_good")
		//fmt.Println(most_profitable_trade_route_with_inventory_good)

		if IsShipAlreadyAtWaypoint(ship, most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol) {
			fmt.Println("[DEBUG] Already at sell marketplace")
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}

			// this doesnt account for TradeVolume < Cargo.Capacity
			units_in_cargo_hold := CountTradeGoodCargo(ship, most_profitable_trade_route_with_inventory_good.TradeGoodSymbol)
			sell_market_trade_volume := most_profitable_trade_route_with_inventory_good.SellMarketTradeGood.TradeVolume
			if units_in_cargo_hold > sell_market_trade_volume {
				number_of_sales_required := units_in_cargo_hold / sell_market_trade_volume

				units_to_sell := sell_market_trade_volume
				for i := 0; i < int(number_of_sales_required); i++ {
					sell_cargo_result := SellCargo(ship.Symbol, most_profitable_trade_route_with_inventory_good.TradeGoodSymbol, units_to_sell)
					if sell_cargo_result.Transaction.Units < sell_market_trade_volume {
						units_to_sell = sell_market_trade_volume
					}
				}
			} else {
				SellCargo(ship.Symbol, most_profitable_trade_route_with_inventory_good.TradeGoodSymbol, units_in_cargo_hold)
			}
			RefuelShip(ship.Symbol)
			OrbitShip(ship.Symbol)
			NavigateShip(ship.Symbol, most_profitable_trade_route.BuyMarketplaceWaypointSymbol)
		} else {
			fmt.Println("[DEBUG] Not yet at SellMarketplaceWaypointSymbol")
			// this is nasty, this whole function is now nasty. it needs to be chopped up into bitesize chunks
			if !MarketScanComplete(trade_routes) {
				fmt.Println("[DEBUG] Market scan not yet complete. Waiting for data")
				return
			}
			if IsShipDocked(ship) {
				RefuelShip(ship.Symbol)
				OrbitShip(ship.Symbol)
			}
			println("most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol")
			println(most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol)

			// dirty
			if most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol == "" {
				println("[ERROR] most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol null. Would have navigated!")
				return
			}

			NavigateShip(ship.Symbol, most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol)
		}
	}
}
