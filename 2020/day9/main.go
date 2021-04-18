package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	defer timeTrack(time.Now(), "day9")
	input, _ := readLines("cypher.txt")

	fmt.Println(len(input))

	var cur_num int
	var valid bool
	var sliceSum int

End:
	for i := 25; i < len(input); i++ {
		var window = input[i-25 : i]
		cur_num = input[i]

		valid = false
		valid = findSum(window, cur_num)

		if !valid {
			fmt.Println("Invalid Number: ", cur_num)

			sliceSum = sumOverRange(input[:i-1], cur_num)
			if sliceSum == 0 {
				sliceSum = sumOverRange(input[i+1:], cur_num)
			}
			if sliceSum == cur_num {
				break End
			}
			if sliceSum > cur_num {
				continue End
			}
		}
	}

}

func sumOverRange(intSlice []int, checkVal int) int {
	for i := range intSlice {
		for j := i + 1; j <= len(intSlice)+1; j++ {
			summed := sliceSum(intSlice[i:j])

			fmt.Println("Invalid Number: ", checkVal)
			fmt.Println("Current Sum: ", summed)
			if summed == checkVal {
				min, max := minMax(intSlice[i:j])
				fmt.Println("Min: ", min)
				fmt.Println("Max: ", max)
				fmt.Println("SUM(min,max): ", max+min)
				return summed
			}
			if summed > checkVal {
				break
			}
		}
	}
	return 0
}

func findSum(intSlice []int, cur_num int) bool {
	var val_sum int
	for w1, v := range intSlice {
		for w2 := w1 + 1; w2 <= len(intSlice)-1; w2++ {
			val_sum = v + intSlice[w2]
			if cur_num == val_sum {
				return true
			}
		}
	}
	return false
}

// sliceSum Sums slice of varying lengths
func sliceSum(intSlice []int) int {
	var total int = 0
	for _, v := range intSlice {
		total = total + v
	}
	return total
}

func minMax(values []int) (int, int) {
	min := values[0]
	max := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
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
