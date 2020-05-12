---
id: 874d25d5-2f36-4cf1-4f8a-08d6d081bda9
title: Declarative Coding Approach to Caching in .NET Core with PostSharp
abstract: PostSharp is a .NET library which gives you ability to program in a declarative
  style and allows you perform many cross-cutting concerns with a minimum amount of
  code by abstracting away the complexity from you. In this post, I will be looking
  into how PostSharp helps us for caching to speed up the performance of our applications
  drastically.
created_at: 2019-05-04 11:14:56.9057917 +0000 UTC
tags:
- .NET Core
- Aspect Oriented Programming
- Caching
- PostSharp
slugs:
- declarative-coding-approach-to-caching-in--net-core-with-postsharp
---

<p>One of the first criteria of effective code is that it does its job with as few lines of code as possible. Effective code does not repeat itself. Less code in our codebases increases our chances of having less bugs. So, how do we avoid repeating ourselves? We apply our intelligence and abstraction skills to generalize behaviors into methods and classes, the constructs offered by C# to implement abstraction which we call encapsulation. However, some features such as logging or caching cannot be properly encapsulated into a class or method. That’s why you end up having code repetition. C# alone is simply not able to properly encapsulate features like logging, caching, security, INotifyPropertyChanged, undo/redo, etc.</p><p>I have been meaning to look into <a href="https://en.wikipedia.org/wiki/Aspect-oriented_programming" target="_blank">Aspect-oriented programming</a> for a while to help my code to be less noisy without sacrificing the application's acceptable performance and observability. This would help cut right to the business logic, allowing me to care about what's more important. When the topic is Aspect-oriented programming, first software comes to my mind is obviously <a href="https://www.postsharp.net/?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">PostSharp</a> in .NET world and in this post, I will be looking at how PostSharp can help us cut the noise out of our code and showcase this with a sample on data caching.</p><h3>Getting Started with PostSharp</h3><p>First of all, let's create our project structure and install PostSharp. I have .NET Core SDK 2.2.202 installed and ran the below commands to create the empty project structure.</p>

<p>
</p><pre>dotnet new web --no-https
dotnet new sln
dotnet sln 1-sample-web.sln add 1-sample-web.csproj
dotnet new globaljson
</pre>
<p></p>

<p>In order to give you an idea about the value proposition of PostSharp, I created this little ASP.NET Core sample which exposes HTTP APIs to read, write and modify the Cars in our system. Some of the code here is contrived such as sleeping for half a second, etc. but we will see why this will be useful for us to see the PostSharp in action.</p>

<p>
</p><pre>using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.DependencyInjection;

namespace _1_sample_web
{
    public class Startup
    {
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddMvc();
        }

