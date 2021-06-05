package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// graphTestCase represents a single graph and contains multiple testcases for
// different functionalities implemented for the type Graph.
type graphTestCase struct {
	name                       string
	input                      string
	edgeCount                  int
	nodeCount                  int
	shortestPathTCS            []shortestPathTestCase
	routeLengthTCS             []routeLengthTestCase
	allRoutesLengthLessThanTCS []allRoutesLessThanLengthTestCase
	allRoutesSizeTCS           []allRoutesSize
}

// shortestPathTestCase the shortest path between source and target must be
// of length length, or returns error err
type shortestPathTestCase struct {
	source string
	target string
	length int
	err    error
}

// routeLengthTestCase length of the route must be equal to length
// or there should be an error equal to err
type routeLengthTestCase struct {
	route  string
	length int
	err    error
}

// allRoutesLessThanLengthTestCase number of routes with lengths less than lessThan must
// be equal to count, or there should be an error equal to err
type allRoutesLessThanLengthTestCase struct {
	source   string
	target   string
	lessThan int
	count    int
	err      error
}

// allRoutesSize number of routes between source and target with exact size
// of exactSize or max size of maxSize must be equal to count or return an error equal to err
// the one that is not being used should be -1
type allRoutesSize struct {
	source    string
	target    string
	maxSize   int
	exactSize int
	count     int
	err       error
}

func TestGraphFunctions(t *testing.T) {
	tc := graphTestCase{
		name:            "test-case-1",
		input:           "AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7",
		nodeCount:       5,
		edgeCount:       9,
		shortestPathTCS: []shortestPathTestCase{{"A", "C", 9, nil}, {"B", "B", 9, nil}},
		routeLengthTCS: []routeLengthTestCase{{"A-B-C", 9, nil}, {"A-D", 5, nil}, {"A-D-C", 13, nil},
			{"A-E-B-C-D", 22, nil}, {"A-E-D", -1, ErrNoSuchRoute}},
		allRoutesSizeTCS:           []allRoutesSize{{"C", "C", 4, -1, 2, nil}, {"A", "C", -1, 5, 3, nil}},
		allRoutesLengthLessThanTCS: []allRoutesLessThanLengthTestCase{{"C", "C", 30, 7, nil}},
	}

	graphTester(t, tc)
}

func graphTester(t *testing.T, tc graphTestCase) {
	g, err := NewGraphFromReader(strings.NewReader(tc.input))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, g.GetNodeCount(), tc.nodeCount)
	assert.Equal(t, g.GetEdgeCount(), tc.edgeCount)

	for i, sp := range tc.shortestPathTCS {
		d, err := g.GetMinDistanceBetweenNodes(sp.source, sp.target)
		if sp.err == nil {
			assert.NoError(t, err, fmt.Sprintf("min distance: %d", i))
			assert.Equal(t, sp.length, d, fmt.Sprintf("min distance: %d", i))
		} else {
			assert.EqualError(t, sp.err, err.Error())
		}
	}

	for i, rl := range tc.routeLengthTCS {
		d, err := g.GetLengthOfRoute(rl.route)
		if rl.err == nil {
			assert.NoError(t, err, fmt.Sprintf("route length: %d", i))
			assert.Equal(t, rl.length, d, fmt.Sprintf("route length: %d", i))
		} else {
			assert.EqualError(t, rl.err, err.Error())
		}
	}

	for i, arl := range tc.allRoutesLengthLessThanTCS {

		d, err := g.GetAllRoutesWithLengthLessThan(arl.source, arl.target, arl.lessThan)
		if arl.err == nil {
			assert.NoError(t, err, fmt.Sprintf("all routes less than: %d", i))
			assert.Equal(t, arl.count, len(d), fmt.Sprintf("all routes less than: %d:\n+%v", i, d))
		} else {
			assert.EqualError(t, arl.err, err.Error())
		}
	}

	for i, ares := range tc.allRoutesSizeTCS {
		var d [][]int
		var err error
		if ares.maxSize == -1 {
			d, err = g.GetAllRoutesWithExactSize(ares.source, ares.target, ares.exactSize)
			if ares.err == nil {
				assert.NoError(t, err, fmt.Sprintf("all routes exact size: %d", i))
				assert.Equal(t, ares.count, len(d), fmt.Sprintf("all routes exact size: %d:\n%+v", i, d))
			} else {
				assert.EqualError(t, ares.err, err.Error())
			}
		} else if ares.exactSize == -1 {
			d, err = g.GetAllRoutesWithMaxSize(ares.source, ares.target, ares.maxSize)
			if ares.err == nil {
				assert.NoError(t, err, fmt.Sprintf("all routes max size: %d", i))
				assert.Equal(t, ares.count, len(d), fmt.Sprintf("all routes max size: %d:\n%+v", i, d))
			} else {
				assert.EqualError(t, ares.err, err.Error())
			}
		}
	}
}
