package src

import (
	"sort"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

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

	//idx := sort.Search(len(a), func(i int) bool { return a[i] >= value })
	idx := sort.Search(len(a), func(i int) bool { return a[i] >= value })
	a = append(a, value)
	copy(a[idx+1:], a[idx:])
	a[idx] = value

	return a

	//return insertInt(a, idx, value)

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
