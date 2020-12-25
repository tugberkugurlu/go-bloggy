---
id: b9717ed7-45d6-473b-b8f4-cbb5e875a5ba
title: Versioning Software Builds Based on Git Tags and Semantic Versioning (SemVer)
abstract: I have been using a technique to set the build version on my CI (continuous
  integration) system, Travis CI, based on Git tags and semantic versioning (SemVer).
  In this post, I want to share this with you and give you an implementation of that
  in Bash.
created_at: 2016-04-16 14:20:00 +0000 UTC
tags:
- Bash
- Continuous Delivery
- Continuous Integration
- Git
- SemVer
- Travis CI
slugs:
- versioning-software-builds-based-on-git-tags-and-semantic-versioning-semver
---

<p>Let's start this post by setting the stage first and then move onto the problem. When a build is kicked off for your application/library/etc. on a CI (continuous integration) system like <a href="https://travis-ci.org/">Travis CI</a> or <a href="https://www.appveyor.com/">AppVeyor</a>, you are most probably flowing a version number for that build no matter what type of tech stack you use. This is mostly to relate the artifacts, which the build will produce (e.g. <a href="https://docs.docker.com/engine/reference/commandline/images/">Docker images</a>, <a href="https://docs.nuget.org/create/creating-and-publishing-a-package">NuGet packages</a>, <a href="https://msdn.microsoft.com/en-us/library/hk5f40ct(v=vs.90).aspx">.NET assemblies</a>, etc.), with a particular context. This is really useful to be able to communicate and correlate stuff. A few scenarios: <ul> <li>Hey Mark, please take a look at foobar-1.2.3-rc.657 from our CI Docker registry. That has the issue I have mentioned. You can check it on that image. <li>Ow, barfoo-2.2.3-beta.362 NuGet package content misses a few assemblies that should have been there. Let's go back to build logs for this and check what went wrong.</li></ul> <p>Convinced? Good :) Otherwise, you won't find the rest of the article useful. <p>The other case is to flow a version number when you actually want to produce a release for your defined environments (e.g. acceptance, staging, production). In this case, you usually don't want to give an arbitrary version to your artifacts because the version will carry the high level information about the changes. There are three important intentions you can give here: <ul> <li>I am releasing something which has no behavior changes <li>I am releasing a new feature which doesn't break my existing consumers <li>Dude, brace yourself! I will break the World into half! </li></ul> <p>You can see <a href="http://semver.org/spec/v2.0.0.html">Semantic Versioning 2.0.0</a> for more information about this. <p>So, what happens here is that we want to let the CI system decide on the version at some cases and take control over which version number to flow in some other cases. Actually, the first statement is not quite correct because you still want to have partial control over what version number to flow for your non-release builds. Here is an example case to highlight what I mean: <ul> <li>You started developing your application and shipped version 1.0.0. <li>Your CI system started flowing prerelease version based on 1.0.0 and also attached the build number to that version (e.g. 1.0.0-beta.54). Notice that it's wrong at this stage because you already shipped v1.0.0. So, it should really be something like 1.0.1-beta-54. <li>Now, you are shipping version 1.1.0 as you introduced a new feature. <li>After that change, you keep building the software and CI system keeps flowing version 1.1.0 based versions. This is a bit bad as you now don't have the chronological order and version order correlation.</li></ul> <p>So, what we want here is to assign a version based on the latest release version, which means that you want to have control over this process of assigning a version number. I have seen people having a text file inside the repository to hold the latest release version but that's a bit manual. I assume you kick a release somewhere and you already assign a version at that stage for releases. So, wouldn't it be bad to leverage this?</p> <p>So, you probably understood my problems here :) Now, let me introduce a few key pieces which will play a role to solve this problem and then later, I will move onto the actual implementation to solve the problem.</p> <h3>Git Tags</h3> <p><a href="https://git-scm.com/book/en/v2/Git-Basics-Tagging">Tagging is a feature of Git</a> which allows you to mark specific points in repository's history. As the Git manual also states, people typically use this functionality to mark release points. This is super convenient for our needs here and gets two important things sorted for us: <ul> <li>A kick-off point for releases. Ultimately, release process will be kicked off when you tag a repository and push that tag to your remote. <li>Deciding the base version based on the latest release version.</li></ul> <h3>Semantic Versioning</h3> <p>So, we have the tags. However, it doesn't mean that every tag is a valid version and you can also use Git's tagging feature for some other purposes. This is where SemVer comes into picture and you can safely assume that any tag which is a valid SemVer is for a release. This makes your life so much easier as you can rely on built-in tools like <a href="https://www.tugberkugurlu.com/archive/node-semver-cli-tool-for-semantic-versioning-2-0-0">node-semver</a> to help you out (as we will see shortly).</p> <p>The other thing we have in the mix is to be able to increment the build version after a release. For example, we release version 2.5.6. The next build right after the release should have the version number bigger than 2.5.6. Seems easy as you can just increment the patch version, right? No! 2.5.6-beta is also a valid SemVer. We can go further with 2.5.6-beta.5+736287 which is also a valid SemVer. So, there is a pre-defined spec here and we can again leverage tools like node-semver to work with this domain nicely.</p> <h3>Solution and Bash Implementation</h3> <p>OK, all this information is super useful but how to make it work? Let me walk you through a solution I have introduced recently on a few of the projects I am working on. It's very trivial but that useful at the same time. However, keep in mind that there might be a few things I might have missed as I have been applying this not for a long time. <strong>In fact, here might even be better techniques on this that you know. If so, please comment here. I would love to hear them!</strong></p> <p>I want to example this in two stages and bring them together at the end.</p> <h4>Deciding on a Base Version</h4> <p>When the build is kicked off, one of the first things to do is to decide a base version. This is fairly trivial and here is the flow chart to describe this decision making process:</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6ec918c2-37e8-4e0f-890a-2ffb4e0f3041.jpg"><img title="base-version" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="base-version" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/4322b563-d61a-4057-ad1a-5db3bdee0225.jpg" width="530" height="484"></a></p> <p>Here is how the implementation looks like in <a href="https://en.wikipedia.org/wiki/Bash_(Unix_shell)">Bash</a>:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: green">#!/bin/bash</span>

