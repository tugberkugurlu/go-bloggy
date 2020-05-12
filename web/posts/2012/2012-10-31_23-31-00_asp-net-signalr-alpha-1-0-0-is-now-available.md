---
id: d7bc98fe-3f0a-4294-a6c9-f2ba33f4ec83
title: ASP.NET SignalR Alpha 1.0.0 is Now Available!
abstract: Couple of hours ago, @DamianEdwards has announced that ASP.NET SignalR Alpha
  1.0.0 release is now publicly available! Even better! SignalR has just shipped with
  ASP.NET Fall 2012 Update!
created_at: 2012-10-31 23:31:00 +0000 UTC
tags:
- ASP.Net
- async
- SignalR
- TPL
slugs:
- asp-net-signalr-alpha-1-0-0-is-now-available
---

<p>Couple of hours ago, <a href="https://twitter.com/DamianEdwards">@DamianEdwards</a> has announced that ASP.NET SignalR Alpha 1.0.0 release is now publicly available.</p>
<blockquote class="twitter-tweet">
<p>ASP(.)NET SignalR 1.0.0 Alpha 1 now live on <a href="https://twitter.com/search/%23nuget">#nuget</a>! <a title="http://nuget.org/packages/microsoft.aspnet.signalr" href="http://t.co/PCqXYMMW">nuget.org/packages/micro&hellip;</a></p>
&mdash; Damian Edwards (@DamianEdwards) <a data-datetime="2012-10-31T19:17:10+00:00" href="https://twitter.com/DamianEdwards/status/263721292845428736">October 31, 2012</a></blockquote>
<script src="//platform.twitter.com/widgets.js"></script>
<p>You can now get the SignalR into your web project through NuGet with the following command:</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package Microsoft.AspNet.SignalR -Pre </code></p>
</div>
<p>Even better! SignalR has just shipped right out of the box with <a href="http://www.asp.net/vnext">ASP.NET Fall 2012 Update</a>! I tried to have a quick view of what has added and changed. In this post, I will share just a few of them.</p>
<p>When you install the package, you will get the most of the usual stuff.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image.png"><img height="244" width="214" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>There is one more thing that I haven&rsquo;t see before (not sure if this has been there with 0.5.3 release). The project is no more registering itself invisibly and RegisterHubs class accomplishes for us.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image_3.png"><img height="333" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>The one other thing that I fell in love with is to be able to return a Task or Task&lt;T&gt; from the hub method! This is a killer feature! Again, I am not sure if this was on 0.5.3 release but I am glad this is now there!</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image_4.png"><img height="157" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image_thumb_4.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>And here it is! We have been all waiting for this one <img src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" /> We now have an AuthorizeAttribute <img src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" /></p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image_5.png"><img height="154" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/image_thumb_5.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>This attribute implements IAuthorizeHubConnection and IAuthorizeHubMethodInvocation interfaces to be recognized as an authorization attribute. So, this means that you can provide your own! If you are familiar with <a href="http://www.asp.net/mvc">ASP.NET MVC</a> or <a href="http://www.asp.net/web-api">ASP.NET Web API</a>, the concept here is the same. <span style="text-decoration: line-through;">However, the interface methods return bool to signal the caller if the call is authorized or not. I would really love to be able to return Task&lt;bool&gt; here or have a similar filter model as ASP.NET Web API.</span>&nbsp;Keep in mind that these are authorization points and they are not meant to be used to perform authantication. SignalR completely leaves the authantication to the underlying hosting layer.</p>
<p>I&rsquo;m sure there are other features but it is 03:24 AM here and my eyes are closing <img src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/ASP.NET-SignalR_2008/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" /> So, it is enough for now. Big thanks and kudos to <a href="https://twitter.com/davidfowl">@davidfowl</a>&nbsp;and <a href="https://twitter.com/DamianEdwards">@DamianEdwards</a> for the Alpha release and for bringing this such a great framework to life.</p>
<p>If I were you, I would go to <a href="https://github.com/SignalR/SignalR">SignalR Github repository</a> and <a href="https://github.com/SignalR/SignalR/tree/master/samples">start exploring the samples</a>. They are awesome and cover the new stuff. Also, <a href="https://twitter.com/DamianEdwards">@DamianEdwards</a>&nbsp;and&nbsp;<a href="https://twitter.com/davidfowl">@davidfowl</a>&nbsp;has a //Build talk tomorrow which will be streamed live: <a href="http://channel9.msdn.com/Events/Build/2012/3-034">http://channel9.msdn.com/Events/Build/2012/3-034</a> Don&rsquo;t miss that one!</p>
<h3>More Information</h3>
<ul>
<li><a href="http://weblogs.asp.net/davidfowler/archive/2012/11/11/microsoft-asp-net-signalr.aspx">Blog: Microsoft ASP.NET SignalR by David Fowler</a></li>
<li><a href="https://github.com/SignalR/SignalR">SignalR on Github</a></li>
<li><a href="http://signalr.net/">SignalR Website</a></li>
</ul>