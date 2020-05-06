---
title: 'Simple OAuth Server: Implementing a Simple OAuth Server with Katana OAuth
  Authorization Server Components (Part 1)'
abstract: 'In my previous post, I emphasized a few important facts on my journey of
  building an OAuth authorization server. As great people say: "Talk is cheap. Show
  me the code." It is exactly what I''m trying to do in this blog post. Also, this
  post is the first one in the "Simple OAuth Server" series. '
created_at: 2014-04-01 14:30:38.1533333 +0000 UTC
tags:
- ASP.NET Web API
- HTTP
- Katana
- OAuth
slugs:
- simple-oauth-server-implementing-a-simple-oauth-server-with-katana-oauth-authorization-server-components-part-1
- simple-oauth-server-implementing-a-simple-oauth-server-with-kat
---

<p>In my previous post, I emphasized a few important facts on my journey of <a href="http://www.tugberkugurlu.com/archive/my-baby-steps-to-oauth-2-0-hell-or-should-i-call-it-heaven">building an OAuth authorization server</a>. As great people say: "<a href="http://www.goodreads.com/quotes/437173-talk-is-cheap-show-me-the-code">Talk is cheap. Show me the code.</a>" It is exactly what I'm trying to do in this blog post. Also, this post is the first one in the "Simple OAuth Server" series.  <h3>What are We Trying to Solve Here?</h3> <p>What we want to achieve at the end of the next two blog posts is actually very doable. We want to have a console application where we handle calls to our protected web service endpoints and access them in a delegated manner which means that the client will actually access the resources on behalf of a user (in other words, resource owner). However, we won't be accessing the web service with resource owner's credentials (username and password). Instead, we will use the credentials to obtain an access token through the resource owner credentials grant and use that token to access the resources from that point on. After this blog post, we will expend our needs and build on top of our existing solution with the upcoming posts. That's why this post will be a little bit detailed about how you could set up the project and we will only cover building the OAuth server part.  <h3>Building the Application Infrastructure</h3> <p>I'll start by creating the ASP.NET Web API application. As mentioned, our application will evolve over time with the upcoming posts. So, this post will only cover the minimum requirements. So, bare this in mind just in case. I used the provided project templates in Visual Studio 2013 to create the project. For this blog post content, we only need ASP.NET Web API components to create our project.  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/4e7b9e3e-bbe6-4307-bb10-d6c853b60bb5.png"><img title="Screenshot 2014-03-31 11.19.58" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-03-31 11.19.58" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d63faed7-b220-4bc1-b9ef-32872a23410a.png" width="644" height="394"></a>&nbsp;</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/23e9b51a-f734-41df-a816-32fdeb7e8311.png"><img title="Screenshot 2014-03-31 11.22.19" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-03-31 11.22.19" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ccb274ec-8995-4f15-abab-d5c075896b99.png" width="644" height="453"></a></p> <p>At the time of writing this post, visual Studio 2013 had the old ASP.NET Web API bits and it's worth updating the package before we continue:</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/03f043e2-ee79-4243-948f-c2006a2d141d.png"><img title="Screenshot 2014-03-31 11.26.09" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-03-31 11.26.09" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/cebfedd4-1215-4f94-a8d3-b36340fc1df9.png" width="644" height="264"></a></p> <blockquote> <p>The OAuth authorization server and the ASP.NET Web API endpoints will be hosted inside the same host in our application here. In your production application, you would probably don't want to do this but for our demo purposes, this will be simpler.</p></blockquote> <p>Now we are ready to build on top of the project template. First thing we need is a membership storage system. Nothing would be better than new <a href="http://www.asp.net/identity">ASP.NET Identity</a> components. I will use the official <a href="https://www.nuget.org/packages/Microsoft.Aspnet.Identity.EntityFramework/">Entity Framework port of the ASP.NET Identity</a> for our application here. However, you are free to choose your own data storage engine. <a href="http://odetocode.com">Scott Allen</a> has a great blog post about the <a href="http://odetocode.com/blogs/scott/archive/2014/01/20/implementing-asp-net-identity.aspx">extensibility of ASP.NET Identity</a> and he listed available open source projects which provide additional storage options for ASP.NET Identity such as <a href="https://github.com/tugberkugurlu/AspNet.Identity.RavenDB">AspNet.Identity.RavenDB</a>.</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/24cb6e74-d77b-4101-b1f6-74b9f9572f0b.png"><img title="Screenshot 2014-03-31 11.36.12" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-03-31 11.36.12" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/30161bd6-ed0d-47c6-a93b-98f631fcedb0.png" width="644" height="264"></a></p> <p>There are two more packages that you need to install. One of them is <a href="https://www.nuget.org/packages/Microsoft.AspNet.Identity.Owin">Microsoft.AspNet.Identity.Owin</a>. This package provides several useful extensions you will use while working with ASP.NET Identity on top of OWIN. The other one is <a href="https://www.nuget.org/packages/Microsoft.Owin.Host.SystemWeb">Microsoft.Owin.Host.SystemWeb</a> package which enables OWIN-based applications to run on IIS using the ASP.NET request pipeline.</p> <blockquote> <p>The packages we just installed (Microsoft.AspNet.Identity.Owin) also brought down some other packages as its dependencies. One of those dependency packages is <a href="http://www.nuget.org/packages/Microsoft.Owin.Security.OAuth">Microsoft.Owin.Security.OAuth</a> and this is the core package that includes the components to support any standard OAuth 2.0 authentication workflow. Just wanted to highlight this fact as this is an important part of the project.</p></blockquote> <p>I will create the Entity Framework DbContext which will hold membership and OAuth client data. ASP.NET Identity Entity Framework package already has the DbContext implementation for the membership storage and our context class will be derived from that.</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> OAuthDbContext : IdentityDbContext
{
    <span style="color: blue">public</span> OAuthDbContext()
        : <span style="color: blue">base</span>(<span style="color: #a31515">"OAuthDbContext"</span>)
    {
    }

    <span style="color: blue">public</span> DbSet&lt;Client&gt; Clients { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
}</pre></div></div>
<p>OAuthDbContext class is derived from IdentityDbContext class as you see. Also notice that we have another DbSet property for clients. That will represent the information of the clients. The Client class is a shown below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> Client
{
    <span style="color: blue">public</span> <span style="color: blue">string</span> Id { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">string</span> Name { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">string</span> ClientSecretHash { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> OAuthGrant AllowedGrant { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }

    <span style="color: blue">public</span> DateTimeOffset CreatedOn { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
}</pre></div></div>
<p>This is the minimum that we need from the client to register in our authorization server. For certain grants, the client doesn't need to have a secret but for "Resource Owner Password Credentials Grant", it's mandatory. The client is also allowed for only one grant, that's all. This is not inside the <a href="http://tools.ietf.org/html/rfc6749">OAuth 2.0 specification</a> but it's the recommended approach. OAuthGrant is an enum and has the following values:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">enum</span> OAuthGrant
{
    Code = 1,
    Implicit = 2,
    ResourceOwner = 3,
    Client = 4
}</pre></div></div>
<p>These are all we need for now and we are ready to create the database. I will use <a href="http://msdn.microsoft.com/en-us/data/jj591621.aspx">Entity Framework Migrations</a> feature to stand up the database and seed some data for demo purposes. As a one time process, I need to enable migrations first by running the "Enable-Migrations" command from the <a href="https://docs.nuget.org/docs/start-here/using-the-package-manager-console">Package Manager Console</a>.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c195c37f-6658-4037-bfa1-d387272e2bd1.png"><img title="Screenshot 2014-03-31 14.39.34" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-03-31 14.39.34" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6fd123ea-6e8b-4938-bedd-8d60cddc86ac.png" width="644" height="182"></a></p>
<p>I will run the another command to add a migration code to reflect my context to a database schema: Add-Migration:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/f26ad46b-f609-46cf-a085-b977ccc8cefc.png"><img title="Screenshot 2014-03-31 14.41.44" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-03-31 14.41.44" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1c338697-8220-48ad-9ba0-a5efce22cd97.png" width="644" height="182"></a></p>
<p>Enable-Migration command created an internal class called Configuration and it contains a Seed method. I can use that seed method to inject some data during the database creation process:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">protected</span> <span style="color: blue">override</span> <span style="color: blue">void</span> Seed(SimpleOAuthSample.Models.OAuthDbContext context)
{
    context.Clients.AddOrUpdate(
        client =&gt; client.Name,
        <span style="color: blue">new</span> Client
        {
            Id = <span style="color: #a31515">"42ff5dad3c274c97a3a7c3d44b67bb42"</span>,
            Name = <span style="color: #a31515">"Demo Resource Owner Password Credentials Grant Client"</span>,
            ClientSecretHash = <span style="color: blue">new</span> PasswordHasher().HashPassword(<span style="color: #a31515">"client123456"</span>),
            AllowedGrant = OAuthGrant.ResourceOwner,
            CreatedOn = DateTimeOffset.UtcNow
        });

    context.Users.AddOrUpdate(
        user =&gt; user.UserName,
        <span style="color: blue">new</span> IdentityUser(<span style="color: #a31515">"Tugberk"</span>)
        {
            Id = Guid.NewGuid().ToString(<span style="color: #a31515">"N"</span>),
            PasswordHash = <span style="color: blue">new</span> PasswordHasher().HashPassword(<span style="color: #a31515">"user123456"</span>),
            SecurityStamp = Guid.NewGuid().ToString(),
            Email = <span style="color: #a31515">"tugberk@example.com"</span>,
            EmailConfirmed = <span style="color: blue">true</span>
        });
}</pre></div></div>
<p>Now, I will use the Update-Database command to create my database:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/58c127f2-ccb2-4662-8cbf-9d169fd936f1.png"><img title="Screenshot 2014-03-31 17.04.38" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-03-31 17.04.38" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/45da2c8c-cc3c-46cc-aae3-5b033e5321a7.png" width="644" height="179"></a></p>
<p>This command just created the database with the seed data on my SQL Express:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c5de7ec5-c2e3-4e86-bde2-1b823dd0be65.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a0fef4b9-bbda-4ea6-9ed1-b897622f0ea9.png" width="244" height="215"></a></p>
<p>We will interact with our database mostly through the UserManager class which ASP.NET Identity core library provides. However, we will still use the OAuthDbContext directly. To use those classes efficiently, we need to write some setup code. I'll do this inside the <a href="http://www.asp.net/aspnet/overview/owin-and-katana/owin-startup-class-detection">OWIN Startup class</a>:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> Startup
{
    <span style="color: blue">public</span> <span style="color: blue">void</span> Configuration(IAppBuilder app)
    {
        app.CreatePerOwinContext&lt;OAuthDbContext&gt;(() =&gt; <span style="color: blue">new</span> OAuthDbContext());
        app.CreatePerOwinContext&lt;UserManager&lt;IdentityUser&gt;&gt;(CreateManager);
    }

    <span style="color: blue">private</span> <span style="color: blue">static</span> UserManager&lt;IdentityUser&gt; CreateManager(
        IdentityFactoryOptions&lt;UserManager&lt;IdentityUser&gt;&gt; options,
        IOwinContext context)
    {
        <span style="color: blue">var</span> userStore =
            <span style="color: blue">new</span> UserStore&lt;IdentityUser&gt;(context.Get&lt;OAuthDbContext&gt;());

        <span style="color: blue">var</span> manager =
            <span style="color: blue">new</span> UserManager&lt;IdentityUser&gt;(userStore);

        <span style="color: blue">return</span> manager;
    }
}</pre></div></div>
<p>This is the minimum code that we can write to use the <a href="http://blogs.msdn.com/b/webdev/archive/2014/02/12/per-request-lifetime-management-for-usermanager-class-in-asp-net-identity.aspx">UserManager class inside our OWIN components</a> efficiently. Although I'm not fan of this approach, I chose to do it this way since doing it in <a href="http://www.tugberkugurlu.com/archive/owin-dependencies--an-ioc-container-adapter-into-owin-pipeline">my way</a> would complicate the post.</p>
<h3>OAuth Authorization Server Application with Katana OAuthAuthorizationServerMiddleware</h3>
<p>Here we come to the real meat of the post. I will now set up the OAuth 2.0 token endpoint to support <a href="http://tools.ietf.org/html/rfc6749#section-4.3">Resource Owner Password Credentials Grant</a> by using the <a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.oauthauthorizationservermiddleware(v=vs.111).aspx">OAuthAuthorizationServerMiddleware</a> which comes with the <a href="http://www.nuget.org/packages/Microsoft.Owin.Security.OAuth/2.1.0">Microsoft.Owin.Security.OAuth</a> library. There is a shorthand extension method on IAppBuilder to use this middleware: <a href="http://msdn.microsoft.com/en-us/library/owin.oauthauthorizationserverextensions.useoauthauthorizationserver(v=vs.111).aspx">UseOAuthAuthorizationServer</a>. I will use this extension method to configure my OAuth 2.0 endpoints through the Configuration method of my Startup class:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">void</span> Configuration(IAppBuilder app)
{
    <span style="color: green">//... </span>
	
    app.UseOAuthAuthorizationServer(<span style="color: blue">new</span> OAuthAuthorizationServerOptions
    {
        TokenEndpointPath = <span style="color: blue">new</span> PathString(<span style="color: #a31515">"/oauth/token"</span>),
        Provider = <span style="color: blue">new</span> MyOAuthAuthorizationServerProvider(),
        AccessTokenExpireTimeSpan = TimeSpan.FromMinutes(30),
<span style="color: blue">#if</span> DEBUG
        AllowInsecureHttp = <span style="color: blue">true</span>,
<span style="color: blue">#endif</span>
    });
}</pre></div></div>
<p>I'm passing an instance of <a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.oauthauthorizationserveroptions(v=vs.111).aspx">OAuthAuthorizationServerOptions</a> here and setting a few of its properties. Everything is pretty much self explanatory except of <a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.oauthauthorizationserveroptions.provider(v=vs.111).aspx">Provider</a> property. I'm setting an implementation of <a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.ioauthauthorizationserverprovider(v=vs.111).aspx">IOAuthAuthorizationServerProvider</a> to Provider property to handle the request at the specific places. Fortunately, I didn't have to implement this interface from top to bottom as there is a default implementation of it (<a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.oauthauthorizationserverprovider(v=vs.111).aspx">OAuthAuthorizationServerProvider</a>) and I just needed to override the methods that I needed.</p>
<blockquote>
<p>Spare some time to read the documentation of the <a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.oauthauthorizationserverprovider_methods(v=vs.111).aspx">OAuthAuthorizationServerProvider's methods</a>. Those are pretty detailed and should give you a great head start.</p></blockquote>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> MyOAuthAuthorizationServerProvider : OAuthAuthorizationServerProvider
{
    <span style="color: blue">public</span> <span style="color: blue">override</span> async Task ValidateClientAuthentication(
        OAuthValidateClientAuthenticationContext context)
    {
        <span style="color: blue">string</span> clientId;
        <span style="color: blue">string</span> clientSecret;

        <span style="color: blue">if</span> (context.TryGetBasicCredentials(<span style="color: blue">out</span> clientId, <span style="color: blue">out</span> clientSecret))
        {
            UserManager&lt;IdentityUser&gt; userManager = 
                context.OwinContext.GetUserManager&lt;UserManager&lt;IdentityUser&gt;&gt;();
            OAuthDbContext dbContext = 
                context.OwinContext.Get&lt;OAuthDbContext&gt;();

            <span style="color: blue">try</span>
            {
                Client client = await dbContext
                    .Clients
                    .FirstOrDefaultAsync(clientEntity =&gt; clientEntity.Id == clientId);

                <span style="color: blue">if</span> (client != <span style="color: blue">null</span> &amp;&amp;
                    userManager.PasswordHasher.VerifyHashedPassword(
                        client.ClientSecretHash, clientSecret) == PasswordVerificationResult.Success)
                {
                    <span style="color: green">// Client has been verified.</span>
                    context.OwinContext.Set&lt;Client&gt;(<span style="color: #a31515">"oauth:client"</span>, client);
                    context.Validated(clientId);
                }
                <span style="color: blue">else</span>
                {
                    <span style="color: green">// Client could not be validated.</span>
                    context.SetError(<span style="color: #a31515">"invalid_client"</span>, <span style="color: #a31515">"Client credentials are invalid."</span>);
                    context.Rejected();
                }
            }
            <span style="color: blue">catch</span>
            {
                <span style="color: green">// Could not get the client through the IClientManager implementation.</span>
                context.SetError(<span style="color: #a31515">"server_error"</span>);
                context.Rejected();
            }
        }
        <span style="color: blue">else</span>
        {
            <span style="color: green">// The client credentials could not be retrieved.</span>
            context.SetError(
                <span style="color: #a31515">"invalid_client"</span>, 
                <span style="color: #a31515">"Client credentials could not be retrieved through the Authorization header."</span>);

            context.Rejected();
        }
    }

    <span style="color: blue">public</span> <span style="color: blue">override</span> async Task GrantResourceOwnerCredentials(
        OAuthGrantResourceOwnerCredentialsContext context)
    {
        Client client = context.OwinContext.Get&lt;Client&gt;(<span style="color: #a31515">"oauth:client"</span>);
        <span style="color: blue">if</span> (client.AllowedGrant == OAuthGrant.ResourceOwner)
        {
            <span style="color: green">// Client flow matches the requested flow. Continue...</span>
            UserManager&lt;IdentityUser&gt; userManager = 
                context.OwinContext.GetUserManager&lt;UserManager&lt;IdentityUser&gt;&gt;();

            IdentityUser user;
            <span style="color: blue">try</span>
            {
                user = await userManager.FindAsync(context.UserName, context.Password);
            }
            <span style="color: blue">catch</span>
            {
                <span style="color: green">// Could not retrieve the user.</span>
                context.SetError(<span style="color: #a31515">"server_error"</span>);
                context.Rejected();

                <span style="color: green">// Return here so that we don't process further. Not ideal but needed to be done here.</span>
                <span style="color: blue">return</span>;
            }

            <span style="color: blue">if</span> (user != <span style="color: blue">null</span>)
            {
                <span style="color: blue">try</span>
                {
                    <span style="color: green">// User is found. Signal this by calling context.Validated</span>
                    ClaimsIdentity identity = await userManager.CreateIdentityAsync(
                        user, 
                        DefaultAuthenticationTypes.ExternalBearer);

                    context.Validated(identity);
                }
                <span style="color: blue">catch</span>
                {
                    <span style="color: green">// The ClaimsIdentity could not be created by the UserManager.</span>
                    context.SetError(<span style="color: #a31515">"server_error"</span>);
                    context.Rejected();
                }
            }
            <span style="color: blue">else</span>
            {
                <span style="color: green">// The resource owner credentials are invalid or resource owner does not exist.</span>
                context.SetError(
                    <span style="color: #a31515">"access_denied"</span>, 
                    <span style="color: #a31515">"The resource owner credentials are invalid or resource owner does not exist."</span>);

                context.Rejected();
            }
        }
        <span style="color: blue">else</span>
        {
            <span style="color: green">// Client is not allowed for the 'Resource Owner Password Credentials Grant'.</span>
            context.SetError(
                <span style="color: #a31515">"invalid_grant"</span>, 
                <span style="color: #a31515">"Client is not allowed for the 'Resource Owner Password Credentials Grant'"</span>);

            context.Rejected();
        }
    }
}</pre></div></div>
<p>Petty much all the methods you will implement, you will be given a context class and you can signal the validity of the request at any point by calling the Validated and Rejected method with their provided signatures. I implemented two methods above (<a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.oauthauthorizationserverprovider.validateclientauthentication(v=vs.111).aspx">ValidateClientAuthentication</a> and <a href="http://msdn.microsoft.com/en-us/library/microsoft.owin.security.oauth.oauthauthorizationserverprovider.grantresourceownercredentials(v=vs.111).aspx">GrantResourceOwnerCredentials</a>) and I performed Validated and Rejected at several points as I have seen it fit. </p>
<p>An HTTP POST request made to "/oauth/token" endpoint with response_type parameter set to "password" will first arrive at the ValidateClientAuthentication method. This is the place where you should retrieve the client credentials and validate it. According to OAuth 2.0 specification, the client credentials can also be sent as request parameters. However, I don't think this is such a good idea comparing to sending the credentials through basic authentication. That's why I only tried to get it from the "Authorization" header. If the client credentials are valid, the request will continue. If not, it will not process further and the error response will be returned as described inside the OAuth 2.0 specification.</p>
<p>If the client credentials are valid and the "response_type" parameter is set to password, the request will arrive at the GrantResourceOwnerCredentials method. Inside this method there are three things we will essentials do:</p>
<ul>
<li>Validate the client's allowed grant. I's check if it's set to ResourceOwner. 
<li>If the client's grant type is valid, validate the resource owner credentials. 
<li>If resource owner credentials are valid, generate a claims identity for the resource owner and pass it to the Validated method.</li></ul>
<p>If all goes as expected, the middleware will issue the access token.</p>
<h3>Calling the OAuth Token Endpoint and Getting the Access Token</h3>
<p>Let's try out the pieces that we have built. As you see previously, I have seeded a sample client and a sample user when during the database creation process. I will use those information to generate a valid OAuth 2.0 "Resource Owner Password Credentials Grant" request.</p>
<p><strong>Request:</strong></p><pre>POST http://localhost:53523/oauth/token HTTP/1.1
User-Agent: Fiddler
Content-Type: application/x-www-form-urlencoded
Authorization: Basic NDJmZjVkYWQzYzI3NGM5N2EzYTdjM2Q0NGI2N2JiNDI6Y2xpZW50MTIzNDU2
Host: localhost:53523
Content-Length: 56

grant_type=password&amp;username=Tugberk&amp;password=user123456</pre>
<p><strong>Response:</strong></p><pre>HTTP/1.1 200 OK
Cache-Control: no-cache
Pragma: no-cache
Content-Length: 550
Content-Type: application/json;charset=UTF-8
Expires: -1
Server: Microsoft-IIS/8.0
X-SourceFiles: =?UTF-8?B?RDpcRHJvcGJveFxBcHBzXFNhbXBsZXNcQXNwTmV0SWRlbnRpdHlTYW1wbGVzXFNpbXBsZU9BdXRoU2FtcGxlXFNpbXBsZU9BdXRoU2FtcGxlXG9hdXRoXHRva2Vu?=
X-Powered-By: ASP.NET
Date: Tue, 01 Apr 2014 13:56:32 GMT

{"access_token":"ydbP24rMOATt7TK3dBCjluD2F5LcLkoX8ud39X135x0a1LEvOgsPf0ekm4Lyu2a06Rv_Z105GRZT_NoclgTTf7Slt5_WNfe68zOUq22j6MqW4Fh__Abzjm6I8otDzxvCJpt5d73R-Um6GwTui3LDbcOk5bH2BZuQLTJsNLknbLPu_FdpgkYfBodUoyPiFhv5-gNBEsfp4gCZYfdKtlhaK0wtloZiIzH1_sNPhBt9FavSfThM5BeoWkz8PFxkv_cOsOhOIzK66nSx7B2XL7K9aLqPSJLxus2ud8GBZyteSeFi26L9oX9do7MyCL1nXa8D9DRWfcIXiQi1v19AwyhoupP3L-k89xOK6_NTSzYOVhSMG9Juz8VYHWGkJeYTmekmnVkCvQe7KMQ6PceeUFJnA88TkiHNhai0hV8j012OUxPpUN5zRPJOU81XywSkQ7oKE0UsX3hQamgFrXV9eA-TSwZd4Qr-P9w6a82OM66Te9E","token_type":"bearer","expires_in":1799}</pre>
<p>We successfully retrieved the response and it contains the JSON response body which includes the access token in the format described inside the OAuth 2.0 specification.</p>
<h3>Summary and What is Next</h3>
<p>In this post, we have set up our authorization server and we have a working OAuth 2.0 token endpoint which only supports "Resource Owner Password Credentials Grant" for now. <a href="https://github.com/tugberkugurlu/AspNetIdentitySamples/tree/master/SimpleOAuthSample">The code is available on GitHub</a> if you are interested in. In the next post, we will create our web service and protect it using our authorization server. We will also see how we can call this web service successfully from a typical .NET application.</p>
<h3>Resources</h3>
<ul>
<li><a href="http://leastprivilege.com/2013/11/25/dissecting-the-web-api-individual-accounts-templatepart-1-overview/">Dissecting the Web API Individual Accounts Template–Part 1: Overview</a> 
<li><a href="http://leastprivilege.com/2013/11/26/dissecting-the-web-api-individual-accounts-templatepart-2-local-accounts/">Dissecting the Web API Individual Accounts Template–Part 2: Local Accounts</a> 
<li><a href="http://leastprivilege.com/2014/03/24/the-web-api-v2-oauth2-authorization-server-middlewareis-it-worth-it/">The Web API v2 OAuth2 Authorization Server Middleware–Is it worth it?</a></li></ul>  