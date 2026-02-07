package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

var apiLimiter = time.NewTicker(2000 * time.Millisecond)

var url_base string = "https://api.spacetraders.io/v2/"
var account_token = "Bearer "
var account_token_filename = "ACCOUNTTOKENDONOTEXPOSE.txt"
var agent_token = "Bearer "
var base_system_symbol = ""
var http_calls = 0
var turn_length = 120
var callsign = os.Args[1]

var trade_routes []TradeRoute
var all_waypoints_in_system []Waypoint
var all_markets_in_system []Market
var probe_shipyard_waypoints []Waypoint

var agent Agent // does this need to be global or should it be a pointer
var ship_list []Ship
var runningShips = make(map[string]bool)
var markets_to_cover = make(map[string]string)
var probe_shipard_waypoints []Waypoint
//var shuttle_shipyard_waypoints []Waypoint
//var mining_drone_shipyard_waypoints []Waypoint
//var siphon_drone_shipyard_waypoints []Waypoint
//var surveyor_shipyard_waypoints []Waypoint

type ShipActionType string

const (
    ActionNavigate   		ShipActionType = "NAVIGATE"
    ActionBuy        		ShipActionType = "BUY"
    ActionSell       		ShipActionType = "SELL"
    ActionExtract    		ShipActionType = "EXTRACT"
    ActionWait       		ShipActionType = "WAIT"
	ActionDock		 		ShipActionType = "DOCK"
	ActionOrbit		 		ShipActionType = "ORBIT"
	ActionUpdateMarketData  ShipActionType = "UPDATE MARKET DATA"
	ActionPurchaseCargo 	ShipActionType = "PURCHASE CARGO"
	ActionRefuel			ShipActionType = "REFUEL"
	ActionFollowPath 		ShipActionType = "FOLLOW PATH"
	ActionSellCargo 		ShipActionType = "SELL CARGO"
	ActionPurchaseShip		ShipActionType = "PURCHASE SHIP"
)

type ShipAction struct {
    Type       ShipActionType
    ShipSymbol string

    // Optional fields depending on Type
    TradeGoodSymbol     string
    Units     	   		int64
	WaypointSymbol 		string
	Path		   		[]string
	ShipType	   		string

    // When should this action be executed?
    NotBefore time.Time
}

type WorldState struct {
    Markets map[string]*MarketState
}

type MarketState struct {
    WaypointSymbol string
    LastSeen time.Time

    // Raw API response
    Market Market
}

var World *WorldState

func (w *WorldState) UpdateFromMarket(m Market) {
    w.Markets[m.Symbol] = &MarketState{
        WaypointSymbol: m.Symbol,
        LastSeen: time.Now(),
        Market:   m,
    }
}

func (w *WorldState) IsMarketStale(waypoint string) bool {
    m, ok := w.Markets[waypoint]
    if !ok {
        return true // unknown == stale
    }

    return time.Since(m.LastSeen) > 1*time.Minute
}

func (w *WorldState) InvalidateMarket(waypoint_symbol string) {
    if m, ok := w.Markets[waypoint_symbol]; ok {
        m.LastSeen = time.Time{} // zero time = definitely stale
    } else {
		fmt.Println("[ERROR] Failed to InvalidateMarket " + waypoint_symbol)
	}
}

func ensureShipRunning(ship Ship) {
    if runningShips[ship.Symbol] {
        return
    }

    runningShips[ship.Symbol] = true

    go runShip(ship)
}

func ExecuteAction(action ShipAction, ship *Ship) (time.Time) {
    switch action.Type {

    case ActionWait:
        return action.NotBefore

	case ActionFollowPath:
		resp := FollowPath(ship, action.Path)
		ship.Nav = resp.Nav
		ship.Fuel = resp.Fuel
		return StringToTimestamp(resp.Nav.Route.Arrival)

    case ActionNavigate:
        resp := NavigateShip(action.ShipSymbol, action.WaypointSymbol)
		ship.Nav = resp.Nav
        return StringToTimestamp(resp.Nav.Route.Arrival)

	case ActionDock:
		resp := DockShip(action.ShipSymbol)
		ship.Nav = resp.Nav
		return time.Now()

	case ActionRefuel:
		resp := RefuelShip(action.ShipSymbol, 1, false)
		ship.Fuel = resp.Fuel
		agent = resp.Agent
		return time.Now()

	case ActionOrbit:
		resp := OrbitShip(action.ShipSymbol)
		ship.Nav = resp.Nav
		return time.Now()

	// TODO
	// This calls GetMarket twice. One can be removed once we're fully using World MarketState
	case ActionUpdateMarketData:
		UpdateTradeRoutesIncludingThisWaypoint(action.WaypointSymbol)
		resp := GetMarket(base_system_symbol, action.WaypointSymbol)
		World.UpdateFromMarket(resp)
		SaveWorldState(callsign, World)
		return time.Now()
	
	case ActionPurchaseCargo:
		resp := PurchaseCargo(action.ShipSymbol,
			action.TradeGoodSymbol, 
			action.Units,
		)
		ship.Cargo = resp.Cargo
		agent = resp.Agent
		World.InvalidateMarket(ship.Nav.WaypointSymbol)
		return time.Now()

	case ActionSellCargo:
		resp := SellCargo(action.ShipSymbol, action.TradeGoodSymbol, action.Units)
		ship.Cargo = resp.Cargo
		agent = resp.Agent
		World.InvalidateMarket(ship.Nav.WaypointSymbol)
		return time.Now()
	
	case ActionPurchaseShip:
		resp := PurchaseShip(action.ShipType, action.WaypointSymbol)
		agent = resp.Agent
		ship_list = append(ship_list, resp.Ship)
		ensureShipRunning(resp.Ship)
		return time.Now()
	
	}

    panic("unknown action")
}

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

