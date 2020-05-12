---
id: 08280ed4-866b-4087-be4a-fbaf4f96587e
title: How to Delete a Previously Created Cookie With C# ASP.Net / Deleting Cookie
  ASP.Net
abstract: You created a cookie on you asp.net forms application now you would like
  to delete it. This quick article show how to do the trick...
created_at: 2010-10-03 08:13:12 +0000 UTC
tags:
- .net
- ASP.Net
- C#
slugs:
- how-to-delete-a-previously-created-cookie-with-c-sharp-asp-net-deleting-cookie-asp-net
---

<p>I am sure that you leave cookies into client computer on purpose or&nbsp;without knowing if you're an asp.net developer. Also if you create a sales website, there must be steps on this sales applicatiion and you should get the data from a previous page to another pages. The best way is always cookie if the information is not so sensetive.&nbsp;<br /> <br /> But what if we wanna delete them ? It is not as simple as creating them (essentially it is but not by directly deleting it) We cannot directly delete a file from user's computer. So what will we do?&nbsp;<br /> <br /> I assume there are another ways but the one way is setting the expiry date on a previous date. Here is the example;</p>
<p>&nbsp;</p>
<pre class="brush: c-sharp">if (Request.Cookies["UserSettings"] != null) {

    HttpCookie myCookie = new HttpCookie("UserSettings");
    //Here, we are setting the time to a previous time.
    //When the browser detect it next time, it will be deleted automatically.
    myCookie.Expires = DateTime.Now.AddDays(-1d);
    Response.Cookies.Add(myCookie);

} </pre>
<p>&nbsp;</p>
<div></div>
<div>Hope this helps :)</div>