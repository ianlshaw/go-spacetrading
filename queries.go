package main

import (
	"fmt"
	"math"
	"slices"
	"time"
)

func IsShipInTransit(ship Ship) bool {
	return (ship.Nav.Status == "IN_TRANSIT")
}

func IsShipAlreadyAtWaypoint(ship_to_test Ship, waypoint_symbol string) bool {
	return (ship_to_test.Nav.WaypointSymbol == waypoint_symbol)
}

func IsShipDocked(ship Ship) bool {
	return ship.Nav.Status == "DOCKED"
}

func IsShipCargoEmpty(ship Ship) bool {
	return ship.Cargo.Units == 0
}

func IsShipCargoFull(ship Ship) bool {
	return ship.Cargo.Units == ship.Cargo.Capacity
}

func IsASatelliteDockedAtMarketplace(list_ships_result []Ship, waypoint_symbol string) (answer bool) {
	for _, ship := range list_ships_result {
		if ship.Registration.Role == "SATELLITE" {
			if ship.Nav.WaypointSymbol == waypoint_symbol {
				if ship.Nav.Status == "DOCKED" {
					return true
				}
			}
		}
	}
	return false
}

func DistanceBetweenTwoWaypoints(waypoint1 Waypoint, waypoint2 Waypoint) int {
	return DistanceBetweenTwoCoordinates(waypoint1.X, waypoint1.Y, waypoint2.X, waypoint2.Y)
}

func IsWaypointWithinDistanceOfWaypoint(waypoint1 Waypoint, waypoint2 Waypoint, distance int) bool {
	return DistanceBetweenTwoWaypoints(waypoint1, waypoint2) < distance
}

func IsWaypointWithinDistanceOfTwoWaypoints(waypoint_to_test Waypoint, origin_waypoint Waypoint, destination_waypoint Waypoint, max_distance int) bool {
	return IsWaypointWithinDistanceOfWaypoint(waypoint_to_test, origin_waypoint, max_distance) && IsWaypointWithinDistanceOfWaypoint(waypoint_to_test, destination_waypoint, max_distance)
}

func DistanceBetweenTwoCoordinates(waypoint1X int64, waypoint1Y int64, waypoint2X int64, waypoint2Y int64) (resultant_distance int) {
	//fmt.Println("[DEBUG] distance_between_two_coordinates")
	XIntermediate := waypoint1X - waypoint2X
	YIntermediate := waypoint1Y - waypoint2Y
	XSquared := XIntermediate * XIntermediate
	YSquared := YIntermediate * YIntermediate
	XPlusY := XSquared + YSquared
	XPlusYFloat := float64(XPlusY)

	distance_float := math.Sqrt(XPlusYFloat)
	resultant_distance = int(distance_float)
	return resultant_distance
}

func GetWaypointCoordinate(waypoint Waypoint) (waypointX int64, waypointY int64) {
	return waypoint.X, waypoint.Y
}

func CountTradeGoodCargo(ship Ship, trade_good_symbol string) int64 {
	inventory := ship.Cargo.Inventory
	for _, trade_good := range inventory {
		if trade_good_symbol == trade_good.Symbol {
			return trade_good.Units
		}
	}
	return 0
}

func MarketplacesWhichSellTradeGood(markets []Market, trade_good_symbol string) (markets_selling_trade_good []Market) {
	for _, market := range markets {
		exports := market.Exports
		for _, export := range exports {
			if export.Symbol == trade_good_symbol {
				markets_selling_trade_good = append(markets_selling_trade_good, market)
			}
		}
	}
	return markets_selling_trade_good
}

func CountShipsByFrame(world *WorldState, frame string) int {
	count := 0
	for _, ship_state := range world.Ships {
		if ship_state.Ship.Frame.Symbol == frame {
			count++
		}
	}
	return count
}

func GetShipsByFrame(world *WorldState, frame string) []Ship {
	ships := []Ship{}
	for _, ship_state := range world.Ships {
		if ship_state.Ship.Frame.Symbol == frame {
			ships = append(ships, ship_state.Ship)
		}
	}
	return ships
}

func IsFuelFull(ship Ship) bool {
	if ship.Fuel.Current == ship.Fuel.Capacity {
		return true
	}
	return false
}

// TODO replace this with a map of ships indexed by frame symbol.
func FindPurchaseableShipByType(world *WorldState, frame string) ([]Shipyard, []Waypoint) {
	shipyards := []Shipyard{}
	shipyard_waypoints := []Waypoint{}

	for _, shipyard_state := range world.Shipyards {
		for _, ship := range shipyard_state.Shipyard.ShipTypes {
			if ship.Type == frame {
				shipyards = append(shipyards, shipyard_state.Shipyard)
				for _, waypoint := range world.Waypoints {
					if waypoint.Symbol == shipyard_state.Shipyard.Symbol {
						shipyard_waypoints = append(shipyard_waypoints, *waypoint)
					}
				}
			}
		}
	}
	return shipyards, shipyard_waypoints
}

