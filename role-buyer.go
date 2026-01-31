package main

import (
	"fmt"
	"time"
)

var desired_number_of_ship_shuttle = 1
var desired_number_of_ship_mining_drone = 1
var desired_number_of_ship_siphon_drone = 1
var desired_number_of_ship_surveyor = 1

func ApplyRoleBuyer(
	ship Ship,
	ship_list []Ship,
	markets_to_cover map[string]string,
	probe_shipard_waypoints []Waypoint,
	shuttle_shipyard_waypoints []Waypoint,
	mining_ship_shipyard_waypoints []Waypoint,
	siphon_ship_shipyard_waypoints []Waypoint,
	survey_ship_shipyard_waypoints []Waypoint,
	agent Agent) time.Time {

	fmt.Println("[DEBUG] ApplyRoleBuyer")



	if ship.Nav.Status == "IN_TRANSIT" {
		fmt.Println("[DEBUG] IN_TRANSIT TO " + ship.Nav.Route.Destination.Symbol)
		fmt.Println("[DEBUG] Arrival " + ship.Nav.Route.Arrival)
		fmt.Println("[ERROR] ApplyRoleBuyer Nav status IN_TRANSIT - THIS SHOULD NOT HAPPEN")
		return time.Now().Add(5 * time.Minute)
	}

	if agent.Credits < 500000 {
		fmt.Println("[INFO] Buyer not enough credits. Sleeping...")
		return time.Now().Add(3 * time.Hour)
	}

	current_waypoint := GetWaypoint(base_system_symbol, ship.Nav.WaypointSymbol)
	var closest_shipyard_waypoint Waypoint
	number_of_ship_shuttle := CountShipsByFrame(ship_list, "FRAME_SHUTTLE")
	if number_of_ship_shuttle < desired_number_of_ship_shuttle {

		fmt.Println("number_of_ship_shuttle")
		fmt.Println(number_of_ship_shuttle)
		fmt.Println("desired_number_of_ship_shuttle")
		fmt.Println(desired_number_of_ship_shuttle)

		// Buy shuttle
		fmt.Println("BUY SHUTTLE")

		closest_shipyard_waypoint = ClosestWaypointFromSliceToWaypoint(shuttle_shipyard_waypoints, current_waypoint)

		if IsShipAlreadyAtWaypoint(ship, closest_shipyard_waypoint.Symbol) {
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}
			PurchaseShip("SHIP_LIGHT_SHUTTLE", closest_shipyard_waypoint.Symbol)
		} else {
			if IsShipDocked(ship) {
				OrbitShip(ship.Symbol)
			}
			_, arrival_time := NavigateShip(ship.Symbol, closest_shipyard_waypoint.Symbol)
			return arrival_time
		}
	}

	//
	return time.Now().Add(1 * time.Minute)
	//

	number_of_ship_mining_drone := CountShipsByMount(ship_list, "MOUNT_MINING_LASER_I")

	if number_of_ship_mining_drone < desired_number_of_ship_mining_drone {

		fmt.Println("number_of_ship_mining_drone")
		fmt.Println(number_of_ship_mining_drone)
		fmt.Println("desired_number_of_ship_mining_drone")
		fmt.Println(desired_number_of_ship_mining_drone)

		// Buy mining drone
		fmt.Println("BUY MINING DRONE")

		closest_shipyard_waypoint = ClosestWaypointFromSliceToWaypoint(mining_ship_shipyard_waypoints, current_waypoint)

		if IsShipAlreadyAtWaypoint(ship, closest_shipyard_waypoint.Symbol) {
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}
			PurchaseShip("SHIP_MINING_DRONE", closest_shipyard_waypoint.Symbol)
		} else {
			if IsShipDocked(ship) {
				OrbitShip(ship.Symbol)
			}
			_, arrival_time := NavigateShip(ship.Symbol, closest_shipyard_waypoint.Symbol)
			return arrival_time
		}
	}

	number_of_ship_surveyor := CountShipsByMount(ship_list, "MOUNT_SURVEYOR_I")
	if number_of_ship_surveyor < desired_number_of_ship_surveyor {
		fmt.Println("number_of_ship_surveyor")
		fmt.Println(number_of_ship_surveyor)
		fmt.Println("desired_number_of_ship_surveyor")
		fmt.Println(desired_number_of_ship_surveyor)
		closest_shipyard_waypoint = ClosestWaypointFromSliceToWaypoint(survey_ship_shipyard_waypoints, current_waypoint)

		if IsShipAlreadyAtWaypoint(ship, closest_shipyard_waypoint.Symbol) {
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}
			PurchaseShip("SHIP_SURVEYOR", closest_shipyard_waypoint.Symbol)
		} else {
			if IsShipDocked(ship) {
				OrbitShip(ship.Symbol)
			}
			_, arrival_time := NavigateShip(ship.Symbol, closest_shipyard_waypoint.Symbol)
			return arrival_time
		}
	}

	number_of_ship_siphon_drone := CountShipsByMount(ship_list, "MOUNT_GAS_SIPHON_I")
	if number_of_ship_siphon_drone < desired_number_of_ship_siphon_drone {
		fmt.Println("number_of_ship_siphon_drone")
		fmt.Println(number_of_ship_siphon_drone)
		fmt.Println("desired_number_of_ship_siphon_drone")
		fmt.Println(desired_number_of_ship_siphon_drone)
		closest_shipyard_waypoint = ClosestWaypointFromSliceToWaypoint(siphon_ship_shipyard_waypoints, current_waypoint)

		if IsShipAlreadyAtWaypoint(ship, closest_shipyard_waypoint.Symbol) {
			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}
			PurchaseShip("SHIP_SIPHON_DRONE", closest_shipyard_waypoint.Symbol)
		} else {
			if IsShipDocked(ship) {
				OrbitShip(ship.Symbol)
			}
			_, arrival_time := NavigateShip(ship.Symbol, closest_shipyard_waypoint.Symbol)
			return arrival_time
		}
	}

	// We dont nessesarily want a limit on this unlike the others
	number_of_ship_probe := CountShipsByFrame(ship_list, "SHIP_PROBE")
	credits := agent.Credits
	if credits < 100000 {
		fmt.Println("[INFO] Not enough money to buy more satellites")
		return time.Now().Add(1 * time.Hour)
	}

	if number_of_ship_probe < len(markets_to_cover) {
		fmt.Println("[INFO] We need more satellites, boss")
		fmt.Print("[DEBUG] number_of_satellites = ")
		fmt.Println(number_of_ship_probe)
		fmt.Print("[DEBUG] markets_to_cover length = ")
		fmt.Println(len(markets_to_cover))

		best_distance := 99999999.9999999
		var probe_ship_shipyard_waypoint_symbol string
		for _, shipyard := range probe_shipard_waypoints {
			distance := DistanceBetweenTwoCoordinates(shipyard.X, shipyard.Y, current_waypoint.X, current_waypoint.Y)
			if distance < int(best_distance) {
				probe_ship_shipyard_waypoint_symbol = shipyard.Symbol
			}
		}

		fmt.Println("[DEBUG] buyer_ship_destination_symbol:")
		fmt.Println(probe_ship_shipyard_waypoint_symbol)

		fmt.Println("[DEBUG] command ship current location")
		fmt.Println(ship.Nav.WaypointSymbol)

		if IsShipAlreadyAtWaypoint(ship, probe_ship_shipyard_waypoint_symbol) {

			if !IsShipDocked(ship) {
				DockShip(ship.Symbol)
			}

			// This will only purchase one ship per turn. We can buy more per turn but we need to update the satellite count afterwards
			PurchaseShip("SHIP_PROBE", ship.Nav.WaypointSymbol)

			fmt.Println("[INFO] command ship is at probe_ship_shipyard_waypoint_symbol BUY SATELLITES")

		} else {
			if IsShipDocked(ship) {
				OrbitShip(ship.Symbol)
			}
			fmt.Println("[INFO] " + ship.Symbol + " Heading to probe shipyard")
			NavigateShip(ship.Symbol, probe_ship_shipyard_waypoint_symbol)
		}
	}
	fmt.Print("[ERROR] ApplyRoleBuyer" + ship.Symbol)
	fmt.Print(" uncaught branch, returning default 5 minute delay")
	return time.Now().Add(5 * time.Minute)
}
