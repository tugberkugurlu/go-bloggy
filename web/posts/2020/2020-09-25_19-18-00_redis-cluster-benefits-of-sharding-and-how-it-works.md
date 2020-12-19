---
id: 01ESS2QPGH34KNKDXWCH689M5X
title: Redis Cluster - Benefits of Sharding and How It Works
abstract: Redis is one of the good friends of a backend engineer, and its versatility and easy of use make it so easy to get started. That said, when it comes to scaling it horizontally for writes, it gets a bit more tricky with different level of trade-offs you need to make. In this post, I want to touch on the basics of Redis Cluster, out of the box solution of Redis to the gnarly write scaling problem.
created_at: 2020-12-17 16:29:00.0000000 +0000 UTC
tags:
- Redis
- Databases
- Distributed Systems
- Sharding
slugs:
- redis-cluster-benefits-of-sharding-and-how-it-works
---

<p>
<a href="https://redis.io/">Redis</a> is by far one of the most frequently used data stores. It's fascinating how much our our software-developer-minds go to Redis when we are faced with a data storage problem that requires some level of scale. Even if this might make us feel guilty, I have a somewhat confident assumption that this's the case, and there is probably a relation here to its simplicity: e.g. Redis is 'just' a data structure server, a hash table in the 'cloud', etc. (I know I am a bit exaggerating here, but hopefully you get the idea). Redis also makes digestible and reasonable trade-offs, and it allows us to solve many problems which require certain degree of scale.
</p>

<p>
For a long time, Redis has come with <a href="https://redis.io/topics/replication">an out-of-the-box replication functionality</a>, which allows for a high availability (HA) setup as well as allowing us to scale the reads by distributing the load across replicas with the cost of eventual consistency. However, it was only in April, 2015 that Redis added support for a built-in <a href="https://en.wikipedia.org/wiki/Shard_(database_architecture)">sharding functionality</a> with <a href="https://raw.githubusercontent.com/antirez/redis/3.0/00-RELEASENOTES">its version 3 release</a>. In this post, my aim is to give you more understanding on why you need such a setup, what it actually is, and most important details you need to know about its configuration and implementation details.
</p>

<h2 href="#the-problem">The Problem</h2>
<p>
When designing a software system, we have somewhat of an idea what the scale of usage is going to be on that system. This could be based off of previous usage patterns on the same or similar functionality, based on the data you collected over an experiment that has been run with a rudimentary functionality in a smaller scale, or based on just a pure guess. If you are mature enough as a business, you should also be able to project how much the expected growth is going to be for the forseeable future (e.g. next 12, 24 months). All of this data should help on determining a baseline number, where you can then be able to extrapolate to understand the load estimations for the system that you are designing.
</p>

<p>
Being a software engineer, I bet you also have the urge to boil these estimates down to peakiest number of writes/reads per second so that you can reason about these numbers in a relatable way, and can test your system accordingly before going to production. The ideal scenario is also that you want to be on the comfortable side, and will likely want to have 20% over scaling here in case your estimation turns out to be wrong.
</p>

<p>
So far so good, and this is exactly what I would expect from a software engineer who knows what they are doing and have proper critical thinking skills. The reason is that these numbers will help you choose the shape and size of the resources you want to set up (e.g. the node size of your Elasticache Redis instance, etc.), which will help you optimize your resources. That said, we still have problems with this:
</p>

<p>
<ul>
<li>These estimations are just estimations, and they will almost certainly turn out to be wrong. When they are higher than you expected, you will struggle with the load. When lower, you will likely burn money unnecessarily and will be overscaled more than you really like it to be.</li>
<li>There will always (I actually mean 'always' here) be unforeseen business activities or external events which will impact the load on your system (e.g. marketing campaigns, etc.). These activities may actually have dramatic impact on the per-second based load. In those circumstances, you need to find a way to accommodate the needs of the new load without actually having any downtime.</li>
</ul>
</p>

<p>
Why am I talking about these? These problems are actually what makes Redis Cluster as the suitable candidate for your needs when those problems are especially centered around the writes. For reads, you might still be able to get away with a single master setup by wiring up as many replicas as you need. This should allow you to distribute the read load across replicas at the cost of data consistency gap depending on the replica lag, which would take the pressure off from the master. When the load is lower and you don't need all the replica, you can tear those down to save some £££. All of these operations shouldn't really require too much logic on the clients, and you should really be able to get away with by only employing a logic to figure out a new Redis replica addition, and start directing requests to it.
</p>

<p>However, the matter is not that simple for writes. One option we have here is scaling up the nodes (i.e. adding more resources). However, that is going to be a complex operation to perform without introducing a downtime. There is also a limit to how much you can scale up to (although for the majority of use cases out there, you may never need to go close to that limit). This could still be an option when the issue is with memory. However, not so much for CPU. When it comes to Redis, your CPU is rarely the issue. It's throughput that ends up becoming the bottleneck.</p>

