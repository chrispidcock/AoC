package main

// https://adventofcode.com/2020/day/15

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// total_iteractions := 5
	// total_iteractions := 2020
	total_iteractions := 30000000

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day15")

	// start_nums := []int{0, 3, 6}
	start_nums := []int{19, 0, 5, 1, 10, 13}

	var nums []int
	var age []int
	var prev_num int
	for i := range start_nums {
		if i != len(start_nums)-1 {
			nums, age = updateSlices(start_nums[i], nums, age, -1, i)
		} else {
			prev_num = start_nums[i]
		}

	}

	iteration := len(nums) - 1

	for iteration < total_iteractions-2 {
		iteration = iteration + 1
		_, mod := divMod(iteration, 10000)
		if mod == 0 {
			logging("iteration", iteration, true)
		}
		num_index := findNum(nums, prev_num)

		if num_index == -1 {
			nums, age = updateSlices(prev_num, nums, age, num_index, iteration)
			prev_num = 0
		} else {
			new_num := iteration - age[num_index]
			nums, age = updateSlices(prev_num, nums, age, num_index, iteration)
			prev_num = new_num
		}
	}
	logging("---Answer--", "", true)
	logging("prev_num", prev_num, true)

}

func updateSlices(new int, nums []int, age []int, index int, iteration int) ([]int, []int) {
	var num_index int
	if index == -1 {
		num_index = index
	} else {
		num_index = findNum(nums, new)
	}
	logging("new", new, false)
	logging("num_index", num_index, false)

	if num_index == -1 {
		nums = append(nums, new)
		num_index = len(nums) - 1
		age = append(age, iteration)
	} else {
		age[num_index] = iteration
	}
	logging("nums", nums, false)
	logging("age", age, false)
	return nums, age
}

func findNum(input []int, num int) int {
	for i := range input {
		if input[i] == num {
			return i
		}
	}
	return -1
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
