package main

import (
	"bufio"
	"os"
	"strconv"
	"unicode"
)

type coordinate struct {
	x int
	y int
}

type position struct {
	coordinate  coordinate
	originalPos coordinate
}

type numberType struct {
	value     int
	size      int
	positions []coordinate
}

type numberString struct {
	value     string
	positions []coordinate
}

type gear struct {
	partNumbers []int
}

var numbers []numberType
var symbols []coordinate
var possiblePositions []position

var gears = map[coordinate]gear{}

var possibleSymbols = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0021, 0x002D, 1},
		{0x002F, 0x002F, 1},
		{0x003A, 0x0040, 1},
		{0x005B, 0x0060, 1},
		{0x007B, 0x00BF, 1},
	},
}

var vectors = []coordinate{
	{-1, -1},
	{0, -1},
	{1, -1},
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
}

var maxX uint
var maxY uint

func main() {
	println(part1())
	println(part2())
}

func part1() int {
	var lines []string

	dat, err := os.Open("./input.txt")
	check(err)
	defer dat.Close()

	scanner := bufio.NewScanner(dat)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	maxX = uint(len(lines[0]))
	maxY = uint(len(lines))
	for idx, line := range lines {
		chars := []rune(line)
		tmpNumber := numberString{"", []coordinate{}}

		for i := 0; i < len(chars); i++ {
			char := chars[i]
			if unicode.IsDigit(char) {
				// Append char
				tmpNumber.value = tmpNumber.value + string(char)
				// Set start Pos

				tmpNumber.positions = append(tmpNumber.positions, coordinate{i, idx})

				// Last digit in line or before '.' -> append tmpNumber
				if i+1 >= len(chars) || unicode.Is(possibleSymbols, chars[i+1]) || chars[i+1] == '.' {
					val, _ := strconv.Atoi(tmpNumber.value)
					numbers = append(numbers, numberType{
						val,
						len(tmpNumber.value),
						tmpNumber.positions,
					})
					tmpNumber = numberString{"", []coordinate{}}
				}
			} else if unicode.Is(possibleSymbols, char) {
				symbols = append(symbols, coordinate{i, idx})
				gears[coordinate{i, idx}] = gear{partNumbers: []int{}}
			}
		}
	}
	return filterAdjacentSymbols()
}

func filterAdjacentSymbols() int {
	result := 0
	tmpNumbers := make([]numberType, len(numbers))
	copy(tmpNumbers, numbers)

	// replace symbol slice with valid spots
	for _, symbol := range symbols {
		for _, vector := range vectors {
			tmp := coordinate{symbol.x + vector.x, symbol.y + vector.y}
			if (uint)(tmp.x) <= maxX && (uint)(tmp.y) <= maxY {
				flag := false
				for _, coordinate := range possiblePositions {
					if coordinate.coordinate.y == tmp.y && coordinate.coordinate.x == tmp.x {
						flag = true
					}
				}
				if !flag {
					possiblePositions = append(possiblePositions, position{coordinate: tmp, originalPos: symbol})
				}
			}
		}
	}
	// iterate and find all number parts
	for _, symbol := range possiblePositions {
		found := false
		for i, number := range tmpNumbers {
			if found {
				break
			}
			for _, position := range number.positions {
				if position.x == symbol.coordinate.x && position.y == symbol.coordinate.y && !found {
					result += number.value
					found = true
					tmpNumbers = remove(tmpNumbers, i)
					break
				}
			}
		}
	}
	return result
}

func part2() int {
	result := 0
	tmpNumbers := make([]numberType, len(numbers))
	copy(tmpNumbers, numbers)

	for _, symbol := range possiblePositions {
		found := false
		for i, number := range tmpNumbers {
			if found {
				break
			}
			for _, position := range number.positions {
				if position.x == symbol.coordinate.x && position.y == symbol.coordinate.y && !found {
					tmp := gears[symbol.originalPos]
					tmp.partNumbers = append(tmp.partNumbers, number.value)

					gears[symbol.originalPos] = tmp
					found = true
					tmpNumbers = remove(tmpNumbers, i)
					break
				}
			}
		}
	}

	for _, gear := range gears {
		if len(gear.partNumbers) == 2 {
			result += gear.partNumbers[0] * gear.partNumbers[1]
		}
	}
	return result
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
