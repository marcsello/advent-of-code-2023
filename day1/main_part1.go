package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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
		last := -1
		for _, c := range line {
			val, err := strconv.Atoi(string(rune(c)))
			if err != nil {
				continue
			}

			if first == -1 {
				first = val
			} else {
				last = val
			}

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
