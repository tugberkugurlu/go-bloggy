---
id: 2ad45143-ba5f-4d67-8d2d-7ab203c780fc
title: 'Running ASP.NET MVC Under IIS 6.0 and IIS 7.0 Classic Mode : Solution to Routing
  Problem'
abstract: In this blog post, we will see how to run ASP.NET MVC application under
  IIS 6.0 and IIS 7.0 classic mode with some configurations on IIS and Global.asax
  file...
created_at: 2011-02-26 07:31:00 +0000 UTC
tags:
- .NET
- ASP.Net
- ASP.NET MVC
- Deployment
- Hosting
- IIS
slugs:
- running-asp-net-mvc-under-iis-6-0-and-iis-7-0-classic-mode---solution-to-routing-problem
---

<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/j0442463.jpg"><img height="173" width="244" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/j0442463_thumb.jpg" align="left" alt="Highway sign" border="0" title="Highway sign" style="background-image: none; margin: 0px 15px 15px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" /></a></p>
<p>I wrote a blog post on <em>&ldquo;<a target="_blank" href="https://www.tugberkugurlu.com/39" title="Deployment of ASP.Net MVC 3 RC 2 Application on a Shared Hosting Environment Without Begging The Hosting Company">Deployment of ASP.Net MVC 3 RC 2 Application on a Shared Hosting Environment Without Begging The Hosting Company</a>&rdquo;</em> couple of months ago. The solution was working for most case scenarios if the server is configured properly for <strong>ASP.NET Routing</strong>. Other working case I have seen was the applications which are running under <strong>IIS 7.0 integrated mode</strong>. Under IIS 7.0 integrated mode, no special configuration necessary to use ASP.NET Routing.</p>
<p>As we know, one of the most beautiful parts of ASP.NET MVC framework is <strong><em>Routing</em></strong>. We have nice, clean, <em><strong>extensionless</strong></em> URLs thanks to routing and this is becoming an issue under IIS 6.0 and IIS 7.0 classic mode.</p>
<p>When we typing the path of a web site page inside the address bar of our web browser, we are making a request against server. If our web application is running under any version of IIS, the request hits ASP.NET framework on certain conditions. Especially, Older versions of IIS only map certain requests to the ASP.NET framework. If the extension of the web request is aspx, ashx, axd or any other extensions which is specific for ASP.NET framework are being mapped to ASP.NET framework. <em><strong>So, in a MVC application the requests are not being mapped to ASP.NET framework.</strong></em>&nbsp;</p>
<p>&nbsp;<a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/image.png"><img height="173" width="644" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/image_thumb.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a> <br /><em><span style="font-size: xx-small;" size="1">(You will be getting this 404 exception when you hit the extensionless URL of your application)</span></em></p>
<p>No, do not throw away your precious, new born ASP.NET MVC application which you have created working along with Nuget PMC <em>(which is perfect), </em>EFCodeFirst and any other cool newbie stuff.&nbsp; There are optional solutions for this problem and they are not like hard things to implement.</p>
<blockquote>
<p><em>Although, if you have a Windows Server 2003 and thinking about going up to IIS 7.0 integrated or classic mode, stop right there my friend ! Because, <strong>IIS 7.0 is not compatible with Windows Server 2003</strong>. So you are stuck with IIS 6.0 for now and keep reading for the solution <img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /></em></p>
</blockquote>
<p>The solutions are optional as indicated. They depend on what kind of powers you have over your server. Here is the list of solution you might be interested;</p>
<p><span style="color: #008040;" color="#008040"><strong>Option 1 : I am the guy with the full control power over my server and I want to keep extensionless URLs</strong></span></p>
<p>If you have full access over your server, you could create so called <strong><em>Wildcard Script Map </em></strong>so that you can use the default ASP.NET MVC route table with IIS 7.0 (in classic mode) or IIS 6.0. This Wildcard Script Map will map all requests to the web server to the ASP.NET framework.</p>
<blockquote>
<p>I have no experience with this option, though. I had a problem like this within this week, I have solved it with the following option and didn&rsquo;t want to use this one even if I have full control over my server. I am not a server pro, so I won&rsquo;t be making any comments on how this will effect the requests flow made against your server. I just thought that <em>&ldquo;man, every single request which made against my server will be mapped to ASP.NET framework. This could effect the speed of the delivery process.&rdquo; </em>and that how I skipped this option.</p>
<p>And, bad news guys <img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/wlEmoticon-sarcasticsmile.png" alt="Sarcastic smile" class="wlEmoticon wlEmoticon-sarcasticsmile" style="border-style: none;" /> Microsoft, also, indicated the following line of sentence one of their web site;</p>
<p><em>&ldquo;Be aware that this option causes IIS to intercept every request made against the web server. This includes requests for images, classic ASP pages, and HTML pages. Therefore, enabling a wildcard script map to ASP.NET does have performance implications.&rdquo;</em></p>
</blockquote>
<p>To implement this feature, you need to follow some steps and here is the text I grabbed from <a href="http://asp.net/">http://asp.net/</a> ;</p>
<p>&nbsp;</p>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/image_thumb5.png"><img height="188" width="248" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/image_thumb5_thumb.png" align="right" alt="image_thumb[5]" border="0" title="image_thumb[5]" style="background-image: none; margin: 0px 0px 0px 10px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Here's how you enable a wildcard script map for IIS 7.0 (classic mode):</p>
<ol>
<li>Select your application in the Connections window </li>
<li><l>Make sure that the <strong>Features</strong> view is selected </l></li>
<li><l>Double-click the <strong>Handler Mappings</strong> button </l></li>
<li><l>Click the <strong>Add Wildcard Script Map</strong> link </l></li>
<li>Enter the path to the aspnet_isapi.dll file (You can copy this path from the PageHandlerFactory script map) </li>
<li><l>Enter the name MVC <l>Click the <strong>OK</strong> button </l></l></li>
</ol>
<p>&nbsp;</p>
<p>&nbsp;</p>
<p>&nbsp;</p>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/image_thumb4.png"><img height="179" width="248" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/image_thumb4_thumb.png" align="right" alt="image_thumb[4]" border="0" title="image_thumb[4]" style="background-image: none; margin: 0px 0px 0px 10px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Follow these steps to create a wildcard script map with IIS 6.0:</p>
<ol>
<li>Right-click a website and select Properties </li>
<li><l>Select the <strong>Home Directory</strong> tab </l></li>
<li><l>Click the <strong>Configuration</strong> button </l></li>
<li><l>Select the <strong>Mappings</strong> tab </l></li>
<li><l>Click the <strong>Insert</strong> button </l></li>
<li><l>Paste the path to the aspnet_isapi.dll into the Executable field (you can copy this path from the script map for .aspx files) </l></li>
<li><l>Uncheck the checkbox labeled <strong>Verify that file exists</strong> </l></li>
<li><l>Click the <strong>OK</strong> button </l></li>
</ol>
<p>&nbsp;</p>
<p><span style="color: #008040;" color="#008040"><strong>Option 2 : I am the guy who has the full control power over my server, cares about performance and not care about URLs</strong></span></p>
<p>This option is the one of the other option I do not like very much but maybe you will <img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Running-ASP.NET-MVC-Under.0-Classic-Mode_14547/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /> So here is the deal;</p>
<p>We will simply add Extensions to the Route Table so that older versions of IIS can pass requests to the ASP.NET framework. This option requires changes inside <span style="font-weight: bold;"><em>Global.asax </em></span>file of you application and some addition work on IIS for modifying the Default route so that it includes a file extension that is mapped to the ASP.NET framework.</p>
<p>Let&rsquo;s see the RegisterRoutes method inside your <strong><em>Global.asax </em></strong>file;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [7]">        public static void RegisterRoutes(RouteCollection routes)
        {
            routes.IgnoreRoute("{resource}.axd/{*pathInfo}");

            routes.MapRoute(
                "Default",                                              // Route name
                "{controller}/{action}/{id}",                           // URL with parameters
                new { controller = "Home", action = "Index", id = "" }  // Parameter defaults
            );

        }</pre>
