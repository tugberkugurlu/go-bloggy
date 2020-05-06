---
title: Turkish I Problem on RavenDB and Solving It with Custom Lucene Analyzers
abstract: Yesterday, I ran into a Turkish I problem on RavenDB and here is how I solved
  It with a custom Lucene analyzer
created_at: 2013-07-16 14:37:00 +0000 UTC
tags:
- Lucene.NET
- RavenDB
slugs:
- turkish-i-problem-on-ravendb-and-solving-it-with-custom-lucene-analyzers
---

<p><a href="http://ravendb.net/">RavenDB</a>, by default, uses a custom <a href="http://lucene.apache.org/core/">Lucene.Net</a> analyzer named <a href="http://ravendb.net/docs/client-api/querying/static-indexes/configuring-index-options#ravendbs-default-analyzer">LowerCaseKeywordAnalyzer</a> and it makes all your queries case-insensitive which is what I would expect. For example, the following query find the User whose name property is set to "TuGbErK":</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">class</span> Program
{
     <span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args)
     {
          <span style="color: blue;">const</span> <span style="color: blue;">string</span> DefaultDatabase = <span style="color: #a31515;">"EqualsTryOut"</span>;
          IDocumentStore store = <span style="color: blue;">new</span> DocumentStore
          {
               Url = <span style="color: #a31515;">"http://localhost:8080"</span>,
               DefaultDatabase = DefaultDatabase
          }.Initialize();
          store.DatabaseCommands.EnsureDatabaseExists(DefaultDatabase);
          <span style="color: blue;">using</span> (<span style="color: blue;">var</span> ses = store.OpenSession())
          {
               <span style="color: blue;">var</span> user = <span style="color: blue;">new</span> User { 
                  Name = <span style="color: #a31515;">"TuGbErK"</span>, 
                  Roles = <span style="color: blue;">new</span> List&lt;<span style="color: blue;">string</span>&gt; { <span style="color: #a31515;">"adMin"</span>, <span style="color: #a31515;">"GuEst"</span> } 
               };
               
               ses.Store(user);
               ses.SaveChanges();
               <span style="color: green;">//this finds name:TuGbErK</span>
               <span style="color: blue;">var</span> user1 = ses.Query&lt;User&gt;()
                  .Where(usr =&gt; usr.Name == <span style="color: #a31515;">"tugberk"</span>)
                  .FirstOrDefault();
          }
     }
}
<span style="color: blue;">public</span> <span style="color: blue;">class</span> User
{
     <span style="color: blue;">public</span> <span style="color: blue;">string</span> Id { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
     <span style="color: blue;">public</span> <span style="color: blue;">string</span> Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
     <span style="color: blue;">public</span> ICollection&lt;<span style="color: blue;">string</span>&gt; Roles { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>The problem starts appearing here where you have a situation that requires you to store Turkish text. You may ask why at this point, which makes sense. The problem is related to <a href="http://www.codinghorror.com/blog/2008/03/whats-wrong-with-turkey.html">well-known</a> <a href="http://www.west-wind.com/weblog/posts/2005/May/23/DataRows-String-Indexes-and-case-sensitivity-with-Turkish-Locale">Turkish "I"</a> <a href="http://www.hanselman.com/blog/UpdateOnTheDasBlogTurkishIBugAndAReminderToMeOnGlobalization.aspx">problem</a>. Let's try to produce this problem with an example on RavenDB.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">class</span> Program
{
    <span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args)
    {
        <span style="color: blue;">const</span> <span style="color: blue;">string</span> DefaultDatabase = <span style="color: #a31515;">"EqualsTryOut"</span>;
        IDocumentStore store = <span style="color: blue;">new</span> DocumentStore
        {
            Url = <span style="color: #a31515;">"http://localhost:8080"</span>,
            DefaultDatabase = DefaultDatabase
        }.Initialize();

        store.DatabaseCommands.EnsureDatabaseExists(DefaultDatabase);

        <span style="color: blue;">using</span> (<span style="color: blue;">var</span> ses = store.OpenSession())
        {
            <span style="color: blue;">var</span> user = <span style="color: blue;">new</span> User { 
                Name = <span style="color: #a31515;">"Irmak"</span>, 
                Roles = <span style="color: blue;">new</span> List&lt;<span style="color: blue;">string</span>&gt; { <span style="color: #a31515;">"adMin"</span>, <span style="color: #a31515;">"GuEst"</span> } 
            };
            
            ses.Store(user);
            ses.SaveChanges();

            <span style="color: green;">//This fails dues to Turkish I</span>
            <span style="color: blue;">var</span> user1 = ses.Query&lt;User&gt;()
                .Where(usr =&gt; usr.Name == <span style="color: #a31515;">"irmak"</span>)
                .FirstOrDefault();

            <span style="color: green;">//this finds name:Irmak</span>
            <span style="color: blue;">var</span> user2 = ses.Query&lt;User&gt;()
                .Where(usr =&gt; usr.Name == <span style="color: #a31515;">"IrMak"</span>)
                .FirstOrDefault();
        }
    }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> User
{
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Id { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> ICollection&lt;<span style="color: blue;">string</span>&gt; Roles { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>Here, we have the same code but the Name value we are storing is different: Irmak. Irmak is a Turkish name which also means "river" in English (which is totally not the point here) and it starts with the Turkish letter "I". The lowercased version of this letter is "ı" and this is where the problem arises because if you are lowercasing this character in an invariant culture manner, it will be transformed as "i", not "ı". This is what RavenDB is doing with its LowerCaseKeywordAnalyzer and that's why we can't find anything with the first query above where we searched against "ırmak". In the second query, we can find the User that we are looking for as it will be lowercased into "irmak".</p>
<h3>The Solution with a Custom Analyzer</h3>
<p>The default analyzer that RavenDB using is the <a href="https://github.com/ravendb/ravendb/blob/8351242c97ba7595ee34edaaa71be85fad76efe2/Raven.Database/Indexing/Analyzers/LowerCaseKeywordAnalyzer.cs">LowerCaseKeywordAnalyzer</a> and it uses the <a href="https://github.com/ravendb/ravendb/blob/8351242c97ba7595ee34edaaa71be85fad76efe2/Raven.Database/Indexing/Analyzers/LowerCaseKeywordTokenizer.cs">LowerCaseKeywordTokenizer</a> as its tokenizer. Inside that tokenizer, you will see the <a href="https://github.com/ravendb/ravendb/blob/8351242c97ba7595ee34edaaa71be85fad76efe2/Raven.Database/Indexing/Analyzers/LowerCaseKeywordTokenizer.cs#L51-L54">Normalize method</a> which is used to lowercase a character in an invariant culture manner which causes a problem on our case here. AFAIK, there is no built in Lucene.Net tokenizer which handles Turkish characters well (I might be wrong here). So, I decided to modify the LowerCaseKeywordTokenizer according to my specific needs. It was a very naive and minor change which worked but not sure if I handled it well. You can find the source code of the <a href="https://github.com/tugberkugurlu/LuceneAnalyzers/blob/ef1d50eb60ed96052baba762287b6dbe8b750bbd/src/LuceneAnalyzers/TutkishLowerCaseKeywordTokenizer.cs">TutkishLowerCaseKeywordTokenizer</a> and <a href="https://github.com/tugberkugurlu/LuceneAnalyzers/blob/ef1d50eb60ed96052baba762287b6dbe8b750bbd/src/LuceneAnalyzers/TurkishLowerCaseKeywordAnalyzer.cs">TurkishLowerCaseKeywordAnalyzer</a> classes on my Github repository.</p>
<h3>Using a Custom Build Analyzer on RavenDB</h3>
<p>RavenDB allows us to use custom analyzers and control the analyzer per-field. If you're going to use the built-in Lucene analyzer for a field, you can need to pass the <a href="http://msdn.microsoft.com/en-us/library/system.type.fullname.aspx">FullName</a> of the analyzer type just like in the below example which is a straight copy and paste <a href="http://ravendb.net/docs/appendixes/lucene-indexes-usage#using-custom-analyzers">from the RavenDB documentation</a>.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>store.DatabaseCommands.PutIndex(
    <span style="color: #a31515;">"Movies"</span>,
    <span style="color: blue;">new</span> IndexDefinition
        {
            Map = <span style="color: #a31515;">"from movie in docs.Movies select new { movie.Name, movie.Tagline }"</span>,
            Analyzers =
                {
                    { <span style="color: #a31515;">"Name"</span>, <span style="color: blue;">typeof</span>(SimpleAnalyzer).FullName },
                    { <span style="color: #a31515;">"Tagline"</span>, <span style="color: blue;">typeof</span>(StopAnalyzer).FullName },
                }
        });</pre>
</div>
</div>
<p>On the other hand, RavenDB also allows us to drop our own custom analyzers in:</p>
<p><i>"You can also create your own custom analyzer, compile it to a dll and drop it in in directory called "Analyzers" under the RavenDB base directory. Afterward, you can then use the fully qualified type name of your custom analyzer as the analyzer for a particular field."</i></p>
<p>There are couple things that you need to be careful of when going down this road:</p>
<ul>
<li>You need to use <a href="https://groups.google.com/d/msg/ravendb/IMTZqboPyYM/_9-9fPBPQygJ">the Lucene.Net assembly that your RavenDB server is using because RavenDB is using a custom build</a>. </li>
<li>Drop your compiled assembly into the directory called "Analyzers" under the RavenDB base directory. </li>
<li>When you are configuring the specific fields to use your analyzer, be sure to pass the <a href="http://msdn.microsoft.com/en-us/library/system.type.assemblyqualifiedname.aspx">AssemblyQualifiedName</a> of your custom analyzer class.</li>
</ul>
<p>After I stopped my RavenDB server, I dropped my assembly, which contains the TurkishLowerCaseKeywordAnalyzer, into the "Analyzers" folder under the RavenDB base directory. At the client side, here is my code which consists of the index creation and the query:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">class</span> Program
{
    <span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args)
    {
        <span style="color: blue;">const</span> <span style="color: blue;">string</span> DefaultDatabase = <span style="color: #a31515;">"EqualsTryOut"</span>;
        IDocumentStore store = <span style="color: blue;">new</span> DocumentStore
        {
            Url = <span style="color: #a31515;">"http://localhost:8080"</span>,
            DefaultDatabase = DefaultDatabase
        }.Initialize();

        IndexCreation.CreateIndexes(<span style="color: blue;">typeof</span>(Users).Assembly, store);
        store.DatabaseCommands.EnsureDatabaseExists(DefaultDatabase);

        <span style="color: blue;">using</span> (<span style="color: blue;">var</span> ses = store.OpenSession())
        {
            <span style="color: blue;">var</span> user = <span style="color: blue;">new</span> User { 
                Name = <span style="color: #a31515;">"Irmak"</span>, 
                Roles = <span style="color: blue;">new</span> List&lt;<span style="color: blue;">string</span>&gt; { <span style="color: #a31515;">"adMin"</span>, <span style="color: #a31515;">"GuEst"</span> } 
            };
            ses.Store(user);
            ses.SaveChanges();

            <span style="color: green;">//this finds name:Irmak</span>
            <span style="color: blue;">var</span> user1 = ses.Query&lt;User, Users&gt;()
                .Where(usr =&gt; usr.Name == <span style="color: #a31515;">"irmak"</span>)
                .FirstOrDefault();
        }
    }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> User
{
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Id { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> ICollection&lt;<span style="color: blue;">string</span>&gt; Roles { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> Users : AbstractIndexCreationTask&lt;User&gt;
{
    <span style="color: blue;">public</span> Users()
    {
        Map = users =&gt; <span style="color: blue;">from</span> user <span style="color: blue;">in</span> users
                       <span style="color: blue;">select</span> <span style="color: blue;">new</span> 
                       {
                          user.Name 
                       };

        Analyzers.Add(
            x =&gt; x.Name, 
            <span style="color: blue;">typeof</span>(LuceneAnalyzers.TurkishLowerCaseKeywordAnalyzer)
                .AssemblyQualifiedName);
    }
}</pre>
</div>
</div>
<p>It worked like a charm. I'm hopping this helps you solve this annoying problem and please post your comment if you know a better way of handling this.</p>
<h3>Resources</h3>
<ul>
<li><a href="https://groups.google.com/forum/#!topic/ravendb/IMTZqboPyYM">Turkish I problem on queries and custom TurkishAnalyzer for Lucene.Net</a> </li>
<li><a href="https://github.com/tugberkugurlu/LuceneAnalyzers/blob/master/src/LuceneAnalyzers/TurkishLowerCaseKeywordAnalyzer.cs">TurkishLowerCaseKeywordAnalyzer</a> </li>
<li><a href="https://groups.google.com/forum/#!msg/ravendb/-__e74nBqDM/1c6W6a3vlzYJ">Can't get a custom Analyzer to work</a></li>
</ul>