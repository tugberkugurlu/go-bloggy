---
id: e71fadd8-aeb2-4b62-8436-865bf990d884
title: 'Exciting Things About ASP.NET vNext Series: MVC View Components'
abstract: 'A few days ago, I started a new blog post series about ASP.NET vNext. Today,
  I would like to talk about something which is MVC specific and takes one of our
  pains away: view components :)'
created_at: 2014-10-06 11:46:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET MVC
- ASP.NET vNext
slugs:
- exciting-things-about-asp-net-vnext-series-mvc-view-components
---

<p>Web development experience with .NET has never seen a drastic change like this since its birth day. Yes, I’m talking about <a href="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-vnext-by-setting-up-the-environment-from-scratch">ASP.NET vNext</a> :) I have been putting my toes into this water for a while now and a few days ago, I started a new blog post series about ASP.NET vNext (with hopes that I will continue this time :)). To be more specific, I’m planning on writing about the things I am actually excited about this new cloud optimized (TM) runtime. Those things could be anything which will come from <a href="http://github.com/aspnet">ASP.NET GitHub account</a>: things I like about the development process, Visual Studio tooling experience for ASP.NET vNext, bowels of <a href="https://github.com/aspnet/kruntime">this new runtime</a>, tiny little things about the frameworks like <a href="http://github.com/aspnet/mvc">MVC</a>, <a href="https://github.com/aspnet/identity">Identity</a>, <a href="https://github.com/aspnet/entityframework">Entity Framework</a>.</p> <p>Today, I would like to talk about something which is MVC specific and takes one of our pains away: view components :)</p> <blockquote> <p><strong>BIG ASS CAUTION!</strong> At the time of this writing, I am using <strong>KRE 1.0.0-beta1-10494 </strong>version. As things are moving really fast in this new world, it’s very likely that the things explained here will have been changed as you read this post. So, be aware of this and try to explore the things that are changed to figure out what are the corresponding new things.</p> <p>Also, inside this post I am referencing a lot of things from ASP.NET GitHub repositories. In order to be sure that the links won’t break in the future, I’m actually referring them by getting permanent links to the files on GitHub. So, these links are actually referring the files from the latest commit at the time of this writing and they have a potential to be changed, too. Read the "<a href="https://help.github.com/articles/getting-permanent-links-to-files/">Getting permanent links to files</a>" post to figure what this actually is.</p></blockquote> <p>Do you remember ugly, nasty <a href="http://www.tugberkugurlu.com/archive/donut-hole-caching-in-asp-net-mvc-by-using-child-actions-and-outputcacheattribute">child actions in ASP.NET MVC</a>? I bet you do. Child actions were pain because it’s something so weird that nobody understood at the first glance. They were sitting inside the controller as an action (hence the name) and can be invoked from the view like below:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">&lt;</span><span style="color: #a31515">div</span><span style="color: blue">&gt;</span>
    @Html.Action("widget")
<span style="color: blue">&lt;/</span><span style="color: #a31515">div</span><span style="color: blue">&gt;</span></pre></div></div>
<p>In MVC, if your HTTP request reaches that view and starts rendering, it means (most of the time with the default flow at least) that you already have gone through a controller and an action. Now, we are also invoking a child action here and this basically means that we will go through the pipeline again to pick the necessary controller and action. So, you can no longer tell that only HTTP requests will hit your controller because the child actions are not HTTP requests. They are basically method calls to your action. Not to mention <a href="https://aspnetwebstack.codeplex.com/workitem/601">the lack of asynchronous processing support inside the child actions</a>. I’m guessing that a lot of people have seen <a href="https://github.com/ravendb/ravendb/pull/545">a deadlock while blocking an asynchronous call inside a child action</a>. In a nutshell, child actions are cumbersome to me.</p>
<h3>View Components in ASP.NET vNext</h3>
<p>In ASP.NET vNext, all of the ugly behaviors of child actions are gone. To be more accurate, child actions are gone :) Instead, we have something called view components which allows you to render a view and it can be called inside a view. It has the same features of child actions without all the ugliness.</p>
<p>Let me walk you though a scenario where view components might be useful to you. In my website, you can see that I’m listing links to my external profiles at the right side.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/3700676d-84c1-421d-88a0-2d143c63cf42.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/2cadc6cf-523e-4af5-b8ed-d8288e0626a4.png" width="644" height="331"></a></p>
<p>Possibly, the links are stored in my database and I’m retrieving them by performing an I/O operation. Also, I am viewing these links in all my pages at the right side. This’s a perfect candidate for a view component where you would implement your links retrieval and display logic once (separately), and use it whenever you need it. At its core, your view component is nothing but a class which is derived from the <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponent.cs">ViewComponent</a> base class. It has a few helper methods on itself like <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponent.cs#L66-L75">View</a>, <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponent.cs#L41-L44">Content</a> and <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponent.cs#L46-L49">Json</a> methods. However, you are not restricted to this. Thanks to <a href="http://blogs.msdn.com/b/webdev/archive/2014/06/17/dependency-injection-in-asp-net-vnext.aspx">new great DI system in K Runtime</a>, we can pass dependencies into view components as well. Let’s build this up.</p>
<p>First of all, I have the following manager class which is responsible for retrieving the profile link list. For the demo purposes, it gets the list from an in-memory collection:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">interface</span> IProfileLinkManager
{
    Task&lt;IEnumerable&lt;ProfileLink&gt;&gt; GetAllAsync();
}

