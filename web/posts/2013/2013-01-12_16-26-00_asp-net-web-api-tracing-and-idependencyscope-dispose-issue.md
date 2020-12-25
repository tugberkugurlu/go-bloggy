---
id: 62263419-c69e-41e8-a769-936bcc83933f
title: ASP.NET Web API Tracing and IDependencyScope Dispose Issue
abstract: If you enabled tracing on your ASP.NET Web API application, you may see
  a dispose issue for IDependencyScope. Here is why and how you can workaround it.
created_at: 2013-01-12 16:26:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET Web API
slugs:
- asp-net-web-api-tracing-and-idependencyscope-dispose-issue
---

<p><a href="http://www.asp.net/web-api">ASP.NET Web API</a> has a very cool built-in tracing mechanism. The coolest part about this feature is that none of the tracing code is being run if you don&rsquo;t enable it. The mechanism makes use of well-known <a href="http://en.wikipedia.org/wiki/Facade_pattern">Facade Pattern</a> and if you enable tracing by providing your custom <a href="http://msdn.microsoft.com/en-us/library/system.web.http.tracing.itracewriter(v=vs.108).aspx">ITraceWriter</a> implementation and don&rsquo;t replace the default <a href="http://msdn.microsoft.com/en-us/library/system.web.http.tracing.itracemanager(v=vs.108).aspx">ITraceManager</a> implementation, several ASP.NET Web API components (Message Handlers, Controllers, Filters, Formatters, etc.) will be wrapped up inside their tracer implementations (these are internal classes inside System.Web.Http assembly). You can learn more about tracing from <a href="http://www.asp.net/web-api/overview/testing-and-debugging/tracing-in-aspnet-web-api">Tracing in ASP.NET Web API</a> article.</p>
<p>ASP.NET Web API also has this concept of carrying disposable objects inside the request properties bag and two objects are added to this disposable list by the framework as shown below (in the same order):</p>
<ul>
<li>The <a href="http://msdn.microsoft.com/en-us/library/system.web.http.dependencies.idependencyscope(v=vs.108).aspx">IDependencyScope</a> implementation for the request.</li>
<li>The selected <a href="http://msdn.microsoft.com/en-us/library/system.web.http.controllers.ihttpcontroller(v=vs.108).aspx">IHttpController</a> implementation for the request.</li>
</ul>
<p>The <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessageextensions.disposerequestresources(v=vs.108).aspx">DisposeRequestResources</a> extension method for the <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage.aspx">HttpRequestMessage</a> object is invoked by the hosting layer at the end of each request to dispose the registered disposable objects. The invoker is the internal ConvertResponse method of <a href="http://msdn.microsoft.com/en-us/library/system.web.http.webhost.httpcontrollerhandler.aspx">HttpControllerHandler</a> in case of ASP.NET host. The implementation of the DisposeRequestResources extension method is exactly as shown below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> DisposeRequestResources(<span style="color: blue;">this</span> HttpRequestMessage request) {

    <span style="color: blue;">if</span> (request == <span style="color: blue;">null</span>) {
        <span style="color: blue;">throw</span> Error.ArgumentNull(<span style="color: #a31515;">"request"</span>);
    }

    List&lt;IDisposable&gt; resourcesToDispose;
    <span style="color: blue;">if</span> (request.Properties.TryGetValue(HttpPropertyKeys.DisposableRequestResourcesKey, <span style="color: blue;">out</span> resourcesToDispose)) {
        <span style="color: blue;">foreach</span> (IDisposable resource <span style="color: blue;">in</span> resourcesToDispose) {
            <span style="color: blue;">try</span> {
                resource.Dispose();
            }
            <span style="color: blue;">catch</span> {
                <span style="color: green;">// ignore exceptions</span>
            }
        }
        resourcesToDispose.Clear();
    }
}</pre>
</div>
</div>
<p>It iterates through the list and disposes the registered objects one by one. So far so good. However, there is one slightly problem here depending on your ITraceWriter implementation. I&rsquo;m not sure if it&rsquo;s fair to call this a bug, but to me, it really is a bug. Let me explain what it is.</p>
<p>As we know, the IDependencyScope implementation and the selected IHttpController implementation for the request are added to the disposables list in order to be disposed at the end of the request. As we also know, if we enable tracing in ASP.NET Web API, several ASP.NET Web API components wrapped inside their tracer implementations. In case of controllers, this is HttpControllerTracer. Since the IHttpController is also disposable, the HttpControllerTracer overrides the Dispose method to write a begin/end trace record for the dispose action. However, the disposable objects are disposed in order and when the controller&rsquo;s dispose method is called, the dependency scope is already too far gone and if your custom ITraceWriter implementation tries to use the dependency scope, which is bound to that request, you will get an exception there. This doesn&rsquo;t effect your application that much as this exception is swollen by the underlying infrastructure but this is not good. I wrote a little application to demonstrate this (actually, I went a little far and used a few more things for this *little* application <img src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Tracing-and_F529/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" />) and show how to workaround it for now. <a href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/DependencyScopeTracingDisposeBug">The code for this project is available on GitHub</a>.</p>
<p>I created a custom tracer and registered through the global <a href="http://msdn.microsoft.com/en-us/library/system.web.http.httpconfiguration(v=vs.108).aspx">HttpConfiguration</a> instance. This tracer tries to reach the ILoggerService implementation through the dependency scope. The code for my ITraceWriter implementation is as shown below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> WebApiTracer : ITraceWriter {

    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Trace(
        HttpRequestMessage request, 
        <span style="color: blue;">string</span> category, 
        TraceLevel level, 
        Action&lt;TraceRecord&gt; traceAction) {

        <span style="color: blue;">if</span> (level != TraceLevel.Off) {

            TraceRecord record = <span style="color: blue;">new</span> TraceRecord(request, category, level);
            traceAction(record);
            Log(record);
        }
    }

    <span style="color: blue;">private</span> <span style="color: blue;">void</span> Log(TraceRecord traceRecord) {

        IDependencyScope dependencyScope = 
            traceRecord.Request.GetDependencyScope();
            
        ILoggerService loggerService = 
            dependencyScope.GetService(<span style="color: blue;">typeof</span>(ILoggerService)) <span style="color: blue;">as</span> ILoggerService;
            
        <span style="color: green;">// Log the trace data here using loggerService</span>
        
        <span style="color: green;">// Lines omitted for brevity</span>
    }
}</pre>
</div>
</div>
<p>When we run this application in debug mode and send a request against a valid resource which will eventually go inside the controller pipeline (for this application /api/cars), we will see an exception as below:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Tracing-and_F529/1-12-2013-4-41-30-PM.png"><img height="229" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Tracing-and_F529/1-12-2013-4-41-30-PM_thumb.png" alt="1-12-2013 4-41-30 PM" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="1-12-2013 4-41-30 PM" /></a></p>
<p>If we are curios enough and decide to dig a little deeper, we will actually see what is causing this exception.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Tracing-and_F529/1-12-2013-4-42-58-PM_thumb.png"><img height="338" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Tracing-and_F529/1-12-2013-4-42-58-PM_thumb.png" alt="1-12-2013 4-42-58 PM" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="1-12-2013 4-42-58 PM" /></a></p>
<blockquote>
<p>{"Instances cannot be resolved and nested lifetimes cannot be created from this LifetimeScope as it has already been disposed."}</p>
<p>&nbsp;&nbsp; at Autofac.Core.Lifetime.LifetimeScope.CheckNotDisposed()<br />&nbsp;&nbsp; at Autofac.Core.Lifetime.LifetimeScope.ResolveComponent(IComponentRegistration registration, IEnumerable`1 parameters)<br />&nbsp;&nbsp; at Autofac.ResolutionExtensions.TryResolveService(IComponentContext context, Service service, IEnumerable`1 parameters, Object&amp; instance)<br />&nbsp;&nbsp; at Autofac.ResolutionExtensions.ResolveOptionalService(IComponentContext context, Service service, IEnumerable`1 parameters)<br />&nbsp;&nbsp; at Autofac.ResolutionExtensions.ResolveOptional(IComponentContext context, Type serviceType, IEnumerable`1 parameters)<br />&nbsp;&nbsp; at Autofac.ResolutionExtensions.ResolveOptional(IComponentContext context, Type serviceType)<br />&nbsp;&nbsp; at Autofac.Integration.WebApi.AutofacWebApiDependencyScope.GetService(Type serviceType)<br />&nbsp;&nbsp; at DependencyScopeTracingDisposeBug.Tracing.WebApiTracer.Log(TraceRecord traceRecord) in e:\Apps\DependencyScopeTracingDisposeBug\Tracing\WebApiTracer.cs:line 25<br />&nbsp;&nbsp; at DependencyScopeTracingDisposeBug.Tracing.WebApiTracer.Trace(HttpRequestMessage request, String category, TraceLevel level, Action`1 traceAction) in e:\Apps\DependencyScopeTracingDisposeBug\Tracing\WebApiTracer.cs:line 18<br />&nbsp;&nbsp; at System.Web.Http.Tracing.ITraceWriterExtensions.TraceBeginEnd(ITraceWriter traceWriter, HttpRequestMessage request, String category, TraceLevel level, String operatorName, String operationName, Action`1 beginTrace, Action execute, Action`1 endTrace, Action`1 errorTrace)<br />&nbsp;&nbsp; at System.Web.Http.Tracing.Tracers.HttpControllerTracer.System.IDisposable.Dispose()<br />&nbsp;&nbsp; at System.Net.Http.HttpRequestMessageExtensions.DisposeRequestResources(HttpRequestMessage request)</p>
</blockquote>
<p>It wasn&rsquo;t hard to diagnose the problem and the question was how to workaround it for now as the issue is caused by a code which is deep inside the bowel of the hosting layer. There are several workarounds for this problem which I can come up with quickly such as replacing the HttpControllerTracer but the one I applied was very dirty: adding a message handler as the first message handler (so that it runs last (just in case)) and reordering the disposables on the way out. Here is the message handler which performs this operation:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> DisposableRequestResourcesReorderHandler : DelegatingHandler {

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request, CancellationToken cancellationToken) {

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, cancellationToken).Finally(() =&gt; {

            List&lt;IDisposable&gt; disposableResources = 
                request.Properties[HttpPropertyKeys.DisposableRequestResourcesKey] <span style="color: blue;">as</span> List&lt;IDisposable&gt;;
                
            <span style="color: blue;">if</span> (disposableResources != <span style="color: blue;">null</span> &amp;&amp; disposableResources.Count &gt; 1) {

                <span style="color: green;">// 1-) Get the first one (which I know is AutofacWebApiDependencyScope).</span>
                <span style="color: green;">// 2-) Remove it from the list.</span>
                <span style="color: green;">// 3-) Push it at the end of the list.</span>

                IDisposable dependencyScope = disposableResources[0];
                disposableResources.RemoveAt(0);
                disposableResources.Add(dependencyScope);
            }
        }, runSynchronously: <span style="color: blue;">true</span>);
    }
}</pre>
</div>
</div>
<p>Notice that I used an extension method for Task object called Finally. This is a method which you can get by installing the <a href="http://nuget.org/packages/TaskHelpers.Sources">TaskHelpers NuGet package</a>. This allows you to run the continuation no matter what the status of the completed Task. Finally method will also propagate the proper Task back to the caller, which is kind of nice and clean in our case here as we want to run this code no matter what the status of the Task is.</p>
<p>When you run the application after registering this message handler, you will see the controller&rsquo;s dispose trace record being logged successfully.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Tracing-and_F529/1-12-2013-7-35-21-PM.png"><img height="161" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-Web-API-Tracing-and_F529/1-12-2013-7-35-21-PM_thumb.png" alt="1-12-2013 7-35-21 PM" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="1-12-2013 7-35-21 PM" /></a></p>
<p>If you believe this is a bug and should be fixed, please vote for this issue: <a href="http://aspnetwebstack.codeplex.com/workitem/768" title="http://aspnetwebstack.codeplex.com/workitem/768">http://aspnetwebstack.codeplex.com/workitem/768</a>.</p>