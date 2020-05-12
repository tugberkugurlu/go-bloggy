---
id: 5e5f0d62-e491-4b64-bc90-dc012d9b4383
title: Setting IHostingEnvironment.IsDevelopment as True in an ASP.NET 5 Application
abstract: Wondering why IHostingEnvironment.IsDevelopment returns false even when
  you are on you development machine? I did indeed wonder and here is why :)
created_at: 2015-09-13 15:35:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET 5
slugs:
- setting-ihostingenvironment-isdevelopment-as-true-in-an-asp-net-5-application
---

<p>I am now in <a href="https://de.wikipedia.org/wiki/Frankfurt_am_Main">Frankfurt</a>, sipping my coffee in a Starbucks shop and enjoying its rubbish internet connection (as usual). I have just thrown away 10 minutes from my life by trying to figure out why my <a href="www.tugberkugurlu.com/tags/asp-net-5">ASP.NET 5</a> application wasn't showing the error page. So, I wanted to write this quick and dirty blog post for people who will potentially have the same problem :)</p> <p>Here is the piece of code I have inside the Configure method of my Startup class:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">void</span> Configure(IApplicationBuilder app, IHostingEnvironment env, ILoggerFactory loggerFactory)
{
    <span style="color: green">// ...</span>

    <span style="color: green">// Add the following to the request pipeline only in development environment.</span>
    <span style="color: blue">if</span> (env.IsDevelopment())
    {
        app.UseErrorPage();
    }
    <span style="color: blue">else</span>
    {
        <span style="color: green">// Add Error handling middleware which catches all application specific errors and</span>
        <span style="color: green">// send the request to the following path or controller action.</span>
        app.UseErrorHandler(<span style="color: #a31515">"/Home/Error"</span>);
    }

    <span style="color: green">// ...</span></pre></div></div>
<p>So, I should see <a href="https://github.com/aspnet/Diagnostics/blob/17a0fc7c2d5ce3f0ce56f27c14b2eefb279fec91/src/Microsoft.AspNet.Diagnostics/ErrorPageExtensions.cs#L20-L23">the beautiful and detailed error page</a> when I am in development. However, all I get is nothing but an empty response body when I run the application:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/4cc71789-ebe9-407b-8423-6bb5ee22f022.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/73ad9120-de8f-4e3c-9ad7-aae96739a319.png" width="644" height="328"></a></p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5a1a5532-65d1-4b19-afeb-4bb400c518ce.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/94eff907-275e-4aff-b1c0-fb15d54c70a7.png" width="644" height="156"></a></p>
<p>With <a href="https://github.com/aspnet/Hosting/blob/dfe8c39fe0858abec9d72e5582b4551bfca456ba/src/Microsoft.AspNet.Hosting.Abstractions/HostingEnvironmentExtensions.cs#L17-L20">a little bit of digging</a>, I remembered that your environment is being determined through an environment variable which is <strong>ASPNET_ENV</strong>. Setting this to <strong>Development</strong> will return true from IHostingEnvironment.IsDevelopment. Also, the <a href="https://github.com/aspnet/Hosting/blob/3e6585dcc81777676665c844af988ddfec87aba7/src/Microsoft.AspNet.Hosting.Abstractions/IHostingEnvironment.cs#L11">IHostingEnvironment.EnvironmentName</a> will get you the value of this environment variable. You can set this environment variable per process, per user or per machine. Whatever floats your boat. I have set this for process on windows with the below script and I was able to get the lovely error page:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>set ASPNET_ENV<span style="color: gray">=</span><span style="color: #a31515">Development</span></pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/53bbf9f9-6e95-424a-aed4-99e6e92bd5ef.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/430d1348-0bfd-4ac5-aa47-b2353f2dfc4d.png" width="644" height="226"></a></p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/3bc29c9d-f24a-4182-b5a7-c49d05e8fa93.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/739e1ee4-2cd4-42fc-8c28-df7bc0a8a7d1.png" width="644" height="406"></a></p>
<p>When you are on Visual Studio 2015, you can handle this better by adding a <strong>launchSettings.json </strong>file as <a href="https://github.com/aspnet/live.asp.net/blob/c1a75fd398ac85b03c6fdd153f6e9713e40b67bd/src/live.asp.net/Properties/launchSettings.json">here</a>. VS will pick this up and set the environment variable for IIS Express process.</p>  