---
id: de260b3d-384c-4c7e-a798-724aa5708597
title: Pro ASP.NET Web API's Sample Project (PingYourPackage) Source Code is Up on
  GitHub
abstract: We wanted to give you an early glimpse on the Pro ASP.NET Web API's sample
  application(PingYourPackage) and its source code is now up on GitHub.
created_at: 2012-12-20 21:32:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET Web API
- GitHub
slugs:
- pro-asp-net-web-api-s-sample-project-pingyourpackage-source-code-is-up-on-github
---

<p>You may know that I have been co-authoring a <a href="http://www.tugberkugurlu.com/archive/pro-asp-net-web-api-book-is-available-on-amazon-for-pre-order">book on ASP.NET Web API</a> for a while now and <a href="http://www.tugberkugurlu.com/archive/pro-asp-net-web-api-book-is-available-through-apress-alpha-program">it is available as an alpha book</a> if you would like to get the early bits and pieces. One of the aim&rsquo;s of this book is to provide a well-structured resource for <a href="http://www.asp.net/web-api">ASP.NET Web API</a> and also to give you a hint on how you would go and build a real world application with this super-awesome framework. In order to make our second goal easy, we have dedicated 3 chapters to build a real world HTTP API from scratch for a small city delivery company and we called it PingYourPackage. These three chapters will be covering the application structure, building the data layer, building and testing the core API layer, creating the .NET wrapper around this API and consuming this through an <a href="http://www.asp.net/mvc">ASP.NET MVC</a> 4 web application.</p>
<p>Although the writing process for those chapters are still not completed, we wanted to give you an early glimpse on the sample application and <a href="https://github.com/tugberkugurlu/PingYourPackage">the source code of the project is now up on GitHub</a>.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/57282e5d8a19_DAC/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/57282e5d8a19_DAC/image_thumb.png" width="644" height="408" /></a></p>
<p>If you are not familiar with GitHub and Git, I recommend watching the <a href="http://haacked.com">Phil Haack</a>&rsquo;s talk on <a href="http://vimeo.com/43612883">Git and GitHub for Developers on Windows</a> and its more hardcore version from <a href="https://twitter.com/dahlbyk">Keith Dahlby</a>: <a href="http://vimeo.com/43659036">Git More Done</a>.</p>
<p>After cloning the source code, you can either directly build the project through the Visual Studio 2012 (or Web Developer Express) or you can run the build script. You can open up the entire solution by double clicking the PingYourPackage.sln file and build it from there.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/57282e5d8a19_DAC/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/57282e5d8a19_DAC/image_thumb_3.png" width="353" height="484" /></a></p>
<p>To run the build script, open up a PowerShell console and <a href="http://technet.microsoft.com/en-us/library/ee176961.aspx">make sure your execution policy is not restricted</a>. When you run the .\scripts\build.ps1 script, it will install the missing NuGet packages, build the entire application and runs the unit/integration tests. The output will be put under .\artifacts folder.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/57282e5d8a19_DAC/image_4.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/57282e5d8a19_DAC/image_4.png" width="640" height="469" /></a></p>
<p>As mentioned, it&rsquo;s not completely done yet but the server level code is complete. The .NET wrapper and web client will be there so soon.</p>
<p>Enjoy <img class="wlEmoticon wlEmoticon-winkingsmile" style="border-style: none;" alt="Winking smile" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/57282e5d8a19_DAC/wlEmoticon-winkingsmile.png" /></p>