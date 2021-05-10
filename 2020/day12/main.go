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
	part := "B"

	// PartA inputs
	// waypoint := []int{0, 0} // {N, E}
	////
	// PartB inputs
	waypoint := []int{1, 10} // {N, E}
	////
	north := 0
	east := 0
	facing := "E"

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

	for i := range dir {
		logging("---------", "", true)
		logging("dir[i]", dir[i], true)
		logging("dist[i]", dist[i], true)
		logging("north", north, true)
		logging("east", east, true)
		if part == "A" {
			logging("facing", facing, true)
			north, east, facing = navSwitch(dir[i], dist[i], north, east, facing)
		} else if part == "B" {
			logging("waypoint", waypoint, true)
			north, east, waypoint = navSwitchWaypoint(dir[i], dist[i], north, east, waypoint)
			logging("new_waypoint", waypoint, true)
			logging("new_north", north, true)
			logging("new_east", east, true)
		}
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

func navSwitchWaypoint(nav string, value int, north int, east int, waypoint []int) (int, int, []int) {
	switch nav {
	case "N":
		waypoint[0] = waypoint[0] + value
	case "S":
		waypoint[0] = waypoint[0] - value
	case "E":
		waypoint[1] = waypoint[1] + value
	case "W":
		waypoint[1] = waypoint[1] - value
	case "L":
		waypoint = compassWaypoint(waypoint, -value)
	case "R":
		waypoint = compassWaypoint(waypoint, value)
	case "F":
		north = north + waypoint[0]*value
		east = east + waypoint[1]*value
	}
	return north, east, waypoint
}

func compassWaypoint(waypoint []int, degrees int) []int {
	_, change := divMod((360+degrees)/90, 4)
	if change == 0 {
		return waypoint
	}
	compass := []string{"N", "E", "S", "W"}
	compass_mag := []int{1, 1, -1, -1}
	var direction string
	var new_waypoint []int
	var compass_index_n int
	var compass_index_e int

	for i := range compass {
		if waypoint[0] >= 0 {
			direction = "N"
		} else {
			direction = "S"
		}
		if compass[i] == direction {
			compass_index_n = i
			break
		}
	}
	for i := range compass {
		if waypoint[1] >= 0 {
			direction = "E"
		} else {
			direction = "W"
		}
		if compass[i] == direction {
			compass_index_e = i
			break
		}
	}

	_, compass_index_n = divMod(compass_index_n+change, 4)
	_, compass_index_e = divMod(compass_index_e+change, 4)

	_, odd_even := divMod(change, 2)
	if odd_even == 0 {
		new_waypoint = append(new_waypoint, compass_mag[compass_index_n]*Abs(waypoint[0]))
		new_waypoint = append(new_waypoint, compass_mag[compass_index_e]*Abs(waypoint[1]))
	} else {
		new_waypoint = append(new_waypoint, compass_mag[compass_index_e]*Abs(waypoint[1]))
		new_waypoint = append(new_waypoint, compass_mag[compass_index_n]*Abs(waypoint[0]))
	}

	return new_waypoint
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
	_, change := divMod((360+degrees)/90, 4)
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
