package main

import (
	"math"
)

func IsShipAlreadyAtWaypoint(ship_to_test Ship, waypoint_symbol string) bool {
	return (ship_to_test.Nav.WaypointSymbol == waypoint_symbol && ship_to_test.Nav.Status != "IN_TRANSIT")
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

func CountShipsByFrame(ship_list []Ship, frame string) int {
	count := 0
	for _, ship := range ship_list {
		if ship.Frame.Symbol == frame {
			count++
		}
	}
	return count
}

func FindPurcahseableShipByFrame(all_waypoints_in_system []Waypoint, all_shipyards_in_system []Shipyard, frame string) ([]Shipyard, []Waypoint) {
	shipyards := []Shipyard{}
	shipyard_waypoints := []Waypoint{}
	for _, shipyard := range all_shipyards_in_system {
		for _, ship := range shipyard.ShipTypes {
			if ship.Type == frame {
				shipyards = append(shipyards, shipyard)
				for _, waypoint := range all_waypoints_in_system {
					if waypoint.Symbol == shipyard.Symbol {
						shipyard_waypoints = append(shipyard_waypoints, waypoint)
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

func WaypointsWithTrait(waypoint_slice []Waypoint, trait_to_check string) []Waypoint {
	waypoints_with_trait := []Waypoint{}
	for _, waypoint := range waypoint_slice {
		for _, trait := range waypoint.Traits {
			if trait.Symbol == trait_to_check {
				waypoints_with_trait = append(waypoints_with_trait, waypoint)
			}
		}
	}
	return waypoints_with_trait
}

func WaypointFromWaypointSymbol(waypoint_slice []Waypoint, waypoint_symbol_to_check string) Waypoint {
	default_waypoint := Waypoint{}
	for _, waypoint := range waypoint_slice {
		if waypoint.Symbol == waypoint_symbol_to_check {
			return waypoint
		}
	}
	return default_waypoint
}

func IsContractNegotiated(contracts []Contract) bool {
	return len(contracts) > 0
}

func IsContractAccepted(contract Contract) bool {
	return contract.Accepted
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