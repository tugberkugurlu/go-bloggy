---
title: A Gentle Introduction to Azure Search
abstract: Microsoft Azure team released Azure Search as a preview product a few days
  ago, an hosted search service solution by Microsoft. Azure Search is a suitable
  product if you are dealing with high volume of data (millions of records) and want
  to have efficient, complex and clever search on those chunk of data. In this post,
  I will try to lay out some fundamentals about this service with a very high level
  introduction.
created_at: 2014-09-10 12:02:00 +0000 UTC
tags:
- Azure Search
- Microsoft Azure
slugs:
- a-gentle-introduction-to-azure-search
---

<p>With many of the applications we build as software developers, we need our data to be exposed and we want that data to be in an easy reach so that the user of the application can find what they are looking for easily. This task is especially tricky if you have high amount of data (millions, even billions) in your system. At that point, the application needs to give user a great and flawless experience so that the user can filter down the results based on what they are actually looking for. Don't we have solutions to address this problems? Of course, we do and solutions such as <a href="http://www.elasticsearch.org/">Elasticsearch</a> and <a href="http://lucene.apache.org/solr/">Apache Solr</a> are top notch problem solvers for this matter. However, hosting these products on your environment and making them scalable is completely another job.  <p>To address these problems, <a href="http://blogs.technet.com/b/dataplatforminsider/archive/2014/08/21/azure-previews-fully-managed-nosql-database-and-search-services.aspx">Microsoft Azure team released Azure Search</a> as a preview product a few days ago, an hosted search service solution by Microsoft. Azure Search is a suitable product if you are dealing with high volume of data (millions of records) and want to have efficient, complex and clever search on those chunk of data. If you have worked with a search engine product (such as Elasticsearch, Apache Solr, etc.) before, you will be too much comfortable with Azure Search as it has some many similar features. In fact, Azure Search is on top of Elasticsearch to provide its full-text search function. However, you shouldn't see this brand-new product as hosted Elasticsearch service on Azure because it has its completely different public interface.  <p>In this post, I will try to lay out some fundamentals about this service with a very high level introduction. I’m hoping that it’s also going to be a starting point for me on Azure Search blog posts :)  <h3>Getting to Know Azure Search Service</h3> <p>When I look at Azure Search service, I see it as four pieces which gives us the whole experience:  <ul> <li>Search Service  <li>Search Unit  <li>Index  <li>Document</li></ul> <p><strong>Search service</strong> is the highest level of the hierarchy and it contains Provisioned search unit(s). Also, a few concepts are targeting the search service such as authentication and scaling.  <p><strong>Search units</strong> allow for scaling of QPS (Queries per second), Document Count and Document Size. This also means that search units are the key concept for high availability and throughput. As a side note, high availability requires at least 3 replicas for the preview.  <p><strong>Index</strong> is the holder for a collection of documents based on a defined schema which specifies the capabilities of the Index (we will touch on this schema later). A search service can contain multiple indexes.  <p>Lastly, <strong>Document</strong> is the actual holder for the data, based on the index schema, which the document itself lives in. A document has a key and this key needs to be unique within the index. A document also has fields to represent the data. Fields of a document contain attributes and those attributes define the capabilities of the field such as whether it can be used to filter the results, etc. Also note that number of documents an index can contain is limited based on the search units the service has.  <h3>Windows Azure Portal Experience</h3> <p>Let's first have a look at the portal experience and how we can get a search service ready for our use. Azure Search is not available through the <a href="https://manage.windowsazure.com/">current Microsoft Azure portal</a>. It's only available through <a href="https://portal.azure.com/">the preview portal</a>. Inside the new portal, click the big plus sign at the bottom left and then click "Everything".  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/91c2d717-9a84-4159-88f4-1357e261869d.png"><img title="Screenshot 2014-09-06 11.55.45" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-09-06 11.55.45" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/e7355ba1-ffe2-494d-919b-23168e2c47f6.png" width="644" height="363"></a>  <p>This is going to get you to "Gallery". From there click "Data, storage, cache + backup" and then click "Search" from the new section.  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/488102ed-f6c5-4e77-b745-eb7eb7d14635.png"><img title="Screenshot 2014-09-06 11.59.16" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-09-06 11.59.16" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d2115ba2-ef40-42ca-aa61-a455a46d1a3d.png" width="644" height="363"></a>  <p>You will have a nice intro about the Microsoft Azure Search service within the new window. Hit "Create" there.  <blockquote> <p>Keep in mind that service name must only contain lowercase letters, digits or dashes, cannot use dash as the first two or last one characters, cannot contain consecutive dashes, and is limited between 2 and 15 characters in length. Other naming conventions about the service has been laid out <a href="http://msdn.microsoft.com/en-us/library/azure/dn798935.aspx">here</a> under Naming Conventions section.</p></blockquote> <p>When you come to selecting the Pricing Tier, it's time to make a decision about your usage scenario.  <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5922f7c0-3909-4435-9b39-7e2b33467eca.png"><img title="Screenshot 2014-09-06 12.06.52" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-09-06 12.06.52" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/5d143b0c-899d-46a5-afb5-d8646a9cb6af.png" width="644" height="391"></a>  <p>Now, there two options: Standard and Free. Free one should be considered as the sandbox experience because it's too limiting in terms of both performance and storage space. You shouldn't try to evaluate the Azure Search service with the free tier. It's, however, great for evaluating the HTTP API. You can create a free service and use this service to run your HTTP requests against.  <p>The standard tier is the one you would like to choose for production use. It can be scaled both in terms of QPS (Queries per Second) and document size through shards and replicas. Head to "<a href="http://azure.microsoft.com/en-us/documentation/articles/search-configure/">Configure Search in the Azure Preview portal</a>" article for more in depth information about scaling.  <p>When you are done setting up your service, you can now get the admin key or the query key from the portal and start hitting the <a href="http://msdn.microsoft.com/library/azure/dn798935.aspx">Azure Search HTTP (or REST, if you want to call it that) API</a>.  <h3>Azure Search HTTP API</h3> <p>Azure Search service is managed through its HTTP API and it's not hard to guess that even the Azure Portal is using its API to manage the service. It's a lightweight API which understands JSON as the content type. When we look at it, we can divide this HTTP API into three parts:  <ul> <li><a href="http://msdn.microsoft.com/en-us/library/azure/dn798918.aspx">Index Management</a>  <li><a href="http://msdn.microsoft.com/en-us/library/azure/dn800962.aspx">Index Population</a>  <li><a href="http://msdn.microsoft.com/en-us/library/azure/dn798927.aspx">Query</a></li></ul> <p><strong>Index Management</strong> part of the API allows us managing the indexes with various operations such as creating, deleting and listing the indexes. It also allow us to see some <a href="http://msdn.microsoft.com/en-us/library/azure/dn798942.aspx">index statistics</a>. Creating the index is probably going to be the first operation you will perform and it has the following structure:  <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>POST https://{search-service-name}.search.windows.net/indexes?api-version=2014-07-31-Preview HTTP/1.1
User-Agent: Fiddler
api-key: {your-api-key}
Content-Type: application/json
Host:{search-service-name}.search.windows.net

