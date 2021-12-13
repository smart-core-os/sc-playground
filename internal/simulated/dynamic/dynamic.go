// Package dynamic implements values that change themselves over time in a specified way.
// A Float32 can be interpolated, or run a Profile, which is a time series of Segments.
package dynamic

import "time"

const DefaultUpdateInterval = 100 * time.Millisecond
