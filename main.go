package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

	//agent := GetAgent()
	//contracts := ListContracts()

	// do waypoint files exist?

	if !DoesWaypointsFileExist(CALLSIGN) {
		fmt.Println("[INFO] Gathering waypoint data...")
		all_waypoints_in_system := []Waypoint{}

		list_waypoints_result := ListWaypointsInSystem(base_system_symbol, "1")

		total_waypoints := list_waypoints_result.Meta.Total

		limit := list_waypoints_result.Meta.Limit

		loop_iterations_required := total_waypoints / int64(limit)

		for i := 1; i < int(loop_iterations_required+2); i++ {
			a_page_of_waypoints := ListWaypointsInSystem(base_system_symbol, strconv.FormatInt(int64(i), 10))
			all_waypoints_in_system = append(all_waypoints_in_system, a_page_of_waypoints.Data...)
			fmt.Println(len(all_waypoints_in_system))
		}
		WriteWaypointsToFile(all_waypoints_in_system, CALLSIGN)
	}

	// read all waypoint data from file
	all_waypoints_in_system := []Waypoint{}
	all_waypoints_in_system = ReadWaypointsFromFile(CALLSIGN, all_waypoints_in_system)

	if !DoesShipyardsFileExist(CALLSIGN) {
		shipyard_waypoints := []Waypoint{}

		for _, waypoint := range all_waypoints_in_system {
			for _, trait := range waypoint.Traits {
				if trait.Symbol == "SHIPYARD" {
					shipyard_waypoints = append(shipyard_waypoints, waypoint)
				}
			}
		}

		all_shipyards_in_system := []Shipyard{}

		for _, shipyard_waypoint := range shipyard_waypoints {
			get_shipyard_result := GetShipyard(base_system_symbol, shipyard_waypoint.Symbol)
			all_shipyards_in_system = append(all_shipyards_in_system, get_shipyard_result)
		}

		WriteShipyardsToFile(all_shipyards_in_system, CALLSIGN)
	}

	all_shipyards_in_system := []Shipyard{}
	ReadShipyardsFromFile(CALLSIGN, all_shipyards_in_system)

	if !DoesMarketsFileExist(CALLSIGN) {
		all_markets_in_system := []Market{}
		for _, waypoint := range all_waypoints_in_system {
			for _, trait := range waypoint.Traits {
				if trait.Symbol == "MARKETPLACE" {
					get_market_result := GetMarket(base_system_symbol, waypoint.Symbol)
					all_markets_in_system = append(all_markets_in_system, get_market_result)
					time.Sleep(2 * time.Second)
				}
			}
		}
		WriteMarketsToFile(all_markets_in_system, CALLSIGN)
	}

	all_markets_in_system := []Market{}
	ReadMarketsFromFile(CALLSIGN, all_markets_in_system)

	marketplace_waypoints := []Waypoint{}

	for _, market := range all_markets_in_system {
		for _, waypoint := range all_waypoints_in_system {
			if market.Symbol == waypoint.Symbol {
				marketplace_waypoints = append(marketplace_waypoints, waypoint)
			}
		}
	}

	// each unique market waypoint symbol (unordered)
	markets_to_cover := make(map[string]string)

	// association for places to BUY and SELL TradeGoods
	trade_routes := []TradeRoute{}

	if !DoesTradeRouteFileExist(CALLSIGN) {
		fmt.Println("[INFO] Trade route file does not exist. Initializing...")
		trade_routes = IdentifyTradeRoutes(CALLSIGN, markets_to_cover, all_markets_in_system, marketplace_waypoints)
		WriteTradeRoutesToFile(trade_routes, CALLSIGN)
	} else {
		fmt.Println("[INFO] Trade file exists. Reading from file...")
		trade_routes = ReadTradeRoutesFromFile(CALLSIGN, trade_routes)
	}

	markets_to_cover = PopulateMarketsToCover(trade_routes)

	number_of_markets_to_cover := len(markets_to_cover)

	fmt.Println("[DEBUG] number_of_markets_to_cover = ")
	fmt.Println(number_of_markets_to_cover)

	number_of_satellites := HowManySatellitesDoIOwn()

	if number_of_satellites >= number_of_markets_to_cover {
		fmt.Println("[INFO] Enough satellites")
		if !SatelliteToMarketAssignmentComplete(markets_to_cover) {
			AssignSatellitesToMarkets(markets_to_cover)
		}
	} else {
		fmt.Println("[INFO] Not enough satellites, postponing market assignments...")
	}

	probe_shipyards := []Shipyard{}
	probe_shipyard_waypoints := []Waypoint{}
	for _, shipyard := range all_shipyards_in_system {
		for _, ship := range shipyard.ShipTypes {
			if ship.Type == "SHIP_PROBE" {
				probe_shipyards = append(probe_shipyards, shipyard)
				for _, waypoint := range all_waypoints_in_system {
					if waypoint.Symbol == shipyard.Symbol {
						probe_shipyard_waypoints = append(probe_shipyard_waypoints, waypoint)
					}
				}
			}
		}
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
			ShipRoleDecider(ship, markets_to_cover, probe_shipyard_waypoints, trade_routes, CALLSIGN)

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
