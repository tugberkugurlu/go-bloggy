---
id: 01ESYWJ3NHDFF5G7H4JKCY3DS0
title: Redis Cluster - Benefits of Sharding and How It Works
abstract: Redis is one of the good friends of a backend engineer, and its versatility and ease of use make it convenient to get started. That said, when it comes to scaling it horizontally for writes, it gets a bit more tricky with different level of trade-offs you need to make. In this post, I want to touch on the basics of Redis Cluster, out of the box solution of Redis to the gnarly write scaling problem.
created_at: 2020-12-20 01:34:00.0000000 +0000 UTC
tags:
- Redis
- Databases
- Distributed Systems
- Sharding
slugs:
- redis-cluster-benefits-of-sharding-and-how-it-works
---

<blockquote>
<h3>Content</h3>
<ul>
<li><a href="#the-problem">The Problem</a></li>
<li><a href="#redis-cluster-enter">Redis Cluster: Enter</a></li>
<ul>
    <li><a href="#key-distribution">Key Distribution</a></li>
    <li><a href="#hash-tags">Hash Tags: Getting back into control of your sharding strategy</a></li>
    <li><a href="#redirection">Redirection</a></li>
    <li><a href="#distributing-reads">Distributing Reads</a></li>
</ul>
<li><a href="#conclusion">Conclusion</a></li>
</ul>
</blockquote>

<p>
<a href="https://redis.io/">Redis</a> is by far one of the most frequently used data stores. It's fascinating how much our our software-developer-minds go to Redis when we are faced with a data storage problem that requires some level of scale. Even if this might make us feel guilty, I have a somewhat confident assumption that this's the case, and there is probably a relation here to its simplicity: e.g. Redis is 'just' a data structure server, a hash table in the 'cloud', etc. (I know I am a bit exaggerating here, but hopefully you get the idea). Redis also makes digestible and reasonable trade-offs, and it allows us to solve many problems which require certain degree of scale.
</p>

<p>
For a long time, Redis has come with <a href="https://redis.io/topics/replication">an out-of-the-box replication functionality</a>, which allows for a high availability (HA) setup as well as allowing us to scale the reads by distributing the load across replicas with the cost of eventual consistency. However, it was only in April, 2015 that Redis added support for a built-in <a href="https://en.wikipedia.org/wiki/Shard_(database_architecture)">sharding functionality</a> with <a href="https://raw.githubusercontent.com/antirez/redis/3.0/00-RELEASENOTES">its version 3 release</a>. I have been working with several Redis Cluster setups for a while, and have probably read the Redis Cluster spec at least couple of times. In this post, my aim is to give you more understanding on what problem Redis Cluster actually solves, why you need such a setup, and most important details you need to know about its configuration and implementation details based on my own experience.
</p>

<h2 id="the-problem">The Problem</h2>
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
Why am I talking about these? These problems are actually what makes Redis Cluster as the suitable candidate for your needs when those problems are especially centered around the writes. For reads, you might still be able to get away with a single master setup by wiring up as many replicas as you need. This should allow you to distribute the read load across replicas at the cost of data consistency gap depending on the replication lag, which would take the pressure off from the master. When the load is lower and you don't need all the replica, you can tear those down to save some £££. All of these operations shouldn't really require too much logic on the clients, and you should really be able to get away with by only employing a logic to figure out a new Redis replica addition, and start directing requests to it.
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

<h2 id="redis-cluster-enter">Redis Cluster: Enter</h2>
<p>
Since v3.0, Redis has included an out of the box support for a data sharding solution, which is called <a href="https://redis.io/topics/cluster-tutorial">Redis Cluster</a>. It provides a way to run a Redis installation where data is sharded across multiple Redis nodes as well as providing tools to manage the setup. These Redis nodes still have the same capabilities as a normal Redis node, and they can have their own replica sets. The only difference is that each node will be only holding the subset of your data, which will depend on the shape of the data and Redis' key distribution model (don't worry about this now, we will get to this concept shortly).
</p>

<p>I have configured a local Redis cluster setup to use throughout this blog post, and with the help of <a href="https://redis.io/commands/cluster-nodes"><code>CLUSTER NODES</code></a> command, I can see its high level structure:</p>

