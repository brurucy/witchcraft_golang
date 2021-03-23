package src

import (
	"math"
	"math/rand"
	"sort"
)

const (
	capBucket = 2000
)

type BucketNew struct {
	Indexes []int
	Min     int
	Max     int
}

type ListOfBuckets struct {
	Buckets []*BucketNew
	Height  int
}

type SplitList struct {
	ListOfBucketLists []*ListOfBuckets
	CurrentHeight     int
	Length            int
	Load              int
}

func (l *ListOfBuckets) Balance(idx, load int) {
	candidate := l.Buckets[idx]
	halfLoad := load / 2
	newIndexes := make([]int, halfLoad, load)

	for i := halfLoad - 1; i >= 0; i-- {
		tmpIndexes := candidate.Indexes[len(candidate.Indexes)-1]
		candidate.Indexes = candidate.Indexes[:len(candidate.Indexes)-1]
		newIndexes[i] = tmpIndexes
	}

	newBucket := &BucketNew{
		Indexes: newIndexes,
		Min:     newIndexes[0],
		Max:     candidate.Max,
	}

	candidate.Max = candidate.Indexes[len(candidate.Indexes)-1]

	idxB := sort.Search(len(l.Buckets), func(i int) bool {
		return l.Buckets[i].Max >= newBucket.Max
	})

	l.Buckets = append(l.Buckets, &BucketNew{})
	copy(l.Buckets[idxB+1:], l.Buckets[idxB:])
	l.Buckets[idxB] = newBucket
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

func getRandomHeight() int {
	return int(math.Abs(math.Log2(rand.Float64())))
}

func (s *SplitList) Add(key int) {
	height := getRandomHeight()
	heightDiff := height - s.CurrentHeight

	if height > s.CurrentHeight {
		s.CurrentHeight = height
	}

	for heightDiff > 0 {
		newListOfBuckets := &ListOfBuckets{
			Buckets: []*BucketNew{{
				Max:     math.MinInt64,
				Min:     math.MaxInt64,
				Indexes: make([]int, 0, capBucket),
			}},
			Height: 0,
		}

		s.ListOfBucketLists = append(s.ListOfBucketLists, newListOfBuckets)
		heightDiff--
	}

	// take the last list of buckets
	buckets := s.ListOfBucketLists[height].Buckets

	// bound index to insert/search
	idxB := sort.Search(len(buckets), func(i int) bool {
		return buckets[i].Max >= key
	})

	if idxB == len(buckets) {
		idxB = len(buckets) - 1
	}

	// get the last candidate bucket
	candidate := buckets[idxB]
	idxI := sort.Search(len(candidate.Indexes), func(i int) bool {
		return candidate.Indexes[i] >= key
	})

	candidate.Indexes = append(candidate.Indexes, -1)
	copy(candidate.Indexes[idxI+1:], candidate.Indexes[idxI:])

	// key insertion
	candidate.Indexes[idxI] = key

	// update the boundaries
	candidate.Max = candidate.Indexes[len(candidate.Indexes)-1]
	candidate.Min = candidate.Indexes[0]

	// if reached the load of indexes, re-balance
	if len(candidate.Indexes) == s.Load-1 {
		s.ListOfBucketLists[height].Balance(idxB, s.Load)
	}

	s.Length++
}

func (s *SplitList) Find(key int) bool {
	return s.Lookup(key, nil)
}

func (s *SplitList) Delete(key int) bool {
	return s.Lookup(key, func(idxI, idxB int, buckets []*BucketNew) {
		indexes := buckets[idxB].Indexes

		if len(indexes) == 1 {
			buckets = buckets[:idxB+copy(buckets[idxB:], buckets[idxB+1:])]
		} else {
			indexes = indexes[:idxI+copy(indexes[idxI:], indexes[idxI+1:])]
			buckets[idxB].Max = indexes[len(indexes)-1]
			buckets[idxB].Min = indexes[0]
		}

		s.Length--
	})
}

func (s *SplitList) Lookup(key int, f func(int, int, []*BucketNew)) bool {
	for _, list := range s.ListOfBucketLists {
		listBuckets := list.Buckets

		if len(listBuckets) == 0 {
			continue
		}

		idxB := sort.Search(len(listBuckets), func(i int) bool {
			return listBuckets[i].Max >= key
		})

		if key <= listBuckets[len(listBuckets)-1].Max &&
			key >= listBuckets[0].Min && listBuckets[idxB].Min <= key {
			indexes := listBuckets[idxB].Indexes
			idxI := sort.Search(len(indexes), func(i int) bool {
				return indexes[i] >= key
			})

			if idxI < len(indexes) && indexes[idxI] == key {
				if f != nil {
					f(idxI, idxB, listBuckets)
				}

				return true
			}
		}
	}

	return false
}
