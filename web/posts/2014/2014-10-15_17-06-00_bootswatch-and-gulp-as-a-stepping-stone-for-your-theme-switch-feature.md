---
id: 92c0119c-1ef5-4463-8565-0b627de3b664
title: Bootswatch and Gulp as a Stepping Stone for Your Theme Switch Feature
abstract: This short blog posts shows you a way of combining bootswatch and gulp together
  to have an easily useable theme switching support for your web application.
created_at: 2014-10-15 17:06:00 +0000 UTC
tags:
- gulp
- JavaScript
slugs:
- bootswatch-and-gulp-as-a-stepping-stone-for-your-theme-switch-feature
---

<p><a href="http://bootswatch.com/">Bootswatch</a> has awesome themes which are built on top of <a href="http://getbootstrap.com/">Bootstrap</a>. To integrate those themes to your web application, you have a few options. One of them is to build the css file by combining bootstrap.less file with the less files provided by bootswatch. This is the way I chose to go with. However, what I actually wanted was a little bit more complex:</p> <ul> <li>Install bootswatch with bower and have every theme available for me to use.</li> <li>Compile each bootswtch theme to CSS.</li> <li>Concatenate and minify each generated bootstrap file with my other CSS files and have separate single CSS file available for me for each theme.</li></ul> <p>With this way, switching between themes would be just as easy as changing the filename suffix. I achieved this with <a href="http://gulpjs.com/">gulp</a> and you can find the sample <a href="https://github.com/tugberkugurlu">in my GitHub account</a>: <a href="https://github.com/tugberkugurlu/gulp-bootswatch-sample">gulp-bootswatch-sample</a>. To get the sample up and running, run the following commands in order:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: gray">/</span><span style="color: gray">/</span> install bower globally <span style="color: blue">if</span> you don't have it already
npm install bower <span style="color: gray">-</span>g

<span style="color: gray">/</span><span style="color: gray">/</span> install gulp globally <span style="color: blue">if</span> you don't have it already
npm install gulp <span style="color: gray">-</span>g

<span style="color: gray">/</span><span style="color: gray">/</span> navigate to my sample projects root folder and run the following commands
npm install
bower install

<span style="color: gray">/</span><span style="color: gray">/</span> finally run gulp the tasks
gulp</pre></div></div>
<p>Here is how my gulpfile.js file looks like:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">var</span> gulp = require(<span style="color: #a31515">'gulp'</span>),
    minifycss = require(<span style="color: #a31515">'gulp-minify-css'</span>),
    concat = require(<span style="color: #a31515">'gulp-concat'</span>),
    less = require(<span style="color: #a31515">'gulp-less'</span>),
    gulpif = require(<span style="color: #a31515">'gulp-if'</span>),
    order = require(<span style="color: #a31515">'gulp-order'</span>),
    gutil = require(<span style="color: #a31515">'gulp-util'</span>),
    rename = require(<span style="color: #a31515">'gulp-rename'</span>),
    foreach = require(<span style="color: #a31515">'gulp-foreach'</span>),
    debug = require(<span style="color: #a31515">'gulp-debug'</span>),
    path =require(<span style="color: #a31515">'path'</span>),
    merge = require(<span style="color: #a31515">'merge-stream'</span>),
    del = require(<span style="color: #a31515">'del'</span>);

gulp.task(<span style="color: #a31515">'default'</span>, [<span style="color: #a31515">'clean'</span>], <span style="color: blue">function</span>() {
    gulp.start(<span style="color: #a31515">'fonts'</span>, <span style="color: #a31515">'styles'</span>);
});

gulp.task(<span style="color: #a31515">'clean'</span>, <span style="color: blue">function</span>(cb) {
    del([<span style="color: #a31515">'assets/css'</span>, <span style="color: #a31515">'assets/js'</span>, <span style="color: #a31515">'assets/less'</span>, <span style="color: #a31515">'assets/img'</span>, <span style="color: #a31515">'assets/fonts'</span>], cb)
});

gulp.task(<span style="color: #a31515">'fonts'</span>, <span style="color: blue">function</span>() {
    
    <span style="color: blue">var</span> fileList = [
        <span style="color: #a31515">'bower_components/bootstrap/dist/fonts/*'</span>, 
        <span style="color: #a31515">'bower_components/fontawesome/fonts/*'</span>
    ];
    
    <span style="color: blue">return</span> gulp.src(fileList)
        .pipe(gulp.dest(<span style="color: #a31515">'assets/fonts'</span>));
});

