package src

import (
	"math"
	"math/rand"
	"sort"
)

type Bucket struct {
	Indexes []int
	Min     int
	Max     int
}

type ListOfBuckets struct {
	Buckets []Bucket
	Height  int
}

type SplitList struct {
	ListOfBucketLists []*ListOfBuckets
	CurrentHeight     int
	Length            int
	Load              int
}

func (l *ListOfBuckets) Balance(i int, load int) {

	B := Bucket{}

	candidate_bucket := &l.Buckets[i]

	half_load := load / 2

	new_bucket_indexes := make([]int, half_load, load)

	for i := half_load - 1; i >= 0; i-- {

		temp := candidate_bucket.Indexes[len((*candidate_bucket).Indexes)-1]

		(*candidate_bucket).Indexes = (*candidate_bucket).Indexes[0 : len((*candidate_bucket).Indexes)-1]

		new_bucket_indexes[i] = temp

	}

	B.Indexes = new_bucket_indexes
	B.Min = B.Indexes[0]
	B.Max = (*candidate_bucket).Max
	(*candidate_bucket).Max = (*candidate_bucket).Indexes[len((*candidate_bucket).Indexes)-1]

	idx := sort.Search(len(l.Buckets), func(i int) bool { return l.Buckets[i].Max >= B.Max })

	l.Buckets = append(l.Buckets, Bucket{})
	copy(l.Buckets[idx+1:], l.Buckets[idx:])
	l.Buckets[idx] = B

}

func NewSplitList(load int) SplitList {

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = load
	splitList.Length = 0

	return splitList

}

func (s *SplitList) Add(key int) {
	candidateHeight := int(math.Abs(math.Log2(rand.Float64())))

	if candidateHeight > s.CurrentHeight {

		heightDiff := candidateHeight - s.CurrentHeight

		for i := 0; i < heightDiff; i++ {

			newListOfBuckets := ListOfBuckets{}
			newBucket := Bucket{}
			newBucket.Max = math.MinInt64
			newBucket.Min = math.MaxInt64
			newBucket.Indexes = make([]int, 0, 2000)
			newListOfBuckets.Buckets = append(newListOfBuckets.Buckets, newBucket)

			s.ListOfBucketLists = append(s.ListOfBucketLists, &newListOfBuckets)

		}

		s.CurrentHeight = candidateHeight

	}

	candidateHeightListOfBucketsBuckets := &s.ListOfBucketLists[candidateHeight].Buckets

	BucketsBucketsLen := len(*candidateHeightListOfBucketsBuckets)

	candidateBucketIndex := sort.Search(BucketsBucketsLen, func(i int) bool { return (*candidateHeightListOfBucketsBuckets)[i].Max >= key })

	if candidateBucketIndex == len(*candidateHeightListOfBucketsBuckets) {

		candidateBucket := &s.ListOfBucketLists[candidateHeight].Buckets[BucketsBucketsLen-1]

		candidateBucketBucketIndex := sort.Search(len((*candidateBucket).Indexes), func(i int) bool { return (*candidateBucket).Indexes[i] >= key })

		(*candidateBucket).Indexes = append((*candidateBucket).Indexes, -1)
		copy((*candidateBucket).Indexes[candidateBucketBucketIndex+1:], (*candidateBucket).Indexes[candidateBucketBucketIndex:])
		(*candidateBucket).Indexes[candidateBucketBucketIndex] = key

		(*candidateBucket).Max = (*candidateBucket).Indexes[len((*candidateBucket).Indexes)-1]
		(*candidateBucket).Min = (*candidateBucket).Indexes[0]

		if len((*candidateBucket).Indexes) == s.Load-1 {

			s.ListOfBucketLists[candidateHeight].Balance(BucketsBucketsLen-1, s.Load)

		}

	} else {

		candidateBucket := &s.ListOfBucketLists[candidateHeight].Buckets[candidateBucketIndex]
		candidateBucketBucketIndex := sort.Search(len((*candidateBucket).Indexes), func(i int) bool { return (*candidateBucket).Indexes[i] >= key })

		(*candidateBucket).Indexes = append((*candidateBucket).Indexes, -1)
		copy((*candidateBucket).Indexes[candidateBucketBucketIndex+1:], (*candidateBucket).Indexes[candidateBucketBucketIndex:])
		(*candidateBucket).Indexes[candidateBucketBucketIndex] = key

		(*candidateBucket).Max = (*candidateBucket).Indexes[len((*candidateBucket).Indexes)-1]
		(*candidateBucket).Min = (*candidateBucket).Indexes[0]

		if len((*candidateBucket).Indexes) == s.Load-1 {

			s.ListOfBucketLists[candidateHeight].Balance(candidateBucketIndex, s.Load)

		}

	}

	s.Length += 1

}

func (s *SplitList) Find(key int) bool {

	//len_lob := len(s.ListOfBucketLists)

	for i := range s.ListOfBucketLists {

		lb := &s.ListOfBucketLists[i].Buckets
		len_lb := len(*lb)

		if len_lb != 0 {

			if !(key > (*lb)[len_lb-1].Max || key < (*lb)[0].Min) {

				i = sort.Search(len_lb, func(i int) bool { return (*lb)[i].Max >= key })

				if (i != len_lb) && (!((*lb)[i].Min > key)) {

					if BinarySearchInt((*lb)[i].Indexes, key) {

						return true

					}

				}

			}

		}

	}

	return false

}

func (s *SplitList) Delete(key int) bool {

	for i := range s.ListOfBucketLists {

		lb := &s.ListOfBucketLists[i].Buckets
		len_lb := len(*lb)

		if len_lb != 0 {

			if !(key > (*lb)[len_lb-1].Max || key < (*lb)[0].Min) {

				i = sort.Search(len_lb, func(i int) bool { return (*lb)[i].Max >= key })

				if (i != len_lb) && (!((*lb)[i].Min > key)) {

					len_idx := len((*lb)[i].Indexes)

					j := sort.Search(len_idx, func(k int) bool { return (*lb)[i].Indexes[k] >= key })

					if j < len_idx && (*lb)[i].Indexes[j] == key {

						if len_idx == 1 {

							*lb = (*lb)[:i+copy((*lb)[i:], (*lb)[i+1:])]

						} else {

							(*lb)[i].Indexes = (*lb)[i].Indexes[:j+copy((*lb)[i].Indexes[j:], (*lb)[i].Indexes[j+1:])]
							(*lb)[i].Max = (*lb)[i].Indexes[len((*lb)[i].Indexes)-1]
							(*lb)[i].Min = (*lb)[i].Indexes[0]

						}

						s.Length -= 1

						return true

					}

				}

			}

		}

	}
	return false
}
