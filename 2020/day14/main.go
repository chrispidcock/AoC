package main

// https://adventofcode.com/2020/day/14

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

func main() {
	// part := "A"

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day14")
	input, _ := readLines("initialization.txt")
	// input, _ := readLines("initialization-test.txt")

	var masks []string
	var mem_init [][]int64
	var mem_init_val [][]int64
	var mem_mask []int

	for i := range input {
		s, err := getMask(input[i])
		if err == nil {
			masks = append(masks, s)
			mem_mask = append(mem_mask, len(masks)-1)
			mem_init = append(mem_init, make([]int64, 0))
			mem_init_val = append(mem_init_val, make([]int64, 0))

		} else {
			val, _ := getMemory(input[i])
			mem_init[mem_mask[len(mem_mask)-1]] = append(mem_init[mem_mask[len(mem_mask)-1]], int64(val))
			val, _ = getMemVal(input[i])
			mem_init_val[mem_mask[len(mem_mask)-1]] = append(mem_init_val[mem_mask[len(mem_mask)-1]], int64(val))
		}
	}

	memory := make(map[int64]int64)

	for mask_i := range masks {
		logging("masks[mask_i]", masks[mask_i], false)
		for i := range mem_init[mask_i] {
			logging("mem_init_val[mask_i][i]", mem_init_val[mask_i][i], false)
			val := fmt.Sprintf("%036b", mem_init_val[mask_i][i])
			logging("val", val, false)
			masked_val := applyMask(masks[mask_i], val)
			logging("masked_val", masked_val, false)
			v_int, _ := strconv.ParseInt(masked_val, 2, 64)
			logging("v_int", v_int, false)
			memory[mem_init[mask_i][i]] = v_int
		}
	}

	// logging("memory", memory, true)
	mem_val_sum := mapSum(memory)
	logging("---Part 1 Answer---", "", true)
	logging("mem_val_sum", mem_val_sum, true)
}

func applyMask(mask string, val string) string {
	val_rune := []rune(val)
	mask_rune := []rune(mask)
	for i := range mask_rune {
		if mask_rune[i] != 'X' {
			val_rune[i] = mask_rune[i]
		}
	}
	return string(val_rune)
}

func applyMaskMem(mask string, val string) []int64 {
	val_rune := []rune(val)
	mask_rune := []rune(mask)
	bit_len := -1
	var masked_mem []int64

	for i := range mask_rune {
		if mask_rune[i] != 'X' {
			val_rune[i] = mask_rune[i]
		} else {
			bit_len = bit_len + 1
		}
	}

	for i := 0; i < 2^bit_len; i++ {
		val_rune_copy := []rune(string(val_rune))
		bit := fmt.Sprintf("%0%vb", i, 2^bit_len)
		bit_rune := rune(bit)
		for bit := range bit_rune {
			for j := range val_rune_copy {
				if val_rune_copy[j] == 'X' {
					val_rune_copy[j] = bit_rune[bit]
					break
				}
			}
		}

		v_int, _ := strconv.ParseInt(masked_val, 2, 64)
		masked_mem = append(masked_mem, v_int)
	}
	return masked_mem
}

func mapSum(m map[int64]int64) int64 {
	var m_sum int64 = 0
	for i := range m {
		if m[i] != 0 {
			m_sum = m_sum + m[i]
		}
	}
	return m_sum
}

func getMask(s string) (string, error) {
	re := regexp.MustCompile(`mask = (.{36})`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		return match[1], nil
	}
	return s, errors.New("error")
}

func getMemory(s string) (int64, error) {
	re := regexp.MustCompile(`mem\[(.*)\]`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			return int64(v), nil
		} else {
			fmt.Println(err)
		}
	}
	return 0, errors.New("error")
}

func getMemVal(s string) (int64, error) {
	re := regexp.MustCompile(`\] = ([0-9]*)$`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			return int64(v), nil
		} else {
			fmt.Println(err)
		}
	}
	return 0, errors.New("error")
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
