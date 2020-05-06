---
title: 'ASP.NET MVC Code Review #2 - A Way of Working with Html Select Element (AKA
  DropDownList) In ASP.NET MVC'
abstract: 'This is the #2 of the series of blog posts which is about some core scenarios
  on ASP.NET MVC: A Way of Working with Html Select Element (AKA DropDownList) In
  ASP.NET MVC'
created_at: 2011-12-30 14:28:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
- C#
- Code Review
- Razor
slugs:
- asp-net-mvc-code-review-sharp2-a-way-of-working-with-html-select-element-aka-dropdownlist-in-asp-net-mvc
---

<p>This is the #2 of the series of blog posts which is about some core scenarios on <a href="http://asp.net/mvc">ASP.NET MVC</a>. In this one, code review #2, I will try to show you the best way of working with <strong>Html Select element </strong>(AKA <strong>DropDownList</strong>) in ASP.NET MVC. Here is the code:</p>
<p><strong>Repository Class which generates the dummy data for the demo:</strong></p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> ProductCategory {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> CategoryId { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> CategoryName { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> ProductCategoryRepo {

    <span style="color: blue;">public</span> List&lt;ProductCategory&gt; GetAll() {

        List&lt;ProductCategory&gt; categories = <span style="color: blue;">new</span> List&lt;ProductCategory&gt;();

        <span style="color: blue;">for</span> (<span style="color: blue;">int</span> i = 1; i &lt;= 10; i++) {

            categories.Add(<span style="color: blue;">new</span> ProductCategory { 
                CategoryId = i,
                CategoryName = <span style="color: blue;">string</span>.Format(<span style="color: #a31515;">"Category {0}"</span>, i)
            });
        }

        <span style="color: blue;">return</span> categories;
    }
}</pre>
</div>
</div>
<p><strong>Controller:</strong></p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> SampleController : Controller {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> ProductCategoryRepo productCategoryRepo = 
         <span style="color: blue;">new</span> ProductCategoryRepo();

    <span style="color: blue;">public</span> ActionResult Index() {

        registerProductCategorySelectListViewBag();
        <span style="color: blue;">return</span> View();
    }

    <span style="color: blue;">private</span> <span style="color: blue;">void</span> registerProductCategorySelectListViewBag() {

        ViewBag.ProductCategorySelectList = 
            productCategoryRepo.GetAll().Select(
                c =&gt; <span style="color: blue;">new</span> SelectListItem { 
                    Text = c.CategoryName,
                    Value = c.CategoryId.ToString()
                }
            );
    }
}</pre>
</div>
</div>
<p><strong>View:</strong></p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>@{
    ViewBag.Title = <span style="color: #a31515;">"HTML Select List Sample"</span>;
}

&lt;h2&gt;HTML Select List Sample&lt;/h2&gt;

&lt;p&gt;
    &lt;strong&gt;Product Categories:&lt;/strong&gt;
&lt;/p&gt;
&lt;p&gt;
    @Html.DropDownList(
        <span style="color: #a31515;">"ProductCategoryId"</span>, 
        (IEnumerable&lt;SelectListItem&gt;)ViewBag.ProductCategorySelectList
    )
&lt;/p&gt;</pre>
</div>
</div>
<p>Let&rsquo;s see what it generates for us:</p>
<p><strong>Output HTML of Select element:</strong></p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">&lt;</span><span style="color: #a31515;">select</span> <span style="color: red;">id</span><span style="color: blue;">=</span><span style="color: blue;">"ProductCategoryId"</span> <span style="color: red;">name</span><span style="color: blue;">=</span><span style="color: blue;">"ProductCategoryId"</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"1"</span><span style="color: blue;">&gt;</span>Category 1<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"2"</span><span style="color: blue;">&gt;</span>Category 2<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"3"</span><span style="color: blue;">&gt;</span>Category 3<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"4"</span><span style="color: blue;">&gt;</span>Category 4<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"5"</span><span style="color: blue;">&gt;</span>Category 5<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"6"</span><span style="color: blue;">&gt;</span>Category 6<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"7"</span><span style="color: blue;">&gt;</span>Category 7<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"8"</span><span style="color: blue;">&gt;</span>Category 8<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"9"</span><span style="color: blue;">&gt;</span>Category 9<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
    <span style="color: blue;">&lt;</span><span style="color: #a31515;">option</span> <span style="color: red;">value</span><span style="color: blue;">=</span><span style="color: blue;">"10"</span><span style="color: blue;">&gt;</span>Category 10<span style="color: blue;">&lt;/</span><span style="color: #a31515;">option</span><span style="color: blue;">&gt;</span>
<span style="color: blue;">&lt;/</span><span style="color: #a31515;">select</span><span style="color: blue;">&gt;</span></pre>
</div>
</div>
<p>Just like what I expected. On the controller, you see that I newed up <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.selectlistitem(v=VS.98).aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.selectlistitem(v=VS.98).aspx">SelectListItem</a> classes and pass them to the view through the <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.viewdatadictionary(v=VS.98).aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.viewdatadictionary(v=VS.98).aspx">ViewDataDictionary</a>. SelectListItem class works nicely with <a target="_blank" href="http://msdn.microsoft.com/en-us/library/dd492738(v=VS.98).aspx" title="http://msdn.microsoft.com/en-us/library/dd492738(v=VS.98).aspx">DropDownList</a> and <a target="_blank" href="http://msdn.microsoft.com/en-us/library/ee703462(v=VS.98).aspx" title="http://msdn.microsoft.com/en-us/library/ee703462(v=VS.98).aspx">DropDownListFor</a> Html helpers within ASP.NET MVC as you see inside the view code above.</p>
<p>So, don&rsquo;t try to reinvent the wheel and don&rsquo;t get into for and foreach loops in order just to create a select element. SelectListItem has been put inside the library just for this purpose and leverage it.</p>