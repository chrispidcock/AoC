package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

// TicketCalc allows methods to update the Ticket details
type TicketCalc interface {
	NewLowerRow() int
	NewTopRow() int
	NewLowerCol() int
	NewTopCol() int
	SeatID() int
}

// Ticket holds detils of a ticket as we calculate its row/col and seatID
type Ticket struct {
	binaryCode string
	frontRow   int
	backRow    int
	leftCol    int
	rightCol   int
	seatID     int
}

// NewLowerRow returns a higher value for the lowest possible row
func (t *Ticket) NewLowerRow() int {
	return t.frontRow + ((t.backRow-t.frontRow)+1)/2
}

// NewTopRow returns a lower value for the highest possible row
func (t *Ticket) NewTopRow() int {
	return t.backRow - ((t.backRow-t.frontRow)+1)/2
}

// NewLowerCol returns a higher value for the lowest possible col
func (t *Ticket) NewLowerCol() int {
	return t.leftCol + ((t.rightCol-t.leftCol)+1)/2
}

// NewTopCol returns a lower value for the highest possible col
func (t *Ticket) NewTopCol() int {
	return t.rightCol - ((t.rightCol-t.leftCol)+1)/2
}

// SeatID calculates the seatid
func (t *Ticket) SeatID() int {
	return t.backRow*8 + t.rightCol
}

func main() {
	defer timeTrack(time.Now(), "day5")
	tickets, _ := readLines("tickets.txt")
	var parsedTickets []Ticket
	var topSeatID int

	for _, code := range tickets {
		var t = Ticket{code, 0, 127, 0, 7, 0}
		// fmt.Println(code)
		for _, c := range code {
			switch c {
			case 'B':
				t.frontRow = t.NewLowerRow()
			case 'F':
				t.backRow = t.NewTopRow()
			case 'R':
				t.leftCol = t.NewLowerCol()
			case 'L':
				t.rightCol = t.NewTopCol()
			}
		}
		t.seatID = t.SeatID()
		if topSeatID < t.seatID {
			topSeatID = t.seatID
		}
		// fmt.Println(t)
		parsedTickets = append(parsedTickets, t)
	}
	fmt.Println("topSeatID:", topSeatID)
	yourID := yourSeatID(parsedTickets)
	fmt.Println("yourSeatID:", yourID)

}

func yourSeatID(pt []Ticket) int {
	seatIDS := make([]int, len(pt))
	for i, s := range pt {
		seatIDS[i] = s.seatID
	}
	sort.Ints(seatIDS[:])
	for i, id := range seatIDS {
		if i > 1 && seatIDS[i-1] == id-2 {
			return id - 1
		}
	}
	return 0
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
