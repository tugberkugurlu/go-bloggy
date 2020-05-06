---
title: Windows Azure Management Client Libraries for .NET and It Supports Portable
  Class Library
abstract: One of the missing pieces of the Windows Azure story is within our reach
  now! A few days ago Azure folks have released Windows Azure .NET Management Client
  Libraries
created_at: 2013-10-30 06:17:00 +0000 UTC
tags:
- .net
- Windows Azure
slugs:
- windows-azure-management-client-libraries-fornet-and-it-supports-portable-class-library
---

<p>One of the missing pieces of the <a href="http://windowsazure.com">Windows Azure</a> story is within our reach now! <a href="https://github.com/WindowsAzure?tab=members">Awesome Azure folks</a> have released Windows Azure .NET Management Client Libraries which also support portable class library (PCL). It is now possible to build custom Windows Azure management applications in .NET ecosystem with possible minimum effort. <a href="http://twitter.com/bradygaster">Brady Gaster</a> has a <a href="http://www.bradygaster.com/post/getting-started-with-the-windows-azure-management-libraries">very detailed blog post about the first release of these libraries</a> if you would like to get started with it really quickly. In this post, I'll look at the these libraries in a web developer point of view.</p>
<p>Management libraries are available as NuGet packages and they are at the pre-release stage for now (ignore the HdInsight package below):</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Windows-Azure-Management_ECAA/image.png"><img height="228" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Windows-Azure-Management_ECAA/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Some packages are there for the infrastructural purposes such as tracing, exception handling, etc. and other packages such as Storage and Compute has .NET clients to manage specific services which is good and bad at the same time. When we take a first look at the libraries, we can see a few great things:</p>
<ul>
<li>All operations are based on the <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httpclient(v=vs.110).aspx">HttpClient</a>. This gives us the opportunity to be very flexible (if the client infrastructure allows us) by injecting our inner handler to process the requests.</li>
<li>All client operations are asynchronous from top to bottom.</li>
<li>It supports Portable Class Library.</li>
</ul>
<p>Getting started is also very easy. The following code is all I needed to do to query the storage accounts I have inside my subscription:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args)
{
	<span style="color: blue;">const</span> <span style="color: blue;">string</span> CertThumbprint = <span style="color: #a31515;">"your-cert-thumbprint"</span>;
	<span style="color: blue;">const</span> <span style="color: blue;">string</span> SubscriptionId = <span style="color: #a31515;">"your-subscription-id"</span>;
	X509Certificate2 cert = Utils.FindX509Certificate(CertThumbprint);
	SubscriptionCloudCredentials creds = 
		<span style="color: blue;">new</span> CertificateCloudCredentials(SubscriptionId, cert);

	StorageManagementClient storageClient = 
		CloudContext.Clients.CreateStorageManagementClient(creds);

	StorageServiceListResponse response = storageClient.StorageAccounts.List();
	<span style="color: blue;">foreach</span> (StorageServiceListResponse.StorageService storageService <span style="color: blue;">in</span> response)
	{
		Console.WriteLine(storageService.Uri);
	}

	Console.ReadLine();
}</pre>
</div>
</div>
<p>The client certificate authentication is performed to authenticate our requests. Notice that how easy it was to get a hold of a storage client. By using the CloudContext class, I was able to create the storage client by using the CreateStorageManagementClient extension method which lives under the Microsoft.WindowsAzure.Management.Storage assembly.</p>
<h3>How HttpClient is Used Underneath</h3>
<p>If we look under the carpet, we will see that each time I call the CreateStorageManagementClient, we'll use a different instance of the HttpClient which is not we would want. However, I assume that StorageManagementClient's public members are thread safe which means that you can use one instance of that class throughout our web applications but I'm not sure. Nevertheless, I would prefer an approach which is similar to the following one:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: green;">// single instance that I would hold onto throughout my application lifecycle. I would most</span>
<span style="color: green;">// probably handle this instance through my IoC container.</span>
<span style="color: blue;">using</span> (CloudContext cloudContex = <span style="color: blue;">new</span> CloudContext(creds))
{
    <span style="color: green;">// These StorageManagementClient instances are either the same reference or </span>
    <span style="color: green;">// use the same HttpClient underneath.</span>
    StorageManagementClient storageClient = 
        cloudContex.CreateStorageManagementClient();
}</pre>
</div>
</div>
<p>On the other hand, we need to perform a similar operation to create a compute management client to manage our could service and virtual machines: call the CreateComputeManagementClient extension method on the CloudContext.Clients. Here, we have the same behavior in terms of HttpClient instance.</p>
<h3>A Sample Usage in an ASP.NET MVC Application</h3>
<p>I created a very simple <a href="http://aspnetwebstack.codeplex.com">ASP.NET MVC</a> application which only lists storage services, hosted services and web spaces. I used an IoC container (Autofac) to inject the clients through the controller constructor. Below is the registration code I have in my application.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">void</span> Application_Start()
{
    <span style="color: green;">// Lines removed for brevity...</span>

    <span style="color: green;">// Get a hold of the credentials.</span>
    <span style="color: blue;">const</span> <span style="color: blue;">string</span> CertThumbprint = <span style="color: #a31515;">"your-cert-thumbprint"</span>;
    <span style="color: blue;">const</span> <span style="color: blue;">string</span> SubscriptionId = <span style="color: #a31515;">"your-subscription-id"</span>;
    X509Certificate2 cert = FindX509Certificate(CertThumbprint);
    SubscriptionCloudCredentials creds =
        <span style="color: blue;">new</span> CertificateCloudCredentials(SubscriptionId, cert);

    <span style="color: green;">// Get the ContainerBuilder ready and register the MVC controllers.</span>
    ContainerBuilder builder = <span style="color: blue;">new</span> ContainerBuilder();
    builder.RegisterControllers(Assembly.GetExecutingAssembly());

    <span style="color: green;">// Register the clients</span>
    builder.Register(c =&gt; CloudContext.Clients.CreateComputeManagementClient(creds))
        .As&lt;IComputeManagementClient&gt;().InstancePerHttpRequest();
    builder.Register(c =&gt; CloudContext.Clients.CreateStorageManagementClient(creds))
        .As&lt;IStorageManagementClient&gt;().InstancePerHttpRequest();
    builder.Register(c =&gt; CloudContext.Clients.CreateWebSiteManagementClient(creds))
        .As&lt;IWebSiteManagementClient&gt;().InstancePerHttpRequest();

    <span style="color: green;">// Set the dependency resolver.</span>
    AutofacDependencyResolver dependencyResolver = 
        <span style="color: blue;">new</span> AutofacDependencyResolver(builder.Build());

    DependencyResolver.SetResolver(dependencyResolver);
}</pre>
</div>
</div>
<p>Noticed that I registered the management client classes as per-request instance. This option will create separate client instances per each request. As the underlying architecture creates separate HttpClient instances per each client creating, I'll not be using the same HttpClient throughout my application lifecycle. I'm also not sure whether the client classes are thread safe. That's why I went for the per-request option.</p>
<p>My controller is even simpler.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller
{
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IComputeManagementClient _computeClient;
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IStorageManagementClient _storageClient;
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IWebSiteManagementClient _webSiteClient;

    <span style="color: blue;">public</span> HomeController(
        IComputeManagementClient computeClient, 
        IStorageManagementClient storageClient, 
        IWebSiteManagementClient webSiteClient)
    {
        _computeClient = computeClient;
        _storageClient = storageClient;
        _webSiteClient = webSiteClient;
    }

    <span style="color: blue;">public</span> async Task&lt;ActionResult&gt; Index()
    {
        Task&lt;StorageServiceListResponse&gt; storageServiceResponseTask = 
            _storageClient.StorageAccounts.ListAsync();
        Task&lt;HostedServiceListResponse&gt; hostedServiceResponseTask = 
            _computeClient.HostedServices.ListAsync();
        Task&lt;WebSpacesListResponse&gt; webSpaceResponseTask = 
            _webSiteClient.WebSpaces.ListAsync();

        await Task.WhenAll(storageServiceResponseTask, 
            hostedServiceResponseTask, webSpaceResponseTask);

        <span style="color: blue;">return</span> View(<span style="color: blue;">new</span> HomeViewModel 
        { 
            StorageServices = storageServiceResponseTask.Result,
            HostedServices = hostedServiceResponseTask.Result,
            WebSpaces = webSpaceResponseTask.Result
        });
    }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeViewModel
{
    <span style="color: blue;">public</span> IEnumerable&lt;StorageServiceListResponse.StorageService&gt; StorageServices { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> IEnumerable&lt;HostedServiceListResponse.HostedService&gt; HostedServices { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> IEnumerable&lt;WebSpacesListResponse.WebSpace&gt; WebSpaces { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>Notice that I used Task.WhenAll to fire all the asynchronous work in parallel which is best of both worlds. If you also look at the client classes, you will see that operations are divided into logical groups. For example, here is the IComputeManagementClient interface:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">interface</span> IComputeManagementClient
{
    Uri BaseUri { <span style="color: blue;">get</span>; }
    SubscriptionCloudCredentials Credentials { <span style="color: blue;">get</span>; }

    IDeploymentOperations Deployments { <span style="color: blue;">get</span>; }
    IHostedServiceOperations HostedServices { <span style="color: blue;">get</span>; }
    IOperatingSystemOperations OperatingSystems { <span style="color: blue;">get</span>; }
    IServiceCertificateOperations ServiceCertificates { <span style="color: blue;">get</span>; }
    IVirtualMachineDiskOperations VirtualMachineDisks { <span style="color: blue;">get</span>; }
    IVirtualMachineImageOperations VirtualMachineImages { <span style="color: blue;">get</span>; }
    IVirtualMachineOperations VirtualMachines { <span style="color: blue;">get</span>; }

    Task&lt;Models.ComputeOperationStatusResponse&gt; GetOperationStatusAsync(
        <span style="color: blue;">string</span> requestId, 
        CancellationToken cancellationToken);
}</pre>
</div>
</div>
<p>HostedServices, OperatingSystems, VirtualMachines and so on. Each logical group of operations are separated as class properties.</p>
<h3>Getting Inside the HTTP Pipeline</h3>
<p>The current implementation of the libraries also allow us to get into the HTTP pipeline really easily. However, the way we do it today is a little ugly but it depends on the mentioned HttpClient usage underneath. Nevertheless, this extensibility is really promising. Every management client class inherits the ServiceClient&lt;T&gt; class and that class has a method called WithHandler. Every available client provides a method with the same name by calling the underlying WithHandler method from ServiceClient&lt;T&gt;. By calling this method with a DelegatingHandler instance, we can get into the HTTP pipeline and have a chance to manipulate the request processing. For example, we can inject an handler which has a request retry logic if certain error cases are met. Windows Azure Management Library provides an abstract class for a retry handler which has the basic functionality: LinearRetryHandler. As this is an abstract class, we can inherit this class to create our own retry handler. If we need, we can always override the ShouldRetry method but for our case, we don't need it.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CustomRetryHandler : LinearRetryHandler
{
}</pre>
</div>
</div>
<p>Now, we can use our retry handler to create the management clients:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>RetryHandler retryHandler = <span style="color: blue;">new</span> CustomRetryHandler();
IComputeManagementClient client = 
    CloudContext.Clients.CreateComputeManagementClient(creds).WithHandler(retryHandler)</pre>
</div>
</div>
<p>Now, if a request to a Windows Azure Management API fails with some of the certain status codes, the handler will retry the request.</p>
<h3>There are More&hellip;</h3>
<p>There are more features that are worth mentioning but it will make a very lengthy blog post :) Especially the unit testing and tracing are the topics that I'm looking forward to blogging about.</p>