---
id: 01EGWWJYXSEWRZQ709XW24NP3Z
title: Implementing OrderedMap in Go 2.0 by Using Generics with Delete Operation in O(1) Time Complexity
abstract: I stumbled upon rocketlaunchr.cloud's great post on using generics in Go 2, and the post shows how you can implement ordered maps. The post is very informative, and shows you how powerful Generics will be for Go. However, I noticed a performance issues with the implementation of Delete operation on the ordered map, which can be a significant performance hit where there is a need to store large collections while making use of the Delete operation frequently. In this post, I want to implement that part with a Doubly Linked List data structure to reduce its time complexity from O(N) to O(1).
created_at: 2020-08-29 11:33:00.0000000 +0000 UTC
tags:
- Go
- Golang
- Data Structures
- Algorithms
slugs:
- implementing-ordered-map-in-go-2-0-by-using-generics-with-delete-operation-in-o-1-time-complexity
- implementing-ordered-map-in-go-2-0-by-using-generics-with-delete-operation-in-0-1-time-complexity
---

<p>
Probably the most sought after feature of Go programming language, Generics, is on its way and is expected to land with v2. You can check out the proposal <a href="https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md">here</a>, and have a play with it in <a href="https://go2goplay.golang.org/">Go playground for v2</a>. I stumbled upon <a href="https://medium.com/@rocketlaunchr.cloud">rocketlaunchr.cloud</a>'s great <a href="https://medium.com/swlh/ordered-maps-for-go-using-generics-875ef3816c71">post on using generics in Go 2</a>, and the post shows how you can implement ordered maps. The post is very informative, and shows you how powerful Generics will be for Go.
</p>

<p>
However, I noticed a performance issues with the implementation of <code>Delete</code> operation on the <code>OrderedMap</code> struct, and in this post, I want to show a much better implementation in terms of time complexity with a Doubly Linked List data structure, and show the impact of the change.
</p>

<h3>The Problem</h3>

<p>
To summarize the current approach in rocketlaunchr.cloud's post, it essentially exposes the below signature: 
</p>

<p>
<pre>
type OrderedMap[type K comparable, V any] struct {
	store map[K]V
	keys  []K
}

func (o *OrderedMap[K, V]) Get(key K) (V, bool) {
    // ...
}

func (o *OrderedMap[K, V]) Set(key K, val V) {
    // ...
}

func (o *OrderedMap[K, V]) Delete(key K) {
    // ...
}

func (o *OrderedMap[K, V]) Iterator() func() (*int, *K, V) {
    // ...
}
</pre>
</p>

<p>
The implementation also gives the FIFO guarantee through the iterator and maintaining the order of the list even after the delete operation (e.g. if the map has <code>1,2,3,4</code>, and <code>3</code> is then deleted, the iterator will output the data with the following order <code>1,2,4</code>.).
</p>

<p>
The performance problem with the implementation is with the <code>Delete</code> method, which has the below implementation:
</p>

<p>
<pre>
func (o *OrderedMap[K, V]) Delete(key K) {
	delete(o.store, key)

	// Find key in slice
	var idx *int

	for i, val := range o.keys {
		if val == key {
			idx = &[]int{i}[0]
			break
		}
	}
	if idx != nil {
		o.keys = append(o.keys[:*idx], o.keys[*idx+1:]...)
	}
}
</pre>
</p>

<p>
The implementation here is iterating over the entire keys slice to perform the delete operation, which has <code>O(N)</code> time complexity and this can be a significant performance hit where there is a need to store large collections while making use of the Delete operation frequently. We can also observe that we are performing a shift in the keys slice as well, with the below operation
</p>

<p>
<pre>
if idx != nil {
    o.keys = append(o.keys[:*idx], o.keys[*idx+1:]...)
}
</pre>
</p>

<p>
<code>o.keys[:*idx]</code>, <code>o.keys[*idx+1:]</code> and <code>append</code> all here have its own time complexity, and <a href="https://stackoverflow.com/questions/15702419/append-complexity">depending on how much the backing array needs to grow</a>, the complexity of this operation can grow. That said, I have to admit that the language features here make it hard to reason about the exact time complexity of each operation here.
</p>

<p>
In the post, I realized that rocketlaunchr.cloud actually calls for an action for readers to address this known issue with the below quote for this :)
</p>

<blockquote>
<p>
Currently, when you delete a key-value pair, you need to iterate over the keys slice to find the index of the key you want to delete. You can use another map that associates the key to the index in the keys slice. I’ll leave it to the reader to implement.
</p>
</blockquote>

<p>
Thinking about this now and it won't actually be enough for us to store the index as it will require us to perform either shift the slice where the keys are stored, or keep the storage for deleted keys and skip them during iteration. Both has its own disadvantages where the first option will have time complexity hit per each delete operation, whereas the second one will both increase the time complexity of iterator as well as increasing the space complexity.
</p>

