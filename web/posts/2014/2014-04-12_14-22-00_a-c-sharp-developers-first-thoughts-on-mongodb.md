---
id: edc9a740-d18e-4606-b204-8a786a12f7dc
title: A C# Developer's First Thoughts on MongoDB
abstract: After working with RavenDB over the year, I just started looking into MongoDB.
  I worked with MongoDB a year ago or so in a small project but my knowledge was mostly
  rusty and I don't want that to happen again :) So, here I'm, documenting what my
  second thoughts are :)
created_at: 2014-04-12 14:22:00 +0000 UTC
tags:
- .net
- JavaScript
- MongoDB
slugs:
- a-c-sharp-developers-first-thoughts-on-mongodb
---

<p>After working with <a href="https://ravendb.net/">RavenDB</a> over the year, I just started looking into <a href="https://www.mongodb.org/&lrm;">MongoDB</a>. I worked with MongoDB a year ago or so in a small project but my knowledge was mostly rusty and I don't want that to happen again :) So, here I'm, documenting what my second thoughts are :) The TL;DR is that: I'm loving it but the lack of transaction sometimes drifts on a vast dark sea. It's OK through. The advantages generally overcomes this disadvantage.  <h3>Loved the Mongo Shell</h3> <p>First thing I liked about MongoDB is its shell (<a href="http://docs.mongodb.org/v2.2/mongo/">Mongo Shell</a>).&nbsp; <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ba742551-7c62-4f4c-8bcb-699381458a2b.png"><img title="3d6f5fdd4be53096b6992d5b84b0d7be" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="3d6f5fdd4be53096b6992d5b84b0d7be" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a66f4f13-1cfc-465a-9bd2-ef5a0f13e877.png" width="644" height="328"></a>  <p>It makes it insanely easy for you to get used to MongoDB. After running the mongod.exe, I fired up a command prompt and navigated to mongo.exe directory and entered the mongo shell. Mongo Shell runs pure JavaScript code. That's right! Anything you know about JavaScript is completely valid inside the mongo shell. Let's see a few things that you can do with mongo shell.  <p>You can list the databases on your server: show dbs  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a24de0b0-cc2c-4e8f-9620-f1b68c7ba269.png"><img title="2" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="2" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6c9c6cf5-58bf-4173-a346-176a1fbefc11.png" width="644" height="192"></a>  <p>You can see which database you are connected to: db  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a887a7cb-56a6-4872-a52c-3afb2e456a31.png"><img title="3" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="3" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/21c16677-8c88-4274-b1e0-63a6506cbec9.png" width="644" height="192"></a>  <p>You can switch to a different database: use &lt;database name here&gt;  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b3047f21-30e7-4f40-8543-8c2233de8e3e.png"><img title="4" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="4" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ab51adb6-5501-469f-a68e-26b231f5680c.png" width="644" height="192"></a>  <p>You can see the collections inside the database you are on: show collections  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c9217a99-4c69-47d3-b2de-2a7d94c693d0.png"><img title="5" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="5" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/1f934369-a96b-4804-b6cc-10eea73d714b.png" width="644" height="192"></a>  <p>You can save a document inside a collection: db.users.save({ _id: "tugberk", userName: "Tugberk" })  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/50040953-5da7-4712-aa0f-7e5874ad9f0a.png"><img title="6" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="6" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/91dda03b-2e75-480f-9971-1fa667d3e9f1.png" width="644" height="192"></a>  <p>You can list the documents inside a collection: db.users.find().pretty()  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/72cd9050-4a18-45ef-b899-d2585012718b.png"><img title="7" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="7" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e5045864-9a8b-48fa-8191-da665be83dfc.png" width="644" height="192"></a>  <p>You can run a for loop:  <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre><span style="color: blue">for</span>(<span style="color: blue">var</span> i = 0; i &lt; 10; i++) { 
	db.users.save({ 
		_id: <span style="color: #a31515">"tugberk"</span> + i.toString(), 
		userName: <span style="color: #a31515">"Tugberk"</span> + i.toString() 
	}) 
}</pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/49d65187-adf1-4749-8cc4-e966fd356341.png"><img title="8" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="8" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6bafbaac-e901-4c04-aea6-5f032d497a35.png" width="644" height="192"></a> 
<p>You can <a href="http://docs.mongodb.org/manual/tutorial/write-scripts-for-the-mongo-shell/">run the code inside a js file</a>:&nbsp; <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ed426bc4-8fd1-4295-aa5a-f20eddec8158.png"><img title="9" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="9" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/3775db78-7ed5-4f4b-9efd-943a4284318e.png" width="644" height="192"></a> 
<p>saveCount.js contains the following code and it just gets the count of documents inside the users collection and logs it inside another collection: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>(<span style="color: blue">function</span>() {
     <span style="color: blue">var</span> myDb = db.getSiblingDB(<span style="color: #a31515">'myDatabase'</span>),
         usersCount = myDb.users.find().count();
         
     myDb.countLogs.save({
          count: usersCount,
          createdOn: <span style="color: blue">new</span> Date()
     });
}());</pre></div></div>
<p>All of those and more can be done using the mongo shell. It's a full blown MongoDB client and probably the best one. I don't know if it's just me but big missing feature of RavenDB is this type of shell. 
<h3>Loved the Updates</h3>
<p><a href="http://docs.mongodb.org/manual/reference/method/db.collection.update/">Update operations in MongoDB</a> is just fabulous. You can construct many kinds of updates, the engine allows you to do this fairly easily. The one I found most useful is the increment updates. Increment updates allows you to increment a field and this operation will be performed concurrency in-mind: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>db.books.update(
   { item: <span style="color: #a31515">"Divine Comedy"</span> },
   {
      $inc: { stock: 5 }
   }
)</pre></div></div>
<p>The above query will update the stock filed by 5 safely. 
<h3>Not much Love for the .NET Client</h3>
<p><a href="http://docs.mongodb.org/ecosystem/drivers/csharp/">MongoDB has an official .NET client</a> but MongoDB guys decided to call this "C# driver". This is so strange because it works with any other .NET languages as well. I have to say that MongoDB .NET client is not so great in my opinion. After coming from the RavenDB .NET client, using the MongoDB .NET client just feels uncomfortable (However, I’m most certainly sure that I’d love its Node.Js client as it would feel very natural). 
<p>First of all, it doesn't support asynchronous requests to MongoDB server. All TCP requests are being done synchronously. Also, there is no embedded server support. RavenDB has this and it makes testing a joy. Let's look at the below code which gets an instance of a database: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>MongoClient client = <span style="color: blue">new</span> MongoClient(<span style="color: #a31515">"mongodb://localhost"</span>);
MongoServer server = client.GetServer();
MongoDatabase db = server.GetDatabase(<span style="color: #a31515">"mongodemo"</span>);</pre></div></div>
<p>There is too much noise going on here. For example, what is GetServer method there? Instead, I would want to see something like below: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>MongoClient client = <span style="color: blue">new</span> MongoClient(<span style="color: #a31515">"mongodb://localhost"</span>);
<span style="color: blue">using</span>(<span style="color: blue">var</span> session = client.OpenSession(<span style="color: #a31515">"myDatabase"</span>))
{
     <span style="color: green">// work with the session here...</span>
}</pre></div></div>
<p>Looks familiar :) I bet it does! Other than the above issues, creating map/reduce jobs just feels weird as well because MongoDB supports JavaScript to perform map/reduce operations. 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">var</span> map =
    <span style="color: #a31515">"function() {"</span> +
    <span style="color: #a31515">"    for (var key in this) {"</span> +
    <span style="color: #a31515">"        emit(key, { count : 1 });"</span> +
    <span style="color: #a31515">"    }"</span> +
    <span style="color: #a31515">"}"</span>;

