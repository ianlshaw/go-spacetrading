package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

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

func DoesWorldStateFileExist(callsign string) (result bool) {
	var filename = WorldStateFilename(callsign)
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Println("[INFO] World State file does not exist")
		return false
	}
	fmt.Println("[INFO] World State file exists")
	return true
}

func LoadWorldState(callsign string) *WorldState {
    filename := WorldStateFilename(callsign)

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return &WorldState{
            Markets: make(map[string]*MarketState),
			Waypoints: make(map[string]*Waypoint),
			Shipyards: make(map[string]*ShipyardState),
			Ships: make(map[string]*ShipState),
			ConstructionSites: make(map[string]*ConstructionSiteState),
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

    if ws.Waypoints == nil {
		fmt.Println("ws.Waypoints is nil")
        ws.Waypoints = make(map[string]*Waypoint)
    }

    if ws.Shipyards == nil {
		fmt.Println("ws.Shipyards is nil")
        ws.Shipyards = make(map[string]*ShipyardState)
    }

    if ws.Ships == nil {
		fmt.Println("ws.Ships is nil")
        ws.Ships = make(map[string]*ShipState)
    }

    if ws.ConstructionSites == nil {
		fmt.Println("ws.ConstructionSites is nil")
        ws.ConstructionSites = make(map[string]*ConstructionSiteState)
    }

    return &ws
}

func SaveWorldState(callsign string, ws *WorldState) {
    data, err := json.MarshalIndent(ws, "", "  ")
    PanicOnError(err)

    err = os.WriteFile(WorldStateFilename(callsign), data, 0644)
    PanicOnError(err)
}