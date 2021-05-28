package main

// https://adventofcode.com/2020/day/22

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Game struct {
	rounds []Round
}

type Round struct {
	players [2]Player
}

type Player struct {
	hand []int
}

func main() {
	part := "B"

	defer timeTrack(time.Now(), "day22")
	input, _ := readLines("hands.txt")
	// input, _ := readLines("hands-test.txt")
	// input, _ := readLines("hands-test-2.txt")

	var games []Game
	// var g_r_h_c [][][2][50]int
	// var round [][2][50]int
	// var hand [50]int
	var err error
	var winner = -1
	var g = 0
	var r = 0
	var p = -1
	var c int
	games = append(games, Game{})
	games[g].rounds = append(games[g].rounds, Round{})
	for i, s := range input {
		if i == 0 {
			_, err = getInt(s)
			if err == nil {
				c = 0
				p = p + 1
			}
		} else if input[i-1] == "" {
			_, err = getInt(s)
			if err == nil {
				c = 0
				p = p + 1
			}
		} else if input[i] == "" {
			continue
		} else {
			card, _ := getInt(s)
			games[g].rounds[r].players[p].hand = append(games[g].rounds[r].players[p].hand, card)
			c = c + 1
		}
	}

	games, winner, _, r = cardGameRecur(g, games, part)

	fmt.Println("\n--- Answer Part 1 ---")
	for j := range games {
		fmt.Println("Game", j)
		for i := range games[j].rounds {
			fmt.Println("Round", i, games[j].rounds[i])
		}
	}
	score := 0
	for i := range games[0].rounds[len(games[0].rounds)-1].players {
		if len(games[0].rounds[len(games[0].rounds)-1].players[i].hand) > 0 {
			for j := len(games[g].rounds[r].players[i].hand); j >= 0; j-- {
				fmt.Println(j)
				fmt.Println(games[0].rounds[len(games[0].rounds)-1].players[i].hand[j], "*", (len(games[g].rounds[r].players[i].hand) - j + 1))
				score = score + games[0].rounds[len(games[0].rounds)-1].players[i].hand[j]*(len(games[g].rounds[r].players[i].hand)-j+1)
			}
		}
	}
	fmt.Println("Player", winner, "score=", score)
}

func cardGameRecur(g int, games []Game, part string) ([]Game, int, int, int) {
	if g > len(games)-1 {
		games = append(games, Game{})
		games[g].rounds = append(games[g].rounds, Round{})
		for i := range games[g].rounds[0].players {
			for j := 1; j <= games[g-1].rounds[len(games[g-1].rounds)-2].players[i].hand[0]; j++ {
				games[g].rounds[0].players[i].hand = append(games[g].rounds[0].players[i].hand, games[g-1].rounds[len(games[g-1].rounds)-1].players[i].hand[j-1])
			}
		}
	}

	var winner int
	r := 0
	end := false

	fmt.Println("\n --- GAME", g, "---")
	for !end {
		// fmt.Println("\n -- ROUND", r, "--")
		min_cards := 100

		// Check size of hand
		for i := range games[g].rounds[r].players {
			// fmt.Println("Player", i, games[g].rounds[r].players[i])
			if len(games[g].rounds[r].players[i].hand) < min_cards {
				min_cards = len(games[g].rounds[r].players[i].hand)
			}
		}
		if min_cards == 0 {
			end = true
			break
		}

		if part == "B" {
			// Instant win player 1 if hands are same from current game previous round
		Next_Round:
			for i := range games[g].rounds {
				if i == r {
					continue
				}
				for player := range games[g].rounds[i].players {
					for j := range games[g].rounds[i].players[player].hand {
						if len(games[g].rounds[i].players[player].hand) != len(games[g].rounds[r].players[player].hand) {
							continue Next_Round
						}
						if games[g].rounds[i].players[player].hand[j] != games[g].rounds[r].players[player].hand[j] {
							continue Next_Round
						}
					}
				}
				return games, 0, g - 1, r
			}
		}

		games[g].rounds = append(games[g].rounds, Round{})
		// Compare top card
		max_card := -1
		winner := -1
		var r_cards []int
		for i := range games[g].rounds[r].players {
			r_cards = append(r_cards, games[g].rounds[r].players[i].hand[0])

			if games[g].rounds[r].players[i].hand[0] > max_card {
				max_card = games[g].rounds[r].players[i].hand[0]
				winner = i
			}

			for j := 1; j < len(games[g].rounds[r].players[i].hand); j++ {
				games[g].rounds[r+1].players[i].hand = append(games[g].rounds[r+1].players[i].hand, games[g].rounds[r].players[i].hand[j])
			}
		}
		// fmt.Println("winner=", winner)
		// fmt.Println(r_cards)

		// If len(deck) > current drawn card, determine winner of round by recursive Game
		if part == "B" {
			start_recursive := true
			for i, v := range r_cards {
				if len(games[g].rounds[r].players[i].hand)-1 < v {
					start_recursive = false
				}
			}
			if start_recursive {
				fmt.Println("Start New Recursive Game")
				games, winner, g, _ = cardGameRecur(g+1, games, part)
				fmt.Println("Returning to Game", g)
			}
		}

		// sort.Ints(r_cards) // Winners card goes on top of losers card

		if winner == 1 {
			for j := len(r_cards) - 1; j >= 0; j-- {
				games[g].rounds[r+1].players[winner].hand = append(games[g].rounds[r+1].players[winner].hand, r_cards[j])
			}
		} else {
			for j := range r_cards {
				games[g].rounds[r+1].players[winner].hand = append(games[g].rounds[r+1].players[winner].hand, r_cards[j])
			}
		}

		for i := range games[g].rounds[r+1].players {
			if len(games[g].rounds[r+1].players[i].hand) == 0 {
				end = true
			}
		}
		if g < len(games)-1 {
			games = games[:len(games)-1]
		}
		if end && winner >= 0 {
			return games, winner, g - 1, r
		}
		r = r + 1
	}
	return games, winner, g - 1, r
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
