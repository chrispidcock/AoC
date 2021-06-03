package main

// https://adventofcode.com/2020/day/25

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
	input, _ := readLines("public-keys.txt")
	// input, _ := readLines("public-keys-test.txt")

	var public_keys []int

	for i := range input {
		v, err := strconv.Atoi(input[i])
		if err != nil {
			panic("public key is not an integer")
		}
		public_keys = append(public_keys, v)
	}

	fmt.Println("public_keys=", public_keys)
	// var c_d [][]int
	var loops []int

	for pk := range public_keys {
		fmt.Println("pk=", pk)
		// c_d = append(c_d, make([]int, 0))
		// c_d[pk] = append(c_d[pk], public_keys[len(public_keys)-1-pk])
		done := false
		loop_size := 0
		sub_num := 7
		val := 1

		for !done {
			val = loop(val, sub_num)
			loop_size = loop_size + 1
			if val == public_keys[pk] {
				loops = append(loops, loop_size)
				done = true
			}
		}
	}
	fmt.Println("loops=", loops)

	var encr_key int
	for pk := range public_keys {
		fmt.Println("pk=", pk)
		loop_size := loops[pk]
		sub_num := public_keys[len(public_keys)-1-pk]
		val := 1

		for i := 0; i < loop_size; i++ {
			val = loop(val, sub_num)
		}
		encr_key = val
	}
	fmt.Println("loops=", loops)

	// done := false
	// var encr_key int
	// for !done {
	// 	for pk := range public_keys {
	// 		sub_num := c_d[pk][len(c_d[pk])-1]
	// 		val := loops[pk]

	// 		val = loop(val, sub_num)
	// 		c_d[pk] = append(c_d[pk], val)
	// 		// fmt.Println(c_d)
	// 		for i := range c_d[len(public_keys)-1-pk] {
	// 			if val == c_d[len(public_keys)-1-pk][i] {
	// 				encr_key = val
	// 				done = true
	// 			}
	// 		}
	// 	}
	// }

	fmt.Println("\n--- Answer Part 1 ---")
	fmt.Println("encr_key =", encr_key)
}

func loop(i int, sub int) int {
	i = sub * i
	_, i = divMod(i, 20201227)
	return i
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
