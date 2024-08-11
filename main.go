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
	"sort"
	"time"
)

var url_base string = "https://api.spacetraders.io/v2/"
var bearer_token = "Bearer "
var base_system_symbol = ""
var http_calls = 0
var turn_length = 120

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func basic_get(endpoint string) (response_body string) {
	url := url_base + endpoint

	// DEBUG
	fmt.Println("[DEBUG] " + url)
	// DEBUG

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

	// DEBUG
	fmt.Println("[DEBUG] " + posturl)
	// DEBUG

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

func GetAgent() Agent {
	endpoint := "my/agent"
	response_string := basic_get(endpoint)

	data_container := GetAgentResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("failed to unmarshal")
	}

	return data_container.Data
}

func ListShips() (ships []Ship) {
	//fmt.Println("[DEBUG] list_ships")
	endpoint := "my/ships"
	response_string := basic_get(endpoint)

	data_container := ListShipsResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("failed to unmarshal")
		//fmt.Println(response_string)
		fmt.Println(err)
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
	return resultant_distance
}

func distance_between_two_waypoints(waypoint1 Waypoint, waypoint2 Waypoint) float64 {
	return distance_between_two_coordinates(waypoint1.X, waypoint1.Y, waypoint2.X, waypoint2.Y)
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

func GetMarket(system_symbol string, waypoint_symbol string) Market {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/market"
	response_string := basic_get(endpoint)
	data_container := GetMarketResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
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

func is_a_satellite_docked_at_marketplace(list_ships_result []Ship, waypoint_symbol string) (answer bool) {
	for _, ship := range list_ships_result {
		if ship.Registration.Role == "SATELLITE" {
			if ship.Nav.WaypointSymbol == waypoint_symbol {
				if ship.Nav.Status == "DOCKED" {
					return true
				}
			}
		}
	}
	return false
}

func is_ship_already_at_waypoint(ship_to_test Ship, waypoint_symbol string) bool {
	return (ship_to_test.Nav.WaypointSymbol == waypoint_symbol && ship_to_test.Nav.Status != "IN_TRANSIT")
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

func is_ship_docked(ship Ship) bool {
	return ship.Nav.Status == "DOCKED"
}

func is_ship_cargo_empty(ship Ship) bool {
	return ship.Cargo.Units == 0
}

func OrbitShip(ship_symbol string) OrbitShipResponse {
	//fmt.Println("[DEBUG] OrbitShip")
	endpoint := "my/ships/" + ship_symbol + "/orbit"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := OrbitShipResponseData{}
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

func PurchaseShip(ship_type string, waypoint_symbol string) PurchaseShipResponse {
	fmt.Println("[DEBUG] PurchaseShip")
	endpoint := "my/ships/"
	payload := &PurchaseShipPayload{}
	payload.WaypointSymbol = waypoint_symbol
	payload.ShipType = ship_type
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := PurchaseShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func PurchaseCargo(ship_symbol string, trade_good_symbol string, units int64) PurchaseCargoResponse {
	fmt.Println("[DEBUG] PurchaseCargo")
	endpoint := "my/ships/" + ship_symbol + "/purchase"
	payload := &PurchaseCargoPayload{}
	payload.Symbol = trade_good_symbol
	payload.Units = units
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := PurchaseCargoResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func SellCargo(ship_symbol string, trade_good_symbol string, units int64) SellCargoResponse {
	fmt.Println("[DEBUG] SellCargo")
	endpoint := "my/ships/" + ship_symbol + "/sell"
	payload := &SellCargoPayload{}
	payload.Symbol = trade_good_symbol
	payload.Units = units
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := SellCargoResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func RefuelShip(ship_symbol string) RefuelShipResponse {
	fmt.Println("[DEBUG] RefuelShip " + ship_symbol)
	endpoint := "my/ships/" + ship_symbol + "/refuel"
	payload := &RefuelShipPayload{}
	//payload.Units = units
	//payload.FromCargo = false
	payloadJSON, err := json.Marshal(payload)
	check(err)
	response_string := basic_post(endpoint, payloadJSON)
	data_container := RefuelShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] failed to unmarshal")
	}
	return data_container.Data
}

func MostProfitableTradeRoute(trade_routes []TradeRoute) TradeRoute {
	most_profitable_trade_route := TradeRoute{}
	var best_profitability_score = 0.0
	for _, trade_route := range trade_routes {
		if trade_route.ProfitabilityRating > best_profitability_score {
			most_profitable_trade_route = trade_route
			best_profitability_score = trade_route.ProfitabilityRating
		}
	}
	return most_profitable_trade_route
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

func UpdateTradeRoutesIncludingThisWaypoint(waypoint_symbol string, trade_routes []TradeRoute) {
	market := GetMarket(base_system_symbol, waypoint_symbol)
	for i, trade_route := range trade_routes {
		if waypoint_symbol == trade_route.BuyWaypoint.Symbol {
			for _, trade_good := range market.TradeGoods {
				if trade_route.TradeGoodSymbol == trade_good.Symbol {
					trade_routes[i].BuyMarketTradeGood = trade_good
				}
			}
		}

		if waypoint_symbol == trade_route.SellWaypoint.Symbol {
			for _, trade_good := range market.TradeGoods {
				if trade_route.TradeGoodSymbol == trade_good.Symbol {
					trade_routes[i].SellMarketTradeGood = trade_good
				}
			}
		}
	}
}

func MarketScanComplete(trade_routes []TradeRoute) bool {
	for _, trade_route := range trade_routes {
		if trade_route.BuyMarketTradeGood.PurchasePrice == 0 {
			fmt.Println("[INFO] MARKET DATA INCOMPLETE, WAIT FOR INPUT")
			return false
		}
	}
	fmt.Println("[INFO] MARKET DATA COMPLETE. LETS TRADE")
	return true
}

func PopulateTradeRoutesWithWaypointData(trade_routes []TradeRoute, markets_to_cover map[string]string) {
	fmt.Println("PopulateTradeRoutesWithWaypointData")

	for market_waypoint := range markets_to_cover {
		get_waypoint_result := GetWaypoint(base_system_symbol, market_waypoint)
		for i, trade_route := range trade_routes {
			if market_waypoint == trade_route.BuyMarketplaceWaypointSymbol {
				trade_routes[i].BuyWaypoint = get_waypoint_result
			}
			if market_waypoint == trade_route.SellMarketplaceWaypointSymbol {
				trade_routes[i].SellWaypoint = get_waypoint_result
			}
		}
	}
}

func PopulateTradeRoutesWithDistances(trade_routes []TradeRoute) {
	for i, trade_route := range trade_routes {
		distance := distance_between_two_waypoints(trade_route.BuyWaypoint, trade_route.SellWaypoint)
		trade_routes[i].Distance = distance
	}
}

func CalculateProfitPerUnit(trade_route TradeRoute) float64 {
	//fmt.Println("[DEBUG] CalculateProfitPerUnit")
	sell_price := float64(trade_route.SellMarketTradeGood.SellPrice)
	buy_price := float64(trade_route.BuyMarketTradeGood.PurchasePrice)
	profit_per_unit := sell_price - buy_price
	return profit_per_unit
}

func PopulateTradeRoutesProfitPerUnit(trade_routes []TradeRoute) {
	for i, trade_route := range trade_routes {
		profit_per_unit := CalculateProfitPerUnit(trade_route)
		trade_routes[i].ProfitPerUnit = int64(profit_per_unit)
		profit_per_unit_divide_by_distance_times_two := float64(profit_per_unit) / (trade_route.Distance * 2)
		trade_routes[i].ProfitabilityRating = profit_per_unit_divide_by_distance_times_two
	}
}

// calculate credits/second generated by running this trade route
func PrintTradeRoutes(ship_list []Ship, trade_routes []TradeRoute) {
	for _, trade_route := range trade_routes {
		fmt.Print("[INFO] BUY ")
		fmt.Print(trade_route.TradeGoodSymbol)
		fmt.Print(" AT ")
		fmt.Print(trade_route.BuyMarketplaceWaypointSymbol)
		fmt.Print(" FOR ")
		fmt.Print(trade_route.BuyMarketTradeGood.PurchasePrice)
		fmt.Print(" SELL AT ")
		fmt.Print(trade_route.SellMarketplaceWaypointSymbol)
		fmt.Print(" FOR ")
		fmt.Print(trade_route.SellMarketTradeGood.SellPrice)
		fmt.Print(" PPU ")
		fmt.Print(trade_route.ProfitPerUnit)
		fmt.Print(" DISTANCE ")
		fmt.Print(trade_route.Distance)
		fmt.Print(" SCORE ")
		fmt.Print(trade_route.ProfitabilityRating)
		fmt.Println()
	}
}

func CountTradeGoodCargo(ship Ship, trade_good_symbol string) int64 {
	inventory := ship.Cargo.Inventory
	for _, trade_good := range inventory {
		if trade_good_symbol == trade_good.Symbol {
			return trade_good.Units
		}
	}
	return 0
}

func HowManyTradeGoodCanIAfford(agent Agent, trade_good TradeGood) int64 {
	fmt.Println("[DEBUG] HowManyTradeGoodCanIAfford")
	fmt.Print("[DEBUG] agent.Credits = ")
	fmt.Print(agent.Credits)
	fmt.Println()
	fmt.Print("[DEBUG] trade_good.PurchasePrice = ")
	fmt.Print(trade_good.PurchasePrice)
	fmt.Println()
	max_buy_count := agent.Credits / trade_good.PurchasePrice
	return max_buy_count
}

func ApplyRoleCommand(ship Ship, markets_to_cover map[string]string, probe_shipyards []Waypoint, trade_routes []TradeRoute) {

	fmt.Println("[INFO] " + ship.Symbol)

	//fmt.Println("[DEBUG] ApplyRoleCommand")

	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[DEBUG] IN_TRANSIT TO")
		fmt.Println(ship.Nav.Route.Destination)
		fmt.Println("[DEBUG] Arrival")
		fmt.Println(ship.Nav.Route.Arrival)
		return
	}

	fmt.Print("[DEBUG] ship.Cargo.Inventory")
	fmt.Print(ship.Cargo.Inventory)
	fmt.Println()

	// count number of satellites
	var number_of_satellites int

	ship_list := ListShips()

	// TODO: not sure this needs to be here or exist
	for _, a_ship := range ship_list {
		if a_ship.Registration.Role == "SATELLITE" {
			number_of_satellites++
		}
	}

	// we need the X and Y coord of the command ship to figure out which shipyard is closest
	current_waypoint := GetWaypoint(base_system_symbol, ship.Nav.WaypointSymbol)

	if number_of_satellites < len(markets_to_cover) {
		fmt.Println("[INFO] We need more satellites, boss")

		best_distance := 99999999.9999999
		var probe_ship_shipyard_waypoint_symbol string
		for _, shipyard := range probe_shipyards {
			distance := distance_between_two_coordinates(shipyard.X, shipyard.Y, current_waypoint.X, current_waypoint.Y)
			if distance < best_distance {
				probe_ship_shipyard_waypoint_symbol = shipyard.Symbol
			}
		}

		fmt.Println("[DEBUG] buyer_ship_destination_symbol:")
		fmt.Println(probe_ship_shipyard_waypoint_symbol)

		fmt.Println("[DEBUG] command ship current location")
		fmt.Println(ship.Nav.WaypointSymbol)

		if is_ship_already_at_waypoint(ship, probe_ship_shipyard_waypoint_symbol) {

			if !is_ship_docked(ship) {
				DockShip(ship.Symbol)
			}

			// This will only purchase one ship per turn. We can buy more per turn but we need to update the satellite count afterwards
			PurchaseShip("SHIP_PROBE", ship.Nav.WaypointSymbol)

			// TODO: buy satellites upto len(markets_to_cover)
			fmt.Println("[INFO] command ship is at probe_ship_shipyard_waypoint_symbol BUY SATELLITES")

		} else {
			// TODO: send command ship to shipyard which sells satellites
			if is_ship_docked(ship) {
				OrbitShip(ship.Symbol)
			}
			navigate_ship_result := NavigateShip(ship.Symbol, probe_ship_shipyard_waypoint_symbol)
			fmt.Println(navigate_ship_result)
		}
	} else {
		// we have enough satellites
		//fmt.Println("[INFO] We have enough satellites, boss. It's time to start trading!")
		AssignSatellitesToMarkets(markets_to_cover)

		if MarketScanComplete(trade_routes) {
			PopulateTradeRoutesProfitPerUnit(trade_routes)
		}

		PrintTradeRoutes(ship_list, trade_routes)

		most_profitable_trade_route := MostProfitableTradeRoute(trade_routes)

		if is_ship_cargo_empty(ship) {
			fmt.Println("[INFO] Cargo hold empty")
			// This is flimsy because the MostProfitableTradeRoute will change, and if it does so while we have cargo this will malfunction
			if is_ship_already_at_waypoint(ship, most_profitable_trade_route.BuyMarketplaceWaypointSymbol) {
				fmt.Println("[DEBUG] Already at waypoint")
				if !is_ship_docked(ship) {
					DockShip(ship.Symbol)
				}

				RefuelShip(ship.Symbol)

				// BUY STUFF
				fmt.Println("[DEBUG] This is where we buy stuff")

				maximum_affordable_units := HowManyTradeGoodCanIAfford(GetAgent(), most_profitable_trade_route.BuyMarketTradeGood)
				units_to_purchase := maximum_affordable_units
				space_in_cargo_hold := ship.Cargo.Capacity - ship.Cargo.Units

				if space_in_cargo_hold < maximum_affordable_units {
					units_to_purchase = space_in_cargo_hold
				}

				buy_market_trade_volume := most_profitable_trade_route.BuyMarketTradeGood.TradeVolume
				if buy_market_trade_volume < units_to_purchase {
					number_of_purchases_required := units_to_purchase / buy_market_trade_volume
					units_to_purchase = most_profitable_trade_route.BuyMarketTradeGood.TradeVolume

					for i := 1; i > int(number_of_purchases_required); i++ {
						purchase_cargo_result := PurchaseCargo(ship.Symbol, most_profitable_trade_route.TradeGoodSymbol, units_to_purchase)
						if purchase_cargo_result.Transaction.Units < buy_market_trade_volume {
							units_to_purchase = buy_market_trade_volume
						}
					}
				}

				fmt.Print("[DEBUG] units_to_purchase = ")
				fmt.Print(units_to_purchase)
				fmt.Println()

				OrbitShip(ship.Symbol)
				NavigateShip(ship.Symbol, most_profitable_trade_route.SellMarketplaceWaypointSymbol)
				return
			}

			// TODO: check if all markets are covered here and return early if not

			if MarketScanComplete(trade_routes) {
				fmt.Println()
				fmt.Println("[INFO] Heading to buy marketplace")
				if is_ship_docked(ship) {
					RefuelShip(ship.Symbol)
					OrbitShip(ship.Symbol)
				}
				NavigateShip(ship.Symbol, most_profitable_trade_route.BuyMarketplaceWaypointSymbol)

			}
		} else {
			fmt.Println("[DEBUG] Cargo not empty")

			first_item_in_inventory := ship.Cargo.Inventory[0]

			trade_routes_with_inventory_good := TradeRoutesWithTradeGood(trade_routes, first_item_in_inventory.Symbol)

			fmt.Print("[DEBUG] first_item_in_inventory.Symbol ")
			fmt.Println(first_item_in_inventory.Symbol)

			fmt.Print("[DEBUG] trade_routes_with_inventory_good length ")
			fmt.Println(len(trade_routes_with_inventory_good))

			most_profitable_trade_route_with_inventory_good := MostProfitableTradeRoute(trade_routes_with_inventory_good)

			fmt.Print("[DEBUG] most_profitable_trade_route_with_inventory_good")
			fmt.Println(most_profitable_trade_route_with_inventory_good)

			if is_ship_already_at_waypoint(ship, most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol) {
				fmt.Println("[DEBUG] Already at sell marketplace")
				if !is_ship_docked(ship) {
					DockShip(ship.Symbol)
				}
				// this doesnt account for TradeVolume < Cargo.Capacity
				units_in_cargo_hold := CountTradeGoodCargo(ship, most_profitable_trade_route_with_inventory_good.TradeGoodSymbol)
				sell_market_trade_volume := most_profitable_trade_route_with_inventory_good.SellMarketTradeGood.TradeVolume
				if units_in_cargo_hold > sell_market_trade_volume {
					number_of_sales_required := units_in_cargo_hold / sell_market_trade_volume

					units_to_sell := sell_market_trade_volume
					for i := 1; i > int(number_of_sales_required); i++ {
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
				if is_ship_docked(ship) {
					RefuelShip(ship.Symbol)
					OrbitShip(ship.Symbol)
				}
				NavigateShip(ship.Symbol, most_profitable_trade_route_with_inventory_good.SellMarketplaceWaypointSymbol)
			}
		}
	}
}

func AssignSatellitesToMarkets(markets_to_cover map[string]string) {

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

func ApplyRoleSatellite(ship Ship, markets_to_cover map[string]string, trade_routes []TradeRoute) {
	fmt.Println("[INFO] " + ship.Symbol)

	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[DEBUG] IN_TRANSIT TO")
		fmt.Println(ship.Nav.Route.Destination)
		fmt.Println("[DEBUG] Arrival")
		fmt.Println(ship.Nav.Route.Arrival)
		return
	}

	var assigned_market_waypoint string

	// find the name of this satellite as a value in the markets_to_cover map, return the key of that value as assigned_market_waypoint
	for market_symbol, assigned_satellite := range markets_to_cover {
		if ship.Symbol == assigned_satellite {
			assigned_market_waypoint = market_symbol
			break
		}
	}

	if is_ship_already_at_waypoint(ship, assigned_market_waypoint) {
		fmt.Println("[DEBUG] Already at assigned market waypoint")
		if !is_ship_docked(ship) {
			DockShip(ship.Symbol)
		}
		UpdateTradeRoutesIncludingThisWaypoint(assigned_market_waypoint, trade_routes)
	} else {
		fmt.Println("[INFO] Not at assigned market waypoint, heading there now")
		if is_ship_docked(ship) {
			OrbitShip(ship.Symbol)
		}
		NavigateShip(ship.Symbol, assigned_market_waypoint)
	}

	// find my assignment waypoint
	// am i there?
	// dock and get market
	// if no orbit and navigate there

}

func ShipRoleDecider(ship Ship, markets_to_cover map[string]string, probe_shipyards []Waypoint, trade_routes []TradeRoute) {
	if ship.Registration.Role == "COMMAND" {
		ApplyRoleCommand(ship, markets_to_cover, probe_shipyards, trade_routes)
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
		all_market_results = append(all_market_results, get_market_result)
	}

	// association for places to BUY and SELL TradeGoods
	trade_routes := []TradeRoute{}

	// each unique market waypoint symbol (unordered)
	markets_to_cover := make(map[string]string)

	// populate both trade_routes and markets_to_cover by iterating through A) every MARKETPLACE B) each of their Imports and C) their Exports
	// associations of import/export are added to trade_routes, and we keep one copy of each waypoint_symbol in markets_to_cover
	for _, each_market_result := range all_market_results {
		if len(each_market_result.Exports) > 0 {
			for _, each_export := range each_market_result.Exports {
				for _, each_market_result_inner := range all_market_results {
					if len(each_market_result_inner.Imports) > 0 {
						for _, each_market_result_imports := range each_market_result_inner.Imports {
							if each_export.Symbol == each_market_result_imports.Symbol {

								fmt.Print("[DEBUG] TRADE ROUTE FOUND BUY ")
								fmt.Print(each_export.Symbol)
								fmt.Print(" AT ")
								fmt.Print(each_market_result.Symbol)
								fmt.Print(" SELL AT ")
								fmt.Print(each_market_result_inner.Symbol)
								fmt.Println()

								trade_route := TradeRoute{}
								trade_route.TradeGoodSymbol = each_export.Symbol
								trade_route.BuyMarketplaceWaypointSymbol = each_market_result.Symbol
								trade_route.SellMarketplaceWaypointSymbol = each_market_result_inner.Symbol
								trade_routes = append(trade_routes, trade_route)
								markets_to_cover[trade_route.BuyMarketplaceWaypointSymbol] = ""
								markets_to_cover[trade_route.SellMarketplaceWaypointSymbol] = ""
							}
						}
					}
				}
			}
		}
	}

	PopulateTradeRoutesWithWaypointData(trade_routes, markets_to_cover)
	PopulateTradeRoutesWithDistances(trade_routes)

	// there can be multiple SHIPYARDs which sell SHIP_PROBE
	probe_shipyards := []Waypoint{}

	// populate probe_shipyards with Waypoints which have SHIPYARDs which sell SHIP_PROBEs
	shipyards_in_system := list_waypoints_in_system_by_trait(base_system_symbol, "SHIPYARD")
	for _, shipyard_waypoint := range shipyards_in_system {
		get_shipyard_result := GetShipyard(base_system_symbol, shipyard_waypoint.Symbol)
		for _, ship := range get_shipyard_result.ShipTypes {
			if ship.Type == "SHIP_PROBE" {
				fmt.Println("[INFO] shipyard with satellites for sale found: ")
				fmt.Println("[INFO] " + get_shipyard_result.Symbol)
				probe_shipyards = append(probe_shipyards, shipyard_waypoint)
			}
		}
	}

	fmt.Println("[DEBUG] markets to cover:")
	fmt.Println(markets_to_cover)

	for market := range markets_to_cover {
		fmt.Println(market)
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
			ShipRoleDecider(ship, markets_to_cover, probe_shipyards, trade_routes)

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
