---
title: How to Use Windows Azure Blob Storage Service With ASP.NET MVC Web Application
abstract: This blog post provides a Windows Azure Blob Storage example which walks
  you through on how to use Blob Storage service with ASP.NET MVC.
created_at: 2012-01-18 07:17:00 +0000 UTC
tags:
- ASP.NET MVC
- Blob Storage
- Windows Azure
slugs:
- how-to-use-windows-azure-blob-storage-service-with-asp-net-mvc-web-application
- how-to-use-windows-azure-blob-storage-service-with-asp-net-mvc-
---

<p>I have been digging into <a target="_blank" href="http://www.windowsazure.com" title="http://www.windowsazure.com">Windows Azure</a> more and more lately. I wish that it would be supported in Turkey but anyway, emulator is my cloud for now <img src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/0310308e59b3_10C04/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /> Yesterday, I thought I should note some of things down and told myself "There is no better place than your blog for that, my friend" and here I am.</p>
<p>One feature of Windows Azure platform which I will be using is <strong>Blob Storage</strong>. Blob storage enables you to store your unstructured data (like pictures, word docs, excel file, etc.) inside Windows Azure servers and access them over HTTP or HTTPS. <a target="_blank" href="https://github.com/WindowsAzure/azure-sdk-for-net" title="https://github.com/WindowsAzure/azure-sdk-for-net">With Windows Azure .Net SDK</a>, you have full control over your blobs and program against that easily. How? Let&rsquo;s see.</p>
<blockquote>
<p>Before starting, make sure that you have installed Windows Azure SDK for .Net which will bring down Windows Azure Tools for Microsoft Visual Studio and Windows Azure Client Libraries for .Net. You can find the information on how to install the SDK from <a href="http://www.windowsazure.com/en-us/develop/net/">http://www.windowsazure.com/en-us/develop/net/</a></p>
</blockquote>
<p>First thing is first. We need an ASP.NET MVC project to simulate this (but it doesn&rsquo;t have to be ASP.NET MVC project). We have two options to make our application azure-cloudy:</p>
<p><strong>Directly Create a Cloud Application</strong></p>
<p>Inside the new project dialog box on Visual Studio, choose Windows Azure Project as indicated below:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image.png"><img height="394" width="644" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>When you hit OK, you will see a dialog as below:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_3.png"><img height="404" width="643" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_3.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>This dialog box is like a open buffet, you can choose which project you need for your application here. But we will choose ASP.NET MVC 3 Web Role and then hit OK:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_4.png"><img height="404" width="643" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_4.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Then we will see a dialog box which is owned by ASP.NET MVC. From now on everything is same like it is a standard ASP.NET MVC project.</p>
<p><strong>Make Your Application Azure-Cloudy Later</strong></p>
<p>Assuming that we have an existing ASP.NET MVC application and we want to run this application on Windows Azure. What we need to do is to right click on our project and choose <em>"Add Windows Azure Deployment Project Option" </em>as below:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_5.png"><img height="303" width="596" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_5.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Either way, our solution will look something like this:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_6.png"><img height="431" width="324" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_6.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p><strong>Configure to Work With Blob Storage</strong></p>
<p>We have a few steps to make before we can start developing. In a real world scenario, you need a Windows Azure storage account to use the blob storage service and you can create this account from <a href="http://windows.azure.com/">Windows Azure Management Portal</a>. After you configure your account, you will have your access keys to that storage account which you will need on your development process.</p>
<p>One thing to mention before going further is that you will be able to access your files through HTTP or HTTPS as motioned before and the URL for your blobs will look like this:</p>
<p><em>http://&lt;storage account&gt;.blob.core.windows.net/&lt;container&gt;/&lt;blob&gt;</em></p>
<p>There is a way to change this so that you can use your own domain. In our example, we will be reaching out our blobs through localhost. we will get to that later in this post.</p>
<p>In order to develop locally with emulator, we do not need a storage account which means that we don&rsquo;t need access keys. In order to configure it so, right click on the web role file and choose <em>Properties</em>.</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_7.png"><img height="378" width="644" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_7.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>This action will brings up the properties windows. From there, go to <em>Settings</em> tab and click add settings and on the new created node, select Connection String as Type and click "..." which stands right hand side.</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_8.png"><img height="378" width="644" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_8.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>This action will also open up a new window for you and this is where you configure your storage account. But we will choose "<em>Use the Windows Azure storage emulator</em>" option and give this configuration a new friendly name:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_9.png"><img height="351" width="644" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_9.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Those actions are not required but viewed as a best practice while working windows azure cloud projects. Also, you can use Windows Azure Storage and run your application inside your own servers. It is totally fine. So, in that case you won&rsquo;t need a cloud project. What .Net SDK provides is a wrapper around Windows Azure REST APIs which makes it easy to program against.</p>
<p>Now we are all set and finally we can write some code.</p>
<p>First of all, we need two additional libraries to develop against Windows Azure with .Net:</p>
<ul>
<li>Microsoft.WindowsAzure.StorageClient.dll </li>
<li>Microsoft.WindowsAzure.ServiceRuntime.dll</li>
</ul>
<p>Those two will give us everything we need. For the sake of simplicity, I created a simple project which uploads images through Windows Azure Storage service and list those images on a page. In order to do that so, I created a service class (which is a standard class, nothing further than that) called <strong>MyBlobStorageService</strong>. Let&rsquo;s see the code first:</p>
<p>&nbsp;</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<p>&nbsp;</p>
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> MyBlobStorageService {

    <span style="color: blue;">public</span> CloudBlobContainer GetCloudBlobContainer() {

        <span style="color: green;">// Retrieve storage account from connection-string</span>
        CloudStorageAccount storageAccount = CloudStorageAccount.Parse(
                RoleEnvironment.GetConfigurationSettingValue(<span style="color: #a31515;">"StorageConnectionString"</span>)
            );

        <span style="color: green;">// Create the blob client </span>
        CloudBlobClient blobClient = storageAccount.CreateCloudBlobClient();

        <span style="color: green;">// Retrieve a reference to a container </span>
        CloudBlobContainer blobContainer = blobClient.GetContainerReference(<span style="color: #a31515;">"albums"</span>);

        <span style="color: green;">// Create the container if it doesn't already exist</span>
        <span style="color: blue;">if</span> (blobContainer.CreateIfNotExist()) {

            blobContainer.SetPermissions(
               <span style="color: blue;">new</span> BlobContainerPermissions { PublicAccess = BlobContainerPublicAccessType.Blob }
            );
        }

        <span style="color: blue;">return</span> blobContainer;
    }

}</pre>
</div>
</div>
<p>What this code does it fairly simple. It gets the storage account information form the connection string that we have configured and creates a <strong>CloudStorageAccount</strong> class according to that. Then, We create a blob storage client (<strong>CloudBlobClient</strong>) over that storage account. Finally, we create a container (<strong>CloudBlobContainer</strong>) for our blobs and check if it exists of not. If not, then we simply create it and set the public access permission to it because we will store images inside that container and we want to display those images on our web page.</p>
<p>In order to demonstrate this, I have two controller actions. One is for HTTP GET and the other is for HTTP POST. Here is the complete controller code:</p>
<p>&nbsp;</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller {

    <span style="color: blue;">public</span> MyBlobStorageService 
        _myBlobStorageService = <span style="color: blue;">new</span> MyBlobStorageService();

    <span style="color: blue;">public</span> ActionResult Index() {

        <span style="color: green;">// Retrieve a reference to a container </span>
        CloudBlobContainer blobContainer = 
            _myBlobStorageService.GetCloudBlobContainer();

        List&lt;<span style="color: blue;">string</span>&gt; blobs = <span style="color: blue;">new</span> List&lt;<span style="color: blue;">string</span>&gt;();

        <span style="color: green;">// Loop over blobs within the container and output the URI to each of them</span>
        <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> blobItem <span style="color: blue;">in</span> blobContainer.ListBlobs())
            blobs.Add(blobItem.Uri.ToString());

        <span style="color: blue;">return</span> View(blobs);
    }

    [HttpPost]
    [ActionName(<span style="color: #a31515;">"index"</span>)]
    <span style="color: blue;">public</span> ActionResult Index_post(HttpPostedFileBase fileBase) {

        <span style="color: blue;">if</span> (fileBase.ContentLength &gt; 0) {

            <span style="color: green;">// Retrieve a reference to a container </span>
            CloudBlobContainer blobContainer = 
                _myBlobStorageService.GetCloudBlobContainer();

            CloudBlob blob = 
                blobContainer.GetBlobReference(fileBase.FileName);
            
            <span style="color: green;">// Create or overwrite the "myblob" blob with contents from a local file</span>
            blob.UploadFromStream(fileBase.InputStream);

        }

        <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"index"</span>);
    }   
}</pre>
</div>
</div>
<p>And the view code is simple as well:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@model List<span style="color: blue;">&lt;</span><span style="color: #a31515;">string</span><span style="color: blue;">&gt;</span>           
@{
    ViewBag.Title = "My Cloudy Album";
}

