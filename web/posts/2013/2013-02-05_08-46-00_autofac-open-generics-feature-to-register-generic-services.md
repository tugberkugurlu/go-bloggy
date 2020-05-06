---
title: Autofac Open Generics Feature to Register Generic Services
abstract: Autofac, an awesome IoC container for .NET platform, has an out of the box
  generic service registration feature which we will quickly cover in this blog post.
created_at: 2013-02-05 08:46:00 +0000 UTC
tags:
- .net
- Autofac
slugs:
- autofac-open-generics-feature-to-register-generic-services
---

<p>This is going to be a quick and dirty blog post but hopefully, will take this giant stupidity out of me. <a href="http://code.google.com/p/autofac/">Autofac</a>, an awesome <a href="http://martinfowler.com/articles/injection.html">IoC container</a> for .NET platform, has an out of the box <a href="http://code.google.com/p/autofac/wiki/OpenGenerics">generic service registration feature</a> and I assume nearly all IoC containers have this today which makes me feel stupid because I have been knowing this for only a month or so <img src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Autofac-Open-Generics_AE48/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" /> I was doing something like below before.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> RegisterRepositories(ContainerBuilder builder) {
 
    Type baseEntityType = <span style="color: blue;">typeof</span>(BaseEntity);
    Assembly assembly = baseEntityType.Assembly;
    IEnumerable&lt;Type&gt; entityTypes = assembly.GetTypes().Where(
        x =&gt; x.IsSubclassOf(baseEntityType));
        
    <span style="color: blue;">foreach</span> (Type type <span style="color: blue;">in</span> entityTypes) {
 
        builder.RegisterType(<span style="color: blue;">typeof</span>(EntityRepository&lt;&gt;)
               .MakeGenericType(type))
               .As(<span style="color: blue;">typeof</span>(IEntityRepository&lt;&gt;).MakeGenericType(type))
               .InstancePerApiRequest();
    }
}</pre>
</div>
</div>
<p>Then, <a href="https://twitter.com/benfosterdev">Ben Foster</a> pinged me on twitter:</p>
<blockquote class="twitter-tweet">
<p>@<a href="https://twitter.com/tourismgeek">tourismgeek</a> what, no For(type of(IFoo&lt;&gt;).Use(typeof(Foo&lt;&gt;)? Even Ninject supports that.</p>
&mdash; Ben Foster (@benfosterdev) <a href="https://twitter.com/benfosterdev/status/281315791910080513">December 19, 2012</a></blockquote>
<script src="//platform.twitter.com/widgets.js"></script>
<p>This tweet made me look for alternative approaches and I found out the Autofac's generic service registration feature. Here is how it looks like now:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> RegisterRepositories(ContainerBuilder builder) {
 
    builder.RegisterGeneric(<span style="color: blue;">typeof</span>(EntityRepository&lt;&gt;))
           .As(<span style="color: blue;">typeof</span>(IEntityRepository&lt;&gt;))
           .InstancePerApiRequest();
}</pre>
</div>
</div>
<p>Way better! Autofac also respects generic type constraints. Here is a quote from the <a href="http://code.google.com/p/autofac/wiki/OpenGenerics">Autofac documentation</a>:</p>
<blockquote>
<p>Autofac respects generic type constraints. If a constraint on the implementation type makes it unable to provide a service the implementation type will be ignored.</p>
</blockquote>
<p>If you didn't know this feature before, you do know it now <img src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Autofac-Open-Generics_AE48/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" /> Enjoy it!</p>