// Graph represents a data structure with a set of weighted edges
// that connect pair of nodes together. Weights, a n by n matrix,
// with n being the number of nodes. Weight[i][j] shows the weight
// of the edge between node i and node j.

// author: Vahid Mostofi, Jun 5 2021

package main

import (
	"container/heap"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Graph is the structure representing a graph with the n by n Weights matrix.
// With n being the number of nodes. Each node, internall has a numberic
// representation using type int. It also has a string representation that can
// be mapped to numberic representation using nodeToId and back with idToNode.
type Graph struct {
	weights  [][]int
	nodeToId map[string]int
	idToNode map[int]string
}

// ErrNoNodeFound happens when name of a not existing node is provided
var ErrNoNodeFound = fmt.Errorf("no node found")

// ErrNoSuchRoute happens when the route uses a non-existing edge
var ErrNoSuchRoute = fmt.Errorf("no such route")

// ErrInvalidRoute happens when the length of the route has less than 2 nodes
var ErrInvalidRoute = fmt.Errorf("invalid route error")

// ErrInvalidRouteInputFormat happens the route is provided as string but it has
// incorrect format
var ErrInvalidRouteInputFormat = fmt.Errorf("invalid format for the route")

// ErrInvalidGraphInputFormat happens when the input provided to build the
// graph, has incorrect format, expected format: AB3, EF5, EG10, AD1
var ErrInvalidGraphInputFormat = fmt.Errorf("invalid format for graph input (edge list with weight)")

//edge is an internal type, used for parsing data of the graph. Represents
//an edge between source and destination with a specific weight.
type edge struct {
	source      int
	destination int
	weight      int
}

// NewGraphFromReader generates a new graph based on string data extracted
// from an io.Reader. The input data must be edge list with each edge as:
// NodeName1NodeName2Weight.
// Example of valid input:
// AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7
// Assumption is the graph is directed and weighted
func NewGraphFromReader(r io.Reader) (*Graph, error) {
	g := &Graph{
		nodeToId: make(map[string]int),
		idToNode: make(map[int]string),
	}

	// read input from the reader
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from provided reader: %w", err)
	}
	iStr := string(b)

	// validate the provided input againset the required schema
	re := regexp.MustCompile(`([a-zA-Z][a-zA-Z][0-9]*, )*[a-zA-Z][a-zA-Z][0-9]*`)
	if !re.MatchString(iStr) {
		return nil, ErrInvalidGraphInputFormat
	}

	// parse each edge data and store them in a list, we don't know the number of
	// nodes, so we use this approach to build up the graph 2d representation
	edges := make([]*edge, 0)
	splits := strings.Split(iStr, ",")
	for _, split := range splits {
		split = strings.Trim(split, " ")
		s := string(split[0])
		d := string(split[1])
		wStr := split[2:]
		w, err := strconv.Atoi(wStr)
		if err != nil {
			panic(err)
		}

		// track nodes, if source or destination are new, add them to the map
		if _, exists := g.nodeToId[s]; !exists {
			g.idToNode[len(g.nodeToId)] = s
			g.nodeToId[s] = len(g.nodeToId)
		}

		if _, exists := g.nodeToId[d]; !exists {
			g.idToNode[len(g.nodeToId)] = d
			g.nodeToId[d] = len(g.nodeToId)
		}

		edges = append(edges, &edge{g.nodeToId[s], g.nodeToId[d], w})
	}

	g.weights = make([][]int, len(g.nodeToId))
	for i := range g.weights {
		g.weights[i] = make([]int, len(g.nodeToId))
	}

	// if we want to keep the edge list we can use build it here and augment the
	// graph.
	for _, e := range edges {
		g.weights[e.source][e.destination] = e.weight
	}
	return g, nil
}

// GetMinDistanceBetweenNodes returns the minimum distance between two nodes
// using Dijkstra algorithm.
func (g *Graph) GetMinDistanceBetweenNodes(src string, destination string) (int, error) {
	sourceNode, sourceExists := g.nodeToId[src]
	destinationNode, destinationExists := g.nodeToId[destination]

	if !sourceExists || !destinationExists {
		return -1, ErrNoNodeFound
	}

	length, _ := g.shortestPath(sourceNode, destinationNode)
	return length, nil
}

// GetNodeCount returns number of nodes in the graph
func (g *Graph) GetNodeCount() int {
	return len(g.weights)
}

// GetEdgeCount returns number of edges in the graph
func (g *Graph) GetEdgeCount() int {
	e := 0
	for i := range g.weights {
		for j := range g.weights[i] {
			if g.weights[i][j] > 0 {
				e++
			}
		}
	}
	return e
}

// GetLengthOfRoute returns the length of the provided route (sumo of the weights)
// provided input is a string with the format: node1-node2-node3, example: A-B-C
func (g *Graph) GetLengthOfRoute(route string) (int, error) {
	// validate the provided input againset the required schema
	re := regexp.MustCompile(`(([a-zA-Z]-)*)([A-Z])`)
	if !re.MatchString(route) {
		return -1, ErrInvalidRouteInputFormat
	}

	splits := strings.Split(route, "-")
	return g.GetLengthOfRouteStringSlice(splits)
}

// GetLengthOfRouteStringSlice returns the length of the provided route
// (sum of the weights) provided input is a slice of node names, example:
// []string{"A", "B", "C"}
func (g *Graph) GetLengthOfRouteStringSlice(route []string) (int, error) {
	if len(route) < 2 {
		return -1, ErrInvalidRoute
	}
	routeInts := make([]int, len(route))
	for i, p := range route {
		if v, exists := g.nodeToId[p]; exists {
			routeInts[i] = v
		} else {
			return -1, ErrNoNodeFound
		}
	}

	return g.getLengthOfRouteInts(routeInts)
}

