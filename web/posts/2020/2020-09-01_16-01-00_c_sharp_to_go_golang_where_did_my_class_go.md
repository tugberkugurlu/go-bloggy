---
id: 01EH2GPB0WW8PPK6M01B5MR3N5
title: "C# to Go (Golang): Where Did My Class Go?"
abstract: In C#, classes allow us to model our domain, and attach behaviour to the that model. It also allows us to make the decision upfront in terms of the storage and flowing behaviour for our model, by making it a reference type. How does this work in Go programming language? What are the ways for us to effectively define the behaviours for a particular model, and confidently flow it throughout the application lifecycle?
created_at: 2020-09-01 16:01:00.0000000 +0000 UTC
tags:
- Go
- Golang
- C#
slugs:
- c-sharp-to-go-golang-where-did-my-class-go
---

<p>
<pre>
</pre>
</p>
https://docs.microsoft.com/en-us/dotnet/csharp/programming-guide/classes-and-structs/classes
 (which you can see <a href="https://play.golang.org/p/bB5ndRYD3hz">here</a>)

<h3>Fundamentals</h3>

<p>
Let's understand how fundamentals of Go first which will help us hold strong assumptions about how Go code executes throughout this post as someone who has a C# knowledge. First of all, in Go, a way for us to structure a domain model is through <a href="A struct is a collection of fields.">a struct</a>, which holds a collection of fields.
</p>

<p>
<pre>
type Person struct {
    Name      string
    Surname   string
}
</pre>
</p>

<p>
Another fundamental aspect of Go code is that <a href="https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go">it's always passed by value</a> when you pass any value to a function. However, you can point to the address of a struct instance through <a href="https://dave.cheney.net/2017/04/26/understand-go-pointers-in-less-than-800-words-or-your-money-back">a pointer</a>, which allows you to mutate the struct fields across functions. The below code for example will print out <code>Bob</code> to the console, even if the name is updated inside the <code>mutate</code> function:
</p>

<p>
<pre>
package main

import (
	"fmt"
)

type Person struct {
    Name      string
    Surname   string
}

func main() {
	p := Person{Name: "Bob", Surname: "Smith"}
	mutate(p)
	fmt.Println(p.Name)
}

func mutate(val Person) {
	val.Name = "Alice"
}
</pre>
</p>

<p>
However, the below code will print out <code>Alice</code> to the console, as we are passing the struct value as pointer:
</p>

<p>
<pre>
package main

import (
	"fmt"
)

type Person struct {
    Name      string
    Surname   string
}

func main() {
	p := Person{Name: "Bob", Surname: "Smith"}
	mutate(&p)
	fmt.Println(p.Name)
}

func mutate(val *Person) {
	val.Name = "Alice"
}
</pre>
</p>

<p>
However, it's important to emphasize here that this is not equivalent of <a href="https://docs.microsoft.com/en-us/dotnet/csharp/language-reference/keywords/ref"><code>ref</code> parameter in C#</a> which allows a value to be passed by reference. The reason for that is the the lack of pass-by-reference semantics in Go, which is explained perfectly in the referenced Dave Cheney's post. In other words, a pointer won't allow you to mutate the value of the variable you are passing. If we execute the below code examples, C# version will print <code>0</code>, whereas the Go version will still print <code>1</code> to the console.
</p>

<p>
<pre>
class Program
{
    static void Main(string[] args)
    {
        var foo = 1;
        mutate(ref foo);
        Console.WriteLine(foo.ToString());
    }

    static void mutate(ref int val) 
    {
        val = 0;
    }
}
</pre>
</p>

<p>
<pre>
package main

import (
	"fmt"
)

func main() {
	foo := 1
	mutate(&foo)
	fmt.Println(foo)
}

func mutate(val *int) {
	newVal := 0
	val = &newVal
}
</pre>
</p>

<p>
The last thing to mention for us is <a href="https://tour.golang.org/methods/1">methods</a>, which is very similar to a <a href="https://tour.golang.org/basics/4">function declaration</a> but it's with a special receiver argument. The methods will allow you to attach behavior to your structs, as well as also enabling you to mutate its internal state. The below is a good example of attaching three methods to the <code>Person</code> struct, where <code>Earn</code> and <code>Spend</code> mutates the internal <code>wallet</code> field, and <code>Balance</code> is used as a reader of the internal state.
</p>

