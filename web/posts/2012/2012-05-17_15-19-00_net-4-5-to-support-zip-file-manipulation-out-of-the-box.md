---
id: 29c80d1b-b6c7-40e3-af6c-38bdc5ce64bc
title: .NET 4.5 to Support Zip File Manipulation Out of the Box
abstract: One of the missing feature of .NET framework was a support for Zip file
  manipulation. In .NET 4.5, we have an extensive support for manipulating zip archives.
created_at: 2012-05-17 15:19:00 +0000 UTC
tags:
- .net
- C#
slugs:
- net-4-5-to-support-zip-file-manipulation-out-of-the-box
---

<p>One of the missing feature of .NET framework was a support for Zip file manipulation such as reading the zip archive, adding files, extracting files, etc. and we were using some third party libraries such as excellent the&nbsp;<a href="http://dotnetzip.codeplex.com/" title="http://dotnetzip.codeplex.com/">DotNetZip</a>. In .NET 4.5, we have an extensive support for manipulating .zip files.</p>
<p>First thing that you should do is to add <strong>System.IO.Compression</strong> assembly as reference to your project. You may also want to reference <strong>System.IO.Compression.FileSystem</strong> assembly to access three extension methods (from the <a href="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions(v=vs.110)" title="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions(v=vs.110)">ZipFileExtensions</a> class) for the ZipArchive class: <a href="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions.createentryfromfile(v=vs.110)" title="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions.createentryfromfile(v=vs.110)">CreateEntryFromFile</a>, <a href="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions.createentryfromfile(v=vs.110)" title="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions.createentryfromfile(v=vs.110)">CreateEntryFromFile</a>, and <a href="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions.extracttodirectory(v=vs.110)" title="http://msdn.microsoft.com/en-us/library/system.io.compression.zipfileextensions.createentryfromfile(v=vs.110)">ExtractToDirectory</a>. These extension methods enable you to compress and decompress the contents of the entry to a file.</p>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/a5e3a3f40b2a_CA41/image.png"><img height="444" width="644" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/a5e3a3f40b2a_CA41/image_thumb.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>Let&rsquo;s cover the bits and pieces that we get from <strong>System.IO.Compression</strong> assembly at first. The below sample shows how to read a zip archive easily with ZipArchive class:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args) {

    <span style="color: blue;">const</span> <span style="color: blue;">string</span> zipFilePath = <span style="color: #a31515;">@"C:\apps\Sample Pictures.zip"</span>;

    <span style="color: blue;">using</span> (FileStream zipFileToOpen = <span style="color: blue;">new</span> FileStream(zipFilePath, FileMode.Open))
    <span style="color: blue;">using</span> (ZipArchive archive = <span style="color: blue;">new</span> ZipArchive(zipFileToOpen, ZipArchiveMode.Read)) {

        <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> zipArchiveEntry <span style="color: blue;">in</span> archive.Entries)
            Console.WriteLine(
                <span style="color: #a31515;">"FullName of the Zip Archive Entry: {0}"</span>, zipArchiveEntry.FullName
            );
    }
}</pre>
</div>
</div>
<p>In this sample, we are opening the zip archive and iterate through the collection of entries. When we run the application, we should see the list of files inside the zip archive:</p>
<p><a href="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/a5e3a3f40b2a_CA41/image_3.png"><img height="327" width="644" src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/a5e3a3f40b2a_CA41/image_thumb_3.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" /></a></p>
<p>It&rsquo;s also so easy to add a new file to the zip archive:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args) {

    <span style="color: blue;">const</span> <span style="color: blue;">string</span> zipFilePath = <span style="color: #a31515;">@"C:\apps\Sample Pictures.zip"</span>;

    <span style="color: blue;">using</span> (FileStream zipFileToOpen = <span style="color: blue;">new</span> FileStream(zipFilePath, FileMode.Open))
    <span style="color: blue;">using</span> (ZipArchive archive = <span style="color: blue;">new</span> ZipArchive(zipFileToOpen, ZipArchiveMode.Update)) {

        ZipArchiveEntry readMeEntry = archive.CreateEntry(<span style="color: #a31515;">"ReadMe.txt"</span>);
        <span style="color: blue;">using</span> (StreamWriter writer = <span style="color: blue;">new</span> StreamWriter(readMeEntry.Open())) {
            writer.WriteLine(<span style="color: #a31515;">"Lorem ipsum dolor sit amet..."</span>);
            writer.Write(<span style="color: #a31515;">"Proin rutrum, massa sed molestie porta, urna..."</span>);
        }

        <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> zipArchiveEntry <span style="color: blue;">in</span> archive.Entries)
            Console.WriteLine(
                <span style="color: #a31515;">"FullName of the Zip Archive Entry: {0}"</span>, zipArchiveEntry.FullName
            );
    }
}</pre>
</div>
</div>
<p>In this sample, we are adding a file named ReadMe.txt at the root of archive and then we are writing some text into that file.</p>
<p>Extracting files is into a folder is so easy as well. You need reference the System.IO.Compression.FileSystem assembly along with <strong>System.IO.Compression </strong>assembly as mentioned before for this sample:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args) {

    <span style="color: blue;">const</span> <span style="color: blue;">string</span> zipFilePath = <span style="color: #a31515;">@"C:\apps\Sample Pictures.zip"</span>;
    <span style="color: blue;">const</span> <span style="color: blue;">string</span> dirToExtract = <span style="color: #a31515;">@"C:\apps\Sample Pictures\"</span>;

    <span style="color: blue;">using</span> (FileStream zipFileToOpen = <span style="color: blue;">new</span> FileStream(zipFilePath, FileMode.Open))
    <span style="color: blue;">using</span> (ZipArchive archive = <span style="color: blue;">new</span> ZipArchive(zipFileToOpen, ZipArchiveMode.Update))
        archive.ExtractToDirectory(dirToExtract);
}</pre>
</div>
</div>
<p>There are some other handy APIs as well but it is so easy to discover them by yourself. Enjoy <img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/a5e3a3f40b2a_CA41/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /></p>
<h3>Resources</h3>
<ul>
<li><a href="http://msdn.microsoft.com/en-us/library/system.io.compression.ziparchive(v=vs.110)" title="http://msdn.microsoft.com/en-us/library/system.io.compression.ziparchive(v=vs.110)">ZipArchive Class</a> </li>
<li><a href="http://msdn.microsoft.com/en-us/library/3z72378a(v=vs.110)" title="http://msdn.microsoft.com/en-us/library/3z72378a(v=vs.110)">System.IO.Compression Namespace</a> </li>
<li><a href="http://blogs.msdn.com/b/somasegar/archive/2012/05/16/net-improvements-for-cloud-and-server-applications.aspx" title="http://blogs.msdn.com/b/somasegar/archive/2012/05/16/net-improvements-for-cloud-and-server-applications.aspx">.NET 4.5 Improvements for Cloud and Server Applications</a> </li>
<li><a href="http://msdn.microsoft.com/en-us/library/w0x726c2(v=vs.110)" title="http://msdn.microsoft.com/en-us/library/w0x726c2(v=vs.110)">.NET Framework 4.5 Beta</a></li>
</ul>