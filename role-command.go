package main

import (
	"fmt"
	"math"
)

func BuyAsManyTradeGoodAsPossible(ship Ship, trade_good TradeGood, market Market) {
	maximum_affordable_units := HowManyTradeGoodCanIAfford(GetAgent(), trade_good)
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

	buy_market_trade_volume := trade_good.TradeVolume

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

			buy_cargo_result := PurchaseCargo(ship.Symbol, trade_good.Symbol, units_to_purchase)

			space_in_cargo_hold = buy_cargo_result.Cargo.Capacity - buy_cargo_result.Cargo.Units

			fmt.Println("[DEBUG] space_in_cargo_hold = ")
			fmt.Println(space_in_cargo_hold)

			if space_in_cargo_hold < units_to_purchase {
				units_to_purchase = space_in_cargo_hold
			}

		}
	} else {
		PurchaseCargo(ship.Symbol, trade_good.Symbol, units_to_purchase)
	}
}

func ApplyRoleCommand(ship Ship, all_waypoints_in_system []Waypoint, all_markets_in_system []Market, markets_to_cover map[string]string, trade_routes []TradeRoute, callsign string) {

	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role)

	//fmt.Println("[DEBUG] ApplyRoleCommand")

	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[DEBUG] IN_TRANSIT TO " + ship.Nav.Route.Destination.Symbol)
		fmt.Println("[DEBUG] Arrival " + ship.Nav.Route.Arrival)
		return
	}

	ship_list := ListShips()
	contracts := ListContracts()

	for _, contract := range contracts {
		fmt.Println(contract)
	}

	ship_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)

	// do we need to negotiate a contract?
	// do we have 0 contracts
	// is our current contract fulfilled

	if IsNewContractRequired(contracts) {
		if IsShipAtFactionWaypoint(ship, all_waypoints_in_system, ship.Registration.FactionSymbol) {
			if !IsShipDocked(ship){
				DockShip(ship.Symbol)
			}
			NegotiateContract(ship.Symbol)
		} else {
			// ship is not at faction waypoint
			// find closest faction waypoint
			clostest_faction_waypoint := ClosestFactionWaypointToShip(ship, all_waypoints_in_system)
			if IsShipDocked(ship) {
				RefuelShip(ship.Symbol)
				OrbitShip(ship.Symbol)
			}
			NavigateShip(ship.Symbol, clostest_faction_waypoint.Symbol)
			return
		}
	}

	contract := ActiveContract(contracts)
	contract_id := contract.ID
	contract_delivery_waypoint_symbol := contract.Terms.Deliver[0].DestinationSymbol
	contract_delivery_trade_good_symbol := contract.Terms.Deliver[0].TradeSymbol

	if !IsContractAccepted(contract) {
		fmt.Println("[INFO] Contract is negotiated but not accepted")
		fmt.Println("[INFO] Terms:")
		fmt.Println(contract)
		accept_contract_result := AcceptContract(contract.ID)
		fmt.Println(accept_contract_result)
	}

	fmt.Println("[INFO] Contract is accepted")
	if len(contract.Terms.Deliver) > 1 {
		fmt.Println("[ERROR] CONTRACT DELIVER OBJECT HAS MORE THAN ONE ELEMENT")
	}
	fmt.Print(contract.Type + " [")
	fmt.Print(contract.Terms.Deliver[0].UnitsFulfilled)
	fmt.Print("/")
	fmt.Print(contract.Terms.Deliver[0].UnitsRequired)
	fmt.Print("] ")
	fmt.Print(contract_delivery_trade_good_symbol + " to ")
	fmt.Print(contract_delivery_waypoint_symbol + " by ")
	fmt.Print(contract.Terms.Deadline + " for ")
	fmt.Println(contract.Terms.Payment.OnFulfilled)

	if CanContractBeCompleted(contract) {
		fulfill_contract_result := FulfillContract(contract_id)
		fmt.Println(fulfill_contract_result)
		// Again we waste a turn here, we should be navigating by now.
	} else {
		// Do we have any contract good in our hold?
		contract_good_in_hold := CountTradeGoodCargo(ship, contract_delivery_trade_good_symbol)
		if contract_good_in_hold > 0 {
			fmt.Println("[INFO] We have contract goods in our hold")
			if IsShipAlreadyAtWaypoint(ship, contract_delivery_waypoint_symbol) {
				fmt.Println("[INFO] Already at contract delivery waypoint")
				if !IsShipDocked(ship) {
					DockShip(ship.Symbol)
				}
				DeliverCargoToContract(contract_id, ship.Symbol, contract_delivery_trade_good_symbol, contract_good_in_hold)
				// This wastes a turn, it could Navigate immidiately after Delivering
			} else {
				// have contract goods in hold
				// not already at delivery waypoint
				if IsShipDocked(ship) {
					RefuelShip(ship.Symbol)
					OrbitShip(ship.Symbol)
				}
				fmt.Println("[INFO] Heading to contract destination.")
				NavigateShip(ship.Symbol, contract_delivery_waypoint_symbol)
				return
			}
		} else {
			// does not have any contract goods in hold
			// go pick up contract trade good
			// find closest out of a []Market
			markets_with_contract_trade_good := MarketplacesWhichSellTradeGood(all_markets_in_system, contract_delivery_trade_good_symbol)
			fmt.Println("[INFO] The following markets sell " + contract_delivery_trade_good_symbol + ":")
			for _, market_symbol := range markets_with_contract_trade_good {
				fmt.Println("[INFO] " + market_symbol.Symbol)
			}
			closest_market := ClosestMarketSellingTradeGood(ship, contract_delivery_trade_good_symbol, markets_with_contract_trade_good)
			closest_market_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, closest_market.Symbol)
			if IsShipAlreadyAtWaypoint(ship, closest_market.Symbol){
				if IsShipCargoEmpty(ship) {
					if !IsShipDocked(ship) {
						DockShip(ship.Symbol)
					}
					fmt.Println("[INFO] At market selling cargo contract")
					fmt.Println("[INFO] Purchase contract cargo.")
					market := GetMarket(base_system_symbol, closest_market.Symbol)
					trade_good := TradeGoodFromMarket(contract_delivery_trade_good_symbol, market)

					//TODO
					// Buy only as many as is required by contract
					required_units := ContractRemainingRequired(contract)

					if required_units < ship.Cargo.Capacity {
						PurchaseCargo(ship.Symbol, contract_delivery_trade_good_symbol, required_units)  
						// TODO account for trade volume limitations
					} else {
						BuyAsManyTradeGoodAsPossible(ship, trade_good, closest_market)
					}
				} else {
					if IsShipDocked(ship) {
						OrbitShip(ship.Symbol)
					}
					NavigateShip(ship.Symbol, contract_delivery_waypoint_symbol)
					return
				}
			} else {
				fmt.Println("[INFO] Heading to marketplace selling contract goods")
				if IsShipDocked(ship) {
					RefuelShip(ship.Symbol)
					OrbitShip(ship.Symbol)
				}
				// calculate path
				fmt.Println(ship.Nav.WaypointSymbol)
				fmt.Println(len(SystemGraph))
				fmt.Println(len(MarketplaceGraph))
				for _, graph_row := range MarketplaceGraph {
					fmt.Println(graph_row)
					for _, v := range graph_row {
						fmt.Println(v)
					}
				}
				path, _ := CalculateShortestPathBetweenTwoWaypoints(MarketplaceGraph, ship_waypoint, closest_market_waypoint)
				FollowPath(ship, path)				
			}
		}
	}


	// TESTING

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

			//TODO
			// Replace this with BuyAsManyTradeGoodAsPossible

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

