---
id: ae0648d8-929f-4ddf-b2de-5564451a3754
title: ASP.NET Core Authentication in a Load Balanced Environment with HAProxy and
  Redis
abstract: Token based authentication is a fairly common way of authenticating a user
  for an HTTP application. However, handling this in a load balanced environment has
  always involved extra caring. In this post, I will show you how this is handled
  in ASP.NET Core by demonstrating it with HAProxy and Redis through the help of Docker.
created_at: 2016-11-28 23:31:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Core
- Docker
- HTTP
- Security
slugs:
- asp-net-core-authentication-in-a-load-balanced-environment-with-haproxy-and-redis
- asp-net-core-authentication-in-a-load-balancer-environment-with-haproxy-and-redis
---

<p>Token based authentication is a fairly common way of authenticating a user for an HTTP application. ASP.NET and its frameworks had support for implementing this out of the box without much effort with different type of authentication approaches such as cookie based authentication, bearer token authentication, etc. ASP.NET Core is a no exception to this and it got even better (which we will see in a while).</p> <p>However, handling this in a load balanced environment has always involved extra caring as all of the nodes should be able to read the valid authentication token even if that token has been written by another node. Old-school ASP.NET solution to this is to keep the <a href="https://msdn.microsoft.com/en-us/library/w8h3skw9(v=vs.100).aspx">Machine Key</a> in sync with all the nodes. Machine key, for those who are not familiar with it, is used to encrypt and decrypt the authentication tokens under ASP.NET and each machine by default has its own unique one. However, you can override this and put your own one in place per application through a setting inside the Web.config file. This approach had <a href="https://twitter.com/blowdart/status/796821271853965312">its own problems</a> and <a href="https://docs.microsoft.com/en-us/aspnet/core/security/data-protection/introduction">with ASP.NET Core, all data protection APIs have been revamped</a> which cleared a room for big improvements in this area such as <a href="https://docs.microsoft.com/en-us/aspnet/core/security/data-protection/implementation/key-management#data-protection-implementation-key-management-expiration">key expiration and rolling</a>, <a href="https://docs.microsoft.com/en-us/aspnet/core/security/data-protection/implementation/key-encryption-at-rest">key encryption at rest</a>, etc. One of those improvements is the ability to <a href="https://docs.microsoft.com/en-us/aspnet/core/security/data-protection/implementation/key-storage-providers">store keys in different storage systems</a>, which is what I am going to touch on in this post.</p> <h3>The Problem</h3> <p>Imagine a case where we have an ASP.NET Core application which uses cookie based authentication and stores their user data in MongoDB, which has been implemented using <a href="https://github.com/aspnet/Identity/">ASP.NET Core Identity</a> and <a href="https://github.com/tugberkugurlu/AspNetCore.Identity.MongoDB">its MongoDB provider</a>.</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/620c8c35-02a0-4601-95b0-a6e5fcd03cb1.jpg"><img title="1-ok-wo-lb" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="1-ok-wo-lb" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a65f7f62-5ccd-4502-9d3a-9c84b8f10b1a.jpg" width="644" height="224"></a></p> <p>This setup is all fine and our application should function perfectly. However, if we put this application behind HAProxy and scale it up to two nodes, we will start seeing problems like below:</p><pre>System.Security.Cryptography.CryptographicException: The key {3470d9c3-e59d-4cd8-8668-56ba709e759d} was not found in the key ring.
   at Microsoft.AspNetCore.DataProtection.KeyManagement.KeyRingBasedDataProtector.UnprotectCore(Byte[] protectedData, Boolean allowOperationsOnRevokedKeys, UnprotectStatus&amp; status)
   at Microsoft.AspNetCore.DataProtection.KeyManagement.KeyRingBasedDataProtector.DangerousUnprotect(Byte[] protectedData, Boolean ignoreRevocationErrors, Boolean&amp; requiresMigration, Boolean&amp; wasRevoked)
   at Microsoft.AspNetCore.DataProtection.KeyManagement.KeyRingBasedDataProtector.Unprotect(Byte[] protectedData)
   at Microsoft.AspNetCore.Antiforgery.Internal.DefaultAntiforgeryTokenSerializer.Deserialize(String serializedToken)</pre>