gulp.task(<span style="color: #a31515">'styles'</span>, <span style="color: blue">function</span>() {
    
    <span style="color: blue">var</span> baseContent = <span style="color: #a31515">'@import "bower_components/bootstrap/less/bootstrap.less";@import "bower_components/bootswatch/$theme$/variables.less";@import "bower_components/bootswatch/$theme$/bootswatch.less";@import "bower_components/bootstrap/less/utilities.less";'</span>;

    <span style="color: blue">var</span> isBootswatchFile = <span style="color: blue">function</span>(file) {
        <span style="color: blue">var</span> suffix = <span style="color: #a31515">'bootswatch.less'</span>;
        <span style="color: blue">return</span> file.path.indexOf(suffix, file.path.length - suffix.length) !== -1;
    }
    
    <span style="color: blue">var</span> isBootstrapFile = <span style="color: blue">function</span>(file) {
        <span style="color: blue">var</span> suffix = <span style="color: #a31515">'bootstrap-'</span>,
            fileName = path.basename(file.path);
        
        <span style="color: blue">return</span> fileName.indexOf(suffix) == 0;
    }
    
    <span style="color: blue">var</span> fileList = [
        <span style="color: #a31515">'client/less/main.less'</span>, 
        <span style="color: #a31515">'bower_components/bootswatch/**/bootswatch.less'</span>, 
        <span style="color: #a31515">'bower_components/fontawesome/css/font-awesome.css'</span>
    ];
    
    <span style="color: blue">return</span> gulp.src(fileList)
        .pipe(gulpif(isBootswatchFile, foreach(<span style="color: blue">function</span>(stream, file) {
            <span style="color: blue">var</span> themeName = path.basename(path.dirname(file.path)),
                content = replaceAll(baseContent, <span style="color: #a31515">'$theme$'</span>, themeName),
                file = string_src(<span style="color: #a31515">'bootstrap-'</span> +  themeName + <span style="color: #a31515">'.less'</span>, content);

            <span style="color: blue">return</span> file;
        })))
        .pipe(less())
        .pipe(gulp.dest(<span style="color: #a31515">'assets/css'</span>))
        .pipe(gulpif(isBootstrapFile, foreach(<span style="color: blue">function</span>(stream, file) {
            <span style="color: blue">var</span> fileName = path.basename(file.path),
                themeName = fileName.substring(fileName.indexOf(<span style="color: #a31515">'-'</span>) + 1, fileName.indexOf(<span style="color: #a31515">'.'</span>));
            
            <span style="color: blue">return</span> merge(stream, gulp.src([<span style="color: #a31515">'assets/css/font-awesome.css'</span>, <span style="color: #a31515">'assets/css/main.css'</span>]))
                .pipe(concat(<span style="color: #a31515">'style-'</span> + themeName + <span style="color: #a31515">".css"</span>))
                .pipe(gulp.dest(<span style="color: #a31515">'assets/css'</span>))
                .pipe(rename({suffix: <span style="color: #a31515">'.min'</span>}))
                .pipe(minifycss())
                .pipe(gulp.dest(<span style="color: #a31515">'assets/css'</span>));
        })))
});

<span style="color: blue">function</span> escapeRegExp(string) {
    <span style="color: blue">return</span> string.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, <span style="color: #a31515">"\\$1"</span>);
}

<span style="color: blue">function</span> replaceAll(string, find, replace) {
  <span style="color: blue">return</span> string.replace(<span style="color: blue">new</span> RegExp(escapeRegExp(find), <span style="color: #a31515">'g'</span>), replace);
}

<span style="color: blue">function</span> string_src(filename, string) {
  <span style="color: blue">var</span> src = require(<span style="color: #a31515">'stream'</span>).Readable({ objectMode: <span style="color: blue">true</span> })
  src._read = <span style="color: blue">function</span> () {
    <span style="color: blue">this</span>.push(<span style="color: blue">new</span> gutil.File({ cwd: <span style="color: #a31515">""</span>, base: <span style="color: #a31515">""</span>, path: filename, contents: <span style="color: blue">new</span> Buffer(string) }))
    <span style="color: blue">this</span>.push(<span style="color: blue">null</span>)
  }
  <span style="color: blue">return</span> src
}</pre></div></div>
<p>End result :)</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/31076db0-3bbe-40b8-8f51-907040333201.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/38476835-0fc6-416a-a627-f42ddd0ca14b.png" width="461" height="484"></a></p>
<p>I wished I had something like this when I started doing this so I hope this was helpful to you somehow. Enjoy!</p>  