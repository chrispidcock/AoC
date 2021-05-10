package main

// https://adventofcode.com/2020/day/12

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	part := "B"

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day13")
	input, _ := readLines("bus-times.txt")
	// input, _ := readLines("bus-test.txt")

	var buses [9]int
	var bus_times [9]int
	var req_mins []int
	arv_time, _ := strconv.Atoi(input[0])

	next := 0
	s := strings.Split(input[1], ",")
	for a := range buses {
		for b := next; b < len(s); b++ {
			next = next + 1
			i, err := strconv.Atoi(s[b])
			if err == nil {
				buses[a] = i
				req_mins = append(req_mins, b)
				break
			}
		}
	}

	logging("arv_time", arv_time, true)
	logging("buses", buses, true)
	logging("req_mins", req_mins, true)
	t := 0

	answer_found := false
	first_bus := 0

	for !answer_found && part == "A" {
		t = t + 1
		for b := range buses {
			_, mod := divMod(t, buses[b])
			if mod == 0 {
				bus_times[b] = t
				if t > arv_time {
					first_bus = buses[b]
					answer_found = true
					break
				}
			}
		}
	}

	max_time := 0
	max_time_i := 0
	for i := range buses {
		if buses[i] > max_time {
			max_time = buses[i]
			max_time_i = i
		}
	}
	logging("max_time", max_time, true)
	logging("max_time_i", max_time_i, true)

	for !answer_found && part == "B" {
		logging("bus_times", bus_times, false)
		bus_times[max_time_i] = bus_times[max_time_i] + buses[max_time_i]

		for b, t := range req_mins {
			compare := bus_times[max_time_i] - (req_mins[max_time_i] - t)
			_, mod := divMod(compare, buses[b])
			if mod == 0 {
				bus_times[b] = compare
			} else {
				break
			}
			if b == len(buses)-1 {
				answer_found = true
			}
		}
	}

	// Part1 answer
	logging("---Part1 answer---", "", true)
	logging("bus_times", bus_times, true)
	logging("first_bus", first_bus, true)
	logging("arv_time", arv_time, true)
	wait_time := t - arv_time
	logging("wait_time", wait_time, true)
	part1_answer := first_bus * wait_time
	logging("part1_answer", part1_answer, true)

	// Part2 Answer
	logging("---Part2 answer---", "", true)
	logging("part2_answer", bus_times[0], true)
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