func runShip(ship Ship){

	for {
		//var expiration time.Time
		expiration := ThreeHoursFromNow()

		// DEBUG

		fmt.Print("[DEBUG] " + ship.Symbol + " " + ship.Registration.Role + " "  + ship.Frame.Symbol + " Fuel [")
		fmt.Print(ship.Fuel.Current)
		fmt.Print("/")
		fmt.Print(ship.Fuel.Capacity)
		fmt.Print("] Cargo [")
		fmt.Print(ship.Cargo.Units)
		fmt.Print("/")
		fmt.Print(ship.Cargo.Capacity)
		fmt.Println("]")

		if ship.Registration.Role == "COMMAND" {
			action := DecideTraderAction(ship, all_waypoints_in_system)
			fmt.Println(action)
			expiration = ExecuteAction(action, &ship)
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
				action := DecideBuyerAction(ship)
				fmt.Println(action)
				expiration = ExecuteAction(action, &ship)
			}
		}
		if len(all_probes) > 1 {
			market_bootstrap_probe := all_probes[1]
			if ship.Symbol == market_bootstrap_probe.Symbol {
				action := DecideSatelliteAction(ship)
				fmt.Println(action)
				expiration = ExecuteAction(action, &ship)
			} else {
				//expiration = ApplyRoleSatellite(ship, trade_routes)
				//action := DecideSatelliteAction(ship, trade_routes)
				//fmt.Println(action)
				//expiration = ExecuteAction(action, &ship)
			}
		}

		if len(all_shuttles) >= 1 {
			if ship.Symbol == all_shuttles[0].Symbol {
				// DEBUG
				action := DecideTraderAction(ship, all_waypoints_in_system)
				fmt.Println(action)
				expiration = ExecuteAction(action, &ship)
				// DEBUG
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

	//World = WorldState{
    //    Markets: make(map[string]*MarketState),
    //}
	
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
		//all_waypoints_in_system := []Waypoint{}
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

	all_waypoints_in_system = ReadWaypointsFromFile(CALLSIGN)

	World = LoadWorldState(CALLSIGN)

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
		PopulateGraphDistancesForWaypointWithMaximum(ShuttleMarketplaceGraph, marketplace_waypoints, waypoint, 300)

		//PopulateGraphDistancesForWaypointWithMaximum(SiphonerMarketplaceGraph, marketplace_waypoints, waypoint, 80)
	}

	// each unique market waypoint symbol (unordered)
	markets_to_cover := make(map[string]string)

	// association for places to BUY and SELL TradeGoods
	//trade_routes := []TradeRoute{}

	if !DoesTradeRouteFileExist(CALLSIGN) {
		fmt.Println("[INFO] Trade route file does not exist. Initializing...")
		trade_routes = IdentifyTradeRoutes(CALLSIGN, markets_to_cover, all_markets_in_system, marketplace_waypoints)
		WriteTradeRoutesToFile(trade_routes, CALLSIGN)
	} else {
		fmt.Println("[INFO] Trade file exists. Reading from file...")
		trade_routes = ReadTradeRoutesFromFile(CALLSIGN)
	}

	markets_to_cover = PopulateMarketsToCover(trade_routes)
	_, probe_shipyard_waypoints = FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_PROBE")
	fmt.Print("[DEBUG] probe shipyards:")

	//_, shuttle_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_LIGHT_SHUTTLE")
	//fmt.Print("[DEBUG] shuttle shipyards:")
//
	//mining_drone_shipyards, mining_drone_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_MINING_DRONE")
	//fmt.Print("[DEBUG] mining_drone shipyards:")
	//fmt.Println(len(mining_drone_shipyards))
//
	//siphon_drone_shipyards, siphon_drone_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_SIPHON_DRONE")
	//fmt.Print("[DEBUG] siphon_drone shipyards:")
	//fmt.Println(len(siphon_drone_shipyards))
//
	//surveyor_shipyards, surveyor_shipyard_waypoints := FindPurcahseableShipByFrame(all_waypoints_in_system, all_shipyards_in_system, "SHIP_SURVEYOR")
	//fmt.Print("[DEBUG] surveyor shipyards:")
	//fmt.Println(len(surveyor_shipyards))

	//turn_number := 1

	fmt.Print("[INFO] http calls: ")
	fmt.Print(http_calls)
	http_calls = 0
	fmt.Println()



	agent = GetAgent()
	fmt.Print("[INFO] ShipCount: ")
	fmt.Print(agent.ShipCount)
	fmt.Println()
	fmt.Print("[INFO] Credits: ")
	fmt.Print(agent.Credits)
	fmt.Println()



	ship_list = ListShips()

	//wait_between_ships := turn_length / len(ship_list)

	for _, ship := range ship_list {
    	ensureShipRunning(ship)
	}

	//for _, ship := range ship_list {
	//	go runShip(
	//		ship,
	//		all_waypoints_in_system,
	//		all_markets_in_system,
	//		markets_to_cover,
	//		shuttle_shipyard_waypoints,
	//		mining_drone_shipyard_waypoints,
	//		siphon_drone_shipyard_waypoints,
	//		surveyor_shipyard_waypoints,
	//		trade_routes,
	//		CALLSIGN)
	//		fmt.Print("[INFO] http calls:")
	//		fmt.Println(http_calls)
	//}

	select {}
}