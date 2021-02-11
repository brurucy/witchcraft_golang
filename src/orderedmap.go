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
	Indexes []KeyValue
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

func NewTeleportList() TeleportList {
	teleportList := TeleportList{}
	teleportList.CurrentHeight = -1
	return teleportList
}

func insertKeyValue(a []KeyValue, index int, value int) []KeyValue {
	if len(a) == index {
		return append(a, KeyValue{key: value})
	}
	a = append(a[:index+1], a[index:]...)
	a[index].key = value
	return a
}

func insertInt(a []int, index int, value int) []int {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func insortKeyValue(a []KeyValue, value int) []KeyValue {

	idx := sort.Search(len(a), func(i int) bool { return a[i].key >= value })

	return insertKeyValue(a, idx, value)

}

func insortInt(a []int, value int) []int {

	idx := sort.Search(len(a), func(i int) bool { return a[i] >= value })

	return insertInt(a, idx, value)

}

func BinarySearchKeyValue(a []KeyValue, value int) bool {

	idx := sort.Search(len(a), func(i int) bool { return a[i].key >= value })

	if idx == len(a) {

		return false

	} else {

		return true
	}

}

func BinarySearchInt(a []int, value int) bool {

	idx := sort.Search(len(a), func(i int) bool { return a[i] >= value })

	if idx == len(a) {

		return false

	} else {

		return true
	}

}

func InterpolationSearchInt(array []int, key int) int {

	min, max := array[0], array[len(array)-1]

	low, high := 0, len(array)-1

	for {
		if key < min {
			return low
		}

		if key > max {
			return high + 1
		}

		// make a guess of the location
		var guess int
		if high == low {
			guess = high
		} else {
			size := high - low
			offset := int(float64(size-1) * (float64(key-min) / float64(max-min)))
			guess = low + offset
		}

		// maybe we found it?
		if array[guess] == key {
			// scan backwards for start of value range
			for guess > 0 && array[guess-1] == key {
				guess--
			}
			return guess
		}

		// if we guessed to high, guess lower or vice versa
		if array[guess] > key {
			high = guess - 1
			max = array[high]
		} else {
			low = guess + 1
			min = array[low]
		}
	}
}

func InterpolationSearchKeyValue(array []KeyValue, key int) int {

	min, max := array[0].key, array[len(array)-1].key

	low, high := 0, len(array)-1

	for {
		if key < min {
			return low
		}

		if key > max {
			return high + 1
		}

		// make a guess of the location
		var guess int
		if high == low {
			guess = high
		} else {
			size := high - low
			offset := int(float64(size-1) * (float64(key-min) / float64(max-min)))
			guess = low + offset
		}

		// maybe we found it?
		if array[guess].key == key {
			// scan backwards for start of value range
			for guess > 0 && array[guess-1].key == key {
				guess--
			}
			return guess
		}

		// if we guessed to high, guess lower or vice versa
		if array[guess].key > key {
			high = guess - 1
			max = array[high].key
		} else {
			low = guess + 1
			min = array[low].key
		}
	}
}

func InterpolatedBinarySearchInt(a []int, key int) bool {

	idx := InterpolationSearchInt(a, key)

	if a[idx] == key {

		return true

	} else {

		return false

	}

}

func InterpolatedBinarySearchKeyValue(a []KeyValue, key int) bool {

	idx := InterpolationSearchKeyValue(a, key)

	if a[idx].key == key {

		return true

	} else {

		return false

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

	candidateHeightSublist.Indexes = insortKeyValue(candidateHeightSublist.Indexes, value)

	t.Length += 1
	if candidateHeightSublist.Max < value {
		candidateHeightSublist.Max = value
	}
	if candidateHeightSublist.Min > value {
		candidateHeightSublist.Min = value
	}

	//sort.Sort(t.SSublists)
	t.Data = insortInt(t.Data, value)

}

func (t TeleportList) Find(value int) bool {

	if last := len(t.Sublists) - 1; last >= 0 {

		for i := last; i >= 0; i-- {

			if t.Sublists[i].Min <= value {

				if t.Sublists[i].Max >= value {

					bsearch := BinarySearchKeyValue(t.Sublists[i].Indexes, value)
					//bsearch := InterpolatedBinarySearchKeyValue(t.Sublists[i].Indexes, value)

					if bsearch == true {

						return true

					}

				}

			}

		}
	}

	return false

}

func (t TeleportList) Index(value int) bool {

	return t.Find(t.Data[value-1])

}

func main() {

	var nums []int

	n := 10_000_000

	fmt.Println("Rng generation started")

	for i := 1; i < n; i++ {

		nums = append(nums, rand.Intn(20_000_000))
	}

	fmt.Println("Rng generation ended")

	tlist := NewTeleportList()

	start := time.Now()

	for i := 1; i < n; i++ {

		tlist.Add(i)

	}

	elapsed := time.Since(start)

	fmt.Println("TList Elapsed add: ", elapsed)

	start = time.Now()

	for i := 1; i < n; i++ {

		tlist.Find(i)

	}

	elapsed = time.Since(start)

	fmt.Println("Tlist Elapsed find: ", elapsed)

	start = time.Now()

	for i := 1; i < n; i++ {

		BinarySearchInt(tlist.Data, i)
		//InterpolatedBinarySearchInt(tlist.Data, i)
	}

	elapsed = time.Since(start)

	fmt.Println("Sorted List Elapsed find slice: ", elapsed)

	start = time.Now()

	for i := 1; i < n; i++ {

		tlist.Index(i)

	}

	elapsed = time.Since(start)

	fmt.Println("Tlist Elapsed index: ", elapsed)

	start = time.Now()

	mapp := make(map[int]bool)

	for i := 1; i < n; i++ {

		mapp[i] = true

	}

	elapsed = time.Since(start)

	fmt.Println("Hashmap Elapsed add: ", elapsed)

	start = time.Now()

	for i := 1; i < n; i++ {

		if mapp[i] {

		}

	}

	elapsed = time.Since(start)

	fmt.Println("Hashmap elapsed find: ", elapsed)

}
