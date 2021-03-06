# Kiwiland railway assistant <!-- omit in toc -->

- [To run with original example](#to-run-with-original-example)
- [How to build and run](#how-to-build-and-run)
- [How to use](#how-to-use)
  - [Command Line Interface](#command-line-interface)
  - [Changing Input/Output mediums](#changing-inputoutput-mediums)
- [Problem Statement/ Initial Requirements](#problem-statement-initial-requirements)
  - [Input features/assumptions](#input-featuresassumptions)
  - [Primary Requirements](#primary-requirements)

## To run with original example
```
script/build
./kiwiland -f samples/original.txt
```

## How to build and run
- ### How to install dependencies
  ```
  script/bootstrap
  ```

- ### How to build
  ```
  script/build
  ```
- ### How to test
  ```
  script/test
  ```
- ### How to run
  After build:
  ```
  ./kiwiland
  ```
  or 
  ```
  ./kiwiland -h
  ```
  Without build:
  ```
  go run *.go
  ```
  or 
  ```
  go run main.go graph.go priorityqueue.go
  ```

## How to use
### Command Line Interface
- Use in interactive mode:
  ```
  ./kiwiland -i
  ```
  [screenshot of interactive mode](https://raw.githubusercontent.com/VahidMostofi/kiwiland/master/kiwkiland.png)
- To pass a file use: (make sure the last command has a new line at the end)
  ```
  ./kiwiland -f <filename>
  ```

- Use `kiwiland -h` to see this message:
  ```
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

  - To provide input using a file, pass it with -f option. The file needs to have input in the first line and can have any number of commands after that, each in one line.
  $> kiwiland -f sample-input-file.txt
  ```

### Changing Input/Output mediums
The function that handles input and output works with `io.Reader` and `io.Writer`, so adding other sources for reading input from and writing output to would be easy.
## Problem Statement/ Initial Requirements
### Input features/assumptions
 - The graph is a directed, weighted graph.

### Primary Requirements
 1. The distance along a certain route.
 2. The number of different routes between two towns (can be the same) with max length. Length means, the distance between two towns matter.
 3. The shortest path between two towns.
 4. The number of trips between two towns (can be the same) with exact or max stops.
    Trips mean, the distance between two towns doesn't matter.