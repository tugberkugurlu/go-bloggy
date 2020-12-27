---
id: 8a0611a5-6f37-48da-9ca5-a2f52a70629e
title: Mapping ASP.NET SignalR Connections to Real Application Users
abstract: One of the common questions about SignalR is how to broadcast a message
  to specific users and the mapping the SignalR connections to your real application
  users is the key component for this.
created_at: 2013-01-01 15:32:00 +0000 UTC
tags:
- .NET
- ASP.Net
- SignalR
slugs:
- mapping-asp-net-signalr-connections-to-real-application-users
---

<p><a href="http://signalr.net">SignalR</a>; the incredible real-time web framework for .NET. You all probably heard of it, maybe played with it and certainly loved it. If you haven&rsquo;t, why don&rsquo;t you start by reading the <a href="https://github.com/SignalR/SignalR/wiki">SignalR docs</a> and <a href="https://twitter.com/davidfowl">@davidfowl</a>&rsquo;s blog post on <a href="http://weblogs.asp.net/davidfowler/archive/2012/11/11/microsoft-asp-net-signalr.aspx">Microsoft ASP.NET SignalR</a>? Yes, you heard me right: <a href="http://www.asp.net/signalr">it&rsquo;s now officially a Microsoft product</a>, too (and you may see as a bad or good thing).</p>
<p>One of the common questions about SignalR is how to broadcast a message to specific users and the answer depends on what you are really trying to do. If you are working with <a href="https://github.com/SignalR/SignalR/wiki/Hubs">Hubs</a>, there is this notion of Groups which you can add connections to. Then, you can send messages to particular groups. It&rsquo;s also very straight forward to work with Groups with the latest SignalR server API:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> MyHub : Hub
{
    <span style="color: blue;">public</span> Task Join()
    {
        <span style="color: blue;">return</span> Groups.Add(Context.ConnectionId, <span style="color: #a31515;">"foo"</span>);
    }

    <span style="color: blue;">public</span> Task Send(<span style="color: blue;">string</span> message)
    {
        <span style="color: blue;">return</span> Clients.Group(<span style="color: #a31515;">"foo"</span>).addMessage(message);
    }
}</pre>
</div>
</div>
<p>You also have a chance to exclude some connections within a group for that particular message. However, if you have more specific needs such as broadcasting a message to user x, Groups are not your best bet. There are couple of reasons:</p>
<ul>
<li>SignalR is not aware of any of your business logic. SignalR knows about currently connected connections and their connection ids. That&rsquo;s all.</li>
<li>Assume that you have some sort of authentication on your application (forms authentication, etc.). In this case, your user can have multiple connection ids by consuming the application with multiple ways. So, you cannot just assume that your user has only one connection id.</li>
</ul>
<p>By considering these factors, mapping the SignalR connections to actual application users is the best way to solve this particular problem here. To demonstrate how we can actually solve this problem with code, I&rsquo;ve put together a simple chat application and <a href="https://github.com/tugberkugurlu/SignalRSamples/tree/master/ConnectionMappingSample">the source code is also available on GitHub</a>.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/56d41d8c6575_E917/image.png"><img height="367" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/56d41d8c6575_E917/image_thumb.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>This application obviously not a production ready application and the purpose here is to show how to achieve connection mapping. I didn&rsquo;t even use a persistent data storage technology for this demo.</p>
<p>The scenarios I needed to the above sample are below ones:</p>
<ul>
<li>A user can log in with his/her username and can access to the chat page. A user also can sign out whenever they want.</li>
<li>A user can see other connected users on the right hand side at the screen.</li>
<li>A user can send messages to all connected users.</li>
<li>A user can send private messages to a particular user.</li>
</ul>
<p>To achieve the first goal, I have a very simple <a href="http://www.asp.net/mvc">ASP.NET MVC</a> controller:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> AccountController : Controller {

    <span style="color: blue;">public</span> ViewResult Login() {

        <span style="color: blue;">return</span> View();
    }

    [HttpPost]
    [ActionName(<span style="color: #a31515;">"Login"</span>)]
    <span style="color: blue;">public</span> ActionResult PostLogin(LoginModel loginModel) {

        <span style="color: blue;">if</span> (ModelState.IsValid) {

            FormsAuthentication.SetAuthCookie(loginModel.Name, <span style="color: blue;">true</span>);
            <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"index"</span>, <span style="color: #a31515;">"home"</span>);
        }

        <span style="color: blue;">return</span> View(loginModel);
    }

    [HttpPost]
    [ActionName(<span style="color: #a31515;">"SignOut"</span>)]
    <span style="color: blue;">public</span> ActionResult PostSignOut() {

        FormsAuthentication.SignOut();
        <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"index"</span>, <span style="color: #a31515;">"home"</span>);
    }
}</pre>
</div>
</div>
<p>When you hit the home page as an unauthenticated user, you will get redirected to login page to log yourself in. As you can see from the PostLogin action method, everybody can authenticate themselves by simply entering their name which is obviously not what you would want in a real world application.</p>
<p>As I am hosting my SignalR application under the same process with my ASP.NET MVC application, the authenticated users will flow through the SignalR pipeline, too. So, I protected my Hub and its methods with the <a href="https://github.com/SignalR/SignalR/blob/master/src/Microsoft.AspNet.SignalR.Core/Hubs/Pipeline/Auth/AuthorizeAttribute.cs">Microsoft.AspNet.SignalR.Hubs.AuthorizeAttribute</a>.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>[Authorize]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> ChatHub : Hub { 

    <span style="color: green;">//...</span>
}</pre>
</div>
</div>
<p>As we are done with the authorization and authentication pieces, we can now move on and implement our Hub. What I want to do first is to keep track of connected users with a static dictionary. Now, keep in mind again here that you would not want to use a static dictionary on a real world application, especially when you have a web farm scenario. You would want to keep track of the connected users with a persistent storage system such as MongoDB, RavenDB, SQL Server, etc. However, for our demo purposes, a static dictionary will just work fine.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> User {

    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> HashSet&lt;<span style="color: blue;">string</span>&gt; ConnectionIds { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

[Authorize]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> ChatHub : Hub {

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">readonly</span> ConcurrentDictionary&lt;<span style="color: blue;">string</span>, User&gt; Users 
        = <span style="color: blue;">new</span> ConcurrentDictionary&lt;<span style="color: blue;">string</span>, User&gt;();
        
    <span style="color: green;">// ...</span>
}</pre>
</div>
</div>
<p>Each user will have a name and associated connection ids. Now the question is how to add and remove values to this dictionary. SignalR raises three particular events on your hub: OnConnected, OnDisconnected, OnReconnected and the purposes of these events are very obvious.</p>
<p>During OnConnected event, we need to add the current connection id to the user&rsquo;s connection id collection (we need to create the User object first if it doesn&rsquo;t exist inside the dictionary). We also want to broadcast this information to all clients so that they can update their connected users list. Here is how I implemented the OnConnected method:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>[Authorize]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> ChatHub : Hub {

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">readonly</span> ConcurrentDictionary&lt;<span style="color: blue;">string</span>, User&gt; Users 
        = <span style="color: blue;">new</span> ConcurrentDictionary&lt;<span style="color: blue;">string</span>, User&gt;();
        
    <span style="color: blue;">public</span> <span style="color: blue;">override</span> Task OnConnected() {

        <span style="color: blue;">string</span> userName = Context.User.Identity.Name;
        <span style="color: blue;">string</span> connectionId = Context.ConnectionId;

        <span style="color: blue;">var</span> user = Users.GetOrAdd(userName, _ =&gt; <span style="color: blue;">new</span> User {
            Name = userName,
            ConnectionIds = <span style="color: blue;">new</span> HashSet&lt;<span style="color: blue;">string</span>&gt;()
        });

        <span style="color: blue;">lock</span> (user.ConnectionIds) {

            user.ConnectionIds.Add(connectionId);
            
            <span style="color: green;">// TODO: Broadcast the connected user</span>
        }

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.OnConnected();
    }
}</pre>
</div>
</div>
<p>First of all, we have gathered the currently authenticated user name and connected user&rsquo;s connection id. Then, we look inside the dictionary to get the user based on the user name. If it doesn&rsquo;t exist inside the dictionary, we create one and set it to the local variable named user. Lastly, we add the connection id and updated the dictionary.</p>
<p>Notice that we have a TODO comment at the end telling that we need to broadcast the connected user&rsquo;s name. Obviously, we don&rsquo;t want to broadcast this information to the caller itself. However, we still have two options here and which one you would choose may depend on your case. As the user might have multiple connections, broadcasting this information over Clients.Others API is not a way to follow. Instead, we can use Clients.AllExcept method which takes a list of connection ids as parameter to exclude. So, we can pass the all connection ids of the user and we are good to go.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">override</span> Task OnConnected() {

    <span style="color: green;">// Lines omitted for brevity</span>
    
    Clients.AllExcept(user.ConnectionIds.ToArray()).userConnected(userName);

    <span style="color: blue;">return</span> <span style="color: blue;">base</span>.OnConnected();
}</pre>
</div>
</div>
<p>This is a fine approach if we want to broadcast each connection of the user to every client other than the user itself. However, we may only want to broadcast the first connection. Doing so is very straight forward, too. We just need to inspect the count of the connection ids and if it equals to one, we can broadcast the information. This approach is the one that I ended up taking for this demo.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">override</span> Task OnConnected() {

    <span style="color: green;">// Lines omitted for brevity</span>

    <span style="color: blue;">lock</span> (user.ConnectionIds) {

        <span style="color: green;">// Lines omitted for brevity</span>
        
        <span style="color: blue;">if</span> (user.ConnectionIds.Count == 1) {

            Clients.Others.userConnected(userName);
        }
    }

    <span style="color: blue;">return</span> <span style="color: blue;">base</span>.OnConnected();
}</pre>
</div>
</div>
<p>When the disconnect event is fired, OnDisconnected method will be called and we need to remove the current connection id from the users dictionary. Similar to what we have done inside the OnConnected method, we need to handle the fact that user can have multiple connections and if there is no connection left, we want to remove the user from Users dictionary completely. As we did when a user connection arrives, we need to broadcast the disconnected users, too and we have the same two options here as well. I added both to the below code and commented out the one that we don&rsquo;t need for our demo.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>[Authorize]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> ChatHub : Hub {

    <span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">readonly</span> ConcurrentDictionary&lt;<span style="color: blue;">string</span>, User&gt; Users 
        = <span style="color: blue;">new</span> ConcurrentDictionary&lt;<span style="color: blue;">string</span>, User&gt;();

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> Task OnDisconnected() {

        <span style="color: blue;">string</span> userName = Context.User.Identity.Name;
        <span style="color: blue;">string</span> connectionId = Context.ConnectionId;
        
        User user;
        Users.TryGetValue(userName, <span style="color: blue;">out</span> user);
        
        <span style="color: blue;">if</span> (user != <span style="color: blue;">null</span>) {

            <span style="color: blue;">lock</span> (user.ConnectionIds) {

                user.ConnectionIds.RemoveWhere(cid =&gt; cid.Equals(connectionId));

                <span style="color: blue;">if</span> (!user.ConnectionIds.Any()) {

                    User removedUser;
                    Users.TryRemove(userName, <span style="color: blue;">out</span> removedUser);

                    <span style="color: green;">// You might want to only broadcast this info if this </span>
                    <span style="color: green;">// is the last connection of the user and the user actual is </span>
                    <span style="color: green;">// now disconnected from all connections.</span>
                    Clients.Others.userDisconnected(userName);
                }
            }
        }

        <span style="color: blue;">return</span> <span style="color: blue;">base</span>.OnDisconnected();
    }
}</pre>
</div>
</div>
<p>When the OnReconnected method is invoked, we don&rsquo;t need to perform any special logic here as the connection id will be the same. With these implementations, we are now keeping track of the connected users and we have mapped the connections to real application users.</p>
<p>Going back to our scenarios list above, we have the 4th requirement: a user sending private messages to a particular user. This is where we actually need the connection mapping functionality. As an high level explanation, the client will send the name of the user that s/he wants to send the message to privately. So, server needs to make sure that it is only sending the message to the designated user. I am not going to go through all the client code (as you can check them out from the source code and they are not that much related to the topic here) but the piece of JavaScript code that actually decides whether to send a public or private message is as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre>$sendBtn.click(<span style="color: blue;">function</span> (e) {

    <span style="color: blue;">var</span> msgValue = $msgTxt.val();
    <span style="color: blue;">if</span> (msgValue !== <span style="color: blue;">null</span> &amp;&amp; msgValue.length &gt; 0) {

        <span style="color: blue;">if</span> (viewModel.isInPrivateChat()) {

            chatHub.server.send(msgValue, viewModel.privateChatUser()).fail(<span style="color: blue;">function</span> (err) {
                console.log(<span style="color: #a31515;">'Send method failed: '</span> + err);
            });
        }
        <span style="color: blue;">else</span> {
            chatHub.server.send(msgValue).fail(<span style="color: blue;">function</span> (err) {
                console.log(<span style="color: #a31515;">'Send method failed: '</span> + err);
            });
        }
    }
    e.preventDefault();
});</pre>
</div>
</div>
<p>The above code inspects the KnockoutJS view model to see if the sender is at the private chat mode. If s/he is, it invokes the send hub method on the sever with two parameters which means that this will be a private message. If the sender is not at the private chat mode, we will just invoke the send hub method by passing only one parameter for the message. Let&rsquo;s first look at Send Hub method that takes one parameter:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">void</span> Send(<span style="color: blue;">string</span> message) {

    <span style="color: blue;">string</span> sender = Context.User.Identity.Name;

    Clients.All.received(<span style="color: blue;">new</span> { 
        sender = sender, 
        message = message, 
        isPrivate = <span style="color: blue;">false</span> 
    });
}</pre>
</div>
</div>
<p>Inside the send method above, we first retrieved the sender's name through the authenticated user principal. Then, we are broadcasting the message to all clients with a few more information such as the sender name and the privacy state of the message. Let&rsquo;s now look at the second Send method inside the Hub whose job is to send private messages:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">void</span> Send(<span style="color: blue;">string</span> message, <span style="color: blue;">string</span> to) {

    User receiver;
    <span style="color: blue;">if</span> (Users.TryGetValue(to, <span style="color: blue;">out</span> receiver)) {

        User sender = GetUser(Context.User.Identity.Name);

        IEnumerable&lt;<span style="color: blue;">string</span>&gt; allReceivers;
        <span style="color: blue;">lock</span> (receiver.ConnectionIds) {
            <span style="color: blue;">lock</span> (sender.ConnectionIds) {

                allReceivers = receiver.ConnectionIds.Concat(
                    sender.ConnectionIds);
            }
        }

        <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> cid <span style="color: blue;">in</span> allReceivers) {
        
            Clients.Client(cid).received(<span style="color: blue;">new</span> { 
                sender = sender.Name, 
                message = message, 
                isPrivate = <span style="color: blue;">true</span> 
            });
        }
    }
}

