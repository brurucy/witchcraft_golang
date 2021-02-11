package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

type KeyValue struct {
	key   int
	value string
}

type MinMaxList struct {
	Indexes []int
	Height  int
	Min     int
	Max     int
}

type TeleportList struct {
	Sublists      []MinMaxList
	Data          []int
	CurrentHeight int
	Length        int
}

type SortedTeleportListSublist []MinMaxList

func (s SortedTeleportListSublist) Len() int {
	return len(s)
}

func (s SortedTeleportListSublist) Swap(i, j int) {
	s[i].Min, s[i].Min = s[j].Min, s[i].Min
}
func (s SortedTeleportListSublist) Less(i, j int) bool {
	return s[i].Min < s[j].Min
}

func NewTeleportList() TeleportList {
	teleportList := TeleportList{}
	teleportList.CurrentHeight = -1
	return teleportList
}

func insert(a []int, index int, value int) []int {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func insort(a []int, value int) []int {

	idx := sort.Search(len(a), func(i int) bool { return a[i] >= value })

	return insert(a, idx, value)

}

func BinarySearch(a []int, value int) bool {

	idx := sort.Search(len(a), func(i int) bool { return a[i] >= value })

	if idx == len(a) {

		return false

	} else {

		return true
	}

}

func (t *TeleportList) Add(value int) {
	candidateHeight := int(math.Abs(math.Log2(rand.Float64())))

	if candidateHeight > t.CurrentHeight {

		heightDiff := candidateHeight - t.CurrentHeight

		for i := 0; i < heightDiff; i++ {

			newMMList := MinMaxList{}
			newMMList.Max = math.MinInt64
			newMMList.Min = math.MaxInt64
			newMMList.Height = t.CurrentHeight + i + 1

			t.Sublists = append(t.Sublists, newMMList)

		}

		t.CurrentHeight = candidateHeight

	}

	candidateHeightSublist := &t.Sublists[candidateHeight]

	candidateHeightSublist.Indexes = insort(candidateHeightSublist.Indexes, value)

	t.Length += 1
	if candidateHeightSublist.Max < value {
		candidateHeightSublist.Max = value
	}
	if candidateHeightSublist.Min > value {
		candidateHeightSublist.Min = value
	}

	t.Data = insort(t.Data, value)

}

func (t TeleportList) Find(value int) bool {

	if last := len(t.Sublists) - 1; last >= 0 {

		for i := last; i >= 0; i-- {

			if t.Sublists[i].Min <= value {

				if t.Sublists[i].Max >= value {

					bsearch := BinarySearch(t.Sublists[i].Indexes, value)

					if bsearch == true {

						return true

					}

				}

			}

		}
	}

	return false

}

func (t TeleportList) Index(value int) int {

	return t.Data[value-1]

}

func main() {

	tlist := NewTeleportList()

	start := time.Now()

	for i := 1; i < 1_000_000; i++ {

		tlist.Add(i)

	}

	elapsed := time.Since(start)

	fmt.Println("Elapsed add: ", elapsed)

	start = time.Now()

	for i := 1; i < 1_000_000; i++ {

		tlist.Find(i)

	}

	elapsed = time.Since(start)

	fmt.Println("Elapsed find: ", elapsed)

	start = time.Now()

	for i := 1; i < 1_000_000; i++ {

		BinarySearch(tlist.Data, i)
	}

	elapsed = time.Since(start)

	fmt.Println("Elapsed find slice: ", elapsed)

	start = time.Now()

	for i := 1; i < 1_000_000; i++ {

		tlist.Index(i)

	}

	elapsed = time.Since(start)

	fmt.Println("Elapsed index: ", elapsed)

	start = time.Now()

	mapp := make(map[int]bool)

	for i := 1; i < 1_000_000; i++ {

		mapp[i] = true

	}

	elapsed = time.Since(start)

	fmt.Println("Elapsed map add: ", elapsed)

	start = time.Now()

	for i := 1; i < 1_000_000; i++ {

		if mapp[i] {

		}

	}

	elapsed = time.Since(start)

	fmt.Println("elapsed map find: ", elapsed)
}
