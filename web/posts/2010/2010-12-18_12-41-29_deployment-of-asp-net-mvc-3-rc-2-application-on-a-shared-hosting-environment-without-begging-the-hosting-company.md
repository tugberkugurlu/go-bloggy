---
id: e722e728-ef06-4408-b032-4000a94a59b0
title: Deployment of ASP.Net MVC 3 RC 2 Application on a Shared Hosting Environment
  Without Begging The Hosting Company
abstract: After the release of ASP.Net MVC RC 2, we are now waiting for the RTM release
  but some of us wanna use RC 2 already... But how to deploy it on a shared hosting
  acount is the mind-exploding problem...
created_at: 2010-12-18 12:41:29 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET MVC
- Deployment
- Hosting
slugs:
- deployment-of-asp-net-mvc-3-rc-2-application-on-a-shared-hosting-environment-without-begging-the-hosting-company
---

<blockquote>
<p><strong>UPDATE on 2011, 02.26</strong></p>
<p>I have wrote another blog post on ASP.NET MVC Deployment problems you might have related to your server. If you are still having problems (especially, if you are getting 404 exceptions for extensionless URLs), you might want to have a look at on "<em><a target="_blank" title="http://tugberkugurlu.com/archive/running-asp-net-mvc-under-iis-6-0-and-iis-7-0-classic-mode---solution-to-routing-problem" href="http://tugberkugurlu.com/47">Running ASP.NET MVC Under IIS 6.0 and IIS 7.0 Classic Mode : Solution to Routing Problem</a></em>"</p>
</blockquote>
<p><img src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/asp-net-mvc-3.gif" width="270" height="210" alt="asp-net-mvc-3.gif" title="asp-net-mvc-3.gif" style="float: left; margin: 0 10px 10px 0;" />On the 10<sup>th</sup> of December in 2010, Microsoft ASP.Net MVC team has released the <a target="_blank" title="http://www.asp.net/mvc/mvc3" href="http://www.asp.net/mvc/mvc3">MVC 3 RC 2</a>&nbsp;and it has some good stuff inside which RC 1 didn't have. It is not so much different but there are some breaking changes, especially <em>ViewBag</em>&nbsp;thing.</p>
<p>You could find more information about ASP.Net MVC RC 3 goodies on <a href="http://weblogs.asp.net/scottgu/archive/2010/12/10/announcing-asp-net-mvc-3-release-candidate-2.aspx" title="http://weblogs.asp.net/scottgu/archive/2010/12/10/announcing-asp-net-mvc-3-release-candidate-2.aspx" target="_blank">ScootGu's blog post</a> or <a target="_blank" title="http://haacked.com/archive/2010/12/10/asp-net-mvc-3-release-candidate-2.aspx" href="http://haacked.com/archive/2010/12/10/asp-net-mvc-3-release-candidate-2.aspx">Phil Haacked blog post</a>.</p>
<p>The MVC 3 RC 2 has come from go-live-license&nbsp;so you can use this in production if you wish. But the main problem rises here if you are in a shared hosting environment. As you know, shared hosting providers are not willing to install the new releases unless it is for sure that there is no detected bug in the package.</p>
<p>I thought that would be a big problem for me [because I am still in shared hosting environment :)] but deploying the necessary assemblies as manually is the solution. Scott Hanselman has a great post on <a href="http://www.hanselman.com/blog/BINDeployingASPNETMVC3WithRazorToAWindowsServerWithoutMVCInstalled.aspx" title="http://www.hanselman.com/blog/BINDeployingASPNETMVC3WithRazorToAWindowsServerWithoutMVCInstalled.aspx" target="_blank">how to deploy an MVC 3 application into a shared hosting environment</a>. The article covers all the necessary steps. One problem is that if you do exactly as it is there, you will sure have a problem if you are deploying the RC 2 of the MVC 3.</p>
<p>The yellow screen of death will give you the following error;</p>
<p><span style="color: #ff0000;"><strong>Could not load file or assembly 'System.Web.WebPages.Deployment, Version=1.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35' or one of its dependencies. The system cannot find the file specified.</strong></span></p>
<p>The error will be thrown by the system because it needs System.Web.WebPages.Deployment.dll as well. The solution of this little problem is simple;</p>
<p>Navigate to&nbsp;<strong><em>C:\Program Files\Microsoft ASP.NET\ASP.NET Web Pages\v1.0\Assemblies</em></strong> (this could be <em>C:\Program Files (x86)\Microsoft ASP.NET\ASP.NET Web Pages\v1.0\Assemblies </em>on Windows 7) inside the windows explorer and you will see some files inside the folder;</p>
<p><img title="asp.net-mvc-3-rc-2-bin-deployment-shared-hosting-environment-full.PNG" alt="asp.net-mvc-3-rc-2-bin-deployment-shared-hosting-environment-full.PNG" height="601" width="803" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/asp.net-mvc-3-rc-2-bin-deployment-shared-hosting-environment-full.PNG" /></p>
<p>You need those 6 dll files. In addition to that you will also need&nbsp;System.Web.Mvc.dll (version 3.0.11209.0). You will be able to find that dll file by navigating to&nbsp;<strong><em>C:\Program Files\Microsoft ASP.NET\ASP.NET MVC 3\Assemblies</em></strong> (that should be <em>C:\Program Files (86x)\Microsoft ASP.NET\ASP.NET MVC 3\Assemblies</em>&nbsp;in Windows 7 [I'm not sure though])</p>
<p><img src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/asp.net-mvc-3-rc-2-bin-deployment-shared-hosting-environment-mvc-folder.PNG" width="805" height="379" alt="asp.net-mvc-3-rc-2-bin-deployment-shared-hosting-environment-mvc-folder.PNG" title="asp.net-mvc-3-rc-2-bin-deployment-shared-hosting-environment-mvc-folder.PNG" /></p>
<p>Finally you should all have the following dlls in hand;</p>
<ul>
<li>Microsoft.Web.Infrastructure</li>
<li>System.Web.Razor</li>
<li>System.Web.WebPages</li>
<li>System.Web.WebPages.Razor</li>
<li>System.Web.Helpers</li>
<li>System.Web.WebPages.Deployment (If you are deploying MVC RC 2, this assembly is necessary to deploy)</li>
<li>System.Web.Mvc</li>
</ul>
<p>I made a copy of those dlls and put them together inside a folder so that I could reach them easily whenever I need them;</p>
<p><img src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/essential-MVC-3-0-RC-2-with-razor-dlls.PNG" width="749" height="439" alt="essential-MVC-3-0-RC-2-with-razor-dlls.PNG" title="essential-MVC-3-0-RC-2-with-razor-dlls.PNG" /></p>
<p>We have all the necessary files in our hands and so what now !</p>
<p>There are some conventions here that you could choose. The way I follow is that;</p>
<ul>
<li>I used built in Visual Studio Publish tool to publish my application to the production side.</li>
<li>After publishing was complated, I simply copied those 7 dll files into the bin folder inside the root directory of my application.</li>
</ul>
<p><img src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/visual-studio-2010-publish-tool-goodies.png" width="314" height="232" alt="visual-studio-2010-publish-tool-goodies.png" title="visual-studio-2010-publish-tool-goodies.png" /><br /> <span style="font-size: xx-small;"><em>Visual Studio 2010 Publish Tool</em></span></p>
<p>That was it ! [Of cource I didn't put the System.Web.WebPages.Dployment.dll into the production side and I got the error firstly :)]</p>
<p>Finally, my bin folder has the necessary assembly files to run the application;</p>
<p><img src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/magic-bin-folder-of-asp-net-mvc-3-rc-2-application.png" width="258" height="443" alt="magic-bin-folder-of-asp-net-mvc-3-rc-2-application.png" title="magic-bin-folder-of-asp-net-mvc-3-rc-2-application.png" /></p>