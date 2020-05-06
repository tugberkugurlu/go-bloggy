---
title: ASP.NET Web API and Handling ModelState Validation
abstract: How to handle ModelState Validation errors in ASP.NET Web API with an Action
  Filter and HttpError object
created_at: 2012-09-11 10:39:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- asp-net-web-api-and-handling-modelstate-validation
---

<p><a href="https://twitter.com/amirrajan">@amirrajan</a>, who thinks that Web API sucks :), wanted to see how model validation works in ASP.NET Web API. Instead of responding him directly, I decided to blog about it. <a href="http://amirrajan.net/blog/a-usability-issue-with-asp-net-web-api-with-a-solution">In his post</a>, he indicated that we shouldn&rsquo;t use ASP.NET Web API because it has usability issues but I am sure that he hasn&rsquo;t taken the time to really understand the whole framework. It is really obvious that the blog post is written by someone who doesn&rsquo;t know much about the framework. This is not a bad thing. You don&rsquo;t have to know about every single technology but if *you* express that the technology bad, you better know about it. Anyway, this is not the topic.</p>
<p>ASP.NET Web API supports validation attributes from .NET Data Annotations&nbsp;out of the box and same ModelState concept which is present all across the ASP.NET is applicable here as well. You can inspect the <a href="http://msdn.microsoft.com/en-us/library/system.web.http.apicontroller.modelstate(v=vs.108).aspx">ModelState</a> property to see the errors inside your controller action but I bet that you want to just terminate the request and return back a error response if the ModelState is invalid. To do that, you just need to create and register an action filer as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>[AttributeUsage(AttributeTargets.Class | AttributeTargets.Method, AllowMultiple = <span style="color: blue;">false</span>, Inherited = <span style="color: blue;">true</span>)]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> InvalidModelStateFilterAttribute : ActionFilterAttribute {

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> <span style="color: blue;">void</span> OnActionExecuting(HttpActionContext actionContext) {

        <span style="color: blue;">if</span> (!actionContext.ModelState.IsValid) {

            actionContext.Response = actionContext.Request.CreateErrorResponse(
                HttpStatusCode.BadRequest, actionContext.ModelState);
        }
    }
}</pre>
</div>
</div>
<p>This is one time code and you don&rsquo;t need to write this over and over again. <a href="http://msdn.microsoft.com/en-us/library/jj127078(v=vs.108)">CreateErrorResponse</a> extension method for <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage.aspx">HttpRequestMessage</a> uses <a href="http://msdn.microsoft.com/en-us/library/system.web.http.httperror(v=vs.108)">HttpError</a> class under the covers to create a consistent looking error response. The key point here is that the content negotiation will be handled here as well. the conneg will be handled according to the content negotiator service that you registered. If you don&rsquo;t register one, the default one (<a href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.defaultcontentnegotiator(v=vs.108)">DefaultContentNegotiator</a>) will kick in which will be fine for 90% of the time.</p>
<p>I created a tiny controller which simply exposes two GET and one POST endpoint:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CarsController : ApiController {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> CarsContext _carsCtx = <span style="color: blue;">new</span> CarsContext();

    <span style="color: blue;">public</span> IEnumerable&lt;Car&gt; GetCars() {

        <span style="color: blue;">return</span> _carsCtx.All;
    }

    <span style="color: blue;">public</span> Car GetCar(<span style="color: blue;">int</span> id) {

        Tuple&lt;<span style="color: blue;">bool</span>, Car&gt; carResult = _carsCtx.GetSingle(id);
        <span style="color: blue;">if</span> (!carResult.Item1) {

            <span style="color: blue;">throw</span> <span style="color: blue;">new</span> HttpResponseException(HttpStatusCode.NotFound);
        }

        <span style="color: blue;">return</span> carResult.Item2;
    }

    <span style="color: blue;">public</span> HttpResponseMessage PostCar(Car car) {

        <span style="color: blue;">var</span> addedCar = _carsCtx.Add(car);
        <span style="color: blue;">var</span> response = Request.CreateResponse(HttpStatusCode.Created, car);
        response.Headers.Location = <span style="color: blue;">new</span> Uri(
            Url.Link(<span style="color: #a31515;">"DefaultApiRoute"</span>, <span style="color: blue;">new</span> { id = addedCar.Id }));

        <span style="color: blue;">return</span> response;
    }
}</pre>
</div>
</div>
<p>Here is what our Car object looks like with the validation attributes:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Car {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Id { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    [Required]
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Make { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    [Required]
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Model { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Year { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    [Range(0, 200000)]
    <span style="color: blue;">public</span> <span style="color: blue;">float</span> Price { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>And I registered the filter globally inside the Application_Start method:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    <span style="color: blue;">var</span> config = GlobalConfiguration.Configuration;
    config.Routes.MapHttpRoute(
        <span style="color: #a31515;">"DefaultApiRoute"</span>,
        <span style="color: #a31515;">"api/{controller}/{id}"</span>,
        <span style="color: blue;">new</span> { id = RouteParameter.Optional }
    );

    config.Filters.Add(<span style="color: blue;">new</span> InvalidModelStateFilterAttribute());
}</pre>
</div>
</div>
<p>Let&rsquo;s retrieve the list of cars first:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_thumb.png" width="546" height="484" /></a></p>
<p>Let&rsquo;s add a new car:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_thumb_3.png" width="542" height="484" /></a></p>
<p>When we make a new request, it is now on the list:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_4.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_thumb_4.png" width="540" height="484" /></a></p>
<p>Now, let&rsquo;s try to add a new car again but this time, our object won&rsquo;t fit the requirements and we should see the error response:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_5.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_thumb_5.png" width="539" height="484" /></a></p>
<p>Here is the better look of the error message that we got:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image61ae3dd3-c6bd-4b85-a18d-68cc0af23d9f.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ef0431a533f6_DE12/image_thumb_6.png" width="644" height="305" /></a></p>
<p>Here is the plain text version of this message:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>{
    <span style="color: #a31515;">"Message"</span>: <span style="color: #a31515;">"The request is invalid."</span>,
    <span style="color: #a31515;">"ModelState"</span>: { 
        <span style="color: #a31515;">"car"</span>: [
            <span style="color: #a31515;">"Required property 'Make' not found in JSON. Path '', line 1, position 57."</span>
        ],
        <span style="color: #a31515;">"car.Make"</span> : [
            <span style="color: #a31515;">"The Make field is required."</span>
        ], 
        <span style="color: #a31515;">"car.Price"</span>: [
            <span style="color: #a31515;">"The field Price must be between 0 and 200000."</span>
        ]
    }
}</pre>
</div>
</div>
<p>This doesn&rsquo;t have to be in this format. You can play with the ModelState and iterate over the error messages to simply create your own error message format.</p>
<p>To be honest, this is no magic and this, by itself, doesn&rsquo;t make ASP.NET Web API super awesome but please don&rsquo;t make accusations on a technology without playing with it first at least enough amount of time to grasp the big picture. Otherwise, IMO, it would be a huge disrespect for the people who made the technology happen.</p>