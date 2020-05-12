---
id: 3c3c0733-f919-403d-8c43-6216cd7d8966
title: Getting Started with Neo4j in .NET with Neo4jClient Library
abstract: I have been looking into Neo4j, a graph database, for a while and here is
  what impressed me the most while trying to work with it through the Neo4jClient
  .NET library.
created_at: 2015-12-13 19:07:00 +0000 UTC
tags:
- Neo4j
slugs:
- getting-started-with-neo4j-in--net-with-neo4jclient-library
---

<p>I am really in love with the side project I am working on now. It is broken down to little "micro" applications (a.k.a. <a href="http://martinfowler.com/articles/microservices.html">microservices</a>), uses multiple data storage technologies and being brought together through <a href="http://www.tugberkugurlu.com/archive/playing-around-with-docker-hello-world-development-environment-and-your-application">Docker</a>. As a result, the entire solution feels very natural, not-restricted and feels so manageable.</p> <p>One part of this solution requires to answer a question which involves going very deep inside the data hierarchy. To illustrate what I mean, have a look at the below graph:</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/f92ca593-a571-4208-8517-e7fbfbd48472.jpg"><img title="movies-with-only-agency-employees-2" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="movies-with-only-agency-employees-2" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/7cf13573-fbd3-4d56-ba64-0f048d7a001e.jpg" width="644" height="286"></a></p> <p>Here, we have an agency which has acquired some actors. Also, we have some movies which employed some actors. You can model this in various data storage systems in various ways but the question I want to answer is the following: "What are the movies which employed all of its actors from Agency-A?". Even thinking about the query you would write in T-SQL is enough to melt your brain for this one. It doesn’t mean that SQL Server, MySQL, etc. are bad data storage systems. It’s just that this type of questions are not among those data storage systems' strengths.</p> <h3>Enters: Neo4j</h3> <p><a href="http://neo4j.com/">Neo4j</a> is an open-source <a href="https://en.wikipedia.org/wiki/Graph_database">graph database</a> implemented in Java and accessible from software written in other languages using the Cypher query language through a transactional HTTP endpoint (<a href="https://en.wikipedia.org/wiki/Neo4j">Wikipedia says</a>). In Neo4j, your data set consists of nodes and relationships between these nodes which you can interact with through the <a href="http://neo4j.com/developer/cypher-query-language/">Cypher query language</a>. Cypher is a very powerful, declarative, SQL-inspired language for describing patterns in graphs. The biggest thing that stands out when working with Cypher is the <a href="http://neo4j.com/docs/stable/query-match.html#_relationship_basics">relationships</a>. Relationships are first class citizens in Cypher. Consider the following Cypher query which is brought from the movie sample in Neo4j web client:</p> <blockquote> <p>You can bring up the this movie sample by just running ":play movie graph" from the Neo4j web client and walk through it.</p></blockquote> <div class="code-wrapper border-shadow-1" style="color: black; background-color: white"><pre>MATCH (tom:Person {<span style="color: blue">name</span>: <span style="color: #a31515">"Tom Hanks"</span>})-[:ACTED_IN]-&gt;(tomHanksMovies) <span style="color: blue">RETURN</span> tom,tomHanksMovies</pre></div>
<p>This will list all Tom Hanks movies. However, when you read it from left to right, you will pretty much understand what it will do anyway. The interesting part here is the ACTED_IN relationship inside the query. You may think at this point that this is not a big deal as it can probably translate the below T-SQL query:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">SELECT</span> * <span style="color: blue">FROM</span> Movies m
<span style="color: blue">INNER</span> <span style="color: blue">JOIN</span> MovieActors ma <span style="color: blue">ON</span> ma.MovieId = m.Id
<span style="color: blue">WHERE</span> ma.ActorId = 1;</pre></div></div>
<p>However, you will start seeing the power as the questions get interesting. For example, let’s find out Tom Hanks’ co-actors from the every movie he acted in (again, from the same sample):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>MATCH (tom:Person {<span style="color: blue">name</span>:<span style="color: #a31515">"Tom Hanks"</span>})-[:ACTED_IN]-&gt;(m)&lt;-[:ACTED_IN]-(coActors) <span style="color: blue">RETURN</span> coActors.name</pre></div></div>
<p>It’s just mind-blowingly complicated to retrieve this from a relational database but with Cypher, it is dead easy. You can start to see that it’s all about building up nodes and declaring the relationships to get the answer to your question in Neo4j.</p>
<h2>Neo4j in .NET</h2>
<p>As Neo4j communicates through HTTP, you can pretty much find a client implementation in every ecosystem and .NET is not an exception. Amazing people from <a href="https://github.com/Readify">Readify</a> is maintaining <a href="https://github.com/Readify/Neo4jClient">Neo4jClient</a> OSS project. It’s extremely easy to use this and <a href="https://github.com/Readify/Neo4jClient/wiki">the library has a very good documentation</a>. I especially liked the part where they have <a href="https://github.com/Readify/Neo4jClient/wiki/connecting#threading-and-lifestyles">documented the thread safety concerns of GraphClient</a>. It is the first thing I wanted to find out and there it was.</p>
<p>Going back to my example which I mentioned at the beginning of this post, I tried to handle this through the .NET Client. Let’s walk through what I did.</p>
<blockquote>
<p>You can find <a href="https://github.com/tugberkugurlu/DotNetSamples/blob/7a15226fbcc883effc05537b59a203d71a51c490/DotNetSpecific/Neo4jClientSample/Neo4jClientSample">the below sample under my DotNetSamples GitHub repository</a>.</p></blockquote>
<p>First, I initiated the GraphClient and made some adjustments:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">var</span> client = <span style="color: blue">new</span> GraphClient(<span style="color: blue">new</span> Uri(<span style="color: #a31515">"http://localhost:7474/db/data"</span>), <span style="color: #a31515">"neo4j"</span>, <span style="color: #a31515">"1234567890"</span>)
{
    JsonContractResolver = <span style="color: blue">new</span> CamelCasePropertyNamesContractResolver()
};

