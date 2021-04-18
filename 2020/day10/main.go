package main

// https://adventofcode.com/2020/day/10

import (
	"bufio"
	"log"
	"os"
	"time"
)

func main() {
	defer timeTrack(time.Now(), "day10")
	input, _ := readLines("adapters.txt")
	outlet := 0
	maxFoltDiff := 3

}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
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
