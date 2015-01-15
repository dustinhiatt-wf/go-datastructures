/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rangetree

import (
	"sort"
	"sync"
)

var entriesPool = sync.Pool{
	New: func() interface{} {
		return make(Entries, 0, 10)
	},
}

// Entries is a typed list of Entry that can be reused if Dispose
// is called.
type Entries []Entry

// Dispose will free the resources consumed by this list and
// allow the list to be reused.
func (entries *Entries) Dispose() {
	for i := 0; i < len(*entries); i++ {
		(*entries)[i] = nil
	}

	*entries = (*entries)[:0]
	entriesPool.Put(*entries)
}

func (entries Entries) search(dimension uint64, value int64) int {
	return sort.Search(
		len(entries),
		func(i int) bool { return entries[i].ValueAtDimension(dimension) >= value },
	)
}

// addAt will add the provided node at the provided index.  Returns
// a node if one was overwritten.
func (entries *Entries) addAt(i int, dimension uint64, entry Entry) Entry {
	if i == len(*entries) {
		*entries = append(*entries, entry)
		return nil
	}

	if (*entries)[i].ValueAtDimension(dimension) == entry.ValueAtDimension(dimension) {
		overwritten := (*entries)[i]
		// this is a duplicate, there can't be a duplicate
		// point in the last dimension
		(*entries)[i] = entry
		return overwritten
	}

	*entries = append(*entries, nil)
	copy((*entries)[i+1:], (*entries)[i:])
	(*entries)[i] = entry
	return nil
}

func (entries *Entries) add(dimension uint64, entry Entry) Entry {
	i := entries.search(dimension, entry.ValueAtDimension(dimension))
	return entries.addAt(i, dimension, entry)
}

func (entries *Entries) deleteAt(i int) {
	if i >= len(*entries) { // no matching found
		return
	}

	copy((*entries)[i:], (*entries)[i+1:])
	(*entries)[len(*entries)-1] = nil
	*entries = (*entries)[:len(*entries)-1]
}

func (entries *Entries) delete(dimension uint64, value int64) {
	i := entries.search(dimension, value)

	if (*entries)[i].ValueAtDimension(dimension) != value || i == len(*entries) {
		return
	}

	entries.deleteAt(i)
}

func (entries Entries) apply(low, high int64, dimension uint64, fn func(Entry) bool) bool {
	index := entries.search(dimension, low)
	if index == len(entries) {
		return true
	}

	for ; index < len(entries); index++ {
		if entries[index].ValueAtDimension(dimension) >= high {
			break
		}

		if !fn(entries[index]) {
			return false
		}
	}

	return true
}

func (entries Entries) get(dimension uint64, value int64) (Entry, int) {
	i := entries.search(dimension, value)
	if i == len(entries) {
		return nil, i
	}

	if entries[i].ValueAtDimension(dimension) == value {
		return entries[i], i
	}

	return nil, i
}

func (entries *Entries) getOrAdd(entry Entry,
	dimension, lastDimension uint64) (Entry, bool) {

	value := entry.ValueAtDimension(dimension)

	i := entries.search(dimension, value)
	if i == len(*entries) {
		*entries = append(*entries, entry)
		return entry, true
	}

	if (*entries)[i].ValueAtDimension(dimension) == value {
		return (*entries)[i], false
	}

	*entries = append(*entries, nil)
	copy((*entries)[i+1:], (*entries)[i:])
	(*entries)[i] = entry
	return entry, true
}

func (es Entries) flatten(entries *Entries) {
	*entries = append(*entries, es...)
}

func (entries *Entries) insert(insertDimension, dimension, maxDimension uint64,
	index, number int64, modified, deleted *Entries) {

	lastDimension := isLastDimension(maxDimension, dimension)

	if insertDimension == dimension {
		i := entries.search(dimension, index)
		var toDelete []int

		for j := i; j < len(*entries); j++ {
			if (*entries)[j].ValueAtDimension(dimension)+number < index {
				toDelete = append(toDelete, j)
				if lastDimension {
					*deleted = append(*deleted, (*entries)[j])
				}

				continue
			}
			if lastDimension {
				*modified = append(*modified, (*entries)[j])
			}
		}

		for i, index := range toDelete {
			entries.deleteAt(index - i)
		}

		return
	}
}

func (entries Entries) immutableInsert(insertDimension, dimension, maxDimension uint64,
	index, number int64, modified, deleted *Entries) dimensionalList {

	lastDimension := isLastDimension(maxDimension, dimension)

	cp := make(Entries, len(entries))
	copy(cp, entries)

	if insertDimension == dimension {
		i := cp.search(dimension, index)
		var toDelete []int

		for j := i; j < len(cp); j++ {
			if cp[j].ValueAtDimension(dimension)+number < index {
				toDelete = append(toDelete, j)
				if lastDimension {
					*deleted = append(*deleted, cp[j])
				}
				continue
			}
			if lastDimension {
				*modified = append(*modified, cp[j])
			}
		}

		for _, index := range toDelete {
			cp.deleteAt(index)
		}

		return cp
	}

	return cp
}

// NewEntries will return a reused list of entries.
func NewEntries() Entries {
	return entriesPool.Get().(Entries)
}