<p>If we want to approach this problem the same way we have approached the read scaling issue, there are some questions that really deserve an upfront answer:</p>

<p>
<ul>
<li>How the clients are going to know which node to write data into, and read data from?</li>
<li>What will happen when we add a new node to scale the writes?</li>
<li>What will happen when we remove a new node to scale down?</li>
<li>How can we distribute the load evenly across the nodes?</li>
<li>If we are making multi-command operations (e.g. pipeline requests, <a href="https://redis.io/commands/mget"><code>MGET</code></a>, etc.), how are those going to work with this model?</li>
</ul>
</p>

<p>Don't get me wrong here: these are not unique Redis problems. Any data storage system that needs to scale the writes face the same challenges, and there are some common techniques such as <a href="https://en.wikipedia.org/wiki/Shard_(database_architecture)">data sharding</a>, and we are now about to see how Redis tackles these problems through the same technique, with some spice added on top to cater for its unique needs.</p>

<h2 href="#redis-cluster-enter">Redis Cluster: Enter</h2>
<p>
Since v3.0, Redis has included an out of the box support for a data sharding solution, which is called <a href="https://redis.io/topics/cluster-tutorial">Redis Cluster</a>. It provides a way to run a Redis installation where data is sharded across multiple Redis nodes as well as providing tools to manage the setup. These Redis nodes still have the same capabilities as a normal Redis node, and they can have their own replica sets. The only difference is that each node will be only holding the subset of your data, which will depend on the shape of the data and Redis' key distribution model (don't worry about this now, we will get to this concept shortly).
</p>

<p>
At this point, you should have more questions in your head compared to when you have started reading this post, which is not good :) So, I am hoping to guess what those questions are and try answer at least some of them proactively.
</p>

<p>
However, note that <a href="https://redis.io/topics/cluster-spec">Redis Cluster Specification</a> already does a pretty good job on the details. With that in mind, my aim is not to duplicate that documentation here. That said, I want to still highlight the most impactful parts that are valuable to focus based on my own experience working with Redis cluster.
</p>

<h3 href="#key-distribution">Key Distribution</h3>
<p>
This section is all about essentially answering our first question above regarding which node holds which data. Redis has an interesting way of making this work which seemed to have worked for the use cases I have experienced with. Here is the very high level summary of how it works:
</p>

<ul>
<li>Redis assigns "slot" ranges for each master node within the cluster. These slots are also referred as "hash slots"</li>
<li>These slots are between <code>0</code> and <code>16384</code>, which means each master node in a cluster handles a subset of the <code>16384</code> hash slots.</li>
<li>Redis clients can query which node is assigned to which slot range by using the <a href="https://redis.io/commands/cluster-slots"><code>CLUSTER SLOTS</code></a> command. This gives clients a way to be able to directly talk to the correct node for the majority of cases.</li>
<li>For a given Redis key, the hash slot for that key is the result of <code>CRC16(key)</code> modulo <code>16384</code>, where CRC16 here is the implementation of the CRC16 hash function. I am no expect when it comes to cryptography and hashing, but <a href="https://play.golang.org/p/mEmbtCibk_o">here</a> is how this can be done in Go by using the <a href="https://github.com/snksoft/crc">snksoft/crc</a> library (note that Redis also has a handy command called <a href="https://redis.io/commands/cluster-keyslot"><code>CLUSTER KEYSLOT</code></a> which performs this operation for you). The clients are expected to embed this logic so that they can directly communicate with the correct node with the help of <code>CLUSTER SLOTS</code> command mentioned above.</li>
<li>Same as the single node Redis setup, Redis Cluster uses asynchronous replication between nodes. So, each shard can have its own set of replicas which would be responsible for the same subset of the hash slots as its master. These replicas can be used for failover scenarios as well as distributing the read load (which we will touch on later).</li>
</ul>

<p>
For example, if you have a setup of 3 master nodes with each having 3 replicas, it would look something like the following:
</p>

<img src="https://tugberkugurlu-blog.s3.us-east-2.amazonaws.com/post-images/01ESSD5TSBTEAK416T46X6MHF0-redis-cluster-1.jpg" alt="redis cluster slot assignment" />

<p>
The specific ranges of the hash slots doesn't matter here too much, even the fact that they might be balanced fairly (as we will touch later, we can have influence over slot allocation if we need to). What matters is that it's clear which master node owns.
</p>

<p>
As an example, I have a local Redis cluster setup which has 3 master nodes, and I am connected to one of them (<code>172.19.197.2</code>) through redis-cli. When I run the <code>CLUSTER SLOTS</code> command, I can see that the node I am connected to handles hash slot range between <code>0</code> and <code>5460</code>:
</p>

