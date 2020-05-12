---
id: dd331015-dcde-49a2-94e0-229be30a8b5f
title: 'Dependency Injection: Inject Your Dependencies, Period!'
abstract: Reasons on why I prefer dependency injection over static accessors.
created_at: 2014-11-18 12:39:00 +0000 UTC
tags:
- .net
- ASP.Net
- C#
slugs:
- dependency-injection-inject-your-dependencies-period
---

<p><a href="https://github.com/aspnet/Logging/issues/57">There is a discussion going on</a> inside <a href="https://github.com/aspnet/Logging">ASP.NET/Logging repository</a> whether we should have static Logger available everywhere or not. I am quite against it and I stated why with a few comments there but hopefully with this post, I will try to address why I think the <a href="http://martinfowler.com/articles/injection.html">dependency injection</a> path is better for most of the cases.</p> <p>Let’s take an example and explain it further based on that.</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> FooManager : IFooManager
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> IClock _clock;
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> ILogger _logger;
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> IEnvironmentInfo _environmentInfo;

    <span style="color: blue">public</span> FooManager(IClock clock, ILogger logger, IEnvironmentInfo environmentInfo)
    {
        _clock = clock;
        _logger = logger;
        _environmentInfo = environmentInfo;
    }

    <span style="color: green">// ...</span>
}</pre></div></div>
<p>This is a simple class&nbsp; whose job is not important in our context. The way I see this class inside Visual Studio is as below:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/57c2c18a-dab5-40af-85c4-48edc6851f30.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/152d8477-cf43-4bf2-9ffa-aef7f68ec171.png" width="644" height="290"></a></p>
<p>In any other text editor:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/2fce05e1-572a-46f3-b1b1-48b4d4881e84.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/78f3dd3f-a0b9-405a-8115-d98376e645cc.png" width="644" height="301"></a></p>
<p>Inside Visual Studio when I am writing tests:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/81e0134d-e318-465c-8629-cbe452fdfcbe.png"><img title="Screenshot 2014-11-18 13.07.38" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Screenshot 2014-11-18 13.07.38" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/cfde73c3-4568-4303-9eb8-aafc9ce678bc.png" width="644" height="218"></a></p>
<p>There are two things fundamental things why I would like this approach.</p>
<h3>The Requirements of the Component is Exposed Clearly</h3>
<p>When I look at this class inside any editor, I can say that this class cannot function properly without IClock, ILogger and IEnvironmentInfo implementation instances as there is no other constructor (preferably, I would also do null checks for constructor parameters but I skipped that to keep the example clean). Instead of the above implementation, imagine that I have the following one:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> FooManager : IFooManager
{
    <span style="color: blue">public</span> <span style="color: blue">void</span> Run()
    {
        <span style="color: blue">if</span>(DateTime.UtcNow &gt; EnvironmentInfo.Instance.LicanceExpiresInUtc)
        {
            Logger.Instance.Log(<span style="color: #a31515">"Cannot run as licance has expired."</span>);
        }

        <span style="color: green">// Do other stuff here...</span>
    }
}</pre></div></div>
<p>With this approach, we are relying on static instances of the components (and they are possibly thread-safe and singleton). What is wrong with this approach? I’m not sure what’s the general idea for this but here are my reasons. </p>
<p>First thing for me to do when I open up a C# code file inside visual studio is to press CTRL+M, O to get an idea about the component. It has been kind of an habit for me. Here how it looks like when I do that for this class:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/95dffde0-177d-4fa9-a98e-03b27161be72.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/bc955379-3caf-4727-b6ee-8af43513175f.png" width="644" height="216"></a></p>
<p>I have no idea what this class needs to function properly. I have also no idea what type of environmental context it relies on. Please note that this issue is not that big of a problem for a class which is this simple but I imagine that your component will have a few other methods, possible other interface implementations, private fields, etc.</p>
<h3>I Don’t Need to Look at the Implementation to See What It is Using</h3>
<p>When I try to construct the class inside a test method in Visual Studio, I can super easily see what I need to supply to make this class function the way I want. If you are using another editor which doesn’t have a tooling support to show you constructor parameters, you are still safe as you would get compile errors if you don’t supply any required parameters. With the above static instance approach, however, you are on your own with this issue in my opinion. It doesn’t have any constructor parameters for you to easily see what it needs. It relies on static instances which are hard and dirty to mock. For me, it’s horrible to see a C# code written like that. </p>
<p>When you try to write a test against this class now, here is what it looks like if you are inside Visual Studio:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c8c2ba3c-30b7-4c8b-863e-40fe683e4607.png"><img title="Screenshot 2014-11-18 14.17.13" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Screenshot 2014-11-18 14.17.13" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/fa0f1c2d-d766-4dc9-acf5-1afad85abaad.png" width="644" height="205"></a></p>
<p>You have no idea what it needs. F12 and try to figure out what it needs by looking at the implementation and <strong>good luck with mocking the static read-only members and DateTime.Now</strong> :)</p>
<blockquote>
<p>I’m intentionally skipping why you shouldn’t use DateTime.Now directly inside your library (even, the whole DateTime API). The reasons are variable depending on the context. However, here are a few further readings for you on this subject:</p>
<ul>
<li><a href="http://codeofmatt.com/2013/04/25/the-case-against-datetime-now/">The case against DateTime.Now</a></li>
<li><a href="http://nodatime.org/">Noda Time: A better date and time API for .NET</a></li>
<li><a href="http://nodatime.org/unstable/userguide/testing.html">Unit testing with Noda Time</a></li>
<li><a href="https://www.nuget.org/packages/NodaTime.Testing">NodaTime.Testing</a></li></ul></blockquote>
<h3>It is Still Bad Even If It is not Static</h3>
<p>Yes, it’s still bad. Let me give you a hint: <a href="http://blog.ploeh.dk/2010/02/03/ServiceLocatorisanAnti-Pattern/">Service Locator Pattern</a>.</p>  