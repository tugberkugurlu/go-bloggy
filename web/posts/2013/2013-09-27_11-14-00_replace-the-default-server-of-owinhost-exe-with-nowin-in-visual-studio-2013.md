---
id: 1c2f8dbc-4ec4-43e4-995a-b45810b6ae22
title: Replace the Default Server of OwinHost.exe with Nowin in Visual Studio 2013
abstract: This post will show you how to you can replace the default server of OwinHost.exe
  with Nowin in Visual Studio 2013
created_at: 2013-09-27 11:14:00 +0000 UTC
tags:
- Katana
- OWIN
- Visual Studio
slugs:
- replace-the-default-server-of-owinhost-exe-with-nowin-in-visual-studio-2013
---

<p>Wow, I cannot believe I'm writing this post :) It's really exciting to see that .NET web stack has evolved this much. First, we had a web framework shipping out of band: <a href="http://asp.net/mvc">ASP.NET MVC</a>. Then, we were given a chance to play with hosting agnostic web frameworks such as <a href="http://asp.net/web-api">ASP.NET Web API</a> and <a href="http://asp.net/signalr">ASP.NET SignalR</a> (I know, their names are still confusing and give an impression that those frameworks are still bound to ASP.NET but <a href="http://www.tugberkugurlu.com/archive/newsflash-asp-net-web-api-does-not-sit-on-top-of-asp-net-mvc-in-fact-it-does-not-sit-on-top-of-anything">they are actually not</a>). Recently, we have been embracing the idea of separating the application, server and the host from each other by making our applications <a href="http://www.tugberkugurlu.com/archive/getting-started-with-owin-and-the-katana-stack">OWIN</a> compliant.</p>
<p>By this way, we can easily build our pipeline on our own and switch the underlying host or server easily. Visual Studio 2013 even has a really nice <a href="http://www.tugberkugurlu.com/archive/good-old-f5-experience-with-owinhost-exe-on-visual-studio-2013">extensibility point to switch the underlying host and still have the F5 experience</a>. <a href="http://nuget.org/packages/OwinHost">OwinHost.exe</a> is one of those hosts that we can use as an alternative to IIS. Today, we can even take this further and completely replace the underlying server (which is HttpListener by default with OwinHost.exe) by preserving the host. There is an OWIN compliant server implementation called <a href="http://www.nuget.org/packages/Nowin">Nowin</a> developed by a very clever guy, <a href="https://github.com/Bobris">Boris Letocha</a>. This server uses a raw <a href="http://msdn.microsoft.com/en-us/library/System.Net.Sockets.aspx">.NET socket</a> to listen to the HTTP requests coming in through the wire and will respond to them. By looking at the <a href="https://github.com/Bobris/Nowin#readme">repository's readme file</a>, I can say that this component is not production ready but will work out just fine for demonstration purposes.</p>
<p>I created a mini sample application to show you this cool feature. You can find it on my <a href="https://github.com/tugberkugurlu/OwinSamples/tree/master/NowinSample">GitHub OwinSamples repository</a>, too. It only contains the following startup class and initially have two NuGet packages installed: <a href="http://nuget.org/packages/Owin">Owin</a> and <a href="http://nuget.org/packages/OwinHost">OwinHost</a>.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">using</span> AppFunc = Func&lt;IDictionary&lt;<span style="color: blue;">string</span>, <span style="color: blue;">object</span>&gt;, Task&gt;;

