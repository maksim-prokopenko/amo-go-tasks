package main

import (
	"cmp"
	"container/heap"
	"iter"
	"slices"
)

// Эта сложность избыточна, но мне было весело.

func mergeChs[S ~[]<-chan E, E cmp.Ordered](chs S) iter.Seq[E] {
	return mergeChsFunc(chs, func(a, b E) bool { return a < b })
}

func mergeChsFunc[S ~[]<-chan E, E any](chs S, less func(a, b E) bool) iter.Seq[E] {
	nodes := make([]*node[E], 0, len(chs))

	for _, ch := range chs {
		n := &node[E]{ch: ch}
		if n.update() {
			nodes = append(nodes, n)
		}
	}

	hn := heapNode[E]{nodes: nodes, less: less}
	heap.Init(&hn)

	return func(yield func(E) bool) {
		for hn.Len() > 0 {

			if !yield(hn.nodes[0].val) {
				return
			}

			if hn.nodes[0].update() {
				heap.Fix(&hn, 0)
			} else {
				heap.Pop(&hn)
			}

		}
	}
}

type node[E any] struct {
	ch  <-chan E
	val E
}

func (n *node[E]) update() (ok bool) {
	n.val, ok = <-n.ch
	return
}

type heapNode[E any] struct {
	nodes []*node[E]
	less  func(a, b E) bool
}

func (hn heapNode[E]) Len() int           { return len(hn.nodes) }
func (hn heapNode[E]) Less(i, j int) bool { return hn.less(hn.nodes[i].val, hn.nodes[j].val) }
func (hn heapNode[E]) Swap(i, j int)      { hn.nodes[i], hn.nodes[j] = hn.nodes[j], hn.nodes[i] }

func (hn heapNode[E]) Push(any) {} // not used in this case
func (hn *heapNode[E]) Pop() any {
	hn.nodes = slices.Delete(hn.nodes, hn.Len()-1, hn.Len())
	return nil // value not used in this case
}
