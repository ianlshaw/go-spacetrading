package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var waypoints_filename = ".waypoints.json"
var shipyards_filename = ".shipyards.json"
var markets_filename = ".markets.json"
var world_state_filename = ".world.json"

func DoesAgentTokenFileExist(callsign string) (result bool) {
	var filename = callsign + ".token"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("[ERROR] Token file does not exist")
		return false
	}
	fmt.Println("[INFO] Token file exists")
	return true
}

func DoesWaypointsFileExist(callsign string) (result bool) {
	var filename = callsign + waypoints_filename
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("[INFO] Waypoints file does not exist")
		return false
	}
	fmt.Println("[INFO] Waypoints file exists")
	return true
}

func WriteWaypointsToFile(waypoints []Waypoint, callsign string) {
	file_content := ""
	f, err := os.Create(callsign + waypoints_filename)
	PanicOnError(err)
	defer f.Close()

	for _, waypoint := range waypoints {
		marshalled_waypoint, err := json.Marshal(waypoint)
		PanicOnError(err)
		file_content = file_content + string(marshalled_waypoint) + "\n"
	}

	write_result, err := f.WriteString(file_content)
	PanicOnError(err)
	fmt.Printf("[DEBUG] WriteWaypointsToFile wrote %d bytes\n", write_result)
}

func ReadWaypointsFromFile(callsign string) []Waypoint {
	fmt.Println("[DEBUG] ReadWaypointsFromFile")
	waypoints := []Waypoint{}
	f, err := os.Open(callsign + waypoints_filename)
	PanicOnError(err)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		waypoint := Waypoint{}
		err := json.Unmarshal([]byte(scanner.Text()), &waypoint)
		PanicOnError(err)
		waypoints = append(waypoints, waypoint)
	}
	return waypoints
}

func WriteAuthTokenToFile(auth_token string, filename string) {
	f, err := os.Create(filename)
	PanicOnError(err)
	defer f.Close()
	write_string_result, err := f.WriteString(auth_token)
	PanicOnError(err)
	fmt.Printf("[DEBUG] WriteAuthTokenToFile wrote %d bytes\n", write_string_result)
}

func ReadAccountTokenFromFile() {
	f, err := os.ReadFile(account_token_filename)
	PanicOnError(err)
	account_token += (string(f))
}

func ReadAgentTokenFromFile(callsign string) {
	f, err := os.ReadFile(callsign + ".token")
	PanicOnError(err)
	agent_token += (string(f))
}

func WorldStateFilename(callsign string) string {
    return callsign + ".world.json"
}

func LoadWorldState(callsign string) *WorldState {
    filename := WorldStateFilename(callsign)

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return &WorldState{
            Markets: make(map[string]*MarketState),
        }
    }

    data, err := os.ReadFile(filename)
    PanicOnError(err)

    var ws WorldState
    PanicOnError(json.Unmarshal(data, &ws))

    if ws.Markets == nil {
		fmt.Println("ws.Markets is nil")
        ws.Markets = make(map[string]*MarketState)
    }

    return &ws
}

func SaveWorldState(callsign string, ws *WorldState) {
    data, err := json.MarshalIndent(ws, "", "  ")
    PanicOnError(err)

    err = os.WriteFile(WorldStateFilename(callsign), data, 0644)
    PanicOnError(err)
}

func DoesShipyardsFileExist(callsign string) (result bool) {
	var filename = callsign + shipyards_filename
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Println("[INFO] Shipyards file does not exist")
		return false
	}
	fmt.Println("[INFO] Shipyards file exists")
	return true
}

func WriteShipyardsToFile(shipyards []Shipyard, callsign string) {
	file_content := ""
	f, err := os.Create(callsign + shipyards_filename)
	PanicOnError(err)
	defer f.Close()

	for _, shipyard := range shipyards {
		marshalled_shipyard, err := json.Marshal(shipyard)
		PanicOnError(err)
		file_content = file_content + string(marshalled_shipyard) + "\n"
	}

	write_result, err := f.WriteString(file_content)
	PanicOnError(err)
	fmt.Printf("[DEBUG] WriteShipyardsToFile wrote %d bytes\n", write_result)
}

func ReadShipyardsFromFile(callsign string) []Shipyard {
	fmt.Println("[DEBUG] ReadShipyardsFromFile")
	shipyards := []Shipyard{}
	f, err := os.Open(callsign + shipyards_filename)
	PanicOnError(err)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		shipyard := Shipyard{}
		err := json.Unmarshal([]byte(scanner.Text()), &shipyard)
		PanicOnError(err)
		shipyards = append(shipyards, shipyard)
	}
	return shipyards
}

func DoesMarketsFileExist(callsign string) (result bool) {
	var filename = callsign + markets_filename
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Println("[INFO] Markets file does not exist")
		return false
	}
	fmt.Println("[INFO] Markets file exists")
	return true
}

func DoesWorldStateFileExist(callsign string) (result bool) {
	var filename = callsign + world_state_filename
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Println("[INFO] World State file does not exist")
		return false
	}
	fmt.Println("[INFO] World State file exists")
	return true
}

func WriteMarketsToFile(markets []Market, callsign string) {
	file_content := ""
	f, err := os.Create(callsign + markets_filename)
	PanicOnError(err)
	defer f.Close()

	for _, market := range markets {
		marshalled_market, err := json.Marshal(market)
		PanicOnError(err)
		file_content = file_content + string(marshalled_market) + "\n"
	}

	write_result, err := f.WriteString(file_content)
	PanicOnError(err)
	fmt.Printf("[DEBUG] WriteMarketsToFile wrote %d bytes\n", write_result)
}

func ReadMarketsFromFile(callsign string) []Market {
	fmt.Println("[DEBUG] ReadMarketsFromFile")
	markets := []Market{}
	f, err := os.Open(callsign + markets_filename)
	PanicOnError(err)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		market := Market{}
		err := json.Unmarshal([]byte(scanner.Text()), &market)
		PanicOnError(err)
		markets = append(markets, market)
	}
	return markets
}
