---
title: Running the Content Negotiation (Conneg) Manually in ASP.NET Web API
abstract: In this post we will see how run the Content Negotiation (Conneg) manually
  in an ASP.NET Web API easily.
created_at: 2012-06-21 09:50:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- running-the-content-negotiation-conneg-manually-in-asp-net-web-api
---

<p>I just wanted to put this out to the World to show how we can run content negotiation manually. As far as I remember, Daniel Roth showed this on his TechEd NA 2012 talk (<a title="http://channel9.msdn.com/Events/TechEd/NorthAmerica/2012/DEV309" href="http://channel9.msdn.com/Events/TechEd/NorthAmerica/2012/DEV309">Building HTTP Services with ASP.NET Web API</a>) as well.</p>
<p>Why on god's green earth do we want do this? So far, I haven&rsquo;t been able to find any reason since all the necessary places have either methods or extension methods which does this in behalf of us. Anyway, here is the code that does the trick.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController {

    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] Get() {

        <span style="color: blue;">var</span> cars = <span style="color: blue;">new</span> <span style="color: blue;">string</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>,
            <span style="color: #a31515;">"Car 3"</span>
        };

        <span style="color: blue;">var</span> contentNegotiator = Configuration.Services.GetContentNegotiator();
        <span style="color: blue;">var</span> connegResult = contentNegotiator.Negotiate(
            <span style="color: blue;">typeof</span>(<span style="color: blue;">string</span>[]), Request, Configuration.Formatters
        );

        <span style="color: green;">//you have the proper formatter for your request in your hand</span>
        <span style="color: blue;">var</span> properFormatter = connegResult.Formatter;

        <span style="color: green;">//you have the proper MediaType for your request in your hand</span>
        <span style="color: blue;">var</span> properMediaType = connegResult.MediaType;

        <span style="color: blue;">return</span> cars;
    }
}</pre>
</div>
</div>
<p align="left">First of all, we grab the registered <a title="http://msdn.microsoft.com/en-us/library/hh944843(v=vs.108)" href="http://msdn.microsoft.com/en-us/library/hh944843(v=vs.108)">IContentNegotiator</a> from the DefaultServices. The dafault content negotiator is <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.defaultcontentnegotiator(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.defaultcontentnegotiator(v=vs.108).aspx">DefaultContentNegotiator</a>&nbsp;and you can replace it with your own implementation easily if you need to. Then, we run the <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.icontentnegotiator.negotiate(v=vs.108)" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.icontentnegotiator.negotiate(v=vs.108)">Negotiate</a> method of IContentNegotiator by providing the necessary parameters as described below:</p>
<ul>
<li>
<div align="left">First parameter we pass is the type of the object that we would like to send over the wire.</div>
</li>
<li>
<div align="left">The second one is the <a title="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage(v=vs.110)" href="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage(v=vs.110)">HttpRequestMessage</a> and we grab that value through the <a title="http://msdn.microsoft.com/en-us/library/system.web.http.apicontroller.request(v=vs.108)" href="http://msdn.microsoft.com/en-us/library/system.web.http.apicontroller.request(v=vs.108)">Request</a> property of <a title="http://msdn.microsoft.com/en-us/library/system.web.http.apicontroller(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.http.apicontroller(v=vs.108).aspx">ApiController</a>.</div>
</li>
<li>
<div align="left">The last parameter is the list of formatters. We grab the registered formatters from the <a title="http://msdn.microsoft.com/en-us/library/system.web.http.httpconfiguration(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.http.httpconfiguration(v=vs.108).aspx">HttpConfiguration</a> object.</div>
</li>
</ul>
<p align="left">There we have it. This method returns us a <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.contentnegotiationresult(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.contentnegotiationresult(v=vs.108).aspx">ContentNegotiationResult</a> instance which carries the chosen proper <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter(v=vs.108).aspx" target="_blank">MediaTypeFormatter</a> and the <a title="http://msdn.microsoft.com/en-us/library/system.net.http.headers.mediatypeheadervalue(v=vs.108)" href="http://msdn.microsoft.com/en-us/library/system.net.http.headers.mediatypeheadervalue(v=vs.108)">MediaTypeHeaderValue</a>. In the above sample, I don&rsquo;t do anything with it and it doesn&rsquo;t effect the behavior in any way but you might have an edge case where you might need something like this.</p>