<span style="color: blue;">public</span> <span style="color: blue;">class</span> Startup
{
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Configuration(IAppBuilder app)
    {
        app.Use(<span style="color: blue;">new</span> Func&lt;AppFunc, AppFunc&gt;(ignoreNext =&gt; Invoke));
    }

    <span style="color: blue;">public</span> async Task Invoke(IDictionary&lt;<span style="color: blue;">string</span>, <span style="color: blue;">object</span>&gt; env)
    {
        <span style="color: green;">// retrieve the Request Data from the environment</span>
        <span style="color: blue;">string</span> path = env[<span style="color: #a31515;">"owin.RequestPath"</span>] <span style="color: blue;">as</span> <span style="color: blue;">string</span>;

        <span style="color: blue;">if</span> (path.Equals(<span style="color: #a31515;">"/"</span>, StringComparison.OrdinalIgnoreCase))
        {
            <span style="color: green;">// Prepare the message</span>
            <span style="color: blue;">const</span> <span style="color: blue;">string</span> Message = <span style="color: #a31515;">"Hello World!"</span>;
            <span style="color: blue;">byte</span>[] bytes = Encoding.UTF8.GetBytes(Message);

            <span style="color: green;">// retrieve the Response Data from the environment</span>
            Stream responseBody = env[<span style="color: #a31515;">"owin.ResponseBody"</span>] <span style="color: blue;">as</span> Stream;
            IDictionary&lt;<span style="color: blue;">string</span>, <span style="color: blue;">string</span>[]&gt; responseHeaders = 
                env[<span style="color: #a31515;">"owin.ResponseHeaders"</span>] <span style="color: blue;">as</span> IDictionary&lt;<span style="color: blue;">string</span>, <span style="color: blue;">string</span>[]&gt;;

            <span style="color: green;">// write the headers, response body</span>
            responseHeaders[<span style="color: #a31515;">"Content-Type"</span>] = <span style="color: blue;">new</span>[] { <span style="color: #a31515;">"text/plain"</span> };
            await responseBody.WriteAsync(bytes, 0, bytes.Length);
        }
    }
}</pre>
</div>
</div>
<blockquote>
<p>I am not sure if you have noticed this but Visual Studio 2013 even has an item template for creating an OWIN startup class:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image.png"><img height="394" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>This seems nice but it secretly installs the Microsoft.Owin package, which is unnecessary if you ask me. Owin NuGet package should be enough IMHO.</p>
</blockquote>
<p>I also applied the steps explained in my "<a href="http://www.tugberkugurlu.com/archive/good-old-f5-experience-with-owinhost-exe-on-visual-studio-2013">OwinHost.exe on Visual Studio 2013</a>" post to get my application running on top of OwinHost and here is the result:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_3.png"><img height="346" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Here, my host is OwinHost.exe and it uses the <a href="http://www.nuget.org/packages/Microsoft.Owin.Host.HttpListener/">Microsoft.Owin.Host.HttpListener</a> as the server by default. At the application level, we don't need to know or care about which server we are on but most of the OWIN server implementations expose their names through the server.Capabilities dictionary:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image79cd6c0d-cafe-475d-a142-60963baaa75c.png"><img height="224" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_thumb_4.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>What we want to accomplish here is to keep the host as it is and only replace the server component that it uses underneath. As our host (OwinHost.exe) is OWIN compliant, it can work with any other type of OWIN compliant server implementations and one of them is Nowin. Installing Nowin over NuGet into your project is the first step that you need to do.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_4.png"><img height="238" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_thumb_5.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>However, this's not enough by itself as OwinHost.exe has no idea that we want to use Nowin as our server. Luckily, we can manipulate the arguments we pass onto OwinHost.exe. You can configure these arguments through the Web pane inside the Visual Studio Project Properties window. Besides that, OwinHost.exe accepts a few command line switches and one of them is the &ndash;s (or &ndash;-server) switch to load the specified server factory type or assembly. These are all we needed.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image.png"><img height="200" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_thumb_6.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>After saving the changes we have made, we can run the application and get the same result on top of a completely different server implementation:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_5.png"><img height="345" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_thumb_7.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Also with the same way as we did earlier, we can see that the switch has been made and Nowin is in use:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image0434bf5d-2eae-4c6b-97b6-f9a1983e65fa.png"><img height="232" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Replace-the-Default-Server_9F8C/image_thumb_8.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Being able to do all of this is a very big deal; especially if you think that we have been tied to IIS for very long time. I'm so happy to see .NET web stack moving towards this flexible direction.</p>