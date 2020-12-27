---
id: b8bf99f0-a1c5-440c-a3fb-4cb70ea2e494
title: ASP.NET 5 and Log Correlation by Request Id
abstract: ASP.NET 5 is full of big new features and enhancements but besides these,
  I am mostly impressed by little, tiny features of ASP.NET 5 Log Correlation which
  is provided out of the box. Let me quickly show you what it is in this post.
created_at: 2015-10-28 00:44:00 +0000 UTC
tags:
- .NET
- ASP.Net
- ASP.NET 5
- Elasticsearch
slugs:
- asp-net-5-and-log-correlation-by-request-id
---

<p><a href="https://www.tugberkugurlu.com/tags/asp-net-5">ASP.NET 5</a> is full of big new features and enhancements like being able to run on multiple operating systems, incredible CLI tools, hassle-free building for multiple framework targets, <a href="https://www.tugberkugurlu.com/archive/exciting-things-about-asp-net-5-series-build-only-dependencies">build only dependencies</a> and many more. Besides these, I am mostly impressed by little, tiny features of ASP.NET 5 because these generally tend to be ignored in this type of rearchitecturing works. One of these little features is <a href="https://github.com/aspnet/Home/issues/1015">log correlation</a>. Let me quickly show you what it is and why it made me smile. </p> <blockquote> <p><strong>BIG ASS CAUTION!</strong> At the time of this writing, I am using <strong>DNX 1.0.0-beta8 </strong>version. As things are moving really fast in this new world, it’s very likely that the things explained here will have been changed as you read this post. So, be aware of this and try to explore the things that are changed to figure out what are the corresponding new things.  <p>Also, inside this post I am referencing a lot of things from ASP.NET GitHub repositories. In order to be sure that the links won’t break in the future, I’m actually referring them by getting permanent links to the files on GitHub. So, these links are actually referring the files from the latest commit at the time of this writing and they have a potential to be changed, too. Read the "<a href="https://help.github.com/articles/getting-permanent-links-to-files/">Getting permanent links to files</a>" post to figure what this actually is.</p></blockquote> <h3>Brief Introduction to Logging in ASP.NET 5 World</h3> <blockquote> <p>If you want to skip this part, you can directly go to "Log Correlation" section below.</p></blockquote> <p>As you probably know, ASP.NET 5 also has a great support for <a href="https://github.com/aspnet/Logging/tree/1.0.0-beta8">logging</a>. The nicest thing about this new logging abstraction is that it’s the only logging abstraction which every provided library and framework is relying on. So, when you enable logging in your application, you will enable it in all components (which is perfect)! Here is a sample in my MVC 6 application. I am just adding MVC to pipeline here, enabling logging by hooking <a href="https://github.com/serilog/serilog">Serilog</a> and configuring it to write the logs to console:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">using</span> System;
<span style="color: blue">using</span> Microsoft.AspNet.Builder;
<span style="color: blue">using</span> Microsoft.Framework.DependencyInjection;
<span style="color: blue">using</span> Microsoft.Framework.Logging;
<span style="color: blue">using</span> Serilog;

