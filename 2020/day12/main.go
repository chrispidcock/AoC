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
	input, _ := readLines("nav-i-part-a-test.txt")
	// input, _ := readLines("nav-instructions.txt")

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
	manhattan_dist := Abs(north) + Abs(east)

	fmt.Println("manhattan_dist ", manhattan_dist)
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
	_, change := divMod(degrees/90, 360)
	_, new_index := divMod(compass_index+change, 4)
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
