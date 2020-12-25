---
id: 176adb7e-b112-4416-8b2c-1000e0977944
title: Asynchronous .NET Client Libraries for Your HTTP API and Awareness of async/await's
  Bad Effects
abstract: Writing asynchronous .NET Client libraries for your HTTP API and using asynchronous
  language features (aka async/await) and some deadlock issue you might face.
created_at: 2012-09-21 06:34:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- ASP.NET Web API
- async
- TPL
slugs:
- asynchronousnet-client-libraries-for-your-http-api-and-awareness-of-async-await-s-bad-effects
---

<p>Haven&rsquo;t you shot yourself in the foot yet with async/await? If not, you are about to if you are writing a client library for your newly created ASP.NET Web API application with .NET 4.5 using new asynchronous language features.</p>
<p>I wrote a blog post couple of months ago on the <a href="https://www.tugberkugurlu.com/archive/the-perfect-recipe-to-shoot-yourself-in-the-foot-ending-up-with-a-deadlock-using-the-c-sharp-5-0-asynchronous-language-features">importance of Current SynchronizationContext and the new C# 5.0 asynchronous language Features</a> (aka async/await). I wrote that post just after I watched the <a href="http://channel9.msdn.com/Events/BUILD/BUILD2011/TOOL-829T">The zen of async: Best practices for best performance</a> talk of <a href="http://blogs.msdn.com/b/pfxteam/">Stephen Toub</a> on //Build 2011 and that was one of the best sessions that I was and still am glad to watch. I learnt so many things from that session and some of them was amazingly important.</p>
<blockquote>
<p>Filip Ekberg also has a very nice and useful blog post which has an identical title as my previous post on the topic: <a href="http://blog.filipekberg.se/2012/09/20/avoid-shooting-yourself-in-the-foot-with-tasks-and-async/">Avoid shooting yourself in the foot with Tasks and Async</a>.</p>
</blockquote>
<p>The post was pointing out that it is extremely easy to end up with a deadlock if you are not careful enough. The post explains every details but if want to recap shortly, here it is:</p>
<blockquote>
<p>When you are awaiting on a method with await keyword, compiler generates bunch of code in behalf of you. One of the purposes of this action is to handle synchronization with the UI thread. The key component of this feature is the<a href="http://msdn.microsoft.com/en-us/library/system.threading.synchronizationcontext.current.aspx">SynchronizationContext.Current</a> which gets the synchronization context for the current thread. SynchronizationContext.Current is populated depending on the environment you are in. The GetAwaiter method of Task looks up for SynchronizationContext.Current. If current synchronization context is not null, the continuation that gets passed to that awaiter will get posted back to that synchronization context.</p>
<p>When consuming a method, which uses the new asynchronous language features, in a blocking fashion, you will end up with a deadlock if you have an available SynchronizationContext. When you are consuming such methods in a blocking fashion (waiting on the Task with Wait method or taking the result directly from the Result property of the Task), you will block the main thread at the same time. When eventually the Task completes inside that method in the threadpool, it is going to invoke the continuation to post back to the main thread because SynchronizationContext.Current is available and captured. But there is a problem here: the UI thread is blocked and you have a deadlock!</p>
</blockquote>
<p>The post also has a sample but I want dig deep into this. As you all know, <a href="http://www.asp.net/web-api">ASP.NET Web API</a> RTW version has been shipped couple of weeks ago and we may have started to take advantage of this framework to build our HTTP-based lightweight APIs. You might be using another framework to create your APIs or you might be maintaining an existing one. No matter which category you put yourself in, you are likely to build platform specific client libraries for your API. If you are going to target .NET framework, I am 90% sure that you want your client library to be asynchronous and new <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httpclient.aspx">HttpClient</a> will just give you that option.</p>
<p>I have started writing a simple blog engine for myself called <a href="https://github.com/tugberkugurlu/MvcBloggy">MvcBloggy</a> nearly a year ago and after the RTW release of the ASP.NET MVC 4 and ASP.NET Web API, I decided to try something different and expose the data through HTTP. So, my ASP.NET MVC application won&rsquo;t know anything about where my data is store or how I retrieve it. It is just going to consume the HTTP APIs. The application is shaping up nicely and I even started to build <a href="https://github.com/tugberkugurlu/MvcBloggy/tree/develop/src/MvcBloggy.API.Client">my .NET client for the API</a> as well.</p>
<p>Couple of days ago, <a href="http://ben.onfabrik.com/">Ben Foster</a> (a brilliant developer) raised a question on Twitter about consuming asynchronous methods on ASP.NET Web Pages application and I also wondered about that. Then, we looked around a bit and also contacted with Erik Porter (<a href="https://twitter.com/HumanCompiler">@HumanCompiler</a>), Program Manager on the ASP.NET team, and he confirmed that it is not possible today. He encouraged us to file an issue and we did (<a href="http://aspnetwebstack.codeplex.com/workitem/418">#418</a>). This got me thinking about my little project&rsquo;s .NET client, though. What if I wanted to create a blog web site with ASP.NET Web Pages by consuming the .NET client of my blog engine? I have no option rather than consuming the methods in a blocking fashion but with my current implementation, I am not able to do that because I will end up with deadlocks. Let me prove that to you with a little example which you can also find up on GitHub: <a href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/AsyncAwaitForLibraryAuthors">AsyncAwaitForLibraryAuthors</a>.</p>
<p>Assuming we have a little API which returns back a list of cars and we want to build a .NET client to consume the HTTP API and abstracts all the lower level HTTP stuff away. The one that I created is as below (it is a simple one for the demo purposes):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Car {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Id { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Make { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Model { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Year { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">float</span> Price { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> SampleAPIClient {

    <span style="color: blue;">private</span> <span style="color: blue;">const</span> <span style="color: blue;">string</span> ApiUri = <span style="color: #a31515;">"http://localhost:17257/api/cars"</span>;

    <span style="color: blue;">public</span> async Task&lt;IEnumerable&lt;Car&gt;&gt; GetCarsAsync() {

        <span style="color: blue;">using</span> (HttpClient client = <span style="color: blue;">new</span> HttpClient()) {

            <span style="color: blue;">var</span> response = await client.GetAsync(ApiUri);

            <span style="color: green;">// Not the best way to handle it but will do the work for demo purposes</span>
            response.EnsureSuccessStatusCode();
            <span style="color: blue;">return</span> await response.Content.ReadAsAsync&lt;IEnumerable&lt;Car&gt;&gt;();
        }
    }
}</pre>
</div>
</div>
<p>We leveraged the new asynchronous language features and we were able to write a small amount of code to get the job done. More importantly, we are making the network call and the deserialization&nbsp;asynchronously. As we know from the earlier sentences that if there is a SynchronizationContext available for us, the code that the compiler is generating for us will capture that and post the continuation back to that context to be executed. Keep this part in mind. I put this little class inside a separate project called SampleAPI.Client and reference this from my web clients. I have created two clients: one ASP.NET MVC and one ASP.NET Web Pages applications.</p>
<p>In my ASP.NET MVC 4 application, I have a controller which has two actions. One of these actions will call the API asynchronously and one of will do the same by blocking:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">public</span> async Task&lt;ViewResult&gt; CarsAsync() {

        SampleAPIClient client = <span style="color: blue;">new</span> SampleAPIClient();
        <span style="color: blue;">var</span> cars = await client.GetCarsAsync();

        <span style="color: blue;">return</span> View(<span style="color: #a31515;">"Index"</span>, model: cars);
    }

    <span style="color: blue;">public</span> ViewResult CarsSync() {

        SampleAPIClient client = <span style="color: blue;">new</span> SampleAPIClient();
        <span style="color: blue;">var</span> cars = client.GetCarsAsync().Result;

        <span style="color: blue;">return</span> View(<span style="color: #a31515;">"Index"</span>, model: cars);
    }
}</pre>
</div>
</div>
<p>Our view is also so simple as follows:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>@model IEnumerable&lt;SampleAPI.Client.Car&gt;
@{
    ViewBag.Title = "Home Page";
}

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>Cars List<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>

<span style="color: blue;">&lt;</span><span style="color: #a31515;">ul</span><span style="color: blue;">&gt;</span>
    @foreach (var car in Model) {
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">li</span><span style="color: blue;">&gt;</span>
            @car.Make, @car.Model (@car.Year) - @car.Price.ToString("C")
        <span style="color: blue;">&lt;/</span><span style="color: #a31515;">li</span><span style="color: blue;">&gt;</span>    
    }
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">ul</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>When we navigate to /home/CarsAsync, we will get back the result.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/af4459ecb4fd_9B63/image.png"><img height="388" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/af4459ecb4fd_9B63/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="image" /></a></p>
<p>However, when we navigate to /home/CarsSync to invoke the CarsSync method, we will see that the page will never come back because we just introduced a deadlock due to the reasons we have explained earlier. Let&rsquo;s have a look at the Web Pages sample:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>@{
    Layout = "~/_SiteLayout.cshtml";
    Page.Title = "Home Page";

    SampleAPI.Client.SampleAPIClient client = new SampleAPI.Client.SampleAPIClient();
    var cars = client.GetCarsAsync().Result;
}

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>Cars List<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>

<span style="color: blue;">&lt;</span><span style="color: #a31515;">ul</span><span style="color: blue;">&gt;</span>
    @foreach (var car in cars) {
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">li</span><span style="color: blue;">&gt;</span>
            @car.Make, @car.Model (@car.Year) - @car.Price.ToString("C")
        <span style="color: blue;">&lt;/</span><span style="color: #a31515;">li</span><span style="color: blue;">&gt;</span>    
    }
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">ul</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>This page is also not going to respond because Web Pages runs under ASP.NET and ASP.NET has a SynchronizationContext available.</p>
<p>When we take a look at our GetCarsAsync method implementation, we will see that it is completely unnecessary for us to get back to current SynchronizationContext because we don&rsquo;t need anything from the current context. This is good because it is not our (I mean our .NET client&rsquo;s) concern to do anything under the current SynchronizationContext. It is, on the other hand, our consumer&rsquo;s responsibility. Stephen Toub said something in his talk on //Build 2011 and the words not the same but it expresses the meaning of the below sentences:</p>
<blockquote>
<p>If you are a library developer, the default behavior which await gives you is nearly never what you want. However, if you are a application developer, the default behavior will nearly always what you want.</p>
</blockquote>
<p>I, again, encourage you to check that video out.</p>
<p>The solution here is simple. When we are creating our libraries, we just need to be more careful and think about the usage scenarios. In our case here, we need to suppress the default SynchronizationContext behavior that the compiler is generating for us. We can achieve this with <a href="http://msdn.microsoft.com/en-us/library/system.threading.tasks.task.configureawait.aspx">ConfigureAwait</a> method of the Task class which was introduced with .NET 4.5. The ConfigureAwait method accepts a Boolean parameter named as continueOnCapturedContext. We can pass false into this method not to marshal the continuation back to the original context captured and our problem would be solved. Here is the new look of our .NET client for our HTTP API.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> SampleAPIClient {

    <span style="color: blue;">private</span> <span style="color: blue;">const</span> <span style="color: blue;">string</span> ApiUri = <span style="color: #a31515;">"http://localhost:17257/api/cars"</span>;

    <span style="color: blue;">public</span> async Task&lt;IEnumerable&lt;Car&gt;&gt; GetCarsAsync() {

        <span style="color: blue;">using</span> (HttpClient client = <span style="color: blue;">new</span> HttpClient()) {

            <span style="color: blue;">var</span> response = await client.GetAsync(ApiUri)
                .ConfigureAwait(continueOnCapturedContext: <span style="color: blue;">false</span>);

            <span style="color: green;">// Not the best way to handle it but will do the work for demo purposes</span>
            response.EnsureSuccessStatusCode();
            <span style="color: blue;">return</span> await response.Content.ReadAsAsync&lt;IEnumerable&lt;Car&gt;&gt;()
                .ConfigureAwait(continueOnCapturedContext: <span style="color: blue;">false</span>);
        }
    }
}</pre>
</div>
</div>
<p>When we now run our Web Pages application, we will see the web site working nicely (same is also applicable for the CarsSync action method of our ASP.NET MVC application).</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/af4459ecb4fd_9B63/image_3.png"><img height="396" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/af4459ecb4fd_9B63/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="image" /></a></p>
<p>If you are going to write a .NET client for your company&rsquo;s big HTTP API using new asynchronous language features, you might want to consider these facts before moving on. Otherwise, your consumers will have hard time understanding what is really going wrong.</p>