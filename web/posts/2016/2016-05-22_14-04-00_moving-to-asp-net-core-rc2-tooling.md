---
id: 03645b78-b584-45e2-b15d-8db2cb8daa2d
title: 'Moving to ASP.NET Core RC2: Tooling'
abstract: .NET Core Runtime RC2 has been released a few days ago along with .NET Core
  SDK Preview 1. At the same time of .NET Core release, ASP.NET Core RC2 has also
  been released. While I am migrating my projects to RC2, I wanted to write about
  how I am getting each stage done. In this post, I will show you the tooling aspect
  of the changes.
created_at: 2016-05-22 14:04:00 +0000 UTC
tags:
- .net
- .NET Core
- ASP.Net
- ASP.NET Core
slugs:
- moving-to-asp-net-core-rc2-tooling
---

<p><a href="https://blogs.msdn.microsoft.com/dotnet/2016/05/16/announcing-net-core-rc2/">.NET Core Runtime RC2 has been released a few days ago along with .NET Core SDK Preview 1</a>. At the same time of .NET Core release, <a href="https://blogs.msdn.microsoft.com/webdev/2016/05/16/announcing-asp-net-core-rc2/">ASP.NET Core RC2 has also been released</a>. Today, I started doing the transition from RC1 to RC2 and I wanted to write about how I am getting each stage done. Hopefully, it will be somewhat useful to you as well. In this post, I want to talk about the tooling aspect of the transition.</p> <h3>Get the dotnet CLI Ready</h3> <p>One of the biggest shift from RC1 and RC2 is the tooling. Before, we had DNVM, DNX and DNU as command line tools. All of them are now gone (RIP). Instead, we have one command line tool: <a href="https://github.com/dotnet/cli">dotnet CLI</a>. First, I installed dotnet CLI on my Ubuntu 14.04 VM by running the following script as explained <a href="https://www.microsoft.com/net/core#ubuntu">here</a>:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>sudo sh <span style="color: gray">-</span>c <span style="color: #a31515">'echo "deb [arch=amd64] https://apt-mo.trafficmanager.net/repos/dotnet/ trusty main" &gt; /etc/apt/sources.list.d/dotnetdev.list'</span>
sudo apt<span style="color: gray">-</span>key adv <span style="color: gray">--</span>keyserver apt<span style="color: gray">-</span>mo.trafficmanager.net <span style="color: gray">--</span>recv<span style="color: gray">-</span>keys 417A0893
sudo apt<span style="color: gray">-ge</span>t update
sudo apt<span style="color: gray">-ge</span>t install dotnet<span style="color: gray">-</span>dev<span style="color: gray">-</span>1.0.0<span style="color: gray">-</span>preview1<span style="color: gray">-</span>002702</pre></div></div>
<p>This installed dotnet-dev-1.0.0-preview1-002702 package on my machine and I am off to go:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/76d4308b-d312-4b07-a06e-c860b5f772f6.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/850e5f07-423c-4d7a-a175-97996453c4dc.png" width="644" height="286"></a></p>
<p>You can also use apt-cache to see all available versions:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/460c0a2f-72ca-4f26-a1b5-0cb734bf0f52.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/9f2d36c5-a48b-448e-aa39-85f82f8925a8.png" width="644" height="445"></a></p>
<p>Just to make things clear, I deleted ~/.dnx folder entirely to get rid of all RC1 garbage.</p>
<h3>Get Visual Studio Code Ready</h3>
<p>At this stage, I had the <a href="https://marketplace.visualstudio.com/items?itemName=ms-vscode.csharp">C# extension</a> installed on my VS Code instance on my Ubuntu VM which was only working for DNX based projects. So, I opened up VS Code and <a href="https://code.visualstudio.com/Docs/editor/extension-gallery#_update-an-extension">updated my C# extension</a> to latest version (which is <a href="https://github.com/OmniSharp/omnisharp-vscode/releases/tag/v1.0.11">1.0.11</a>). After the upgrade, I opened up a project which was on RC1 and VS Code started downloading .NET Core Debugger.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/faddecc9-51d9-4cf7-a401-e7e87e554263.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/2e68b12b-71ff-46f1-85e0-1d30714054f7.png" width="644" height="119"></a></p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/93a28d08-e9a7-45be-a997-8f5c3d275fc8.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ff842886-b113-4c3b-be28-6121288d272e.png" width="644" height="83"></a></p>

<p>That was a good experience, I didn't dig into how to do that but I am not sure at this point why it didn't come inside the extension itself. There is probably a reason for that but not my priority to dig into right now :)</p>
<h3>Try out the Setup</h3>
<p>Now, I am ready to blast off with .NET Core. I used dotnet CLI to create a sample project by typing <em>"dotnet new --type=console" </em>and opened up project folder with VS Code. As soon as VS Code is launched, it asked me to set the stage.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/66e7708b-3775-4f13-9912-97ea187c6bc4.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/3b715c7b-4212-4809-bae6-57c5ca7f5cb3.png" width="644" height="114"></a></p>
<p>Which got me a few new files under .vscode folder.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/efa3d81b-2e18-4e8d-9945-30c7ee53d420.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/f63c1054-a690-48e0-9207-f9d394a8be62.png" width="644" height="370"></a></p>
<p>I jumped into the debug pane, selected the correct option and hit the play button after putting a breakpoint inside the Program.cs file.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/82dfb8b9-1b5f-4149-ab8e-1438793f9b38.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/2a7f2ce3-b851-48ff-acc9-3595147e5c8d.png" width="644" height="370"></a></p>
<p>Boom! I am in business.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/0e033229-d5e7-4a37-9a23-ed1557f2b54b.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ea79fc82-417e-479e-81a1-9b413707dd34.png" width="644" height="370"></a></p>
<p>Now, I am moving to code changes which will involve more hair pulling I suppose.</p>
<h3>Resources</h3>
<ul>
<li><a href="https://docs.asp.net/en/latest/">ASP.NET Core Docs</a></li>
<li><a href="http://dotnet.github.io/docs/">.NET Core Documentation</a></li>
<li><a href="http://dotnet.github.io/docs/core-concepts/core-sdk/index.html">.NET Core SDK Reference</a></li>
<ul></ul>
<li><a href="https://docs.asp.net/en/latest/conceptual-overview/aspnet.html">Introduction to ASP.NET Core</a></li>
<li><a href="https://docs.asp.net/en/latest/tutorials/your-first-mac-aspnet.html">Your First ASP.NET Core Application on a Mac Using Visual Studio Code</a></li>
<li><a href="https://github.com/dotnet/cli/blob/rel/1.0.0/Documentation/intro-to-cli.md">Intro to .NET Core CLI</a></li></ul>  