<p>
<pre>
172.19.197.2:6379> CLUSTER NODES
b7366bdbb09dbb20dcf0d4f8b7281c98f7e3b78e 172.19.197.7:6379@16379 master - 0 1608418117542 10 connected 10923-16383
164dc6aaf77aa0530490f0c9fbf5c8eb9f653a53 172.19.197.5:6379@16379 slave fdf56116c8b8f322561c7189574e6092101fa718 0 1608418118557 12 connected
f75939944d18ee12995c60d4cc9fcc1e53458d32 172.19.197.3:6379@16379 slave 88875e065f5ecf24b5adde973223a7799aee4521 0 1608418117949 11 connected
fdf56116c8b8f322561c7189574e6092101fa718 172.19.197.2:6379@16379 myself,master - 0 1608418118000 12 connected 0-5460
1c822510aa0f349a9b12cba1c68bc98feab5433e 172.19.197.4:6379@16379 slave b7366bdbb09dbb20dcf0d4f8b7281c98f7e3b78e 0 1608418118000 10 connected
88875e065f5ecf24b5adde973223a7799aee4521 172.19.197.6:6379@16379 master - 0 1608418118963 11 connected 5461-10922
</pre>
</p>

<p>
You can learn more about the serialization format of this output from <a href="https://redis.io/commands/cluster-nodes#serialization-format">the doc</a>, but let me take a stab at summarizing it:
</p>

<ul>
<li>We have setup of 3 master nodes with each having one replica.</li>
<li>We are currently connected to the node at <code>172.19.197.2:6379</code>, and its node ID is <code>fdf56116c8b8f322561c7189574e6092101fa718</code>. We know this is the node we are connected as the <code>myself</code> flag indicates the the node you are contacted. This node is also one of the master nodes.</li>
<li>The node that we are connected is shown to be responsible for <code>0</code>-<code>5460</code> slot range (don't worry about what exactly this is now, we will shortly get to this).</li>
<li>The node at <code>172.19.197.5:6379</code> is the replica of the current node which we are connected to. We know this as the node ID of <code>fdf56116c8b8f322561c7189574e6092101fa718</code> is shown under the <code>master</code> column and we know that this the ID of the node that we are connected to.</li>
</ul>

<p>
At this point, you should have more questions in your head compared to when you have started reading this post, which is not good :) So, I am hoping to guess what those questions are and try answer at least some of them proactively.
</p>

<p>
However, note that <a href="https://redis.io/topics/cluster-spec">Redis Cluster Specification</a> already does a pretty good job on the details. With that in mind, my aim is not to duplicate that documentation here. That said, I want to still highlight the most impactful parts that are valuable to focus based on my own experience working with Redis cluster.
</p>

<h3 id="key-distribution">Key Distribution</h3>
<p>
This section is all about essentially answering our first question above regarding which node holds which data. Redis has an interesting way of making this work which seemed to have worked for the use cases I have experienced with. Here is the very high level summary of how it works:
</p>

<ul>
<li>Redis assigns "slot" ranges for each master node within the cluster. These slots are also referred as "hash slots"</li>
<li>These slots are between <code>0</code> and <code>16384</code>, which means each master node in a cluster handles a subset of the <code>16384</code> hash slots.</li>
<li>Redis clients can query which node is assigned to which slot range by using the <a href="https://redis.io/commands/cluster-slots"><code>CLUSTER SLOTS</code></a> command. This gives clients a way to be able to directly talk to the correct node for the majority of cases.</li>
<li>For a given Redis key, the hash slot for that key is the result of <code>CRC16(key)</code> modulo <code>16384</code>, where <code>CRC16</code> here is the implementation of the CRC16 hash function. I am no expect when it comes to cryptography and hashing, but <a href="https://play.golang.org/p/mEmbtCibk_o">here</a> is how this can be done in Go by using the <a href="https://github.com/snksoft/crc">snksoft/crc</a> library. Note that Redis also has a handy command called <a href="https://redis.io/commands/cluster-keyslot"><code>CLUSTER KEYSLOT</code></a> which performs this operation for you per given Redis key. <a href="https://redis.io/topics/cluster-spec#clients-first-connection-and-handling-of-redirections">The clients are expected to embed this logic</a> so that they can directly communicate with the correct node with the help of <code>CLUSTER SLOTS</code> command mentioned above.</li>
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

<h3 id="hash-tags">Hash Tags: Getting back into control of your sharding strategy</h3>
<p>
In certain cases, we would like to influence which node our data is stored at. This is to be able to group certain keys together so that we can later access them together through a multi-key operation, or through <a href="https://redis.io/topics/pipelining">pipelining</a>.
</p>

<p>
One use case here would be to satisfy the access pattern of retrieving the status of multiple coffee shops within the same city, where we don't have a way to group these together during write time. Therefore, it makes sense to write the status of each coffee shop under their individual keys, and access the ones that we care about through pipelining, or <code>MGET</code>.</p> 

