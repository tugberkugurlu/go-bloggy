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
When designing a software system, we have somewhat of an idea what the scale of usage is going to be on that system. This could be based off of previous usage patterns, based on the experiment that has been run with a rudimentary functionality, or based on just a pure guess. If you are lucky and mature enough as a business, you should also be able to project how much you are expected to grow as a business for the forseeable future (e.g. next 12, 24 months). All of this data should give you a way to have a baseline number, and be able to extrapolate to understand the load estimations for your system. Being a software engineer, you also have the urge to boil these estimates down to peakiest number of writes/reads per second so that you can reason about these numbers in a relatable way, and can test your system accordingly before going to production.
</p>

<p>
Long story short, what I am trying to get at here is that you have an estimation on the highest number of reads/writes that your system is expected to receive per second. The ideal scenario is also that you want to be on the comfortable side, and will likely want to have 20% over scaling here in case your estimation turns out to be wrong.
</p>

<p>
So far so good, and this is exactly what I would expect an engineer who knows what they are doing and have proper critical thinking skills. These numbers will help you choose the shape and size of the resources you want to set up (e.g. the node size of your Elasticache Redis instance, etc.). That said, we still have problems with this:
</p>

<p>
<ul>
<li>These estimations are just estimations, and they will turn out to be wrong. When they turn out to be higher than you expected, you will struggle with the load. When lower, you will likely burn money unnecessarily and will be overscaled more than you really like it to be.</li>
<li>There will always (I actually mean 'always' here) be unforeseen business activities or external events which will impact the load on your system (e.g. marketing campaigns, etc.). These activities may actually have dramatic impact on the per-second based load. In those circumstances, you need to find a way to accommodate the needs of the new load without actually having any downtime.</li>
</ul>
</p>

<p>
Why am I talking about these? These problems are actually what makes Redis Cluster as the suitable candidate for your needs when they are centered around the writes. For reads, our solution is simple: wire up another Redis node as a replica, and let it serve some portion of the reads, which would take the pressure off from the existing nodes. When the load is lower and you don't need that replica, tear it down to save some £££. All of these operations shouldn't really require too much logic on the clients, and you should be able to get away with by only employing a logic to figure out a new Redis replica addition, and start directing requests to it.
</p>

<p>However, the matter is not that simple for writes. One option we have here is scaling up the nodes (i.e. adding more resources). However, that is going to be a complex operation to perform without introducing a downtime. Besides, when it comes to Redis, your CPU is rarely the issue, and it's throughput that becomes the bottleneck.</p>

<p>If we want to approach this problem the same way we have approach the read scaling issue, there are some questions that really deserve an upfront answer:</p>

<p>
<ul>
<li>How the clients are going to know which node to write data into, and read data from?</li>
<li>What will happen when we add a new node to scale the writes?</li>
<li>What will happen when we remove a new node to scale down?</li>
<li>How can we distribute the load evenly across the nodes?</li>
<li>If we are making multi-command operations (e.g. pipeline requests, <code>MGET</code>, etc.), how are those going to work with this model?</li>
</ul>
</p>

<p>Don't get me wrong here: these are not unique Redis problems. Any data storage system that needs to scale the writes faces these, and there are some common techniques such as <a href="https://en.wikipedia.org/wiki/Shard_(database_architecture)">data sharding</a>, and we are now about to see how Redis tackled these problems through the same technique but with some spice added on top of it for its unique needs.</p>

<h2 href="#redis-cluster-enter">Redis Cluster: Enter</h2>
<p>
Since v3.0, Redis has included an out of the box support for a data sharding solution, which is called Redis Cluster. It provides a way to run a Redis installation where data is automatically sharded across multiple Redis nodes. These Redis nodes are still has the same capabilities as a normal Redis node, and they can have their own replica sets. The only difference is that each node will be only holding the subset of your data, which will depend on your data and Redis' key distribution model.
</p>

<p>
At this point, you should have more questions in your head compared to when you have started reading this post, which is not good :) So, I am hoping to guess what those questions are and try answer them proactively.
</p>

