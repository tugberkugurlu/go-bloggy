---
id: 94769503-cb58-4ff9-b446-dd769e47f8f3
title: Introduction to WCF Web API - New REST Face of .NET
abstract: This blog post will give you an introduction to WCF Web API and show you
  how to get started with WCF Web API along with Dependency Inject support with Ninject.
created_at: 2011-11-21 08:09:00 +0000 UTC
tags:
- .NET
- C#
- WCF Web API
slugs:
- introduction-to-wcf-web-api-new-rest-face-ofnet
---

<blockquote>
<p><strong>24 February 2012</strong></p>
WCF Web API is now <a href="http://asp.net/web-api" title="http://asp.net/web-api">ASP.NET Web API</a>&nbsp;and has changed a lot. The beta version is now available.&nbsp;For more information: <a href="https://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples" title="https://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples">Getting Started With ASP.NET Web API - Tutorials, Videos, Samples</a>.</blockquote>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image.png"><img style="background-image: none; margin: 0px 15px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" align="left" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb.png" width="244" height="166" /></a>Microsoft Web Platform is evolving. I mean really evolving. <a title="http://hanselman.com/" href="http://hanselman.com/" target="_blank">Scott Hanselman</a>, <a title="http://haacked.com" href="http://haacked.com" target="_blank">Phil Haack</a>, <a title="http://blogs.msdn.com/gblock" href="http://blogs.msdn.com/gblock" target="_blank">Glenn Block</a> and <a title="http://damianedwards.wordpress.com" href="http://damianedwards.wordpress.com" target="_blank">Damian Edwards</a> are the main actors for this evolution.</p>
<p>One of the biggest frustration we had as web developers was to face with the endless configurations with WCF. I mean, WCF is great but hasn&rsquo;t been embraced the REST since the <a title="http://wcf.codeplex.com/wikipage?title=WCF%20HTTP" href="http://wcf.codeplex.com/wikipage?title=WCF%20HTTP" target="_blank">WCF Web API</a> framework. There was something called REST Starter Kit but it ended up dead.</p>
<p>Yes, WCF Web API. WCF Web API is the new way of exposing and consuming APIs over http in .NET. The unofficial (not sure, maybe the official one) slogan of WCF Web API is this:</p>
<p><em>"Making REST a first class citizen in .NET"</em></p>
<p>I stole (or quoted would be nicer) the real sentences which explain the WCF Web API from the official WCF Web API page (I hope they don&rsquo;t mind):</p>
<blockquote>
<h4>What is WCF Web API?</h4>
<p><em>Applications are continually evolving to expose their functionality over the web for example social services like </em><a href="http://www.flickr.com/"><em>Flickr</em></a><em>, </em><a href="http://www.twitter.com/"><em>Twitter</em></a><em> and </em><a href="http://www.facebook.com/"><em>Facebook</em></a><em>. Aside from social applications, organizations are also looking to surface their core enterprise business functionality to an ever expanding array of client platforms. WCF Web API allows developers to expose their applications, data and services to the web directly over </em><a href="http://en.wikipedia.org/wiki/HTTP"><em>HTTP</em></a><em>. This allows developers to fully harness the richness of the HTTP as an application layer protocol. Applications can communicate with a very broad set of clients whether they be browsers, mobile devices, desktop applications or other backend services. They can also take advantage of the caching and proxy infrastructure of the web through providing proper control and entity headers. We are designing specifically to support applications built with a </em><a href="http://en.wikipedia.org/wiki/Representational_State_Transfer"><em>RESTful</em></a><em> architecture style though it does not force developers to use REST. The benefits of REST for your applications include discoverability, evolvability and scalability.</em></p>
</blockquote>
<p>The project is still at the preview stage and we are swimming in the dark sea. There are lots of rumors going around about it and most of them are nearly certain to be true. One of them is that WCF Web API and <a title="http://asp.net/mvc" href="http://asp.net/mvc" target="_blank">ASP.NEt MVC</a> will be blood brothers in near future. In plain English, they will be merged together. This will be exciting and we won&rsquo;t feel ourselves in a fork in the road when we need to pick ASP.NET MVC or WCF Web API in order to expose our data over http.</p>
<p><strong>Cut the crap and show me the code</strong></p>
<p>Well, when I first saw the WCF Web API, I told myself that this&rsquo;s it! Why I told that? Because it is extremely easy to get started and going from there. I am interested in WCF as well but its configuration is endless so I haven&rsquo;t been able to develop a decent project with WCF so far (maybe the problem is me, who knows!). Be careful here though, WCF is not gone! It is still the way to go with for SOAP based services.</p>
<p>Let&rsquo;s see how we get started developing a <strong>Web API</strong> with <strong>WCF Web API </strong>(this sentence is like a poem <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/wlEmoticon-smile.png" />).</p>
<p>First of all, open up your VS and create a new <strong>ASP.NET Empty Web Application</strong> (this one is real <strong>empty</strong> guys unlike <strong>ASP.NET MVC 3 Empty Web Application</strong>):</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_3.png"><img style="background-image: none; margin: 0px 15px 0px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_3.png" width="244" height="159" /></a><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_4.png"><img style="background-image: none; margin: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_4.png" width="244" height="160" /></a></p>
<p>I told you it is real empty. Anyway, let&rsquo;s stick to the point here<img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/wlEmoticon-smile.png" /> Now, bring up the <a title="http://nuget.org" href="http://nuget.org" target="_blank">NuGet</a> PMC (<strong>P</strong>ackage <strong>M</strong>anager <strong>C</strong>onsole) and install the package called <strong>WebApi.All</strong>:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_5.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_5.png" width="640" height="75" /></a></p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_6.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_6.png" width="644" height="238" /></a></p>
<p>Let&rsquo;s look what we got:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_7.png" width="644" height="187" /></a></p>
<p>All the packages we have pulled is individual packages. For example, if you are consuming web services on the server side, <strong>HttpClient</strong> package will help you a lot. I think all of those packages will be baked-in for .NET 4.5 so we will see a lot samples for those.</p>
<p>One thing that I looked when I first bring up this package is Web.Config file because I wondered how giant it was. It comes to me as a shock and it will for you, too:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_8.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_8.png" width="644" height="229" /></a>&nbsp;</p>
<p>3 lines of code which is special to WCF Web API. This is awesome. So, where do we configure the stuff. First of all, if you would like to get started you do not need to make any configuration. We will see how in a minute. Web API comes with default configuration and this can be overridden in any steps. You can set your default configuration. One of your APIs needs different configuration? Don&rsquo;t change the default one. Configure it separately. So, it is really a <strong>convention</strong>. Best part is that it is all done with code, I mean inside <strong>Global.asax</strong>.</p>
<p>For the sake of this demo, I created a dummy data to play with. It is really simple as follows (I will put up the source code online, you can check out what is in there):</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_9.png"><img style="background-image: none; margin: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_9.png" width="220" height="244" /></a></p>
<p>I implemented the repository pattern here with an interface. I did that because I would like to show you how easy is to get going along with Dependency Injection (DI) here.</p>
<p>In order to create our API, we need to create a separate class. I put it under <strong>People</strong> folder and named it <strong>PeopleApi </strong>but where it stands and what name it carries don&rsquo;t matter here.</p>
<p>I would like to go with the simplest approach firstly. Here is how <strong>PeopleApi</strong> class looks like:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">using</span> System;
<span style="color: blue;">using</span> System.Collections.Generic;
<span style="color: blue;">using</span> System.Linq;
<span style="color: blue;">using</span> System.Net.Http;
<span style="color: blue;">using</span> System.ServiceModel;
<span style="color: blue;">using</span> System.ServiceModel.Web;
<span style="color: blue;">using</span> System.Web;
<span style="color: blue;">using</span> VeryFirstWcfWebAPI.People.Infrastructure;
<span style="color: blue;">using</span> VeryFirstWcfWebAPI.People.Models;

<span style="color: blue;">namespace</span> VeryFirstWcfWebAPI.People {

    [ServiceContract]
    <span style="color: blue;">public</span> <span style="color: blue;">class</span> PeopleApi {

        <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IPeopleRepository _repo = <span style="color: blue;">new</span> PeopleRepository();

        [WebGet]
        <span style="color: blue;">public</span> HttpResponseMessage&lt;IQueryable&lt;Person&gt;&gt; Get() {

            <span style="color: blue;">var</span> model = _repo.GetAll();

            <span style="color: blue;">var</span> responseMessage = <span style="color: blue;">new</span> HttpResponseMessage&lt;IQueryable&lt;Person&gt;&gt;(model);
            responseMessage.Content.Headers.Expires = <span style="color: blue;">new</span> DateTimeOffset(DateTime.Now.AddHours(6));

            <span style="color: blue;">return</span> responseMessage;
        }

    }
}</pre>
</div>
</div>
<p>When we observe this class a little bit carefully, we see some staff going on there:</p>
<ul>
<li>You API class needs to be annotated with <a title="http://msdn.microsoft.com/en-us/library/system.servicemodel.servicecontractattribute.aspx" href="http://msdn.microsoft.com/en-us/library/system.servicemodel.servicecontractattribute.aspx" target="_blank">ServiceContractAttribute</a>. This is must to do (<a title="http://twitter.com/gblock" href="http://twitter.com/gblock" target="_blank">@gblock</a> said at the //Build conference that we are still in love with attributes but we are trying to get rid of them...).</li>
<li>The methods inside your class needs to have some special attributes like <a title="http://msdn.microsoft.com/en-us/library/system.servicemodel.web.webgetattribute.aspx" href="http://msdn.microsoft.com/en-us/library/system.servicemodel.web.webgetattribute.aspx" target="_blank">WebGetAttribute</a> and <a title="http://msdn.microsoft.com/en-us/library/system.servicemodel.web.webinvokeattribute.aspx" href="http://msdn.microsoft.com/en-us/library/system.servicemodel.web.webinvokeattribute.aspx" target="_blank">WebInvokeAttribute</a>. If you put them without <a title="http://msdn.microsoft.com/en-us/library/system.servicemodel.web.webgetattribute.uritemplate.aspx" href="http://msdn.microsoft.com/en-us/library/system.servicemodel.web.webgetattribute.uritemplate.aspx" target="_blank">UriTemplate</a> property, it assumes that the method is for root of the URL.</li>
</ul>
<p>Also, another thing to notice here is we are returning our model by wrapping it up with <a title="http://msdn.microsoft.com/en-us/library/system.net.http.httpresponsemessage(v=vs.110).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.httpresponsemessage(v=vs.110).aspx" target="_blank">HttpResponseMessage</a> class. You don&rsquo;t have to do that. You can just return your object but if you need to add some special headers or response message code, it is nice way to do it that way as we set here our expires header.</p>
<p>As I mentioned before, there is no configuration at all in the web.config but we still need to do some configuration to tell the system to figure out what to do.</p>
<p>Add a Global Application Class under the root of your project:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_10.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_10.png" width="644" height="418" /></a></p>
<p>Inside the Application_Start method, here is our initial set up to get going:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">using</span> System;
<span style="color: blue;">using</span> System.Collections.Generic;
<span style="color: blue;">using</span> System.Linq;
<span style="color: blue;">using</span> System.Web;
<span style="color: blue;">using</span> System.Web.Routing;
<span style="color: blue;">using</span> System.Web.Security;
<span style="color: blue;">using</span> System.Web.SessionState;

<span style="color: blue;">namespace</span> VeryFirstWcfWebAPI {

    <span style="color: blue;">public</span> <span style="color: blue;">class</span> Global : System.Web.HttpApplication {

        <span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

            RouteTable.Routes.SetDefaultHttpConfiguration(<span style="color: blue;">new</span> Microsoft.ApplicationServer.Http.WebApiConfiguration() { 
            });

            RouteTable.Routes.MapServiceRoute&lt;People.PeopleApi&gt;(<span style="color: #a31515;">"Api/People"</span>);
        }

    }
}</pre>
</div>
</div>
<p>What we do here is so simple:</p>
<ul>
<li>We are initializing the WCF Web API with default configuration.</li>
<li>We specifically register our API with a base URL structure as route prefix.</li>
</ul>
<blockquote>
<p>As we have folder called <strong>People</strong> under the root of out application, if you put <strong>People</strong> as route prefix there, you will see that your API won&rsquo;t work and you will get 404. I am nearly sure that it is related to routing.</p>
<p>I haven&rsquo;t figured out how to solve this issue and I tried to Ignore that folder but it didn&rsquo;t work either.</p>
</blockquote>
<p>We are ready to run:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_11.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_11.png" width="644" height="433" /></a></p>
<p>This is good for a start. We didn&rsquo;t suffer much. Now, it is time for us to think of the possible enhancements.</p>
<p><strong>Dependency Injection and IoC Container Support</strong></p>
<p>As you probably noticed, I "newed" up the repository inside our API class:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IPeopleRepository _repo = <span style="color: blue;">new</span> PeopleRepository();</pre>
</div>
</div>
<p>This thing made my application tightly-coupled and it is not good. Here I have a static data resource and it is not much of a problem but if we had database related data structure here, this would make it hard for us to unit test our application.</p>
<p>In order to get around for this, we need to figure out a way to new up the resource outside of our context and WCF Web API offers really good extensibility point here. I won&rsquo;t really extend and customize our configuration much here in order to stick with the basics but I will probably blog about that either.</p>
<p>First of all, let&rsquo;s make our API class a little DI friendly:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>    [ServiceContract]
    <span style="color: blue;">public</span> <span style="color: blue;">class</span> PeopleApi {

        <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IPeopleRepository _repo;

        <span style="color: blue;">public</span> PeopleApi(IPeopleRepository repo) {
            _repo = repo;
        }

        [WebGet]
        <span style="color: blue;">public</span> HttpResponseMessage&lt;IQueryable&lt;Person&gt;&gt; Get() {

            <span style="color: blue;">var</span> model = _repo.GetAll();

            <span style="color: blue;">var</span> responseMessage = <span style="color: blue;">new</span> HttpResponseMessage&lt;IQueryable&lt;Person&gt;&gt;(model);
            responseMessage.Content.Headers.Expires = <span style="color: blue;">new</span> DateTimeOffset(DateTime.Now.AddHours(6));

            <span style="color: blue;">return</span> responseMessage;
        }

    }</pre>
</div>
</div>
<p>What we do here enables someone else to provide the resource so our API class is now loosely-coupled. I would like to see what happens when we run the app like this:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_12.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_12.png" width="644" height="433" /></a></p>
<p>We got an error:</p>
<blockquote>
<p><span color="#ff0000" style="color: #ff0000;"><strong>The service type provided could not be loaded as a service because it does not have a default (parameter-less) constructor. To fix the problem, add a default constructor to the type, or pass an instance of the type to the host.</strong></span></p>
</blockquote>
<p>The system needs a parameter-less constructor as default. Let&rsquo;s see how we get around with this.</p>
<p>Now, we need to provide those resources outside of our context but where? Let&rsquo;s first bring down a IoC container. I am fan of Ninject so I will use it here as well:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_13.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_13.png" width="644" height="191" /></a></p>
<p>After we install Ninject through NuGet, I added the following code inside my Global.asax file:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">private</span> IKernel GetKernel() { 
    
    IKernel kernel = <span style="color: blue;">new</span> StandardKernel();

    kernel.Bind&lt;People.Infrastructure.IPeopleRepository&gt;().
        To&lt;People.Models.PeopleRepository&gt;();

    <span style="color: blue;">return</span> kernel;
}</pre>
</div>
</div>
<p>This will provide us the resources we need. Now, we need to tell WCF Web API to use this provider to create instances. Believe it or not, it is extremely easy to do that. Remember our default configuration object, WebApiConfiguration class? We will register our IoC container there:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>RouteTable.Routes.SetDefaultHttpConfiguration(<span style="color: blue;">new</span> Microsoft.ApplicationServer.Http.WebApiConfiguration() { 
    CreateInstance = (serviceType, context, request) =&gt; GetKernel().Get(serviceType)
});</pre>
</div>
</div>
<p>We are passing a delegate here for CreateInstance property of our WebApiConfiguration class. When we run the application, we should see it working:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_14.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/image_thumb_14.png" width="644" height="433" /></a></p>
<p>So nice to do something like that without much effort.</p>
<p>There is so much to show but I think this is enough for an intro (which I write on the stage of preview 5, I am little late <img style="border-style: none;" class="wlEmoticon wlEmoticon-confusedsmile" alt="Confused smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ec9739380e27_9F4D/wlEmoticon-confusedsmile.png" />). I am sure that you get the idea here.</p>
<p>You can find the full code on <a title="https://github.com" href="https://github.com" target="_blank">GitHub</a>: <a href="https://github.com/tugberkugurlu/VeryFirstWcfWebAPI" target="_blank">https://github.com/tugberkugurlu/VeryFirstWcfWebAPI</a></p>