        public void Configure(IApplicationBuilder app, IHostingEnvironment env)
        {
            app.UseMvcWithDefaultRoute();
        }
    }

    public class CarsController : Controller
    {
        private static readonly CarsContext _carsCtx = new CarsContext();

        [HttpGet("cars")]
        public IEnumerable<car> Get()
        {
            return _carsCtx.GetAll();
        }

        [HttpGet("cars/{id}")]
        public IActionResult GetCar(int id) 
        {
            var carTuple = _carsCtx.GetSingle(id);
            if (!carTuple.Item1) 
            {
                return NotFound();
            }

            return Ok(carTuple.Item2);
        }

        [HttpPost("cars/{id}")]
        public IActionResult PostCar(Car car) 
        {
            var createdCar = _carsCtx.Add(car);
            return CreatedAtAction(nameof(GetCar), 
                new { id = createdCar.Id }, 
                createdCar);
        }

        [HttpPut("cars/{id}")]
        public IActionResult PutCar(int id, Car car) 
        {
            car.Id = id;
            if (!_carsCtx.TryUpdate(car)) 
            {
                return NotFound();
            }

            return Ok(car);
        }

        [HttpDelete("cars/{id}")]
        public IActionResult DeleteCar(int id) 
        {
            if (!_carsCtx.TryRemove(id)) 
            {
                return NotFound();
            }

            return NoContent();
        }
    }

    public class Car 
    {
        public int Id { get; set; }

        [Required]
        [StringLength(20)]
        public string Make { get; set; }

        [Required]
        [StringLength(20)]
        public string Model { get; set; }

        public int Year { get; set; }

        [Range(0, 500000)]
        public float Price { get; set; }
    }

    public class CarsContext
    {
        private int _nextId = 9;
        private object _idLock = new object();

        private readonly ConcurrentDictionary<int, car=""> _database = new ConcurrentDictionary<int, car="">(new HashSet<keyvaluepair<int, car="">&gt; 
        { 
            new KeyValuePair<int, car="">(1, new Car { Id = 1, Make = "Make1", Model = "Model1", Year = 2010, Price = 10732.2F }),
            new KeyValuePair<int, car="">(2, new Car { Id = 2, Make = "Make2", Model = "Model2", Year = 2008, Price = 27233.1F }),
            new KeyValuePair<int, car="">(3, new Car { Id = 3, Make = "Make3", Model = "Model1", Year = 2009, Price = 67437.0F }),
            new KeyValuePair<int, car="">(4, new Car { Id = 4, Make = "Make4", Model = "Model3", Year = 2007, Price = 78984.2F }),
            new KeyValuePair<int, car="">(5, new Car { Id = 5, Make = "Make5", Model = "Model1", Year = 1987, Price = 56200.89F }),
            new KeyValuePair<int, car="">(6, new Car { Id = 6, Make = "Make6", Model = "Model4", Year = 1997, Price = 46003.2F }),
            new KeyValuePair<int, car="">(7, new Car { Id = 7, Make = "Make7", Model = "Model5", Year = 2001, Price = 78355.92F }),
            new KeyValuePair<int, car="">(8, new Car { Id = 8, Make = "Make8", Model = "Model1", Year = 2011, Price = 1823223.23F })
        });
        
        public IEnumerable<car> GetAll()
        {
            Thread.Sleep(500);
            return _database.Values;
        }

        public IEnumerable<car> Get(Func<car, bool=""> predicate) 
        {
            Thread.Sleep(500);
            return _database.Values.Where(predicate);
        }

        public Tuple<bool, car=""> GetSingle(int id) 
        {
            Thread.Sleep(500);

            Car car;
            var doesExist = _database.TryGetValue(id, out car);
            return new Tuple<bool, car="">(doesExist, car);
        }

        public Car GetSingle(Func<car, bool=""> predicate) 
        {
            Thread.Sleep(500);
            return _database.Values.FirstOrDefault(predicate);
        }

        public Car Add(Car car) 
        {
            Thread.Sleep(500);
            lock(_idLock) 
            {
                car.Id = _nextId;
                _database.TryAdd(car.Id, car);
                _nextId++;
            }

            return car;
        }

        public bool TryRemove(int id) 
        {
            Thread.Sleep(500);

            Car removedCar;
            return _database.TryRemove(id, out removedCar);
        }

        public bool TryUpdate(Car car) 
        {
            Thread.Sleep(500);

            Car oldCar;
            if (_database.TryGetValue(car.Id, out oldCar)) {

                return _database.TryUpdate(car.Id, car, oldCar);
            }

            return false;
        }
    }
}<font color="#000000" face="-apple-system, system-ui, Segoe UI, Roboto, Helvetica Neue, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji"><span style="font-size: 16px; white-space: normal;">
</span></font></car,></bool,></bool,></car,></car></car></int,></int,></int,></int,></int,></int,></int,></int,></keyvaluepair<int,></int,></int,></car></pre>
<p></p>

<p>Before going further, let's install PostSharp through NuGet. The first thing you want to install is <a href="https://www.nuget.org/packages/PostSharp" target="_blank">PostSharp NuGet package</a> which magically hooks into the compilation step thanks to its custom MSBuild scripts. The other package here will be <a href="https://www.nuget.org/packages/PostSharp.Patterns.Diagnostics" target="_blank">PostSharp.Patterns.Diagnostics</a> as I want to show you a logging example first.</p><p>
</p><pre>dotnet add package PostSharp
dotnet add package PostSharp.Patterns.Diagnostics
</pre>
<p></p><p>Let's get the sample code from the <a href="https://doc.postsharp.net/add-logging?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">logging documentation</a>.</p><p>
</p><pre>using PostSharp.Patterns.Diagnostics;
using PostSharp.Extensibility;

