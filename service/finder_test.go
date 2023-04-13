package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dimat/volume-airports-finder/finder"
	"github.com/dimat/volume-airports-finder/service"
)

// MockPathFinder is a mock implementation of the finder.PathFinder interface
type MockPathFinder struct {
	mock.Mock
}

func (m *MockPathFinder) FindFlightPath(flights []finder.Flight) (finder.Path, error) {
	args := m.Called(flights)
	return args.Get(0).(finder.Path), args.Error(1)
}

func TestFinderFindPath(t *testing.T) {
	testCases := []struct {
		name          string
		input         service.FindPathRequest
		pathFinderRet finder.Path
		pathFinderErr error
		expectedReply service.FindPathResponse
		expectedErr   string
	}{
		{
			name:        "Invalid input",
			input:       [][]string{{"SFO"}},
			expectedErr: "each flight record should have two airports",
		},
		{
			name:          "PathFinder error",
			input:         [][]string{{"SFO", "EWR"}},
			pathFinderErr: errors.New("pathfinder error"),
			expectedErr:   "pathfinder error",
		},
		{
			name:          "Successful result",
			input:         [][]string{{"SFO", "EWR"}},
			pathFinderRet: finder.Path{Start: "SFO", Finish: "EWR"},
			expectedReply: []string{"SFO", "EWR"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPathFinder := new(MockPathFinder)
			finderService := &service.Finder{PathFinder: mockPathFinder}

			if tc.pathFinderErr != nil || (tc.pathFinderRet != finder.Path{}) {
				mockPathFinder.On("FindFlightPath", mock.Anything).Return(tc.pathFinderRet, tc.pathFinderErr)
			}

			reply, err := finderService.Call(tc.input)

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedReply, reply)
			mockPathFinder.AssertExpectations(t)
		})
	}
}