client.Connect();</pre></div></div>
<p>I started with creating the agency.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">var</span> agencyA = <span style="color: blue">new</span> Agency { Name = <span style="color: #a31515">"Agency-A"</span> };
client.Cypher
    .Create(<span style="color: #a31515">"(agency:Agency {agencyA})"</span>)
    .WithParam(<span style="color: #a31515">"agencyA"</span>, agencyA)
    .ExecuteWithoutResultsAsync()
    .Wait();</pre></div></div>
<p>Next is to <a href="https://github.com/Readify/Neo4jClient/wiki/cypher-examples#create-a-user-and-relate-them-to-an-existing-one">create the actors and ACQUIRED relationship between the agency</a> and some actors (in below case, only the odd numbered actors):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">for</span> (<span style="color: blue">int</span> i = 1; i &lt;= 5; i++)
{
    <span style="color: blue">var</span> actor = <span style="color: blue">new</span> Person { Name = $<span style="color: #a31515">"Actor-{i}"</span> };

    <span style="color: blue">if</span> ((i % 2) == 0)
    {
        client.Cypher
            .Create(<span style="color: #a31515">"(actor:Person {newActor})"</span>)
            .WithParam(<span style="color: #a31515">"newActor"</span>, actor)
            .ExecuteWithoutResultsAsync()
            .Wait();
    }
    <span style="color: blue">else</span>
    {
        client.Cypher
            .Match(<span style="color: #a31515">"(agency:Agency)"</span>)
            .Where((Agency agency) =&gt; agency.Name == agencyA.Name)
            .Create(<span style="color: #a31515">"agency-[:ACQUIRED]-&gt;(actor:Person {newActor})"</span>)
            .WithParam(<span style="color: #a31515">"newActor"</span>, actor)
            .ExecuteWithoutResultsAsync()
            .Wait();
    }
}</pre></div></div>
<p>Then, I have created the movies :</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">char</span>[] chars = Enumerable.Range(<span style="color: #a31515">'a'</span>, <span style="color: #a31515">'z'</span> - <span style="color: #a31515">'a'</span> + 1).Select(i =&gt; (Char)i).ToArray();
<span style="color: blue">for</span> (<span style="color: blue">int</span> i = 0; i &lt; 3; i++)
{
    <span style="color: blue">var</span> movie = <span style="color: blue">new</span> Movie { Name = $<span style="color: #a31515">"Movie-{chars[i]}"</span> };

    client.Cypher
        .Create(<span style="color: #a31515">"(movie:Movie {newMovie})"</span>)
        .WithParam(<span style="color: #a31515">"newMovie"</span>, movie)
        .ExecuteWithoutResultsAsync()
        .Wait();
}</pre></div></div>
<p>Lastly, I have related existing movies and actors through the EMPLOYED relationship.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>client.Cypher
    .Match(<span style="color: #a31515">"(movie:Movie)"</span>, <span style="color: #a31515">"(actor1:Person)"</span>, <span style="color: #a31515">"(actor5:Person)"</span>)
    .Where((Movie movie) =&gt; movie.Name == <span style="color: #a31515">"Movie-a"</span>)
    .AndWhere((Person actor1) =&gt; actor1.Name == <span style="color: #a31515">"Actor-1"</span>)
    .AndWhere((Person actor5) =&gt; actor5.Name == <span style="color: #a31515">"Actor-5"</span>)
    .Create(<span style="color: #a31515">"(movie)-[:EMPLOYED]-&gt;(actor1), (movie)-[:EMPLOYED]-&gt;(actor5)"</span>)
    .ExecuteWithoutResultsAsync()
    .Wait();

