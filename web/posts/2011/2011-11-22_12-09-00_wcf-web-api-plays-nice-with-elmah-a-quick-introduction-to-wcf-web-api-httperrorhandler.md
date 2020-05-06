---
title: WCF Web API Plays Nice With ELMAH - A Quick Introduction to WCF Web API HttpErrorHandler
abstract: See how WCF Web API Plays Nice With ELMAH. This blog post is a Quick introduction
  to WCF Web API HttpErrorHandler
created_at: 2011-11-22 12:09:00 +0000 UTC
tags:
- .net
- ASP.Net
- WCF Web API
slugs:
- wcf-web-api-plays-nice-with-elmah-a-quick-introduction-to-wcf-web-api-httperrorhandler
---

<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image.png"><img height="165" width="244" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb.png" align="right" alt="image" border="0" title="image" style="background-image: none; margin: 0px 0px 10px 15px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border: 0px;" /></a>How many times you face an error message like this one when you are working with WCF? Many times I guess.</p>
<p>On WCF Web API, this is not a big deal on the development stage because Web API runs on core ASP.NET and debugging is not a big deal (maybe it is not a big deal on WCF as well but it always for me, anybody knows a good debugging scenario on WCF, please let me know).</p>
<p>But when you expose your data and your customers starts to consume your service, you will pull your hairs when you see this screen.</p>
<p>WCF Web API has been built extensibility and testability in mind so it is real easy to plug things into the framework. One of the extensibility points is <strong>ErrorHandlers</strong> and in this quick blog post I will show you how to handle error nicely on WCF Web API.</p>
<blockquote>
<p>If you haven&rsquo;t seen my previous blog post on <a target="_blank" href="http://www.tugberkugurlu.com/archive/introduction-to-wcf-web-api-new-rest-face-ofnet" title="http://www.tugberkugurlu.com/archive/introduction-to-wcf-web-api-new-rest-face-ofnet">Introduction to WCF Web API - New REST Face of .NET</a>, I encourage you to check that out. What I do here will be addition to that.</p>
</blockquote>
<p>When somebody tells me the words &ldquo;Error Handling&rdquo; and &ldquo;.NET&rdquo; in the same sentence, I tell him/her <a target="_blank" href="http://code.google.com/p/elmah/" title="http://code.google.com/p/elmah/">ELMAH</a> in response. <em>ELMAH (Error Logging Modules and Handlers) is is an application-wide error logging facility that is completely pluggable. It can be dynamically added to a running </em><a href="http://www.asp.net/"><em>ASP.NET</em></a><em> web application, or even all ASP.NET web applications on a machine, without any need for re-compilation or re-deployment. </em>These are the official words.</p>
<p>In order to integrate ELMAH with our Web API application, we need to bring down the ELMAH via <a target="_blank" href="http://nuget.org" title="http://nuget.org">NuGet</a> as we always do for open source libraries.</p>
<blockquote>
<p><a target="_blank" href="http://hanselman.com/" title="http://hanselman.com/">Scott Hanselman</a> has a great blog post on how ELMAH gets into your application via NuGet Package Manager. You should check that out if you are interested:</p>
<p><a href="http://www.hanselman.com/blog/IntroducingNuGetPackageManagementForNETAnotherPieceOfTheWebStack.aspx">http://www.hanselman.com/blog/IntroducingNuGetPackageManagementForNETAnotherPieceOfTheWebStack.aspx</a></p>
</blockquote>
<p>Open up you PMC and type <em><strong>Install-Package ELMAH</strong></em>:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_3.png"><img height="253" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb_3.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>After you install the package successfully, it doesn&rsquo;t need any extra configuration to work but you would probably want to secure your error logging page. Check it out how on <a target="_blank" href="http://haacked.com" title="http://haacked.com">Phil Haack</a>&rsquo;s blog post: <a href="http://haacked.com/archive/2007/07/24/securely-implement-elmah-for-plug-and-play-error-logging.aspx">http://haacked.com/archive/2007/07/24/securely-implement-elmah-for-plug-and-play-error-logging.aspx</a></p>
<p>When you go to http://localhost:{port_number}/elmah.axd, you will see the error list page:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_4.png"><img height="433" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb_4.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>I added a new method for PeopleApi in order to be able to reach single person data:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[WebGet(UriTemplate = <span style="color: #a31515;">"{id}"</span>)]
<span style="color: blue;">public</span> HttpResponseMessage&lt;Person&gt; GetSingle(<span style="color: blue;">int</span> id) {

    <span style="color: blue;">var</span> person = _repo.GetSingle(id);

    <span style="color: blue;">if</span> (person == <span style="color: blue;">null</span>) {
        <span style="color: blue;">var</span> response = <span style="color: blue;">new</span> HttpResponseMessage();
        response.StatusCode = HttpStatusCode.NotFound;
        response.Content = <span style="color: blue;">new</span> StringContent(<span style="color: #a31515;">"Country not found"</span>);
        <span style="color: blue;">throw</span> <span style="color: blue;">new</span> HttpResponseException(response);
    }

    <span style="color: blue;">var</span> personResponse = <span style="color: blue;">new</span> HttpResponseMessage&lt;Models.Person&gt;(person);
    personResponse.Content.Headers.Expires = <span style="color: blue;">new</span> DateTimeOffset(DateTime.Now.AddHours(6));
    <span style="color: blue;">return</span> personResponse;
}</pre>
</div>
</div>
<p>It is a simple method. It returns a single Person data wrapped up inside <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.net.http.httpresponsemessage(v=vs.110).aspx" title="http://msdn.microsoft.com/en-us/library/system.net.http.httpresponsemessage(v=vs.110).aspx">HttpResponseMessage</a> if there is one and returns 404 if there is no person with the given id value.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_5.png"><img height="165" width="244" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb_5.png" alt="image" border="0" title="image" style="background-image: none; margin: 0px 15px 0px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_6.png"><img height="165" width="244" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb_6.png" alt="image" border="0" title="image" style="background-image: none; margin: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Let&rsquo;s send a string value instead of int32 and see what happens:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_7.png"><img height="433" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb_7.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>It is an error so ELMAH should have logged this, right? It didn&rsquo;t because WCF Web API handles exceptions on its own but good news is you can get in there and plug your own stuff into it.</p>
<blockquote>
<p><strong>Disclaimer: </strong></p>
<p>I am still a newbie on WCF Web API and learning it day by day. Also, the framework is at the preview stage (not even Alpha) so the things I explain and show might be not entirely the best case scenarios.</p>
</blockquote>
<p>In WCF Web API preview 5, you can use the <strong>ErrorHandler</strong>, which is the recommended way to do error handling. In order to implement your own error handler, your class needs to derived from HttpErrorHandler class. Here is the implementation I use:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">using</span> System;
<span style="color: blue;">using</span> System.Collections.Generic;
<span style="color: blue;">using</span> System.Linq;
<span style="color: blue;">using</span> System.Net;
<span style="color: blue;">using</span> System.Net.Http;
<span style="color: blue;">using</span> System.Web;
<span style="color: blue;">using</span> Microsoft.ApplicationServer.Http.Dispatcher;