<blockquote>
⚠️ I am mentioning <code>MGET</code> as an option here as it is technically a viable option. However, keep in mind that <a href="https://stackoverflow.com/a/61532233/463785"><code>MGET</code> blocks other clients</a> till the whole read operation completes, whereas pipelining doesn't since it's just a way of batching commands. Although you may not see the difference with just a few keys, it's not a good idea to use <code>MGET</code> for too many keys. I suggest for you to perform your own benchmarks for your own use case to see what the threshold might be here.
</blockquote>

<p>Idea is solid but there is still a question: how can we make sure that coffee shops under the same city are co-located within the same node? For example, if we also have the coffee shops with ID <code>1</code> and <code>4</code>, they are not going to be stored within the same node as coffee shops with ID <code>2</code>, <code>3</code>, <code>6</code> and <code>7</code> based on our current setup (remember: the node at <code>172.19.197.2</code> is responsible for hash slot range of <code>0</code>-<code>5460</code>):
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
We can also see the same behavior even if we remove <code>coffee_shop_branch.status.1</code> and <code>coffee_shop_branch.status.4</code> from the list of keys. This is because the fact that <code>MGET</code> can only succeed if all of the keys belong to same slot as the error message suggests.
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
We can also see that <code>MGET</code> will start working as expected with these keys:
</p>

<p>
<pre>
172.19.197.2:6379> MGET coffee_shop_branch.{city_4}.status.1 coffee_shop_branch.{city_4}.status.4 coffee_shop_branch.{city_4}.status.2 coffee_shop_branch.{city_4}.status.3 coffee_shop_branch.{city_4}.status.6 coffee_shop_branch.{city_4}.status.7
1) "OPEN"
2) "CLOSED"
3) "OPEN"
4) "CLOSED"
5) "PERMANENTLY-CLOSED"
6) "PERMANENTLY-CLOSED"
</pre>
</p>

<p>So, hash tags are great, and we should use them all the time, right? Not so fast! This approach can make a notable positive impact on the latency of your application, and resource utilization of your redis nodes. However, there is a drawback here which might be a big worry for you depending on your load and data distribution: the Hot Shard problem (a.k.a. Hot Key problem). In our use case for instance, this can be a significant problem when certain cities hold way more coffee shops than the others, or the access for certain cities are significantly higher even if the data sizes are the same. I will leave <a href="http://highscalability.com/blog/2010/10/15/troubles-with-sharding-what-can-we-learn-from-the-foursquare.html">this super informative post</a> from 2010 here, which is about one of the Foursquare outages. You will quickly realise after reading <a href="https://web.archive.org/web/20131114075609/http://blog.foursquare.com/2010/10/05/so-that-was-a-bummer/">the post-mortem</a> that it was caused by the exact same problem.</p>

<p>
Hash tags is a tool that can help you, but there is unfortunately no magic bullet here. You need to understand your use case, data distribution, and test different setups to understand what might work for you the best.
</p>

<h3 id="redirection">Redirection</h3>
<p>Apart from the <code>MGET</code> example above, we have been playing it by the rules so far: knowingly issuing commands against the nodes that actually hold the data for the given keys. We were able to do this through the couple of cluster commands that Redis provides such as <code>CLUSTER SLOTS</code> and <code>CLUSTER KEYSLOT</code>.</p>

<p>
What would happen if we do the opposite though: issuing a command against a Redis node which doesn't actually own the hash slot for the given key? Here is the answer:
</p>

<p>
<pre>
172.19.197.2:6379> get coffee_shop_branch.status.1
(error) MOVED 8715 172.19.197.6:6379
</pre>
</p>

<p>Redis is erroring, but erroring in a more clever way than you probably have guessed. The error itself includes the hash slot of the key, and the ip:port of the instance that owns that hash slot and can serve the query. This is called <a href="https://redis.io/topics/cluster-spec#moved-redirection">MOVED redirection</a> in Redis spec, and all the Redis Cluster clients are expected to handle this error appropriately so that they can eventually succeed the request by connecting to the correct node and issuing the command there.</p>

<p><a href="https://redis.io/topics/rediscli">redis-cli</a>, as being one of the Redis clients, also knows how to handle <code>MOVED</code> redirection. The CLI utility implements basic cluster support when started with the <code>-c</code> switch.</p>

<p>
<pre>
➜ docker run -it --rm \
    --net redis-cluster_redis_cluster_network \
    redis \
    redis-cli -h redis_1
redis_1:6379> get coffee_shop_branch.status.1
(error) MOVED 8715 172.19.197.6:6379
redis_1:6379> exit

➜ docker run -it --rm \
    --net redis-cluster_redis_cluster_network \
    redis \
    redis-cli -c -h redis_1
