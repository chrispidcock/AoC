package main

// https://adventofcode.com/2020/day/18

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type value struct {
	Index int
	Value int
	Depth int
}

type operation struct {
	Index int
	Op    string
	Depth int
}

type bracket struct {
	MinI  int
	MaxI  int
	Depth int
}

type elements struct {
	index   int
	element string
	Depth   int
}

func main() {
	part := "B"
	defer timeTrack(time.Now(), "day18")
	input, _ := readLines("homework.txt")
	// input, _ := readLines("homework-test.txt")

	result_sum := 0
	for row := range input {
		eq := []elements{}
		brackets := [][]bracket{}
		values := [][]value{}
		operations := [][]operation{}
		s := strings.Split(input[row], "")
		depth := 0
		brackets = append(brackets, []bracket{})
		brackets[depth] = append(brackets[depth], bracket{-1, len(s), depth})
		i := -1
		for i < len(s)-1 {
			i = i + 1
			if depth > len(values)-1 {
				values = append(values, []value{})
			}
			if depth > len(operations)-1 {
				operations = append(operations, []operation{})
			}
			switch s[i] {
			case " ":
				continue
			case "(":
				depth = depth + 1
				for len(brackets)-1 < depth {
					brackets = append(brackets, []bracket{})
				}
				brackets[depth] = append(brackets[depth], bracket{i, -1, depth})
				eq = append(eq, elements{i, "(", depth})
			case ")":
				brackets[depth][len(brackets[depth])-1].MaxI = i
				depth = depth - 1
			case "+":
				operations[depth] = append(operations[depth], operation{i, "+", depth})
				eq = append(eq, elements{i, "+", depth})
			case "*":
				operations[depth] = append(operations[depth], operation{i, "*", depth})
				eq = append(eq, elements{i, "*", depth})
			default:
				start_idx := i
				end_idx := start_idx - 1
				vs := ""
				for end_idx < len(s)-1 {
					end_idx = end_idx + 1
					_, err := strconv.Atoi(s[end_idx])
					if err != nil {
						end_idx = end_idx - 1
						break
					}
				}
				for j := start_idx; j <= end_idx; j++ {
					vs = vs + s[j]
				}
				v, _ := strconv.Atoi(vs)
				i = end_idx
				values[depth] = append(values[depth], value{start_idx, v, depth})
				eq = append(eq, elements{start_idx, "i", depth})
			}
		}
		logging("", "", true)
		logging("-----s------", s, true)
		logging("len(s)", len(s), true)
		logging("eq", eq, true)
		logging("brackets", brackets, true)
		logging("values", values, true)
		logging("operations", operations, true)

		result := part1Result(brackets, values, operations, part)
		result_sum = result_sum + result
		logging(s, result, true)
		logging("result_sum", result_sum, true)
	}
	logging("result_sum", result_sum, true)
}

func part1Result(brackets [][]bracket, values [][]value, operations [][]operation, part string) int {
	result := 0
	for i := len(brackets) - 1; i >= 0; i-- {
		for j := range brackets[i] {
			var op string
			var v2 int
			var err error
			var v int
			var index int
			if part == "B" {
				for len(values[i]) > 0 && len(operations[i]) > 0 {
					var v3 int
					var sum_idx int
					operations[i], op, sum_idx, err = lowestIndexAdd(operations[i], brackets[i][j].MinI, brackets[i][j].MaxI)
					if err != nil {
						fmt.Println(err)
						break
					}
					values[i], v2, sum_idx, err = valesToSum(values[i], brackets[i][j].MinI, brackets[i][j].MaxI, sum_idx, -1)
					if err != nil {
						fmt.Println(err)
					}
					values[i], v3, sum_idx, err = valesToSum(values[i], brackets[i][j].MinI, brackets[i][j].MaxI, sum_idx, 1)
					if err != nil {
						fmt.Println(err)
					}
					v, _ = sumOrMulti(v2, v3, "+")
					values[i] = append(values[i], value{sum_idx, v, i})
					v = 0
				}
			}
			for len(values[i]) > 0 && len(operations[i]) > 0 {
				if v == 0 {
					values[i], v, index, err = lowestIndexValue(values[i], brackets[i][j].MinI, brackets[i][j].MaxI)
					if err != nil {
						fmt.Println(err)
					}
				}
				operations[i], op, index, err = lowestIndexOp(operations[i], index, brackets[i][j].MaxI)
				if err != nil {
					fmt.Println(err)
					break
				}
				values[i], v2, index, err = lowestIndexValue(values[i], index, brackets[i][j].MaxI)
				if err != nil {
					fmt.Println(err)
					break
				}
				v, _ = sumOrMulti(v, v2, op)
			}
			if len(values[i]) > 0 && len(operations[i]) == 0 && i != 0 {
				for remain := range values[i] {
					values[i-1] = append(values[i-1], value{values[i][remain].Index, values[i][remain].Value, i - 1})
				}
			}
			if i != 0 && v != 0 {
				values[i-1] = append(values[i-1], value{brackets[i][j].MinI, v, i - 1})
			} else if v != 0 {
				result = v
				break
			} else if v == 0 {
				result = values[i][0].Value
				break
			} else {
				fmt.Println("ERROR: Return No Value!!!!!")
			}
		}
	}
	return result
}

