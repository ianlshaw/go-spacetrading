package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var url_base string = "https://api.spacetraders.io/v2/"
var bearer_token = "Bearer "
var base_system_symbol = ""
var http_calls = 0

func main() {

	// Ensure the CALLSIGN is provided as a command line argument
	if len(os.Args) != 2 {
		fmt.Println("go-spacetrade CALLSIGN")
		os.Exit(1)
	}

	CALLSIGN := os.Args[1]

	// Check if an auth token file is present for the CALLSIGN provided
	if !does_auth_file_exist(CALLSIGN) {
		fmt.Println("this is where we would register agent")
		//register_agent(CALLSIGN)
	}

	read_auth_token_from_file(CALLSIGN)

	// TODO: globals are bad, this should be removed
	populate_base_system_symbol()

	// association for places to BUY and SELL TradeGoods
	trade_routes := []TradeRoute{}

	// cache for full response from every get_market call
	all_market_results := []Market{}

	// each unique market waypoint symbol
	markets_to_cover := make(map[string]bool)

	marketplaces_in_system := list_waypoints_in_system_by_trait(base_system_symbol, "MARKETPLACE")
	for _, marketplace := range marketplaces_in_system.Data {
		get_market_result := get_market(base_system_symbol, marketplace.Symbol)
		all_market_results = append(all_market_results, get_market_result.Data)
	}

	for _, each_market_result := range all_market_results {
		if len(each_market_result.Exports) > 0 {
			for _, each_export := range each_market_result.Exports {
				for _, each_market_result_inner := range all_market_results {
					if len(each_market_result_inner.Imports) > 0 {
						for _, each_market_result_imports := range each_market_result_inner.Imports {
							if each_export.Symbol == each_market_result_imports.Symbol {
								trade_route := TradeRoute{}
								trade_route.TradeGoodSymbol = each_export.Symbol
								trade_route.BuyWaypointSymbol = each_market_result.Symbol
								trade_route.SellWaypointSymbol = each_market_result_inner.Symbol
								trade_routes = append(trade_routes, trade_route)
								markets_to_cover[trade_route.BuyWaypointSymbol] = true
								markets_to_cover[trade_route.SellWaypointSymbol] = true
							}
						}
					}
				}
			}
		}
	}

	fmt.Println("Shipyard finder...")
	shipyards_in_system := list_waypoints_in_system_by_trait(base_system_symbol, "SHIPYARD")
	for _, shipyard_waypoint := range shipyards_in_system.Data {
		get_shipyard_result := get_shipyard(base_system_symbol, shipyard_waypoint.Symbol)
		for _, ship := range get_shipyard_result.ShipTypes {
			if ship.Type == "SHIP_PROBE" {
				fmt.Println("Shipyard with probes for sale found: ")
				fmt.Println(get_shipyard_result.Symbol)

			}
		}
	}

	list_ships_result := list_ships()

	for _, each_trade_route := range trade_routes {
		fmt.Println("BUY " + each_trade_route.TradeGoodSymbol + " AT " + each_trade_route.BuyWaypointSymbol + " SELL AT " + each_trade_route.SellWaypointSymbol)
		if is_satellite_present_at_marketplace(list_ships_result, each_trade_route.BuyWaypointSymbol) {
			fmt.Println("satellite present at BUY waypoint")
		} else {
			//fmt.Println("No satellite present at BUY waypoint")
		}
		if is_satellite_present_at_marketplace(list_ships_result, each_trade_route.SellWaypointSymbol) {
			fmt.Println("satellite present at SELL waypoint")
		} else {
			//fmt.Println("No satellite present at SELL waypoint")
		}
	}

	fmt.Println("markets to cover:")
	//fmt.Println(markets_to_cover)

	for market := range markets_to_cover {
		fmt.Println(market)
	}

	fmt.Println("Number of markets to cover:")
	fmt.Println(len(markets_to_cover))

	// count number of satellites
	number_of_satellites := 0

	for _, ship := range list_ships_result.Data {
		if ship.Registration.Role == "SATELLITE" {
			number_of_satellites++
		}
	}

	fmt.Println("Number of satellites")
	fmt.Println(number_of_satellites)

	if number_of_satellites < len(markets_to_cover) {
		fmt.Println("We need more satellites, boss")

		// send command ship to shipyard which sells satellites

		// buy satellites
	}

	fmt.Println("http calls:")
	fmt.Println(http_calls)

	// match imports to an exports

	//jumpgates_in_system := list_waypoints_in_system_by_type(base_system_symbol, "JUMP_GATE")
	//fmt.Println("JUMP_GATEs")
	//for _, jumpgate_waypoint := range jumpgates_in_system.Data {
	//	fmt.Println(jumpgate_waypoint.Symbol)
	//	jumpgate := get_jump_gate(base_system_symbol, jumpgate_waypoint.Symbol)
	//	fmt.Println(jumpgate.Data.Connections)
	//}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type RegisterAgentPayload struct {
	Symbol  string `json:"symbol"`
	Faction string `json:"faction"`
}

