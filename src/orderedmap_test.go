package src

import (
	"fmt"
	mauriceSkipList "github.com/MauriceGit/skiplist"
	//chavezSkipList "github.com/mtchavez/skiplist"
	//ryszardSkipList "github.com/ryszard/goskiplist/skiplist"
	btree "github.com/google/btree"
	seanSkipList "github.com/sean-public/fast-skiplist"
	"math/rand"
	"testing"
	"time"
)

func TestBalance(t *testing.T) {

	lob := &ListOfBuckets{Buckets: []*Bucket{{Indexes: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
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

func TestAddFindTiny(t *testing.T) {
	splitList := NewSplitList(1024)

	for i := 0; i < 10; i++ {
		splitList.Add(i)
	}

	if !splitList.Find(55) || splitList.Find(909) {
		t.Fail()
	}
}

func TestAddFindTinyInBetween(t *testing.T) {
	splitList := NewSplitList(1024)
	n := 1_000_000

	for i := 0; i < n; i++ {
		splitList.Add(i)
	}

	for i := 0; i < n; i++ {
		if !splitList.Find(i) {
			t.Fail()
		}
	}
}

func TestAddFindSmall(t *testing.T) {
	splitList := NewSplitList(3)

	var bla []int
	n := 10

	for i := 0; i < n; i++ {
		bla = append(bla, rand.Intn(100))
	}

	fmt.Println("To be inserted: ", bla)

	for i := 0; i < n; i++ {
		splitList.Add(bla[i])
		fmt.Println("Attempting to find", i, ", success: ", splitList.Find(i))
	}

	fmt.Println("Testing if it finds elements that are NOT in")

	for i := 19; i > 9; i-- {
		fmt.Println("Attempting to find", i, ", success: ", splitList.Find(i))
	}

	for i := 0; i < n; i++ {
		fmt.Println("Before removing")
		for j := range splitList.ListOfBucketLists {
			fmt.Println(*splitList.ListOfBucketLists[j])
		}
		fmt.Println("Attempting to remove ", bla[i])
		splitList.Delete(bla[i])
		for j := range splitList.ListOfBucketLists {
			fmt.Println(*splitList.ListOfBucketLists[j])
		}
	}
}

func TestAddFindSequential(t *testing.T) {
	splitList := NewSplitList(1024)
	fmt.Println(splitList)

	var bla []int
	n := 10_000_000

	for i := n; i >= 0; i-- {
		bla = append(bla, i) //rand.Intn(300))
	}

	start := time.Now()
	for i := n; i >= 0; i-- {
		splitList.Add(bla[i])
	}

	elapsed := time.Since(start)
	fmt.Printf("Time to insert %d : %d seconds \n", n, elapsed/1_000_000_000)
	summ := 0
	for i := 0; i <= splitList.CurrentHeight; i++ {
		for j := range splitList.ListOfBucketLists[i].Buckets {
			summ += len(splitList.ListOfBucketLists[i].Buckets[j].Indexes)
		}
	}

	if (summ - 1) != n {
		t.Errorf("Inserted elements: %d, found ones: %d", n, summ-1)
	}

	start = time.Now()

	for i := n; i >= 0; i-- {
		if splitList.Find(bla[i]) != true {
			t.Errorf("Could not find %d", i)
		}
	}

	elapsed = time.Since(start)
	fmt.Printf("Time to find %d : %d seconds \n", n, elapsed/1_000_000_000)

}

func TestAddFindRandom(t *testing.T) {
	splitList := NewSplitList(1024)
	fmt.Println(splitList)
	var bla []int

	n := 10_000_000
	fmt.Println("about to append")
	for i := n; i >= 0; i-- {
		bla = append(bla, rand.Intn(100_000_000)) //rand.Intn(300))
	}

	fmt.Println("random done")
	start := time.Now()
	for i := n; i >= 0; i-- {
		splitList.Add(bla[i])
	}

	elapsed := time.Since(start)
	fmt.Printf("Time to insert %d : %d seconds \n", n, elapsed/1_000_000_000)
	summ := 0

	for i := 0; i <= splitList.CurrentHeight; i++ {
		for j := range splitList.ListOfBucketLists[i].Buckets {
			summ += len(splitList.ListOfBucketLists[i].Buckets[j].Indexes)
		}
	}

	if (summ - 1) != n {
		t.Errorf("Inserted elements: %d, found ones: %d", n, summ-1)
	}

	start = time.Now()

	for i := n; i >= 0; i-- {
		if splitList.Find(bla[i]) != true {
			t.Errorf("Could not find %d", bla[i])
		}
	}

	elapsed = time.Since(start)

	fmt.Printf("Time to find %d : %d seconds \n", n, elapsed/1_000_000_000)
	fmt.Println("Elements before deleting: ", splitList.Length)
	start = time.Now()
	for i := n; i >= 0; i-- {
		if splitList.Delete(bla[i]) != true {
			t.Errorf("Could not delete %d", bla[i])
		}
	}

	elapsed = time.Since(start)

	fmt.Printf("Time to delete %d : %d seconds \n", n, elapsed/1_000_000_000)

	fmt.Println("Elements after deleting: ", splitList.Length)

	for j := range splitList.ListOfBucketLists {

		fmt.Println(*splitList.ListOfBucketLists[j])

	}

}

func BenchmarkSplitListRandAdd(b *testing.B) {
	var bla []int
	splitList := NewSplitList(1024)
	n := 20_000_000
	for i := 0; i < n; i++ {
		bla = append(bla, rand.Intn(100_000_000))
	}

	for i := 0; i < (n / 2); i++ {
		splitList.Add(bla[i])
	}

	b.ResetTimer()
	for i := (n / 2); i < b.N; i++ {
		splitList.Add(bla[i])
	}
}

func BenchmarkSplitListIncAdd(b *testing.B) {
	var bla []int
	splitList := NewSplitList(1024)
	n := 20_000_000
	for i := 0; i < n; i++ {
		bla = append(bla, i)
	}

	for i := 0; i < (n / 2); i++ {
		splitList.Add(bla[i])
	}

	b.ResetTimer()
	for i := (n / 2); i < b.N; i++ {
		splitList.Add(bla[i])
	}
}

func BenchmarkSplitListRandDelete(b *testing.B) {
	var bla []int
	splitList := NewSplitList(1024)
	n := 20_000_000
	for i := 0; i < n; i++ {
		bla = append(bla, rand.Intn(100_000_000))
	}

	for i := 0; i < n; i++ {
		splitList.Add(bla[i])
	}

	b.ResetTimer()
	for i := (n / 2); i < b.N; i++ {
		splitList.Delete(bla[i])
	}
}

func BenchmarkSplitListIncDelete(b *testing.B) {
	var bla []int
	splitList := NewSplitList(1024)
	n := 20_000_000
	for i := 0; i < n; i++ {
		bla = append(bla, i)
	}

	for i := 0; i < n; i++ {
		splitList.Add(bla[i])
	}

	b.ResetTimer()
	for i := (n / 2); i < b.N; i++ {
		splitList.Delete(bla[i])
	}
}

func BenchmarkSplitListRandFind(b *testing.B) {
	var bla []int
	splitList := NewSplitList(1024)
	//skipList := New()
	n := 20_000_000
	for i := 0; i < n; i++ {
		bla = append(bla, rand.Intn(100_000_000))
	}

	//var a []int
	for i := 0; i < n/2; i++ {
		splitList.Add(bla[i])
	}

	b.ResetTimer()

	for i := n / 2; i < b.N; i++ {
		splitList.Find(bla[i])
	}
}

func BenchmarkSplitListIncFind(b *testing.B) {

	var bla []int

	splitList := NewSplitList(1024)

	//skipList := New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, i)
	}

	//var a []int

	for i := 0; i < n/2; i++ {

		splitList.Add(bla[i])

	}

	b.ResetTimer()

	for i := n / 2; i < b.N; i++ {

		splitList.Find(bla[i])

	}

}

type Element int

// Implement the interface used in skiplist
func (e Element) ExtractKey() float64 {
	return float64(e)
}
func (e Element) String() string {
	return fmt.Sprintf("%03d", e)
}

func BenchmarkMauriceSkipListRandAdd(b *testing.B) {

	var bla []int

	//splitList := SplitList{}
	//splitList.CurrentHeight = -1
	//splitList.Load = 2000

	skipList := mauriceSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, rand.Intn(100_000_000))
	}

	for i := 0; i < (n / 2); i++ {

		skipList.Insert(Element(bla[i]))

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		skipList.Insert(Element(bla[i]))

	}

}

func BenchmarkMauriceSkipListIncAdd(b *testing.B) {

	var bla []int

	//splitList := SplitList{}
	//splitList.CurrentHeight = -1
	//splitList.Load = 2000

	skipList := mauriceSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, i)
	}

	for i := 0; i < (n / 2); i++ {

		skipList.Insert(Element(bla[i]))

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		skipList.Insert(Element(bla[i]))

	}

}

