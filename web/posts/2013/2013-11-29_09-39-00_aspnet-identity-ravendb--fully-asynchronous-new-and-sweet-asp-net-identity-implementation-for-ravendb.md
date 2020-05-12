---
id: abeba914-b4ba-4bb8-ab4e-ab8af2635bd6
title: 'AspNet.Identity.RavenDB: Fully asynchronous, new and sweet ASP.NET Identity
  implementation for RavenDB'
abstract: 'A while back, ASP.NET team has introduced ASP.NET Identity, a membership
  system for ASP.NET applications. Today, I''m introducing you its RavenDB implementation:
  AspNet.Identity.RavenDB.'
created_at: 2013-11-29 09:39:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET MVC
- ASP.NET Web API
- RavenDB
- SignalR
slugs:
- aspnet-identity-ravendb--fully-asynchronous-new-and-sweet-asp-net-identity-implementation-for-ravendb
---

<p>A while back, ASP.NET team has <a href="http://blogs.msdn.com/b/webdev/archive/2013/06/27/introducing-asp-net-identity-membership-system-for-asp-net-applications.aspx">introduced ASP.NET Identity</a>, a membership system for ASP.NET applications. If you have Visual Studio 2013 installed on you box, you will see <a href="http://www.nuget.org/packages/Microsoft.AspNet.Identity.Core/">ASP.NET Identity Core</a> and <a href="http://www.nuget.org/packages/Microsoft.AspNet.Identity.EntityFramework/">ASP.NET Identity Entity Framework</a> libraries are pulled down when you create a new Web Application by configuring the authentication as "Individual User Accounts".</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP_955C/image.png"><img height="455" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP_955C/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="image" /></a></p>
<p>After creating your MVC project, you will see that you have an AccountController which a completely different code from the previous project templates as it uses ASP.NET Identity.</p>
<blockquote>
<p>You can find tons of information about this new membership system from <a href="http://www.asp.net/identity/overview/getting-started">ASP.NET Identity section on official ASP.NET web site</a>. Also, <a href="https://twitter.com/rustd">Pranav Rastogi</a> (a.k.a @rustd) has a great <a href="http://channel9.msdn.com/Shows/Web+Camps+TV/Special-Movember-Episode-ASPNET-Authentication-Provider">introduction video on ASP.NET Identity</a> which you shouldn't miss for sure.</p>
</blockquote>
<p>One of the great features of ASP.NET Identity system is the fact that it is super extensible. The core layer and the implementation layer (which is Entity Framework by default) are decouple from each other. This means that you are not bound to Entity Framework and SQL Server for storage. You can implement ASP.NET Identity on your choice of storage system. This is exactly what I did and I created <a href="https://github.com/tugberkugurlu/AspNet.Identity.RavenDB">AspNet.Identity.RavenDB</a> project which is fully asynchronous, new and sweet ASP.NET Identity implementation for <a href="http://ravendb.net/">RavenDB</a>. You can install this library from NuGet:</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package AspNet.Identity.RavenDB</code></p>
</div>
<p>Getting started with AspNet.Identity.RavenDB is also really easy. Just create an <a href="http://aspnetwebstack.codeplex.com">ASP.NET MVC</a> application from scratch by configuring the authentication as "Individual User Accounts". Then, install the AspNet.Identity.RavenDB package. As the default project is set to work with ASP.NET Identity Entity Framework implementation, you need to make a few more tweak here and there to make it work with RavenDB.</p>
<p>First, open the IdentityModels.cs file under the "Models" folder and delete the two classes you see there. Only class you need is the following ApplicationUser class:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> ApplicationUser : RavenUser
{
}</pre>
</div>
</div>
<p>As the second step, open up the AccountController.cs file under the "Controllers" folder and have delete the first constructor you see there. Only constructor you need is the following one:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> AccountController(UserManager&lt;ApplicationUser&gt; userManager)
{
    UserManager = userManager;
}</pre>
</div>
</div>
<p>Now you should be able to build the project successfully and from that point on, you can uninstall the Microsoft.AspNet.Identity.EntityFramework package which you don't need anymore. Lastly, we need to provide an instance of UserManager&lt;ApplicationUser&gt; to our account controller. I'm going to use Autofac IoC container for that operation to inject the dependency into my project. However, you can choose any IoC container you like. After I install the Autofac.Mvc5 package, here how my Global class looks like inside Global.asax.cs file:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">using</span> AspNet.Identity.RavenDB.Stores;
<span style="color: blue;">using</span> AspNetIndetityRavenDb.Models;
<span style="color: blue;">using</span> Autofac;
<span style="color: blue;">using</span> Autofac.Integration.Mvc;
<span style="color: blue;">using</span> Microsoft.AspNet.Identity;
<span style="color: blue;">using</span> Raven.Client;
<span style="color: blue;">using</span> Raven.Client.Document;
<span style="color: blue;">using</span> Raven.Client.Extensions;
<span style="color: blue;">using</span> System.Reflection;
<span style="color: blue;">using</span> System.Web.Mvc;
<span style="color: blue;">using</span> System.Web.Optimization;
<span style="color: blue;">using</span> System.Web.Routing;

