---
title: NGINX Reverse Proxy and Load Balancing for ASP.NET 5 Applications with Docker
  Compose
abstract: In this post, I want to show you how it would look like to expose ASP.NET
  5 through NGINX, provide a simple load balancing mechanism running locally and orchestrate
  this through Docker Compose.
created_at: 2016-01-17 14:39:00 +0000 UTC
tags:
- ASP.NET 5
- Docker
- Docker Compose
- NGINX
slugs:
- nginx-reverse-proxy-and-load-balancing-for-asp-net-5-applications-with-docker-compose
- nginx-reverse-proxy-and-load-balancing-for-asp-net-5-applicatio
---

<p>We have a lot of hosting options with ASP.NET 5 under different operating systems and under different web servers like <a href="https://www.iis.net/">IIS</a>. <a href="https://twitter.com/filip_woj">Filip W</a> has a great blog post on <a href="http://www.strathweb.com/2015/12/running-asp-net-5-website-on-iis/">Running ASP.NET 5 website under IIS</a>. Here, I want to show you how it would look like to expose <a href="http://www.tugberkugurlu.com/tags/asp-net-5">ASP.NET 5</a> through <a href="https://www.nginx.com/resources/wiki/">NGINX</a>, provide a simple load balancing mechanism running locally and orchestrate this through <a href="https://docs.docker.com/compose/">Docker Compose</a>.</p> <blockquote lang="en" class="twitter-tweet"> <p lang="en" dir="ltr">sweet, finally! <a href="https://twitter.com/hashtag/aspnet5?src=hash">#aspnet5</a> + <a href="https://twitter.com/hashtag/nginx?src=hash">#nginx</a> + <a href="https://twitter.com/hashtag/docker?src=hash">#docker</a> + <a href="https://twitter.com/hashtag/dockercompose?src=hash">#dockercompose</a> and simple load balancing :) <a href="https://t.co/YnJamDubIS">https://t.co/YnJamDubIS</a> <a href="https://t.co/pBOWDDnVHR">pic.twitter.com/pBOWDDnVHR</a></p>— Tugberk Ugurlu (@tourismgeek) <a href="https://twitter.com/tourismgeek/status/688319154084909056">January 16, 2016</a></blockquote><script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script> <blockquote> <p>It is not like we didn't have these options before in .NET web development world. To give you an example, you can perfectly run an ASP.NET Web API application under mono and expose it to outside world behind NGINX. However, ASP.NET 5 makes these options really straight forward to adopt.</p></blockquote> <p>The end result we will achieve here will have the below look and you can see the sample I have put together for this <a href="https://github.com/tugberkugurlu/aspnet-5-samples/tree/b5ff85b35cd4d24474a575f07a97bad98bd4e5d5/nginx-lb-sample">here</a>:</p> <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a0eac3a1-6842-4ed8-8694-d120e6ae8ffd.jpg"><img title="arch-diagram" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="arch-diagram" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/f51236ee-c58a-4b01-975d-1962f0286076.jpg" width="644" height="323"></a></p> <h3>ASP.NET 5 Application on RC1</h3> <p>For this sample, I have a very simple APS.NET 5 application which gives you an hello message and lists the environment variables available under the machine. The project structure looks like this:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>tugberk@ubuntu:~<span style="color: gray">/</span>apps<span style="color: gray">/</span>aspnet<span style="color: gray">-</span>5<span style="color: gray">-</span>samples<span style="color: gray">/</span>nginx<span style="color: gray">-</span>lb<span style="color: gray">-</span>sample$ tree
.
├── docker<span style="color: gray">-</span>compose.yml
├── docker<span style="color: gray">-</span>nginx.dockerfile
├── docker<span style="color: gray">-</span>webapp.dockerfile
├── global.json
├── nginx.conf
├── NuGet.Config
├── README.md
└── WebApp
    ├── hosting.json
    ├── project.json
    └── Startup.cs</pre></div></div>
