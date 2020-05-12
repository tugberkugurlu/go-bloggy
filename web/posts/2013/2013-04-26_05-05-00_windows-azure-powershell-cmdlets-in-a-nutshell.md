---
id: 81e4d746-b2c9-49ee-9462-989762ab2771
title: Windows Azure PowerShell Cmdlets In a Nutshell
abstract: Windows Azure PowerShell Cmdlets is a great tool to manage your Windows
  Azure services but if you are like me, you would wanna know where all the stuff
  is going. This post is all about it.
created_at: 2013-04-26 05:05:00 +0000 UTC
tags:
- PowerShell
- Windows Azure
slugs:
- windows-azure-powershell-cmdlets-in-a-nutshell
---

<p>Managing your <a href="http://www.windowsazure.com">Windows Azure</a> services is super easy with the various management options and my favorite among these options is <a href="http://msdn.microsoft.com/en-us/library/windowsazure/jj156055.aspx">Windows Azure PowerShell Cmdlets</a>. It's very well-documented and if you know PowerShell enough, Windows Azure PowerShell Cmdlets are really easy to grasp. In this post I would like to give a few details about this management option and hopefully, it'll give you a head start.</p>
<h3>Install it and get going</h3>
<p>Installation of the Windows Azure PowerShell Cmdlets is very easy. <a href="http://msdn.microsoft.com/en-us/library/windowsazure/jj554332.aspx">It's also well-documented</a>. You can reach the download link from the <a href="http://www.windowsazure.com/en-us/downloads/?fb=en-us">"Downloads" section on Windows Azure web site</a>. From there, all you need to do is follow the instructions to install the Cmdlets through <a href="http://www.microsoft.com/web/downloads/platform.aspx">Web Platform Installer</a>.</p>
<p>After the installation, we can view that we have the Windows Azure PowerShell Cmdlets installed on our system by running "Get-Module -ListAvailable" command:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image.png"><img height="289" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>To get you started using Cmdlets, you can see the "<a href="http://msdn.microsoft.com/en-us/library/windowsazure/jj554332.aspx">Get Started with Windows Azure Cmdlets</a>" article which explains how you will set up the trust between your machine and your Windows Azure account. However, I will cover some steps here as well.</p>
<p>First thing you need to do is download your publish settings file. There is a handy cmdlet to do this for you: Get-AzurePublishSettingsFile. By running this command, it will go to your account and create a management certificate for your account and download the necessary publish settings file.</p>
<p>Next step is to import this publish settings file through another handy cmdlet: Import-AzurePublishSettingsFile &lt;publishsettings-file-path&gt;. This command is actually setting up lots of stuff on your machine.</p>
<ul>
<li>Under "%appdata%\Windows Azure Powershell", it creates necessary configuration files for the cmdlets to get the authentication information.</li>
</ul>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_3.png"><img height="308" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; margin: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<ul>
<li>These configuration files don't actually contain certificate information on its own; they just hold the thumbprint of our management certificate and your subscription id.</li>
<li>Actual certificate is imported inside your certificate store. You can view the installed certificate by running "dir Cert:\CurrentUser\My" command.</li>
</ul>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_4.png"><img height="272" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_4.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Now you are ready to go. Run "Get-AzureSubscription" command to see your subscription details and you will see that it's set as your default subscription. So, from now on, you don't need to do anything with your subscription. You can just run the commands without worrying about your credentials (of course, this maybe a good or bad thing; depends on your situation). For example, I ran the Get-AzureVM command to view my VMs:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_5.png"><img height="281" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_5.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<h3>So, where is my stuff?</h3>
<p>We installed the stuff and we just saw that it's working. So, where did all the stuff go and how does this thing even work? Well, if you know PowerShell, you also know that modules are stored under specific folders. You can view these folders by running the '<em>$env:PSModulePath.split(';')</em>' command:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/imageba8b64ba-dacb-48c8-b906-38e38283c022.png"><img height="228" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_6.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Notice that there is a path for Windows Azure PowerShell Cmdlets, too. Without knowing this stuff, we could also view the module definition and get to its location from there:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>Get<span style="color: gray;">-</span>Module <span style="color: gray;">-</span>ListAvailable <span style="color: gray;">-</span>Name Azure</pre>
</div>
</div>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_6.png"><img height="228" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_7.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>"C:\Program Files (x86)\Microsoft SDKs\Windows Azure\PowerShell\Azure" directory is where you will find the real meat:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image.png"><img height="430" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_8.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>On the other hand, when we imported the publish settings, it puts a few stuff about my subscription under "%appdata%\Windows Azure Powershell". The certificate is also installed under my certificate store as mentioned before.</p>
<h3>Clean Up</h3>
<p>When you start managing your Windows Azure services through PowerShell Cmdlets, you have your Windows Azure account information and management certificate information at various places on your computer. Even if you uninstall your Windows Azure PowerShell Cmdlets from your machine, you are not basically cleaning up everything. Let's start by uninstalling the Cmdlets from your computer.</p>
<p>Simply go to Control Panel &gt; Programs &gt; Programs and Features and find the installed program named as Windows Azure PowerShell and uninstall it. You will be done.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_7.png"><img height="401" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_9.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Next step is to go to "%appdata%\Windows Azure Powershell" directory and delete the folder completely. One more step to go now: delete your certificate. Find out what the thumbprint of your certificate is:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image97faf2c6-2bd0-4d7b-8132-e8acdc240dcc.png"><img height="246" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/image_thumb_10.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Then, run the Remove-Item command to remove the certificate:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>Remove<span style="color: gray;">-</span>Item Cert:\CurrentUser\My\507DAAF6F285C4A72A45909ACCEE552B4E2AE916 &ndash;DeleteKey</pre>
</div>
</div>
<p>You are all done uninstalling Windows Azure PowerShell Cmdlets. Remember, Windows Azure is powerful but it's more powerful when you manage it through PowerShell <img src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/3477395e1f4d_135A1/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" /></p>
<h3>References</h3>
<ul>
<li><a href="http://msdn.microsoft.com/en-us/library/windowsazure/jj156055.aspx">Windows Azure PowerShell (MSDN)</a></li>
<li><a href="http://msdn.microsoft.com/en-us/library/windowsazure/jj554332.aspx">Get Started with Windows Azure Cmdlets</a></li>
<li><a href="http://www.windowsazure.com/en-us/develop/nodejs/how-to-guides/powershell-cmdlets">How to use Windows Azure PowerShell</a></li>
<li><a href="http://michaelwasham.com/tag/windows-azure-powershell-cmdlets">Michael Washam&rsquo;s Blog</a></li>
<li><a href="http://michaelwasham.com/2013/04/16/windows-azure-powershell-updates-for-iaas-ga/">Windows Azure PowerShell Updates for IaaS GA</a></li>
<li><a href="https://channel9.msdn.com/Shows/Cloud+Cover/Episode-105-General-Availability-of-Windows-Azure-Infrastructure-as-a-Service-IaaS">Cloud Cover Show - General Availability of Windows Azure Infrastructure as a Service (IaaS)</a></li>
</ul>