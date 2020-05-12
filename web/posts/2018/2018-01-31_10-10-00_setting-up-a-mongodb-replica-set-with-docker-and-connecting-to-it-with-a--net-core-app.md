---
id: 55fa0e98-4a31-48e2-9473-919b82fb9232
title: Setting up a MongoDB Replica Set with Docker and Connecting to It With a .NET
  Core App
abstract: 'Easily setting up realistic non-production (e.g. dev, test, QA, etc.) environments
  is really critical in order to reduce the feedback loop. In this blog post, I want
  to talk about how you can achieve this if your application relies on MongoDB Replica
  Set by showing you how to set it up with Docker for non-production environments. '
created_at: 2018-01-31 10:10:00 +0000 UTC
tags:
- ASP.NET Core
- Docker
- MongoDB
slugs:
- setting-up-a-mongodb-replica-set-with-docker-and-connecting-to-it-with-a--net-core-app
---

<p>Easily setting up realistic non-production (e.g. dev, test, QA, etc.) environments is really critical in order to reduce the feedback loop. In this blog post, I want to talk about how you can achieve this if your application relies on MongoDB Replica Set by showing you how to set it up with Docker for non-production environments.  <blockquote> <h3>Hold on! I want to watch, not read!</h3> <p>I got you covered there! I have also recorded a ~5m covering the content of this blog post, where I also walks you through the steps visually. If you find this option useful, let me know through the comments below and I can aim harder to repeat that :)  <p><iframe height="315" src="https://www.youtube.com/embed/1lAjmJ1ht1o" frameborder="0" width="560" allowfullscreen allow="autoplay; encrypted-media"></iframe></p></blockquote> <h3>What are we trying to do here and why?</h3> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/67b0afbb-825b-4cd3-b1b0-522d2f6869f8.png"><img title="Picture2" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Picture2" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b5111dd7-6e7e-4246-a2ec-76421a7dfd98.png" width="640" height="222"></a></p> <p>&nbsp;</p> <p>If you have an application which works against a MongoDB database, it’s very common to have a replica set in production. This approach ensures the high availability of the data, especially for read scenarios. However, applications mostly end up working against a single MongoDB instance, because setting up a Replica Set in isolation is a tedious process. As mentioned at the beginning of the post, we want to reflect the production environment to the process of developing or testing the software applications as much as possible. The reason for that is to catch unexpected behaviour which may only occur under a production environment. This approach is valuable because it would allow us to reduce the feedback loop on those exceptional cases.</p> <h3>Docker makes this all easy!</h3> <p>This is where Docker enters into the picture! Docker is containerization technology and it allows us to have repeatable process to provision environments in a declarative way. It also gives us a try and tear down model where we can experiment and easily start again from the initial state. Docker can also help us with easily setting up a MongoDB Replica Set. Within our Docker Host, we can create Docker Network which would give us the isolated DNS resolution across containers. Then we can start creating the MongoDB docker containers. They would initially be unaware of each other. However, we can initialise the replication by connecting to one of the containers and running the replica set initialisation command. Finally,&nbsp; we can deploy our application container under the same docker network. <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/34165076-d303-419d-a7c1-505c5fc277ec.png"><img title="Picture1" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Picture1" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ec36359f-82aa-4d8e-9fc7-a025e8ea4986.png" width="644" height="358"></a> <p>There are a handful of advantages to setting up this with Docker and I want to specifically touch on some of them: <ul> <li>&nbsp;<strong>It can be automated easily</strong>. This is especially crucial for test environments which are provisioned on demand.</li> <li><strong>It’s repeatable!</strong> The declarative nature of the Dockerfile makes it possible to end up with the same environment setup even if you run the scripts months later after your initial setup.</li> <li><strong>Familiarity!</strong> Docker is a widely known and used tool for lots of other purposes and familiarity to the tool is high. Of course, this may depend on your development environment</li></ul> <h3>Let’s make it work!</h3> <p>First of all, I need to create a docker network. I can achieve this by running the "docker network create” command and giving it a unique name. <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>docker network create my<span style="color: gray">-</span>mongo<span style="color: gray">-</span>cluster</pre></div></div>

