---
title: Efficiently Streaming Large HTTP Responses With HttpClient
abstract: Downloading large files with HttpClient and you see that it takes lots of
  memory space? This post is probably for you. Let's see how to efficiently streaming
  large HTTP responses with HttpClient.
created_at: 2014-05-11 12:56:00 +0000 UTC
tags:
- .net
- ASP.NET Web API
- HTTP
slugs:
- efficiently-streaming-large-http-responses-with-httpclient
---

<p>I see common scenarios where people need to download large files (images, PDF files, etc.) on their .NET projects. What I mean by large files here is probably not what you think. It should be enough to call it large if it’s 500 KB as you will hit a memory limit once you try to download lots of files concurrently in a wrong way as below:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">static</span> async Task HttpGetForLargeFileInWrongWay()
{
    <span style="color: blue">using</span> (HttpClient client = <span style="color: blue">new</span> HttpClient())
    {
        <span style="color: blue">const</span> <span style="color: blue">string</span> url = <span style="color: #a31515">"https://github.com/tugberkugurlu/ASPNETWebAPISamples/archive/master.zip"</span>;
        <span style="color: blue">using</span> (HttpResponseMessage response = await client.GetAsync(url))
        <span style="color: blue">using</span> (Stream streamToReadFrom = await response.Content.ReadAsStreamAsync())
        {
            <span style="color: blue">string</span> fileToWriteTo = Path.GetTempFileName();
            <span style="color: blue">using</span> (Stream streamToWriteTo = File.Open(fileToWriteTo, FileMode.Create))
            {
                await streamToReadFrom.CopyToAsync(streamToWriteTo);
            }

            response.Content = <span style="color: blue">null</span>;
        }
    }
}</pre></div></div>
<p>By calling <a href="http://msdn.microsoft.com/en-us/library/hh158944(v=vs.118).aspx">GetAsync</a> method directly there, we are loading every single byte into memory. You can see this happening in a simple way by opening the <a href="http://windows.microsoft.com/en-us/windows7/open-task-manager">Task Manager</a> and observing the memory of the process.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6a12db4f-7633-467e-80db-e3fb3d789163.gif"><img title="2" style="display: inline" alt="2" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/41ad073a-aeae-432b-91d1-f133de0c58d7.gif" width="640" height="342"></a></p>
<p>We are calling <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httpcontent.readasstreamasync">ReadAsStreamAsync</a> on <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httpcontent">HttpContent</a> after the GetAsync method is completed. This will just get us the <a href="http://msdn.microsoft.com/en-us/library/system.io.memorystream.aspx">MemoryStream</a>, so there is no point there:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/aaa80cca-2091-45ef-b31a-009706ca98b2.png"><img title="Screenshot 2014-05-11 15.18.14" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Screenshot 2014-05-11 15.18.14" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/f65e5357-5536-452e-804d-0bf4ef961f7e.png" width="644" height="158"></a></p>
<p>We need a way not to load the response body into memory and have the raw network stream so that we can pass the bytes into another stream without hitting the memory too hard. We can do it by just reading the headers of the response and then getting a handle for the network stream as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">static</span> async Task HttpGetForLargeFileInRightWay()
{
    <span style="color: blue">using</span> (HttpClient client = <span style="color: blue">new</span> HttpClient())
    {
        <span style="color: blue">const</span> <span style="color: blue">string</span> url = <span style="color: #a31515">"https://github.com/tugberkugurlu/ASPNETWebAPISamples/archive/master.zip"</span>;
        <span style="color: blue">using</span> (HttpResponseMessage response = await client.GetAsync(url, HttpCompletionOption.ResponseHeadersRead))
        <span style="color: blue">using</span> (Stream streamToReadFrom = await response.Content.ReadAsStreamAsync())
        {
            <span style="color: blue">string</span> fileToWriteTo = Path.GetTempFileName();
            <span style="color: blue">using</span> (Stream streamToWriteTo = File.Open(fileToWriteTo, FileMode.Create))
            {
                await streamToReadFrom.CopyToAsync(streamToWriteTo);
            }
        }
    }
}</pre></div></div>
<p>Notice that we are calling <a href="http://msdn.microsoft.com/en-us/library/hh551757(v=vs.118).aspx">another overload of the GetAsync</a> method by passing the <a href="http://msdn.microsoft.com/en-us/library/system.net.http.httpcompletionoption.aspx">HttpCompletionOption</a> enumeration value as ResponseHeadersRead. This switch tells the HttpClient not to buffer the response. In other words, it will just read the headers and return the control back. This means that the HttpContent is not ready at the time when you get the control back. Afterwards, we are getting the stream and calling the <a href="http://msdn.microsoft.com/en-us/library/hh159084(v=vs.110).aspx">CopyToAsync</a> method on it by passing our <a href="http://msdn.microsoft.com/en-us/library/system.io.filestream.aspx">FileStream</a>. The result is much better:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/8dd27967-5f41-46b3-882c-c1be21d4658c.gif"><img title="3" style="display: inline" alt="3" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a4f9e338-947c-4afc-8d69-328c5a05672c.gif" width="640" height="342"></a></p>
<h3>Resources</h3>
<ul>
<li><a href="http://www.tugberkugurlu.com/archive/streaming-with-newnet-httpclient-and-httpcompletionoption-responseheadersread">Streaming with New .NET HttpClient and HttpCompletionOption.ResponseHeadersRead</a></li>
<li><a href="http://stackoverflow.com/questions/12533533/async-reading-chunked-content-with-httpclient-from-asp-net-webapi">Async reading chunked content with HttpClient from ASP.NET WebApi</a></li></ul>  