redis_1:6379> get coffee_shop_branch.status.1
-> Redirected to slot [8715] located at 172.19.197.6:6379
"OPEN"
172.19.197.6:6379> 
</pre>
</p>

<p>
You can see that on the first case when we connected to a Redis node through redis-cli without the <code>-c</code> switch, we got the <code>MOVED</code> redirection. However, in the case where we used the <code>-c</code> switch, the client handled the redirection transparently by connecting to the given Redis node, and issuing the command there.
</p>

<p>However, Redis already gives a way to identify which master node is responsible for which hash slot range, and Redis cluster clients should also be able to generate the hash of a given key to figure out which node to connect to. So, why is this feature useful? There are two main key reasons that I am aware of:</p>

<p>First one is that Redis cluster specification doesn't require Redis Clsuter clients to be clever about routing, meaning that clients don't need to keep track of which master nodes serve for which hash slot range. Instead, they can just have the logic to be able to handle the redirection to be considered a complete Redis Cluster client. I don't exactly know what the reason was for this, but I presume this made it easier for existing Redis clients to adopt to be a Redis Cluster client at the time. That said, these clients have a major drawback that they are so much inefficient compared to their clever counterparts since these clients have a high change of making at least twice the number of requests than they need to for the majority of the operations they perform.</p>

<p>
Another reason why we have the <code>MOVED</code> redirection in place (probably the most important one) is related to <a href="https://redis.io/topics/cluster-spec#cluster-live-reconfiguration">resharding</a>. For instance, when a new master node is added to the Redis Cluster to offload some of the pressure from the existing nodes, it's expected to perform some of the cluster reconfiguration operations to move certain hash slot ranges from the existing nodes to the new node. This would trigger a what-is-commonly-known-as resharding operation, and Redis aims to handle this without causing a disruption. However, when this happens and certain hash slot ranges are being moved from one node to another, there is a chance that the client can have the stale information about the cluster during this phase. This might cause the client to connect to the old node which used to be responsible for a given hash slot, instead of the correct node which took charge of that slot after the client retrieved the latest state of the cluster. This is where the <code>MOVED</code> redirection is handy, and it also hints to the client to reload its cluster configuration.
</p>

<p>I am aware that we haven't touched on the resharding point in depth yet (and we won't be in this post), but redirection is such a fundamental concept of the Redis Cluster specification that I wanted briefly to go over at a high level. Also note that there is another type of redirection which is known as <a href="https://redis.io/topics/cluster-spec#ask-redirection">ASK redirection</a>, and we won't be covering that here at all since it's fundamentally related to resharding and that one really deserves its own post.</p>

<h3 id="distributing-reads">Distributing Reads</h3>
<p>
The last point I want to touch on is around scaling reads, where we can make use of the replicas to distribute the load. For example, with the setup that we have been working with in this post, we have a replica per each master node. Considering we have 3 master nodes, by default, 3 nodes are serving reads and writes. However, we can utilize the replicas to serve the read commands which would essentially double the number of nodes that can serve reads.
</p>

<p>This is great but it's at the cost of data consistency since Redis uses by default asynchronous replication unless you are using the <a href="https://redis.io/commands/wait">WAIT</a> command to enforce a synchronous replication during write time.</p>

<p>Let's assume that we are OK with the data inconsistency, and we are monitoring the replication lag. How can we utilize these replicas for reads? We can start by exploring this through redis-cli. From <a href="#redis-cluster-enter">our previous exploration</a>, we know that the node at <code>172.19.197.5:6379</code> is the replica of the node at <code>172.19.197.2:6379</code>. So, let's connect to that node directly, and issue a <code>GET</code> command there:</p>

<p>
<pre>
➜ docker run -it --rm \
    --net redis-cluster_redis_cluster_network \
    redis \
    redis-cli -c -h 172.19.197.5
172.19.197.5:6379> get coffee_shop_branch.{city_4}.status.4
-> Redirected to slot [1555] located at 172.19.197.2:6379
"CLOSED"
172.19.197.2:6379> 
</pre>
</p>

<p>That's a surprising outcome as we were being redirected to the node at <code>172.19.197.2:6379</code> which is the master node of the replica that we were connected to. From this, it seems like the replica either doesn't hold the data that we need, or it doesn't allow any read operations.</p>

<p>
Let's first check whether it actually holds the data. Looking at the <a href="https://redis.io/commands/keys">KEYS</a> stored at that node, it seems like it has the data that we need:
</p>

