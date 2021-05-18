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
	part := "A"
	defer timeTrack(time.Now(), "day19")
	input, _ := readLines("msgs.txt")
	// input, _ := readLines("msgs-test.txt")
	// input, _ := readLines("msgs-test2.txt")
	// input, _ := readLines("msgs-test-part2.txt")
	rule_max := 0
	max_msg_len := 0
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
				links, err := getLinks(si[i])
				if err != nil {
					a := strings.Trim(si[i], "\"")
					rules[num].matches = append(rules[num].matches, a)
				} else {
					rules[num].links = append(rules[num].links, links)
				}
			}
			if num > rule_max {
				rule_max = num
			}
		}
		if rules_parsed {
			msgs = append(msgs, input[row])
			if max_msg_len < len(input[row]) {
				max_msg_len = len(input[row])
			}
		}
	}

	logging("msgs", msgs, true)
	logging("---rules---", "", true)
	for i := 0; i <= rule_max; i++ {
		if rules[i].num != i {
			continue
		}
		logging(i, rules[i], true)
	}

	logging("---max_msg_len---", max_msg_len, true)
	logging("---matches---", "", true)
	for i := rule_max; i >= 0; i-- {
		if rules[i].num != i {
			continue
		}
		logging(i, rules[i].links, true)
		if part == "B" {
			// min 42 42 31 // len = 8 * 3
			// max // len = 8 * 12 (max_msg_len = 96)
			if rules[i].num == 11 || rules[i].num == 8 || rules[i].num == 0 {
				continue
			}
		}
		GenMatches(i, &rules)
	}

	valid_msgs := 0
	logging("---Valid msgs---", "", true)
	if part == "A" {
		for m := range msgs {
			for i := range rules[0].matches {
				if msgs[m] == rules[0].matches[i] {
					valid_msgs = valid_msgs + 1
					logging(msgs[m], valid_msgs, true)
					break
				}
			}
		}
	}
	if part == "B" {
		msg_seg_len := len(rules[42].matches[0])
		for i := range msgs {
			valid := ValidateMsgSegments(msgs[i], &rules, msg_seg_len)
			if valid {
				valid_msgs = valid_msgs + 1
			}
		}
	}

	logging("valid_msgs", valid_msgs, true)

}

func GenMatches(num int, ra *[150]rule) []string {
	// logging("GenMatches", num, true)
	var new_matches []string
	// logging("-- Getting matches from links ", ra[num].links, true)
	if len(ra[num].matches) == 0 {
		// eg. m.links = 2 3 | 3 2
		for r := range ra[num].links {
			var p_matches [][]string
			// p_matches eg. [aaa aba bbb]
			for i := range ra[num].links[r] {
				// r[i] = [2 3]
				matches := GenMatches(ra[num].links[r][i], ra)
				// logging("matches", matches, true)
				p_matches = append(p_matches, matches)
			}
			// logging("p_matches", p_matches, true)
			new_matches = ListRecursive(0, p_matches)
			// logging("new_matches", new_matches, true)
			ra[num].matches = append(ra[num].matches, new_matches...)
		}

	} else {
		return ra[num].matches
	}
	return ra[num].matches
}

// ListRecursive combines all combinations from each list, taking 1 element from each list
func ListRecursive(i int, matches [][]string) []string {
	var new_matches []string
	for j := range matches[i] {
		if i+1 < len(matches) {
			other := ListRecursive(i+1, matches)
			for m := range other {
				new_matches = append(new_matches, string(matches[i][j]+other[m]))
			}
		} else {
			new_matches = append(new_matches, matches[i][j])
		}
	}
	return new_matches
}

// ListRecursive combines all combinations from each list, taking 1 element from each list
func ValidateMsgSegments(msg string, ra *[150]rule, msg_seg_len int) bool {
	logging("Checking:", msg, true)
	segments, mod := divMod(len(msg), msg_seg_len)
	if mod != 0 {
		return false
	}
	for i := 1; i < segments; i++ {
		s31, mod := divMod(segments-i, 2)
		s42 := i + mod + s31

		// Some obvious violation checking
		if s31+s42 != segments {
			fmt.Println("s31 + s42 != segments", s31+s42, " != ", segments)
			break
		} else if s31 < 1 {
			fmt.Println("s31 < 1", s31)
			break
		} else if s42 < 1 {
			fmt.Println("s31 < 1", s31)
			break
		}

		var r_strs [][]string
		var r_num []int
		for i := 1; i <= s42; i++ {
			r_strs = append(r_strs, ra[42].matches)
			r_num = append(r_num, 42)
		}
		for i := 1; i <= s31; i++ {
			r_strs = append(r_strs, ra[31].matches)
			r_num = append(r_num, 31)
		}
		logging("r_num", r_num, true)

		var match_index [20]int
		var match_rule [20]int

		for j := range r_strs {
			if match_rule[j] == 42 {
				continue
			}
			for p := range r_strs[j] {
				if msg[msg_seg_len*j:msg_seg_len*(j+1)] == r_strs[j][p] {
					match_index[j] = p
					match_rule[j] = r_num[j]
					if j != len(r_strs)-1 {
						break
					} else {
						match_str := ""
						for k := range r_num {
							match_str = match_str + ra[match_rule[k]].matches[match_index[k]]
						}
						fmt.Println("MATCH: ", match_str, msg)
						return true
					}
				}
			}
			if match_rule[j] == 0 {
				break
			}
		}
	}
	return false
}

func RecursiveMsgCheck(str string, level int, sss [][]string, msg string) bool {
	// logging("level", level, true)
	for j := range sss[level-1] {
		if level == len(sss) {
			// logging("str+sss[level-1][j]", str+sss[level-1][j], true)
			if msg == str+sss[level-1][j] {
				logging("match", msg, true)
				return true
			}
		} else {
			RecursiveMsgCheck(str+sss[level-1][j], level+1, sss, msg)
		}
	}
	return false
}

func getLinks(s string) ([]int, error) {
	matches := 20
	re := regexp.MustCompile(`[0-9]+`)
	match := re.FindAllStringSubmatch(s, matches)
	var vs []int
	var err error
	if len(match) == matches {
		return vs, errors.New("Suggest increasing matches:" + string(rune(matches)))
	}
	if len(match) > 0 {
		for i := range match {
			v, err := strconv.Atoi(match[i][0])
			if err != nil {
				break
			} else {
				vs = append(vs, v)
			}
		}
		return vs, err
	}
	return vs, errors.New("No Integers")
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

func loggingSlice(name interface{}, value []interface{}, active bool) {
	for i := range value {
		logging(i, value[i], true)
	}
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
