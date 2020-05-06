---
title: Donut Hole Caching In ASP.NET MVC by Using Child Actions and OutputCacheAttribute
abstract: This blog post demonstrates how to implement Donut Hole Caching in ASP.NET
  MVC by Using Child Actions and OutputCacheAttribute
created_at: 2011-10-30 09:40:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- C#
- Caching
- Razor
slugs:
- donut-hole-caching-in-asp-net-mvc-by-using-child-actions-and-outputcacheattribute
---

<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Donout_991E/homer-and-donut.jpg"><img style="background-image: none; margin: 0px 0px 15px 15px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border: 0px;" title="homer-and-donut" border="0" alt="homer-and-donut" align="right" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Donout_991E/homer-and-donut_thumb.jpg" width="174" height="244" /></a></p>
<p>One of the common issues of web applications is performance and serving the same content over and over again for hours, days, even months is certainly effecting the performance of our web applications. This is where <strong>server side caching </strong>comes in handy. But sometimes the whole page we are scripting on the server side is going to be static for hours and sometimes some portion of it.</p>
<p>The first one, caching the whole page, is sort of easy with almost any web development frameworks. The second one, caching a portion of your web page, is the tricky one.</p>
<p><a title="http://haacked.com" href="http://haacked.com" target="_blank">Phil Haack</a> blogged about <a title="http://haacked.com/archive/2009/05/12/donut-hole-caching.aspx" href="http://haacked.com/archive/2009/05/12/donut-hole-caching.aspx" target="_blank">Donut Hole Caching in ASP.NET MVC</a> a while back but it is a little outdated. In this quick blog post, I will try to show you the easiest way of implementing this feature with <a title="http://asp.net/mvc" href="http://asp.net/mvc" target="_blank">ASP.NET MVC</a>.</p>
<p><strong>Use Cases</strong></p>
<p>So, where might we need this Donut Hole Caching thing? For example, if you are listing categories on the side bar, you probably gathering those categories from a database. Add the fact that you render this part on every single page of your web application, it is a waste of time if your category list is not so dynamic. So, you may want to cache the category list part of the page so that you don&rsquo;t go to your database every single time.</p>
<p><strong>How to Do</strong></p>
<p>In ASP.NET MVC framework, we can cache the whole controller action with <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.outputcacheattribute.aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.outputcacheattribute.aspx" target="_blank">OutputCacheAttribute</a>&nbsp;<em>which represents an attribute that is used to mark an action method whose output will be cached</em> according to <a title="http://msdn.microsoft.com" href="http://msdn.microsoft.com" target="_blank">MSDN</a>. This is perfect but wait a second. We do not want to cache the entire action, we just want to cache the portion of it. This is where <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.html.childactionextensions.action(v=VS.98).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.html.childactionextensions.action(v=VS.98).aspx" target="_blank">ChildActionExtensions.Action</a> method, which <em>invokes a child action method and returns the result as an HTML string</em>, plays a huge role. Let me show you how these parts of the framework play nice together.</p>
<p>Have a look at the below controller action :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[ChildActionOnly]
[OutputCache(Duration=60)]
<span style="color: blue;">public</span> ActionResult sampleChildAction() {

    <span style="color: green;">//Put the thread at sleep for 3 seconds to see the difference.</span>
    System.Threading.Thread.Sleep(3000);

    <span style="color: green;">//Also pass the date time from here just to see that it is comming from here.</span>
    ViewBag.DateTime = DateTime.Now.ToString(<span style="color: #a31515;">"dd.MM.yyyy HH:mm.ss"</span>);

    <span style="color: blue;">return</span> View();
}</pre>
</div>
</div>
<p>A simple controller action method which returns ActionResult, nothing fancy going on except for <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.childactiononlyattribute.aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.childactiononlyattribute.aspx" target="_blank">ChildActionOnlyAttribute</a> which <em>represents an attribute that is used to indicate that an action method should be called only as a child action. </em>Let&rsquo;s look at the <strong>sampleChildAction </strong>view and I will try to explain ChildActionOnlyAttribute function after that.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@{
    Layout = null;
}

