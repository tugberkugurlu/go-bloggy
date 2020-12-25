---
id: 705ee6c3-ec95-4b80-9910-0d8e8a95c81a
title: ASP.NET Web API Catch-All Route Parameter Binding
abstract: ASP.NET Web API has a concept of Catch-All routes but the frameowk doesn't
  automatically bind catch-all route values to a string array. Let's customize it
  with a custom HttpParameterBinding.
created_at: 2012-08-29 13:14:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- asp-net-web-api-catch-all-route-parameter-binding
---

<p>I just realized that ASP.NET Web API doesn&rsquo;t bind catch-all route values as ASP.NET MVC does. If you are not familiar with catch all routing, <a href="http://stephenwalther.com">Stephen Walter</a> has a great explanation on <a href="http://stephenwalther.com/archive/2009/02/06/chapter-2-understanding-routing.aspx">his article under the "Using Catch-All Routes" section</a>.</p>
<p>In ASP.NET MVC, when you have a route as below, you can retrieve the values of the catch all parameter as string array.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>RouteTable.Routes.MapRoute(
    <span style="color: #a31515;">"CatchAllRoute"</span>,
    <span style="color: #a31515;">"blog/tags/{*tags}"</span>,
    <span style="color: blue;">new</span> { controller = <span style="color: #a31515;">"blog"</span>, action = <span style="color: #a31515;">"tags"</span> }
);</pre>
</div>
</div>
<p>The controller action would look like as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> BlogController : Controller {

    <span style="color: blue;">public</span> ActionResult Tags(<span style="color: blue;">string</span>[] tags) { 

        <span style="color: green;">//...</span>
    }
}</pre>
</div>
</div>
<p>In ASP.NET Web API, we don&rsquo;t have that capability. If we have a catch-all route, we could retrieve it as string and parse it manually but that would be so lame to do it inside the controller, isn&rsquo;t it? There must be a better way. Well, there is! We can create a custom <a href="http://msdn.microsoft.com/en-us/library/system.web.http.controllers.httpparameterbinding(v=vs.108).aspx">HttpParameterBinding</a> and register it globally for string arrays. If you are interested in learning more about parameter binding in ASP.NET Web API, you might wanna have a look at <a href="http://blogs.msdn.com/b/jmstall/">Mike Stall</a>&rsquo;s <a href="http://blogs.msdn.com/b/jmstall/archive/2012/05/11/webapi-parameter-binding-under-the-hood.aspx">WebAPI Parameter binding under the hood blog post</a>. In our case, the custom HttpParameterBinding we want to create looks like as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CatchAllRouteParameterBinding : HttpParameterBinding {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">string</span> _parameterName;
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">char</span> _delimiter;

    <span style="color: blue;">public</span> CatchAllRouteParameterBinding(
        HttpParameterDescriptor descriptor, <span style="color: blue;">char</span> delimiter) : <span style="color: blue;">base</span>(descriptor) {

        _parameterName = descriptor.ParameterName;
        _delimiter = delimiter;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> Task ExecuteBindingAsync(
        System.Web.Http.Metadata.ModelMetadataProvider metadataProvider,
        HttpActionContext actionContext,
        CancellationToken cancellationToken) {

        <span style="color: blue;">var</span> routeValues = actionContext.ControllerContext.RouteData.Values;
            
        <span style="color: blue;">if</span> (routeValues[_parameterName] != <span style="color: blue;">null</span>) {

            <span style="color: blue;">string</span>[] catchAllValues = 
                routeValues[_parameterName].ToString().Split(_delimiter);

            actionContext.ActionArguments.Add(_parameterName, catchAllValues);
        }
        <span style="color: blue;">else</span> {

            actionContext.ActionArguments.Add(_parameterName, <span style="color: blue;">new</span> <span style="color: blue;">string</span>[0]);
        }

        <span style="color: blue;">return</span> Task.FromResult(0);
    }
}</pre>
</div>
</div>
<p>All the necessary information has been provided to us inside the <a href="http://msdn.microsoft.com/en-us/library/system.web.http.controllers.httpparameterbinding.executebindingasync(v=vs.108).aspx">ExecuteBindingAsync</a> method. From there, we simply grab the values from the RouteData and see if there is any route value whose route parameter name is the same as the action method parameter name. If there is one, we go ahead and split the values using the delimiter char provided to us. If there is no, we just attach an empty string array for the parameter. At the end, we let our caller know that we are done by returning a pre-completed Task object. I was using .NET 4.5, so I simply used <a href="http://msdn.microsoft.com/en-us/library/hh194922.aspx">FromResult</a> method of Task class. If you are on .NET 4.0, you can return a completed task by using <a href="http://msdn.microsoft.com/en-us/library/dd449174(v=vs.100).aspx">TaskCompletionSource</a> class.</p>
<p>The following code is the our catch-all route.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start() {

    <span style="color: blue;">var</span> config = GlobalConfiguration.Configuration;

    config.Routes.MapHttpRoute(
        <span style="color: #a31515;">"BlogpostTagsHttpApiRoute"</span>,
        <span style="color: #a31515;">"api/blogposts/tags/{*tags}"</span>,
        <span style="color: blue;">new</span> { controller = <span style="color: #a31515;">"blogposttags"</span> }
    );
}</pre>
</div>
</div>
<p>The last thing is that we need to register a rule telling that if there is an action method parameter which is a type of string array, go ahead and use our custom HttpParameterBinding.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start() {

    <span style="color: blue;">var</span> config = GlobalConfiguration.Configuration;

    <span style="color: green;">//...</span>

    config.ParameterBindingRules.Add(<span style="color: blue;">typeof</span>(<span style="color: blue;">string</span>[]),
        descriptor =&gt; <span style="color: blue;">new</span> CatchAllRouteParameterBinding(descriptor, <span style="color: #a31515;">'/'</span>));
}</pre>
</div>
</div>
<p>Now, if we send a request to /api/blogposts/tags/asp-net/asp-net-web-api, we would see that our action method parameter is bound.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Catch-All-Route-Paramete_FD94/image.png"><img height="154" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Catch-All-Route-Paramete_FD94/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>So far so good but we might not want to register our HttpParameterBinding rule globally. Instead, we might want to specify it manually when we require it. Well, we can do that as well. We just need to create a ParameterBindingAttribute to get our custom HttpParameterBinding so that it will be used to bind the action method parameter.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> BindCatchAllRouteAttribute : ParameterBindingAttribute {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">char</span> _delimiter;

    <span style="color: blue;">public</span> BindCatchAllRouteAttribute(<span style="color: blue;">char</span> delimiter) {

        _delimiter = delimiter;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> HttpParameterBinding GetBinding(HttpParameterDescriptor parameter) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span> CatchAllRouteParameterBinding(parameter, _delimiter);
    }
}</pre>
</div>
</div>
<p>As you can see, it is dead simple. The only thing we need to do now is to apply this attribute to our action parameter:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> BlogPostTagsController : ApiController {

    <span style="color: green;">//GET /api/blogposts/tags/asp-net/asp-net-web-api</span>
    <span style="color: blue;">public</span> HttpResponseMessage Get([BindCatchAllRoute(<span style="color: #a31515;">'/'</span>)]<span style="color: blue;">string</span>[] tags) {

        <span style="color: green;">//TODO: Do your thing here...</span>

        <span style="color: blue;">return</span> <span style="color: blue;">new</span> HttpResponseMessage(HttpStatusCode.OK);
    }
}</pre>
</div>
</div>
<p>When we send a request to /api/blogposts/tags/asp-net/asp-net-web-api, we shouldn&rsquo;t see any difference.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Catch-All-Route-Paramete_FD94/image_3.png"><img height="154" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Catch-All-Route-Paramete_FD94/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>I am still discovering how parameter and model binding works inside the ASP.NET Web API. So, there is a good chance that I did something wrong here :) If you spot it, please let me know :)</p>