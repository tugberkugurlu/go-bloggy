---
id: 62137996-bb8c-4596-ae8a-e8302075462a
title: Securing MongoDB Access with Username and Password
abstract: 'My MongoDb journey continues :) and I had my first attempt to put a username
  and password protection against a MongoDB instance. It went OK besides some hiccups
  along the way :) Let''s see what I did. '
created_at: 2014-05-01 10:44:00 +0000 UTC
tags:
- MongoDB
slugs:
- securing-mongodb-access-with-username-and-password
---

<p>My <a href="https://www.tugberkugurlu.com/archive/a-c-sharp-developers-first-thoughts-on-mongodb">MongoDb journey</a> <a href="https://www.tugberkugurlu.com/tags/mongodb">continues</a> :) and I had my first attempt to put a username and password protection against a <a href="http://mongodb.org">MongoDB</a> instance. It went OK besides some hiccups along the way :) Let's see what I did.  <p>First, I downloaded the latest (v2.6.0) MongoDB binaries as zip file and unzipped them. I put all MongoDB related stuff inside the c:\mongo directory for my development environment on windows and the structure of my c:\mongo directory is a little different: <p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/089be820-d661-4be0-9d9f-5828dc5790a4.png"><img title="Screenshot 2014-04-30 13.50.16" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="Screenshot 2014-04-30 13.50.16" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/24801bbb-f07b-46c3-bd06-e0ed289aee9d.png" width="644" height="477"></a>  <p>In order to set up the username and password authentication, first I need to get the <a href="http://docs.mongodb.org/manual/reference/program/mongod/">mongod</a> instance up and running with the <a href="http://docs.mongodb.org/manual/reference/configuration-options/#security.authorization">authorization on</a>. I achieved that by configuring it with the <a href="http://docs.mongodb.org/manual/reference/configuration-options/">config file</a>. Here is how it looks like:  <p><pre>dbpath = c:\mongo\data\db
port = 27017
logpath = c:\mongo\data\logs\mongo.log
auth = true</pre>
<p>With this config file in place, I can get the mongod instance up:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/44c7bf75-9666-4c07-8a1e-b57467a2c8d9.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/eaedc9f4-d0e0-4e86-b6d6-de728a14e1a1.png" width="644" height="181"></a></p>
<p>First, I need to connect to this mongod instance and create the admin user. As you can see inside my config file, the server requires authentication. However, there is a <a href="http://docs.mongodb.org/manual/core/authentication/#localhost-exception">localhost exception</a> if there is no user defined inside the system. So, I can connect to my instance anonymously (as I’m running on port 27017 on localhost, I don’t need to define anything while firing up the <a href="http://docs.mongodb.org/v2.2/mongo/">mongo shell</a>): </p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/70aa951f-d55f-4455-b3f9-fc809f727aa2.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/ba4f9ee3-5733-43e6-9646-4d5431dd7732.png" width="644" height="169"></a></p>
<p>All great! Let’s create the system user administrator. As everything else, <a href="http://docs.mongodb.org/manual/tutorial/enable-authentication/#create-the-system-user-administrator">this chore is nicely documented</a>, too:</p><pre>use admin
db.createUser(
  {
    user: "tugberk",
    pwd: "12345678",
    roles:
    [
      {
        role: "userAdminAnyDatabase",
        db: "admin"
      }
    ]
  }
)</pre>
<p>We are pretty much done. We have a user to administer our server now. Let’s disconnect from the mongo shell and reconnect to our mongod instance with our credentials:</p><pre>mongo --host localhost --port 27017 -u tugberk -p 12345678 --authenticationDatabase admin</pre>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/158e3df8-3964-478a-ac6e-320461f5db95.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d5743353-524c-42a6-ae44-a64a6b65c8fe.png" width="539" height="484"></a></p>
<p>I’m all there and I can see what my privileges at the server with this user are.</p>
<blockquote>
<p>If you try to connect to this MonogDB server anonymously, <a href="http://stackoverflow.com/questions/23387689/mongodb-server-can-still-be-accessed-without-credentials">you will see that you are still able to connect to it</a>. This’s really bad, isn’t it? Not at the level that you think it would be at. The real story is that MongoDB still allows you to connect to, but you won’t be able to do anything as the anonymous access is fully disabled. </p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/fe13b848-07fb-41e2-a156-5aaeb3a0e9ab.png"><img title="image" style="border-left-width: 0px; border-right-width: 0px; background-image: none; border-bottom-width: 0px; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border-top-width: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/23c83efd-e9cd-443f-9e7e-616db419f1ce.png" width="644" height="294"></a></p>
<p>The bad thing here is that your server existence is exposed which is still an important issue. Just be aware of this fact before getting started.</p></blockquote>
<p>The user we created still has restricted access to MonogDB server. If you want to <a href="http://docs.mongodb.org/manual/tutorial/add-admin-user/">have a user with unrestricted access</a>, you can create a user with <a href="http://docs.mongodb.org/manual/reference/built-in-roles/#superuser-roles">root role</a> assigned. In our case here, I will assign myself the root role:</p><pre>use admin
db.grantRolesToUser("tugberk", ["root"])</pre>
<h3>Resources</h3>
<ul>
<li><a href="http://docs.mongodb.org/manual/core/authorization/">MongoDB Authorization</a></li>
<li><a href="http://docs.mongodb.org/manual/administration/security-user-role-management/">User and Role Management Tutorials</a>
<li><a href="http://docs.mongodb.org/manual/administration/security-access-control/">Access Control Tutorials</a>
<li><a href="http://docs.mongodb.org/manual/tutorial/add-user-to-database/">Add a User to a Database</a></li></ul>  