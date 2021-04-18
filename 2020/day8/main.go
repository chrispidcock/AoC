package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	"strconv"
)

func main() {
	defer timeTrack(time.Now(), "day8")
	input, _ := readLines("cypher.txt")

	fmt.Println(len(input))

	var cur_num int
	var val_sum int
	var valid bool

	for i := 25; i < len(input); i++ {
		var window = input[i-25:i]
		fmt.Println(window)
		fmt.Println(len(window))

		cur_num = input[i]
		fmt.Println("cur_num: ", cur_num)
		valid = false
		Find:
			for w1, v := range window {
				for w2 := w1+1; w2 <= len(window)-1; w2++ {
					val_sum = v + window[w2]
					fmt.Println(v, " + ", window[w2])
					fmt.Println(val_sum)
					if cur_num == val_sum {
						valid = true
						fmt.Println("Match")
						break Find
					}
				}
			}
		if valid == false {
			fmt.Println("Invalid Number: ", cur_num)

			for back := i-2; back >= 0; back-- {
				curSlice := input[back:i-1]
				fmt.Println(curSlice)
				sumOut := sliceSum(curSlice)
				fmt.Println("Invalid Number: ", cur_num)
				fmt.Println("Current Sum: ", sumOut)
				if sumOut == cur_num {
					min, max := minMax(curSlice)
					fmt.Println("Min: ", min)
					fmt.Println("Max: ", max)
					break
				}
				if sumOut > cur_num {
					break
				}
			}
			break
		}
	}

}

// sliceSum Sums slice of varying lengths
func sliceSum(intSlice []int) (int) {
	var total int = 0
	for _, v := range intSlice {
		total = total + v
	}
	return total
}


func minMax(values []int)(int, int) {
	min := values[0]
	max := values[0]
	for _, v := range values {
			if (v < min) {
				min = v
			}
			if (v > max) {
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
