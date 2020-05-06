---
title: Scaling out SignalR with a Redis Backplane and Testing It with IIS Express
abstract: Learn how easy to scale out SignalR with a Redis backplane and simulate
  a local web farm scenario with IIS Express
created_at: 2013-07-02 11:01:00 +0000 UTC
tags:
- ASP.Net
- Redis
- SignalR
slugs:
- scaling-out-signalr-with-a-redis-backplane-and-testing-it-with-iis-express
---

<p>SignalR was built with scale out in mind from day one and they ship some scale out providers such as <a href="http://nuget.org/packages/Microsoft.AspNet.SignalR.Redis">Redis</a>, <a href="http://nuget.org/packages/Microsoft.AspNet.SignalR.SqlServer">SQL Server</a> and <a href="http://nuget.org/packages/Microsoft.AspNet.SignalR.ServiceBus">Windows Azure Service Bus</a>. There is <a href="http://www.asp.net/signalr/overview/performance-and-scaling">a really nice documentation series on this at official ASP.NET SignalR web site</a> and you can find <a href="http://www.asp.net/signalr/overview/performance-and-scaling/scaleout-with-redis">Redis</a>, <a href="http://www.asp.net/signalr/overview/performance-and-scaling/scaleout-with-windows-azure-service-bus">Windows Azure Service Bus</a> and <a href="http://www.asp.net/signalr/overview/performance-and-scaling/scaleout-with-sql-server">SQL Server</a> samples there. In this quick post, I would like to show you how easy is to get SignalR up and running in a scale out scenario with a <a href="http://redis.io/">Redis</a> backplane.</p>
<h3>Sample Chat Application</h3>
<p>First of all, I have a very simple and stupid real-time web application. The source code is also available on GitHub if you are interested in: <a href="https://github.com/tugberkugurlu/SignalRSamples/tree/master/RedisScaleOutSample">RedisScaleOutSample</a>. Guess what it is? Yes, you&rsquo;re right. It&rsquo;s a chat application :) I&rsquo;m using <a href="http://nuget.org/packages/Microsoft.AspNet.SignalR/2.0.0-beta2">SignalR 2.0.0-beta2</a> on this sample and here is how my hub looks like:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> ChatHub : Hub
{
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Send(<span style="color: blue;">string</span> message)
    {
        Clients.All.messageReceived(message);
    }
}</pre>
</div>
</div>
<p>A very simple hub implementation. Now, let&rsquo;s look at the entire HTML and JavaScript code that I have:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">&lt;!</span><span style="color: #a31515;">DOCTYPE</span> <span style="color: red;">html</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;</span><span style="color: #a31515;">html</span> <span style="color: red;">xmlns</span><span style="color: blue;">=</span><span style="color: blue;">"http://www.w3.org/1999/xhtml"</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;</span><span style="color: #a31515;">head</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">title</span><span style="color: blue;">&gt;</span>Chat Sample<span style="color: blue;">&lt;/</span><span style="color: #a31515;">title</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">head</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;</span><span style="color: #a31515;">body</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">input</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: blue;">"text"</span> <span style="color: red;">id</span><span style="color: blue;">=</span><span style="color: blue;">"msg"</span> <span style="color: blue;">/&gt;</span> 
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">button</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: blue;">"button"</span> <span style="color: red;">id</span><span style="color: blue;">=</span><span style="color: blue;">"send"</span><span style="color: blue;">&gt;</span>Send<span style="color: blue;">&lt;/</span><span style="color: #a31515;">button</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">ul</span> <span style="color: red;">id</span><span style="color: blue;">=</span><span style="color: blue;">"messages"</span><span style="color: blue;">&gt;</span><span style="color: blue;">&lt;/</span><span style="color: #a31515;">ul</span><span style="color: blue;">&gt;</span>

    <span style="color: blue;">&lt;</span><span style="color: #a31515;">script</span> <span style="color: red;">src</span><span style="color: blue;">=</span><span style="color: blue;">"/Scripts/jquery-1.6.4.min.js"</span><span style="color: blue;">&gt;</span><span style="color: blue;">&lt;/</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">script</span> <span style="color: red;">src</span><span style="color: blue;">=</span><span style="color: blue;">"/Scripts/jquery.signalR-2.0.0-beta2.min.js"</span><span style="color: blue;">&gt;</span><span style="color: blue;">&lt;/</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">script</span> <span style="color: red;">src</span><span style="color: blue;">=</span><span style="color: blue;">"/signalr/hubs"</span><span style="color: blue;">&gt;</span><span style="color: blue;">&lt;/</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
        (<span style="color: blue;">function</span> () {

            <span style="color: blue;">var</span> chatHub = $.connection.chatHub,
                msgContainer = $(<span style="color: #a31515;">'#messages'</span>);

            chatHub.client.messageReceived = <span style="color: blue;">function</span> (msg) {
                $(<span style="color: #a31515;">'&lt;li&gt;'</span>).text(msg).appendTo(msgContainer);
            };

            $.connection.hub.start().done(<span style="color: blue;">function</span> () {

                $(<span style="color: #a31515;">'#send'</span>).click(<span style="color: blue;">function</span> () {
                    <span style="color: blue;">var</span> msg = $(<span style="color: #a31515;">'#msg'</span>).val();
                    chatHub.server.send(msg);
                });
            });
        }());
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">body</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">html</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>When I run the application, I can see that it works like a charm:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/1.png"><img height="395" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/1_thumb.png" alt="1" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="1" /></a></p>
<p>This&rsquo;s a single machine scenario and if we want to move this application to multiple VMs, a Web Farm or whatever your choice of scaling out your application, you will see that your application misbehaving. The reason is very simple to understand actually. Let&rsquo;s try to understand why.</p>
<h3>Understanding the Need of Backplane</h3>
<p>Assume that you have two VMs for your super chat application: VM-1 and VM-2. The client-a comes to your application and your load balancer routes that request to VM-1. As your SignalR connection will be persisted as long as it can be, you will be connected to VM-1 for any messages you receive (assuming you are not on Long Pooling transport) and send (if you are on Web Sockets). Then, client-b comes to your application and the load balancer routes that request to VM-2 this time. What happens now? Any messages that client-a sends will not be received by client-b because they are on different nodes and SignalR has no idea about any other node except that it&rsquo;s executing on.</p>
<p>To demonstrate this scenario easily in our development environment, I will fire up the same application in different ports through IIS Express with the following script:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">function</span> programfiles<span style="color: gray;">-</span>dir {
    <span style="color: blue;">if</span> (is64bit <span style="color: gray;">-eq</span> <span style="color: orangered;">$true</span>) {
        (Get<span style="color: gray;">-</span>Item <span style="color: #a31515;">"Env:ProgramFiles(x86)"</span>).Value
    } <span style="color: blue;">else</span> {
        (Get<span style="color: gray;">-</span>Item <span style="color: #a31515;">"Env:ProgramFiles"</span>).Value
    }
}

