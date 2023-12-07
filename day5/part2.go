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
	TypeSeed = iota // seed is lowest
	TypeSoil
	TypeFertilizer
	TypeWater
	TypeLight
	TypeTemperature
	TypeHumidity
	TypeLocation // location is highest
)

func intOrPanic(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return i
}

type Range struct {
	SelfType uint8
	RelShift int
	DstStart int
	Start    int
	DstEnd   int
	End      int
	Length   int
}

func NewRange(t uint8, dstStart, start, len int) Range {
	relShift := dstStart - start
	inclEnd := start + len - 1 // the last still mappable address
	return Range{
		SelfType: t,
		DstStart: dstStart,
		RelShift: relShift, // add this to an absolute value to shift it
		Start:    start,
		DstEnd:   inclEnd + relShift,
		End:      inclEnd,
		Length:   len,
	}
}

func NewIdentityRange(t uint8, start, len int) Range {
	relShift := 0
	inclEnd := start + len - 1 // the last still mappable address
	return Range{
		SelfType: t,
		DstStart: start,
		RelShift: relShift, // add this to an absolute value to shift it
		Start:    start,
		DstEnd:   inclEnd + relShift,
		End:      inclEnd,
		Length:   len,
	}
}

func (r *Range) String() string {
	return fmt.Sprintf("%d:[%d - %d (%d)] -{%d}-> %d:[%d - %d (%d)]", r.SelfType, r.Start, r.End, r.Length, r.RelShift, r.SelfType+1, r.DstStart, r.DstEnd, r.Length)
}

func (r *Range) IntersectMapNext(next *Range) *Range {
	// r is upper
	// r.SelfType must be always smaller than r2.SelfType

	if r.SelfType >= next.SelfType {
		panic("nemjo")
	}

	altStart := max(r.DstStart, next.Start)
	altEnd := min(r.DstEnd, next.End)

	if altEnd < altStart {
		// this would be an empty, or invalid range...
		return nil // no intersection
	}

	intersect := NewRange(next.SelfType, next.DstStart+(altStart-next.Start), altStart, altEnd-altStart+1)

	return &intersect
}

type MapThing struct {
	Type   uint8
	Ranges []Range
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
		mt.Ranges = append(mt.Ranges, NewRange(typ, shift, start, length))
	}

	return &mt
}

func LoadAllMaps() []*MapThing {
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

	maps := make([]*MapThing, 8)

	seedMap := &MapThing{
		Type:   TypeSeed,
		Ranges: make([]Range, seedCnt/2),
	}
	j := 0
	for i := 0; i < seedCnt; i += 2 {
		seedMap.Ranges[j] = NewIdentityRange(TypeSeed, intOrPanic(seedsStrs[i]), intOrPanic(seedsStrs[i+1]))
		j++
	}
	maps[0] = seedMap

	scanner.Scan() // skip empty line

	var t uint8
	for t = TypeSoil; t <= TypeLocation; t++ {
		scanner.Scan()                // skip header
		maps[t] = LoadMap(scanner, t) // scans until the first empty line
	}

	return maps
}

var FoundLocationRanges = make([]*Range, 0)

func BlackMagic(maps []*MapThing, upperRange *Range, lowerType uint8, end uint8) {
	fmt.Printf("->\n")

	if upperRange.SelfType >= lowerType {
		panic("upper type must have lower id")
	}

	lastRangeEnd := -1

	for _, lowerRange := range maps[lowerType].Ranges {

		// consider gap since the last range
		if lastRangeEnd+1 < lowerRange.Start {
			gapRange := NewIdentityRange(lowerType, lastRangeEnd+1, lowerRange.Start-lastRangeEnd-1)
			fmt.Println(" Gap test: ", gapRange.String())
			intersect := upperRange.IntersectMapNext(&gapRange)
			if intersect != nil {
				if intersect.SelfType == end {
					fmt.Printf("Found gap Location range! %s\n", intersect.String())
					FoundLocationRanges = append(FoundLocationRanges, intersect)
				} else {
					fmt.Println("  Step Into: ", intersect.String())
					BlackMagic(maps, intersect, lowerType+1, end)
				}
			}
		}

		// consider the current range
		fmt.Println(" Testing: ", lowerRange.String())
		intersect := upperRange.IntersectMapNext(&lowerRange)

		if intersect != nil {
			if intersect.SelfType == end {
				fmt.Printf("Found valid Location range! %s\n", intersect.String())
				FoundLocationRanges = append(FoundLocationRanges, intersect)
			} else {
				fmt.Println("  Step Into: ", intersect.String())
				BlackMagic(maps, intersect, lowerType+1, end)
			}
		}

		lastRangeEnd = lowerRange.End

	}

	// Tail test
	if lastRangeEnd <= upperRange.End {
		tailRange := NewIdentityRange(lowerType, lastRangeEnd+1, upperRange.End-lastRangeEnd+1)
		fmt.Println(" Tail test: ", tailRange.String())
		intersect := upperRange.IntersectMapNext(&tailRange)
		if intersect != nil {
			if intersect.SelfType == end {
				fmt.Printf("Found tailed Location range! %s\n", intersect.String())
				FoundLocationRanges = append(FoundLocationRanges, intersect)
			} else {
				fmt.Println("  Step Into: ", intersect.String())
				BlackMagic(maps, intersect, lowerType+1, end)
			}
		}
	}

	fmt.Println("<-")
}

func main() {
	maps := LoadAllMaps()

	for _, mp := range maps {
		fmt.Printf("%d\t%d\n", mp.Type, len(mp.Ranges))
	}

	for _, seedRange := range maps[TypeSeed].Ranges {
		fmt.Println(" Testing: ", seedRange.String())
		BlackMagic(maps, &seedRange, TypeSeed+1, TypeLocation)
	}

	slices.SortFunc(FoundLocationRanges, func(a, b *Range) int {
		return a.DstStart - b.DstStart
	})

	fmt.Println(FoundLocationRanges[0].DstStart)

}
