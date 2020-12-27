---
id: b0924573-7c78-4e52-bcb9-d62f24e76760
title: 'Exciting Things About ASP.NET 5 Series: Build Only Dependencies'
abstract: In this very exciting post, I would like to talk about build only dependencies
  whose code can be compiled into target project and the dependency won’t be shown
  as a dependency.
created_at: 2015-04-28 07:48:00 +0000 UTC
tags:
- .NET
- ASP.NET 5
slugs:
- exciting-things-about-asp-net-5-series-build-only-dependencies
---

<p>Web development experience with .NET has never seen a drastic change like this since its birth day. Yes, I’m talking about <a href="https://www.tugberkugurlu.com/archive/getting-started-with-asp-net-vnext-by-setting-up-the-environment-from-scratch">ASP.NET 5</a> :) I have been putting my toes into this water for a while now and a few days ago, I started a new blog post series about ASP.NET 5 (with hopes that I will continue this time :)). To be more specific, I’m planning on writing about the things I am actually excited about this new cloud optimized (TM) runtime. Those things could be anything which will come from <a href="http://github.com/aspnet">ASP.NET GitHub account</a>: things I like about the development process, Visual Studio tooling experience for ASP.NET 5, bowels of <a href="https://github.com/aspnet/dnx">.NET Execution Runtime</a>, tiny little things about the frameworks like <a href="http://github.com/aspnet/mvc">MVC</a>, <a href="https://github.com/aspnet/identity">Identity</a>, <a href="https://github.com/aspnet/entityframework">Entity Framework</a>.</p> <p>In this very exciting post, I would like to talk about build only dependencies whose code can be compiled into target project.</p> <blockquote> <p><strong>BIG ASS CAUTION!</strong> At the time of this writing, I am using <strong>DNX 1.0.0-beta5-11611 </strong>version. As things are moving really fast in this new world, it’s very likely that the things explained here will have been changed as you read this post. So, be aware of this and try to explore the things that are changed to figure out what are the corresponding new things. <p>Also, inside this post I am referencing a lot of things from ASP.NET GitHub repositories. In order to be sure that the links won’t break in the future, I’m actually referring them by getting permanent links to the files on GitHub. So, these links are actually referring the files from the latest commit at the time of this writing and they have a potential to be changed, too. Read the "<a href="https://help.github.com/articles/getting-permanent-links-to-files/">Getting permanent links to files</a>" post to figure what this actually is.</p></blockquote> <h3>The Problem</h3> <p>From the start of NuGet, it has been a real pain to have source file dependencies. There are some examples of this like <a href="https://www.nuget.org/packages/TaskHelpers.Sources/">TaskHelpers.Sources</a>. When you install this package, it will end up inside your codebase. </p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c7c7bee5-5d91-4e9e-b8d0-77de93ac1c4d.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a73ec837-e85d-4446-80fb-2d7037661e1f.png" width="644" height="414"></a></p> <p>The nice thing about this type of source dependencies is that you don’t need to fight with DLL hell. You can have one version of this package and your consumer can have another version of it. As the source files you pull down from NuGet has no public members, there will be no problems whatsoever as the code is compiled into their assembly separately. However, there are several problems with the way we are getting them in:</p> <ul> <li>I am committing this code into source control system which is weird.</li> <li>How about updates? What happens if I make a change to that file?</li></ul> <p>So, it wasn’t that good of an approach we had there but ASP.NET 5 has a top notch solution this problem: build only dependencies. </p> <h3>Consuming Build Only Dependencies</h3> <p>These are the kind of dependencies that you can pull in and it will just be compiled into your stuff. As you can also guess, it won’t be shown as a dependency. Let’s see an example!</p> <p>One of the packages that support this concept is <a href="https://github.com/aspnet/dnx/tree/671232faf23e0608037cde145af7281aa0d29dd1/src/Microsoft.Framework.CommandLineUtils">Microsoft.Framework.CommandLineUtils</a> package. You can pull this down as a build-only dependency by declaring it inside your <a href="https://github.com/tugberkugurlu/AspNet5BuildOnlyDependenciesSample/blob/99b5f89a3e1aa8258e715a83ecdcf0b5f93bf860/src/AspNet5CommandLineSample/project.json">project.json</a> file as below:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>{
    <span style="color: #a31515">"version"</span>: <span style="color: #a31515">"1.0.0-*"</span>,

    <span style="color: #a31515">"dependencies"</span>: {
        <span style="color: #a31515">"Microsoft.Framework.CommandLineUtils"</span>: { 
            <span style="color: #a31515">"version"</span>: <span style="color: #a31515">"1.0.0-beta5-11611"</span>, <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"build"</span> 
        }
    },

    <span style="color: green">// ...</span>
}</pre></div></div>
<p>Notice the type field there.&nbsp; That indicates the type of the dependency. Let’s stop here and without doing anything else further, run dnu pack to get a NuGet package out. When we look at the manifest of the generated NuGet package, we won’t see any sign of the build dependency there:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/35604c66-a0a9-4afc-8efc-d8bf083c148f.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/41a70d9c-f3d7-4706-a7c4-fac5033e1058.png" width="644" height="468"></a></p>
<p>Makes sense. Let’s peak inside the assembly now.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/03d15209-e183-46dd-824d-33bbbade77f5.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/fdc5a7ad-5239-438a-a682-3309dc97dca0.png" width="644" height="330"></a></p>
<p>That’s what I expected to see. All the stuff distributed with that packages is compiled into my target assembly. As you can guess, I can use these stuff inside my project without any problems:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">using</span> Microsoft.Framework.Runtime.Common.CommandLine;

