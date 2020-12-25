---
id: 0d3b8414-09b3-4017-9550-617c2fd383d4
title: ASP.NET MVC Remote Validation For Multiple Fields With AdditionalFields Property
abstract: This post shows the implementation of ASP.NET MVC Remote Validation for
  multiple fields with AdditionalFields property and we will validate the uniqueness
  of a product name under a chosen category.
created_at: 2011-10-24 06:44:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- C#
- JQuery
slugs:
- asp-net-mvc-remote-validation-for-multiple-fields-with-additionalfields-property
---

<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/why-wont-you-validate-me-hrrrr.png"><img style="background-image: none; margin: 0px 15px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="why-wont-you-validate-me-hrrrr" border="0" alt="why-wont-you-validate-me-hrrrr" align="left" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/why-wont-you-validate-me-hrrrr_thumb.png" width="244" height="200" /></a></p>
<p>I blogged about <a title="https://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation" href="https://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation" target="_blank">ASP.NET MVC Remote Validation</a> couple of days ago by showing a simple scenario for how to use it. That blog post simply covers everything about remote validation except for properties of <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.remoteattribute(v=vs.98).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.remoteattribute(v=vs.98).aspx" target="_blank">RemoteAttribute Class</a>.</p>
<p>In that blog post, we covered validating user inputs uniqueness for one field in order to alert the user before the form is posted. This is a quite nice built-in feature which takes basically not more 5 minutes to implement.</p>
<p>But we as in all human beings tend not to be happy with the current situation and want more. It is in our nature, we can&rsquo;t help that. So, as soon as I use this feature I ask myself <em><strong>"What if I would like to check the uniqueness of product name inside the chosen category?"</strong></em>.<em><strong> </strong></em>Well, it turns out that <a title="http://asp.net/mvc" href="http://asp.net/mvc" target="_blank">ASP.NET MVC</a> team asked this question themselves way earlier than me and added a support for this kind of situations.</p>
<p><a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.remoteattribute.additionalfields(v=vs.98).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.remoteattribute.additionalfields(v=vs.98).aspx" target="_blank">RemoteAttribute.AdditionalFields Property</a> gets or sets the additional fields that are required for validation. <strong>AdditionalFields</strong> property is <strong>string</strong> property and can be provided for multiple fields.</p>
<p>Probably another question pops up on your mind : <em>"How does ASP.NET MVC Framework relate the fields, which we will provide for <strong>AdditionalFields</strong> property, with <strong>form elements</strong> and pass them our <strong>controller</strong>?" </em>If we need to ask this question in plain English, that would be this : <em>"How the heck does it know in which category I trying to check the uniqueness of the product name?"</em>. Here, ASP.NET MVC Model Binding feature plays a huge role here as much as <a title="http://jquery.com/" href="http://jquery.com/" target="_blank">JQuery</a>.</p>
<p><strong>Simple Scenario</strong></p>
<p>Assume that we need to store products under a category inside our database and we need every each product name to be unique under the chosen category.</p>
<blockquote>
<p>In order to stick with the topic here, I have already built the necessary stuff (database, repos, etc.). I used SQL Server Compact Version 4 here. As sample database, you won&rsquo;t believe what I used : <strong>Northwind </strong><img style="border-style: none;" class="wlEmoticon wlEmoticon-winkingsmile" alt="Winking smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/wlEmoticon-winkingsmile.png" /> I also used Entity Framework to reach out to the data. You can see all of the code from my <a title="https://github.com" href="https://github.com" target="_blank">GitHub</a> Repo : <a href="https://github.com/tugberkugurlu/MvcRemoteValidationSample">https://github.com/tugberkugurlu/MvcRemoteValidationSample</a></p>
</blockquote>
<p>Our create form looks like below :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb.png" width="644" height="383" /></a></p>
<p>The product model class looks like this :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">partial</span> <span style="color: blue;">class</span> Product
{
    <span style="color: blue;">public</span> Product()
    {
        <span style="color: blue;">this</span>.Order_Details = <span style="color: blue;">new</span> HashSet&lt;Order_Detail&gt;();
    }

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Product_ID { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> Nullable&lt;<span style="color: blue;">int</span>&gt; Supplier_ID { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> Nullable&lt;<span style="color: blue;">int</span>&gt; Category_ID { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Product_Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> English_Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Quantity_Per_Unit { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> Nullable&lt;<span style="color: blue;">decimal</span>&gt; Unit_Price { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> Nullable&lt;<span style="color: blue;">short</span>&gt; Units_In_Stock { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> Nullable&lt;<span style="color: blue;">short</span>&gt; Units_On_Order { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> Nullable&lt;<span style="color: blue;">short</span>&gt; Reorder_Level { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">bool</span> Discontinued { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> Category Category { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> ICollection&lt;Order_Detail&gt; Order_Details { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">virtual</span> Supplier Supplier { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>Here is my additional partial class under the same namespace of my <strong>Product</strong> class :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[MetadataType(<span style="color: blue;">typeof</span>(Product.MetaData))]
<span style="color: blue;">public</span> <span style="color: blue;">partial</span> <span style="color: blue;">class</span> Product {

    <span style="color: blue;">private</span> <span style="color: blue;">class</span> MetaData {

        [Remote(
            <span style="color: #a31515;">"doesProductNameExistUnderCategory"</span>, 
            <span style="color: #a31515;">"Northwind"</span>, 
            AdditionalFields = <span style="color: #a31515;">"Category_ID"</span>,
            ErrorMessage = <span style="color: #a31515;">"Product name already exists under the chosen category. Please enter a different product name."</span>,
            HttpMethod = <span style="color: #a31515;">"POST"</span>
        )]
        [Required]
        <span style="color: blue;">public</span> <span style="color: blue;">string</span> Product_Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
        
        [Required]
        <span style="color: blue;">public</span> <span style="color: blue;">int</span>? Supplier_ID { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

        [Required]
        <span style="color: blue;">public</span> <span style="color: blue;">int</span>? Category_ID { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    }

}</pre>
</div>
</div>
<p>As you see there, we have decorated <strong>Product_Name</strong> property with <a title="http://msdn.microsoft.com/en-us/library/system.web.mvc.remoteattribute(v=vs.98).aspx" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.remoteattribute(v=vs.98).aspx" target="_blank">RemoteAttribute</a>. What we need to point out here is AdditionalFields property of RemoteAttribute. Simply we are saying here that : '<em>Pass Category_ID value to the controller action when you try validate the Product_Name</em>'. So, <strong>Ajax Request </strong>will send two form inputs to our action : <em>Product_Name</em> and <em>Category_ID</em>.<em> </em>As we have already seen on my previous blog post on <a title="https://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation" href="https://www.tugberkugurlu.com/archive/check-instantly-if-username-exists-asp-net-mvc-remote-validation" target="_blank">ASP.NET MVC Remote Validation</a>, Model Binder will bind those inputs to our action parameters in order to allow us to easily pick them.</p>
<p>Our main goal here is to check the uniqueness of the product name under the chosen category and that&rsquo;s why we need <em>Category_ID </em>value inside our controller action. Check controller action code below and you will get what I am trying to say :</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> JsonResult doesProductNameExistUnderCategory(<span style="color: blue;">int</span>? Category_ID, <span style="color: blue;">string</span> Product_Name) {

    <span style="color: blue;">var</span> model = _entities.Products.Where(x =&gt; (Category_ID.HasValue) ? 
            (x.Category_ID == Category_ID &amp;&amp; x.Product_Name == Product_Name) : 
            (x.Product_Name == Product_Name)
        );

    <span style="color: blue;">return</span> Json(model.Count() == 0);

}</pre>
</div>
</div>
<p>Inside this action, we are checking if the user name already exists under the chosen category or not. For the real purpose of this blog post, we only need to focus what we are returning here. The other part of this code mostly depends on your application&rsquo;s business logic.</p>
<p>We are basically done here. I will run the app and see it in action but please firstly see the below screenshot :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb_3.png" width="644" height="383" /></a></p>
<p>This table show the list of products inside our database. On the first row, we see that there is a product named <strong>Chai</strong> under Category_ID&nbsp;<strong>1</strong> whose Category_Name is <strong>Beverages</strong> (you don&rsquo;t see the Category_Name here but don&rsquo;t worry, just trust me). We will demonstrate our sample with this product values.</p>
<p>First, select the category from the select list :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb_4.png" width="644" height="383" /></a></p>
<p>Then, type Chai for <strong>Product_Name </strong>:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_5.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb_5.png" width="644" height="383" /></a></p>
<p>Then, simply hit Tab key :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_6.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb_6.png" width="644" height="383" /></a></p>
<p>All done but I would like to go a little deeper here. Let&rsquo;s put a break point on <strong>doesProductNameExistUnderCategory </strong>method and start a debug session.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb_7.png" width="644" height="379" /></a></p>
<p>When we follow the above steps again, we should end up like this :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_8.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb_8.png" width="644" height="379" /></a></p>
<p>Let&rsquo;s look what we got here :</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_9.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/image_thumb_9.png" width="644" height="379" /></a></p>
<p>We have all the necessary values to check if the value is legitimate or not.</p>
<p>Again, you can see all of the code from my GitHub Repo : <a href="https://github.com/tugberkugurlu/MvcRemoteValidationSample">https://github.com/tugberkugurlu/MvcRemoteValidationSample</a></p>
<p>Hope you enjoy the post <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/560581a7261c_92FF/wlEmoticon-smile.png" /></p>