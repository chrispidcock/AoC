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
	part := "B"

	logsActive := false
	fmt.Println("logsActive ", logsActive)
	defer timeTrack(time.Now(), "day14")
	input, _ := readLines("initialization.txt")
	// input, _ := readLines("initialization-test-p1.txt")
	// input, _ := readLines("initialization-test-p2.txt")

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
	logging("masks", masks, false)
	logging("mem_init", mem_init, false)
	logging("mem_init_val", mem_init_val, false)
	memory := make(map[int64]int64)

	for j := range masks {
		fmt.Println("------------ ", masks[j], " == ", j)
		for i := range mem_init[j] {
			if part == "A" {
				logging("mem_init_val[j][i]", mem_init_val[j][i], false)
				val := fmt.Sprintf("%036b", mem_init_val[j][i])
				masked_val := applyMask(masks[j], val)
				v_int, _ := strconv.ParseInt(masked_val, 2, 64)
				memory[mem_init[j][i]] = v_int
			}
			if part == "B" {
				init_mem := fmt.Sprintf("%036b", mem_init[j][i])
				mem_slice := applyMaskMem(masks[j], init_mem)
				for mem := range mem_slice {
					memory[mem_slice[mem]] = mem_init_val[j][i]
				}
			}
		}
	}

	mem_val_sum := mapSum(memory)
	logging("---Answer---", "", true)
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

func applyMaskMem(m string, v string) []int64 {
	val_rune := []rune(v)
	mask_rune := []rune(m)
	var bit_len int = 0
	var masked_mem []int64

	for i := range mask_rune {
		if mask_rune[i] != '0' {
			val_rune[i] = mask_rune[i]
		}
		if mask_rune[i] == 'X' {
			bit_len = bit_len + 1
		}
	}

	for i := 0; i <= pow(2, bit_len)+1; i++ {
		val_rune_copy := []rune(string(val_rune))
		bits := fmt.Sprintf("%036b", i)
		bit_rune := []rune(bits)
		bit_correct := []rune(string(bit_rune[(len(bit_rune) - bit_len):]))
		for bit := range bit_correct {
			logging("bit", string(bit_correct[bit]), false)
			for j := range val_rune_copy {
				if val_rune_copy[j] == 'X' {
					val_rune_copy[j] = bit_correct[bit]
					break
				}
			}
		}

		v_int, _ := strconv.ParseInt(string(val_rune_copy), 2, 64)
		masked_mem = append(masked_mem, v_int)
	}
	return masked_mem
}

func pow(x int, y int) int {
	var result int = 1
	for i := 1; i <= y; i++ {
		result = result * i
	}
	return result
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