<span style="color: blue">public</span> <span style="color: blue">class</span> ProfileLinkManager : IProfileLinkManager
{
    <span style="color: blue">private</span> <span style="color: blue">static</span> <span style="color: blue">readonly</span> IEnumerable&lt;ProfileLink&gt; _profileLinks = <span style="color: blue">new</span> List&lt;ProfileLink&gt; 
    {
        <span style="color: blue">new</span> ProfileLink { Name = <span style="color: #a31515">"Twitter"</span>, Url = <span style="color: #a31515">"http://twitter.com/tourismgeek"</span>, FaName = <span style="color: #a31515">"twitter"</span> },
        <span style="color: blue">new</span> ProfileLink { Name = <span style="color: #a31515">"linkedIn"</span>, Url = <span style="color: #a31515">"http://www.linkedin.com/in/tugberk"</span>, FaName = <span style="color: #a31515">"linkedin"</span> },
        <span style="color: blue">new</span> ProfileLink { Name = <span style="color: #a31515">"GitHub"</span>, Url = <span style="color: #a31515">"http://github.com/tugberkugurlu"</span>, FaName = <span style="color: #a31515">"github"</span> },
        <span style="color: blue">new</span> ProfileLink { Name = <span style="color: #a31515">"Stackoverflow"</span>, Url = <span style="color: #a31515">"http://stackoverflow.com/users/463785/tugberk"</span>, FaName = <span style="color: #a31515">"stack-exchange"</span> }
    };

    <span style="color: blue">public</span> Task&lt;IEnumerable&lt;ProfileLink&gt;&gt; GetAllAsync()
    {
        <span style="color: blue">return</span> Task.FromResult&lt;IEnumerable&lt;ProfileLink&gt;&gt;(_profileLinks);
    }
}