baseVersion<span style="color: gray">=</span>0.0.0<span style="color: gray">-</span>0
<span style="color: blue">if</span> semver <span style="color: #a31515">"ignorethis"</span> $(git tag <span style="color: gray">-</span>l) &amp;<span style="color: gray">&gt;</span><span style="color: gray">/</span>dev<span style="color: gray">/</span>null
then
    baseVersion<span style="color: gray">=</span>$(semver $((semver $(git tag <span style="color: gray">-</span>l)) | tail <span style="color: gray">-</span>n1) <span style="color: gray">-</span>i prerelease)
fi</pre></div></div>
<blockquote>
<p>Keep in mind that I am fairly new to Bash. So, there might be wrong/bad usages here.</p></blockquote>
<p>To explain what happens here with a bit more details:</p>
<ul>
<li>We get all the tags for the repository as a list by running git tag -l
<li>We pass this list to semver command-line tool to filter the invalid SemVer strings. Notice that there is another parameter we pass to semver here called "ignorethis". It's just there to cover cases when there is no tag so that semver command-line tool can return non-zero exit code.
<li>If semver command-line tool exits with 0, we know that there is at least one tag which is a valid SemVer. So, we run tail -n1 on the semver output to retrieve the latest version and we increment it on its prerelease identifier. This is now our base version.
<li>If there are no valid SemVer tags on the repository, we set 0.0.0-0 as the base version.</li></ul>

<h4>Decide on a Build Version</h4>
<p>Now we have a base version and we now need to decide on a build version based on that. This is a bit more involved but again, very trivial to implement. Here is another flow chart to describe this decision making process:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d0b999d7-cd88-4a75-9677-3d574755ae78.jpg"><img title="build-version" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="build-version" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c7a3cd2a-5375-4889-b0ca-adbed6e5e709.jpg" width="469" height="484"></a></p>
<p>And, here is how the implementation looks like in Bash (specific to Travis CI as it uses Travis CI specific environment variables):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">if</span> <span style="color: gray">[</span><span style="color: teal"> -z "$TRAVIS_TAG" </span><span style="color: gray">]</span>;
then
    <span style="color: blue">if</span> <span style="color: gray">[</span><span style="color: teal"> -z "$TRAVIS_BRANCH" </span><span style="color: gray">]</span>;
    then
        <span style="color: green"># can add the build metadata to indicate this is pull request build</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span><span style="color: #a31515">"$baseVersion.$TRAVIS_BUILD_NUMBER"</span>;
    <span style="color: blue">else</span>
        <span style="color: green"># can add the build metadata to indicate this is a branch build</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span><span style="color: #a31515">"$baseVersion.$TRAVIS_BUILD_NUMBER"</span>;
    fi
