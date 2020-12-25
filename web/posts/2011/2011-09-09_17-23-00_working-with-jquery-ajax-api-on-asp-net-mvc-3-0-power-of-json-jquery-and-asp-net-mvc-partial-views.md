---
id: 590d5018-3890-4eac-ba4c-916a7ab38f07
title: Working With JQuery Ajax API on ASP.NET MVC 3.0 - Power of JSON, JQuery and
  ASP.NET MVC Partial Views
abstract: In this post, we'll see how easy to work with JQuery AJAX API on ASP.NET
  MVC and how we make use of Partial Views to transfer chunk of html from server to
  client with a toggle button example.
created_at: 2011-09-09 17:23:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET MVC
- C#
- JQuery
slugs:
- working-with-jquery-ajax-api-on-asp-net-mvc-3-0-power-of-json-jquery-and-asp-net-mvc-partial-views
---

<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/JQuery-logo.png"><img height="244" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/JQuery-logo_thumb.png" align="right" alt="JQuery-logo" border="0" title="JQuery-logo" style="background-image: none; margin: 0px 0px 10px 10px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border-width: 0px;" /></a></p>
<p>So, here it is. Nearly every blog I occasionally read has one <a target="_blank" href="http://jquery.com/" title="http://jquery.com/">JQuery</a> related post which has JQuery<em> &lsquo;write less, do more&rsquo;</em> logo on it. Now, I have it either <img src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /> This is entirely not the intention of this blog post. What the real intention of this blog post is to lay down some useful stuff done with JQuery Ajax API on ASP.NET MVC 3.0 framework.</p>
<p>I am no expert on neither JQuery nor JavaScript but I am working on it and I hit the resources so hard nowadays. I am no expert on ASP.NET MVC either but I am sure I am pretty damn good at it. It is my favorite platform since v1.0.</p>
<p>ASP.NET MVC embraces the web (html, http, plain JavaScript) and enables you to move which direction you want on your layout. I mean you can real have your way. So, it is so important here to be able to use JavaScript <em>(for most, it is not JavaScript anymore. It is JQuery. Don&rsquo;t be surprised if you see Content-Type: text/jquery response header) </em>as Ninja in order to get the most out of ASP.NET MVC&rsquo;s V part <em>(V as in View).</em></p>
<blockquote>
<p>Easiest way to get JQuery package into your ASP.NET MVC 3.0 project is to do exactly nothing. No kidding, it is already there if you install the <a target="_blank" href="http://haacked.com/archive/2011/04/12/introducing-asp-net-mvc-3-tools-update.aspx" title="http://haacked.com/archive/2011/04/12/introducing-asp-net-mvc-3-tools-update.aspx">ASP.NET MVC 3 Tools Update</a></p>
<p>If you would like to get the latest version of it (which the below samples has been written on), simply type Update-Package JQuery inside PMC (Package Manager Console) and hit enter and you will be done !</p>
</blockquote>
<p>Let&rsquo;s get it started then.</p>
<p><strong>Overview of the Process</strong></p>
<p>Setting expectations correctly is the real deal. We need to set them correctly so that we will be as satisfied as we planned. In this quick blog post, I will walk you through a scenario: a basic to-do list application. We will:</p>
<ul>
<li>Make AJAX calls to our server, </li>
<li>Display loading messages while processing, </li>
<li>Make use of partial views when we need to update our web page seamlessly. </li>
</ul>
<p><strong>File New &gt; Project &gt; ASP.NET MVC 3 Web Application &gt; Internet Application &gt; Razor View Engine &gt; No thanks, I don&rsquo;t want the test project &gt; OK</strong></p>
<p>So the mini title above actual explains what we need to do first, briefly <img src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /> What I would like to do first is to set up my model, ORM (I will use Entity Framework for that) and database.</p>
<blockquote>
<p>The main concern of this blog post is AJAX API usage on ASP.NET MVC so I am not go through all the steps of creation our model, database and all that stuff. I will put the sample code at the bottom of this blog post so feel free to dig into that</p>
</blockquote>
<p>Then I will make sure to update my JQuery <a target="_blank" href="http://nuget.org" title="http://nuget.org">Nuget</a> package to the latest one. I will tweak the default application a little bit as well.</p>
<p><strong>Go get him tiger</strong></p>
<p>Now we are all set up and ready to rock. Let&rsquo;s begin.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image.png"><img height="199" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb.png" align="left" alt="image" border="0" title="image" style="background-image: none; margin: 0px 10px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" /></a>Our database here is so simple as you can see on the left hand side. I have created the database with <strong>SQL Server Compact</strong> edition because database is not the point we are emphasizing on this post.</p>
<p>Our application will only consist one page and we will do all the work there.</p>
<p>First, I will create the JQuery structure. As I mentioned, we will display user friendly loading message to the end user. In order to do that, we will use <strong>JQuery.UI dialog</strong> api. We will have it inside our <strong>_layout.cshtml</strong> page (the reason why we do that is to be able to use it everywhere without writing it again. This app might contains only one page but think about multiple page apps for this purpose).</p>
<p>The code for that is really simple actually. We will only have a section defined as html code inside the DOM and we will turn that into a JQuery.UI dialog for loading messages.</p>
<p>Below is the html code for our loading box :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>        <span style="color: blue;">&lt;</span><span style="color: #a31515;">div</span> <span style="color: red;">id</span><span style="color: blue;">=</span><span style="color: blue;">"ajax-progress-dialog"</span><span style="color: blue;">&gt;</span>
            <span style="color: blue;">&lt;</span><span style="color: #a31515;">div</span> <span style="color: red;">style</span><span style="color: blue;">=</span><span style="color: blue;">"margin:10px 0 0 0;text-align:center;"</span><span style="color: blue;">&gt;</span>
                <span style="color: blue;">&lt;</span><span style="color: #a31515;">img</span> <span style="color: red;">src</span><span style="color: blue;">=</span><span style="color: blue;">"@Url.Content("~/Content/img/facebook.gif")"</span> <span style="color: red;">alt</span><span style="color: blue;">=</span><span style="color: blue;">"Loading..."</span> <span style="color: blue;">/&gt;</span>
            <span style="color: blue;">&lt;/</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;/</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>Very simple html (please forgive me for putting the style code inline <img src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" />). The below JQuery code will handle this html code :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>        (<span style="color: blue;">function</span> () {
            $(<span style="color: blue;">function</span> () {
                <span style="color: green;">//Global ajax progress dialog box</span>
                <span style="color: green;">//Simply run $("#ajax-progress-dialog").dialog("open"); script before the ajax post and</span>
                <span style="color: green;">//$("#ajax-progress-dialog").dialog("close"); on the ajax post complate</span>
                $(<span style="color: #a31515;">"#ajax-progress-dialog"</span>).dialog({
                    autoOpen: <span style="color: blue;">false</span>,
                    draggable: <span style="color: blue;">false</span>,
                    modal: <span style="color: blue;">true</span>,
                    height: 80,
                    resizable: <span style="color: blue;">false</span>,
                    title: <span style="color: #a31515;">"Processing, please wait..."</span>,
                    closeOnEscape: <span style="color: blue;">false</span>,
                    open: <span style="color: blue;">function</span> () { $(<span style="color: #a31515;">".ui-dialog-titlebar-close"</span>).hide(); } <span style="color: green;">// Hide close button</span>
                });
            });
        })();</pre>
