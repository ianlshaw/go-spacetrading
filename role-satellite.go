package main

import (
	"fmt"
	"math"
)

// TODO
// Either prioritize known trade routes or 
// Use Dikjsra to create an efficient loop or
// Both

func DecideSatelliteAction(ship Ship, world *WorldState) ShipAction {
	fmt.Println("[INFO] " + ship.Symbol + " " + ship.Registration.Role + " " + "DecideSatelliteAction")

	if IsShipInTransit(ship) {
		fmt.Printf("[WARN] %s in transit, waiting one minute...\n", ship.Symbol)
		return ShipAction{
			Type: ActionWait,
			NotBefore: OneMinuteFromNow(),
		}
	}

	// This should only happen once when the ship is first purchased.
	if IsShipDocked(ship) {
		return ShipAction{
			Type: ActionOrbit,
			ShipSymbol: ship.Symbol,
		}
	}

	current_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, ship.Nav.WaypointSymbol)

	markets_yet_to_be_visited := MarketsYetToBeVisited(world)

	target := OldestMarket(world)

	if len(markets_yet_to_be_visited) > 0 {
		fmt.Printf("[INFO] %d markets have not yet been visited\n", len(markets_yet_to_be_visited))
		best_distance := math.MaxInt

		// find closest of markets yet to be visited
		closest_market := Market{}
		for _, market := range markets_yet_to_be_visited {
			market_waypoint := WaypointFromWaypointSymbol(all_waypoints_in_system, market.Symbol)
			distance := DistanceBetweenTwoWaypoints(current_waypoint, market_waypoint)
			if distance < best_distance {
				closest_market = market
				best_distance = distance
			}
		}
		target = world.Markets[closest_market.Symbol]
	}
	
    if target == nil {
		fmt.Println("[ERROR] DecideSatelliteAction target is nil")
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