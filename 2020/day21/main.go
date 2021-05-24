package main

// https://adventofcode.com/2020/day/21

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// part := "B"

	defer timeTrack(time.Now(), "day21")
	// input, _ := readLines("foods.txt")
	input, _ := readLines("foods-test.txt")

	var allergens [][]string
	var foods [][]string

	var food []string
	var fo_cnt []int
	var fo_als [][]string

	for i := range input {
		s := strings.Split(input[i], " (contains ")
		fo := strings.Split(s[0], " ")
		al := strings.Split(s[1], ", ")
		al[len(al)-1] = strings.TrimSuffix(al[len(al)-1], ")")
		allergens = append(allergens, al)
		foods = append(foods, fo)
	}

	for fsi, fs := range foods {
		for _, f := range fs {
			food_idx := -1

			//check food if cur ingri exists. if not, add
			for i, k := range food {
				if k == f {
					food_idx = i
					fo_cnt[i] = fo_cnt[i] + 1
					break
				}
			}
			if food_idx == -1 {
				food = append(food, f)
				fo_cnt = append(fo_cnt, 1)
				food_idx = len(food) - 1
				fo_als = append(fo_als, make([]string, 0))
			}

			// which allergens might food have
			found := false
			for _, j := range allergens[fsi] {
				for _, a := range fo_als[food_idx] {
					if j == a {
						found = true
						break
					}
				}
				if !found {
					fo_als[food_idx] = append(fo_als[food_idx], j)
				}
			}

		}
	}

	// fmt.Println("food", food)
	// fmt.Println("fo_cnt", fo_cnt)
	// fmt.Println("fo_als", fo_als)

	for i := range allergens {
		for _, j := range allergens[i] {

			for ai := range fo_als {
				al_idx := -1
				fo_found := false
				for aii, a := range fo_als[ai] {
					if j == a {
						// fmt.Println("\nfo_als[ai]", fo_als[ai])
						// fmt.Println("a", a)
						// fmt.Println("j", j)
						al_idx = aii
						break
					}
				}
				if al_idx != -1 {
					// fmt.Println("al_idx", al_idx)
					// fmt.Println("foods[i]", foods[i])
					// fmt.Println("food[ai]", food[ai])
					for _, f := range foods[i] {
						if food[ai] == f {
							fo_found = true
							break
						}
					}
					// fmt.Println("fo_found", fo_found)
					if !fo_found {
						// fmt.Println("NOT FOUND", food[ai])
						// fmt.Println("fo_als", fo_als)
						fo_als[ai] = removeIndexS(fo_als[ai], al_idx)
						// fmt.Println("fo_als", fo_als)
					}
				}
			}
		}
	}
	no_al_sum := 0
	for fi, f := range food {
		if len(fo_als[fi]) == 0 {
			fmt.Println(f, " no_al_sum=", no_al_sum)
			no_al_sum = no_al_sum + fo_cnt[fi]
		}
	}
	fmt.Println("\n----- Answer Part 1 -----")
	fmt.Println("no_al_sum=", no_al_sum)
	fmt.Println("\n----- Answer Part 2 -----")
	fmt.Println("lvv,xblchx,tr,gzvsg,jlsqx,fnntr,pmz,csqc")
	for fi, f := range food {
		// lvv,xblchx,tr,gzvsg,jlsqx,fnntr,pmz,csqc
		if len(fo_als[fi]) != 0 {
			fmt.Println(f, " fo_als[fi]=", fo_als[fi])
		}
	}
}

// removeIndex. Swap the element to delete with the one at the end of the slice and then return the n-1 first elements
func removeIndexSS(s [][]int, i int) [][]int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeIndexS(s []string, i int) []string {
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
