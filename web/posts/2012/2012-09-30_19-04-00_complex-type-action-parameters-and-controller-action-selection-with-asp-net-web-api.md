---
id: 6be30851-6c42-4573-a3a8-f5c0ba9f6b10
title: Complex Type Action Parameters and Controller Action Selection with ASP.NET
  Web API
abstract: How to use complex type action parameters in ASP.NET Web API and involve
  them inside the controller action selection logic
created_at: 2012-09-30 19:04:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- complex-type-action-parameters-and-controller-action-selection-with-asp-net-web-api
---

<p>If you are familiar with <a href="http://www.asp.net/mvc">ASP.NET MVC</a> and trying to find your way with <a href="http://www.asp.net/web-api">ASP.NET Web API</a>, you may have noticed that the default action selection logic with ASP.NET Web API is pretty different than the ASP.NET MVC's. First of all, the action parameters play a huge role on action selection in ASP.NET Web API. Consider the following controller and its two action methods:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController { 

    <span style="color: green;">//GET /api/cars?categoryId=10</span>
    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByCategoryId(<span style="color: blue;">int</span> categoryId) { 
        
        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>,
            <span style="color: #a31515;">"Car 3"</span>
        };
    }
    
    <span style="color: green;">//GET /api/cars?colorId=10</span>
    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByColorId(<span style="color: blue;">int</span> colorId) { 
        
        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>
        };
    }
}</pre>
</div>
</div>
<p>This doesn&rsquo;t going to cause the action ambiguity because the action parameter names are different. The default action selector (<a href="http://msdn.microsoft.com/en-us/library/system.web.http.controllers.apicontrolleractionselector(v=vs.108).aspx">ApiControllerActionSelector</a>) going to extract the action parameter names and try to match those with the URI parameters such as query string and route values. So if a GET request comes to /api/cars?categoryId=10, the GetCarsByCategoryId action method will be invoked. If a GET request comes to /api/cars?colorId=10 in this case, the GetCarsByColorId action method will be called.</p>
<p>It's possible to use complex types as action parameters for GET requests and bind the route and query string values by marking the complex type parameters with <a href="http://msdn.microsoft.com/en-us/library/system.web.http.fromuriattribute(v=vs.108).aspx">FromUriAttribute</a>. However, the default action selection logic only considers simple types which are System.String, System.DateTime, System.Decimal, System.Guid, System.DateTimeOffset and System.TimeSpan. For example, if you have GetCars(Foo foo) and GetCars(Bar bar) methods inside your controller, you will get the ambiguous action error as the complex types are completely ignored by the ApiControllerActionSelector.</p>
<p>Let&rsquo;s take the following as example here:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsByCategoryRequestCommand {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> CategoryId { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Page { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    [Range(1, 50)]
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Take { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsByColorRequestCommand {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> ColorId { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Page { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    [Range(1, 50)]
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Take { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController {

    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByCategoryId(
        [FromUri]CarsByCategoryRequestCommand cmd) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>,
            <span style="color: #a31515;">"Car 3"</span>
        };
    }

    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByColorId(
        [FromUri]CarsByColorRequestCommand cmd) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>
        };
    }
}</pre>
</div>
</div>
<p>We are not performing any logic inside the action here but you can understand from the action parameter types that we are aiming to perform pagination here. So, we are receiving the inputs from the consumer. We can use simple types directy as action parameters but there is no built-in way to validate the simple types and I haven&rsquo;t found an elegant way to hook something up for that. As a result, complex type action parameters comes in handy in such cases.</p>
<p>If we now send a GET request to /api/cars?colorId=23&amp;page=2&amp;take=12, we would get the ambiguity error message:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action-Paramete.NET-Web-API_11EF9/image.png"><img height="371" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action-Paramete.NET-Web-API_11EF9/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>To workaround this issue, I created a new action selector which has the same implementation as the ApiControllerActionSelector and a few tweaks to make this feature work. It wasn&rsquo;t easy at all. The ApiControllerActionSelector is not so extensible and I had to manually rewrite it (honestly, I didn&rsquo;t directly copy-paste. I rewrote the every single line). I also thought that this could make it into the framework. So, I sent a pull request which got rejected:&nbsp; <a href="http://aspnetwebstack.codeplex.com/SourceControl/network/forks/tugberk/aspnetwebstack/contribution/3338">3338</a>. There is also an issue open to make the default action selector more extensible: <a href="http://aspnetwebstack.codeplex.com/workitem/277">#277</a>. I encourage you to go and vote!</p>
<p>So, what can we do for now to make this work? Go and install the latest <a href="https://github.com/WebAPIDoodle/WebAPIDoodle">WebAPIDoodle</a> package from the Official NuGet feed:</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package WebAPIDoodle </code></p>
</div>
<p>This package has a few useful components for ASP.NET Web API and one of them is the ComplexTypeAwareActionSelector. First of all, we need to replace the default action selector with our ComplexTypeAwareActionSelector as below. Note that ComplexTypeAwareActionSelector preserves all the features of the ApiControllerAction selector.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    <span style="color: blue;">var</span> config = GlobalConfiguration.Configuration;
    config.Routes.MapHttpRoute(
        <span style="color: #a31515;">"DefaultApiRoute"</span>,
        <span style="color: #a31515;">"api/{controller}/{id}"</span>,
        <span style="color: blue;">new</span> { id = RouteParameter.Optional }
    );

    <span style="color: green;">// Replace the default action IHttpActionSelector with</span>
    <span style="color: green;">// WebAPIDoodle.Controllers.ComplexTypeAwareActionSelector</span>
    config.Services.Replace(
        <span style="color: blue;">typeof</span>(IHttpActionSelector),
        <span style="color: blue;">new</span> ComplexTypeAwareActionSelector());
}</pre>
</div>
</div>
<p>This package also contains an attribute named UriParametersAttribute which accepts a params string[] parameter. We can apply this attribute to action methods and pass the parameters that we want to be considered during the action selection. The below one shows the sample usage for our above case:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController {

    [UriParameters(<span style="color: #a31515;">"CategoryId"</span>, <span style="color: #a31515;">"Page"</span>, <span style="color: #a31515;">"Take"</span>)]
    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByCategoryId(
        [FromUri]CarsByCategoryRequestCommand cmd) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>,
            <span style="color: #a31515;">"Car 3"</span>
        };
    }

    [UriParameters(<span style="color: #a31515;">"ColorId"</span>, <span style="color: #a31515;">"Page"</span>, <span style="color: #a31515;">"Take"</span>)]
    <span style="color: blue;">public</span> <span style="color: blue;">string</span>[] GetCarsByColorId(
        [FromUri]CarsByColorRequestCommand cmd) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span>[] { 
            <span style="color: #a31515;">"Car 1"</span>,
            <span style="color: #a31515;">"Car 2"</span>
        };
    }
}</pre>
</div>
</div>
<p>If we now send the proper GET requests as below, we should see it working:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action-Paramete.NET-Web-API_11EF9/image_3.png"><img height="243" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action-Paramete.NET-Web-API_11EF9/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action-Paramete.NET-Web-API_11EF9/image_4.png"><img height="243" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action-Paramete.NET-Web-API_11EF9/image_thumb_4.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>You can also grab <a href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/ComplexTypeParamActionSelection">the sample</a> to see this in action. I should also mention that I am not saying that this is the way to go. Clearly, this generates a lot of noise and we can do better here. The one solution would be to inspect the simple type properties of the complex type action parameter without needing the UriParametersAttribute.</p>