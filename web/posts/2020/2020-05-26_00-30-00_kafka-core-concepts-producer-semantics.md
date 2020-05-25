---
id: 01E972TEACJE3TQW6Q59C0X4NJ
title: Kafka Core Concepts and Producer Semantics 
abstract: Understanding the intrinsic behaviors of a component your system is making use of will make you fear less about it as you will have a better understanding on what might happen under which circumstances. In this post, we will start to understand the core concepts of Kafka as well as diving deep into publishing semantics.
created_at: 2020-05-26 00:30:00.0000000 +0000 UTC
tags:
- Kafka
- Distributed Systems
- Messaging
slugs:
- kafka-core-concepts-and-producer-semantics
---

<p id="91d84e99-35a9-4e5b-9914-280efbee84a6" class="">Being able to pass data around within a distributed
    system is the one of the the most crucial aspects of the success for your business, especially when you
    are dealing with large number of users, reads and writes. It&#x27;s usual that for a given data write
    for an entity, you will have N number of read patterns, not just one. <a
        href="https://kafka.apache.org/">Apache Kafka</a> is one of the most effective ways to enable that
    data distribution within a complex system. I have had the chance to use Kafka at work for more than a
    year now. However, it has always been implicit and I never needed to understand its intrinsic semantics
    (standing on the shoulders of giants). I have spent this extended weekend reading the Kafka
    documentation and running some local examples with Kafka to understand it in details, not just at a high
    level.</p>
<p id="831c386d-05a1-46bc-909a-14ed649663be" class="">Kafka already has <a
        href="https://kafka.apache.org/documentation">a great documentation</a>, which is very detailed and
    clear. The intention with this post is not to replicate that document. Instead, it&#x27;s to pull out
    bits and pieces which helped me understand Kafka better, and increased by trust. As it has been said in
    Batman Begins movie (which is one of my all-time favourites): &quot;<a
        href="https://www.magicalquote.com/moviequotes/you-always-fear-what-you-dont-understand/">You always
        fear what you don&#x27;t understand</a>&quot;, and the main outcome here is to remove that fear :)
    The post is written by a someone, which is me, who has previous experience with messaging systems such
    as RabbitMQ, Amazon SQS, and Azure Service Bus. So, I might be overlooking some important aspects which
    you may also need if you don&#x27;t have this background. If that&#x27;s the case, it might be useful to
    first understand <a href="https://kafka.apache.org/documentation/#uses">some use cases where Kafka might
        fit in</a>.</p>
<h2 id="ee274949-d1ca-447e-aa21-db03ad51d1af" class="">Concepts</h2>
<p id="2e13d642-9e0a-46a7-81f7-36f504b75e0f" class="">
</p>
<p>
<img src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/01E972A8Z95PB98E3082YJ1P2M-log_anatomy.png" alt="Kafka Topic" />
</p>
<p id="91f725c1-83dc-47f0-90e2-55c0ba86db27" class="">Let&#x27;s first understand some of the high level
    concepts of Kafka, which will allow us to get started and work on a sample later on. This is by all
    means not an exhaustive list of concepts in Kafka but will be enough to get us going by allowing us to
    extract some facts as well as allowing us to make some assumptions with high confidence.</p>
<p id="e53c9b17-6500-4ca6-9cee-bc21f99dd908" class="">The most important concept of Kafka is a
    <strong>Topic</strong>. Topics in Kafka is a place where you can logically group your messages into.
    When I say logically, I don&#x27;t mean a schema or anything. You can think of it as just a bucket where
    your data will end up in the order they appear, and can also be retrieved in the same order (i.e.
    continually appended to a structured commit log). Topics can be subscribed to by one or more consumers,
    which we will touch on that a few points later, but this means that Kafka doesn&#x27;t have exact
    message queue semantics, which ensures that the data is gone as soon as one consumer processes the data.