[assembly: Log(AttributePriority = 1, AttributeTargetMemberAttributes = MulticastAttributes.Protected | MulticastAttributes.Internal | MulticastAttributes.Public)]
[assembly: Log(AttributePriority = 2, AttributeExclude = true, AttributeTargetMembers = "get_*" )]<font color="#000000" face="-apple-system, system-ui, Segoe UI, Roboto, Helvetica Neue, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji"><span style="font-size: 16px; white-space: normal;">
</span></font></pre>
<p></p><p>When you run the application now, you will be impressed and probably also be blown away by how much value and observability you get with a very little work!</p><p><img src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/20190504104919-266c1533-a4c2-4e39-bb68-b03d7a67c788-image.png"><br></p><p><img src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/20190504105012-f2342e91-3912-41af-8491-e8576b9c7b6c-image.png"><br></p><h3>PostSharp Caching Example</h3><p>The main reason for me to explore PostSharp is for caching and this is where <a href="https://www.postsharp.net/caching?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">PostSharp Caching</a> shines really. Let's run our sample application again and perform a mini load test on it.</p><p>
</p><pre>1..10 | foreach {write-host "$([Math]::Round((Measure-Command -Expression { Invoke-WebRequest -Uri http://localhost:5000/cars }).TotalMilliseconds, 1))"}<font color="#000000" face="-apple-system, system-ui, Segoe UI, Roboto, Helvetica Neue, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji"><span style="font-size: 16px; white-space: normal;">
</span></font></pre>
<p></p><p>You will notice that each call to the "/cars" endpoint takes more than 500ms, which is fair due to us sleeping that amount of time on purpose. However, this could well be the case when you connect to a data store in a real world example. Even if your data store is performant and gets the result instantly, we are still wasting resources here because the data hasn't changed and we would be doing an unnecessary trip to the database to get the data which we have already retrieved previously.</p><p>Caching is the solution to this problem. However, it's not really easy to get right on your own in a web application which is multithreaded in its nature. You can use built-in APIs such as <a href="https://docs.microsoft.com/en-us/aspnet/core/performance/caching/memory?view=aspnetcore-2.2" target="_blank">the ones come from ASP.NET Core</a> but you then need to express your caching requirements in code, in a verbose way which will make it hard to understand the business logic behind a cluttered codebase and suddenly, you will be struggling to add or modify functionality in an existing software.</p><p>Let's see how PostSharp can help us here. First, we need to add the caching support by installing <a href="https://www.nuget.org/packages/PostSharp.Patterns.Caching" target="_blank">PostSharp.Patterns.Caching</a> NuGet package.</p><p>
</p><pre>dotnet add package PostSharp.Patterns.Caching<font color="#000000" face="-apple-system, system-ui, Segoe UI, Roboto, Helvetica Neue, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji"><span style="font-size: 16px; white-space: normal;">
</span></font></pre>
<p></p><p>Then, we need to make some changes to our code to enable caching. Here is the git patch which shows you what exactly I have changed:<br></p><p>
</p><pre>From a20fc8e95ffd9bf5d424467e0e1283ae5891454a Mon Sep 17 00:00:00 2001
From: Tugberk Ugurlu
Date: Tue, 9 Apr 2019 23:38:32 +0100
Subject: [PATCH] add caching

---
 postsharp/0-caching/1-sample-web/1-sample-web.csproj | 1 +
 postsharp/0-caching/1-sample-web/Program.cs          | 3 +++
 postsharp/0-caching/1-sample-web/Startup.cs          | 4 +++-
 3 files changed, 7 insertions(+), 1 deletion(-)

diff --git a/postsharp/0-caching/1-sample-web/1-sample-web.csproj b/postsharp/0-caching/1-sample-web/1-sample-web.csproj
index bd55b6c..008c486 100644
--- a/postsharp/0-caching/1-sample-web/1-sample-web.csproj
+++ b/postsharp/0-caching/1-sample-web/1-sample-web.csproj
@@ -10,6 +10,7 @@
     <packagereference include="Microsoft.AspNetCore.App"></packagereference>
     <packagereference include="Microsoft.AspNetCore.Razor.Design" version="2.2.0" privateassets="All"></packagereference>
     <packagereference include="PostSharp" version="6.1.17"></packagereference>
+    <packagereference include="PostSharp.Patterns.Caching" version="6.1.17"></packagereference>
     <packagereference include="PostSharp.Patterns.Diagnostics" version="6.1.17"></packagereference>
   
 
diff --git a/postsharp/0-caching/1-sample-web/Program.cs b/postsharp/0-caching/1-sample-web/Program.cs
index 3dcae2c..9d241eb 100644
--- a/postsharp/0-caching/1-sample-web/Program.cs
+++ b/postsharp/0-caching/1-sample-web/Program.cs
@@ -7,6 +7,8 @@ using Microsoft.AspNetCore;
 using Microsoft.AspNetCore.Hosting;
 using Microsoft.Extensions.Configuration;
 using Microsoft.Extensions.Logging;