</div>
</div>
<p>With this approach, we set up an environment for ourselves to open this dialog when we start the AJAX call and close it when the action is completed.</p>
<blockquote>
<p>Some of you guys probably think that <a target="_blank" href="http://api.jquery.com/ajaxStart/" title="http://api.jquery.com/ajaxStart/">.ajaxStart()</a> and <a target="_blank" href="http://api.jquery.com/ajaxComplete/" title="http://api.jquery.com/ajaxComplete/">.ajaxComplate()</a> would be a better fit here but I don&rsquo;t. They will hook up the functions for all the AJAX calls but we might need separate loading messages. With the above approach we will call the global loading message only when we need it.</p>
</blockquote>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_3.png"><img height="152" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_3.png" align="left" alt="image" border="0" title="image" style="background-image: none; margin: 0px 20px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" /></a>Now we can move up. I have added some items to the list manually and list then on the page and we have the result as the picture on the left side.</p>
<p>It is quite simple yet but don&rsquo;t worry, it will get messy in a moment. First thing that we will do here is to add a functionality to enable toggling <strong>IsDone</strong> property of the items. When we click the anchor, the To-Do item will be marked as completed if it is not and vice versa.</p>
<p>Let&rsquo;s look at the view code here before doing that because there is some fundamental structure for our main purpose here :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@model IEnumerable&lt;JQueryAJAXToDoListMVCApp.Models.ToDoTB&gt;

@{
    ViewBag.Title = "Tugberk's To-Do List";
}

