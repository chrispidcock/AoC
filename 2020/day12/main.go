package main

// https://adventofcode.com/2020/day/12

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// PartA inputs

	// PartB inputs

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day12")
	// input, _ := readLines("nav-i-part-a-test.txt")
	input, _ := readLines("nav-instructions.txt")

	var dir []string
	var dist []int
	for i := range input {
		dir = append(dir, string(input[i][0]))
		distance, _ := strconv.Atoi(input[i][1:])
		dist = append(dist, distance)
	}

	north := 0
	east := 0
	facing := "E"

	for i := range dir {
		logging("dir[i]", dir[i], false)
		logging("dist[i]", dist[i], false)
		logging("north", north, false)
		logging("east", east, false)
		logging("facing", facing, false)
		north, east, facing = navSwitch(dir[i], dist[i], north, east, facing)
	}

	manhattan_dist := Abs(north) + Abs(east)
	fmt.Println("manhattan_dist ", manhattan_dist)
}

func navSwitch(nav string, value int, north int, east int, facing string) (int, int, string) {
	switch nav {
	case "N":
		north = north + value
	case "S":
		north = north - value
	case "E":
		east = east + value
	case "W":
		east = east - value
	case "L":
		facing = compass(facing, -value)
	case "R":
		facing = compass(facing, value)
	case "F":
		north, east, facing = navSwitch(facing, value, north, east, facing)
	}
	return north, east, facing
}

func compass(facing string, degrees int) (new_facing string) {
	compass := []string{"N", "E", "S", "W"}
	var compass_index int
	for i := range compass {
		if compass[i] == facing {
			compass_index = i
			break
		}
	}
	_, change := divMod(360+degrees/90, 4)
	_, new_index := divMod(compass_index+change, 4)
	logging("compass_index", compass_index, true)
	logging("change", change, true)
	logging("new_index", new_index, true)
	return compass[new_index]
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func divMod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func logging(name string, value interface{}, active bool) {
	if active {
		fmt.Println(name, " ", value)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