client.Cypher
    .Match(<span style="color: #a31515">"(movie:Movie)"</span>, <span style="color: #a31515">"(actor1:Person)"</span>, <span style="color: #a31515">"(actor3:Person)"</span>, <span style="color: #a31515">"(actor5:Person)"</span>)
    .Where((Movie movie) =&gt; movie.Name == <span style="color: #a31515">"Movie-b"</span>)
    .AndWhere((Person actor1) =&gt; actor1.Name == <span style="color: #a31515">"Actor-1"</span>)
    .AndWhere((Person actor3) =&gt; actor3.Name == <span style="color: #a31515">"Actor-3"</span>)
    .AndWhere((Person actor5) =&gt; actor5.Name == <span style="color: #a31515">"Actor-5"</span>)
    .Create(<span style="color: #a31515">"(movie)-[:EMPLOYED]-&gt;(actor1), (movie)-[:EMPLOYED]-&gt;(actor3), (movie)-[:EMPLOYED]-&gt;(actor5)"</span>)
    .ExecuteWithoutResultsAsync()
    .Wait();

client.Cypher
    .Match(<span style="color: #a31515">"(movie:Movie)"</span>, <span style="color: #a31515">"(actor2:Person)"</span>, <span style="color: #a31515">"(actor5:Person)"</span>)
    .Where((Movie movie) =&gt; movie.Name == <span style="color: #a31515">"Movie-c"</span>)
    .AndWhere((Person actor2) =&gt; actor2.Name == <span style="color: #a31515">"Actor-2"</span>)
    .AndWhere((Person actor5) =&gt; actor5.Name == <span style="color: #a31515">"Actor-5"</span>)
    .Create(<span style="color: #a31515">"(movie)-[:EMPLOYED]-&gt;(actor2), (movie)-[:EMPLOYED]-&gt;(actor5)"</span>)
    .ExecuteWithoutResultsAsync()
    .Wait();</pre></div></div>
