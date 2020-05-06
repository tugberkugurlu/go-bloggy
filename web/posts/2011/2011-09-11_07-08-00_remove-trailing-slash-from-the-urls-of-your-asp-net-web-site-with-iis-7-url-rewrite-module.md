---
title: Remove Trailing Slash From the URLs of Your ASP.NET Web Site With IIS 7 URL
  Rewrite Module
abstract: One of the aspect of SEO (Search Engine Optimization) is canonicalization.
  In this blog post, we will see how easy to work with IIS Rewrite Module in order
  to remove evil trailing slash from our URLs
created_at: 2011-09-11 07:08:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- Deployment
- IIS
- Tips
slugs:
- remove-trailing-slash-from-the-urls-of-your-asp-net-web-site-with-iis-7-url-rewrite-module
---

<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/redirect-me-baby_3.gif"><img height="240" width="216" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/redirect-me-baby_thumb_3.gif" align="right" alt="redirect-me-baby" border="0" title="redirect-me-baby" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border-width: 0px;" /></a></p>
<p>One of the aspect of <a target="_blank" href="http://en.wikipedia.org/wiki/Search_engine_optimization" title="http://en.wikipedia.org/wiki/Search_engine_optimization">SEO</a> (Search Engine Optimization) is canonicalization. Canonicalization is the process of picking the best URL when there are several choices according to <a target="_blank" href="http://www.mattcutts.com" title="http://www.mattcutts.com">Matt Cutts</a>, the head of Google&rsquo;s Webspam team.</p>
<p>Here is how Matt Cutts explains what canonical URL is :</p>
<p><em>&ldquo;Sorry that it&rsquo;s a strange word; that&rsquo;s what we call it around Google. Canonicalization is the process of picking the best URL when there are several choices, and it usually refers to home pages. For example, most people would consider these the same URLs:</em></p>
<ul>
<li><em>www.example.com</em> </li>
<li><em>example.com/</em> </li>
<li><em>www.example.com/index.html</em> </li>
<li><em>example.com/home.asp</em> </li>
</ul>
<p><em>But technically all of these URLs are different. A web server could return completely different content for all the URLs above.&rdquo;</em></p>
<p>If you have multiple ways of reaching your web page (as above), then you need to sit down because it is time to make some decisions my friends.</p>
<p><strong>Trailing Slash is Evil</strong></p>
<p>Let&rsquo;s assume that we have created a web application, an ASP.NET MVC app because we are so cool. We have our pretty URLs as well.</p>
<p>Let&rsquo;s go to a page on our web site :</p>
<p>http://localhost:55050/Home/About</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_8.png"><img height="323" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_thumb_8.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>And another page :</p>
<p>http://localhost:55050/Home/About/</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_9.png"><img height="323" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_thumb_9.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>We have got the same page content. As we have mentioned before, these two will be treated as two different web page and it will confuse the search engine a bit (even if they are so smart today).</p>
<p>The solution is pretty simple : when a page is requested with trailing slash, then make a 301 (permanent) redirect to the non-trailing-slash version.</p>
<p><strong>IIS URL Rewrite Module</strong></p>
<p>There are several ways of doing that with ASP.NET architecture :</p>
<ul>
<li>You could write your own HttpModule to handle this. </li>
<li>You could do a poor man&rsquo;s redirection on your controller (on your page load if the application is a web forms application). </li>
<li>You could use IIS URL Rewrite Module to easily handle this. </li>
<li>And so on&hellip; </li>
</ul>
<p>In this quick blog post, I will show how we can implement this feature for our whole web site with IIS Rewrite Module.</p>
<p><a target="_blank" href="http://www.iis.net/download/urlrewrite" title="http://www.iis.net/download/urlrewrite">URL Rewrite Module</a> is an awesome extension to IIS. Installing it to your web server is also pretty easy if you haven&rsquo;t got it yet. Just run the <a target="_blank" href="http://www.microsoft.com/web/downloads/platform.aspx" title="http://www.microsoft.com/web/downloads/platform.aspx">Web Platform Installer</a> on your server, and make a search for &ldquo;url rewrite&rdquo;. Then the filtered result will appear and you will see if it is installed or not :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_10.png"><img height="449" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_thumb_10.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>After you have it, you will see the management section inside your IIS Manager under IIS section :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_11.png"><img height="371" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_thumb_11.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p><strong>Cut the crap and show me the code</strong></p>
<p>Now, we are all set up and ready to implement this feature. As it is usual nearly for all Microsoft products, there are thousands (ok, not thousand but still) of way to approach this feature but the easiest way of implementing it is to write the logic inside your web.config file.</p>
<p>As you already know, there is a node called <a target="_blank" href="http://msdn.microsoft.com/en-us/library/bb763179.aspx" title="http://msdn.microsoft.com/en-us/library/bb763179.aspx">system.webServer</a> under the root configuration node. IIS Rewrite Module reserves a node under system.webServer section and allow us to configure the settings there pretty easily. What we will do is to only write the following code under system.webServer node :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">rewrite</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;</span><span style="color: #a31515;">rules</span><span style="color: blue;">&gt;</span>
  
    <span style="color: green;">&lt;!--To always remove trailing slash from the URL--&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">rule</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Remove trailing slash</span><span style="color: black;">"</span> <span style="color: red;">stopProcessing</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">true</span><span style="color: black;">"</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">match</span> <span style="color: red;">url</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">(.*)/$</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">conditions</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">add</span> <span style="color: red;">input</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">{REQUEST_FILENAME}</span><span style="color: black;">"</span> <span style="color: red;">matchType</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">IsFile</span><span style="color: black;">"</span> <span style="color: red;">negate</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">true</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">add</span> <span style="color: red;">input</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">{REQUEST_FILENAME}</span><span style="color: black;">"</span> <span style="color: red;">matchType</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">IsDirectory</span><span style="color: black;">"</span> <span style="color: red;">negate</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">true</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
      <span style="color: blue;">&lt;/</span><span style="color: #a31515;">conditions</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">action</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Redirect</span><span style="color: black;">"</span> <span style="color: red;">redirectType</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Permanent</span><span style="color: black;">"</span> <span style="color: red;">url</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">{R:1}</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">rule</span><span style="color: blue;">&gt;</span>
    
  <span style="color: blue;">&lt;/</span><span style="color: #a31515;">rules</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">rewrite</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>What this code does is to tell the module to remove the trailing slash from the url if there is one and make 301 permanent redirect to the new URL.</p>
