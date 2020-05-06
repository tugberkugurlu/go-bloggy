---
title: My Take on Task-base Asynchronous Programming in C# 5.0 and ASP.NET MVC Web
  Applications
abstract: I'm trying to show you what new C# 5.0 can bring us in terms of asynchronous
  programming with await keyword. Especially on ASP.NET MVC 4 Web Applications.
created_at: 2012-02-26 17:21:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- ASP.NET Web API
- async
- C#
- TPL
slugs:
- my-take-on-task-base-asynchronous-programming-in-c-sharp-5-0-and-asp-net-mvc-web-applications
---

<p>I have been playing with <a title="http://msdn.microsoft.com/en-us/vstudio/gg316360" href="http://msdn.microsoft.com/en-us/vstudio/gg316360" target="_blank">Visual Studio Async CTP</a> for a while now and I had enough idea on how it works and how I can use it for my benefit and that means that I am able to blog about it.</p>
<p>On most of the softwares we build, the main problem occurs on long running operations. If a user starts a long running operation, that operation blocks the main thread till it completes. Doing so will result pretty unusable applications and unhappy customers for your business. I am sure that we all have used that kind of applications in our lives.</p>
<p>Asynchronous programming model is not a new paradigm and it has been available on .NET since v1.0. I wasn&rsquo;t very interested in programming 3 years ago so this concept is fairly new to me even now. The concept has evolved a lot on .NET in terms of asynchronous programming as far as I read and I think it has reached its best stage.</p>
<p>On the other hand, you can get confused about asynchronous programming pretty much easily as I did several times. I went back and forth about how it works and where to use it. As a result, I believe I finally got it right. On the other hand,&nbsp;It is much easier to get confused if you&rsquo;d like to leverage asynchrony on web applications because you do not have a UI thread to flush your results out immediately. Asynchronous or not, the user has to wait the amount of time operation takes to complete. This makes asynchronous programming undesirable for web applications but if you think that way, as I did, you are missing the point.</p>
<p>If you start an operation synchronously and it takes long time, you have no option rather than blocking that thread. If you do the same operation asynchronously, what happens is that you start an operation and come back later when it finishes without waiting it to be finished. Between that duration, your thread is free and can do other stuff as well. Most of the confusion happens here I think. Creating additional threads is not cheap and may cause issues but it is all in past I think. .NET 4.0 has increases the number of thread limits very dramatically. It is still not good to walk around and create threads but it is good not to worry about them. There are couple of debates going around about asynchronous programming (or should I have said went around?):</p>
<ul>
<li><a title="http://blogs.msdn.com/b/rickandy/archive/2009/11/14/should-my-database-calls-be-asynchronous.aspx" href="http://blogs.msdn.com/b/rickandy/archive/2009/11/14/should-my-database-calls-be-asynchronous.aspx" target="_blank">Should my database calls be Asynchronous?</a> </li>
<li><a title="http://blogs.msdn.com/b/rickandy/archive/2011/07/19/should-my-database-calls-be-asynchronous-part-ii.aspx" href="http://blogs.msdn.com/b/rickandy/archive/2011/07/19/should-my-database-calls-be-asynchronous-part-ii.aspx" target="_blank">Should my database calls be Asynchronous Part II</a> </li>
<li><a title="http://stackoverflow.com/questions/8743067/do-asynchronous-operations-in-asp-net-mvc-use-a-thread-from-threadpool-on-net-4" href="http://stackoverflow.com/questions/8743067/do-asynchronous-operations-in-asp-net-mvc-use-a-thread-from-threadpool-on-net-4" target="_blank">Do asynchronous operations in ASP.NET MVC use a thread from ThreadPool on .NET 4</a> </li>
<li><a title="http://stackoverflow.com/questions/9432647/any-disadvantage-of-using-executereaderasync-from-c-sharp-asyncctp" href="http://stackoverflow.com/questions/9432647/any-disadvantage-of-using-executereaderasync-from-c-sharp-asyncctp" target="_blank">Any disadvantage of using ExecuteReaderAsync from C# AsyncCTP</a></li>
</ul>
<p>I am not sure what is the real answer of those questions. But here are two examples:</p>
<ul>
<li>Windows Runtime (WinRT) is designed to be async friendly. Anything longer than 40ms (operations related to network, file system. Mainly I/O bound operations), that is async. </li>
<li><a title="http://www.asp.net/web-api" href="http://www.asp.net/web-api" target="_blank">ASP.NET Web API</a> introduced a new way of exposing your data to the World with different formats. One thing to notice is that there is nearly no synchronous method on their API. Not metaphorically, literally there is no synchronous versions of the some methods.</li>
</ul>
<p>That&rsquo;s being said, I think we have a point here that we should take advantage of asynchrony one way or another in our applications. So, what is the problem? The problem is the way how async programming model works. As <a title="http://wikipedia.org/wiki/Anders_Hejlsberg" href="http://wikipedia.org/wiki/Anders_Hejlsberg" target="_blank">Anders Hejlsberg</a> always says, that programing model turns our code inside out. Not to mention that it is extremely hard to do nested async operations and exception handling.</p>
<p><strong>Visual Studio Async CTP and How It Works</strong></p>
<p>In C# 5.0, we will have a new asynchronous programming model which looks a lot like synchronous. On the other hand, C# team has released a CTP version of those features and it has a go-live license. Here is quote from the AsyncCTP spec:</p>
<blockquote>
<p><strong>Asynchronous functions</strong> is a new feature in C# which provides an easy means for expressing asynchronous operations. Inside asynchronous functions <strong>await expressions</strong> can <strong>await</strong> ongoing tasks, which causes the rest of the execution of the asynchronous function to be transparently signed up as a continuation of the awaited task. In other words, it becomes the job of the programming language, not the programmer, to express and sign up continuations. As a result, asynchronous code can retain its logical structure.</p>
</blockquote>
<p>An <b><i>asynchronous function</i></b> is a method or anonymous function which is marked with the <strong>async</strong> modifier. An asynchronous function can either return Task or Task&lt;T&gt; for some T and both of them can be awaited. Those kind of functions can also return void but it cannot be awaited. On the other hand, you can await any type if that type satisfies a certain pattern.</p>
<p><strong>ASP.NET MVC 4 and C# Async Features</strong></p>
<p>Benefiting from asynchrony in a right way on ASP.NET MVC applications can result huge positive performance impact. Believe it or not it&rsquo;s true. I&rsquo;ll show you how.</p>
<p>So why don't we use it much? Because it is hard and error prone. In the long run, it is hard to maintain the application as well. But with new asynchronous programming model, it is about to change.</p>
<p>In ASP.NET MVC 4, asynchronous programming model has been changed a lot as well. As you probably know, in ASP.NET MVC 3, our controller has to be derived from <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.asynccontroller.aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.asynccontroller.aspx" target="_blank">AsyncController</a> and must satisfy a certain pattern to work with. You can see the <a title="http://msdn.microsoft.com/en-us/library/ee728598.aspx" href="http://msdn.microsoft.com/en-us/library/ee728598.aspx" target="_blank">Using an Asynchronous Controller in ASP.NET MVC</a> article if you would like to see how it works.</p>
<p>In ASP.NET MVC 4, we do not need AsyncController to leverage asynchrony in our applications. Our controller actions can be marked with <strong>async </strong>keyword and return Task or Task&lt;T&gt; where the T is usually the type of <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionresult(v=vs.98).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionresult(v=vs.98).aspx" target="_blank">ActionResult</a>.</p>
<p>I put together a sample application which does the same thing both asynchronously and synchronously. I also did a load test and the end result was shocking. Let&rsquo;s see what the code looks like:</p>
<p>Firstly, I created a simple REST service endpoint with new REST hotness of .NET: <a title="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples" href="http://www.asp.net/web-api" target="_blank">ASP.NET Web API</a>. I have a simple model and collection which I store it in memory. This would be normally a database instead of memory.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Car {

    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Make;
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Model;
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Year;
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Doors;
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Colour;
    <span style="color: blue;">public</span> <span style="color: blue;">float</span> Price;
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Mileage;
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> CarService {

    <span style="color: blue;">public</span> List&lt;Car&gt; GetCars() {

        List&lt;Car&gt; Cars = <span style="color: blue;">new</span> List&lt;Car&gt; {

            <span style="color: blue;">new</span> Car{Make=<span style="color: #a31515;">"Audi"</span>,Model=<span style="color: #a31515;">"A4"</span>,Year=1995,Doors=4,Colour=<span style="color: #a31515;">"Red"</span>,Price=2995f,Mileage=122458},
            <span style="color: blue;">new</span> Car{Make=<span style="color: #a31515;">"Ford"</span>,Model=<span style="color: #a31515;">"Focus"</span>,Year=2002,Doors=5,Colour=<span style="color: #a31515;">"Black"</span>,Price=3250f,Mileage=68500},
            <span style="color: blue;">new</span> Car{Make=<span style="color: #a31515;">"BMW"</span>,Model=<span style="color: #a31515;">"5 Series"</span>,Year=2006,Doors=4,Colour=<span style="color: #a31515;">"Grey"</span>,Price=24950f,Mileage=19500}
            <span style="color: green;">//This keeps going like that</span>
        };

        <span style="color: blue;">return</span> Cars;
    }
}</pre>
</div>
</div>
<p>And here is my Web API:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController { 

    <span style="color: blue;">public</span> IEnumerable&lt;Car&gt; Get() {

        <span style="color: blue;">var</span> service = <span style="color: blue;">new</span> CarService();

        <span style="color: blue;">return</span> service.GetCars();
    }
}</pre>
</div>
</div>
<p>I have my service now. In my web application, I will get the data from this service and display on the web page. To do that, I created a service class which gets the data from that endpoint and deserialize the string into an object. I used the new <a title="http://blogs.msdn.com/b/henrikn/archive/2012/02/16/httpclient-is-here.aspx" href="http://blogs.msdn.com/b/henrikn/archive/2012/02/16/httpclient-is-here.aspx" target="_blank">HttpClient</a> for asynchronous version of GetCars operation and <a title="http://msdn.microsoft.com/en-us/library/system.net.webclient.aspx" href="http://msdn.microsoft.com/en-us/library/system.net.webclient.aspx" target="_blank">WebClient</a> for synchronous version of it. I also used <a title="http://nuget.org/packages/newtonsoft.json" href="http://nuget.org/packages/newtonsoft.json" target="_blank">Json.NET</a> for working with JSON payload.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarRESTService {

    <span style="color: blue;">readonly</span> <span style="color: blue;">string</span> uri = <span style="color: #a31515;">"http://localhost:2236/api/cars"</span>;

    <span style="color: blue;">public</span> List&lt;Car&gt; GetCars() { 

        <span style="color: blue;">using</span> (WebClient webClient = <span style="color: blue;">new</span> WebClient()) {
            
            <span style="color: blue;">return</span> JsonConvert.DeserializeObject&lt;List&lt;Car&gt;&gt;(
                webClient.DownloadString(uri)
            );
        }
    }

    <span style="color: blue;">public</span> async Task&lt;List&lt;Car&gt;&gt; GetCarsAsync() {

        <span style="color: blue;">using</span> (HttpClient httpClient = <span style="color: blue;">new</span> HttpClient()) {
            
            <span style="color: blue;">return</span> JsonConvert.DeserializeObject&lt;List&lt;Car&gt;&gt;(
                await httpClient.GetStringAsync(uri)    
            );
        }
    }
}</pre>
</div>
</div>
<p>Above GetCars method is very boring as you see. Nothing to talk about. The real deal is in the second method which is <strong>GetCarsAsync</strong>:</p>
<ul>
<li>The method is marked with <strong>async</strong> keyword which indicates that the method has some asynchronous code. </li>
<li>We used await keyword before <strong>HttpClient.GetStringAsync</strong> method which returns <strong>Task&lt;string&gt;.</strong> But notice here that we use it as <strong>string. </strong>The await keyword enables that. </li>
</ul>
<p>Lastly, here is my controller:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">private</span> CarRESTService service = <span style="color: blue;">new</span> CarRESTService();

    <span style="color: blue;">public</span> async Task&lt;ActionResult&gt; Index() {

        <span style="color: blue;">return</span> View(<span style="color: #a31515;">"index"</span>,
            await service.GetCarsAsync()
        );
    }

    <span style="color: blue;">public</span> ActionResult IndexSync() {

        <span style="color: blue;">return</span> View(<span style="color: #a31515;">"index"</span>,
            service.GetCars()
        );
    }
}</pre>
</div>
</div>
<p>We have two actions here, one is Index which is an asynchronous function and returns <strong>Task&lt;ActionResult&gt; </strong>and the second one is <strong>IndexSync</strong> which is a typical ASP.NET MVC controller action. When we navigate to <strong>/home/index</strong> and <strong>/home/indexsync</strong>, we cannot really see the difference. It takes approx. the same time.</p>
<p>In order to measure the difference, I configured a load test with Visual Studio 2010 Ultimate Load Testing features. I hit the two pages for 2 minutes. I started with 50 users and it incremented by 20 users per 5 seconds and max user limit was 500. The result was really shocking in terms of page response time.</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/77d0b156dd2b_D494/asyncFTW.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="asyncFTW" border="0" alt="asyncFTW" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/77d0b156dd2b_D494/asyncFTW_thumb.png" width="644" height="344" /></a></p>
<p>While average response time for the synchronous one is about <strong>11.2</strong> seconds, it is <strong>3.65 </strong>for asynchronous one. I think the difference is pretty compelling and overwhelming.</p>
<p>From now on, I am adopting the Windows approach: <strong>"Anything longer than 40ms (operations related to network, file system. Mainly I/O bound operations), that is async!"</strong></p>
<p>If you believe that some of the information here is wrong or misleading, please make me suffer and post a comment which would bring me down. Also, do the same if you have any additional information that is worth mentioning on this post but I missed.<strong><br /></strong></p>
<p><strong>Resources</strong></p>
<ul>
<li><a title="http://msdn.microsoft.com/en-us/vstudio/gg316360" href="http://msdn.microsoft.com/en-us/vstudio/gg316360" target="_blank">Visual Studio Asynchronous Programming</a></li>
<li><a title="http://channel9.msdn.com/Events/BUILD/BUILD2011/TOOL-816T" href="http://channel9.msdn.com/Events/BUILD/BUILD2011/TOOL-816T" target="_blank">Future directions for C# and Visual Basic (Build 2011 by&nbsp; Anders Hejlsberg)</a></li>
<li><a title="http://www.wischik.com/lu/AsyncSilverlight/AsyncSamples.html" href="http://www.wischik.com/lu/AsyncSilverlight/AsyncSamples.html" target="_blank">101 C# Async Samples</a></li>
<li><a title="http://channel9.msdn.com/Shows/AppFabric-tv/AppFabrictv-Threading-with-Jeff-Richter" href="http://channel9.msdn.com/Shows/AppFabric-tv/AppFabrictv-Threading-with-Jeff-Richter" target="_blank">AppFabric.tv - Threading with Jeff Richter</a></li>
<li><a title="http://channel9.msdn.com/Events/TechDays/Techdays-2012-the-Netherlands/2287" href="http://channel9.msdn.com/Events/TechDays/Techdays-2012-the-Netherlands/2287" target="_blank">C#5, ASP.NET MVC 4, and asynchronous Web applications</a></li>
<li><a title="http://www.youtube.com/watch?v=yhkHtXcgWUc" href="http://www.youtube.com/watch?v=yhkHtXcgWUc" target="_blank">Web performance Test and Load Test in Visual Studio 2010</a></li>
<li><a target="_blank" title="http://stackoverflow.com/questions/1453283/threadpool-in-iis-context" href="http://stackoverflow.com/questions/1453283/threadpool-in-iis-context">Threadpool in IIS context</a></li>
</ul>