package main

// https://adventofcode.com/2020/day/20

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type tile struct {
	ID       int
	Sides    [][]int // [[Top][Right][Bottom][Left]]
	Centre   [][]int8
	SMatched []int
	IDMatch  []int // [Top][Right][Bottom][Left]
	Rotation int
	Flip     bool
	XY       []int
}

func main() {
	part := "B"

	defer timeTrack(time.Now(), "day20")
	// input, _ := readLines("tiles.txt")
	input, _ := readLines("tiles-test.txt")

	var side_len int
	var tiles []tile
	var cur_tile tile
	pix_cnt := 0
	pix_cnt_added := 0
	for i := range input {
		_, mod := divMod(i, side_len+2)
		if i == 0 || mod == 0 {
			if i != 0 {
				tiles = append(tiles, cur_tile)
				for i := range cur_tile.Sides {
					pix_cnt_added = pix_cnt_added + len(cur_tile.Sides[i])
				}
				if pix_cnt_added != pix_cnt {
					fmt.Println("\nERROR: pix_cnt_added != pix_cnt", pix_cnt_added, "!=", pix_cnt)
				}
				pix_cnt = 0
				pix_cnt_added = 0
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
			pix = getNonEmptyPix(s, true)
			if mod == 1 {
				cur_tile.Sides = append(cur_tile.Sides, pix)
				cur_tile.Sides = append(cur_tile.Sides, make([]int, 0, side_len))
				cur_tile.Sides = append(cur_tile.Sides, make([]int, 0, side_len))
				cur_tile.Sides = append(cur_tile.Sides, make([]int, 0, side_len))
				cur_tile.IDMatch = make([]int, 4)
			} else if mod == side_len {
				for i := range pix {
					cur_tile.Sides[2] = append(cur_tile.Sides[2], side_len-1-pix[i])
				}
			} else {
			}
			if len(pix) > 0 {
				pix_cnt = pix_cnt + len(pix)
				if pix[0] == 0 {
					cur_tile.Sides[3] = append(cur_tile.Sides[3], side_len-mod)
				}
				if pix[len(pix)-1] == side_len-1 {
					fmt.Println(side_len - mod)
					cur_tile.Sides[1] = append(cur_tile.Sides[1], mod-1)
				}
			}
			if mod > 1 && mod < side_len {
				cur_tile.Centre = append(cur_tile.Centre, make([]int8, 0))
				for idx := 1; idx < side_len-1; idx++ {
					found := false
					for j := range pix {
						if pix[j] == idx {
							found = true
							cur_tile.Centre[len(cur_tile.Centre)-1] = append(cur_tile.Centre[len(cur_tile.Centre)-1], 1)
							pix_cnt_added = pix_cnt_added + 1
						}
					}
					if !found {
						cur_tile.Centre[len(cur_tile.Centre)-1] = append(cur_tile.Centre[len(cur_tile.Centre)-1], 0)
					}
				}
			}
		}
		if i == len(input)-1 {
			tiles = append(tiles, cur_tile)
			pix_cnt_added := (side_len - 2) * 2
			for i := range cur_tile.Sides {
				pix_cnt_added = pix_cnt_added + len(cur_tile.Sides[i])
			}
			if pix_cnt_added != pix_cnt {
				fmt.Println("\nERROR: pix_cnt_added != pix_cnt", pix_cnt_added, "!=", pix_cnt)
			}
			pix_cnt = 0
			pix_cnt_added = 0
		}
	}

	fmt.Println("tile side_len=", side_len)
	for i := range tiles {
		fmt.Println(tiles[i])
	}

	tiles = findAdjacent(tiles, side_len)

	corner_multi := 1
	corner_cnt := 0
	for i := range tiles {
		if len(tiles[i].SMatched) == 2 {
			corner_multi = corner_multi * tiles[i].ID
			corner_cnt = corner_cnt + 1
		}
	}
	if corner_cnt > 4 {
		fmt.Println("ERROR")
		fmt.Println("corner_cnt=", corner_cnt, "> 4")
	}
	fmt.Println("\n------ Answer Part 1 ------")
	fmt.Println("corner_multi=", corner_multi)
	fmt.Println("\n------ Working Part 2 ------")

	if part == "B" {
		tile_cnt := len(tiles)
		image_size_f := math.Sqrt(float64(tile_cnt))
		image_size := int(image_size_f)
		fmt.Println("tile_cnt=", tile_cnt)
		fmt.Println("image_size=", image_size)
		pos_m := alignTiles(tiles, side_len, image_size, tile_cnt)
		for i := range pos_m {
			fmt.Println(pos_m[i])
		}

		// Sea Monster
		var s_m [3][20]int8
		s_m[0] = [...]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
		s_m[1] = [...]int8{1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1}
		s_m[2] = [...]int8{0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0}
		s_m_i := 0
		for i := range s_m {
			for j := range s_m[i] {
				if s_m[i][j] == 1 {
					s_m_i = s_m_i + 1
				}
			}
		}

		var f_map [][]int8

		for i := 0; i < image_size*(side_len-2); i++ {
			f_map = append(f_map, make([]int8, 0))
		}

		for i := range pos_m {
			for j := range pos_m[i] {
				for t := range tiles {
					if tiles[t].ID == pos_m[i][j] {
						for ci := range tiles[t].Centre {
							f_map[ci+i*(side_len-2)] = append(f_map[ci+i*(side_len-2)], tiles[t].Centre[ci]...)
						}
					}
				}
			}
		}
		printMap(f_map)

		// iterate over map with the monster
		flip := false
		rot := 0
		s_m_f := 0
		for s_m_f == 0 {
			i_i := 0
			for i_i <= len(f_map)-len(s_m) {
				j_j := 0
				for j_j <= len(f_map)-len(s_m[0]) {
					cur_s_m := 0
					for i := range s_m {
						for j := range s_m[i] {
							if s_m[i][j] != 1 {
								continue
							}
							if s_m[i][j] == f_map[i_i+i][j_j+j] {
								cur_s_m = cur_s_m + 1
							}
						}
					}
					if cur_s_m == s_m_i {
						fmt.Println("i_i", i_i, "j_j", j_j, s_m_i, "?", cur_s_m)
						s_m_f = s_m_f + 1
						for i := range s_m {
							for j := range s_m[i] {
								if s_m[i][j] != 1 {
									continue
								}
								if s_m[i][j] == f_map[i_i+i][j_j+j] {
									f_map[i_i+i][j_j+j] = 2
								}
							}
						}
					}
					j_j = j_j + 1
				}
				i_i = i_i + 1
			}
			if s_m_f == 0 && rot < 3 {
				_, f_map = rotate([]int{}, f_map, 1)
				rot = rot + 1
				fmt.Println("\nROTATE=", rot, "  Flipped?=", flip)
				printMap(f_map)
			}
			if s_m_f == 0 && rot == 3 && !flip {
				_, f_map = transpose([]int{}, f_map)
				rot = 0
				flip = true
				fmt.Println("\n-FLIP-", "ROTATE=", rot)
				printMap(f_map)
			}
			if rot == 3 && flip {
				break
			}
		}
		printMap(f_map)
		fmt.Println("Sea Monsters Found=", s_m_f)
		rough_water := 0
		for i := range f_map {
			for j := range f_map[i] {
				if f_map[i][j] == 1 {
					rough_water = rough_water + 1
				}
			}
		}
		fmt.Println("\n------ Answer Part 2 ------")
		fmt.Println("rough_water=", rough_water)
	}
}

func printMap(f_map [][]int8) {
	for i := range f_map {
		fmt.Print("\n")
		for j := range f_map[i] {
			if f_map[i][j] == 1 {
				fmt.Print("#")
			} else if f_map[i][j] == 0 {
				fmt.Print(".")
			} else {
				fmt.Print("O")
			}
		}
	}
	fmt.Print("\n")
	fmt.Println(len(f_map), "x", len(f_map[0]))
}

func findAdjacent(tiles []tile, side_len int) []tile {
	for i := range tiles {
		for j := range tiles {
			if i == j {
				continue
			}
			for si := range tiles[i].Sides {
				for sj := range tiles[j].Sides {
					s1 := tiles[i].Sides[si]
					var s2 []int
					for x := range tiles[j].Sides[sj] {
						s2 = append(s2, side_len-1-tiles[j].Sides[sj][x])
					}
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
					} else {
						match1 = false
						match2 = false
					}
					if match1 {
						tiles[i].SMatched = append(tiles[i].SMatched, si)
						tiles[i].IDMatch[si] = tiles[j].ID
						fmt.Println(tiles[i].ID, "match", tiles[j].ID, "on", s1, "==", s2)
					}
					if match2 {
						tiles[i].SMatched = append(tiles[i].SMatched, -si)
						tiles[i].IDMatch[si] = -tiles[j].ID
						fmt.Println(tiles[i].ID, "match", tiles[j].ID, "on", s1, "==", s2)
					}
				}
			}
		}
	}
	return tiles
}

