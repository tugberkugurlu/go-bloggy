---
id: c7de3200-6d07-498f-a7d1-9ad47d6d770b
title: Order of Fields Matters on MongoDB Indexes
abstract: Order of Fields Matters on MongoDB Indexes. Let's see how with an example.
created_at: 2014-04-12 18:52:00 +0000 UTC
tags:
- .net
- JavaScript
- MongoDB
slugs:
- order-of-fields-matters-on-mongodb-indexes
---

<p>As my <a href="http://mongodb.org">MongoDB</a> journey continues, I discover new stuff along the way and one of them is about <a href="http://docs.mongodb.org/manual/indexes/">indexes in MongoDB</a>. Let me try to explain it with a sample.</p> <p>First, create the below four documents inside our users collection:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>db.users.save({ 
	_id: 1, 
	name: <span style="color: #a31515">"tugberk1"</span>, 
	login: [
		{ProviderName: <span style="color: #a31515">"twitter"</span>, ProviderKey: <span style="color: #a31515">"232"</span>}, 
		{ProviderName: <span style="color: #a31515">"facebook"</span>, ProviderKey: <span style="color: #a31515">"423"</span>}
	]
});
	
db.users.save({ 
	_id: 2, 
	name: <span style="color: #a31515">"tugberk23"</span>, 
	login: [
		{ProviderName: <span style="color: #a31515">"twitter"</span>, ProviderKey: <span style="color: #a31515">"3443"</span>}
	]
});

db.users.save({ 
	_id: 3, 
	name: <span style="color: #a31515">"tugberk4343"</span>, 
	login: [
		{ProviderName: <span style="color: #a31515">"dropbox"</span>, ProviderKey: <span style="color: #a31515">"445345"</span>}
	]
});

db.users.save({ 
	_id: 4, 
	name: <span style="color: #a31515">"tugberk98"</span>, 
	login: [
		{ProviderName: <span style="color: #a31515">"dropbox"</span>, ProviderKey: <span style="color: #a31515">"3443"</span>}, 
		{ProviderName: <span style="color: #a31515">"facebook"</span>, ProviderKey: <span style="color: #a31515">"768"</span>}
	]
});</pre></div></div>
<p>Let’s query the users collection by login.ProviderKey and login.ProviderName:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>db.users.find({
	<span style="color: #a31515">"login.ProviderKey"</span>: <span style="color: #a31515">"232"</span>, 
	<span style="color: #a31515">"login.ProviderName"</span>: <span style="color: #a31515">"twitter"</span>
}).pretty();</pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/acf9f3f2-a69a-4ac4-a427-95d47e73fb19.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b07208e4-af75-4750-b2de-53d82d46776f.png" width="644" height="306"></a></p>
<p>It found the document we wanted. Let’s see how it performed:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>db.users.find({
	<span style="color: #a31515">"login.ProviderKey"</span>: <span style="color: #a31515">"232"</span>, 
	<span style="color: #a31515">"login.ProviderName"</span>: <span style="color: #a31515">"twitter"</span>
}).explain();</pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5ed832a0-b5ca-4c52-84dc-89da3aa50129.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/570862b5-61d4-4ba8-be31-c8dab40d2e66.png" width="644" height="306"></a></p>
<p>Result is actually pretty bad. It scanned all four documents to find the one that we wanted to get. Let’s put an index to ProviderName and ProviderKey fields:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>db.users.ensureIndex({
	<span style="color: #a31515">"login.ProviderName"</span>: 1, 
	<span style="color: #a31515">"login.ProviderKey"</span>: 1
});</pre></div></div>
<p>Now, let’s see how it performs the query:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/9c9bba37-e1bc-4def-a934-b3922a09306a.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/02257a87-9c87-4fdf-8ed5-f888da80fc4f.png" width="644" height="476"></a></p>
<p>It’s better as it scanned only two documents. However, we had only one matching document for our query. As the chances that the providerKey will be more unique than the ProviderName, I want it to first look for the ProviderKey. To do that, I need to change the index:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>db.users.dropIndex({
	<span style="color: #a31515">"login.ProviderName"</span>: 1, 
	<span style="color: #a31515">"login.ProviderKey"</span>: 1 
});
db.users.ensureIndex({ 
	<span style="color: #a31515">"login.ProviderKey"</span>: 1, 
	<span style="color: #a31515">"login.ProviderName"</span>: 1 
});</pre></div></div>
<p>Let’s now see how it’s trying to find the matching documents:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>db.users.find({
	<span style="color: #a31515">"login.ProviderKey"</span>: <span style="color: #a31515">"232"</span>, 
	<span style="color: #a31515">"login.ProviderName"</span>: <span style="color: #a31515">"twitter"</span>
}).explain();</pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/244ca85a-9d88-4f65-a8d5-0b057102714d.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/f87bc4b9-e732-4dce-a1ef-61dd89166d61.png" width="644" height="476"></a></p>
<p>Boom! Only one document was scanned. This shows us how it’s important to put the fields in right order for our queries.</p>
<h3>Resources</h3>
<ul>
<li><a href="http://emptysqua.re/blog/optimizing-mongodb-compound-indexes/">Optimizing MongoDB Compound Indexes</a></li>
<li><a href="http://docs.mongodb.org/manual/core/indexes-introduction/">Index Introduction</a></li></ul>  