package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	TypeSeed = iota
	TypeSoil
	TypeFertilizer
	TypeWater
	TypeLight
	TypeTemperature
	TypeHumidity
	TypeLocation
)

func intOrPanic(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return i
}

type Element struct {
	Type   uint8
	SelfID int
	Next   *Element
	LastID int // quick and dirty lol
}

type Range struct {
	Shift  int
	Start  int
	End    int
	Length int
}

func NewRange(Shift, Start, Length int) Range {
	return Range{
		Shift:  Shift,
		Start:  Start,
		End:    Start + Length - 1, // the last still mappable address
		Length: Length,
	}
}

func (r *Range) Innit(val int) bool {
	return val >= r.Start && val <= r.End // End is the last mappable addr, so we have to check inclusively here
}

func (r *Range) MaybeMap(val int) int {
	if r.Innit(val) {
		return (val - r.Start) + r.Shift
	} else {
		return val
	}
}

type MapThing struct {
	Type   uint8
	Ranges []Range
}

func (mt *MapThing) MapVal(val int) int {
	for _, mp := range mt.Ranges { // this requires ranges to be sorted
		if mp.End >= val {
			return mp.MaybeMap(val)
		}
	}
	return val // lol?
}

func (mt *MapThing) Prepare() {
	slices.SortFunc(mt.Ranges, func(a, b Range) int {
		return a.Start - b.Start
	})
	for i := 2; i < len(mt.Ranges); i++ {
		if mt.Ranges[i-1].Start > mt.Ranges[i].Start { // check for order
			fmt.Printf("%d > %d\n", mt.Ranges[i-1].Start, mt.Ranges[i].Start)
			panic("nemjo1")
		}
		if mt.Ranges[i-1].End > mt.Ranges[i].Start { // check for overlaps
			fmt.Printf("%d >= %d\n", mt.Ranges[i-1].End, mt.Ranges[i].Start)
			panic("nemjo2")
		}
	}
}

func LoadMap(scanner *bufio.Scanner, typ uint8) *MapThing {

	mt := MapThing{
		Type:   typ,
		Ranges: make([]Range, 0),
	}

	defer mt.Prepare()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			return &mt
		}
		strs := strings.Split(line, " ")
		shift := intOrPanic(strs[0])
		start := intOrPanic(strs[1])
		length := intOrPanic(strs[2])
		mt.Ranges = append(mt.Ranges, NewRange(shift, start, length))
	}

	return &mt
}

func LoadChains() map[int]*Element {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	seedsLine := scanner.Text()
	seedsStrs := strings.Split(seedsLine[7:], " ")

	seedCnt := len(seedsStrs)
	seedsMap := make(map[int]*Element, seedCnt)
	lastSet := make([]*Element, seedCnt)

	for i, seedStr := range seedsStrs {
		seedID := intOrPanic(seedStr)
		seedElm := &Element{
			Type:   TypeSeed,
			SelfID: seedID,
			Next:   nil,
		}
		seedsMap[seedID] = seedElm
		lastSet[i] = seedElm
	}

	scanner.Scan() // skip empty line

	var t uint8
	for t = TypeSoil; t <= TypeLocation; t++ {
		scanner.Scan()                // skip header
		mapper := LoadMap(scanner, t) // scans until the first empty line

		newSet := make([]*Element, seedCnt)
		for i, lElm := range lastSet {

			nElmID := mapper.MapVal(lElm.SelfID)

			nElm := &Element{
				Type:   t,
				SelfID: nElmID,
				Next:   nil,
			}

			lElm.Next = nElm
			newSet[i] = nElm

		}

		lastSet = newSet

	}

	// debug print
	for seedID, seed := range seedsMap {
		fmt.Printf("%d ", seedID)

		elm := seed.Next
		lastID := elm.SelfID
		for elm != nil {
			lastID = elm.SelfID
			fmt.Printf("-> %d ", elm.SelfID)
			elm = elm.Next
		}
		fmt.Printf("\n")
		seed.LastID = lastID

	}

	return seedsMap
}

func main() {
	seedsMap := LoadChains()

	const MaxUint = ^uint(0)
	var smallest = int(MaxUint >> 1)
	for _, v := range seedsMap {
		if v.LastID < smallest {
			smallest = v.LastID
		}
	}
	fmt.Println(smallest)
}
