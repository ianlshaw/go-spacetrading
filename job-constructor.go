package main

import "fmt"

func DecideConstructorAction(ship Ship, world *WorldState) ShipAction {
	fmt.Println("[DEBUG] DecideConstructorAction")

	// go-to closest incomplete jump gate

	first_construction_site := ConstructionSite{}

	if len(world.ConstructionSites) > 0 {
		for _, construction_site_state := range world.ConstructionSites {
			first_construction_site = construction_site_state.ConstructionSite
			break
		}

	} else {
		fmt.Printf("[WARN] No unfinished construction sites. Waiting...\n")
		return ShipAction{
			Type:      ActionWait,
			NotBefore: ThreeHoursFromNow(),
		}
	}

	materials := first_construction_site.Materials

	for _, material := range materials {
		if IsMaterialFulfilled(material) {
			continue
		} else {
			trade_good_symbol := material.TradeSymbol
			fmt.Println(trade_good_symbol)
			remaining := material.Required - material.Fulfilled
			fmt.Println(remaining)
		}
	}

	fmt.Printf("[WARN] DecideConstructorAction uncaught branch. Waiting...\n")
	return ShipAction{
		Type:      ActionWait,
		NotBefore: ThreeHoursFromNow(),
	}
	//construction_site, ok := world.ConstructionSites[first_under_construction_jump_gate_waypoint.Symbol]

	// update construction site data

	// identify buy location for missing materials

	// go to marketplace selling missing materials

	// purchase missing materials

	// go to construction site

	// deliver construction materials

	// do nothing
}
