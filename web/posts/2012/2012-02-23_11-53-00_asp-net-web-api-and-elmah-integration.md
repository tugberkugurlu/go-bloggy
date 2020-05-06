---
title: ASP.NET Web API and ELMAH Integration
abstract: See how ASP.NET Web API Plays Nice With ELMAH. This blog post is a Quick
  introduction to ASP.NET Web API and System.Web.Http.Filters.IExceptionFilter
created_at: 2012-02-23 11:53:00 +0000 UTC
tags:
- ASP.NET Web API
slugs:
- asp-net-web-api-and-elmah-integration
---

<p>As you all probably heard, <a target="_blank" href="http://asp.net/mvc/mvc4">ASP.NET MVC 4</a> Beta is now available and has new features in it. One of them is that ASP.NET MVC has shipped with <a target="_blank" href="http://www.asp.net/web-api" title="http://www.asp.net/web-api">ASP.NET Web API</a> (which was previously know as <a target="_blank" href="http://www.tugberkugurlu.com/archive/introduction-to-wcf-web-api-new-rest-face-ofnet" title="http://www.tugberkugurlu.com/archive/introduction-to-wcf-web-api-new-rest-face-ofnet">WCF Web API</a>). Here is the quote from ASP.NET web site which explains what Web API Framework is all about shortly:</p>
<p><em>ASP.NET Web API is a framework that makes it easy to build HTTP services that reach a broad range of clients, including browsers and mobile devices. ASP.NET Web API is an ideal platform for building RESTful applications on the .NET Framework.</em></p>
<p>I am not going to give an into on ASP.NET Web API. There are great into tutorials and videos on ASP.NET web site for ASP.NET Web API. Instead, I will give you an example of a custom filter implementation.</p>
<p>A couple of months ago, I wrote a blog post about <em>WCF Web API HttpErrorHandlers</em>: <a target="_blank" href="http://www.tugberkugurlu.com/archive/wcf-web-api-plays-nice-with-elmah-a-quick-introduction-to-wcf-web-api-httperrorhandler" title="http://www.tugberkugurlu.com/archive/wcf-web-api-plays-nice-with-elmah-a-quick-introduction-to-wcf-web-api-httperrorhandler">WCF Web API Plays Nice With ELMAH - A Quick Introduction to WCF Web API HttpErrorHandler</a>. It works nicely on WCF Web API but this version of the framework, things a little bit changed. Instead of ErrorHandlers, we now have Filters which is more generic.</p>
<p>In order to make ELMAH work with ASP.NET Web API, we need to create a new Attribute which implements IExceptionFilter interface. Since we have ExceptionFilterAttribute (which implements the IExceptionFilter interface) available at the framework, we will derived from that class instead. Here is the whole implementation:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> ElmahErrorAttribute : 
    System.Web.Http.Filters.ExceptionFilterAttribute {

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> <span style="color: blue;">void</span> OnException(
        System.Web.Http.Filters.HttpActionExecutedContext actionExecutedContext) {

        <span style="color: blue;">if</span>(actionExecutedContext.Exception != <span style="color: blue;">null</span>)
            Elmah.ErrorSignal.FromCurrentContext().Raise(actionExecutedContext.Exception);

        <span style="color: blue;">base</span>.OnException(actionExecutedContext);
    }
}</pre>
</div>
</div>
<p>Now we have our attribute, we need to tell the framework to use it. It is very straight forward as well. Here is how my Global.asax (Global.asax.cs) looks like:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> WebApiApplication : System.Web.HttpApplication {

    <span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start() {

        Configure(
            System.Web.Http.GlobalConfiguration.Configuration
        );
    }

    <span style="color: blue;">private</span> <span style="color: blue;">void</span> Configure(HttpConfiguration httpConfiguration) {

        httpConfiguration.Filters.Add(
            <span style="color: blue;">new</span> ElmahErrorAttribute()
        );

        httpConfiguration.Routes.MapHttpRoute(
            name: <span style="color: #a31515;">"DefaultApi"</span>,
            routeTemplate: <span style="color: #a31515;">"api/{controller}/{id}"</span>,
            defaults: <span style="color: blue;">new</span> { id = RouteParameter.Optional }
        );
    }
}</pre>
</div>
</div>
<p>This will make Elmah log the errors. Now, you can configure ELMAH to send you an e-mail when an error occurred or you can log the error inside an XML file, SQL Server Database, wherever you what.</p>
<p>One more thing to mention about is that when you hit an error, the response will carry the exception details at the body of the response if you run your application locally. You can configure this option as well with the following configuration:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>httpConfiguration.IncludeErrorDetailPolicy = IncludeErrorDetailPolicy.Never;</pre>
</div>
</div>