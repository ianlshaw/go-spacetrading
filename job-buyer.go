package main

import (
	"fmt"
)

// TODO
// Once we have some floating cash
// Purchase enough satellites to cover all marketplaces in system

var desired_number_of_ship_probe = 2
var desired_number_of_ship_light_freighter = 0
//var desired_number_of_ship_shuttle = 1
//var desired_number_of_ship_mining_drone = 1
//var desired_number_of_ship_siphon_drone = 1
//var desired_number_of_ship_surveyor = 1

func DecideBuyerAction(ship_ptr *Ship, world *WorldState) ShipAction {
	ship := *ship_ptr

	if IsShipDocked(ship){
		if world.IsShipyardStale(ship.Nav.WaypointSymbol) {
			return ShipAction{
				Type: ActionUpdateShipyardData,
				ShipSymbol: ship.Symbol,
				WaypointSymbol: ship.Nav.WaypointSymbol,
			}
		}
	}

	#if world.Agent.Credits > 2000000 {
	#	desired_number_of_ship_light_freighter = 1
	#}

	if world.Agent.Credits > 2100000 {
		desired_number_of_ship_probe = len(world.Markets)
	}

	_, probe_shipyard_waypoints := FindPurchaseableShipByType(world, "SHIP_PROBE")

	number_of_ship_probe := CountShipsByFrame(world, "FRAME_PROBE") 
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

	number_of_ship_light_hauler := CountShipsByFrame(world, "FRAME_LIGHT_FREIGHTER")
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
