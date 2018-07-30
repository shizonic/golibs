// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package graph

import (
	"math/bits"
	"unsafe"

	"github.com/guns/golibs/bitslice"
	"github.com/guns/golibs/generic/impl"
)

// A Workspace provides general-purpose scratch storage for Graph methods.
type Workspace struct {
	len, cap int // Logical len/cap, not buffer len/cap
	a, b, c  []int
}

// NewWorkspace returns a new Workspace for a Graph of a given size.
func NewWorkspace(size int) *Workspace {
	// Single shared int buffer
	n := 1 << uint(bits.Len(uint(size*3-1)))
	buf := make([]int, n)

	return &Workspace{
		len: size,
		cap: n / 3,
		a:   buf[:size],
		b:   buf[size : size*2],
		c:   buf[size*2 : size*3],
	}
}

// workspaceField values represent fields of a Workspace.
type workspaceField uint

const (
	wA    workspaceField = 1 << iota // Select or reset (*Workspace).a
	wB                               // Select or reset (*Workspace).b
	wC                               // Select or reset (*Workspace).c
	wANeg                            // Fill (*Workspace).a with undefined
	wBNeg                            // Fill (*Workspace).b with undefined
	wCNeg                            // Fill (*Workspace).c with undefined
)

func (w *Workspace) selectSlice(field workspaceField) []int {
	switch field {
	case wA:
		return w.a
	case wB:
		return w.b
	case wC:
		return w.c
	default:
		return nil // panic() defeats inlining [go1.11]
	}

}

// makeQueue returns an empty IntQueue with the specified field as a backing slice.
func (w *Workspace) makeQueue(field workspaceField) impl.IntQueue {
	q := *impl.NewIntQueueWithBuffer(w.selectSlice(field))
	q.SetAutoGrow(false)
	q.Reset()
	return q
}

// makeStack returns an empty IntStack with the specified field as a backing slice.
func (w *Workspace) makeStack(field workspaceField) impl.IntStack {
	s := *impl.NewIntStackWithBuffer(w.selectSlice(field))
	s.SetAutoGrow(false)
	s.Reset()
	return s
}

// makeBitsliceN returns a slice of n empty bitslice.T with the specified
// field as a backing slice. Each bitslice has a capacity equal to the current
// size of the workspace. The maximum number of bitslices that can be returned
// is equal to:
//
//	currentWorkspaceLen / bitslice.UintLen(currentWorkspaceLen)
//
func (w *Workspace) makeBitsliceN(n int, field workspaceField) []bitslice.T {
	buf := w.selectSlice(field)[:w.cap]
	bs := make([]bitslice.T, n)
	blen := bitslice.UintLen(w.len)
	offset := 0

	for i := 0; i < n; i++ {
		b := buf[offset : offset+blen]
		bs[i] = *(*bitslice.T)(unsafe.Pointer(&b))
		offset += blen
	}

	s := buf[:offset]
	for i := range s {
		s[i] = 0
	}

	return bs
}

// makeAutoPromotingStack returns an autoPromotingStack with the specified
// fields as a backing slice. The fields parameter must specify two contiguous
// internal fields.
func (w *Workspace) makeAutoPromotingStack(fields workspaceField) autoPromotingStack {
	var buf []int

	switch fields {
	case wA | wB:
		buf = w.a[:w.len*2]
	case wB | wC:
		buf = w.b[:w.len*2]
	}

	for i := range buf {
		buf[i] = undefined
	}

	return *newAutoPromotingStack(makeListNodeSlice(buf))
}

// reset prepares a Workspace for a Graph of a given size. The fields
// parameter is a bitfield of workspaceField values that specify which
// fields to clear.
func (w *Workspace) reset(size int, fields workspaceField) {
	switch {
	case size == w.len:
		// No resize necessary
	case size <= w.cap:
		w.len = size
		w.a = w.a[:size]
		w.b = w.b[:size]
		w.c = w.c[:size]
	default:
		*w = *NewWorkspace(size)
		// New workspaces are zero-filled, so avoid unnecessary work.
		fields &= ^(wA | wB | wC)
	}

	w.clear(fields)
}

// clear specific fields a Workspace.
func (w *Workspace) clear(fields workspaceField) {
	if fields == 0 {
		return
	}

	if fields&wA > 0 {
		for i := range w.a {
			w.a[i] = 0
		}
	} else if fields&wANeg > 0 {
		for i := range w.a {
			w.a[i] = undefined
		}
	}

	if fields&wB > 0 {
		for i := range w.b {
			w.b[i] = 0
		}
	} else if fields&wBNeg > 0 {
		for i := range w.b {
			w.b[i] = undefined
		}
	}

	if fields&wC > 0 {
		for i := range w.c {
			w.c[i] = 0
		}
	} else if fields&wCNeg > 0 {
		for i := range w.c {
			w.c[i] = undefined
		}
	}
}
