---
id: 35028787-9eb7-41e7-92c9-8b0c6661c33c
title: 'Local IIS 7.0 - CS0016: Could not write to output file / Microsoft.Net > Framework
  > v4.0.30319 > Temporary ASP.NET Files'
abstract: Solution to an annoying error message! You are getting 'Could not write
  to output file 'c:\Windows\Microsoft.NET\Framework\....' message? You are at the
  right place.
created_at: 2011-03-05 08:18:00 +0000 UTC
tags:
- .NET
- ASP.Net
- Deployment
- IIS
- Tips
slugs:
- local-iis-7-0-cs0016--could-not-write-to-output-file-microsoft-net-framework-v4-0-30319-temporary-asp-net-files
---

<p>This week I went nuts over my local IIS. I have never swore to a machine that much in my whole life. I am sure of that! The problem is not that big and probably not worth to be written on a blog post I am going to write it anyway because the solution was hard to find on the internet. Maybe this post will help you to fix the problem as I did and you will stop swearing to you machine as I did <img height="19" width="19" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/CS0016-Could-not-write-to-output-file_10928/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" /></p>
<p>Let&rsquo;s get to the point. I am no IIS guy! Seriously!&nbsp; The so called <strong>Cassini</strong> <em>(the tiny web server which pops up when you run a web application on Visual Studio)</em> was so enough for me for over 2 years. But no more enough. I figured that IIS can not be ignored by me anymore. How can I get that point? That&rsquo;s not the issue here. The issue is that I tried to run an ASP.NET MVC 3.0 application under my local IIS 7.0 and got a very annoying error. Which is;</p>
<blockquote>
<p><em><b>Compiler Error Message: </b>CS0016: Could not write to output file 'c:\Windows\Microsoft.NET\Framework\v4.0.30319\Temporary ASP.NET Files\root\62d43c41\27d749ca\App_Code.7lodcznm.dll' &ndash; 'Access denied.'</em></p>
</blockquote>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/CS0016-Could-not-write-to-output-file_10928/image.png"><img height="281" width="220" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/CS0016-Could-not-write-to-output-file_10928/image_thumb.png" align="left" alt="image" border="0" title="image" style="background-image: none; margin: 0px 20px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" /></a>Temporary ASP.NET Files folder <em>(c:\Windows\Microsoft.NET\Framework\v4.0.30319\Temporary ASP.NET Files)</em> was the the one with the problem here. First, I thought that the problem is related to security permissions on the folder and I was right.</p>
<p>I right clicked on the <em>Temporary ASP.NET Files</em> folder and go to the security tab. I noticed that there is user called <em>IIS_IUSRS</em> and that guy has the full control permission. But apparently that was not enough.</p>
<p>The <em><strong>Temporary ASP.NET Files</strong>&nbsp; </em>and <em><strong>C:\Windows\temp</strong></em> folders should have <strong>IIS_WPG</strong> and <strong>NETWORK SERVICE</strong> users with the <strong>full control permission</strong>. I have no idea why <em>C:\Windows\temp</em> folder needs that but I have no effort left to try to find that. Instead, I am writing a blog post about the problem. Maybe latter I will get to that and figure it out, too <img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/CS0016-Could-not-write-to-output-file_10928/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /></p>
<p>Once you applied those setting, restart your IIS and try to run your application again. The error should be gone by now.</p>
<p>I suffered a lot by trying to find the right method for the problem and I hope you didn&rsquo;t have to go through hell over this.</p>
<p>Hope this helps <img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/CS0016-Could-not-write-to-output-file_10928/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /></p>