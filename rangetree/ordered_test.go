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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedAdd(t *testing.T) {
	nodes := make(orderedNodes, 0)

	e1 := constructMockEntry(1, 4)
	e2 := constructMockEntry(2, 1)

	nodes.add(1, e1)
	nodes.add(1, e2)

	assert.Equal(t, Entries{e2, e1}, nodes.toEntries())
}

func TestOrderedDelete(t *testing.T) {
	nodes := make(orderedNodes, 0)

	e1 := constructMockEntry(1, 4)
	e2 := constructMockEntry(2, 1)

	nodes.add(1, e1)
	nodes.add(1, e2)

	nodes.delete(1)

	assert.Equal(t, Entries{e1}, nodes.toEntries())

	nodes.delete(4)

	assert.Len(t, nodes, 0)
}

func TestApply(t *testing.T) {
	ns := make(orderedNodes, 0)

	e1 := constructMockEntry(1, 4)
	e2 := constructMockEntry(2, 1)

	ns.add(1, e1)
	ns.add(1, e2)

	results := make(Entries, 0, 2)

	ns.apply(1, 2, func(n *node) bool {
		results = append(results, n.entry)
		return true
	})

	assert.Equal(t, Entries{e2}, results)

	results = results[:0]

	ns.apply(0, 1, func(n *node) bool {
		results = append(results, n.entry)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(2, 4, func(n *node) bool {
		results = append(results, n.entry)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(4, 5, func(n *node) bool {
		results = append(results, n.entry)
		return true
	})

	assert.Equal(t, Entries{e1}, results)
	results = results[:0]

	ns.apply(0, 5, func(n *node) bool {
		results = append(results, n.entry)
		return true
	})

	assert.Equal(t, Entries{e2, e1}, results)
	results = results[:0]

	ns.apply(5, 10, func(n *node) bool {
		results = append(results, n.entry)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(0, 100, func(n *node) bool {
		results = append(results, n.entry)
		return false
	})

	assert.Equal(t, Entries{e2}, results)
}

func TestInsertDelete(t *testing.T) {
	ns := make(orderedNodes, 0)

	e1 := constructMockEntry(1, 4)
	e2 := constructMockEntry(2, 1)
	e3 := constructMockEntry(3, 2)

	ns.add(1, e1)
	ns.add(1, e2)
	ns.add(1, e3)

	modified := make(Entries, 0, 1)
	deleted := make(Entries, 0, 1)

	ns.insert(2, 2, 2, 0, -5, &modified, &deleted)

	assert.Len(t, ns, 0)
	assert.Equal(t, Entries{e2, e3, e1}, deleted)
}
