package src

import (
	"sort"
)

type KeyValue struct {
	key   int
	value string
}

func BinarySearchInt(a []int, value int) bool {

	idx := sort.Search(len(a), func(i int) bool { return (a)[i] >= value })

	if idx < len(a) && a[idx] == value {

		return true

	} else {

		return false

	}

}