+using PostSharp.Patterns.Caching;
+using PostSharp.Patterns.Caching.Backends;
 using PostSharp.Patterns.Diagnostics;
 using PostSharp.Patterns.Diagnostics.Backends.Console;
 
@@ -18,6 +20,7 @@ namespace _1_sample_web
         public static void Main(string[] args)
         {
             LoggingServices.DefaultBackend = new ConsoleLoggingBackend();
+            CachingServices.DefaultBackend = new MemoryCachingBackend();
             CreateWebHostBuilder(args).Build().Run();
         }
 
diff --git a/postsharp/0-caching/1-sample-web/Startup.cs b/postsharp/0-caching/1-sample-web/Startup.cs
index 18b3dbc..bed37ca 100644
--- a/postsharp/0-caching/1-sample-web/Startup.cs
+++ b/postsharp/0-caching/1-sample-web/Startup.cs
@@ -10,6 +10,7 @@ using Microsoft.AspNetCore.Hosting;
 using Microsoft.AspNetCore.Http;
 using Microsoft.AspNetCore.Mvc;
 using Microsoft.Extensions.DependencyInjection;
+using PostSharp.Patterns.Caching;
 
 namespace _1_sample_web
 {
@@ -115,7 +116,8 @@ namespace _1_sample_web
             new KeyValuePair<int, car="">(7, new Car { Id = 7, Make = "Make7", Model = "Model5", Year = 2001, Price = 78355.92F }),
             new KeyValuePair<int, car="">(8, new Car { Id = 8, Make = "Make8", Model = "Model1", Year = 2011, Price = 1823223.23F })
         });

+        [Cache]
         public IEnumerable<car> GetAll()
         {
             Thread.Sleep(500);
-- 
2.15.2 (Apple Git-101.1)<font color="#000000" face="-apple-system, system-ui, Segoe UI, Roboto, Helvetica Neue, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji"><span style="font-size: 16px; white-space: normal;">
</span></font></car></int,></int,></pre>
<p></p><p>Couple of things we have done here:<br></p><ul><li>In our entry point, we configured the cache backend we wanted to use which in our case is the <a href="https://doc.postsharp.net/caching-memory?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">MemoryCache</a>.</li><li>We marked the CarContext.GetAll method with the <a href="https://doc.postsharp.net/t_postsharp_patterns_caching_cacheattribute?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">CacheAttribute</a>.<br></li></ul><p>Believe it or not, this is pretty much it! When we run the sample mini load test, you will see the dramatic difference even if we are seeing a higher response time on the first load.<br></p><p><img src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/20190504110040-34404b0c-2b63-430f-921c-6eae974388d5-image.png"><br></p><p>Again, very little work but tremendous gain in terms of value!</p><p>We have improved our performance drastically but introduced a very nasty problem now: serving stale data. Thankfully, <a href="https://doc.postsharp.net/caching-invalidation?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">PostSharp has a solution to cache invalidation</a> out of the box without losing our declarative nature for simple cases. For this, we need to use <a href="https://doc.postsharp.net/t_postsharp_patterns_caching_invalidatecacheattribute?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">InvalidateCacheAttribute</a> aspect. When this attribute is applied to a method, it causes any call to this method to remove from the cache the value of one or more other methods. It’s worth noting that the cached methods are matched, by type and name, against the parameters of the invalidating method. PostSharp compilation takes care of the rest during the build step to set up all the invalidation logic.</p><p>For example, the below changes makes it possible for us to invalidate the cache of a single car entity for example when it’s updated.</p><p>
</p><pre>From f0889e68e55298e43360e01dd3b0e8b1cf6468e3 Mon Sep 17 00:00:00 2001
From: Tugberk Ugurlu
Date: Tue, 30 Apr 2019 09:40:21 +0100
Subject: [PATCH] cache invalidation, declarative

---
 postsharp/0-caching/1-sample-web/Startup.cs | 6 ++++--
 1 file changed, 4 insertions(+), 2 deletions(-)

diff --git a/postsharp/0-caching/1-sample-web/Startup.cs b/postsharp/0-caching/1-sample-web/Startup.cs
index bed37ca..ec95d1e 100644
--- a/postsharp/0-caching/1-sample-web/Startup.cs
+++ b/postsharp/0-caching/1-sample-web/Startup.cs
@@ -62,7 +62,7 @@ namespace _1_sample_web
         public IActionResult PutCar(int id, Car car) 
         {
             car.Id = id;
-            if (!_carsCtx.TryUpdate(car)) 
+            if (!_carsCtx.TryUpdate(id, car)) 
             {
                 return NotFound();
             }
@@ -130,6 +130,7 @@ namespace _1_sample_web
             return _database.Values.Where(predicate);
         }
 
+        [Cache]
         public Tuple<bool, car=""> GetSingle(int id) 
         {
             Thread.Sleep(500);
@@ -166,7 +167,8 @@ namespace _1_sample_web
             return _database.TryRemove(id, out removedCar);
         }
 
-        public bool TryUpdate(Car car) 
+        [InvalidateCache(nameof(GetSingle))]
+        public bool TryUpdate(int id, Car car) 
         {
             Thread.Sleep(500);
 
-- 
2.20.1 (Apple Git-117)<font color="#000000" face="-apple-system, system-ui, Segoe UI, Roboto, Helvetica Neue, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji"><span style="font-size: 16px; white-space: normal;">
</span></font></bool,></pre>
<p></p><p>However, this only invalidates the GetSingle method and we still have problem of serving stale data from GetAll method. There is also an ability out of the box to to <a href="https://doc.postsharp.net/caching-invalidation#ID2RBSection?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">imperatively invalidate an item from the cache</a> which is very handy for cases where we cannot simply invalidate the cache purely based on method signature. You can see below an example of how this looks like.&nbsp;<br></p><p>
</p><pre>From f629b295fc8f9bbd44904284cb0ec832d51185be Mon Sep 17 00:00:00 2001
From: Tugberk Ugurlu
Date: Tue, 30 Apr 2019 09:55:44 +0100
Subject: [PATCH] cache invalidation, imperatively

---
 postsharp/0-caching/1-sample-web/Startup.cs | 4 ++++
 1 file changed, 4 insertions(+)

diff --git a/postsharp/0-caching/1-sample-web/Startup.cs b/postsharp/0-caching/1-sample-web/Startup.cs
index ec95d1e..8ee6652 100644
--- a/postsharp/0-caching/1-sample-web/Startup.cs
+++ b/postsharp/0-caching/1-sample-web/Startup.cs
@@ -67,6 +67,10 @@ namespace _1_sample_web
                 return NotFound();
             }
 
+            CachingServices.Invalidation.Invalidate(
+                typeof(CarsContext).GetMethod(nameof(CarsContext.GetAll)), 
+                _carsCtx);
+                
             return Ok(car);
         }
 