</p>
<p id="94138dad-c82c-41e8-89ef-fc56d65c9609" class="">These message are called <strong>Records</strong>,
    which are durably persisted in the Kafka cluster regardless of the fact that they have been consumed or
    not. This differentiates Kafka from queuing systems such as RabbitMQ or SQS, where messages vanish after
    they are being consumed and processed. Using Kafka for storing records permanently is a perfectly valid
    choice. However, if this is not desired, Kafka also give you a retention configuration options to
    specific how long you want to hold onto records per topic basis. </p>
<p id="d9dba195-478a-497e-9d03-3a04e571265b" class="">The records gets into (i.e. written) a topic through a
    <a href="https://kafka.apache.org/documentation/#intro_producers"><strong>producer</strong></a>, who are
    also responsible for choosing which record to assign to which partition within the topic. In other
    words, data sharding is handled by the clients which publish data to a particular topic. Depending on
    what client you use, you may have different options on how to distribute data across the partitions,
    e.g. round robin, your custom sharding strategy, etc.</p>
<p id="cbeac82e-a422-4047-aeab-f29ceaee0a73" class="">The records within a specific topic are consumed (i.e.
    read) by a <a
        href="https://kafka.apache.org/documentation/#intro_consumers"><strong>consumer</strong></a>, which
    is part of a <strong>consumer group</strong>. Consumer groups allow records to be processed in parallel
    by the consumer instances (which are associated to that group, and can live in separate processes or
    machines) with a guarantee that a record is only delivered to one consumer instance. A consumer instance
    within a consumer group will own one or more partitions exclusively, which means that you can have at
    max N number of consumer if you have N partitions.</p>
<p id="c879ef20-895f-43c1-a44f-835dea800f6a" class="">So, based on these, here are some take aways which I
    was able to further unpack by following up:</p>
<ul id="a636bdfa-af19-42f1-8f94-d31454f22312" class="bulleted-list">
    <li>Data stored in Kafka is immutable, meaning that it cannot be updated. In other words, Kafka is
        working with an append-only data structure and all you can do with it is to ask for the next record
        and reset to current pointer.</li>
</ul>
<ul id="dd72b22e-8079-4871-aef2-eeea83300e13" class="bulleted-list">
    <li>Kafka has a distributed nature to cater your scalability and high availability needs.</li>
</ul>
<ul id="0543d407-ec23-4f2c-a8f0-451701cbd62c" class="bulleted-list">
    <li>Kafka guarantees ordering for the records but this is only per partition basis and how you retry
        messages can also have an impact on this order. Therefore, it&#x27;s safest to assume at the
        consumption level that Kafka won&#x27;t give you a message ordering guarantees, and you may need to
        understanding the details of this further depending on how your messages are distributed across the
        partitions, and how you plan to process that data. </li>
</ul>
<ul id="44915cbf-9b6c-4123-aaad-7a24aba376e3" class="bulleted-list">
    <li>Kafka is consumer driven, which means that consumer is in charge of determining reading the data
        from which position they like. In practical terms, this means that the consumer can reset the offset
        and start from wherever it wants to. Check out the <a
            href="https://kafka.apache.org/documentation/#impl_offsettracking">Offset Tracking</a> and <a
            href="https://kafka.apache.org/documentation/#basic_ops_consumer_group">Consumer Group
            Management</a> sections for more info on this.</li>
</ul>
<ul id="205e5f65-a123-42c3-8f19-81ace05b8e98" class="bulleted-list">
    <li><a href="https://kafka.apache.org/documentation/#basic_ops_cluster_expansion">It&#x27;s possible to
            add new nodes to your cluster</a>. The data distribution to this node though <a
            href="https://kafka.apache.org/documentation/#basic_ops_automigrate">needs to be triggered
            manually</a>.</li>
</ul>
<ul id="311e25fe-5b27-4a0a-b2e6-2165a19d7987" class="bulleted-list">
    <li>Related to above, <a href="https://kafka.apache.org/documentation.html#basic_ops_modify_topic">you
            can increase the number of partitions for a given topic</a>. However, this is an operation you
        do not want to perform without proactively thinking through the consequences since the way you
        publish data to Kafka might be impacted by this, if, for example, your sharding strategy is rely on
        knowing the partition count (i.e. <code>hash(key) % number_of_partitions</code>). It&#x27;s also
        important to know that Kafka will not attempt to automatically redistribute data in any way. So,
        this onus is also on you, too.</li>