func CountShipsByModule(ship_list []Ship, module_to_check string) int {
	count := 0
	for _, ship := range ship_list {
		for _, ship_module := range ship.Modules {
			if module_to_check == ship_module.Symbol {
				count++
			}
		}
	}
	return count
}

func CountShipsByMount(ship_list []Ship, mount_to_check string) int {
	count := 0
	for _, ship := range ship_list {
		for _, ship_mount := range ship.Mounts {
			if mount_to_check == ship_mount.Symbol {
				count++
			}
		}
	}
	return count
}

func ClosestWaypointFromSliceToWaypoint(waypoint_slice []Waypoint, singular_waypoint Waypoint) Waypoint {
	best_distance := 99999999
	closest_waypoint := Waypoint{}
	for _, slice_waypoint := range waypoint_slice {
		distance := DistanceBetweenTwoCoordinates(slice_waypoint.X, slice_waypoint.Y, singular_waypoint.X, singular_waypoint.Y)
		if distance < best_distance {
			closest_waypoint = slice_waypoint
			best_distance = distance
		}
	}
	return closest_waypoint
}

func WaypointsWithTrait(world *WorldState, trait_to_check string) map[string]*Waypoint {
	waypoints_with_trait := make(map[string]*Waypoint)
	for waypoint_symbol, waypoint := range world.Waypoints {
		for _, trait := range waypoint.Traits {
			if trait.Symbol == trait_to_check {
				waypoints_with_trait[waypoint_symbol] = waypoint
			}
		}
	}
	return waypoints_with_trait
}

func WaypointsOfType(world *WorldState, type_to_check string) map[string]*Waypoint {
	waypoints_of_type := make(map[string]*Waypoint)
	for waypoint_symbol, waypoint := range world.Waypoints {
		if waypoint.Type == type_to_check {
			waypoints_of_type[waypoint_symbol] = waypoint
		}
	}
	return waypoints_of_type
}

// This is depreciated in favour of world.Waypoints[waypoint_symbol]
func WaypointFromWaypointSymbol(waypoint_slice []Waypoint, waypoint_symbol_to_check string) Waypoint {
	fmt.Println("[DEPRECIATED] WaypointFromWaypointSymbol")
	default_waypoint := Waypoint{}
	for _, waypoint := range waypoint_slice {
		if waypoint.Symbol == waypoint_symbol_to_check {
			return waypoint
		}
	}
	return default_waypoint
}

func IsNewContractRequired(contracts []Contract) bool {
	// if we do not have any contracts
	if len(contracts) == 0 {
		fmt.Println("[DEBUG] we have 0 contracts")
		return true
	}

	// or all contracts are fulfilled
	number_of_unfulfilled_contracts := 0
	for _, contract := range contracts {
		if !IsContractFulfilled(contract) {
			//fmt.Println("[DEBUG] unfulfilled contract found")
			number_of_unfulfilled_contracts++
		}
	}

	if number_of_unfulfilled_contracts == 0 {
		return true
	}

	//fmt.Println("[DEBUG] new contract is not required")
	return false
}

func IsContractAccepted(contract Contract) bool {
	return contract.Accepted
}

func IsContractFulfilled(contract Contract) bool {
	return contract.Fulfilled
}

func CanContractBeCompleted(contract Contract) bool {
	return contract.Terms.Deliver[0].UnitsFulfilled >= contract.Terms.Deliver[0].UnitsRequired
}

func ClosestMarketSellingTradeGood(ship Ship, trade_good string, markets []Market) Market {
	shortest_distance := 9999
	closest_market := Market{}
	ship_waypoint := GetWaypoint(base_system_symbol, ship.Nav.WaypointSymbol)
	for _, market := range markets {
		market_waypoint := GetWaypoint(base_system_symbol, market.Symbol)
		distance := DistanceBetweenTwoWaypoints(ship_waypoint, market_waypoint)
		if distance < shortest_distance {
			closest_market = market
		}
	}
	return closest_market
}

// TODO both all_waypoints and markets should be replaced with world state equivilents.
func ClosestMarketToWaypoint(target_waypoint Waypoint, all_waypoints []Waypoint, markets []Market) Market {
	fmt.Println("[WARN] ClosestMarketToWaypoint needs rewrite see TODO")
	shortest_distance := 9001
	closest_market := Market{}
	for _, market := range markets {
		market_waypoint := WaypointFromWaypointSymbol(all_waypoints, market.Symbol)
		distance := DistanceBetweenTwoWaypoints(target_waypoint, market_waypoint)
		if distance < shortest_distance {
			shortest_distance = distance
			closest_market = market
		}
	}
	fmt.Print("[DEBUG] Closest market to " + target_waypoint.Symbol + " is ")
	fmt.Print(closest_market.Symbol + "  which is ")
	fmt.Print(shortest_distance)
	fmt.Println(" away from it.")
	return closest_market
}

