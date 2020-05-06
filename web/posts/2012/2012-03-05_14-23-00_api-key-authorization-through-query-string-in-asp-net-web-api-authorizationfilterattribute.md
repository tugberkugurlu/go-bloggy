---
title: API Key Authorization Through Query String In ASP.NET Web API AuthorizationFilterAttribute
abstract: We will see how API key authorization (verification) through query string
  would be implemented In ASP.NET Web API AuthorizationFilterAttribute
created_at: 2012-03-05 14:23:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- api-key-authorization-through-query-string-in-asp-net-web-api-authorizationfilterattribute
---

<blockquote>
<p><strong>Update on the 29th of June, 2012:</strong></p>
<p>The nuget package I use inside this post is not compatable with ASP.NET Web API RC and will not be with ASP.NET Web API RTM. I have another package named <a title="http://nuget.org/packages/WebApiDoodle" href="http://nuget.org/packages/WebApiDoodle">WebAPIDoodle</a>&nbsp;which has the same funtionality as here. The source code for WebAPIDoodle:&nbsp;<a title="https://github.com/tugberkugurlu/WebAPIDoodle" href="https://github.com/tugberkugurlu/WebAPIDoodle">https://github.com/tugberkugurlu/WebAPIDoodle</a></p>
</blockquote>
<blockquote>
<p><strong>Update on the 17th of October, 2012:</strong></p>
<h3>Confession</h3>
<p>@SoyUnEmilio has pointed out a very good topic <a href="http://www.tugberkugurlu.com/archive/api-key-authorization-through-query-string-in-asp-net-web-api-authorizationfilterattribute#2673">in his comment</a> and <a href="http://www.tugberkugurlu.com/archive/api-key-authorization-through-query-string-in-asp-net-web-api-authorizationfilterattribute#2674">I replied back to him</a>&nbsp;but I would like to make it clear inside the post itself, too.</p>
<p>I wrote this article at the early stages of the framework and I now think that I mixed the concepts of authantication and authorization a lot.&nbsp;So, this blog post does not point you to &nbsp;a good way of implementing authantication and authorization inside you ASP.NET Web API application. I don't want to delete the stuff that I don't like anymore. So, this blog post will stay as it is but I don't want to give the wrong impression as well. As the framework is now more mature, my thoughts on authentication and authorization is shaped better.</p>
<p>In my opinion, you should implement your authantication through a message handler and you almost always don't want to try to perform authorization inside that handler. You may only return "Unauthorized" response inside the message handler (depending on your situation) if the user is not authanticated at all but that should be futher than that.&nbsp;Here is a message handler sample for the API Key authentication: <a title="https://github.com/WebAPIDoodle/WebAPIDoodle/blob/dev/src/apps/WebAPIDoodle/Http/Handlers/ApiKeyAuthenticationHandler.cs" href="https://github.com/WebAPIDoodle/WebAPIDoodle/blob/dev/src/apps/WebAPIDoodle/Http/Handlers/ApiKeyAuthenticationHandler.cs">ApiKeyAuthenticationHandler.cs</a>.</p>
<p>As for the authorization part, it can be handled by the System.Web.Http.AuthorizeAttribute. The AuthorizeAttribute checks against the Thread.CurrentPrincipal. So, the principal that you have supplided inside your message handler will be checked against. You also have a change to perform role or user name based authorization through the AuthorizeAttribute.</p>
</blockquote>
<p>If you are going to build a REST API, you don&rsquo;t probably want to expose all the bits and pieces to everyone. Even so, you would like to see who are making request for various reasons. One of the best ways is the API Key verification to enable that.</p>
<p>In a nutshell, I'll try to explain how it works. You give each user an API key (a GUID would be suitable) and ask them to concatenate this key on the request Uri as query string value. This business logic of assigning the keys works like this. But IMO, one thing is important. No matter what you do, do not manage the API key as admin. You need to find a way to make users manage their keys. User should be able to reset their API key whenever they want and by doing that, old key must gets invalid. As much as this part of the process is important, how you very them on your application is another issue.</p>
<p>With <a title="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples" href="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-web-api-tutorials-videos-samples" target="_blank">ASP.NET Web API</a>, it is pretty easy to intercept a request and change the behavior of that request in any level. For API Key verification, we have two options: 1) Creating a <a title="http://msdn.microsoft.com/en-us/library/hh193679(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/hh193679(v=vs.108).aspx" target="_blank">DelegetingHandler</a> and register it as a message handler. 2) Creating an Authorization filter which will be derived from <a title="http://msdn.microsoft.com/en-us/library/system.web.http.filters.authorizationfilterattribute(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.http.filters.authorizationfilterattribute(v=vs.108).aspx" target="_blank">AuthorizationFilterAttribute</a> class. With one of those two ways, we can verify the user according to API Key supplied.</p>
<p>Honestly, I am not sure which one would be the best option. But, it is certain that if you don&rsquo;t want the whole application to be API key verified, a filter is the best option and it can be applied to an action, controller and globally for entire application. On the other hand, message handlers are involved before the filters. So, that might seem better if you would like your whole app to be API key verified.</p>
<p>I have created an API key verification filter and I tried to make it generic so that it can be applied for all different verification scenarios. Let me show you what I mean.</p>
<p>First of all, go get the bits and pieces through <a title="http://nuget.org" href="http://nuget.org" target="_blank">Nuget</a>. The package is <a title="http://nuget.org/packages/TugberkUg.Web.Http" href="http://nuget.org/packages/TugberkUg.Web.Http" target="_blank">TugberkUg.Web.Http</a> and it is a prerelease package for now:</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package TugberkUg.Web.Http -Pre </code></p>
</div>
<p>This package contains other stuff related to ASP.NET Web API. You can check out the source code on <a href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http">https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http</a>.</p>
<p>For the purpose of this post, what we are interested in is <a title="https://github.com/tugberkugurlu/ASPNETWebAPISamples/blob/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http/Filters/ApiKeyAuthAttribute.cs" href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/blob/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http/Filters/ApiKeyAuthAttribute.cs" target="_blank">ApiKeyAuthAttribute</a> class and <a title="https://github.com/tugberkugurlu/ASPNETWebAPISamples/blob/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http/IApiKeyAuthorizer.cs" href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/blob/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http/IApiKeyAuthorizer.cs" target="_blank">IApiKeyAuthorizer</a> interface. The logic works in a very simple way:</p>
<p>We need a class which will be derived from IApiKeyAuthorizer interface.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">interface</span> IApiKeyAuthorizer {

    <span style="color: blue;">bool</span> IsAuthorized(<span style="color: blue;">string</span> apiKey);
    <span style="color: blue;">bool</span> IsAuthorized(<span style="color: blue;">string</span> apiKey, <span style="color: blue;">string</span>[] roles);
}</pre>
</div>
</div>
<p>As you can see, there are two methods here. First one takes only one parameter which is the API key. This method will be invoked if you try to verify the request only based on API key. Unlike the first one, the second method takes two parameters: API key as string and roles as array of string. This one will be invoked if you try to verify the request based on API key and roles.</p>
<p>I have created an in memory API key authorizer to try this out:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> InMemoryApiKeyAuthorizer : IApiKeyAuthorizer {

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> IList&lt;User&gt; _validApiUsers = <span style="color: blue;">new</span> List&lt;User&gt; { 

        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"d9c99318-53b6-4846-8613-e5aecb473066"</span>, 
            Roles = <span style="color: blue;">new</span> List&lt;Role&gt;() { 
                <span style="color: blue;">new</span> Role { Name = <span style="color: #a31515;">"Admin"</span> } 
            }
        },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"dd97a5aa-704e-4c9e-9bd5-5e2828392eee"</span>, 
            Roles = <span style="color: blue;">new</span> List&lt;Role&gt;() { 
                <span style="color: blue;">new</span> Role { Name = <span style="color: #a31515;">"Customer"</span> } 
            }
        },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"b2e684d7-8807-4232-b5fc-1a6e80c175c0"</span>, 
            Roles = <span style="color: blue;">new</span> List&lt;Role&gt;() { 
                <span style="color: blue;">new</span> Role { Name = <span style="color: #a31515;">"Admin"</span> } 
            }
        },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"36171dc0-4925-4b12-a162-0d6d193acb75"</span> },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"c8028fae-4887-4e91-8fa5-9655adae6ec1"</span> },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"c4bdb227-095a-4fde-8db5-1c96d86e897a"</span> },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"ff10e537-44d5-49b3-add2-6011f54de996"</span> },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"3dcd18cf-e373-4436-9171-aa7f20dae23c"</span> },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"17b2663d-df81-4f63-b10e-5ed918a920cf"</span> },
        <span style="color: blue;">new</span> User { ApiKey = <span style="color: #a31515;">"44fffbf2-8b32-4c4c-834a-518dd0279efa"</span> }
    };

    <span style="color: blue;">public</span> <span style="color: blue;">bool</span> IsAuthorized(<span style="color: blue;">string</span> apiKey) {

        <span style="color: blue;">return</span>
            _validApiUsers.Any(x =&gt; x.ApiKey == apiKey);
    }

    <span style="color: blue;">public</span> <span style="color: blue;">bool</span> IsAuthorized(<span style="color: blue;">string</span> apiKey, <span style="color: blue;">string</span>[] roles) {

        <span style="color: blue;">if</span>(_validApiUsers.Any(x =&gt; 
            x.ApiKey == apiKey &amp;&amp; x.Roles.Where(r =&gt; 
                roles.Contains(r.Name)).Count() &gt; 0)) {

            <span style="color: blue;">return</span> <span style="color: blue;">true</span>;
        }

        <span style="color: blue;">return</span> <span style="color: blue;">false</span>;
    }
}</pre>
</div>
</div>
<p>Here, normally you would look at the supplied information about the request and return either true or false. As you can imagine, the request will be verified if true is returned. If not, then we will handle the unauthorized request in a specific way. We will get to there in a minute.</p>
<blockquote>
<p>I am not sure these two parameters enough to see if the request is legitimate or not. I had a chance to easily supply the <a title="http://msdn.microsoft.com/en-us/library/system.web.http.controllers.httpactioncontext(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.http.controllers.httpactioncontext(v=vs.108).aspx" target="_blank">System.Web.Http.Controllers.HttpActionContext</a> as parameter to these methods but I didn&rsquo;t. Let me know what you think.</p>
</blockquote>
<p>Now we have our logic implemented, we can now apply the filter to our application. As you know, filters are attributes and our attribute is as follows:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> ApiKeyAuthAttribute : AuthorizationFilterAttribute {

    <span style="color: blue;">public</span> ApiKeyAuthAttribute(<span style="color: blue;">string</span> apiKeyQueryParameter, Type apiKeyAuthorizerType);

    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Roles { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">protected</span> <span style="color: blue;">virtual</span> <span style="color: blue;">void</span> HandleUnauthorizedRequest(HttpActionContext actionContext);
    <span style="color: blue;">public</span> <span style="color: blue;">override</span> <span style="color: blue;">void</span> OnAuthorization(HttpActionContext actionContext);
    
}</pre>
</div>
</div>
<p>ApiKeyAuthAttribute has only one constructor and takes two parameters: apiKeyQueryParameter as string for the query string parameter name which will carry the API key and apiKeyAuthorizerType as Type which is the type of Api Key Authorizer which implements IApiKeyAuthorizer interface. We also have a public string property called Roles which accepts comma separated list of roles which user needs to be in.</p>
<p>The usage of the filter is simple as follows:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[ApiKeyAuth(<span style="color: #a31515;">"apiKey"</span>, <span style="color: blue;">typeof</span>(InMemoryApiKeyAuthorizer), Roles = <span style="color: #a31515;">"Admin"</span>)]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController {

    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCars() {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span> <span style="color: blue;">string</span>[] { 
            <span style="color: #a31515;">"BMW"</span>,
            <span style="color: #a31515;">"FIAT"</span>,
            <span style="color: #a31515;">"Mercedes"</span>
        };
    }
}</pre>
</div>
</div>
<p>Now, when we hit the site without API Key and with the legitimate user API key <strong>which is not under the Admin role</strong>, we will get 401.0 Unauthorized response back:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/API-Key-Authorization-Throug.NET-Web-API_1A5/apiKey.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="apiKey" border="0" alt="apiKey" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/API-Key-Authorization-Throug.NET-Web-API_1A5/apiKey_thumb.png" width="644" height="345" /></a></p>
<p>This is the default behavior when an unauthorized user sends a request but can be overridden. You need to simply override the HandleUnauthorizedRequest method and implement your own logic.</p>
<p>When we send a request with a proper API key, we will get the expected result:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/API-Key-Authorization-Throug.NET-Web-API_1A5/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/API-Key-Authorization-Throug.NET-Web-API_1A5/image_thumb.png" width="644" height="383" /></a></p>
<p>The sample I used here is also on <a title="http://github.com" href="http://github.com" target="_blank">GitHub</a>: <a href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/samples/ApiKeyAuthAttributeSample">https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/samples/ApiKeyAuthAttributeSample</a></p>
<p>If you have any advice, please comment or fork me on GitHub.</p>