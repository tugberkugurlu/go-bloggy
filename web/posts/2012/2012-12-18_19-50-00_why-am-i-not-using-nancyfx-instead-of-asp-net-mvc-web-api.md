---
id: 69079453-6382-499f-902f-b10c5ec1bcf0
title: Why am I not Using NancyFx Instead of ASP.NET MVC / Web API
abstract: 'Why am I not using NancyFx instead of ASP.NET MVC / Web API? Because of
  a very important and vital missing part with NancyFx: asynchrony!'
created_at: 2012-12-18 19:50:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET Web API
- async
- TPL
slugs:
- why-am-i-not-using-nancyfx-instead-of-asp-net-mvc-web-api
---

<blockquote>
<p><strong>Note:</strong> If you don't understand asynchrony and its benefits, you are better off without this blog post. If you are willing to learn the importance of asynchrony for server applications, here are two videos for you to start:</p>
<ul>
<li><a href="http://channel9.msdn.com/Events/TechDays/Techdays-2012-the-Netherlands/2287">C#5, ASP.NET MVC 4, and asynchronous Web applications</a></li>
<li><a href="http://channel9.msdn.com/Events/aspConf/aspConf/Async-in-ASP-NET">Async in ASP.NET</a></li>
</ul>
</blockquote>
<blockquote>
<p><strong>Update 2013-10-08</strong>:</p>
<p>Great news! NancyFx now has asyncronous processing ability. See&nbsp;Jonathan Channon's post for more details: <a href="http://blog.jonathanchannon.com/2013/08/24/async-route-handling-with-nancy/">Async Route Handling with Nancy</a>&nbsp;</p>
</blockquote>
<p>Server side .NET ecosystem is not just about the frameworks that Microsoft builds and ships such as <a href="http://www.asp.net/mvc">ASP.NET MVC</a>, <a href="http://www.asp.net/web-api">ASP.NET Web API</a>, etc. There are other open source frameworks out there such as <a href="http://nancyfx.org/">NancyFx</a> which has been embraced by the community. I haven&rsquo;t played with NancyFx that much but I really wanted to understand how it is more special than its equivalent frameworks. As far as I can tell, it&rsquo;s not but I am nowhere near to make that judgment because I only spared 10 minutes to figure that out, which is obviously not enough to judge a framework. Why and where did I stop? Let&rsquo;s get to that later.</p>
<p>You may have also noticed that I have been kind of opinionated about ASP.NET Web API lately. It&rsquo;s not because I am writing a <a href="http://www.amazon.com/dp/1430247258/ref=as_li_ss_til?tag=tugsblo0c-20&amp;camp=213381&amp;creative=390973&amp;linkCode=as4&amp;creativeASIN=1430247258&amp;adid=033D6R8D59X3A92YAMMZ&amp;&amp;ref-refURL=http%3A%2F%2Fwww.tugberkugurlu.com%2Farchive%2Fpro-asp-net-web-api-book-is-available-on-amazon-for-pre-order">book on ASP.NET Web API</a>. The reason of my passion is related to fact that I know how it works internally and I can see that none of the other equivalent frameworks hasn&rsquo;t been able to get asynchrony right as this framework did (<a href="http://signalr.net">SignalR</a> also gets asynchrony right but it&rsquo;s kind of a different area)! Please do inform me and make me happy if you think that there is any other framework which we can put under this category. I didn&rsquo;t put ASP.NET MVC under this category either because it&rsquo;s not possible to go fully asynchronous throughout the pipeline with ASP.NET MVC. However, it&rsquo;s has a really nice asynchronous support for most of the cases.</p>
<p>I&rsquo;m not going to get into the details of how asynchrony is important for a server application because <a href="https://www.tugberkugurlu.com/blog/tags/async">I have tried to do that several times</a>. Let&rsquo;s stick with the topic: why am I not using NancyFx instead of ASP.NET MVC / Web API? Don&rsquo;t get me wrong here. NancyFx is a great framework as far as I can see and I&rsquo;m all supportive for any .NET open source project which makes it easy to get our job done. <a href="https://twitter.com/TheCodeJunkie">@TheCodeJunkie</a> and the other NancyFx members have been doing a great job on maintaining this awesome open source framework. However, there is a very important missing part with this framework and it&rsquo;s vital to any server application: asynchrony! We are now about to see why this is important.</p>
<p>I have created two applications which do the same thing in a different way: downloading the data over HTTP and return it back as a text. I used NancyFx to create one of these applications and this application will make the HTTP call synchronously because there is no way to do it asynchronously with NancyFx. The other application has been implemented using ASP.NET Web API and the only difference will be that it will make the HTTP call asynchronously.</p>
<p>It&rsquo;s fairly straight forward to get started with NancyFx. I jumped to <a href="https://github.com/NancyFx/Nancy/wiki/Introduction">one of the framework&rsquo;s documentation page</a> and I was good to go by installing the NuGet package and adding the following module. It was a nice and quick start.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HelloModule : NancyModule {

    <span style="color: blue;">public</span> HelloModule() {

        Get[<span style="color: #a31515;">"/"</span>] = (p) =&gt; {

            <span style="color: blue;">using</span> (<span style="color: blue;">var</span> client = <span style="color: blue;">new</span> WebClient()) {
                <span style="color: blue;">var</span> textData = client.DownloadString(<span style="color: #a31515;">"http://localhost:52907/api/cars"</span>);
                <span style="color: blue;">return</span> textData;
            }
        };
    }
}</pre>
</div>
</div>
<p>I used <a href="http://msdn.microsoft.com/en-us/library/system.net.webclient.aspx">WebClient</a> on both applications here because <a href="http://msdn.microsoft.com/en-us/library/system.net.http.HttpClient.aspx">HttpClient</a> has no synchronous methods and making blocking calls with HttpClient has some overheads. I also created an ASP.NET Web API application which does the same thing but asynchronously.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> DefaultController : ApiController {

    <span style="color: blue;">public</span> async Task&lt;<span style="color: blue;">string</span>&gt; Get() {

        <span style="color: blue;">using</span> (<span style="color: blue;">var</span> client = <span style="color: blue;">new</span> WebClient()) {

            <span style="color: blue;">var</span> textData = await client.DownloadStringTaskAsync(
                <span style="color: #a31515;">"http://localhost:52907/api/cars"</span>);

            <span style="color: blue;">return</span> textData;
        }
    }
}</pre>
</div>
</div>
<p>As you can see, it&rsquo;s extremely easy to write asynchronous code with the ".NET 4.5 + C# 5.0 + ASP.NET Web API" combination. You probably noticed that I am making a call against http://localhost:52907/api/cars. This is a simple HTTP API which takes minimum 500 milliseconds to return. This wait is is on purpose to simulate a a-bit-long-running call. Let me also be clear on why I&rsquo;ve chosen to make a comparison by making an HTTP call. This really doesn&rsquo;t matter here. Any type of I/O operation could have been a perfect case here but an HTTP call scenario is simple to demonstrate.</p>
<p>Finally, I now have an endpoint on http://localhost:47835 for NancyFx application and another on http://localhost:49167 for ASP.NET Web API application. I have opened up a command line console and get the <a href="http://httpd.apache.org/docs/2.2/programs/ab.html">ab.exe (Apache HTTP Server Benchmarking Tool)</a> under my path. When I generate a request for each one (only one request), I will have a result which is similar to the below screenshot (NancyFx is the lest side, ASP.NET Web API is the right side).</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Why-am-I-not-using-Nanc.NET-MVC--Web-API_116FA/image.png"><img height="291" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Why-am-I-not-using-Nanc.NET-MVC--Web-API_116FA/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="image" /></a></p>
<p>Notice that both have taken approx. the same amount of time to complete the request but this&rsquo;s not the place where asynchrony shines. Let&rsquo;s now generate totally 200 request (30 concurrent requests at a time) to each endpoint and see the result (NancyFx is the lest side, ASP.NET Web API is the right side).</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Why-am-I-not-using-Nanc.NET-MVC--Web-API_116FA/image_3.png"><img height="313" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Why-am-I-not-using-Nanc.NET-MVC--Web-API_116FA/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="image" /></a></p>
<p>15.31 requests per second with NancyFx + synchronous processing and 49.35 requests per second with ASP.NET Web API + asynchronous processing. Besides the fact that this&rsquo;s totally an on-the-fly benchmarking (I&rsquo;m not even sure whether we can call this a benchmark) on my Windows 8 machine, the result is pretty compelling.</p>
<p>Let&rsquo;s now have a look at the reason why I stopped discovering NancyFx further. Go to <a href="https://github.com/NancyFx/Nancy">NancyFx GitHub repository</a> and navigate to <a href="https://github.com/NancyFx/Nancy/blob/master/src/Nancy.Hosting.Aspnet/NancyHttpRequestHandler.cs">Nancy.Hosting.Aspnet.NancyHttpRequestHandler</a> class file. This is located under the ASP.NET host project for NancyFx and actually the HTTP handler the framework uses to process requests. Earlier today I tweeted after I had a look at this class:</p>
<blockquote class="twitter-tweet tw-align-left">
<p>If your <a href="https://twitter.com/search/%23aspnet">#aspnet</a> framework implements IHttpHandler.ProcessRequest method, you are doing it wrong. <a href="https://twitter.com/search/%23JustSayin">#JustSayin</a></p>
&mdash; Tugberk Ugurlu (@tourismgeek) <a data-datetime="2012-12-18T13:22:28+00:00" href="https://twitter.com/tourismgeek/status/281026645199040512">December 18, 2012</a></blockquote>
<script src="//platform.twitter.com/widgets.js"></script>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Why-am-I-not-using-Nanc.NET-MVC--Web-API_116FA/image_4.png"><img height="321" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Why-am-I-not-using-Nanc.NET-MVC--Web-API_116FA/image_thumb_4.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Yes, it&rsquo;s sitting right there. By exposing the framework with a normal handler implementation instead of the <a href="http://msdn.microsoft.com/en-us/library/system.web.IHttpAsyncHandler.aspx">IHttpAsyncHandler</a> implementation, you are making your ASP.NET framework suffer. It&rsquo;s possible for you to handle synchronous calls with an asynchronous handler but it&rsquo;s impossible (or either not right and easy) to do this with a normal handler. On the other hand, we know that NancyFx doesn&rsquo;t only have an ASP.NET host option. You can even host this framework under <a href="http://katanaproject.codeplex.com/">katana</a> with its OWIN Host implementation but it&rsquo;s not going to matter unless the framework itself has asynchronous processing capabilities and by looking at the above handler and <a href="https://github.com/NancyFx/Nancy/issues/148">this issue</a>, I can see that it doesn&rsquo;t have that.</p>
<p>This is not an I-hate-you-and-I-will-stay-that-way-forever-please-go-the-hell-away blog post. Again: <strong>NancyFx is a great project and it deserves a good amount of respect.</strong> However, this post explains my reason on why I am not using NancyFx instead of ASP.NET Web API or ASP.NET MVC. Will it going to change? Unless <a href="https://github.com/NancyFx/Nancy/issues/148">#148</a> is not fixed, it is certainly not. Even if it is fixed, I already have applications running on ASP.NET Web API and ASP.NET MVC. Why bother then? We will see that in the future.</p>