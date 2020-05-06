---
title: How to Work With Generic Repositories on ASP.NET MVC and Unit Testing Them
  By Mocking
abstract: In this blog post we will see how to work with generic repositories on ASP.NET
  MVC and unit testing them by mocking with moq
created_at: 2011-12-22 17:05:00 +0000 UTC
tags:
- ASP.NET MVC
- C#
- DbContext
- Entity Framework
- Unit Testing
slugs:
- how-to-work-with-generic-repositories-on-asp-net-mvc-and-unit-testing-them-by-mocking
- how-to-work-with-generic-repositories-on-asp-net-mvc-and-unit-t
---

<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/How-to-Wor.NET-MVC-and-Unit-Testing-Them_FE36/image.png"><img height="244" width="155" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/How-to-Wor.NET-MVC-and-Unit-Testing-Them_FE36/image_thumb.png" align="left" alt="image" border="0" title="image" style="background-image: none; margin: 0px 15px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border: 0px;" /></a></p>
<p>I have blogged about <a target="_blank" href="http://www.tugberkugurlu.com/archive/generic-repository-pattern-entity-framework-asp-net-mvc-and-unit-testing-triangle" title="http://www.tugberkugurlu.com/archive/generic-repository-pattern-entity-framework-asp-net-mvc-and-unit-testing-triangle">Generic Repository Pattern - Entity Framework, ASP.NET MVC and Unit Testing Triangle</a> but that blog post only contains how we can implement that pattern on our DAL (Data Access Layer) project. Now, let&rsquo;s see how it fits into our <a target="_blank" href="http://asp.net/mvc" title="http://asp.net/mvc">ASP.NET MVC</a> application.</p>
<p>I created a sample project which basically has CRUD operations for <strong>Foo</strong> class and I put the complete code of the project on <a target="_blank" href="https://github.com" title="https://github.com">GitHub</a>. <a target="_blank" href="https://github.com/tugberkugurlu/GenericRepoWebApp" title="https://github.com/tugberkugurlu/GenericRepoWebApp">https://github.com/tugberkugurlu/GenericRepoWebApp</a></p>
<p>So, in the previous blog post we have created our DAL project and repository classes based on our generic repository. I have created an ASP.NET MVC 3 Web Application project and a Visual Studio Test Project besides the DAL project. On the left had side, you can see how the project structure looks like inside solution explorer.</p>
<p>When we think about the project, we have repository classes which implements <strong>GenericRepositroy&lt;C, T&gt;</strong> abstract class and individual repository interfaces. So, we need a way to inject our concrete repository classes into our ASP.NET MVC application so that we make our application loosely coupled which means unit testing friendly.</p>
<p>Fortunately, ASP.NET MVC has been built with unit testing in mind so it makes dependency injection very easy with <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.dependencyresolver(v=vs.98).aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.dependencyresolver(v=vs.98).aspx">DependencyResolver</a> class. DependencyResolver class Provides a registration point for dependency resolvers that implement <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.mvc.idependencyresolver(v=vs.98).aspx" title="http://msdn.microsoft.com/en-us/library/system.web.mvc.idependencyresolver(v=vs.98).aspx">IDependencyResolver</a> or the Common Service Locator IServiceLocator interface. But we won&rsquo;t be dealing with that class at all. Instead, we will use a third party dependency injector called <a target="_blank" href="http://ninject.org/" title="http://ninject.org/">Ninject</a>. Ninject also has a package which benefits from DependencyResolver of ASP.NET MVC. We will use <a target="_blank" href="http://nuget.org" title="http://nuget.org">Nuget</a> to bring down Ninject.</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/How-to-Wor.NET-MVC-and-Unit-Testing-Them_FE36/image_3.png"><img height="67" width="640" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/How-to-Wor.NET-MVC-and-Unit-Testing-Them_FE36/image_thumb_3.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>After we install the package, we will see a folder named App_Start added to our project. Inside that folder, open up the NinjectMVC3.cs file and go to RegisterServices method. In our case, here what we do:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: gray;">///</span> <span style="color: gray;">&lt;summary&gt;</span>
<span style="color: gray;">///</span><span style="color: green;"> Load your modules or register your services here!</span>
<span style="color: gray;">///</span> <span style="color: gray;">&lt;/summary&gt;</span>
<span style="color: gray;">///</span> <span style="color: gray;">&lt;param name="kernel"&gt;</span><span style="color: green;">The kernel.&lt;/param&gt;</span>
<span style="color: blue;">private</span> <span style="color: blue;">static</span> <span style="color: blue;">void</span> RegisterServices(IKernel kernel) {

    kernel.Bind&lt;IFooRepository&gt;().To&lt;FooRepository&gt;();
    kernel.Bind&lt;IBarRepository&gt;().To&lt;BarReposiltory&gt;();
}</pre>
</div>
</div>
<p>This does some stuff behind the scenes and I won&rsquo;t go into details here but I really recommend you to go and take a look at series of blog posts which <a target="_blank" href="http://bradwilson.typepad.com" title="http://bradwilson.typepad.com">Brad Wilson</a> has done on <a target="_blank" href="http://bradwilson.typepad.com/blog/2010/07/service-location-pt1-introduction.html" title="http://bradwilson.typepad.com/blog/2010/07/service-location-pt1-introduction.html">ASP.NET MVC 3 Service Location</a>. However, if we try to explain it with simplest words, Ninject news up the controller classes with parameters which we specify. In this case, if a controller constructor accepts a parameter which is type of IFooRepository, Ninject will give it FooRepository class and news it up. We will see why this is useful on unit testing stage.</p>
<p>When we look inside the RegisterServices method, we don&rsquo;t see neither GenericRepository nor IGenericRepository because the way we implement them enables them to work behind the scenes.</p>
<p>As for the implementation, here how the controller looks like:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> FooController : Controller {

    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IFooRepository _fooRepo;
    
    <span style="color: blue;">public</span> FooController(IFooRepository fooRepo) {
        _fooRepo = fooRepo;
    }

    <span style="color: blue;">public</span> ViewResult Index() {

        <span style="color: blue;">var</span> model = _fooRepo.GetAll();
        <span style="color: blue;">return</span> View(model);
    }

    <span style="color: blue;">public</span> ActionResult Details(<span style="color: blue;">int</span> id) {

        <span style="color: blue;">var</span> model = _fooRepo.GetSingle(id);
        <span style="color: blue;">if</span> (model == <span style="color: blue;">null</span>)
            <span style="color: blue;">return</span> HttpNotFound();

        <span style="color: blue;">return</span> View(model);
    }

    <span style="color: blue;">public</span> ActionResult Edit(<span style="color: blue;">int</span> id) {

        <span style="color: blue;">var</span> model = _fooRepo.GetSingle(id);
        <span style="color: blue;">if</span> (model == <span style="color: blue;">null</span>)
            <span style="color: blue;">return</span> HttpNotFound();

        <span style="color: blue;">return</span> View(model);
    }

    [ActionName(<span style="color: #a31515;">"Edit"</span>), HttpPost]
    <span style="color: blue;">public</span> ActionResult Edit_post(Foo foo) {

        <span style="color: blue;">if</span> (ModelState.IsValid) {

            <span style="color: blue;">try</span> {
                _fooRepo.Edit(foo);
                _fooRepo.Save();

                <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"details"</span>, <span style="color: blue;">new</span> { id = foo.FooId });

            } <span style="color: blue;">catch</span> (Exception ex) {
                ModelState.AddModelError(<span style="color: blue;">string</span>.Empty, <span style="color: #a31515;">"Something went wrong. Message: "</span> + ex.Message);
            }
        }

        <span style="color: green;">//If we come here, something went wrong. Return it back.</span>
        <span style="color: blue;">return</span> View(foo);
    }

    <span style="color: blue;">public</span> ActionResult Create() {

        <span style="color: blue;">return</span> View();
    }

    [ActionName(<span style="color: #a31515;">"Create"</span>), HttpPost]
    <span style="color: blue;">public</span> ActionResult Create_post(Foo foo) {

        <span style="color: blue;">if</span> (ModelState.IsValid) {

            <span style="color: blue;">try</span> {
                _fooRepo.Add(foo);
                _fooRepo.Save();

                <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"details"</span>, <span style="color: blue;">new</span> { id = foo.FooId });

            } <span style="color: blue;">catch</span> (Exception ex) {
                ModelState.AddModelError(<span style="color: blue;">string</span>.Empty, <span style="color: #a31515;">"Something went wrong. Message: "</span> + ex.Message);
            }
        }

        <span style="color: green;">//If we come here, something went wrong. Return it back.</span>
        <span style="color: blue;">return</span> View(foo);

    }

    <span style="color: blue;">public</span> ActionResult Delete(<span style="color: blue;">int</span> id) {

        <span style="color: blue;">var</span> model = _fooRepo.GetSingle(id);
        <span style="color: blue;">if</span> (model == <span style="color: blue;">null</span>)
            <span style="color: blue;">return</span> HttpNotFound();

        <span style="color: blue;">return</span> View(model);
    }

    [ActionName(<span style="color: #a31515;">"Delete"</span>), HttpPost]
    <span style="color: blue;">public</span> ActionResult Delete_post(<span style="color: blue;">int</span> id) {

        <span style="color: blue;">var</span> model = _fooRepo.GetSingle(id);
        <span style="color: blue;">if</span> (model == <span style="color: blue;">null</span>)
            <span style="color: blue;">return</span> HttpNotFound();

        _fooRepo.Delete(model);
        _fooRepo.Save();

        <span style="color: blue;">return</span> RedirectToAction(<span style="color: #a31515;">"Index"</span>);
    }
}</pre>
</div>
</div>
<p>All of the action methods do some very basic stuff but one thing to notice here is the below code:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>    <span style="color: blue;">private</span> <span style="color: blue;">readonly</span> IFooRepository _fooRepo;
    
    <span style="color: blue;">public</span> FooController(IFooRepository fooRepo) {
        _fooRepo = fooRepo;
    }</pre>