<span style="color: blue;">function</span> is64bit() {
    <span style="color: blue;">return</span> (<span style="color: gray;">[</span><span style="color: teal;">IntPtr</span><span style="color: gray;">]</span><span style="color: gray;">::</span>Size <span style="color: gray;">-eq</span> 8)
}

<span style="color: orangered;">$executingPath</span> <span style="color: gray;">=</span> (Split<span style="color: gray;">-</span>Path <span style="color: gray;">-</span>parent <span style="color: orangered;">$MyInvocation</span>.MyCommand.Definition)
<span style="color: orangered;">$appPPath</span> <span style="color: gray;">=</span> (join<span style="color: gray;">-</span>path <span style="color: orangered;">$executingPath</span> <span style="color: #a31515;">"RedisScaleOutSample"</span>)
<span style="color: orangered;">$iisExpress</span> <span style="color: gray;">=</span> <span style="color: #a31515;">"$(programfiles-dir)\IIS Express\iisexpress.exe"</span>
<span style="color: orangered;">$args1</span> <span style="color: gray;">=</span> <span style="color: #a31515;">"/path:$appPPath /port:9090 /clr:v4.0"</span>
<span style="color: orangered;">$args2</span> <span style="color: gray;">=</span> <span style="color: #a31515;">"/path:$appPPath /port:9091 /clr:v4.0"</span>