@section Head { 
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
        (<span style="color: blue;">function</span> () {
            $(<span style="color: blue;">function</span> () {
                <span style="color: green;">//JQuery code will be put here</span>
            });
        })();
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">script</span><span style="color: blue;">&gt;</span>
}

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h2</span><span style="color: blue;">&gt;</span>Tugberk's To-Do List<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h2</span><span style="color: blue;">&gt;</span>

<span style="color: blue;">&lt;</span><span style="color: #a31515;">div</span> <span style="color: red;">id</span><span style="color: blue;">=</span><span style="color: blue;">"to-do-db-list-container"</span><span style="color: blue;">&gt;</span>
    @Html.Partial("_ToDoDBListPartial")
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">div</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>As you can see have rendered a partial view here. The complete partial view code is as follows :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@model IEnumerable&lt;JQueryAJAXToDoListMVCApp.Models.ToDoTB&gt;

<span style="color: blue;">&lt;</span><span style="color: #a31515;">table</span> <span style="color: red;">style</span><span style="color: blue;">=</span><span style="color: blue;">"width:100%;"</span><span style="color: blue;">&gt;</span>

    <span style="color: blue;">&lt;</span><span style="color: #a31515;">tr</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">th</span><span style="color: blue;">&gt;</span>Item<span style="color: blue;">&lt;/</span><span style="color: #a31515;">th</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">th</span><span style="color: blue;">&gt;</span>Creation Date<span style="color: blue;">&lt;/</span><span style="color: #a31515;">th</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">th</span><span style="color: blue;">&gt;</span>IsDone<span style="color: blue;">&lt;/</span><span style="color: #a31515;">th</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">tr</span><span style="color: blue;">&gt;</span>

@foreach (var item in Model) {
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">tr</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>@item.Item<span style="color: blue;">&lt;/</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>@item.CreationDate<span style="color: blue;">&lt;/</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>@item.IsDone<span style="color: blue;">&lt;/</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">tr</span><span style="color: blue;">&gt;</span>   
}

<span style="color: blue;">&lt;/</span><span style="color: #a31515;">table</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>So why we did that is a good question in this case. When we need to make a change to our list view through an AJAX call, we will render this partial view on our controller method and send it back to the client as JSON result. This partial view will be our template here. When we make a change to this partial view, this change will be made all of the places that use this partial view which is kind of nice.</p>
<p>Also, this action will show us one of the main fundamentals of ASP.NET MVC which is <strong>separation of concerns</strong>. View doesn&rsquo;t care how the model arrives to it. It will just displays it. The same is applicable for controller as well. It does not care about how the view will display the model. It will just pass it.</p>
<p>Let&rsquo;s put the above paragraph in action. First, we need to make a little change to our partial view for toggle function :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@foreach (var item in Model) {
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">tr</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>@item.Item<span style="color: blue;">&lt;/</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>@item.CreationDate<span style="color: blue;">&lt;/</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>
            @if (@item.IsDone) { 
                <span style="color: blue;">&lt;</span><span style="color: #a31515;">a</span> <span style="color: red;">class</span><span style="color: blue;">=</span><span style="color: blue;">"isDone"</span> <span style="color: red;">href</span><span style="color: blue;">=</span><span style="color: blue;">"#"</span> <span style="color: red;">data-tododb-itemid</span><span style="color: blue;">=</span><span style="color: blue;">"@item.ToDoItemID"</span><span style="color: blue;">&gt;</span>Complated (Mark as Not Complated)<span style="color: blue;">&lt;/</span><span style="color: #a31515;">a</span><span style="color: blue;">&gt;</span>
            } else { 
                <span style="color: blue;">&lt;</span><span style="color: #a31515;">a</span> <span style="color: red;">class</span><span style="color: blue;">=</span><span style="color: blue;">"isDone"</span> <span style="color: red;">href</span><span style="color: blue;">=</span><span style="color: blue;">"#"</span> <span style="color: red;">data-tododb-itemid</span><span style="color: blue;">=</span><span style="color: blue;">"@item.ToDoItemID"</span><span style="color: blue;">&gt;</span>Not Complated (Mark as Complated)<span style="color: blue;">&lt;/</span><span style="color: #a31515;">a</span><span style="color: blue;">&gt;</span>
            }
        <span style="color: blue;">&lt;/</span><span style="color: #a31515;">td</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">tr</span><span style="color: blue;">&gt;</span>   
}</pre>
</div>
</div>
<p>With this change, we make the DOM ready for the JQuery. Let&rsquo;s look at the JQuery code for the toggle action :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>        (<span style="color: blue;">function</span> () {
            $(<span style="color: blue;">function</span> () {

                <span style="color: blue;">function</span> toggleIsDone(e, element) {

                    <span style="color: blue;">var</span> itemId = element.attr(<span style="color: #a31515;">"data-tododb-itemid"</span>);
                    <span style="color: blue;">var</span> d = <span style="color: #a31515;">"itemId="</span> + itemId;
                    <span style="color: blue;">var</span> actionURL = <span style="color: #a31515;">'@Url.Action("toogleIsDone", "ToDo")'</span>;

                    $(<span style="color: #a31515;">"#ajax-progress-dialog"</span>).dialog(<span style="color: #a31515;">"open"</span>);

                    $.ajax({
                        type: <span style="color: #a31515;">"POST"</span>,
                        url: actionURL,
                        data: d,
                        success: <span style="color: blue;">function</span> (r) {
                            $(<span style="color: #a31515;">"#to-do-db-list-container"</span>).html(r.data);
                        },
                        complete: <span style="color: blue;">function</span> () {
                            $(<span style="color: #a31515;">"#ajax-progress-dialog"</span>).dialog(<span style="color: #a31515;">"close"</span>);
                            $(<span style="color: #a31515;">".isDone"</span>).bind(<span style="color: #a31515;">"click"</span>, <span style="color: blue;">function</span> (event) {
                                toggleIsDone(event, $(<span style="color: blue;">this</span>));
                            });
                        },
                        error: <span style="color: blue;">function</span> (req, status, error) {
                            <span style="color: green;">//do what you need to do here if an error occurs</span>
                            $(<span style="color: #a31515;">"#ajax-progress-dialog"</span>).dialog(<span style="color: #a31515;">"close"</span>);
                        }
                    });

                    e.preventDefault();
                }

                $(<span style="color: #a31515;">".isDone"</span>).bind(<span style="color: #a31515;">"click"</span>, <span style="color: blue;">function</span>(e) {
                        toggleIsDone(e, $(<span style="color: blue;">this</span>));
                });

            });
        })();</pre>
