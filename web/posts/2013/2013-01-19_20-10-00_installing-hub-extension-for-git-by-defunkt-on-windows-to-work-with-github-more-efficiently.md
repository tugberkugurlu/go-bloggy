---
id: 44c8806f-5066-49aa-8c3e-bdea9de6d959
title: Installing hub Extension for Git (by defunkt) on Windows to Work With GitHub
  More Efficiently
abstract: This post walks you through on how you can install hub extension for Git
  (by defunkt) on Windows to work with GitHub more efficiently.
created_at: 2013-01-19 20:10:00 +0000 UTC
tags:
- Git
- GitHub
- Tips
slugs:
- installing-hub-extension-for-git-by-defunkt-on-windows-to-work-with-github-more-efficiently
---

<p>We are all in love with <a href="http://git-scm.com/">Git</a> but without <a href="https://github.com">GitHub</a>, we love Git less. On GitHub, we can maintain our projects very efficiently. Pull Request&rdquo; and "Issues" features of GitHub are the key factors for that IMO. You can even send yourself a pull request from one branch to another and discuss that particular change with your team. As your discussion flows, your code can flow accordingly, too. This is just one of the many coolest features of GitHub.</p>
<p>There is a cool Git extension for GitHub which is maintained by one of the founders of GitHub: <a href="https://github.com/defunkt">Chris Wanstrath</a>. This cool extension named <a href="https://github.com/defunkt/hub">hub</a> lets us work with GitHub more efficiently from the command line and perform GitHub specific operations easily like sending pull requests, forking repositories, etc. It&rsquo;s fairly easy to install it on other platforms as far as I can see but it&rsquo;s not that straight forward for Windows.</p>
<p>You should first go and install <a href="http://msysgit.github.com/">msysgit</a> on Windows and I am assuming most of us using this on Windows for Git. Secondly, we should install Ruby on windows. You can install Ruby on windows through <a href="http://rubyinstaller.org/">RubyInstaller</a> easily.</p>
<p>After installing ruby on our machine successfully, we should add the bin path of Ruby to our system PATH variable. In order to do this, press Windows Key + PAUSE BREAK to open up the Windows System window and click "Advanced system settings" link on the left hand side of the window.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/SNAGHTML25551807.png"><img height="384" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/SNAGHTML25551807_thumb.png" alt="SNAGHTML25551807" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="SNAGHTML25551807" /></a></p>
<p>A new window should appear. From there, click "Environment Variables..." button to open up the Environment Variables window.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/image.png"><img height="244" width="219" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>From there, you should see "System variables" section. Find the Path variable and concatenate the proper ruby bin path to that semicolon-separated list.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/image_3.png"><img height="114" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/SNAGHTML255c49a4.png"><img height="108" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/SNAGHTML255c49a4_thumb.png" alt="SNAGHTML255c49a4" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; margin: 0px 0px 0px 10px; display: inline; padding-right: 0px; border: 0px;" title="SNAGHTML255c49a4" /></a></p>
<p>Last step is actually installing the hub. You should grab <a href="http://defunkt.io/hub/standalone">the standalone file</a> and then rename it to "hub". Then, put it under the Git\bin folder. The full path of my Git\bin folder on my 64x machine is "C:\Program Files (x86)\Git\bin".</p>
<p>Now you should be able to run hub command from Git Bash:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/image_4.png"><img height="349" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/image_thumb_4.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Special GitHub commands you get through hub extension is nicely documented on <a href="https://github.com/defunkt/hub#readme">the "Readme" file of the project</a>. I think the coolest feature of hub is the pull-request feature. On GitHub, You can send pull requests to another repository through GitHub web site or GitHub API and hub extension uses GitHub API under the covers to send pull requests. You can even attach your pull request to an existing issue. For example, the following command sends a pull request to master branch of the tugberkugurlu&rsquo;s repository from the branch that I am currently on and attaches this to an existing issue #1.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>hub pull-request -i 1 -b tugberkugurlu:master</pre>
</div>
</div>
<p>Have fun <img src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Installing_14478/wlEmoticon-winkingsmile.png" alt="Winking smile" style="border-style: none;" class="wlEmoticon wlEmoticon-winkingsmile" /></p>