<p>The next step is to create the MongoDB docker containers and start them. I can use “docker run” command for this. Also, MongoDB has an official image on Docker Hub. So, I can reuse that to simplify the acqusition of MongoDB. For convenience, I will name the container with a number suffix. The container also needs to be tied to the network we have previously created. Finally, I need to specify the name of the replica set for each container.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>docker run <span style="color: gray">--</span>name mongo<span style="color: gray">-</span>node1 <span style="color: gray">-</span>d <span style="color: gray">--</span>net my<span style="color: gray">-</span>mongo<span style="color: gray">-</span>cluster mongo <span style="color: gray">--</span>replSet “rs0"</pre></div></div>
<p>First container is created and I need to run the same command to create two more MongoDB containers. The only difference is with the container names.
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>docker run <span style="color: gray">--</span>name mongo<span style="color: gray">-</span>node2 <span style="color: gray">-</span>d <span style="color: gray">--</span>net my<span style="color: gray">-</span>mongo<span style="color: gray">-</span>cluster mongo <span style="color: gray">--</span>replSet <span style="color: #a31515">"rs0"</span>
docker run <span style="color: gray">--</span>name mongo<span style="color: gray">-</span>node3 <span style="color: gray">-</span>d <span style="color: gray">--</span>net my<span style="color: gray">-</span>mongo<span style="color: gray">-</span>cluster mongo <span style="color: gray">--</span>replSet “rs0"</pre></div></div>
<p>I can see that all of my MongoDB containers are at the running state by executing the “docker ps” command.
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/0ccafe12-d243-48c3-93ea-0930fbaa2b99.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c328fad0-9fe9-42f4-923d-aafb9d39ccc9.png" width="644" height="74"></a>
<p>In order to form a replica set, I need to initialise the replication. I will do that by connecting to one of the containers through the “docker exec” command and starting the mongo shell client.
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>docker exec <span style="color: gray">-</span>it mongo<span style="color: gray">-</span>node1 mongo</pre></div></div>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a2b23cf7-54c7-42d9-9e8e-67f79f1fdb44.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ae3ae9c5-9098-4885-9e47-ea88e245811f.png" width="644" height="206"></a>
<p>As I now have a connection to the server, I can initialise the replication. This requires me to declare a config object which will include connection details of all the servers.
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>config = {
      <span style="color: #a31515">"_id"</span> : <span style="color: #a31515">"rs0"</span>,
      <span style="color: #a31515">"members"</span> : [
          {
              <span style="color: #a31515">"_id"</span> : 0,
              <span style="color: #a31515">"host"</span> : <span style="color: #a31515">"mongo-node1:27017"</span>
          },
          {
              <span style="color: #a31515">"_id"</span> : 1,
              <span style="color: #a31515">"host"</span> : <span style="color: #a31515">"mongo-node2:27017"</span>
          },
          {
              <span style="color: #a31515">"_id"</span> : 2,
              <span style="color: #a31515">"host"</span> : <span style="color: #a31515">"mongo-node3:27017"</span>
          }
      ]
  }</pre></div></div>