<p>
However, note that <a href="https://redis.io/topics/cluster-spec">Redis Cluster Specification</a> already does a pretty good job on the details. With that in mind, my aim is not to duplicate that documentation here. That said, I want to still highlight the most impactful parts that are valuable to focus more than the others.
</p>

<h3 href="#key-distribution">Key Distribution</a>
<p>
This section is all about essentially answering our first question above regarding which node holds which data. Redis has an interesting way of making this work which seemed to have worked for the use cases I have experienced with. Here is the very high level summary of how it works:
</p>

<ul>
<li>Redis assigns "slot" ranges for each master node within the cluster. These slots are also referred as "hash slots"</li>
<li>These slots are between 0 and 16384, which means each master node in a cluster handles a subset of the 16384 hash slots.</li>
<li>Redis clients can query which node is assigned to which slot range by using the <a href="https://redis.io/commands/cluster-slots"><code>CLUSTER SLOTS</code></a> command. This gives clients a way to be able to directly talk to the correct node for the majority of cases.</li>
<li>For a given Redis key, the hash slot for that key is the result of <code>CRC16(key)</code> modulo <code>16384</code>, where CRC16 here is the implementation of the CRC16 hash function. I am no expect when it comes to cryptography and hashing, but <a href="https://play.golang.org/p/mEmbtCibk_o">here</a> is how this can be done in Go by using the <a href="https://github.com/snksoft/crc">snksoft/crc</a> library. The clients are expected to embed this logic so that they can directly communicate with the correct node with the help of <code>CLUSTER SLOTS</code> command mentioned above.</li>
<li>Same as the single node Redis setup, Redis Cluster uses asynchronous replication between nodes. So, each shard can have its own set of replicas which would be responsible for the same subset of the hash slots as its master. These replicas can be used for failover scenarios as well as distributing the read load (which we will touch on later).</li>
</ul>

<p>
For example, if you have a setup of 3 master nodes with each having 3 replicas, it would look something like the following:
</p>

<img src="https://tugberkugurlu-blog.s3.us-east-2.amazonaws.com/post-images/01ESSD5TSBTEAK416T46X6MHF0-redis-cluster-1.jpg" alt="redis cluster slot assignment" />

<p>
The specific ranges of the hash slots doesn't matter here too much, even the fact that they might be balanced fairly (as we will touch later, we can have influence over slot allocation if we need to). What matters is that it's clear which master node owns.
</p>

<h3 href="#hash-tags">Hash Tags: Getting back into control for your sharding strategy</h3>

<h3 href="#distributing-reads">Distributing Reads</h3>

<h3 href="#redirection">Redirection</h3>

<h3 href="#resharding">Resharding</h3>

<h2 href="#conclusion">Conclusion</h2>

<p>
Redis cluster gives us more ability to scale our systems, especially for the write heavy workloads where we cannot easily predict the demand ahead of time. The sharding model Redis is offering us here is also very interesting where it has the mix of both client and server level logic on where your data is, and how to find it. This gives us an easy way to get started with a rudimentary sharding setup as well as allowing us to optimize our system further by making our clients a bit more claver.
</p>

<p>
I am aware that there are still further unknowns in terms of how clients interact with a Redis cluster setup, how resharding works in practice, etc. However, this post is already too long (there you go, my excuse!) and I hope to cover those in the upcoming posts one by one. If you have any specific areas that you are wondering about Redis Cluster, drop a comment below and I will try to cover them (if I have any experience about those areas). 
</p>

<h2 href="#resources">Resources</a>

<ul>
<li><a href=""></a></li>

<li><a href="https://redis.io/topics/cluster-spec">Redis Cluster Specification</a></li>
<li><a href="https://redis.io/topics/cluster-tutorial">Redis cluster tutorial</a></li>
<li><a href="https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/scaling-redis-cluster-mode-enabled.html">Elasticache: Scaling Clusters in Redis (Cluster Mode Enabled)</a></li>
<li><a href="https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Replication.Redis-RedisCluster.html">Elasticashe: Replication: Redis (Cluster Mode Disabled) vs. Redis (Cluster Mode Enabled)</a></li>
</ul>