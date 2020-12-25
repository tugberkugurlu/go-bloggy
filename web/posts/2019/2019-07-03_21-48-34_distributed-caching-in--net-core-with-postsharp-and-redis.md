---
id: c6d7befc-7404-4d92-b8c7-08d7000032b9
title: Distributed Caching in .NET Core with PostSharp and Redis
abstract: On my previous post, I walked through the benefits of using PostSharp for
  caching in a .NET Core server application, by making it work on a single node application.
  In this post, we will see how we can enable Redis as the caching backend through
  PostSharp's modular nature.
created_at: 2019-07-03 21:48:34.6800844 +0000 UTC
tags:
- .NET Core
- Aspect Oriented Programming
- Caching
- PostSharp
- Redis
slugs:
- distributed-caching-in--net-core-with-postsharp-and-redis
---

<p><a href="https://www.tugberkugurlu.com/archive/declarative-coding-approach-to-caching-in--net-core-with-postsharp" target="_blank">On my previous post</a>, I walked through the benefits of using <a href="https://www.postsharp.net/?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=06_2019">PostSharp</a> for caching in a .NET Core server application. However, the example I have showed there would work on a single node application but as we know, probably no application today works on a single node. The benefits of deploying into multiple nodes are multiple such as providing further fault tolerance, and load distribution.</p><p>Luckily for us, <a href="https://doc.postsharp.net/caching?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=06_2019" target="_blank">PostSharp caching</a> backend is modular and the default in-memory one I have used in my previous post can be swapped. One of the out of the box implementations is based on <a href="https://doc.postsharp.net/caching-redis/?utm_source=blog&amp;utm_medium=tugberk&amp;utm_campaign=06_2019" target="_blank">Redis</a>, which is a highly scalable, distributed data structure server solution. One of the widely use cases of Redis is to be used as a ephemeral key/value store to power the caching needs of the apps.</p><h3>Run Redis Locally</h3><p>The best way to <a href="https://hub.docker.com/_/redis/" target="_blank">run Redis locally is through Docker</a>. Let’s run the below code to do this:</p>
<p>
</p><pre>docker run --name postsharp-redis -p 6379:6379 -d redis a30f1c1e991e0159fb5f96dfb053f50c50726101907c7f76d319d5e987a6cf3a
</pre>
<p></p>
<p>We have just got a redis instance up and running on our local environment and exposed it through TCP port mapping to the host machine to be available at port 6379. The final thing we need to do to get this ready for PostSharp usage is to set up the key-space notification to include the AKE events. You can see <a href="https://redis.io/topics/notifications#configuration" target="_blank">the Redis notifications document</a> for details on this.</p><h3>Configure for Redis Cache</h3><p>First thing to do is to install the NuGet package which contains the Redis&nbsp; caching backend implementation for PostSharp.</p><p>
</p><pre>dotnet add package PostSharp.Patterns.Caching.Redis --version 6.2.8
</pre>
<p></p><p>Then, all we need to do is to change the caching backend to be the Redis implementation, which we have configured inside our Program.Main method in the previous post:</p><p>
</p><pre>string connectionConfiguration = "127.0.0.1";
var connection = ConnectionMultiplexer.Connect(connectionConfiguration);
var redisCachingConfiguration = new RedisCachingBackendConfiguration();
CachingServices.DefaultBackend = RedisCachingBackend.Create(connection, redisCachingConfiguration);
</pre>
<p></p><p>Notice the server address we have entered, that points to the Redis instance we have got up and running through Docker and exposed to the host through port mapping. As we used the default Redis port, we didn’t need to state it explicitly.</p><p>From this point forward, our app is all ready to run with Redis caching enabled, without a single line of code change on the app components. Only change we had to do was on the configuration side.</p><p>For production, it’s worth getting a hold of the Redis server address through <a href="https://docs.microsoft.com/en-us/aspnet/core/fundamentals/configuration/?view=aspnetcore-2.2" target="_blank">a configuration system such as the one provided with ASP.NET Core</a> so that you can swap it based on your environment.</p>