<span style="color: blue">public</span> <span style="color: blue">class</span> ProfileLink
{
    <span style="color: blue">public</span> <span style="color: blue">string</span> Name { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">string</span> Url { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
    <span style="color: blue">public</span> <span style="color: blue">string</span> FaName { <span style="color: blue">get</span>; <span style="color: blue">set</span>; }
}</pre></div></div>
<p>Inside the <a href="https://github.com/tugberkugurlu/AspNetVNextSamples/blob/9d3eecc94d5121fa5158744160a90f010b2c0399/src/ViewComponentSample/Startup.cs">Startup.cs</a> file, we need to <a href="https://github.com/tugberkugurlu/AspNetVNextSamples/blob/9d3eecc94d5121fa5158744160a90f010b2c0399/src/ViewComponentSample/Startup.cs#L17">register the implementation of IProfileLinkManager</a>:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>app.UseServices(services =&gt; 
{
    services.AddMvc();
    services.AddScoped&lt;IProfileLinkManager, ProfileLinkManager&gt;();
});</pre></div></div>
<h3>View Components and How They Work</h3>
<p>We can now create our view component which has a dependency on ProfileLinkManager:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> ProfileLinksViewComponent : ViewComponent
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> IProfileLinkManager _profileLinkManager;

    <span style="color: blue">public</span> ProfileLinksViewComponent(IProfileLinkManager profileLinkManager)
    {
        <span style="color: blue">if</span> (profileLinkManager == <span style="color: blue">null</span>)
        {
            <span style="color: blue">throw</span> <span style="color: blue">new</span> ArgumentNullException(<span style="color: #a31515">"profileLinkManager"</span>);
        }

        _profileLinkManager = profileLinkManager;
    }

    <span style="color: blue">public</span> async Task&lt;IViewComponentResult&gt; InvokeAsync()
    {
        <span style="color: blue">var</span> profileLinks = await _profileLinkManager.GetAllAsync();            
        <span style="color: blue">return</span> View(profileLinks);
    }
}</pre></div></div>
<p>There are couple of things to highlight here. Let’s start with the name of the view component which is very important as we need to know its name so that we can refer to it when we need to process it. The name of the view component can be inferred in two ways:</p>
<ul>
<li>The name can be inferred from the name of the view component class. If the class name has a <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponentConventions.cs#L11">suffix as "ViewComponent"</a>, the view component name will be the class name with "ViewComponent" suffix. If it doesn’t have that suffix, the class name is the view component name as is.</li>
<li>You can specifically give it a name by applying the <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponentAttribute.cs">ViewComponentAttribute</a> to the view component class and setting its Name property.</li></ul>
<p>The other thing that is worth mentioning is our ability to inject dependencies into our view component. Any dependency that we have inside the request scope at the time of invocation can be injected into the view component class. Have a look at <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/DefaultViewComponentInvoker.cs#L75-L81">CreateComponent private</a> method on the DefaultViewComponentInvoker (we will touch on this later) to see how the view component class is activated by default.</p>
<p>The last thing I want to mention is the method of the view component that will be called. By default, you can have two method names here: <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponentMethodSelector.cs#L14">InvokeAsync</a> or <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/ViewComponentMethodSelector.cs#L15">Invoke</a>. As you can guess, InvokeAsync is the one you can use for any type of asynchronous processing. As these two methods are not part of the base class (ViewComponent), we can have as many parameters as we want here. Those parameters will be passed when we are actually invoking the view component inside a view and the registered <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/IViewComponentInvoker.cs">IViewComponentInvoker</a> (<a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/DefaultViewComponentInvoker.cs">DefaultViewComponentInvoker</a> by default) is responsible for calling the InvokeAsync or Invoke method (if you have both InvokeAsync and Invoke method on your view component, the InvokeAsync will be the one that’s called). From the method, you can return three types (if it’s InvokeAsync, the following types are for the generic parameter of the Task&lt;T&gt; class): <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/ViewComponents/IViewComponentResult.cs">IViewComponentResult</a>, string and <a href="https://github.com/aspnet/Mvc/blob/8802c831a0f2131d904d0b416466d285d87f7139/src/Microsoft.AspNet.Mvc.Core/Rendering/HtmlString.cs">HtmlString</a>. As you can see, I am using the View method of the ViewComponent base class above which returns IViewComponentResult.</p>
<h3>Using the View Components</h3>
<p>We covered most of things that we need to know about view components except for how we can actually use them. Using a view component is actually very similar to how we used child actions. We need to work with the Component property inside the view (which is an implementation of <a href="https://github.com/aspnet/Mvc/blob/12477c9f52be0b6f7cdf9eb71953b6dbbbd3b5da/src/Microsoft.AspNet.Mvc.Core/ViewComponents/IViewComponentHelper.cs">IViewComponentHelper</a>). There are a few methods for IViewComponentHelper that we can call. You can see below that I am calling the InvokeAsync by only passing the component name. </p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">&lt;</span><span style="color: #a31515">body</span><span style="color: blue">&gt;</span>
    <span style="color: blue">&lt;</span><span style="color: #a31515">div</span> <span style="color: red">class</span><span style="color: blue">=</span><span style="color: blue">"col-md-8"</span><span style="color: blue">&gt;</span>
        <span style="color: blue">&lt;</span><span style="color: #a31515">hr</span> <span style="color: blue">/&gt;</span>
        <span style="color: blue">&lt;</span><span style="color: #a31515">div</span><span style="color: blue">&gt;</span>
            Main section...
        <span style="color: blue">&lt;/</span><span style="color: #a31515">div</span><span style="color: blue">&gt;</span>
    <span style="color: blue">&lt;/</span><span style="color: #a31515">div</span><span style="color: blue">&gt;</span>
    <span style="color: blue">&lt;</span><span style="color: #a31515">div</span> <span style="color: red">class</span><span style="color: blue">=</span><span style="color: blue">"col-md-4"</span><span style="color: blue">&gt;</span>
        <span style="color: blue">&lt;</span><span style="color: #a31515">hr</span> <span style="color: blue">/&gt;</span>
        @await Component.InvokeAsync("ProfileLinks")
    <span style="color: blue">&lt;/</span><span style="color: #a31515">div</span><span style="color: blue">&gt;</span>
