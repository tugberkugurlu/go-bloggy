---
id: 01EJ1FYNB4CB3B13JVAEHB5908
title: Working with Slices in Go (Golang) - Slicing Responsibly and Understanding How append, copy and Slicing Syntax Work
abstract: Slices in Go programming language gives us greater flexibility over working with arrays, but this flexibility comes with some trade-offs. Go optimizes towards the most frequent uses cases, and you are often working with them correctly. However, in certain cases, some of the implicit hidden behaviors of Go slices can create unclear issues which can be hard to diagnose at the first place. In this post, we will go over some of the implicit behaviors while recapping how slices work in Go in general.
created_at: 2020-09-12 16:55:00.0000000 +0000 UTC
tags:
- Go
- Golang
- Data Structures
- Algorithms
slugs:
- working-with-slices-in-go-golang-slicing-responsibly-and-understanding-how-append-copy-and-slicing-syntax-work
---

<p>
Go programming language, has two fundamental types at the language level to enable working with numbered sequence of elements: <a href="https://golang.org/ref/spec#Array_types">array</a> and <a href="https://golang.org/ref/spec#Slice_types">slice</a>. At the syntax level, they may look like the same but they are fundamentally very different in terms of their behavior. The most critical fundamental differences are:
</p>

<ul>
<li>Size of the array is fixed and determined at the construction (as you may expect). However, slices can dynamically grow in size (you may wonder how? We will touch on this soon, be patient!)</li>
<li>An array with a specific length is a distinct type based on its length (check <a href="https://play.golang.org/p/gmF99PSNhiX">this</a> out). Whereas the slice can be represented as one type (e.g. <code>[]int</code>)</li>
<li>The in-memory representation of an array type is values laid out sequentially. A slice is a descriptor of an array segment. It consists of a pointer to the array, the length of the segment, and its capacity (we will shortly see what this actually means).</li>
<li>Go's arrays are values, which means that the entire content of the array will be copied when you start passing it around. Slices, on the other hand, a pointer to the underlying along with the length of the segment. So,  when we started passing around a slice, it creates a new slice value that points to the original array, which will be much cheaper to pass around.</li>
</ul>

<p>
Above points highlight some characteristics of slices and how they differ from arrays, but these are mostly differences in terms of how they are structured. More interesting and unobvious differences of slices are around their behaviors around manipulations.
</p>

<h3>How Slices Work</h3>

<p>
To be able to understand how slices works, we first need to have a good understanding of how arrays work in Go, and you can check out <a href="https://blog.golang.org/slices#TOC_2.">this</a> informative description to gain more understanding on arrays than the above summary.
</p>

<p>
<a href="https://blog.golang.org/slices#TOC_3.">Slices</a> are the constructs in Go which give us the flexibility to work with dynamically sized collections. A slice is an abstraction of an array, and it points to a contiguous section of an array stored separately from the slice variable itself. Internally, a slice is a descriptor which holds the following values:
</p>

<ul>
<li>pointer to the backing array</li>
<li>the length of the segment it's referring to</li>
<li>its capacity (the maximum length of the segment)</li>
</ul>

<p>There are various ways how you can define a slice in Go, and all of the following ways leads to the same outcome: a slice with a zero length and capacity</p>

<p><pre>
package main

import (
	"fmt"
)

func main() {
	var a []int 
	b := []int{}
	c := make([]int, 0)
	fmt.Printf("a: %v, len %d, cap: %d\n", a, len(a), cap(a))
	fmt.Printf("b: %v, len %d, cap: %d\n", b, len(b), cap(b))
	fmt.Printf("c: %v, len %d, cap: %d\n", c, len(c), cap(c))
}
</pre></p>

<p><pre>
a: [], len 0, cap: 0
b: [], len 0, cap: 0
c: [], len 0, cap: 0
</pre></p>

<p>
You can also initialize a slice with seed values, and the length of the values here will also be the capacity of the backing array:
</p>

<p><pre>
func main() {
	a := []int{1,2,3}
	fmt.Printf("a: %v, len %d, cap: %d\n", a, len(a), cap(a))
}
</pre></p>

<p><pre>
a: [1 2 3], len 3, cap: 3
</pre></p>

<p>
In case you know the maximum capacity that a slice can grow to, it's best to initialize the slice by hinting the capacity so that you don't have to grow the backing array as you add new values to slice (which we will see how to in the next section). You can do so by passing it to the <a href="https://golang.org/pkg/builtin/#make"><code>make</code> builtin function</a>.
</p>

<p><pre>
func main() {
	a := make([]int, 0, 10)
	fmt.Printf("a: %v, len %d, cap: %d\n", a, len(a), cap(a))
}
</pre></p>

