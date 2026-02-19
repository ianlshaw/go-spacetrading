package main

import (
	"fmt"
)

// TODO

var desired_number_of_ship_probe = 2
var desired_number_of_ship_light_freighter = 1
var desired_number_of_ship_shuttle = 1
var desired_number_of_ship_mining_drone = 1
var desired_number_of_ship_siphon_drone = 1
var desired_number_of_ship_surveyor = 1

func DecideBuyerAction(ship Ship, world *WorldState) ShipAction {

	_, probe_shipyard_waypoints := FindPurchaseableShipByType(world, "SHIP_PROBE")

	number_of_ship_probe := CountShipsByFrame(ship_list, "FRAME_PROBE") // This would need to -1 since the buyer is now a probe
	if number_of_ship_probe < desired_number_of_ship_probe {
		if !IsShipAlreadyAtWaypoint(ship, probe_shipyard_waypoints[0].Symbol) {
			if IsShipDocked(ship) {
				return ShipAction{
					Type: ActionOrbit,
					ShipSymbol: ship.Symbol,
				}
			} else {
				return ShipAction{
					Type: ActionNavigate,
					ShipSymbol: ship.Symbol,
					WaypointSymbol: probe_shipyard_waypoints[0].Symbol,
				}
			}
		} else {
			// already at probe shipyard waypoint
			if !IsShipDocked(ship) {
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			} else {
				return ShipAction{
					Type: ActionPurchaseShip,
					WaypointSymbol: probe_shipyard_waypoints[0].Symbol,
					ShipType: "SHIP_PROBE",
				}
			}
		}
	}

	if agent.Credits < 1000000 {
		fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + ship.Frame.Symbol + " DecideBuyerAction less than 1000000 credits. Waiting.")
		return ShipAction{
			Type: ActionWait,
			NotBefore: FifteenMinutesFromNow(),
		}
	}

	number_of_ship_light_hauler := CountShipsByFrame(ship_list, "FRAME_LIGHT_FREIGHTER")
	_, light_hauler_shipyard_waypoints := FindPurchaseableShipByType(world, "SHIP_LIGHT_HAULER")

	if number_of_ship_light_hauler < desired_number_of_ship_light_freighter {
		if !IsShipAlreadyAtWaypoint(ship, light_hauler_shipyard_waypoints[0].Symbol) {
			if IsShipDocked(ship) {
				return ShipAction{
					Type: ActionOrbit,
					ShipSymbol: ship.Symbol,
				}
			} else {
				return ShipAction{
					Type: ActionNavigate,
					ShipSymbol: ship.Symbol,
					WaypointSymbol: light_hauler_shipyard_waypoints[0].Symbol,
				}
			}
		} else {
			// already at probe shipyard waypoint
			if !IsShipDocked(ship) {
				return ShipAction{
					Type: ActionDock,
					ShipSymbol: ship.Symbol,
				}
			} else {
				return ShipAction{
					Type: ActionPurchaseShip,
					WaypointSymbol: light_hauler_shipyard_waypoints[0].Symbol,
					ShipType: "SHIP_LIGHT_HAULER",
				}
			}
		}
	}


	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + ship.Frame.Symbol + " DecideBuyerAction uncaught branch. Waiting.")
	return ShipAction{
		Type: ActionWait,
		NotBefore: FifteenMinutesFromNow(),
	}
}