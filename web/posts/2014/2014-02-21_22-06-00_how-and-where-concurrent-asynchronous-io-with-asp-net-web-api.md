---
title: How and Where Concurrent Asynchronous I/O with ASP.NET Web API
abstract: When we have uncorrelated multiple I/O operations that need to be kicked
  off, we have quite a few ways to fire them off and which way you choose makes a
  great amount of difference on a .NET server side application. In this post, we will
  see how we can handle the different approaches in ASP.NET Web API.
created_at: 2014-02-21 22:06:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
- async
- C#
- Concurrency
slugs:
- how-and-where-concurrent-asynchronous-io-with-asp-net-web-api
---

<p>When we have uncorrelated multiple I/O operations that need to be kicked off, we have quite a few ways to fire them off and which way you choose makes a great amount of difference on a .NET server side application. <a href="http://weblogs.asp.net/cibrax/default.aspx">Pablo Cibraro</a> already has a great post on this topic (<a href="http://weblogs.asp.net/cibrax/archive/2012/11/15/await-whenall-waitall-oh-my.aspx">await, WhenAll, WaitAll, oh my!!</a>) which I recommend you to check that out. In this article, I would like to touch on a few more points. Let's look at the options one by one. I will use a multiple HTTP request scenario here which will be consumed by an <a href="http://www.asp.net/web-api">ASP.NET Web API</a> application but this is applicable for any sort of I/O operations (long-running database calls, file system operations, etc.).</p> <p>We will have two different endpoint which will hit to consume the data:</p> <ul> <li>http://localhost:2700/api/cars/cheap  <li>http://localhost:2700/api/cars/expensive</li></ul> <p>As we can infer from the URI, one of them will get us the cheap cars and the other one will get us the expensive ones. I created a separate ASP.NET Web API application to simulate these endpoints. Each one takes more than 500ms to complete and in our target ASP.NET Web API application, we will aggregate these two resources together and return the result. Sounds like a very common scenario.  <p>Inside our target API controller, we have the following initial structure:  <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> Car 
{
    <span style="color: blue">public</span> <span style="color: blue">int</span> Id { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">string</span> Make { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">string</span> Model { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">int</span> Year { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">float</span> Price { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
}

<span style="color: blue">public</span> <span style="color: blue">class</span> CarsController : BaseController 
{
    <span style="color: blue">private</span> <span style="color: blue">static</span> <span style="color: blue">readonly</span> <span style="color: blue">string</span>[] PayloadSources = <span style="color: blue">new</span>[] { 
        <span style="color: #a31515">"http://localhost:2700/api/cars/cheap"</span>,
        <span style="color: #a31515">"http://localhost:2700/api/cars/expensive"</span>
    };

    <span style="color: blue">private</span> async Task&lt;IEnumerable&lt;Car&gt;&gt; GetCarsAsync(<span style="color: blue">string</span> uri) 
    {
        <span style="color: blue">using</span> (HttpClient client = <span style="color: blue">new</span> HttpClient()) 
        {
            <span style="color: blue">var</span> response = await client.GetAsync(uri).ConfigureAwait(<span style="color: blue">false</span>);
            <span style="color: blue">var</span> content = await response.Content
                .ReadAsAsync&lt;IEnumerable&lt;Car&gt;&gt;().ConfigureAwait(<span style="color: blue">false</span>);

            <span style="color: blue">return</span> content;
        }
    }

    <span style="color: blue">private</span> IEnumerable&lt;Car&gt; GetCars(<span style="color: blue">string</span> uri) 
    {
        <span style="color: blue">using</span> (WebClient client = <span style="color: blue">new</span> WebClient()) 
        {    
            <span style="color: blue">string</span> carsJson = client.DownloadString(uri);
            IEnumerable&lt;Car&gt; cars = JsonConvert
                .DeserializeObject&lt;IEnumerable&lt;Car&gt;&gt;(carsJson);
                
            <span style="color: blue">return</span> cars;
        }
    }
}
</pre></div></div>
<p>We have a Car class which will represent a car object that we are going to deserialize from the JSON payload. Inside the controller, we have our list of endpoints and two private methods which are responsible to make HTTP GET requests against the specified URI. GetCarsAsync method uses the System.Net.Http.<a href="http://msdn.microsoft.com/en-us/library/system.net.http.httpclient(v=vs.110).aspx">HttpClient</a> class, which has been introduces with .NET 4.5, to make the HTTP calls asynchronously. With the new C# 5.0 asynchronous language features (A.K.A async modifier and await operator), it is pretty straight forward to write the asynchronous code as you can see. Note that we used ConfigureAwait method here by passing the false Boolean value for <a href="http://msdn.microsoft.com/en-us/library/hh194876(v=vs.110).aspx">continueOnCapturedContext</a> parameter. It’s a quite long topic why we need to do this here but briefly, one of our samples, which we are about to go deep into, would introduce deadlock if we didn’t use this method. 
<p>To be able to measure the performance, we will use a little utility tool from <a href="http://httpd.apache.org/docs/2.2/programs/ab.html">Apache Benchmarking Tool (A.K.A ab.exe)</a>. This comes with Apache Web Server installation but you don’t actually need to install it. When you download the necessary ZIP file for the installation and extract it, you will find the ab.exe inside. Alternatively, you may use <a href="http://www.iis.net/downloads/community/2007/05/wcat-63-(x64)">Web Capacity Analysis Tool (WCAT)</a> from IIS team. It’s a lightweight HTTP load generation tool primarily designed to measure the performance of a web server within a controlled environment. However, WCAT is a bit hard to grasp and set up. That’s why we used ab.exe here for simple load tests. 
<blockquote>
<p>Please, note that the below compressions are poor and don't indicate any real benchmarking. These are just compressions for demo purposes and they indicate the points that we are looking for.</p></blockquote>
<h3>Synchronous and not In Parallel</h3>
<p>First, we will look at all synchronous and not in parallel version of the code. This operation will block the running the thread for the amount of time which takes to complete two network I/O operations. The code is very simple thanks to LINQ. 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>[HttpGet]
<span style="color: blue">public</span> IEnumerable&lt;Car&gt; AllCarsSync() {

    IEnumerable&lt;Car&gt; cars =
        PayloadSources.SelectMany(x =&gt; GetCars(x));

    <span style="color: blue">return</span> cars;
}
</pre></div></div>
<p>For a single request, we expect this to complete for about a second. 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/90f2ab25-dbe7-4855-9b9c-4bd5c6bcfdea.png"><img title="AllCarsSync" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsSync" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/fbe980a7-6f07-436a-89f0-73759fd38ea1.png" width="452" height="484"></a> 
<p>The result is not surprising. However, when you have multiple concurrent requests against this endpoint, you will see that the blocking threads will be the bottleneck for your application. The following screenshot shows the 200 requests to this endpoint in 50 requests blocks. 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/cde5aa76-bcf9-4022-be51-9b7d0d1f7ac9.png"><img title="AllCarsSync_200" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsSync_200" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ea917c63-60d9-4058-89dd-0e3a7054af17.png" width="450" height="484"></a> 
<p>The result is now worse and we are paying the price for blowing the threads for long running I/O operations. You may think that running these in-parallel will reduce the single request time and you are not wrong but this has its own caveats, which is our next section. 
<h3>Synchronous and In Parallel</h3>
<p>This option is mostly never good for your application. With this option, you will perform the I/O operations in parallel and the request time will be significantly reduced if you try to measure only with one request. However, in our sample case here, you will be consuming two threads instead of one to process the request and you will block both of them while waiting for the HTTP requests to complete. Although this reduces the overall request processing time for a single request, it consumes more resources and you will see that the overall request time increases while your request count increases. Let’s look at the code of the ASP.NET Web API controller action method. 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>[HttpGet]
<span style="color: blue">public</span> IEnumerable&lt;Car&gt; AllCarsInParallelSync() {

    IEnumerable&lt;Car&gt; cars = PayloadSources.AsParallel()
        .SelectMany(uri =&gt; GetCars(uri)).AsEnumerable();

    <span style="color: blue">return</span> cars;
}
</pre></div></div>
<p>We used “Parallel LINQ (PLINQ)” feature of .NET framework here to process the HTTP requests in parallel. As you can, it was just too easy; in fact, it was only one line of digestible code. I tent to see a relationship between the above code and tasty donuts. They all look tasty but they will work as hard as possible to clog our carotid arteries. Same applies to above code: it looks really sweet but can make our server application miserable. How so? Let’s send a request to this endpoint to start seeing how. 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ce25e957-59ed-45a8-a23c-8dd17eeddb6a.png"><img title="AllCarsInParallelSync" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsInParallelSync" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/95acf6c9-0523-4bff-a31e-cc1b272c4fc0.png" width="450" height="484"></a> 
<p>As you can see, the overall request time has been reduced in half. This must be good, right? Not completely. As mentioned before, this is going to hurt us if we see too many requests coming to this endpoint. Let’s simulate this with ab.exe and send 200 requests to this endpoint in 50 requests blocks. 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/49d64535-cac0-406d-b42b-3cbc86f82e42.png"><img title="AllCarsInParallelSync_200" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsInParallelSync_200" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/590f550d-c7bf-496a-b164-2a9e7496e904.png" width="452" height="484"></a> 
<p>The overall performance is now significantly reduced. So, where would this type of implementation make sense? If your server application has small number of users (for example, an HTTP API which consumed by the internal applications within your small team), this type of implementation may give benefits. However, as it’s now annoyingly simple to write asynchronous code with built-in language features, I’d suggest you to choose our last option here: “Asynchronous and In Parallel (In a Non-Blocking Fashion)”. 
<h3>Asynchronous and not In Parallel</h3>
<p>Here, we won’t introduce any concurrent operations and we will go through each request one by one but in an asynchronous manner so that the processing thread will be freed up during the dead waiting period. 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>[HttpGet]
<span style="color: blue">public</span> async Task&lt;IEnumerable&lt;Car&gt;&gt; AllCarsAsync() {

    List&lt;Car&gt; carsResult = <span style="color: blue">new</span> List&lt;Car&gt;();
    <span style="color: blue">foreach</span> (<span style="color: blue">var</span> uri <span style="color: blue">in</span> PayloadSources) {

        IEnumerable&lt;Car&gt; cars = await GetCarsAsync(uri);
        carsResult.AddRange(cars);
    }

    <span style="color: blue">return</span> carsResult;
}
</pre></div></div>
<p>What we do here is quite simple: we are iterating through the URI array and making the asynchronous HTTP call for each one. Notice that we were able to use the await keyword inside the foreach loop. This is all fine. The compiler will do the right thing and handle this for us. One thing to keep in mind here is that the asynchronous operations won’t run in parallel here. So, we won’t see a difference when we send a single request to this endpoint as we are going through the each request one by one. 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/06248b96-18e6-491a-adde-d5713b01572b.png"><img title="AllCarsAsync" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsAsync" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c4f304f0-ffee-4648-8b72-1dddad27bf4c.png" width="450" height="484"></a> 
<p>As expected, it took around a second. When we increase the number of requests and concurrency level, we will see that the average request time still stays around a second to perform. 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5d5b1e24-3faf-45d6-8f2c-4fc930f562fb.png"><img title="AllCarsAsync_200" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsAsync_200" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/646d0c68-7409-4c47-b364-22e14d53a1f7.png" width="449" height="484"></a> 
<p>This option is certainly better than the previous ones. However, we can still do better in some certain cases where we have limited number of concurrent I/O operations. The last option will look into this solution but before moving onto that, we will look at one other option which should be avoided where possible. 
<p>Asynchronous and In Parallel (In a Blocking Fashion)</p>
<p>Among these options shown here, this is the worst one that one can choose. When we have multiple Task returning asynchronous methods in our hand, we can wait all of them to finish with <a href="http://msdn.microsoft.com/en-us/library/system.threading.tasks.task.waitall(v=vs.110).aspx">WaitAll</a> static method on Task object. This results several overheads: you will be consuming the asynchronous operations in a blocking fashion and if these asynchronous methods is not implemented right, you will end up with deadlocks. At the beginning of this article, we have pointed out the usage of <a href="http://msdn.microsoft.com/en-us/library/hh194876(v=vs.110).aspx">ConfigureAwait</a> method. This was for preventing the deadlocks here. You can learn more about this from the following blog post: <a href="http://www.tugberkugurlu.com/archive/asynchronousnet-client-libraries-for-your-http-api-and-awareness-of-async-await-s-bad-effects">Asynchronous .NET Client Libraries for Your HTTP API and Awareness of async/await's Bad Effects</a>. 
<p>Let’s look at the code: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>[HttpGet]
<span style="color: blue">public</span> IEnumerable&lt;Car&gt; AllCarsInParallelBlockingAsync() {
    
    IEnumerable&lt;Task&lt;IEnumerable&lt;Car&gt;&gt;&gt; allTasks = 
        PayloadSources.Select(uri =&gt; GetCarsAsync(uri));

    Task.WaitAll(allTasks.ToArray());
    <span style="color: blue">return</span> allTasks.SelectMany(task =&gt; task.Result);
}
</pre></div></div>
<p>Let's send a request to this endpoint to see how it performs: 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6b0aa65e-2cc2-419a-a177-deb4d02532f2.png"><img title="AllCarsInParallelBlockingAsync" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsInParallelBlockingAsync" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1d73aa7e-732c-4548-9c4e-a0e14c326d90.png" width="451" height="484"></a> 
<p>It performed really bad but it gets worse as soon as you increase the concurrency rate: 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ac0d65e0-cd6f-4d43-826f-4147e79c2d69.png"><img title="AllCarsInParallelBlockingAsync_200" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsInParallelBlockingAsync_200" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/934120b0-5929-4378-8e4e-633d82bc13f9.png" width="451" height="484"></a> 
<p>Never, ever think about implementing this solution. No further discussion is needed here in my opinion. 
<p>Asynchronous and In Parallel (In a Non-Blocking Fashion)</p>
<p>Finally, the best solution: <strong>Asynchronous and In Parallel (In a Non-Blocking Fashion)</strong>. The below code snippet indicates it all but just to go through it quickly, we are bundling the Tasks together and await on the <a href="http://msdn.microsoft.com/en-us/library/system.threading.tasks.task.whenall(v=vs.110).aspx">Task.WhenAll</a> utility method. This will perform the operations asynchronously in Parallel.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white">
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>[HttpGet]
<span style="color: blue">public</span> async Task&lt;IEnumerable&lt;Car&gt;&gt; AllCarsInParallelNonBlockingAsync() {

    IEnumerable&lt;Task&lt;IEnumerable&lt;Car&gt;&gt;&gt; allTasks = PayloadSources.Select(uri =&gt; GetCarsAsync(uri));
    IEnumerable&lt;Car&gt;[] allResults = await Task.WhenAll(allTasks);

    <span style="color: blue">return</span> allResults.SelectMany(cars =&gt; cars);
}
</pre></div></div></div></div>
<p>If we make a request to the endpoint to execute this piece of code, the result will be similar to the previous one:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d4012a32-7e6c-44e6-9c4c-8c8a435321ab.png"><img title="AllCarsInParallelNonBlockingAsync" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsInParallelNonBlockingAsync" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/fc6712b7-e283-410e-89d7-a31ce01dd37c.png" width="452" height="484"></a></p>
<p>However, when we make 50 concurrent requests 4 times, the result will shine and lays out the advantages of asynchronous I/O handling:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a8b134fd-7cc1-4382-a920-9ed6a677e043.png"><img title="AllCarsInParallelNonBlockingAsync_200" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="AllCarsInParallelNonBlockingAsync_200" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c6c10e08-8f71-4615-8727-10472a8172f7.png" width="450" height="484"></a></p>
<h3>Conclusion</h3>
<p>At the very basic level, what we can get out from this article is this: do perform load tests against your server applications based on your estimated consumption rates if you have any sort of multiple I/O operations. Two of the above options are what you would want in case of multiple I/O operations. One of them is <strong>"Asynchronous but not In Parallel"</strong>,<strong> </strong>which is the safest option in my personal opinion, and the other is <strong>"Asynchronous and In Parallel (In a Non-Blocking Fashion)"</strong>. The latter option significantly reduces the request time depending on the hardware and number of I/O operations you have but as our small benchmarking results showed, it may not be a good fit to process a request many concurrent I/O asynchronous operations in one just to reduce a single request time. The result we would see will most probably be different under high load. 
<h3>References</h3>
<ul>
<li><a href="http://msdn.microsoft.com/en-us/library/windows/desktop/aa365683(v=vs.85).aspx">Synchronous and Asynchronous I/O (Windows)</a>
<li><a href="http://stackoverflow.com/questions/12337671/using-async-await-for-multiple-tasks">Using async/await for multiple tasks</a>
<li><a href="http://msdn.microsoft.com/en-us/library/windows/desktop/aa365683(v=vs.85).aspx">Synchronous and Asynchronous I/O (Windows)</a>
<li><a href="http://msdn.microsoft.com/en-us/library/hh156548.aspx">Parallel Processing and Concurrency in the .NET Framework</a>
<li><a href="http://httpd.apache.org/docs/2.2/programs/ab.html">ab - Apache HTTP server benchmarking tool</a>
<li><a href="http://www.iis.net/downloads/community/2007/05/wcat-63-(x64)">WCat 6.3 (x64)</a></li></ul>  