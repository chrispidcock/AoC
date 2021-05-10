package main

// https://adventofcode.com/2020/day/12

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	part := "A"

	// PartA inputs

	////
	// PartB inputs

	////

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day13")
	input, _ := readLines("bus-times.txt")
	// input, _ := readLines("bus-test.txt")
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