func alignTiles(tiles []tile, side_len int, image_size int, tile_cnt int) [][]int {
	var pos_m [][]int
	for i := 0; i < image_size; i++ {
		pos_m = append(pos_m, make([]int, image_size))
	}
	for i := range pos_m {
		fmt.Println(pos_m[i])
	}

	arranged := 0
	corner := true
	side := false
	cor_xy := []int{0, 0}
	var last_tile *tile
Arrange:
	for arranged < tile_cnt {
		for y := range pos_m {
		Next:
			for x := range pos_m[y] {
				cor_xy[0] = x
				cor_xy[1] = y
				fmt.Println("\n", pos_m, "\n")
				fmt.Println(x)
				fmt.Println(y)
				fmt.Println("corner", corner)
				fmt.Println("side", side)
				if cor_xy[0] == 0 && cor_xy[1] == 0 {

				} else if corner {
					t := getCompareTile(&tiles, pos_m, x-1, 0)
					last_tile = &tiles[t]
				} else if side {
					t := getCompareTile(&tiles, pos_m, x-1, y)
					last_tile = &tiles[t]
				} else {
					t := getCompareTile(&tiles, pos_m, x, y-1)
					last_tile = &tiles[t]
				}
				fmt.Println("\nlast_tile=", last_tile)

				found := false
				if pos_m[y][x] != 0 {
					continue
				}
				if side && y != 0 && (x == 0 || x == image_size-1) {
					continue
				}
				if corner && !((x == 0 && x == y) || (x == image_size-1 && y == 0)) {
					continue
				}

				for i := range tiles {
					if len(tiles[i].XY) != 0 {
						continue
					}
					if corner && len(tiles[i].SMatched) == 2 && len(tiles[i].XY) == 0 {
						min_side := 3
						max_side := 0
						for sm := range tiles[i].SMatched {
							if tiles[i].SMatched[sm] < min_side {
								min_side = tiles[i].SMatched[sm]
							}
							if tiles[i].SMatched[sm] > max_side {
								max_side = tiles[i].SMatched[sm]
							}
						}
						d := Abs(max_side - min_side)

						if cor_xy[0] == 0 && cor_xy[1] == 0 {
							fmt.Println("tiles[i]=", tiles[i])
							if d == 3 {
								tileRotate(&tiles[i], max_side, 1)
							} else {
								tileRotate(&tiles[i], min_side, 1)
							}
							found = true
							side = true
							fmt.Println("tiles[i]=", tiles[i])
						} else {
							if d == 1 {
								found = findTileJoin(last_tile, &tiles[i], min_side)
							} else {
								found = findTileJoin(last_tile, &tiles[i], max_side)
							}
							if found {
								side = false
								corner = false
							}
						}

						if !found {
							continue
						} else {
							tiles[i].XY = cor_xy
							pos_m[y][x] = tiles[i].ID
							corner = false
							arranged = arranged + 1
						}

					} else if side && len(tiles[i].SMatched) == 3 && len(tiles[i].XY) == 0 {
						found = findTileJoin(last_tile, &tiles[i], 1)
						if !found {
							continue
						} else {
							tiles[i].XY = cor_xy
							pos_m[y][x] = tiles[i].ID
							arranged = arranged + 1
							if x == image_size-1 {
								for t := range tiles {
									if tiles[t].ID == pos_m[y][0] {
										last_tile = &tiles[t]
									}
								}
							} else {
								last_tile = &tiles[i]
							}
						}

						if cor_xy[0] == 0 && cor_xy[1] == 1 {
							side = false
						} else if cor_xy[0] == image_size-2 && cor_xy[1] == 0 {
							corner = true
							side = false
						}

					} else if !corner && !side && len(tiles[i].XY) == 0 {
						found = findTileJoin(last_tile, &tiles[i], 2)
						if !found {
							continue
						} else {
							fmt.Println("\n=========findTileJoin=", &tiles[i])
							tiles[i].XY = cor_xy
							pos_m[y][x] = tiles[i].ID
							arranged = arranged + 1
						}
					}

					if arranged == tile_cnt {
						fmt.Println("THE END")
						break Arrange
					}

					if found {
						continue Next
					}
				}
			}
		}
	}
	return pos_m
}

