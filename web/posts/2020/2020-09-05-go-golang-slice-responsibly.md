<h3>How Slices Work</h3>

<p>
To be able to understand how slices works, we first need to have a good understanding of how arrays work in Go, and you can check out <a href="https://blog.golang.org/slices#TOC_2.">this</a> informative description to gain that understanding.
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

<h3>How Slicing Works</h3>

<h3>Modifying a Sliced Slice Modifies the Original Slice</h3>

<!-- https://play.golang.org/p/aCoF1zJf18y -->

<p>
As we have just seen by looking at the internals of how slicing works, the new slice that's returned by slicing an existing slice is still referring to the same backing array as the original sliced slice. This introduces a very interesting implication that modifying data on the indices of the newly sliced slice will also cause the same modification on the original slice. This can actually cause very hard to track down bugs, and the below is showing how this can happen:
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
In this given example, the issue might be obvious to us. However, this unobvious behavior of how slicing works underneath (to be fair, for the right performance reasons) can make some issues more obfuscated when the slicing and modification is done different places. For instance, with the below example (which you can also see <a href="https://play.golang.org/p/2OUIE7s69mi">here</a>), we can see that <code>Result</code> method on the <code>race</code> instance is not returning the expected result anymore due to the modifications done to the slice returned by the <code>Top10Finishers</code> method, because <code>sort.Strings</code> call modified the array which is actually backing the both slices
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
There is no one size fit all the problems solution here. It will really depend on your usage, and what type of contract you are exposing from your defined type. If you are after creating a domain model encapsulation where you don't want to allow unmodified access to the state of that model, you should instead make a copy of the slice that you want to return, with the cost of extra time and space complexity you are introducing. The following code shows the only modification we would have done to the above example to make this work:
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

https://play.golang.org/p/iDnB1ZrG4lj

<p></p>

<p><pre></pre></p>

https://play.golang.org/p/g3nfRo8kXll

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
</ul>