<span style="color: blue">namespace</span> LoggingCorrelationSample
{
    <span style="color: blue">public</span> <span style="color: blue">class</span> Startup
    {
        <span style="color: blue">public</span> Startup(ILoggerFactory loggerFactory)
        {
            <span style="color: blue">var</span> serilogLogger = <span style="color: blue">new</span> LoggerConfiguration()
                .WriteTo
                .TextWriter(Console.Out)
                .MinimumLevel.Verbose()
                .CreateLogger();

            loggerFactory.MinimumLevel = LogLevel.Debug;
            loggerFactory.AddSerilog(serilogLogger);
        }

        <span style="color: blue">public</span> <span style="color: blue">void</span> ConfigureServices(IServiceCollection services)
        {
            services.AddMvc();
        }

        <span style="color: blue">public</span> <span style="color: blue">void</span> Configure(IApplicationBuilder app)
        {
            app.UseMvc();
        }
    }
}</pre></div></div>
<p>When I run the application and hit a valid endpoint, I will see bunch of things being logged to console:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/eee008b4-bc00-4bc5-9633-271a61fbf066.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/24112851-301a-422f-8d7a-26a3885a811b.png" width="644" height="476"></a></p>
<p>Remember, I haven’t logged anything myself yet. It’s just the stuff I hooked in which were already relying on ASP.NET 5 logging infrastructure. This doesn’t mean I can’t though. Hooking into logging is super easy since an instance of <a href="https://github.com/aspnet/Logging/blob/1.0.0-beta8/src/Microsoft.Framework.Logging.Abstractions/ILoggerFactory.cs">ILoggerFactory</a> is already inside the DI system. Here is an example class which I have for my application and it is responsible for getting the cars (forgive the stupid example here but I am sure you will get the idea):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> CarsContext : IDisposable
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> ILogger _logger;

    <span style="color: blue">public</span> CarsContext(ILoggerFactory loggerFactory)
    {
        _logger = loggerFactory.CreateLogger&lt;CarsContext&gt;();
        _logger.LogDebug(<span style="color: #a31515">"Constructing CarsContext"</span>);
    }

    <span style="color: blue">public</span> IEnumerable&lt;<span style="color: blue">string</span>&gt; GetCars()
    {
        _logger.LogInformation(<span style="color: #a31515">"Found 3 cars."</span>);
        
        <span style="color: blue">return</span> <span style="color: blue">new</span>[]
        {
            <span style="color: #a31515">"Car 1"</span>,
            <span style="color: #a31515">"Car 2"</span>,
            <span style="color: #a31515">"Car 3"</span>
        };
    }
    
    <span style="color: blue">public</span> <span style="color: blue">void</span> Dispose()
    {
        _logger.LogDebug(<span style="color: #a31515">"Disposing CarsContext"</span>);
    }
}</pre></div></div>
<p>I will register this class so that it can get the dependencies it needs and also, it can be injected into other places:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">void</span> ConfigureServices(IServiceCollection services)
{
    services.AddMvc();
    services.AddScoped&lt;CarsContext, CarsContext&gt;();
}</pre></div></div>
<p>Finally, I will use it inside my controller:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> CarsController : Controller
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> CarsContext _carsContext;
    
    <span style="color: blue">public</span> CarsController(CarsContext carsContext)
    {
        _carsContext = carsContext;
    }
    
    [Route(<span style="color: #a31515">"cars"</span>)]
    <span style="color: blue">public</span> IActionResult Get()
    {
        <span style="color: blue">var</span> cars = _carsContext.GetCars();
        <span style="color: blue">return</span> Ok(cars);
    }
}</pre></div></div>
<p>Just seeing how beautifully things are coming together is really great! When I run the application and hit the /cars endpoint now, I will see my logs appearing along side the framework and library logs:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/626bd661-84b5-4900-b06b-6f800dcfd877.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ba08b08e-b448-4667-8699-ff6c86b3e167.png" width="644" height="476"></a></p>
<p>Same goes for your middlewares. You can naturally hook into logging system from your middleware thanks to <a href="https://www.tugberkugurlu.com/archive/exciting-things-about-asp-net-vnext-series-middlewares-and-per-request-dependency-injection">first class middleware DI support</a>.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> RequestUrlLoggerMiddleware 
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> RequestDelegate _next;
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> Microsoft.Framework.Logging.ILogger _logger;
    
    <span style="color: blue">public</span> RequestUrlLoggerMiddleware(RequestDelegate next, ILoggerFactory loggerFactory) 
    {
        _next = next;
        _logger = loggerFactory.CreateLogger&lt;RequestUrlLoggerMiddleware&gt;();
    }
    
    <span style="color: blue">public</span> Task Invoke (HttpContext context)
    {
        _logger.LogInformation(<span style="color: #a31515">"{Method}: {Url}"</span>, context.Request.Method, context.Request.Path);
        <span style="color: blue">return</span> _next(context);
    }
}</pre></div></div>
<blockquote>
<p>Notice that we have a log <strong>message template</strong> rather the <strong>actual log message</strong> here. This is another great feature of the new logging system which is pretty much the same as <a href="https://github.com/serilog/serilog/wiki/Structured-Data">what Serilog have had for log time</a>.</p></blockquote>
<p>When we run this, we should see the middleware log appear, too:</p>