<p><pre>
a: [], len 0, cap: 10
</pre></p>

<p>
This still doesn't mean that you can access the backing array freely by index, as the length is still <code>0</code>. If you attempt to do so, you will get "index out of range" runtime error.
</p>

<h3>How append and copy Works</h3>

<p>
When we want to add a new value to an existing slice which will mean growing its length, we can use <a href="https://golang.org/pkg/builtin/#append">append</a>, which is a built-in and <a href="https://golang.org/ref/spec#Function_types">variadic</a> function. This function appends elements to the end of a slice, and returns the updated slice.
</p>

<!-- https://play.golang.org/p/YvzmKIIfY0U -->

<p><pre>
func main() {
	var result []int
	for i := 0; i < 10; i++ {
		if i % 2 == 0 {
			result = append(result, i)
		}
	}
	fmt.Println(result)
}
</pre></p>

<p>
As you may expect, this prints <code>[0 2 4 6 8]</code> to the console as a result. However, it's not clear here what exactly is happening underneath as a result of invocation of the <code>append</code> function, and what the time complexity of the call is. When we run the below code, things will be a bit more clear to us:
</p>

<!-- https://play.golang.org/p/uQQFBbJWhAS -->

<p><pre>
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	var result []int
	for i := 0; i < 10; i++ {
		if i % 2 == 0 {
			fmt.Printf("appending '%d': %s\n", i, getSliceHeader(&result))
			result = append(result, i)
			fmt.Printf("appended '%d':  %s\n", i, getSliceHeader(&result))
		}
	}
	fmt.Println(result)
}

// https://stackoverflow.com/a/54196005/463785
func getSliceHeader(slice *[]int) string {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(slice))
	return fmt.Sprintf("%+v", sh)
}
</pre></p>

<p><pre>
appending '0': &{Data:0 Len:0 Cap:0}
appended '0':  &{Data:824633901184 Len:1 Cap:1}
appending '2': &{Data:824633901184 Len:1 Cap:1}
appended '2':  &{Data:824633901296 Len:2 Cap:2}
appending '4': &{Data:824633901296 Len:2 Cap:2}
appended '4':  &{Data:824633803136 Len:3 Cap:4}
appending '6': &{Data:824633803136 Len:3 Cap:4}
appended '6':  &{Data:824633803136 Len:4 Cap:4}
appending '8': &{Data:824633803136 Len:4 Cap:4}
appended '8':  &{Data:824634228800 Len:5 Cap:8}
[0 2 4 6 8]
</pre></p>

<p>
We can extract the following facts from this result:
</p>

<ul>
<li><code>nil</code> slice starts off with empty capacity, nothing surprising with that</li>
<li>The capacity of the slice doubles while attempting to append a new item when its capacity and length are equal</li>
<li>When the capacity is doubled, we can also observe that the pointer to the backing array (i.e. the <code>Data</code> field value of <code>reflect.SliceHeader</code> struct) changes</li>
</ul>

<p>
In summary, it's a fair to assume from these facts that the content of the backing array of the slice is copied into a new array which has double capacity than the itself when it's being attempted to append a new item to it while its capacity is full. It should go without saying that the implementation is a bit more complicated than this as you may expect, and <a href="https://medium.com/vendasta/golang-the-time-complexity-of-append-2177dcfb6bad">this post</a> from <a href="https://medium.com/@glucn">Gary Lu</a> does a good job on explaining the implementation details in more details. You can also check out the <a href="https://github.com/golang/go/blob/b3ef90ec7304a28b89f616ced20b09f56be30cc4/src/runtime/slice.go#L125-L240">growSlice</a> function which is used by the compiler generated code to grow the capacity of the slice when needed.
</p>

<p>
In a nutshell, this is not a good news to us since we are doing too much more work than we desired to do. In this cases, initializing the array with the make built-in function is a far better option with a capacity hint based on the max capacity that the slice can grow to:
</p>

<!-- https://play.golang.org/p/rAE7KWP9yfk -->

<p><pre>
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	maxValue := 10
	result := make([]int, 0, maxValue)
	for i := 0; i < maxValue; i++ {
		if i % 2 == 0 {
			fmt.Printf("appending '%d': %s\n", i, getSliceHeader(&result))
			result = append(result, i)
			fmt.Printf("appended '%d':  %s\n", i, getSliceHeader(&result))
		}
	}
	fmt.Println(result)
}

// https://stackoverflow.com/a/54196005/463785
func getSliceHeader(slice *[]int) string {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(slice))
	return fmt.Sprintf("%+v", sh)
}
</pre></p>

