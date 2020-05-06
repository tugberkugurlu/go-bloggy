---
title: ASP.NET Web API Custom RouteDataMapping (MediaTypeMapping)
abstract: In this post, we will create RouteDataMapping. This custom MediaTypeMapping
  will allow us to involve the decision-making process about the response format according
  to RouteData values.
created_at: 2012-03-03 21:05:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- asp-net-web-api-custom-routedatamapping-mediatypemapping
---

<p>I have talked about on <a target="_blank" href="http://www.tugberkugurlu.com/archive/asp-net-web-api-mediatypeformatters-with-mediatypemappings" title="http://www.tugberkugurlu.com/archive/asp-net-web-api-mediatypeformatters-with-mediatypemappings">ASP.NET Web API Content-Negotiation algorithm and MediaTypeMapping</a> on my previous post. As I said there, creating one custom MediaTypeMapping is fairly simple.</p>
<p>In this post, we will create RouteDataMapping. This custom MediaTypeMapping will allow us to involve the decision-making process about the response format according to RouteData values. Here is the complete implementation:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> RouteDataMapping : MediaTypeMapping {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">string</span> _routeDataValueName;
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">string</span> _routeDataValueValue;

    <span style="color: blue;">public</span> RouteDataMapping(
        <span style="color: blue;">string</span> routeDataValueName, 
        <span style="color: blue;">string</span> routeDataValueValue, 
        MediaTypeHeaderValue mediaType) : <span style="color: blue;">base</span>(mediaType) {

        _routeDataValueName = routeDataValueName;
        _routeDataValueValue = routeDataValueValue;
    }

    <span style="color: blue;">public</span> RouteDataMapping(
        <span style="color: blue;">string</span> routeDataValueName, 
        <span style="color: blue;">string</span> routeDataValueValue, 
        <span style="color: blue;">string</span> mediaType) : <span style="color: blue;">base</span>(mediaType) {

        _routeDataValueName = routeDataValueName;
        _routeDataValueValue = routeDataValueValue;
    }

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> <span style="color: blue;">double</span> OnTryMatchMediaType(
        System.Net.Http.HttpResponseMessage response) {

        <span style="color: blue;">return</span> (
            response.RequestMessage.GetRouteData().
            Values[_routeDataValueName].ToString() == _routeDataValueValue
        ) ? 1.0 : 0.0;
    }

    <span style="color: green;">//Don't use this</span>
    <span style="color: green;">//This will be removed on the first RC (according to team members)</span>
    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> <span style="color: blue;">double</span> OnTryMatchMediaType(
        System.Net.Http.HttpRequestMessage request) {

        <span style="color: blue;">throw</span> <span style="color: blue;">new</span> NotImplementedException();
    }
}</pre>
</div>
</div>
<p>The implementation is fairly simple.One thing that you might notice is that we are returning double in order to tell the framework if it is a match or not. Here is the reason why:</p>
<blockquote>
<p><em>The returned double value is used by the conneg algorithm to find the appropriate formatter to write. Its similar to how you can set the quality value in Accept header. - Kiran Challa</em></p>
</blockquote>
<p>Let&rsquo;s try this out. we will use the same sample on my previous post but we will make a few changes. First of all, we will change our route a little to add an extension. Keep in mind that you do not need to make this an extension. It will work for every RouteData value.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>GlobalConfiguration.Configuration.Routes.MapHttpRoute(
    <span style="color: #a31515;">"defaultHttpRoute"</span>,
    routeTemplate: <span style="color: #a31515;">"api/{controller}.{extension}"</span>,
    defaults: <span style="color: blue;">new</span> { },
    constraints: <span style="color: blue;">new</span> { extension = <span style="color: #a31515;">"json|xml"</span> }
);</pre>
</div>
</div>
<p>Then we will make sure that we have the following registry inside our Web.config file in order for our Urls with extensions to work.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">system.webServer</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">modules</span> <span style="color: red;">runAllManagedModulesForAllRequests</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">true</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">system.webServer</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>And finally, we will hook up our RouteDataMapping to the formatters:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>GlobalConfiguration.Configuration.Formatters.JsonFormatter.
    MediaTypeMappings.Add(
        <span style="color: blue;">new</span> RouteDataMapping(
            <span style="color: #a31515;">"extension"</span>, <span style="color: #a31515;">"json"</span>, <span style="color: #a31515;">"application/json"</span>
    )
);

GlobalConfiguration.Configuration.Formatters.XmlFormatter.
    MediaTypeMappings.Add(
        <span style="color: blue;">new</span> RouteDataMapping(
            <span style="color: #a31515;">"extension"</span>, <span style="color: #a31515;">"xml"</span>, <span style="color: #a31515;">"application/xml"</span>
    )
);</pre>
</div>
</div>
<p>Now, when you navigate to <strong>/api/cars.json</strong>, you will get the data as <strong>json</strong>. If you navigate to <strong>/api/cars.xml</strong>, you will get the result as <strong>xml</strong> as below.</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/638163d39ce1_685/routeDataMapping.png"><img height="345" width="644" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/638163d39ce1_685/routeDataMapping_thumb.png" alt="routeDataMapping" border="0" title="routeDataMapping" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>