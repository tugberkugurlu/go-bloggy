---
title: How to Detect Errors of Our ASP.NET MVC Views on Compile Time - Blow up In
  My Face Theory
abstract: We will see How to detect errors of our ASP.NET MVC views on compile time
  inside Visual Studio.
created_at: 2011-09-29 05:04:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- Tips
- Visual Studio
slugs:
- how-to-detect-errors-of-our-asp-net-mvc-views-on-compile-time-blow-up-in-my-face-theory
---

<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/head-in-a-cake.png"><img style="background-image: none; margin: 0px 15px 15px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border: 0px;" title="head-in-a-cake" border="0" alt="head-in-a-cake" align="left" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/head-in-a-cake_thumb.png" width="244" height="176" /></a></p>
<p>Yesterday, I saw a question on <a title="http://stackoverflow.com" href="http://stackoverflow.com" target="_blank">stackoverflow.com</a> asking why Visual Studio doesn&rsquo;t know that the app is going to fail when there is an error on your code inside your views. This is a good question and it brings up an philosophical question :</p>
<p><strong><em>Do we trust compile time check?</em></strong></p>
<p>In my opinion, no. If this was the case, there would be no point for TDD, even Unit Testing.</p>
<p>Compile time check is no more useful than your Microsoft Word&rsquo;s spell checker. It helps a lot but it is basically a spell checker.</p>
<p>&nbsp;</p>
<blockquote>
<p>In this blog post, I am trying a new way of blogging which I just learnt from <a title="http://haacked.com" href="http://haacked.com" target="_blank">Phil Haack</a>&nbsp;<img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/wlEmoticon-smile.png" /> Put unrelated photo to your blog post. This approach is among <a title="http://haacked.com/archive/2011/01/02/top-ten-blogging-cliches.aspx" href="http://haacked.com/archive/2011/01/02/top-ten-blogging-cliches.aspx" target="_blank">Top 10 Blogging Clich&eacute;s of 2010</a>. It was a not-to-do but here I am doing it.</p>
</blockquote>
<p>The basic problem here is that Visual Studio is not even useful as a spell checker. How so? Let me show you an example.</p>
<blockquote><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/works-on-my-machine-seal-of-approval.png"><img style="background-image: none; margin: 0px 0px 15px 15px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border: 0px;" title="works-on-my-machine-seal-of-approval" border="0" alt="works-on-my-machine-seal-of-approval" align="right" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/works-on-my-machine-seal-of-approval_thumb.png" width="200" height="193" /></a>
<p><strong>NOTE</strong></p>
<p>Couple of days ago, I saw something cool from <a title="http://www.hanselman.com/blog/GuideToInstallingAndBootingWindows8DeveloperPreviewOffAVHDVirtualHardDisk.aspx" href="http://www.hanselman.com/blog/GuideToInstallingAndBootingWindows8DeveloperPreviewOffAVHDVirtualHardDisk.aspx" target="_blank">one of Scott Hanselman&lsquo;s blog posts</a> : <a title="http://www.codinghorror.com/blog/2007/03/the-works-on-my-machine-certification-program.html" href="http://www.codinghorror.com/blog/2007/03/the-works-on-my-machine-certification-program.html" target="_blank">The "Works on My Machine" Certification Program</a>. This blog post fits very well in this program so here <strong>Works on My Machine Seal of Approval</strong>.</p>
<div style="clear: both;"></div>
</blockquote>
<p><strong>Silently Blow up</strong></p>
<p>I have simple ASP.NET MVC 4 internet application which we get out of the box (Things that I will show can be applied with ASP.NET MVC 3 as well). Then I will put a bug inside the index view of my home controller :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@{
    var poo = "bar"
    }

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span>We suggest the following:<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h3</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>As you can see it is not a valid code. So it should fail on both compile time and runtime, right? Let&rsquo;s build the project first :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb.png" width="644" height="166" /></a></p>
<p>Successful! Let&rsquo;s run it :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb_3.png" width="644" height="419" /></a></p>
<p>Boom! Now you have a new, polished <a title="http://weblogs.asp.net/scottgu/" href="http://weblogs.asp.net/scottgu/" target="_blank">ScottGu</a>&rsquo;s YSOD (Yellow Screen of Death). So, the question here is why it wasn&rsquo;t caught by VS.</p>
<p>By default, on ASP.NET MVC projects, your views aren&rsquo;t compiled on build process of your application. VS treats your views just like it treats CSS files. This doesn&rsquo;t mean that VS isn&rsquo;t doing its job, it certainly does. Here is a proof of it :</p>
<p>I went into my <strong>index.cshtml</strong> file and and press <strong>CTRL+W, E</strong> to bring up the Error List window. Here is the result :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb_4.png" width="644" height="138" /></a></p>
<p>Everything is there. But how many of us really check this window on regular basis unless you work on a team project and your team has a habit of putting <strong><a title="http://msdn.microsoft.com/en-us/library/963th5x3(v=vs.71).aspx" href="http://msdn.microsoft.com/en-us/library/963th5x3(v=vs.71).aspx" target="_blank">#warning</a></strong> blocks inside your C# codes. Probably a few of us.</p>
<p>What we can do here is a fairly simple action to get this blow up on compile time.</p>
<p><strong>Blow Up In My Face so I can See you</strong></p>
<p>Right click on your project inside solution explorer and click <strong>Unload Project</strong> section as follows :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_5.png"><img style="background-image: none; margin: 0px 20px 15px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" align="left" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb_5.png" width="427" height="484" /></a></p>
<p>After that your project will be unloaded.</p>
<p>Then, right click it again and it will bring up much shorter. From that list, click on <strong>Edit {yourProjectName}.csproj </strong>section (If your project is a VB project, then .csproj should be .vbproj) as follows :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_6.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb_6.png" width="244" height="135" /></a></p>
<p>As you will see, there is so much going on there. We won&rsquo;t dive into detail there. What we will simply do is to toggle <strong>MvcBuildViews </strong>node from false to true :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb_7.png" width="244" height="116" /></a></p>
<p>Save the file and close it. Then, right click on your project inside solution explorer again. This time click <strong>Reload Project</strong> section :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_8.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb_8.png" width="644" height="242" /></a></p>
<p>Finally, press <strong>CTRL+W, O</strong> to bring up the Output window, build your project and watch it fail :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_9.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/cda9a6d3f884_7739/image_thumb_9.png" width="644" height="166" /></a></p>
<p>Now, when we correct the bug, we will see that we will have a successful build. Note that this won&rsquo;t compile your views into a .dll file, this action will only check them for any compile time errors.</p>
<p>Hope that this helps.</p>