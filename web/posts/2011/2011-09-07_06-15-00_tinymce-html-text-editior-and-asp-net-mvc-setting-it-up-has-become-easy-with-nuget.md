---
id: 1f7cf9db-2a48-4719-8b19-9cc04182144a
title: TinyMCE HTML Text Editior & ASP.NET MVC - Setting It Up Has Become Easy With
  Nuget
abstract: One of the best Javascript WYSIWYG Editors TinyMCE is now up on Nuget live
  feed. How to get TinyMCE through Nuget and get it working is documented in this
  blog post.
created_at: 2011-09-07 06:15:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- C#
- JQuery
- NuGet
- Razor
slugs:
- tinymce-html-text-editior-and-asp-net-mvc-setting-it-up-has-become-easy-with-nuget
- tinymce-html-text-editior-and-asp-net-mvc-setting-it-up-has-bec
---

<p><strong><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/tinymce-logo.jpg"><img style="background-image: none; margin: 0px 20px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="tinymce-logo" border="0" alt="tinymce-logo" align="left" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/tinymce-logo_thumb.jpg" width="236" height="121" /></a>Overview of the Process</strong></p>
<p>Couple of weeks ago, I was setting up a new project and I have created a folder called &lsquo;lib&rsquo; inside the <a title="Macros for Build Commands and Properties" href="http://msdn.microsoft.com/en-us/library/c02as0cs(v=vs.71).aspx" target="_blank">$(SolutionDir)</a> as usual to put all the dependencies that I am going to be using for my project (.Net libraries, JavaScript libraries, etc.).</p>
<p>Then, I go to <a href="http://www.tinymce.com/">http://www.tinymce.com/</a> to grab the latest version of <a title="TinyMCE - Javascript WYSIWYG Editor" href="http://www.tinymce.com/" target="_blank">TinyMCE</a> which is an awesome Javascript <a title="http://en.wikipedia.org/wiki/WYSIWYG" href="http://en.wikipedia.org/wiki/WYSIWYG" target="_blank">WYSIWYG</a> Editor. This action popped up a balloon inside my head that tells me I had been doing this stuff over and over again for 80% of the projects that I have created. So this hits me hard and I thought I should automate this process. In order to do that, I&rsquo;ve created an internal <a title="http://Nuget.org" href="http://Nuget.org" target="_blank">Nuget</a> package just for myself to pull the TinyMCE package with plugins. Even, I have created EditorTemplates so that I could just give my model a hint to use the template.</p>
<p>That worked pretty well for <strong>me</strong> but don&rsquo;t miss the point here. That worked pretty well <strong>only for me</strong>. Not for John, Denis, Adam and so on. After I thought this, I have come across a blog post about <a title="http://www.yakupbugra.com/2011/07/asp-net-mvc-3-tinymce-editor-kullanimi.html" href="http://www.yakupbugra.com/2011/07/asp-net-mvc-3-tinymce-editor-kullanimi.html" target="_blank">TinyMCE integration on ASP.NET MVC 3</a> (The blog post is in Turkish). Remember the balloon which popped up inside my head? This blog post widened that balloon and it hit me harder this time. The process which has been documented there is a very well working sample but well&hellip;, it looked more like a poor man&rsquo;s TinyMCE integration for me. (The article was nice but it wasn&rsquo;t following the DRY way)</p>
<p>After that, I have decided to enhance my packages to push them to live Nuget feed. So, I have contacted with the package admins on their forum :</p>
<p><a title="http://www.tinymce.com/forum/viewtopic.php?id=26568" href="http://www.tinymce.com/forum/viewtopic.php?id=26568" target="_blank">TinyMCE Nuget Package</a></p>
<p><a title="http://twitter.com/Spocke" href="http://twitter.com/Spocke" target="_blank">@Spocke</a> has replied my post in a short time and gave me a go ahead to push the package to live Nuget feed. I am not going to get into details of how I created the package. Mostly, I will show how to set up your environment to get TinyMCE working in a short time.</p>
<p><strong>TinyMCE into an ASP.NET MVC project</strong></p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/tinymce-Nuget-packages-of-mine.png"><img style="background-image: none; margin: 0px 20px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border: 0px;" title="tinymce-Nuget-packages-of-mine" border="0" alt="tinymce-Nuget-packages-of-mine" align="left" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/tinymce-Nuget-packages-of-mine_thumb.png" width="229" height="244" /></a></p>
<p>It is now easy to get TinyMCE package to your project with Nuget. It&rsquo;s even easier to get it working with ASP.NET MVC. In this post, I am going show you the easiest way to get it work ASP.NET MVC but I hope I am going to cover for Web Forms as well in near future.</p>
<p>There are several packages which are related to TinyMCE that I have pushed to live Nuget feed as you can see in the picture left hand side. (This list might extend in the future) Briefly the packages are :</p>
<p><a title="TinyMCE - Javascript WYSIWYG Editor" href="http://Nuget.org/List/Packages/TinyMCE" target="_blank">TinyMCE</a> : The main package which holds the infrastructure .js files of the library and plugings. This package has <strong>no dependency</strong> at all.</p>
<p><a title="http://Nuget.org/List/Packages/TinyMCE.JQuery" href="http://Nuget.org/List/Packages/TinyMCE.JQuery" target="_blank">TinyMCE.JQuery</a> : Holds the <a title="http://jquery.com/" href="http://jquery.com/" target="_blank">JQuery</a> integration files for TinyMCE. This package depends on <strong>TinyMCE</strong> and <strong>JQuery</strong> packages.</p>
<p><a title="http://Nuget.org/List/Packages/TinyMCE.MVC" href="http://Nuget.org/List/Packages/TinyMCE.MVC" target="_blank">TinyMCE.MVC</a> : This package holds the ASP.NET MVC EditorTemplates for TinyMCE. This package depends on <strong>TinyMCE</strong> package.</p>
<p><a title="http://Nuget.org/List/Packages/TinyMCE.MVC.JQuery" href="http://Nuget.org/List/Packages/TinyMCE.MVC.JQuery" target="_blank">TinyMCE.MVC.JQuery</a> : Holds the ASP.NET MVC EditorTemplates for TinyMCE.JQuery. This package depends on <strong>TinyMCE.JQuery</strong> package.</p>
<p><a title="http://Nuget.org/List/Packages/TinyMCE.MVC.Sample" href="http://Nuget.org/List/Packages/TinyMCE.MVC.Sample" target="_blank">TinyMCE.MVC.Sample</a> : Holds a sample ASP.NET MVC mini Application (it is all just a model, a controller and a view) for TinyMCE. This package depends on <strong>TinyMCE.MVC</strong> package so it uses EditorTemplates to illustrate a sample.</p>
<p><a title="http://Nuget.org/List/Packages/TinyMCE.MVC.JQuery.Sample" href="http://Nuget.org/List/Packages/TinyMCE.MVC.JQuery.Sample" target="_blank">TinyMCE.MVC.JQuery.Sample</a> : This package holds a sample ASP.NET MVC mini Application (it is all just a model, a controller and a view) for TinyMCE.JQuery. This package depends on <strong>TinyMCE.MVC.JQuery</strong> package so it uses EditorTemplates to illustrate a sample.</p>
<p><strong>How to get TinyMCE work on and ASP.NET MVC project</strong></p>
<p>I would like set the boundaries here :</p>
<ul>
<li>This will demonstrate the process of integrating TinyMCE to an ASP.NET MVC application.</li>
<li>I will use the JQuery one because it has more downloads on Nuget <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/wlEmoticon-smile.png" /></li>
<li>I will work with the TinyMCE.JQuery.MVC package to pull down all the TinyMCE stuff along with ASP.NET MVC EditorTemplates.</li>
</ul>
<p>Always first thing to go is File &gt; New Project &gt; ASP.NET MVC 3 Web Application.</p>
<p>In order to make it more clear, I have created a very simple BlogPost model to work with. First, I am going to show you the way without TinyMCE then later we will bring it down through Nuget.</p>
<p>The model looks like as below :</p>
<pre class="brush: c-sharp; toolbar: false">    public class BlogPost {

        public string Title { get; set; }
        public DateTime PostedOn { get; set; }
        public string Tags { get; set; }
        public string Content { get; set; }

    }</pre>
