package main

// https://adventofcode.com/2020/day/16

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type field struct {
	Name string
	Min1 int
	Max1 int
	Min2 int
	Max2 int
}

func main() {
	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day16")
	input, _ := readLines("notes.txt")
	// input, _ := readLines("notes-test-p1.txt")
	// input, _ := readLines("notes-test-p2.txt")

	var fieldRules []field

	var your_ticket []int
	var other_tickets [][]int

	prev_row := ""
	process_tickets := false

	for row := range input {
		if prev_row == "nearby tickets:" {
			process_tickets = true
		}
		if prev_row == "your ticket:" {
			s := strings.Split(input[row], ",")
			for i := range s {
				v, err := strconv.Atoi(s[i])
				if err == nil {
					your_ticket = append(your_ticket, v)
				}
			}
		} else if input[row] == "" || input[row] == "your ticket:" || input[row] == "nearby tickets:" {
			prev_row = input[row]
			continue
		} else if process_tickets {
			s := strings.Split(input[row], ",")
			var cur_ticket []int
			for i := range s {
				v, err := strconv.Atoi(s[i])
				if err == nil {
					cur_ticket = append(cur_ticket, v)
				}
			}
			other_tickets = append(other_tickets, cur_ticket)
		} else if !process_tickets {
			newfield := parseFields(input[row])
			fieldRules = append(fieldRules, newfield)
		}
		prev_row = input[row]
	}

	logging("fieldRules", fieldRules, true)

	invalid_sum := 0
	for i := range other_tickets {
		for j := range other_tickets[i] {
			valid := false
			v := other_tickets[i][j]
			for f := range fieldRules {
				if (v >= fieldRules[f].Min1 && v <= fieldRules[f].Max1) || (v >= fieldRules[f].Min2 && v <= fieldRules[f].Max2) {
					valid = true
				}
			}
			if !valid {
				invalid_sum = invalid_sum + v
			}
		}
	}

	logging("---Answer---", "", true)
	logging("invalid_sum", invalid_sum, true)
}

func parseFields(s string) field {
	var newfield = field{}
	newfield.Name, _ = getFieldName(s)
	newfield.Min1, newfield.Max1, newfield.Min2, newfield.Max2 = getFieldMinMax(s)
	return newfield
}

func getFieldName(s string) (string, error) {
	re := regexp.MustCompile(`(.*): `)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		logging("match", s, false)
		return match[1], nil
	}
	return s, errors.New("error")
}

func getFieldMinMax(s string) (int, int, int, int) {
	re := regexp.MustCompile(`[0-9]+`)
	match := re.FindAllStringSubmatch(s, 4)
	logging("-getFieldMinMax-", "", false)
	logging("s", s, false)
	logging("match", match, false)
	if len(match) > 2 {
		logging("match[0][0]", match[0][0], false)
		logging("match[1][0]", match[1][0], false)
		logging("match[2][0]", match[2][0], false)
		logging("match[3][0]", match[3][0], false)
		v1, _ := strconv.Atoi(match[0][0])
		v2, _ := strconv.Atoi(match[1][0])
		v3, _ := strconv.Atoi(match[2][0])
		v4, _ := strconv.Atoi(match[3][0])
		return v1, v2, v3, v4
	}
	return 0, 0, 0, 0
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
