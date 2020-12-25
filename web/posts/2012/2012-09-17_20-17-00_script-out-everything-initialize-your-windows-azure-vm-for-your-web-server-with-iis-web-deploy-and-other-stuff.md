---
id: 6839d3bc-2fe9-4f30-a0cf-b1a9a0865885
title: Script Out Everything - Initialize Your Windows Azure VM for Your Web Server
  with IIS, Web Deploy and Other Stuff
abstract: Script Out Everything - Initialize Your Windows Azure VM for Your Web Server
  with IIS, Web Deploy and Other Stuff
created_at: 2012-09-17 20:17:00 +0000 UTC
tags:
- IIS
- PowerShell
- Windows Azure
slugs:
- script-out-everything-initialize-your-windows-azure-vm-for-your-web-server-with-iis-web-deploy-and-other-stuff
---

<p>Today, I am officially sick and tired of initializing a new Web Server every time so I decided to script it all out as much as I can. I created a new Windows Azure VM running Windows Server 2012 and installed&nbsp;<a href="http://www.iis.net/learn/install/web-platform-installer/web-platform-installer-v4-command-line-webpicmdexe-rtw-release">Web Platform Installer v4 Command Line</a> tool (aka WebPICMD.exe).</p>
<p>Next thing I need to do is to install the IIS Web Server Role inside my VM. To do that, I opened up the Server Manager and closed it instantly because I didn't wanna do that through a GUI either. It turned out that it is fairly easy to manage your server roles and features through PowerShell thanks to ServerManager PowerShell module. This module has couple of handy Cmdlets which enable you to manage your server&rsquo;s roles and features.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image.png"><img height="230" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>To install the IIS Web Server, I run the following PowerShell command. It installs IIS Web Server (along with the necessary dependencies and management tools) and logs the output under the TEMP folder.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: orangered;">$logLabel</span> <span style="color: gray;">=</span> $((get<span style="color: gray;">-</span>date).ToString(<span style="color: #a31515;">"yyyyMMddHHmmss"</span>))
Import<span style="color: gray;">-</span>Module <span style="color: gray;">-</span>Name ServerManager
Install<span style="color: gray;">-</span>WindowsFeature <span style="color: gray;">-</span>Name Web<span style="color: gray;">-</span>Server <span style="color: gray;">-</span>IncludeManagementTools <span style="color: gray;">-</span>LogPath <span style="color: #a31515;">"$env:TEMP\init-webservervm_webserver_install_log_$logLabel.txt"</span></pre>
</div>
</div>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_3.png"><img height="152" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>When we look at the installed features, we should see that IIS Web Server is now listed there:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_4.png"><img height="369" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_thumb_4.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>At this stage, we have a few more necessary features to install such as ASP.NET 4.5 and Management Service. Optionally, I always want to install Dynamic Content Compression and&nbsp;IIS Management Scripts and Tools features.&nbsp;To install those features, we will run the following script:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: green;">#add additional windows features</span>
<span style="color: orangered;">$additionalFeatures</span> <span style="color: gray;">=</span> @(<span style="color: #a31515;">'Web-Mgmt-Service'</span>, <span style="color: #a31515;">'Web-Asp-Net45'</span>, <span style="color: #a31515;">'Web-Dyn-Compression'</span>, <span style="color: #a31515;">'Web-Scripting-Tools'</span>)
<span style="color: blue;">foreach</span>(<span style="color: orangered;">$feature</span> <span style="color: blue;">in</span> <span style="color: orangered;">$additionalFeatures</span>) { 
    
    <span style="color: blue;">if</span>(!(Get<span style="color: gray;">-</span>WindowsFeature | where { <span style="color: orangered;">$_</span>.Name <span style="color: gray;">-eq</span> <span style="color: orangered;">$feature</span> }).Installed) { 

        Install<span style="color: gray;">-</span>WindowsFeature <span style="color: gray;">-</span>Name <span style="color: orangered;">$feature</span> <span style="color: gray;">-</span>LogPath <span style="color: #a31515;">"$env:TEMP\init-webservervm_feature_$($feature)_install_log_$((get-date).ToString("</span>yyyyMMddHHmmss<span style="color: #a31515;">")).txt"</span>   
    }
}
<span style="color: green;">#Set WMSvc to Automatic Startup</span>
Set<span style="color: gray;">-</span>Service <span style="color: gray;">-</span>Name WMSvc <span style="color: gray;">-</span>StartupType Automatic