<p>
<pre>
172.19.197.5:6379> KEYS *
 1) "coffee_shop_branch.status.3"
 2) "coffee_shop_branch.status.6"
 3) "coffee_shop_branch.status.7"
 4) "coffee_shop_branch.{city_4}.status.2"
 5) "coffee_shop_branch.{city_4}.status.4"
 6) "coffee_shop_branch.status.2"
 7) "coffee_shop_branch.{city_4}.status.7"
 8) "coffee_shop_branch.{city_4}.status.3"
 9) "coffee_shop_branch.{city_4}.status.6"
10) "coffee_shop_branch.{city_4}.status.1"
</pre>
</p>

<p>When we check the replica status, we can also see that the replica is up-to-date:</p>

<p>
<pre>
172.19.197.5:6379> INFO replication
# Replication
role:slave
master_host:172.19.197.2
master_port:6379
master_link_status:up
master_last_io_seconds_ago:8
master_sync_in_progress:0
...
...
</pre>
</p>

<p>It seems like the replica doesn't allow us to perform any read operations, and this is expected which is also documented inside the <a href="https://redis.io/topics/cluster-spec#scaling-reads-using-slave-nodes">Redis Cluster spec</a>:</p>

<blockquote>
Normally slave nodes will redirect clients to the authoritative master for the hash slot involved in a given command, however clients can use slaves in order to scale reads using the <code>READONLY</code> command.
</blockquote>

<p>
<a href="https://redis.io/commands/readonly"><code>READONLY</code></a> command enables read queries for a connection to a Redis Cluster replica node. This command hints to the server that the client is OK with the potential data inconsistency. This command needs to be sent per each connection to the replica nodes and ideally should be sent right after the connection is established. 
</p>

<p>
<pre>
➜ docker run -it --rm \
    --net redis-cluster_redis_cluster_network \
    redis \
    redis-cli -c -h 172.19.197.5
172.19.197.5:6379> READONLY
OK
172.19.197.5:6379> get coffee_shop_branch.{city_4}.status.4
"CLOSED"
172.19.197.5:6379> 
</pre>
</p>

<p>To be honest, I remember that this threw me off when I first realized this behavior. However, it makes sort of a sense to be explicit when it comes to reading stale data. My only gripe about it is the name of the command which is sort of confusing. That said, you get used to it after a while, and it's well supported by the clients (e.g. <a href="https://github.com/go-redis">go-redis</a> client has a way for you to <a href="https://github.com/go-redis/redis/blob/143859e34596a8e80ee858b5842d503d86572249/cluster.go#L38-L45">configure this as well as being able to configure the replica routing behavior</a>).</p>

<h2 id="conclusion">Conclusion</h2>

<p>
Redis cluster gives us the ability to scale our Redis setup horizontally not just for reads but also for writes, and you should consider it especially if you have a write heavy workload where you cannot easily predict the demand ahead of time. The sharding model Redis is offering us is also very interesting where it has the mix of both client and server level logic on where your data is, and how to find it. This gives us an easy way to get started with a rudimentary sharding setup as well as allowing us to optimize our system further by making our clients a bit more clever.
</p>

<p>
I am aware that there are still further unknowns in terms of how to actually initialize a Redis cluster setup from scratch, details of how clients interact with a Redis cluster setup, how maintenance/operational side of the cluster setup actually works (e.g. resharding), etc. However, this post is already too long (there you go, my excuse!), and I hope to cover those in the upcoming posts one by one. If you have any specific areas that you are wondering about Redis Cluster, drop a comment below and I will try to cover them if I have any experience around those areas. 
</p>

<h2 href="#resources">Resources</h2>

<ul>
<li><a href="https://redis.io/topics/cluster-spec">Redis Cluster Specification</a></li>
<li><a href="https://redis.io/topics/cluster-tutorial">Redis cluster tutorial</a></li>
<li><a href="https://lmgtfy.app/?q=site%3Aredis.io%2Fcommands%2Fcluster-">Redis cluster commands</a></li>
<li><a href="https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/scaling-redis-cluster-mode-enabled.html">Elasticache: Scaling Clusters in Redis (Cluster Mode Enabled)</a></li>
<li><a href="https://redis.io/topics/pipelining">Using pipelining to speedup Redis queries</a></li>
<li><a href="https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Replication.Redis-RedisCluster.html">Elasticashe: Replication: Redis (Cluster Mode Disabled) vs. Redis (Cluster Mode Enabled)</a></li>
<li><a href="http://highscalability.com/blog/2010/10/15/troubles-with-sharding-what-can-we-learn-from-the-foursquare.html">Troubles With Sharding - What Can We Learn From The Foursquare Incident?</a></li>
</ul>