<span style="color: blue;">namespace</span> VeryFirstWcfWebAPI.Handlers {

    <span style="color: blue;">public</span> <span style="color: blue;">class</span> GlobalErrorHandler : HttpErrorHandler {

        <span style="color: blue;">protected</span> <span style="color: blue;">override</span> <span style="color: blue;">bool</span> OnTryProvideResponse(Exception exception, <span style="color: blue;">ref</span> System.Net.Http.HttpResponseMessage message) {

            <span style="color: blue;">if</span>(exception != <span style="color: blue;">null</span>) <span style="color: green;">// Notify ELMAH</span>
                Elmah.ErrorSignal.FromCurrentContext().Raise(exception);

            message = <span style="color: blue;">new</span> HttpResponseMessage {
                StatusCode = HttpStatusCode.InternalServerError
            };

            <span style="color: blue;">return</span> <span style="color: blue;">true</span>;
        }
    }
}</pre>
</div>
</div>
<p>Pretty straight forward implementation. The last step is to register this handler. We will do that inside the <strong>Global.asax</strong> file as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    RouteTable.Routes.SetDefaultHttpConfiguration(<span style="color: blue;">new</span> Microsoft.ApplicationServer.Http.WebApiConfiguration() { 
        CreateInstance = (serviceType, context, request) =&gt; GetKernel().Get(serviceType),
        ErrorHandlers = (handlers, endpoint, description) =&gt; handlers.Add(<span style="color: blue;">new</span> GlobalErrorHandler())
    });

    RouteTable.Routes.MapServiceRoute&lt;People.PeopleApi&gt;(<span style="color: #a31515;">"Api/People"</span>);
}</pre>
</div>
</div>
<p>Now, hit your service and cause an error. After that go back to elmah.axd page to see what we got there:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_8.png"><img height="433" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb_8.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_9.png"><img height="433" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/c64ad9c94e51_D270/image_thumb_9.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>We totally nailed it! Now, you can configure ELMAH to send you an e-mail when an error occurred or you can log the error inside an XML file, SQL Server Database, wherever you what.</p>
<p>As I said at the end of my previous post, there is so much to cover about WCF Web API. I hope I will blog about it more.</p>
<blockquote>
<p>You can find the full code on <a title="https://github.com" href="https://github.com" target="_blank">GitHub</a>: <a href="https://github.com/tugberkugurlu/VeryFirstWcfWebAPI" target="_blank">https://github.com/tugberkugurlu/VeryFirstWcfWebAPI</a></p>
</blockquote>