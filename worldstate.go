package main

import (
	"fmt"
	"time"
)

var World *WorldState

type WorldState struct {
    Markets map[string]*MarketState
	Waypoints map[string]*Waypoint
	Shipyards map[string]*ShipyardState
	Ships map[string]*ShipState
	ConstructionSites map[string]*ConstructionSiteState
}

type ConstructionSiteState struct {
	WaypointSymbol string
	LastSeen time.Time
	ConstructionSite ConstructionSite
}

type MarketState struct {
    WaypointSymbol string
    LastSeen time.Time
    Market Market
}

type ShipyardState struct {
    WaypointSymbol string
    LastSeen time.Time
    Shipyard Shipyard
}

type ShipState struct {
    BusyUntil time.Time
    Ship Ship
    Job ShipJob
}

func (w *WorldState) UpdateFromShip(s Ship) {
    _, ok := w.Ships[s.Symbol]
    if !ok {
        w.Ships[s.Symbol] = &ShipState{
            Ship: s,
        }
    } else {
        w.Ships[s.Symbol].Ship = s
    }
    //w.Ships[s.Symbol] = &ShipState{
    //    Ship: s,
		// TODO ensure this covers both navigation arrival times and mining cooldown durations
        //BusyUntil:  StringToTimestamp(s.Cooldown.Expiration),
    //}
}

func (w *WorldState) UpdateFromMarket(m Market) {
    w.Markets[m.Symbol] = &MarketState{
        WaypointSymbol: m.Symbol,
        LastSeen: time.Now(),
        Market:   m,
    }
}

func (w *WorldState) UpdateFromShipyard(s Shipyard) {
    w.Shipyards[s.Symbol] = &ShipyardState{
        WaypointSymbol: s.Symbol,
        LastSeen: time.Now(),
        Shipyard:   s,
    }
}

func (w *WorldState) AssignShipJob(ship_symbol string, job ShipJob) {
    w.Ships[ship_symbol].Job = job
}

func (w *WorldState) UpdateFromConstructionSite(cs ConstructionSite) {
    w.ConstructionSites[cs.Symbol] = &ConstructionSiteState{
        WaypointSymbol: cs.Symbol,
        LastSeen: time.Now(),
        ConstructionSite:   cs,
    }
}

func (w *WorldState) UpdateFromWaypoint(wp Waypoint) {
    w.Waypoints[wp.Symbol] = &wp
}

func (w *WorldState) IsMarketStale(waypoint string) bool {
    m, ok := w.Markets[waypoint]
    if !ok {
        return true // unknown == stale
    }
    return time.Since(m.LastSeen) > 1*time.Minute
}

func (w *WorldState) IsShipyardStale(waypoint string) bool {
    s, ok := w.Shipyards[waypoint]
    if !ok {
        return true // unknown == stale
    }
    return time.Since(s.LastSeen) > 1*time.Minute
}

func (w *WorldState) InvalidateMarket(waypoint_symbol string) {
    if m, ok := w.Markets[waypoint_symbol]; ok {
        m.LastSeen = time.Time{} // zero time = definitely stale
    } else {
		fmt.Println("[ERROR] Failed to InvalidateMarket " + waypoint_symbol)
	}
}

func (w *WorldState) InvalidateShipyard(waypoint_symbol string) {
    if s, ok := w.Shipyards[waypoint_symbol]; ok {
        s.LastSeen = time.Time{} // zero time = definitely stale
    } else {
		fmt.Println("[ERROR] Failed to InvalidateShipyard " + waypoint_symbol)
	}
}

func OldestMarket(world *WorldState) *MarketState {
    var oldest *MarketState

    for _, m := range world.Markets {
        if oldest == nil ||
           m.LastSeen.IsZero() ||
           m.LastSeen.Before(oldest.LastSeen) {
            oldest = m
        }
    }

    return oldest
}