<p>
<pre>
172.19.197.2:6379> CLUSTER SLOTS
...
...
2) 1) (integer) 0
   2) (integer) 5460
   3) 1) "172.19.197.2"
      2) (integer) 6379
      3) "fdf56116c8b8f322561c7189574e6092101fa718"
   4) 1) "172.19.197.5"
      2) (integer) 6379
      3) "164dc6aaf77aa0530490f0c9fbf5c8eb9f653a53"
...
...
</pre>
</p>

<p>
I want to set 4 keys, which I already know that falls into the slot range of this node:
</p>

<p>
<pre>
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.7
(integer) 717
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.6
(integer) 4844
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.2
(integer) 4712
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.3
(integer) 585
172.19.197.2:6379> SET coffee_shop_branch.status.7 PERMANENTLY-CLOSED
OK
172.19.197.2:6379> SET coffee_shop_branch.status.6 PERMANENTLY-CLOSED
OK
172.19.197.2:6379> SET coffee_shop_branch.status.2 OPEN
OK
172.19.197.2:6379> SET coffee_shop_branch.status.3 CLOSED
OK
172.19.197.2:6379> KEYS *
1) "coffee_shop_branch.status.7"
2) "coffee_shop_branch.status.6"
3) "coffee_shop_branch.status.2"
4) "coffee_shop_branch.status.3"
</pre>
</p>

<p>
I can also successfully read these the same way I would have done with a single node Redis setup: 
</p>

<p>
<pre>
172.19.197.2:6379> GET coffee_shop_branch.status.7
"PERMANENTLY-CLOSED"
172.19.197.2:6379> GET coffee_shop_branch.status.6
"PERMANENTLY-CLOSED"
172.19.197.2:6379> GET coffee_shop_branch.status.2
"OPEN"
172.19.197.2:6379> GET coffee_shop_branch.status.3
"CLOSED"
</pre>
</p>

<h3 href="#hash-tags">Hash Tags: Getting back into control of your sharding strategy</h3>
<p>
In certain cases, we would like to influence which node our data is stored at. This to be able to group certain keys together so that we can later be able to access them together through a multi-key operation, or a pipeline request.
</p>

<p>
One use case here would be to satisfy the access pattern of retrieving the status of multiple coffee shops within the same city, where we don't have a way to group these together during write time. Therefore, it makes sense to write the status of each coffee shop under their individual keys, and access the ones that we care about through a <a href="https://redis.io/topics/pipelining">pipelining</a>, or <code>MGET</code>.</p> 

<blockquote>
⚠️ I am mentioning <code>MGET</code> as an option here as it is technically a viable option. However, keep in mind that <a href="https://stackoverflow.com/a/61532233/463785"><code>MGET</code> blocks other clients</a> till the whole read operation completes, whereas pipelining doesn't since it's just a way of batching commands. Although you may not see the difference with just a few keys, it's not a good idea to use <code>MGET</code> for too many keys from Redis. I suggest you to perform your own benchmarks for your own use case to see what the threshold is.
</blockquote>

<p>That said, how can we make sure that coffee shops under the same city are co-located within the same node? For example, if we also have the coffee shops with ID <code>1</code> and <code>4</code>, they are not going to be stored within the same node as coffee shops with ID <code>2</code>, <code>3</code>, <code>6</code> and <code>7</code> based on our current setup (remember: the node at <code>172.19.197.2</code> is responsible for hash slot range of <code>0</code>-<code>5460</code>):
</p>

<p>
<pre>
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.1
(integer) 8715
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.4
(integer) 12974

172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.2
(integer) 4712
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.3
(integer) 585
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.6
(integer) 4844
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.status.7
(integer) 717
</pre>
</p>

<p>
You can also see that Redis will also complain when we try to <code>MGET</code> all of these keys:
</p>

<p>
<pre>
172.19.197.2:6379> MGET coffee_shop_branch.status.1 coffee_shop_branch.status.2 coffee_shop_branch.status.3 coffee_shop_branch.status.4 coffee_shop_branch.status.6 coffee_shop_branch.status.7
(error) CROSSSLOT Keys in request don't hash to the same slot
</pre>
</p>

<p>
We can also see the same behavior even if we remove <code>coffee_shop_branch.status.1</code> and <code>coffee_shop_branch.status.4</code> from the list of keys. This is because the <code>MGET</code> can only succeed if all of the keys belong to same slot as the error message suggests.
</p>

<p>
<pre>
172.19.197.2:6379> MGET coffee_shop_branch.status.2 coffee_shop_branch.status.3 coffee_shop_branch.status.6 coffee_shop_branch.status.7
(error) CROSSSLOT Keys in request don't hash to the same slot
</pre>
</p>