</ul>
<ul id="724bf185-22ce-4ce5-a796-986961796e9e" class="bulleted-list">
    <li>There is currently no support for reducing the number of partitions for a topic.</li>
</ul>
<h2 id="cc2239f7-0c89-4dd1-bdcb-d9bde4f1cf3a" class="">Semantics of Data Producing</h2>
<p id="b45d7cab-9342-48be-b927-c6ea20b9508c" class="">On the data producing side, we need to know the topic
    name and the approach we need to use to distribute data across partitions (which is likely that your
    client will help on this with some out-of-the-box strategies, such as round-robin <a
        href="https://docs.confluent.io/4.0.0/clients/producer.html">as guaranteed by Confluent
        clients</a>). Apart from this, we have quite a few <a
        href="https://kafka.apache.org/documentation.html#producerconfigs">producer level configuration</a>
    we can apply to influence the semantics of data publishing.</p>
<p id="751e8cc8-4b69-4b51-8df3-f4b2d8dc6ed7" class="">When I am working with messaging systems, the first
    thing I want to understand is how the message delivery and durability guarantees are influenced, and
    what the default behaviour is for these. In Kafka, I found that this story a bit more confusing that it
    should probably be, which is due to a few configuration settings to be aligned to make it work in favour
    of durability to prevent message loss. Here are some important configuration for this:</p>
