---
id: 1cf88d24-ed88-48a6-b616-26014de6bd02
title: jQuery.simpleJSlide-1.1 / Simple jQuery Slide Show Plugin With Unobtrusive
  JavaScript Support
abstract: An awesome and super simple jQuery slide show plugin. This plugin also takes
  advantage of unobtrusive JavaScript method as well.
created_at: 2011-05-24 12:44:00 +0000 UTC
tags:
- JQuery
- NuGet
slugs:
- jquery-simplejslide-1-1-simple-jquery-slide-show-plugin-with-unobtrusive-javascript-support
---

<p>Well, the time actually came for me for both my first <a title="http://nuget.org" href="http://nuget.org" target="_blank">Nuget</a> package and <a title="http://jquery.com/" href="http://jquery.com/" target="_blank">jQuery</a> plugin : <a title="http://nuget.org/List/Packages/jquery.simpleJSlide" href="http://nuget.org/List/Packages/jquery.simpleJSlide" target="_blank">jQuery.simpleJSlide-1.1</a></p>
<p>In this blog post I am not going to go through how I created the Nuget package. Instead, I am going to explain how this small and simple slide show plugin works.</p>
<blockquote>
<p>If you are interested in Nuget, there are many great learning stuff on <a title="http://docs.nuget.org/" href="http://docs.nuget.org/" target="_blank">Nuget Docs</a> and you can also watch some basic and extreme videos on Nuget;</p>
<p>&nbsp;<a href="http://docs.nuget.org/docs/start-here/videos">http://docs.nuget.org/docs/start-here/videos</a></p>
</blockquote>
<p>jQuery.simpleJSlide is an awesome and super simple jQuery slide show plugin. This plugin also takes advantage of unobtrusive JavaScript method as well. As I explained above, you can reach this plugin easily via Nuget;</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb.png" width="640" height="72" /></a></p>
<p>Also, there is a sample package which has a basic sample example inside it. If you would like to install this package and also would like to have an example, I recommend you to install <a title="http://nuget.org/List/Packages/jquery.simpleJSlide.Sample" href="http://nuget.org/List/Packages/jquery.simpleJSlide.Sample" target="_blank">jQuery.simpleJSlide.Sample</a>;</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image4_thumb.png" width="640" height="75" /></a></p>
<p><em>jQuery.simpleJSlide.Sample</em> package has a dependency on <em>jQuery.simpleJSlide</em> package and <em>jQuery.simpleJSlide</em> has a dependency on <em>jQuery</em> and <a title="http://nuget.org/List/Packages/jquery.ui.unobtrusively" href="http://nuget.org/List/Packages/jquery.ui.unobtrusively" target="_blank"><em>jQuery.ui.unobtrusively</em></a> package. jQuery UI Unobtrusively package is a simple package which contains only one .js file that enables you to wire-up jQuery UI unobtrusively. <em>This code iterates the page elements, finds elements with the custom attribute <strong>data-ui-fn</strong> and then apply the desired JavaScript / jQuery attribute.</em> (<a title="http://www.msjoe.com/2011/05/unobtrusive-javascript-in-your-asp-net-pages/" href="http://www.msjoe.com/2011/05/unobtrusive-javascript-in-your-asp-net-pages/" target="_blank">Quoted from MSJoe</a>) For Example;</p>
<pre class="brush: xhtml; toolbar: false">&lt;asp:TextBox runat="server" ID="startDate" data-ui-fn="datepicker" /&gt;</pre>
<p>When you reference the proper JavaScript files (<em>jQeury</em>, <em>jQuery UI</em> and <em>jQuery.UI.Unobtrusively</em>) to a page, you are enabling <a title="http://jqueryui.com/demos/datepicker/" href="http://jqueryui.com/demos/datepicker/" target="_blank">jQeury UI Datepicker</a> by writing above code. Yes, you are not writing one line of JavaScript code here. Very simple and clean. But what if we want to override or set some options. We&rsquo;ll get to that latter but I am not going to dive into this package more in this blog post though.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb_3.png" width="640" height="74" /></a></p>
<p>I pushed this package into Nuget live feed but this code actually written by <a title="http://damianedwards.wordpress.com" href="http://damianedwards.wordpress.com" target="_blank">Damian Edwards</a> who is a web guy at Microsoft. After I pushed the package to Nuget live feed, I told him about the package and he is now the admin of the package.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb_4.png" width="644" height="397" /></a></p>
<p>As I motioned, this jQuery plugin and the others can be used by this way.</p>
<p>Let&rsquo;s look at what this plugin can give us. For the sake of simplicity I&rsquo;m going to create a new ASP.NET Web Site with Visual Studio 2010;</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_5.png"><img style="background-image: none; margin: 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb_5.png" width="644" height="315" /></a><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_6.png"><img style="background-image: none; margin: 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb_6.png" width="644" height="364" /></a></p>
<p>I&rsquo;m going to fire up PMC (Package Manager Console) and install <em>jQuery.simpleJSlide.Sample </em>package;</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb_7.png" width="644" height="307" /></a></p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_8.png"><img style="background-image: none; margin: 0px 10px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" align="left" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb_8.png" width="189" height="484" /></a></p>
<p>As you can see, it gets bunch of things and set it up a sample for us. <span style="text-decoration: line-through;">Your solution should look like as it is on left hand side</span>. I made a little change on the folder structure so it doen't have look like as it is on left side. Let&rsquo;s open up the index.htm file and see how the sample code looks like;</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_9.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px initial initial;" title="image" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/image_thumb_9.png" width="539" height="423" /></a></p>
<p>Don&rsquo;t try to read the code. It is not the case here. Just looked at it and see if there is any JavaScript code inside the DOM. There is not. Then, let&rsquo;s see what this page is doing;</p>
<p>&nbsp;</p>
<div style="padding: 10px; border: 1px solid gray; width: 90%;"><iframe src="https://www.tugberkugurlu.com/Content/Docs/WebSite7/jquery.simpleJSlide-1.1/index.htm" height="380px" width="100%" frameborder="0"></iframe></div>
<p>&nbsp;</p>
<p>As you can see it is fully functional, simple image slider with play and stop buttons. Let&rsquo;s look at the code a little closer;</p>
<pre class="brush: xhtml; toolbar: false">&lt;div class="simpleJSlide" data-ui-fn="simpleJSlide" data-ui-simpleJSlide-actionactive="true" 
    data-ui-simpleJSlide-playbuttonid="play" data-ui-simpleJSlide-stopbuttonid="stop" data-ui-simpleJSlide-notificationactive="true"&gt;

    &lt;h2&gt;simpleJSlide Example Page&lt;/h2&gt;

    &lt;div id="mainImageContainer"&gt;
        &lt;img alt="image" id="containerImage" src="" /&gt;
    &lt;/div&gt;

    &lt;div id="thumbContainer"&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/1.jpg" data-ui-simpleJSlide-imgSrc="img/1.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/2.jpg" data-ui-simpleJSlide-imgSrc="img/2.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/3.jpg" data-ui-simpleJSlide-imgSrc="img/3.jpg" /&gt;&lt;br /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/4.jpg" data-ui-simpleJSlide-imgSrc="img/4.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/5.jpg" data-ui-simpleJSlide-imgSrc="img/5.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/6.jpg" data-ui-simpleJSlide-imgSrc="img/6.jpg" /&gt;&lt;br /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/7.jpg" data-ui-simpleJSlide-imgSrc="img/7.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/8.jpg" data-ui-simpleJSlide-imgSrc="img/8.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/9.jpg" data-ui-simpleJSlide-imgSrc="img/9.jpg" /&gt;&lt;br /&gt;&lt;br /&gt;

        &lt;span class="funcBtn" id="play"&gt;Play&lt;/span&gt; &lt;span class="funcBtn" id="stop"&gt;Stop&lt;/span&gt; &lt;span id="notificationContainer"&gt;&lt;/span&gt;
    &lt;/div&gt;

    &lt;div class="clearFix"&gt;&lt;/div&gt;

