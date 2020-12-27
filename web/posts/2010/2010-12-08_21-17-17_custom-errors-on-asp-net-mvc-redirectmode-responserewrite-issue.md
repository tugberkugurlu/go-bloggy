---
id: 48011de8-13fb-4dc8-81b0-d2b67c0f6277
title: Custom Errors on ASP.Net MVC - redirectMode="ResponseRewrite" Issue
abstract: I assume that some of you folks have tried that in your ASP.Net MVC applications
  and try to figure out why it doesn't work. Well, I have figured it out...
created_at: 2010-12-08 21:17:17 +0000 UTC
tags:
- .NET
- ASP.Net
- ASP.NET MVC
- Web.Config
slugs:
- custom-errors-on-asp-net-mvc-redirectmode-responserewrite-issue
---

<p>Today, I was wraping up an asp.net mvc project. Make it prettier&nbsp;and safer. I realised that there is a thing which doesn't quite work with asp.net mvc. CustomErrors !</p>
<p>I wanted to refill the page, which has one or multiple errors, with a custom error page so I implemented the following code on web.config file;</p>
<p><img style="border: 1px solid gray;" title="UploadedByAuthors/customErrors-redirectMode-ResponseRewrite-web-config.PNG" alt="UploadedByAuthors/customErrors-redirectMode-ResponseRewrite-web-config.PNG" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/customErrors-redirectMode-ResponseRewrite-web-config.PNG" /></p>
<p>&nbsp;</p>
<p>Then, I hit an error on pupose just in case to see if it works or not and <strong>Boom...</strong> !! It failed ! It gave me the <em>famous ASP.Net yellow screen of death</em>;</p>
<p><img title="yellow-screen-of-death-asp.net.PNG" alt="yellow-screen-of-death-asp.net.PNG" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/yellow-screen-of-death-asp.net.PNG" /></p>
<p>I was a little surprised about that and I wonder why that thing happened. So I made a little research and I opened a thread on <a target="_blank" title="http://forums.asp.net/1146.aspx" href="http://forums.asp.net/1146.aspx">ASP.Net MVC Forums</a>&nbsp;</p>
<p>I found out that MVC Routes are not compatible&nbsp;with ResponseRewrite. <a target="_blank" title="http://stackoverflow.com/questions/781861/customerrors-does-not-work-when-setting-redirectmoderesponserewrite" href="http://stackoverflow.com/questions/781861/customerrors-does-not-work-when-setting-redirectmoderesponserewrite">A smilar thread</a> was opened on <a target="_blank" title="http://stackoverflow.com" href="http://stackoverflow.com">Stackoverflow.com</a>&nbsp;and the answer is there as appears below;</p>
<blockquote>
<p>It is important to note for anyone trying to do this in an MVC application that ResponseRewrite uses Server.Transfer behind the scenes. Therefore, the defaultRedirect must correspond to a legitimate file on the file system. Apparently, Server.Transfer is not compatible with MVC routes, therefore, if your error page is served by a controller action, Server.Transfer is going to look for /Error/Whatever, not find it on the file system, and return a generic 404 error page!</p>
</blockquote>
<p>The answer is pretty reasonable for me and I changed the RedirectMode to ResponseRedirect which is the default one.</p>
<p>But I still wonder that ASP.Net Mvc team will fix it in next versions or not....</p>