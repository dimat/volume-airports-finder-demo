package service

import (
	"errors"

	"github.com/dimat/volume-airports-finder/finder"
)

type Finder struct {
	PathFinder finder.PathFinder
}

type FindPathRequest [][]string

type FindPathResponse []string

func (f *Finder) Call(args FindPathRequest) (FindPathResponse, error) {
	flights := make([]finder.Flight, len(args), len(args))
	for idx, flight := range args {
		if len(flight) != 2 {
			return nil, errors.New("each flight record should have two airports")
		}
		flights[idx] = finder.Flight{
			Source:      finder.Airport(flight[0]),
			Destination: finder.Airport(flight[1]),
		}
	}
	path, err := f.PathFinder.FindFlightPath(flights)
	if err != nil {
		return nil, err
	}
	return []string{string(path.Start), string(path.Finish)}, err
}