<p>Let’s look at the below diagram to understand why we are having this problem: 
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/3e7d151a-dd5d-4521-a049-c253a4589ca6.jpg"><img title="2-not-ok-w-lb" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="2-not-ok-w-lb" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e6eb5b9c-84b2-4abf-a49d-2eccb027768a.jpg" width="644" height="252"></a> 
<p>By default, ASP.NET Core Data Protection is wired up to store its keys under the file system. If you have your application running under multiple nodes as shown in above diagram, each node will have its own keys to protect and unprotect the sensitive information like authentication cookie data. As you can guess, this behaviour is problematic with the above structure since one node cannot read the protected data which the other node protected. 
<h3>The Solution</h3>
<p>As I mentioned before, one of the extensibility points of ASP.NET Core Data Protection stack is the storage of the data protection keys. This place can be a central place where all the nodes of our web application can reach out to. Let’s look at the below diagram to understand what we mean by this:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a20192f3-e690-4d12-bdb7-13efa1babb80.jpg"><img title="4-ok-w-lb" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="4-ok-w-lb" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/11d9bcd2-53ea-483e-b712-a5a3eb63c011.jpg" width="644" height="259"></a></p>
<p>Here, we have <a href="https://redis.io/">Redis</a> as our Data Protection key storage. Redis is a good choice here as it’s a well-suited for key-value storage and that’s what we need. With this setup, it will be possible for both nodes of our application to read protected data regardless of which node has written it.</p>
<h3>Wiring up Redis Data Protection Key Storage</h3>
<p>With ASP.NET Core 1.0.0, we had to write the implementation by ourselves to make ASP.NET Core to store Data Protection keys on Redis but with 1.1.0 release, the team has simultaneously shipped a NuGet package which makes it really easy to wire this up: <a href="https://www.nuget.org/packages/Microsoft.AspNetCore.DataProtection.Redis/0.1.0">Microsoft.AspNetCore.DataProtection.Redis</a>. This package easily allows us to swap the data protection storage destination to be Redis. We can do this while we are configuring services as part of ConfigureServices:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">void</span> ConfigureServices(IServiceCollection services)
{
    <span style="color: green">// sad but a giant hack :(</span>
    <span style="color: green">// https://github.com/StackExchange/StackExchange.Redis/issues/410#issuecomment-220829614</span>
    <span style="color: blue">var</span> redisHost = Configuration.GetValue&lt;<span style="color: blue">string</span>&gt;(<span style="color: #a31515">"Redis:Host"</span>);
    <span style="color: blue">var</span> redisPort = Configuration.GetValue&lt;<span style="color: blue">int</span>&gt;(<span style="color: #a31515">"Redis:Port"</span>);
    <span style="color: blue">var</span> redisIpAddress = Dns.GetHostEntryAsync(redisHost).Result.AddressList.Last();
    <span style="color: blue">var</span> redis = ConnectionMultiplexer.Connect($<span style="color: #a31515">"{redisIpAddress}:{redisPort}"</span>);

    services.AddDataProtection().PersistKeysToRedis(redis, <span style="color: #a31515">"DataProtection-Keys"</span>);
    services.AddOptions();

    <span style="color: green">// ...</span>
}</pre></div></div>
<p><a href="https://github.com/tugberkugurlu/AspNetCoreSamples/blob/haproxy-redis-auth/haproxy-redis-auth/src/Startup.cs#L52-L62">I have wired it up exactly like this</a> in <a href="https://github.com/tugberkugurlu/AspNetCoreSamples/blob/haproxy-redis-auth/haproxy-redis-auth">my sample application</a> in order to show you a working example. It’s an example taken from ASP.NET Identity repository but slightly changed to make it work with <a href="https://github.com/tugberkugurlu/AspNetCore.Identity.MongoDB">MongoDB Identity store provider</a>.</p>
<blockquote>
<p>Note here that configuration values above are specific to my implementation and it doesn’t have to be that way. <a href="https://github.com/tugberkugurlu/AspNetCoreSamples/blob/haproxy-redis-auth/haproxy-redis-auth/docker-compose.yml#L23-L24">See these lines inside my Docker Compose file</a> and <a href="https://github.com/tugberkugurlu/AspNetCoreSamples/blob/haproxy-redis-auth/haproxy-redis-auth/src/Startup.cs#L36-L43">these inside my Startup class</a> to understand how it’s being passed and hooked up.</p></blockquote>
<p>The sample application can be run on <a href="https://www.docker.com/">Docker</a> through <a href="https://docs.docker.com/compose/">Docker Compose</a> and it will get a few things up and running: 
<ul>
<li>Two nodes of the application 
<li>A MongoDB instance 
<li>A Redis instance</li></ul>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1c5ebc5e-adc8-46da-9344-d5d1202002ac.png"><img title="Image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1dd326e1-35c8-42e5-9b80-fa79e3576fb4.png" width="640" height="360"></a></p>
<p>You can see <a href="https://github.com/tugberkugurlu/AspNetCoreSamples/blob/haproxy-redis-auth/haproxy-redis-auth/docker-compose.yml">my docker-compose.yml file</a> to understand how I hooked things together:</p><pre>mongo:
    build: .
    dockerfile: mongo.dockerfile
    container_name: haproxy_redis_auth_mongodb
    ports:
      - "27017:27017"

redis:
    build: .
    dockerfile: redis.dockerfile
    container_name: haproxy_redis_auth_redis
    ports:
      - "6379:6379"

