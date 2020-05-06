---
title: Asynchronous Database Calls With Task-based Asynchronous Programming Model
  (TAP) in ASP.NET MVC 4
abstract: Asynchronous Database Calls With Task-based Asynchronous Programming Model
  (TAP) in ASP.NET MVC 4 and its performance impacts.
created_at: 2012-04-09 11:23:00 +0000 UTC
tags:
- ASP.NET MVC
- async
- C#
- MS SQL
- TPL
slugs:
- asynchronous-database-calls-with-task-based-asynchronous-programming-model-tap-in-asp-net-mvc-4
- asynchronous-database-calls-with-task-based-asynchronous-progra
---

<p>You have probably seen some people who are against asynchronous database calls in ASP.NET web applications so far. They mostly right but there are still some cases that processing database queries asynchronous has very important impact on.</p>
<blockquote>
<p>If you are unfamiliar with asynchronous programming model on ASP.NET MVC 4, you might want to read one of my previous posts: <a href="http://www.tugberkugurlu.com/archive/my-take-on-task-base-asynchronous-programming-in-c-sharp-5-0-and-asp-net-mvc-web-applications" title="http://www.tugberkugurlu.com/archive/my-take-on-task-base-asynchronous-programming-in-c-sharp-5-0-and-asp-net-mvc-web-applications">Asynchronous Programming in C# 5.0 and ASP.NET MVC Web Applications</a>.</p>
<p><strong>Update on the 11th of April, 2012:</strong></p>
<p><a href="http://twitter.com/bradwilson" title="http://twitter.com/bradwilson">@BradWilson</a>&nbsp;has just&nbsp;started a new blog post&nbsp;series on using Task Parallel Library when writing server applications, especially ASP.NET MVC and ASP.NET Web API applications. You should certainly check them out:&nbsp;<a href="http://bradwilson.typepad.com/blog/2012/04/tpl-and-servers-pt1.html" title="http://bradwilson.typepad.com/blog/2012/04/tpl-and-servers-pt1.html">Task Parallel Library and Servers, Part 1: Introduction</a></p>
<p><strong>Update on the 1st of July, 2012:</strong></p>
<p><a title="https://twitter.com/#!/RickAndMSFT" href="https://twitter.com/#!/RickAndMSFT">@RickAndMSFT</a> has a new tutorial on <a title="http://www.asp.net/mvc/tutorials/mvc-4/using-asynchronous-methods-in-aspnet-mvc-4" href="http://www.asp.net/mvc/tutorials/mvc-4/using-asynchronous-methods-in-aspnet-mvc-4">Using Asynchronous Methods in ASP.NET MVC 4</a>. Definitely&nbsp;check that out.</p>
</blockquote>
<p>One of the reasons why asynchronous programming is not recommended for database calls is that it is extremely hard to get it right, even if we adopt Task-based Asynchronous Programming in .NET 4.0. But, with the new async / await features of C# 5.0, it is easier and still complex at the same time.</p>
<p>When you call a method which returns <strong>Task</strong> or <strong>Task&lt;T&gt;</strong> for some T and await on that, the compiler does the heavy lifting by assigning continuations, handling exceptions and so on. Because of this fact, it adds a little overhead and you might notice this when you are dealing with short running operations. At that point, asynchrony will do more harm than good to your application.</p>
<p>Let&rsquo;s assume we have a SQL Server database out there somewhere and we want to query against that database in order to get the cars list that we have. I have created a class which will do the query operations and hand us the results as C# CLR objects.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> GalleryContext : IGalleryContext {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">string</span> selectStatement = <span style="color: #a31515;">"SELECT * FROM Cars"</span>;

    <span style="color: blue;">public</span> IEnumerable&lt;Car&gt; GetCars() {

        <span style="color: blue;">var</span> connectionString = 
            ConfigurationManager.ConnectionStrings[<span style="color: #a31515;">"CarGalleryConnStr"</span>].ConnectionString;

        <span style="color: blue;">using</span> (<span style="color: blue;">var</span> conn = <span style="color: blue;">new</span> SqlConnection(connectionString)) {
            <span style="color: blue;">using</span> (<span style="color: blue;">var</span> cmd = <span style="color: blue;">new</span> SqlCommand()) {

                cmd.Connection = conn;
                cmd.CommandText = selectStatement;
                cmd.CommandType = CommandType.Text;

                conn.Open();

                <span style="color: blue;">using</span> (<span style="color: blue;">var</span> reader = cmd.ExecuteReader()) {

                    <span style="color: blue;">return</span> reader.Select(r =&gt; carBuilder(r)).ToList();
                }
            }
        }
    }

    <span style="color: blue;">public</span> async Task&lt;IEnumerable&lt;Car&gt;&gt; GetCarsAsync() {

        <span style="color: blue;">var</span> connectionString = 
            ConfigurationManager.ConnectionStrings[<span style="color: #a31515;">"CarGalleryConnStr"</span>].ConnectionString;
            
        <span style="color: blue;">var</span> asyncConnectionString = <span style="color: blue;">new</span> SqlConnectionStringBuilder(connectionString) {
            AsynchronousProcessing = <span style="color: blue;">true</span>
        }.ToString();

        <span style="color: blue;">using</span> (<span style="color: blue;">var</span> conn = <span style="color: blue;">new</span> SqlConnection(asyncConnectionString)) {
            <span style="color: blue;">using</span> (<span style="color: blue;">var</span> cmd = <span style="color: blue;">new</span> SqlCommand()) {

                cmd.Connection = conn;
                cmd.CommandText = selectStatement;
                cmd.CommandType = CommandType.Text;

                conn.Open();

                <span style="color: blue;">using</span> (<span style="color: blue;">var</span> reader = await cmd.ExecuteReaderAsync()) {

                    <span style="color: blue;">return</span> reader.Select(r =&gt; carBuilder(r)).ToList();
                }
            }
        }
    }

    <span style="color: green;">//private helpers</span>
    <span style="color: blue;">private</span> Car carBuilder(SqlDataReader reader) {

        <span style="color: blue;">return</span> <span style="color: blue;">new</span> Car {

            Id = <span style="color: blue;">int</span>.Parse(reader[<span style="color: #a31515;">"Id"</span>].ToString()),
            Make = reader[<span style="color: #a31515;">"Make"</span>] <span style="color: blue;">is</span> DBNull ? <span style="color: blue;">null</span> : reader[<span style="color: #a31515;">"Make"</span>].ToString(),
            Model = reader[<span style="color: #a31515;">"Model"</span>] <span style="color: blue;">is</span> DBNull ? <span style="color: blue;">null</span> : reader[<span style="color: #a31515;">"Model"</span>].ToString(),
            Year = <span style="color: blue;">int</span>.Parse(reader[<span style="color: #a31515;">"Year"</span>].ToString()),
            Doors = <span style="color: blue;">int</span>.Parse(reader[<span style="color: #a31515;">"Doors"</span>].ToString()),
            Colour = reader[<span style="color: #a31515;">"Colour"</span>] <span style="color: blue;">is</span> DBNull ? <span style="color: blue;">null</span> : reader[<span style="color: #a31515;">"Colour"</span>].ToString(),
            Price = <span style="color: blue;">float</span>.Parse(reader[<span style="color: #a31515;">"Price"</span>].ToString()),
            Mileage = <span style="color: blue;">int</span>.Parse(reader[<span style="color: #a31515;">"Mileage"</span>].ToString())
        };
    }
}</pre>
</div>
</div>
<blockquote>
<p>You might notice that I used a Select method on SqlDataReader which does not exist. It is a small extension method which makes it look prettier.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">static</span> <span style="color: blue;">class</span> Extensions {

    <span style="color: blue;">public</span> <span style="color: blue;">static</span> IEnumerable&lt;T&gt; Select&lt;T&gt;(
        <span style="color: blue;">this</span> SqlDataReader reader, Func&lt;SqlDataReader, T&gt; projection) {

        <span style="color: blue;">while</span> (reader.Read()) {
            yield <span style="color: blue;">return</span> projection(reader);
        }
    }
}</pre>
</div>
</div>
</blockquote>
<p>As you can see, it has two public methods to query the database which does the same thing but one of them doing it as synchronously and the other one as asynchronously.</p>
<p>Inside the GetCarsAsync method, you can see that we append the <a href="http://msdn.microsoft.com/en-us/library/system.data.sqlclient.sqlconnectionstringbuilder.asynchronousprocessing.aspx" title="http://msdn.microsoft.com/en-us/library/system.data.sqlclient.sqlconnectionstringbuilder.asynchronousprocessing.aspx">AsynchronousProcessing</a> property and set it to true in order to run the operation asynchronously. Otherwise, no matter how you implement it, your query will be processed synchronously.</p>
<blockquote>
<p>When you look behind the curtain, you will notice that <strong>ExecuteReaderAsync</strong> method is really using the old Asynchronous Programming Model (APM) under the covers.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">static</span> Task&lt;SqlDataReader&gt; ExecuteReaderAsync(<span style="color: blue;">this</span> SqlCommand source)
{
    <span style="color: blue;">return</span> Task&lt;SqlDataReader&gt;.Factory.FromAsync(
        <span style="color: blue;">new</span> Func&lt;AsyncCallback, <span style="color: blue;">object</span>, IAsyncResult&gt;(source.BeginExecuteReader), 
        <span style="color: blue;">new</span> Func&lt;IAsyncResult, SqlDataReader&gt;(source.EndExecuteReader), 
        <span style="color: blue;">null</span>
    );
}</pre>
</div>
</div>
</blockquote>
<p>When we try to consume these methods inside our controller, we will have the following implementation.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> GalleryContext ctx = <span style="color: blue;">new</span> GalleryContext();

    <span style="color: blue;">public</span> ViewResult Index() {

        <span style="color: blue;">return</span> View(ctx.GetCars());
    }

    <span style="color: blue;">public</span> async Task&lt;ViewResult&gt; IndexAsync() {

        <span style="color: green;">//workaround: http://aspnetwebstack.codeplex.com/workitem/22</span>
        await TaskEx.Yield();

        <span style="color: blue;">return</span> View(<span style="color: #a31515;">"Index"</span>, await ctx.GetCarsAsync());
    }
}</pre>
</div>
</div>
<p>Now, when we hit <strong>/home/Index</strong>, we will be querying our database as synchronously. If we navigate to <strong>/home/IndexAsync</strong>, we will be doing the same thing but asynchronously this time. Let&rsquo;s do a little benchmarking with <a href="http://httpd.apache.org/docs/2.0/programs/ab.html" title="http://httpd.apache.org/docs/2.0/programs/ab.html">Apache HTTP server benchmarking tool</a>.</p>
<p>First, we will simulate 50 concurrent requests on synchronously running wev page:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_syncdb_short_1.png"><img height="484" width="530" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_syncdb_short_1_thumb.png" alt="ab_syncdb_short_1" border="0" title="ab_syncdb_short_1" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Let&rsquo;s do the same thing for asynchronous one:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_asyncdb_short_1.png"><img height="484" width="540" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_asyncdb_short_1_thumb.png" alt="ab_asyncdb_short_1" border="0" title="ab_asyncdb_short_1" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Did you notice? We have nearly got the same result. In fact, you will see that synchronous version of the operation completes faster than the asynchronous one at some points. The reason is that the SQL query takes small amount of time (approx. 8ms) here to complete.</p>
<p>Let&rsquo;s take another scenario. Now, we will try to get the same data through a Stored Procedure but this time, the database call will be slow (approx. 1 second). Here are two methods which will nearly the same as others:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">private</span> <span style="color: blue;">readonly</span> <span style="color: blue;">string</span> spName = <span style="color: #a31515;">"sp$GetCars"</span>;