<p>Then I created a controller to create a new blog post (assume here this model is related to database). The controller is simple :</p>
<pre class="brush: c-sharp; toolbar: false">using System.Web.Mvc;
using TinyMCEJQueryMVCNuget.Models;

namespace TinyMCEJQueryMVCNuget.Controllers {

    public class BlogPostController : Controller {

        public ActionResult Create() {

            return View();
        }

        [HttpPost, ActionName("Create")]
        public ActionResult Create_post(BlogPost model) {

            ViewBag.HtmlContent = model.Content;

            return View(model);
        }

    }
}</pre>
<p>It basically does noting. Just getting the model on the http post event and pushing it back to the view with a ViewBag property.</p>
<p>Here is what some portion of our view looks like :</p>
<pre class="brush: xhtml; toolbar: false">@using (Html.BeginForm()) {
    
    @Html.ValidationSummary(true)

    &lt;fieldset&gt;
        &lt;legend&gt;BlogPost&lt;/legend&gt;

        &lt;div class="editor-label"&gt;
            @Html.LabelFor(model =&gt; model.Title)
        &lt;/div&gt;
        &lt;div class="editor-field"&gt;
            @Html.EditorFor(model =&gt; model.Title)
            @Html.ValidationMessageFor(model =&gt; model.Title)
        &lt;/div&gt;

        &lt;div class="editor-label"&gt;
            @Html.LabelFor(model =&gt; model.PostedOn)
        &lt;/div&gt;
        &lt;div class="editor-field"&gt;
            @Html.EditorFor(model =&gt; model.PostedOn)
            @Html.ValidationMessageFor(model =&gt; model.PostedOn)
        &lt;/div&gt;

        &lt;div class="editor-label"&gt;
            @Html.LabelFor(model =&gt; model.Tags)
        &lt;/div&gt;
        &lt;div class="editor-field"&gt;
            @Html.EditorFor(model =&gt; model.Tags)
            @Html.ValidationMessageFor(model =&gt; model.Tags)
        &lt;/div&gt;

        &lt;div class="editor-label"&gt;
            @Html.LabelFor(model =&gt; model.Content)
        &lt;/div&gt;
        &lt;div class="editor-field"&gt;
            @Html.EditorFor(model =&gt; model.Content)
            @Html.ValidationMessageFor(model =&gt; model.Content)
        &lt;/div&gt;

        &lt;p&gt;
            &lt;input type="submit" value="Create" /&gt;
        &lt;/p&gt;

        &lt;p&gt;
            Posted Content : @ViewBag.HtmlContent
        &lt;/p&gt;

    &lt;/fieldset&gt;
}</pre>
<p>When we open it up, we will see nothing fancy there :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_thumb.png" width="644" height="438" /></a></p>
<p>When we post it, we will get this result :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_thumb_3.png" width="644" height="438" /></a></p>
<p>So, we got this baby working. Let&rsquo;s improve it. What we need here is to give the blogger a real blogging experience. In order to do that, we need a text editor so that we could enter our content pretty easily.</p>
<p>First thing is to pull down the TinyMCE.JQuery.MVC package over the wire.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/install-package.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="install-package" border="0" alt="install-package" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/install-package_thumb.png" width="644" height="75" /></a></p>
<p>When you start installing the package, your solution will have a movement so don&rsquo;t freak out <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/wlEmoticon-smile.png" /></p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_thumb_4.png" width="644" height="349" /></a></p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_5.png"><img style="background-image: none; margin: 0px 20px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" align="left" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_thumb_5.png" width="217" height="244" /></a>So now we have our package installed, we are ready to make some changes on our model. When you go to ~/Views/Shared/EditorTemplates folder on your solution, you will see that there is a cshtml file there called tinymce_jquery_full.cshtml. This partial view enables you to view your model property with TinyMCE editor.</p>
<p>I am not going to go inside this file and explain how to do this (however, it is pretty simple). It is entirely an another topic.</p>
<p>What I would like to point out here is this : if you are working with this package (TinyMCE.JQuery.MVC), you need to have JQuery referenced on your file. We have our JQuery referenced on our _Layout.cshtml file so we do not need to do anything else.</p>
<p>As I got the information from the project admins, TinyMCE can work with JQuery 1.4.3 and after. You don&rsquo;t need to worry about that as well. Nuget will resolve the dependency without you knowing.</p>
<p><strong>Change to Our Model</strong></p>
<p>To get it work, we need to change our model as follows :</p>
<pre class="brush: c-sharp; toolbar: false">using System;
using System.ComponentModel.DataAnnotations;
using System.Web.Mvc;

