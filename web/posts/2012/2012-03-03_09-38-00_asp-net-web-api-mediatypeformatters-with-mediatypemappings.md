---
id: ea3827a4-7923-423e-9ff5-25edb9550f9c
title: ASP.NET Web API MediaTypeFormatters With MediaTypeMappings
abstract: We will see how Content-Negotiation (Conneg) Algorithm works on ASP.NET
  Web API with MediaTypeFormatters and MediaTypeMappings
created_at: 2012-03-03 09:38:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- asp-net-web-api-mediatypeformatters-with-mediatypemappings
---

<p>If you have read my post on <a title="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples" href="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples" target="_blank">Getting Started With ASP.NET Web API</a>, you probably saw me talking about exposing your data to the world with various types of formats. This feature has been made possible by formatters. Formatters handles serializing and deserializing strongly-typed objects.</p>
<p>For two days, I have been really looking into formatters and explored a lot of useful stuff and I thought that sharing those would be great.</p>
<blockquote>
<p>Before diving deeply, ASP.NET Web API team moderating <a title="http://forums.asp.net/1246.aspx/1?Web+API" href="http://forums.asp.net/1246.aspx/1?Web+API" target="_blank">ASP.NET Web API Forum</a> pretty often and provide you a way to solve your problems. Also, community is heavily involved. Hence the framework is so new, I was stuck at a few places but team and the community helped me out a lot.</p>
<p>One other place you should be aware of is <a title="http://msdn.microsoft.com/en-us/library/hh849329(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/hh849329(v=vs.108).aspx" target="_blank">ASP.NET Web API Reference</a> on MSDN.</p>
</blockquote>
<p><a title="http://asp.net" href="http://asp.net" target="_blank">ASP.NET Web API</a>&nbsp;beta was shiped with 3 different formatters: <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.jsonmediatypeformatter(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.jsonmediatypeformatter(v=vs.108).aspx" target="_blank">JsonMediaTypeFormatter</a>, <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.xmlmediatypeformatter(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.xmlmediatypeformatter(v=vs.108).aspx" target="_blank">XmlMediaTypeFormatter</a>, <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.formurlencodedmediatypeformatter(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.formurlencodedmediatypeformatter(v=vs.108).aspx" target="_blank">FormUrlEncodedMediaTypeFormatter</a>. All these classes are derived from <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter(v=vs.108).aspx" target="_blank">MediaTypeFormatter</a> abstract class. As you might guess, it is fairly easy to create one and hook it up but in this post, I won&rsquo;t talk about custom MediaTypeFormatters. I would like to talk about how they are being chosen and assigned to process the request by the framework, especially on MediaTypeMappings.</p>
<p>ASP.NET Web API decides which formatter to process request with according to its Content-Negotiation (<em>Conneg</em>) Algorithm. <a title="http://blogs.msdn.com/b/kiranchalla/" href="http://blogs.msdn.com/b/kiranchalla/" target="_blank">Kiran Challa</a> has two great blog posts on this:</p>
<ul>
<li><a title="http://blogs.msdn.com/b/kiranchalla/archive/2012/02/25/content-negotiation-in-asp-net-mvc4-web-api-beta-part-1.aspx" href="http://blogs.msdn.com/b/kiranchalla/archive/2012/02/25/content-negotiation-in-asp-net-mvc4-web-api-beta-part-1.aspx" target="_blank">Content Negotiation in ASP.NET MVC4 Web API Beta &ndash; Part 1</a></li>
<li><a title="http://blogs.msdn.com/b/kiranchalla/archive/2012/02/27/content-negotiation-in-asp-net-mvc4-web-api-beta-part-2.aspx" href="http://blogs.msdn.com/b/kiranchalla/archive/2012/02/27/content-negotiation-in-asp-net-mvc4-web-api-beta-part-2.aspx" target="_blank">Content Negotiation in ASP.NET MVC4 Web API Beta &ndash; Part 2</a></li>
</ul>
<p>On these posts, you will find how the Conneg algorithm works inside the framework. It has various options and as default it looks at the http headers to decide the most suitable format.</p>
<p>For this post, I have created a very simple Web API project. I did that by creating an empty ASP.NET Web Application, installing AspNetWebApi nuget package (had to install System.Json package separately). Then I registered my route:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    GlobalConfiguration.Configuration.Routes.MapHttpRoute(
        <span style="color: #a31515;">"defaultHttpRoute"</span>,
        routeTemplate: <span style="color: #a31515;">"api/{controller}/{id}"</span>,
        defaults: <span style="color: blue;">new</span> { id = RouteParameter.Optional }
    );
}</pre>
</div>
</div>
<p>Finally, I created a simple API:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController {

    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] Get() {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span> <span style="color: blue;">string</span>[] { 
            <span style="color: #a31515;">"BMW"</span>,
            <span style="color: #a31515;">"Ferrari"</span>,
            <span style="color: #a31515;">"FIAT"</span>
        };
    }
}</pre>
</div>
</div>
<p>When I fire up the development web server IIS Express and navigate to <strong>/api/cars</strong>, I get the list of cars as expected. This is not the clearest way of explaining it, is it? Let&rsquo;s see the headers:</p>
<p><strong>Request:</strong></p>
<p><em>GET http://localhost:4446/api/cars HTTP/1.1<br />User-Agent: Fiddler<br />Host: localhost:4446</em></p>
<p><strong>Response:</strong></p>
<p><em>HTTP/1.1 200 OK<br />Cache-Control: no-cache<br />Pragma: no-cache<br />Transfer-Encoding: chunked<br />Content-Type: application/json; charset=utf-8<br />Expires: -1<br />Server: Microsoft-IIS/7.5<br />X-AspNet-Version: 4.0.30319<br />X-SourceFiles: =?UTF-8?B?RDpcRHJvcGJveFxBcHBzXEFTUE5FVFdlYkFQSVNhbXBsZXNcQ29ubmVnQWxnb3JpdGhtU2FtcGxlXHNyY1xDb25uZWdBbGdvcml0aG1TYW1wbGVcYXBpXGNhcnM=?=<br />X-Powered-By: ASP.NET<br />Date: Sat, 03 Mar 2012 10:52:49 GMT</em></p>
<p><em>18<br />["BMW","Ferrari","FIAT"]<br />0</em></p>
<p>As you see, we have the response back as json because it is the first formatter registered (yes, order matters) by default and we didn&rsquo;t specify which format we are interested in. When you add <em>"Accept: application/xml"</em> to your request, you will see that you will be getting the response back as xml. Approve it or not, this is the RESTFul way of negotiating between client and server. But sometimes we would like to decide the format according to QueryString. If so, you have an OOB support for this.</p>
<p><strong>Intro to MediaTypeMappings</strong></p>
<p>By default, Accept and Request Content-Type headers play role on deciding which format you serve. One other way of involving a formatter to process your request is MediaTypeMapping.</p>
<p>MediaTypeMapping provides a way for us to participate the Conneg algorithm decision making process and decide if we would like the formatter to take part in writing the response. There are several built in MediaTypeMappings (actually 4) supported out of the box. These are <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.querystringmapping(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.querystringmapping(v=vs.108).aspx" target="_blank"><em>QueryStringMapping</em></a><em>, </em><a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.requestheadermapping(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.requestheadermapping(v=vs.108).aspx" target="_blank"><em>RequestHeaderMapping</em></a><em>, </em><a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.uripathextensionmapping(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.uripathextensionmapping(v=vs.108).aspx" target="_blank"><em>UriPathExtensionMapping</em></a><em>, </em><a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediarangemapping(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediarangemapping(v=vs.108).aspx" target="_blank"><em>MediaRangeMapping</em></a><em>. </em>All these classes are derived from <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypemapping(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypemapping(v=vs.108).aspx" target="_blank">MediaTypeMapping</a> abstract class (yes, creating a custom one is tedious and I plan on writing a post on that as well).<em>&nbsp;</em>We have these mappings and the other great stuff is that all default formatters has a hook up point in order to register mappings.</p>
<p>Let&rsquo;s assume that we would like to decide the format of response based on a query string value as well. As we have QuesryStringMapping, we can use this and can provide our data on json format if request comes with ?format=json quesry string and xml format if it is ?format=xml. Here is the configuration in order to enable this:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    GlobalConfiguration.Configuration.Routes.MapHttpRoute(
        <span style="color: #a31515;">"defaultHttpRoute"</span>,
        routeTemplate: <span style="color: #a31515;">"api/{controller}"</span>
    );

    GlobalConfiguration.Configuration.Formatters.JsonFormatter.
        MediaTypeMappings.Add(
            <span style="color: blue;">new</span> QueryStringMapping(
                <span style="color: #a31515;">"format"</span>, <span style="color: #a31515;">"json"</span>, <span style="color: #a31515;">"application/json"</span>
        )
    );

    GlobalConfiguration.Configuration.Formatters.XmlFormatter.
        MediaTypeMappings.Add(
            <span style="color: blue;">new</span> QueryStringMapping(
                <span style="color: #a31515;">"format"</span>, <span style="color: #a31515;">"xml"</span>, <span style="color: #a31515;">"application/xml"</span>
        )
    );
}</pre>
</div>
</div>
<p>When we make a request with accept header and format query string, we will see that framework honors our mapping registrations:</p>
<p><strong>Request:</strong></p>
<p><em>GET http://localhost:4446/api/cars?<strong>format=xml</strong> HTTP/1.1<br />User-Agent: Fiddler<br />Host: localhost:4446<br />Accept: <strong>appication/json</strong></em></p>
<p><strong>Response:</strong></p>
<p><em>HTTP/1.1 200 OK<br />Cache-Control: no-cache<br />Pragma: no-cache<br />Transfer-Encoding: chunked<br />Content-Type: </em><em><strong>application/xml<br /></strong>Expires: -1<br />Server: Microsoft-IIS/7.5<br />X-AspNet-Version: 4.0.30319<br />X-SourceFiles: =?UTF-8?B?RDpcRHJvcGJveFxBcHBzXEFTUE5FVFdlYkFQSVNhbXBsZXNcQ29ubmVnQWxnb3JpdGhtU2FtcGxlXHNyY1xDb25uZWdBbGdvcml0aG1TYW1wbGVcYXBpXGNhcnM=?=<br />X-Powered-By: ASP.NET<br />Date: Sat, 03 Mar 2012 11:36:04 GMT</em></p>
<p><em>e9<br />&lt;?xml version="1.0" encoding="utf-8"?&gt;&lt;ArrayOfString xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema"&gt;&lt;string&gt;BMW&lt;/string&gt;&lt;string&gt;Ferrari&lt;/string&gt;&lt;string&gt;FIAT&lt;/string&gt;&lt;/ArrayOfString&gt;<br />0<br /></em></p>
<p>Pretty powerful stuff. Enjoy <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/6ff79dbb50c1_A02A/wlEmoticon-smile.png" /></p>