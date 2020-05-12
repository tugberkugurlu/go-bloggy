---
id: 49b13388-5a40-4a75-86e6-a36ab440e210
title: 'Owin.Dependencies: An IoC Container Adapter Into OWIN Pipeline'
abstract: Owin.Dependencies is an IoC container adapter into OWIN pipeline. This post
  will walk you through the Autofac IoC container implementation and ASP.NET Web API
  framework adapter for OWIN dependencies.
created_at: 2013-09-06 07:33:00 +0000 UTC
tags:
- .net
- ASP.NET Web API
- Autofac
- OWIN
slugs:
- owin-dependencies--an-ioc-container-adapter-into-owin-pipeline
---

<blockquote>
<p><strong>Update on 2013-11-05:</strong></p>
<p>The project has been renamed from&nbsp;Owin.Dependencies to&nbsp;DotNetDoodle.Owin.Dependencies.</p>
</blockquote>
<p>Dependency inject is the design pattern that most frameworks in .NET ecosystem have a first class support: <a href="http://www.asp.net/web-api">ASP.NET Web API</a>, <a href="http://www.asp.net/mvc">ASP.NET MVC</a>, <a href="http://asp.net/signalr">ASP.NET SignalR</a>, <a href="http://nancyfx.org/">NancyFx</a> and so on. Also, the most IoC container implementations for these frameworks allow us to have a per-request lifetime scope which is nearly always what we want. This's great till we stay inside the borders of our little framework. However, with <a href="http://owin.org/">OWIN</a>, our pipeline is much more extended. I have tried to give <a href="http://www.tugberkugurlu.com/archive/getting-started-with-owin-and-the-katana-stack">some insight on OWIN and its Microsoft implementation: Katana</a>. Also, <a href="http://byterot.blogspot.com/2013/08/OWIN-Katana-challenges-blues-library-developer-aspnetwebapi-nancyfx.html">Ali has pointed out some really important facts about OWIN on his post</a>. <a href="http://www.strathweb.com/category/owin/">Flip also has really great content on OWIN</a>. Don't miss the good stuff from <a href="https://twitter.com/randompunter">Damian Hickey</a> <a href="http://dhickey.ie/">on his blog</a> and <a href="https://github.com/damianh">his GitHub repository</a>.</p>
<p>Inside this extended pipeline, the request goes through several middlewares and one of these middlewares may get us into a specific framework. At that point, we will be tied to framework's dependency injection layer. This means that all of your instances you have created inside the previous middlewares (assuming you have a different IoC container there) basically are no good inside your framework's pipeline. This issue had been bitten me and I was searching for a way to at least work around this annoying problem. Then I ended up talking to <a href="https://twitter.com/randompunter">@randompunter</a> and the idea cam out:</p>
<blockquote class="twitter-tweet">
<p><a href="https://twitter.com/tourismgeek">@tourismgeek</a> needed it *across* middleware you can put a request level container in the env dict from a custom middleware.</p>
&mdash; Damian Hickey (@randompunter) <a href="https://twitter.com/randompunter/statuses/373809233449742336">August 31, 2013</a></blockquote>
<script src="//platform.twitter.com/widgets.js"></script>
<p>Actually, the idea is so simple:</p>
<ul>
<li>Have a middleware which starts a per-request lifetime scope at the very beginning inside the OWIN pipeline. </li>
<li>Then, dispose the scope at the very end when the pipeline ends. </li>
</ul>
<p>This idea made me create the <a href="https://github.com/DotNetDoodle/DotNetDoodle.Owin.Dependencies">Owin.Dependencies</a> project. Owin.Dependencies is an IoC container adapter into OWIN pipeline. The core assembly just includes a few interfaces and extension methods. That's all. As expected, <a href="http://www.nuget.org/packages/DotNetDoodle.Owin.Dependencies">it has a pre-release NuGet package</a>, too:</p>
<p class="nuget-badge"><code>PM&gt; Install-Package DotNetDoodle.Owin.Dependencies -Pre </code></p>
<p>Besides the Owin.Dependencies package, I have pushed two different packages and I expect more to come (from and from you, the bright .NET developer). Basically, we have two things to worry about here and these two parts are decoupled from each other:</p>
<ul>
<li>IoC container implementation for Owin.Dependencies </li>
<li>The framework adapter for Owin.Dependencies</li>
</ul>
<p>The IoC container implementation is very straightforward. You just implement two interfaces for that: IOwinDependencyResolver and IOwinDependencyScope. I stole the dependency resolver pattern that ASP.NET Web API has been using but this may change over time as I get feedback. <a href="https://github.com/DotNetDoodle/DotNetDoodle.Owin.Dependencies/tree/master/src/DotNetDoodle.Owin.Dependencies.Autofac">My first OWIN IoC container implementation</a> is for <a href="https://code.google.com/p/autofac/">Autofac</a> and <a href="http://www.nuget.org/packages/DotNetDoodle.Owin.Dependencies.Autofac">it has a seperate NuGet package</a>, too:</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package DotNetDoodle.Owin.Dependencies.Autofac -Pre </code></p>
</div>
<p>The usage of the OWIN IoC container implementation is very simple. You just need to have an instance of the IOwinDependencyResolver implementation and you need to pass that along the UseDependencyResolver extension method on the IAppBuilder interface. Make sure to call the UseDependencyResolver before registering any other middlewares if you want to have the per-request dependency scope available on all middlewares. Here is a very simple Startup class after installing the Owin.Dependencies.Autofac package:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Startup
{
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Configuration(IAppBuilder app)
    {
        IContainer container = RegisterServices();
        AutofacOwinDependencyResolver resolver = <span style="color: blue;">new</span> AutofacOwinDependencyResolver(container);

        app.UseDependencyResolver(resolver)
            .Use&lt;RandomTextMiddleware&gt;();
    }

    <span style="color: blue;">public</span> IContainer RegisterServices()
    {
        ContainerBuilder builder = <span style="color: blue;">new</span> ContainerBuilder();

        builder.RegisterType&lt;Repository&gt;()
                .As&lt;IRepository&gt;()
                .InstancePerLifetimeScope();

        <span style="color: blue;">return</span> builder.Build();
    }
}</pre>
</div>
</div>
<p>After registering the dependency resolver into the OWIN pipeline, we registered our custom middleware called "RandomTextMiddleware". This middleware has been built using a handy abstract class called "OwinMiddleware" from Microsoft.Owin package. The Invoke method of the OwinMiddleware class will be invoked on each request and we can decide there whether to handle the request, pass the request to the next middleware or do the both. The Invoke method gets an IOwinContext instance and we can get to the per-request dependency scope through the IOwinContext instance. Here is the code:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> RandomTextMiddleware : OwinMiddleware
{
    <span style="color: blue;">public</span> RandomTextMiddleware(OwinMiddleware next)
        : <span style="color: blue;">base</span>(next)
    {
    }

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> async Task Invoke(IOwinContext context)
    {
        IOwinDependencyScope dependencyScope = 
            context.GetRequestDependencyScope();
            
        IRepository repository = 
            dependencyScope.GetService(<span style="color: blue;">typeof</span>(IRepository)) 
                <span style="color: blue;">as</span> IRepository;

        <span style="color: blue;">if</span> (context.Request.Path == <span style="color: #a31515;">"/random"</span>)
        {
            await context.Response.WriteAsync(
                repository.GetRandomText()
            );
        }
        <span style="color: blue;">else</span>
        {
            context.Response.Headers.Add(
                <span style="color: #a31515;">"X-Random-Sentence"</span>, 
                <span style="color: blue;">new</span>[] { repository.GetRandomText() });
                
            await Next.Invoke(context);
        }
    }
}</pre>
</div>
</div>
<p>I accept that we kind of have an anti-pattern here: <a href="http://blog.ploeh.dk/2010/02/03/ServiceLocatorisanAnti-Pattern/">Service Locator</a>. It would be really nice to get the object instances injected as the method parameters but couldn't manage to figure out how to do that for now. However, this will be my next try. Here, we get an implementation of the IRepository and it is disposable. When we invoke this little pipeline now, we will see that the disposition will be handled by the infrastructure provided by the Owin.Dependencies implementation.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/b3e73c2776bd_F83E/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/b3e73c2776bd_F83E/image_thumb.png" width="644" height="179" /></a></p>
<p>So far so good but what happens when we need to integrate with a specific framework which has its own dependency injection implementation such as ASP.NET Web API? This is where the framework specific adapters come. I provided the APS.NET Web API adapter and <a href="http://www.nuget.org/packages/DotNetDoodle.Owin.Dependencies.Adapters.WebApi">it has its own NuGet package</a> which depends on the <a href="http://www.nuget.org/packages/Microsoft.AspNet.WebApi.Owin/">Microsoft.AspNet.WebApi.Owin</a> package.</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package DotNetDoodle.Owin.Dependencies.Adapters.WebApi -Pre </code></p>
</div>
<p>Beside this package, I also installed the <a href="https://www.nuget.org/packages/Autofac.WebApi5">Autofac.WebApi5</a> pre-release package to register the API controllers inside the Autofac ContainerBuilder. Here is the modified Startup class to integrate with ASP.NET Web API:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Startup
{
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Configuration(IAppBuilder app)
    {
        IContainer container = RegisterServices();
        AutofacOwinDependencyResolver resolver = 
            <span style="color: blue;">new</span> AutofacOwinDependencyResolver(container);

        HttpConfiguration config = <span style="color: blue;">new</span> HttpConfiguration();
        config.Routes.MapHttpRoute(<span style="color: #a31515;">"DefaultHttpRoute"</span>, <span style="color: #a31515;">"api/{controller}"</span>);

        app.UseDependencyResolver(resolver)
           .Use&lt;RandomTextMiddleware&gt;()
           .UseWebApiWithOwinDependencyResolver(resolver, config);
    }

    <span style="color: blue;">public</span> IContainer RegisterServices()
    {
        ContainerBuilder builder = <span style="color: blue;">new</span> ContainerBuilder();

        builder.RegisterApiControllers(Assembly.GetExecutingAssembly());
        builder.RegisterType&lt;Repository&gt;()
               .As&lt;IRepository&gt;()
               .InstancePerLifetimeScope();

        <span style="color: blue;">return</span> builder.Build();
    }
}</pre>
</div>
</div>
<p>I also added an ASP.NET Web API controller class to serve all the texts I have. As ASP.NET Web API has a fist class DI support and our Web API adapter handles all the things for us, we can inject our dependencies through the controller constructor.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> TextsController : ApiController
{
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IRepository _repo;

    <span style="color: blue;">public</span> TextsController(IRepository repo)
    {
        _repo = repo;
    }

    <span style="color: blue;">public</span> IEnumerable&lt;<span style="color: blue;">string</span>&gt; Get()
    {
        <span style="color: blue;">return</span> _repo.GetTexts();
    }
}</pre>
</div>
</div>
<p>Now, when we send a request to /api/texts, the IRepository implementation is called twice: once from our custom middleware and once from the Web API's TextsController. At the end of the request, the instance is disposed.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/b3e73c2776bd_F83E/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/b3e73c2776bd_F83E/image_thumb_3.png" width="644" height="241" /></a></p>
<p>The source code of this sample is available in the <a href="https://github.com/DotNetDoodle/DotNetDoodle.Owin.Dependencies">Owin.Dependencies repository on GitHub</a>. I agree that there might be some hiccups with this implementation and I expect to have more stable solution in near future. Have fun with this little gem and give your feedback :)</p>