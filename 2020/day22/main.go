package main

// https://adventofcode.com/2020/day/22

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func main() {
	// part := "B"

	defer timeTrack(time.Now(), "day22")
	input, _ := readLines("hands.txt")
	// input, _ := readLines("hands-test.txt")

	var players []int
	var cards [][]int

	var player int
	var card int

	for i, s := range input {
		fmt.Println("s", s)
		if i == 0 {
			player, _ = getInt(s)
			players = append(players, player)
			cards = append(cards, make([]int, 0))
		} else if input[i-1] == "" {
			fmt.Println("HIT")
			player, _ = getInt(s)
			players = append(players, player)
			cards = append(cards, make([]int, 0))
		} else if input[i] == "" {
			continue
		} else {
			card, _ = getInt(s)
			cards[len(players)-1] = append(cards[len(players)-1], card)
		}
	}

	fmt.Println("players", players)
	fmt.Println("cards", cards)

	end := false
	round := 0
	var m_c_p int
	for !end {
		round = round + 1
		fmt.Println("\n -- ROUND", round, "--")
		min_cards := 100
		max_card := -1
		for i := range cards {
			fmt.Println("Player", i, cards[i])
			if len(cards[i]) < min_cards {
				min_cards = len(cards[i])
			}
		}
		if min_cards == 0 {
			end = true
			break
		}
		var r_cards []int
		for i := range players {
			r_cards = append(r_cards, cards[i][0])
			if cards[i][0] > max_card {
				max_card = cards[i][0]
				m_c_p = i
			}
		}
		sort.Ints(r_cards)
		fmt.Println(r_cards)
		for j := len(r_cards) - 1; j >= 0; j-- {
			cards[m_c_p] = append(cards[m_c_p], r_cards[j])
		}
		for i := range players {
			fmt.Println("Player", i, cards[i])
			cards[i] = cards[i][1:len(cards[i])]
			fmt.Println("Player", i, cards[i])
		}

	}
	fmt.Println("\n--- Answer Part 1 ---")
	score := 0
	for j := len(cards[m_c_p]) - 1; j >= 0; j-- {
		score = score + cards[m_c_p][j]*(len(cards[m_c_p])-j)
		fmt.Println(cards[m_c_p][j], "*", (len(cards[m_c_p]) - j))
	}
	fmt.Println("Player", m_c_p, "score=", score)
}

func getInt(s string) (int, error) {
	matches := 20
	re := regexp.MustCompile(`[0-9]+`)
	match := re.FindAllStringSubmatch(s, matches)
	var v int
	var err error
	if len(match) == 1 {
		for i := range match {
			v, err = strconv.Atoi(match[i][0])
			if err != nil {
				break
			} else {
				return v, nil
			}
		}
	}
	return v, errors.New("No Integers")
}

// removeIndex. Swap the element to delete with the one at the end of the slice and then return the n-1 first elements
func removeIndexSS(s [][]int, i int) [][]int {
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
