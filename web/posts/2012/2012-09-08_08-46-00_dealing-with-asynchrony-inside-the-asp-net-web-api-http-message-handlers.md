---
id: afb2b340-f594-4c72-99d1-85848a29eee3
title: Dealing with Asynchrony inside the ASP.NET Web API HTTP Message Handlers
abstract: How to most efficiently deal with asynchrony inside ASP.NET Web API HTTP
  Message Handlers with TaskHelpers NuGet package or C# 5.0 asynchronous language
  features
created_at: 2012-09-08 08:46:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
- async
- TPL
slugs:
- dealing-with-asynchrony-inside-the-asp-net-web-api-http-message-handlers
---

<p>ASP.NET Web API has a concept of Message Handlers for processing HTTP messages on both the client and server. If we take look under the hood inside the framework, we will see that Initialize method of the HttpServer instance is invoking CreatePipeline method of System.Net.Http.HttpClientFactory&rsquo;s. CreatePipeline method accepts two parameters: HttpMessageHandler and IEnumerable&lt;DelegatingHandler&gt;.HttpServer.Initialize method is passing System.Web.Http.Dispatcher.HttpControllerDispatcher for HttpMessageHandlerparameter as the last HttpMessageHandler inside the chain and HttpConfiguration.MessageHandlers for IEnumerable&lt;DelegatingHandler&gt; parameter.</p>
<p>What happens inside the CreatePipeline method is very clever IMO:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">static</span> HttpMessageHandler CreatePipeline(
    HttpMessageHandler innerHandler, 
    IEnumerable&lt;DelegatingHandler&gt; handlers) {
    
    <span style="color: blue;">if</span> (innerHandler == <span style="color: blue;">null</span>)
    {
        <span style="color: blue;">throw</span> Error.ArgumentNull(<span style="color: #a31515;">"innerHandler"</span>);
    }

    <span style="color: blue;">if</span> (handlers == <span style="color: blue;">null</span>)
    {
        <span style="color: blue;">return</span> innerHandler;
    }

    <span style="color: green;">// Wire handlers up in reverse order starting with the inner handler</span>
    HttpMessageHandler pipeline = innerHandler;
    IEnumerable&lt;DelegatingHandler&gt; reversedHandlers = handlers.Reverse();
    <span style="color: blue;">foreach</span> (DelegatingHandler handler <span style="color: blue;">in</span> reversedHandlers)
    {
        <span style="color: blue;">if</span> (handler == <span style="color: blue;">null</span>)
        {
            <span style="color: blue;">throw</span> Error.Argument(<span style="color: #a31515;">"handlers"</span>, 
                Properties.Resources.DelegatingHandlerArrayContainsNullItem, 
                <span style="color: blue;">typeof</span>(DelegatingHandler).Name);
        }

        <span style="color: blue;">if</span> (handler.InnerHandler != <span style="color: blue;">null</span>)
        {
            <span style="color: blue;">throw</span> Error.Argument(<span style="color: #a31515;">"handlers"</span>, 
                Properties.Resources
                .DelegatingHandlerArrayHasNonNullInnerHandler, 
                <span style="color: blue;">typeof</span>(DelegatingHandler).Name, 
                <span style="color: #a31515;">"InnerHandler"</span>, 
                handler.GetType().Name);
        }

        handler.InnerHandler = pipeline;
        pipeline = handler;
    }

    <span style="color: blue;">return</span> pipeline;
}</pre>
</div>
</div>
<p>As you can see, the message handler order is reversed and the <a href="http://en.wikipedia.org/wiki/Matryoshka_doll">Matryoshka Doll</a> is created but be careful here: it is ensured that HttpControllerDispatcher is the last message handler to run inside the chain. Also there is a bit of a misunderstanding about message handlers, IMHO, about calling the handler again on the way back to the client. This is actually not quite correct because the message handler&rsquo;s SendAsync method itself won&rsquo;t be called twice, the continuation delegate that you will chain onto SendAsync method will be invoked on the way back to client with the generated response message that you can play with.</p>
<p>For example, let's assume that the following two are the message handler:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> XMagicMessageHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request, 
        CancellationToken cancellationToken) {
            
        <span style="color: green;">//Play with the request here</span>

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, cancellationToken)
            .ContinueWith(task =&gt; {

                <span style="color: green;">//inspect the generated response</span>
                <span style="color: blue;">var</span> response = task.Result;

                <span style="color: green;">//Add the X-Magic header</span>
                response.Headers.Add(<span style="color: #a31515;">"X-Magic"</span>, <span style="color: #a31515;">"ThisIsMagic"</span>);

                <span style="color: blue;">return</span> response;
        });
    }
}</pre>
</div>
</div>
<p>Besides this I have a tiny controller which serves cars list if the request is a GET request and I registered this message handler through the GlobalConfiguration object as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    <span style="color: blue;">var</span> config = GlobalConfiguration.Configuration;
    config.Routes.MapHttpRoute(<span style="color: #a31515;">"DefaultHttpRoute"</span>, <span style="color: #a31515;">"api/{controller}"</span>);
    config.MessageHandlers.Add(<span style="color: blue;">new</span> XMagicMessageHandler());
    config.MessageHandlers.Add(<span style="color: blue;">new</span> SecondMessageHandler());
}</pre>
</div>
</div>
<p>When we send a request to /api/cars, the result should be similar to below one:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image_thumb.png" width="644" height="254" /></a></p>
<p>Worked great! But, there is very nasty bug inside our message handler? Can you spot it? OK, I&rsquo;ll wait a minute for you...</p>
<p>Did u spot it? No! OK, no harm there. Let&rsquo;s try to understand what it was. I added another message handler and that message handler runs after the XMagicMessageHandler but on the way back, it runs first. The continuation delegate for the message handler will throw DivideByZeroException exception on purpose.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> SecondMessageHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request,
        CancellationToken cancellationToken) {

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, cancellationToken).ContinueWith(task =&gt; {

            <span style="color: blue;">throw</span> <span style="color: blue;">new</span> DivideByZeroException();

            <span style="color: blue;">return</span> <span style="color: blue;">new</span> HttpResponseMessage();
        });
    }
}</pre>
</div>
</div>
<p>Now let&rsquo;s put a breakpoint inside the continuation method of our XMagicMessageHandler and see what we got.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image_thumb_3.png" width="644" height="187" /></a></p>
<p>As you might expect, the Result property of the task doesn&rsquo;t contain the response message generated by the controller. When we try to reach out to the Result property, the exception thrown by the SecondMessageHandler will be re-thrown here again and more importantly, IMO, your code doesn&rsquo;t do what it is supposed to do. So, how do we get around this? Surely, you can put a try/catch block around the task.Result but that'd be a lame solution. The answer depends on what version of .NET Framework are you running on.</p>
<p>If you are running on .NET v4.0, things are a bit harder for you as you need to deal with TPL face to face. Thanks to ASP.NET team, it is now easier with <a href="http://nuget.org/packages/TaskHelpers.Sources">TaskHelpers</a> NuGet package. The TaskHelpers is actively being used by some major ASP.NET Frameworks internally at Microsoft such as ASP.NET Web API which embraces TPL from top to bottom all the way through the pipeline.</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package TaskHelpers.Sources </code></p>
</div>
<p>If you would like to learn more about TaskHelpers class and how it helps you, <a href="http://twitter.com/bradwilson">@bradwilson</a> has a nice blog posts series on this topic and I listed the links to those posts under resources section at the bottom of this post. After installing the TaskHelpers package from NuGet, we need to fix the bugs inside our message handlers.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> SecondMessageHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request,
        CancellationToken cancellationToken) {

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, cancellationToken).Then(response =&gt; {

            <span style="color: blue;">int</span> left = 10, right = 0;
            <span style="color: blue;">var</span> result = left / right;

            <span style="color: blue;">return</span> response;

        }).Catch&lt;HttpResponseMessage&gt;(info =&gt; {

            <span style="color: blue;">var</span> cacthResult = 
                <span style="color: blue;">new</span> CatchInfoBase&lt;Task&lt;HttpResponseMessage&gt;&gt;.CatchResult();

            cacthResult.Task = TaskHelpers.FromResult(
                request.CreateErrorResponse(
                    HttpStatusCode.InternalServerError, info.Exception));

            <span style="color: blue;">return</span> cacthResult;
        });
    }
}</pre>
</div>
</div>
<p>So, what did we actually change? We now use Then method from the TaskHelpers package instead of directly using ContinueWith method. What Then method does is that it only runs the continuation if the Task&rsquo;s status is RanToCompletion. So, if the base.SendAsync method here returns a faulted or cancelled task, the continuation delegate that we pass into the Then method won&rsquo;t be invoked.</p>
<p>Secondly, we chain another method called Catch which only runs the continuation delegate if the task&rsquo;s status is faulted. If the status is cancelled, the Catch method will return back a Task whose status is set to cancelled. Inside the continuation delegate, we construct a new HttpResponseMessage through <a href="http://msdn.microsoft.com/en-us/library/jj130610(v=vs.108).aspx">CreateErrorResponse</a> extension method for HttpRequestMessage by passing the response status and the exception. The exception details are only sent over the wire if the following conditions are met:</p>
<ul>
<li>If you set GlobalConfiguration.Configuration.IncludeErrorDetailPolicy to IncludeErrorDetailPolicy.Always. </li>
<li>if you set GlobalConfiguration.Configuration.IncludeErrorDetailPolicy to IncludeErrorDetailPolicy.LocalOnly and you run your application locally. </li>
<li>If you set GlobalConfiguration.Configuration.IncludeErrorDetailPolicy to IncludeErrorDetailPolicy.Default and your host environment&rsquo;s error policy allows you to expose error details (for ASP.NET hosting, customErrors element in the Web.config file).</li>
</ul>
<p>The XMagicMessageHandler has been changed nearly the same way as SecondMessageHandler.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> XMagicMessageHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request, 
        CancellationToken cancellationToken) {
            
        <span style="color: green;">//Play with the request here</span>

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, cancellationToken)
            .Then(response =&gt; {

                <span style="color: green;">//Add the X-Magic header</span>
                response.Headers.Add(<span style="color: #a31515;">"X-Magic"</span>, <span style="color: #a31515;">"ThisIsMagic"</span>);
                <span style="color: blue;">return</span> response;

            }).Catch&lt;HttpResponseMessage&gt;(info =&gt; {

                <span style="color: blue;">var</span> cacthResult = 
                    <span style="color: blue;">new</span> CatchInfoBase&lt;Task&lt;HttpResponseMessage&gt;&gt;.CatchResult();

                cacthResult.Task = TaskHelpers.FromResult(
                    request.CreateErrorResponse(
                        HttpStatusCode.InternalServerError, info.Exception));

                <span style="color: blue;">return</span> cacthResult;
            });
    }
}</pre>
</div>
</div>
<p>As I am running the application locally, I should see the error details if I send a request.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image7d388eeb-4c90-4a5d-a60b-ff3db99b4859.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image_thumb_4.png" width="644" height="281" /></a></p>
<p>If I were to send a request with application/xml Accept header, I would get the error response back in xml format.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image_4.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image_thumb_5.png" width="644" height="281" /></a></p>
<p>If you are on .NET v4.5, you will get lots of the things, which we&rsquo;ve just done, out of the box. The Then and Catch extensions method from TaskHelpers package sort of mimic&nbsp;<span style="text-decoration: line-through;">the new async/await compiler features</span>&nbsp;some parts of the new async/await compiler features. So, the .NET v4.5 version of our message handlers look as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> SecondMessageHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> async Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request,
        CancellationToken cancellationToken) {

        <span style="color: blue;">try</span> {

            <span style="color: blue;">int</span> left = 10, right = 0;
            <span style="color: blue;">var</span> result = left / right;

            <span style="color: blue;">var</span> response = await <span style="color: blue;">base</span>.SendAsync(request, cancellationToken);
            <span style="color: blue;">return</span> response;
        }
        <span style="color: blue;">catch</span> (Exception ex) {

            <span style="color: blue;">return</span> request.CreateErrorResponse(
                HttpStatusCode.InternalServerError, ex);
        }
    }
}</pre>
</div>
</div>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> XMagicMessageHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> async Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request, 
        CancellationToken cancellationToken) {

        <span style="color: blue;">try</span> {

            <span style="color: blue;">var</span> response = await <span style="color: blue;">base</span>.SendAsync(request, cancellationToken);
            response.Headers.Add(<span style="color: #a31515;">"X-Magic"</span>, <span style="color: #a31515;">"ThisIsMagic"</span>);
            <span style="color: blue;">return</span> response;
        }
        <span style="color: blue;">catch</span> (Exception ex) {

            <span style="color: blue;">return</span> request.CreateErrorResponse(
                HttpStatusCode.InternalServerError, ex);
        }
    }
}</pre>
</div>
</div>
<p>I would like to clear out that the code that we wrote for .NET v4.0 would work just fine under .NET v4.5 but as you can see we get rid of all the noise here and have a nice and more readable code. The functionality also stays the same.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/HTTP-Message-Handlers-and-AS.NET-Web-API_98DA/image_thumb_6.png" width="644" height="281" /></a></p>
<p>As expected stack trace is different but the result is the same. No matter which way you choose, just choose one. Don&rsquo;t leave things to chance with message handlers!</p>
<h3>Resources</h3>
<ul>
<li><a href="http://www.asp.net/web-api/overview/working-with-http/http-message-handlers">HTTP Message Handlers</a> </li>
<li><a href="http://stackoverflow.com/questions/10251062/asp-net-web-api-message-handlers">ASP.NET Web API Message Handlers</a> </li>
<li><a href="http://channel9.msdn.com/Events/aspConf/aspConf/Async-in-ASP-NET">Async in ASP.NET</a> </li>
<li><a href="http://bradwilson.typepad.com/blog/2012/04/tpl-and-servers-pt1.html">Task Parallel Library and Servers, Part 1: Introduction</a> </li>
<li><a href="http://bradwilson.typepad.com/blog/2012/04/tpl-and-servers-pt2.html">Task Parallel Library and Servers, Part 2: SynchronizationContext</a> </li>
<li><a href="http://bradwilson.typepad.com/blog/2012/04/tpl-and-servers-pt3.html">Task Parallel Library and Servers, Part 3: ContinueWith</a> </li>
<li><a href="http://bradwilson.typepad.com/blog/2012/04/tpl-and-servers-pt4.html">Task Parallel Library and Servers, Part 4: TaskHelpers</a> </li>
<li><a href="http://www.asp.net/mvc/tutorials/mvc-4/using-asynchronous-methods-in-aspnet-mvc-4">Using Asynchronous Methods in ASP.NET MVC 4</a></li>
</ul>