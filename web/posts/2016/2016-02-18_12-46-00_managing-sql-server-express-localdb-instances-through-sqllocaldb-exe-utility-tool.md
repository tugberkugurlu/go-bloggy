---
id: 78065935-a7ab-48c4-9d4a-71c10a56b21a
title: Managing SQL Server Express LocalDB Instances Through SqlLocalDB.exe Utility
  Tool
abstract: I love the feeling when I discover a tiny, hidden tool and SqlLocalDB.exe,
  a management utility tool for Microsoft SQL Server LocalDB which allows you to manage
  the LocalDB instances on your machine, is one of them. Let me show you what it is.
created_at: 2016-02-18 12:46:00 +0000 UTC
tags:
- SQL Server
slugs:
- managing-sql-server-express-localdb-instances-through-sqllocaldb-exe-utility-tool
---

<p>I love the feeling when I discover a tiny, hidden tool which I can put into my daily software toolbox. I started to sense this feeling more and more lately with some amazing command line tools and I want to write about those here (famous last words) like I did for <a href="http://www.tugberkugurlu.com/archive/quickly-hosting-static-files-in-your-development-environment-with-node-http-server">http-server</a> and <a href="http://www.tugberkugurlu.com/archive/using-azure-storage-emulator-command-line-tool-wastorageemulator-exe">WAStorageEmulator.exe</a> a while back. Today, I want to start this by writing about <a href="https://msdn.microsoft.com/en-gb/library/hh247716.aspx">SqlLocalDB.exe</a>, a management utility tool for <a href="https://msdn.microsoft.com/en-gb/library/hh510202.aspx">Microsoft SQL Server Express LocalDB</a> which allows you to manage the LocalDB instances on your machine such as creating, starting and stopping them. This is really handy if you are after creating lightweight SQL Server instances for temporary processing like we do as part of some <a href="https://www.red-gate.com/products/dlm/dlm-automation-suite/">DLM Automation Suite</a> <a href="https://documentation.red-gate.com/display/SR1/Cmdlet+reference">PowerShell cmdlets</a> (e.g. <a href="https://documentation.red-gate.com/display/SR1/Invoke-DlmDatabaseSchemaValidation">Invoke-DlmDatabaseSchemaValidation</a>).</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/cd54e987-117b-4037-81f8-7a6dc398a86c.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c2abdc0e-45fc-40d7-a16e-9a00abe813fb.png" width="644" height="464"></a></p> <blockquote> <p>I am not entirely sure hot you would install this tool and I failed to find out exactly how. However, it seems like this comes with the LocalDB installation and you can acquire that through <a href="https://www.microsoft.com/en-us/download/details.aspx?id=42299">here</a> (also check out <a href="https://twitter.com/shanselman">Scott Hanselman</a>'s ironic "<a href="http://www.hanselman.com/blog/DownloadSqlServerExpress.aspx">Download SQL Server Express</a>" blog post).</p></blockquote> <p>This may depend on the version you installed but I can locate the SqlLocalDB.exe under "<em>C:\Program Files\Microsoft SQL Server\120\Tools\Binn</em>" on my machine and after that, it's just executing the commands. For example, I can execute the info command to see the LocalDB instances I have in my machine:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>C:\Program Files\Microsoft SQL Server\120\Tools\Binn<span style="color: gray">&gt;</span>SqlLocalDB.exe info

MSSQLLocalDB
RedGateTemp</pre></div></div>
<p>I can also list which versions of LocalDB I have installed on my machine:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>C:\Program Files\Microsoft SQL Server\120\Tools\Binn<span style="color: gray">&gt;</span>SqlLocalDB.exe versions

Microsoft SQL Server 2014 (12.0.2000.8)</pre></div></div>
<p>Let's create our own instance through the create command:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>C:\Program Files\Microsoft SQL Server\120\Tools\Binn<span style="color: gray">&gt;</span>SqlLocalDB.exe create tugberk

LocalDB instance "tugberk" created with version 12.0.2000.8.</pre></div></div>
<p>You can view the status of an instance through the info command by passing the instance name as an argument:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>C:\Program Files\Microsoft SQL Server\120\Tools\Binn<span style="color: gray">&gt;</span>SqlLocalDB.exe info tugberk

Name:               tugberk
Version:            12.0.2000.8
Shared name:
Owner:              TUGBERKPC\Tugberk
Auto<span style="color: gray">-</span>create:        No
State:              Stopped
Last start time:    2<span style="color: gray">/</span>18<span style="color: gray">/</span>2016 12:06:00 PM
Instance pipe name:</pre></div></div>
<p>Let's start the instance through the start command:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>C:\Program Files\Microsoft SQL Server\120\Tools\Binn<span style="color: gray">&gt;</span>SqlLocalDB.exe start tugberk

LocalDB instance "tugberk" started.</pre></div></div>
<p>We can see that it has started:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>C:\Program Files\Microsoft SQL Server\120\Tools\Binn<span style="color: gray">&gt;</span>SqlLocalDB.exe info tugberk

Name:               tugberk
Version:            12.0.2000.8
Shared name:
Owner:              TUGBERKPC\Tugberk
Auto<span style="color: gray">-</span>create:        No
State:              Running
Last start time:    2<span style="color: gray">/</span>18<span style="color: gray">/</span>2016 12:09:28 PM
Instance pipe name: np:\\.\pipe\LOCALDB#7F6D2993\tsql\query</pre></div></div>
<p>Great, now I can connect to this instance using the provided magical, special instance pipe name:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/58bb58a5-2014-41ee-aa02-af16f2fef81e.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ade9194e-186c-43f8-a7df-35a220488f84.png" width="331" height="484"></a></p>
<p>Or, LocalDB server name by prefixing the instance name with "(localdb)\":</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e88ee341-486c-441a-bc2f-7e86ae6d54b9.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/bf5410c8-51df-436d-9ba7-5f72a9d082b7.png" width="472" height="484"></a></p>
<blockquote>
<p>If you dig a little deeper, you will see that an instance of "sqlservr.exe" has been started for your LocalDB instance:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/59747e1a-7cce-4321-a47d-b288eab0ec5a.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/75b3bc37-051f-4e93-a25e-3bd0bfba0972.png" width="644" height="78"></a></p>
<p>You can also see that all LocalDB data is stored under "<em>%LOCALAPPDATA%\Microsoft\Microsoft SQL Server Local DB\Instances</em>".</p></blockquote>
<p>From there, you can treat this instance as a usual SQL Server instance and perform whatever operation you want to perform on it. Be aware that LocalDB has the same limitations as SQL Server Express.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/10daaf3d-dfe1-445e-aed6-3f9f52ae4c5c.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/449f3ed6-ed06-4bb1-9a25-146f9b39ad25.png" width="244" height="167"></a></p>
<p>Finally, you can stop the instance through the stop command:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>C:\Program Files\Microsoft SQL Server\120\Tools\Binn<span style="color: gray">&gt;</span>SqlLocalDB.exe stop tugberk

LocalDB instance "tugberk" stopped.</pre></div></div>
<h3>Further Information</h3>
<ul>
<li><a href="https://blogs.msdn.microsoft.com/sqlexpress/2011/07/12/introducing-localdb-an-improved-sql-express/">Introducing LocalDB, an improved SQL Express</a></li>
<li><a href="https://msdn.microsoft.com/en-gb/library/hh212961.aspx">SqlLocalDB Utility</a></li>
<li><a href="http://blogs.msdn.com/b/jerrynixon/archive/2012/02/26/sql-express-v-localdb-v-sql-compact-edition.aspx">SQL Express v LocalDB v SQL Compact Edition</a></li></ul>  