<span style="color: blue;">&lt;</span><span style="color: #a31515;">p</span><span style="color: blue;">&gt;</span>
    This portion of the web page was scripted on <span style="color: blue;">&lt;</span><span style="color: #a31515;">strong</span><span style="color: blue;">&gt;</span>@ViewBag.DateTime<span style="color: blue;">&lt;/</span><span style="color: #a31515;">strong</span><span style="color: blue;">&gt;</span> and I will be cached for <span style="color: blue;">&lt;</span><span style="color: #a31515;">strong</span><span style="color: blue;">&gt;</span>60<span style="color: blue;">&lt;/</span><span style="color: #a31515;">strong</span><span style="color: blue;">&gt;</span> seconds.
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">p</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>This html will be a part of our web page which will be cached. It doesn&rsquo;t mean anything by itself but we have created an action for this view which means that we can call this page directly from a browser. <strong>ChildActionOnlyAttribute</strong> exactly prevent users to call this kind of actions. You do not need to implement this attribute there but it is nice to know that it is there for us.</p>
<p>The controller action which will render the whole page is so simple as below and doesn&rsquo;t require any special thing for us to implement in order caching to work.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> ActionResult Index() {

    <span style="color: blue;">return</span> View();
}</pre>
</div>
</div>
<p>Let&rsquo;s also look at the view implementation :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@{
    ViewBag.Title = "Donut Hole Caching Sample";
}

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h2</span><span style="color: blue;">&gt;</span>Donut Hole Caching Sample<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h2</span><span style="color: blue;">&gt;</span>

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>Cached<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>

<span style="color: blue;">&lt;</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span>
    @Html.Action("sampleChildAction", 
      new { controller = "Sample", Area = "DonutHoleCaching" }
    )
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span>

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>Normal<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>

<span style="color: blue;">&lt;</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">p</span><span style="color: blue;">&gt;</span>
        This portion of the web page was scripted on <span style="color: blue;">&lt;</span><span style="color: #a31515;">strong</span><span style="color: blue;">&gt;</span>@DateTime.Now.ToString("dd.MM.yyy HH:mm.ss")<span style="color: blue;">&lt;/</span><span style="color: #a31515;">strong</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">p</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>What we are doing here is that rather than putting the part, which we would like to cache, directly here, we are calling it as child action. So, the framework will treat the child action as it does for normal action methods.</p>
<p>When we first call it, we will see something like below :</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Donout_991E/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Donout_991E/image_thumb.png" width="644" height="419" /></a></p>
<p>While I was calling this page, I was on hold for 3 seconds because we have put the thread at sleep for 3 seconds on our child action method to feel the difference as you can see on above code in order.</p>
<p>When I make the second call, I got something like this and I wasn&rsquo;t on hold for 3 seconds :</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Donout_991E/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Donout_991E/image_thumb_3.png" width="644" height="419" /></a></p>
<p>Did you notice the time difference? This proves that if the cache is valid, our child action method won&rsquo;t be rendered again. It will serve from the cache. Awesome, ha?</p>
<p>I decided to create an ASP.NET MVC project called <a title="http://mvcmiracleworker.tugberkugurlu.com" href="http://mvcmiracleworker.tugberkugurlu.com" target="_blank">MvcMiracleWorker</a> for this kind of small samples. You can find the complete source code from <a title="https://github.com" href="https://github.com" target="_blank">GitHub</a> : <a title="http://mvcmiracleworker.tugberkugurlu.com" href="http://mvcmiracleworker.tugberkugurlu.com"><a href="https://github.com/tugberkugurlu/MvcMiracleWorker">https://github.com/tugberkugurlu/MvcMiracleWorker</a></a></p>
<p>Behave well, use ASP.NET MVC <img style="border-style: none;" class="wlEmoticon wlEmoticon-winkingsmile" alt="Winking smile" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Donout_991E/wlEmoticon-winkingsmile.png" /></p>