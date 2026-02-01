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

func BuyAsManyTradeGoodAsPossible(ship Ship, trade_good TradeGood, market Market) {
	maximum_affordable_units := HowManyTradeGoodCanIAfford(GetAgent(), trade_good)
	fmt.Print("[DEBUG] maximum_affordable_units = ")
	fmt.Println(maximum_affordable_units)

	units_to_purchase := maximum_affordable_units
	fmt.Print("[DEBUG] units_to_purchase = ")
	fmt.Println(units_to_purchase)

	space_in_cargo_hold := ship.Cargo.Capacity - ship.Cargo.Units
	fmt.Print("[DEBUG] space_in_cargo_hold = ")
	fmt.Println(space_in_cargo_hold)

	if space_in_cargo_hold < units_to_purchase {
		units_to_purchase = space_in_cargo_hold
		fmt.Print("[DEBUG] units_to_purchase = ")
		fmt.Println(space_in_cargo_hold)
	}

	buy_market_trade_volume := trade_good.TradeVolume

	fmt.Print("[DEBUG] buy_market_trade_volume = ")
	fmt.Println(buy_market_trade_volume)

	if space_in_cargo_hold > buy_market_trade_volume {
		fmt.Println("[DEBUG] space_in_cargo_hold > buy_market_trade_volume")

		number_of_purchases_required := float64(space_in_cargo_hold) / float64(buy_market_trade_volume)
		rounded_number_of_purchases_required := math.Ceil(number_of_purchases_required)
		fmt.Println("[DEBUG] rounded_number_of_purchases_required = ")
		fmt.Println(rounded_number_of_purchases_required)

		units_to_purchase = buy_market_trade_volume
		fmt.Println("[DEBUG] units_to_purchase = ")
		fmt.Println(units_to_purchase)

		for i := 0; float64(i) < rounded_number_of_purchases_required; i++ {
			fmt.Println("[DEBUG] units_to_purchase = ")
			fmt.Println(units_to_purchase)

			buy_cargo_result := PurchaseCargo(ship.Symbol, trade_good.Symbol, units_to_purchase)

			space_in_cargo_hold = buy_cargo_result.Cargo.Capacity - buy_cargo_result.Cargo.Units

			fmt.Println("[DEBUG] space_in_cargo_hold = ")
			fmt.Println(space_in_cargo_hold)

			if space_in_cargo_hold < units_to_purchase {
				units_to_purchase = space_in_cargo_hold
			}
		}
	} else {
		PurchaseCargo(ship.Symbol, trade_good.Symbol, units_to_purchase)
	}
}

func BuyX(ship Ship, trade_good TradeGood, market Market, x int64) {
	units_to_purchase 					 := x
	remaining_units_to_purchase 		 := x
	space_in_cargo_hold 	    		 := ship.Cargo.Capacity - ship.Cargo.Units
	buy_market_trade_volume 			 := trade_good.TradeVolume
	number_of_purchases_required 		 := float64(units_to_purchase) / float64(buy_market_trade_volume)
	rounded_number_of_purchases_required := math.Ceil(number_of_purchases_required)

	fmt.Print("number_of_purchases_required ")
	fmt.Println(number_of_purchases_required)
	fmt.Print("rounded_number_of_purchases_required ")
	fmt.Println(rounded_number_of_purchases_required)

	for i := 0; float64(i) < rounded_number_of_purchases_required; i++ {
		fmt.Print("i ")
		fmt.Println(i)
		fmt.Print("rounded_number_of_purchases_required ")
		fmt.Println(rounded_number_of_purchases_required)
    	if remaining_units_to_purchase <= space_in_cargo_hold && remaining_units_to_purchase <= buy_market_trade_volume {
    	    units_to_purchase = remaining_units_to_purchase
    	} else if space_in_cargo_hold <= remaining_units_to_purchase && space_in_cargo_hold <= buy_market_trade_volume {
    	    units_to_purchase = space_in_cargo_hold
    	} else {
    	    units_to_purchase = buy_market_trade_volume
    	}

		fmt.Print("[DEBUG] units_to_purchase = ")
		fmt.Println(units_to_purchase)
		buy_cargo_result := PurchaseCargo(ship.Symbol, trade_good.Symbol, units_to_purchase)
		space_in_cargo_hold = buy_cargo_result.Cargo.Capacity - buy_cargo_result.Cargo.Units
		remaining_units_to_purchase -= buy_cargo_result.Cargo.Units
		fmt.Print("[DEBUG] space_in_cargo_hold = ")
		fmt.Println(space_in_cargo_hold)

		// This shouldn't really be required if this loop makes the correct number of purchases.
		if (space_in_cargo_hold == 0) || (remaining_units_to_purchase == 0) {
			fmt.Println("BuyX returning because space or remaining cargo == 0")
			return
		}
	}
	fmt.Println("BuyX escaped its buy loop")
}

func ThreeHoursFromNow() time.Time {
	return time.Now().Add(3 * time.Hour)
}