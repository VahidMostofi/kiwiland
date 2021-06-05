package main

import (
	"bytes"
	"strings"
	"testing"
)

const tc1 = `AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7
distance of route A-B-C
distance of route A-D
distance of route A-D-C
distance of route A-E-B-C-D
distance of route A-E-D
all trips C C steps <= 3
all trips A C steps = 4
shortest route A C
shortest route B B
all routes C C distance < 30
exit`
const tc1Out = `9
5
13
22
no such route ;if you need help type help
2
3
9
9
7
`

func TestInteractiveCommandLine(t *testing.T) {
	tcInput := tc1
	tcOutput := tc1Out

	r := strings.NewReader(tcInput)
	buf := bytes.NewBufferString("")
	handleInput(r, buf)
	actualOutput := strings.Trim(buf.String(), "\n")
	tcOutput = strings.Trim(buf.String(), "\n")

	if actualOutput != tcOutput {
		t.Fail()
		inputLines := strings.Split(tcInput, "\n")
		outputLines := strings.Split(tcOutput, "\n")
		actualOutputLines := strings.Split(actualOutput, "\n")
		for i := range outputLines {
			if outputLines[i] != actualOutputLines[i] {
				t.Logf("\ninput:%s\nactual output:%s\nexpected output:%s", inputLines[i+1], actualOutputLines[i], outputLines[i])
			}
		}
	}
}
