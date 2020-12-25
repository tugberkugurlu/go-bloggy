---
id: 6e2aa50f-2360-4e5d-a67e-b55aac946969
title: Getting Started With OWIN and the Katana Stack
abstract: OWIN and Katana is best way to build web server indipendent web applications
  in .NET ecosystem and this post will help you a bit on getting started.
created_at: 2013-05-04 18:04:00 +0000 UTC
tags:
- .net
- Hosting
- Katana
- OWIN
slugs:
- getting-started-with-owin-and-the-katana-stack
---

<p>As usual, I tweeted about my excitement which finally led me to this blog post:</p>
<blockquote class="twitter-tweet">
<p>It was amazing to be able to pull assemblies through NuGet w/o being tied to GAC. Now it's super awesome to do the same for web servers <a href="https://twitter.com/search/%23win">#win</a></p>
&mdash; Tugberk Ugurlu (@tourismgeek) <a href="https://twitter.com/tourismgeek/status/330721178765385728">May 4, 2013</a></blockquote>
<script src="//platform.twitter.com/widgets.js"></script>
<p><a href="http://owin.org/">OWIN</a> and <a href="https://katanaproject.codeplex.com/">Katana project</a> is the effort moving towards that direction and even through it's at its early stages, the adoption is promising. Katana is a flexible set of components for building and hosting OWIN-based web applications. This is the definition that I pulled from its Codeplex site. The stack includes several so-called OWIN middlewares, server and host implementations which work with OWIN-based web applications. You can literally get your application going within no time without needing any installations on the machine other than .NET itself. The other benefit is that your application is not tied to one web server; you can choose any of the web server and host implementations at any time without needing to recompile your project's code.</p>
<p>To get started with Katana today, the best way is to jump to your Visual Studio and navigate to Tools &gt; Library Package Manager &gt; Package Manager Settings. From there, navigate to Package Sources and add the following sources:</p>
<ul>
<li>MyGet Katana: <a title="http://myget.org/F/katana/" href="http://myget.org/F/katana/">http://myget.org/F/katana/</a> </li>
<li>MyGet Owin: <a title="http://myget.org/F/owin/" href="http://myget.org/F/owin/">http://myget.org/F/owin/</a></li>
</ul>
<p>Now, we should be able to see those sources through PMC (Package Manager Console):</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb.png" width="644" height="240" /></a></p>
<p>We performed these actions because latest bits of the Katana and <a href="https://github.com/owin/owin-hosting">OWIN Hosting</a> project are pushed to MyGet and those packages are what you want to work with for now. As you can guess, those packages are not stable and not meant to be for production use but good for demonstration cases :) Let's start writing some code and see the beauty.</p>
<p>I started by creating an empty C# Class Library project. Before moving forward, I would like to take a step back and see what packages I have. I selected the MyGet Owin as the current package source and executed the following command: Get-Package -ListAvailable -pre</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb_3.png" width="644" height="331" /></a></p>
<p>These packages are coming from the <a href="https://github.com/owin/owin-hosting">OWIN Hosting project</a> and I encourage you to check the source code out. Let's do the same for the MyGet Katana source:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_4.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb_4.png" width="644" height="465" /></a></p>
<p>We got more packages this time and these packages are coming from the Katana project which is hosted on Codeplex. These packages consist of OWIN middlewares, host and server implementations which we will have a chance to use some of them now.</p>
<p>Let's start installing a host implementation: Microsoft.Owin.Host.HttpListener pre-release package. Now, change the current package source selection to MyGet OWIN and install Owin.Extensions package to be able to get the necessary bits and pieces to complete our demo.</p>
<blockquote>
<p>The Owin.Extensions package will bring down another package named <a href="http://owin.org/">Owin</a> and that Owin package is the only necessary package to have actually. The others are just there to help us but as you can understand, there is no assembly hell involved when working with OWIN. In fact, the Owin.dll only contains one interface which is IAppBuilder. You may wonder how this thing even works then. The answer is simple actually: by convention and discoverability on pure .NET types. To get a more in depth answer on that question, check out <a href="http://vimeo.com/57007898">Louis DeJardin's awesome talk on OWIN</a>.</p>
</blockquote>
<p>What we need to do now is have a class called Startup and that class will have a method called Configuration which takes an IAppBuilder implementation as a parameter.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">partial</span> <span style="color: blue;">class</span> Startup {

    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Configuration(IAppBuilder app) {

        app.UseHandler(async (request, response) =&gt; {

            response.ContentType = <span style="color: #a31515;">"text/html"</span>;
            await response.WriteAsync(<span style="color: #a31515;">"OWIN Hello World!!"</span>);
        });
    }
}</pre>
</div>
</div>
<p>For the demonstration purposes, I used the UseHandler extension method to handle the requests and return the responses. In our case above, all paths will return the same response which is kind of silly but OK for demonstration purposes. To run this application, we need to some sort of a glue which needs to tie our Startup class with the host implementation that we have brought down: Microsoft.Owin.Host.HttpListener. That glue is the OwinHost.exe which we can install from the MyGet Katana NuGet feed.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_5.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb_5.png" width="644" height="161" /></a></p>
<p>OwinHost.exe is going to prepare the settings to host the application and give them to the hosting implementation. Then, it will get out of the way. To make it run, execute the OwinHost.exe without any arguments under the root of your project and you should see screen as below:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image27ea34fb-bc6e-4c13-baba-eb5f1725f337.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb_6.png" width="644" height="193" /></a></p>
<p>We got this unfriendly error message because the OwinHost.exe was unable to locate our assemblies as it looks under the bin directory but our project outputs the compiled assemblies under bin\Debug or bin\Release; depending on the configuration. Change the output directory to bin through the Properties menu, rebuild the solution and run the OwinHost.exe again. This time there should be no error and if we navigate to localhost:5000 (as 5000 is the default port), we should see that the response that we have prepared:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_6.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb_7.png" width="644" height="390" /></a></p>
<p>Cool! You may wonder how it knows to which host implementation to use. The default behavior is the auto-detect but we can explicitly specify the server type as well (however, it's <a href="https://katanaproject.codeplex.com/workitem/29">kind of confusing how to do it today</a>):</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb_8.png" width="644" height="184" /></a></p>
<p>OwinHost.exe is great but as you can guess, it's not our only option. Thanks to Katana project, we can easily host our application on our own process. This option is particularly useful if you would like to deploy your application as a Windows Service or host your application on <a href="http://www.windowsazure.com">Windows Azure</a> Worker Role. To demonstrate this option, I created a console application and referenced our assembly. Then, I installed Microsoft.Owin.Hosting package to be able to host it easily. Here is the code to do that:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">class</span> Program {

    <span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args) {

        <span style="color: blue;">using</span> (IDisposable server = WebApp.Start&lt;Startup&gt;()) {

            Console.WriteLine(<span style="color: #a31515;">"Started..."</span>);
            Console.ReadKey();
        }
    }
}</pre>
</div>
</div>
<p>When I run the console application now, I can navigate to localhost:5000 again and see my response:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_7.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Getting-Started-With-the-Katana-Stack_11A36/image_thumb_9.png" width="644" height="375" /></a></p>
<p>I plan on writing a few more posts on OWIN-based applications but for now, that's all I can give away :) The code I demoed here is available on GitHub as well: <a href="https://github.com/tugberkugurlu/OwinSamples/tree/master/HelloWorldOwin">https://github.com/tugberkugurlu/OwinSamples/tree/master/HelloWorldOwin</a>. but run the Build.cmd file before starting.</p>
<p>Note: Earlier today, I was planning to play with <a href="http://www.asp.net/web-api">ASP.NET Web API</a>'s official OWIN host implementation (which will come out in vNext) and writing a blog post on that. <a href="https://github.com/tugberkugurlu/OwinSamples/tree/master/WebApi50OwinSample">The playing part went well</a> but Flip W. was again a Sergeant Buzzkill :s</p>
<blockquote class="twitter-tweet">
<p>@<a href="https://twitter.com/tourismgeek">tourismgeek</a> writing a post now</p>
&mdash; Filip W (@filip_woj) <a href="https://twitter.com/filip_woj/status/330668689953288192">May 4, 2013</a></blockquote>
<script src="//platform.twitter.com/widgets.js"></script>
<h3>&nbsp;</h3>
<h3>References</h3>
<ul>
<li><a href="https://katanaproject.codeplex.com/">Katana Project</a> </li>
<li><a href="https://github.com/owin/owin-hosting">Official OWIN Hosting Project which includes OWIN hosting components</a> </li>
<li><a href="http://vimeo.com/57007898">OWIN - Run your C# for the web anywhere by Louis DeJardin</a> </li>
<li><a href="http://whereslou.com/2012/05/14/owin-compile-once-and-run-on-any-server">OWIN &ndash; compile once and run on any server</a> (kind of an old but still very useful blog post by Louis DeJardin)</li>
</ul>