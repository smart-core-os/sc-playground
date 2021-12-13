// Package dynamic implements values that change themselves over time in a specified way.
// A Value can be interpolated, or run a Profile, which is a time series of Segments.
package dynamic

import "time"

const UpdateRate = 100 * time.Millisecond