<p>
This is where the concept of <a href="https://redis.io/topics/cluster-spec#keys-hash-tags">hash tags</a> comes in. Hash tags allow us to force certain keys to be stored in the same hash slot. I encourage you the read the linked section of the spec to understand better how hash tags work as I am going to skip some corner cases here, but in a nutshell, the concept is really simple from the usage point of view: when the Redis key contains <code>"{...}"</code> pattern only the substring between <code>{</code> and <code>}</code> is hashed in order to obtain the hash slot.
</p>

<p>
For our use case, this means that we can change our key structure from <code>coffee_shop_branch.status.COFFEE-SHOP-ID</code> to something like <code>coffee_shop_branch.{city_CITY-ID}.status.COFFEE-SHOP-ID</code>. The exact shape of the key is not important here. What's important is that the value between curly braces which is the city ID prefixed with <code>city_</code> for readability purposes.
</p>

<p>
For the example that we have been working with, and with the assumption that the coffee shops with ID <code>1</code>, <code>4</code>, <code>2</code>, <code>3</code>, <code>6</code> and <code>7</code> are all with the same city, let's say that it's the city with ID <code>4</code>, the keys will shape up as following, and we can see from the <code>CLUSTER KEYSLOT</code> command outcome that all of these keys are hashed to the same slot:
</p>

<p>
<pre>
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.{city_4}.status.1
(integer) 1555
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.{city_4}.status.4
(integer) 1555
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.{city_4}.status.2
(integer) 1555
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.{city_4}.status.3
(integer) 1555
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.{city_4}.status.6
(integer) 1555
172.19.197.2:6379> CLUSTER KEYSLOT coffee_shop_branch.{city_4}.status.7
(integer) 1555
</pre>
</p>

<p>
We can also see that <code>MGET</code> will start working:
</p>

<p>
<pre>
redis_1:6379> MGET coffee_shop_branch.{city_4}.status.1 coffee_shop_branch.{city_4}.status.4 coffee_shop_branch.{city_4}.status.2 coffee_shop_branch.{city_4}.status.3 coffee_shop_branch.{city_4}.status.6 coffee_shop_branch.{city_4}.status.7
1) "OPEN"
2) "CLOSED"
3) "OPEN"
4) "CLOSED"
5) "PERMANENTLY-CLOSED"
6) "PERMANENTLY-CLOSED"
</pre>
</p>

<p>So, hash tags are great, and we should use them all the time, right? Not so fast! This approach can make a notable positive impact on the latency of your application, and resource utilization of your redis nodes. However, there is a drawback here which might be a big worry for you depending on your load and data distribution: the Hot Shard problem (a.k.a. Hot Key problem). In our use case, this can be a significant problem when certain cities hold way more coffee shops than the others, or the access to certain cities are significantly higher even if the data sizes are the same. I will leave <a href="http://highscalability.com/blog/2010/10/15/troubles-with-sharding-what-can-we-learn-from-the-foursquare.html">this super informative post</a> from 2010 about one of the Foursquare outages which was caused by the exact same problem.</p>

<p>
Hash tags is a tool that can help you, but there is unfortunately no magic bullet here. You need to understand your use case, data distribution, and test different setups to understand what might work for you the best.
</p>

<h3 href="#redirection">Redirection</h3>

<h3 href="#distributing-reads">Distributing Reads</h3>

<h2 href="#conclusion">Conclusion</h2>

<p>
Redis cluster gives us more ability to scale our systems, especially for the write heavy workloads where we cannot easily predict the demand ahead of time. The sharding model Redis is offering us here is also very interesting where it has the mix of both client and server level logic on where your data is, and how to find it. This gives us an easy way to get started with a rudimentary sharding setup as well as allowing us to optimize our system further by making our clients a bit more claver.
</p>

<p>
I am aware that there are still further unknowns in terms of how to actually initialize a Redis cluster setup from scratch, details of how clients interact with a Redis cluster setup, how maintenance/operational side of the cluster setup actually works (e.g. resharding), etc. However, this post is already too long (there you go, my excuse!), and I hope to cover those in the upcoming posts one by one. If you have any specific areas that you are wondering about Redis Cluster, drop a comment below and I will try to cover them (if I have any experience about those areas). 
</p>

<h2 href="#resources">Resources</h2>

<ul>
<li><a href="https://redis.io/topics/cluster-spec">Redis Cluster Specification</a></li>
<li><a href="https://redis.io/topics/cluster-tutorial">Redis cluster tutorial</a></li>
<li><a href="https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/scaling-redis-cluster-mode-enabled.html">Elasticache: Scaling Clusters in Redis (Cluster Mode Enabled)</a></li>
<li><a href="https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Replication.Redis-RedisCluster.html">Elasticashe: Replication: Redis (Cluster Mode Disabled) vs. Redis (Cluster Mode Enabled)</a></li>
</ul>