<span style="color: blue">&lt;/</span><span style="color: #a31515">body</span><span style="color: blue">&gt;</span></pre></div></div>
<blockquote>
<p>And yes, you can await inside the razor view itself :) </p></blockquote>
<p>InvokeAsync method here also accepts "params object[]" as its second parameter and you can pass the view component parameters there (if there are any). Let’s run the application and see what we are getting:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6f7cd2cc-f642-486f-b422-aa3da4639d3b.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a044f248-2a0a-4514-8978-11ffcdf674b3.png" width="644" height="383"></a></p>
<blockquote>
<p>To be able to see this rich error message, you need to pull down the <a href="https://github.com/aspnet/Diagnostics">Microsoft.AspNet.Diagnostics</a> package and <a href="https://github.com/tugberkugurlu/AspNetVNextSamples/blob/9d3eecc94d5121fa5158744160a90f010b2c0399/src/ViewComponentSample/Startup.cs#L11">activate the diagnostic page for your application as I did</a>. Unlike the old ASP.NET world, we don’t get any special feature for free in this world. Remember, pay-for-play model :)</p></blockquote>
<p>We got an error which is expected:</p>
<blockquote>
<h3>An unhandled exception occurred while processing the request.</h3>
<p>InvalidOperationException: The view 'Components/ProfileLinks/Default' was not found. The following locations were searched: /Views/Home/Components/ProfileLinks/Default.cshtml /Views/Shared/Components/ProfileLinks/Default.cshtml.</p>
<p>Microsoft.AspNet.Mvc.ViewViewComponentResult.FindView(ActionContext context, String viewName)</p></blockquote>
<p>This’s actually very good that we are getting this error message because it explains what we need to do next. As we didn’t pass any view name, the Default.cshtml the view file is the one our component is looking for. The location of the view needs to be either under the controller that we now rendering (which is Views/Home/ProfileLinks) or the Shared folder (which is Views/Shared/ProfileLinks). Let’s put the view under the Views/Shared/Components/ProfileLinks directory:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>@using ViewComponentSample
@model IEnumerable<span style="color: blue">&lt;</span><span style="color: #a31515">ProfileLink</span><span style="color: blue">&gt;</span>

<span style="color: blue">&lt;</span><span style="color: #a31515">ul</span> <span style="color: red">class</span><span style="color: blue">=</span><span style="color: blue">"profilelink-list"</span><span style="color: blue">&gt;</span>
    @foreach(ProfileLink profileLink in Model)
    {
        <span style="color: blue">&lt;</span><span style="color: #a31515">li</span><span style="color: blue">&gt;</span><span style="color: blue">&lt;</span><span style="color: #a31515">a</span> <span style="color: red">class</span><span style="color: blue">=</span><span style="color: blue">"btn btn-info"</span> <span style="color: red">title</span><span style="color: blue">=</span><span style="color: blue">"@profileLink.Name"</span> <span style="color: red">href</span><span style="color: blue">=</span><span style="color: blue">"@profileLink.Url"</span><span style="color: blue">&gt;</span><span style="color: blue">&lt;</span><span style="color: #a31515">i</span> <span style="color: red">class</span><span style="color: blue">=</span><span style="color: blue">"fa fa-@(profileLink.FaName) fa-lg"</span><span style="color: blue">&gt;</span><span style="color: blue">&lt;/</span><span style="color: #a31515">i</span><span style="color: blue">&gt;</span><span style="color: blue">&lt;/</span><span style="color: #a31515">a</span><span style="color: blue">&gt;</span><span style="color: blue">&lt;/</span><span style="color: #a31515">li</span><span style="color: blue">&gt;</span>
    }
<span style="color: blue">&lt;/</span><span style="color: #a31515">ul</span><span style="color: blue">&gt;</span></pre></div></div>
<p>When we now run the application, we should see that the component is rendered successfully:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/bce6b0ef-f164-4397-a16a-078d82bccdb0.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/68de3cba-aff2-49e1-b14d-53ffe3904be9.png" width="644" height="196"></a></p>
<p>Nice and shiny! You can find the sample I have gone through here inside the <a href="https://github.com/tugberkugurlu/AspNetVNextSamples">ASP.NET vNext samples</a> repository: <a href="https://github.com/tugberkugurlu/AspNetVNextSamples/tree/9d3eecc94d5121fa5158744160a90f010b2c0399/src/ViewComponentSample">ViewComponentSample</a>.</p>  