package main

import (
	"fmt"
)

func DecideSaturationSatelliteAction(ship_ptr *Ship, world *WorldState) ShipAction {

	ship := *ship_ptr

	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + "DecideSaturationSatelliteAction")

	// Is it a shipyard?
	if world.Waypoints["ship.Nav.WaypointSymbol"].WaypointHasTrait("SHIPYARD") {
		if !IsShipDocked(ship) {
			return ShipAction{
				Type: ActionDock,
				ShipSymbol: ship.Symbol,
			}
		}
		if world.IsShipyardStale(ship.Nav.WaypointSymbol){
			return ShipAction{
				Type: ActionUpdateShipyardData,
					WaypointSymbol: ship.Nav.WaypointSymbol,
			}
		}
	}

	// Is it a marketplace?
	if world.Waypoints["ship.Nav.WaypointSymbol"].WaypointHasTrait("MARKETPLACE") {
		if world.IsMarketStale(ship.Nav.WaypointSymbol) {
			return ShipAction{
				Type: ActionUpdateMarketData,
					WaypointSymbol: ship.Nav.WaypointSymbol,
			}
		}
	}

	fmt.Printf("[WARN] DecideSaturationSatelliteAction uncaught branch. Waiting 15 minutes.\n")
	return ShipAction{
		Type: ActionWait,
		ShipSymbol: ship.Symbol,
		NotBefore: FifteenMinutesFromNow(),
	}
}