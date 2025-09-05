package main

import (
	"fmt"

	"github.com/albertorestifo/dijkstra"
)

var system_graph dijkstra.Graph = make(dijkstra.Graph)
var marketplace_graph dijkstra.Graph = make(dijkstra.Graph)

func AddWaypointToSystemGraph(waypoint Waypoint) {
	fmt.Println("AddWaypointToSystemGraph")
	system_graph_row := make(map[string]int)
	system_graph[waypoint.Symbol] = system_graph_row
	fmt.Println(system_graph)
}

func PopulateSystemGraphDistancesForWaypoint(all_waypoints_in_system []Waypoint, waypoint_to_populate Waypoint) {
	system_graph_row := make(map[string]int)
	for _, other_waypoint := range all_waypoints_in_system {
		distance := DistanceBetweenTwoWaypoints(other_waypoint, waypoint_to_populate)
		system_graph_row[other_waypoint.Symbol] = distance
	}
	system_graph[waypoint_to_populate.Symbol] = system_graph_row
	fmt.Println(len(system_graph))
}

func PopulateMarketplaceGraphDistancesForWaypoint(marketplace_waypoints []Waypoint, waypoint_to_populate Waypoint) {
	max_distance := 400
	marketplace_graph_row := make(map[string]int)
	for _, other_waypoint := range marketplace_waypoints {
		distance := DistanceBetweenTwoWaypoints(other_waypoint, waypoint_to_populate)
		if distance < max_distance {
			marketplace_graph_row[other_waypoint.Symbol] = distance
		}
	}
	marketplace_graph[waypoint_to_populate.Symbol] = marketplace_graph_row
	fmt.Println(len(marketplace_graph))
}

func ProduceGraphWithMaxDistanceBetweenHops(input_graph dijkstra.Graph, max_distance int) dijkstra.Graph {
	//new_system_graph := make(dijkstra.Graph)
	for _, distances_map := range input_graph {
		for destination_waypoint_symbol, distance := range distances_map {
			if distance > max_distance {
				delete(distances_map, destination_waypoint_symbol)
			}
		}
	}
	return input_graph
}

func CalculateShortestPathBetweenTwoWaypoints(SourceWaypoint Waypoint, DestinationWaypoint Waypoint) (resultant_path []string, resultant_cost int) {
	path, cost, _ := system_graph.Path(SourceWaypoint.Symbol, DestinationWaypoint.Symbol) // skipping error handling
	fmt.Printf("path: %v, cost: %v", path, cost)
	return path, cost
}
