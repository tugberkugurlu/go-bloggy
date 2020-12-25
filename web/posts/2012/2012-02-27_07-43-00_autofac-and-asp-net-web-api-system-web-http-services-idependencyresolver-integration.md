---
id: 6af246b5-048d-4d47-85ed-0f975ea1aa9b
title: Autofac and ASP.NET Web API System.Web.Http.Services.IDependencyResolver Integration
abstract: In this post, you can make Autofac work with ASP.NET Web API System.Web.Http.Services.IDependencyResolver.
  Solution to the 'controller has no parameterless public constructor' error.
created_at: 2012-02-27 07:43:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
- Unit Testing
slugs:
- autofac-and-asp-net-web-api-system-web-http-services-idependencyresolver-integration
---

<p>I have been using <a target="_blank" href="http://www.ninject.org/" title="http://www.ninject.org/">Ninject</a> in all of my ASP.NET MVC applications till <a target="_blank" href="http://dotnetchris.wordpress.com/" title="http://dotnetchris.wordpress.com/">Chris Marisic</a> (<a target="_blank" href="https://twitter.com/dotnetchris" title="https://twitter.com/dotnetchris">@dotnetchris</a>) poked me on twitter about how slow Ninject is. After a little bit of research, I changed it to <a target="_blank" href="http://code.google.com/p/autofac/" title="http://code.google.com/p/autofac/">Autofac</a> which works pretty perfectly.</p>
<p><a target="_blank" href="http://www.asp.net/web-api" title="https://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples">ASP.NET Web API</a> has nearly the same Dependency Injection support as ASP.NET MVC. As we do not have a built in support for ASP.NET Web API in Autofac (yet), I created a simple one. The implementation is not as straight forward as Ninject and you probably saw the below error if you tried to make it work:</p>
<blockquote>
<p><em>System.InvalidOperationException:</em></p>
<p><em>An error occurred when trying to create a controller of type '<strong>TourismDictionary.APIs.Controllers.WordsController</strong>'. Make sure that the controller has a parameterless public constructor.</em></p>
</blockquote>
<p>Let&rsquo;s see how it works:</p>
<p>First I created a class which implements <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.http.services.idependencyresolver(v=vs.108).aspx" title="http://msdn.microsoft.com/en-us/library/system.web.http.services.idependencyresolver(v=vs.108).aspx">System.Web.Http.Services.IDependencyResolver</a> interface.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">internal</span> <span style="color: blue;">class</span> AutofacWebAPIDependencyResolver : 
    System.Web.Http.Services.IDependencyResolver {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IContainer _container;

    <span style="color: blue;">public</span> AutofacWebAPIDependencyResolver(IContainer container) {

        _container = container;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">object</span> GetService(Type serviceType) {

        <span style="color: blue;">return</span> 
            _container.IsRegistered(serviceType) ? 
            _container.Resolve(serviceType) : <span style="color: blue;">null</span>;
    }

    <span style="color: blue;">public</span> IEnumerable&lt;<span style="color: blue;">object</span>&gt; GetServices(Type serviceType) {

        Type enumerableServiceType = 
            <span style="color: blue;">typeof</span>(IEnumerable&lt;&gt;).MakeGenericType(serviceType);
            
        <span style="color: blue;">object</span> instance = 
            _container.Resolve(enumerableServiceType);
            
        <span style="color: blue;">return</span> ((IEnumerable)instance).Cast&lt;<span style="color: blue;">object</span>&gt;();
    }
}</pre>
</div>
</div>
<p>And I have another class which holds my registrations and sets the ServiceResolver with <strong>GlobalConfiguration.Configuration.ServiceResolver.SetResolver </strong>method which is equivalent to <a target="_blank" href="http://msdn.microsoft.com/en-us/library/hh835218(v=vs.108).aspx" title="http://msdn.microsoft.com/en-us/library/hh835218(v=vs.108).aspx">DependencyResolver.SetResolver</a> method:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">internal</span> <span style="color: blue;">class</span> AutofacWebAPI {

    <span style="color: blue;">public</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> Initialize() {
        <span style="color: blue;">var</span> builder = <span style="color: blue;">new</span> ContainerBuilder();
        GlobalConfiguration.Configuration.ServiceResolver.SetResolver(
            <span style="color: blue;">new</span> AutofacWebAPIDependencyResolver(RegisterServices(builder))
        );
    }

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> IContainer RegisterServices(ContainerBuilder builder) {

        builder.RegisterAssemblyTypes(
            <span style="color: blue;">typeof</span>(MvcApplication).Assembly
        ).PropertiesAutowired();

        <span style="color: green;">//deal with your dependencies here</span>
        builder.RegisterType&lt;WordRepository&gt;().As&lt;IWordRepository&gt;();
        builder.RegisterType&lt;MeaningRepository&gt;().As&lt;IMeaningRepository&gt;();

        <span style="color: blue;">return</span>
            builder.Build();
    }
}</pre>
</div>
</div>
<p>Then, initialize it at <strong>Application_Start</strong>:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> MvcApplication : System.Web.HttpApplication {

    <span style="color: blue;">private</span> <span style="color: blue;">void</span> Configure(HttpConfiguration httpConfiguration) {

        httpConfiguration.Routes.MapHttpRoute(
            name: <span style="color: #a31515;">"DefaultApi"</span>,
            routeTemplate: <span style="color: #a31515;">"api/{controller}/{id}"</span>,
            defaults: <span style="color: blue;">new</span> { id = RouteParameter.Optional }
        );
    }

    <span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start() {

        Configure(GlobalConfiguration.Configuration);
        AutofacWebAPI.Initialize();
    }

}</pre>
</div>
</div>
<p>There are different ways of doing it but I like this approach.</p>