func TradeGoodFromMarket(trade_good_symbol string, market Market) (bool, TradeGood) {
	//fmt.Println(trade_good_symbol + " " + market.Symbol)
	default_trade_good := TradeGood{}
	for _, trade_good := range market.TradeGoods {
		if trade_good.Symbol == trade_good_symbol {
			return true, trade_good
		}
	}
	//fmt.Println("[ERROR] TradeGoodFromMarket market does not contain trade good")
	//fmt.Println(trade_good_symbol)
	//fmt.Println(market.TradeGoods)
	//fmt.Println(market)
	return false, default_trade_good
}

func ContractRemainingRequired(contract Contract) int64 {
	units_required := contract.Terms.Deliver[0].UnitsRequired
	units_fulfilled := contract.Terms.Deliver[0].UnitsFulfilled
	return units_required - units_fulfilled
}

func IsShipAtFactionWaypoint(ship Ship, waypoints []Waypoint, faction string) bool {
	waypoint := WaypointFromWaypointSymbol(waypoints, ship.Nav.WaypointSymbol)
	return waypoint.Faction.Symbol == faction
}

func ClosestFactionWaypointToShip(ship Ship, waypoints []Waypoint) Waypoint {
	faction := ship.Registration.FactionSymbol
	faction_waypoints := make([]Waypoint, 0)
	for _, waypoint := range waypoints {
		if waypoint.Faction.Symbol == faction {
			faction_waypoints = append(faction_waypoints, waypoint)
		}
	}
	ship_waypoint := WaypointFromWaypointSymbol(waypoints, ship.Nav.WaypointSymbol)
	return ClosestWaypointFromSliceToWaypoint(faction_waypoints, ship_waypoint)
}

func ActiveContract(contracts []Contract) Contract {
	default_contract := Contract{}
	for _, contract := range contracts {
		if contract.Fulfilled == false {
			return contract
		}
	}
	fmt.Println("[ERROR] returning null contract - this should never happen - Please call Theo on 1800-bug")
	return default_contract
}

func StringToTimestamp(input_string string) time.Time {
	t, err := time.Parse(time.RFC3339, input_string)
	if err != nil {
		fmt.Println(err)
	}
	return (t)
}

func IsWaypointUnderConstruction(waypoint Waypoint) bool {
	return waypoint.IsUnderConstruction
}

func IsMaterialFulfilled(material Material) bool {
	return material.Fulfilled == material.Required
}

func IsContractDeliverble(contract Contract, all_markets_in_system []Market, mineable_goods []string, siphonable_goods []string) bool {
	contract_delivery_trade_good_symbol := contract.Terms.Deliver[0].TradeSymbol
	markets_with_contract_trade_good := MarketplacesWhichSellTradeGood(all_markets_in_system, contract_delivery_trade_good_symbol)
	if len(markets_with_contract_trade_good) > 0 {
		return true
	}
	if slices.Contains(mineable_goods, contract_delivery_trade_good_symbol) {
		return true
	}
	if slices.Contains(siphonable_goods, contract_delivery_trade_good_symbol) {
		return true
	}
	return false
}

func HaveAtLeastOneBuyerShip(world *WorldState) bool {
	for _, ship_state := range world.Ships {
		if ship_state.Job == JobBuyer {
			return true
		}
	}
	return false
}

func HaveAtLeastOneMarketBoostrap(world *WorldState) bool {
	for _, ship_state := range world.Ships {
		if ship_state.Job == JobMarketBootstrap {
			return true
		}
	}
	return false
}

func UnassignedShipOfRole(world *WorldState, role string) (bool, string) {
	for _, ship_state := range world.Ships {
		if ship_state.Ship.Registration.Role == role {
			if ship_state.Job == "" {
				return true, ship_state.Ship.Symbol
			}
		}
	}
	return false, ""
}

func IsShipStateJobUnassigned(ship_state *ShipState) bool {
	if ship_state.Job == "" {
		return true
	}
	return false
}

func (w *Waypoint) WaypointHasTrait(trait_to_check string) bool {
	for _, wp_trait := range w.Traits {
		if wp_trait.Symbol == trait_to_check {
			return true
		}
	}
	return false
}

func (w *WorldState) CountShipsAtWaypointByFrame(waypoint *Waypoint, frame string) int {
	count := 0
	for _, ship_state := range w.Ships {
		if ship_state.Ship.Nav.WaypointSymbol == waypoint.Symbol {
			if ship_state.Ship.Frame.Symbol == frame {
				count++
			}
		}
	}
	return count
}