<span style="color: blue;">public</span> IEnumerable&lt;Car&gt; GetCarsViaSP() {

    <span style="color: blue;">var</span> connectionString = ConfigurationManager.ConnectionStrings[<span style="color: #a31515;">"CarGalleryConnStr"</span>].ConnectionString;

    <span style="color: blue;">using</span> (<span style="color: blue;">var</span> conn = <span style="color: blue;">new</span> SqlConnection(connectionString)) {
        <span style="color: blue;">using</span> (<span style="color: blue;">var</span> cmd = <span style="color: blue;">new</span> SqlCommand()) {

            cmd.Connection = conn;
            cmd.CommandText = spName;
            cmd.CommandType = CommandType.StoredProcedure;

            conn.Open();

            <span style="color: blue;">using</span> (<span style="color: blue;">var</span> reader = cmd.ExecuteReader()) {

                <span style="color: blue;">return</span> reader.Select(r =&gt; carBuilder(r)).ToList();
            }
        }
    }
}

<span style="color: blue;">public</span> async Task&lt;IEnumerable&lt;Car&gt;&gt; GetCarsViaSPAsync() {

    <span style="color: blue;">var</span> connectionString = ConfigurationManager.ConnectionStrings[<span style="color: #a31515;">"CarGalleryConnStr"</span>].ConnectionString;
    <span style="color: blue;">var</span> asyncConnectionString = <span style="color: blue;">new</span> SqlConnectionStringBuilder(connectionString) {
        AsynchronousProcessing = <span style="color: blue;">true</span>
    }.ToString();

    <span style="color: blue;">using</span> (<span style="color: blue;">var</span> conn = <span style="color: blue;">new</span> SqlConnection(asyncConnectionString)) {
        <span style="color: blue;">using</span> (<span style="color: blue;">var</span> cmd = <span style="color: blue;">new</span> SqlCommand()) {

            cmd.Connection = conn;
            cmd.CommandText = spName;
            cmd.CommandType = CommandType.StoredProcedure;

            conn.Open();

            <span style="color: blue;">using</span> (<span style="color: blue;">var</span> reader = await cmd.ExecuteReaderAsync()) {

                <span style="color: blue;">return</span> reader.Select(r =&gt; carBuilder(r)).ToList();
            }
        }
    }
}</pre>
</div>
</div>
<blockquote>
<p>I was able to make the SQL query long running by waiting inside the stored procedure for 1 second:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">ALTER</span> <span style="color: blue;">PROCEDURE</span> dbo.sp$GetCars
<span style="color: blue;">AS</span>

