package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
	"slices"
)

var SystemGraph dijkstra.Graph = make(dijkstra.Graph)
var MarketplaceGraph dijkstra.Graph = make(dijkstra.Graph)


func AddWaypointToGraph(graph dijkstra.Graph, waypoint Waypoint) {
	fmt.Println("[DEBUG] AddWaypointToGraph")
	graph_row := make(map[string]int)
	graph[waypoint.Symbol] = graph_row
}

func PopulateGraphDistancesForWaypoint(graph dijkstra.Graph, all_waypoints_in_system []Waypoint, waypoint_to_populate Waypoint) {
	graph_row := make(map[string]int)
	for _, other_waypoint := range all_waypoints_in_system {
		distance := DistanceBetweenTwoWaypoints(other_waypoint, waypoint_to_populate)
		graph_row[other_waypoint.Symbol] = distance
	}
	graph[waypoint_to_populate.Symbol] = graph_row
}

func PopulateGraphDistancesForWaypointWithMaximum(graph dijkstra.Graph, all_waypoints_in_system []Waypoint, waypoint_to_populate Waypoint, max_distance int) {
	graph_row := make(map[string]int)
	for _, other_waypoint := range all_waypoints_in_system {
		distance := DistanceBetweenTwoWaypoints(other_waypoint, waypoint_to_populate)
		if distance < max_distance {
			graph_row[other_waypoint.Symbol] = distance
		}
	}
	graph[waypoint_to_populate.Symbol] = graph_row
}

func CalculateShortestPathBetweenTwoWaypoints(Graph dijkstra.Graph, SourceWaypoint Waypoint, DestinationWaypoint Waypoint) (resultant_path []string, resultant_cost int) {
	path, cost, _ := Graph.Path(SourceWaypoint.Symbol, DestinationWaypoint.Symbol) // skipping error handling
	//fmt.Printf("path: %v, cost: %v", path, cost)
	return path, cost
}

func FollowPath(ship *Ship, path []string) NavigateShipResponse {
	ship_waypoint_symbol := ship.Nav.WaypointSymbol
	last_waypoint_in_path := path[len(path)-1]
	if ship_waypoint_symbol == last_waypoint_in_path {
		fmt.Println("[ERROR] Already at path final destination. FollowPath shouldnt have been called")
		return NavigateShipResponse{}
	}
	current_waypoint_path_index := slices.Index(path, ship_waypoint_symbol)
	target_waypoint_path_index := current_waypoint_path_index + 1
	target_waypoint := path[target_waypoint_path_index]
	resp := NavigateShip(ship.Symbol, target_waypoint)
	return resp
}

func IsShipLocatedInGraph(ship Ship, graph dijkstra.Graph) bool {
	for waypoint_symbol := range graph {
		if ship.Nav.WaypointSymbol == waypoint_symbol {
			return true
		}
	}
	return false
}

func IsWaypointInGraph(waypoint Waypoint, graph dijkstra.Graph) bool {
	for waypoint_symbol := range graph {
		if waypoint.Symbol == waypoint_symbol {
			return true
		}
	}
	return false
}