<span style="color: green;">#Check if WMSvc (Web Management Service) is running</span>
<span style="color: blue;">if</span>((Get<span style="color: gray;">-</span>Service WMSvc).Status <span style="color: gray;">-ne</span> <span style="color: #a31515;">'Running'</span>) { 
    Start<span style="color: gray;">-</span>Service WMSvc
}</pre>
</div>
</div>
<p>As you can see at the end of the script, we are also setting the&nbsp;Management Service's startup type to automatic and finally, we are starting the service.</p>
<p>Now the IIS Web Server is ready, we need to get a few more bits through the Web Platform Installer v4 Command Line tool. With all of my web servers, there are two concrete tools I would like to have: <a href="http://www.iis.net/downloads/microsoft/web-deploy">Web Deploy</a> and <a href="https://www.tugberkugurlu.com/archive/remove-trailing-slash-from-the-urls-of-your-asp-net-web-site-with-iis-7-url-rewrite-module">URL Rewrite Module</a>. We can certainly install those manually but we can also script out their installation. WebPICMD.exe allows us to install products through command line and the following command will work for us:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: orangered;">$webPiProducts</span> <span style="color: gray;">=</span> @(<span style="color: #a31515;">'WDeployPS'</span>, <span style="color: #a31515;">'UrlRewrite2'</span>)
.\WebPICMD.exe <span style="color: gray;">/</span>Install <span style="color: gray;">/</span>Products:<span style="color: #a31515;">"$($webPiProducts -join ',')"</span> <span style="color: gray;">/</span>AcceptEULA <span style="color: gray;">/</span>Log:<span style="color: #a31515;">"$env:TEMP\webpi_products_install_log_$((get-date).ToString("</span>yyyyMMddHHmmss<span style="color: #a31515;">")).txt"</span></pre>
</div>
</div>
<p>I assume that WebPICMD.exe is under your path here (if you have installed the 64x version of the product, you can find the executable file under C:\Program Files\Microsoft\Web Platform Installer). When the installation is complete, we should see the success message:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_5.png"><img height="343" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_thumb_5.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Well, we have completed most of the work but there are still couple of things to do. First of all, we need to allow incoming connections through TCP port 8172 because this is the port that Web Deploy will talk through. To enable that, we can go to Windows Firewall with Advanced Security window but that would be lame. Are we gonna use <a href="http://www.microsoft.com/resources/documentation/windows/xp/all/proddocs/en-us/netsh.mspx?mfr=true">netsh</a>? Certainly not :) With PowerShell 3.0, we can now control the <a href="http://technet.microsoft.com/en-us/library/hh831755.aspx">Windows Firewall with Advanced Security Administration</a>. This functionality is provided through NetSecurity PowerShell module but with the new dynamic module loading feature of PowerShell 3.0, we don&rsquo;t need to separately import this. The following command will add the proper firewall rule to our server.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>New<span style="color: gray;">-</span>NetFirewallRule <span style="color: gray;">-</span>DisplayName <span style="color: #a31515;">"Allow IIS Management Service In"</span> <span style="color: gray;">-</span>Direction Inbound <span style="color: gray;">-</span>LocalPort 8172 <span style="color: gray;">-</span>Protocol TCP <span style="color: gray;">-</span>Action Allow</pre>
</div>
</div>
<p>Lastly, we need to inform windows azure about this firewall rule as well because all the requests, which come from outside, will go through the load balancer and it doesn&rsquo;t open up any ports by default (except for Remote Desktop ports). You can the endpoints through Windows Azure Portal but <a href="http://msdn.microsoft.com/en-us/library/windowsazure/jj156055.aspx">Windows Azure PowerShell Cmdlets</a> to add this as well. The following command will add the proper rule. Just change the $serviceName and $vmName according to your credentials:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>Get<span style="color: gray;">-</span>AzureVM <span style="color: gray;">-</span>ServiceName <span style="color: orangered;">$serviceName</span> <span style="color: gray;">-</span>Name <span style="color: orangered;">$vmName</span> | Add<span style="color: gray;">-</span>AzureEndpoint <span style="color: gray;">-</span>Name <span style="color: #a31515;">"WebDeploy"</span> <span style="color: gray;">-</span>Protocol TCP <span style="color: gray;">-</span>LocalPort 8172 <span style="color: gray;">-</span>PublicPort 8172 | Update<span style="color: gray;">-</span>AzureVM</pre>
</div>
</div>
<p>When we look at the portal, we should see that our endpoint was created.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/AzureEndpoint.png"><img height="191" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/AzureEndpoint_thumb.png" alt="AzureEndpoint" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="AzureEndpoint" /></a></p>
<p>Now, everything should be working perfectly. Of course you also need to add TCP port 80 to your Endpoint lists for your VM in order for your web sites to be reachable through HTTP (assuming you will only use PORT 80 for your web applications). To test everything out, I created a web application under IIS. Then, I right clicked on it and navigate to Deploy &gt; Configure Web Deploy Publishing.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_6.png"><img height="484" width="581" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_thumb_6.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>This will bring up another dialog. This is the place where we can configure web deploy settings. There are other ways to do this as well.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_7.png"><img height="240" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/image_thumb_7.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Just change the VM name with your VIP address and this will generate a publish profile. We can now use this publish profile file to push our web application to the server.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/WebDeployPublishVS.png"><img height="191" width="244" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/WebDeployPublishVS_thumb.png" alt="WebDeployPublishVS" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="WebDeployPublishVS" /></a></p>
<p>When the publish is completed, we will see the complete result inside the Output window.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/WebDeployPublishVS_2.png"><img height="177" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Script-Out-Everything----Initialize-Your_150B7/WebDeployPublishVS_2_thumb.png" alt="WebDeployPublishVS_2" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="WebDeployPublishVS_2" /></a></p>
<p>We jumped through lots of hoops to get this done but how is this better than doing it manually? Imagine that you bring all these scripts together and run them through multiple VMs. It&rsquo;s going to save a lot of time for you. I am not an IT pro. So, this is enough to make me happy because I just proved that nearly everything that a web developer needs can be automated through PowerShell or various command line tools.</p>
<h3>Update:</h3>
<p>I put these together inside one script file. Follow the instructions inside the script and you will be good to go.</p>
<p><a href="https://gist.github.com/3742921#file_init_web_server_vm.ps1">https://gist.github.com/3742921#file_init_web_server_vm.ps1</a></p>
<script type="text/javascript"></script>
<p>This doesn't add the endpoints for your VM. For this, you have another script as well:</p>
<p><a href="https://gist.github.com/3742921#file_init_web_server_azure_vm_endpoints.ps1">https://gist.github.com/3742921#file_init_web_server_azure_vm_endpoints.ps1</a></p>