webapp1:
    build: .
    dockerfile: app.dockerfile
    container_name: haproxy_redis_auth_webapp1
    environment:
      - ASPNETCORE_ENVIRONMENT=Development
      - ASPNETCORE_server.urls=http://0.0.0.0:6000
      - WebApp_MongoDb__ConnectionString=mongodb://mongo:27017
      - WebApp_Redis__Host=redis
      - WebApp_Redis__Port=6379
    links:
      - mongo
      - redis

webapp2:
    build: .
    dockerfile: app.dockerfile
    container_name: haproxy_redis_auth_webapp2
    environment:
      - ASPNETCORE_ENVIRONMENT=Development
      - ASPNETCORE_server.urls=http://0.0.0.0:6000
      - WebApp_MongoDb__ConnectionString=mongodb://mongo:27017
      - WebApp_Redis__Host=redis
      - WebApp_Redis__Port=6379
    links:
      - mongo
      - redis

app_lb:
    build: .
    dockerfile: haproxy.dockerfile
    container_name: app_lb
    ports:
      - "5000:80"
    links:
      - webapp1
      - webapp2</pre>
<p>HAProxy is also configured to balance the load between two application nodes as you can see inside <a href="https://github.com/tugberkugurlu/AspNetCoreSamples/blob/haproxy-redis-auth/haproxy-redis-auth/haproxy.cfg">the haproxy.cfg file</a>, <a href="https://github.com/tugberkugurlu/AspNetCoreSamples/blob/haproxy-redis-auth/haproxy-redis-auth/haproxy.dockerfile#L2">which we copy under the relevant path inside our dockerfile</a>:</p><pre>global
  log 127.0.0.1 local0
  log 127.0.0.1 local1 notice

defaults
  log global
  mode http
  option httplog
  option dontlognull
  timeout connect 5000
  timeout client 10000
  timeout server 10000

frontend balancer
  bind 0.0.0.0:80
  mode http
  default_backend app_nodes

backend app_nodes
  mode http
  balance roundrobin
  option forwardfor
  http-request set-header X-Forwarded-Port %[dst_port]
  http-request set-header Connection keep-alive
  http-request add-header X-Forwarded-Proto https if { ssl_fc }
  option httpchk GET / HTTP/1.1\r\nHost:localhost
  server webapp1 webapp1:6000 check
  server webapp2 webapp2:6000 check</pre>
<p>All of these are some details on how I wired up the sample to work. If we now look closely at the header of the web page, you should see the server name written inside the parenthesis. If you refresh enough, you will see that part alternating between two server names:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/9ae005d8-6e6e-4fa6-bdb5-3738a1aea104.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/68afaadf-0a64-43af-9698-55a4b21fdab9.png" width="644" height="106"></a></p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/087be567-c5bd-4198-bda8-cbe5312d1dff.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e3db0ca9-d287-4bb3-8cdb-e23d8eab89e4.png" width="644" height="100"></a></p>
<p>This confirms that our load is balanced between the two application nodes. The rest of the demo is actually very boring. It should just work as you expect it to work. Go to “Register” page and register for an account, log out and log back in. All of those interactions should just work. If we look inside the Redis instance, we should also see that Data Protection key has been written there:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>docker run <span style="color: gray">-</span>it <span style="color: gray">--</span>link haproxy_redis_auth_redis:redis <span style="color: gray">--</span>rm redis redis<span style="color: gray">-</span>cli <span style="color: gray">-</span>h redis <span style="color: gray">-</span>p 6379
LRANGE DataProtection<span style="color: gray">-</span>Keys 0 10</pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ceda7940-9916-4a5a-83a4-23b0d6b946a2.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/107d9136-3858-404b-951a-c31564a29fc4.png" width="640" height="312"></a></p>
<h3>Conclusion and Going Further</h3>
<p>I believe that I was able to show you what you need to care about in terms of authentication when you scale our your application nodes to multiple servers. However, do not take my sample as is and apply to your production application :) There are a few important things that suck on my sample, like the fact that my application nodes talk to Redis in an unencrypted fashion. You may want to consider <a href="https://redis.io/topics/encryption">exposing Redis over a proxy which supports encryption</a>.</p>
<p>The other important bit with my implementation is that all of the nodes of my application act as Data Protection key generators. Even if I haven’t seen much problems with this in practice so far, you may want to restrict only one node to be responsible for key generation. You can achieve this by calling DisableAutomaticKeyGeneration like below during <a href="https://docs.microsoft.com/en-us/aspnet/core/security/data-protection/configuration/overview">the configuration stage</a> on your secondary nodes:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">void</span> ConfigureServices(IServiceCollection services)
{
     services.AddDataProtection().DisableAutomaticKeyGeneration();
}</pre></div></div>
<p>I would suggest determining whether a node is primary or not through a configuration value so that you can override this through an environment variable for example.</p>  