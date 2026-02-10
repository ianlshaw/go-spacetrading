package main

import (
	"fmt"
)

// TODO
// Either prioritize known trade routes or 
// Use Dikjsra to create an efficient loop or
// Both

func DecideSatelliteAction(
	ship Ship,
	world *WorldState,
) ShipAction {
	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + "DecideSatelliteAction")

	if IsShipDocked(ship) {
		return ShipAction{
			Type: ActionOrbit,
			ShipSymbol: ship.Symbol,
		}
	}
	
    target := OldestMarket(world)
    if target == nil {
		return ShipAction{
			Type: ActionWait,
			ShipSymbol: ship.Symbol,
			NotBefore: FifteenMinutesFromNow(),
		}
    }

    // If not at market, go there
    if ship.Nav.WaypointSymbol != target.WaypointSymbol {
        return ShipAction{
			Type: ActionNavigate,
            ShipSymbol:     ship.Symbol,
            WaypointSymbol: target.WaypointSymbol,
        }
    }

    // Already there → scan market
    return ShipAction{
		Type: ActionUpdateMarketData,
        ShipSymbol:     ship.Symbol,
        WaypointSymbol: target.WaypointSymbol,
    }
	
	fmt.Println(" uncaught branch, returning 15 minute delay")
	return ShipAction{
		Type: ActionWait,
		ShipSymbol: ship.Symbol,
		NotBefore: FifteenMinutesFromNow(),
	}
}