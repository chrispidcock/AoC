package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	defer timeTrack(time.Now(), "day3")

	var day1A bool = false

	rows, _ := readLines("tree_map.txt")
	var rowLen int = len(rows[0])

	if day1A == true {
		trees := descend(1, 3, rowLen, rows)
		fmt.Println(trees)
	} else {
		var treesMultiple int = 1
		downs := [...]int{1, 1, 1, 1, 2}
		rights := [...]int{1, 3, 5, 7, 1}

		for i := 0; i < len(downs); i++ {
			treesMultiple = treesMultiple * descend(downs[i], rights[i], rowLen, rows)
		}
		fmt.Println(treesMultiple)
	}
}

func descend(down int, right int, rowLen int, rows []string) int {
	var trees int = 0
	var x int = 0
	var y int = 0
	for y < len(rows)-1 {
		y = y + down
		x = (x + right) % (rowLen)
		if string(rows[y][x]) == "#" {
			trees = trees + 1
		}
	}
	return trees
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