&lt;/div&gt;</pre>
<p>There is some conventions here;</p>
<ol>
<li>You need to specify a <strong>div</strong> element and declare <strong><em>data-ui-fn</em></strong> attribute and give it <strong>simpleJSlide</strong> as value. </li>
<li>Inside this div element, there must be two div elements whose ids are <strong><em>mainImageContainer</em></strong> and <strong><em>thumbContainer</em></strong>.</li>
<li>mainImageContainer div element must contain an <strong>img</strong> element whose id is <strong>containerImage</strong> and this image will be the container for your bigger sized images. </li>
<li>Every <strong>img</strong> element inside <strong>thumbContainer</strong> div element will be recognized as a thumbnail image. There is no way to excape that for now but maybe this can be done in the future or you can do it by yourself. </li>
<li>Every single image inside <strong>thumbContainer</strong> div element must contain <strong><em>data-ui-simpleJSlide-imgSrc</em></strong> attribute. The value of this attribute should contain bigger sized image <strong>src</strong> attribute for the particular thumbnail image as you can see above.</li>
</ol>
<p>Also, there are some options that you can override. You can specify them by simply declaring the option name after <strong><em>data- ui-simpleJSlide-</em></strong> attribute on the div element whose <strong>data-ui-fn</strong> value is simpleJSlide. Here are the options that you can override;</p>
<ul>
<li><strong>actionactive</strong> (data- ui-simpleJSlide-actionactive) as Boolean : Specifies if the images can be iterated or not. It is set to <strong>false</strong> by default.</li>
<li><strong>notificationactive</strong> (data- ui-simpleJSlide-notificationactive) as Boolean : Specifies if you want to display a massage after clicking play and stop buttons. It is set to <strong>false</strong> by default.</li>
<li><strong>playbuttonid</strong> (data- ui-simpleJSlide-playbuttonid) as String : Id of the element which will start play function by clicking for iterating images. It is set to <strong>simpleJSlidePlayButton</strong> by default.</li>
<li><strong>stopbuttonid</strong> (data- ui-simpleJSlide-stopbuttonid) as String : Id of the element which will start stop function by clicking for stopping iterating images. It is set to <strong>simpleJSlideStopButton</strong> by default.</li>
<li><strong>notificationcontainerid</strong> (data- ui-simpleJSlide-notificationcontainerid) as String : Id of the element which will hold the notification massage for play and stop notification. It is set to <strong>notificationContainer </strong>by default.</li>
<li><strong>playingtext </strong>(data- ui-simpleJSlide-playingtext) as String : Value of the notification massage which will be displayed for play function start. It is set to <strong>Now Playing&hellip;</strong> by default.</li>
<li><strong>stopedtext </strong>(data- ui-simpleJSlide-stopedtext) as String : Value of the notification massage which will be displayed for stop function start. It is set to <strong>Stopped&hellip;</strong> by default.</li>
</ul>
<p>As you see on the example, I&rsquo;ve overridden some options inside the DOM.</p>
<p><strong>Obtrusive and Total Lame Way</strong></p>
<p>Well, if you would like to use this with JavaScript instead of this unobtrusive way, you totally can. The following code shows how the above example can be done that way;</p>
<pre class="brush: xhtml; toolbar: false">&lt;!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"&gt;
&lt;html xmlns="http://www.w3.org/1999/xhtml"&gt;