start<span style="color: gray;">-</span><span style="color: blue;">process</span> <span style="color: orangered;">$iisExpress</span> <span style="color: orangered;">$args1</span> <span style="color: gray;">-</span>windowstyle Normal
start<span style="color: gray;">-</span><span style="color: blue;">process</span> <span style="color: orangered;">$iisExpress</span> <span style="color: orangered;">$args2</span> <span style="color: gray;">-</span>windowstyle Normal</pre>
</div>
</div>
<p>I&rsquo;m <a href="http://www.iis.net/learn/extensions/using-iis-express/running-iis-express-from-the-command-line">running IIS Express here from the command line</a>&nbsp;and it&rsquo;s a very powerful feature if you ask me. When you execute the following script (which is run.ps1 in my sample application), you will have the chat application running on localhost:9090 and localhost:9091:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/2.png"><img height="419" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/2_thumb.png" alt="2" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="2" /></a></p>
<p>When we try to same scenario now by connecting both endpoints, you will see that it&rsquo;s not working as it should be:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/3.png"><img height="394" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/3_thumb.png" alt="3" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="3" /></a></p>
<p>SignalR makes it really easy to solve this type of problems. In its core architecture, SignalR uses a pub/sub mechanism to broadcast the messages. Every message in SignalR goes through the message bus and by default, SignalR uses <a href="http://msdn.microsoft.com/en-us/library/microsoft.aspnet.signalr.messaging.messagebus(v=vs.111).aspx">Microsoft.AspNet.SignalR.Messaging.MessageBus</a> which implements <a href="http://msdn.microsoft.com/en-us/library/microsoft.aspnet.signalr.messaging.imessagebus(v=vs.100).aspx">IMessageBus</a> as its in-memory message bus. However, this&rsquo;s fully replaceable and it&rsquo;s where you need to plug your own message bus implementation for your scale out scenarios. SignalR team provides bunch of backplanes for you to work with but if you can totally implement your own if none of the scale-out providers that SignalR team is providing is not enough for you. For instance, the community has a <a href="http://www.rabbitmq.com/">RabbitMQ</a> message bus implementation for SignalR: <a href="https://github.com/mdevilliers/SignalR.RabbitMq">SignalR.RabbitMq</a>.</p>
<h3>Hooking up Redis Backplane to Your SignalR Application</h3>
<p>In order to test configure using Redis as the backplane for SignalR, we need to have a Redis server up and running somewhere. The Redis project does not directly support Windows but <a href="http://msopentech.com/">Microsoft Open Tech</a> provides the <a href="https://github.com/MSOpenTech/redis">Redis Windows port</a> which targets both x86 and x64 bit architectures. The better news is that they distribute the binaries through NuGet: <a href="http://nuget.org/packages/Redis-64">http://nuget.org/packages/Redis-64</a>.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/4.png"><img height="240" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/4_thumb.png" alt="4" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="4" /></a></p>
<p>Now I have Redis binaries, I can get the Redis server up. For our demonstration purposes, running the redis-server.exe without any arguments with the default configuration should be enough:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/5.png"><img height="422" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/5_thumb.png" alt="5" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="5" /></a></p>
<p>The Redis server is running on port 6379 now and we can configure SignalR to use Redis as its backplane. First thing to do is to install the <a href="http://nuget.org/packages/Microsoft.AspNet.SignalR.Redis">SignalR Redis Messaging Backplane NuGet package</a>. As I&rsquo;m using the SignalR 2.0.0-beta2, I will install the <a href="http://nuget.org/packages/Microsoft.AspNet.SignalR.Redis/2.0.0-beta2">version 2.0.0-beta2 of Microsoft.AspNet.SignalR.Redis package</a>.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/6.png"><img height="260" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/6_thumb.png" alt="6" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="6" /></a></p>
<p>Last thing to do is to write a one line of code to replace the IMessageBus implementation:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Startup
{
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Configuration(IAppBuilder app)
    {
        GlobalHost.DependencyResolver
            .UseRedis(<span style="color: #a31515;">"localhost"</span>, 6379, <span style="color: blue;">string</span>.Empty, <span style="color: #a31515;">"myApp"</span>);

        app.MapHubs();
    }
}</pre>
</div>
</div>
<p>The parameters we are passing into the <a href="http://msdn.microsoft.com/en-us/library/jj907714(v=vs.111).aspx">UseRedis</a> method are related to your Redis server. For our case here, we don&rsquo;t have any password and that&rsquo;s why we passed string.Empty. Now, let&rsquo;s compile the application and run the same PowerShell script now to stand up two endpoints which simulates a web farm scenario in your development environment. When we navigate the both endpoints, we will see that messages are broadcasted to all nodes no matter which node they arrive:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/7.png"><img height="419" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Scaling-out-SignalR-with-a-Redis-Backpla_EA5C/7_thumb.png" alt="7" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-width: 0px;" title="7" /></a></p>
<p>That was insanely easy to implement, isn&rsquo;t it.</p>
<h3>A Few Things to Keep in Mind</h3>
<p>The purpose of the SignalR&rsquo;s backplane approach is to enable you to serve more clients in cases where one server is becoming your bottleneck. As you can imagine, having a backplane for your SignalR application can affect the message throughput as your messages need to go through the backplane first and distributed from there to all subscribers. For high-frequency real-time applications, such as real-time games, a backplane is not recommended. For those cases, cleverer load balancers are what you would want. <a href="https://twitter.com/DamianEdwards">Damian Edwards</a> has talked about SignalR and different scale out cases <a href="http://channel9.msdn.com/Events/Build/2013/3-502">on his Build 2013 talk</a> and I strongly recommend you to check that out if you are interested in.</p>
<div><iframe scrolling="no" frameborder="0" src="http://channel9.msdn.com/Events/Build/2013/3-502/player?w=560&amp;h=320" style="height: 320px; width: 560px;"></iframe></div>