<p><pre>
appending '0': &{Data:824633794640 Len:0 Cap:10}
appended '0':  &{Data:824633794640 Len:1 Cap:10}
appending '2': &{Data:824633794640 Len:1 Cap:10}
appended '2':  &{Data:824633794640 Len:2 Cap:10}
appending '4': &{Data:824633794640 Len:2 Cap:10}
appended '4':  &{Data:824633794640 Len:3 Cap:10}
appending '6': &{Data:824633794640 Len:3 Cap:10}
appended '6':  &{Data:824633794640 Len:4 Cap:10}
appending '8': &{Data:824633794640 Len:4 Cap:10}
appended '8':  &{Data:824633794640 Len:5 Cap:10}
[0 2 4 6 8]
</pre></p>

<p>
We can observe from the result here that we have been operating over the same backing array with the size of 10, which means that all the append operations have run in <code>O(1)</code> time.
</p>

<h3>How Slicing Works</h3>



<h3>Modifying a Sliced-slice Modifies the Original Slice</h3>

<!-- https://play.golang.org/p/aCoF1zJf18y -->

<p>
As we have just seen, <code>append</code> function appends elements to the end of a slice, and returns the updated slice. This can sort of gives you the impression that the append function is a pure function, and doesn't modify your state. However, by looking at the internals of how slicing works, we have seen that the new slice that's returned by slicing an existing slice is still referring to the same backing array as the original sliced-slice. This introduces a very interesting implication that modifying data on the indices of the newly sliced slice will also cause the same modification on the original slice. This can actually cause very hard to track down bugs, and the below is showing how this can happen:
</p>

<p><pre>
package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3, 4, 5}
	b := a[2:4]
	b[0] = 10
	fmt.Println(b)
	fmt.Println(a)
}
</pre></p>

<p><pre>
[10 4]
[1 2 10 4 5]
</pre></p>

<!-- https://play.golang.org/p/AZR2Qbucmpy -->

<p>
In this given example, the issue might be obvious to us. However, this unobvious behavior of how slicing works underneath (to be fair, for the right performance reasons) can make some issues more obfuscated when the slicing and modification is done different places. For instance, with the below example (which you can also see <a href="https://play.golang.org/p/2OUIE7s69mi">here</a>), we can see that <code>Result</code> method on the <code>race</code> instance is not returning the expected result anymore due to the modifications done to the slice returned by the <code>Top10Finishers</code> method, because <code>sort.Strings</code> call modified the array which is actually backing the both slices.
</p>

<p><pre>
package main

import (
	"fmt"
	"play.ground/race"
	"sort"
)

func main() {
	belgium2020Race := race.New("Belgian", []string{
		"Hamilton", "Bottas", "Verstappen", "Ricciardo", "Ocon",
		"Albon", "Norris", "Gasly", "Stroll", "Perez",
		"Kvyat", "Räikkönen", "Vettel", "Leclerc", "Grosjean",
		"Latifi", "Magnussen", "Giovinazzi", "Russell", "Sainz",
	})
	top10Finishers := belgium2020Race.Top10Finishers()
	sort.Strings(top10Finishers)
	fmt.Printf("%s GP top 10 finishers, in alphabetical order: %v\n", belgium2020Race.Name(), top10Finishers)
	fmt.Printf("%s GP result: %v\n", belgium2020Race.Name(), belgium2020Race.Result())
}

-- go.mod --
module play.ground

-- race/race.go --
package race

type race struct {
	name   string
	result []string
}

func (r race) Name() string {
	return r.name
}

func (r race) Result() []string {
	return r.result
}

func (r race) Top10Finishers() []string {
	return r.result[:10]
}

func New(name string, result []string) race {
	return race{
		name:   name,
		result: result,
	}
}
</pre></p>

<p><pre>
Belgian GP top 10 finishers: [Hamilton Bottas Verstappen Ricciardo Ocon Albon Norris Gasly Stroll Perez]
Belgian GP top 10 finishers, in alphabetical order: [Albon Bottas Gasly Hamilton Norris Ocon Perez Ricciardo Stroll Verstappen]
Belgian GP result: [Albon Bottas Gasly Hamilton Norris Ocon Perez Ricciardo Stroll Verstappen Kvyat Räikkönen Vettel Leclerc Grosjean Latifi Magnussen Giovinazzi Russell Sainz]
</pre></p>

<p>
There is no one-size-fits-all solution to the the problem here. It will really depend on your usage, and what type of contract you are exposing from your defined type. If you are after creating a domain model encapsulation where you don't want to allow unmodified access to the state of that model, you should instead make a copy of the slice that you want to return, with the cost of extra time and space complexity you are introducing. The following code shows the only modification we would have done to the above example to make this work:
</p>

<p><pre>
func (r race) Top10Finishers() []string {
	top10 := r.result[:10]
	result := make([]string, len(top10))
	copy(result, top10)
	return result
}
</pre></p>

