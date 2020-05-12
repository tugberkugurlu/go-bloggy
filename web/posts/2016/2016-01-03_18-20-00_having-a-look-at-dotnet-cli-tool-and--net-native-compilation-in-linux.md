---
id: 4949a72d-bf91-483e-9f2a-b4c4480da7f7
title: Having a Look at dotnet CLI Tool and .NET Native Compilation in Linux
abstract: dotnet CLI tool can be used for building .NET Core apps and for building
  libraries through your development flow (compiling, NuGet package management, running,
  testing, etc.) on various operating systems. Today, I will be looking at this tool
  in Linux, specifically its native compilation feature.
created_at: 2016-01-03 18:20:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET 5
- Linux
slugs:
- having-a-look-at-dotnet-cli-tool-and--net-native-compilation-in-linux
---

<p>I have been following ASP.NET 5 development from the very start and it has been an amazing experience so far. This new platform has seen so many changes both on libraries and concepts throughout but the biggest of all is about to come. The new command line tools that ASP.NET 5 brought to us like dnx and dnu will vanish soon. However, this doesn’t mean that we won’t have a command line first experience. Concepts of these tools will be carried over by a new command line tool: <a href="https://github.com/dotnet/cli">dotnet CLI</a>.</p> <blockquote> <p>Note that dotnet CLI is not even a beta yet. It’s very natural that some of the stuff that I show below may change or even be removed. So, be cautious.</p></blockquote> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/34f855f2-d604-423a-816b-7798b682bfa9.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a2fc83bd-f967-4f43-a4bc-63bdf3657ef2.png" width="644" height="414"></a></p> <p><a href="https://twitter.com/shanselman">Scott Hanselman</a> gave <a href="http://www.hanselman.com/blog/ExploringTheNewNETDotnetCommandLineInterfaceCLI.aspx">an amazing introduction to this tool in his blog post</a>. As indicated in that post, new dotnet CLI tool will give a very similar experience to us compared to other platforms like <a href="https://golang.org/">Go</a>, <a href="https://www.ruby-lang.org/en/">Ruby</a>, <a href="https://www.python.org/">Python</a>. This is very important because, again, this will remove another entry barrier for the newcomers.</p> <p>You can think of this new CLI tool as combination of following three in terms of concepts:</p> <ul> <li>csc.exe</li> <li>msbuild.exe</li> <li>nuget.exe</li></ul> <p>Of course, this is an understatement but it will help you get a grasp of what that tools can do. One other important aspect of the tool is to be able to bootstrap your code and execute it. Here is how:</p> <blockquote> <p>In order to install dotnet CLI tool into my Ubuntu machine, I just followed the steps on <a href="http://dotnet.github.io/getting-started/">the Getting Started guide</a> for Ubuntu.</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/662669c5-a62b-4aed-bdb1-ba43b2068bff.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ee166171-7d6d-46ff-9339-f6460c8dfd13.png" width="644" height="288"></a></p></blockquote> <p>Step one is to create a project structure. My project has two files under "hello-dotnet" folder. <strong>Program.cs</strong>:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">using</span> System;

