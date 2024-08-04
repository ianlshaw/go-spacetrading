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

func main() {

	// Ensure the CALLSIGN is provided as a command line argument
	if len(os.Args) != 2 {
		fmt.Println("go-spacetrade CALLSIGN")
		os.Exit(1)
	}

	CALLSIGN := os.Args[1]

	if !does_auth_file_exist(CALLSIGN) {
		register_agent(CALLSIGN)
	}

	// Check if an auth token file is present for the CALLSIGN provided

	//pretty_print_json(dumb_get(""))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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

type RegisterAgentPayload struct {
	Symbol  string `json:"symbol"`
	Faction string `json:"faction"`
}

type ResponseDataContainer struct {
	Container map[string]interface{}
}

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

func basic_post(endpoint string, payload []byte) (returnedJSON map[string]interface{}) {
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

	sb := string(body)
	mapResult := jsonToMap(sb)
	//fmt.Printf("Map with keys :%+v\n", mapResult)

	return mapResult

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

func register_agent(callsign string) {
	fmt.Println("register_agent")

	payload := &RegisterAgentPayload{}

	payload.Faction = "COSMIC"
	payload.Symbol = callsign

	payloadJSON, err := json.Marshal(payload)
	check(err)

	result := basic_post("register", payloadJSON)

	//data := result["data"]

	token := result["data"].(map[string]interface{})["token"]

	fmt.Println(token)

	auth_token := token.(string)
	write_auth_token_to_file(auth_token, callsign+".token")
}

func write_auth_token_to_file(auth_token string, filename string) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()
	write_string_result, err := f.WriteString(auth_token)
	check(err)
	fmt.Printf("wrote %d bytes\n", write_string_result)
}
