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

func (t *tree[T]) hasNodeWith(n T) bool {
	checked := make([]*xtreeNode[T], 0)

	_, ok := search(t.start, n, t.dir, -1, checked)
	return ok
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

func search[T comparable](nd *xtreeNode[T], val T, dir vdirection, depth int, cache []*xtreeNode[T]) (*xtreeNode[T], bool) {
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

	if nd.value == val {
		return nd, true
	}

	////

	searchTop := func() (*xtreeNode[T], bool) {
		if nd.topleft != nil {
			if node, ok := search(nd.topleft, val, dir, nxtlim, cache); ok {
				return node, ok
			}
		}

		if nd.topright != nil {
			if node, ok := search(nd.topright, val, dir, nxtlim, cache); ok {
				return node, ok
			}
		}

		return nil, false
	}

	searchBottom := func() (*xtreeNode[T], bool) {
		if nd.bottomleft != nil {
			if node, ok := search(nd.bottomleft, val, dir, nxtlim, cache); ok {
				return node, ok
			}
		}

		if nd.bottomright != nil {
			if node, ok := search(nd.bottomright, val, dir, nxtlim, cache); ok {
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
