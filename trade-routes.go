package main

import (
	"fmt"
	"sort"
)

func HowManyTradeGoodCanIAfford(agent Agent, trade_good TradeGood) int64 {
	// explicitly avoid /0 errors
	if trade_good.PurchasePrice == 0 {
		return 0
	}
	max_buy_count := agent.Credits / trade_good.PurchasePrice
	return max_buy_count
}

// TODO Remove
func CalculateProfitPerUnit(trade_route TradeRoute) float64 {
	//fmt.Println("[DEBUG] CalculateProfitPerUnit")
	if trade_route.BuyMarketTradeGood.PurchasePrice == 0 {
		fmt.Println("[WARNING] CalculateProfitPerUnit trade_route.BuyMarketTradeGood.PurchasePrice == 0")
		return 0
	}
	if trade_route.SellMarketTradeGood.SellPrice == 0 {
		fmt.Println("[WARNING] CalculateProfitPerUnit trade_route.SellMarketTradeGoodPurchasePrice == 0")
		return 0
	}
	sell_price := float64(trade_route.SellMarketTradeGood.SellPrice)
	buy_price := float64(trade_route.BuyMarketTradeGood.PurchasePrice)
	profit_per_unit := sell_price - buy_price
	return profit_per_unit
}

func MostProfitableTradeRoute(trade_routes []TradeRoute) TradeRoute {
	most_profitable_trade_route := TradeRoute{}
	best_profitability_score := -99.00
	for _, trade_route := range trade_routes {
		//fmt.Printf("%.2f", trade_route.ProfitabilityRating)
		if trade_route.ProfitabilityRating > best_profitability_score {
			most_profitable_trade_route = trade_route
			best_profitability_score = trade_route.ProfitabilityRating
		}
	}

	return most_profitable_trade_route
}

// TODO repurpose this for derived trade routes
func ShortestTradeRoute(trade_routes []TradeRoute) TradeRoute {
	shortest_trade_route := TradeRoute{}
	best_distance := 9999
	for _, trade_route := range trade_routes {
		if trade_route.Distance < best_distance {
			shortest_trade_route = trade_route
			best_distance = trade_route.Distance
		}
	}
	return shortest_trade_route
}

func TradeRoutesWithTradeGood(trade_routes []TradeRoute, trade_good_symbol string) []TradeRoute {
	var trade_routes_with_trade_good = []TradeRoute{}
	for _, trade_route := range trade_routes {
		if trade_route.TradeGoodSymbol == trade_good_symbol {
			trade_routes_with_trade_good = append(trade_routes_with_trade_good, trade_route)
		}
	}
	return trade_routes_with_trade_good
}

func SatelliteToMarketAssignmentComplete(markets_to_cover map[string]string) bool {
	for _, v := range markets_to_cover {

		if len(v) == 0 {
			//fmt.Println("[DEBUG] SatelliteToMarketAssignment Incomplete")
			return false
		}
	}
	//fmt.Println("[DEBUG] SatelliteToMarketAssignment Complete")
	return true
}

func MarketScanComplete(trade_routes []TradeRoute) bool {

	if len(trade_routes) == 0 {
		fmt.Println("[WARN] NO TRADE ROUTES")
		return false
	}

	for _, trade_route := range trade_routes {
		if trade_route.BuyMarketTradeGood.PurchasePrice == 0 {
			fmt.Println("[INFO] Market data incomplete, wait for input...")
			return false
		}
	}
	fmt.Println("[INFO] Market data complete, let's trade...")
	return true
}

// TODO repurpose this for derived trade routes
func PopulateTradeRoutesWithDistances(trade_routes []TradeRoute) {
	for i, trade_route := range trade_routes {
		distance := DistanceBetweenTwoWaypoints(trade_route.BuyWaypoint, trade_route.SellWaypoint)
		trade_routes[i].Distance = distance
	}
}

