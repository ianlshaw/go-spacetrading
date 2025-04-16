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

func DoesTradeRouteFileExist(callsign string) (result bool) {
	var filename = callsign + ".trade_routes"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("[INFO] Trade route file does not exist")
		return false
	}
	fmt.Println("[INFO] Trade route file exists")
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

func WriteTradeRoutesToFile(trade_routes []TradeRoute, callsign string) {
	file_content := ""
	f, err := os.Create(callsign + ".trade_routes")
	PanicOnError(err)
	defer f.Close()

	for _, trade_route := range trade_routes {
		marshalled_trade_route, err := json.Marshal(trade_route)
		PanicOnError(err)
		file_content = file_content + string(marshalled_trade_route) + "\n"
	}

	write_result, err := f.WriteString(file_content)
	PanicOnError(err)
	fmt.Printf("[DEBUG] WriteTradeRoutesToFile wrote %d bytes\n", write_result)
}

func ReadTradeRoutesFromFile(callsign string, trade_routes []TradeRoute) []TradeRoute {
	fmt.Println("[DEBUG] ReadTradeRoutesFromFile")
	f, err := os.Open(callsign + ".trade_routes")
	PanicOnError(err)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		trade_route := TradeRoute{}
		err := json.Unmarshal([]byte(scanner.Text()), &trade_route)
		PanicOnError(err)
		trade_routes = append(trade_routes, trade_route)
	}
	return trade_routes
}
