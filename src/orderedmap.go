package src

import (
	"fmt"
	gbtree "github.com/google/btree"
	"math"
	"math/rand"
	"sort"
)

const (
	capBucket = 2000
)

type Bucket struct {
	Indexes []gbtree.Item
	Min     gbtree.Item
	Max     gbtree.Item
}

// find returns the index where the given item should be inserted into this
// list.  'found' is true if the item already exists in the list at the given
// index.
func (b Bucket) find(item gbtree.Item) (index int, found bool) {
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
	CachedFlatItems   flatIndexes
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
	newIndexes := make([]gbtree.Item, halfLoad, load)

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
		return !l.Buckets[i].Max.Less(newBucket.Max)
	})

	l.Buckets = append(l.Buckets, &Bucket{})
	copy(l.Buckets[idxB+1:], l.Buckets[idxB:])
	l.Buckets[idxB] = newBucket
}

func getRandomHeight() int {
	return int(math.Abs(math.Log2(rand.Float64())))
}

func (s *SplitList) Add(item gbtree.Item) {
	if item == nil {
		panic("nil item being added to Split List")
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
				Indexes: make([]gbtree.Item, 0, capBucket),
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
		return !buckets[i].Max.Less(item) //!item.Less(buckets[i].Max)
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

type flatIndexes []gbtree.Item

func (f flatIndexes) Len() int {
	return len(f)
}

func (f flatIndexes) Less(i, j int) bool {
	return f[i].Less(f[j])
}

func (f flatIndexes) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Select finds the kth smallest element across all splitlist levels
func (s *SplitList) Select(kth int) gbtree.Item {
	if s.CachedFlatItems == nil || len(s.CachedFlatItems) != s.Length {
		s.CachedFlatItems = make([]gbtree.Item, s.Length)

		offset := 0
		for _, l := range s.ListOfBucketLists {
			for _, b := range l.Buckets {
				copy(s.CachedFlatItems[offset:], b.Indexes)
				offset += len(b.Indexes)
			}
		}

		sort.Sort(s.CachedFlatItems)
	}

	return s.CachedFlatItems[kth]
}

func (s *SplitList) Select1(kth int) gbtree.Item {
	flat := make([]gbtree.Item, s.Length)

	offset := 0
	for _, l := range s.ListOfBucketLists {
		for _, b := range l.Buckets {
			copy(flat[offset:], b.Indexes)
			offset += len(b.Indexes)
		}
	}

	flatIdx := flatIndexes(flat)

	first, last := 0, flatIdx.Len()-1
	for {
		flatIdx.Swap(first, rand.Intn(last-first+1)+first)
		left := first + 1
		right := last
		for left <= right {
			for left <= last && flatIdx.Less(left, first) {
				left++
			}
			for right >= first && flatIdx.Less(first, right) {
				right--
			}
			if left <= right {
				flatIdx.Swap(left, right)
				left++
				right--
			}
		}
		flatIdx.Swap(first, right)

		if kth == right {
			return flatIdx[right]
		} else if kth < right {
			last = right - 1
		} else {
			first = right + 1
		}
	}
}

// Rank outputs the rank in the sorted union of all lists, that the given value would occupy
func (s *SplitList) Rank(item gbtree.Item) (kth int) {

	for _, list := range s.ListOfBucketLists {

		if !list.ready || len(list.Buckets) == 0 {
			continue
		}

		bucketsWithMaximumLessThanItem := sort.Search(len(list.Buckets), func(i int) bool {
			return !list.Buckets[i].Max.Less(item)
		})

		switch bucketsWithMaximumLessThanItem {
		// If it isn't smaller than anything, continue iterating
		case 0:
			{

				valuesWithMaximumLessThanItemInTheFirstBucket, _ := list.Buckets[0].find(item)

				kth += valuesWithMaximumLessThanItemInTheFirstBucket

			}
		// Else...
		default:
			{
				// Count the lengths of all buckets UP TO the one that supposedly could contain the item that we want
				for i := 0; i < bucketsWithMaximumLessThanItem; i++ {
					kth += len(list.Buckets[i].Indexes)
				}

				if len(list.Buckets) == bucketsWithMaximumLessThanItem {

					continue

				} else {

					// Now, count the values in the bucket that we want that are LESS than the item
					valuesWithMaximumLessThanItemInTheLastBucket, _ := list.Buckets[bucketsWithMaximumLessThanItem].find(item)

					kth += valuesWithMaximumLessThanItemInTheLastBucket

				}

			}
		}

	}

	return kth
}

func (s *SplitList) Find(item gbtree.Item) bool {
	return s.Lookup(item, nil)
}

func (s *SplitList) Delete(item gbtree.Item) bool {
	return s.Lookup(item, func(idxI, idxB int, buckets *[]*Bucket) {
		indexes := &((*buckets)[idxB].Indexes)

		if len(*indexes) == 1 {
			*buckets = (*buckets)[:idxB+copy((*buckets)[idxB:], (*buckets)[idxB+1:])]
		} else {
			*indexes = (*indexes)[:idxI+copy((*indexes)[idxI:], (*indexes)[idxI+1:])]
			(*buckets)[idxB].Max = (*indexes)[len(*indexes)-1]
			(*buckets)[idxB].Min = (*indexes)[0]
		}

		s.Length--
	})
}

func (s *SplitList) LookupReverse(item gbtree.Item) bool {

	for i := s.CurrentHeight; i >= 0; i-- {
		listBuckets := s.ListOfBucketLists[i].Buckets

		if !s.ListOfBucketLists[i].ready || len(listBuckets) == 0 {
			continue
		}

		if item.Less(listBuckets[len(listBuckets)-1].Max) ||
			(!item.Less(listBuckets[len(listBuckets)-1].Max) &&
				!listBuckets[len(listBuckets)-1].Max.Less(item)) {

			idxB := sort.Search(len(listBuckets), func(i int) bool {
				return !listBuckets[i].Max.Less(item)
			})

			_, found := listBuckets[idxB].find(item)

			if found {

				return true
			}
		}
	}

	return false

}

func (s *SplitList) Lookup(item gbtree.Item, f func(int, int, *[]*Bucket)) bool {
	for _, list := range s.ListOfBucketLists {
		listBuckets := list.Buckets

		if !list.ready || len(listBuckets) == 0 {
			continue
		}

		if item.Less(listBuckets[len(listBuckets)-1].Max) ||
			(!item.Less(listBuckets[len(listBuckets)-1].Max) &&
				!listBuckets[len(listBuckets)-1].Max.Less(item)) {

			idxB := sort.Search(len(listBuckets), func(i int) bool {
				return !listBuckets[i].Max.Less(item)
			})

			idxI, found := listBuckets[idxB].find(item)

			if found {
				if f != nil {
					f(idxI, idxB, &list.Buckets)
				}

				return true
			}
		}
	}

	return false
}

func (s *SplitList) GetMin() gbtree.Item {

	runningMinimum := s.ListOfBucketLists[0].Buckets[0].Min

	for _, list := range s.ListOfBucketLists {

		if !list.ready {

			continue

		}

		if runningMinimum == nil {
			runningMinimum = list.Buckets[0].Min
			continue

		}

		if list.Buckets[0].Min.Less(runningMinimum) || (!list.Buckets[0].Min.Less(runningMinimum) && !runningMinimum.Less(list.Buckets[0].Min)) {

			runningMinimum = list.Buckets[0].Min

		}
	}

	return runningMinimum
}

func (s *SplitList) GetMax() gbtree.Item {
	runningMaximum := s.ListOfBucketLists[0].Buckets[len(s.ListOfBucketLists[0].Buckets)-1].Max

	for _, list := range s.ListOfBucketLists {

		if !list.ready {
			continue

		}

		if runningMaximum == nil {
			runningMaximum = list.Buckets[len(list.Buckets)-1].Max
			continue

		}

		if runningMaximum.Less(list.Buckets[len(list.Buckets)-1].Max) || (!list.Buckets[len(list.Buckets)-1].Max.Less(runningMaximum) && !runningMaximum.Less(list.Buckets[len(list.Buckets)-1].Max)) {
			runningMaximum = list.Buckets[len(list.Buckets)-1].Max
		}
	}

	return runningMaximum

}

func (s *SplitList) PopMin() gbtree.Item {
	if s.Length == 0 {
		return nil
	}

	min := s.GetMin()

	for _, list := range s.ListOfBucketLists {
		if !list.ready {
			continue
		}
		if !list.Buckets[0].Min.Less(min) && !min.Less(list.Buckets[0].Min) {
			list.Buckets[0].Indexes = list.Buckets[0].Indexes[1:]

			if len(list.Buckets[0].Indexes) == 0 {
				if len(list.Buckets) > 1 {
					list.Buckets = list.Buckets[1:]
				} else {
					list.Buckets[0].Min = nil
					list.Buckets[0].Max = nil
					list.ready = false
				}
			} else {
				list.Buckets[0].Min = list.Buckets[0].Indexes[0]
				list.Buckets[0].Max = list.Buckets[0].Indexes[len(list.Buckets[0].Indexes)-1]
			}

			s.Length--

			return min
		}
	}

	return min
}

func (s *SplitList) PopMax() gbtree.Item {
	if s.Length == 0 {
		return nil
	}

	max := s.GetMax()

	for _, list := range s.ListOfBucketLists {
		lastBucketIndex := len(list.Buckets) - 1

		if !list.ready {
			continue
		}
		if !list.Buckets[lastBucketIndex].Max.Less(max) && !max.Less(list.Buckets[lastBucketIndex].Max) {

			list.Buckets[lastBucketIndex].Indexes = list.Buckets[lastBucketIndex].Indexes[:len(list.Buckets[lastBucketIndex].Indexes)-1]

			if len(list.Buckets[lastBucketIndex].Indexes) == 0 {
				if len(list.Buckets) > 1 {
					list.Buckets = list.Buckets[:lastBucketIndex]
				} else {
					list.Buckets[lastBucketIndex].Min = nil
					list.Buckets[lastBucketIndex].Max = nil
					list.ready = false
				}
			} else {
				list.Buckets[lastBucketIndex].Min = list.Buckets[lastBucketIndex].Indexes[0]
				list.Buckets[lastBucketIndex].Max = list.Buckets[lastBucketIndex].Indexes[len(list.Buckets[lastBucketIndex].Indexes)-1]
			}

			s.Length--

			return max

		}

	}

	return max
}
