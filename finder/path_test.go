package finder

import (
	"testing"
)

func TestFindFlightPath(t *testing.T) {
	testCases := []struct {
		name          string
		flights       []Flight
		expected      Path
		expectedError error
	}{
		{
			name: "Single flight",
			flights: []Flight{
				{"SFO", "EWR"},
			},
			expected: Path{Start: "SFO", Finish: "EWR"},
		},
		{
			name: "Disconnected flights",
			flights: []Flight{
				{"SFO", "EWR"},
				{"ATL", "BNM"},
			},
			expectedError: ErrMultiplePaths,
		},
		{
			name: "Round trip",
			flights: []Flight{
				{"SFO", "EWR"},
				{"EWR", "SFO"},
			},
			expectedError: ErrNoPath,
		},
		{
			name: "Two flights",
			flights: []Flight{
				{"ATL", "EWR"},
				{"SFO", "ATL"},
			},
			expected: Path{Start: "SFO", Finish: "EWR"},
		},
		{
			name: "Four flights",
			flights: []Flight{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
			expected: Path{"SFO", "EWR"},
		},
		{
			name: "Invalid route",
			flights: []Flight{
				{"ATL", "ATL"},
			},
			expectedError: ErrNoPath,
		},
		{
			name: "Multiple final destinations",
			flights: []Flight{
				{"ATL", "SFO"},
				{"ATL", "LGW"},
			},
			expectedError: ErrMultiplePaths,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pathFinder := NewPathFinder()
			result, err := pathFinder.FindFlightPath(tc.flights)
			if err != nil {
				if tc.expectedError != err {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
				return
			}
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
