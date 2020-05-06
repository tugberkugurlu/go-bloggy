---
title: Getting the Client’s IP Address in ASP.NET vNext Web Applications
abstract: I was wondering about how to get client’s IP address inside an ASP.NET vNext
  web application. It’s a little tricky than it should be but I finally figured it
  out :)
created_at: 2014-10-30 10:04:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET vNext
slugs:
- getting-the-clients-ip-address-in-asp-net-vnext-web-applications
---

<p>I was wondering about how to get client’s IP address inside an <a href="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-vnext-by-setting-up-the-environment-from-scratch">ASP.NET vNext</a> web application. It’s a little tricky than it should be but I finally figured it out. So, I decided to write it down here so that anyone else can see how it’s done right now.</p> <blockquote> <p><strong>BIG ASS CAUTION!</strong> At the time of this writing, I am using <strong>KRE 1.0.0-beta2-10648</strong> version. As things are moving really fast in this new world, it’s very likely that the things explained here will have been changed as you read this post. So, be aware of this and try to explore the things that are changed to figure out what are the corresponding new things.  <p>Also, inside this post I am referencing a lot of things from ASP.NET GitHub repositories. In order to be sure that the links won’t break in the future, I’m actually referring them by getting permanent links to the files on GitHub. So, these links are actually referring the files from the latest commit at the time of this writing and they have a potential to be changed, too. Read the "<a href="https://help.github.com/articles/getting-permanent-links-to-files/">Getting permanent links to files</a>" post to figure what this actually is.</p></blockquote> <p>First thing to highlight here is "<a href="https://github.com/aspnet/HttpAbstractions/tree/dev/src/Microsoft.AspNet.FeatureModel">HTTP Features</a>" feature (mouthful, I know :)). It’s a really exciting design if you ask me and basically, the HTTP server underneath your application will implement certain features and your application can work with these implementations. Some of the features that a server might implement can be found inside the <a href="https://github.com/aspnet/HttpAbstractions/tree/dev/src/Microsoft.AspNet.HttpFeature">Microsoft.AspNet.HttpFeature</a> package. However, it’s not limited to these features only. You watch <a href="http://channel9.msdn.com/Shows/Web+Camps+TV/ASP-NET-vNext-with-Chris-Ross">Chris Ross talking about HTTP Features on Web Camps TV</a> to find out more about this as he gives a nice overview on HTTP Features design.</p> <p>One of these HTTP features is <a href="https://github.com/aspnet/HttpAbstractions/blob/fee220569aa108078ab0e231080724eb74ec8b2d/src/Microsoft.AspNet.HttpFeature/IHttpConnectionFeature.cs">HTTP connection feature</a> (IHttpConnectionFeature if you want to be more specific). This feature is responsible for giving local and remote connection information for current request and this is where we can get client’s IP Address. You can try to get IHttpConnectionFeature from the <a href="https://github.com/aspnet/HttpAbstractions/blob/fee220569aa108078ab0e231080724eb74ec8b2d/src/Microsoft.AspNet.Http/HttpContext.cs">HttpContext</a> instance which you get per each request. HttpContext has two methods called <a href="https://github.com/aspnet/HttpAbstractions/blob/fee220569aa108078ab0e231080724eb74ec8b2d/src/Microsoft.AspNet.Http/HttpContext.cs#L41">GetFeature</a> and <a href="https://github.com/aspnet/HttpAbstractions/blob/fee220569aa108078ab0e231080724eb74ec8b2d/src/Microsoft.AspNet.Http/HttpContext.cs#L45-L48">one of them accepts a generic parameter</a>. You can use this method to retrieve a specific feature you want to reach out to. The below Startup class shows you how you can get the client’s IP Address through the IHttpConnectionFeature.</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">using</span> Microsoft.AspNet.Builder;
<span style="color: blue">using</span> Microsoft.AspNet.Http;
<span style="color: blue">using</span> System.Threading.Tasks;
<span style="color: blue">using</span> Microsoft.AspNet.HttpFeature;

<span style="color: blue">namespace</span> RemoteIPAddressSample
{
    <span style="color: blue">public</span> <span style="color: blue">class</span> Startup
    {
        <span style="color: blue">public</span> <span style="color: blue">void</span> Configure(IApplicationBuilder app)
        {
            app.Run(async (ctx) =&gt;
            {
                IHttpConnectionFeature connection = ctx.GetFeature&lt;IHttpConnectionFeature&gt;();
                <span style="color: blue">string</span> ipAddress = connection != <span style="color: blue">null</span>
                    ? connection.RemoteIpAddress.ToString()
                    : <span style="color: blue">null</span>;

                await ctx.Response.WriteAsync(<span style="color: #a31515">"IP Address: "</span> + ipAddress);
            });
        }
    }
}</pre></div></div>
<p>When I hit the HTTP endpoint now, I’ll see my client’s IP address inside the response body:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/dd7c0b7c-7aa3-41b8-9187-8ce7966930f2.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/dc1f9767-7b24-4fb4-845c-1142ca9cb89f.png" width="644" height="190"></a></p>
<p>We can even make this look nicer by refactoring this into an extension method.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">static</span> <span style="color: blue">class</span> HttpContextExtensions
{
    <span style="color: blue">public</span> <span style="color: blue">static</span> <span style="color: blue">string</span> GetClientIPAddress(<span style="color: blue">this</span> HttpContext context)
    {
        <span style="color: blue">if</span>(context == <span style="color: blue">null</span>)
        {
            <span style="color: blue">throw</span> <span style="color: blue">new</span> ArgumentNullException(<span style="color: #a31515">"context"</span>);
        }

        IHttpConnectionFeature connection = context.GetFeature&lt;IHttpConnectionFeature&gt;();

        <span style="color: blue">return</span> connection != <span style="color: blue">null</span>
            ? connection.RemoteIpAddress.ToString()
            : <span style="color: blue">null</span>;
    }
}</pre></div></div>
<p>Way better. We can now use this nice looking HttpContext extension instead:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> Startup
{
    <span style="color: blue">public</span> <span style="color: blue">void</span> Configure(IApplicationBuilder app)
    {
        app.Run(async (ctx) =&gt;
        {
            await ctx.Response.WriteAsync(<span style="color: #a31515">"IP Address: "</span> + ctx.GetClientIPAddress());
        });
    }
}</pre></div></div>
<p>One thing to mention here is that the HTTP server which your application is running on top of may choose not to implement this feature. It sounds like a very unreasonable situation as something this fundamental should be there all the time but it’s what it is. So, you should be always careful when reaching out to properties of this feature as you may get null values when you call GetFeature&lt;T&gt;. This rule probably applies most of the HTTP features but I wasn’t able to find any certain information whether there is a check at the very early stages of the pipeline to refuse to use the server if specific features are not implemented by the server. I’ll probably update the post when I find this out.</p>  