<span style="color: blue">var</span> reduce =
    <span style="color: #a31515">"function(key, emits) {"</span> +
    <span style="color: #a31515">"    total = 0;"</span> +
    <span style="color: #a31515">"    for (var i in emits) {"</span> +
    <span style="color: #a31515">"        total += emits[i].count;"</span> +
    <span style="color: #a31515">"    }"</span> +
    <span style="color: #a31515">"    return { count : total };"</span> +
    <span style="color: #a31515">"}"</span>;

<span style="color: blue">var</span> mr = collection.MapReduce(map, reduce);
<span style="color: blue">foreach</span> (<span style="color: blue">var</span> document <span style="color: blue">in</span> mr.GetResults()) {
    Console.WriteLine(document.ToJson());
}</pre></div></div>
<p>The above code is directly taken from the <a href="http://docs.mongodb.org/ecosystem/tutorial/use-csharp-driver/#mapreduce-method">MongoDB documentation</a>. 
<h3>Explore Yourself</h3>
<p><a href="http://docs.mongodb.org/manual/">MongoDB has a nice documentation</a> and you can explore it yourself. Besides that, Pluralsight has a pretty nice course on MongoDB: <a href="http://www.pluralsight.com/training/Courses/TableOfContents/mongodb-introduction">Introduction to MongoDB</a> by Nuri Halperin. Also, don't miss the Ben's post on <a href="http://benfoster.io/blog/map-reduce-in-mongodb-and-ravendb">the comparison of Map-Reduce in MongoDB and RavenDB</a>.</p>  