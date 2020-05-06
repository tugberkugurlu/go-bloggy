---
title: 'Logging in the OWIN World with Microsoft.Owin: Introduction'
abstract: Microsoft implementation of OWIN (Microsoft.Owin or Katana for now) has
  a great logging infrastructure and this post will walk you through the basics of
  this component.
created_at: 2013-09-18 16:14:00 +0000 UTC
tags:
- .net
- Katana
- OWIN
slugs:
- logging-in-the-owin-world-with-microsoft-owin--introduction
---

<p>Microsoft implementation of <a href="http://owin.org/">OWIN</a> (called <a href="http://katanaproject.codeplex.com">Microsoft.Owin or Katana</a> for now) has a great infrastructure for logging under Microsoft.Owin.Logging namespace inside the <a href="http://www.nuget.org/packages/Microsoft.Owin/">Microsoft.Owin</a> assembly. Same as all other Microsoft.Owin components, the logging infrastructure is built on top of interfaces and there are two of them: ILoggerFactory and ILogger. The default implementation is using .NET Framework tracing components heavily such as <a href="http://msdn.microsoft.com/en-us/library/system.diagnostics.tracesource.aspx">System.Diagnostics.TraceSource</a> and <a href="http://msdn.microsoft.com/en-us/library/system.diagnostics.traceswitch.aspx">System.Diagnostics.TraceSwitch</a>. You can learn about these components and how to instrument your applications with them through the <a href="http://msdn.microsoft.com/en-us/library/zs6s4h68.aspx">.NET Development Guide's Tracing section on MSDN</a>.</p>
<p>The responsibility of the logger factory is to construct a new logger when needed. You can see how ILoggerFactory interface looks like below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">namespace</span> Microsoft.Owin.Logging
{
    <span style="color: gray;">///</span> <span style="color: gray;">&lt;summary&gt;</span>
    <span style="color: gray;">///</span><span style="color: green;"> Used to create logger instances of the given name.</span>
    <span style="color: gray;">///</span> <span style="color: gray;">&lt;/summary&gt;</span>
    <span style="color: blue;">public</span> <span style="color: blue;">interface</span> ILoggerFactory
    {
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;summary&gt;</span>
        <span style="color: gray;">///</span><span style="color: green;"> Creates a new ILogger instance of the given name.</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;/summary&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;param name="name"&gt;</span><span style="color: gray;">&lt;/param&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;returns&gt;</span><span style="color: gray;">&lt;/returns&gt;</span>
        ILogger Create(<span style="color: blue;">string</span> name);
    }
}</pre>
</div>
</div>
<p>As you can see, the Create method returns an ILogger implementation. The ILogger implementation is solely responsible for writing the logs through its one and only method: WriteCore.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">namespace</span> Microsoft.Owin.Logging
{
    <span style="color: gray;">///</span> <span style="color: gray;">&lt;summary&gt;</span>
    <span style="color: gray;">///</span><span style="color: green;"> A generic interface for logging.</span>
    <span style="color: gray;">///</span> <span style="color: gray;">&lt;/summary&gt;</span>
    <span style="color: blue;">public</span> <span style="color: blue;">interface</span> ILogger
    {
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;summary&gt;</span>
        <span style="color: gray;">///</span><span style="color: green;"> Aggregates most logging patterns to a single method.  This must be compatible with the Func representation in the OWIN environment.</span>
        <span style="color: gray;">///</span><span style="color: green;"> </span>
        <span style="color: gray;">///</span><span style="color: green;"> To check IsEnabled call WriteCore with only TraceEventType and check the return value, no event will be written.</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;/summary&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;param name="eventType"&gt;</span><span style="color: gray;">&lt;/param&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;param name="eventId"&gt;</span><span style="color: gray;">&lt;/param&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;param name="state"&gt;</span><span style="color: gray;">&lt;/param&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;param name="exception"&gt;</span><span style="color: gray;">&lt;/param&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;param name="formatter"&gt;</span><span style="color: gray;">&lt;/param&gt;</span>
        <span style="color: gray;">///</span> <span style="color: gray;">&lt;returns&gt;</span><span style="color: gray;">&lt;/returns&gt;</span>
        <span style="color: blue;">bool</span> WriteCore(TraceEventType eventType, <span style="color: blue;">int</span> eventId, <span style="color: blue;">object</span> state, 
            Exception exception, Func&lt;<span style="color: blue;">object</span>, Exception, <span style="color: blue;">string</span>&gt; formatter);
    }
}</pre>
</div>
</div>
<p>For these interfaces, there are two implementations provided by the Microsoft.Owin assembly for logging: DiagnosticsLoggerFactory (publicly exposed) and DiagnosticsLogger (internally exposed). However, these are the types provided by the Microsoft's OWIN implementation and don't belong to core .NET framework. As it's discouraged to put non-framework types inside the IAppBuilder properties dictionary and request's environment, a Func is put inside the IAppBuilder properties dictionary with a key named "server.LoggerFactory". Here is the signature of that Func.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">using</span> TraceFactoryDelegate = 
    Func
    &lt;
        <span style="color: blue;">string</span>, 
        Func
        &lt;
            TraceEventType, 
            <span style="color: blue;">int</span>, 
            <span style="color: blue;">object</span>, 
            Exception, 
            Func
            &lt;
                <span style="color: blue;">object</span>, 
                Exception, 
                <span style="color: blue;">string</span>
            &gt;, 
            <span style="color: blue;">bool</span>
        &gt;
    &gt;;</pre>
