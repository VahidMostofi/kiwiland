package main

import (
	"fmt"
	"strconv"
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

	tc2 := graphTestCase{
		name:                       "test-case-2",
		input:                      "AB1, BD2, BC10, DC3, CA4, FG2, GF10, GE3, EF4",
		nodeCount:                  7,
		edgeCount:                  9,
		shortestPathTCS:            []shortestPathTestCase{{"A", "C", 6, nil}, {"C", "A", 4, nil}, {"A", "F", -1, ErrNoSuchRoute}, {"F", "F", 9, nil}, {"G", "F", 7, nil}},
		routeLengthTCS:             []routeLengthTestCase{{"A-B-C", 11, nil}, {"A-B-F", -1, ErrNoSuchRoute}, {"G-F", 10, nil}, {"G-E-F", 7, nil}},
		allRoutesSizeTCS:           []allRoutesSize{{"A", "A", 5, -1, 2, nil}, {"A", "F", 10, -1, 0, nil}, {"F", "G", -1, 5, 1, nil}, {"F", "G", 6, -1, 4, nil}},
		allRoutesLengthLessThanTCS: []allRoutesLessThanLengthTestCase{{"A", "A", 15, 1, nil}, {"A", "A", 10, 0, nil}, {"A", "F", 10, 0, nil}, {"A", "C", 12, 2, nil}},
	}

	graphTester(t, tc)
	graphTester(t, tc2)
}

func TestCompleteGraph(t *testing.T) {
	k := 5
	nodes := []string{"A", "B", "C", "D", "E"}
	l := 1
	input := ""
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			if i != j {
				input += nodes[i] + nodes[j] + strconv.Itoa(l) + ", "
			}
		}
	}
	input = input[:len(input)-2]
	gtc := graphTestCase{
		name:                       "complete-graph-" + strconv.Itoa(k),
		input:                      input,
		nodeCount:                  k,
		edgeCount:                  (k * (k - 1)),
		shortestPathTCS:            make([]shortestPathTestCase, 0),
		routeLengthTCS:             make([]routeLengthTestCase, 0),
		allRoutesSizeTCS:           make([]allRoutesSize, 0),
		allRoutesLengthLessThanTCS: make([]allRoutesLessThanLengthTestCase, 0),
	}

	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			if i == j {
				gtc.shortestPathTCS = append(gtc.shortestPathTCS, shortestPathTestCase{nodes[i], nodes[j], l * 2, nil})
			} else {
				gtc.shortestPathTCS = append(gtc.shortestPathTCS, shortestPathTestCase{nodes[i], nodes[j], l, nil})
			}
		}
	}

	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			if i != j {
				gtc.allRoutesSizeTCS = append(gtc.allRoutesSizeTCS, allRoutesSize{nodes[i], nodes[j], 52, 5, 4096, nil})
			}
		}
	}

	graphTester(t, gtc)
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
			assert.Equal(t, sp.length, d, fmt.Sprintf("min distance: %d; %v", i, sp))
		} else {
			assert.EqualError(t, sp.err, err.Error())
		}
	}

	for i, rl := range tc.routeLengthTCS {
		d, err := g.GetLengthOfRoute(rl.route)
		if rl.err == nil {
			assert.NoError(t, err, fmt.Sprintf("route length: %d", i))
			assert.Equal(t, rl.length, d, fmt.Sprintf("route length: %d; %v", i, rl))
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
