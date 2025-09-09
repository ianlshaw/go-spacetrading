package main

import (
	"strconv"
	"math"
	"time"
	"fmt"
)

func ListAllContracts() []Contract {
	all_contracts := []Contract{}
	contracts, meta := ListContracts("1", "20")
	total := meta.Total
	page := meta.Page
	total_float := float64(total)
	pages_float := total_float / 20.00
	pages_ceil := math.Ceil(pages_float)
	for _, contract := range contracts {
		all_contracts = append(all_contracts, contract)
	}
	page++
	for ; page <= int64(pages_ceil); page++ {
		page_string := strconv.FormatInt(page, 10)
		contracts, meta = ListContracts(page_string, "20")
		for _, contract := range contracts {
			all_contracts = append(all_contracts, contract)
		}
	}
	return all_contracts
}

func Log(log_level string, message string) {
	fmt.Print("[")
	fmt.Print(log_level)
	fmt.Print("] ")
	fmt.Print(time.Now().Format(time.RFC3339))
	fmt.Print(" ")
	fmt.Println(message)
}