namespace TinyMCEJQueryMVCNuget.Models {

    public class BlogPost {

        public string Title { get; set; }
        public DateTime PostedOn { get; set; }
        public string Tags { get; set; }

        [UIHint("tinymce_jquery_full"), AllowHtml]
        public string Content { get; set; }

    }
}</pre>
<p>What UIHint attribute does here is to tell the framework to use tinymce_jquery_full editor template instead of the default one. AllowHtml is also there to allow html string to pass from your browser to the server. Build you project and run it again and then you should see a page similar to following one:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_6.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_thumb_6.png" width="644" height="438" /></a></p>
<blockquote>
<p>If you are unable to see this when you run your application, this can be related to several things. First thing I would do is to check if the Content property is used with EditorFor instead of another Html Helpers :</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [5]">        &lt;div class="editor-label"&gt;
            @Html.LabelFor(model =&gt; model.Content)
        &lt;/div&gt;
        &lt;div class="editor-field"&gt;
            @Html.EditorFor(model =&gt; model.Content)
            @Html.ValidationMessageFor(model =&gt; model.Content)
        &lt;/div&gt;</pre>
<p>The second thing would be a JavaScript error related to misconfiguration or a wrong library reverence. Check your browser console for any JavaScript errors for that.</p>
</blockquote>
<p>Now when we post it back to server, we will get it back as :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/91cdfb7fb32a_6EB8/image_thumb_7.png" width="644" height="438" /></a></p>
<p>Very nice. Html has been created for us behind the scenes.</p>
<blockquote>
<p>It displays it as encoded html because we didn&rsquo;t specify to view as raw html. In order to do that, just use <a title="http://haacked.com/archive/2011/01/06/razor-syntax-quick-reference.aspx" href="http://haacked.com/archive/2011/01/06/razor-syntax-quick-reference.aspx" target="_blank">@Html.Raw</a>.</p>
</blockquote>
<p>I hope this will help you to automate some portion of your project as well.</p>