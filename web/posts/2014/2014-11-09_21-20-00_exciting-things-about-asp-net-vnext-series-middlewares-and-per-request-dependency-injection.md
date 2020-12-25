---
id: 95510d2d-7608-4f72-8023-895a6be5b70c
title: 'Exciting Things About ASP.NET vNext Series: Middlewares and Per Request Dependency
  Injection'
abstract: From the very first day of ASP.NET vNext, per request dependencies feature
  is a first class citizen inside the pipeline. In this post, I'd like to show you
  how you can use this feature inside your middlewares.
created_at: 2014-11-09 21:20:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET vNext
slugs:
- exciting-things-about-asp-net-vnext-series-middlewares-and-per-request-dependency-injection
---

<p>Web development experience with .NET has never seen a drastic change like this since its birth day. Yes, I’m talking about <a href="https://www.tugberkugurlu.com/archive/getting-started-with-asp-net-vnext-by-setting-up-the-environment-from-scratch">ASP.NET vNext</a> :) I have been putting my toes into this water for a while now and a few weeks ago, I started a new blog post series about ASP.NET vNext. To be more specific, I’m planning on writing about the things I am actually excited about this new cloud optimized (TM) runtime. Those things could be anything which will come from <a href="http://github.com/aspnet">ASP.NET GitHub account</a>: things I like about the development process, Visual Studio tooling experience for ASP.NET vNext, bowels of <a href="https://github.com/aspnet/kruntime">this new runtime</a>, tiny little things about the frameworks like <a href="http://github.com/aspnet/mvc">MVC</a>, <a href="https://github.com/aspnet/identity">Identity</a>, <a href="https://github.com/aspnet/entityframework">Entity Framework</a>.</p> <p>Today,&nbsp; I would like to show you one of my favorite features in ASP.NET vNext: per-request dependencies.</p> <blockquote> <p><strong>BIG ASS CAUTION!</strong> At the time of this writing, I am using <strong>KRE 1.0.0-beta2-10679 </strong>version. As things are moving really fast in this new world, it’s very likely that the things explained here will have been changed as you read this post. So, be aware of this and try to explore the things that are changed to figure out what are the corresponding new things. <p>Also, inside this post I am referencing a lot of things from ASP.NET GitHub repositories. In order to be sure that the links won’t break in the future, I’m actually referring them by getting permanent links to the files on GitHub. So, these links are actually referring the files from the latest commit at the time of this writing and they have a potential to be changed, too. Read the "<a href="https://help.github.com/articles/getting-permanent-links-to-files/">Getting permanent links to files</a>" post to figure what this actually is.</p></blockquote> <p>If you follow my blog, you probably know that I have written about <a href="https://www.tugberkugurlu.com/archive/owin-dependencies--an-ioc-container-adapter-into-owin-pipeline">OWIN and dependency injection</a> before. I also have a little library called <a href="https://github.com/DotNetDoodle/DotNetDoodle.Owin.Dependencies">DotNetDoodle.Owin.Dependencies</a> which acts as an IoC container adapter into OWIN pipeline. One of the basic ideas with that library is to be able to reach out to a service provider and get required services inside your middleware. Agree or disagree, you will sometimes need external dependency inside your middleware :) and some of those dependencies will need to be constructed per each request separately and disposed at the end of that request. This is a very common approach and nearly all .NET IoC containers has a concept of per-request lifetime for <a href="http://asp.net/mvc">ASP.NET MVC</a>, <a href="http://asp.net/web-api">ASP.NET Web API</a> and possibly for other web frameworks like <a href="http://mvc.fubu-project.org/">FubuMVC</a>.&nbsp; </p> <p>From the very first day of ASP.NET vNext, per request dependencies feature is a first class citizen inside the pipeline. There are a few ways to work with per request dependencies and I’ll show you two of the ways that you can take advantage of inside your middleware.</p> <h3>Basics</h3> <p>For you ASP.NET vNext web application, you would register your services using one of the overloads of <a href="https://github.com/aspnet/Hosting/blob/8f16060f941b71551be09015d76efb86770d84d7/src/Microsoft.AspNet.RequestContainer/ContainerExtensions.cs#L32-L39">UseServices</a> extension method on <a href="https://github.com/aspnet/HttpAbstractions/blob/dev/src/Microsoft.AspNet.Http/IApplicationBuilder.cs">IApplicationBuilder</a>. Using this extension method, you are registering your services and adding the <a href="https://github.com/aspnet/Hosting/blob/8f16060f941b71551be09015d76efb86770d84d7/src/Microsoft.AspNet.RequestContainer/ContainerMiddleware.cs">ContainerMiddlware</a> into the middleware pipeline. ContainerMiddlware is responsible for creating you a dependency scope behind the scenes and any dependency you register using the <a href="https://github.com/aspnet/DependencyInjection/blob/b2734278f5a0dc456d04941c0954145eff567c39/src/Microsoft.Framework.DependencyInjection/ServiceCollectionExtensions.cs#L24-L36">AddScoped</a> method (in other words, any dependency which is marked with <a href="https://github.com/aspnet/DependencyInjection/blob/b2734278f5a0dc456d04941c0954145eff567c39/src/Microsoft.Framework.DependencyInjection/LifecycleKind.cs#L14">LifecycleKind.Scoped</a>) on the <a href="https://github.com/aspnet/DependencyInjection/blob/b2734278f5a0dc456d04941c0954145eff567c39/src/Microsoft.Framework.DependencyInjection/IServiceCollection.cs">IServiceCollection</a> implementation will end up being a per-request dependency (of course it depends how you use it). The following code is a sample of how you would register your dependencies <a href="https://github.com/tugberkugurlu/AspNetvNext-ConfigOptionsSample/blob/4d0fe78d43c6b8993fd6906b0a73e8cf7ea46f33/src/Farticus.App/Startup.cs">inside your Startup class</a>:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">void</span> Configure(IApplicationBuilder app)
{
    app.UseServices(services =&gt;
    {
        services.AddScoped&lt;IFarticusRepository, InMemoryFarticusRepository&gt;();
        services.AddTransient&lt;IConfigureOptions&lt;FarticusOptions&gt;, FarticusOptionsSetup&gt;();
    });

    app.UseMiddleware&lt;FarticusMiddleware&gt;();
}</pre></div></div>
<p>The way you reach out to per-request dependencies also vary and we will go through this topic in a few words later but for now, we should know that request dependencies are hanging off of <a href="https://github.com/aspnet/HttpAbstractions/blob/389e27e46055a95648efe48a512614fc4f8ff08e/src/Microsoft.AspNet.Http/HttpContext.cs#L27">HttpContext.RequestServices</a> property.</p>
<h3>Per-request Dependencies Inside a Middleware (Wrong Way)</h3>
<p>Inside the above dependency registration sample code snippet, I also registered FarticusMiddleware which hijacks all the requests coming to the web server and well..., farts :) FarticusMiddleware is actually very ignorant and it doesn’t know what to say when it farts. So, it gets the fart message from an external service: an IFarticusRepository implementation. IFarticusRepository implementation can depend on other services to do some I/O to get a random message and those services might be constructed per request lifecycle. So, FarticusMiddleware needs to be aware of this fact and consume the service in this manner. You can see <a href="https://github.com/tugberkugurlu/AspNetvNext-ConfigOptionsSample/blob/d19ef70f3562cd6b0a106d97c57156038af8ead7/src/Farticus/FarticusMiddleware.cs">the implementation of this middleware</a> below (trimmed-down version):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> FarticusMiddleware
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> FarticusOptions _options;

    <span style="color: blue">public</span> FarticusMiddleware(
        RequestDelegate next, IOptions&lt;FarticusOptions&gt; options)
    {
        _options = options.Options;
    }

    <span style="color: blue">public</span> async Task Invoke(HttpContext context)
    {
        IFarticusRepository repository = context
            .RequestServices
            .GetService(<span style="color: blue">typeof</span>(IFarticusRepository)) <span style="color: blue">as</span> IFarticusRepository;

        <span style="color: blue">if</span>(repository == <span style="color: blue">null</span>)
        {
            <span style="color: blue">throw</span> <span style="color: blue">new</span> InvalidOperationException(
                <span style="color: #a31515">"IFarticusRepository is not available."</span>);
        }

        <span style="color: blue">var</span> builder = <span style="color: blue">new</span> StringBuilder();
        builder.Append(<span style="color: #a31515">"&lt;div&gt;&lt;strong&gt;Farting...&lt;/strong&gt;&lt;/div&gt;"</span>);
        <span style="color: blue">for</span>(<span style="color: blue">int</span> i = 0; i &lt; _options.NumberOfMessages; i++)
        {
            <span style="color: blue">string</span> message = await repository.GetFartMessageAsync();
            builder.AppendFormat(<span style="color: #a31515">"&lt;div&gt;{0}&lt;/div&gt;"</span>, message);
        }

        context.Response.ContentType = <span style="color: #a31515">"text/html"</span>;
        await context.Response.WriteAsync(builder.ToString());
    }
}</pre></div></div>
<p>So, look at what we are doing here: </p>
<ul>
<li>We are getting the IOptions&lt;FarticusOptions&gt; implementation through the middleware constructor which is constructed per-pipeline instance which basically means per application lifetime.</li>
<li>When the Invoke method is called per each request, we reach out to HttpContext.RequestServices and try to retrieve the IFarticusRepository from there.</li>
<li>When we have the IFarticusRepository implementation, we are creating the response message and writing it to the response stream.</li></ul>
<p>The only problem here is that we are being friends with the <a href="http://blog.ploeh.dk/2010/02/03/ServiceLocatorisanAnti-Pattern/">notorious service locator pattern</a> inside our application code which is pretty bad. Invoke method tries to resolved the dependency itself through the RequestServices property. Imagine that you would like to write unit tests for this code. At that time, you need know about the implementation of the Invoke method here to write unit tests against it because there is no other way for you to know that this method relies on IFarticusRepository.</p>
<h3>Per-request Dependencies Inside a Middleware (Correct Way)</h3>
<p>I was talking to a few ASP.NET team members if it was possible to inject per-request dependencies into Invoke method of my middleware. A few thoughts, <a href="https://twitter.com/loudej">Louis DeJardin</a> sent <a href="https://github.com/aspnet/HttpAbstractions/pull/145">a pull request</a> to enable this feature in a very simple and easy way. If you register your middleware using the <a href="https://github.com/aspnet/HttpAbstractions/blob/b7d9e11a8442994cf5c4b8292d72b1efbd7e5a42/src/Microsoft.AspNet.Http.Extensions/UseMiddlewareExtensions.cs#L20">UseMiddleware</a> extension method on IApplicationBuilder, it checks whether you are expecting any other parameters through the Invoke method rather than an HttpContext instance and inject those if you have any. This new behavior allows us to <a href="https://github.com/tugberkugurlu/AspNetvNext-ConfigOptionsSample/blob/4d0fe78d43c6b8993fd6906b0a73e8cf7ea46f33/src/Farticus/FarticusMiddleware.cs#L43-L67">change the above Invoke method with the below one</a>:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> async Task Invoke(HttpContext context, IFarticusRepository repository)
{
    <span style="color: blue">if</span>(repository == <span style="color: blue">null</span>)
    {
        <span style="color: blue">throw</span> <span style="color: blue">new</span> InvalidOperationException(<span style="color: #a31515">"IFarticusRepository is not available."</span>);
    }

    <span style="color: blue">var</span> builder = <span style="color: blue">new</span> StringBuilder();
    builder.Append(<span style="color: #a31515">"&lt;div&gt;&lt;strong&gt;Farting...&lt;/strong&gt;&lt;/div&gt;"</span>);
    <span style="color: blue">for</span>(<span style="color: blue">int</span> i = 0; i &lt; _options.NumberOfMessages; i++)
    {
        <span style="color: blue">string</span> message = await repository.GetFartMessageAsync();
        builder.AppendFormat(<span style="color: #a31515">"&lt;div&gt;{0}&lt;/div&gt;"</span>, message);
    }

    context.Response.ContentType = <span style="color: #a31515">"text/html"</span>;
    await context.Response.WriteAsync(builder.ToString());
}</pre></div></div>
<p>Nice and very clean! It’s now very obvious that the Invoke method relies on IFarticusRepository. You can find <a href="https://github.com/tugberkugurlu/AspNetvNext-ConfigOptionsSample">the sample application</a> I referred here under my GitHub account. When you run the application and send an HTTP request to localhost:5001, you will see the log which shows you when the instances are created and disposed:</p>
<blockquote>
<p>I am still adding new things to this sample and it's possible that the code is changed when you read this post. However, <a href="https://github.com/tugberkugurlu/AspNetvNext-ConfigOptionsSample/tree/4d0fe78d43c6b8993fd6906b0a73e8cf7ea46f33">here</a> is latest state of the sample application at the time of this writing.</p></blockquote>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/14aabbe0-d678-4397-8509-ee488f2086d0.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/58243d52-9c6b-410f-b308-3993f98ad124.png" width="644" height="328"></a></p>  