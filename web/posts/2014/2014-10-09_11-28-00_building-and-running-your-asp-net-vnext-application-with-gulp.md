---
id: ed485177-9a35-485a-8e2f-b7eaeec725dd
title: Building and Running Your ASP.NET vNext Application with Gulp
abstract: Wanna see ASP.NET vNext and Gulp working together? You are at the right
  place :) Let's have look at gulp-aspnet-k, a little plugin that I have created for
  ASP.NET vNext gulp integration.
created_at: 2014-10-09 11:28:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET vNext
- gulp
slugs:
- building-and-running-your-asp-net-vnext-application-with-gulp
---

<p>I love <a href="http://gulpjs.com/">gulp</a>! It has been only a few weeks since I started getting my hands dirty with gulp but it’s ridiculously simple and good. I am actually using gulp with one of my <a href="http://www.tugberkugurlu.com/archive/getting-started-with-asp-net-vnext-by-setting-up-the-environment-from-scratch">ASP.NET vNext</a> applications to see how they fit together. I am <a href="https://github.com/plus3network/gulp-less">compiling my less files</a>, doing <a href="https://github.com/wearefractal/gulp-concat">concatenation</a> and minification for <a href="https://www.npmjs.org/package/gulp-uglify">scripts</a>/<a href="https://www.npmjs.org/package/gulp-minify-css">styles</a> files with gulp. It’s also extremely comforting that gulp has <a href="https://github.com/gulpjs/gulp/blob/master/docs/API.md#gulpwatchglob-opts-cb">file watch capability</a>. For example, I can change my less files during the development and they are being recompiled as I save them:</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1d418de3-0578-4256-9e5b-5da1751f827d.gif"><img title="gulp-watch" style="display: inline" alt="gulp-watch" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5a9e45f8-2a80-48c6-a486-ce4a2cff8468.gif" width="640" height="361"></a></p> <p>As I am working with my ASP.NET vNext application in this project, I found myself jumping between command prompt windows. I imagined for a second that it would be cool to have a gulp plugin for ASP.NET vNext to build and run the application. It would also take advantage of k --watch so that it would restart the host when any code files are changed. Then, I started digging into it and finally, I managed to get <a href="https://github.com/tugberkugurlu/gulp-aspnet-k">gulp-aspnet-k</a> (ASP.NET vNext Gulp Plugin) out :) <a href="https://www.npmjs.org/package/gulp-aspnet-k">gulp-aspnet-k is also available on npm</a> and you can install it right from there. Check out <a href="https://github.com/tugberkugurlu/gulp-aspnet-k#readme">the readme</a> for further info about its usage.</p> <p>This is an insanely simple plugin which wraps <a href="https://github.com/aspnet/KRuntime/blob/dev/scripts/kpm.cmd">kpm</a> and <a href="https://github.com/aspnet/KRuntime/blob/dev/scripts/k.cmd">k</a> commands for you. In its simplest form, its usage is as below:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">var</span> gulp = require(<span style="color: #a31515">'gulp'</span>),
    aspnetk = require(<span style="color: #a31515">"gulp-aspnet-k"</span>);

gulp.task(<span style="color: #a31515">'default'</span>, <span style="color: blue">function</span>(cb) {
    <span style="color: blue">return</span> gulp.start(<span style="color: #a31515">'aspnet-run'</span>);
});

gulp.task(<span style="color: #a31515">'aspnet-run'</span>, aspnetk());</pre></div></div>
<p>You can find a sample application that uses gulp-aspnet-k plugin in my <a href="https://github.com/tugberkugurlu/AspNetVNextSamples">ASP.NET vNext samples</a> repository: <a href="https://github.com/tugberkugurlu/AspNetVNextSamples/tree/7ec8a005d95e120347ef2b2ecc76995d461cac46/src/GulpSample">GulpSample</a>. It is also working with gulp watch in peace.</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/86490f28-b664-4aa3-bf93-92d4405ffdc9.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/0fa0fe89-7efc-4b91-bd1c-6f5cf92a89de.png" width="451" height="484"></a></p>
<p>Be sure to watch the following short video of mine to see this tiny plugin in action.</p><iframe height="281" src="//player.vimeo.com/video/108458085" frameborder="0" width="500" allowfullscreen mozallowfullscreen webkitallowfullscreen></iframe>
<p><a href="http://vimeo.com/108458085">ASP.NET vNext with Gulp</a> from <a href="http://vimeo.com/user6670252">Tugberk Ugurlu</a> on <a href="https://vimeo.com">Vimeo</a>.</p>
<p>Keep in mind that, currently, this only works on windows :s Also remember that little things matter in your daily life and this thing is one of them :)</p>  