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
	Sides    [][]uint8 // [[Top][Right][Bottom][Left]]
	Centre   [][]uint8
	SMatched []int
	IDMatch  []int // [Top Right Bottom Left]
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
	for i := range input {
		_, mod := divMod(i, side_len+2)
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
			var pix []uint8
			s := strings.Split(input[i], "")
			side_len = len(s)
			pix = getPixVal(s)
			if mod == 1 {
				cur_tile.Sides = append(cur_tile.Sides, pix)
				cur_tile.Sides = append(cur_tile.Sides, make([]uint8, side_len, side_len))
				cur_tile.Sides = append(cur_tile.Sides, make([]uint8, side_len, side_len))
				cur_tile.Sides = append(cur_tile.Sides, make([]uint8, side_len, side_len))
				cur_tile.IDMatch = make([]int, 4)
			} else if mod == side_len {
				for x := range pix {
					cur_tile.Sides[2][len(pix)-1-x] = pix[x]
				}
			}
			if len(pix) > 0 {
				cur_tile.Sides[3][side_len-mod] = pix[0]
				cur_tile.Sides[1][mod-1] = pix[len(pix)-1]
			}
			if mod > 1 && mod < side_len {
				cur_tile.Centre = append(cur_tile.Centre, pix[1:side_len-1])
			}
		}
		if i == len(input)-1 {
			tiles = append(tiles, cur_tile)
		}
	}

	fmt.Println("tile side_len=", side_len)
	for i := range tiles {
		fmt.Println("\n", tiles[i].ID)
		for s := range tiles[i].Sides {
			fmt.Println(tiles[i].Sides[s])
		}
	}
	tiles = findAdjacent(tiles, side_len)
	for i := range tiles {
		fmt.Println(tiles[i].IDMatch)
	}

	tile_cnt := len(tiles)
	image_size_f := math.Sqrt(float64(tile_cnt))
	image_size := int(image_size_f)

	len_2_t := 0
	len_3_t := 0
	len_4_t := 0
	fmt.Println("tile side_len=", side_len)
	for i := range tiles {
		matches := 0
		for j := range tiles[i].IDMatch {
			if tiles[i].IDMatch[j] != 0 {
				matches = matches + 1
			}
		}
		if matches == 2 {
			len_2_t = len_2_t + 1
		} else if matches == 3 {
			len_3_t = len_3_t + 1
		} else if matches == 4 {
			len_4_t = len_4_t + 1
		} else {
			fmt.Println("INCORRECT MATCHES", tiles[i])
			panic("INCORRECT MATCHES")
		}
	}
	fmt.Println("** len_2_t=", len_2_t)
	fmt.Println("** len_3_t=", len_3_t)
	fmt.Println("** len_4_t=", len_4_t)
	if len_2_t != 4 {
		fmt.Println("** len_2_t=", len_2_t)
		panic("INCORRECT MATCH COUNT")
	} else if len_3_t != 4*(image_size-2) {
		fmt.Println("** len_3_t=", len_3_t)
		panic("INCORRECT MATCH COUNT")
	} else if len_4_t != (image_size-2)*(image_size-2) {
		fmt.Println("** len_4_t=", len_4_t)
		panic("INCORRECT MATCH COUNT")
	}

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
		fmt.Println("tile_cnt=", tile_cnt)
		fmt.Println("image_size=", image_size)
		pos_m := alignTiles(tiles, side_len, image_size, tile_cnt)
		for i := range pos_m {
			fmt.Println(pos_m[i])
		}
		valid := validateJoin(tiles, pos_m)
		if valid != 0 {
			// return
		}

		// Sea Monster
		var s_m [3][20]uint8
		s_m[0] = [...]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
		s_m[1] = [...]uint8{1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1}
		s_m[2] = [...]uint8{0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0}
		s_m_i := 0
		for i := range s_m {
			for j := range s_m[i] {
				if s_m[i][j] == 1 {
					s_m_i = s_m_i + 1
				}
			}
		}

		var f_map [][]uint8

		for i := 0; i < image_size*(side_len-2); i++ {
			f_map = append(f_map, make([]uint8, 0))
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

func printMap(f_map [][]uint8) {
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
					s11 := tiles[i].Sides[si]
					var s12 []uint8
					for x := len(tiles[i].Sides[si]) - 1; x >= 0; x-- {
						s12 = append(s12, tiles[i].Sides[si][x])
					}
					s2 := tiles[j].Sides[sj]

					if len(s11) != len(s2) {
						continue
					}

					match1 := true
					match2 := true
					if len(s11) > 0 && 0 < len(s2) {
						for ii := range s11 {
							if s11[ii] != s2[ii] {
								match1 = false
								break
							}
						}
						for ii := range s12 {
							if s12[ii] != s2[ii] {
								match2 = false
								break
							}
						}
					} else {
						match1 = false
						match2 = false
					}
					if match1 {
						tiles[i].SMatched = append(tiles[i].SMatched, -si)
						tiles[i].IDMatch[si] = -tiles[j].ID
					}
					if match2 {
						tiles[i].SMatched = append(tiles[i].SMatched, si)
						tiles[i].IDMatch[si] = tiles[j].ID
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
				var t int
				var b []int
				if y == 0 {
					b = append(b, 0)
				}
				if y == len(pos_m)-1 {
					b = append(b, 2)
				}
				if x == 0 {
					b = append(b, 3)
				}
				if x == len(pos_m)-1 {
					b = append(b, 1)
				}
				cor_xy[0] = x
				cor_xy[1] = y
				for i := range pos_m {
					fmt.Println(pos_m[i])
				}
				fmt.Println("(x y)", x, y)
				if cor_xy[0] == 0 && cor_xy[1] == 0 {
				} else if corner {
					t = getCompareTile(&tiles, pos_m, x-1, 0)
					last_tile = &tiles[t]
				} else if side {
					t = getCompareTile(&tiles, pos_m, x-1, y)
					last_tile = &tiles[t]
				} else {
					t = getCompareTile(&tiles, pos_m, x, y-1)
					last_tile = &tiles[t]
				}
				fmt.Println("\nlast_tile=", tiles[t].ID, ":", tiles[t].IDMatch)

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

						if cor_xy[0] == 0 && cor_xy[1] == 0 {
							max_match := 0
							for z := range tiles[i].IDMatch {
								if max_match < tiles[i].IDMatch[z] {
									max_match = tiles[i].IDMatch[z]
								}
							}
							if max_match == 0 {
								tileTranspose(&tiles[i])
							}
							min_side := 3
							max_side := 0
							for sm := range tiles[i].IDMatch {
								if tiles[i].IDMatch[sm] == 0 {
									continue
								}
								if sm < min_side {
									min_side = sm
								}
								if sm > max_side {
									max_side = sm
								}
							}
							d := Abs(max_side - min_side)
							fmt.Println("tiles[i]=", tiles[i].ID, ":", tiles[i].IDMatch)
							if d == 3 {
								tileRotate(&tiles[i], max_side, 1)
							} else {
								tileRotate(&tiles[i], min_side, 1)
							}

							found = true
							side = true
							fmt.Println("tiles[i]=", tiles[i].ID, ":", tiles[i].IDMatch)
						} else {
							found = findTileJoin(last_tile, &tiles[i], 1, b, pos_m, y, x)
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
						found = findTileJoin(last_tile, &tiles[i], 1, b, pos_m, y, x)
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
						found = findTileJoin(last_tile, &tiles[i], 2, b, pos_m, y, x)
						if !found {
							continue
						} else {
							fmt.Println("\n=========findTileJoin=", tiles[i].ID, ":", tiles[i].IDMatch)
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

func validateJoin(ts []tile, m [][]int) int {
	inv_cnt := 0
	for i := range m {
		for j := range m[i] {
			for t := range ts {
				inv_cnt = validatePos(ts[t], m, i, j, inv_cnt)
			}
		}
	}
	return inv_cnt
}

func validatePos(ts tile, m [][]int, i int, j int, inv_cnt int) int {
	if ts.ID == m[i][j] {
		if i == 0 && ts.IDMatch[0] != 0 {
			fmt.Println("i == 0 && ts.IDMatch[0] != 0", ts.IDMatch[0], "!=", 0, "i=", i, "j=", j)
			inv_cnt = inv_cnt + 1
		}
		if i == len(m)-1 && ts.IDMatch[2] != 0 {
			fmt.Println("i == len(m)-1 && ts.IDMatch[2] != 0", ts.IDMatch[1], "!=", 0, "i=", i, "j=", j)
			inv_cnt = inv_cnt + 1
		}
		if j == 0 && ts.IDMatch[3] != 0 {
			fmt.Println("j == 0 && ts.IDMatch[3] != 0", ts.IDMatch[3], "!=", 0, "i=", i, "j=", j)
			inv_cnt = inv_cnt + 1
		}
		if j == len(m)-1 && ts.IDMatch[1] != 0 {
			fmt.Println("j == len(m)-1 && ts.IDMatch[1] != 0", ts.IDMatch[2], "!=", 0, "i=", i, "j=", j)
			inv_cnt = inv_cnt + 1
		}
		for x := range ts.IDMatch {
			if ts.IDMatch[x] == 0 {
				continue
			} else if x == 0 && i > 0 {
				if Abs(ts.IDMatch[x]) != Abs(m[i-1][j]) {
					fmt.Println("ts.IDMatch[x] != m[i-1][j]", ts.IDMatch[x], "!=", m[i-1][j])
					inv_cnt = inv_cnt + 1
				}
			} else if x == 1 && j < len(m)-1 {
				if Abs(ts.IDMatch[x]) != Abs(m[i][j+1]) {
					fmt.Println("ts.IDMatch[x] != m[i][j+1]", ts.IDMatch[x], "!=", m[i][j+1])
					inv_cnt = inv_cnt + 1
				}
			} else if x == 2 && i < len(m)-1 {
				if Abs(ts.IDMatch[x]) != Abs(m[i+1][j]) {
					fmt.Println("ts.IDMatch[x] != m[i+1][j]", ts.IDMatch[x], "!=", m[i+1][j])
					inv_cnt = inv_cnt + 1
				}
			} else if x == 3 && j > 0 {
				if Abs(ts.IDMatch[x]) != Abs(m[i][j-1]) {
					fmt.Println("ts.IDMatch[x] != m[i][j-1]", ts.IDMatch[x], "!=", m[i][j-1])
					inv_cnt = inv_cnt + 1
				}
			}
		}
	}
	return inv_cnt
}

// findTileJoin finds the corresponding side join and rotates/transpose to the correct orientation
func findTileJoin(lt *tile, tm *tile, m_side int, o_side []int, pos_m [][]int, y int, x int) bool {
	if Abs(lt.IDMatch[m_side]) == 0 {
		fmt.Println("\nERROR: findTileJoin\n")
		panic("\nERROR: findTileJoin\n")
	}
	if Abs(lt.IDMatch[m_side]) == Abs(tm.ID) {
		for idx := range tm.IDMatch {
			valid := true
			flip := false
			_, req_idx := divMod(m_side+2, 4)
			if lt.IDMatch[m_side] > 0 && tm.IDMatch[idx] < 0 {
				tm.IDMatch[idx] = Abs(tm.IDMatch[idx])
			}
			if Abs(tm.IDMatch[idx]) == lt.ID {
				tileRotate(tm, idx, req_idx)
				for _, v := range o_side {
					if tm.IDMatch[v] != 0 {
						fmt.Println("VALIDATE FAIL", tm.ID, ":", o_side, tm.IDMatch)
						valid = false
						flip = true
					}
				}
				valid_cnt := validatePos(*tm, pos_m, y, x, 0)
				if valid_cnt != 0 {
					valid = false
				}
				if valid {
					fmt.Println("++++++++++++++++findTileJoin", Abs(lt.IDMatch[m_side]), "==", tm.ID)
					return true
				}

			}
			if (tm.IDMatch[idx] == -lt.ID && valid) || flip {
				valid = true
				tileTranspose(tm)
				for idx := range tm.IDMatch {
					if Abs(tm.IDMatch[idx]) == lt.ID {
						tileRotate(tm, idx, req_idx)
						for _, v := range o_side {
							if tm.IDMatch[v] != 0 {
								fmt.Println("VALIDATE FAIL", tm.ID, ":", o_side, tm.IDMatch)
								valid = false
							}
						}
						if valid {
							fmt.Println("^^^^^^^^^^^^^^^^^^^^^findTileJoin", Abs(lt.IDMatch[m_side]), "==", tm.ID, " >>", tm)
							return true
						}
					}
				}
			}
		}
	}
	fmt.Println("!!!!!!!!!!!!!!!!!!!!findTileJoin", lt.IDMatch[m_side], "!=", tm.ID, " >>", tm)
	return false
}

// titleRotate rotates the tile to the required orientation
func tileRotate(t *tile, s int, m int) {
	// var arr [4]int
	// copy(arr[:], t.IDMatch)
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
		t.Rotation = change
	}
	return
}

func rotate(t []int, pix [][]uint8, change int) ([4]int, [][]uint8) {
	var arr [4]int
	fmt.Println("rotation:", change)
	for i := range t {
		_, new_index := divMod(i+change, 4)
		arr[new_index] = t[i]
	}

	var t_copy [][]uint8
	for i := range pix {
		t_copy = append(t_copy, make([]uint8, 0))
		for j := range pix[i] {
			t_copy[i] = append(t_copy[i], pix[i][j])
		}
	}

	if change == 1 {
		for i := range t_copy {
			for j := range t_copy[i] {
				pix[j][len(t_copy)-1-i] = t_copy[i][j]
			}
		}
	} else if change == 2 {
		for i := range t_copy {
			for j := range t_copy[i] {
				pix[len(t_copy)-1-i][len(t_copy)-1-j] = t_copy[i][j]
			}
		}
	} else if change == 3 {
		for i := range t_copy {
			for j := range t_copy[i] {
				pix[len(t_copy)-1-j][i] = t_copy[i][j]
			}
		}
	} else {
		fmt.Println("\nERROR\n")
	}
	return arr, pix
}

// tileTranspose transpose the tile
func tileTranspose(t *tile) {
	fmt.Println("tileTranspose--t", t.ID, ":", t.IDMatch)
	arr, pix := transpose(t.IDMatch, t.Centre)
	copy(t.IDMatch, arr[:])
	copy(t.Centre, pix)
	t.Flip = true

	fmt.Println("tileTranspose--t", t.ID, ":", t.IDMatch)
	return
}
func transpose(t []int, pix [][]uint8) ([4]int, [][]uint8) {
	var t_copy [][]uint8
	var arr [4]int

	for i := range pix {
		t_copy = append(t_copy, make([]uint8, 0))
		for j := range pix[i] {
			t_copy[i] = append(t_copy[i], pix[i][j])
		}
	}

	for i := range t_copy {
		for j := range t_copy[i] {
			pix[j][i] = t_copy[i][j]
		}
	}

	if len(t) > 0 {
		arr[0] = -t[3]
		arr[1] = -t[2]
		arr[2] = -t[1]
		arr[3] = -t[0]
	}
	return arr, pix
}

func getPixVal(s []string) []uint8 {
	var results []uint8
	for i := range s {
		if s[i] == "#" {
			results = append(results, 1)
		} else if s[i] == "." {
			results = append(results, 0)
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
