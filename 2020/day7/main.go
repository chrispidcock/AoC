package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Bags interface {
	children()
}
type Bag struct {
	bagName     string
	childNames  []string
	childCounts []int
}

func main() {
	defer timeTrack(time.Now(), "day7")
	input, _ := readLines("bags.txt")

	bagAllMap := make(map[string]map[string]int)
	for _, row := range input {
		s := strings.Fields(row)
		var parentBag string
		var parentFound bool
		bagMap := make(map[string]int)
		var childBagCount int
		var childBag string

		for _, v := range s {
			if v == "contain" {
				continue
			}

			if parentFound == false {
				if strings.Contains(v, "bag") {
					parentFound = true
					continue
				}
				parentBag = parentBag + v
				continue
			}

			if strings.Contains(v, "bag") {
				bagMap[childBag] = childBagCount
				childBag = ""
				continue
			}

			i, err := strconv.Atoi(v)
			if err == nil {
				childBagCount = i
				continue
			}
			childBag = childBag + v
		}
		bagAllMap[parentBag] = bagMap
	}

	var bagsWithShinyGold int
	for k := range bagAllMap {
		bag := &Bag{}
		bag.bagName = k
		bag.children(k, 1, bagAllMap)
		_, found := find(&bag.childNames, "shinygold")
		if found == true {
			bagsWithShinyGold = bagsWithShinyGold + 1
		}
		if k == "shinygold" {
			var bagsInShinyGold int
			for _, v := range bag.childCounts {
				bagsInShinyGold = bagsInShinyGold + v
			}
			fmt.Println("bagsInShinyGold:", bagsInShinyGold)
		}
	}
	fmt.Println("bagsWithShinyGold:", bagsWithShinyGold)
}

func (b *Bag) children(bag string, count int, bagAllMap map[string]map[string]int) {
	var newChildren []string
	var newChildrenCount []int
	for k, v := range bagAllMap[bag] {
		i, found := find(&b.childNames, k)
		if found == false {
			b.childNames = append(b.childNames, k)
			b.childCounts = append(b.childCounts, v*count)
		} else {
			b.childCounts[i] = b.childCounts[i] + v*count
		}
		newChildren = append(newChildren, k)
		newChildrenCount = append(newChildrenCount, v*count)
	}
	for i, v := range newChildren {
		b.children(v, newChildrenCount[i], bagAllMap)
	}
}

func find(slice *[]string, val string) (int, bool) {
	for i, item := range *slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
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
