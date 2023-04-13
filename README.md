# Senior Software Engineer Take-Home Programming Assignment for Golang

## Algorithm
The algorithm doesn't try to connect all the flights together but calculates indegree and outdegree for each airport.

Based on the difference of indegree and outdegree, we can detect the start and end airports:
- The starting airport will have one more outgoing edge than incoming edges (outgoing - incoming = 1).
- The final airport will have one more incoming edge than outgoing edges (incoming - outgoing = 1).

If it was required to return the path, it would be possible by iterating through the map using the starting airport 
and following the flights until the final airport is reached.

This algorithm has a linear time complexity of O(V+E), where V is the number of nodes (airports) and E is the 
number of edges (flights).

## Implementation
The service is split into a few abstraction layers:
- [`server`](server/) - the JSON-RPC server. It is a wrapper around the `net/http` package that adds a graceful shutdown.
- [`finder`](finder/) - the core algorithm
- [`service`](service/) - handles the JSON-RPC requests and responses and calls the finder

Each package is covered with tests that can be run with: 
```sh
go test ./...
```

I tried to use as least of external dependencies as possible, although I used `testify` to make
unit tests more compact.

## Running the server

### With Go installed
```sh
go run main.go
```

### With Docker
```sh
docker build -t pathfinder-server
docker run -p 8080:8080 pathfinder-server
```

## API

### Endpoint: `/calculate`

**HTTP method**: `POST`

**Body** : `[][]string` - a list of flight codes

For example,
```json
[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]
```

**Response**: `[]string` - an array that contains the start and end airports

For example,
```json
["SFO", "EWR"]
```

## Testing

### From the command line
```sh
curl -v -X POST -H 'Content-type: application/json' \
    -d '[["SWO", "ATL"], ["ATL", "LGW"]]' localhost:8080/calculate
```

Response:
```json
["SWO", "LGW"]
```
