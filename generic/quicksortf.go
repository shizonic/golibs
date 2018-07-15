// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package generic

func QuicksortTypeSlice(v []Type, less func(a, b *Type) bool) {
	switch len(v) {
	case 0, 1:
		return
	case 2:
		if less(&v[1], &v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		return
	case 3:
		if less(&v[1], &v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		if less(&v[2], &v[1]) {
			v[1], v[2] = v[2], v[1]
		}
		if less(&v[1], &v[0]) {
			v[0], v[1] = v[1], v[0]
		}
		return
	}

	i := PartitionTypeSlice(v, less)
	QuicksortTypeSlice(v[:i+1], less)
	QuicksortTypeSlice(v[i+1:], less)
}

// Hoare's partitioning with median of first, middle, and last as pivot
func PartitionTypeSlice(v []Type, less func(a, b *Type) bool) int {
	var pivot Type

	if len(v) > 16 {
		pivot = MedianOfThreeType(v, less)
	} else {
		pivot = v[(len(v)-1)/2]
	}

	i, j := -1, len(v)

	for {
		for {
			i++
			if !less(&v[i], &pivot) {
				break
			}
		}

		for {
			j--
			if !less(&pivot, &v[j]) {
				break
			}
		}

		if i < j {
			v[i], v[j] = v[j], v[i]
		} else {
			return j
		}
	}
}

func MedianOfThreeType(v []Type, less func(a, b *Type) bool) Type {
	a := v[0]
	b := v[(len(v)-1)/2]
	c := v[len(v)-1]

	if less(&b, &a) {
		a, b = b, a
	}
	if less(&c, &b) {
		b, c = c, b
	}
	if less(&b, &a) {
		a, b = b, a
	}

	return b
}
