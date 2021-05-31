package main

// https://adventofcode.com/2020/day/24

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
	N     int
	E     int
	Black bool
}

func main() {
	part := "B"

	defer timeTrack(time.Now(), "day24")
	// input, _ := readLines("tileflip-test.txt")
	input, _ := readLines("tileflip.txt")

	// {{e},{se},{sw},{w},{nw},{ne}}
	var adj = [][]int{{0, 2}, {-1, 1}, {-1, -1}, {0, -2}, {1, -1}, {1, 1}}
	fmt.Println(adj)
	var tiles []tile

	for i := range input {
		idx := 0
		var cur_t [2]int
		s := strings.Split(input[i], "")
		for idx < len(s) {
			switch s[idx] {
			case "s":
				if s[idx+1] == "e" {
					cur_t[0] = cur_t[0] - 1
					cur_t[1] = cur_t[1] + 1
				} else {
					cur_t[0] = cur_t[0] - 1
					cur_t[1] = cur_t[1] - 1
				}
				idx = idx + 1
			case "n":
				if s[idx+1] == "e" {
					cur_t[0] = cur_t[0] + 1
					cur_t[1] = cur_t[1] + 1
				} else {
					cur_t[0] = cur_t[0] + 1
					cur_t[1] = cur_t[1] - 1
				}
				idx = idx + 1
			case "e":
				cur_t[1] = cur_t[1] + 2
			case "w":
				cur_t[1] = cur_t[1] - 2
			}
			idx = idx + 1
		}
		found := false
		for ti, t := range tiles {
			if t.N == cur_t[0] && t.E == cur_t[1] {
				if t.Black {
					tiles[ti].Black = false
				} else if !t.Black {
					tiles[ti].Black = true
				}
				found = true
				break
			}
		}
		if !found {
			tiles = append(tiles, tile{cur_t[0], cur_t[1], true})
		}
	}
	black_tiles := 0
	for _, t := range tiles {
		fmt.Println(t)
		if t.Black {
			black_tiles = black_tiles + 1
		}
	}

	fmt.Println("\n--- Answer Part 1 ---")
	fmt.Println("black_tiles=", black_tiles)

	if part == "B" {
		fmt.Println("\n--- Answer Part 2 ---")
		days := 100
		for d := 1; d <= days; d++ {
			var update_tiles []tile
			black_tiles = 0
			tile_cnt := len(tiles)

			for _, t := range tiles {
				var n_found []int

				for i, tn := range tiles {
					if t.N == tn.N && t.E == tn.E {
						continue
					}
					if (Abs(t.N-tn.N) == 1 && Abs(t.E-tn.E) == 1) || (Abs(t.N-tn.N) == 0 && Abs(t.E-tn.E) == 2) {
						n_found = append(n_found, i)
					}
				}

				if len(n_found) > 6 {
					fmt.Println("n_found is greater than 6=", n_found)
					panic("n_found is greater than 6")
				}

				adj_black := 0
				for _, a := range adj {
					new_tile := true
					for _, j := range n_found {
						if tiles[j].N+a[0] == t.N && tiles[j].E+a[1] == t.E {
							new_tile = false
							if tiles[j].Black {
								adj_black = adj_black + 1
							}
							break
						}
					}
					if new_tile {
						tiles = append(tiles, tile{t.N - a[0], t.E - a[1], false})
					}
				}

				if adj_black == 2 && !t.Black {
					update_tiles = append(update_tiles, tile{t.N, t.E, true})
				} else if (adj_black == 0 || adj_black > 2) && t.Black {
					update_tiles = append(update_tiles, tile{t.N, t.E, false})
				}
			}

			for _, t := range tiles[tile_cnt:] {
				var n_found []int
				for i, tn := range tiles {
					if t.N == tn.N && t.E == tn.E {
						continue
					}
					if (Abs(t.N-tn.N) == 1 && Abs(t.E-tn.E) == 1) || (Abs(t.N-tn.N) == 0 && Abs(t.E-tn.E) == 2) {
						n_found = append(n_found, i)
					}
				}

				if len(n_found) > 6 {
					fmt.Println("n_found is greater than 6=", n_found)
					panic("n_found is greater than 6")
				}

				adj_black := 0
				for _, a := range adj {
					for _, j := range n_found {
						if tiles[j].N+a[0] == t.N && tiles[j].E+a[1] == t.E {
							if tiles[j].Black {
								adj_black = adj_black + 1
							}
							break
						}
					}
				}

				if adj_black == 2 && !t.Black {
					update_tiles = append(update_tiles, tile{t.N, t.E, true})
				} else if (adj_black == 0 || adj_black > 2) && t.Black {
					update_tiles = append(update_tiles, tile{t.N, t.E, false})
				}
			}

			for _, ut := range update_tiles {
				for ti := range tiles {
					if ut.N == tiles[ti].N && ut.E == tiles[ti].E {
						tiles[ti].Black = ut.Black
					}
				}
			}
			for _, t := range tiles {
				if t.Black {
					black_tiles = black_tiles + 1
				}
			}
			fmt.Println("Day", d, ":", black_tiles)
		}
	}
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
