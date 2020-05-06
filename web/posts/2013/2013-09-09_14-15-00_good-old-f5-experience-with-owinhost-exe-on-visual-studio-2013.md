---
title: Good Old F5 Experience With OwinHost.exe on Visual Studio 2013
abstract: 'With Visual Studio 2013 RC, we are introduced to a new extensiblity point:
  External Host. This gives us the F5 experience Wwth OwinHost.exe on VS 2013 and
  this post walks you through this feature.'
created_at: 2013-09-09 14:15:00 +0000 UTC
tags:
- ASP.Net
- Katana
- OWIN
- Visual Studio
slugs:
- good-old-f5-experience-with-owinhost-exe-on-visual-studio-2013
---

<p>Since the first day I started developing software solutions for web on Visual Studio, I have had the 'F5' experience. It's nice, clean and easy to get up a web application running without digging through a specific exe file to invoke. After Cassini, our development environments have change a little with <a href="http://www.iis.net/learn/extensions/introduction-to-iis-express/iis-express-overview">IIS Express</a>, in a nice way of course. With the introduction of IIS Express, <a href="http://www.iis.net/learn/extensions/using-iis-express/running-iis-express-from-the-command-line">running it through the command line</a> is also possible but we still have the ultimate F5 experience to get it up.</p>
<p>Today, we are still in love with IIS but with the introduction of <a href="http://owin.org/">OWIN</a>, <a href="http://www.tugberkugurlu.com/archive/getting-started-with-owin-and-the-katana-stack">it's now much easier to build your application in an hosting agnostic way</a> and <a href="https://katanaproject.codeplex.com/">Katana Project</a> is providing components for hosting your application outside of IIS with a very thin wrapper on top of <a href="http://msdn.microsoft.com/en-us/library/system.net.httplistener.aspx">HttpListener</a>. Actually, the HttpListener is one of the options in terms of receiving the HTTP request and sending back the response at the operating system level. The <a href="http://www.nuget.org/packages/OwinHost">OwinHost.exe</a> even abstracts that part for us. The OwinHost.exe is a very simple tool which will bring your application startup and the host layer together. Then, it will get out of your way.</p>
<p>However, the development environment with a custom executable tool can be very frustrating, especially in .NET ecosystem where we have really great tooling all over the place. Typically, if you would like to host your application with OwinHost.exe in your development environment, here is what you need to do:</p>
<ol>
<li>Create a class library project and install Microsoft.Owin package from <a href="http://www.myget.org/F/aspnetwebstacknightly/">MyGet ASP.NET Web Stack Nightlies</a> feed. </li>
<li>Add your startup class and build your pipeline with several OWIN middlewares inside the Configuration method. </li>
<li>Add an assembly level attribute for OwinStartupAttribute to indicate your start-up class. </li>
<li>Optionally, install the OwinHost.exe from the same feed into your project or install it any other place. </li>
<li>Go to project properties window and navigate the "Build" pane. Change the "Output Path" from "bin\Debug\" or "bin\Release\" to "bin". </li>
<li>Open up a command prompt window and navigate the root of your project. </li>
<li>Execute the OwinHost.exe from the command prompt window you have just opened up and here you have the running HTTP application.</li>
</ol>
<p>First four steps are the ones that we are OK with but the rest is generally frustrating. What we really want there is to have the F5 experience that we are used to. With Visual Studio 2013 RC, we are introduced to a new extensibility point in Visual Studio: External Host. With this feature, we are now able to replace the host layer completely. Let's see how. First, let's create a new ASP.NET Web Application:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb.png" width="644" height="394" /></a></p>
<p>If you are new to this "One ASP.NET" idea, this window may seem a bit awkward at first but think this as an open buffet of frameworks and modules. After selecting this, we will see the open buffet window:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_3.png" width="644" height="402" /></a></p>
<p>For our sample application, we can continue our merry way with the empty project template without any core references. However, the Visual Studio project will still include System.Web and its related references as shown in the below picture.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_4.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_4.png" width="368" height="484" /></a></p>
<p>Including the web.config file, we can remove all of this System.Web related stuff now to run our pure OWIN based application. In some cases, we may need web.config though as <a href="https://katanaproject.codeplex.com/discussions/455148">OwinHost.exe uses it for some cases such as binding redirects</a>. After removing these junks, we have a project as seen below:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_5.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_5.png" width="244" height="225" /></a></p>
<p>As we mentioned before, Visual Studio 2013 RC has the external host extensibility support for web applications and we can see this extensibility point by opening the Properties window of the project.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image9764a5e1-dc32-461d-b613-0636cdf8b988.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_6.png" width="644" height="421" /></a></p>
<p>You get the IIS Express and Local IIS out of the box with Visual Studio 2013 RC (and yes, Cassini is dead for God's sake). However, it's really easy to get in there and OwinHost NuGet package has the specific installation scripts to register itself as an external host. To get this feature from OwinHost, we just need to add it into our project through NuGet. Today, OwinHost 2.0.0-rc1 package is available on NuGet.org and when we install the package, we get a warning popup from Visual Studio:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_6.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_7.png" width="644" height="233" /></a></p>
<p>When we accept the integration, we can see the OwinHost showing up inside the external servers list.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_8.png" width="644" height="203" /></a></p>
<p>If you select the OwinHost option, we can see the option which we can configure if we need to.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_7.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_9.png" width="644" height="218" /></a></p>
<p>All these information is stored inside the project file (csproj if you are on C#):</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image26a099db-2c78-4467-9b6d-f3ecebb286e6.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_10.png" width="644" height="248" /></a></p>
<p>Now, I can build my OWIN pipeline and have it running on top of OwinHost.exe just by pressing F5 inside the Visual Studio 2013. For demonstration purposes, I wrote up the following Startup class:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">using</span> Microsoft.Owin;
<span style="color: blue;">using</span> Owin;
<span style="color: blue;">using</span> OwinHostVsIntegrationSample;
<span style="color: blue;">using</span> System.IO;

[<span style="color: blue;">assembly</span>: OwinStartup(typeof(Startup))]
<span style="color: blue;">namespace</span> OwinHostVsIntegrationSample
{
    <span style="color: blue;">public</span> <span style="color: blue;">class</span> Startup
    {
        <span style="color: blue;">public</span> <span style="color: blue;">void</span> Configuration(IAppBuilder app) 
        {
            app.Use(async (ctx, next) =&gt; 
            {
                TextWriter output = ctx.Get&lt;TextWriter&gt;(<span style="color: #a31515;">"host.TraceOutput"</span>);
                output.WriteLine(<span style="color: #a31515;">"{0} {1}: {2}"</span>, 
                    ctx.Request.Scheme, ctx.Request.Method, ctx.Request.Path);

                await ctx.Response.WriteAsync(<span style="color: #a31515;">"Hello world!"</span>);
            });
        }
    }
}</pre>
</div>
</div>
<p>When I press CTRL + F5, I have my Web application running without debugging on top of OwinHost.exe:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_8.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_11.png" width="644" height="330" /></a></p>
<p>Very nice and elegant way of running your application with a custom server! Besides this, we can certainly run our application in debug mode and debug our OWIN pipeline. If you press F5 to run your web application in debug mode, the VS will warn you if you don't have a Web.config file:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_12.png" width="644" height="303" /></a></p>
<p>Not sure why this is still needed but after selecting the first choice, VS will add a web.config file with the necessary configuration and we can now debug our middlewares:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_9.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Good-Old_7762/image_thumb_13.png" width="644" height="177" /></a></p>
<p>Very neat. This sample I demonstrated is also <a href="https://github.com/tugberkugurlu/OwinSamples/tree/master/OwinHostVsIntegrationSample">available on GitHub</a>. Enjoy the goodness :)</p>