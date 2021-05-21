package main

// https://adventofcode.com/2020/day/20

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

type tile struct {
	ID       int
	Sides    [][]int // [[Top][Right][Bottom][Left]]
	Centre   [][]string
	SMatched []int
	Rotation int
	Flip     bool
}

func main() {
	// part := "B"

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day20")
	input, _ := readLines("tiles.txt")
	// input, _ := readLines("tiles-test.txt")

	var side_len int
	var tiles []tile
	var cur_tile tile
	for i := range input {
		_, mod := divMod(i, side_len+2)
		fmt.Println("input[i], i=", i, "mod=", mod, input[i])
		if i == 0 || mod == 0 {
			if i != 0 {
				tiles = append(tiles, cur_tile)
			}
			cur_tile = tile{}
			id, err := getID(input[i])
			if err != nil {
				fmt.Println("ERROR at getID, i=", i)
				continue
			} else {
				cur_tile.ID = id
			}
		}
		if mod > 0 {
			var pix []int
			s := strings.Split(input[i], "")
			side_len = len(s)
			if mod == 1 {
				pix = getNonEmptyPix(s, true)
				cur_tile.Sides = append(cur_tile.Sides, pix)
				cur_tile.Sides = append(cur_tile.Sides, make([]int, 0))
				cur_tile.Sides = append(cur_tile.Sides, make([]int, 0))
				cur_tile.Sides = append(cur_tile.Sides, make([]int, 0))
			} else if mod == side_len {
				pix = getNonEmptyPix(s, true)
				cur_tile.Sides[2] = append(cur_tile.Sides[2], pix...)
			} else {
				pix = getNonEmptyPix(s, false)
			}
			if len(pix) > 0 {
				if pix[0] == 0 {
					cur_tile.Sides[1] = append(cur_tile.Sides[1], mod-1)
				}
				if pix[len(pix)-1] == side_len-1 {
					cur_tile.Sides[3] = append(cur_tile.Sides[3], mod-1)
				}
			}
		}
		if i == len(input)-1 {
			tiles = append(tiles, cur_tile)
		}
	}
	fmt.Println("side_len=", side_len)
	for i := range tiles {
		fmt.Println(tiles[i])
	}

	// var adjacent_cnt []int
	// var pos_m [][]int
	for i := range tiles {
		for j := range tiles {
			if i == j {
				continue
			}
			for si := range tiles[i].Sides {
				for sj := range tiles[j].Sides {
					s1 := tiles[i].Sides[si]
					s2 := tiles[j].Sides[sj]
					if len(s1) != len(s2) {
						continue
					}
					match1 := true
					match2 := true
					if len(s1) > 0 && 0 < len(s2) {
						for ii := range s1 {
							if s1[ii] != s2[ii] {
								match1 = false
								break
							}
						}
						for ii := range s1 {
							if s1[ii] != Abs(s2[len(s2)-1-ii]-side_len+1) {
								match2 = false
								break
							}
						}
					}
					if match1 {
						fmt.Println(s1, "==", s2)
						tiles[i].SMatched = append(tiles[i].SMatched, si)
					}
					if match2 {
						fmt.Println(s1, "== -", s2)
						tiles[i].SMatched = append(tiles[i].SMatched, -si)
					}
				}
			}
		}
	}
	corner_multi := 1
	corner_cnt := 0
	for i := range tiles {
		fmt.Println(tiles[i])
		if len(tiles[i].SMatched) == 2 {
			corner_multi = corner_multi * tiles[i].ID
			corner_cnt = corner_cnt + 1
		}
	}
	if corner_cnt > 4 {
		fmt.Println("ERROR")
		fmt.Println("corner_cnt=", corner_cnt, "> 4")
	}
	fmt.Println("------ Answer Part 1 ------")
	fmt.Println("corner_multi=", corner_multi)
}

func getNonEmptyPix(s []string, full bool) []int {
	var results []int
	for i := range s {
		if !full && i > 0 && i < len(s)-1 {
			continue
		}
		if s[i] == "#" {
			results = append(results, i)
		}
	}
	return results
}

func getID(s string) (int, error) {
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
