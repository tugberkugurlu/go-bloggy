---
id: 78be3dca-bc86-4f51-a7b5-7c901517916a
title: Getting Started with Windows Server AppFabric Caching PowerShell Commands
abstract: I started to use Windows Server AppFabric for its distributed caching feature
  and I wanted to take a note of the useful PowerShell commands to manage the service
  configuration and administration.
created_at: 2012-05-23 05:51:00 +0000 UTC
tags:
- ASP.Net
- PowerShell
- Windows Server AppFabric
slugs:
- getting-started-with-windows-server-appfabric-caching-powershell-commands
---

<p>I started to use <a title="http://msdn.microsoft.com/en-us/windowsserver/ee695849.aspx" href="http://msdn.microsoft.com/en-us/windowsserver/ee695849.aspx">Windows Server AppFabric</a> for its distributed caching feature and I wanted to take a note of the useful PowerShell commands to manage the service configuration and administration. When I write this blog post, <a title="http://www.microsoft.com/en-us/download/details.aspx?id=27115" href="http://www.microsoft.com/en-us/download/details.aspx?id=27115">Microsoft AppFabric 1.1 for Windows Server</a> is available. I installed the caching service and caching management features firstly. Then, I configured it properly. See <a title="http://msdn.microsoft.com/en-us/library/hh351248" href="http://msdn.microsoft.com/en-us/library/hh351248">Configure AppFabric</a> page for more information on configuration.</p>
<p>First of all, import necessary modules:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>Import<span style="color: gray;">-</span>Module DistributedCacheConfiguration
Import<span style="color: gray;">-</span>Module DistributedCacheAdministration</pre>
</div>
</div>
<p>You can view the commands with Get-Command Cmdlet. For example, we can see all commands which are available for DistributedCacheAdministration with the following one line of code.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>Get<span style="color: gray;">-</span>Command <span style="color: gray;">-</span>Module DistributedCacheAdministration</pre>
</div>
</div>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Server-AppFabric-Caching-PowerSh_975D/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Server-AppFabric-Caching-PowerSh_975D/image_thumb.png" width="644" height="413" /></a></p>
<p><strong>Use-CacheCluster</strong> command sets the context of your PowerShell session to a particular cache cluster. <span style="text-decoration: underline;"><strong>Note that you must run this command before using any other Cache Administration commands in PowerShell</strong></span>.</p>
<p><strong>Get-CacheHost</strong> command lists all cache host services that are members of the cache cluster.</p>
<p>Your cache host is running as a windows service and you can retrieve the status of your caching host with Get-CacheHost command. It will output a result which is similar to following:</p>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Server-AppFabric-Caching-PowerSh_975D/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Server-AppFabric-Caching-PowerSh_975D/image_thumb_3.png" width="644" height="413" /></a></p>
<p>My service is up and running but if the status indicates that the service is down, you can start it with the following command:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>Start<span style="color: gray;">-</span>CacheHost <span style="color: gray;">-</span>HostName TugberkWin08R2 <span style="color: gray;">-</span>CachePort 22233</pre>
</div>
</div>
<p>You should be able to see your service running after this:</p>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Server-AppFabric-Caching-PowerSh_975D/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Server-AppFabric-Caching-PowerSh_975D/image_thumb_4.png" width="644" height="413" /></a></p>
<p><strong>Get-CacheClusterInfo</strong> command from configuration module returns the cache cluster information, including details on its initialization status and its size. The following one line of code shows a sample usage.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>Get<span style="color: gray;">-</span>CacheClusterInfo <span style="color: gray;">-</span>Provider XML <span style="color: gray;">-</span>ConnectionString \\TUGBERKWIN08R2\Caching</pre>
</div>
</div>
<p>More information on this command: <a href="http://msdn.microsoft.com/en-us/library/ff631076(v=ws.10).aspx">http://msdn.microsoft.com/en-us/library/ff631076(v=ws.10).aspx</a></p>
<p>That&rsquo;s it for now.</p>
<h3>Resources</h3>
<ul>
<li><a title="http://msdn.microsoft.com/en-us/library/hh351318" href="http://msdn.microsoft.com/en-us/library/hh351318">Microsoft AppFabric 1.1 for Windows</a> 
<ul>
<li><a href="http://msdn.microsoft.com/en-us/library/hh334305">AppFabric Caching Features</a></li>
<li><a href="http://msdn.microsoft.com/en-us/library/hh851388">Caching Powershell Cmdlets (AppFabric 1.1)</a> 
<ul>
<li><a href="http://msdn.microsoft.com/en-us/library/hh848898">Microsoft.ApplicationServer.Caching.Commands</a></li>
<li><a href="http://msdn.microsoft.com/en-us/library/hh848869">Microsoft.ApplicationServer.Caching.Configuration.Commands</a></li>
</ul>
</li>
</ul>
</li>
<li><a title="http://www.devtrends.co.uk/blog/asp.net-mvc-output-caching-with-windows-appfabric-cache" href="http://www.devtrends.co.uk/blog/asp.net-mvc-output-caching-with-windows-appfabric-cache">ASP.NET MVC Output Caching with Windows AppFabric Cache</a></li>
<li><a title="http://www.yazgelistir.com/makale/csharp-ile-appfabric-cache-yonetimi" href="http://www.yazgelistir.com/makale/csharp-ile-appfabric-cache-yonetimi">C# &amp; VB.Net / C# ile AppFabric Cache Y&ouml;netimi (in Turkish)</a></li>
</ul>