package main

import (
	"fmt"
)

func DecideSaturationSatelliteAction(ship_state *ShipState, world *WorldState) ShipAction {

	ship := &ship_state.Ship

	fmt.Printf("[INFO %s %s %s DecideSaturationSatelliteAction\n", ship.Symbol, ship.Registration.Role, ship.Nav.WaypointSymbol)

	// do i have a target?
	if ship_state.TargetWaypointSymbol == "" {
		for _, market := range world.Markets {
			for _, other_ship := range world.Ships {
				// exclude self
				if ship.Symbol == other_ship.Symbol {
					continue
				}
				// exclude non satellites
				if other_ship.Registration.Role != "SATELLITE" {
					continue
				}
				if other_ship.Nav.WaypointSymbol == market.WaypointSymbol {
					continue
				}
				// if the iterator survives to this point then the market should not have any other satellites at or on the way to it.
				ship_state.TargetWaypointSymbol = market.WaypointSymbol
			}
		}
	}

	if IsShipAlreadyAtWaypoint(ship, ship_state.TargetWaypointSymbol) {

		// Is it a shipyard?
		if world.Waypoints[ship.Nav.WaypointSymbol].WaypointHasTrait("SHIPYARD") {
			if !IsShipDocked(ship) {
				return ShipAction{
					Type:       ActionDock,
					ShipSymbol: ship.Symbol,
				}
			}
			if world.IsShipyardStale(ship.Nav.WaypointSymbol) {
				return ShipAction{
					Type:           ActionUpdateShipyardData,
					WaypointSymbol: ship.Nav.WaypointSymbol,
				}
			}
		}

		// Is it a marketplace?
		if world.Waypoints[ship.Nav.WaypointSymbol].WaypointHasTrait("MARKETPLACE") {
			if world.IsMarketStale(ship.Nav.WaypointSymbol) {
				return ShipAction{
					Type:           ActionUpdateMarketData,
					WaypointSymbol: ship.Nav.WaypointSymbol,
				}
			}
		}

	}

	return ShipAction{
		Type:       ActionNavigate,
		ShipSymbol: ship.Symbol,
	}

	// Is there another satellite already at this location?

	// already at target

	// Am I the only satellite here

	// this would result in multiple satellites marking the same waypoint, there is nothing fanning them out yet.
	// nor any navigation

	// I think this is unreachable
	fmt.Printf("[WARN] DecideSaturationSatelliteAction uncaught branch. Waiting 15 minutes.\n")
	return ShipAction{
		Type:       ActionWait,
		ShipSymbol: ship.Symbol,
		NotBefore:  FifteenMinutesFromNow(),
	}
}
