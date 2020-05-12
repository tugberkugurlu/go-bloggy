---
id: 56287082-0f17-4b95-b8c1-fc7df92a6f90
title: How To Handle Multiple Checkboxes From Controller In ASP.NET MVC - Sample App
  With New Stuff After MVC 3 Tools Update
abstract: 'In this blog post, we will see how to handle multiple checkboxes inside
  a controller in ASP.NET MVC. We will demonstrate a sample in order to delete multiple
  records from database. '
created_at: 2011-04-27 11:10:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET MVC
- C#
- MvcScaffolding
- NuGet
slugs:
- how-to-handle-multiple-checkboxes-from-controller-in-asp-net-mvc-sample-app-with-new-stuff-after-mvc-3-tools-update
---

<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/checkbox.jpg"><img style="background-image: none; margin: 0px 15px 15px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="checkbox" border="0" alt="checkbox" align="left" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/checkbox_thumb.jpg" width="244" height="185" /></a></p>
<p>Checkboxes are so <strong>handy </strong>and easy to use. But ASP.Net MVC is a framework which stands http verbs so there is no <em>OnClick</em> event for buttons. Actually, there are <em><strong>no buttons</strong></em>! As you know buttons are rendered as <em>&lt;input value=&rdquo;submit&rdquo; /&gt; </em>native html code and .net assign some values for those input to handle them on the http post. Handling checkboxes with web forms are easy along with code behind. There is no headache there <em>(There is no headache in MVC for me but some think that way)</em></p>
<p>But what about ASP.NET MVC? How can we handle checkboxes, which are inside a view, from controller? The solution is easier than you imagine. In this blog post, I will demonstrate a simple scenario about how to delete multiple records with the help of checkboxes.</p>
<p>In our sample app, I am gonna be using new things on .Net such as <a title="http://nuget.org/List/Packages/EFCodeFirst" href="http://nuget.org/List/Packages/EFCodeFirst" target="_blank">EFCodeFirst</a>, <a title="http://nuget.org/List/Packages/EntityFramework.SqlServerCompact" href="http://nuget.org/List/Packages/EntityFramework.SqlServerCompact" target="_blank">EntityFramework.SqlServerCompact</a>, <a title="http://nuget.org/List/Packages/MvcScaffolding" href="http://nuget.org/List/Packages/MvcScaffolding" target="_blank">MvcScaffolding</a>. The idea is simple : we will have a simple contact form and a backend to manage those messages.</p>
<blockquote>
<p>At its core, this blog post intends to show the way of handling multiple checkboxes on ASP.NET MVC but it is also a good example of how to create an ASP.NET MVC app from scratch. As you will see, we will see MvcScaffolding many times. I'll go through the basics of it but for more information, there is a blog post series on <a target="_blank" title="http://blog.stevensanderson.com" href="http://blog.stevensanderson.com">Steve Sanderson</a>'s blog:</p>
<ol>
<li><a target="_blank" href="http://blog.stevensanderson.com/2011/01/13/scaffold-your-aspnet-mvc-3-project-with-the-mvcscaffolding-package/">Introduction: Scaffold your ASP.NET MVC 3 project with the MvcScaffolding package</a> </li>
<li><a target="_blank" href="http://blog.stevensanderson.com/2011/01/13/mvcscaffolding-standard-usage/">Standard usage: Typical use cases and options</a> </li>
<li><a target="_blank" href="http://blog.stevensanderson.com/2011/01/28/mvcscaffolding-one-to-many-relationships/">One-to-Many Relationships</a> </li>
<li><a target="_blank" href="http://blog.stevensanderson.com/2011/03/28/scaffolding-actions-and-unit-tests-with-mvcscaffolding/">Scaffolding Actions and Unit Tests</a></li>
<li><a target="_blank" href="http://blog.stevensanderson.com/2011/04/06/mvcscaffolding-overriding-the-t4-templates/">Overriding the T4 templates</a> </li>
<li><a target="_blank" href="http://blog.stevensanderson.com/2011/04/07/mvcscaffolding-creating-custom-scaffolders/">Creating custom scaffolders</a> </li>
<li><a target="_blank" href="http://blog.stevensanderson.com/2011/04/08/mvcscaffolding-scaffolding-custom-collections-of-files/">Scaffolding custom collections of files</a></li>
</ol></blockquote>
<p>Let the code begin&hellip;</p>
<p>Fist thing is always first and this thing is File &gt; New &gt; Project on Visual Studio world. I am gonna create a new ASP.NET MVC 3 Internet Application with Razor View Engine. I am also going to create a Unit Test project related to my MVC project because I am not a bad person (on <a title="http://hanselman.com/" href="http://hanselman.com/" target="_blank">Scott Hanselman</a>&rsquo;s point of view <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/wlEmoticon-smile.png" />)</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image.png"><img style="background-image: none; margin: 0px 20px 0px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb.png" width="244" height="220" /></a><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_3.png" width="149" height="244" /></a></p>
<p>BTW, Don&rsquo;t worry, I am not gonna put the three circles here which explain <strong>CLEARLY</strong> what MVC framework is and blows our minds<img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/wlEmoticon-smile.png" /></p>
<p><strong>Nuget Packages is First Thing to Go With</strong></p>
<p>But before getting our hands dirty, bring up the PMC (which stands for Package Manager Console [<em>I made it up but it might be the actual abbreviation. <img style="border-style: none;" class="wlEmoticon wlEmoticon-openmouthedsmile" alt="Open-mouthed smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/wlEmoticon-openmouthedsmile.png" /> If it is not, use it to make it legit!</em>]) on Visual Studio 2010 to get the following packages;</p>
<p><strong><em>MvcScaffolding</em> : </strong>A fast and customizable way to add controllers, views, and other items to your ASP.NET MVC application</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_4.png" width="640" height="80" /></a></p>
<p><strong>EntityFramework.SqlServerCompact : </strong>Allows SQL Server Compact 4.0 to be used with Entity Framework..</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_5.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_5.png" width="640" height="83" /></a></p>
<blockquote>
<p>I am here using <a title="http://haacked.com/archive/2011/04/12/introducing-asp-net-mvc-3-tools-update.aspx" href="http://haacked.com/archive/2011/04/12/introducing-asp-net-mvc-3-tools-update.aspx" target="_blank">ASP.NET MVC 3 Tools Update</a> so I got <a title="http://nuget.org/List/Packages/EntityFramework" href="http://nuget.org/List/Packages/EntityFramework" target="_blank">EntityFramework</a> Nuget Package out of the box for support of Code First workflow for ADO.NET Entity Framework. If you are not using this update, simple install EntityFramework package from PMC or Add Library Package Reference dialog box on VS.</p>
</blockquote>
<p>Now, we have built our project environment and ready to hit the code.</p>
<p><strong>Be Cool and Use the PowerShell</strong></p>
<p>I am no PowerShell expert. Even, I am no good at that thing. But I am enthusiastic about that and I am learning it. So in our sample project, I will use MvcScaffolding to create nearly everything that I can. I have two options here;</p>
<p>1 &ndash; To use GUI to scaffold my crap. This option comes with ASP.NET MVC 3 Tools Update. I am not going to be using that option.</p>
<p>2 &ndash; By using PMC which is a lot cooler than the GUI. I&rsquo;ll be using this option along the way.</p>
<p>Firstly, we need our model classes. Below code includes all of model classes that we need for this small project;</p>
<pre class="brush: c-sharp; toolbar: false">using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;

namespace ContactFormWithMultipleCheckboxApp.Models {

    public class Product {

        public int ProductId { get; set; }
        [Required, StringLength(50)]
        public string ProductName { get; set; }
        public string Description { get; set; }

        public virtual ICollection&lt;Message&gt; Messages { get; set; }

    }

    public class Message {

        public int MessageId { get; set; }
        public string From { get; set; }
        [Required]
        //below one is to validate whether the e-mail address is legit or not
        [RegularExpression(@"\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*", ErrorMessage = "Email address is not valid")]
        public string Email { get; set; }
        [StringLength(100)]
        public string Subject { get; set; }
        public string Content { get; set; }

        public int ProductId { get; set; }
        public Product Product { get; set; }

    }
}</pre>
<p><strong>Scaffolding Controllers, Views and Repositories, DBContext (Entity Framework Code First Database Context Class)</strong></p>
<p>To do so, bring up the PMC and type the following code;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_6.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_6.png" width="644" height="206" /></a></p>
<p>Nice ! We have what we need so far. But not so useful. When we open up the product controller file, we will see that controller actions directly talking to Entity Framework which is not so handy in terms of Unit Testing. What can we do here. MvcScaffolding has a nice solution for that. We will just rescaffold the controller for Product class and at the same time we will create repositories. But by default, MvcScaffolding doesn&rsquo;t overwrite any files that already exist. If you do want to overwrite things, you need to pass &ndash;force switch. To do so, we will do the following PS command on PMC;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_7.png" width="644" height="206" /></a></p>
<p>Cool. We have our core logic for backend system. Do the same for Message class with following one line of PS code;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_8.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_8.png" width="644" height="206" /></a></p>
<p>Rock on ! Perfect so far. When we fire up our app, Sql Compact database will be created for us inside our App_Data folder. Now go to the <strong>Products </strong>directory and add some products and then go to the <strong>Messages</strong> directory on the address bar.</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_9.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_9.png" width="644" height="232" /></a></p>
<blockquote>
<p>There is one thing that I would like to share here. When you go to your <strong>_CreateOrEdit.cshtml</strong> file inside <strong>Views/Messages</strong> folder, you will see that the following code used for creating relational dropdownlist.</p>
<pre class="brush: xhtml; toolbar: false">&lt;div class="editor-field"&gt;
    @Html.DropDownListFor(model =&gt; model.ProductId, ((IEnumerable&lt;ContactFormWithMultipleCheckboxApp.Models.Product&gt;)ViewBag.PossibleProducts).Select(option =&gt; new SelectListItem {
        Text = Html.DisplayTextFor(_ =&gt; option).ToString(), 
        Value = option.ProductId.ToString(),
        Selected = (Model != null) &amp;&amp; (option.ProductId == Model.ProductId)
    }), "Choose...")
    @Html.ValidationMessageFor(model =&gt; model.ProductId)
&lt;/div&gt;</pre>
<p>Well, this code does not work quite fine. I don&rsquo;t know why and I didn&rsquo;t dig into that.&nbsp; I just changed the code with the following one and now it works quite well. Nevertheless, MvcScaffolding has done its part very well.</p>
<pre class="brush: xhtml; toolbar: false">&lt;div class="editor-field"&gt;
    @Html.DropDownListFor(model =&gt; model.ProductId, ((IEnumerable&lt;ContactFormWithMultipleCheckboxApp.Models.Product&gt;)ViewBag.PossibleProducts).Select(option =&gt; new SelectListItem {
        Text = option.ProductName, 
        Value = option.ProductId.ToString(),
        Selected = (Model != null) &amp;&amp; (option.ProductId == Model.ProductId)
    }), "Choose...")
    @Html.ValidationMessageFor(model =&gt; model.ProductId)
&lt;/div&gt;</pre>
</blockquote>
<p>Now, we are at the part where we should implement the logic of our core idea : <strong>deleting multiple records from the database</strong>. I have created several Messages now and our list looks like as below;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_10.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_10.png" width="644" height="359" /></a></p>
<p>Hmmm, looks like Fraud has been spaming us with very old-fashioned way. Well, better delete them from database. They are garbage, right? But that&rsquo;d be nice to delete them all by one click. Let&rsquo;s make that happen <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/wlEmoticon-smile.png" />&nbsp;</p>
<p>Open up the <strong>Index.cshtml</strong> file inside <strong>Views/Messages</strong> folder. See the change what I have done on the lines 4,5,6;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [4,5,6]">@foreach (var item in Model) {
    
    &lt;tr&gt;
        &lt;td&gt;
            &lt;input type="checkbox" name="deleteInputs" value="@item.MessageId" /&gt;
        &lt;/td&gt;
        &lt;td&gt;
            @Html.ActionLink("Edit", "Edit", new { id=item.MessageId }) |
            @Html.ActionLink("Details", "Details", new { id=item.MessageId }) |
            @Html.ActionLink("Delete", "Delete", new { id=item.MessageId })
        &lt;/td&gt;
        &lt;td&gt;
			@item.From
        &lt;/td&gt;
        &lt;td&gt;
			@item.Email
        &lt;/td&gt;
        &lt;td&gt;
			@item.Subject
        &lt;/td&gt;
        &lt;td&gt;
			@item.Content
        &lt;/td&gt;
        &lt;td&gt;
			@Html.DisplayTextFor(_ =&gt; item.Product.ProductName)
        &lt;/td&gt;
    &lt;/tr&gt;
    
}</pre>
<p>I have added an input whose type attribute is <strong><em>checkbox</em></strong> for every generated record. So what does this think do? It enable us to catch the selected checkboxes and relate them to a record with the <strong>MessageId</strong> property. But we are not done yet. we will create a PostAction method for Index view but we are not gonna do that by typing. We will use PMC again because it&rsquo;s cooler <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/wlEmoticon-smile.png" /> Type the following;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_11.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_11.png" width="644" height="206" /></a></p>
<p>Now we have our <strong>Post Action Method. </strong>The last view of this method will look like this;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [26]">        [HttpPost, ActionName("Index")]
        public ActionResult IndexPost(int[] deleteInputs) {

            var model = messageRepository.AllIncluding(message =&gt; message.Product);

            if (deleteInputs == null) {

                ModelState.AddModelError("", "None of the reconds has been selected for delete action !");
                return View(model);
            }

            foreach (var item in deleteInputs) {

                try {

                    messageRepository.Delete(item);

                } catch (Exception err) {

                    ModelState.AddModelError("", "Failed On Id " + item.ToString() + " : " + err.Message);
                    return View(model);
                }
            }

            messageRepository.Save();
            ModelState.AddModelError("", deleteInputs.Count().ToString() + " record(s) has been successfully deleted!");

            return View(model);
        }</pre>