<p>As you see here on line 7, we are aiming to get URLs without extension. Let&rsquo;s see what it needs to be look like after we change it in order to implement this option;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [7]">        public static void RegisterRoutes(RouteCollection routes)
        {
            routes.IgnoreRoute("{resource}.axd/{*pathInfo}");

            routes.MapRoute(
                "Default",
                "{controller}.mvc/{action}/{id}",
                new { action = "Index", id = "" }
              );

            routes.MapRoute(
              "Root",
              "",
              new { controller = "Home", action = "Index", id = "" }
            );


        }</pre>
<p>Here, we are assigning .mvc extension for every URL for <em>controller</em> name.</p>
<blockquote>
<p>I want to warn you about your links inside your views. If you created them by hard coding, now you are officially screwed my friend. Because, you either have to change all of them by hand or need to implement option 1 in order to keep them unbreakable. If you used <em>ActionLink</em> or <em>RouteLink</em> kind of way, then your are good to go. Changes will be handled by MVC Framework for you.</p>
</blockquote>
<p>Therefore, to get ASP.NET Routing to work, we must modify the Default route so that it includes a file extension that is mapped to the ASP.NET framework.</p>
<p>This is done using a script named <code>registermvc.wsf</code>. It was included with the ASP.NET MVC 1 release in <code>C:\Program Files\Microsoft ASP.NET\ASP.NET MVC 1.0\Scripts</code>, but as of ASP.NET 2 this script has been moved to the ASP.NET Futures, available at <a href="http://aspnet.codeplex.com/releases/view/39978">http://aspnet.codeplex.com/releases/view/39978</a>.</p>
<p>Executing this script registers a new .mvc extension with IIS. After you register the .mvc extension, you can modify your routes in the Global.asax file so that the routes use the .mvc extension.</p>
<p>After this implementation, your URLs will look like this;</p>
<blockquote>
<p><em>/Home.mvc/Index/</em></p>
<p><em>/Product.mvc/Details/3</em></p>
<p><em>/Product.mvc/</em></p>
</blockquote>
<p><span style="color: #008040;" color="#008040"><strong>Option 3 : I am the guy who has no control power over my server, not care about URLs (If you care, it does not matter. You have no choice)</strong></span></p>
<p>This option is the easiest way of make your application up and running within minutes depending on your application structure. Only you need to do here is; making some changes inside your <strong><em>Global.asax</em></strong> file, recompiling your application and publishing it into your server. That&rsquo;s all. Let&rsquo;s see how our new <em>RegisterRoutes</em> method inside <strong><em>Global.asax</em></strong> file needs to be look like;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [7]">        public static void RegisterRoutes(RouteCollection routes)
        {
            routes.IgnoreRoute("{resource}.axd/{*pathInfo}");

            routes.MapRoute(
                "Default",
                "{controller}.ashx/{action}/{id}",
                new { action = "Index", id = "" }
              );

            routes.MapRoute(
              "Root",
              "",
              new { controller = "Home", action = "Index", id = "" }
            );


        }</pre>
<p>Look at line 7. Isn&rsquo;t it familiar? It is one of the ASP.NET framework extensions. You could add <em>.aspx</em> or whatever you want from ASP.NET framework extensions. It is already registered into the Default route so <strong>the requests will be mapped to the ASP.NET framework.</strong></p>
<p>I hope that this is the solution you are looking for your problem.</p>