</div>
</div>
<p>As I mentioned before, controller constructor accepts a parameter which is type of IFooRepository and inside the constructor method, we expose the parameter for internal use of that controller class.</p>
<p>I have some views which corresponds to each action method and they work as expected.</p>
<p><strong>Unit Testing</strong></p>
<p>So, how do we unit test that controller without connecting to our database? When we think theoretically, what we need is a fake repository which implements IFooRepository interface so that we can pass that fake repository into our controller as a constructor parameter. Pay attention here that we still has no need for neither generic repository interface nor generic repository abstract class. We just need to fake FooRepository interface with fake data.</p>
<p>We will do this by mocking and creating <a target="_blank" href="http://en.wikipedia.org/wiki/Mock_object" title="http://en.wikipedia.org/wiki/Mock_object">mock objects</a>. In order to do that, we will benefit from an awesome library called <a target="_blank" href="http://code.google.com/p/moq" title="http://code.google.com/p/moq">Moq</a>.</p>
<blockquote>
<p>As you can see inside the project on GitHub, I didn&rsquo;t use NuGet to bring down the Moq because I tried and it failed over and over again. So, I put that inside the lib folder under the root directory and reference it from there.</p>
</blockquote>
<p>After you reference the Moq library inside the Test application, create a class named FooControllerTest.cs. Here how it should look like at first:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[TestClass]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> FooControllerTest {

}</pre>
</div>
</div>
<p>So empty. We will start to fill it in with a mock of IFooReporsitory. Below, you can see the complete code which enables that:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[TestClass]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> FooControllerTest {

    <span style="color: blue;">private</span> IFooRepository fooRepo;

    [TestInitialize]
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Initialize() {

        <span style="color: green;">//Mock repository creation</span>
        Mock&lt;IFooRepository&gt; mock = <span style="color: blue;">new</span> Mock&lt;IFooRepository&gt;();
        mock.Setup(m =&gt; m.GetAll()).Returns(<span style="color: blue;">new</span>[] { 
            <span style="color: blue;">new</span> Foo { FooId = 1, FooName = <span style="color: #a31515;">"Fake Foo 1"</span> },
            <span style="color: blue;">new</span> Foo { FooId = 2, FooName = <span style="color: #a31515;">"Fake Foo 2"</span> },
            <span style="color: blue;">new</span> Foo { FooId = 3, FooName = <span style="color: #a31515;">"Fake Foo 3"</span> },
            <span style="color: blue;">new</span> Foo { FooId = 4, FooName = <span style="color: #a31515;">"Fake Foo 4"</span> }
        }.AsQueryable());

        mock.Setup(m =&gt; 
            m.GetSingle(
                It.Is&lt;<span style="color: blue;">int</span>&gt;(i =&gt; 
                    i == 1 || i == 2 || i == 3 || i == 4
                )
            )
        ).Returns&lt;<span style="color: blue;">int</span>&gt;(r =&gt; <span style="color: blue;">new</span> Foo { 
            FooId = r,
            FooName = <span style="color: blue;">string</span>.Format(<span style="color: #a31515;">"Fake Foo {0}"</span>, r)
        });

        fooRepo = mock.Object;
    }

}</pre>
</div>
</div>
<p>This Initialize method will run before all of the test methods run so we can work inside that method in order to mock our object.</p>
<p>In here, first I have setup a mock for <strong>GetAll</strong> method result and it returns 4 instances or Foo class.</p>
<p>Second, I do the same thing for <strong>GetSingle</strong> method. It looks a little different because it accepts a parameter type of <strong>int</strong>. What I am telling there is that: there are 4 instances I have here and if the parameter matches one of those instances, it will returns a Foo class which is a type of <strong>GetSingle</strong> method returns.</p>
<p>Lastly, I expose the mock object for internal use for the test class.</p>
<p>Now we have an IFooRepository instance we can work with. We have completed the hard part and now we can start writing our unit tests. Here is some of possible unit tests that we can have:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>[TestClass]
<span style="color: blue;">public</span> <span style="color: blue;">class</span> FooControllerTest {

    <span style="color: blue;">private</span> IFooRepository fooRepo;

    [TestInitialize]
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> Initialize() {

        <span style="color: green;">//Mock repository creation</span>
        Mock&lt;IFooRepository&gt; mock = <span style="color: blue;">new</span> Mock&lt;IFooRepository&gt;();
        mock.Setup(m =&gt; m.GetAll()).Returns(<span style="color: blue;">new</span>[] { 
            <span style="color: blue;">new</span> Foo { FooId = 1, FooName = <span style="color: #a31515;">"Fake Foo 1"</span> },
            <span style="color: blue;">new</span> Foo { FooId = 2, FooName = <span style="color: #a31515;">"Fake Foo 2"</span> },
            <span style="color: blue;">new</span> Foo { FooId = 3, FooName = <span style="color: #a31515;">"Fake Foo 3"</span> },
            <span style="color: blue;">new</span> Foo { FooId = 4, FooName = <span style="color: #a31515;">"Fake Foo 4"</span> }
        }.AsQueryable());

        mock.Setup(m =&gt; 
            m.GetSingle(
                It.Is&lt;<span style="color: blue;">int</span>&gt;(i =&gt; 
                    i == 1 || i == 2 || i == 3 || i == 4
                )
            )
        ).Returns&lt;<span style="color: blue;">int</span>&gt;(r =&gt; <span style="color: blue;">new</span> Foo { 
            FooId = r,
            FooName = <span style="color: blue;">string</span>.Format(<span style="color: #a31515;">"Fake Foo {0}"</span>, r)
        });

        fooRepo = mock.Object;
    }

    [TestMethod]
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> is_index_returns_model_type_of_iqueryable_foo() {
        
        <span style="color: green;">//Arrange</span>
        <span style="color: green;">//Create the controller instance</span>
        FooController fooController = <span style="color: blue;">new</span> FooController(fooRepo);

        <span style="color: green;">//Act</span>
        <span style="color: blue;">var</span> indexModel = fooController.Index().Model;

        <span style="color: green;">//Assert</span>
        Assert.IsInstanceOfType(indexModel, <span style="color: blue;">typeof</span>(IQueryable&lt;Foo&gt;));
    }

    [TestMethod]
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> is_index_returns_iqueryable_foo_count_of_4() {

        <span style="color: green;">//Arrange</span>
        <span style="color: green;">//Create the controller instance</span>
        FooController fooController = <span style="color: blue;">new</span> FooController(fooRepo);

        <span style="color: green;">//Act</span>
        <span style="color: blue;">var</span> indexModel = (IQueryable&lt;<span style="color: blue;">object</span>&gt;)fooController.Index().Model;

        <span style="color: green;">//Assert</span>
        Assert.AreEqual&lt;<span style="color: blue;">int</span>&gt;(4, indexModel.Count());
    }

    [TestMethod]
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> is_details_returns_type_of_ViewResult() {

        <span style="color: green;">//Arrange</span>
        <span style="color: green;">//Create the controller instance</span>
        FooController fooController = <span style="color: blue;">new</span> FooController(fooRepo);

        <span style="color: green;">//Act</span>
        <span style="color: blue;">var</span> detailsResult = fooController.Details(1);

        <span style="color: green;">//Assert</span>
        Assert.IsInstanceOfType(detailsResult, <span style="color: blue;">typeof</span>(ViewResult));
    }

    [TestMethod]
    <span style="color: blue;">public</span> <span style="color: blue;">void</span> is_details_returns_type_of_HttpNotFoundResult() { 

        <span style="color: green;">//Arrange</span>
        <span style="color: green;">//Create the controller instance</span>
        FooController fooController = <span style="color: blue;">new</span> FooController(fooRepo);

        <span style="color: green;">//Act</span>
        <span style="color: blue;">var</span> detailsResult = fooController.Details(5);

        <span style="color: green;">//Assert</span>
        Assert.IsInstanceOfType(detailsResult, <span style="color: blue;">typeof</span>(HttpNotFoundResult));
    }
}</pre>
</div>
</div>
<p>When I run the all the test, I should see all of them pass and I do:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/How-to-Wor.NET-MVC-and-Unit-Testing-Them_FE36/image_4.png"><img height="196" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/How-to-Wor.NET-MVC-and-Unit-Testing-Them_FE36/image_thumb_4.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>I hope this blog post gives you an idea. Stay tuned for others <img src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/How-to-Wor.NET-MVC-and-Unit-Testing-Them_FE36/wlEmoticon-winkingsmile.png" alt="Winking smile" class="wlEmoticon wlEmoticon-winkingsmile" style="border-style: none;" /></p>
<blockquote>
<p><strong>NOTE:</strong></p>
<p>Entity Framework DbContext Generic Repository Implementation Is Now On Nuget and GitHub:&nbsp;<a target="_blank" href="http://www.tugberkugurlu.com/archive/entity-framework-dbcontext-generic-repository-implementation-is-now-on-nuget-and-github">http://www.tugberkugurlu.com/archive/entity-framework-dbcontext-generic-repository-implementation-is-now-on-nuget-and-github</a><strong><br /></strong></p>
</blockquote>