{
	<span style="color: #a31515">"name"</span>: <span style="color: #a31515">"employees"</span>,
	<span style="color: #a31515">"fields"</span>: [{
		<span style="color: #a31515">"name"</span>: <span style="color: #a31515">"employeeId"</span>,
		<span style="color: #a31515">"type"</span>: <span style="color: #a31515">"Edm.String"</span>,
		<span style="color: #a31515">"key"</span>: <span style="color: blue">true</span>,
		<span style="color: #a31515">"searchable"</span>: <span style="color: blue">false</span>
	},
	{
		<span style="color: #a31515">"name"</span>: <span style="color: #a31515">"firstName"</span>,
		<span style="color: #a31515">"type"</span>: <span style="color: #a31515">"Edm.String"</span>
	},
	{
		<span style="color: #a31515">"name"</span>: <span style="color: #a31515">"lastName"</span>,
		<span style="color: #a31515">"type"</span>: <span style="color: #a31515">"Edm.String"</span>
	},
	{
		<span style="color: #a31515">"name"</span>: <span style="color: #a31515">"age"</span>,
		<span style="color: #a31515">"type"</span>: <span style="color: #a31515">"Edm.Int32"</span>
	},
	{
		<span style="color: #a31515">"name"</span>: <span style="color: #a31515">"about"</span>,
		<span style="color: #a31515">"type"</span>: <span style="color: #a31515">"Edm.String"</span>,
		<span style="color: #a31515">"filterable"</span>: <span style="color: blue">false</span>,
		<span style="color: #a31515">"facetable"</span>: <span style="color: blue">false</span>
	},
	{
		<span style="color: #a31515">"name"</span>: <span style="color: #a31515">"interests"</span>,
		<span style="color: #a31515">"type"</span>: <span style="color: #a31515">"Collection(Edm.String)"</span>
	}]
}</pre></div></div>
<blockquote>
<p>With the above request, you can also spot a few more things which are applied to every API call we make. There is a header we are sending with the request: api-key. This is where you are supposed to put your api-key. Also, we are passing the API version through a query string parameter called api-version. Have a look at the <a href="http://msdn.microsoft.com/en-us/library/azure/dn798935.aspx">Azure Search REST API MSDN documentation</a> for further detailed information.</p></blockquote>
<p>With this request, we are specifying the schema of the index. Keep in mind that schema updates are limited at the time of this writing. Although existing fields cannot be changed or deleted, new fields can be added at any time. When a new field is added, all existing documents in the index will automatically have a null value for that field. No additional storage space will be consumed until new documents are added to the index. Have a look at the <a href="http://msdn.microsoft.com/en-us/library/azure/dn800964.aspx">Update Index API documentation</a> for further information on index schema update. 
<p>After you have your index schema defined, you can now start populating your index with <strong>Index Population</strong> API. Index Population API is a little bit different and I honestly don’t like it (I have a feeling that <a href="https://twitter.com/darrel_miller">Darrel Miller</a> won’t like it, too :)). The reason why I don’t like it is the way we define the operation. With this HTTP API, we can add new document, update and remove an existing one. However, we are defining the type of the operation inside the request body which is so weird if you ask me. The other weird thing about this API is that you can send multiple operations in one HTTP request by putting them inside a JSON array. The important fact here is that those operations don’t run in transaction which means that some of them may succeed and some of them may fail. So, how do we know which one actually failed? The response will contain a JSON array indicating each operation’s status. Nothing wrong with that but why do we reinvent the World? :) I would be more happy to send batch request using <a href="http://www.w3.org/Protocols/rfc1341/7_2_Multipart.html">the multipart content-type</a>. Anyway, enough bitching about the API :) Here is a sample request to add a new document to the index: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>POST https:<span style="color: green">//{search-service-name}.search.windows.net/indexes/employees/docs/index?api-version=2014-07-31-Preview HTTP/1.1</span>
User-Agent: Fiddler
api-key: {your-api-key}
Content-Type: application/json
Host: {search-service-name}.search.windows.net

