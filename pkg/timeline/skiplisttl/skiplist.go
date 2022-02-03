package skiplisttl

import "github.com/MauriceGit/skiplist"

// fixedSkipList fixes some of the functionality of skiplist.SkipList, specifically around duplicate keys.
type fixedSkipList struct {
	*skiplist.SkipList
}

// Delete is like SkipList.Delete except it only removes the element(s) that is/are equal to e,
// SkipList.Delete removes only elements who's key is equal to e's key.
// Delete has time complexity O(M.log(N)) where M is the number of elements with e's key.
func (sl *fixedSkipList) Delete(e skiplist.ListElement) {
	candidates := sl.FindAll(e)
	if len(candidates) == 0 {
		return // nothing to remove
	}

	for i, candidate := range candidates {
		if candidate.GetValue() == e {
			candidates[i] = nil
		}
	}

	// remove all items that share a key with e
	sl.DeleteAll(e)

	// add back in those candidates that are != e
	for _, candidate := range candidates {
		if candidate != nil {
			sl.Insert(candidate.GetValue())
		}
	}
}

// DeleteAll is like SkipList.Delete except it removes all elements of the list with a matching key to e.
func (sl *fixedSkipList) DeleteAll(e skiplist.ListElement) {
	if sl.IsEmpty() {
		return
	}
	for {
		oldLen := sl.GetNodeCount()
		sl.SkipList.Delete(e)
		newLen := sl.GetNodeCount()
		if oldLen == newLen || newLen == 0 {
			return
		}
	}
}

// Find is like SkipList.Find except it guarantees that if there are duplicate keys,
// this returns the first matching element in the list.
func (sl *fixedSkipList) Find(e skiplist.ListElement) (*skiplist.SkipListElement, bool) {
	item, ok := sl.SkipList.Find(e)
	if !ok {
		return item, false
	}
	// It's really annoying that the list loops when calling Next or Prev,
	// we don't want to do that here
	first := sl.GetSmallestNode()
	key := e.ExtractKey()

	// find the first item that is still at this time
	for item != first {
		prev := sl.Prev(item)
		if prev == nil {
			break
		}
		if prev.GetValue().ExtractKey() != key {
			break
		}
		item = prev
	}

	return item, true
}

// FindAll returns all elements whose key is equal to e's key.
func (sl *fixedSkipList) FindAll(e skiplist.ListElement) []*skiplist.SkipListElement {
	item, ok := sl.Find(e)
	if !ok {
		return nil
	}

	// It's really annoying that the list loops when calling Next or Prev,
	// we don't want to do that here
	first := sl.GetSmallestNode()
	key := e.ExtractKey()

	var result []*skiplist.SkipListElement
	for item.GetValue().ExtractKey() == key {
		result = append(result, item)
		item = sl.Next(item)
		if item == first {
			// we have looped around all the items
			break
		}
	}
	return result
}
