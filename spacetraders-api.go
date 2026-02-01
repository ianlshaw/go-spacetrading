package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Agents

func RegisterAgent(callsign string) (result RegisterAgentResponse) {
	fmt.Println("RegisterAgent")
	payload := &RegisterAgentPayload{}
	payload.Faction = "COSMIC"
	payload.Symbol = callsign
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost("register", payloadJSON)
	pretty_print_json(response_string)
	data_container := RegisterAgentResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] RegisterAgent failed to unmarshal")
	}
	token := data_container.Data.Token
	auth_token := token
	WriteAuthTokenToFile(auth_token, callsign+".token")
	return data_container.Data
}

func GetAgent() Agent {
	endpoint := "my/agent"
	response_string := BasicGet(endpoint)

	data_container := GetAgentResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetAgent failed to unmarshal")
	}

	return data_container.Data
}

// Contracts
func ListContracts(page string, limit string) (contracts []Contract, meta Meta) {
	endpoint := "my/contracts?page=" + page + "&limit=" + limit
	response_string := BasicGet(endpoint)
	data_container := ListContractsResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] ListContracts failed to unmarshal")
		//fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data, data_container.Meta
}

func NegotiateContract(ship_symbol string) Contract {
	fmt.Println("[DEBUG] NegotiateContract")
	endpoint := "my/ships/" + ship_symbol + "/negotiate/contract"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)

	data_container := NegotiateContractResponseData{}

	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] NegotiateContract failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}

	return data_container.Data
}

func AcceptContract(contract_id string) Contract {
	fmt.Println("[INFO] AcceptContract")
	endpoint := "my/contracts/" + contract_id + "/accept"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := AcceptContractResponseData{}

	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] AcceptContract failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}

	return data_container.Data
}

// Fleet
func ListShips() (ships []Ship) {
	//fmt.Println("[DEBUG] list_ships")
	endpoint := "my/ships"
	response_string := BasicGet(endpoint)
	data_container := ListShipsResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] ListShips failed to unmarshal")
		//fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func GetShip(ship_symbol string) Ship {
	endpoint := "my/ships/" + ship_symbol
	response_string := BasicGet(endpoint)
	data_container := GetShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetShip failed to unmarshal")
	}
	return data_container.Data
}

func NavigateShip(ship_symbol string, waypoint_symbol string) (NavigateShipResponse, time.Time) {
	Log("INFO", "NavigateShip " + ship_symbol + " " + waypoint_symbol)
	fmt.Println("[INFO] NavigateShip " + ship_symbol + " " + waypoint_symbol)
	endpoint := "my/ships/" + ship_symbol + "/navigate"
	payload := &NavigateShipPayload{}
	payload.WaypointSymbol = waypoint_symbol
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := NavigateShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] NavigateShip failed to unmarshal")
	}
	arrival_time := StringToTimestamp(data_container.Data.Nav.Route.Arrival)
	return data_container.Data, arrival_time
}

func OrbitShip(ship_symbol string) OrbitShipResponse {
	//Log("DEBUG", "OrbitShip " + ship_symbol)
	endpoint := "my/ships/" + ship_symbol + "/orbit"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := OrbitShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] OrbitShip failed to unmarshal")
	}
	return data_container.Data
}

func DockShip(ship_symbol string) DockShipResponse {
	//Log("DEBUG", "DockShip " + ship_symbol)
	endpoint := "my/ships/" + ship_symbol + "/dock"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := DockShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] DockShip failed to unmarshal")
	}
	return data_container.Data
}

