---
title: ASP.NET MVC Server Side Remote Validation
abstract: In this quick post, I will show you a way of implementing ASP.NET MVC Server
  Side Remote Validation just like ASP.NET MVC Remote Validation
created_at: 2011-11-06 17:10:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- C#
slugs:
- asp-net-mvc-server-side-remote-validation
---

<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/go-to-server-and-come-back.png"><img height="183" width="244" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/go-to-server-and-come-back_thumb.png" align="left" alt="go-to-server-and-come-back" border="0" title="go-to-server-and-come-back" style="background-image: none; margin: 0px 15px 15px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" /></a></p>
<p>For nearly 15 days, I have been poking around on the internet, books, etc. for anything related to <a target="_blank" href="http://asp.net/mvc" title="http://asp.net/mvc">ASP.NET MVC</a> validation. I have to say that validation features in ASP.NET MVC framework is really outstanding.</p>
<p>I have blogged about <a target="_blank" href="http://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation" title="http://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation">Remote Validation</a> twice so far: <a target="_blank" href="http://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation" title="http://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation">Check Instantly If Username Exists - ASP.NET MVC Remote Validation</a> and <a target="_blank" href="http://www.tugberkugurlu.com/archive/asp-net-mvc-remote-validation-for-multiple-fields-with-additionalfields-property" title="http://www.tugberkugurlu.com/archive/asp-net-mvc-remote-validation-for-multiple-fields-with-additionalfields-property">ASP.NET MVC Remote Validation For Multiple Fields With AdditionalFields Property</a>.</p>
<p>Remote validation is one of the useful stuff which is baked into MVC framework. It is easy to implement this feature manually with <a target="_blank" href="http://jquery.com/" title="http://jquery.com/">JQuery</a> but it is kind of nice to have it out of the box. So, thanks a lot ASP.NET MVC team.</p>
<p>One thing which ASP.NET MVC Remove Validation is missing is no support for server side validation of the chosen property. There is probably a good reason not to support it out of the box though.</p>
<p>So, I decided to build one for myself. When I implemented it, I thought it worked fine and I should blog about it. So, I am writing this <strong>short </strong>blog post. But, on the other hand I am also sure that no one should use this on production at least for now<img src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /> I don&rsquo;t know why but it feels like it is not stable. Maybe you can use it if you see a positive comment from an MSFT person on this post.</p>
<p>Here how it works:</p>
<p>First of all, we need to use Nuget here to bring down a very small package called <strong><a target="_blank" href="http://nuget.org/List/Packages/TugberkUg.MVC" title="http://nuget.org/List/Packages/TugberkUg.MVC">TugberkUg.MVC</a></strong> which will have the necessary stuff for server side remote validation to work.</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_4.png"><img original="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_4.png" height="75" width="640" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_4.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>After that, here how we can use it:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>  [Required]
  [Display(Name = <span style="color: #a31515;">"User name"</span>)]
  [Remote(<span style="color: #a31515;">"doesUserNameExist"</span>, <span style="color: #a31515;">"Account"</span>, HttpMethod = <span style="color: #a31515;">"POST"</span>, ErrorMessage = <span style="color: #a31515;">"User name already exists. Please enter a different user name."</span>)]
  [ServerSideRemote(<span style="color: #a31515;">"Account"</span>, <span style="color: #a31515;">"doesUserNameExistGet"</span>)]
  <span style="color: blue;">public</span> <span style="color: blue;">string</span> UserName { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }</pre>
</div>
</div>
<blockquote>
<p>You do not need to use <a target="_blank" href="http://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation" title="http://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation">ASP.NET MVC Remote validation</a> along with ServerSideRemote validation but in a real world scenario, we probably would use both.</p>
</blockquote>
<p>As you can see, there is not much to specify but of course, you get the <span style="text-decoration: line-through;">options</span> properties of <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.componentmodel.dataannotations.validationattribute.aspx" title="http://msdn.microsoft.com/en-us/library/system.componentmodel.dataannotations.validationattribute.aspx">ValidationAttribute</a> class along with it. there is not much to specify there because I have built it for like 30 minutes. So, spare me for that.</p>
<p>When we look at our controller action which holds the validation logic, we will see that it is not much different that Remote validation&rsquo;s:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> JsonResult doesUserNameExistGet(<span style="color: blue;">string</span> term) {

    <span style="color: blue;">var</span> user = Membership.GetUser(term);

    <span style="color: blue;">return</span> Json(user == <span style="color: blue;">null</span>, JsonRequestBehavior.AllowGet);
}</pre>
</div>
</div>
<p>When I fire up the register page, I am getting the following screen :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image.png"><img height="398" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image_thumb.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>I have excluded the JavaScript files for client side validation because we do not want to see client validation working. I already have a user registered whose user name is <strong>User1</strong>. When I completed the fields, I am getting no warning:</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image_3.png"><img height="441" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image_thumb_3.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Let&rsquo;s push the <em>Register</em> button :</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image_4.png"><img height="441" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image_thumb_4.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Bingo! We got it working. Let&rsquo;s try a legitimate one:</p>
<p><a href="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image_5.png"><img height="441" width="644" src="http://tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/de656a5ad88d_11FE4/image_thumb_5.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>That worked as well.</p>
<p>I also put the sample project on <a href="https://github.com/">GitHub</a> so you can get the working sample code if you want :</p>
<p><a target="_blank" href="https://github.com/tugberkugurlu/MvcRemoteValidationSample" title="https://github.com/tugberkugurlu/MvcRemoteValidationSample">https://github.com/tugberkugurlu/MvcRemoteValidationSample</a></p>
<p>Also, you can find the whole <a target="_blank" href="http://nuget.org/List/Packages/TugberkUg.MVC" title="http://nuget.org/List/Packages/TugberkUg.MVC">TugberkUg.MVC</a> package&rsquo;s code on <a target="_blank" href="https://bitbucket.org" title="https://bitbucket.org">BitBucket</a>:</p>
<p><a target="_blank" href="https://bitbucket.org/tugberk/tugberkug.mvc/src" title="https://bitbucket.org/tugberk/tugberkug.mvc/src">https://bitbucket.org/tugberk/tugberkug.mvc/src</a></p>