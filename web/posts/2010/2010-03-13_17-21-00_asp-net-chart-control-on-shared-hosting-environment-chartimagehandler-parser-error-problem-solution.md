---
id: 4f4a26ca-e44e-4573-9717-0d0434a63c61
title: ASP.Net Chart Control On Shared Hosting Environment, Chartimagehandler / Parser
  Error Problem Solution
abstract: ASP.Net Chart Control On Shared Hosting Envronment (Parser Error) - Deployment
  of Chart Controls On Shared Hosting Environment Properly...
created_at: 2010-03-13 17:21:00 +0000 UTC
tags:
- .NET
- ASP.Net
- Deployment
- Hosting
slugs:
- asp-net-chart-control-on-shared-hosting-environment-chartimagehandler-parser-error-problem-solution
---

<p><span style="font-style: italic; color: #339966;">[This Article is for those who can manage to run Chart Control on Local Server but not on Web Server]</span><br /> <br /> Chart Control of ASP.Net 3.5 is very handy way to get statictical data as charts. There are variety of ways to display data with Chart Control. Also it is extramely simple to adapt it to your project. But as I hear from so many ASP.Net user nowadays, there might be some problems when it comes to publish it on your web server, especially on <span style="font-weight: bold;">shared hosting servers</span> <span style="font-weight: bold;">!<br /> <br /> </span>After you install the Chart Control,&nbsp;in order to&nbsp;run ASP.Net chart control on your project, you need to configure your web.config file&nbsp;on your root directory which you have already&nbsp;done (I guess !). If you do that properly, you are able to run it on your local server perfectly. But it is not enough for you to be good to go on your web server. <br /> <br /> <span style="font-size: 14pt;">You need to have Chart Control installed on your web server in order to run it on web !<br /> <br /> </span>It is easy to do that if you have your own server but for those who have their website files on a shared hosting service, it is a bit hard. But you do not need to beg your hosting provider to make it be installed on server. Only you need to do is to make proper changes on your web.config file and believe me&nbsp;those changes are much more&nbsp;simple than you think ! Of Course, some references is needed to be added to the <span style="font-style: italic;">Bin </span>Folder on your root directory !<br /> <br /> <span style="color: #009900; font-weight: bold;">Solution :</span><br /> <br /></p>
<ol>
<li>Follow this directory on windows explorer:<br /> <br /> <span style="color: black; font-size: 12pt;">C:\Program Files\Microsoft Chart Controls\Assemblies</span><br /> <br /> </li>
<li>You will see 4 <span style="font-style: italic;">DLL </span>files inside the folder. Two of those files are for Windows Applications and two of them for <span style="font-weight: bold;">Web Applications</span>. We need web application dll files so copy the dll files which are named <span style="font-style: italic;">'System.Web.DataVisualization.Design'</span> and <span style="font-style: italic;">'System.Web.DataVisualization'</span><br /> <br /> </li>
<li>Paste&nbsp;those dll files&nbsp;into the <span style="font-style: italic;">Bin </span>folder on the root directory of your Web Application.<br /> <br /> </li>
<li>Secondly, open the Web.Config file on the root directory of your web application.</li>
<li>Find the <span size="2" color="#0000ff" style="color: #0000ff; font-size: x-small;"><span size="2" color="#0000ff" style="color: #0000ff; font-size: x-small;">&lt;</span></span><span size="2" color="#a31515" style="color: #a31515; font-size: x-small;"><span size="2" color="#a31515" style="color: #a31515; font-size: x-small;">appSettings</span></span><span size="2" color="#0000ff" style="color: #0000ff; font-size: x-small;"><span size="2" color="#0000ff" style="color: #0000ff; font-size: x-small;">&gt; </span></span>node. You have a key add inside this&nbsp;just like&nbsp;below;<br /> <br /> <span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">&lt;</span></span><span color="#a31515" style="color: #a31515;"><span color="#a31515" style="color: #a31515;">add</span></span><span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;"> </span></span><span color="#ff0000" style="color: #ff0000;"><span color="#ff0000" style="color: #ff0000;">key</span></span><span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">=</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">ChartImageHandler</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;"> </span></span><span color="#ff0000" style="color: #ff0000;"><span color="#ff0000" style="color: #ff0000;">value</span></span><span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">=</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">storage=file;timeout=20;dir=c:\TempImageFiles\;</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;"> /&gt;<br /> </span></span><br /> Replace&nbsp;this tag with the new one as below;<br /> <br /> <span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">&lt;</span></span><span color="#a31515" style="color: #a31515;"><span color="#a31515" style="color: #a31515;">add</span></span><span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;"> </span></span><span color="#ff0000" style="color: #ff0000;"><span color="#ff0000" style="color: #ff0000;">key</span></span><span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">=</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">ChartImageHandler</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;"> </span></span><span color="#ff0000" style="color: #ff0000;"><span color="#ff0000" style="color: #ff0000;">value</span></span><span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">=</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">storage=file;timeout=20;</span></span>"<span color="#0000ff" style="color: #0000ff;"><span color="#0000ff" style="color: #0000ff;">/&gt;<br /> </span></span><br /> Save the changes on your <span style="font-style: italic;">Web.Config</span> file and close it. Now copy the two dll inside the bin folder and replace the Web.Config file on your server with your new Web.Config file.</li>
</ol>
<p><br /> <span style="font-size: 12pt;">That's it !<br /> <br /> </span>Now you should be able to run the Chart Control on your Web Server without begging your hosting provider :) If you have any problem with this, I recommend you to check the following codes if they exist on your Web.Config file or not;</p>
<p>&nbsp;</p>
<pre class="brush: xhtml">...

<system.web>

  ...

  <pages>
    <controls>
      <add tagprefix="asp" namespace="System.Web.UI.DataVisualization.Charting" assembly="System.Web.DataVisualization, Version=3.5.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35">
    </add></controls>
  </pages>

  ...

  <httphandlers>
    <add path="ChartImg.axd" verb="GET,HEAD" type="System.Web.UI.DataVisualization.Charting.ChartHttpHandler, System.Web.DataVisualization, Version=3.5.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" validate="false">
  </add></httphandlers>

  ...

</system.web>

...

<system.webserver>
  <handlers>
    <remove name="ChartImageHandler">
    <add name="ChartImageHandler" precondition="integratedMode" verb="GET,HEAD" path="ChartImg.axd" type="System.Web.UI.DataVisualization.Charting.ChartHttpHandler, System.Web.DataVisualization, Version=3.5.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35">
  </add></remove></handlers>
</system.webserver>

... </pre>
<p>&nbsp;</p>