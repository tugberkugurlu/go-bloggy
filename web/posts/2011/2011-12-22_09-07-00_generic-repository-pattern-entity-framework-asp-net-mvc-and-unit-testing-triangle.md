---
id: 381086f1-18c0-4ad1-bb92-8b926bdc6c12
title: Generic Repository Pattern - Entity Framework, ASP.NET MVC and Unit Testing
  Triangle
abstract: We will see how we can implement Generic Repository Pattern with Entity
  Framework and how to benefit from that.
created_at: 2011-12-22 09:07:00 +0000 UTC
tags:
- ASP.NET MVC
- C#
- DbContext
- Entity Framework
- Unit Testing
slugs:
- generic-repository-pattern-entity-framework-asp-net-mvc-and-unit-testing-triangle
- generic-repository-pattern-entity-framework-asp-net-mvc-and-uni
---

<blockquote>
<h3>IMPORTANT NOTE:</h3>
I have a new blog post about Generic Repository implementation for Entity Framework. Please check it out instead of this one: <a href="http://www.tugberkugurlu.com/archive/clean-better-and-sexier-generic-repository-implementation-for-entity-framework">Clean, Better, and Sexier Generic Repository Implementation for Entity Framework </a></blockquote>
<blockquote>
<p><strong>NOTE:</strong></p>
<p>Entity Framework DbContext Generic Repository Implementation Is Now On Nuget and GitHub:&nbsp;<a href="http://www.tugberkugurlu.com/archive/entity-framework-dbcontext-generic-repository-implementation-is-now-on-nuget-and-github" target="_blank">http://www.tugberkugurlu.com/archive/entity-framework-dbcontext-generic-repository-implementation-is-now-on-nuget-and-github</a><strong><br /></strong></p>
</blockquote>
<p><strong>DRY</strong>: <a target="_blank" href="http://en.wikipedia.org/wiki/Don't_repeat_yourself" title="http://en.wikipedia.org/wiki/Don't_repeat_yourself">Don&rsquo;t repeat yourself</a> which is a <strong>principle</strong> of software development aimed at reducing repetition of information of all kinds, especially useful in multi-tier architectures. That&rsquo;s what <a target="_blank" href="http://en.wikipedia.org" title="http://en.wikipedia.org">Wikipedia</a> says. In my words, if you are writing the same code twice, follow these steps:</p>
<ul>
<li>Step back. </li>
<li>Sit. </li>
<li>Think about it and dwell on that. </li>
</ul>
<p>That&rsquo;s what I have done for repository classes on my DAL projects. I nearly all the time use Entity Framework to reach out my database and I create repositories in order to query and manipulate data inside that database. There are some specific methods which I use for every single repository. As you can assume, those are <em>FindBy</em>, <em>Add</em>, <em>Edit</em>, <em>Delete</em>, <em>Save</em>. Let&rsquo;s see on code what is my story here.</p>
<p><strong>First approach (worst approach)</strong></p>
<p>At first, long time ago, I have been creating all the single methods for each interface. For example below one is one of my repository interfaces:</p>
<blockquote>
<p>I am giving the examples here with EF 4.2 but I was following this approach with EF 4 which does not contain DbContext class.</p>
</blockquote>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">interface</span> IFooRepository {
        
    IQueryable&lt;Foo&gt; GetAll();
    Foo GetSingle(int fooId);
    IQueryable&lt;Foo&gt; FindBy(Expression&lt;Func&lt;Foo, <span style="color: blue;">bool</span>&gt;&gt; predicate);
    <span style="color: blue;">void</span> Add(Foo entity);
    <span style="color: blue;">void</span> Delete(Foo entity);
    <span style="color: blue;">void</span> Edit(Foo entity);
    <span style="color: blue;">void</span> Save();
}</pre>
</div>
</div>
<p>This repo is for <strong>Foo</strong> class I have (imaginary). Let's see the implementation for <strong>Bar</strong> class.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">interface</span> IBarRepository {
    
    IQueryable&lt;Bar&gt; GetAll();
    Bar GetSingle(<span style="color: blue;">int</span> barId);
    IQueryable&lt;Bar&gt; FindBy(Expression&lt;Func&lt;Bar, <span style="color: blue;">bool</span>&gt;&gt; predicate);
    <span style="color: blue;">void</span> Add(Bar entity);
    <span style="color: blue;">void</span> Delete(Bar entity);
    <span style="color: blue;">void</span> Edit(Bar entity);
    <span style="color: blue;">void</span> Save();
}</pre>
</div>
</div>
<p>Implementation nearly exactly the same here. Here is also an example of implementing one of these interfaces:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> FooRepository : IFooRepository {

    <span style="color: blue;">private</span> readonly FooBarEntities context = <span style="color: blue;">new</span> FooBarEntities();

    <span style="color: blue;">public</span> IQueryable&lt;Foo&gt; GetAll() {

        IQueryable&lt;Foo&gt; query = context.Foos;
        <span style="color: blue;">return</span> query;
    }

    <span style="color: blue;">public</span> Foo GetSingle(<span style="color: blue;">int</span> fooId) {

        <span style="color: blue;">var</span> query = <span style="color: blue;">this</span>.GetAll().FirstOrDefault(x =&gt; x.FooId == fooId);
        <span style="color: blue;">return</span> query;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Add(Foo entity) {

        context.Foos.Add(entity);
    }

    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Delete(Foo entity) {

        context.Foos.Remove(entity);
    }

    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Edit(Foo entity) {

        context.Entry&lt;Foo&gt;(entity).State = System.Data.EntityState.Modified;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Save() {

        context.SaveChanges();
    }
}</pre>
</div>
</div>
<p>Also imagine this implementation for BarRepository as well. Indeed, there would be probably more repository classes for your project. After playing like that for a while I decided to do something different which still sucked but better.</p>
<p><strong>A better approach but still sucks</strong></p>
<p>I created a generic interface which saves me a lot of keystrokes. Here how it looks like:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">interface</span> IGenericRepository&lt;T&gt; <span style="color: blue;">where</span> T : <span style="color: blue;">class</span> {
    
    IQueryable&lt;T&gt; GetAll();
    IQueryable&lt;T&gt; FindBy(Expression&lt;Func&lt;T, <span style="color: blue;">bool</span>&gt;&gt; predicate);
    <span style="color: blue;">void</span> Add(T entity);
    <span style="color: blue;">void</span> Delete(T entity);
    <span style="color: blue;">void</span> Edit(T entity);
    <span style="color: blue;">void</span> Save();
}</pre>
</div>
</div>
<p>And how I implemented in on repository interfaces:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">interface</span> IFooRepository : IGenericRepository&lt;Foo&gt; {
    
    Foo GetSingle(<span style="color: blue;">int</span> fooId);
}</pre>
</div>
</div>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">interface</span> IBarRepository : IGenericRepository&lt;Bar&gt; {
    
    Bar GetSingle(<span style="color: blue;">int</span> barId);
}</pre>
</div>
</div>
<p>You can see that I only needed to implement <strong>GetSingle</strong> method here and others come with <strong>IGenericRepositoy&lt;T&gt;</strong> interface.</p>
<p>Where I implement these repository interfaces to my concrete classes, I still need to go over all the methods and create them individually. The repository class looked like as the same. So it leads me to a final solution which is the best one I can come up with so far.</p>
<p><strong>Best approach</strong></p>
<p>The generic interface I have created is still legitimate and usable here. In fact, I won&rsquo;t touch the repository interfaces at all. What I did here first is to create an abstract class which implements IGenericReposity&lt;T&gt; interface but also accepts another type parameter defined in a generic declaration which is a type of DbConetxt class. Here is how it looks like:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">abstract</span> <span style="color: blue;">class</span> GenericRepository&lt;C, T&gt; : 
    IGenericRepository&lt;T&gt; <span style="color: blue;">where</span> T : <span style="color: blue;">class</span> <span style="color: blue;">where</span> C : DbContext, <span style="color: blue;">new</span>() {

    <span style="color: blue;">private</span> C _entities = <span style="color: blue;">new</span> C();
    <span style="color: blue;">public</span> C Context {

        <span style="color: blue;">get</span> { <span style="color: blue;">return</span> _entities; }
        <span style="color: blue;">set</span> { _entities = value; }
    }

    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> IQueryable&lt;T&gt; GetAll() {

        IQueryable&lt;T&gt; query = _entities.Set&lt;T&gt;();
        <span style="color: blue;">return</span> query;
    }

    <span style="color: blue;">public</span> IQueryable&lt;T&gt; FindBy(System.Linq.Expressions.Expression&lt;Func&lt;T, <span style="color: blue;">bool</span>&gt;&gt; predicate) {

        IQueryable&lt;T&gt; query = _entities.Set&lt;T&gt;().Where(predicate);
        <span style="color: blue;">return</span> query;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> <span style="color: blue;">void</span> Add(T entity) {
        _entities.Set&lt;T&gt;().Add(entity);
    }

    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> <span style="color: blue;">void</span> Delete(T entity) {
        _entities.Set&lt;T&gt;().Remove(entity);
    }

    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> <span style="color: blue;">void</span> Edit(T entity) {
        _entities.Entry(entity).State = System.Data.EntityState.Modified;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> <span style="color: blue;">void</span> Save() {
        _entities.SaveChanges();
    }
}</pre>
</div>
</div>
<p>This is so nice because of some factors I like:</p>
<ul>
<li>This implements so basic and ordinary methods</li>
<li>If necessary, those methods can be overridden because each method is virtual.</li>
<li>As we newed up the DbContext class here and expose it public with a public property, we have flexibility of extend the individual repositories for our needs.</li>
<li>As we only implement this abstract class only to our repository classes, it won&rsquo;t effect unit testing at all. DbContext is not in the picture in terms of unit testing.</li>
</ul>
<p>So, when we need to implement these changes to our concrete repository classes, we will end up with following result:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> FooRepository :
    GenericRepository&lt;FooBarEntities, Foo&gt;, IFooRepository {

    <span style="color: blue;">public</span> Foo GetSingle(<span style="color: blue;">int</span> fooId) {

        <span style="color: blue;">var</span> query = GetAll().FirstOrDefault(x =&gt; x.FooId == fooId);
        <span style="color: blue;">return</span> query;
    }
}</pre>
</div>
</div>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> BarReposiltory : 
    GenericRepository&lt;FooBarEntities, Bar&gt;, IBarRepository  {

    <span style="color: blue;">public</span> Bar GetSingle(<span style="color: blue;">int</span> barId) {

        <span style="color: blue;">var</span> query = Context.Bars.FirstOrDefault(x =&gt; x.BarId == barId);
        <span style="color: blue;">return</span> query;
    }
}</pre>
</div>
</div>
<p>Very nice and clean. Inside <strong>BarRepository</strong> <strong>GetSingle</strong> method, as you see I use <strong>Context</strong> property of <strong>GenericRepository&lt;C, T&gt;</strong> abstract class to access an instance of DbContext.</p>
<p>So, how the things work inside our <a target="_blank" href="http://asp.net/mvc" title="http://asp.net/mvc">ASP.NET MVC</a> project? this is another story but no so complicated. I will continue right from here on my next post.</p>
<p><strong>UPDATE:</strong></p>
<p>Here is the next post:</p>
<p><a href="http://www.tugberkugurlu.com/archive/how-to-work-with-generic-repositories-on-asp-net-mvc-and-unit-testing-them-by-mocking" title="http://www.tugberkugurlu.com/archive/how-to-work-with-generic-repositories-on-asp-net-mvc-and-unit-testing-them-by-mocking" target="_blank">How to Work With Generic Repositories on ASP.NET MVC and Unit Testing Them By Mocking</a></p>