<p>When I run this, I now have the data set that I can play with. I have jumped back to web client and ran the below query to retrieve the relations:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>MATCH (agency:Agency)-[:ACQUIRED]-&gt;(actor:Person)&lt;-[:EMPLOYED]-(movie:Movie)
<span style="color: blue">RETURN</span> agency, actor, movie</pre></div></div>
<p>One of the greatest features of the web client is that you can view your query result in a graph representation. How cool is that? You can exactly see the smilarity between the below result and the graph I have put together above:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1e121e36-1824-4988-b752-cc8cdc1fb0b8.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/2026169c-966e-4225-b89f-e755d76f279f.png" width="644" height="308"></a></p>
<p>Of course, we can run the same above query through the .NET client and grab the results:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">var</span> results = client.Cypher
    .Match(<span style="color: #a31515">"(agency:Agency)-[:ACQUIRED]-&gt;(actor:Person)&lt;-[:EMPLOYED]-(movie:Movie)"</span>)
    .Return((agency, actor, movie) =&gt; <span style="color: blue">new</span>
    {
        Agency = agency.As&lt;Agency&gt;(),
        Actor = actor.As&lt;Person&gt;(),
        Movie = movie.As&lt;Movie&gt;()
    }).Results;</pre></div></div>
<h3>Going Beyond</h3>
<p>However, how can we answer my "What are the movies which employed all of its actors from Agency-A?" question? As I am very new to Neo4j, I struggled a lot with this. In fact, I was not even sure whether this was possible to do in Neo4J. <a href="http://stackoverflow.com/questions/34233230/neo4j-and-multiple-items-nested-relationships">I asked this as a question in Stackoverflow</a> (as every stuck developer do) and <a href="http://stackoverflow.com/users/2662355/christophe-willemsen">Christophe Willemsen</a> gave an amazing answer which literally blew my mind. I warn you now as the below query is a bit complex and I am still going through it piece by piece to try to understand it but it does the trick:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>MATCH (agency:Agency { <span style="color: blue">name</span>:<span style="color: #a31515">"Agency-A"</span> })-[:ACQUIRED]-&gt;(actor:Person)&lt;-[:EMPLOYED]-(movie:Movie)
<span style="color: blue">WITH</span> <span style="color: blue">DISTINCT</span> movie, collect(actor) <span style="color: blue">AS</span> actors
MATCH (movie)-[:EMPLOYED]-&gt;(allemployees:Person)
<span style="color: blue">WITH</span> movie, actors, <span style="color: magenta">count</span>(allemployees) <span style="color: blue">AS</span> c
<span style="color: blue">WHERE</span> c = <span style="color: blue">size</span>(actors)
<span style="color: blue">RETURN</span> movie.name</pre></div></div>

<p>The result is as you would expect:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/117e9cec-1108-4735-892a-84c530c82eb2.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c93fe936-e36d-4d94-af86-16e0d905583c.png" width="644" height="243"></a></p>
<h3>Still Dipping My Toes</h3>
<p>I am hooked but this doesn’t mean that Neo4j is the solution to my problems. I am still evaluating it by implementing a few features on top of it. There are a few parts which I haven’t been able to answer exactly yet:</p>
<ul>
<li>How does this scale with large data sets?</li>
<li>Can I shard the data across servers?</li>
<li>Want are the hosted options?</li>
<li>What is the story on geo location queries?</li></ul>
<p>However, the architecture I have in my solution allows me to evaluate this type of technologies. At worst case scenario, Neo4j will not work for me but I will be able to replace it with something else (which I doubt that it will be the case).</p>
<h3>Resources</h3>
<ul>
<li><a href="http://neo4j.com/developer/cypher-query-language/">Intro to Cypher</a></li>
<li><a href="http://graphdatabases.com/">Free e-book: The Definitive Book on Graph Databases</a></li>
<li><a href="http://neo4j.com/docs/stable/query-aggregation.html">Cypher: Aggregation</a></li>
<li><a href="http://neo4j.com/docs/stable/query-functions-mathematical.html">Cypher:&nbsp; Mathematical functions</a></li>
<li><a href="http://neo4j.com/docs/stable/performance-guide.html">Neo4j Performance Guide</a></li>
<li><a href="http://neo4j.com/use-cases/real-time-recommendation-engine/">Use Case: Real-Time Recommendation Engine</a></li>
<li><a href="http://neo4j.com/use-cases/social-network/">Use Case: Social Network</a></li></ul>  