</div>
</div>
<p>Let&rsquo;s bring that code down into pieces. First, we bind <strong>toggleIsDone(e, element)</strong> as a click function to every element whose class attribute is &lsquo;<strong>isDone</strong>&rsquo; and we know that they are our toggle anchors. On the <strong>toggleIsDone </strong>function, we will grab the itemId of the clicked item from <strong>data-tododb-itemid</strong> attribute of clicked anchor. Then, we will set the itemId as parameter on <strong>d</strong> variable.</p>
<p>Look at the actionURL variable. What happens there is that : we assigning the URL of <strong>toogleIsDone</strong> action of <strong>ToDo </strong>controller.</p>
<p>Before we begin our AJAX call, we simply fire up our <strong>loading dialog</strong> to display that we are doing some stuff.</p>
<p>On the AJAX call, we make a post request to <strong>actionURL </strong>and we are passing the <strong>d</strong> variable as data. We know that we will get back a JSON result which contains a <strong>data</strong> property and on success event of the call we simply change the content of <strong>to-do-db-list-container</strong> element with the new one which we have received from the server as JSON result.</p>
<p>Just before we step out the success event, we do very tricky stuff there. It is tricky in this case and hard to figure it out if you are new to JQuery. I will try to explain what we did there. We bind <strong>toggleIsDone(e, element)</strong> as a click function to every element whose class attribute is &lsquo;<strong>isDone</strong>&rsquo;. The weird thing is the fact that we have done the same just after the DOM has loaded. So, why the heck do we do that? We have bind the specified function to specified elements just after the DOM has loaded. That&rsquo;s fine. Then, we have update the entire content of to-do-db-list-container div and the all click event that we have bind has been cleared out. In order not to lose the functionality, we have bind them again.</p>
<blockquote>
<p>This kind of an Inception way is the best way that I come up with for this functionality and if you thing there is a better one, let me know.</p>
</blockquote>
<p>On complete event, we made sure to close our loading dialog.</p>
<p>At the end of the code, we call a function called <a target="_blank" href="http://api.jquery.com/event.preventDefault/" title="http://api.jquery.com/event.preventDefault/">preventDefault()</a> for the event. What this does is to prevent the anchor to do its default function which would be the append the # to the URL. Not necessary here but it is kind of nice to use here though.</p>
<p>So far, we have completed our work on client side code and lastly, we need to implement the server side function which updates the database according to parameters and sends a JSON result back to the client.</p>
<p>Before we do that we need to use Nuget here to bring down a very small package called <strong><a target="_blank" href="http://nuget.org/List/Packages/TugberkUg.MVC" title="http://nuget.org/List/Packages/TugberkUg.MVC">TugberkUg.MVC</a></strong> which will have a Controller extension for us to convert partial views to string inside the controller.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_4.png"><img height="75" width="640" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_4.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_5.png"><img height="180" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_5.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>The complete code of out ToDoController is as indicated below :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">using</span> System.Linq;
<span style="color: blue;">using</span> System.Web.Mvc;
<span style="color: blue;">using</span> TugberkUg.MVC.Helpers;
<span style="color: blue;">using</span> System.Data.Objects;
<span style="color: blue;">using</span> System.Data;

