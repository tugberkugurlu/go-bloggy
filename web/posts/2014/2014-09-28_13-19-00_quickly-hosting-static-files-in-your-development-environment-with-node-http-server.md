---
title: Quickly Hosting Static Files In Your Development Environment with Node http-server
abstract: 'Yesterday, I was looking for something to have a really quick test space
  on my machine to play with AngularJS and I found http-server: a simple zero-configuration
  command-line http server.'
created_at: 2014-09-28 13:19:00 +0000 UTC
tags:
- HTTP
- nodejs
- Tips
slugs:
- quickly-hosting-static-files-in-your-development-environment-with-node-http-server
---

<p>Yesterday, I was looking into <a href="https://angularjs.org/">AngularJS</a> which is a long overdue for me. I wanted to have a really quick test space on my machine to play with AngularJS. I could just go with <a href="http://plnkr.co/">Plunker</a> or <a href="http://jsfiddle.net/">JSFiddle</a> but I wasn’t in the mood for an online editor.</p> <p>So, I first installed <a href="http://bower.io/">Bower</a> and then I installed a few libraries to get me going. I made sure to save them inside my <a href="https://github.com/bower/bower.json-spec">bower.json</a> file, too:</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/4f8a9a62-4e15-4917-b90f-a38bff41e746.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/368b969a-819d-45d3-9b09-2bc0f6d89e62.png" width="640" height="484"></a></p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>bower install angular <span style="color: gray">-</span>S
bower install bootstrap <span style="color: gray">-</span>S
bower install underscore –S</pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a8e17db8-5dfc-47e2-973a-f459f2e7a387.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/7ab66c3e-5835-44e8-9288-2ff397991d7d.png" width="644" height="408"></a></p>
<p>Then, I installed my node module to help me automate a few processes with <a href="http://gulpjs.com/">gulp</a>. </p>
<blockquote>
<p>BTW, do you know VS has support for running you gulp and grunt tasks? Check out <a href="http://www.hanselman.com/">Scott Hanselman</a>’s post if you don’t: <a href="http://www.hanselman.com/blog/introducinggulpgruntbowerandnpmsupportforvisualstudio.aspx">Introducing Gulp, Grunt, Bower, and npm support for Visual Studio</a></p></blockquote>
<p>Finally, I have written a few lines of code to get me going. The state the this tiny app can be found <a href="https://github.com/tugberkugurlu/angularjs-getting-started/tree/181acbc22e3ec6463156073ccce19473260476ec">under my GtHub account</a>. I was at the point where I would like to see it working in my browser:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/11806dad-b86d-403e-b506-0748d7e8d4c9.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e25aaa97-83c3-4e5d-9aa1-84e3bad216a5.png" width="644" height="332"></a></p>
<p>As a developer who have been heavily working with .NET for a while, I wanted to hit F5 :) Jokes aside, what I only want here is something so simple that would host my static files inside my directory with no configuration required. While searching a great option for this, I found <a href="https://github.com/nodeapps/http-server">http-server</a>: a simple zero-configuration command-line http server. </p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5c90437f-00e0-4bc6-b5a0-103670c3907b.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/13182d4a-862e-44a7-b885-ade0368b7f8f.png" width="590" height="484"></a></p>
<p>I simply installed it as an global <a href="https://www.npmjs.org/">npm</a> module:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/7f68d764-4473-49cb-80fe-df84662c5945.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ceb8856e-0756-409d-9025-1233260cfbf6.png" width="644" height="215"></a></p>
<p>All done! It’s now under my path and I can simply cd into my root web site directory and run http-server there.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/bc1f491d-f193-4a08-be0e-899d501ed5b4.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/560c500e-bdbc-4552-93f5-abb055fcb01c.png" width="598" height="484"></a></p>
<p>Super simple! It also observes the changes. So, you can modify your files as the http-server serves them. You can even combine this with <a href="http://rhumaric.com/2014/01/livereload-magic-gulp-style/">LiveReload + a gulp task</a>.</p>  