&lt;head&gt;
    &lt;title&gt;&lt;/title&gt;
    &lt;link href="Content/site.css" rel="stylesheet" type="text/css" /&gt;

    &lt;script src="../Scripts/jquery-1.6.1.js" type="text/javascript"&gt;&lt;/script&gt;
    &lt;script src="../Scripts/jquery.simplejslide-1.1.min.js" type="text/javascript"&gt;&lt;/script&gt;

    &lt;script type="text/javascript"&gt;

        $(function () {

            $("#simpleJSliderDiv").simpleJSlide({
                actionactive: true,
                playbuttonid: "play",
                stopbuttonid: "stop",
                notificationactive: true
            });

        });

    &lt;/script&gt;

&lt;/head&gt;

&lt;body&gt;

&lt;div class="simpleJSlide" id="simpleJSliderDiv"&gt;

    &lt;h2&gt;simpleJSlide Example Page&lt;/h2&gt;

    &lt;div id="mainImageContainer"&gt;
        &lt;img alt="image" id="containerImage" src="" /&gt;
    &lt;/div&gt;

    &lt;div id="thumbContainer"&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/1.jpg" data-ui-simpleJSlide-imgSrc="img/1.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/2.jpg" data-ui-simpleJSlide-imgSrc="img/2.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/3.jpg" data-ui-simpleJSlide-imgSrc="img/3.jpg" /&gt;&lt;br /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/4.jpg" data-ui-simpleJSlide-imgSrc="img/4.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/5.jpg" data-ui-simpleJSlide-imgSrc="img/5.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/6.jpg" data-ui-simpleJSlide-imgSrc="img/6.jpg" /&gt;&lt;br /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/7.jpg" data-ui-simpleJSlide-imgSrc="img/7.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/8.jpg" data-ui-simpleJSlide-imgSrc="img/8.jpg" /&gt;
        &lt;img alt="image" title="Click to See Bigger" class="thumbhotel" src="img/thumb/9.jpg" data-ui-simpleJSlide-imgSrc="img/9.jpg" /&gt;&lt;br /&gt;&lt;br /&gt;

        &lt;span class="funcBtn" id="play"&gt;Play&lt;/span&gt; &lt;span class="funcBtn" id="stop"&gt;Stop&lt;/span&gt; &lt;span id="notificationContainer"&gt;&lt;/span&gt;
    &lt;/div&gt;

    &lt;div class="clearFix"&gt;&lt;/div&gt;

&lt;/div&gt;

&lt;/body&gt;
&lt;/html&gt;</pre>
<p><strong>I am not using Nuget !</strong></p>
<p>Don&rsquo;t worry. I have pushed the plugin to GitHub as well. Here is the project link : <a href="https://github.com/tugberkugurlu/jQuery-simpleJSlide">https://github.com/tugberkugurlu/jQuery-simpleJSlide</a></p>
<p>Enjoy <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/7999cbb37a60_B324/wlEmoticon-smile.png" /></p>