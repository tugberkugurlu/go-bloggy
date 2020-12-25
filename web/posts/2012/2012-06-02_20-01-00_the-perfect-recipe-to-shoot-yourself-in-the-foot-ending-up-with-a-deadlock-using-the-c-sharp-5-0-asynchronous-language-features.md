---
id: f46c8b51-5370-4529-9e0e-ba997bc9796d
title: The Perfect Recipe to Shoot Yourself in The Foot - Ending up with a Deadlock
  Using the C# 5.0 Asynchronous Language Features
abstract: Let's see how we can end up with a deadlock using the C# 5.0 asynchronous
  language features (AKA async/await) in our ASP.NET applications and how to prevent
  these kinds of scenarios.
created_at: 2012-06-02 20:01:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- ASP.NET Web API
- async
- C#
- TPL
slugs:
- the-perfect-recipe-to-shoot-yourself-in-the-foot-ending-up-with-a-deadlock-using-the-c-sharp-5-0-asynchronous-language-features
- the-perfect-recipe-to-shot-yourself-in-the-foot-ending-up-with-a-deadlock-using-the-c-sharp-5-0-asynchronous-language-features
---

<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/How-to-Shot-Yourself-in.NET-and-Deadlock_128/shoot_yourself_in_the_foot.jpg"><img style="background-image: none; margin: 0px 10px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border: 0px;" title="shoot_yourself_in_the_foot" border="0" alt="shoot_yourself_in_the_foot" align="left" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/How-to-Shot-Yourself-in.NET-and-Deadlock_128/shoot_yourself_in_the_foot_thumb.jpg" width="244" height="163" /></a>I just finished <a title="http://channel9.msdn.com/Events/BUILD/BUILD2011/TOOL-829T" href="http://channel9.msdn.com/Events/BUILD/BUILD2011/TOOL-829T">The zen of async: Best practices for best performance</a>&nbsp;talk of <a href="http://blogs.msdn.com/b/pfxteam/" title="http://blogs.msdn.com/b/pfxteam/">Stephen Toub</a>&rsquo;s on //Build&nbsp;and learnt how easy is to end up with a deadlock while writing asynchronous code with new C# 5.0 language features (AKA async/await). Here is the quick background first which I was aware of before I watched this talk. The talk that Toub has given was mostly about how you create better reusable libraries with asynchronous language features but there are lots of great information concerning application level code.</p>
<p>When you are awaiting on a method with await keyword, compiler generates bunch of code in behalf of you. One of the purposes of this action is to handle synchronization with the UI thread. The key component of this feature is the <a title="http://msdn.microsoft.com/en-us/library/system.threading.synchronizationcontext.current.aspx" href="http://msdn.microsoft.com/en-us/library/system.threading.synchronizationcontext.current.aspx">SynchronizationContext.Current</a> which gets the synchronization context for the current thread. SynchronizationContext.Current is populated depending on the environment you are in. For example, if you are on a WPF application, SynchronizationContext.Current will be a type of <a title="http://msdn.microsoft.com/en-us/library/system.windows.threading.dispatchersynchronizationcontext.aspx" href="http://msdn.microsoft.com/en-us/library/system.windows.threading.dispatchersynchronizationcontext.aspx">DispatcherSynchronizationContext</a>. In an ASP.NET application, this will be an instance of AspNetSynchronizationContext which is not publicly available and not meant for external consumptions. SynchronizationContext has virtual <a title="http://msdn.microsoft.com/en-us/library/system.threading.synchronizationcontext.post" href="http://msdn.microsoft.com/en-us/library/system.threading.synchronizationcontext.post">Post</a> method which dispatches an asynchronous message to a synchronization context when overridden in a derived class. When you give this method a delegate, this delegate will be marshaled back to the UI thread and invoked on that UI thread.</p>
<p>The GetAwaiter method (yes, "waiter" joke is proper here) of Task looks up for SynchronizationContext.Current. If current synchronization context is not null, the continuation that gets passed to that awaiter will get posted back to that synchronization context. This feature is a great feature when we are writing application level code but it turns out that if we are not careful enough, we might end up with a deadlock because of this. Here is how:</p>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/How-to-Shot-Yourself-in.NET-and-Deadlock_128/build_task.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="build_task" border="0" alt="build_task" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/How-to-Shot-Yourself-in.NET-and-Deadlock_128/build_task_thumb.png" width="644" height="259" /></a></p>
<p>This is a part of the slide from that talk and shows how we can end up with a deadlock. Task has a method named Wait which hangs on the task till it completes. This method is an evil method in my opinion but I am sure that there is a good reason why it is provided. You never ever want to use this method. But, assume that you did use it on a method which returns a Task or Task&lt;T&gt; for some T and uses await inside it to await another task, you will end up with a deadlock if you are not careful. The picture explains how that can happen but let&rsquo;s recap with an example.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">public</span> ActionResult Index() {
            
        doWorkAsync().Wait();
        <span style="color: blue;">return</span> View();
    }

    <span style="color: blue;">private</span> async Task doWorkAsync() {

        await Task.Delay(500);
    }

}</pre>
</div>
</div>
<p>The above code has an ASP.NET MVC 4 asynchronous action method (I am using .NET 4.5 here) which consumes the private doWorkAsync method. I used ASP.NET MVC here but everything applies for ASP.NET Web API asynchronous action methods as well. Inside the doWorkAsync method, we used Delay method of Task to demonstrate but this could be any method which returns a Task or Task&lt;T&gt;. inside the Index action method, we invoke the Wait method on the task which we got from doWorkAsync method to ensure that we won&rsquo;t go further unless the operation completes and the interesting part happens right here. At that point we block the UI thread at the same time. When eventually the Task.Delay method completes in the threadpool, it is going to invoke the continuation to post back to the UI thread because SynchronizationContext.Current is available and captured. But there is a problem here: the UI thread is blocked. Say hello to our deadlock!</p>
<p>When you run the application, you will see that the web page will never come back. Solution to this problem is so simple: don&rsquo;t use the Wait method. The code you should be writing should be as follows:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">public</span> async Task&lt;ActionResult&gt; Index() {
        
        await doWorkAsync();
        <span style="color: blue;">return</span> View();
    }

    <span style="color: blue;">private</span> async Task doWorkAsync() {

        <span style="color: blue;">var</span> task = Task.Delay(500);
    }
}</pre>
</div>
</div>
<p>But if you have to use it for some weird reasons, there is another method named ConfigureAwait on Task which you can configure not to use the captured context.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">public</span> ActionResult Index() {
        
        doWorkAsync().Wait();
        <span style="color: blue;">return</span> View();
    }

    <span style="color: blue;">private</span> async Task doWorkAsync() {

        <span style="color: blue;">var</span> task = Task.Delay(500);
        await task.ConfigureAwait(continueOnCapturedContext: <span style="color: blue;">false</span>);
    }
}</pre>
</div>
</div>
<p>Very nice little and tricky information but a lifesaver.</p>
<h3>Resources</h3>
<ul>
<li><a title="http://blogs.msdn.com/b/pfxteam/archive/2012/01/20/10259049.aspx" href="http://blogs.msdn.com/b/pfxteam/archive/2012/01/20/10259049.aspx">Await, SynchronizationContext, and Console Apps</a></li>
</ul>