-- 
2.20.1 (Apple Git-117)<font color="#000000" face="-apple-system, system-ui, Segoe UI, Roboto, Helvetica Neue, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji"><span style="font-size: 16px; white-space: normal;">
</span></font></pre>
<p></p><p>We Invalidate the GetAll method cache on the given CarsContext instance when we have an update on any of the items.<br></p><p>This is all I want to cover on this post in terms of the API surface area of PostSharp and I hope this gives you taste of how simple it’s to get going with PostSharp. <a href="https://doc.postsharp.net/caching?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">PostSharp Caching documentation</a> is also very comprehensive and I recommend you to check that out for further details.<br></p><h3>Limitations</h3><p>The biggest limitation I have seen with PostSharp is its lack of .NET Core compilation support outside of Windows at the time of writing (you may check the current status <a href="https://support.postsharp.net/request/20871-support-for-coreclr-as-a-build" target="_blank">here</a>). You can run PostSharp on .NET Core, even outside of Windows. However, you first need a Windows machine to be able to compile your code.</p><p>Apart from this, there is also a trade off for you to make with PostSharp which is the increased build time. However, with incremental builds, this additional increase can become noticeable. Besides this, compared to the value you got from the tool, I think this is trade-off which is well worth to be made.</p><h3>Conclusion</h3><p>This post just touches the surface on what you can achieve with PostSharp. In terms of caching for example, there is even a support for <a href="https://doc.postsharp.net/caching-redis?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">Redis</a> which is very suitable for horizontally scaled web applications where multiple nodes serve HTTP requests.</p><p>PostSharp provides help on many other various patterns such as <a href="https://doc.postsharp.net/threading-library?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">mutithreading</a>. You can get started with PostSharp with <a href="https://www.postsharp.net/essentials?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=3_2019" target="_blank">PostSharp Essentials, the free but project-size-limited edition</a>.<br></p>