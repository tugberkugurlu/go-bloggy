---
id: 96897472-5818-411d-af54-c0d97f9e7f73
title: 'Challenge of Solving an ASP.NET Web API Self-Hosting Problem: No HTTP resource
  was found that matches the request URI'
abstract: 'Couple of weeks ago, one of my former coworkers ran across a very weird
  problem when he was prototyping on some of his ideas with ASP.NET Web API: No HTTP
  resource was found that matches the request URI. Let''s see what this issue was
  all about and what is the solution.'
created_at: 2014-04-16 08:28:00 +0000 UTC
tags:
- ASP.NET Web API
- Katana
- OWIN
slugs:
- challenge-of-solving-an-asp-net-web-api-self-hosting-problem-no-http-resource-was-found-that-matches-the-request-uri
---

<p>Couple of weeks ago, one of my former coworkers ran across a very weird problem when he was prototyping on some of his ideas with ASP.NET Web API. He was hosting his ASP.NET Web API application on a console application using the <a href="https://www.nuget.org/packages/Microsoft.Owin.Hosting">Microsoft.Owin.Hosting</a> components and <a href="https://www.nuget.org/packages/Microsoft.Owin.Host.HttpListener">Microsoft.Owin.Host.HttpListener</a> host. His solution structure was also very simple. He put all of his controllers, message handlers, filters, etc. in one class library and all the hosting logic inside the console application. The below structure was pretty similar to what he did:  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/9eb56e3f-e625-4b71-ba4f-73589fc8eab0.png"><img title="Screenshot 2014-04-16 10.36.01" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-04-16 10.36.01" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/846550e2-22a7-444a-8f7a-70fe4935dfab.png" width="319" height="484"></a>  <p>Console application also has very little amount of code:  <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">class</span> Program
{
    <span style="color: blue">static</span> <span style="color: blue">void</span> Main(<span style="color: blue">string</span>[] args)
    {
        <span style="color: blue">using</span> (WebApp.Start(<span style="color: #a31515">"http://localhost:5555/"</span>, Start))
        {
            Console.WriteLine(<span style="color: #a31515">"Started listening on localhost:5555"</span>);
            Console.ReadLine();
            Console.WriteLine(<span style="color: #a31515">"Shutting down..."</span>);
        }
    }

    <span style="color: blue">static</span> <span style="color: blue">void</span> Start(IAppBuilder app)
    {
        HttpConfiguration config = <span style="color: blue">new</span> HttpConfiguration();
        config.Routes.MapHttpRoute(<span style="color: #a31515">"DefaultHttpRoute"</span>, <span style="color: #a31515">"api/{controller}"</span>);
        app.UseWebApi(config);
    }
}</pre></div></div>
<p>As you can see, <a href="http://www.asp.net/web-api/overview/hosting-aspnet-web-api/use-owin-to-self-host-web-api">it's all done by the book</a>. However, he was constantly getting 404 when he fired up the application and sent a request to /api/cars:
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ed1d73e6-7d25-4a99-afaf-e51f1fe3fb19.png"><img title="Screenshot 2014-04-16 10.41.50" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Screenshot 2014-04-16 10.41.50" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/52aaa751-e615-4d3f-8c44-7e06335e77a8.png" width="644" height="412"></a>
<p>"No HTTP resource was found that matches the request URI". It's pretty strange. After I looked into the issue for a while, I was able to figure what the problem is:
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/4717a1a9-10e7-4aad-b3ec-b28748db1b84.jpg"><img title="Screenshot 2014-04-16 10.45.50" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Screenshot 2014-04-16 10.45.50" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b43f5764-30bd-4727-9bfc-e34829c636c8.jpg" width="244" height="244"></a>
<p>Let's make this a little bit interesting and have a look at the modules loaded into the AppDomain :)
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/58d144ee-51cc-4f62-a8f2-95c8e47a5f84.png"><img title="Screenshot 2014-04-16 10.49.30" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Screenshot 2014-04-16 10.49.30" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c1cd0338-c708-4854-ac95-767f002d19d7.png" width="644" height="255"></a>
<p>Notice that the WebApiStrangeConsoleHostSample.dll was never loaded into the AppDomain because we never used it even if it's referenced. As ASP.NET Web API uses reflection to determine the controller and the action, it never finds the CarsController. To prove our point here, I'll load the assembly manually:
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">static</span> <span style="color: blue">void</span> Main(<span style="color: blue">string</span>[] args)
{
    Assembly.LoadFrom(Path.Combine(Environment.CurrentDirectory, <span style="color: #a31515">"WebApiStrangeConsoleHostSample.dll"</span>));
    <span style="color: blue">using</span> (WebApp.Start(<span style="color: #a31515">"http://localhost:5555/"</span>, Start))
    {
        Console.WriteLine(<span style="color: #a31515">"Started listening on localhost:5555"</span>);
        Console.ReadLine();
        Console.WriteLine(<span style="color: #a31515">"Shutting down..."</span>);
    }
}</pre></div></div>
<p>The result is a success:
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/7de51c74-94ee-490b-9a9f-7984ff15983f.png"><img title="Screenshot 2014-04-16 10.56.56" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Screenshot 2014-04-16 10.56.56" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/8a1ea5a9-4cfc-49f2-a4a1-c918609d6210.png" width="644" height="318"></a>
<p>However, this is not an ideal solution and I bet that you never run into this issue before. Why? Because, you wise developer keep your hosting agnostic bootstrap code inside the same assembly with you core ASP.NET Web API layer and you call this inside the host application. As soon as you call a method from the core layer assembly, that assembly will be loaded into your AppDomain.
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">static</span> <span style="color: blue">class</span> WebApiConfig
{
    <span style="color: blue">public</span> <span style="color: blue">static</span> <span style="color: blue">void</span> Configure(HttpConfiguration config)
    {
        config.Routes.MapHttpRoute(<span style="color: #a31515">"DefaultHttpRoute"</span>, <span style="color: #a31515">"api/{controller}"</span>);
    }
}</pre></div></div>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">class</span> Program
{
    <span style="color: blue">static</span> <span style="color: blue">void</span> Main(<span style="color: blue">string</span>[] args)
    {
        <span style="color: blue">using</span> (WebApp.Start(<span style="color: #a31515">"http://localhost:5555/"</span>, Start))
        {
            Console.WriteLine(<span style="color: #a31515">"Started listening on localhost:5555"</span>);
            Console.ReadLine();
            Console.WriteLine(<span style="color: #a31515">"Shutting down..."</span>);
        }
    }

    <span style="color: blue">static</span> <span style="color: blue">void</span> Start(IAppBuilder app)
    {
        HttpConfiguration config = <span style="color: blue">new</span> HttpConfiguration();
        WebApiConfig.Configure(config);
        app.UseWebApi(config);
    }
}</pre></div></div>
<p>There is one other potential solution to a problem which is similar to this one. That is to <a href="http://www.strathweb.com/2012/06/using-controllers-from-an-external-assembly-in-asp-net-web-api/">replace the IAssembliesResolver service</a> as Filip did in this post.
  