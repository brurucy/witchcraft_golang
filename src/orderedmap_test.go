package src

import (
	"fmt"
	gbtree "github.com/google/btree"
	"testing"
)

type intT struct {
	val int
}

func (i *intT) Less(other gbtree.Item) bool {
	return i.val < other.(*intT).val
}

// Don't trust this.
func TestAddFindTiny(t *testing.T) {
	splitList := NewSplitList(1024)

	for i := 0; i < 10; i++ {
		splitList.Add(&intT{i})
	}

	for _, list := range splitList.ListOfBucketLists {
		for _, bucket := range list.Buckets {
			for _, index := range bucket.Indexes {
				fmt.Print(index, "\t")
			}
			fmt.Println("---")
		}
		fmt.Println()
	}

	for i := 0; i < 10; i++ {
		if !splitList.Find(&intT{i}) {
			//t.Log(i)
			t.Fail()
		}
	}
}

func TestGetMin(t *testing.T) {

	splitList := NewSplitList(1024)

	//skipList := New()

	n := 1_000_000

	for i := n; i >= 0; i-- {
		ints := &intT{i}

		splitList.Add(ints)

		if (splitList.Find(ints) == false || splitList.LookupReverse(ints) == false) || splitList.GetMin() != ints {
			t.Fatal()
		}

	}

}

func TestPopMin(t *testing.T) {

	splitList := NewSplitList(1024)

	//skipList := New()

	n := 1_000_000

	for i := n; i >= 0; i-- {
		ints := &intT{i}

		splitList.Add(ints)

		if splitList.Find(ints) == false || splitList.GetMin() != ints {
			t.Fatal()
		}

	}

	for i := 0; i <= n; i++ {
		ints := &intT{i}

		poppity := splitList.PopMin()

		if (!poppity.Less(ints) && !ints.Less(poppity)) != true {
			fmt.Println(poppity, ints)
			t.Fatal()
		}

	}

	if splitList.Length != 0 {
		t.Fatal()
	}

}

func TestGetMax(t *testing.T) {

	splitList := NewSplitList(1024)

	//skipList := New()

	n := 1_000_000

	for i := 0; i <= n; i++ {
		ints := &intT{i}

		splitList.Add(ints)
		//fmt.Println(splitList.GetMax())

		if splitList.Find(ints) == false || splitList.GetMax() != ints {
			t.Fatal()
		}

	}

}

func TestPopMax(t *testing.T) {

	splitList := NewSplitList(1024)

	//skipList := New()

	n := 1_000_000

	for i := n; i >= 0; i-- {
		ints := &intT{i}

		splitList.Add(ints)

	}

	for i := n; i >= 0; i-- {
		ints := &intT{i}
		poppity := splitList.PopMax()

		if (!poppity.Less(ints) && !ints.Less(poppity)) != true {
			t.Fatal()
		}
	}

	if splitList.Length != 0 {
		fmt.Println(splitList.Length)
		t.Fatal()
	}

}

func TestDelete(t *testing.T) {

	splitList := NewSplitList(1024)

	n := 1_000_000

	for i := n; i >= 0; i-- {
		ints := &intT{i}

		splitList.Add(ints)

	}

	for i := n; i >= 0; i-- {
		ints := &intT{i}

		splitList.Delete(ints)

		if splitList.Find(ints) == true {
			t.Fatal()
		}

	}

	if splitList.Length != 0 {
		fmt.Println(splitList.Length)
		t.Fatal()
	}

}