<p>Finally, we can run “rs.initialize" command to complete the set up.
<p>You will notice that the server I am connected to will be elected as the primary in the replica set shortly. By running “rs.status()”, I can view the status of other MongoDB servers within the replica set. We can see that there are two secondaries and one primary in the replica set.
<h3>.NET Core Application</h3>
<p>As a scenario, I want to run <a href="https://github.com/tugberkugurlu/mongodb-replica-set">my .NET Core application</a> which writes data to a MongoDB database and start reading it in a loop. This application will be connecting to the MongoDB replica set which we have just created.&nbsp; This is a standard .NET Core console application which you can create by running the following script:
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>dotnet new console</pre></div></div>
<p>The csproj file for this application looks like below.
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">&lt;</span><span style="color: #a31515">Project</span> <span style="color: red">Sdk</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">Microsoft.NET.Sdk</span><span style="color: black">"</span><span style="color: blue">&gt;</span>
  <span style="color: blue">&lt;</span><span style="color: #a31515">PropertyGroup</span><span style="color: blue">&gt;</span>
    <span style="color: blue">&lt;</span><span style="color: #a31515">OutputType</span><span style="color: blue">&gt;</span>Exe<span style="color: blue">&lt;/</span><span style="color: #a31515">OutputType</span><span style="color: blue">&gt;</span>
    <span style="color: blue">&lt;</span><span style="color: #a31515">TargetFramework</span><span style="color: blue">&gt;</span>netcoreapp2.0<span style="color: blue">&lt;/</span><span style="color: #a31515">TargetFramework</span><span style="color: blue">&gt;</span>
  <span style="color: blue">&lt;/</span><span style="color: #a31515">PropertyGroup</span><span style="color: blue">&gt;</span>
  <span style="color: blue">&lt;</span><span style="color: #a31515">ItemGroup</span><span style="color: blue">&gt;</span>
    <span style="color: blue">&lt;</span><span style="color: #a31515">PackageReference</span> <span style="color: red">Include</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">Bogus</span><span style="color: black">"</span> <span style="color: red">Version</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">18.0.2</span><span style="color: black">"</span> <span style="color: blue">/&gt;</span>
    <span style="color: blue">&lt;</span><span style="color: #a31515">PackageReference</span> <span style="color: red">Include</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">MongoDB.Driver</span><span style="color: black">"</span> <span style="color: red">Version</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">2.4.4</span><span style="color: black">"</span> <span style="color: blue">/&gt;</span>
    <span style="color: blue">&lt;</span><span style="color: #a31515">PackageReference</span> <span style="color: red">Include</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">Polly</span><span style="color: black">"</span> <span style="color: red">Version</span><span style="color: blue">=</span><span style="color: black">"</span><span style="color: blue">5.3.1</span><span style="color: black">"</span> <span style="color: blue">/&gt;</span>
  <span style="color: blue">&lt;/</span><span style="color: #a31515">ItemGroup</span><span style="color: blue">&gt;</span>