{
	<span style="color: #a31515">"value"</span>: [{
		<span style="color: #a31515">"@search.action"</span>: <span style="color: #a31515">"upload"</span>,
		<span style="color: #a31515">"employeeId"</span>: <span style="color: #a31515">"1"</span>,
		<span style="color: #a31515">"firstName"</span>: <span style="color: #a31515">"Jane"</span>,
		<span style="color: #a31515">"lastName"</span>: <span style="color: #a31515">"Smith"</span>,
		<span style="color: #a31515">"age"</span>: 32,
		<span style="color: #a31515">"about"</span>: <span style="color: #a31515">"I like to collect rock albums"</span>,
		<span style="color: #a31515">"interests"</span>: [<span style="color: #a31515">"music"</span>]
	}]
}</pre></div></div>
<p>As said, you can send the operations in batch: 
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>POST https:<span style="color: green">//{search-service-name}.search.windows.net/indexes/employees/docs/index?api-version=2014-07-31-Preview HTTP/1.1</span>
User-Agent: Fiddler
api-key: {your-api-key}
Content-Type: application/json
Host: {search-service-name}.search.windows.net

{
	<span style="color: #a31515">"value"</span>: [{
		<span style="color: #a31515">"@search.action"</span>: <span style="color: #a31515">"upload"</span>,
		<span style="color: #a31515">"employeeId"</span>: <span style="color: #a31515">"2"</span>,
		<span style="color: #a31515">"firstName"</span>: <span style="color: #a31515">"Douglas"</span>,
		<span style="color: #a31515">"lastName"</span>: <span style="color: #a31515">"Fir"</span>,
		<span style="color: #a31515">"age"</span>: 35,
		<span style="color: #a31515">"about"</span>: <span style="color: #a31515">"I like to build cabinets"</span>,
		<span style="color: #a31515">"interests"</span>: [<span style="color: #a31515">"forestry"</span>]
	},
	{
		<span style="color: #a31515">"@search.action"</span>: <span style="color: #a31515">"upload"</span>,
		<span style="color: #a31515">"employeeId"</span>: <span style="color: #a31515">"3"</span>,
		<span style="color: #a31515">"firstName"</span>: <span style="color: #a31515">"John"</span>,
		<span style="color: #a31515">"lastName"</span>: <span style="color: #a31515">"Fir"</span>,
		<span style="color: #a31515">"age"</span>: 25,
		<span style="color: #a31515">"about"</span>: <span style="color: #a31515">"I love to go rock climbing"</span>,
		<span style="color: #a31515">"interests"</span>: [<span style="color: #a31515">"sports"</span>, <span style="color: #a31515">"music"</span>]
	}]
}</pre></div></div>
<p>Check out the <a href="http://msdn.microsoft.com/en-us/library/azure/dn800962.aspx">great documentation about index population API</a> to learn about it more.</p>
<p>Lastly, there are <strong>query and lookup APIs</strong> where you can use <a href="http://www.odata.org/documentation/odata-version-4-0/">OData 4.0</a> expression syntax to define your query. Go and check out <a href="http://msdn.microsoft.com/en-us/library/azure/dn798927.aspx">its documentation</a> as well.</p>
<p>Even if the service is so new, there are already great things happening around it. <a href="https://twitter.com/sandrinodm">Sandrino Di Mattia</a> has two cool open source projects on Azure Search. One is <a href="http://fabriccontroller.net/blog/posts/introducing-microsoft-azure-search-and-the-reddog-search-client-sdk/">RedDog.Search .NET Client</a> and the other is the <a href="http://fabriccontroller.net/blog/posts/managing-microsoft-azure-search-with-the-reddog-search-portal/">RedDog Search Portal</a> which is a web based UI tool to manage your Azure Search service. The other one is from <a href="https://twitter.com/richorama">Richard Astbury</a>: <a href="http://coderead.wordpress.com/2014/09/02/using-the-azure-search-service-from-javascript/">Azure Search node.js / JavaScript client</a>. I strongly encourage you to check them out. There are also two great video presentations about Azure Search by <a href="https://twitter.com/liamca">Liam Cavanagh</a>, a Senior Program Manager in the Azure Data Platform Incubation team at Microsoft.</p>
<ul>
<li><a href="http://channel9.msdn.com/Shows/Data-Exposed/Introduction-To-Azure-Search">Introduction To Azure Search</a> 
<li><a href="http://channel9.msdn.com/Shows/Cloud+Cover/Cloud-Cover-152-Azure-Search-with-Liam-Cavanagh">Cloud Cover 152: Azure Search with Liam Cavanagh</a></li></ul>
<p>Stop what you are doing and go watch them if you care about Azure Search. It will give you a nice overview about the product and those videos could be your starting point.</p>
<p>You can also view my talk on AzureConf 2014 about Azure Search:</p><iframe style="height: 320px; width: 540px" src="//channel9.msdn.com/Events/Microsoft-Azure/AzureConf-2014/Search-Like-a-Pro-with-Azure-Search/player?h=320&amp;w=540" frameborder="0" scrolling="no" allowfullscreen></iframe>  