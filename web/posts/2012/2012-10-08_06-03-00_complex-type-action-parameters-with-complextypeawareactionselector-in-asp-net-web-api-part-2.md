---
title: Complex Type Action Parameters with ComplexTypeAwareActionSelector in ASP.NET
  Web API - Part 2
abstract: In this post, we will see how ComplexTypeAwareActionSelector behaves under
  the covers to involve complex type action parameters during the action selection
  process.
created_at: 2012-10-08 06:03:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
slugs:
- complex-type-action-parameters-with-complextypeawareactionselector-in-asp-net-web-api-part-2
---

<p>In my previous post on <a href="http://www.tugberkugurlu.com/archive/complex-type-action-parameters-with-complextypeawareactionselector-in-asp-net-web-api-part-1">complex type action parameters with ComplexTypeAwareActionSelector in ASP.NET Web API</a>, I showed how to leverage the ComplexTypeAwareActionSelector from <a href="https://github.com/WebAPIDoodle/WebAPIDoodle">WebAPIDoodle</a> project to involve complex type action parameters during the action selection process. In this post, we will go into details to see how ComplexTypeAwareActionSelector behaves under the covers.</p>
<p>Assuming that we have a Person class as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> Person {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> FooBar;

    <span style="color: blue;">public</span> Nullable&lt;<span style="color: blue;">int</span>&gt; Id { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Surname { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
        
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> FullName { 

        <span style="color: blue;">get</span> {
            <span style="color: blue;">return</span> <span style="color: blue;">string</span>.Format(<span style="color: #a31515;">"{0} {1}"</span>, Name, Surname);
        }
    }

    [BindingInfo(NoBinding = <span style="color: blue;">true</span>)]
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> FullName2 { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">public</span> Country Country { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">internal</span> <span style="color: blue;">int</span> Foo { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }

    <span style="color: blue;">public</span> <span style="color: blue;">bool</span> IsLegitPerson() {

        <span style="color: blue;">return</span> Name.Equals(<span style="color: #a31515;">"tugberk"</span>, StringComparison.OrdinalIgnoreCase);
    }
}

<span style="color: blue;">public</span> <span style="color: blue;">class</span> Country {

    <span style="color: blue;">public</span> <span style="color: blue;">int</span> Id { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> Name { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
    <span style="color: blue;">public</span> <span style="color: blue;">string</span> ISOCode { <span style="color: blue;">get</span>; <span style="color: blue;">set</span>; }
}</pre>
</div>
</div>
<p>In a real world scenario, we wouldn't use the Person class to bind its values from the URI but bare with me for sake of this demo. Person class has five publicly-settable properties: Id, Name, Surname, FullName2 and Country. It also has internally-settable property called Foo. There are also one read-only property (FullName) and one public field (FooBar). Besides those, the FullName2 property has been marked with WebAPIDoodle.BindingInfoAttribute by setting its NoBinding property to true.</p>
<p>Let&rsquo;s also assume that we have the below controller and action.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> FooBarController : ApiController {

    <span style="color: blue;">public</span> IEnumerable&lt;FooBar&gt; Get([FromUri]Person person) { 

        <span style="color: green;">//...</span>
    }
}</pre>
</div>
</div>
<p>Now, the question is how ComplexTypeAwareActionSelector behaves here and which members of the Person class are going to be involved to perform the action selection. To perform this logic, the ComplexTypeAwareActionSelector uses two helper methods as below:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">internal</span> <span style="color: blue;">static</span> <span style="color: blue;">class</span> TypeHelper {

    <span style="color: blue;">internal</span> <span style="color: blue;">static</span> <span style="color: blue;">bool</span> IsSimpleType(Type type) {

        <span style="color: blue;">return</span> type.IsPrimitive ||
                type.Equals(<span style="color: blue;">typeof</span>(<span style="color: blue;">string</span>)) ||
                type.Equals(<span style="color: blue;">typeof</span>(DateTime)) ||
                type.Equals(<span style="color: blue;">typeof</span>(Decimal)) ||
                type.Equals(<span style="color: blue;">typeof</span>(Guid)) ||
                type.Equals(<span style="color: blue;">typeof</span>(DateTimeOffset)) ||
                type.Equals(<span style="color: blue;">typeof</span>(TimeSpan));
    }

    <span style="color: blue;">internal</span> <span style="color: blue;">static</span> <span style="color: blue;">bool</span> IsSimpleUnderlyingType(Type type) {

        Type underlyingType = Nullable.GetUnderlyingType(type);
        <span style="color: blue;">if</span> (underlyingType != <span style="color: blue;">null</span>) {
            type = underlyingType;
        }

        <span style="color: blue;">return</span> TypeHelper.IsSimpleType(type);
    }
}</pre>
</div>
</div>
<p>These two methods belong to <a href="http://www.asp.net/web-api">ASP.NET Web API</a> source code but they are internal. So, I ported them to my project as they are. As you can see, IsSimpleType method accepts a Type parameter and determines if the type is a simple or primitive type. The IsSimpleUnderlyingType method, on the other hand, looks if the Type is Nullable type. If so, it looks at the underlying type to see if it is a simple type or not. This is how the ComplexTypeAwareActionSelector determines if a parameter is simple type or not.</p>
<p>When the ComplexTypeAwareActionSelector sees a complex type action parameter, it hands that type to another private method to get the useable properties. To mimic how that private helper method filters the properties, I created the a little console application which holds the actual filter logic.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white;">
<pre><span style="color: blue;">class</span> Program {

    <span style="color: blue;">static</span> <span style="color: blue;">void</span> Main(<span style="color: blue;">string</span>[] args) {

        Console.Write(Environment.NewLine);
        Console.WriteLine(<span style="color: #a31515;">"==============================================="</span>);
        Console.WriteLine(<span style="color: #a31515;">"========This is the actual logic in use========"</span>);
        Console.WriteLine(<span style="color: #a31515;">"==============================================="</span>);

        <span style="color: blue;">var</span> propInfos = <span style="color: blue;">from</span> propInfo <span style="color: blue;">in</span> <span style="color: blue;">typeof</span>(Person).GetProperties()
                        <span style="color: blue;">where</span> TypeHelper
                              .IsSimpleUnderlyingType(propInfo.PropertyType) &amp;&amp; 
                              propInfo.GetSetMethod(<span style="color: blue;">false</span>) != <span style="color: blue;">null</span>
                              
                        <span style="color: blue;">let</span> noBindingAttr = propInfo
                            .GetCustomAttributes().FirstOrDefault(attr =&gt; 
                                attr.GetType() == <span style="color: blue;">typeof</span>(BindingInfoAttribute)) 
                                    <span style="color: blue;">as</span> BindingInfoAttribute
                                    
                        <span style="color: blue;">where</span> (noBindingAttr != <span style="color: blue;">null</span>) 
                              ? noBindingAttr.NoBinding == <span style="color: blue;">false</span> 
                              : <span style="color: blue;">true</span>
                              
                        <span style="color: blue;">select</span> propInfo;

        <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> _propInfo <span style="color: blue;">in</span> propInfos) {

            Console.WriteLine(_propInfo.Name);
        }

        Console.ReadLine();
    }
}</pre>
</div>
</div>
<p>Here is what it is doing here for the Person type:</p>
<ul>
<li>It first gets all the public properties of the Person class with <a href="http://msdn.microsoft.com/en-us/library/aky14axb.aspx">GetProperties method of the Type class</a>. So, the ForBar field and the Foo property is ignored.</li>
<li>Secondly, it looks if the property is simple underlying type and publicly-settable. If any one of them is not applicable, it ignores them. In our case here, the FullName property, which is a read-only property, and the Country property, which is a complex type property, are ignored.</li>
<li>As a last step, it looks at the attributes of the each filtered property. If the property is marked with WebAPIDoodle.BindingInfoAttribute and the BindingInfoAttribute&rsquo;s NoBinding property is set to true, the property gets ignored. In our case here, the FullName2 property will be ignored.</li>
</ul>
<p>If we run this little console application, we will see the following result:</p>
<p><a href="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action.NET-Web-API---Part-2_F18F/SNAGHTML12a04389.png"><img height="167" width="644" src="http://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Complex-Type-Action.NET-Web-API---Part-2_F18F/SNAGHTML12a04389_thumb.png" alt="SNAGHTML12a04389" border="0" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" title="SNAGHTML12a04389" /></a></p>
<p>As a result, the Id, Name and Surname properties will be considered during the action selection. I would like to point out couple of things before finishing up this post:</p>
<ul>
<li>This is a one-time operation per controller action. For example, when we first fire up our ASP.NET Web API application and send a request which is eventually going to correspond to FooBarController, the action selector will look at all the actions under the FooBarController and performs the above logic along with others and cache lots of stuff including the above logic. So, when you hit the FooBarController next time, this process won&rsquo;t be run and the result will be pulled directly from the cache.</li>
<li>The WebAPIDoodle.BindingInfoAttribute lives inside a separate package named as <a href="http://nuget.org/packages/WebAPIDoodle.Meta">WebAPIDoodle.Meta</a>. The WebAPIDoodle.Meta package contains some runtime components such as Attributes, Interfaces. This assembly has no dependency on ASP.NET Web API so that it would be easy to reference this on Model or Domain Layer projects.</li>
</ul>
<p>Happy coding and give feedback for this little feature :)</p>