package main

// https://adventofcode.com/2020/day/14

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	// part := "A"

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day14")
	// input, _ := readLines("initialization.txt")
	input, _ := readLines("initialization-test.txt")

	var masks []string
	var mem_init [][]int
	var mem_init_val [][]int
	var mem_mask []int

	for i := range input {
		s, err := getMask(input[i])
		if err == nil {
			masks = append(masks, s)
			mem_mask = append(mem_mask, len(masks)-1)
			mem_init = append(mem_init, make([]int, 0))
			mem_init_val = append(mem_init_val, make([]int, 0))

		} else {
			val, _ := getMemory(input[i])
			mem_init[mem_mask[len(mem_mask)-1]] = append(mem_init[mem_mask[len(mem_mask)-1]], val)
			val, _ = getMemVal(input[i])
			mem_init_val[mem_mask[len(mem_mask)-1]] = append(mem_init_val[mem_mask[len(mem_mask)-1]], val)
		}
	}

	fmt.Println(masks)
	fmt.Println(mem_init)
	fmt.Println(mem_init_val)
	fmt.Println(mem_mask)

	// var memory []int
	// var memory_values []int

	// s := fmt.Sprintf("%036b", 123)
	// i, err := strconv.ParseInt(s, 2, 64)

}

func getMask(s string) (string, error) {
	re := regexp.MustCompile(`mask = ([X10]+?)`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		return match[1], nil
	}
	return s, errors.New("error")
}

func getMemory(s string) (int, error) {
	re := regexp.MustCompile(`(?:mem[)([0-9]+?)(?:])`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			return v, nil
		} else {
			fmt.Println(err)
		}
	}
	return 0, errors.New("error")
}

func getMemVal(s string) (int, error) {
	re := regexp.MustCompile(`(?:] = )([0-9]+?)$`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			return v, nil
		} else {
			fmt.Println(err)
		}
	}
	return 0, errors.New("error")
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