func PurchaseShip(ship_type string, waypoint_symbol string) PurchaseShipResponse {
	fmt.Println("[DEBUG] PurchaseShip")
	endpoint := "my/ships"
	payload := &PurchaseShipPayload{}
	payload.WaypointSymbol = waypoint_symbol
	payload.ShipType = ship_type
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := PurchaseShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] PurchaseShip failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func PurchaseCargo(ship_symbol string, trade_good_symbol string, units int64) PurchaseCargoResponse {
	units_as_string := strconv.FormatInt(units, 10)
	Log("DEBUG", ship_symbol + " PurchaseCargo " + units_as_string + " " + trade_good_symbol)
	endpoint := "my/ships/" + ship_symbol + "/purchase"
	payload := &PurchaseCargoPayload{}
	payload.Symbol = trade_good_symbol
	payload.Units = units
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := PurchaseCargoResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] PurchaseCargo failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}

	fmt.Print("[INFO] ")
	fmt.Print(ship_symbol)
	fmt.Print(" Purchased ")
	fmt.Print(data_container.Data.Transaction.Units)
	fmt.Print(" ")
	fmt.Print(trade_good_symbol)
	fmt.Print(" for ")
	fmt.Print(data_container.Data.Transaction.TotalPrice)
	fmt.Print("\n")

	return data_container.Data
}

func SellCargo(ship_symbol string, trade_good_symbol string, units int64) SellCargoResponse {
	fmt.Println("[DEBUG] SellCargo")
	endpoint := "my/ships/" + ship_symbol + "/sell"
	payload := &SellCargoPayload{}
	payload.Symbol = trade_good_symbol
	payload.Units = units
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := SellCargoResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] SellCargo failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	fmt.Print("[INFO] ")
	fmt.Print(ship_symbol)
	fmt.Print(" Sold ")
	fmt.Print(data_container.Data.Transaction.Units)
	fmt.Print(" ")
	fmt.Print(trade_good_symbol)
	fmt.Print(" for ")
	fmt.Print(data_container.Data.Transaction.TotalPrice)
	fmt.Print("\n")

	fmt.Print("[INFO] New balance: ")
	fmt.Print(data_container.Data.Agent.Credits)
	fmt.Print("\n")

	return data_container.Data
}

func RefuelShip(ship_symbol string, units int64, from_cargo bool) RefuelShipResponse {
	//fmt.Print("[DEBUG] RefuelShip ")
	//fmt.Print(ship_symbol + " ")
	//fmt.Print(units)
	//fmt.Print(" ")
	//fmt.Println(from_cargo)
	endpoint := "my/ships/" + ship_symbol + "/refuel"
	payload := &RefuelShipPayload{}
	//payload.Units = units // this is nulled in the model to force "fill completely" behaviour
	payload.FromCargo = from_cargo
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := RefuelShipResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] RefuelShip failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func TransferCargo(source_ship_symbol string, target_ship_symbol string, cargo_symbol string, units int64) TransferCargoResponse {
	Log("DEBUG", "TransferCargo " + source_ship_symbol + " " + target_ship_symbol + " " + cargo_symbol + string(units) )
	endpoint := "my/ships/" + source_ship_symbol + "/transfer"
	payload := &TransferCargoPayload{}
	payload.TradeSymbol = cargo_symbol
	payload.Units = units
	payload.ShipSymbol = target_ship_symbol
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := TransferCargoResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] TransferCargo failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	fmt.Print("[INFO] Transfered ")
	fmt.Print(units)
	fmt.Println(" " + cargo_symbol + " from " + source_ship_symbol + " to " + target_ship_symbol)
	return data_container.Data
}

