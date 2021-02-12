package src

import (
	"math"
	"math/rand"
)

type KeyValue struct {
	key   int
	value string
}

type MinMaxList struct {
	Indexes []KeyValue
	Height  int
	Min     int
	Max     int
}

type TeleportList struct {
	Sublists      [30]MinMaxList
	Data          []int
	CurrentHeight int
	Length        int
}

func NewTeleportList() TeleportList {
	teleportList := TeleportList{}
	teleportList.CurrentHeight = -1
	for i := 0; i < 30; i++ {

		newMMList := MinMaxList{}
		newMMList.Max = math.MinInt64
		newMMList.Min = math.MaxInt64
		newMMList.Height = -1

		teleportList.Sublists[i] = newMMList
	}
	return teleportList
}

func (t *TeleportList) Add(value int) {
	candidateHeight := int(math.Abs(math.Log2(rand.Float64())))

	if t.Sublists[candidateHeight].Height == -1 {

		t.Sublists[candidateHeight].Height = candidateHeight

	}

	candidateHeightSublist := &t.Sublists[candidateHeight]

	candidateHeightSublist.Indexes = insortKeyValue(candidateHeightSublist.Indexes, value)

	t.Length += 1
	if candidateHeightSublist.Max < value {
		candidateHeightSublist.Max = value
	}
	if candidateHeightSublist.Min > value {
		candidateHeightSublist.Min = value
	}

	//sort.Sort(t.SSublists)
	//	t.Data = insortInt(t.Data, value)

}

func (t TeleportList) Find(value int) bool {

	if last := len(t.Sublists) - 1; last >= 0 {

		for i := last; i >= 0; i-- {

			if t.Sublists[i].Min <= value {

				if t.Sublists[i].Max >= value {

					bsearch := BinarySearchKeyValue(t.Sublists[i].Indexes, value)
					//bsearch := InterpolatedBinarySearchKeyValue(t.Sublists[i].Indexes, value)

					if bsearch == true {

						return true

					}

				}

			}

		}
	}

	return false

}

func (t TeleportList) Index(value int) bool {

	return t.Find(t.Data[value-1])

}