// getLengthOfRouteInts computes and returns the total length (sum of weights)
// of route, the provided input must be a slice of node Ids.
func (g *Graph) getLengthOfRouteInts(route []int) (int, error) {
	length := 0
	for i := 0; i < len(route)-1; i++ {
		w := g.weights[route[i]][route[i+1]]
		if w == 0 {
			return -1, ErrNoSuchRoute
		}
		length += w
	}
	return length, nil
}

// shortestPath finds the shortest path between two nodes using Dijkstra
// algorithm
// src: https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
func (g *Graph) shortestPath(src int, target int) (int, []int) {
	const infinity = math.MaxInt32
	pq := new(PriorityQueue)
	heap.Init(pq)

	visited := make([]bool, len(g.weights))
	distance := make([]int, len(g.weights))
	parent := make([]int, len(g.weights))

	for i := range g.weights {
		visited[i] = false
		distance[i] = infinity
		parent[i] = -1
	}

	pq.Push(NewItem(src, 0))
	distance[src] = 0
	parent[src] = -1
	iter := 0
	for pq.Len() > 0 {
		u := heap.Pop(pq).(int)
		if visited[u] {
			continue
		}
		visited[u] = true
		if u == target && iter > 0 {
			break
		}
		iter++
		for v, w := range g.weights[u] {
			if visited[v] || w <= 0 {
				continue
			}
			if distance[v] > distance[u]+w {
				distance[v] = distance[u] + w
				parent[v] = u
				heap.Push(pq, NewItem(v, distance[v]))
			}
			// }
		}
		if src == target && u == target { // first iteration
			visited[u] = false
			distance[u] = infinity
		}
	}

	path := make([]int, target)
	p := parent[target]
	for p != -1 {
		path = append([]int{p}, path...)
		p = parent[p]
		if src == target && p == src {
			break
		}
	}
	return distance[target], path
}

// GetAllRoutesWithExactSize finds all the Routes between source and target that
// have a size exactly equal to size. The returning results, can have cycles.
// Warning: be carefull with the value of size, a high value can lead
// to consuming too much memory
func (g *Graph) GetAllRoutesWithExactSize(source, target string, size int) ([][]int, error) {
	sourceNode, sourceExists := g.nodeToId[source]
	targetNode, targetExists := g.nodeToId[target]

	if !sourceExists || !targetExists {
		return nil, ErrNoNodeFound
	}

	return g.allRoutesSourceTarget(sourceNode, targetNode,
		func(i []int) bool { return len(i) <= size },
		func(i []int) bool { return len(i) == size }), nil
}

// GetAllRoutesWithMaxSize finds all the Routes between source and target that
// have a size exactly equal to size. The returning results, can have cycles.
// Warning: be carefull with the value of size, a high value can lead
// to consuming too much memory
func (g *Graph) GetAllRoutesWithMaxSize(source, target string, size int) ([][]int, error) {
	sourceNode, sourceExists := g.nodeToId[source]
	targetNode, targetExists := g.nodeToId[target]

	if !sourceExists || !targetExists {
		return nil, ErrNoNodeFound
	}

	return g.allRoutesSourceTarget(sourceNode, targetNode,
		func(i []int) bool { return len(i) <= size },
		func(i []int) bool { return len(i) <= size && len(i) > 1 }), nil
}

// GetAllRoutesWithLengthLessThan finds all the Routes between source and target that
// have a length less than maxRouteLength. The returning results, can have cycles.
// Warning: be carefull with the value of maxRouteLength, a high value can lead
// to consuming too much memory
func (g *Graph) GetAllRoutesWithLengthLessThan(source, target string, lengthLessThan int) ([][]int, error) {
	sourceNode, sourceExists := g.nodeToId[source]
	targetNode, targetExists := g.nodeToId[target]

	if !sourceExists || !targetExists {
		return nil, ErrNoNodeFound
	}

	return g.allRoutesSourceTarget(sourceNode, targetNode, func(i []int) bool {

		// get the length of the new route
		l, _ := g.getLengthOfRouteInts(i)
		return l < lengthLessThan

	}, func(i []int) bool { return len(i) > 1 }), nil
}

// allRoutesSourceTarget finds all the Routes between source and target that have
// a length less than lengthLessThan. The returning results, can have cycles.
// Warning: be carefull with the value of lengthLessThan, a high value can lead
// to consuming too much memory
func (g *Graph) allRoutesSourceTarget(source, target int, checkSubRoute func([]int) bool, checkRoute func([]int) bool) [][]int {
	routes := make([][]int, 0)

	q := make([][]int, 0)
	q = append(q, []int{source})

	var x []int
	for len(q) > 0 {
		x, q = q[0], q[1:]

		last := x[len(x)-1]
		if last == target && checkRoute(x) {
			routes = append(routes, x)
		}

		for y, w := range g.weights[last] {
			if w == 0 {
				continue
			}
			// build the new route
			newRoute := make([]int, len(x))
			copy(newRoute, x)
			newRoute = append(newRoute, y)

			// add the new route if it valid
			if checkSubRoute(newRoute) {
				q = append(q, newRoute)
			}
		}
	}

	return routes
}

// routeToString convers a route of integers to a string, representing a route
// of nodes
func (g *Graph) routeToString(route []int) string {
	var sb strings.Builder
	for _, p := range route {
		sb.WriteString(g.idToNode[p])
	}
	return sb.String()
}
