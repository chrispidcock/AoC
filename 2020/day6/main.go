package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	defer timeTrack(time.Now(), "day6")
	input, _ := readLines("answers.txt")

	var groupSize []int
	var size int
	var yesSum int
	var allSum int

	var textBlocks []string
	var textBlock string
	// fmt.Println(len(input))
	for _, text := range input {
		if text == "" {
			// fmt.Println(textBlock)
			textBlocks = append(textBlocks, textBlock)
			textBlock = ""
			groupSize = append(groupSize, size)
			size = 0
		} else {
			size = size + 1
			textBlock = textBlock + text
		}
	}

	for i, group := range textBlocks {
		var groupAnswers = make(map[rune]int)
		for _, char := range group {
			if val, ok := groupAnswers[char]; ok {
				groupAnswers[char] = 1 + val
			} else {
				groupAnswers[char] = 1
			}
		}
		yesSum = yesSum + len(groupAnswers)
		for _, yes := range groupAnswers {
			if yes == groupSize[i] {
				allSum = allSum + 1
			}
		}
	}

	fmt.Println("partA: answer sum:", yesSum)
	fmt.Println("partB: whole group answer sum:", allSum)
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