func getCompareTile(tiles *[]tile, m [][]int, x int, y int) int {
	for t, v := range *tiles {
		if v.ID == m[y][x] {
			return t
		}
	}
	return -1
}

// findTileJoin finds the corresponding side join and rotates/transpose to the correct orientation
func findTileJoin(lt *tile, tm *tile, side int) bool {
	fmt.Println("--lt", lt)
	fmt.Println("--side", side)
	fmt.Println("--tm", tm)
	fmt.Println("lt.IDMatch[side]", lt.IDMatch[side])
	fmt.Println("tm.ID", tm.ID)
	if Abs(lt.IDMatch[side]) == tm.ID {
		for idx := range tm.IDMatch {
			_, req_idx := divMod(side+2, 4)
			if tm.IDMatch[idx] == lt.ID {
				tileRotate(tm, idx, req_idx)
				fmt.Println("++++++++++++++++findTileJoin", Abs(lt.IDMatch[side]), "==", tm.ID)
				return true
			}
			if tm.IDMatch[idx] == -lt.ID {
				tileTranspose(tm)
				lt.IDMatch[side] = Abs(lt.IDMatch[side])
				for idx := range tm.IDMatch {
					_, req_idx := divMod(side+2, 4)
					if tm.IDMatch[idx] == lt.ID {
						tileRotate(tm, idx, req_idx)
						fmt.Println("++++++++++++++++findTileJoin", Abs(lt.IDMatch[side]), "==", tm.ID)
						return true
					}
				}
			}
		}
	}
	fmt.Println("!!!!!!!!!!!!!!!!!!!!findTileJoin", Abs(lt.IDMatch[side]), "!=", tm.ID)
	return false
}

