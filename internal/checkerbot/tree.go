package checkerbot

import (
	"fmt"

	"github.com/ifebles/chalhub/pkg/util"
)

type vdirection uint8

const (
	up vdirection = iota
	down
	vboth
)

type hdirection uint8

const (
	left hdirection = iota
	right
	hboth
)

type tree[T comparable] struct {
	dir   vdirection
	start *xtreeNode[T]
}

type xtreeNode[T comparable] struct {
	topleft, topright,
	bottomleft, bottomright *xtreeNode[T]

	value T
}

type pathMarker[T comparable] struct {
	id   int
	path []int
	val  T
}

func (t tree[T]) getPathCollection(stpif func(T) bool, ignif func(T, T) bool) []pathMarker[T] {
	result := []pathMarker[T]{}

	if t.start == nil {
		return result
	}

	result = getPaths(t.start, new(int), []int{}, stpif, ignif, []*xtreeNode[T]{})

	return result
}

func getPaths[T comparable](
	nd *xtreeNode[T], cnt *int, histids []int, stp func(T) bool, ign func(T, T) bool, cache []*xtreeNode[T],
) []pathMarker[T] {
	result := []pathMarker[T]{}

	if nd == nil {
		return result
	}

	if _, ok := util.Find(cache, func(i *xtreeNode[T]) bool { return i == nd }); ok {
		return result
	} else {
		cache = append(cache, nd)
	}

	////

	*cnt++
	id := *cnt
	upthist := make([]int, len(histids)+1)
	copy(upthist, histids)
	upthist[len(upthist)-1] = id

	result = append(result, pathMarker[T]{id, upthist, nd.value})

	if stp(nd.value) {
		return result
	}

	nodes := []*xtreeNode[T]{
		nd.topleft, nd.topright,
		nd.bottomleft, nd.bottomright,
	}

	////

	for _, a := range nodes {
		if a == nil || ign(nd.value, a.value) {
			continue
		}

		nucache := make([]*xtreeNode[T], len(cache))
		copy(nucache, cache)

		resp := getPaths(a, cnt, upthist, stp, ign, nucache)
		result = append(result, resp...)
	}

	return result
}

func (t *tree[T]) hasNodeWith(v T) bool {
	fn := func(vl T) bool { return vl == v }
	_, ok := search(t.start, fn, t.dir, -1, []*xtreeNode[T]{})

	return ok
}

func (xn *xtreeNode[T]) getNodeWith(v T, dir vdirection, ign ...*xtreeNode[T]) *xtreeNode[T] {
	fn := func(vl T) bool { return vl == v }
	ignore := make([]*xtreeNode[T], len(ign))
	copy(ignore, ign)

	nd, _ := search(xn, fn, dir, -1, ignore)

	return nd
}

func (xn *xtreeNode[T]) findNode(fn func(T) bool) *xtreeNode[T] {
	nd, _ := search(xn, fn, vboth, -1, []*xtreeNode[T]{})

	return nd
}

func (xn *xtreeNode[T]) set(n *xtreeNode[T], vdir vdirection, hdir hdirection) {
	if hdir == hboth || vdir == vboth {
		panic("cannot set in both directions")
	}

	if vdir == up {
		if hdir == left {
			xn.topleft = n
			n.bottomright = xn
		} else {
			xn.topright = n
			n.bottomleft = xn
		}
	} else {
		if hdir == left {
			xn.bottomleft = n
			n.topright = xn
		} else {
			xn.bottomright = n
			n.topleft = xn
		}
	}
}

func (xn *xtreeNode[T]) add(n *xtreeNode[T], vdir vdirection, hdir hdirection) error {
	if hdir == hboth || vdir == vboth {
		panic("cannot add in both directions")
	}

	if vdir == up {
		if hdir == left {
			if xn.topleft != nil {
				return fmt.Errorf("current node has a 'topleft' reference")
			}

			if n.bottomright != nil {
				return fmt.Errorf("new node has a 'bottomright' reference")
			}

			xn.topleft = n
			n.bottomright = xn

			return nil
		} else {
			if xn.topright != nil {
				return fmt.Errorf("current node has a 'topright' reference")
			}

			if n.bottomleft != nil {
				return fmt.Errorf("new node has a 'bottomleft' reference")
			}

			xn.topright = n
			n.bottomleft = xn

			return nil
		}
	} else {
		if hdir == left {
			if xn.bottomleft != nil {
				return fmt.Errorf("current node has a 'bottomleft' reference")
			}

			if n.topright != nil {
				return fmt.Errorf("new node has a 'topright' reference")
			}

			xn.bottomleft = n
			n.topright = xn

			return nil
		} else {
			if xn.bottomright != nil {
				return fmt.Errorf("current node has a 'bottomright' reference")
			}

			if n.topleft != nil {
				return fmt.Errorf("new node has a 'topleft' reference")
			}

			xn.bottomright = n
			n.topleft = xn

			return nil
		}
	}
}

func search[T comparable](nd *xtreeNode[T], fn func(T) bool, dir vdirection, depth int, cache []*xtreeNode[T]) (*xtreeNode[T], bool) {
	if depth < -1 {
		panic(fmt.Sprintf("invalid depth: %d", depth))
	}

	if depth == 0 || nd == nil {
		return nil, false
	}

	nxtlim := -1

	if depth != -1 {
		nxtlim = depth - 1
	}

	if _, ok := util.Find(cache, func(i *xtreeNode[T]) bool { return i == nd }); ok {
		return nil, false
	} else {
		cache = append(cache, nd)
	}

	////

	if fn(nd.value) {
		return nd, true
	}

	////

	searchTop := func() (*xtreeNode[T], bool) {
		if nd.topleft != nil {
			if node, ok := search(nd.topleft, fn, dir, nxtlim, cache); ok {
				return node, ok
			}
		}

		if nd.topright != nil {
			if node, ok := search(nd.topright, fn, dir, nxtlim, cache); ok {
				return node, ok
			}
		}

		return nil, false
	}

	searchBottom := func() (*xtreeNode[T], bool) {
		if nd.bottomleft != nil {
			if node, ok := search(nd.bottomleft, fn, dir, nxtlim, cache); ok {
				return node, ok
			}
		}

		if nd.bottomright != nil {
			if node, ok := search(nd.bottomright, fn, dir, nxtlim, cache); ok {
				return node, ok
			}
		}

		return nil, false
	}

	////

	switch dir {
	case up:
		return searchTop()

	case down:
		return searchBottom()

	default: // both
		if node, ok := searchTop(); ok {
			return node, ok
		}

		if node, ok := searchBottom(); ok {
			return node, ok
		}
	}

	return nil, false
}
