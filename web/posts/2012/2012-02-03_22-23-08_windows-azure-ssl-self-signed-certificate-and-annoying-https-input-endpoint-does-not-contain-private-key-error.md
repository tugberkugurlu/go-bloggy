---
id: 27998ff3-9db8-4b89-956b-b2cee019c995
title: Windows Azure, SSL, Self-Signed Certificate and Annoying HTTPS Input Endpoint
  Does Not Contain Private Key Error
abstract: Deploying a Web Role with HTTPS Endpoint enabled with Self-Signed Certificate
  and a annoying problem of HTTPS input endpoint does not contain private key
created_at: 2012-02-03 22:23:08 +0000 UTC
tags:
- ASP.NET MVC
- PowerShell
- Windows Azure
slugs:
- windows-azure-ssl-self-signed-certificate-and-annoying-https-input-endpoint-does-not-contain-private-key-error
---

<p>While I was trying out the <a target="_blank" href="http://www.windowsazure.com/" title="http://www.windowsazure.com">Windows Azure</a>&nbsp;features  yesterday, I had a deployment problem. The case was to deploying SSL enabled web role. Let&rsquo;s walk  through the steps I have taken.</p>
<p>Since it was a try out, I decided to create a self-signed certificate instead  of buying one. The case of how I created the self-signed certificate was fairly  simple. I opened up the Visual Studio Command Prompt (2010) and <em>cd</em> to directory path where I would like to put the certificate file I was about to create. Then, I used the following <a href="http://msdn.microsoft.com/en-us/library/bfsktky3(v=vs.80).aspx" title="http://msdn.microsoft.com/en-us/library/bfsktky3(v=vs.80).aspx" target="_blank">makecert</a>&nbsp;command line utility to create the certificate. Here is the code I used:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>makecert <span style="color: gray;">-</span>sky exchange <span style="color: gray;">-</span>r <span style="color: gray;">-</span>n <span style="color: #a31515;">"CN=TugberkUgurlu.Com"</span> <span style="color: gray;">-</span>pe <span style="color: gray;">-</span>a sha1 <span style="color: gray;">-le</span>n 2048 <span style="color: gray;">-</span>ss My <span style="color: #a31515;">"Azure.TugberkUgurlu.Com.cer"</span></pre>
</div>
</div>
<p>I had my .cer file under <strong>c:\Azure\Certs</strong> directory. Also, I executed the following script on PowerShell console and I saw the certificate listed there.</p>
<p>&nbsp;</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>dir cert:\CurrentUser\My</pre>
</div>
</div>
<p>&nbsp;</p>
<p>In order to upload the certificate to windows azure, I needed to export the certificate from <em>.cer</em> file to <em>.pfx</em>. Here where things get messy. I used the followig powershell script to create one but I was about to realize that it was the wrong decision. We'll see why in a minute.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: orangered;">$c</span> <span style="color: gray;">=</span> New<span style="color: gray;">-</span>Object System.Security.Cryptography.X509Certificates.X509Certificate2(<span style="color: #a31515;">"c:\azure\certs\Azure.TugberkUgurlu.Com.cer"</span>)
<span style="color: orangered;">$bytes</span> <span style="color: gray;">=</span> <span style="color: orangered;">$c</span>.Export(<span style="color: #a31515;">"Pfx"</span>,<span style="color: #a31515;">"password"</span>)
<span style="color: gray;">[</span><span style="color: teal;">System.IO.File</span><span style="color: gray;">]</span><span style="color: gray;">::</span>WriteAllBytes(<span style="color: #a31515;">"c:\azure\certs\Azure.TugberkUgurlu.Com.pfx"</span>, <span style="color: orangered;">$bytes</span>)
</pre>
</div>
</div>
<p><img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/azure_certs.PNG" width="852" height="175" alt=".cer and .pfx certificates" /></p>
<p>The next step was to deploy this .pfx file to Certificates store of my windows azure hosted service. In order to complate this challange, I went to Windows Azure portal, navigated to my hosted service. Right click on the blue <strong>Certificate </strong>folder (I think it is a folder icon but not sure exactly what it is) and click <strong>Add Certificate</strong>:</p>
<p><img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/azure_hosted_service.PNG" width="838" height="209" alt="Azure Hosted Service, Add Certificate" /></p>
<p>It poped up a dialog box for me to upload that certificate file. I completed the steps and there it was. I had the certificate deployed on my hosted service.</p>
<p>Finally, I was done setting things up and I can jump right to my application. Wait, I wasn't done yet complately! I had to set things up at the application level so that I could hook it up to that certificate I had just uploded.</p>
<p>At that stage, first thing I did was to grab the thumbprint of the certificate. I ran the following PowerShell command to grab the thumbprint of the certificate:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>(New<span style="color: gray;">-</span>Object System.Security.Cryptography.X509Certificates.X509Certificate2(<span style="color: #a31515;">"c:\azure\certs\Azure.TugberkUgurlu.Com.cer"</span>)).Thumbprint
</pre>
</div>
</div>
<p>With that thumbprint, I went to my project and added the following code inside ServiceConfiguration.Clound.cscfg and ServiceConfiguration.Local.cscfg under <strong>Role </strong>node:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">Certificates</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;</span><span style="color: #a31515;">Certificate</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Azure.TugberkUgurlu.Com</span><span style="color: black;">"</span> <span style="color: red;">thumbprint</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">AAD5DDD0CA9B4D3CFEF1652130142020770B8BDF</span><span style="color: black;">"</span> <span style="color: red;">thumbprintAlgorithm</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">sha1</span><span style="color: black;">"</span><span style="color: blue;">/&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">Certificates</span><span style="color: blue;">&gt;</span>
</pre>
</div>
</div>
<p>The ServiceDefinition.csdef file needed a little more touch than the configuration files. Here is the complate csdef file after the set-up:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">ServiceDefinition</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">HttpsEnabledCloudProject</span><span style="color: black;">"</span> <span style="color: red;">xmlns</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">http://schemas.microsoft.com/ServiceHosting/2008/10/ServiceDefinition</span><span style="color: black;">"</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;</span><span style="color: #a31515;">WebRole</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">MvcWebRole1</span><span style="color: black;">"</span> <span style="color: red;">vmsize</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Small</span><span style="color: black;">"</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">Sites</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">Site</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Web</span><span style="color: black;">"</span><span style="color: blue;">&gt;</span>
        <span style="color: blue;">&lt;</span><span style="color: #a31515;">Bindings</span><span style="color: blue;">&gt;</span>
          <span style="color: blue;">&lt;</span><span style="color: #a31515;">Binding</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Endpoint1</span><span style="color: black;">"</span> <span style="color: red;">endpointName</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Endpoint1</span><span style="color: black;">"</span> <span style="color: red;">hostHeader</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">azure.tugberkugurlu.com</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
          <span style="color: blue;">&lt;</span><span style="color: #a31515;">Binding</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">HttpsIn</span><span style="color: black;">"</span> <span style="color: red;">endpointName</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">HttpsIn</span><span style="color: black;">"</span> <span style="color: red;">hostHeader</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">azure.tugberkugurlu.com</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
        <span style="color: blue;">&lt;/</span><span style="color: #a31515;">Bindings</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;/</span><span style="color: #a31515;">Site</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">Sites</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">Endpoints</span><span style="color: blue;">&gt;</span>
      &lt;InputEndpoint name="Endpoint1" protocol="http" port="80" /&gt;
      &lt;InputEndpoint name="HttpsIn" protocol="https" port="443" certificate="Azure.TugberkUgurlu.Com" /&gt;
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">Endpoints</span><span style="color: blue;">&gt;</span>
    &lt;Imports&gt;
      &lt;Import moduleName="Diagnostics" /&gt;
    &lt;/Imports&gt;
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">Certificates</span><span style="color: blue;">&gt;</span>
      <span style="color: blue;">&lt;</span><span style="color: #a31515;">Certificate</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">Azure.TugberkUgurlu.Com</span><span style="color: black;">"</span> <span style="color: red;">storeLocation</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">LocalMachine</span><span style="color: black;">"</span> <span style="color: red;">storeName</span><span style="color: blue;">=</span><span style="color: black;">"</span><span style="color: blue;">My</span><span style="color: black;">"</span> <span style="color: blue;">/&gt;</span>
    <span style="color: blue;">&lt;/</span><span style="color: #a31515;">Certificates</span><span style="color: blue;">&gt;</span>
  <span style="color: blue;">&lt;/</span><span style="color: #a31515;">WebRole</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">ServiceDefinition</span><span style="color: blue;">&gt;</span>
