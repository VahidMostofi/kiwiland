package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		printHelp(os.Stdout)
	} else if len(args) == 1 && args[0] == "-i" {
		handleInput(os.Stdin, os.Stdout)
	} else if len(args) == 2 && args[0] == "-f" { // -f filename
		f, err := os.Open(args[1])
		if err != nil {
			panic(err)
		}
		handleInput(f, os.Stdout)
	} else {
		printHelp(os.Stdout)
	}
}

const DistanceCommanPrefix = "distance of route"
const ShortestPathCommanPrefix = "shortest route"
const AllRoutesCommandPrefix = "all routes"
const AllTripsCommandPrefix = "all trips"

var reader *bufio.Reader

func readSingleLine() (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "exit", nil
		}
		return "", err
	}
	line = strings.Replace(line, "\n", "", 1)
	line = strings.Trim(line, " ")
	return line, nil
}

func handleInput(r io.Reader, w io.Writer) error {
	reader = bufio.NewReader(r)
	inputLine, err := readSingleLine()
	if err != nil {
		return err
	}
	g, err := NewGraphFromReader(strings.NewReader(inputLine))
	if err != nil {
		return err
	}

	for {
		var r int
		line, err := readSingleLine()
		if err != nil {
			return err
		}
		if line == "exit" {
			return nil
		} else if strings.HasPrefix(line, DistanceCommanPrefix) {
			r = handleDistanceCommand(w, line, g)
		} else if strings.HasPrefix(line, ShortestPathCommanPrefix) {
			r = handleShortestPathCommand(w, line, g)
		} else if strings.HasPrefix(line, AllRoutesCommandPrefix) {
			r = handleAllRoutesCommand(w, line, g)
		} else if strings.HasPrefix(line, AllTripsCommandPrefix) {
			r = handleAllTripsCommand(w, line, g)
		} else if line == "help" {
			printHelp(w)
		} else {
			fmt.Fprintln(w, "unknown command, if you need help, type help")
		}

		if r != 0 {
			fmt.Fprintln(w, "error in running command, if you need help, type help")
		}
	}
}

// all trips X Y = w
// all trips X Y <= w
func handleAllTripsCommand(w io.Writer, line string, g *Graph) int {
	var (
		src, dst, operator string
		steps              int
		d                  [][]int
		err                error
	)
	fmt.Sscanf(line, AllTripsCommandPrefix+" %s %s steps %s %d", &src, &dst, &operator, &steps)
	steps++

	if operator == "=" {
		d, err = g.GetAllRoutesWithExactSize(src, dst, steps)
	} else if operator == "<=" {
		d, err = g.GetAllRoutesWithMaxSize(src, dst, steps)
	} else {
		return 1
	}
	if err != nil {
		fmt.Fprintln(w, err.Error(), "; if you need help, type help")
		return 0
	}
	fmt.Fprintln(w, len(d))
	return 0
}

// all routes X Y distance < w
func handleAllRoutesCommand(w io.Writer, line string, g *Graph) int {
	var (
		src, dst, operator string
		value              int
		d                  [][]int
		err                error
	)
	fmt.Sscanf(line, AllRoutesCommandPrefix+" %s %s distance %s %d", &src, &dst, &operator, &value)
	if operator == "<" {
		d, err = g.GetAllRoutesWithLengthLessThan(src, dst, value)
	} else {
		return 1
	}
	if err != nil {
		fmt.Fprintln(w, err.Error(), "; if you need help, type help")
		return 0
	}
	fmt.Fprintln(w, len(d))
	return 0
}

// shortest route X Y
func handleShortestPathCommand(w io.Writer, line string, g *Graph) int {
	var (
		src, dst string
		d        int
		err      error
	)
	fmt.Sscanf(line, ShortestPathCommanPrefix+" %s %s", &src, &dst)
	d, err = g.GetMinDistanceBetweenNodes(src, dst)
	if err != nil {
		fmt.Fprintln(w, err.Error(), ";if you need help type help")
		return 0
	}
	fmt.Fprintln(w, d)
	return 0
}

// distance X-Y-Z
func handleDistanceCommand(w io.Writer, line string, g *Graph) int {
	var (
		route string
		d     int
		err   error
	)
	fmt.Sscanf(line, DistanceCommanPrefix+" %s", &route)
	d, err = g.GetLengthOfRoute(route)
	if err != nil {
		fmt.Fprintln(w, err.Error(), ";if you need help type help")
		return 0
	}
	fmt.Fprintln(w, d)
	return 0
}

func printHelp(w io.Writer) {
	message := `
- To see help use kiwiland -h or kiwiland --help.
- To provide input using stdin, run kiwiland without any arg. The input should be in one line and press new line
  formatted like: "AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7" (without double quotes).
- To provide command to get output using stdin, after entering input use this pattern:
  * distance of route X-Y-Z:
    distance X-Y-Z
  * shortest route of between X and Y:
    shortest route X Y
  * all routes between X and Y with a distance less than w:
    all routes X Y distance < w
  * all trips between X Y with exactly w stops:
    all trips X Y steps = w
  * all trips between X Y with maximum of 3 stops:
    all trips X Y steps <= w 
  * to see this message:
    help
  * exit:
    exit
	`

	fmt.Fprintln(w, message)

	noArgMessage := `- To provide input using a file, pass it with -f option. The file needs to have input
in the first line and can have any number of commands after that, each in one line.
$> kiwiland -f sample-input-file.txt
		`
	fmt.Fprintln(w, noArgMessage)
}
