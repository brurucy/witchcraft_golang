package src

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

const (
	capBucket = 2000
)

// Item represents a single object in the tree.
type Item interface {
	// Less tests whether the current item is less than the given argument.
	//
	// This must provide a strict weak ordering.
	// If !a.Less(b) && !b.Less(a), we treat this to mean a == b (i.e. we can only
	// hold one of either a or b in the tree).
	Less(than Item) bool
}

type Bucket struct {
	Indexes []Item
	Min     Item
	Max     Item
}

// find returns the index where the given item should be inserted into this
// list.  'found' is true if the item already exists in the list at the given
// index.
func (b Bucket) find(item Item) (index int, found bool) {
	i := sort.Search(len(b.Indexes), func(i int) bool {
		return item.Less(b.Indexes[i])
	})
	if i > 0 && !b.Indexes[i-1].Less(item) {
		return i - 1, true
	}
	return i, false
}

type ListOfBuckets struct {
	Buckets []*Bucket
	Height  int
	ready   bool
}

func (l ListOfBuckets) String() string {
	s := ""

	for _, list := range l.Buckets {
		s += fmt.Sprintf("%v\n", list)
	}

	return s
}

type SplitList struct {
	ListOfBucketLists []*ListOfBuckets
	CurrentHeight     int
	Length            int
	Load              int
}

func NewSplitList(load int) *SplitList {
	splitList := &SplitList{
		ListOfBucketLists: make([]*ListOfBuckets, 0),
		CurrentHeight:     -1,
		Length:            0,
		Load:              load,
	}

	return splitList
}

func (l *ListOfBuckets) Balance(idx, load int) {
	candidate := l.Buckets[idx]
	halfLoad := load / 2
	newIndexes := make([]Item, halfLoad, load)

	for i := halfLoad - 1; i >= 0; i-- {
		tmpIndexes := candidate.Indexes[len(candidate.Indexes)-1]
		candidate.Indexes = candidate.Indexes[:len(candidate.Indexes)-1]
		newIndexes[i] = tmpIndexes
	}

	newBucket := &Bucket{
		Indexes: newIndexes,
		Min:     newIndexes[0],
		Max:     candidate.Max,
	}

	candidate.Max = candidate.Indexes[len(candidate.Indexes)-1]

	idxB := sort.Search(len(l.Buckets), func(i int) bool {
		return l.Buckets[i].Max.Less(newBucket.Max)
	})

	l.Buckets = append(l.Buckets, &Bucket{})
	copy(l.Buckets[idxB+1:], l.Buckets[idxB:])
	l.Buckets[idxB] = newBucket
}

func getRandomHeight() int {
	return int(math.Abs(math.Log2(rand.Float64())))
}

func (s *SplitList) Add(item Item) {
	if item == nil {
		panic("nil item being added to BTree")
	}

	height := getRandomHeight()
	heightDiff := height - s.CurrentHeight

	if height > s.CurrentHeight {
		s.CurrentHeight = height
	}

	for heightDiff > 0 {
		newListOfBuckets := &ListOfBuckets{
			Buckets: []*Bucket{{
				Max:     nil,
				Min:     nil,
				Indexes: make([]Item, 0, capBucket),
			}},
			Height: height + heightDiff,
			ready:  false,
		}

		s.ListOfBucketLists = append(s.ListOfBucketLists, newListOfBuckets)
		heightDiff--
	}

	// take the last list of buckets
	buckets := s.ListOfBucketLists[height].Buckets

	if !s.ListOfBucketLists[height].ready {
		s.ListOfBucketLists[height].Buckets[0].Min = item
		s.ListOfBucketLists[height].Buckets[0].Max = item
		s.ListOfBucketLists[height].ready = true
	}

	// bound index to insert/search
	idxB := sort.Search(len(buckets), func(i int) bool {
		return !item.Less(buckets[i].Max)
	})

	if idxB == len(buckets) {
		idxB = len(buckets) - 1
	}

	// get the last candidate bucket
	candidate := buckets[idxB]

	idxI, _ := candidate.find(item)

	candidate.Indexes = append(candidate.Indexes, nil)
	copy(candidate.Indexes[idxI+1:], candidate.Indexes[idxI:])

	// key insertion
	candidate.Indexes[idxI] = item

	// update the boundaries
	candidate.Max = candidate.Indexes[len(candidate.Indexes)-1]
	candidate.Min = candidate.Indexes[0]

	// if reached the load of indexes, re-balance
	if len(candidate.Indexes) == s.Load-1 {
		s.ListOfBucketLists[height].Balance(idxB, s.Load)
	}

	s.Length++
}