<h3>Doubly Linked List Data Structure to Rescue</h3>

<p>
It's actually possible to reduce the time complexity of the <code>Delete</code> operation to <code>O(1)</code> by changing the way how we store the data within the implementation of <code>OrderedMap</code> struct, without increasing the time complexity of other operations, and needing to change any of its public signature. We can do this by storing the ordered data in a <a href="https://en.wikipedia.org/wiki/Doubly_linked_list">Doubly Linked List</a> data structure, and storing each node in the map as the value, instead of the raw value.
</p>

<p>
<img src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/doubly%20linked%20list%20and%20map.jpg" alt="doubly linked list and map storage" />
</p>

<p>
Implementation of a Doubly Linked List data structure should be fairly straight forward to implement. However, Go already has one under <a href="https://golang.org/pkg/container/list/"><code>container/list</code> package</a>. The only caveat for us with this one is that there is no generic version of this in Go 2 at the moment, as far as I am aware. That said, we actually don't need a generic version of this since we will use this internally within the scope of <code>OrderedMap</code>, and we can store the value as <code>interface{}</code> instead.
</p>

<h3>Implementation</h3>

<p>
Let's look at the implementation, going through it step by step. Below we can see how we are changing the internal state storage constructs of the <code>OrderedMap</code> struct as well as its construction:
</p>

<p>
<pre>
type OrderedMap[type K comparable, V any] struct {
	store map[K]*list.Element
	keys  *list.List
}

func NewOrderedMap[type K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		store: map[K]*list.Element{},
		keys:  list.New(),
	}
}
</pre>
</p>

<p>
The biggest changes to notice here are:
</p>

<ul>
<li><code>keys</code> type has changed from <code>[]K</code> to <code>*list.List</code></li>
<li>map type has changed from <code>map[K]V</code> to <code>map[K]*list.Element</code>, where we are now storing the doubly linked list node, instead of the raw value</li>
</ul>

<p>
These changes to the internal storage will impact the implementation of all methods, but without increasing the time complexity, and needing to change the public signature of each method. Below, for example, is how the <code>Set</code> and <code>Get</code> method implementations have changed:
</p>

<p>
<pre>
type keyValueHolder[type K comparable, V any] struct {
	key K
	value V
}

func (o *OrderedMap[K, V]) Set(key K, val V) {
	var e *list.Element
	if _, exists := o.store[key]; !exists {
		e = o.keys.PushBack(keyValueHolder[K, V]{
			key: key,
			value: val,
		})
	} else {
		e = o.store[key]
		e.Value = keyValueHolder[K, V]{
			key: key,
			value: val,
		}
	}
	o.store[key] = e
}

func (o *OrderedMap[K, V]) Get(key K) (V, bool) {
	val, exists := o.store[key]
	if !exists {
		return *new(V), false
	}
	return val.Value.(keyValueHolder[K, V]).value, true
}

func (o *OrderedMap[K, V]) Iterator() func() (*int, *K, V) {
	e := o.keys.Front()
	j := 0
	return func() (_ *int, _ *K, _ V) {
		if e == nil {
			return
		}

		keyVal := e.Value.(keyValueHolder[K, V])
		j++
		e = e.Next()

		return func() *int { v := j-1; return &v }(), &keyVal.key, keyVal.value
	}
}
</pre>
</p>

<p>
In a nutshell, the difference here is that we are storing the doubly linked list node as the value in the map, and we store both the key and value as the value of the doubly linked list node through <code>keyValueHolder</code> struct. This obviously impacts how we set and get the data, but the time complexity of both the methods stay as <code>O(1)</code>. We will actually observe the biggest change with the <code>Delete</code> method here:
</p>

<p>
<pre>
func (o *OrderedMap[K, V]) Delete(key K) {
	e, exists := o.store[key]
	if !exists {
		return
	}

	o.keys.Remove(e)

	delete(o.store, key)
}
</pre>
</p>

<p>
If we were to break this down, here is what we are doing here:
</p>

<ul>
<li>Accessing the doubly linked list node from the map first, based on the given key. This is <code>O(1)</code> in terms of time complexity.</li>
<li>Calling <code>Remove</code> on the doubly linked list by passing the found element. This is <code>O(1)</code> in terms of time complexity.</li>
<li>Deleting node from the map, based on the given key. This is also <code>O(1)</code> in terms of time complexity.</li>
</ul>

<p>
One more thing to unwrap here is how come calling <code>Remove</code> on the doubly linked list is <code>O(1)</code> in terms of time complexity, and showing its implementation might shed some light on the rationale:
</p>

<p>
<pre>
// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List) Remove(e *Element) interface{} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *List) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}
</pre>
</p>