<p>I am not going to put the application code here but you can find the entire code <a href="https://github.com/tugberkugurlu/aspnet-5-samples/tree/b5ff85b35cd4d24474a575f07a97bad98bd4e5d5/nginx-lb-sample/WebApp">here</a>. However, there is one important thing that I want to mention which is the server URL you will expose ASP.NET 5 application through Kestrel. To make Docker happy, we need to expose the application through "0.0.0.0" rather than localhost or 127.0.0.1. <a href="https://twitter.com/markrendle">Mark Rendle</a> has <a href="https://blog.rendle.io/asp-net-5-dnx-beta8-connection-refused-in-docker/">a great resource on this</a> explaining why and I have the following <a href="https://github.com/tugberkugurlu/aspnet-5-samples/blob/b5ff85b35cd4d24474a575f07a97bad98bd4e5d5/nginx-lb-sample/WebApp/hosting.json">hosting.json file</a> which also covers this issue:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>{
    <span style="color: #a31515">"server"</span>: <span style="color: #a31515">"Microsoft.AspNet.Server.Kestrel"</span>,
    <span style="color: #a31515">"server.urls"</span>: <span style="color: #a31515">"http://0.0.0.0:5090"</span>
}</pre></div></div>
<h3>Running ASP.NET 5 Application under Docker</h3>
<p>The next step is to run the ASP.NET 5 application under Docker. With the <a href="https://hub.docker.com/r/microsoft/aspnet/">ASP.NET Docker image</a> on <a href="https://hub.docker.com/">Docker Hub</a>, this is insanely simple. Again, Mark Rendle has three amazing posts on ASP.NET 5, Docker and Linux combination as <a href="https://blog.rendle.io/fun-with-asp-net-5-and-docker/">Part 1</a>, <a href="https://blog.rendle.io/fun-with-asp-net-5-linux-docker-part-2/">Part 2</a> and <a href="https://blog.rendle.io/fun-with-asp-net-5-linux-docker-part-3/">Part 3</a>. I strongly encourage you to check them out. For my sample here, I have the below <a href="https://docs.docker.com/engine/reference/builder/">Dockerfile</a> (<a href="https://github.com/tugberkugurlu/aspnet-5-samples/blob/b5ff85b35cd4d24474a575f07a97bad98bd4e5d5/nginx-lb-sample/docker-webapp.dockerfile">reference to the file</a>):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>FROM microsoft/aspnet:1.0.0-rc1-update1

COPY ./WebApp/project.json /app/WebApp/
COPY ./NuGet.Config /app/
COPY ./global.json /app/
WORKDIR /app/WebApp
RUN [<span style="color: #a31515">"dnu"</span>, <span style="color: #a31515">"restore"</span>]
ADD ./WebApp /app/WebApp/

EXPOSE 5090
ENTRYPOINT [<span style="color: #a31515">"dnx"</span>, <span style="color: #a31515">"run"</span>]</pre></div></div>
<p>That's all I need to be able to run my ASP.NET 5 application under Docker. What I can do now is to build the Docker image and run it:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>docker build <span style="color: gray">-</span>f docker<span style="color: gray">-</span>webapp.dockerfile <span style="color: gray">-</span>t hellowebapp .
docker run <span style="color: gray">-</span>d <span style="color: gray">-</span>p 5090:5090 hellowebapp</pre></div></div>
<p>The container now running in a detached mode and you should be able to hit the HTTP endpoint from your host:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/208324cc-5b7c-4f6a-a506-c10679a2e588.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d14bc016-8d8b-4f4f-a787-90afbec823cf.png" width="644" height="255"></a></p>
<p>From there, you do whatever you want to the container. Rebuild it, stop it, remove it, so and so forth.</p>
<h3>NGINX and Docker Compose</h3>
<p>Last pieces of the puzzle here are NGINX and Docker Compose. For those of who don't know what NGINX is: NGINX is a free, open-source, high-performance HTTP server and reverse proxy. Under production, you really don't want to expose Kestrel to outside world directly. Instead, you should put Kestrel behind a mature web server like NGINX, IIS or <a href="https://httpd.apache.org/">Apache Web Server</a>. </p>
<blockquote>
<p>There are two great videos you can watch on Kestrel and Linux hosting which gives you the reasons why you should put Kestrel behind a web server. I strongly encourage you to check them out before putting your application on production in Linux.</p>
<ul>
<li><a href="https://channel9.msdn.com/Events/ASPNET-Events/ASPNET-Fall-Sessions/ASPNET-5-Kestrel">ASP.NET 5: Kestrel</a> 
<li><a href="https://channel9.msdn.com/Events/ASPNET-Events/ASPNET-Fall-Sessions/ASPNET-5-Considerations-for-Production-Linux-Environments">ASP.NET 5: Considerations for Production Linux Environments</a></li></ul></blockquote>
<p>Docker Compose, on the other hand, is a completely different type of tool. It is a tool for defining and running multi-container Docker applications. With Compose, you use <a href="https://docs.docker.com/compose/compose-file/">a Compose file</a> (which is a YAML file) to configure your application’s services. This is a perfect fit for what we want to achieve here since we will have at least three containers running:</p>
<ul>
<li>ASP.NET 5 application 1 Container: An instance of the ASP.NET 5 application 
<li>ASP.NET 5 application 2 Container: Another instance of the ASP.NET 5 application 
<li>NGINX Container: An NGINX process which will proxy the requests to ASP.NET 5 applications.</li></ul>
<p>Let's start with configuring NGINX first and make it possible to run under Docker. This is going to very easy as NGINX also has <a href="https://hub.docker.com/_/nginx/">an image up on Docker Hub</a>. We will use this image and tell NGINX to read <a href="https://github.com/tugberkugurlu/aspnet-5-samples/blob/b5ff85b35cd4d24474a575f07a97bad98bd4e5d5/nginx-lb-sample/nginx.conf">our config</a> file which looks like this:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>worker_processes 4;

