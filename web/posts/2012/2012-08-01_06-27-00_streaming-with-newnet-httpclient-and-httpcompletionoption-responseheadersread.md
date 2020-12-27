---
id: 7762bb40-b07a-4d07-84fa-d3571e473a4f
title: Streaming with New .NET HttpClient and HttpCompletionOption.ResponseHeadersRead
abstract: How to consume a streaming endpoint with new .NET System.Net.Http.HttpClient
  and the role of HttpCompletionOption.ResponseHeadersRead
created_at: 2012-08-01 06:27:00 +0000 UTC
tags:
- .NET
- ASP.NET Web API
- C#
slugs:
- streaming-with-newnet-httpclient-and-httpcompletionoption-responseheadersread
---

<p>For a while now, I have been playing with a new application I created for fun: <a title="http://www.tweetmapr.net/" href="http://www.tweetmapr.net/">TweetMapR</a> and I am not sure where it is going :) The application itself consumes <a title="https://dev.twitter.com/docs/streaming-apis" href="https://dev.twitter.com/docs/streaming-apis">Twitter Streaming API</a> and broadcasts the retrieved data to the connected clients through <a title="https://github.com/SignalR/SignalR" href="https://github.com/SignalR/SignalR">SignalR</a>. I used new <a title="http://msdn.microsoft.com/en-us/library/system.net.http.httpclient(v=vs.110)" href="http://msdn.microsoft.com/en-us/library/system.net.http.httpclient(v=vs.110)">HttpClient</a> to connect to Twitter Streaming API and I used OAuth as authentication protocol. In fact, I created my own <a title="https://github.com/WebAPIDoodle/TwitterDoodle" href="https://github.com/WebAPIDoodle/TwitterDoodle">OAuth Twitter client</a> to connect to Twitter which is so raw right now but works well.</p>
<p>But how to send a request to Twitter Streaming API and keep it open infinitely was my main question for the whole time. So, I brought up the <a title="http://wiki.sharpdevelop.net/ILSpy.ashx" href="http://wiki.sharpdevelop.net/ILSpy.ashx">ILSpy</a> and decompiled the System.Net.Http.dll to see what is really going on and I realized that it was going to be so easy. First of all, if your request is a GET request, you are covered by the framework. The <a title="http://msdn.microsoft.com/en-us/library/hh551756(v=vs.110)" href="http://msdn.microsoft.com/en-us/library/hh551756(v=vs.110)">GetStreamAsync</a> method of the HttpClient sends a GET request and returns back the stream as soon as it completes reading the response headers. The following example code shows the usage.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">using</span> (HttpClient httpClient = <span style="color: blue;">new</span> HttpClient()) {

    httpClient.Timeout = TimeSpan.FromMilliseconds(Timeout.Infinite);
    <span style="color: blue;">var</span> requestUri = <span style="color: #a31515;">"http://localhost:6797"</span>;
    <span style="color: blue;">var</span> stream = httpClient.GetStreamAsync(requestUri).Result;

    <span style="color: blue;">using</span> (<span style="color: blue;">var</span> reader = <span style="color: blue;">new</span> StreamReader(stream)) {

        <span style="color: blue;">while</span> (!reader.EndOfStream) { 

            <span style="color: green;">//We are ready to read the stream</span>
            <span style="color: blue;">var</span> currentLine = reader.ReadLine();
        }
    }
}</pre>
</div>
</div>
<p>First of all, don&rsquo;t use Result on Task :) It blocks but I used it here for the sake of simplicity to stick with the main point of this post :) Assuming that the localhost:6797 is a streaming endpoint here, the above code should work perfectly fine and shouldn&rsquo;t be timed out as we set the Timeout to <a title="http://msdn.microsoft.com/en-us/library/system.threading.timeout.infinite.aspx" href="http://msdn.microsoft.com/en-us/library/system.threading.timeout.infinite.aspx">System.Threading.Timeout.Infinite</a>. But what if the request we need to send is a POST request? In that case, we have a little more work to do here.</p>
<p>If we just try to send a POST request with the&nbsp;<a title="http://msdn.microsoft.com/en-us/library/hh138245(v=vs.110)" href="http://msdn.microsoft.com/en-us/library/hh138245(v=vs.110)">PostAsync</a> method or any of its variants, the response will never come back unless the server ends the request because it will try to read the response till the end. In order to omit this problem, we can pass a <a title="http://msdn.microsoft.com/en-us/library/system.net.http.httpcompletionoption(v=vs.110).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.httpcompletionoption(v=vs.110).aspx">HttpCompletionOption</a> enumuration value to specify the completion option. The HttpCompletionOption enumeration type has two members and one of them is ResponseHeadersRead which tells the HttpClient to only read the headers and then return back the result immediately. The following code shows a sample example where we need to send a form-urlencoded POST request to a streaming endpoint.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">using</span> (HttpClient httpClient = <span style="color: blue;">new</span> HttpClient()) {

    httpClient.Timeout = TimeSpan.FromMilliseconds(Timeout.Infinite);
    <span style="color: blue;">var</span> requestUri = <span style="color: #a31515;">"http://localhost:6797"</span>;

    <span style="color: blue;">var</span> formUrlEncodedContent = <span style="color: blue;">new</span> FormUrlEncodedContent(
        <span style="color: blue;">new</span> List&lt;KeyValuePair&lt;<span style="color: blue;">string</span>, <span style="color: blue;">string</span>&gt;&gt;() { 
            <span style="color: blue;">new</span> KeyValuePair&lt;<span style="color: blue;">string</span>, <span style="color: blue;">string</span>&gt;(<span style="color: #a31515;">"userId"</span>, <span style="color: #a31515;">"1000"</span>) });

    formUrlEncodedContent.Headers.ContentType = 
        <span style="color: blue;">new</span> MediaTypeHeaderValue(<span style="color: #a31515;">"application/x-www-form-urlencoded"</span>);

    <span style="color: blue;">var</span> request = <span style="color: blue;">new</span> HttpRequestMessage(HttpMethod.Post, requestUri);
    request.Content = formUrlEncodedContent;

    <span style="color: blue;">var</span> response = httpClient.SendAsync(
        request, HttpCompletionOption.ResponseHeadersRead).Result;
    <span style="color: blue;">var</span> stream = response.Content.ReadAsStreamAsync().Result;

    <span style="color: blue;">using</span> (<span style="color: blue;">var</span> reader = <span style="color: blue;">new</span> StreamReader(stream)) {

        <span style="color: blue;">while</span> (!reader.EndOfStream) { 

            <span style="color: green;">//We are ready to read the stream</span>
            <span style="color: blue;">var</span> currentLine = reader.ReadLine();
        }
    }
}</pre>
</div>
</div>
<p>Again, don&rsquo;t use Result as I do here :) We have obviously more noise this time but if this is something that you will use often, there is nothing stopping you to write an extension method.</p>