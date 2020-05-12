---
id: d7181e8b-2ae7-4839-b4a4-e3a3f2ffb3c9
title: Getting Started with ASP.NET vNext by Setting Up the Environment From Scratch
abstract: In this post, I'll walk you through how you can set up your environment
  from scratch to get going with ASP.NET vNext.
created_at: 2014-09-28 19:08:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET vNext
- HTTP
slugs:
- getting-started-with-asp-net-vnext-by-setting-up-the-environment-from-scratch
---

<p>I'm guessing that you already heard the news about <a href="http://asp.net/vnext">ASP.NET vNext</a>.&nbsp; It has been announced publicly a few months back at <a href="http://channel9.msdn.com/Events/TechEd/NorthAmerica/2014/DEV-B411#fbid=">TechEd North America 2014</a> and it's being rewritten from the ground up which means: "Say goodbye to our dear System.Web.dll" :) No kidding, I'm pretty serious :) It brings lots of improvements which will take the application and development performance to the next level for us, .NET developers. ASP.NET vNext is coming so hard and there are already good amount of resources for you to dig into. If you haven't done this yet, navigate to at the bottom of this post to go through the links under Resources.&nbsp; However, I strongly suggest you to check <a href="https://channel9.msdn.com/Events/TechEd/NewZealand/2014/DEV213">Daniel Roth’s talk on ASP.NET vNext at TechEd New Zealand 2014</a> first which is probably the best introduction talk on ASP.NET vNext. <p>What I would like to go through in this post is how you can set up your environment from scratch to get going with ASP.NET vNext. I also tnhink that this post is useful for understanding key concepts behind this new development environment.  <p>First of all, you need to install <a href="https://github.com/aspnet/kvm">kvm</a> which stands for k Version Manager. You can install kvm by running the below command on your command prompt in Windows.</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>@powershell <span style="color: gray">-</span>NoProfile <span style="color: gray">-</span>ExecutionPolicy unrestricted <span style="color: gray">-</span>Command <span style="color: #a31515">"iex ((new-object net.webclient).DownloadString('https://raw.githubusercontent.com/aspnet/Home/master/kvminstall.ps1'))"</span></pre></div></div>
<p>This will invoke an elevated PowerShell command prompt and installs a few stuff on your machine. Actually, all the things that this command installs are under %userprofile%\.kre directory. 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1f67d5e8-30c4-4ae2-abd4-27d3a0b5e5f9.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/bf777dfe-ab38-4549-bc98-8495d188c614.png" width="451" height="484"></a></p>
<p>Now, install the latest available <a href="https://github.com/aspnet/KRuntime">KRuntime</a> environment. You will do this by running kvm upgrade command:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/07a7953f-ca3b-43b1-9ab6-b3f64c9b315b.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/97694b9c-5f72-4990-9075-f6b529ab94e1.png" width="644" height="249"></a></p>
<p>The latest version is installed from the default feed (which is https://www.myget.org/F/aspnetmaster/api/v2 at the time of this writing). We can verify that the K environment is really installed by running kvm list which will list installed k environments with their associated information.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a39396b7-c89f-491f-b406-fa873285c29f.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a9664ae5-2077-4e90-a546-93176cc78363.png" width="644" height="249"></a></p>
<p>Here, we only have good old desktop CLR. If we want to work against CoreCLR (a.k.a K10), we should install it using the –svrc50 switch:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/8d37f64c-a182-4c24-9e10-caf6c5c221a3.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/4d4a46d1-a7a1-4ebe-a581-751049c95382.png" width="644" height="226"></a></p>
<p>You can switch between versions using the "kvm use" command. You can also use –p switch to persist your choice of the runtime and that would allow you to specify your default runtime choice which will live between processes.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d241e9b8-2212-44cb-a83c-1537dc279798.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/11601122-363b-4545-80a3-ba75bc373dee.png" width="644" height="181"></a></p>
<p>Our system is ready for ASP.NET vNext development. I have a tiny working AngularJS application that you can find here in <a href="https://github.com/tugberkugurlu/angularjs-getting-started/tree/35d0998c5ada624afe24b7e09a7fa7e5a5e89d2a">my GitHub repository</a>. This was a pure HTML + CSS + JavaScript web application which required no backend system. However, at some point I needed some backend functionality. So, I integrated with ASP.NET vNext. Here is how I did it:</p>
<p>First, we need to specify the NuGet feeds that we will be using for our application. Go ahead and place the the following content into NuGet.config file inside the root of your solution:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">&lt;?</span><span style="color: #a31515">xml</span> <span style="color: red">version</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">1.0</span><span style="color: black">"</span> <span style="color: red">encoding</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">utf-8</span><span style="color: black">"</span><span style="color: blue">?&gt;</span>
<span style="color: blue">&lt;</span><span style="color: #a31515">configuration</span><span style="color: blue">&gt;</span>
  <span style="color: blue">&lt;</span><span style="color: #a31515">packageSources</span><span style="color: blue">&gt;</span>
      <span style="color: blue">&lt;</span><span style="color: #a31515">clear</span> <span style="color: blue">/&gt;</span>
      <span style="color: blue">&lt;</span><span style="color: #a31515">add</span> <span style="color: red">key</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">AspNetVNext</span><span style="color: black">"</span> <span style="color: red">value</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">https://www.myget.org/F/aspnetrelease/api/v2</span><span style="color: black">"</span> <span style="color: blue">/&gt;</span>
      <span style="color: blue">&lt;</span><span style="color: #a31515">add</span> <span style="color: red">key</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">NuGet.org</span><span style="color: black">"</span> <span style="color: red">value</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">https://nuget.org/api/v2/</span><span style="color: black">"</span> <span style="color: blue">/&gt;</span>
  <span style="color: blue">&lt;/</span><span style="color: #a31515">packageSources</span><span style="color: blue">&gt;</span>
<span style="color: blue">&lt;/</span><span style="color: #a31515">configuration</span><span style="color: blue">&gt;</span></pre></div></div>
<p>Later, have the project.json inside your application directory. This will have your dependencies and commands. It can contain more but we will just have those for this post:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>{
    <span style="color: #a31515">"dependencies"</span>: {
        <span style="color: #a31515">"Kestrel"</span>: <span style="color: #a31515">"1.0.0-alpha3"</span>,
        <span style="color: #a31515">"Microsoft.AspNet.Diagnostics"</span>: <span style="color: #a31515">"1.0.0-alpha3"</span>,
        <span style="color: #a31515">"Microsoft.AspNet.Hosting"</span>: <span style="color: #a31515">"1.0.0-alpha3"</span>,
        <span style="color: #a31515">"Microsoft.AspNet.Mvc"</span>: <span style="color: #a31515">"6.0.0-alpha3"</span>,
        <span style="color: #a31515">"Microsoft.AspNet.Server.WebListener"</span>: <span style="color: #a31515">"1.0.0-alpha3"</span>
    },
    
    <span style="color: #a31515">"commands"</span>: { 
        <span style="color: #a31515">"web"</span>: <span style="color: #a31515">"Microsoft.AspNet.Hosting --server Microsoft.AspNet.Server.WebListener --server.urls http://localhost:5001"</span>,
        <span style="color: #a31515">"kestrel"</span>: <span style="color: #a31515">"Microsoft.AspNet.Hosting --server Kestrel --server.urls http://localhost:5004"</span>
    },
    <span style="color: #a31515">"frameworks"</span>: {
        <span style="color: #a31515">"net45"</span>: {},
        <span style="color: #a31515">"k10"</span>: {}
    }
}</pre></div></div>
<p>With the runtime instillation, you can have a few things and one of them is kpm tool which allows you to manage your packages. You can think of this as NuGet (indeed, it uses NuGet behind the scenes). However, it knows how to read your project.json file and installs the packages according to that. If you call kpm now, you can see the options it give you: 
<p>As you have the project.json ready, you can now run kpm restore:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e9b8c150-3c0c-4db1-985a-9ae762012918.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/708d2c39-513b-4ffd-92a8-d2b30669ce9e.png" width="585" height="484"></a></p>
<blockquote>
<p>Note that the restore output is a little different if you are using an alpha4 release:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d079fa76-38d0-4b38-8f83-86ffbeffc37b.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5aa5c73f-eb0a-4aec-9646-8b58acad8048.png" width="585" height="484"></a></p></blockquote>
<p>Based on the commands available in your project.json file, you can run the command to fire up your application now. </p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d7eb4658-4a87-4b45-97c4-77338baa38c4.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/74c5f671-dabe-4cdf-84e7-acdff79aaf19.png" width="644" height="147"></a></p>
<p>Also, you can run "set KRE_TRACE=1" before running your command to see diagnostic details about the process if you need to:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/3efafc1c-751a-4a98-af1a-165c57b50767.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/97bfe60c-b3ee-4152-b334-be845d798d23.png" width="491" height="484"></a></p>
<p>My little app is now running:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/67807e39-7db8-47d2-bc6b-8de29f327803.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c0f2a8cc-b255-407d-8fac-b21061ce057b.png" width="644" height="164"></a></p>
<h3>Resources</h3>
<ul>
<li><a href="https://github.com/aspnet/Home/wiki">Introduction to ASP.NET vNext (GitHub)</a> 
<li><a href="http://www.asp.net/vnext/overview/aspnet-vnext/overview">Getting Started with ASP.NET vNext (ASP.NET)</a> 
<li><a href="http://davidfowl.com/asp-net-vnext/">ASP.NET vNext by David Fowler</a> 
<li><a href="http://davidfowl.com/asp-net-vnext-architecture/">ASP.NET vNext Overview by David Fowler</a> 
<li><a href="https://github.com/aspnet/Home/wiki/KRuntime-structure">KRuntime Structure</a> 
<li><a href="https://github.com/aspnet/Home/wiki/Project.json-file">Project.json File</a> 
<li><a href="https://github.com/aspnet/Home/wiki/Command-Line">Command Line Usage</a> 
<li><a href="https://github.com/aspnet/Home/wiki/Assembly-Neutral-Interfaces">Assembly Neutral Interfaces</a> 
<li><a href="http://davidfowl.com/assembly-neutral-interfaces/">Assembly Neutral Interfaces by David Fowler</a> 
<li><a href="http://davidfowl.com/assembly-neutral-interfaces-implementation/">Assembly Neutral Types Implementation by David Fowler</a> 
<li><a href="http://www.asp.net/vnext/overview/aspnet-vnext/walkthrough-mvc-music-store">MVC Music Store Sample Application for ASP.NET vNext</a> 
<li><a href="http://www.asp.net/vnext/overview/aspnet-vnext/walkthrough-bugtracker">BugTracker Sample Application for ASP.NET vNext</a></li></ul>  