<span style="color: blue">else</span> 
    <span style="color: blue">if</span> ! semver <span style="color: orangered">$TRAVIS_TAG</span> &amp;<span style="color: gray">&gt;</span><span style="color: gray">/</span>dev<span style="color: gray">/</span>null
    then
        <span style="color: green"># can add the build metadata to indicate this is a tag build which is not a SemVer</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span><span style="color: #a31515">"$baseVersion.$TRAVIS_BUILD_NUMBER"</span>;
    <span style="color: blue">else</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span>$(semver <span style="color: orangered">$TRAVIS_TAG</span>);
    fi 
fi</pre></div></div>
<blockquote>
<p>Notice that I am echoing commands rather than directly calling them. This is because of a fact that Travis CI doesn't flow the exports which happens inside a script file. Maybe it does but I was not able to get it working. Anyways, I am calling this script inside my .travis.yml file by evaluating the output like this: eval $(./scripts/set-build-version.sh)</p></blockquote>
<p>I am not going to separately explain how this works as the flow chart is very easy to grasp (also the Bash script). However, one thing which is worth mentioning is the branch check. After we check if the build is for a branch, we do the same operation no matter what. This is OK for my use case but you can add special metadata to your version in order to indicate which branch the build has happened or whether it was a pull request.</p>
<h3>Conclusion</h3>
<p>I find this solution very straight forward to pick the version of the build and have a central way of kicking of a release process. I applied this on <a href="https://github.com/tugberkugurlu/AspNetCore.Identity.MongoDB">AspNetCore.Identity.MongoDB</a> project, a MongoDB data store adapter for <a href="https://github.com/aspnet/Identity">ASP.NET Core identity</a>. You can also see <a href="https://github.com/tugberkugurlu/AspNetCore.Identity.MongoDB/blob/eb0022b5bf82e6988d792e854f38e04e5447e8a1/.travis.yml#L26">how I am setting the build version</a>, <a href="https://github.com/tugberkugurlu/AspNetCore.Identity.MongoDB/blob/eb0022b5bf82e6988d792e854f38e04e5447e8a1/.travis.yml#L27">how I am using it</a> and <a href="https://github.com/tugberkugurlu/AspNetCore.Identity.MongoDB/blob/eb0022b5bf82e6988d792e854f38e04e5447e8a1/.travis.yml#L36">how I am kicking off a release process</a>.</p>
<p>To bring everything together, here is <a href="https://github.com/tugberkugurlu/AspNetCore.Identity.MongoDB/blob/eb0022b5bf82e6988d792e854f38e04e5447e8a1/scripts/set-build-version.sh">the entire script to set the build version</a>:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: green">#!/bin/bash</span>

baseVersion<span style="color: gray">=</span>0.0.0<span style="color: gray">-</span>0
<span style="color: blue">if</span> semver <span style="color: #a31515">"ignorethis"</span> $(git tag <span style="color: gray">-</span>l) &amp;<span style="color: gray">&gt;</span><span style="color: gray">/</span>dev<span style="color: gray">/</span>null
then
    baseVersion<span style="color: gray">=</span>$(semver $((semver $(git tag <span style="color: gray">-</span>l)) | tail <span style="color: gray">-</span>n1) <span style="color: gray">-</span>i prerelease)
fi

<span style="color: blue">if</span> <span style="color: gray">[</span><span style="color: teal"> -z "$TRAVIS_TAG" </span><span style="color: gray">]</span>;
then
    <span style="color: blue">if</span> <span style="color: gray">[</span><span style="color: teal"> -z "$TRAVIS_BRANCH" </span><span style="color: gray">]</span>;
    then
        <span style="color: green"># can add the build metadata to indicate this is pull request build</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span><span style="color: #a31515">"$baseVersion.$TRAVIS_BUILD_NUMBER"</span>;
    <span style="color: blue">else</span>
        <span style="color: green"># can add the build metadata to indicate this is a branch build</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span><span style="color: #a31515">"$baseVersion.$TRAVIS_BUILD_NUMBER"</span>;
    fi
<span style="color: blue">else</span> 
    <span style="color: blue">if</span> ! semver <span style="color: orangered">$TRAVIS_TAG</span> &amp;<span style="color: gray">&gt;</span><span style="color: gray">/</span>dev<span style="color: gray">/</span>null
    then
        <span style="color: green"># can add the build metadata to indicate this is a tag build which is not a SemVer</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span><span style="color: #a31515">"$baseVersion.$TRAVIS_BUILD_NUMBER"</span>;
    <span style="color: blue">else</span>
        echo export PROJECT_BUILD_VERSION<span style="color: gray">=</span>$(semver <span style="color: orangered">$TRAVIS_TAG</span>);
    fi 
fi</pre></div></div>
<p>I hope this will be useful to you in some way and as said, if you have a similar technique or a practice that you apply for this case, please share it. Now, go and enjoy this spectacular weekend ;)</p>  