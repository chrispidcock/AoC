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
	Name        string
	Min1        int
	Max1        int
	Min2        int
	Max2        int
	TicketIndex int
}

func main() {
	part := "B"

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
	tickets := len(other_tickets)

	for i := 0; i < tickets; i++ {
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
				if part == "B" {
					other_tickets = removeIndex(other_tickets, i)
					i = i - 1
					tickets = tickets - 1
				}
			}
		}
	}

	logging("---Answer Part1---", "", true)
	logging("invalid_sum", invalid_sum, true)

	var ticket_invalid_rules [][]int

	for rule := range fieldRules {
		logging("rule", rule, false)
		ticket_invalid_rules = append(ticket_invalid_rules, make([]int, 0))
	}

	// find invalid rules for each unknown ticket field in other_tickets
	for i := 0; i < tickets; i++ {
		for j := range other_tickets[i] {
			v := other_tickets[i][j]
			for f := range fieldRules {
				if (v >= fieldRules[f].Min1 && v <= fieldRules[f].Max1) || (v >= fieldRules[f].Min2 && v <= fieldRules[f].Max2) {
					continue
				} else {
					num_index := findNum(ticket_invalid_rules[f], j)
					if num_index == -1 {
						ticket_invalid_rules[f] = append(ticket_invalid_rules[f], j)
					}
				}
			}
		}
	}

	logging("ticket_invalid_rules", ticket_invalid_rules, true)

	// order the ticket fields by the number of invalid field rules
	var desc_len_order []int
	for j := len(fieldRules); j >= 0; j-- {
		for i := range ticket_invalid_rules {
			if len(ticket_invalid_rules[i]) == j {
				desc_len_order = append(desc_len_order, i)
				break
			}
		}
	}

	logging("desc_len_order", desc_len_order, true)

	for _, v := range desc_len_order {
		// iterate over ticket_invalid_rules to find rule index with no match
		for j := range fieldRules {
			if fieldRules[j].TicketIndex != -1 {
				continue
			}
			index := findNum(ticket_invalid_rules[v], j)
			if index == -1 {
				fmt.Println("VALID ", j)
				fieldRules[j].TicketIndex = v
				break
			}
		}
	}

	depart_multi := 1
	for i := range fieldRules {
		if strings.Contains(fieldRules[i].Name, "departure") {
			v := your_ticket[fieldRules[i].TicketIndex]
			if (v >= fieldRules[i].Min1 && v <= fieldRules[i].Max1) || (v >= fieldRules[i].Min2 && v <= fieldRules[i].Max2) {
				fmt.Println(depart_multi)
				fmt.Println(v)
				depart_multi = depart_multi * v
			} else {
				fmt.Println("ERRORS")
			}
		}
	}
	logging("---Answer Part2---", "", true)
	logging("ticket_invalid_rules", ticket_invalid_rules, true)
	logging("desc_len_order", desc_len_order, true)
	for i := range fieldRules {
		logging("fieldRules[i]", fieldRules[i], true)
	}
	logging("depart_multi", depart_multi, true)
}

func findNum(input []int, num int) int {
	for i := range input {
		if input[i] == num {
			return i
		}
	}
	return -1
}

func parseFields(s string) field {
	var newfield = field{}
	newfield.Name, _ = getFieldName(s)
	newfield.Min1, newfield.Max1, newfield.Min2, newfield.Max2 = getFieldMinMax(s)
	newfield.TicketIndex = -1
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
	if len(match) > 2 {
		v1, _ := strconv.Atoi(match[0][0])
		v2, _ := strconv.Atoi(match[1][0])
		v3, _ := strconv.Atoi(match[2][0])
		v4, _ := strconv.Atoi(match[3][0])
		return v1, v2, v3, v4
	}
	return 0, 0, 0, 0
}

// removeIndex. Swap the element to delete with the one at the end of the slice and then return the n-1 first elements
func removeIndex(s [][]int, i int) [][]int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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
