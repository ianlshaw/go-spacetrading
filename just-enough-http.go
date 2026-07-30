package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func BasicGet(endpoint string) (response_body string) {
	<-apiLimiter.C

	url := url_base + endpoint

	// DEBUG
	//fmt.Println("[DEBUG] " + url)
	// DEBUG

	payload := EmptyPayload{}
	payloadJSON, err := json.Marshal(payload)
	PanicOnError(err)

	//fmt.Println("payload:")
	//fmt.Println(payload)

	request, _ := http.NewRequest("GET", url, bytes.NewBuffer(payloadJSON))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", agent_token)
	result, err := http.DefaultClient.Do(request)

	//fmt.Println("request:")
	//fmt.Println(request)

	PanicOnError(err)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	PanicOnError(err)

	error_container := ErrorResponse{}
	if err := json.Unmarshal(body, &error_container); err != nil {
		fmt.Println("[ERROR] BasicGet failed to unmarshal")
	}

	// 4113 is returned when the server restarts
	if error_container.Error.Code == 4113 {
		fmt.Println("[ERROR] Server has restarted.")
		EraseState(callsign)
		UpdateAgentTokenSecret("AGENT_TOKEN_EXPIRED")
		os.Exit(42)
	}

	// If the error["message"] field exists, the game returned an error.
	if error_container.Error.Message != "" {
		fmt.Println("[ERROR] response contains error key")
		fmt.Println("Error Code:")
		fmt.Println(error_container.Error.Code)
		fmt.Println(error_container.Error.Message)
		fmt.Println(error_container)
		os.Exit(1)
	}

	sb := string(body)
	http_calls++
	return sb
}

func BasicPost(endpoint string, payload []byte) (response_body string) {

	<-apiLimiter.C

	posturl := url_base + endpoint

	// DEBUG
	//fmt.Println("[DEBUG] " + posturl)
	// DEBUG

	request, err := http.NewRequest("POST", posturl, bytes.NewBuffer(payload))
	PanicOnError(err)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	if endpoint == "register" {
		auth_token := account_token
		request.Header.Add("Authorization", auth_token)
	} else {
		auth_token := agent_token
		request.Header.Add("Authorization", auth_token)
	}
	client := &http.Client{}
	result, err := client.Do(request)
	PanicOnError(err)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	PanicOnError(err)

	error_container := ErrorResponse{}
	if err := json.Unmarshal(body, &error_container); err != nil {
		fmt.Println("[ERROR] BasicPost failed to unmarshal")
	}

	// 4113 is returned when the server restarts
	if error_container.Error.Code == 4113 {
		fmt.Println("[ERROR] Server has restarted.")
		EraseState(callsign)
		os.Exit(4113)
	}

	// If the error["message"] field exists, the game returned an error.
	if error_container.Error.Message != "" {
		fmt.Println("[ERROR] response contains error key")
		fmt.Println("Error Code:")
		fmt.Println(error_container.Error.Code)
		fmt.Println(error_container.Error.Message)
		fmt.Println(error_container.Error.Data)
		os.Exit(1)
	}

	//Convert the body to type string
	sb := string(body)
	http_calls++
	return sb
}