<span style="color: blue;">namespace</span> JQueryAJAXToDoListMVCApp.Controllers
{
    <span style="color: blue;">public</span> <span style="color: blue;">class</span> ToDoController : Controller {

        <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> Models.ToDoDBEntities _entities = <span style="color: blue;">new</span> Models.ToDoDBEntities();

        <span style="color: blue;">public</span> ActionResult Index() {

            <span style="color: blue;">var</span> model = _entities.ToDoTBs;

            <span style="color: blue;">return</span> View(model);
        }

        [HttpPost]
        <span style="color: blue;">public</span> ActionResult toogleIsDone(<span style="color: blue;">int</span> itemId) {

            <span style="color: green;">//Getting the item according to itemId param</span>
            <span style="color: blue;">var</span> model = _entities.ToDoTBs.FirstOrDefault(x =&gt; x.ToDoItemID == itemId);
            <span style="color: green;">//toggling the IsDone property</span>
            model.IsDone = !model.IsDone;

            <span style="color: green;">//Making the change on the db and saving</span>
            ObjectStateEntry osmEntry = _entities.ObjectStateManager.GetObjectStateEntry(model);
            osmEntry.ChangeState(EntityState.Modified);
            _entities.SaveChanges();

            <span style="color: blue;">var</span> updatedModel = _entities.ToDoTBs;

            <span style="color: green;">//returning the new template as json result</span>
            <span style="color: blue;">return</span> Json(<span style="color: blue;">new</span> { data = <span style="color: blue;">this</span>.RenderPartialViewToString(<span style="color: #a31515;">"_ToDoDBListPartial"</span>, updatedModel) });
        }

        <span style="color: blue;">protected</span> <span style="color: blue;">override</span> <span style="color: blue;">void</span> Dispose(<span style="color: blue;">bool</span> disposing) {

            _entities.Dispose();
            <span style="color: blue;">base</span>.Dispose(disposing);
        }

    }
}</pre>
</div>
</div>
<p>When you look inside the <strong>toogleIsDone</strong> method, you will see that after the necessary actions are completed, we will pass a model to our partial view which we have created earlier, render it and finally return it to the client as JSON. Why we return it as JSON instead of content is a totally subjective question in my opinion. For me, the advantage of JSON result is to be able to pass multiple values. For example, we would pass the result as follows if we needed :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>            <span style="color: blue;">return</span> Json(
                <span style="color: blue;">new</span> { 
                    CustomMessage = <span style="color: #a31515;">"My message"</span>, 
                    data = <span style="color: blue;">this</span>.RenderPartialViewToString(<span style="color: #a31515;">"_ToDoDBListPartial"</span>, updatedModel) 
                });</pre>
</div>
</div>
<p>Now when we compile our project and run it, we will see a working application :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_6.png"><img height="142" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_6.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a>&nbsp;<a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_7.png"><img height="142" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_7.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a>&nbsp;<a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_8.png"><img height="142" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/image_thumb_8.png" alt="image" border="0" title="image" style="background-image: none; margin: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Also, here is a quick video of a working example :</p>
<div><iframe width="425" height="349" src="http://www.youtube.com/embed/_AzThF73y5M?hl=en&amp;fs=1" frameborder="0"></iframe></div>
<p><strong>Summary</strong></p>
<p>What we need to get out of from this blog post is to see that what we are capable of doing with a little effort for very usable and user friendly web pages. Also, we can visualize how things fit together and flow on your project. There are lots of way of making calls to your server from client side code and lots of them have its own pros and cons.</p>
<p>I hope that this blog post helps you, even a little, to solve a real world problem with your fingers <img src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/222f24832a52_124A3/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /></p>
<div><iframe title="Preview" scrolling="no" marginheight="0" marginwidth="0" frameborder="0" style="width: 98px; height: 115px; padding: 0; background-color: #fcfcfc;" src="https://skydrive.live.com/embedicon.aspx/Programming/JQueryAJAXToDoListMVCApp.rar?cid=0ee89cb310fe3603&amp;sc=documents"></iframe></div>