<p>
<a href="https://play.golang.org/p/mwWdeuvx20n">When we execute this version of the implementation</a>, we can now see that the <code>sort.Strings</code> call is not implicitly modifying the original slice:
</p>

<p><pre>
Belgian GP top 10 finishers: [Hamilton Bottas Verstappen Ricciardo Ocon Albon Norris Gasly Stroll Perez]
Belgian GP top 10 finishers, in alphabetical order: [Albon Bottas Gasly Hamilton Norris Ocon Perez Ricciardo Stroll Verstappen]
Belgian GP result: [Hamilton Bottas Verstappen Ricciardo Ocon Albon Norris Gasly Stroll Perez Kvyat Räikkönen Vettel Leclerc Grosjean Latifi Magnussen Giovinazzi Russell Sainz]
</pre></p>

<h3>Calling append on a Sliced Slice May Modify the Original Slice</h3>

<!-- https://play.golang.org/p/iDnB1ZrG4lj -->

<p></p>

<p><pre></pre></p>

<!-- https://play.golang.org/p/g3nfRo8kXll -->

<p></p>

<p><pre></pre></p>

<!-- https://play.golang.org/p/sRh6RpAek65 -->

<p></p>

<p><pre>
package main

import (
	"fmt"
)

func main() {
	a := make([]int, 5, 6)
	copy(a, []int{1, 2, 3, 4, 5})
	fmt.Printf("a: %v, cap: %d\n", a, cap(a))

	b := a[3:]
	fmt.Printf("b: %v, cap: %d\n", b, cap(b))
	b = append(b, 10)
	fmt.Printf("b: %v, cap: %d\n", b, cap(b))

	a = append(a, 20)
	fmt.Printf("a: %v, cap: %d\n", a, cap(a))
	fmt.Printf("b: %v, cap: %d\n", b, cap(b))
}
</pre></p>

<p><pre>
a: [1 2 3 4 5], cap: 6
b: [4 5], cap: 3
b: [4 5 10], cap: 3
a: [1 2 3 4 5 20], cap: 6
b: [4 5 20], cap: 3
</pre></p>

<!-- https://play.golang.org/p/5RaMtP0u-wF -->

<p></p>

<p><pre>
package main

import (
	"fmt"
)

func main() {
	a := make([]int, 5, 6)
	copy(a, []int{1, 2, 3, 4, 5})
	fmt.Printf("a: %v, cap: %d\n", a, cap(a))

	// points to the same backing array as 'a'
	b := a[3:]
	fmt.Printf("b: %v, cap: %d\n", b, cap(b))

	// elements in backing array of a are copied to a new 
    // array to be able to store '40', and the capacity 
    // is grown to be double for future growth.
	a = append(a, 30, 40)
	fmt.Printf("a: %v, cap: %d\n", a, cap(a))

	// b still points to the old backing array of 'a'
	b = append(b, 10)
	fmt.Printf("b: %v, cap: %d\n", b, cap(b))
	fmt.Printf("a: %v, cap: %d\n", a, cap(a))

	// a's backing array is different. 
    // So, this also doesn't have any impact on slice 'b'.
	a = append(a, 20)
	fmt.Printf("a: %v, cap: %d\n", a, cap(a))
	fmt.Printf("b: %v, cap: %d\n", b, cap(b))
}
</pre></p>

<p><pre>
a: [1 2 3 4 5], cap: 6
b: [4 5], cap: 3
a: [1 2 3 4 5 30 40], cap: 12
b: [4 5 10], cap: 3
a: [1 2 3 4 5 30 40], cap: 12
a: [1 2 3 4 5 30 40 20], cap: 12
b: [4 5 10], cap: 3
</pre></p>

<h3>Conclusion</h3>

<p>
Slice type in Go is a powerful construct, giving us flexibility over Go's array type with as minimum performance hit as possible. This flexibility introduced with least performance impact comes with some additional cost of being implicit on the implications of modifications performed on slices, and these implications can be significant depending on the use case while also being very hard to track down. This post sheds some light on some of these, but I encourage you to spend time on understand how slices really works in depth before making use of them in anger. Besides this post, following resources should also help you
</p>

<ul>
<li><a href="https://blog.golang.org/slices-intro">Go Slices: usage and internals</a></li>
<li><a href="https://blog.golang.org/slices">Arrays, slices (and strings): The mechanics of 'append'</a></li>
<li><a href="https://tour.golang.org/moretypes/7">A Tour of Go: Slices</a></li>
<li><a href="https://medium.com/@marty.stepien/arrays-vs-slices-bonanza-in-golang-fa8d32cd2b7c">Arrays vs. slices bonanza</a></li>
</ul>