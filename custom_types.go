package main

type ListShipsResponseData struct {
	Data []Ship `json:"data"`
	Meta Meta   `json:"meta"`
}

type Ship struct {
	Symbol       string       `json:"symbol"`
	Nav          Nav          `json:"nav"`
	Crew         Crew         `json:"crew"`
	Fuel         Fuel         `json:"fuel"`
	Cooldown     Cooldown     `json:"cooldown"`
	Frame        Frame        `json:"frame"`
	Reactor      Reactor      `json:"reactor"`
	Engine       Engine       `json:"engine"`
	Modules      []Module     `json:"modules"`
	Mounts       []Mount      `json:"mounts"`
	Registration Registration `json:"registration"`
	Cargo        Cargo        `json:"cargo"`
}

type Cargo struct {
	Capacity  int64         `json:"capacity"`
	Units     int64         `json:"units"`
	Inventory []interface{} `json:"inventory"`
}

type Cooldown struct {
	ShipSymbol       string `json:"shipSymbol"`
	TotalSeconds     int64  `json:"totalSeconds"`
	RemainingSeconds int64  `json:"remainingSeconds"`
}

type Crew struct {
	Current  int64  `json:"current"`
	Capacity int64  `json:"capacity"`
	Required int64  `json:"required"`
	Rotation string `json:"rotation"`
	Morale   int64  `json:"morale"`
	Wages    int64  `json:"wages"`
}

type Engine struct {
	Symbol       string             `json:"symbol"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Condition    int64              `json:"condition"`
	Integrity    int64              `json:"integrity"`
	Speed        int64              `json:"speed"`
	Requirements EngineRequirements `json:"requirements"`
}

type EngineRequirements struct {
	Power int64 `json:"power"`
	Crew  int64 `json:"crew"`
}

type Frame struct {
	Symbol         string             `json:"symbol"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	ModuleSlots    int64              `json:"moduleSlots"`
	MountingPoints int64              `json:"mountingPoints"`
	FuelCapacity   int64              `json:"fuelCapacity"`
	Condition      int64              `json:"condition"`
	Integrity      int64              `json:"integrity"`
	Requirements   EngineRequirements `json:"requirements"`
}

type Fuel struct {
	Current  int64    `json:"current"`
	Capacity int64    `json:"capacity"`
	Consumed Consumed `json:"consumed"`
}

type Consumed struct {
	Amount    int64  `json:"amount"`
	Timestamp string `json:"timestamp"`
}

type Module struct {
	Symbol       string             `json:"symbol"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Capacity     *int64             `json:"capacity,omitempty"`
	Requirements ModuleRequirements `json:"requirements"`
}

type ModuleRequirements struct {
	Crew  int64 `json:"crew"`
	Power int64 `json:"power"`
	Slots int64 `json:"slots"`
}

type Mount struct {
	Symbol       string             `json:"symbol"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Strength     int64              `json:"strength"`
	Requirements EngineRequirements `json:"requirements"`
	Deposits     []string           `json:"deposits,omitempty"`
}

type Nav struct {
	SystemSymbol   string `json:"systemSymbol"`
	WaypointSymbol string `json:"waypointSymbol"`
	Route          Route  `json:"route"`
	Status         string `json:"status"`
	FlightMode     string `json:"flightMode"`
}

type Route struct {
	Origin        Destination `json:"origin"`
	Destination   Destination `json:"destination"`
	Arrival       string      `json:"arrival"`
	DepartureTime string      `json:"departureTime"`
}

type Destination struct {
	Symbol       string `json:"symbol"`
	Type         string `json:"type"`
	SystemSymbol string `json:"systemSymbol"`
	X            int64  `json:"x"`
	Y            int64  `json:"y"`
}

type Reactor struct {
	Symbol       string              `json:"symbol"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Condition    int64               `json:"condition"`
	Integrity    int64               `json:"integrity"`
	PowerOutput  int64               `json:"powerOutput"`
	Requirements ReactorRequirements `json:"requirements"`
}

type ReactorRequirements struct {
	Crew int64 `json:"crew"`
}

type Registration struct {
	Name          string `json:"name"`
	FactionSymbol string `json:"factionSymbol"`
	Role          string `json:"role"`
}

type Meta struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type ListWaypointsInSystemResponseData struct {
	Data []ListWaypointsInSystemResponse `json:"data"`
	Meta Meta                            `json:"meta"`
}

type ListWaypointsInSystemResponse struct {
	SystemSymbol        string        `json:"systemSymbol"`
	Symbol              string        `json:"symbol"`
	Type                string        `json:"type"`
	X                   int64         `json:"x"`
	Y                   int64         `json:"y"`
	Orbitals            []Faction     `json:"orbitals"`
	Traits              []Trait       `json:"traits"`
	Modifiers           []interface{} `json:"modifiers"`
	Chart               Chart         `json:"chart"`
	Faction             Faction       `json:"faction"`
	IsUnderConstruction bool          `json:"isUnderConstruction"`
}

type Chart struct {
	SubmittedBy string `json:"submittedBy"`
	SubmittedOn string `json:"submittedOn"`
}

type Faction struct {
	Symbol string `json:"symbol"`
}

type Trait struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetMarketResponseData struct {
	Data GetMarketResponse `json:"data"`
}

type GetMarketResponse struct {
	Symbol       string        `json:"symbol"`
	Exports      []Exchange    `json:"exports"`
	Imports      []Exchange    `json:"imports"`
	Exchange     []Exchange    `json:"exchange"`
	Transactions []Transaction `json:"transactions"`
	TradeGoods   []TradeGood   `json:"tradeGoods"`
}

type Exchange struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TradeGood struct {
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	TradeVolume   int64  `json:"tradeVolume"`
	Supply        string `json:"supply"`
	Activity      string `json:"activity"`
	PurchasePrice int64  `json:"purchasePrice"`
	SellPrice     int64  `json:"sellPrice"`
}

type Transaction struct {
	WaypointSymbol string `json:"waypointSymbol"`
	ShipSymbol     string `json:"shipSymbol"`
	TradeSymbol    string `json:"tradeSymbol"`
	Type           string `json:"type"`
	Units          int64  `json:"units"`
	PricePerUnit   int64  `json:"pricePerUnit"`
	TotalPrice     int64  `json:"totalPrice"`
	Timestamp      string `json:"timestamp"`
}