func (s *SplitList) Find(item Item) bool {
	return s.Lookup(item, nil)
}

//func (s *SplitList) Delete(key int) bool {
//	return s.Lookup(key, func(idxI, idxB int, buckets []*Bucket) {
//		indexes := buckets[idxB].Indexes
//
//		if len(indexes) == 1 {
//			buckets = buckets[:idxB+copy(buckets[idxB:], buckets[idxB+1:])]
//		} else {
//			indexes = indexes[:idxI+copy(indexes[idxI:], indexes[idxI+1:])]
//			buckets[idxB].Max = indexes[len(indexes)-1]
//			buckets[idxB].Min = indexes[0]
//		}
//
//		s.Length--
//	})
//}

func (s *SplitList) Lookup(item Item, f func(int, int, []*Bucket)) bool {
	for _, list := range s.ListOfBucketLists {
		listBuckets := list.Buckets

		if len(listBuckets) == 0 {
			continue
		}

		idxB := sort.Search(len(listBuckets), func(i int) bool {
			return !listBuckets[i].Max.Less(item)
		})

		if item.Less(listBuckets[len(listBuckets)-1].Max) {
			idxI, found := listBuckets[idxB].find(item)

			if found {
				if f != nil {
					f(idxI, idxB, listBuckets)
				}

				return true
			}
		}

		//if key <= listBuckets[len(listBuckets)-1].Max &&  key >= listBuckets[0].Min && listBuckets[idxB].Min <= key {
		//	idxI, found := listBuckets[idxB].find(item)
		//
		//	if found {
		//		if f != nil {
		//			f(idxI, idxB, listBuckets)
		//		}
		//
		//		return true
		//	}
		//}
	}

	return false
}

// Theta log(n)
func (s *SplitList) GetMin() int {
	runningMinimum := math.MaxInt64

	//for _, list := range s.ListOfBucketLists {
	//	if list.Buckets[0].Min < runningMinimum {
	//		runningMinimum = list.Buckets[0].Min
	//	}
	//}

	return runningMinimum

}

// Theta log(n)
func (s *SplitList) GetMax() int {
	runningMaximum := math.MinInt64

	//for _, list := range s.ListOfBucketLists {
	//	if list.Buckets[len(list.Buckets)-1].Max > runningMaximum {
	//		runningMaximum = list.Buckets[len(list.Buckets)-1].Max
	//	}
	//}

	return runningMaximum
}

func (s *SplitList) PopMin() int {
	if s.Length == 0 {
		return -1
	}

	min := s.GetMin()

	//for _, list := range s.ListOfBucketLists {
	//	if list.Buckets[0].Min == min {
	//		list.Buckets[0].Indexes = list.Buckets[0].Indexes[1:]
	//
	//		if len(list.Buckets[0].Indexes) == 0 {
	//			if len(list.Buckets) > 1 {
	//				list.Buckets = list.Buckets[1:]
	//			} else {
	//				list.Buckets[0].Min = math.MaxInt64
	//				list.Buckets[0].Max = math.MinInt64
	//			}
	//		} else {
	//			list.Buckets[0].Min = list.Buckets[0].Indexes[0]
	//			list.Buckets[0].Max = list.Buckets[0].Indexes[len(list.Buckets[0].Indexes)-1]
	//		}
	//
	//		s.Length--
	//
	//		return min
	//	}
	//}

	return min
}

func (s *SplitList) PopMax() int {
	if s.Length == 0 {
		return -1
	}

	max := s.GetMax()
	//
	//for _, list := range s.ListOfBucketLists {
	//	lastBucketIndex := len(list.Buckets) - 1
	//
	//	if list.Buckets[lastBucketIndex].Max == max {
	//		list.Buckets[lastBucketIndex].Indexes = list.Buckets[lastBucketIndex].Indexes[:len(list.Buckets[lastBucketIndex].Indexes)-1]
	//
	//		if len(list.Buckets[lastBucketIndex].Indexes) == 0 {
	//			if len(list.Buckets) > 1 {
	//				list.Buckets = list.Buckets[:lastBucketIndex]
	//			} else {
	//				list.Buckets[lastBucketIndex].Min = math.MaxInt64
	//				list.Buckets[lastBucketIndex].Max = math.MinInt64
	//			}
	//		} else {
	//			list.Buckets[lastBucketIndex].Min = list.Buckets[lastBucketIndex].Indexes[0]
	//			list.Buckets[lastBucketIndex].Max = list.Buckets[lastBucketIndex].Indexes[len(list.Buckets[lastBucketIndex].Indexes)-1]
	//		}
	//
	//		s.Length--
	//
	//		return max
	//
	//	}
	//
	//}

	return max
}
