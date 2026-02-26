package main

import (
	"time"
)

type ShipActionType string

const (
    ActionNavigate   			 ShipActionType = "NAVIGATE"
    ActionBuy        			 ShipActionType = "BUY"
    ActionSell       			 ShipActionType = "SELL"
    ActionExtract    			 ShipActionType = "EXTRACT"
    ActionWait       			 ShipActionType = "WAIT"
	ActionDock		 			 ShipActionType = "DOCK"
	ActionOrbit		 			 ShipActionType = "ORBIT"
	ActionUpdateMarketData  	 ShipActionType = "UPDATE MARKET DATA"
	ActionUpdateShipyardData  	 ShipActionType = "UPDATE SHIPYARD DATA"
	ActionPurchaseCargo 		 ShipActionType = "PURCHASE CARGO"
	ActionRefuel				 ShipActionType = "REFUEL"
	ActionFollowPath 			 ShipActionType = "FOLLOW PATH"
	ActionSellCargo 			 ShipActionType = "SELL CARGO"
	ActionPurchaseShip			 ShipActionType = "PURCHASE SHIP"
	ActionUpdateConstructionSite ShipActionType = "UPDATE CONSTRUCTION SITE"
)

type ShipAction struct {
    Type       ShipActionType
    ShipSymbol string

    // Optional fields depending on Type
    TradeGoodSymbol     string
    Units     	   		int64
	WaypointSymbol 		string
	Path		   		[]string
	ShipType	   		string

    // When should this action be executed?
    NotBefore time.Time
}

func ExecuteAction(action ShipAction, ship *Ship) (time.Time) {
    switch action.Type {

    case ActionWait:
        return action.NotBefore

	case ActionFollowPath:
		resp := FollowPath(ship, action.Path)
		ship.Nav = resp.Nav
		ship.Fuel = resp.Fuel
		return StringToTimestamp(resp.Nav.Route.Arrival)

    case ActionNavigate:
        resp := NavigateShip(action.ShipSymbol, action.WaypointSymbol)
		ship.Nav = resp.Nav
		ship.Fuel = resp.Fuel
        return StringToTimestamp(resp.Nav.Route.Arrival)

	case ActionDock:
		resp := DockShip(action.ShipSymbol)
		ship.Nav = resp.Nav
		return time.Now()

	case ActionRefuel:
		resp := RefuelShip(action.ShipSymbol, 1, false)
		ship.Fuel = resp.Fuel
		World.UpdateFromAgent(resp.Agent)
		return time.Now()

	case ActionOrbit:
		resp := OrbitShip(action.ShipSymbol)
		ship.Nav = resp.Nav
		return time.Now()

	case ActionUpdateMarketData:
		resp := GetMarket(base_system_symbol, action.WaypointSymbol)
		World.UpdateFromMarket(resp)
		SaveWorldState(callsign, World)
		return time.Now()

	case ActionUpdateConstructionSite:
		resp := GetConstructionSite(base_system_symbol, action.WaypointSymbol)
		World.UpdateFromConstructionSite(resp)
		SaveWorldState(callsign, World)
		return time.Now()

	case ActionUpdateShipyardData:
		resp := GetShipyard(base_system_symbol, action.WaypointSymbol)
		World.UpdateFromShipyard(resp)
		SaveWorldState(callsign, World)
		return time.Now()
	
	case ActionPurchaseCargo:
		resp := PurchaseCargo(action.ShipSymbol, action.TradeGoodSymbol, action.Units)
		ship.Cargo = resp.Cargo
		World.UpdateFromAgent(resp.Agent)
		World.InvalidateMarket(ship.Nav.WaypointSymbol)
		return time.Now()

	case ActionSellCargo:
		resp := SellCargo(action.ShipSymbol, action.TradeGoodSymbol, action.Units)
		ship.Cargo = resp.Cargo
		World.UpdateFromAgent(resp.Agent)
		World.InvalidateMarket(ship.Nav.WaypointSymbol)
		return time.Now()
	
	case ActionPurchaseShip:
		resp := PurchaseShip(action.ShipType, action.WaypointSymbol)
		World.UpdateFromAgent(resp.Agent)
		World.UpdateFromShip(resp.Ship)
		ship_state := World.Ships[resp.Ship.Symbol]
		ensureShipRunning(World, ship_state)
		World.InvalidateShipyard(ship.Nav.WaypointSymbol)
		return time.Now()
	
	}

    panic("unknown action")
}