events { worker_connections 1024; }

http {
    upstream web-app {
        server webapp1:5090;
        server webapp2:5090;
    }

    server {
      listen 80;

      location / {
        proxy_pass http://web-app<span style="color: green">;</span>
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection keep-alive;
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
      }
    }
}</pre></div></div>
<p>This configuration file has some generic stuff in it but the most importantly, it has our load balancing and reverse proxy configuration. This configuration tells NGINX to accept requests on port 80 and proxy those requests to webapp1:5090 and webapp2:5090. Check out the <a href="https://www.nginx.com/resources/admin-guide/reverse-proxy/">NGINX reverse proxy guide</a> and <a href="https://www.nginx.com/resources/admin-guide/load-balancer/">load balancing guide</a> for more information about how you can customize the way you are doing the proxying and load balancing but the above configuration is enough for this sample.</p>
<blockquote>
<p>There is also an important part in this NGINX configuration to make Kestrel happy. <a href="https://github.com/aspnet/KestrelHttpServer/issues/341">Kestrel has an annoying bug in RC1</a> which <a href="https://github.com/aspnet/KestrelHttpServer/commit/e4fd91bb68f535801ca8a79aa453ea3fb3f448fe">has been already fixed for RC2</a>. To work around the issue, you need to set "Connection: keep-alive" header which is what we are doing with "proxy_set_header Connection keep-alive;" declaration in our NGINX configuration.</p></blockquote>
<p>Here is what NGINX Dockerfile looks like (<a href="https://github.com/tugberkugurlu/aspnet-5-samples/blob/b5ff85b35cd4d24474a575f07a97bad98bd4e5d5/nginx-lb-sample/docker-nginx.dockerfile">reference to the file</a>):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>FROM nginx
COPY ./nginx.conf /etc/nginx/nginx.conf</pre></div></div>
<p>You might wonder at this point about what webapp1 and webapp2 (which we have indicated inside the NGINX configuration file) map to. These are the DNS references for the containers which will run our ASP.NET 5 applications and when we link them in our Docker Compose file, the DNS mapping will happen automatically for container names. Finally, here is what out composition looks like inside the Docker Compose file (<a href="https://github.com/tugberkugurlu/aspnet-5-samples/blob/b5ff85b35cd4d24474a575f07a97bad98bd4e5d5/nginx-lb-sample/docker-compose.yml">reference to the file</a>):</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>webapp1:
  build: .
  dockerfile: docker-webapp.dockerfile
  container_name: hasample_webapp1
  ports:
    - <span style="color: #a31515">"5091:5090"</span>
    
webapp2:
  build: .
  dockerfile: docker-webapp.dockerfile
  container_name: hasample_webapp2
  ports:
    - <span style="color: #a31515">"5092:5090"</span>

nginx:
  build: .
  dockerfile: docker-nginx.dockerfile
  container_name: hasample_nginx
  ports:
    - <span style="color: #a31515">"5000:80"</span>
  links:
    - webapp1
    - webapp2</pre></div></div>
<p>You can see under the third container definition, we linked previously defined two containers to NGINX container. Alternatively, you may want to look at <a href="http://progrium.com/blog/2014/07/29/understanding-modern-service-discovery-with-docker/">Service Discovery in context of Docker</a> instead of linking.</p>
<p>Now we have everything in place and all we need to do is to run two docker-compose commands (under the directory where we have Docker Compose file) to get the application up and running:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>docker<span style="color: gray">-</span>compose build
docker<span style="color: gray">-</span>compose up</pre></div></div>
<p>After these, we should see three containers running. We should also be able to hit localhost:5000 from the host machine and see that the load is being distributed to both ASP.NET 5 application containers:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/de2fc6ca-9619-4a57-bf03-072832c84402.gif"><img title="compose2" style="display: inline" alt="compose2" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a178e8d6-767e-4c11-a101-979e22d24717.gif" width="640" height="360"></a></p>
<p>Pretty great! However, this is just sample for demo purposes to show how simple it is to have an environment like this up an running locally. This probably provides no performance gains when you run all containers in one box. My next step is going to be to get <a href="http://www.haproxy.org/">HAProxy</a> in this mix and let it do the load balancing instead.</p>  