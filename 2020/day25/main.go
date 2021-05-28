package main

// https://adventofcode.com/2020/day/25

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Game struct {
	rounds []Round
}

type Round struct {
	players [2]Player
}

type Player struct {
	hand []int
}

func main() {
	// part := "B"

	defer timeTrack(time.Now(), "day25")
	// input, _ := readLines("public-keys.txt")
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

func threeCups(s []int, idx int) ([]int, []int, int) {
	var t [3]int
	var i = idx

	for idx {

	}

	return s, t, idx
}

func destCup(s []int, idx int) ([]int, []int) {
	var t [3]int
	var i = idx

	for i < len(s) {

	}

	return s, t
}

// removeIndex. Swap the element to delete with the one at the end of the slice and then return the n-1 first elements
func removeIndexSS(s [][]int, i int) [][]int {
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
