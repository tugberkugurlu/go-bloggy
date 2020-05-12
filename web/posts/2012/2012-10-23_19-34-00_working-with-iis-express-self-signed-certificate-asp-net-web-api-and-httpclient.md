---
id: 0ef2ff7f-c792-4a99-9578-3fa5919e2936
title: Working with IIS Express Self-signed Certificate, ASP.NET Web API and HttpClient
abstract: We will see how to smoothly work with IIS Express Self-signed Certificate,
  ASP.NET Web API and HttpClient by placing the self-signed certificate in the Trusted
  Root CA store.
created_at: 2012-10-23 19:34:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
- PowerShell
- Visual Studio
slugs:
- working-with-iis-express-self-signed-certificate-asp-net-web-api-and-httpclient
---

<p>If you would like to only support HTTPS with your <a href="http://www.asp.net/web-api">ASP.NET Web API</a> application, you might also want expose your application through HTTPS during the development time. This is not a big problem if you are heavily integration-testing your Web API as you can pass anything you want as a host name but if you are building your HTTP API wrapper simultaneously, you want to sometimes do manual testing to see if it&rsquo;s actually working. There are sevaral ways to sort this problem out and one of them is provided directly by Visual Studio. Visual Studio allows us to create HTTPS bindings to use with IIS Express during development time. Let&rsquo;s see how that works.</p>
<blockquote>
<p><strong>Note:</strong> I am assuming everybody understands that I am talking about ASP.NET Web API web host scenario with ASP.NET here. This blog post is not about self-host scenarios.</p>
</blockquote>
<p>First of all, I created an empty web application through visual studio. Then, I added one of my NuGet packages: <a href="http://nuget.org/packages/WebAPIDoodle.Bootstrap.Sample.Complex">WebAPIDoodle.Bootstrap.Sample.Complex</a>. This package will get you all ASP.NET Web API stuff and a working sample with all CRUD operations.</p>
<p>I also created a message handler which is going to ensure that our API is only going to be exposed over HTTPS.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> RequireHttpsMessageHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request, 
        CancellationToken cancellationToken) {

        <span style="color: blue;">if</span> (request.RequestUri.Scheme != Uri.UriSchemeHttps) {

            <span style="color: blue;">var</span> forbiddenResponse = 
                request.CreateResponse(HttpStatusCode.Forbidden);

            forbiddenResponse.ReasonPhrase = <span style="color: #a31515;">"SSL Required"</span>;
            <span style="color: blue;">return</span> Task.FromResult&lt;HttpResponseMessage&gt;(forbiddenResponse);
        }

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, cancellationToken);
    }
}</pre>
</div>
</div>
<p>Then, I registered this message handler as you can see below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start(<span style="color: blue;">object</span> sender, EventArgs e) {

    <span style="color: blue;">var</span> config = GlobalConfiguration.Configuration;

    <span style="color: green;">//...</span>

    config.MessageHandlers.Add(<span style="color: blue;">new</span> RequireHttpsMessageHandler());
}</pre>
</div>
</div>
<p>To configure the HTTPS endpoint with IIS Express, I simply need to click on the web application project and press F4. This will bring up the project properties and there will be a option there called "SSL Enabled".</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; margin: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb_3.png" width="219" height="244" /></a></p>
<p>By default, this is set to False as you can see. If we change this and set it to True, Visual Studio will create the necessary binding for our application by assigning a new port number and attaching the pre-generated self-signed certificate for that endpoint.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_4.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; margin: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb_4.png" width="244" height="165" /></a></p>
<p>Now, when we open up a browser and navigate to that HTTPS endpoint, we should face a scary looking error:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_5.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb_5.png" width="644" height="334" /></a></p>
<p>As our certificate is a self-signed certificate, the browser doesn&rsquo;t trust that and gives this error. This error is not a blocking issue for us and we can just click the Proceed anyway button to suppress this error and everything will be work just fine.</p>
<p>Let&rsquo;s close the browser and create a very simple and na&iuml;ve .NET console client for our API as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">class</span> Program {

    <span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args) {

        Console.WriteLine(GetStringAsync().Result);
        Console.ReadLine();
    }

    <span style="color: blue;">public</span> <span style="color: blue;">static</span> async Task&lt;<span style="color: blue;">string</span>&gt; GetStringAsync() {

        <span style="color: blue;">using</span> (HttpClient client = <span style="color: blue;">new</span> HttpClient()) {

            <span style="color: blue;">return</span> await client
               .GetStringAsync(<span style="color: #a31515;">"https://localhost:44304/api/cars"</span>);
        }
    }
}</pre>
</div>
</div>
<p>As I said, it is insanely simple <img class="wlEmoticon wlEmoticon-smile" style="border-style: none;" alt="Smile" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/wlEmoticon-smile.png" /> but will serve for our demo purposes. If we run this console application along with our Web API application, we should see an exception thrown as below:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_6.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb_6.png" width="644" height="288" /></a></p>
<p>This time the HttpClient is nagging at us because it didn&rsquo;t trust the server certificate. If we open up the browser again and take a look at the certificate, we will see that it is also complaining that this self-signed certificate is not in the Trusted Root CA store.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_7.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb63146b9a-c220-4c87-93a3-eaed3f9517b0.png" width="522" height="484" /></a></p>
<p>One way to get rid of this problem is to place this self-signed certificate in the Trusted Root CA store and the error will go away. Let&rsquo;s first open up a PowerShell Command window and see where this self-signed certificate lives.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_8.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb_7.png" width="644" height="184" /></a></p>
<p>As we now know where the certificate is, we can grab this certificate and place it in the Trusted Root CA store. There are several ways of doing this but I love PowerShell, so I am going to do this with PowerShell, too. To be able to perform the below operations, we need to open up an elevated PowerShell Command Prompt.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_9.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb.png" width="644" height="228" /></a></p>
<p>Let&rsquo;s explain what we did above:</p>
<ul>
<li>We retrieved our self-signed certificate and stuck it into a variable. </li>
<li>Then, we retrieved the Trusted Root CA store and stuck that into a variable, too. </li>
<li>Lastly, we opened up the store with Read and Write access and added the certificate there. As a final step, we closed the store.</li>
</ul>
<p>After this change, if we take a look at the Trusted Root CA store, we will see our certificate there:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/SNAGHTML1d009a0f.png"><img title="SNAGHTML1d009a0f" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="SNAGHTML1d009a0f" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/SNAGHTML1d009a0f_thumb.png" width="644" height="219" /></a></p>
<p>If we now run our console client application, we should see it working smoothly.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_10.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/image_thumb_8.png" width="644" height="353" /></a></p>
<p>I hope this will help you as much as it helped me <img class="wlEmoticon wlEmoticon-smile" style="border-style: none;" alt="Smile" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Working-with-IIS-Express-Sel.NET-Web-API_453/wlEmoticon-smile.png" /></p>