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
	N := 10_000_000
	var temp intT
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
	sortInts()

	print("mauricesl:  set-seq\t")
	skipList := mauriceSkipList.New()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Insert(Element(keys[i].val))
	})

	print("seansl:  set-seq\t")
	skipList2 := seanSkipList.New()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Set(float64(keys[i].val), "")
	})

	print("google:  set-seq\t")
	tr2 := gbtree.New(256)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.ReplaceOrInsert(&keys[i])
	})

	print("tidwall: set-seq\t")
	tr := tbtree.New(less)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})

	print("tidwall: set-seq-hint\t")
	tr = tbtree.New(less)
	var hint tbtree.PathHint
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.SetHint(&keys[i], &hint)
	})

	print("tidwall: load-seq\t")
	tr = tbtree.New(less)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Load(&keys[i])
	})

	print("go-arr:  append-seq\t")
	var arr []interface{}
	lotsa.Ops(N, 1, func(i, _ int) {
		arr = append(arr, &keys[i])
	})

	print("splitlist: add\t\t")
	tsl := src.NewSplitList(1024)
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Add(keys[i].val)
	})

	print("go-hashmap: set-seq\t")
	hm := make(map[int]intT, 0)
	lotsa.Ops(N, 1, func(i, _ int) {
		hm[i] = keys[i]
	})

	println()
	println("** sequential get **")

	print("mauricesl:  find-seq\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Find(Element(keys[i].val))
	})

	print("seansl:  get-seq\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Get(float64(keys[i].val))
	})

	print("google:  get-seq\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.Get(&keys[i])
	})

	print("tidwall: get-seq\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Get(&keys[i])
	})

	print("tidwall: get-seq-hint\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.GetHint(&keys[i], &hint)
	})

	print("splitlist: find seq\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Find(keys[i].val)
	})

	print("go-hashmap: get-seq\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		temp = hm[i]
	})

	println()
	println("** random set **")
	shuffleInts()

	print("mauricesl: set-rand\t")
	skipList = mauriceSkipList.New()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Insert(Element(keys[i].val))
	})

	print("seansl: set-rand\t")
	skipList2 = seanSkipList.New()
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Set(float64(keys[i].val), "")
	})

	print("google: set-rand\t")
	tr2 = gbtree.New(256)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.ReplaceOrInsert(&keys[i])
	})

	print("tidwall: set-rand\t")
	tr = tbtree.New(less)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})

	print("tidwall: set-rand-hint\t")
	tr = tbtree.New(less)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.SetHint(&keys[i], &hint)
	})

	print("tidwall: set-rand-again\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})

	print("tidwall: set-copy-rand\t")
	tr = tr.Copy()
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Set(&keys[i])
	})

	print("tidwall: load-rand\t")
	tr = tbtree.New(less)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Load(&keys[i])
	})

	print("splitlist: add-rand\t")
	tsl = src.NewSplitList(1024)
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Add(keys[i].val)
	})

	print("go-hashmap: set-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		hm[keys[i].val] = keys[i]
	})

	println()
	println("** random get **")

	print("mauricesl: find-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList.Find(Element(keys[i].val))
	})

	print("seansl: get-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		skipList2.Get(float64(keys[i].val))
	})

	print("google: get-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr2.Get(&keys[i])
	})

	print("tidwall: get-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Get(&keys[i])
	})

	print("tidwall: get-hint-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.GetHint(&keys[i], &hint)
	})

	print("splitlist: find-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		tsl.Find(keys[i].val)
	})

	print("go-hashmap: get-rand\t")
	lotsa.Ops(N, 1, func(i, _ int) {
		temp = hm[keys[i].val]
	})
}
