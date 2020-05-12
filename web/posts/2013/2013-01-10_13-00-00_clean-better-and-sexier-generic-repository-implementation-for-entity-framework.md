---
id: af58f326-e330-4b5f-8cdd-a38052160864
title: Clean, Better, and Sexier Generic Repository Implementation for Entity Framework
abstract: With the new release of the GenericRepository.EntityFramework package, we
  now have clean, better and sexier generic repository implementation for Entity Framework.
  Enjoy!
created_at: 2013-01-10 13:00:00 +0000 UTC
tags:
- .net
- ASP.NET MVC
- ASP.NET Web API
- DbContext
- Entity Framework
- MS SQL
slugs:
- clean-better-and-sexier-generic-repository-implementation-for-entity-framework
---

<p>I have written a few blog posts about my experience <a href="http://www.tugberkugurlu.com/archive/generic-repository-pattern-entity-framework-asp-net-mvc-and-unit-testing-triangle">on applying the Generic Repository pattern to Entity Framework</a> and I even made a <a href="http://nuget.org/packages/GenericRepository.EF">NuGet package for my na&iuml;ve implementation</a>. Even if that looked OK at the time for me, I had troubles about my implementation for couple of reasons:</p>
<ul>
<li>The implementations inside the NuGet package didn&rsquo;t allow the developer to share the DbContext instance between repositories (per-request for example). </li>
<li>When the generic repository methods weren&rsquo;t enough, I was creating new repository interfaces and classes based on the generic ones. This was the biggest failure for me and didn&rsquo;t scale very well as you can imagine. </li>
<li>There were no pagination support. </li>
<li>As each repository take a direct dependency on DbContext and it is impossible to entirely fake the DbContext, generic repositories needed to be mocked so that we could test with them. However, it would be just very useful to pass a fake DbContext instance into the repository implementation itself and use it as fake.</li>
</ul>
<p>With the new release of the <a href="http://nuget.org/packages/GenericRepository.EntityFramework">GenericRepository.EntityFramework</a> package, the all of the above problems have their solutions. The source code for this release is available under the <a href="https://github.com/tugberkugurlu/GenericRepository/tree/master">master branch of the repository</a> and you can also see the ongoing work for the final release under the <a href="https://github.com/tugberkugurlu/GenericRepository/tree/v0.3.0">v0.3.0 branch</a>. The NuGet package is available as pre-release for now. So, you need to use the &ndash;pre switch to install it.</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package GenericRepository.EntityFramework -Pre </code></p>
</div>
<p>The old <a href="http://nuget.org/packages/GenericRepository.EF">GenericRepository.EF</a> package is still around and I will update it, too but it&rsquo;s now unlisted and only thing it does is to install the GenericRepository.EntityFramework package.</p>
<p>I also included <a href="https://github.com/tugberkugurlu/GenericRepository/tree/master/samples">a sample application which shows the usage briefly</a>. I will complete the sample and extend it further for a better view. Definitely check this out!</p>
<p>Let&rsquo;s dive right in and see what is new and cool.</p>
<h3>IEntity and IEntity&lt;TId&gt; Interfaces</h3>
<p>I introduced two new interfaces: <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository/IEntity.cs">IEntity</a> and <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository/IEntity'1.cs">IEntity&lt;TId&gt;</a> and each of your entity classes needs to implement one of these. As you can see from the implementation, IEntity just implements the IEntity&lt;int&gt; and you can use IEntity if you are using integer based Ids. The reason why I added these is make the GetSingle method work.</p>
<h3>Use EntitiesContext Instead of DbContext</h3>
<p>Instead of deriving your context class from DbContext, you now need to take the <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository.EntityFramework/EntitiesContext.cs">EntitiesContext</a> as the base class for your context. If you have an existing context class based on DbContext, changing it to use EntitiesContext should not break it. The EntitiesContext class has all the same constructors as DbContext. So, you can also use those. Here is the sample:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> AccommodationEntities : EntitiesContext {

    <span style="color: green;">// NOTE: You have the same constructors as the DbContext here. E.g:</span>
    <span style="color: green;">// public AccommodationEntities() : base("nameOrConnectionString") { }</span>

    <span style="color: blue;">public</span> IDbSet&lt;Country&gt; Countries { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> IDbSet&lt;Resort&gt; Resorts { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> IDbSet&lt;Hotel&gt; Hotels { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>Then, through your IoC container, you can register your context as a new instance for <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository.EntityFramework/IEntitiesContext.cs">IEntitiesContext</a> per a particular scope. The below example uses Autofac to do that for an <a href="http://www.asp.net/web-api">ASP.NET Web API</a> application:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> RegisterDependencies(HttpConfiguration config) {

    <span style="color: blue;">var</span> builder = <span style="color: blue;">new</span> ContainerBuilder();
    builder.RegisterApiControllers(Assembly.GetExecutingAssembly());

    <span style="color: green;">// Register IEntitiesContext</span>
    builder.Register(_ =&gt; <span style="color: blue;">new</span> AccommodationEntities())
           .As&lt;IEntitiesContext&gt;().InstancePerApiRequest();

    <span style="color: green;">// TODO: Register repositories here</span>

    config.DependencyResolver = 
        <span style="color: blue;">new</span> AutofacWebApiDependencyResolver(builder.Build());
}</pre>
</div>
</div>
<h3>IEntityRepository&lt;TEntity&gt; and EntityRepository&lt;TEntity&gt;</h3>
<p>Here is the real meat of the package: IEntityRepository and EntityRepository. Same as the IEntity and IEntity&lt;TId&gt;, we have two different IEntityRepository generic interfaces: <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository.EntityFramework/IEntityRepository'1.cs">IEntityRepository&lt;TEntity&gt;</a> and <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository.EntityFramework/IEntityRepository'2.cs">IEntityRepository&lt;TEntity, TId&gt;</a>. They have their implementations under the same generic signature: <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository.EntityFramework/EntityRepository'1.cs">EntityRepository&lt;TEntity&gt;</a> and <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository.EntityFramework/EntityRepository'2.cs">EntityRepository&lt;TEntity, TId&gt;</a>. The big improvement now is that EntityRepository generic repository implementation accepts an IEntitiesContext implementation through its constructor. This, for example, enables you to use the same DbContext (IEntitiesContext implementation in our case, which is EntitiesContext by default) instance per-request for your <a href="http://www.asp.net/mvc">ASP.NET MVC</a>, ASP.NET Web API application and share that across your repositories. Note: don&rsquo;t ever use singleton DbContext instance throughout your AppDomain. <a href="http://stackoverflow.com/questions/6126616/is-dbcontext-thread-safe">DbContext is not thread safe</a>.</p>
<p>As we have registered our EntitiesContext instance per request above, we can now register the repositories as well. As our repositories accepts an IEntitiesContext implementation through their constructor, our IoC container will use our previous registration for that automatically. Autofac has this ability as nearly all IoC containers do.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> RegisterDependencies(HttpConfiguration config) {

    <span style="color: blue;">var</span> builder = <span style="color: blue;">new</span> ContainerBuilder();
    builder.RegisterApiControllers(Assembly.GetExecutingAssembly());

    <span style="color: green;">// Register IEntitiesContext</span>
    builder.Register(_ =&gt; <span style="color: blue;">new</span> AccommodationEntities())
           .As&lt;IEntitiesContext&gt;().InstancePerApiRequest();

    <span style="color: green;">// TODO: Register repositories here</span>
    builder.RegisterType&lt;EntityRepository&lt;Country&gt;&gt;()
           .As&lt;IEntityRepository&lt;Country&gt;&gt;().InstancePerApiRequest();
    builder.RegisterType&lt;EntityRepository&lt;Resort&gt;&gt;()
           .As&lt;IEntityRepository&lt;Resort&gt;&gt;().InstancePerApiRequest();
    builder.RegisterType&lt;EntityRepository&lt;Hotel&gt;&gt;()
           .As&lt;IEntityRepository&lt;Hotel&gt;&gt;().InstancePerApiRequest();

    config.DependencyResolver = 
        <span style="color: blue;">new</span> AutofacWebApiDependencyResolver(builder.Build());
}</pre>
</div>
</div>
<h3>Out of the Box Pagination Support</h3>
<p>Best feature with this release is out of the box pagination support with generic repository instances. It doesn&rsquo;t perform the pagination in-memory; it queries the database accordingly and gets only the parts which are needed which is the whole point <img class="wlEmoticon wlEmoticon-smile" style="border-style: none;" alt="Smile" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/8007d0fd3e72_107A4/wlEmoticon-smile.png" /> Here is an ASP.NET Web API controller which uses the pagination support comes with the EntityRepository:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CountriesController : ApiController {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IEntityRepository&lt;Country&gt; _countryRepository;
    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IMappingEngine _mapper;
    <span style="color: blue;">public</span> CountriesController(
        IEntityRepository&lt;Country&gt; countryRepository, 
        IMappingEngine mapper) {

        _countryRepository = countryRepository;
        _mapper = mapper;
    }

    <span style="color: green;">// GET api/countries?pageindex=1&amp;pagesize=5</span>
    <span style="color: blue;">public</span> PaginatedDto&lt;CountryDto&gt; GetCountries(<span style="color: blue;">int</span> pageIndex, <span style="color: blue;">int</span> pageSize) {

        PaginatedList&lt;Country&gt; countries = 
             _countryRepository.Paginate(pageIndex, pageSize);

        PaginatedDto&lt;CountryDto&gt; countryPaginatedDto = 
           _mapper.Map&lt;PaginatedList&lt;Country&gt;, PaginatedDto&lt;CountryDto&gt;&gt;(countries);

        <span style="color: blue;">return</span> countryPaginatedDto;
    }
}

<span style="color: blue;">public</span> <span style="color: blue;">interface</span> IPaginatedDto&lt;<span style="color: blue;">out</span> TDto&gt; <span style="color: blue;">where</span> TDto : IDto {

    <span style="color: blue;">int</span> PageIndex { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">int</span> PageSize { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">int</span> TotalCount { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">int</span> TotalPageCount { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">bool</span> HasNextPage { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">bool</span> HasPreviousPage { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    IEnumerable&lt;TDto&gt; Items { <span style="color: blue;">get</span>; }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> PaginatedDto&lt;TDto&gt; : IPaginatedDto&lt;TDto&gt; <span style="color: blue;">where</span> TDto : IDto {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> PageIndex { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> PageSize { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> TotalCount { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">int</span> TotalPageCount { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">public</span> <span style="color: blue;">bool</span> HasNextPage { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">bool</span> HasPreviousPage { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">public</span> IEnumerable&lt;TDto&gt; Items { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>Paginate method will return us the PaginatedList&lt;TEntity&gt; object back and we can project that into our own Dto object as you can see above. I used AutoMapper for that. If I send a request to this API endpoint and ask for response in JSON format, I get back the below result:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>{
    <span style="color: #a31515;">"PageIndex"</span>:1,
    <span style="color: #a31515;">"PageSize"</span>:2,
    <span style="color: #a31515;">"TotalCount"</span>:6,
    <span style="color: #a31515;">"TotalPageCount"</span>:3,
    <span style="color: #a31515;">"HasNextPage"</span>:<span style="color: blue;">true</span>,
    <span style="color: #a31515;">"HasPreviousPage"</span>:<span style="color: blue;">false</span>,
    <span style="color: #a31515;">"Items"</span>:[
      {
        <span style="color: #a31515;">"Id"</span>:1,
        <span style="color: #a31515;">"Name"</span>:<span style="color: #a31515;">"Turkey"</span>,
        <span style="color: #a31515;">"ISOCode"</span>:<span style="color: #a31515;">"TR"</span>,
        <span style="color: #a31515;">"CreatedOn"</span>:<span style="color: #a31515;">"2013-01-08T21:12:26.5854461+02:00"</span>
      },
      {
        <span style="color: #a31515;">"Id"</span>:2,
        <span style="color: #a31515;">"Name"</span>:<span style="color: #a31515;">"United Kingdom"</span>,
        <span style="color: #a31515;">"ISOCode"</span>:<span style="color: #a31515;">"UK"</span>,
        <span style="color: #a31515;">"CreatedOn"</span>:<span style="color: #a31515;">"2013-01-08T21:12:26.5864465+02:00"</span>
      }
    ]
}</pre>
</div>
</div>
<p>Isn&rsquo;t this perfect <img class="wlEmoticon wlEmoticon-smile" style="border-style: none;" alt="Smile" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/8007d0fd3e72_107A4/wlEmoticon-smile.png" /> There are other pagination method inside the EntityRepository implementation which supports <a href="http://msdn.microsoft.com/en-us/library/gg671236%28v=vs.103%29.aspx">including child or parent entities</a> and sorting. You also have the <a href="https://github.com/tugberkugurlu/GenericRepository/blob/master/src/GenericRepository/QueryableExtensions.cs">ToPaginatedList</a> extension method and you can build your query and call ToPaginatedList on that query to get PaginatedList&lt;TEntity&gt; object back.</p>
<h3>Extending the IEntityRepository&lt;TEntity&gt;</h3>
<p>In my previous blog posts, I kind of sucked at extending the generic repository. So, I wanted to show here the better approach that I have been taking for a while now. This is not a feature of my generic repository, this is the feature of .NET itself: extension methods! If you need extra methods for your specific repository, you can always extend the IEntityRepository&lt;TEntity, TId&gt; which gives you a better way to extend your repositories. Here is an example:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">static</span> <span style="color: blue;">class</span> HotelRepositoryExtensions {

    <span style="color: blue;">public</span> <span style="color: blue;">static</span> IQueryable&lt;Hotel&gt; GetAllByResortId(
        <span style="color: blue;">this</span> IEntityRepository&lt;Hotel, <span style="color: blue;">int</span>&gt; hotelRepository, <span style="color: blue;">int</span> resortId) {

        <span style="color: blue;">return</span> hotelRepository.FindBy(x =&gt; x.ResortId == resortId);
    }
}</pre>
</div>
</div>
<h3>What is Next?</h3>
<p>My first intention is finish writing all the tests for the whole project, fix bugs and inconsistencies for the v0.3.0 release. After that release, I will work on EF6 version for my generic repository implementation which will have sweet asynchronous support. I also plan to release a generic repository implementation for MongoDB.</p>
<p>Stay tuned, <a href="http://nuget.org/packages/GenericRepository.EntityFramework">install the package</a>, play with it and <a href="https://github.com/tugberkugurlu/GenericRepository/issues">give feedback</a>&nbsp;<img class="wlEmoticon wlEmoticon-winkingsmile" style="border-style: none;" alt="Winking smile" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/8007d0fd3e72_107A4/wlEmoticon-winkingsmile.png" /></p>