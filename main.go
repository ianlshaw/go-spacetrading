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
		fmt.Println("[ERROR] populate_base_system_symbol failed to unmarshal")
	}
	base_system_symbol = response_typed.Data[0].Nav.SystemSymbol
}

func runShip(
	ship Ship,
	all_waypoints_in_system []Waypoint,
	all_markets_in_system []Market,
	markets_to_cover map[string]string,
	probe_shipyard_waypoints []Waypoint,
	shuttle_shipyard_waypoints []Waypoint,
	mining_drone_shipyard_waypoints []Waypoint,
	siphon_drone_shipyard_waypoints []Waypoint,
	surveyor_shipyard_waypoints []Waypoint,
	trade_routes []TradeRoute,
	ship_list []Ship,
	agent Agent,
	callsign string) {

	for {

		//var expiration time.Time
		expiration := ThreeHoursFromNow()

		if ship.Registration.Role == "COMMAND" {
			expiration = ApplyRoleCommand(ship, all_waypoints_in_system, all_markets_in_system, markets_to_cover, trade_routes, callsign)
		}

		all_probes := []Ship{}
		all_shuttles := []Ship{}

		for _, ship := range ship_list {
			if ship.Registration.Role == "SATELLITE" {
				all_probes = append(all_probes, ship)
			}
			if ship.Registration.Role == "TRANSPORT" {
				all_shuttles = append(all_shuttles, ship)
			}
		}

		buyer_ship := all_probes[0]

		if ship.Registration.Role == "SATELLITE" {
			if ship.Symbol == buyer_ship.Symbol {
				expiration = ApplyRoleBuyer(
					ship,
					ship_list,
					markets_to_cover,
					probe_shipyard_waypoints,
					shuttle_shipyard_waypoints,
					mining_drone_shipyard_waypoints,
					siphon_drone_shipyard_waypoints,
					surveyor_shipyard_waypoints,
					agent)
			}
			expiration = ApplyRoleSatellite(ship, markets_to_cover, trade_routes)
		}
		//}
		//if ship.Registration.Role == "EXCAVATOR" {
		//	for _, mount := range ship.Mounts {
		//		if mount.Symbol == "MOUNT_MINING_LASER_I" {
		//			ApplyRoleMiner()
		//			return
		//		}
		//		if mount.Symbol == "MOUNT_GAS_SIPHON_I" {
		//			ApplyRoleSiphoner(ship, all_waypoints_in_system, all_markets_in_system)
		//			return
		//		}
		//	}
		//}
		//if ship.Registration.Role == "SURVEYOR" {
		//	ApplyRoleSurveyor()
		//	return
		//}


		if len(all_shuttles) >= 1 {
			if ship.Symbol == all_shuttles[0].Symbol {
				expiration = ApplyRoleTrader(ship)
			}
		}

		if len(all_shuttles) >= 2 {
			if ship.Symbol == all_shuttles[1].Symbol {
				expiration = ApplyRoleTransportOre(ship)
			}
		}

		if len(all_shuttles) >= 3 {
			if ship.Symbol == all_shuttles[2].Symbol {
				expiration = ApplyRoleTransportGas(ship)
			}
		}

		expiration_formatted := expiration.Format(time.RFC3339)
		fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " Sleeping until " + expiration_formatted)
		time.Sleep(time.Until(expiration))

		// Anti-Spam
		//Log("DEBUG", ship.Symbol + " ANTI SPAM ENGAGED")
		//time.Sleep(60 * time.Second)
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
		}
		WriteWaypointsToFile(all_waypoints_in_system, CALLSIGN)
	}

	all_waypoints_in_system := ReadWaypointsFromFile(CALLSIGN)

	for _, waypoint := range all_waypoints_in_system {
		//AddWaypointToSystemGraph(waypoint)
		PopulateGraphDistancesForWaypoint(SystemGraph, all_waypoints_in_system, waypoint)
		PopulateGraphDistancesForWaypointWithMaximum(SystemGraph, all_waypoints_in_system, waypoint, 400)
	}

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

	all_shipyards_in_system := ReadShipyardsFromFile(CALLSIGN)

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

	all_markets_in_system := ReadMarketsFromFile(CALLSIGN)
	marketplace_waypoints := []Waypoint{}

	for _, market := range all_markets_in_system {
		for _, waypoint := range all_waypoints_in_system {
			if market.Symbol == waypoint.Symbol {
				marketplace_waypoints = append(marketplace_waypoints, waypoint)
			}
		}
	}

	for _, waypoint := range marketplace_waypoints {
		//AddWaypointToSystemGraph(waypoint)
		PopulateGraphDistancesForWaypointWithMaximum(MarketplaceGraph, marketplace_waypoints, waypoint, 400)
		PopulateGraphDistancesForWaypointWithMaximum(SiphonerMarketplaceGraph, marketplace_waypoints, waypoint, 80)
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
		trade_routes = ReadTradeRoutesFromFile(CALLSIGN)
	}

	markets_to_cover = PopulateMarketsToCover(trade_routes)

	probe_shipyards, probe_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_PROBE")
	fmt.Print("[DEBUG] probe shipyards:")
	fmt.Println(len(probe_shipyards))

	shuttle_shipyards, shuttle_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_LIGHT_SHUTTLE")
	fmt.Print("[DEBUG] shuttle shipyards:")
	fmt.Println(len(shuttle_shipyards))

	mining_drone_shipyards, mining_drone_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_MINING_DRONE")
	fmt.Print("[DEBUG] mining_drone shipyards:")
	fmt.Println(len(mining_drone_shipyards))

	siphon_drone_shipyards, siphon_drone_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_SIPHON_DRONE")
	fmt.Print("[DEBUG] siphon_drone shipyards:")
	fmt.Println(len(siphon_drone_shipyards))

	surveyor_shipyards, surveyor_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_SURVEYOR")
	fmt.Print("[DEBUG] surveyor shipyards:")
	fmt.Println(len(surveyor_shipyards))

	//turn_number := 1

	fmt.Print("[INFO] http calls: ")
	fmt.Print(http_calls)
	http_calls = 0
	fmt.Println()



	agent := GetAgent()
	fmt.Print("[INFO] ShipCount: ")
	fmt.Print(agent.ShipCount)
	fmt.Println()
	fmt.Print("[INFO] Credits: ")
	fmt.Print(agent.Credits)
	fmt.Println()



	ships_list := ListShips()

	//wait_between_ships := turn_length / len(ships_list)

	for _, ship := range ships_list {
		go runShip(
			ship,
			all_waypoints_in_system,
			all_markets_in_system,
			markets_to_cover,
			probe_shipyard_waypoints,
			shuttle_shipyard_waypoints,
			mining_drone_shipyard_waypoints,
			siphon_drone_shipyard_waypoints,
			surveyor_shipyard_waypoints,
			trade_routes,
			ships_list,
			agent,
			CALLSIGN)
			fmt.Print("[INFO] http calls:")
			fmt.Println(http_calls)
		}

	select {}
}