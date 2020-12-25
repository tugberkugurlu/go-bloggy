---
id: cf5ff6da-9150-4e87-8bd9-bc5961332e53
title: Check Instantly If Username Exists - ASP.NET MVC Remote Validation
abstract: This blog post will walk you through on implementation and usage of ASP.NET
  MVC Remote Validation. As a sample, we will validate the availability of the username
  for membership registration.
created_at: 2011-10-19 04:07:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET MVC
- C#
- JQuery
slugs:
- check-instantly-if-username-exists-asp-net-mvc-remote-validation
- check-instantly-if-username-exists-asp-net-mvc-remote-validatio
---

<p>On some of our web applications, we require users to register our site in order to perform some special actions. This feature kind of a no-brainer thing thanks to built-in ASP.NET Membership API. It works well with either ASP.NET Web Forms or <a target="_blank" href="http://asp.net/mvc" title="http://asp.net/mvc">ASP.NET MVC</a>. In this path, on some cases, we also require each users to have <strong>unique username</strong> or <strong>unique e-mail</strong>. The implementation of this feature also does not require our concern either.</p>
<p>Let&rsquo;s see it in action what I am getting at here.</p>
<p>I created an ASP.NET MVC 3 Web Application project with Visual Studio 2010 (It is also same for Visual Web Developer which is a free version of Visual Studio) and I fired it up by clicking <strong>Ctrl + F5</strong>. When I navigate to <em>/Account/Register</em>, I get the following view :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image.png"><img height="422" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_thumb.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>In order to stick with the topic here, I already registered a user with a username of <em><strong>user1</strong></em>.<em><strong> </strong></em>Look what happens when I try to create another user with a username of <strong><em>user1 </em></strong>:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_3.png"><img height="422" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_thumb_3.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>So, the thing ASP.NET offers us out of the box is clever. But there is a little bit a problem here. Well, not a problem but how can I call it, hmmm&hellip; I can say lack of creativeness. When I push the register button here, it goes all the way to server to register the user. Inside our controller, built-in register action checks if user already exists or not.</p>
<p><em>Wouldn&rsquo;t it be nice if user instantly sees if typed username is valid or not?<strong> </strong></em>It certainly would be!</p>
<p><strong>Introduction to ASP.NET MVC Remote Validation</strong></p>
<p>We know that ASP.NET MVC takes advantage of .NET Framework Data Annotations to validate user&rsquo;s input. By this way, we can check for required fields with <strong>RequiredAttribute</strong>, length of a string with <strong>StringLengthAttribute </strong>and so on. Well, it turns out that there is another useful one called <strong>RemoteAttribute</strong> which we can use, for example, to validate e-mail and username in a registration form and alert the user before the form is posted.</p>
<p>This <strong>Remote Validation</strong> thing is what we are looking for. In order to implement this for our scenario, we need to tweak the <strong>UserName</strong> property of <strong>RegisterModel</strong> class. Here is the whole RegisterModel class after the modification I made :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>    <span style="color: blue;">public</span> <span style="color: blue;">class</span> RegisterModel {

        [Required]
        [Display(Name = <span style="color: #a31515;">"User name"</span>)]
        [Remote(<span style="color: #a31515;">"doesUserNameExist"</span>, <span style="color: #a31515;">"Account"</span>, HttpMethod = <span style="color: #a31515;">"POST"</span>, ErrorMessage = <span style="color: #a31515;">"User name already exists. Please enter a different user name."</span>)]
        <span style="color: blue;">public</span> <span style="color: blue;">string</span> UserName { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

        [Required]
        [DataType(DataType.EmailAddress)]
        [Display(Name = <span style="color: #a31515;">"Email address"</span>)]
        <span style="color: blue;">public</span> <span style="color: blue;">string</span> Email { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

        [Required]
        [StringLength(100, ErrorMessage = <span style="color: #a31515;">"The {0} must be at least {2} characters long."</span>, MinimumLength = 6)]
        [DataType(DataType.Password)]
        [Display(Name = <span style="color: #a31515;">"Password"</span>)]
        <span style="color: blue;">public</span> <span style="color: blue;">string</span> Password { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

        [DataType(DataType.Password)]
        [Display(Name = <span style="color: #a31515;">"Confirm password"</span>)]
        [Compare(<span style="color: #a31515;">"Password"</span>, ErrorMessage = <span style="color: #a31515;">"The password and confirmation password do not match."</span>)]
        <span style="color: blue;">public</span> <span style="color: blue;">string</span> ConfirmPassword { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    }</pre>
</div>
</div>
<p>As you can see, we added <strong>Remote </strong>attribute for UserName property inside RegisterModel class. RemoteAttribute has several overloads but what we entered here is as follows :</p>
<ul>
<li><strong><em>doesUserNameExist</em></strong> : It is the name of the Action which will validates the user input by returning true or false as JSON format. We will see how we do that in a minute. </li>
<li><strong><em>Account</em></strong> : It is the name of the controller which <strong><em>doesUserNameExist</em></strong> lives in. </li>
<li><strong><em>HttpMethod </em></strong>: Simply we are telling the remote validation to communicate our mini service over Http POST method (This is not required but I&rsquo;d like to do it this way). </li>
</ul>
<p>How our controller action <strong><em>doesUserNameExist </em></strong>looks like is as follows :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[HttpPost]
<span style="color: blue;">public</span> JsonResult doesUserNameExist(<span style="color: blue;">string</span> UserName) {

    <span style="color: blue;">var</span> user = Membership.GetUser(UserName);

    <span style="color: blue;">return</span> Json(user == <span style="color: blue;">null</span>);
}</pre>
</div>
</div>
<p>As you see, it is a simple controller action which returns <strong>JsonResult</strong>. What we need to notice here is <strong>UserName</strong> paramater of this action. We put that parameter there to pick up the user input which will be sent through <strong>Ajax request</strong> to our action. Remote Validation will send the user input as Form Data which will have the same name as property. So, in this case it is <strong>UserName</strong>.</p>
<p>Before we give it a try, make sure that you have the following libraries referenced on your registration page :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">script</span> <span style="color: red;">src</span><span style="color: blue;">=</span><span style="color: blue;">"@Url.Content("~/Scripts/jquery.validate.min.js")"</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: blue;">"text/javascript"</span><span style="color: blue;">&gt;</span><span style="color: blue;">&lt;/</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;</span><span style="color: #a31515;">script</span> <span style="color: red;">src</span><span style="color: blue;">=</span><span style="color: blue;">"@Url.Content("~/Scripts/jquery.validate.unobtrusive.min.js")"</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: blue;">"text/javascript"</span><span style="color: blue;">&gt;</span><span style="color: blue;">&lt;/</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>Let&rsquo;s build our project and fire it up to see it in action.</p>
<p><strong>1.</strong> I loaded the page :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_4.png"><img height="419" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_thumb_4.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p><strong>2.</strong> I enter <strong><em>user1 </em></strong>for UserName field (user1 is the username which is already inside our database) :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_5.png"><img height="419" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_thumb_5.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p><strong>3.</strong> Then I just put the mouse cursor inside another text input field and see what happened :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_6.png"><img height="419" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/e349262d0d85_1411C/image_thumb_6.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Just what we needed. I thought showing this in action like that is a lame approach to do that so I recorded a short video in order to demonstrate this. In this video I also tried to show what is happening behind the scenes. See it below :</p>
<p><iframe width="425" height="349" src="http://www.youtube.com/embed/wIZRQBAdkfM?hl=en&amp;fs=1" frameborder="0"></iframe></p>
<p>I also put the sample project on <a target="_blank" href="https://github.com" title="https://github.com">GitHub</a> so you can get the working sample code if you want :</p>
<p><a href="https://github.com/tugberkugurlu/MvcRemoteValidationSample">https://github.com/tugberkugurlu/MvcRemoteValidationSample</a></p>
<p><strong>Note :&nbsp;</strong></p>
<p>But there is a donwside for this validation attribute. The remote validation does not validate the user input on server side. So, you need to make sure on the server side that it is valid as it is being done on our sample. We also make sure on the server side that username which is being passed to our controller is unique.</p>
<p>Hope this will help you even a little bit.</p>