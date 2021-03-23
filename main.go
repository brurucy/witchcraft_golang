package main

import (
	"fmt"
	mauriceSkipList "github.com/MauriceGit/skiplist"
	gbtree "github.com/google/btree"
	seanSkipList "github.com/sean-public/fast-skiplist"
	tbtree "github.com/tidwall/btree"
	"github.com/tidwall/lotsa"
	"math/rand"
	"os"
	"sort"
	"witchcraft/src"
)

type intT struct {
	val int
}

func (i *intT) Less(other gbtree.Item) bool {
	return i.val < other.(*intT).val
}

type Element int

func (e Element) ExtractKey() float64 {
	return float64(e)
}
func (e Element) String() string {
	return fmt.Sprintf("%03d", e)
}

func main() {
	less := func(a, b interface{}) bool {
		return a.(*intT).val < b.(*intT).val
	}
	N := 1_000_000
	keys := make([]intT, N)
	for i := 0; i < N; i++ {
		keys[i] = intT{i}
	}
	lotsa.Output = os.Stdout
	lotsa.MemUsage = true

	sortInts := func() {
		sort.Slice(keys, func(i, j int) bool {
			return less(&keys[i], &keys[j])
		})
	}

	shuffleInts := func() {
		for i := range keys {
			j := rand.Intn(i + 1)
			keys[i], keys[j] = keys[j], keys[i]
		}
	}

	println()
	println("** sequential set **")

	print("mauricesl:  set-seq\t")
	skipList := mauriceSkipList.New()
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Insert(Element(keys[i].val))
	})
	print("seansl:  set-seq\t")
	skipList2 := seanSkipList.New()
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Set(float64(keys[i].val), "")
	})
	print("google:  set-seq\t")
	tr2 := gbtree.New(256)
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.ReplaceOrInsert(&keys[i])
	})
	print("tidwall: set-seq\t")
	tr := tbtree.New(less)
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})
	print("tidwall: set-seq-hint\t")
	tr = tbtree.New(less)
	sortInts()
	var hint tbtree.PathHint
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.SetHint(&keys[i], &hint)
	})
	print("tidwall: load-seq\t")
	tr = tbtree.New(less)
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Load(&keys[i])
	})
	print("go-arr:  append\t\t")
	var arr []interface{}
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		arr = append(arr, &keys[i])
	})
	print("splitlist: add\t\t")
	sortInts()
	tsl := src.NewSplitList(1024)
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Add(keys[i].val)
	})

	println()
	println("** random set **")

	print("mauricesl:  set-seq\t")
	skipList = mauriceSkipList.New()
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Insert(Element(keys[i].val))
	})
	print("seansl:  set-seq\t")
	skipList2 = seanSkipList.New()
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Set(float64(keys[i].val), "")
	})
	print("google:  set-rand\t")
	tr2 = gbtree.New(256)
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.ReplaceOrInsert(&keys[i])
	})
	print("tidwall: set-rand\t")
	tr = tbtree.New(less)
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})
	print("tidwall: set-rand-hint\t")
	tr = tbtree.New(less)
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.SetHint(&keys[i], &hint)
	})
	print("tidwall: set-again\t")
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})
	print("tidwall: set-after-copy\t")
	tr = tr.Copy()
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})
	print("tidwall: load-rand\t")
	tr = tbtree.New(less)
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Load(&keys[i])
	})
	print("splitlist: add\t\t")
	shuffleInts()
	tsl = src.NewSplitList(1024)
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Add(keys[i].val)
	})

	println()
	println("** sequential get **")

	print("mauricesl:  get		")
	skipList = mauriceSkipList.New()
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Find(Element(keys[i].val))
	})
	print("seansl:  get\t\t")
	skipList2 = seanSkipList.New()
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Get(float64(keys[i].val))
	})
	print("google:  get-seq\t")
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.Get(&keys[i])
	})
	print("tidwall: get-seq\t")
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Get(&keys[i])
	})
	print("tidwall: get-seq-hint\t")
	sortInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.GetHint(&keys[i], &hint)
	})
	print("splitlist: find seq\t")
	sortInts()
	tsl = src.NewSplitList(1024)
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Find(keys[i].val)
	})

	println()
	println("** random get **")

	print("mauricesl:  get\t\t")
	skipList = mauriceSkipList.New()
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Find(Element(keys[i].val))
	})
	print("seansl:  get\t\t")
	skipList2 = seanSkipList.New()
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Get(float64(keys[i].val))
	})
	print("google:  get-rand\t")
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.Get(&keys[i])
	})
	print("tidwall: get-rand\t")
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Get(&keys[i])
	})
	print("tidwall: get-rand-hint\t")
	shuffleInts()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.GetHint(&keys[i], &hint)
	})
	print("splitlist: find random\t")
	shuffleInts()
	tsl = src.NewSplitList(1024)
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Find(keys[i].val)
	})
}
