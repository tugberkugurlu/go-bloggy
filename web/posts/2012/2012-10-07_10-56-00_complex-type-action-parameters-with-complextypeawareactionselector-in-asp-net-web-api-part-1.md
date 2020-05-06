---
title: Complex Type Action Parameters with ComplexTypeAwareActionSelector in ASP.NET
  Web API - Part 1
abstract: We will see how to make complex type action parameters play nice with controller
  action selection in ASP.NET Web API by using ComplexTypeAwareActionSelector from
  WebAPIDoodle NuGet package.
created_at: 2012-10-07 10:56:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- complex-type-action-parameters-with-complextypeawareactionselector-in-asp-net-web-api-part-1
---

<p>Couple of days ago, I wrote a blog post on <a href="http://www.tugberkugurlu.com/archive/complex-type-action-parameters-and-controller-action-selection-with-asp-net-web-api">complex type action parameters and controller action selection with ASP.NET Web API</a> with a solution that I&rsquo;ve come up with and I encourage you to check that blog post out to get a sense of what this is really about. However, that solution was sort of noisy. Yesterday, I tweaked the ComplexTypeAwareActionSelector implementation a bit and now it directly supports the complex type action parameters without any additional attributes for the action methods. In this post, we will see how we can use it and in the next post, we will see how it works under the covers.</p>
<p>First of all, install the latest <a href="https://github.com/WebAPIDoodle/WebAPIDoodle">WebAPIDoodle</a> package rom the Official NuGet feed:</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package WebAPIDoodle </code></p>
</div>
<p>Let&rsquo;s go through the scenario briefly. Assume that we have the following controller with two action methods that are expected to serve for different GET requests:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsByCategoryRequestCommand {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> CategoryId { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Page { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    [Range(1, 50)]
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Take { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsByColorRequestCommand {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> ColorId { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Page { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    [Range(1, 50)]
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Take { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

[InvalidModelStateFilter]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController {

    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByCategoryId(
        [FromUri]CarsByCategoryRequestCommand cmd) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>,
            <span style="color: #a31515;">"Car 3"</span>
        };
    }

    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByColorId(
        [FromUri]CarsByColorRequestCommand cmd) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>
        };
    }
}</pre>
</div>
</div>
<p>If we now send a GET request to /api/cars?colorId=23&amp;page=2&amp;take=12 with the default action selector registered, we would get the ambiguity error message because the default action selector doesn&rsquo;t consider the complex type action parameters while performing the action selection.</p>
<p><img height="367" width="640" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action-Paramete.NET-Web-API_11EF9/image.png" /></p>
<p>Let&rsquo;s replace the default action selector with our ComplexTypeAwareActionSelector as below. Note that ComplexTypeAwareActionSelector preserves all the features of the ApiControllerActionSelector.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    <span style="color: blue;">var</span> config = GlobalConfiguration.Configuration;
    config.Routes.MapHttpRoute(
        <span style="color: #a31515;">"DefaultApiRoute"</span>,
        <span style="color: #a31515;">"api/{controller}/{id}"</span>,
        <span style="color: blue;">new</span> { id = RouteParameter.Optional }
    );

    <span style="color: green;">// Replace the default action IHttpActionSelector with</span>
    <span style="color: green;">// WebAPIDoodle.Controllers.ComplexTypeAwareActionSelector</span>
    config.Services.Replace(
        <span style="color: blue;">typeof</span>(IHttpActionSelector),
        <span style="color: blue;">new</span> ComplexTypeAwareActionSelector());
}</pre>
</div>
</div>
<p>As explained inside <a href="http://www.tugberkugurlu.com/archive/complex-type-action-parameters-and-controller-action-selection-with-asp-net-web-api">the previous related post</a>, we previously had to mark the action methods with UriParametersAttribute to give a hint about the action parameters we want to support. However, with the current implementation of the the ComplexTypeAwareActionSelector, it just works as it is. Only thing required to perform is to mark the complex action parameter with <a href="http://msdn.microsoft.com/en-us/library/system.web.http.fromuriattribute(v=vs.108).aspx">FromUriAttribute</a>. By marking the complex type parameters with FromUriAttribute, you are making it possible to bind the route and query string values.</p>
<p>After replacing the default action selector with our own implementation, we will see it working if we send a GET request to /api/cars?colorId=23&amp;page=2&amp;take=12.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/439a78ea23ca_71A4/image.png"><img height="224" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/439a78ea23ca_71A4/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Now, let&rsquo;s send a request to /api/cars?categoryId=23&amp;page=2&amp;take=12 and see what we will get back:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/439a78ea23ca_71A4/imageac93a1be-95c5-4f38-9540-73bc5fe72c35.png"><img height="224" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/439a78ea23ca_71A4/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Working perfectly as expected. The ComplextTypeAwareActionSelector considers simple types inside a complex type parameter which are all primitive .NET types, System.String, System.DateTime, System.Decimal, System.Guid, System.DateTimeOffset, System.TimeSpan and underlying simple types (e.g: Nullable&lt;System.Int32&gt;).</p>
<p>In the next post, we will see how ComplextTypeAwareActionSelector works and behaves with complex type action parameters. Stay tuned! ;)</p>