<p>
<pre>
package main

import (
	"errors"
	"fmt"
)

type Person struct {
	Name    string
	Surname string
	wallet  float64
}

func (p *Person) Earn(val float64) {
	p.wallet += val
}

func (p *Person) Spend(val float64) error {
	if p.wallet < val {
		return errors.New("the balance is not enough")
	}
	p.wallet -= val
	return nil
}

func (p *Person) Balance() float64 {
	return p.wallet
}

func main() {
	p := Person{Name: "Alice", Surname: "Smith"}
	p.Earn(40.5)
	err := p.Spend(10)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Balance: %f\n", p.Balance())

	err = p.Spend(40)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
</pre>
</p>

<h3>Example Scenario</h3>

<p>
Let's define an example domain scenario for us to model. The domain model we will be defining is going to represent a "Trip", which has the following aspects:
</p>

<ul>
<li>Name</li>
<li>Start and end date</li>
<li>Members (which are defined by their name only for simplicity purposes)</li>
<li>Locations to visit during the trip, and each location will carry its name, latitude and longitude</li>
</ul>

<p>
Besides these aspects, the model will have the following behaviors:
</p>

<ul>
<li>Add a new member</li>
<li>Remove an existing member</li>
<li>Check the length of the trip</li>
<li>Check whether a trip overlaps with a given time</li>
<li>Check whether a member is in the trip or not</li>
<li>Retrieve the closest locations to visit in the order of straight line distance based on a given target location</li>
</ul>

<h3>Poor Man's Domain Modeling in Go</h3>

<p>

</p>

<p>
<pre>
package main

import (
	"sort"
	"time"

	"github.com/umahmood/haversine"
)

type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

type LocationWithDistance struct {
	Destination  Location
	DistanceInKM float64
}

type Trip struct {
	Name      string
	TripStart time.Time
	TripEnd   time.Time
	Members   map[string]bool
	Locations []Location
}

func (t *Trip) Length() time.Duration {
	return t.TripEnd.Sub(t.TripStart)
}

func (t *Trip) IsDuringTrip(at time.Time) bool {
	return t.TripStart.Before(at) && t.TripEnd.After(at)
}

func (t *Trip) IsMemberInTrip(name string) bool {
	return t.Members[name]
}

func (t *Trip) GetClosestVisitedLocations(target Location) []LocationWithDistance {
	var distances []LocationWithDistance
	for _, loc := range t.Locations {
		lh := haversine.Coord{Lat: loc.Latitude, Lon: loc.Longitude}
		th := haversine.Coord{Lat: target.Latitude, Lon: target.Longitude}
		_, km := haversine.Distance(lh, th)
		distances = append(distances, LocationWithDistance{
			Destination:  loc,
			DistanceInKM: km,
		})
	}
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].DistanceInKM < distances[j].DistanceInKM
	})
	return distances
}
</pre>
</p>

<h3>Main Problems with This Implementation</h3>

<p>
If we were to list the problems with this approach to modeling our domain, the below would be the key ones based on my experience with Go in large codebases:
</p>

<ul>
<li>It's easy for the consumer to construct this type in an invalid state (i.e. lack of constructor)</li>
<li>The decision of whether to flow the type instance as a value or a pointer is left to the consumer, which is not really useful and can cause unexpected issues</li>
<li>The internal state of the model is exposed to the consumer, which means that it can be modified in an uncontrolled way</li>
<li>Hard to enforce consumers to act on changes to the public signature (i.e. compiler won't fail when you add a new field to the struct)</li>
</ul>

<p>

</p>

<h3>Encapsulating the Construction</h3>

<h3>Scoping the Model Under a Package</h3>

<p>
This will allow us to hide the domain model internals from the outside world by internalizing them to package itself, and only exposing the behavior and data we want. This will indeed require more effort but it means that we will be stopping the uncontrolled mutation of the model's internal state, which is priceless in a large codebase with many contributors who have different level of Go experience.
</p>