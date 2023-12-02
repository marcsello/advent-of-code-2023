package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var tokens = map[string]int{
	"0":     0,
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var sum int
	for scanner.Scan() {
		line := scanner.Text()
		first := -1
		firstPos := len(line) + 1
		last := -1
		lastPos := -1

		for token, val := range tokens {
			idxf := strings.Index(line, token)
			if idxf == -1 {
				continue
			}

			if idxf < firstPos {
				first = val
				firstPos = idxf
			}

			idxl := strings.LastIndex(line, token)

			if idxl > lastPos {
				last = val
				lastPos = idxl
			}

		}
		if first == -1 {
			panic("nemjolet")
		}

		if last == -1 {
			last = first
		}
		num := first*10 + last
		fmt.Printf("%d %d | %d\n", first, last, num)
		sum += num
	}

	fmt.Println(sum)
}
