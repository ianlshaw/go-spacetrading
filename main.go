package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

var apiLimiter = time.NewTicker(550 * time.Millisecond)

var url_base string = "https://api.spacetraders.io/v2/"
var account_token = "Bearer "

// TODO Remove
var account_token_filename = "ACCOUNTTOKENDONOTEXPOSE.txt"

var agent_token = "Bearer "
var base_system_symbol = ""
var http_calls = 0
var callsign = os.Getenv("CALLSIGN")

// maybe put this into WorldState
var runningShips = make(map[string]bool)

func ensureShipRunning(world *WorldState, ship_state *ShipState) {
	if runningShips[ship_state.Ship.Symbol] {
		return
	}
	World.UpdateFromShip(ship_state.Ship)
	runningShips[ship_state.Ship.Symbol] = true
	go runShip(world, ship_state)
}

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
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

type ShipJob string

const (
	JobUnassigned          ShipJob = "UNASSIGNED"
	JobMarketBootstrap     ShipJob = "MARKET_BOOTSTRAP"
	JobScout               ShipJob = "SCOUT"
	JobBuyer               ShipJob = "BUYER"
	JobTrader              ShipJob = "TRADER"
	JobSaturationSatellite ShipJob = "SATURATION SATELLITE"
)

func runShip(world *WorldState, ship_state *ShipState) {

	ship := &ship_state.Ship

	for {
		//var expiration time.Time
		expiration := ThreeHoursFromNow()

		// DEBUG

		fmt.Printf("[DEBUG] %s %s %s Fuel [%d/%d] Cargo [%d/%d]\n",
			ship.Symbol,
			ship.Registration.Role,
			ship.Frame.Symbol,
			ship.Fuel.Current,
			ship.Fuel.Capacity,
			ship.Cargo.Units,
			ship.Cargo.Capacity)

		// This can be set once outside of this loop
		if ship.Registration.Role == "COMMAND" {
			ship_state.Job = JobTrader
			SaveWorldState(callsign, world)
		}

		if ship.Registration.Role == "SATELLITE" {
			if !HaveAtLeastOneBuyerShip(world) {
				fmt.Printf("[INFO] No buyer ships. We need to assign one.\n")
				if IsShipStateJobUnassigned(ship_state) {
					fmt.Printf("[INFO] %s assigned job BUYER\n", ship.Symbol)
					ship_state.Job = JobBuyer
					SaveWorldState(callsign, world)
				}
			}
			if !HaveAtLeastOneMarketBoostrap(world) {
				if IsShipStateJobUnassigned(ship_state) {
					ship_state.Job = JobMarketBootstrap
					SaveWorldState(callsign, world)
				}
			}
		}

		if CountShipsByFrame(world, "FRAME_PROBE") >= len(world.Markets) {
			if ship.Registration.Role == "SATELLITE" {
				ship_state.Job = JobSaturationSatellite
				SaveWorldState(callsign, world)
			}
		}

		if ship.Registration.Role == "HAULER" {
			if IsShipStateJobUnassigned(ship_state) {
				ship_state.Job = JobTrader
				SaveWorldState(callsign, world)
			}
		}

		switch ship_state.Job {
		case JobTrader:
			action := DecideTraderAction(ship, World)
			fmt.Println(action)
			expiration = ExecuteAction(action, ship)

		case JobBuyer:
			action := DecideBuyerAction(ship, World)
			fmt.Println(action)
			expiration = ExecuteAction(action, ship)

		case JobMarketBootstrap:
			action := DecideMarketBootstrapSatellite(ship, World)
			fmt.Println(action)
			expiration = ExecuteAction(action, ship)

		case JobSaturationSatellite:
			action := DecideSaturationSatelliteAction(ship, World)
			fmt.Println(action)
			expiration = ExecuteAction(action, ship)

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

		for {
			sleepFor := time.Until(expiration)
			if sleepFor <= 0 {
				break
			}
			time.Sleep(sleepFor)
			//time.Sleep(sleepFor + 1*time.Second)
		}
	}
}

func main() {

	// Ensure the CALLSIGN is provided as a command line argument
	if len(os.Args) != 1 {
		fmt.Println("go run .")
		os.Exit(1)
	}

	//TODO Remove
	//CALLSIGN := os.Args[1]

	// TODO Remove: depreciating this in favour of aws secrets. RegisterAgent moves to server reset handling?
	// Check if an auth token file is present for the CALLSIGN provided
	//if !DoesAgentTokenFileExist(CALLSIGN) {
	//	ReadAccountTokenFromFile()
	//	RegisterAgent(CALLSIGN)
	//}

	ReadAccountTokenFromEnvironmentVariable()

	if account_token == "Bearer " {
		fmt.Printf("[ERROR] Account token null. Big problem\n")
		os.Exit(1)
	}

	ReadAgentTokenFromEnvironmentVariable()

	if agent_token == "Bearer " {
		fmt.Printf("[WARN] Agent token null. Calling RegisterAgent\n")
		register_agent_result := RegisterAgent(callsign)
		os.Setenv("SPACETRADERS_AGENT_TOKEN", register_agent_result.Token)
		ReadAgentTokenFromEnvironmentVariable()
		UpdateAgentTokenSecret(register_agent_result.Token)
	}

	// TODO remove bearer from this and the above. Dont use functions to create this var, its messy
	if agent_token == "Bearer AGENT_TOKEN_EXPIRED" {
		fmt.Printf("[INFO] Agent token expired. Server must have reset. Regenerating agent token.\n")
		register_agent_result := RegisterAgent(callsign)
		os.Setenv("SPACETRADERS_AGENT_TOKEN", register_agent_result.Token)
		ReadAgentTokenFromEnvironmentVariable()
		UpdateAgentTokenSecret(register_agent_result.Token)
	}

	// TODO Remove
	//ReadAgentTokenFromFile(CALLSIGN)

	// TODO: globals are bad, this should be removed
	populate_base_system_symbol()

	//agent := GetAgent()
	//contracts := ListContracts()

	World = LoadWorldState(callsign)

	if len(World.Waypoints) == 0 {
		fmt.Println("[INFO] Gathering waypoint data...")
		list_waypoints_result := ListWaypointsInSystem(base_system_symbol, "1")
		total_waypoints := list_waypoints_result.Meta.Total
		limit := list_waypoints_result.Meta.Limit
		loop_iterations_required := total_waypoints / int64(limit)
		for i := 1; i < int(loop_iterations_required+2); i++ {
			a_page_of_waypoints := ListWaypointsInSystem(base_system_symbol, strconv.FormatInt(int64(i), 10))
			for _, waypoint := range a_page_of_waypoints.Data {
				World.UpdateFromWaypoint(waypoint)
			}
		}
		SaveWorldState(callsign, World)
	}

	for _, waypoint := range World.Waypoints {
		PopulateGraphDistancesForWaypoint(SystemGraph, World.Waypoints, *waypoint)
		PopulateGraphDistancesForWaypointWithMaximum(SystemGraph, World.Waypoints, *waypoint, 400)
	}

	shipyard_waypoints := WaypointsWithTrait(World, "SHIPYARD")

	if len(World.Shipyards) == 0 {

		for _, shipyard_waypoint := range shipyard_waypoints {
			get_shipyard_result := GetShipyard(base_system_symbol, shipyard_waypoint.Symbol)
			World.UpdateFromShipyard(get_shipyard_result)
		}
		SaveWorldState(callsign, World)
	}

	marketplace_waypoints := WaypointsWithTrait(World, "MARKETPLACE")

	if len(World.Markets) == 0 {
		for marketplace_waypoint_symbol, _ := range marketplace_waypoints {
			get_market_result := GetMarket(base_system_symbol, marketplace_waypoint_symbol)
			World.UpdateFromMarket(get_market_result)
		}
		SaveWorldState(callsign, World)
	}

	construction_site_waypoints := WaypointsOfType(World, "JUMP_GATE")

	if len(World.ConstructionSites) == 0 {
		for construction_site_symbol, _ := range construction_site_waypoints {
			get_construction_site_result := GetConstructionSite(base_system_symbol, construction_site_symbol)
			World.UpdateFromConstructionSite(get_construction_site_result)
		}
		SaveWorldState(callsign, World)
	}

	ship_list := ListShips()
	for _, ship := range ship_list {
		World.UpdateFromShip(ship)
	}
	SaveWorldState(callsign, World)

	for _, waypoint := range marketplace_waypoints {
		PopulateGraphDistancesForWaypointWithMaximum(MarketplaceGraph, marketplace_waypoints, *waypoint, 400)
		PopulateGraphDistancesForWaypointWithMaximum(ShuttleMarketplaceGraph, marketplace_waypoints, *waypoint, 300)
		//PopulateGraphDistancesForWaypointWithMaximum(SiphonerMarketplaceGraph, marketplace_waypoints, waypoint, 80)
	}

	resp := GetAgent()
	World.UpdateFromAgent(resp)

	fmt.Print("[INFO] ShipCount: ")
	fmt.Print(World.Agent.ShipCount)
	fmt.Println()
	fmt.Print("[INFO] Credits: ")
	fmt.Print(World.Agent.Credits)
	fmt.Println()

	for _, ship_state := range World.Ships {
		ensureShipRunning(World, ship_state)
	}

	select {}
}
