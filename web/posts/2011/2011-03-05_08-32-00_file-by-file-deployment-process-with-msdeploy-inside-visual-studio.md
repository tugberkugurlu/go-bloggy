---
id: dee1f7ff-d170-4129-9078-88ab645a2419
title: File By File Deployment Process with MSDeploy Inside Visual Studio
abstract: Have you ever used 'MSDeploy' inside Visual Studio 2010 and wished a nice
  process bar while publishing a web application? There is even a better way!
created_at: 2011-03-05 08:32:00 +0000 UTC
tags:
- .net
- Deployment
- Tips
- Visual Studio
slugs:
- file-by-file-deployment-process-with-msdeploy-inside-visual-studio
---

<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/File-By-File-Deployment-Process-with-Vis_2306/image.png"><img height="112" width="399" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/File-By-File-Deployment-Process-with-Vis_2306/image_thumb.png" alt="image" title="image" style="background-image: none; margin-top: 0px; margin-right: 15px; margin-bottom: 15px; margin-left: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; float: left; border: 0px initial initial;" /></a></p>
<p>Deployment of your web application is not hard anymore. It is so much easier than before with <strong><em>'MSDeploy'</em></strong>. MSDeploy was introduced to us with Visual Studio 2010. It takes all the application and publish it how we want.</p>
<p>I am not going to explain what <strong>MsDeploy</strong> is and how it works. This is bot the topic here. I assumed that you are reading this post because you have had a least one or two experience with <strong>MSDeploy</strong> inside <strong>Visual Studio 2010</strong>.</p>
<p>When we are publishing an application with one click publish button, you probably noticed a little green world icon with coming and going tiny things on the left bottom side. That think is showing the process of your publishing but not so much informative. I thought that would be cool to view the remaining time of the process but Visual Studio has a lot better feature which I didn&rsquo;t know until last night! Indeed, we have a chance to view the process file by file from output window of visual studio 2010 while publishing your application. That is pretty darn cool. But how we can configure this feature so that we could be able to view that.</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/File-By-File-Deployment-Process-with-Vis_2306/image_3.png"><img height="143" width="244" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/File-By-File-Deployment-Process-with-Vis_2306/image_thumb_3.png" align="left" alt="image" border="0" title="image" style="background-image: none; margin: 0px 15px 15px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" /></a>We need to go over to Tools &gt; Options. Then expend the <em><strong>Projects and Solution</strong></em> section from the list on the window. You will find Build and Run section under that which will give you a window which will look like as it is here on the left. You notice that there is an option called '<strong><em>MsBuild project build output verbosity'</em></strong>. There are five options there and it is set to <strong><em>Minimal</em></strong> by default. I changed it to <strong><em>Normal</em></strong> which is enough for our purpose here.</p>
<p>Save your settings and then expend the output window. You can start your publishing process now and you will see which file is now being deployed by MSDeploy as you can also see on the above screenshot. Yes, awesome&hellip; <img src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/File-By-File-Deployment-Process-with-Vis_2306/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /> Now, think about twice before using an FTP client again <img src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/File-By-File-Deployment-Process-with-Vis_2306/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /></p>