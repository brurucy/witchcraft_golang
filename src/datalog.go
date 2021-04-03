package src

import gbtree "github.com/google/btree"

type DirectedEdge struct {
	From int
	To   int
}

type Joinable interface {
	Joinable(DirectedEdge) bool
}

func (d *DirectedEdge) Less(other gbtree.Item) bool {

	if d.From == other.(*DirectedEdge).From && d.To < other.(*DirectedEdge).To {

		return true

	} else if d.From < other.(*DirectedEdge).From && d.To == other.(*DirectedEdge).To {

		return true

	} else if d.From < other.(*DirectedEdge).From && d.To < other.(*DirectedEdge).To {

		return true

	} else if d.From < other.(*DirectedEdge).From && d.To > other.(*DirectedEdge).To {

		return true

	}
	return false

}

func (d *DirectedEdge) Joinable(other DirectedEdge) bool {
	return d.To == other.From
}