func AssignSatellitesToMarkets(markets_to_cover map[string]string) {

	fmt.Println("[DEBUG] AssignSatellitesToMarkets")
	list_of_ships := ListShips()
	list_of_satellites := []Ship{}

	for _, ship := range list_of_ships {
		if ship.Registration.Role == "SATELLITE" {
			list_of_satellites = append(list_of_satellites, ship)
		}
	}

	keys := make([]string, 0, len(markets_to_cover))

	for k := range markets_to_cover {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var market_index int
	for _, market_waypoint := range keys {
		// if this market does not have a satellite assigned:
		if markets_to_cover[market_waypoint] == "" {
			markets_to_cover[market_waypoint] = list_of_satellites[market_index].Symbol
			fmt.Println("[INFO] Assigned satellite " + list_of_satellites[market_index].Symbol + " to market " + market_waypoint)
		}
		market_index++
	}
}

type DerivedRoute struct {
	From                string
	To                  string
	Good                string
	ProfitabilityRating float64
}

func DeriveTradeRoutes(world *WorldState) []DerivedRoute {
	routes := []DerivedRoute{}

	for _, fromMarket := range world.Markets {
		for _, toMarket := range world.Markets {

			if fromMarket.Market.Symbol == toMarket.Market.Symbol {
				continue
			}

			for _, fromGood := range fromMarket.Market.TradeGoods {
				success, toGood := TradeGoodFromMarket(fromGood.Symbol, toMarket.Market)

				if !success {
					continue
				}

				if fromGood.PurchasePrice == 0 || toGood.SellPrice == 0 {
					continue
				}

				profit := toGood.SellPrice - fromGood.PurchasePrice
				if profit <= 0 {
					continue
				}

				// TODO this should be a path cost from a graph by ship type rather than as-the-crow-flies
				distance := DistanceBetweenTwoWaypoints(*world.Waypoints[fromMarket.Market.Symbol], *world.Waypoints[toMarket.Market.Symbol])
				if distance == 0 {
					distance = 1
				}
				profitability_rating := float64(profit) / float64(distance*2)

				routes = append(routes, DerivedRoute{
					From:                fromMarket.Market.Symbol,
					To:                  toMarket.Market.Symbol,
					Good:                fromGood.Symbol,
					ProfitabilityRating: profitability_rating,
				})
			}
		}
	}

	return routes
}

func MostProfitableDerivedTradeRoute(derived_routes []DerivedRoute) DerivedRoute {
	highest_pr := float64(-9999)
	highest_pr_route := DerivedRoute{}
	for _, route := range derived_routes {
		if route.ProfitabilityRating > highest_pr {
			highest_pr_route = route
			highest_pr = route.ProfitabilityRating
		}
	}
	return highest_pr_route
}

func DerivedTradeRoutesWithTradeGood(derived_trade_routes []DerivedRoute, trade_good_symbol string) []DerivedRoute {
	var derived_trade_routes_with_trade_good = []DerivedRoute{}
	for _, trade_route := range derived_trade_routes {
		if trade_route.Good == trade_good_symbol {
			derived_trade_routes_with_trade_good = append(derived_trade_routes_with_trade_good, trade_route)
		}
	}
	return derived_trade_routes_with_trade_good
}

func MarketsYetToBeVisited(world *WorldState) []Market {
	markets_yet_to_be_vistied := []Market{}
	for _, market_state := range world.Markets {
		if market_state.Market.TradeGoods == nil {
			markets_yet_to_be_vistied = append(markets_yet_to_be_vistied, market_state.Market)
		}
	}
	return markets_yet_to_be_vistied
}

func BestMarketToSellGood(world *WorldState, tradeGoodSymbol string) (Market, int64, bool) {

	var bestMarket Market
	bestPrice := int64(-1000)
	found := false

	for _, market_state := range world.Markets {
		for _, good := range market_state.Market.TradeGoods {
			if good.Symbol != tradeGoodSymbol {
				continue
			}

			// If the market isn't buying it, skip
			if good.SellPrice <= int64(0) {
				continue
			}

			if good.SellPrice > bestPrice {
				bestPrice = good.SellPrice
				bestMarket = market_state.Market
				found = true
			}
		}
	}

	return bestMarket, bestPrice, found
}
