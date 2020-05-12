---
id: 80868994-8293-4595-b46e-6e0c5e9c4b32
title: Hierarchical Resource Structure in ASP.NET Web API
abstract: This post explains the concerns behind the hierarchical resource structure
  in ASP.NET Web API such as routing, authorization and ownership.
created_at: 2013-02-04 08:11:00 +0000 UTC
tags:
- .net
- ASP.NET Web API
- C#
slugs:
- hierarchical-resource-structure-in-asp-net-web-api
---

<p>I came across a question on <a href="http://stackoverflow.com">Stackoverflow</a> today about the hierarchical resource structure in <a href="http://www.asp.net/web-api">ASP.NET Web API</a>: <a title="http://stackoverflow.com/questions/14674255" href="http://stackoverflow.com/questions/14674255">http://stackoverflow.com/questions/14674255</a>. The question is basically about the following issue:</p>
<blockquote>
<p>I have the following schema that I'd like to implement in ASP.NET Web API. What is the proper approach?</p>
<p>http://mydomain/api/students<br />http://mydomain/api/students/s123<br />http://mydomain/api/students/s123/classes<br />http://mydomain/api/students/s123/classes/c456</p>
</blockquote>
<p>With this nice hierarchical approach, you have more concerns that routing here in terms of ASP.NET Web API. There is a good sample application which adopts the hierarchical resource structure: <a href="https://github.com/tugberkugurlu/PingYourPackage">PingYourPackage</a>. I definitely suggest you to check it out.</p>
<p>Let me explain the concerns here in details by setting up a sample scenario. This may not be the desired approach for these types of situations but lays out the concerns very well and if you have a better way to eliminate these concerns, I'd be more that happy to hear those.</p>
<p>Let's say you have the below two affiliates inside your database for a shipment company:</p>
<ul>
<li>Affiliate1 (Id: 100) </li>
<li>Affiliate2 (Id: 101)</li>
</ul>
<p>And then assume that these affiliates have some shipments attached to them:</p>
<ul>
<li>Affiliate1 (Key: 100)   
<ul>
<li>Shipment1 (Key: 100) </li>
<li>Shipment2 (Key: 102) </li>
<li>Shipment4 (Key: 104)</li>
</ul>
</li>
<li>Affiliate2 (Key: 101)   
<ul>
<li>Shipment3 (Key: 103) </li>
<li>Shipment5 (Key: 105)</li>
</ul>
</li>
</ul>
<p>Finally, we want to have the following resource structure:</p>
<ul>
<li>GET api/affiliates/{key}/shipments </li>
<li>GET api/affiliates/{key}/shipments/{shipmentKey} </li>
<li>POST api/affiliates/{key}/shipments </li>
<li>PUT api/affiliates/{key}/shipments/{shipmentKey} </li>
<li>DELETE api/affiliates/{key}/shipments/{shipmentKey}</li>
</ul>
<p>In view of ASP.NET Web API, we have three obvious concerns here: routing, authorization and ownership. Let's go through this one by one.</p>
<blockquote>
<p>The below code snippets have been taken from the <a href="https://github.com/tugberkugurlu/PingYourPackage">PingYourPackage</a> source code. They won't probably work if you copy and paste them but you will get the idea.</p>
</blockquote>
<h3>Routing Concerns</h3>
<p>Assume that we are sending a GET request against /api/affiliates/105/shipments/102 (considering our above scenario). Notice that the affiliate key is 105 here which doesn't exist. So, we would want to terminate the request here ASAP. We can achieve this with a per-route message handler as early as possible. The following AffiliateShipmentsDispatcher is responsible for checking the affiliate existence and acting on the result.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> AffiliateShipmentsDispatcher : DelegatingHandler {

  <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
      HttpRequestMessage request, 
      CancellationToken cancellationToken) {

      <span style="color: green;">// We know at this point that the {key} route variable has </span>
      <span style="color: green;">// been supplied. Otherwise, we wouldn't be here. So, just get it.</span>
      IHttpRouteData routeData = request.GetRouteData();
      Guid affiliateKey = Guid.ParseExact(routeData.Values[<span style="color: #a31515;">"key"</span>].ToString(), <span style="color: #a31515;">"D"</span>);

      IShipmentService shipmentService = request.GetShipmentService();
      <span style="color: blue;">if</span> (shipmentService.GetAffiliate(affiliateKey) == <span style="color: blue;">null</span>) {

          <span style="color: blue;">return</span> Task.FromResult(
              request.CreateResponse(HttpStatusCode.NotFound));
      }

      <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, cancellationToken);
  }
}</pre>
</div>
</div>
<p>I am here using a few internal extension methods which are used inside the project but the idea is simple: go to the database and check the existence of the affiliate. If it doesn't exist, terminate the request and return back the "404 Not Found" response. If it exists, continue executing by calling the base.SendAsync method which will invoke the next message handler inside the chain. Which message handler is the next here? Good question, you dear reader! It's going to be the <a href="http://msdn.microsoft.com/en-us/library/system.web.http.dispatcher.httpcontrollerdispatcher(v=vs.108).aspx">HttpControllerDispatcher</a> which basically puts us inside the controller pipeline. To attach this handler to a route, we need to create a pipeline first to include the controller pipeline by chaining AffiliateShipmentsDispatcher and HttpControllerDispatcher together. The following code snippet shows the AffiliateShipmentsHttpRoute registration.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> RouteConfig {

    <span style="color: blue;">public</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> RegisterRoutes(HttpConfiguration config) {

        <span style="color: blue;">var</span> routes = config.Routes;

        <span style="color: green;">// Pipelines</span>
        HttpMessageHandler affiliateShipmentsPipeline =
            HttpClientFactory.CreatePipeline(
                <span style="color: blue;">new</span> HttpControllerDispatcher(config),
                <span style="color: blue;">new</span>[] { <span style="color: blue;">new</span> AffiliateShipmentsDispatcher() });

        <span style="color: green;">// Routes</span>
        routes.MapHttpRoute(
            <span style="color: #a31515;">"AffiliateShipmentsHttpRoute"</span>,
            <span style="color: #a31515;">"api/affiliates/{key}/shipments/{shipmentKey}"</span>,
            defaults: <span style="color: blue;">new</span> { controller = <span style="color: #a31515;">"AffiliateShipments"</span>, shipmentKey = RouteParameter.Optional },
            constraints: <span style="color: blue;">null</span>,
            handler: affiliateShipmentsPipeline);

        routes.MapHttpRoute(
            <span style="color: #a31515;">"DefaultHttpRoute"</span>,
            <span style="color: #a31515;">"api/{controller}/{key}"</span>,
            defaults: <span style="color: blue;">new</span> { key = RouteParameter.Optional },
            constraints: <span style="color: blue;">null</span>);
    }
}</pre>
</div>
</div>
<h3>Authorization Concerns</h3>
<p>If you have some type of authentication in place, you would want to make sure (in our scenario here) that the authenticated user and the requested affiliate resource is related. For example, assume that <strong>Affiliate1</strong> is authenticated under the <strong>Affiliate role</strong> and you have the <a href="http://msdn.microsoft.com/en-us/library/system.web.http.authorizeattribute(v=vs.108).aspx">AuthorizeAttribute</a> registered to check the "Affiliate" role authorization. In this case, you will fail miserably because this means that Affiliate1 can get to the following resource: /api/affiliates/101/shipments which belongs to Affiliate2. We can eliminate this problem with a custom AuthorizeAttribute which is similar to below one:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>[AttributeUsage(AttributeTargets.Class, AllowMultiple = <span style="color: blue;">false</span>)]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> AffiliateShipmentsAuthorizeAttribute : AuthorizeAttribute {

    <span style="color: blue;">public</span> AffiliateShipmentsAuthorizeAttribute() {

        <span style="color: blue;">base</span>.Roles = <span style="color: #a31515;">"Affiliate"</span>;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> <span style="color: blue;">void</span> OnAuthorization(HttpActionContext actionContext) {
        
        <span style="color: blue;">base</span>.OnAuthorization(actionContext);

        <span style="color: green;">// If not authorized at all, don't bother checking for the </span>
        <span style="color: green;">// user - affiliate relation</span>
        <span style="color: blue;">if</span> (actionContext.Response == <span style="color: blue;">null</span>) { 

            <span style="color: green;">// We are here sure that the request has been authorized and </span>
            <span style="color: green;">// the user is in the Affiliate role. We also don't need </span>
            <span style="color: green;">// to check the existence of the affiliate as it has </span>
            <span style="color: green;">// been also already done by AffiliateShipmentsDispatcher.</span>

            HttpRequestMessage request = actionContext.Request;
            Guid affiliateKey = GetAffiliateKey(request.GetRouteData());
            IPrincipal principal = Thread.CurrentPrincipal;
            IShipmentService shipmentService = request.GetShipmentService();
            <span style="color: blue;">bool</span> isAffiliateRelatedToUser =
                shipmentService.IsAffiliateRelatedToUser(
                    affiliateKey, principal.Identity.Name);

            <span style="color: blue;">if</span> (!isAffiliateRelatedToUser) {

                <span style="color: green;">// Set Unauthorized response as the user and </span>
                <span style="color: green;">// affiliate isn't related to each other. You might</span>
                <span style="color: green;">// want to return "404 NotFound" response here if you don't</span>
                <span style="color: green;">// want to expose the existence of the affiliate.</span>
                actionContext.Response = 
                    request.CreateResponse(HttpStatusCode.Unauthorized);
            }
        }
    }

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> Guid GetAffiliateKey(IHttpRouteData routeData) {

        <span style="color: blue;">var</span> affiliateKey = routeData.Values[<span style="color: #a31515;">"key"</span>].ToString();
        <span style="color: blue;">return</span> Guid.ParseExact(affiliateKey, <span style="color: #a31515;">"D"</span>);
    }
}</pre>
</div>
</div>
<p>This will be registered at the controller level for the AffiliateShipmentsController.</p>
<h3>Ownership Concerns</h3>
<p>Consider this URI for an HTTP GET request: /api/affiliates/100/shipments/102. This URI should get us the correct data. However, what would happen for the this URI: /api/affiliates/100/shipments/103? This should get you a "404 Not Found" HTTP response because the <strong>affiliate</strong> whose Id is <strong>100</strong> doesn't own the <strong>shipment</strong> whose id is <strong>103</strong>. Inside the PingYourPackage project, I ensured the ownership of the resource with the following authorization filter which will be applied to proper action methods.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>[AttributeUsage(AttributeTargets.Method, AllowMultiple = <span style="color: blue;">false</span>)]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> EnsureShipmentOwnershipAttribute 
    : Attribute, IAuthorizationFilter {

    <span style="color: blue;">private</span> <span style="color: blue;">const</span> <span style="color: blue;">string</span> ShipmentDictionaryKey = 
        <span style="color: #a31515;">"__AffiliateShipmentsController_Shipment"</span>;
        
    <span style="color: blue;">public</span> <span style="color: blue;">bool</span> AllowMultiple { <span style="color: blue;">get</span> { <span style="color: blue;">return</span> <span style="color: blue;">false</span>; } }

    <span style="color: blue;">public</span> Task&lt;HttpResponseMessage&gt; ExecuteAuthorizationFilterAsync(
        HttpActionContext actionContext,
        CancellationToken cancellationToken,
        Func&lt;Task&lt;HttpResponseMessage&gt;&gt; continuation) {

        <span style="color: green;">// We are here sure that the user is authanticated and request </span>
        <span style="color: green;">// can be kept executing because the AuthorizeAttribute has </span>
        <span style="color: green;">// been invoked before this filter's OnActionExecuting method.</span>
        <span style="color: green;">// Also, we are sure that the affiliate is associated with</span>
        <span style="color: green;">// the currently authanticated user as the previous action filter </span>
        <span style="color: green;">// has checked against this.</span>
        IHttpRouteData routeData = actionContext.Request.GetRouteData();
        Uri requestUri = actionContext.Request.RequestUri;

        Guid affiliateKey = GetAffiliateKey(routeData);
        Guid shipmentKey = GetShipmentKey(routeData, requestUri);

        <span style="color: green;">// Check if the affiliate really owns the shipment</span>
        <span style="color: green;">// whose key came from the request. We don't need to check the </span>
        <span style="color: green;">// existence of the affiliate as this check has been already </span>
        <span style="color: green;">// performed by the AffiliateShipmentsDispatcher.</span>
        IShipmentService shipmentService = 
            actionContext.Request.GetShipmentService();
        Shipment shipment = shipmentService.GetShipment(shipmentKey);

        <span style="color: green;">// Check the shipment existance</span>
        <span style="color: blue;">if</span> (shipment == <span style="color: blue;">null</span>) {

            <span style="color: blue;">return</span> Task.FromResult(
                <span style="color: blue;">new</span> HttpResponseMessage(HttpStatusCode.NotFound));
        }

        <span style="color: green;">// Check the shipment ownership</span>
        <span style="color: blue;">if</span> (shipment.AffiliateKey != affiliateKey) {

            <span style="color: green;">// You might want to return "404 NotFound" response here </span>
            <span style="color: green;">// if you don't want to expose the existence of the shipment.</span>
            <span style="color: blue;">return</span> Task.FromResult(
                <span style="color: blue;">new</span> HttpResponseMessage(HttpStatusCode.Unauthorized));
        }

        <span style="color: green;">// Stick the shipment inside the Properties dictionary so </span>
        <span style="color: green;">// that we won't need to have another trip to database.</span>
        <span style="color: green;">// The ShipmentParameterBinding will bind the Shipment param</span>
        <span style="color: green;">// if needed.</span>
        actionContext.Request
            .Properties[ShipmentDictionaryKey] = shipment;

        <span style="color: green;">// The request is legit, continue executing.</span>
        <span style="color: blue;">return</span> continuation();
    }

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> Guid GetAffiliateKey(IHttpRouteData routeData) {

        <span style="color: blue;">var</span> affiliateKey = routeData.Values[<span style="color: #a31515;">"key"</span>].ToString();
        <span style="color: blue;">return</span> Guid.ParseExact(affiliateKey, <span style="color: #a31515;">"D"</span>);
    }

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> Guid GetShipmentKey(
        IHttpRouteData routeData, Uri requestUri) {

        <span style="color: green;">// We are sure at this point that the shipmentKey value has been</span>
        <span style="color: green;">// supplied (either through route or quesry string) because it </span>
        <span style="color: green;">// wouldn't be possible for the request to arrive here if it wasn't.</span>
        <span style="color: blue;">object</span> shipmentKeyString;
        <span style="color: blue;">if</span> (routeData.Values.TryGetValue(<span style="color: #a31515;">"shipmentKey"</span>, <span style="color: blue;">out</span> shipmentKeyString)) {

            <span style="color: blue;">return</span> Guid.ParseExact(shipmentKeyString.ToString(), <span style="color: #a31515;">"D"</span>);
        }

        <span style="color: green;">// It's now sure that query string has the shipmentKey value</span>
        <span style="color: blue;">var</span> quesryString = requestUri.ParseQueryString();
        <span style="color: blue;">return</span> Guid.ParseExact(quesryString[<span style="color: #a31515;">"shipmentKey"</span>], <span style="color: #a31515;">"D"</span>);
    }
}</pre>
</div>
</div>
<p>Now, this filter can be applied to proper action methods to allow the proper authorization. At the very end, the AffiliateShipmentsController class looks clean and readable:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>[AffiliateShipmentsAuthorize]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> AffiliateShipmentsController : ApiController {

    <span style="color: green;">// We are OK inside this controller in terms of </span>
    <span style="color: green;">// Affiliate existance and its relation with the current </span>
    <span style="color: green;">// authed user has been checked by the handler </span>
    <span style="color: green;">// and AffiliateShipmentsAuthorizeAttribute.</span>

    <span style="color: green;">// The action method which requests the shipment instance:</span>
    <span style="color: green;">// We can just get the shipment as the shipment </span>
    <span style="color: green;">// existance and its ownership by the affiliate has been </span>
    <span style="color: green;">// approved by the EnsureShipmentOwnershipAttribute.</span>
    <span style="color: green;">// The BindShipmentAttribute can bind the shipment from the</span>
    <span style="color: green;">// Properties dictionarty of the HttpRequestMessage instance</span>
    <span style="color: green;">// as it has been put there by the EnsureShipmentOwnershipAttribute.</span>

    <span style="color: blue;">private</span> <span style="color: blue;">const</span> <span style="color: blue;">string</span> RouteName = <span style="color: #a31515;">"AffiliateShipmentsHttpRoute"</span>;
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IShipmentService _shipmentService;

    <span style="color: blue;">public</span> AffiliateShipmentsController(IShipmentService shipmentService) {

        _shipmentService = shipmentService;
    }

    <span style="color: blue;">public</span> PaginatedDto&lt;ShipmentDto&gt; GetShipments(
        Guid key, 
        PaginatedRequestCommand cmd) {

        <span style="color: blue;">var</span> shipments = _shipmentService
            .GetShipments(cmd.Page, cmd.Take, affiliateKey: key);

        <span style="color: blue;">return</span> shipments.ToPaginatedDto(
            shipments.Select(sh =&gt; sh.ToShipmentDto()));
    }

    [EnsureShipmentOwnership]
    <span style="color: blue;">public</span> ShipmentDto GetShipment(
        Guid key, 
        Guid shipmentKey, 
        [BindShipment]Shipment shipment) {

        <span style="color: blue;">return</span> shipment.ToShipmentDto();
    }

    [EmptyParameterFilter(<span style="color: #a31515;">"requestModel"</span>)]
    <span style="color: blue;">public</span> HttpResponseMessage PostShipment(
        Guid key, 
        ShipmentByAffiliateRequestModel requestModel) {

        <span style="color: blue;">var</span> createdShipmentResult =
            _shipmentService.AddShipment(requestModel.ToShipment(key));

        <span style="color: blue;">if</span> (!createdShipmentResult.IsSuccess) {

            <span style="color: blue;">return</span> <span style="color: blue;">new</span> HttpResponseMessage(HttpStatusCode.Conflict);
        }

        <span style="color: blue;">var</span> response = Request.CreateResponse(HttpStatusCode.Created,
            createdShipmentResult.Entity.ToShipmentDto());

        response.Headers.Location = <span style="color: blue;">new</span> Uri(
            Url.Link(RouteName, <span style="color: blue;">new</span> { 
                key = createdShipmentResult.Entity.AffiliateKey,
                shipmentKey = createdShipmentResult.Entity.Key
            })
        );

        <span style="color: blue;">return</span> response;
    }

    [EnsureShipmentOwnership]
    [EmptyParameterFilter(<span style="color: #a31515;">"requestModel"</span>)]
    <span style="color: blue;">public</span> ShipmentDto PutShipment(
        Guid key, 
        Guid shipmentKey,
        ShipmentByAffiliateUpdateRequestModel requestModel,
        [BindShipment]Shipment shipment) {

        <span style="color: blue;">var</span> updatedShipment = _shipmentService.UpdateShipment(
            requestModel.ToShipment(shipment));

        <span style="color: blue;">return</span> updatedShipment.ToShipmentDto();
    }

    [EnsureShipmentOwnership]
    <span style="color: blue;">public</span> HttpResponseMessage DeleteShipment(
        Guid key, 
        Guid shipmentKey,
        [BindShipment]Shipment shipment) {

        <span style="color: blue;">var</span> operationResult = _shipmentService.RemoveShipment(shipment);

        <span style="color: blue;">if</span> (!operationResult.IsSuccess) {

            <span style="color: blue;">return</span> <span style="color: blue;">new</span> HttpResponseMessage(HttpStatusCode.Conflict);
        }

        <span style="color: blue;">return</span> <span style="color: blue;">new</span> HttpResponseMessage(HttpStatusCode.NoContent);
    }
}</pre>
</div>
</div>
<p>As said, I'd love to know how you handle these types of situations in your applications.</p>