<span style="color: blue">namespace</span> ConsoleApplication
{
    <span style="color: blue">public</span> <span style="color: blue">class</span> Program
    {
        <span style="color: blue">public</span> <span style="color: blue">static</span> <span style="color: blue">void</span> Main(<span style="color: blue">string</span>[] args)
        {
            Console.WriteLine(<span style="color: #a31515">"Hello World!"</span>);
        }
    }
}</pre></div></div>
<p><strong>project.json</strong>:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>{
    <span style="color: #a31515">"version"</span>: <span style="color: #a31515">"1.0.0-*"</span>,
    <span style="color: #a31515">"compilationOptions"</span>: {
        <span style="color: #a31515">"emitEntryPoint"</span>: <span style="color: blue">true</span>
    },

    <span style="color: #a31515">"dependencies"</span>: {
        <span style="color: #a31515">"Microsoft.NETCore.Runtime"</span>: <span style="color: #a31515">"1.0.1-beta-*"</span>,
        <span style="color: #a31515">"System.IO"</span>: <span style="color: #a31515">"4.0.11-beta-*"</span>,
        <span style="color: #a31515">"System.Console"</span>: <span style="color: #a31515">"4.0.0-beta-*"</span>,
        <span style="color: #a31515">"System.Runtime"</span>: <span style="color: #a31515">"4.0.21-beta-*"</span>
    },

    <span style="color: #a31515">"frameworks"</span>: {
        <span style="color: #a31515">"dnxcore50"</span>: { }
    }
}</pre></div></div>
<p>These are the bare essentials that I need to get something outputted to my console window. One important piece here is the <strong>emitEntryPoint</strong> bit inside the project.json file which indicates that the module will have an entry point which is the static Main method by default.</p>
<p>The second step here is to restore the defined dependencies. This can be done through the <em>"dotnet restore"</em> command:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/8430ab41-3fe1-46a2-8f8b-6722d1f733e5.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/0d269180-0550-477e-b7e2-09339a2c72bd.png" width="644" height="414"></a></p>
<p>Finally, we can now execute the code that we have written and see that we can actually output some text to console. At the same path, just run <em>"dotnet run"</em> command to do that:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d535e9ca-676a-49db-a279-e7b2a3b674b1.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b0ec87f5-076a-4a89-970f-d629ece3d8cb.png" width="644" height="190"></a></p>
<p>Very straight forward experience! Let’s just try to compile the code through "dotnet compile" command:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b3b377ce-8280-4629-a9ea-ed26a06401b0.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b7419bed-803e-4c93-8880-7d8926faef78.png" width="644" height="366"></a></p>
<p>Notice the <em><strong>"hello-dotnet"</strong></em> file there. You can think of this file as dnx which can just run your app. It’s basically the bootstrapper just for your application.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ad4004d6-d65b-49ce-bc16-78fbaaee94f7.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e47105c9-8fe8-4146-9a29-ffff42b9732b.png" width="644" height="206"></a></p>
<p>So, we understand that we can just run this thing:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c9d75640-ebb3-4b68-9fbd-87a9c45773d0.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/49a6ce0b-33ff-4dcc-b7ed-d41b47dbb177.png" width="644" height="206"></a></p>
<p>Very nice! However, that’s not all! This is still a .NET application which requires a few things to be in place to be executed. What we can also do here is to compile native, standalone executables (just like you can do with Go).</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/0322eb73-9a64-4b48-aab6-cd9b150c87cf.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/f5611d58-52f5-4c4b-944e-6440013d73e3.png" width="593" height="484"></a></p>
<p>Do you see the <em>"--native"</em> switch? That will allow you to compile the native executable binary which will be specific to the acrhitecture that you are compiling on (in my case, it’s Ubuntu 14.04):</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/cfcae7d1-3467-4de5-a8a3-39b5ee594e4e.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/83431dfb-9175-4751-872c-668ec0a5ede4.png" width="644" height="366"></a></p>
<p>"hello-dotnet" file here can be executed same as the previous one but this time, it’s all machine code and everything is embedded (yes, even the .NET runtime). So, it’s very usual that you will see a significant increase in the size:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/56764e85-26a6-445e-929c-f3bbf3e6001d.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/9d27b3b6-abd3-4b7b-ae95-e7573c171e84.png" width="644" height="222"></a></p>
<p>This is a promising start and amazing to see that we have a unified tool to rule them all (famous last words). The name of the tool is also great, it makes it straight forward to understand based on your experiences with other platforms and seeing this type of command line first architecture adopted outside of ASP.NET is also great and will bring consistency throughout the ecosystem. I will be watching this space as I am sure there will be more to come :)</p>
<h3>Resources</h3>
<ul>
<li><a href="http://dotnet.github.io/getting-started/">Getting Started with .NET</a></li>
<li><a href="http://dotnet.github.io/api/index.html">.NET Core API Reference</a></li>
<li><a href="https://channel9.msdn.com/Events/ASPNET-Events/ASPNET-Fall-Sessions/Introducing-the-dotnet-CLI">Introducing the dotnet CLI</a></li></ul>  