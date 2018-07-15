// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package generic provides generic templates for common data structures and
// functions via https://github.com/cheekybits/genny/generic
package generic

import "github.com/cheekybits/genny/generic"

type GenericType generic.Type
type GenericNumber generic.Number
type ComparableType interface{ Less(x *ComparableType) bool } // generic.Type