</pre>
</div>
</div>
<p>Everything looked right at that point but got an ugly error message telling me that my deployment had been failed:</p>
<blockquote>
<p><em>HTTP Status Code: 400. Error Message: Certificate with thumbprint AAD5DDD0CA9B4D3CFEF1652130142020770B8BDF associated with HTTPS input endpoint HttpsIn does not contain private key</em></p>
</blockquote>
<p><img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/windows_azure_vs_deployment_error.PNG" width="743" height="316" alt="Windows Azure Visual Studio Deployment Error" /></p>
<p>After a couple of searched on the internet, I ended up checking the private key of my certificate and here is the result:</p>
<p><img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/ps_cert_hasprivatekey.PNG" width="681" height="398" alt="$cert.HasPrivateKey" /></p>
<p>That looked awkward and might be the problem. Then I grabed the certificate from the certificate store and check the private key. That was the evidence that it was the problem:</p>
<p><img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/ps_cert_hasprivatekey_2.PNG" width="743" height="314" alt="$cert.HasPrivateKey" /></p>
<p>At that point, I created the .pfx file from the certificate in my certificate store with following code:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: orangered;">$cert</span> <span style="color: gray;">=</span> (get<span style="color: gray;">-</span>item cert:\CurrentUser\My\AAD5DDD0CA9B4D3CFEF1652130142020770B8BDF)
<span style="color: orangered;">$bytes</span> <span style="color: gray;">=</span> <span style="color: orangered;">$cert</span>.Export(<span style="color: #a31515;">"Pfx"</span>,<span style="color: #a31515;">"password"</span>)
<span style="color: gray;">[</span><span style="color: teal;">System.IO.File</span><span style="color: gray;">]</span><span style="color: gray;">::</span>WriteAllBytes(<span style="color: #a31515;">"c:\azure\certs\Azure.TugberkUgurlu.Com_2.pfx"</span>, <span style="color: orangered;">$bytes</span>)
</pre>
</div>
</div>
<p>At last, I deleted the certificate which was under my hosted service and reuploaded the new one I had just created.</p>
<p>Lastly, I ran the publish process again without changing anything inside my code, configuration or service definition files and I suceedded!</p>
<p><img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/windows_azure_vs_deployment.PNG" width="743" height="318" alt="Windows Azure Visual Studio Deployment" /></p>
<p>Here is the SSL enabled application running in the cloud:</p>
<p><img src="https://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/azure_ssl_app.PNG" width="686" height="429" alt="Windows Azure SSL Enable Application" /></p>
<p>Please comment if you have the same problem as me so that I won't feel lonely in this small World:)</p>