func lowestIndexValue(s []value, minI int, maxI int) ([]value, int, int, error) {
	index_diff := 999
	var v int
	var index int
	for i := range s {
		if s[i].Index <= minI || s[i].Index >= maxI {
			continue
		}
		cur_diff := s[i].Index - minI
		if cur_diff < index_diff && cur_diff > 0 {
			index_diff = cur_diff
			v = s[i].Value
			index = i
		}
	}
	if index_diff+minI < maxI && index_diff > 0 {
		s = removeIndexValue(s, index)
		return s, v, index_diff + minI, nil
	}
	return s, v, maxI, errors.New("lowestIndexValue: maxI out of range")
}

func lowestIndexOp(s []operation, minI int, maxI int) ([]operation, string, int, error) {
	index_diff := 999
	var op string
	var index int
	for i := range s {
		if s[i].Index <= minI || s[i].Index >= maxI {
			continue
		}
		cur_diff := s[i].Index - minI
		if cur_diff < index_diff && cur_diff > 0 {
			index_diff = cur_diff
			op = s[i].Op
			index = i
		}
	}
	if index_diff+minI < maxI && index_diff > 0 {
		s = removeIndexOperation(s, index)
		return s, op, index_diff + minI, nil
	}
	return s, op, maxI, errors.New("lowestIndexOp: maxI out of range")
}

func valesToSum(s []value, minI int, maxI int, idx int, dir int) ([]value, int, int, error) {
	index_diff := 999
	var v int
	var index int
	for i := range s {
		if s[i].Index <= minI || s[i].Index >= maxI {
			continue
		}
		cur_diff := (s[i].Index - idx) * dir
		if cur_diff < index_diff && cur_diff > 0 {
			index_diff = cur_diff
			v = s[i].Value
			index = i
		}
	}
	if index_diff != 999 {
		s = removeIndexValue(s, index)
		return s, v, idx, nil
	}
	return s, v, idx, errors.New("valesToSum: index out of range")
}

func lowestIndexAdd(s []operation, minI int, maxI int) ([]operation, string, int, error) {
	index_diff := 999
	var op string
	var index int
	for i := range s {
		if s[i].Index <= minI || s[i].Index >= maxI || s[i].Op != "+" {
			continue
		}
		cur_diff := s[i].Index - minI
		if cur_diff < index_diff && cur_diff > 0 {
			index_diff = cur_diff
			op = s[i].Op
			index = i
		}
	}
	if index_diff+minI < maxI && index_diff > 0 && index_diff != 999 {
		s = removeIndexOperation(s, index)
		return s, op, index_diff + minI, nil
	}
	return s, op, minI, errors.New("lowestIndexAdd: index out of range")
}

func sumOrMulti(i1 int, i2 int, op string) (int, error) {
	if op == "*" {
		fmt.Println(i1, " * ", i2)
		return i1 * i2, nil
	} else if op == "+" {
		fmt.Println(i1, " + ", i2)
		return i1 + i2, nil
	}
	return 0, nil
}

// removeIndex. Swap the element to delete with the one at the end of the slice and then return the n-1 first elements
func removeIndexValue(s []value, i int) []value {
	if len(s) == 0 {
		return s
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
func removeIndexOperation(s []operation, i int) []operation {
	if len(s) == 0 {
		return s
	}
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
