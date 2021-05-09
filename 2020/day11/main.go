package main

// https://adventofcode.com/2020/day/11

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// PartA inputs
	// max_dist := 1
	// adj_count_max := 4

	// PartB inputs
	max_dist := 100
	adj_count_max := 5

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day11")
	input, _ := readLines("seat-layout.txt")
	// input, _ := readLines("seat-test-part-b.txt")

	seating := make([]string, 0)
	for _, row := range input {
		r := strings.Split(row, "")
		seating = append(seating, r...)
	}

	// seatingModel := [...]string{seating}
	var seatingModel [8736]string
	var seatingNew [8736]string
	copy(seatingModel[:], seating)
	copy(seatingNew[:], seating)

	rowCount := len(input)
	rowLength := len(seating) / rowCount
	fmt.Println("rowCount ", rowCount)
	fmt.Println("rowLength ", rowLength)
	fmt.Println("spaces ", len(seating))

	// 1 2 3
	// 4 L 5
	// 6 7 8
	adjacent := [][]int{{-1, -1, 1}, {0, -1, 2}, {1, -1, 3}, {-1, 0, 4}, {1, 0, 5}, {-1, 1, 6}, {0, 1, 7}, {1, 1, 8}}

	fmt.Println("max_dist ", max_dist)
	fmt.Println("adj_count_max ", adj_count_max)

	new_occupied := 1
	new_empty := 1
	occupied := 0

	interations := 0

	for new_occupied > 0 || new_empty > 0 {
		interations = interations + 1
		new_occupied = 0
		new_empty = 0
		occupied = 0
		copy(seatingModel[:], seatingNew[:])
		working_slice := arrayToSliceOfSlices(seatingModel, rowLength, rowCount)
		if len(working_slice) > rowCount {
			panic(len(working_slice))
		}
		logging("working_slice", working_slice, false)

		for row := range working_slice {
			logging("row", row, false)
			for column := range working_slice[row] {
				logging("--column--", column, false)
				if working_slice[row][column] == "." {
					continue
				}

				adj_count := 0
				for aj := range adjacent {
					x_axis := column
					y_axis := row
					adj_dist := 0
					adj_found := false
					for x_axis >= 0 && x_axis < rowLength && y_axis >= 0 && y_axis < rowCount && !adj_found && adj_dist < max_dist {
						adj_dist = adj_dist + 1
						x_axis = x_axis + adjacent[aj][0]
						y_axis = y_axis + adjacent[aj][1]
						logging("direction", adjacent[aj][2], false)
						logging("x_axis", x_axis, false)
						logging("y_axis", y_axis, false)
						if x_axis < 0 || x_axis >= rowLength {
							break
						}
						if y_axis < 0 || y_axis >= rowCount {
							break
						}
						if working_slice[y_axis][x_axis] == "L" {
							adj_found = true
						}
						if working_slice[y_axis][x_axis] == "#" {
							adj_count = adj_count + 1
							adj_found = true
							logging("adj_found", adj_found, false)
						}
					}
				}
				logging("adj_count", adj_count, false)

				array_index := row*rowLength + column
				if seatingModel[array_index] == "L" && adj_count == 0 {
					seatingNew[array_index] = "#"
					new_occupied = new_occupied + 1
					// fmt.Println("new_occupied!")

				} else if seatingModel[array_index] == "#" && adj_count >= adj_count_max {
					seatingNew[array_index] = "L"
					new_empty = new_empty + 1
					// fmt.Println("new_empty!")
				}
			}
		}
		for i := range seatingNew {
			if seatingNew[i] == "#" {
				occupied = occupied + 1
			}
		}
		// prettyPrint(seatingNew, rowLength, rowCount)
		// logging("new_occupied", new_occupied, true)
		// logging("new_empty", new_empty, true)
		// logging("occupied", occupied, true)
	}
	prettyPrint(seatingNew, rowLength, rowCount)
	logging("occupied", occupied, true)
	fmt.Println("interations ", interations-1)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func divMod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func prettyPrint(arr [8736]string, columns int, rows int) {
	var printstring []string
	for i := range arr {
		whole, mod := divMod(i, columns)
		if whole == rows {
			break
		}
		printstring = append(printstring, arr[i])
		if mod == columns-1 {
			fmt.Println(printstring)
			printstring = make([]string, 0)
		}
	}
	fmt.Println("")
}

func arrayToSliceOfSlices(arr [8736]string, columns int, rows int) [][]string {
	var slices [][]string
	for i := range arr {
		whole, mod := divMod(i, columns)
		if whole == rows {
			break
		}
		if mod == 0 {
			slices = append(slices, make([]string, 0))
		}
		slices[whole] = append(slices[whole], arr[i])
	}
	return slices
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