<ul id="b9d8ed44-c468-4981-9b0b-13c56f09f36b" class="bulleted-list">
    <li><code><a href="https://kafka.apache.org/documentation.html#acks">acks</a></code>: This setting
        indicates the number of acknowledgments the producer requires for a message publishing to be deemed
        as successful. It can be set to <code>0</code>, meaning that the producer won&#x27;t require an ack
        from any of the servers and this won&#x27;t give us any guarantees that the message is received by
        the server. This option could be preferable for cases where we need high throughput at the producing
        side and the data loss is not critical (e.g. sensor data, where losing a few seconds of data from a
        source won&#x27;t spoil our world). For cases where record durability is important, this can be set
        to <code>all</code>. This means the leader will wait for the full set of in-sync replicas to
        acknowledge the record, where the minimum number of required in-sync replicas is configured
        separately.</li>
</ul>
<ul id="d4459335-5203-457b-a408-266ebd92f5c1" class="bulleted-list">
    <li><code><a href="https://kafka.apache.org/documentation.html#min.insync.replicas">min.insync.replicas</a></code>:
        Quoting from the doc directly: &quot;When a producer sets <code>acks</code> to
        &quot;<code>all</code>&quot; (or &quot;<code>-1</code>&quot;), min.insync.replicas specifies the
        minimum number of replicas that must acknowledge a write for the write to be considered
        successful&quot;. This setting is topic level but can also be specified at the broker level. Setting
        this to the correct amount is really important and it&#x27;s set to <code>1</code> by default, which
        is probably not what you want if you care about durability of your messages and you have replication
        factor of <code>&gt;3</code> for the topic.</li>
</ul>
<ul id="a6dc5996-8065-494b-a542-a00fc6b34f2c" class="bulleted-list">
    <li><code><a href="https://kafka.apache.org/documentation.html#flush.messages">flush.messages</a></code>:
        In Kafka, messages are immediately written to the filesystem but by default we only fsync() to sync
        the OS cache lazily. This means that even if we have set our acks and min.insync.replicas to
        optimise for durability, there is still a theoretical chance that we can lose data with this
        behaviour. I explicitly said &quot;theoretical&quot; here as it&#x27;s quite unlikely to lose data
        with appropriate settings to rely on replication for data durability. For instance, with
        <code>acks=all</code> and <code>min.insync.replicas=2</code> settings for a topic which has
        replication factor of 3, we would be losing data after seeing a data write as successfull in cases
        of 3 machines (1 leader and 2 replicas) to fail at the same time before having a chance to flush
        that particular record to the disk, which is pretty unlikely, and this is why Kafka doesn&#x27;t
        recommend setting this value as well as <a
            href="https://kafka.apache.org/documentation.html#flush.ms"><code>flush.ms</code></a> value. So,
        we need to think a bit harder before setting these configuration values as this has some trade-offs
        to be thought about:<ul id="72970d7c-7179-41c6-8b0f-f510c324ebb2" class="bulleted-list">
            <li><strong>Durability</strong>: Unflushed data may be lost if you are not using replication.
            </li>
        </ul>
        <ul id="b6835546-07bf-4311-a2da-824c6f980acf" class="bulleted-list">
            <li><strong>Latency</strong>: Very large flush intervals may lead to latency spikes when the
                flush does occur as there will be a lot of data to flush.</li>
        </ul>
        <ul id="7e7571ab-c4d6-4a13-94f7-35e20b5c9cef" class="bulleted-list">
            <li><strong>Throughput</strong>: The flush is generally the most expensive operation, and a
                small flush interval may lead to exceessive seeks.</li>
        </ul>
    </li>
</ul>
<p id="6ab2286f-f5db-494b-9cb4-51680e9122cf" class="">So, a lot to think about here just to get message
    durability right. The good side of this complexity here is that Kafka is not trying to provide one way
    to solve all problems, which is not really possible especially when you want to optimise against
    different aspects (e.g. durability, throughput, etc.) depending the problem at hand. There is some
    further information on <a href="https://kafka.apache.org/documentation.html#semantics">message delivery
        guarantees in Kafka documentation</a>.</p>
<p id="945a6992-f1cf-4c81-8e36-d49532ced8ad" class="">There are some other producer semantics that requires
    understanding since the consequences of not understanding these might be costly depending on your needs.
    For example, <a href="https://kafka.apache.org/documentation.html#semantics">producer retries</a> is
    really important to understand correctly as this will have impact on message ordering even within a
    single partition. Another one is the <a
        href="https://kafka.apache.org/documentation.html#batch.size">batch size configuration</a>, which
    influences how many records to batch into one request whenever multiple records are being sent to the
    same partition. This might mean that <a
        href="https://kafka.apache.org/documentation.html#design_asyncsend">the sends will be performed
        asynchronously</a> and it may not be suitable for your needs. Finally, the <a
        href="https://kafka.apache.org/documentation.html#compaction">log compaction</a> is another concept
    which can be really useful to have a prior knowledge on, especially for cases where you publish the
    current state of an entity to a topic instead of publishing fine-grained events.</p>
<h2 id="8a618d84-a3bf-4e35-98e3-03698a63a312" class="">Resources</h2>
<ul id="f875d8a2-fd2d-46b6-899f-787842adafda" class="bulleted-list">
    <li><a href="https://medium.com/better-programming/kafka-acks-explained-c0515b3b707e">Kafka Acks
            Explained</a></li>
</ul>
<ul id="14c2c4c5-4fff-4a15-97ac-3c72ae20663f" class="bulleted-list">
    <li><a href="https://link.medium.com/sfm4jiKbL6">Why fsync is bad for Kafka</a></li>
</ul>
<ul id="75cfc68b-ddad-4541-8970-f45ac94c9071" class="bulleted-list">
    <li><a
            href="https://users.kafka.apache.narkive.com/UUQx7UcG/does-kafka-send-the-acks-response-to-the-producer-after-flush-the-messages-to-the-disk-or-just-keep-">Does
            kafka send the acks response to the producer after flush the messages to the disk or just keep
            them in the memory</a></li>
</ul>
<ul id="50279304-d760-4bdf-8607-6012628be5a2" class="bulleted-list">
    <li><a
            href="https://stackoverflow.com/questions/57987591/can-a-message-loss-occur-in-kafka-even-if-producer-gets-acknowledgement-for-it">Can
            a message loss occur in Kafka even if producer gets acknowledgement for it?</a></li>
</ul>