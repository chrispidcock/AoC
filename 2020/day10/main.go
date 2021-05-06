package main

// https://adventofcode.com/2020/day/10

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	defer timeTrack(time.Now(), "day10")
	input, _ := readLines("adapters.txt")
	sort.Ints(input[:])
	// outlet := 0
	// maxJoltDiff := 3
	dif1 := 0
	dif2 := 0
	dif3 := 1
	var difArray []int
	var diff int
	fmt.Println("input: ", input)

	for i := range input {
		if i == 0 {
			diff = input[i]
		} else {
			diff = input[i] - input[i-1]
			fmt.Println("input[i-1]: ", input[i-1])
			fmt.Println("input[i]: ", input[i])
			fmt.Println("diff: ", diff)
		}

		switch diff {
		case 1:
			dif1 = dif1 + 1
		case 2:
			dif2 = dif2 + 1
		case 3:
			dif3 = dif3 + 1
		default:
			fmt.Println("diff miseed switch:", diff)
		}
		difArray = append(difArray, diff)
	}
	difArray = append(difArray, 3)
	fmt.Println("difArray: ", difArray)
	fmt.Println("dif1: ", dif1)
	fmt.Println("dif2: ", dif2)
	fmt.Println("dif3: ", dif3)
	fmt.Println("dif1 + dif2 + dif3: ", dif1+dif2+dif3)
	fmt.Println("len(input) :", len(input))
	fmt.Println("dif1 * dif3: ", dif1*dif3)

	permSum := difPermutations(difArray)
	fmt.Println(permSum)
}

func difPermutations(diffArray []int) int {
	// Example1, where the int indicates the jolts from the previous socket
	// 1 1 1 1 - 4
	// 1 1 0 2 - 3
	// 1 0 2 1 - 3
	// 0 2 1 1 - 3
	// 1 0 0 3 - 2
	// 0 0 3 1 - 2
	// 0 2 0 2 - 2
	// n = 7

	// Example2
	// 1 1 1 - 3
	// 1 0 2 - 2
	// 0 2 1 - 2
	// 0 0 3 - 1
	// n = 4

	// Example3
	// 1 1 1 1 1 - 5
	// 1 1 0 2 1 - 4
	// 1 0 2 1 1 - 4
	// 0 2 1 1 1 - 4
	// 1 1 1 0 2 - 4
	// 1 1 0 0 3 - 3
	// 1 0 0 3 1 - 3
	// 1 0 2 0 2 - 3
	// 0 2 1 0 2 - 3
	// 0 2 0 2 1 - 3
	// 0 0 3 1 1 - 3
	// 0 2 0 0 3 - 2
	// 0 0 3 0 2 - 2
	// n = 13
	// jolt1 = [...]int{1}
	// jolt2 = [...]int{0, 2}
	// jolt3 = [...]int{0, 0, 3}
	diffs := []int{2, 3}
	startIndex := 0
	permSum := 0
	var permSlice []int
	var endIndex int

	for i1, v1 := range diffArray {
		// Start checking for permutations once we find the end of a block of 1s
		if v1 == 3 {
			permInt := 0
			endIndex = i1
			curSlice := diffArray[startIndex:endIndex]
			fmt.Println("curSlice ", curSlice)
			arrayLen := len(curSlice)
			fmt.Println("arrayLen ", arrayLen)

			// Find how many times we can fit 2 and 3 jolt adapters in the length of 1 jolt space
			// Generate array of 1s [1 1 1 1 1 1 1]
			// Permutations of all 1 jolts = 1
			if arrayLen > 1 {
				permInt = permInt + 1
			}
			if arrayLen < 3 {
				permInt = permInt + 1
				permSlice = append(permSlice, permInt)
				startIndex = endIndex + 1
				fmt.Println("permInt ", permInt)
				continue
			}

			for _, v2 := range diffs {
				fmt.Println("v2 ", v2)
				whole, mod := divMod(arrayLen, v2)
				fmt.Println("whole ", whole)
				fmt.Println("mod ", mod)
				permutations := make([][]int, 0)
				for i := 1; i <= whole; i++ {
					adapters := arrayLen - (i * (v2 - 1))
					fmt.Println("adapters ", adapters)
					for ii := 0; ii < adapters; ii++ {
						workingSlice := IntSlice(adapters, 1)
						workingSlice[ii] = v2

						for i4 := range workingSlice {
							tempSlice := make([]int, adapters)
							copy(tempSlice, workingSlice)
							tempSlice[i4] = v2
							if sliceSum(tempSlice) != arrayLen {
								continue
							}
							if newPerm(permutations, tempSlice) {
								permutations = append(permutations, tempSlice)
								fmt.Println("permutations ", permutations)
								permInt = permInt + 1
								fmt.Println("permInt ", permInt)
							}
						}

					}
				}
			}
			fmt.Println(permInt)
			permSlice = append(permSlice, permInt)
			startIndex = endIndex + 1
		}
	}
	fmt.Println(permSlice)
	permSum = sliceMultiply(permSlice)
	return permSum
}

func valueFit(indexValue, value int) bool {
	if indexValue != value {
		return true
	}
	return false
}

func newPerm(permutations [][]int, sliceCheck []int) bool {
	for i := 0; i < len(permutations); i++ {
		if Equal(permutations[i], sliceCheck) {
			return false
		}
	}
	return true
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func sliceMultiply(intSlice []int) int {
	result := 1
	for _, v := range intSlice {
		if v == 0 {
			continue
		}
		result = result * v
	}
	return result
}

func sliceSum(intSlice []int) int {
	result := 0
	for _, v := range intSlice {
		result = result + v
	}
	return result
}

// OnceSlice creates a slice of 1s
func IntSlice(length int, value int) []int {
	s := make([]int, length)
	for i := range s {
		s[i] = value
	}
	return s
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

func readLines(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		value, err := strconv.Atoi(txt)
		if err != nil {
			fmt.Println("Error reading file")
			break
		}
		lines = append(lines, value)
	}
	return lines, scanner.Err()
}
