package src

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestSplit(t *testing.T) {

	lob := &ListOfBuckets{Buckets: []Bucket{{Indexes: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Min: 1,
		Max: 10,
	},
		{Indexes: []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			Min: 11,
			Max: 20,
		}},
		Height: 5}

	lob.Balance(1, 10)

	fmt.Println(lob.Buckets, "at A")
}

func TestAdd(t *testing.T) {

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = 1000

	fmt.Println("there yet?")
	fmt.Println(splitList)

	var bla []int

	n := 1_000

	for i := n; i >= 0; i-- {

		bla = append(bla, rand.Intn(10_000_000))

	}

	start := time.Now()

	for i := n; i >= 0; i-- {

		splitList.Add(bla[i])

	}

	elapsed := time.Since(start)

	fmt.Printf("Time to insert %d : %d seconds \n", n, elapsed/1_000_000_000)

	summ := 0

	for i := 0; i < splitList.CurrentHeight; i++ {

		for j := range splitList.ListOfBucketLists[i].Buckets {

			summ += len(splitList.ListOfBucketLists[i].Buckets[j].Indexes)

		}

	}

	if summ != n {

		t.Errorf("Inserted elements: %d, found ones: %d", n, summ)

	}

}

func TestFind(t *testing.T) {

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = 1000

	n := 1_000

	for i := 0; i < n; i++ {

		splitList.Add(i)

	}

	start := time.Now()

	for i := 0; i < n; i++ {

		if splitList.Find(i) != true {

			t.Errorf("Could not find %d", i)

		}

	}

	elapsed := time.Since(start)

	fmt.Printf("Time to find %d : %d seconds \n", n, elapsed/1_000_000_000)

}

func TestInsort(t *testing.T) {

	var nums []int
	var bla []int

	n := 100_000

	for i := 1; i < n; i++ {

		nums = append(nums, rand.Intn(20_000_000)) //append(nums, rand.Intn(20_000_000))
	}

	for i := 1; i < n; i++ {

		bla = insortInt(bla, nums[i-1])

	}

	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })

	for i := 1; i < n; i++ {

		if nums[i-1] != bla[i-1] {

			t.Errorf("bla ain't working at %d == %d", nums[i], bla[i])

		}

	}

}

func BenchmarkAdd(b *testing.B) {

	//tlist := NewTeleportList()

	var bla []int

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = 2000

	n := 10_000_000

	for i := 1; i < n; i++ {

		bla = append(bla, rand.Intn(20_000_000))
	}

	b.ResetTimer()

	//var a []int

	for i := 1; i < b.N; i++ {

		splitList.Add(bla[i])

	}

}