<span style="color: blue;">&lt;</span><span style="color: #a31515;">h2</span><span style="color: blue;">&gt;</span>My Cloudy Album<span style="color: blue;">&lt;/</span><span style="color: #a31515;">h2</span><span style="color: blue;">&gt;</span>

@foreach (var item in Model) {
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">img</span> <span style="color: red;">src</span><span style="color: blue;">=</span><span style="color: blue;">"@item"</span> <span style="color: red;">width</span><span style="color: blue;">=</span><span style="color: blue;">"200"</span> <span style="color: red;">height</span><span style="color: blue;">=</span><span style="color: blue;">"100"</span> <span style="color: blue;">/&gt;</span>
}

@using (Html.BeginForm("index", "home",  
    FormMethod.Post, new { enctype = "multipart/form-data" })) {
    
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">input</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: blue;">"file"</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: blue;">"fileBase"</span> <span style="color: blue;">/&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">p</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">input</span> <span style="color: red;">type</span><span style="color: blue;">=</span><span style="color: blue;">"submit"</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"Upload"</span> <span style="color: blue;">/&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">p</span><span style="color: blue;">&gt;</span>    
}</pre>
</div>
</div>
<p>Now we are all set and go to go.</p>
<blockquote>
<p>But (a big but), you need to run Visual Studio with admin rights in order to run the emulator. Otherwise, you won&rsquo;t be able to.</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_10.png"><img height="217" width="311" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_10.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
</blockquote>
<p>When we hit CTRL + F5, VS will create a cloud deployment package and deploy it to emulator and we will see the emulator starting up:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_11.png"><img height="41" width="247" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_11.png" alt="image" border="0" title="image" style="padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>Now we are cloudy!</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_12.png"><img height="358" width="644" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_12.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>When we choose a picture and upload it, we should be able to see it after we got back.</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_13.png"><img height="358" width="644" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_13.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>One more thing to prove that we really run inside the emulator is to right click the emulator icon, choose Show Storage Emulator UI:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_14.png"><img height="191" width="244" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_14.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>This will bring up the Storage Emulator as below:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_15.png"><img height="209" width="607" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/Windows-Azure-Blob-Storage-and-A.NET-MVC_FA59/image_15.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>On the right hand side, there is a reset button. Click it to reset the emulator and go back to your page and refresh. You will see that all the pictures that you have uploaded are gone now.</p>
<p>That&rsquo;s all for now <img src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/0310308e59b3_10C04/wlEmoticon-smile.png" alt="Winking smile" class="wlEmoticon wlEmoticon-winkingsmile" style="border-style: none;" /></p>
<p><strong>Additional Resources</strong></p>
<ul>
<li><a target="_blank" href="http://www.windowsazure.com" title="http://www.windowsazure.com">Windows Azure Web Site</a></li>
<li><a href="http://blogs.msdn.com/b/windowsazurestorage/">Windows Azure Storage Team Blog</a> </li>
<li><a target="_blank" href="http://www.windowsazure.com/en-us/develop/net/how-to-guides/blob-storage/" title="http://www.windowsazure.com/en-us/develop/net/how-to-guides/blob-storage/">How to Use the Blob Storage Service</a></li>
</ul>