</div>
</div>
<p>Through this delegate, you are expected to create a Func object based on the logger name and this Func will be used to write logs. The currently provided extension methods to work with the Katana&rsquo;s logging infrastructure hides these delegates from you. Here are a few extension methods under Microsoft.Owin.Logging namespace for the IAppBuilder interface:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> SetLoggerFactory(
    <span style="color: blue;">this</span> IAppBuilder app, ILoggerFactory loggerFactory);
    
<span style="color: blue;">public</span> <span style="color: blue;">static</span> ILoggerFactory GetLoggerFactory(
    <span style="color: blue;">this</span> IAppBuilder app);
    
<span style="color: blue;">public</span> <span style="color: blue;">static</span> ILogger CreateLogger(
    <span style="color: blue;">this</span> IAppBuilder app, <span style="color: blue;">string</span> name);</pre>
</div>
</div>
<p>All these extension methods work on the TraceFactoryDelegate Func instance shown above but it&rsquo;s all hidden from us as mentioned.</p>
<p>So, how we can take advantage of this provided components? It&rsquo;s quite easy actually. You can create separate loggers for each of your components using one of the CreateLogger extension method overloads for the IAppBuilder interface and use that logger to write log messages. The following is just a little sample to give you a hint:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Startup
{
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Configuration(IAppBuilder app)
    {
        Log1(app);
        app.Use&lt;MyCustomMiddleware&gt;(app);
    }

    <span style="color: blue;">private</span> <span style="color: blue;">void</span> Log1(IAppBuilder app) 
    {
        ILogger logger = app.CreateLogger&lt;Startup&gt;();
        logger.WriteError(<span style="color: #a31515;">"App is starting up"</span>);
        logger.WriteCritical(<span style="color: #a31515;">"App is starting up"</span>);
        logger.WriteWarning(<span style="color: #a31515;">"App is starting up"</span>);
        logger.WriteVerbose(<span style="color: #a31515;">"App is starting up"</span>);
        logger.WriteInformation(<span style="color: #a31515;">"App is starting up"</span>);

        <span style="color: blue;">int</span> foo = 1;
        <span style="color: blue;">int</span> bar = 0;

        <span style="color: blue;">try</span>
        {
            <span style="color: blue;">int</span> fb = foo / bar;
        }
        <span style="color: blue;">catch</span> (Exception ex)
        {
            logger.WriteError(<span style="color: #a31515;">"Error on calculation"</span>, ex);
        }
    }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> MyCustomMiddleware : OwinMiddleware
{
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> ILogger _logger;

    <span style="color: blue;">public</span> MyCustomMiddleware(
        OwinMiddleware next, IAppBuilder app) : <span style="color: blue;">base</span>(next)
    {
        _logger = app.CreateLogger&lt;MyCustomMiddleware&gt;();
    }

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> Task Invoke(IOwinContext context)
    {
        _logger.WriteVerbose(
            <span style="color: blue;">string</span>.Format(<span style="color: #a31515;">"{0} {1}: {2}"</span>, 
            context.Request.Scheme, 
            context.Request.Method, 
            context.Request.Path));

        context.Response.Headers.Add(
            <span style="color: #a31515;">"Content-Type"</span>, <span style="color: blue;">new</span>[] { <span style="color: #a31515;">"text/plain"</span> });
            
        <span style="color: blue;">return</span> context.Response.WriteAsync(
            <span style="color: #a31515;">"Logging sample is runnig!"</span>);
    }
}</pre>
</div>
</div>
<p>There are a few things I would like to touch on here:</p>
<ul>
<li>I used generic CreateLogger method to create loggers. This will create loggers based on the name of the type I&rsquo;m passing in. In view of the default logger implementation provided by the Katana, this means that we will create TraceSource instances named as the full type name.</li>
<li>You can see that I receive an IAppBuilder implementation instance through the constructor of my middleware and I used that instance to create a logger specific to my middleware. I can hold onto that logger throughout the AppDomain lifetime as all members of my logger are thread safe.</li>
<li>Instead of using the WriteCore method to write logs, I used several extension methods such as WriteVerbose to write specific logs.</li>
</ul>
<p>This is not enough by default as we didn&rsquo;t configure what kind of tracing data we are interested in and how we would like to output them. We need to configure <a href="http://msdn.microsoft.com/en-us/library/3at424ac.aspx">Trace Switches</a> and <a href="http://msdn.microsoft.com/en-us/library/4y5y10s7.aspx">Trace Listeners</a> properly to instrument our application. Microsoft.Owin is root switch and if we enable it, we will see all messages we write through our loggers. The following configuration will enable the Microsoft.Owin switch:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">configuration</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;</span><span style="color: #a31515;">system.diagnostics</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">switches</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">add</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Microsoft.Owin</span><span style="color: black;">"</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Verbose</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">switches</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;/</span><span style="color: #a31515;">system.diagnostics</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">configuration</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>When we run our application on Debug mode, we can see that the Output window will show our log data:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Logging-and-Microsoft-Owin_D6C1/image.png"><img height="221" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Logging-and-Microsoft-Owin_D6C1/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>We can take it further and configure different listeners for different switches. The following configuration will enable <a href="http://msdn.microsoft.com/en-us/library/system.diagnostics.textwritertracelistener.aspx">TextWriterTraceListener</a> for LoggingSample.MyCustomMiddleware <a href="http://msdn.microsoft.com/en-us/library/system.diagnostics.sourceswitch.aspx">SourceSwitch</a>:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">configuration</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;</span><span style="color: #a31515;">system.diagnostics</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">switches</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">add</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Microsoft.Owin</span><span style="color: black;">"</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Verbose</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">switches</span><span style="color: blue;">&gt;</span>

    <span style="color: blue;">&lt;</span><span style="color: #a31515;">sharedListeners</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">add</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">console</span><span style="color: black;">"</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">System.Diagnostics.ConsoleTraceListener</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">sharedListeners</span><span style="color: blue;">&gt;</span>

    <span style="color: blue;">&lt;</span><span style="color: #a31515;">trace</span> <span style="color: red;">autoflush</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">true</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
    
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">sources</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">source</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Microsoft.Owin</span><span style="color: black;">"</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">listeners</span><span style="color: blue;">&gt;</span>
          <span style="color: blue;">&lt;</span><span style="color: #a31515;">add</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">console</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
        <span style="color: blue;">&lt;/</span><span style="color: #a31515;">listeners</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;/</span><span style="color: #a31515;">source</span><span style="color: blue;">&gt;</span>

      <span style="color: blue;">&lt;</span><span style="color: #a31515;">source</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">LoggingSample.MyCustomMiddleware</span><span style="color: black;">"</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">listeners</span><span style="color: blue;">&gt;</span>
          <span style="color: blue;">&lt;</span><span style="color: #a31515;">add</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">file</span><span style="color: black;">"</span> 
               <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">System.Diagnostics.TextWriterTraceListener</span><span style="color: black;">"</span> 
               <span style="color: red;">initializeData</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">traces-MyCustomMiddleware.log</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
        <span style="color: blue;">&lt;/</span><span style="color: #a31515;">listeners</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;/</span><span style="color: #a31515;">source</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">sources</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;/</span><span style="color: #a31515;">system.diagnostics</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">configuration</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>Now we can see that our middleware&rsquo;s logs will only be written into traces-MyCustomMiddleware.log file:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Logging-and-Microsoft-Owin_D6C1/image_3.png"><img height="338" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Logging-and-Microsoft-Owin_D6C1/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>I think this post should give you a pretty good understanding of how Katana&rsquo;s logging infrastructure has been set up and works under the hood. I plain on writing a few more posts on this logging infrastructure. So, stay tuned. Also, as I get used to the concepts of OWIN, I started to think that this logging infrastructure is more fitting as an OWIN extension (just like the <a href="http://owin.org/extensions/owin-SendFile-Extension-v0.3.0.htm">SendFile extension</a>) rather than being tied to Katana.</p>
<p>The sample I used here <a href="https://github.com/tugberkugurlu/OwinSamples/tree/master/LoggingSample">also available on GitHub</a> inside my <a href="https://github.com/tugberkugurlu/OwinSamples">OwinSamples</a> repository.</p>