package main

// https://adventofcode.com/2020/day/23

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

func main() {
	part := "B"

	start := time.Now()
	defer timeTrack(time.Now(), "day23")
	input := "562893147"
	// input := "389125467"
	s := strings.Split(input, "")
	var cups []int
	for i := range s {
		v, _ := getInt(s[i])
		cups = append(cups, v)
	}

	if part == "B" {
		iterations := 10000000
		for i := 10; i <= 1000000; i++ {
			cups = append(cups, i)
		}
	} else {
		iterations := 100
	}

	var dest int
	var c3 []int
	c_cnt := len(cups)
	fmt.Println("c_cnt=", c_cnt)

	it := 0
	idx := -1
	for it < iterations {
		it = it + 1
		_, mod := divMod(it, 1000)
		if mod == 0 {
			t := time.Now()
			fmt.Println("it=", it, "  ", (it/iterations)*100, "%", "  ", t.Sub(start))
		}
		// Select next index cup
		if idx == len(cups)-1 {
			idx = 0
		} else {
			idx = idx + 1
		}

		// remove 3 cups after selected cup
		cups, c3, idx = threeCups(cups, idx)

		// find dest cup
		dest = destCup(cups, idx, c_cnt)
		// place 3 removed cups after dest cup
		cups = appendIdxSOrdered(cups, dest, c3)
		if dest < idx {
			_, idx = divMod(idx+3, c_cnt)
		}
		// repeat on next index cup
	}

	var idx_1 int
	for i := range cups {
		if cups[i] == 1 {
			idx_1 = i + 1
			break
		}
	}
	if part == "A" {
		fmt.Println("--- Answer Part 1 ---")
		ans1 := ""
		for len(ans1) < c_cnt {
			if idx_1 == c_cnt {
				idx_1 = 0
			}
			ans1 = ans1 + fmt.Sprint(cups[idx_1])
			idx_1 = idx_1 + 1
		}
		fmt.Println("ans1:", ans1)
	}
	if part == "B" {
		fmt.Println("--- Answer Part 2 ---")
		fmt.Println(cups[idx_1], "*", cups[idx_1+1], "=", cups[idx_1]*cups[idx_1+1])
	}
}

func threeCups(s []int, idx int) ([]int, []int, int) {
	var t []int
	var i = idx + 1

	for len(t) < 3 {
		if i >= len(s) || i < 0 {
			i = 0
		}
		if i < idx {
			idx = idx - 1
		}
		t = append(t, s[i])
		s = removeIdxSOrdered(s, i)
	}

	return s, t, idx
}

func destCup(s []int, idx int, c_cnt int) int {
	var dest = s[idx]
	found := false
	for !found {
		if dest == 1 {
			dest = c_cnt
		} else {
			dest = dest - 1
		}

		for j := range s {
			if s[j] == dest {
				found = true
				return j
			}
		}
	}
	return dest
}

// removeIdxSOrdered. Realocate all elements so that element at index i is discarded
func removeIdxSOrdered(s []int, i int) []int {
	return append(s[:i], s[i+1:]...)
}

// removeIdxSOrdered. Realocate all elements so that slice a is appended at element index i
func appendIdxSOrdered(s []int, i int, a []int) []int {
	var s_new []int
	s_new = append(s_new, s[:i+1]...)
	s_new = append(s_new, a...)
	return append(s_new, s[i+1:]...)
}

// removeIndex. Swap the element to delete with the one at the end of the slice and then return the n-1 first elements
func removeIndexSS(s [][]int, i int) [][]int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getInt(s string) (int, error) {
	matches := 20
	re := regexp.MustCompile(`[0-9]+`)
	match := re.FindAllStringSubmatch(s, matches)
	var v int
	var err error
	if len(match) == 1 {
		for i := range match {
			v, err = strconv.Atoi(match[i][0])
			if err != nil {
				break
			} else {
				return v, nil
			}
		}
	}
	return v, errors.New("No Integers")
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
