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

func TestCopy(t *testing.T) {

	bla := make([]int, 10)

	for i := 0; i < 10; i++ {

		bla[i] = i

	}

	ble := make([]int, 10)

	fmt.Println("bla", bla)

	copy_result := copy(ble[0:5], bla[5:len(bla)])

	fmt.Println("bla", bla)

	fmt.Println("ble", ble)

	fmt.Println("Copy Result", copy_result)

}

func TestAdd(t *testing.T) {

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = 1000

	fmt.Println("there yet?")
	fmt.Println(splitList)

	var bla []int

	n := 10000000

	for i := n; i >= 0; i-- {

		bla = append(bla, rand.Intn(100_000_000))

		//bla = append(bla, i)

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

	/*

		for i := 0; i < splitList.CurrentHeight; i++ {

			for j := range splitList.ListOfBucketLists[i].Buckets {

				fmt.Println("Length: ", len(splitList.ListOfBucketLists[i].Buckets[j].Indexes))

			}

			fmt.Println(splitList.ListOfBucketLists[i])
		}
	*/

}

type Element int

// Implement the interface used in skiplist
func (e Element) ExtractKey() float64 {
	return float64(e)
}
func (e Element) String() string {
	return fmt.Sprintf("%03d", e)
}

/*

func TestSkipListAdd(t *testing.T) {

	var bla []int

	skipList := New()

	n := 10_000_000

	for i := n; i >= 0; i-- {

		bla = append(bla, rand.Intn(100_000_000))

	}

	start := time.Now()

	for i := n; i >= 0; i-- {

		skipList.Insert(Element(bla[i]))

	}

	elapsed := time.Since(start)

	fmt.Printf("Time to insert %d : %d seconds \n", n, elapsed/1_000_000_000)

}
*/
func TestFind(t *testing.T) {

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = 1000

	n := 10_000_000

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

/*
func BenchmarkSkipListAdd(b *testing.B) {

	//tlist := NewTeleportList()

	var bla []int

	//splitList := SplitList{}
	//splitList.CurrentHeight = -1
	//splitList.Load = 2000

	skipList := New()

	n := 10_000_000

	for i := 1; i < n; i++ {

		bla = append(bla, rand.Intn(20_000_000))
	}

	b.ResetTimer()

	//var a []int

	for i := 1; i < b.N; i++ {

		skipList.Insert(Element(bla[i]))

	}

}
*/
func BenchmarkSplitListAdd(b *testing.B) {

	var bla []int

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = 2000

	n := 20_000_000

	for i := 1; i < n; i++ {

		bla = append(bla, rand.Intn(100_000_000))
	}

	for i := 1; i < (n / 2); i++ {

		splitList.Add(bla[i])

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		splitList.Add(bla[i])

	}

}

func BenchmarkSkipListAdd(b *testing.B) {

	var bla []int

	//splitList := SplitList{}
	//splitList.CurrentHeight = -1
	//splitList.Load = 2000

	skipList := New()

	n := 20_000_000

	for i := 1; i < n; i++ {

		bla = append(bla, rand.Intn(100_000_000))
	}

	for i := 1; i < (n / 2); i++ {

		skipList.Insert(Element(bla[i]))

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		skipList.Insert(Element(bla[i]))

	}

}

/*
func BenchmarkSkipListFind(b *testing.B) {

	var bla []int

	skipList := New()

	n := 10_000_000

	for i := 1; i < n; i++ {

		bla = append(bla, rand.Intn(20_000_000))
	}

	for i := 1; i < n; i++ {

		skipList.Insert(Element(bla[i-1]))

	}

	b.ResetTimer()

	for i := 1; i < n; i++ {

		skipList.Find(Element(bla[i-1]))

	}

}
*/
func BenchmarkSplitListFind(b *testing.B) {

	var bla []int

	splitList := SplitList{}
	splitList.CurrentHeight = -1
	splitList.Load = 2000

	//skipList := New()

	n := 20_000_000

	for i := 1; i < n; i++ {

		bla = append(bla, rand.Intn(100_000_000))
	}

	//var a []int

	for i := 1; i < n; i++ {

		splitList.Add(bla[i-1])

	}

	b.ResetTimer()

	for i := 1; i < b.N; i++ {

		splitList.Find(bla[i-1])

	}

}
