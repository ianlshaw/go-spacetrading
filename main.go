package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var url_base string = "https://api.spacetraders.io/v2/"
var account_token = "Bearer "
var account_token_filename = "ACCOUNTTOKENDONOTEXPOSE.txt"
var agent_token = "Bearer "
var base_system_symbol = ""
var http_calls = 0
var turn_length = 120

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func pretty_print_json(json_blob string) {
	byt := []byte(json_blob)

	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(dat, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}

func populate_base_system_symbol() {
	//fmt.Println("[DEBUG] populate_base_system_symbol")
	endpoint := "my/ships"
	response_string := BasicGet(endpoint)

	response_typed := ListShipsResponseData{}
	if err := json.Unmarshal([]byte(response_string), &response_typed); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	base_system_symbol = response_typed.Data[0].Nav.SystemSymbol
}

func ShipRoleDecider(ship Ship, markets_to_cover map[string]string, probe_shipyards []Waypoint, trade_routes []TradeRoute, callsign string) {
	if ship.Registration.Role == "COMMAND" {
		ApplyRoleCommand(ship, markets_to_cover, probe_shipyards, trade_routes, callsign)
	}

	if ship.Registration.Role == "SATELLITE" {
		ApplyRoleSatellite(ship, markets_to_cover, trade_routes)
	}
}

func main() {

	// Ensure the CALLSIGN is provided as a command line argument
	if len(os.Args) != 2 {
		fmt.Println("go-spacetrade CALLSIGN")
		os.Exit(1)
	}

	CALLSIGN := os.Args[1]

	// Check if an auth token file is present for the CALLSIGN provided
	if !DoesAgentTokenFileExist(CALLSIGN) {
		ReadAccountTokenFromFile()
		RegisterAgent(CALLSIGN)
	}

	ReadAgentTokenFromFile(CALLSIGN)

	// TODO: globals are bad, this should be removed
	populate_base_system_symbol()

	//marketplace_waypoints := ListWaypointInSystemByTrait(base_system_symbol, "MARKETPLACE")
	//for _, marketplace_waypoint := range marketplace_waypoints {
	//	marketplace := GetMarket(base_system_symbol, marketplace_waypoint.Symbol)
	//	fmt.Println(marketplace.Symbol)
	//
	//	//fmt.Println("Exchange")
	//	//for _, exchange_good := range marketplace.Exchange {
	//	//	fmt.Println(exchange_good.Symbol)
	//	//}
	//	//
	//	//fmt.Println("Exports")
	//	//for _, export_good := range marketplace.Exports {
	//	//	fmt.Println(export_good.Symbol)
	//	//}
	//
	//	fmt.Println("Imports")
	//	for _, import_good := range marketplace.Imports {
	//		fmt.Println(import_good.Symbol)
	//	}
	//
	//	//fmt.Println("TradeGoods")
	//	//for _, tradegood_good := range marketplace.TradeGoods {
	//	//	fmt.Println(tradegood_good.Symbol)
	//	//}
	//
	//	fmt.Println()
	//
	//}

	// early exit while testing
	//os.Exit(0)

	// each unique market waypoint symbol (unordered)
	markets_to_cover := make(map[string]string)

	// association for places to BUY and SELL TradeGoods
	trade_routes := []TradeRoute{}

	if !DoesTradeRouteFileExist(CALLSIGN) {
		fmt.Println("[INFO] Trade route file does not exist. Initializing...")
		trade_routes = IdentifyTradeRoutes(markets_to_cover)
		WriteTradeRoutesToFile(trade_routes, CALLSIGN)
	} else {
		fmt.Println("[INFO] Trade file exists. Reading from file...")
		trade_routes = ReadTradeRoutesFromFile(CALLSIGN, trade_routes)
	}

	//fmt.Println("[DEBUG] markets_to_cover: " + string(len(markets_to_cover)))

	markets_to_cover = PopulateMarketsToCover(trade_routes)

	number_of_markets_to_cover := len(markets_to_cover)

	number_of_satellites := HowManySatellitesDoIOwn()

	if number_of_satellites >= number_of_markets_to_cover {
		fmt.Println("[INFO] Enough satellites")
		if !SatelliteToMarketAssignmentComplete(markets_to_cover) {
			AssignSatellitesToMarkets(markets_to_cover)
		}
	} else {
		fmt.Println("[INFO] Not enough satellites, postponing market assignments...")
	}

	// there can be multiple SHIPYARDs which sell SHIP_PROBE
	probe_shipyards := []Waypoint{}

	// populate probe_shipyards with Waypoints which have SHIPYARDs which sell SHIP_PROBEs
	shipyards_in_system := ListWaypointInSystemByTrait(base_system_symbol, "SHIPYARD")
	for _, shipyard_waypoint := range shipyards_in_system {
		get_shipyard_result := GetShipyard(base_system_symbol, shipyard_waypoint.Symbol)
		for _, ship := range get_shipyard_result.ShipTypes {
			if ship.Type == "SHIP_PROBE" {
				//fmt.Println("[DEBUG] shipyard with satellites for sale found: ")
				//fmt.Println("[DEBUG] " + get_shipyard_result.Symbol)
				probe_shipyards = append(probe_shipyards, shipyard_waypoint)
			}
		}
	}

	//fmt.Println("[DEBUG] markets to cover:")
	//fmt.Println(markets_to_cover)

	for market := range markets_to_cover {
		fmt.Println("[DEBUG] " + market)
	}

	turn_number := 1

	fmt.Print("[INFO] http calls: ")
	fmt.Print(http_calls)
	http_calls = 0
	fmt.Println()

	// this runs forever
	for {

		fmt.Print("[INFO] START OF TURN ")
		fmt.Print(turn_number)
		fmt.Println()

		agent := GetAgent()

		fmt.Println("[INFO] " + agent.Symbol)
		fmt.Print("[INFO] ShipCount: ")
		fmt.Print(agent.ShipCount)
		fmt.Println()
		fmt.Print("[INFO] Credits: ")
		fmt.Print(agent.Credits)
		fmt.Println()

		ships_list := ListShips()
		wait_between_ships := turn_length / len(ships_list)

		for _, ship := range ships_list {
			ShipRoleDecider(ship, markets_to_cover, probe_shipyards, trade_routes, CALLSIGN)

			// turns are always turn_length (default 2 minutes) but as we add ships they fill the time between turns
			time.Sleep(time.Duration(wait_between_ships) * time.Second)
		}

		// outro

		// inform user of http calls/turn to ease rate limit issues
		fmt.Print("[INFO] http calls: ")
		fmt.Print(http_calls / 2)
		fmt.Print("/m")
		fmt.Println()
		fmt.Println("[INFO] END OF TURN")

		// reset call counter
		http_calls = 0
		turn_number++
	}
}