// titleRotate rotates the tile to the required orientation
func tileRotate(t *tile, s int, m int) {
	change := 0
	test := s
	for test != m {
		change = change + 1
		_, test = divMod(s+change, 4)

	}

	if change > 0 {
		arr, pix := rotate(t.IDMatch, t.Centre, change)
		copy(t.IDMatch, arr[:])
		copy(t.Centre, pix)
	}
	t.Rotation = change
	return
}

func rotate(t []int, pix [][]int8, r int) ([4]int, [][]int8) {
	var arr [4]int
	for i := range t {
		_, new_index := divMod(i+r, 4)
		arr[new_index] = t[i]
	}

	var t_copy [][]int8
	for i := range pix {
		t_copy = append(t_copy, make([]int8, 0))
		for j := range pix[i] {
			t_copy[i] = append(t_copy[i], pix[i][j])
		}
	}

	if r == 1 {
		for i := range t_copy {
			for j := range t_copy[i] {
				pix[j][len(t_copy)-1-i] = t_copy[i][j]
			}
		}
	} else if r == 2 {
		for i := range t_copy {
			for j := range t_copy[i] {
				pix[len(t_copy)-1-i][len(t_copy)-1-j] = t_copy[i][j]
			}
		}
	} else if r == 3 {
		for i := range t_copy {
			for j := range t_copy[i] {
				pix[len(t_copy)-1-j][i] = t_copy[i][j]
			}
		}
	}

	return arr, pix
}

// tileTranspose transpose the tile
func tileTranspose(t *tile) {
	arr, pix := transpose(t.IDMatch, t.Centre)
	copy(t.IDMatch, arr[:])
	copy(t.Centre, pix)
	return
}
func transpose(t []int, pix [][]int8) ([4]int, [][]int8) {
	var t_copy [][]int8
	var arr [4]int

	for i := range pix {
		t_copy = append(t_copy, make([]int8, 0))
		for j := range pix[i] {
			t_copy[i] = append(t_copy[i], pix[i][j])
		}
	}

	for i := range t_copy {
		for j := range t_copy[i] {
			pix[j][i] = t_copy[i][j]
		}
	}

	for i := range t {
		arr[i] = Abs(t[len(arr)-1-i])
	}

	// t.Flip = true
	return arr, pix
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
