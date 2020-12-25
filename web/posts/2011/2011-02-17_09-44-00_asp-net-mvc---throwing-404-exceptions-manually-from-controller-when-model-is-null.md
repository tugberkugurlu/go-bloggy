---
id: 49a12e87-d0ad-4374-a00e-1908317be65f
title: 'ASP.NET MVC : Throwing 404 Exceptions Manually From Controller When Model
  is Null'
abstract: This post is a quick demonstration of how you can throw HttpException of
  404 manually from a controller on ASP.NET MVC when the model you're passing is null
created_at: 2011-02-17 09:44:00 +0000 UTC
tags:
- .net
- ASP.Net
- ASP.NET MVC
- C#
- Tips
slugs:
- asp-net-mvc---throwing-404-exceptions-manually-from-controller-when-model-is-null
---

<p>How do you handle your actions inside a controller when the model you are passing to the view is null? Sometimes it is best to show <em>&lsquo;There is no product as you requested&rsquo;</em> kind of message but sometime it is getting dull. Especially on <strong>CRUD </strong>based actions.</p>
<p><strong>HttpException</strong> is becoming very handy here. We could throw this exception from our controllers and the application will render the status exception page we are defining. In this post I would like to show how we handle this properly. Here is our scenario : we have an MVC app and in one of our edit page, we would like to throw 404 exception if the model we are passing is null.</p>
<p>&nbsp;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [7]">        [Authorize]
        public ActionResult Edit(Guid id) {

            var model = Repository.GetSingle(id);

            if (model == null)
                throw new HttpException(404, "not found");

            return View(model);
        }</pre>
<p>&nbsp;</p>
<p>This code is basically telling that try to get the item whose ID is <strong><em>id </em></strong>and check if it is null. If it is null, throw HttpException whose status code is 404. If not, return the View along with passing <strong><em>model </em></strong>into view page.</p>
<p>When we hit a page which has null model, we will see the 404 page as expected;</p>
<p><a href="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/097842e5fed6_BBB6/image.png" target="_blank"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px initial initial;" title="image" alt="image" src="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/097842e5fed6_BBB6/image_thumb.png" width="644" height="195" /></a></p>
<p>This is that easy ! Give it a try, it won&rsquo;t hurt <img style="border-style: none;" alt="Smile" src="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/097842e5fed6_BBB6/wlEmoticon-smile.png" /></p>