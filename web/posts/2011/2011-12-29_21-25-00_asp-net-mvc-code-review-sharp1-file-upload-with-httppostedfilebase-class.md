---
id: f56075ed-f701-4c14-956f-0df3a46e29fa
title: 'ASP.NET MVC Code Review #1 - File Upload With HttpPostedFileBase Class'
abstract: 'This is #1 of the series of blog posts which is about some core scenarios
  on ASP.NET MVC: File Upload With HttpPostedFileBase Class'
created_at: 2011-12-29 21:25:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- C#
- Code Review
- Razor
slugs:
- asp-net-mvc-code-review-sharp1-file-upload-with-httppostedfilebase-class
---

<p>Today, I decided to change my blogging style a little bit. I always try to write posts which contains detailed information about the subject but I realized that it prevents me from writing much more useful things. So, from now on, I will drop a bunch of code on the screen and talk about that briefly.</p>
<p>This is the beginning of the series of blog posts which is about some core scenarios on <a title="http://asp.net/mvc" href="http://asp.net/mvc" target="_blank">ASP.NET MVC</a>. In this one, code review #1, I will give you an example on how to get file upload functionality working in ASP.NET MVC. Here is the code:</p>
<p><strong>Controller:</strong></p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> SampleController : Controller {

    <span style="color: blue;">public</span> ActionResult Index() {
        <span style="color: blue;">return</span> View();
    }

    [ActionName(<span style="color: #a31515;">"Index"</span>)]
    [ValidateAntiForgeryToken, HttpPost]
    <span style="color: blue;">public</span> ActionResult Index_post(HttpPostedFileBase File) {

        <span style="color: green;">//Check if the file is not null and content length is bigger than 0</span>
        <span style="color: blue;">if</span> (File != <span style="color: blue;">null</span> &amp;&amp; File.ContentLength &gt; 0) {

            <span style="color: green;">//Check if folder is there</span>
            <span style="color: blue;">if</span>(!System.IO.Directory.Exists(Server.MapPath(<span style="color: #a31515;">"~/Content/PostedFiles"</span>)))
                System.IO.Directory.CreateDirectory(
                    Server.MapPath(<span style="color: #a31515;">"~/Content/PostedFiles"</span>)
                );

            <span style="color: green;">//Set the full path</span>
            <span style="color: blue;">string</span> path = System.IO.Path.Combine(
                Server.MapPath(<span style="color: #a31515;">"~/Content/PostedFiles"</span>),
                System.IO.Path.GetFileName(File.FileName)
            );

            <span style="color: green;">//Save the thing</span>
            File.SaveAs(path);

            TempData[<span style="color: #a31515;">"Result"</span>] = <span style="color: #a31515;">"File created successfully!"</span>;
        }

        <span style="color: green;">//RedirectToAction so that we can get rid of so-called "Form Resubmission"</span>
        <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"Index"</span>);
    }

}</pre>
</div>
</div>
<p><strong>View:</strong></p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@{
    ViewBag.Title = <span style="color: #a31515;">"File Upload Sample"</span>;
}
&lt;h2&gt;File Upload Sample&lt;/h2&gt;

@<span style="color: blue;">if</span> (TempData[<span style="color: #a31515;">"Result"</span>] != <span style="color: blue;">null</span>) { 
    &lt;ul&gt;
        &lt;li&gt;@TempData[<span style="color: #a31515;">"Result"</span>]&lt;/li&gt;
    &lt;/ul&gt;
}

@<span style="color: blue;">using</span> (Html.BeginForm(<span style="color: #a31515;">"index"</span>, <span style="color: #a31515;">"sample"</span>, 
    FormMethod.Post, <span style="color: blue;">new</span> { enctype = <span style="color: #a31515;">"multipart/form-data"</span> })) {
 
    @Html.AntiForgeryToken()
    
    &lt;input type=<span style="color: #a31515;">"file"</span> name=<span style="color: #a31515;">"File"</span> /&gt;
    &lt;p&gt;
        &lt;input type=<span style="color: #a31515;">"submit"</span> value=<span style="color: #a31515;">"Upload"</span> /&gt;
    &lt;/p&gt;
}</pre>
</div>
</div>
<p>Actually, code explains itself nicely here but I see one thing to worth pointing out here. <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.httppostedfilebasemodelbinder(v=VS.98).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.httppostedfilebasemodelbinder(v=VS.98).aspx" target="_blank">HttpPostedFileBaseModelBinder</a> class from System.Web.Mvc namespace is the class which binds a model to a posted file. So, the parameter (which is type of <a title="http://msdn.microsoft.com/en-us/library/system.web.httppostedfilebase.aspx" href="http://msdn.microsoft.com/en-us/library/system.web.httppostedfilebase.aspx" target="_blank">HttpPostedFileBase</a>) of the Index_post method receives the posted file.</p>
<p>Pretty neat stuff. Hope you enjoy it.</p>