<span style="color: blue;">private</span> User GetUser(<span style="color: blue;">string</span> username) {

    User user;
    Users.TryGetValue(username, <span style="color: blue;">out</span> user);

    <span style="color: blue;">return</span> user;
}</pre>
</div>
</div>
<p>Here, we are first trying to get the receiver based on the to parameter that we have received. If we find one, we are also retrieving the sender based on the his/her name. Now, we have the sender and the receiver in our hands. What we want is to broadcast this message to the receiver and the sender. So, we are putting the sender&rsquo;s and the receiver&rsquo;s connection ids together first. Finally, we are looping through that connection ids list to send the message to each connection by using the Clients.Client method which takes the connection id as a parameter.</p>
<p>When we try this out, we should see it working as below:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/56d41d8c6575_E917/image_3.png"><img height="335" width="644" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/56d41d8c6575_E917/image_thumb_3.png" alt="image" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="image" /></a></p>
<p>Grab the solution and try it yourself, too. I hope this post helped you to solve your problem <img src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/56d41d8c6575_E917/wlEmoticon-smile.png" alt="Smile" style="border-style: none;" class="wlEmoticon wlEmoticon-smile" /></p>
<h3>References</h3>
<ul>
<li><a href="https://github.com/SignalR/SignalR/wiki/Hubs">SignalR Hubs</a></li>
<li><a href="https://github.com/SignalR/Samples/blob/master/BasicChat/ChatWithTracking.cs">ChatWithTracking.cs Hub Sample</a></li>
<li><a href="https://github.com/SignalR/SignalR/wiki/Talks">SignalR Videos to watch</a></li>
</ul>