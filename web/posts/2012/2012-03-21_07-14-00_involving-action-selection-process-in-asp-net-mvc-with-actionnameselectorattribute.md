---
id: a02fdb76-3f3c-447c-bc56-aee37a54e879
title: Involving Action Selection Process in ASP.NET MVC with ActionNameSelectorAttribute
abstract: We will see how we can involve action selection process in ASP.NET MVC with
  ActionNameSelectorAttribute with a real world use case scenario.
created_at: 2012-03-21 07:14:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET MVC
slugs:
- involving-action-selection-process-in-asp-net-mvc-with-actionnameselectorattribute
---

<p>In ASP.NET MVC, out of the box Controller class provides us a nice way to work with the framework. One of the advantages of using this class as a base controller class is that it provides so much nice functionality. One of those is <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionnameselectorattribute.aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionnameselectorattribute.aspx">ActionNameSelectorAttribute</a> class.</p>
<p>This class <em>represents an attribute that affects the selection of an action method. </em><a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionnameattribute.aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionnameattribute.aspx">ActionNameAttribute</a> class is an implementation of this abstract class and provides an ability to catch requests which comes to a particular action. Here is a sample:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> HomeController : Controller { 

    [ActionName(<span style="color: #a31515;">"FooBar"</span>)]
    <span style="color: blue;">public</span> ViewResult Foo() { 
    
        <span style="color: blue;">return</span> View();
    }
}</pre>
</div>
</div>
<p>We have an action method named <em>Foo</em> and we know that MVC Framework will pick the method which has the same name as the action route parameter. In this case, we expect <em>Foo</em> method to be invoked if we hit <strong>/Home/Foo</strong> but it is not going to be because we supplied the ActionNameAttribute to involve in the action selection process and tell it to pick actions which has the <em>FooBar</em> value.</p>
<p>You might be using ASP.NET MVC for a while now and have never used this feature so far but sometimes this might come in handy. Here is a weird use case which I needed to implement:</p>
<p>In an application, I made use of new <a target="_blank" href="https://developer.mozilla.org/en/DOM/Manipulating_the_browser_history" title="https://developer.mozilla.org/en/DOM/Manipulating_the_browser_history">JavaScript pushstate and popstate</a> features but I wanted to gracefully handle this. I implemented some sort of logic at client side and server side. Finally, I got it right but there was a problem:</p>
<p>I have an Index action which excepts all GET and POST requests but I wanted to invoke some other function if the request comes as POST and is an none-ajax request. there are some ways to handle this, just like checking the request method and if the request is an Ajax request but it felt so dirty to me. So, I decided to take advantage of ActionNameSelectorAttribute.</p>
<p>Here is the implementation:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[AttributeUsage(
    AttributeTargets.Method, AllowMultiple = <span style="color: blue;">false</span>, Inherited = <span style="color: blue;">true</span>)]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> NoneAjaxActionNameAttribute : ActionNameSelectorAttribute {

    <span style="color: blue;">public</span> NoneAjaxActionNameAttribute(<span style="color: blue;">string</span> name) {
        <span style="color: blue;">if</span> (String.IsNullOrEmpty(name)) {
            <span style="color: blue;">throw</span> <span style="color: blue;">new</span> ArgumentException(<span style="color: #a31515;">"Name paramater is null"</span>, <span style="color: #a31515;">"name"</span>);
        }

        Name = name;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Name {
        <span style="color: blue;">get</span>;
        <span style="color: blue;">private</span> <span style="color: blue;">set</span>;
    }

    <span style="color: blue;">public</span> <span style="color: blue;">override</span> <span style="color: blue;">bool</span> IsValidName(ControllerContext controllerContext, 
        <span style="color: blue;">string</span> actionName, MethodInfo methodInfo) {

        <span style="color: blue;">return</span> 
            String.Equals(actionName, Name, StringComparison.OrdinalIgnoreCase) &amp;&amp;
            !controllerContext.HttpContext.Request.IsAjaxRequest();
    }
}</pre>
</div>
</div>
<p>As you see, it implements the <em>ActionNameSelectorAttribute</em> class and overrides only one method which is <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionnameselectorattribute.isvalidname(v=vs.98).aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.actionnameselectorattribute.isvalidname(v=vs.98).aspx">IsValidName</a>. Inside that method, we should decide whether the action is in a valid state to be invoked or not. In our case, it checks the action name against the supplied name and if the request is an Ajax request or not.</p>
<p>Her is my controller which made use of this attribute:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> ContentSearchController : Controller {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IContentSearchService _contentSearchService;
    <span style="color: blue;">const</span> <span style="color: blue;">int</span> pageSize = 10;

    <span style="color: blue;">public</span> ContentSearchController(IContentSearchService contentSearchService) {

        _contentSearchService = contentSearchService;
    }

    <span style="color: blue;">public</span> ActionResult Index(<span style="color: blue;">string</span> q, <span style="color: blue;">int</span> page = 1) {

        <span style="color: blue;">var</span> model = _contentSearchService.Search(q, page, pageSize);

        <span style="color: blue;">if</span> (Request.IsAjaxRequest()) {
        
            <span style="color: blue;">return</span> Json(<span style="color: blue;">new</span> { 
                    data = <span style="color: blue;">this</span>.RenderPartialViewToString(<span style="color: #a31515;">"_SearchResult"</span>, model) 
            });
        }

        <span style="color: blue;">return</span> View(
            model
        );
    }

    [NoneAjaxActionName(<span style="color: #a31515;">"Index"</span>), HttpPost]
    <span style="color: blue;">public</span> RedirectToRouteResult Index_post(<span style="color: blue;">string</span> searchTerm) {

        <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"index"</span>, <span style="color: blue;">new</span> { q = searchTerm });
    }
}</pre>
</div>
</div>
<p>As you see, I also added <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.httppostattribute(v=vs.98).aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.httppostattribute(v=vs.98).aspx">HttpPostAttribute</a> to pick only POST requests.</p>
<p>By the help of a little bit of code, we suddenly find ourselves in the middle of action selection process and I think it is pretty powerful even if it&rsquo;s a cheesy implementation.</p>