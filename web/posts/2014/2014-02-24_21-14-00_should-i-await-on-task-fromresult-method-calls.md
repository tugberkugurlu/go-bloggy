---
title: Should I await on Task.FromResult Method Calls?
abstract: Task class has a static method called FromResult which returns an already
  completed (at the RanToCompletion status) Task object. I have seen a few developers
  "await"ing on Task.FromResult method call and this clearly indicates that there
  is a misunderstanding here. I'm hoping to clear the air a bit with this post.
created_at: 2014-02-24 21:14:00 +0000 UTC
tags:
- async
- C#
- TPL
slugs:
- should-i-await-on-task-fromresult-method-calls
---

<p><a href="http://msdn.microsoft.com/en-us/library/system.threading.tasks.task(v=vs.110).aspx">Task</a> class has a static method called <a href="http://msdn.microsoft.com/en-us/library/hh194922(v=vs.110).aspx">FromResult</a> which returns an already completed (at the <a href="http://msdn.microsoft.com/en-us/library/system.threading.tasks.taskstatus(v=vs.110).aspx">RanToCompletion</a> status) Task object. I have seen a few developers "await"ing on Task.FromResult method call and this clearly indicates that there is a misunderstanding here. I'm hoping to clear the air a bit with this post. <h3>What is the use of Task.FromResult method?</h3> <p>Imagine a situation where you are implementing an interface which has the following signature:  <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">interface</span> IFileManager
{
     Task&lt;IEnumerable&lt;File&gt;&gt; GetFilesAsync();
}</pre></div></div>
<p>Notice that the method is Task returning which allows you to make the return expression represent an ongoing operation and also allows the consumer of this method to call this method in an asynchronous manner without blocking (of course, if the underlying layer supports it). However, depending on the case, your operation may not be asynchronous. For example, you may just have the files inside an in memory collection and want to return it from there, or you can perform an I/O operation to retrieve the files list asynchronously from a particular data store for the first time and cache the results there so that you can just return it from the in-memory cache for the upcoming calls. These are just some scenarios where you need to return a successfully completed Task object. Here is how you can achieve that without the help of Task.FromResult method: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> InMemoryFileManager : IFileManager
{
    IEnumerable&lt;File&gt; files = <span style="color: blue">new</span> List&lt;File&gt;
    {
        <span style="color: green">//...</span>
    };

    <span style="color: blue">public</span> Task&lt;IEnumerable&lt;File&gt;&gt; GetFilesAsync()
    {
        <span style="color: blue">var</span> tcs = <span style="color: blue">new</span> TaskCompletionSource&lt;IEnumerable&lt;File&gt;&gt;();
        tcs.SetResult(files);

        <span style="color: blue">return</span> tcs.Task;
    }
}</pre></div></div>
<p>We here used the <a href="http://msdn.microsoft.com/en-us/library/dd449174(v=vs.110).aspx">TaskCompletionSource</a> to produce a successfully completed Task object with the result. Therefore, the caller of the method will immediately have the result. This was what we had been doing till the introduction of .NET 4.5. If you are on .NET 4.5 or above, you can just use the Task.FromResult to perform the same operation:
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> InMemoryFileManager : IFileManager
{
    IEnumerable&lt;File&gt; files = <span style="color: blue">new</span> List&lt;File&gt;
    {
        <span style="color: green">//...</span>
    };

    <span style="color: blue">public</span> Task&lt;IEnumerable&lt;File&gt;&gt; GetFilesAsync()
    {
        <span style="color: blue">return</span> Task.FromResult&lt;IEnumerable&lt;File&gt;&gt;(files);
    }
}</pre></div>
<h3>Should I await Task.FromResult method calls?</h3>
<p>TL;DR version of the answer: absolutely not! If you find yourself in need to using Task.FromResult, it's clear that you are not performing any asynchronous operation. Therefore, just return the Task from the Task.FromResult output. Is it dangerous to do this? Not completely but it's illogical and has a performance effect.
<p>Long version of the answer is a bit more in depth. Let's first see what happens when you "await" on a method which matches the pattern:
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>IEnumerable&lt;File&gt; files = await fileManager.GetFilesAsync();</pre></div></div></div>
<p>This code will be read by the compiler as follows (well, in a simplest way):
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">var</span> $awaiter = fileManager.GetFilesAsync().GetAwaiter();
<span style="color: blue">if</span>(!$awaiter.IsCompleted) 
{
     DO THE AWAIT/RETURN AND RESUME
}

<span style="color: blue">var</span> files = $awaiter.GetResult();</pre></div>
<p>Here, we can see that if the awaited Task already completed, then it skips all the await/resume work and directly gets the result. Besides this fact, if you put "async" keyword on a method, a bunch of code (including the state machine) is generated regardless of the fact that you use await keyword inside the method or not. Keeping all these facts in mind, implementing the IFileManager as below is going to cause nothing but overhead:
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> InMemoryFileManager : IFileManager
{
    IEnumerable&lt;File&gt; files = <span style="color: blue">new</span> List&lt;File&gt;
    {
        <span style="color: green">//...</span>
    };

    <span style="color: blue">public</span> async Task&lt;IEnumerable&lt;File&gt;&gt; GetFilesAsync()
    {
        <span style="color: blue">return</span> await Task.FromResult&lt;IEnumerable&lt;File&gt;&gt;(files);
    }
}</pre></div>
<p>So, don't ever think about "await"ing on Task.FromResult or I'll hunt you down in your sweet dreams :)</p>
<h3>References</h3>
<ul>
<li><a href="http://dhickey.ie/post/2013/04/11/An-async-await-antipattern-awaiting-tasks-instead-of-just-returning-them.aspx">An async await anti-pattern - awaiting tasks instead of just returning them</a>
<li><a href="http://channel9.msdn.com/Series/Three-Essential-Tips-for-Async/Async-libraries-APIs-should-be-chunky">Tip 5: Async libraries APIs should be chunky</a></li></ul></div></div>  