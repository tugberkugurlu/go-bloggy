---
id: fd9c1aa5-5d02-4a0f-bc71-14fa41be2fd6
title: Unit Testing With xUnit.net for ASP.NET MVC Web Applications - Part 1
abstract: In this blog post, we will see how we set up our environment for xUnit.net
  Unit Testing Framework. This is the first blog post of the blog post series on Unit
  Testing With xUnit.net for ASP.NET MVC.
created_at: 2011-10-16 14:40:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- Unit Testing
- Visual Studio
- xUnit
slugs:
- unit-testing-with-xunit-net-for-asp-net-mvc-web-applications-part-1
---

<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/you-need-some-tests-yo.jpg"><img style="background-image: none; margin: 0px 0px 15px 15px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border: 0px;" title="you-need-some-tests-yo..." border="0" alt="you-need-some-tests-yo..." align="right" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/you-need-some-tests-yo..._thumb.jpg" width="244" height="184" /></a></p>
<p>I have been developing web applications with <a title="http://asp.net/mvc" href="http://asp.net/mvc" target="_blank">ASP.NET MVC</a> for long time and to be honest, none of them has good Unit Testing backing them up. I don&rsquo;t know you have realized but I am not even mentioning anything about <a title="http://en.wikipedia.org/wiki/Test-driven_development" href="http://en.wikipedia.org/wiki/Test-driven_development" target="_blank">TDD (<strong>T</strong>est <strong>D</strong>riven <strong>D</strong>evelopment)</a>.</p>
<p>This isn&rsquo;t a big problem for me. It is mainly because I work in a team, which has only one member inside it [yeah, it&rsquo;s me :)], for my company. But all the cool guys are always talking about how awesome it is to work with Unit Testing.</p>
<p>This made me obsessed about Unit Testing and I stated to spend real quality time on investigating that so that I will be able to sleep without having dreams about a guy who wears an ugly sunglasses and telling me &lsquo;<strong>You need some tests yo!</strong>&rsquo;. (No, I am just kidding, I am not insane. At least, not yet.)</p>
<p>Visual Studio has a built-in test environment but I figured that is not enough. So, from now on, I will be using <a title="http://xunit.codeplex.com" href="http://xunit.codeplex.com" target="_blank">xUnit.net</a> for my unit test applications.</p>
<p><strong>What is xUnit.net?</strong></p>
<p><em>xUnit.net is a unit testing tool for the .NET Framework. Written by the original inventor of NUnit, xUnit.net is the latest technology for unit testing C#, F#, VB.NET and other .NET languages. Works with ReSharper, CodeRush, and TestDriven.NET.</em></p>
<p><strong>Set Up Your Environment</strong></p>
<p>First thing is first. Go to <a href="http://xunit.codeplex.com/releases/view/62840">http://xunit.codeplex.com/releases/view/62840</a> to download the latest version of xUnit.net. At the time that I write this blog post, the latest version of xUnit.net is <strong>1.8</strong>. Unzip the file inside a path and run the <strong>xUnit.installer.exe</strong> file.</p>
<p>You will have a window opened in front of you. For the sake of simplicity, I only check Microsoft <strong>ASP.NET MVC 3 Unit Testing Templates </strong>and click <strong>Apply</strong> button. After that, here is the screen shots of the result :</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image.png"><img style="background-image: none; margin: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb.png" width="244" height="179" /></a><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_3.png"><img style="background-image: none; margin: 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_3.png" width="244" height="179" /></a></p>
<p><strong>File &gt; New &gt; Project</strong></p>
<p>After the installation, go to VS 2010 and follow File &gt; New &gt; Project &gt; ASP.NET MVC 3 Web Application. You will see that xUnit is sitting right there.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_4.png" width="538" height="484" /></a></p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_5.png"><img style="background-image: none; margin: 0px 15px 15px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" align="left" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_5.png" width="161" height="244" /></a>Go ahead and create the project. When you create the project you will see the usual HomeController.cs inside your <strong>Controllers</strong> folder. With xUnit.net project, you will see that we get HomeControllerFacts.cs file out of the box as sample. This file gives us a head start on this new thing for us.</p>
<p>On this blog post, I would like to focus on how things goes for us instead of how we write our unit test code.</p>
<p>As we have seen, we get a sample project and sample test project with some tests written inside it the the sample application. When we build the solution, we should be able to process successfully.</p>
<p>So, the next question is how we run our tests. The xUnit.net framework supports 4 distinct test runners:</p>
<ul>
<li>GUI Test Runner</li>
<li>Console Test Runner</li>
<li>TestDriven.NET Test Runner</li>
<li>Resharper 4.0 Test Runner</li>
</ul>
<p>I will walk you through how we run our tests with <strong>GUI Test Runner </strong>for this blog post.</p>
<p><strong>Running xUnit Tests With GUI Test Runner</strong></p>
<p>Navigate to the path which you have installed you xUnit project files and you will notice 4 files there which are green as follows :</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_6.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_6.png" width="644" height="484" /></a></p>
<p>Firstly, they are not same. I am sure you already understood what is for what but here is the detailed info :</p>
<p><strong>xuit.gui.clr4, <strong>xuit.gui.clr4.x86</strong> : </strong>This is for CLR 4 projects. .x86 is for x86 compatibility which you already know.</p>
<p><strong>xuit.gui, xuit.gui.x86 : </strong>This is for CLR 2 projects. .x86 is for x86 compatibility which you already know.</p>
<p>I have open up the <strong>xuit.gui.exe </strong>which is suitable for me. Then, from the GUI, select the menu option Assembly &gt; Open (<strong>Ctrl + O</strong> for keyboard shortcut) and select your test application output <strong>dll file </strong>(Be careful here. If you don&rsquo;t build your solution before this process, you won&rsquo;t be able to find the output <strong>dll file</strong>).</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_7.png" width="594" height="484" /></a></p>
<p>When you make the selection, the runner will be able to pick up all the tests as you can see below :</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_8.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_8.png" width="594" height="484" /></a></p>
<p>Let&rsquo;s touch the <em>Run All</em> button on the left side at the bottom :</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_9.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_9.png" width="594" height="484" /></a></p>
<p>Looks cool. I break one of the tests on purpose so that we can see it failing :</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_10.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/1934ffeb460c_FB47/image_thumb_10.png" width="594" height="484" /></a></p>
<p>So, this is the simplest as it can get. This blog post is supposed to be the first blog post of the blog post series on <strong>Unit Testing With xUnit.net for ASP.NET MVC </strong>and I will be blogging on further features of xUnit for ASP.NET MVC. Stay tuned!</p>