<span style="color: blue">namespace</span> AspNet5CommandLineSample
{
    <span style="color: blue">public</span> <span style="color: blue">class</span> Program
    {
        <span style="color: blue">public</span> <span style="color: blue">void</span> Main(<span style="color: blue">string</span>[] args)
        {
            <span style="color: blue">var</span> app = <span style="color: blue">new</span> CommandLineApplication();
        }
    }
}</pre></div></div>
<blockquote>
<p>You may ask that ASP.NET 5 applications can work without assemblies on disk. That’s true and at that point, this will end up being compiled into the target assembly in-memory. </p></blockquote>
<p>If you look at <a href="https://github.com/tugberkugurlu/AspNet5BuildOnlyDependenciesSample/tree/99b5f89a3e1aa8258e715a83ecdcf0b5f93bf860/src/AspNet5CommandLineSample">what I committed to my source control system</a>, it’s barely nothing which solves one of the biggest pains of source packages.</p>
<h3>Generating Build Only Dependencies</h3>
<p>Generating libraries which can be consumed as a build only dependency is also fairly simple but there are some little things which doesn’t make sense. Assuming I have a library called <a href="https://github.com/tugberkugurlu/AspNet5BuildOnlyDependenciesSample/tree/99b5f89a3e1aa8258e715a83ecdcf0b5f93bf860/src/AspNet5Utils">AspNet5Utils</a> and it has the following internal type:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">namespace</span> AspNet5Utils
{
    <span style="color: blue">internal</span> <span style="color: blue">static</span> <span style="color: blue">class</span> StringExtensions
    {
        <span style="color: blue">internal</span> <span style="color: blue">static</span> <span style="color: blue">string</span> Suffix(<span style="color: blue">this</span> <span style="color: blue">string</span> value, <span style="color: blue">string</span> suffix)
        {
            <span style="color: blue">return</span> $<span style="color: #a31515">"{value}-{suffix}"</span>;
        }
    }
}</pre></div></div>
<p>If you want this type to end up as a build dependency, you need to declare this as shared inside the project.json file.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>{
    <span style="color: #a31515">"version"</span>: <span style="color: #a31515">"1.0.0-*"</span>,

    <span style="color: #a31515">"shared"</span>: <span style="color: #a31515">"**/*.cs"</span>,

    <span style="color: #a31515">"dependencies"</span>: {
    },

    <span style="color: green">// ...</span>
}</pre></div></div>
<p>Doing this will give a hint to dnu pack command to pack these types into the shared folder inside the NuGet package.</p>

<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/31f2d208-ed2b-4bb3-b94d-7b7247d6bc96.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b3f4128e-ec52-4045-8634-33280831435d.png" width="644" height="407"></a></p>
<p>Notice that there is also an assembly generated there. Maybe there is a reason behind why this is there but as I don’t have any type which ends up inside an assembly, I would expect this to not have one at all. Indeed, if you decompile the assembly, you will see that nothing is there:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/26ae595f-0c34-44ba-89e6-0120a5e3879e.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b8c1bee4-c15e-4873-b3ee-0807b8bcdfaf.png" width="644" height="349"></a></p>
<p>In order to consume this package, you don’t actually need to distribute this through NuGet if you only want to consume this inside the same solution. As the dependency consumption is unified in ASP.NET 5, <a href="https://github.com/tugberkugurlu/AspNet5BuildOnlyDependenciesSample/blob/99b5f89a3e1aa8258e715a83ecdcf0b5f93bf860/src/AspNet5Utils.Consumer/project.json">this can easy be a project dependency</a> as you would expect:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>{
    <span style="color: #a31515">"version"</span>: <span style="color: #a31515">"1.0.0-*"</span>,

    <span style="color: #a31515">"dependencies"</span>: {
        <span style="color: #a31515">"AspNet5Utils"</span>: { <span style="color: #a31515">"version"</span>: <span style="color: #a31515">""</span>, <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"build"</span> }
    },

    <span style="color: green">// ..</span>
}</pre></div></div>
<p>In my opinion, this is one of the many powerful and yet simple concepts that ASP.NET 5 has brought to us. Enjoy!</p>  