<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b7ef6163-1adc-4171-8611-d5abe01dbca1.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/25b24177-d068-4136-b5cf-cbe1558cb455.png" width="644" height="385"></a></p>
<h3>Log Correlation</h3>
<p>Without doing anything else first, let me also write logs to <a href="https://www.elastic.co/products/elasticsearch">Elasticsearch</a> by pulling in <a href="https://github.com/serilog/serilog-sinks-elasticsearch">Serilog Elasticsearch sink</a> and <a href="https://github.com/tugberkugurlu/AspNetVNextSamples/blob/c6fb3c946c48ff0fff3ed83cfd5e34de339efc0b/beta8/LoggingCorrelationSample/Startup.cs#L20-L22">hooking it in</a>. After hitting the same endpoint, I have the below result inside my Elasticsearch index:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/3f465e8d-84d3-4622-b43b-dad56f8d73d9.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a4653a4d-1773-4381-90d3-c9a0855cbaff.png" width="644" height="436"></a></p>
<p>You can see that each log message has got richer and we can see new things like RequestId which will allow you to correlate your logs per request. This information is being logged because <a href="https://github.com/aspnet/Hosting/blob/9be0758c4d1da87aa5afe931e030d59636d33cb6/src/Microsoft.AspNet.Hosting/Internal/HostingEngine.cs#L106-L128">the hosting layer starts a new log scope for each request</a>. RequestId is particularly useful when you have an unexpected behavior with an HTTP request and you want to see what was happening with that request. In order to take advantage this, you should send the the RequestId along side your response (ideally among the response headers). The below is a sample middleware which you can hook into your pipeline in order to add RequestId to your response:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> RequestIdMiddleware
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> RequestDelegate _next;

    <span style="color: blue">public</span> RequestIdMiddleware(RequestDelegate next)
    {
        _next = next;
    }

    <span style="color: blue">public</span> async Task Invoke(HttpContext context)
    {
        <span style="color: blue">var</span> requestIdFeature = context.Features.Get&lt;IHttpRequestIdentifierFeature&gt;();
        <span style="color: blue">if</span> (requestIdFeature?.TraceIdentifier != <span style="color: blue">null</span>)
        {
            context.Response.Headers[<span style="color: #a31515">"RequestId"</span>] = requestIdFeature.TraceIdentifier;
        }

        await _next(context);
    }
}</pre></div></div>
<blockquote>
<p>Note that, <a href="https://github.com/aspnet/HttpAbstractions/blob/1.0.0-beta8/src/Microsoft.AspNet.Http.Features/IHttpRequestIdentifierFeature.cs">IHttpRequestIdentifierFeature</a> is the way to get a hold of RequestId in beta8 but in upcoming versions, <a href="https://github.com/aspnet/HttpAbstractions/commit/e01a05d21439815370a6b82682839390b7c57a24">it’s likely to change</a> to <a href="https://github.com/aspnet/HttpAbstractions/blob/e01a05d21439815370a6b82682839390b7c57a24/src/Microsoft.AspNet.Http.Abstractions/HttpContext.cs#L37">HttpContext.TraceIdentifier</a>.</p></blockquote>
<p>If you look at the response headers now, you should see the RequestId header there:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>HTTP/1.1 200 OK
Date: Wed, 28 Oct 2015 00:32:22 GMT
Content-Type: application/json; charset=utf-8
Server: Kestrel
RequestId: 0b66784c-eb98-4a53-9247-8563fad85857
Transfer-Encoding: chunked</pre></div></div>
<p>Assuming that I have problems with this request and I have been handed the RequestId, I should be able to see what happened in that request by running a simple query on my Elasticsearch index:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/7cc74b79-94b5-4980-9664-11960c3d8afa.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6eb2b934-73f1-41c2-842c-b36bcd0b9c14.png" width="644" height="412"></a></p>
<p>That’s pretty much it and as mentioned, this is one those tiny features which was always possible but painful to get it all right. If you are also interested, you can find the <a href="https://github.com/tugberkugurlu/AspNetVNextSamples/tree/c6fb3c946c48ff0fff3ed83cfd5e34de339efc0b/beta8/LoggingCorrelationSample">full source code of the sample under my GitHub repository</a>.</p>  