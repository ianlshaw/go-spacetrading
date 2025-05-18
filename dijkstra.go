package main

import (
	"fmt"

	"github.com/albertorestifo/dijkstra"
)

var system_graph dijkstra.Graph = make(dijkstra.Graph)

func AddWaypointToSystemGraph(waypoint Waypoint) {
	fmt.Println("AddWaypointToSystemGraph")
	system_graph_row := make(map[string]int)
	system_graph[waypoint.Symbol] = system_graph_row
	fmt.Println(system_graph)
}

func PopulateSystemGraphDistancesForWaypoint(all_waypoints_in_system []Waypoint, waypoint_to_populate Waypoint) {
	system_graph_row := make(map[string]int)
	for _, other_waypoint := range all_waypoints_in_system {
		if waypoint_to_populate.Symbol != other_waypoint.Symbol {
			distance := DistanceBetweenTwoWaypoints(other_waypoint, waypoint_to_populate)
			system_graph_row[other_waypoint.Symbol] = distance
		}
	}
	system_graph[waypoint_to_populate.Symbol] = system_graph_row
	fmt.Println(len(system_graph))
}

func CalculateShortestPathBetweenTwoWaypoints(SourceWaypoint Waypoint, DestinationWaypoint Waypoint) (resultant_path []string, resultant_cost int) {
	g := dijkstra.Graph{
		"a": {"b": 20, "c": 80},
		"b": {"a": 20, "c": 20},
		"c": {"a": 80, "b": 20},
	}

	path, cost, _ := g.Path("a", "c") // skipping error handling

	fmt.Printf("path: %v, cost: %v", path, cost)
	return path, cost
}