<p>So here what we are doing is simple. We are retrieving selected checkboxes and deleting the records from database according to <strong>MessageId</strong> property. Also, I am doing something bizarre here on the 26. I am passing the confirmation message to ModelError to view it afterwards. Probably not the smartest thing on the planet but it works for me.</p>
<p>Finally, add the following code to your index view page;</p>
<pre class="brush: xhtml; toolbar: false">@using (Html.BeginForm()) {
&lt;div&gt;

    @Html.ValidationSummary()

&lt;/div&gt;
&lt;div style="float:right;"&gt;

        
        &lt;input type="submit" value="Delete Selected Values" /&gt;
        

&lt;/div&gt;

@*Add this below curly brace at the end of the page to close the using statement for Html.Beginform()*@
}</pre>
<p>Fire up the project by clicking Ctrl + F5 and go to the <strong>Messages</strong> directory. You will end up with a web page which is similar to following left one. Then, select the messages which have come from Fraud who is our old-fashioned spammer. After selecting the records you would like to delete, click the button on top-right side.</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_12.png"><img style="background-image: none; margin: 0px 30px 0px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_12.png" width="244" height="133" /></a><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_13.png"><img style="background-image: none; margin: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_13.png" width="244" height="134" /></a></p>
<p>Bingo ! We are done here.</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_14.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/image_thumb_14.png" width="644" height="231" /></a></p>
<p>Of course this way is not the nicest way for this type of action. The better solution would be doing this with Ajax or JQuery post to server and removing deleted record rows with nice animation.</p>
<p>So why didn&rsquo;t I do that here? Because I don&rsquo;t have enough time for now <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/wlEmoticon-smile.png" /> I wanna get this post up and running quickly. So, this is it for now <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/2fc35edc2e14_AE5E/wlEmoticon-smile.png" /> Feel free to download the source code and use it anywhere you like.</p>
<p><iframe src="http://cid-0ee89cb310fe3603.office.live.com/embedicon.aspx/Programming/ContactFormWithMultipleCheckboxApp.rar" style="width: 98px; height: 115px; padding: 0; background-color: #fcfcfc;" frameborder="0" marginwidth="0" marginheight="0" scrolling="no" title="Preview"></iframe></p>