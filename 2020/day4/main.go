package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type passport struct {
	BirthYear   int    //byr (Birth Year)
	IssueYear   int    //iyr (Issue Year)
	ExpYear     int    //eyr (Expiration Year)
	Height      int    //hgt (Height)
	HeightUnits string //hgt (Height)
	HairColour  string //hcl (Hair Color)
	EyeColour   string //ecl (Eye Color)
	PassportID  string //pid (Passport ID)
	CountryID   int    //cid (Country ID)
	Valid       bool
}

func main() {
	defer timeTrack(time.Now(), "day4")
	rows, _ := readLines("passports.txt")
	var passportsValid []passport
	var passportsInvalid []passport

	var textBlocks []string
	var textBlock string
	fmt.Println(len(rows))
	for row, text := range rows {
		fmt.Println(row)
		if text == "" {
			fmt.Println(textBlock)
			textBlocks = append(textBlocks, textBlock)
			textBlock = ""
		} else {
			textBlock = textBlock + text + " "
		}
	}

	for _, entry := range textBlocks {
		parsedPassport, valid := parseText(entry)
		if valid == true {
			passportsValid = append(passportsValid, parsedPassport)
		} else {
			passportsInvalid = append(passportsInvalid, parsedPassport)
		}
	}
	fmt.Println("--Valid--")
	for _, entry := range passportsValid {
		fmt.Println(entry)
	}
	fmt.Println("--INVALID--")
	for _, entry := range passportsInvalid {
		fmt.Println(entry)
	}

	fmt.Println(len(passportsValid))
}

func parseText(s string) (passport, bool) {
	var newPassport = passport{}

	re := regexp.MustCompile(`(?:byr:)([12][0-9]{3})`)
	// re := regexp.MustCompile(`byr:(\S*)`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			newPassport.BirthYear = v
		} else {
			fmt.Println(err)
		}
	}

	re = regexp.MustCompile(`iyr:([12][0-9]{3})`)
	// re = regexp.MustCompile(`iyr:(\S*)`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			newPassport.IssueYear = v
		} else {
			fmt.Println(err)
		}
	}

	re = regexp.MustCompile(`eyr:([12][0-9]{3})`)
	// re = regexp.MustCompile(`eyr:(\S*)`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			newPassport.ExpYear = v
		} else {
			fmt.Println(err)
		}
	}

	re = regexp.MustCompile(`hgt:([1-9][0-9]*)`)
	// re = regexp.MustCompile(`hgt:(\S*)`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			newPassport.Height = v
		} else {
			fmt.Println(err)
		}
	}

	re = regexp.MustCompile(`hgt:[1-9][0-9]*([cmin]{2})`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		newPassport.HeightUnits = match[1]
	}

	// re = regexp.MustCompile(`hcl:(#(?:[0-9a-fA-F]{3}){1,2})`)
	re = regexp.MustCompile(`hcl:(#[0-9a-fA-F]{6})`)
	// re = regexp.MustCompile(`hcl:(\S*)`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		newPassport.HairColour = match[1]
	}

	// re = regexp.MustCompile(`ecl:(#(?:[0-9a-fA-F]{3}){1,2})`)
	re = regexp.MustCompile(`ecl:(\S*)`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		newPassport.EyeColour = match[1]
	}

	re = regexp.MustCompile(`pid:([0-9]*)`)
	// re = regexp.MustCompile(`pid:(\S*)`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		newPassport.PassportID = match[1]
	}

	re = regexp.MustCompile(`cid:([0-9]*)`)
	// re = regexp.MustCompile(`cid:(\S*)`)
	match = re.FindStringSubmatch(s)
	if len(match) > 1 {
		v, err := strconv.Atoi(match[1])
		if err == nil {
			newPassport.CountryID = v
		} else {
			fmt.Println(err)
		}
	}

	validEyeColour := map[string]bool{
		"amb": true,
		"blu": true,
		"brn": true,
		"gry": true,
		"grn": true,
		"hzl": true,
		"oth": true,
	}

	if newPassport.BirthYear != 0 &&
		newPassport.BirthYear >= 1920 &&
		newPassport.BirthYear <= 2002 &&
		newPassport.IssueYear != 0 &&
		newPassport.IssueYear >= 2010 &&
		newPassport.IssueYear <= 2020 &&
		newPassport.ExpYear != 0 &&
		newPassport.ExpYear >= 2020 &&
		newPassport.ExpYear <= 2030 &&
		newPassport.Height != 0 &&
		newPassport.HeightUnits != "" &&
		((newPassport.HeightUnits == "cm" && newPassport.Height >= 150 && newPassport.Height <= 193) ||
			(newPassport.HeightUnits == "in" && newPassport.Height >= 59 && newPassport.Height <= 76)) &&
		newPassport.HairColour != "" &&
		newPassport.EyeColour != "" &&
		validEyeColour[newPassport.EyeColour] &&
		// newPassport.CountryID != "CountryID" &&
		len(newPassport.PassportID) == 9 {
		newPassport.Valid = true
		return newPassport, true
	} else {
		fmt.Println(s)
		return newPassport, false
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