func BenchmarkMauriceSkipListRandFind(b *testing.B) {

	var bla []int

	skipList := mauriceSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, rand.Intn(100_000_000))
	}

	for i := 0; i < n/2; i++ {

		skipList.Insert(Element(bla[i]))

	}

	b.ResetTimer()

	for i := n / 2; i < n; i++ {

		skipList.Find(Element(bla[i]))

	}

}

func BenchmarkMauriceSkipListIncFind(b *testing.B) {

	var bla []int

	skipList := mauriceSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, i)
	}

	for i := 0; i < n/2; i++ {

		skipList.Insert(Element(bla[i]))

	}

	b.ResetTimer()

	for i := n / 2; i < n; i++ {

		skipList.Find(Element(bla[i]))

	}

}

func BenchmarkSeanSkipListRandAdd(b *testing.B) {

	var bla []float64

	//splitList := SplitList{}
	//splitList.CurrentHeight = -1
	//splitList.Load = 2000

	skipList := seanSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, rand.Float64()*100_000_000)
	}

	for i := 0; i < (n / 2); i++ {

		skipList.Set(bla[i], "")

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		skipList.Set(bla[i], "")

	}

}

func BenchmarkSeanSkipListIncAdd(b *testing.B) {

	var bla []float64

	skipList := seanSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, float64(i))
	}

	for i := 0; i < (n / 2); i++ {

		skipList.Set(bla[i], "")

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		skipList.Set(bla[i], "")

	}

}