<p>
All that this implementation does is to appropriately remove the links from the node that we want to remove, and wiring up its next node with its previous node (if they exist). You can also see this implementation <a href="https://github.com/golang/go/blob/4fc3896e7933e31822caa50e024d4e139befc75f/src/container/list/list.go#L107-L144">here</a> in Go source code.
</p>

<p>
You can also see the whole implementation working <a href="https://go2goplay.golang.org/p/UJZsQnPRmRh ">here in Go 2 Playground</a>:
</p>

<p>
<pre>
0 1 string1 is a string
1 2 string2 is a string
2 4 string4 is a string

Program exited.
</pre>
</p>

<h3>Benchmarking to Show the Impact</h3>

<p>
To be able to understand the improvement we have made here, we can run a benchmark between the two <code>Delete</code> implementations with <a href="https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go">Go's built-in benchmark tooling</a>. For this, we will go with the below benchmark setup:
</p>

<ul>
<li>Sequential input of an array with 1M items</li>
<li>Seeding the <code>OrderedMap</code> with these 1M items</li>
<li>Deleting a random key for 100K times</li>
</ul>

<p>
One caveat with the benchmark here is that we will run it with non-generic implementations since I wasn't unable to find a way to install Go 2 on my machine to be able to run a benchmark, and I don't believe it's possible to run one with Go playground. That said, this shouldn't change anything for us to be able to understand the difference between two <code>Delete</code> implementations as the logic will stay the same. The only change will be that the usage of the type will not be obvious due to lack of strongly typed signature. I have put the <code>OrderedMap</code> with original <code>Delete</code> implementation <a href="https://github.com/tugberkugurlu/algos-go/blob/a8c519e7c60e5fb838a3e5bd8a97893f96ad834b/ordered-map-generics/ordered_map.go#L27-L42">here</a>, and with the improved one <a href="https://github.com/tugberkugurlu/algos-go/blob/a8c519e7c60e5fb838a3e5bd8a97893f96ad834b/ordered-map-generics/ordered_map_linkedlist_based_delete.go#L47-L56">here</a>.
</p>

<p>
The code for the benchmark itself is as below (which you can also find <a href="https://github.com/tugberkugurlu/algos-go/blob/a8c519e7c60e5fb838a3e5bd8a97893f96ad834b/ordered-map-generics/main_test.go#L136-L160">here</a>):
</p>

<p>
<pre>
package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkOrderedMapLinkedListBasedDelete(b *testing.B) {
	for n := 0; n < b.N; n++ {
		seedCount := 1000000
		m := NewLinkedListBasedOrderedMap()
		for i := 1; i <= seedCount; i++ {
			m.Set(i, fmt.Sprintf("string%d", i))
		}
		for i := 0; i < 100000; i++ {
			m.Delete(rand.Intn(seedCount-1) + 1)
		}
	}
}

func BenchmarkOrderedMapSliceBasedDelete(b *testing.B) {
	for n := 0; n < b.N; n++ {
		seedCount := 1000000
		m := NewOrderedMap()
		for i := 1; i <= seedCount; i++ {
			m.Set(i, fmt.Sprintf("string%d", i))
		}
		for i := 0; i < 100000; i++ {
			m.Delete(rand.Intn(seedCount-1) + 1)
		}
	}
}
</pre>
</p>

<p>
We can run this benchmark with <a href="https://golang.org/pkg/cmd/go/internal/test/"><code>go test</code> command</a>. However, note that the exact time that takes to run each function is not significant here, since it will depend on the machine spec, etc. Also, in the test, we are seeding the map implementations for each run, which means extra time we are adding to the exact time to run the function. That said, what's important to observe here is the difference in time between two functions:
</p>

<p>
<pre>
➜  ordered-map-generics git:(master) ✗ go test -bench=BenchmarkOrderedMap -benchtime=30s
goos: darwin
goarch: amd64
pkg: github.com/tugberkugurlu/algos-go/ordered-map-generics
BenchmarkOrderedMapLinkedListBasedDelete-4   	      42	 724271767 ns/op
BenchmarkOrderedMapSliceBasedDelete-4        	       1	54509739278 ns/op
PASS
ok  	github.com/tugberkugurlu/algos-go/ordered-map-generics	86.361s
</pre>
</p>

<p>
Wow, we can see the orders of magnitude improvement here, and it's very rewarding to see the impact.
</p>

<h3>Conclusion</h3>

<p>
The biggest take away from this post I believe is how Generics will change the way we can implement more complex and powerful data structures, and allow them to be reused more effectively by truly taking advantage of strongly-typed nature of Go, as well as how important it's to use the correct data structures for the expected usage of own types. As a side note, I thoroughly enjoyed writing this post, as I was able to feel the power of having Generics available to use in Go!
</p>