<blockquote>
<p>Don&rsquo;t allow the code to freak you out. It might look complicated but there are good recourses out there to make you feel better. Here is one of them :</p>
<p><a target="_blank" href="http://learn.iis.net/page.aspx/460/using-the-url-rewrite-module/" title="http://learn.iis.net/page.aspx/460/using-the-url-rewrite-module/">Using the URL Rewrite Module by Ruslan Yakushev</a></p>
</blockquote>
<p>When you run your site after this implementation and navigate to <strong>/Home/About/</strong>, watch what is going to happen :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_12.png"><img height="386" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_thumb_12.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_13.png"><img height="386" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Remove-Traling_A9C3/image_thumb_13.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Isn&rsquo;t that awesome? A little effort and perfectly clean way of implementing the 1 of a thousand parts of canonicalization.</p>
<p><strong>Some Gotchas</strong></p>
<ul>
<li>In your development environment, if you run your web site under Visual Studio Development Sever, you won&rsquo;t be able to see this feature working. You need to configure your application to run under at least <a target="_blank" href="http://weblogs.asp.net/scottgu/archive/2010/06/28/introducing-iis-express.aspx" title="http://weblogs.asp.net/scottgu/archive/2010/06/28/introducing-iis-express.aspx">IIS Express</a> to see this feature working. </li>
<li>When you deploy your web site and see this feature not working on your server, it is highly possible that you misconfigured something on your server. One of the misconfiguration you might have done could be setting the <strong>overrideModeDefault</strong> attribute to Deny for rules under <strong>&lt;sectionGroup name="rewrite"&gt;</strong> inside your <a target="_blank" href="http://learn.iis.net/page.aspx/124/introduction-to-applicationhostconfig/" title="http://learn.iis.net/page.aspx/124/introduction-to-applicationhostconfig/">applicationHost.config</a> file. </li>
<li>If you are on a shared hosting environment and you see this feature not working, then ask your provider if they have given you the permission of configuring this part. </li>
</ul>