func BenchmarkSeanSkipListRandFind(b *testing.B) {

	var bla []float64

	skipList := seanSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, rand.Float64()*100_000_000)
	}

	for i := 0; i < n/2; i++ {

		skipList.Set(bla[i], "")

	}

	b.ResetTimer()

	for i := n / 2; i < n; i++ {

		skipList.Get(bla[i])

	}

}

func BenchmarkSeanSkipListIncFind(b *testing.B) {

	var bla []float64

	skipList := seanSkipList.New()

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, float64(i))
	}

	for i := 0; i < n/2; i++ {

		skipList.Set(bla[i], "")

	}

	b.ResetTimer()

	for i := n / 2; i < n; i++ {

		skipList.Get(bla[i])

	}

}

type Int int

func (a Int) Less(b btree.Item) bool {
	return a < b.(Int)
}

func BenchmarkGBTreeSkipListRandAdd(b *testing.B) {

	var bla []Int

	bTree := btree.New(1024)

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, Int(rand.Intn(100_000_000)))
	}

	for i := 0; i < (n / 2); i++ {

		bTree.ReplaceOrInsert(bla[i])

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		bTree.ReplaceOrInsert(bla[i])

	}

}

func BenchmarkGBTreeSkipListIncAdd(b *testing.B) {

	var bla []Int

	bTree := btree.New(1024)

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, Int(i))
	}

	for i := 0; i < (n / 2); i++ {

		bTree.ReplaceOrInsert(bla[i])

	}

	b.ResetTimer()

	for i := (n / 2); i < b.N; i++ {

		bTree.ReplaceOrInsert(bla[i])

	}

}

func BenchmarkGBTreeSkipListRandFind(b *testing.B) {

	var bla []Int

	bTree := btree.New(1024)

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, Int(rand.Intn(100_000_000)))
	}

	for i := 0; i < n/2; i++ {

		bTree.ReplaceOrInsert(bla[i])

	}

	b.ResetTimer()

	for i := n / 2; i < n; i++ {

		bTree.Get(bla[i])

	}

}

func BenchmarkGBTreeSkipListIncFind(b *testing.B) {

	var bla []Int

	bTree := btree.New(1024)

	n := 20_000_000

	for i := 0; i < n; i++ {

		bla = append(bla, Int(i))
	}

	for i := 0; i < n/2; i++ {

		bTree.ReplaceOrInsert(bla[i])

	}

	b.ResetTimer()

	for i := n / 2; i < n; i++ {

		bTree.Get(bla[i])

	}
}

func TestGetMinMax(t *testing.T) {

	splitList := NewSplitList(1024)

	//skipList := New()

	n := 100000

	for i := 500; i < n; i++ {

		splitList.Add(i)

	}

	min := splitList.GetMin()
	max := splitList.GetMax()

	if min != 500 {

		t.Errorf("Failed to get min")

	}

	if max != n-1 {

		t.Errorf("Failed to get max")

	}

	fmt.Println(splitList.GetMin())
	fmt.Println(splitList.GetMax())

}

func TestPopMinMax(t *testing.T) {

	splitList := NewSplitList(1024)

	//skipList := New()

	n := 100000

	for i := 0; i < n; i++ {

		splitList.Add(i)

	}

	for i := 0; i < n; i++ {

		poppity := splitList.PopMin()
		fmt.Println(poppity)

	}

	fmt.Println(splitList.Length)

}