<span style="color: blue;">namespace</span> AspNetIndetityRavenDb
{
    <span style="color: blue;">public</span> <span style="color: blue;">class</span> MvcApplication : System.Web.HttpApplication
    {
        <span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start()
        {
            AreaRegistration.RegisterAllAreas();
            FilterConfig.RegisterGlobalFilters(GlobalFilters.Filters);
            RouteConfig.RegisterRoutes(RouteTable.Routes);
            BundleConfig.RegisterBundles(BundleTable.Bundles);

            <span style="color: blue;">const</span> <span style="color: blue;">string</span> RavenDefaultDatabase = <span style="color: #a31515;">"Users"</span>;
            ContainerBuilder builder = <span style="color: blue;">new</span> ContainerBuilder();
            builder.Register(c =&gt;
            {
                IDocumentStore store = <span style="color: blue;">new</span> DocumentStore
                {
                    Url = <span style="color: #a31515;">"http://localhost:8080"</span>,
                    DefaultDatabase = RavenDefaultDatabase
                }.Initialize();

                store.DatabaseCommands.EnsureDatabaseExists(RavenDefaultDatabase);

                <span style="color: blue;">return</span> store;

            }).As&lt;IDocumentStore&gt;().SingleInstance();

            builder.Register(c =&gt; c.Resolve&lt;IDocumentStore&gt;()
                .OpenAsyncSession()).As&lt;IAsyncDocumentSession&gt;().InstancePerHttpRequest();
            builder.Register(c =&gt; <span style="color: blue;">new</span> RavenUserStore&lt;ApplicationUser&gt;(c.Resolve&lt;IAsyncDocumentSession&gt;(), <span style="color: blue;">false</span>))
                .As&lt;IUserStore&lt;ApplicationUser&gt;&gt;().InstancePerHttpRequest();
            builder.RegisterType&lt;UserManager&lt;ApplicationUser&gt;&gt;().InstancePerHttpRequest();

            builder.RegisterControllers(Assembly.GetExecutingAssembly());

            DependencyResolver.SetResolver(<span style="color: blue;">new</span> AutofacDependencyResolver(builder.Build()));
        }
    }
}</pre>
</div>
</div>
<p>Before starting up my application, I expose my RavenDB engine through http://localhost:8080. and I'm all set to fly right now. The default project template allows me to register and log in the application and we can perform all those actions now.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP_955C/2.png"><img height="345" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP_955C/2_thumb.png" alt="2" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="2" /></a></p>
<p>The same sample application available inside the repository as well if you are interested in: <a href="https://github.com/tugberkugurlu/AspNet.Identity.RavenDB/tree/master/samples/AspNet.Identity.RavenDB.Sample.Mvc">AspNet.Identity.RavenDB.Sample.Mvc</a>.</p>
<p>The current ASP.NET Identity system doesn't provide that many features which we require in real world applications such as account confirmation, password reset but it provides us a really great infrastructure and the UserManager&lt;TUser&gt; class saves us from writing bunch of code. I'm sure we will see all other implementations of ASP.NET Identity such as MongoDB, Windows Azure Table Storage, etc. from the community.</p>