func dumb_get(endpoint string) (response_body string) {
	resp, err := http.Get(url_base + endpoint)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	//log.Printf(sb)

	return sb
}

func basic_get(endpoint string) (response_body string) {
	url := url_base + endpoint

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", bearer_token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//Convert the body to type string
	sb := string(body)
	//log.Printf(sb)

	http_calls++

	return sb
}

func basic_post(endpoint string, payload []byte) (response_body string) {
	// HTTP endpoint
	posturl := url_base + endpoint

	// JSON body

	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(payload))
	check(err)

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	check(err)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	check(err)
	//Convert the body to type string
	sb := string(body)
	http_calls++
	return sb
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

func does_auth_file_exist(callsign string) (result bool) {
	var filename = callsign + ".token"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("[ERROR] Token file does not exist")
		return false
	}
	fmt.Println("[INFO] Token file exists")
	return true
}

func write_auth_token_to_file(auth_token string, filename string) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()
	write_string_result, err := f.WriteString(auth_token)
	check(err)
	fmt.Printf("wrote %d bytes\n", write_string_result)
}

func read_auth_token_from_file(callsign string) {
	f, err := os.ReadFile(callsign + ".token") // just pass the file name
	check(err)
	bearer_token += (string(f))
}

func get_status() {
	result := dumb_get("")
	pretty_print_json(result)
}

func register_agent(callsign string) {
	fmt.Println("register_agent")

	payload := &RegisterAgentPayload{}

	payload.Faction = "COSMIC"
	payload.Symbol = callsign

	payloadJSON, err := json.Marshal(payload)
	check(err)

	result := basic_post("register", payloadJSON)

	fmt.Println(result)

	// TODO: create model for register_agent response and convert this to unmarshal and return
	//token := result["data"].(map[string]interface{})["token"]

	//auth_token := token.(string)
	//write_auth_token_to_file(auth_token, callsign+".token")
}

func list_ships() (list_ships_result ListShipsResponseData) {
	fmt.Println("list_ships")
	endpoint := "my/ships"
	response_string := basic_get(endpoint)

	//response_typed := ListShipsResponseData{}
	if err := json.Unmarshal([]byte(response_string), &list_ships_result); err != nil {
		fmt.Println("failed to unmarshal")
	}

	return list_ships_result

}

func populate_base_system_symbol() {
	//fmt.Println("[DEBUG] populate_base_system_symbol")
	endpoint := "my/ships"
	response_string := basic_get(endpoint)

	response_typed := ListShipsResponseData{}
	if err := json.Unmarshal([]byte(response_string), &response_typed); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	base_system_symbol = response_typed.Data[0].Nav.SystemSymbol
}

func get_waypoint_coordinate(waypoint Waypoint) (waypointX int64, waypointY int64) {
	return waypoint.X, waypoint.Y
}

func distance_between_two_coordinates(waypoint1X int64, waypoint1Y int64, waypoint2X int64, waypoint2Y int64) (resultant_distance int64) {

	return resultant_distance
}

func list_waypoints_in_system_by_trait(system_symbol string, trait string) (list_waypoints_in_system_result ListWaypointsInSystemResponseData) {
	endpoint := "systems/" + system_symbol + "/waypoints?traits=" + trait
	response_string := basic_get(endpoint)
	if err := json.Unmarshal([]byte(response_string), &list_waypoints_in_system_result); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return list_waypoints_in_system_result
}

func list_waypoints_in_system_by_type(system_symbol string, query_type string) (list_waypoints_in_system_result ListWaypointsInSystemResponseData) {
	endpoint := "systems/" + system_symbol + "/waypoints?type=" + query_type
	response_string := basic_get(endpoint)
	if err := json.Unmarshal([]byte(response_string), &list_waypoints_in_system_result); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return list_waypoints_in_system_result
}

func get_market(system_symbol string, waypoint_symbol string) (get_market_result GetMarketResponseData) {

	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/market"
	response_string := basic_get(endpoint)
	if err := json.Unmarshal([]byte(response_string), &get_market_result); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return get_market_result
}

func get_shipyard(system_symbol string, waypoint_symbol string) (get_shipyard_result Shipyard) {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/shipyard"
	response_string := basic_get(endpoint)
	data_container := GetShipyardResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func get_jump_gate(system_symbol string, waypoint_symbol string) (get_jump_gate_result GetJumpGateResponseData) {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "jump-gate"
	response_string := basic_get(endpoint)
	if err := json.Unmarshal([]byte(response_string), &get_jump_gate_result); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return get_jump_gate_result
}

func is_satellite_present_at_marketplace(list_ships_result ListShipsResponseData, waypoint_symbol string) (answer bool) {
	list_ships_data := list_ships_result.Data

	for _, ship := range list_ships_data {
		if ship.Registration.Role == "SATELLITE" {
			if ship.Nav.WaypointSymbol == waypoint_symbol {
				return true
			}
		}
	}
	return false
}
