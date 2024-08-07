package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"
)

var url_base string = "https://api.spacetraders.io/v2/"
var bearer_token = "Bearer "
var base_system_symbol = ""
var http_calls = 0
var turn_length = 120

func main() {

	// Ensure the CALLSIGN is provided as a command line argument
	if len(os.Args) != 2 {
		fmt.Println("go-spacetrade CALLSIGN")
		os.Exit(1)
	}

	CALLSIGN := os.Args[1]

	// Check if an auth token file is present for the CALLSIGN provided
	if !does_auth_file_exist(CALLSIGN) {
		RegisterAgent(CALLSIGN)
	}

	read_auth_token_from_file(CALLSIGN)

	// TODO: globals are bad, this should be removed
	populate_base_system_symbol()

	// cache for full response from every get_market call
	all_market_results := []Market{}

	// populate all_market results with the result of get_market against each waypoint which has a MARKETPLACE
	marketplaces_in_system := list_waypoints_in_system_by_trait(base_system_symbol, "MARKETPLACE")
	for _, marketplace := range marketplaces_in_system {
		get_market_result := GetMarket(base_system_symbol, marketplace.Symbol)
		all_market_results = append(all_market_results, get_market_result.Data)
	}

	// association for places to BUY and SELL TradeGoods
	trade_routes := []TradeRoute{}

	// each unique market waypoint symbol (unordered)
	markets_to_cover := make(map[string]bool)

	// populate both trade_routes and markets_to_cover by iterating through A) every MARKETPLACE B) each of their Imports and C) their Exports
	// associations of import/export are added to trade_routes, and we keep one copy of each waypoint_symbol in markets_to_cover
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

	// there can be multiple SHIPYARDs which sell SHIP_PROBE
	probe_shipyards := []Waypoint{}

	// populate probe_shipyards with waypoints which have SHIPYARDs which sell SHIP_PROBE's
	shipyards_in_system := list_waypoints_in_system_by_trait(base_system_symbol, "SHIPYARD")
	for _, shipyard_waypoint := range shipyards_in_system {
		get_shipyard_result := GetShipyard(base_system_symbol, shipyard_waypoint.Symbol)
		for _, ship := range get_shipyard_result.ShipTypes {
			if ship.Type == "SHIP_PROBE" {
				fmt.Println("[DEBUG] shipyard with probes for sale found: ")
				fmt.Println(get_shipyard_result.Symbol)
				probe_shipyards = append(probe_shipyards, shipyard_waypoint)
			}
		}
	}

	// TODO: i think most of this goes into the turn loop
	list_ships_result := ListShips()

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

	fmt.Println("[DEBUG] markets to cover:")
	//fmt.Println(markets_to_cover)

	for market := range markets_to_cover {
		fmt.Println(market)
	}

	fmt.Println("[DEBUG] number of markets to cover:")
	fmt.Println(len(markets_to_cover))

	// count number of satellites
	number_of_satellites := 0

	var command_ship Ship

	// TODO: this may as well become the role assigner
	for _, ship := range list_ships_result {
		if ship.Registration.Role == "SATELLITE" {
			number_of_satellites++
		}
		if ship.Registration.Role == "COMMAND" {
			command_ship = ship
		}
	}

	// TODO: This goes into the COMMAND_SHIP's role
	command_ship_current_location := GetWaypoint(base_system_symbol, command_ship.Nav.WaypointSymbol)

	fmt.Println("[DEBUG] number of satellites")
	fmt.Println(number_of_satellites)

	if number_of_satellites < len(markets_to_cover) {
		fmt.Println("We need more satellites, boss")

		best_distance := 99999999.9999999
		var buyer_ship_destination_waypoint_symbol string
		for _, shipyard := range probe_shipyards {
			distance := distance_between_two_coordinates(shipyard.X, shipyard.Y, command_ship_current_location.X, command_ship_current_location.Y)
			if distance < best_distance {
				buyer_ship_destination_waypoint_symbol = shipyard.Symbol
			}
		}

		fmt.Println("[DEBUG] buyer_ship_destination_symbol:")
		fmt.Println(buyer_ship_destination_waypoint_symbol)

		fmt.Println("[DEBUG] command ship current location")
		fmt.Println(command_ship.Nav.WaypointSymbol)

		if is_ship_already_at_waypoint(command_ship, buyer_ship_destination_waypoint_symbol) {
			// TODO: buy satellites upto len(markets_to_cover)
			fmt.Println("[INFO] command ship is at buyer_ship_destination_waypoint_symbol BUY SATELLITES")
			if !is_ship_docked(command_ship) {
				DockShip(command_ship.Symbol)
			}
		} else {
			// TODO: send command ship to shipyard which sells satellites
			if is_ship_docked(command_ship) {
				OrbitShip(command_ship.Symbol)
			}
			navigate_ship_result := NavigateShip(command_ship.Symbol, buyer_ship_destination_waypoint_symbol)
			fmt.Println(navigate_ship_result)
		}
	}

	turn_number := 1

	// this runs forever
	for {
		fmt.Print("[INFO] START OF TURN ")
		fmt.Print(turn_number)
		fmt.Println()
		list_ships_result := ListShips()
		wait_between_ships := turn_length / len(list_ships_result)

		// inform user of http calls/turn to ease rate limit issues
		fmt.Print("[INFO] http calls: ")
		fmt.Print(http_calls / 2)
		fmt.Print("/m")
		fmt.Println()
		fmt.Println("[INFO] END OF TURN")

		// reset call counter
		http_calls = 0

		turn_number++

		// turns are always turn_length (default 2 minutes) but as we add ships they fill the time between turns
		time.Sleep(time.Duration(wait_between_ships) * time.Second)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func basic_get(endpoint string) (response_body string) {
	url := url_base + endpoint
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer_token)
	result, err := http.DefaultClient.Do(request)
	check(err)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	check(err)

	error_container := ErrorResponse{}
	if err := json.Unmarshal(body, &error_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}

	// If the error["message"] field exists, the game returned an error.
	if error_container.Error.Message != "" {
		fmt.Println("[ERROR] response contains error key")
		fmt.Println("Error Code:")
		fmt.Println(error_container.Error.Code)
		fmt.Println(error_container.Error.Message)
		fmt.Println(error_container.Error.Message)
	}

	sb := string(body)
	http_calls++
	return sb
}

func basic_post(endpoint string, payload []byte) (response_body string) {
	posturl := url_base + endpoint
	request, err := http.NewRequest("POST", posturl, bytes.NewBuffer(payload))
	check(err)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer_token)
	client := &http.Client{}
	result, err := client.Do(request)
	check(err)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	check(err)

	error_container := ErrorResponse{}
	if err := json.Unmarshal(body, &error_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}

	// If the error["message"] field exists, the game returned an error.
	if error_container.Error.Message != "" {
		fmt.Println("[ERROR] response contains error key")
		fmt.Println("Error Code:")
		fmt.Println(error_container.Error.Code)
		fmt.Println(error_container.Error.Message)
		fmt.Println(error_container.Error.Data)
	}

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

func WriteAuthTokenToFile(auth_token string, filename string) {
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
	result := basic_get("")
	pretty_print_json(result)
}

func RegisterAgent(callsign string) (result RegisterAgentResponse) {
	fmt.Println("RegisterAgent")
	payload := &RegisterAgentPayload{}
	payload.Faction = "COSMIC"
	payload.Symbol = callsign
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post("register", payloadJSON)
	data_container := RegisterAgentResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	token := data_container.Data.Token
	auth_token := token
	WriteAuthTokenToFile(auth_token, callsign+".token")
	return data_container.Data
}

func ListShips() (ships []Ship) {
	//fmt.Println("[DEBUG] list_ships")
	endpoint := "my/ships"
	response_string := basic_get(endpoint)

	data_container := ListShipsResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("failed to unmarshal")
	}

	return data_container.Data

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

func GetWaypoint(system_symbol string, waypoint_symbol string) (resultant_waypoint Waypoint) {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol
	response_string := basic_get(endpoint)
	data_container := GetWaypointResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func get_waypoint_coordinate(waypoint Waypoint) (waypointX int64, waypointY int64) {
	return waypoint.X, waypoint.Y
}

func distance_between_two_coordinates(waypoint1X int64, waypoint1Y int64, waypoint2X int64, waypoint2Y int64) (resultant_distance float64) {
	//fmt.Println("[DEBUG] distance_between_two_coordinates")

	XIntermediate := waypoint1X - waypoint2X
	YIntermediate := waypoint1Y - waypoint2Y
	XSquared := XIntermediate * XIntermediate
	YSquared := YIntermediate * YIntermediate
	XPlusY := XSquared + YSquared
	XPlusYFloat := float64(XPlusY)
	resultant_distance = math.Sqrt(XPlusYFloat)

	fmt.Println("[DEBUG] distance: ")
	fmt.Println(resultant_distance)
	return resultant_distance
}

func list_waypoints_in_system_by_trait(system_symbol string, trait string) []Waypoint {
	endpoint := "systems/" + system_symbol + "/waypoints?traits=" + trait
	response_string := basic_get(endpoint)
	data_container := ListWaypointsInSystemResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func list_waypoints_in_system_by_type(system_symbol string, query_type string) (list_waypoints_in_system_result []Waypoint) {
	endpoint := "systems/" + system_symbol + "/waypoints?type=" + query_type
	response_string := basic_get(endpoint)
	data_container := ListWaypointsInSystemResponseData{}
	if err := json.Unmarshal([]byte(response_string), &list_waypoints_in_system_result); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func GetMarket(system_symbol string, waypoint_symbol string) (get_market_result GetMarketResponseData) {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/market"
	response_string := basic_get(endpoint)
	if err := json.Unmarshal([]byte(response_string), &get_market_result); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return get_market_result
}

func GetShipyard(system_symbol string, waypoint_symbol string) (get_shipyard_result Shipyard) {
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

func is_satellite_present_at_marketplace(list_ships_result []Ship, waypoint_symbol string) (answer bool) {
	list_ships_data := list_ships_result

	for _, ship := range list_ships_data {
		if ship.Registration.Role == "SATELLITE" {
			if ship.Nav.WaypointSymbol == waypoint_symbol {
				return true
			}
		}
	}
	return false
}

func is_ship_already_at_waypoint(ship_to_test Ship, waypoint_symbol string) bool {
	return ship_to_test.Nav.WaypointSymbol == waypoint_symbol
}

func NavigateShip(ship_symbol string, waypoint_symbol string) NavigateShipResponse {
	fmt.Println("[DEBUG] NavigateShip " + ship_symbol + " " + waypoint_symbol)
	endpoint := "my/ships/" + ship_symbol + "/navigate"
	payload := &NavigateShipPayload{}

	payload.WaypointSymbol = waypoint_symbol
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := NavigateShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data

}

func is_ship_docked(ship Ship) (is_docked bool) {
	return ship.Nav.Status == "DOCKED"
}

func OrbitShip(ship_symbol string) NavigateShipResponse {
	fmt.Println("[DEBUG] OrbitShip")
	endpoint := "my/ships/" + ship_symbol + "/orbit"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := NavigateShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func DockShip(ship_symbol string) DockShipResponse {
	fmt.Println("[DEBUG] DockShip")
	endpoint := "my/ships/" + ship_symbol + "/dock"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := DockShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}
