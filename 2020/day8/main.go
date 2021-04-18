package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type operation interface {
	operate()
	prepare()
}

type actions struct {
	action    []int
	inc       []int
	seen      []bool
	index     int
	indexSwap []bool
	acc       int
	swaped    bool
}

func main() {
	defer timeTrack(time.Now(), "day8")
	input, _ := readLines("instructions.txt")
	act := &actions{}
	act.prepare(input)

	// day8A
	// var end bool
	// for end == false {
	// 	end = act.operate()
	// }

	// day8B
	var swapCount int
	var end bool
	for act.index < len(act.action) {
		end = act.operate()
		if end == true {
			swapCount = swapCount + 1
			fmt.Println("swap attempt:", swapCount)
			act.swaped = false
			act.seen = nil
			act.seen = make([]bool, len(act.action), len(act.action))
			act.index = 0
			act.acc = 0
		}
	}
	fmt.Println("end:", act.acc)
}

func (a *actions) operate() bool {
	if a.seen[a.index] == true {
		return true
	}
	a.seen[a.index] = true

	switch a.action[a.index] {
	case 0:
		if a.indexSwap[a.index] == false && a.swaped == false {
			a.swaped = true
			a.indexSwap[a.index] = true
			a.index = a.index + a.inc[a.index]
		} else {
			a.index = a.index + 1
		}
	case 1:
		a.acc = a.acc + a.inc[a.index]
		a.index = a.index + 1
	case 2:
		if a.indexSwap[a.index] == false && a.swaped == false {
			a.swaped = true
			a.indexSwap[a.index] = true
			a.index = a.index + 1
		} else {
			a.index = a.index + a.inc[a.index]
		}
	}
	return false
}

func (a *actions) prepare(input []string) {
	a.indexSwap = make([]bool, len(input), len(input))
	a.seen = make([]bool, len(input), len(input))

	for _, row := range input {
		switch row[:3] {
		case "nop":
			a.action = append(a.action, 0)
		case "acc":
			a.action = append(a.action, 1)
		case "jmp":
			a.action = append(a.action, 2)
		}

		i, err := strconv.Atoi(row[4:])
		if err != nil {
			fmt.Println(err)
		}
		a.inc = append(a.inc, i)
	}
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