<span style="color: green;">-- wait for 1 second</span>
<span style="color: blue;">WAITFOR</span> DELAY <span style="color: #a31515;">'00:00:01'</span>;

<span style="color: blue;">SELECT</span> * <span style="color: blue;">FROM</span> Cars;</pre>
</div>
</div>
</blockquote>
<p>Controller actions are nearly the same as before:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> GalleryContext ctx = <span style="color: blue;">new</span> GalleryContext();

    <span style="color: blue;">public</span> ViewResult IndexSP() {

        <span style="color: blue;">return</span> View(<span style="color: #a31515;">"Index"</span>, ctx.GetCarsViaSP());
    }

    <span style="color: blue;">public</span> async Task&lt;ViewResult&gt; IndexSPAsync() {

        <span style="color: green;">//workaround: http://aspnetwebstack.codeplex.com/workitem/22</span>
        await TaskEx.Yield();

        <span style="color: blue;">return</span> View(<span style="color: #a31515;">"Index"</span>, await ctx.GetCarsViaSPAsync());
    }
}</pre>
</div>
</div>
<p>Let&rsquo;s simulate 50 concurrent requests on synchronously running web page first:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_syncdb_long_1.png"><img height="484" width="452" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_syncdb_long_1_thumb.png" alt="ab_syncdb_long_1" border="0" title="ab_syncdb_long_1" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Ok, that&rsquo;s apparently not good. As you might see, some requests take more than 8 seconds to complete which is very bad for a web page. Remember, the database call takes approx. 1 second to complete and under a particular number of concurrent requests, we experience a very serious bottleneck here.</p>
<p>Let&rsquo;s have a look at the asynchronous implementation and see the difference:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_asyncdb_long_1.png"><img height="484" width="452" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/ASP.NET-MVC-4_E432/ab_asyncdb_long_1_thumb.png" alt="ab_asyncdb_long_1" border="0" title="ab_asyncdb_long_1" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Approximately 1 second to complete for each request, pretty impressive compared to synchronous implementation.</p>
<p>Asynchronous database calls are not as straight forward as other types of asynchronous operations but sometimes it will gain so much more responsiveness to our applications. We just need to get it right and implement them properly.</p>