---
id: 2d9788cc-922f-4fd0-af87-73dadffa1657
title: Disposing Resources At the End of the Request Lifecycle in ASP.NET Web API
abstract: How to disposing resources at the end of the request lifecycle in ASP.NET
  Web API with the RegisterForDispose extension method for the HttpRequestMessage
  class
created_at: 2012-07-27 05:26:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
- C#
slugs:
- disposing-resources-at-the-end-of-the-request-lifecycle-in-asp-net-web-api
---

<p>When we are working with ASP.NET Web API, we always deal with disposable resources and if we are using an IoC container, it deals with disposals in behalf of us.</p>
<p>However, when we are inside the extensibility points and we just want to create a resource there and dispose it as soon as we are done with it, it can be easily overwhelming because the asynchronous structure of the ASP.NET Web API. We should be very careful where to dispose our resources.</p>
<p>However, ASP.NET team has come up with this great idea that they allow us to register resources to be disposed by a host once the request is disposed. The <a title="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessageextensions.registerfordispose(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessageextensions.registerfordispose(v=vs.108).aspx">RegisterForDispose</a> extension method for the <a title="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage(v=vs.110).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage(v=vs.110).aspx">HttpRequestMessage</a> class enables this feature. In fact, all it does is to keep track of all the registered resources inside a List&lt;IDisposable&gt; object. They stick this collection into the <a title="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage.properties(v=vs.110)" href="http://msdn.microsoft.com/en-us/library/system.net.http.httprequestmessage.properties(v=vs.110)">Properties</a> property of the HttpRequestMessage (under the key of <strong>MS_DisposableRequestResources</strong>). Then, at the end of the request, these resources will be disposed along with the request.</p>
<p>The usage is also so simple. For example, the below code is for a message handler which serves as a timeout watcher.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> TimeoutHandler : DelegatingHandler {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">int</span> _milliseconds;
    <span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">readonly</span> TimerCallback s_timerCallback = 
        <span style="color: blue;">new</span> TimerCallback(TimerCallbackLogic);

    <span style="color: blue;">public</span> TimeoutHandler(<span style="color: blue;">int</span> milliseconds) {

        <span style="color: blue;">if</span> (milliseconds &lt; -1) {

            <span style="color: blue;">throw</span> <span style="color: blue;">new</span> ArgumentOutOfRangeException(<span style="color: #a31515;">"milliseconds"</span>);
        }

        _milliseconds = milliseconds;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Timeout {

        <span style="color: blue;">get</span> {
            <span style="color: blue;">return</span> _milliseconds;
        }
    }

    <span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task&lt;HttpResponseMessage&gt; SendAsync(
        HttpRequestMessage request, CancellationToken cancellationToken) {

        <span style="color: blue;">var</span> cts = <span style="color: blue;">new</span> CancellationTokenSource();
        <span style="color: blue;">var</span> linkedTokenSource = CancellationTokenSource.CreateLinkedTokenSource(
            cts.Token, cancellationToken
        );
        <span style="color: blue;">var</span> linkedToken = linkedTokenSource.Token;
        <span style="color: blue;">var</span> timer = <span style="color: blue;">new</span> Timer(s_timerCallback, cts, -1, -1);

        request.RegisterForDispose(timer);
        request.RegisterForDispose(cts);
        request.RegisterForDispose(linkedTokenSource);

        timer.Change(_milliseconds, -1);

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.SendAsync(request, linkedToken).ContinueWith(task =&gt; {

            <span style="color: blue;">if</span> (task.Status == TaskStatus.Canceled) {

                <span style="color: blue;">return</span> request.CreateResponse(HttpStatusCode.RequestTimeout);
            }

            <span style="color: green;">//TODO: Handle faulted task as well</span>

            <span style="color: blue;">return</span> task.Result;

        }, TaskContinuationOptions.ExecuteSynchronously);
    }

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> TimerCallbackLogic(<span style="color: blue;">object</span> obj) {

        CancellationTokenSource cancellationTokenSource = 
            (CancellationTokenSource)obj;
            
        cancellationTokenSource.Cancel();
    }
}</pre>
</div>
</div>
<p>You may see that there are still TODO notes inside the code which means I am still not done with the code, so don't use it in production!:) but it is a good example to illustrate the usage of RegisterForDispose extension method. Notice that I used two <a title="http://msdn.microsoft.com/en-us/library/system.threading.cancellationtokensource.aspx" href="http://msdn.microsoft.com/en-us/library/system.threading.cancellationtokensource.aspx">CancellationTokenSource</a> objects and one <a title="http://msdn.microsoft.com/en-us/library/system.threading.timer.aspx" href="http://msdn.microsoft.com/en-us/library/system.threading.timer.aspx">Timer</a> object (I need a Timer object here to request a cancellation because CancellationTokenSource does not have CancelAfter method in .NET v4.0 and I couldn&rsquo;t find any other easier way to do this). All of these objects need to be disposed once they are no longer in use. However, it is really hard to keep track to figure out which one of them needs to be disposed where inside this handler. So, I just registered them to disposed at the end of the request lifetime.</p>