<span style="color: blue">&lt;/</span><span style="color: #a31515">Project</span><span style="color: blue">&gt;</span></pre></div></div>
<p>Notice that I have two interesting dependencies there. Polly is used to retry the read calls to MongoDB based on defined policies. This bit is interesting as I would expect the MongoDB client to handle that for read calls. However, it might be also a good way of explicitly stating which calls can be retried inside your application. Bogus, on the other hand, is just here to be able to create fake names to make the application a bit more realistic :)
<p> Finally, this is the code to make this application work:
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">partial</span> <span style="color: blue">class</span> Program
{
    <span style="color: blue">static</span> <span style="color: blue">void</span> Main(<span style="color: blue">string</span>[] args)
    {
        <span style="color: blue">var</span> settings = <span style="color: blue">new</span> MongoClientSettings
        {
            Servers = <span style="color: blue">new</span>[]
            {
                <span style="color: blue">new</span> MongoServerAddress(<span style="color: #a31515">"mongo-node1"</span>, 27017),
                <span style="color: blue">new</span> MongoServerAddress(<span style="color: #a31515">"mongo-node2"</span>, 27017),
                <span style="color: blue">new</span> MongoServerAddress(<span style="color: #a31515">"mongo-node3"</span>, 27017)
            },
            ConnectionMode = ConnectionMode.ReplicaSet,
            ReplicaSetName = <span style="color: #a31515">"rs0"</span>
        };

        <span style="color: blue">var</span> client = <span style="color: blue">new</span> MongoClient(settings);
        <span style="color: blue">var</span> database = client.GetDatabase(<span style="color: #a31515">"mydatabase"</span>);
        <span style="color: blue">var</span> collection = database.GetCollection&lt;User&gt;(<span style="color: #a31515">"users"</span>);

        System.Console.WriteLine(<span style="color: #a31515">"Cluster Id: {0}"</span>, client.Cluster.ClusterId);
        client.Cluster.DescriptionChanged += (<span style="color: blue">object</span> sender, ClusterDescriptionChangedEventArgs foo) =&gt; 
        {
            System.Console.WriteLine(<span style="color: #a31515">"New Cluster Id: {0}"</span>, foo.NewClusterDescription.ClusterId);
        };

        <span style="color: blue">for</span> (<span style="color: blue">int</span> i = 0; i &lt; 100; i++)
        {
            <span style="color: blue">var</span> user = <span style="color: blue">new</span> User { Id = ObjectId.GenerateNewId(), Name = <span style="color: blue">new</span> Bogus.Faker().Name.FullName() };
            collection.InsertOne(user);
        }

        <span style="color: blue">while</span> (<span style="color: blue">true</span>)
        {
            <span style="color: blue">var</span> randomUser = collection.GetRandom();
            Console.WriteLine(randomUser.Name);

            Thread.Sleep(500);
        }
    }
}</pre></div></div>
<p>This is not the most beautiful and optimized code ever but should demonstrate what we are trying to achieve by having a replica set. It's actually the GetRandom method on the MongoDB collection object which handles the retry:
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">public</span> <span style="color: blue">static</span> <span style="color: blue">class</span> CollectionExtensions 
{
    <span style="color: blue">private</span> <span style="color: blue">readonly</span> <span style="color: blue">static</span> Random random = <span style="color: blue">new</span> Random();

    <span style="color: blue">public</span> <span style="color: blue">static</span> T GetRandom&lt;T&gt;(<span style="color: blue">this</span> IMongoCollection&lt;T&gt; collection) 
    {
        <span style="color: blue">var</span> retryPolicy = Policy
            .Handle&lt;MongoCommandException&gt;()
            .Or&lt;MongoConnectionException&gt;()
            .WaitAndRetry(2, retryAttempt =&gt; 
                TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)) 
            );

        <span style="color: blue">return</span> retryPolicy.Execute(() =&gt; GetRandomImpl(collection));
    }

    <span style="color: blue">private</span> <span style="color: blue">static</span> T GetRandomImpl&lt;T&gt;(<span style="color: blue">this</span> IMongoCollection&lt;T&gt; collection)

    {
        <span style="color: blue">return</span> collection.Find(FilterDefinition&lt;T&gt;.Empty)
            .Limit(-1)
            .Skip(random.Next(99))
            .First();
    }
}</pre></div></div>
<p>I will run this through docker as well and here is the dockerfile for this:&nbsp; <div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>FROM microsoft/dotnet:2-sdk

COPY ./mongodb-replica-set.csproj /app/
WORKDIR /app/
RUN dotnet --info
RUN dotnet restore
ADD ./ /app/
RUN dotnet publish -c DEBUG -o out
ENTRYPOINT ["dotnet", "out/mongodb-replica-set.dll"]</pre></div></div>
<p>When it starts, we can see that it will output the result to the console:
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/404373db-0d84-4fa2-96f5-5f50a5bc8370.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/c453cae9-dc86-41bf-9ecf-bf1c7449c5bc.png" width="606" height="484"></a>
<h3>Prove that It Works!</h3>
<p>In order to demonstrate the effect of the replica set, I want to take down the primary node. First of all, we need to have look at the output of rs.status command we have previously ran in order to identify the primary node. We can see that it’s node1!&nbsp; <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ea0da0c1-4c20-4242-a766-5f4ff0598be4.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/2a7ea5bd-c60c-41fb-9554-e8271a849492.png" width="644" height="268"></a>
<p>Secondly, we need to get the container id for that node.&nbsp; <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ce3a7018-9e43-48fd-be87-488a0e71e784.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/b814732a-dcc3-4b8b-9ea9-dadefc375a01.png" width="644" height="79"></a>
<p>Finally, we can kill the container by running the “docker stop command”. Once the container is stopped, you will notice that application will gracefully recover and continue reading the data.&nbsp; <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d583a8d9-fbd8-4c85-8c79-668f3ef50f3a.png"><img title="Image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="Image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/6243d043-0d47-4488-9a89-a5747a58a4dc.png" width="644" height="270"></a>  