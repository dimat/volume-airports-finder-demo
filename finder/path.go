package finder

import "fmt"

type Airport string

type Flight struct {
	Source      Airport
	Destination Airport
}

type Path struct {
	Start  Airport
	Finish Airport
}

var ErrNoPath = fmt.Errorf("no path found")
var ErrMultiplePaths = fmt.Errorf("multiple paths found")

// PathFinder is an interface for finding a path between two airports
type PathFinder interface {
	FindFlightPath(flights []Flight) (Path, error)
}

type pathFinder struct {
}

func NewPathFinder() PathFinder {
	return &pathFinder{}
}

// FindFlightPath finds the starting and ending airports for a given set of flights.
// These flights may not be listed in order.
// If there is no path, returns ErrNoPath.
// If there are multiple paths, returns ErrMultiplePaths.
func (f *pathFinder) FindFlightPath(flights []Flight) (Path, error) {
	incomingEdges := make(map[Airport]int)
	outgoingEdges := make(map[Airport]int)

	for _, flight := range flights {
		outgoingEdges[flight.Source]++
		incomingEdges[flight.Destination]++
	}

	var startingAirport, finalAirport Airport

	for airport, outgoing := range outgoingEdges {
		incoming := incomingEdges[airport]
		if outgoing-incoming == 1 {
			if startingAirport != "" {
				return Path{}, ErrMultiplePaths
			}
			startingAirport = airport
		}
	}

	for airport, incoming := range incomingEdges {
		outgoing := outgoingEdges[airport]
		if incoming-outgoing == 1 {
			if finalAirport != "" {
				return Path{}, ErrMultiplePaths
			}
			finalAirport = airport
		}
	}
	if startingAirport == "" || finalAirport == "" {
		return Path{}, ErrNoPath
	}

	return Path{Start: startingAirport, Finish: finalAirport}, nil
}
