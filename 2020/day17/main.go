package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type coord struct {
	X int
	Y int
	Z int
	W int
}

func main() {
	part := "B"
	defer timeTrack(time.Now(), "day17")
	input, _ := readLines("config.txt")
	// input, _ := readLines("config-test.txt")

	active := make([]coord, 0)
	for row := range input {
		s := strings.Split(input[row], "")
		for i := range s {
			if s[i] == "#" {
				var active_coord = coord{i - (len(s) - 2), row - (len(s) - 2), 0, 0}
				active = append(active, active_coord)
			}
		}
	}

	logging("active", active, true)

	var active_len int = len(active)
	logging("active_len", active_len, true)
	var activeArray [10000]coord
	copy(activeArray[:], active)

	total_iterations := 6
	iteration := 0

	for iteration < total_iterations {
		iteration = iteration + 1
		logging("---iteration---", iteration, true)
		new_active := make([]coord, 0)
		cube := getActiveCube(&activeArray, active_len)

		for x := -cube; x <= cube; x++ {
			for y := -cube; y <= cube; y++ {
				for z := -cube; z <= cube; z++ {
					for w := -cube; w <= cube; w++ {
						if part == "A" && w != 0 {
							continue
						}
						n := countActiveN(&activeArray, x, y, z, w, active_len)
						a := isActive(&activeArray, x, y, z, w, active_len)
						if a {
							if n == 2 || n == 3 {
								a = true
							} else {
								a = false
							}
						} else {
							if n == 3 {
								a = true
							} else {
								a = false
							}
						}
						if a {
							new_active = append(new_active, coord{x, y, z, w})
						}
					}
				}
			}
		}
		active_len = len(new_active)
		if active_len == 10000 {
			panic("active_len == 10000")
		}
		logging("active_len", active_len, true)
		// logging("new_active", new_active, true)
		activeArray = [10000]coord{}
		copy(activeArray[:], new_active)
		new_active = nil
		// logging("new_active", new_active, true)
	}
}

// countActiveN counts the number of active neighbours to a co-ord
func countActiveN(ss *[10000]coord, x int, y int, z int, w int, r int) int {
	n := 0
	for i := range ss {
		if i == r {
			break
		}
		for xx := -1; xx <= 1; xx++ {
			for yy := -1; yy <= 1; yy++ {
				for zz := -1; zz <= 1; zz++ {
					for ww := -1; ww <= 1; ww++ {
						if xx == 0 && xx == yy && xx == zz && xx == ww {
							continue
						}
						if Abs(ss[i].X-x) == xx && Abs(ss[i].Y-y) == yy && Abs(ss[i].Z-z) == zz && Abs(ss[i].W-w) == ww {
							n = n + 1
						}
					}
				}
			}
		}
	}
	return n
}

// isActive checks if the current co-ord is active
func isActive(ss *[10000]coord, x int, y int, z int, w int, r int) bool {
	for i := range ss {
		if i == r {
			break
		}
		if ss[i].X == x && ss[i].Y == y && ss[i].Z == z && ss[i].W == w {
			return true
		}
	}
	return false
}

// getActiveCube returns the current range of co-ord + 1 for the next iteration
func getActiveCube(ss *[10000]coord, r int) int {
	var max int
	for i := range ss {
		if i == r {
			break
		}
		var s = []int{ss[i].X, ss[i].Y, ss[i].Z, ss[i].W}
		for j := range s {
			ij := Abs(s[j])
			if ij > max {
				max = ij
			}
		}
	}
	logging("max", max, true)
	return max + 1
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