// Systems
func GetSystem(system_symbol string) (system System) {
	endpoint := "systems/" + system_symbol
	response_string := BasicGet(endpoint)
	data_container := GetSystemResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetSystem failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func GetWaypoint(system_symbol string, waypoint_symbol string) (resultant_waypoint Waypoint) {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol
	response_string := BasicGet(endpoint)
	data_container := GetWaypointResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetWaypoint failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func ListWaypointsInSystem(system_symbol string, offset string) ListWaypointsInSystemResponseData {
	endpoint := "systems/" + system_symbol + "/waypoints?page=" + offset + "&limit=20"
	response_string := BasicGet(endpoint)
	data_container := ListWaypointsInSystemResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] ListWaypointInSystem failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container
}

func ListWaypointInSystemByTrait(system_symbol string, trait string) []Waypoint {
	endpoint := "systems/" + system_symbol + "/waypoints?traits[]=" + trait
	response_string := BasicGet(endpoint)
	data_container := ListWaypointsInSystemResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] ListWaypointInSystemByTrait failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func ListWaypointInSystemByType(system_symbol string, query_type string) []Waypoint {
	endpoint := "systems/" + system_symbol + "/waypoints?type=" + query_type
	response_string := BasicGet(endpoint)
	data_container := ListWaypointsInSystemResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] ListWaypointInSystemByType failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func GetMarket(system_symbol string, waypoint_symbol string) Market {
	fmt.Println("[DEBUG] GetMarket")
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/market"
	response_string := BasicGet(endpoint)
	data_container := GetMarketResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetMarket failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func GetShipyard(system_symbol string, waypoint_symbol string) (get_shipyard_result Shipyard) {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/shipyard"
	response_string := BasicGet(endpoint)
	data_container := GetShipyardResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetShipyard failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func GetJumpGate(system_symbol string, waypoint_symbol string) JumpGate {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/jump-gate"
	response_string := BasicGet(endpoint)
	data_container := GetJumpGateResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetJumpGate failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func GetConstructionSite(system_symbol string, waypoint_symbol string) ConstructionSite {
	endpoint := "systems/" + system_symbol + "/waypoints/" + waypoint_symbol + "/construction"
	response_string := BasicGet(endpoint)
	data_container := GetConstructionSiteResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] GetConstructionSite failed to unmarshal")
	}
	return data_container.Data
}

func SiphonResources(ship_symbol string) SiphonResourcesResponse {
	fmt.Println("[DEBUG] SiphonResources " + ship_symbol)
	endpoint := "my/ships/" + ship_symbol + "/siphon"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := SiphonResourcesResponseData{}

	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] SiphonResources failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	fmt.Print("[INFO] " + ship_symbol + " siphoned ")
	fmt.Print(data_container.Data.Siphon.Yield.Units)
	fmt.Println(" " + data_container.Data.Siphon.Yield.Symbol)	
	return data_container.Data
}

func DeliverCargoToContract(contract_id string, ship_symbol string, trade_symbol string, units int64) DeliverCargoToContractResponse {
	fmt.Println("[DEBUG] DeliverCargoToContract")
	endpoint := "my/contracts/" + contract_id + "/deliver"
	payload := &DeliverCargoToContractPayload{}
	payload.ShipSymbol = ship_symbol
	payload.TradeSymbol = trade_symbol
	payload.Units = units
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)

	data_container := DeliverCargoToContractResponseData{}
	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] DeliverCargoToContract: failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}
	return data_container.Data
}

func FulfillContract(contract_id string) FulfillContractResponse {
	fmt.Println("[DEBUG] FulfillContract")
	endpoint := "my/contracts/" + contract_id + "/fulfill"
	payload := &EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := FulfillContractResponseData{}

	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] FulfillContract failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}

	fmt.Print("[INFO] New balance: ")
	fmt.Println(data_container.Data.Agent.Credits)
	return data_container.Data
}

func JettisonCargo(ship Ship, trade_good_symbol string, units int64) Cargo {
	fmt.Print("[DEBUG] JettisonCargo " + trade_good_symbol + " ")
	fmt.Println(units)
	endpoint := "my/ships/" + ship.Symbol + "/jettison"
	payload := &JettisonCargoPayload{}
	payload.TradeSymbol = trade_good_symbol
	payload.Units = units
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)
	response_string := BasicPost(endpoint, payloadJSON)
	data_container := JettisonCargoResponseData{}

	if err := json.Unmarshal([]byte(response_string), &data_container); err != nil {
		fmt.Println("[ERROR] JettisonCargo failed to unmarshal")
		fmt.Println(response_string)
		fmt.Println(err)
	}

	fmt.Print("[INFO] " + ship.Symbol)
	fmt.Print(" Jettisoned ")
	fmt.Print(units)
	fmt.Println(" " + trade_good_symbol)
	return data_container.Data.Cargo
}