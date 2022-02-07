package timeline

import (
	"fmt"
)

// Equal compares two TLs for equality.
func Equal(a, b TL) bool {
	equal, _ := EqualExplain(a, b)
	return equal
}

// EqualExplain compares two TLs for equality, returning a reason if they aren't equal.
func EqualExplain(a, b TL) (bool, string) {
	if Empty(a) != Empty(b) {
		if Empty(a) {
			return false, "a is empty, b is not"
		} else {
			return false, "b is empty, a is not"
		}
	}
	aFirst, aLast, aExists := Bound(a)
	bFirst, bLast, bExists := Bound(b)
	if !aFirst.Equal(bFirst) || !aLast.Equal(bLast) || aExists != bExists {
		return false, fmt.Sprintf("a bounds { %v %v %v } != b bounds { %v %v %v }", aFirst, aLast, aExists, bFirst, bLast, bExists)
	}

	// check items
	cur, exists := aFirst, aExists
	for exists {
		aItems := a.At(cur)
		bItems := b.At(cur)
		if len(aItems) != len(bItems) {
			return false, fmt.Sprintf("at %v, entry lengths %v != %v", cur, len(aItems), len(bItems))
		}
		for i, aItem := range aItems {
			if aItem != bItems[i] {
				return false, fmt.Sprintf("at %v index %v, %v != %v", cur, i, len(aItems), len(bItems))
			}
		}

		aNext, aExists := a.Next(cur)
		bNext, bExists := b.Next(cur)
		if !aNext.Equal(bNext) || aExists != bExists {
			return false, fmt.Sprintf("at %v, next { %v %v } != { %v %v}", cur, aNext, aExists, bNext, bExists)
		}
		cur, exists = aNext, aExists
	}

	return true, ""
}
