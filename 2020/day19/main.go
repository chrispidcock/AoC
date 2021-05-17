package main

// https://adventofcode.com/2020/day/19

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

type rule struct {
	num     int
	links   [][]int
	matches []string
}

func main() {
	// part := "B"
	defer timeTrack(time.Now(), "day19")
	input, _ := readLines("msgs.txt")
	// input, _ := readLines("msgs-test.txt")
	rule_count := 0
	rules := [150]rule{}
	rules_parsed := false
	var msgs []string
	for row := range input {
		if input[row] == "" {
			rules_parsed = true
			continue
		}
		if !rules_parsed {
			s := strings.Split(input[row], ": ")
			num, _ := strconv.Atoi(s[0])
			rules[num].num = num
			si := strings.Split(s[1], " | ")
			for i := range si {
				v1, v2, err := getLinks(si[i])
				if err != nil {
					a := strings.Trim(si[i], "\"")
					rules[num].matches = append(rules[num].matches, a)
				} else {
					rules[num].links = append(rules[num].links, []int{v1, v2})
				}
			}
			rule_count = rule_count + 1
		}
		if rules_parsed {
			msgs = append(msgs, input[row])
		}
	}
	logging("rules", rules, true)
	logging("msgs", msgs, true)

	for i := 0; i < rule_count; i++ {
		genMatches(i, &rules)
		logging(i, rules[i].matches, true)
	}

	// var valid []string
}

func genMatches(num int, ra *[150]rule) []string {
	var new_matches []string
	if len(ra[num].matches) == 0 {
		for r := range ra[num].links {
			// eg. m.links = 2 3 | 3 2
			// r = [2 3]
			var p_matches []string
			for i := range ra[num].links[r] {
				matches := genMatches(ra[num].links[r][i], ra)
				p_matches = append(p_matches, matches...)
			}
			for i := range p_matches[0] {
				for j := range p_matches[1] {
					var match string
					match = string(p_matches[0][i]) + string(p_matches[1][j])
					new_matches = append(new_matches, match)
				}
			}
			ra[num].matches = append(ra[num].matches, new_matches...)
		}
	}
	return ra[num].matches
}

func getLinks(s string) (int, int, error) {
	re := regexp.MustCompile(`[0-9]+`)
	match := re.FindAllStringSubmatch(s, 2)
	if len(match) > 1 {
		v1, _ := strconv.Atoi(match[0][0])
		v2, _ := strconv.Atoi(match[1][0])
		return v1, v2, nil
	}
	return 0, 0, errors.New("No Integers")
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

func logging(name interface{}, value interface{}, active bool) {
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
