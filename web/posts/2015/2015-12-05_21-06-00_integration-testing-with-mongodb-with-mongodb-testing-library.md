---
id: 72dbfc8d-a117-4a95-a44f-0d3b4c8412d4
title: Integration Testing with MongoDB with MongoDB.Testing Library
abstract: I have put together a library, MongoDB.Testing, which makes it easy to stand
  up a MongoDB server, create a random database and clean up the resources afterwards.
  Here is how you can start using it.
created_at: 2015-12-05 21:06:00 +0000 UTC
tags:
- .net
- MongoDB
slugs:
- integration-testing-with-mongodb-with-mongodb-testing-library
---

<p>Considering the applications we produce today (small, targeted, "micro" applications), I value <a href="https://en.wikipedia.org/wiki/Integration_testing">integration tests</a> way more than unit tests (along with acceptance tests). They provide much more realistic testing on your application with the only downside of being hard to pinpoint which part of your code is the problem when you have failures. I have been writing integration tests for the .NET based HTTP applications which use <a href="https://www.mongodb.org/">MongoDB</a> as the data storage system on same parts and I pulled out a helper into library which makes it easy to stand up a MongoDB server, create a random database and clean up the resources afterwards. The library is called MongoDB.Testing and it’s on <a href="https://www.nuget.org/packages/MongoDB.Testing">NuGet</a>, <a href="https://github.com/tugberkugurlu/MongoDB.Testing">GitHub</a>. Usage is also pretty simple and there is also a <a href="https://github.com/tugberkugurlu/MongoDB.Testing/tree/master/samples">a few samples I have put together</a>.</p> <p>Install the library into your testing project through NuGet:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>Install<span style="color: gray">-</span>Package MongoDB.Testing <span style="color: gray">-</span>pre</pre></div></div>
<p>Write a mongod.exe locator:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">class</span> MongodExeLocator : IMongoExeLocator
{
    <span style="color: blue">public</span> <span style="color: blue">string</span> Locate()
    {
        <span style="color: blue">return</span> <span style="color: #a31515">@"C:\Program Files\MongoDB\Server\3.0\bin\mongod.exe"</span>;
    }
}</pre></div></div>
<p>Finally, integrate this into your tests:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>[Test]
<span style="color: blue">public</span> async Task HasEnoughRating_Should_Throw_When_The_User_Is_Not_Found()
{
    <span style="color: blue">using</span> (MongoTestServer server = MongoTestServer.Start(27017, <span style="color: blue">new</span> MongodExeLocator()))
    {
        <span style="color: green">// ARRANGE</span>
        <span style="color: blue">var</span> collection = server.Database.GetCollection&lt;UserEntity&gt;(<span style="color: #a31515">"users"</span>);
        <span style="color: blue">var</span> service = <span style="color: blue">new</span> MyCounterService(collection);
        await collection.InsertOneAsync(<span style="color: blue">new</span> UserEntity
        {
            Id = ObjectId.GenerateNewId().ToString(),
            Name = <span style="color: #a31515">"foo"</span>,
            Rating = 23
        });

        <span style="color: green">// ACT, ASSERT</span>
        Assert.Throws&lt;InvalidOperationException&gt;(
            () =&gt; service.HasEnoughRating(ObjectId.GenerateNewId().ToString()));
    }
}</pre></div></div>
<p>That’s basically all. <a href="https://github.com/tugberkugurlu/MongoDB.Testing/blob/1.0.0-beta-001/src/MongoDB.Testing/Mongo/MongoTestServer.cs#L93-L96">MongoTestServer.Start</a> will do the following for you:</p>
<ul>
<li>Start a mongod instance and expose it through the specified port.</li>
<li>Creates a randomly named MongoDB database on the started instance and exposes it through the MongoTestServer instance returned from MongoTestServer.Start method.</li>
<li>Cleans up the resources, kills the mongod.exe instance when the MongoTestServer instance is disposed.</li></ul>
<p>If you are doing a similar sort of testing with MongoDB, give this a shot. I want to improve this based on the needs. So, make sure to file <a href="https://github.com/tugberkugurlu/MongoDB.Testing/issues">issues</a> and send some lovely <a href="https://github.com/tugberkugurlu/MongoDB.Testing/pulls">pull requests</a>.</p>  