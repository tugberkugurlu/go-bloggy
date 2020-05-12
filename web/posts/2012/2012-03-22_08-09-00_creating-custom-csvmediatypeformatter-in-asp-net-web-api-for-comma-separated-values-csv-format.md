---
id: 0fcc0286-8087-4380-9a73-caaaf6b52872
title: Creating Custom CSVMediaTypeFormatter In ASP.NET Web API for Comma-Separated
  Values (CSV) Format
abstract: In this post, we will see how to create a custom CSVMediaTypeFormatter in
  ASP.NET Web API for comma-separated values (CSV) format
created_at: 2012-03-22 08:09:00 +0000 UTC
tags:
- ASP.Net
- ASP.NET Web API
- C#
slugs:
- creating-custom-csvmediatypeformatter-in-asp-net-web-api-for-comma-separated-values-csv-format
---

<p>As I tried to explain on my previous <a title="http://www.tugberkugurlu.com/archive/asp-net-web-api-mediatypeformatters-with-mediatypemappings" href="http://www.tugberkugurlu.com/archive/asp-net-web-api-mediatypeformatters-with-mediatypemappings" target="_blank">MediaTypeFormatters With MediaTypeMappings</a> post, formatters play a huge role inside the ASP.NET Web API processing pipeline. As Web API framework programming model is so similar to MVC framework, I kind of want to see formatters as views. Formatters handles serializing and deserializing strongly-typed objects into specific format.</p>
<p>I wanted to create CSVMediaTypeFormatter to hook it up for list of objects and I managed to get it working. After I created it, I saw the great <a title="http://www.asp.net/web-api/overview/formats-and-model-binding/media-formatters" href="http://www.asp.net/web-api/overview/formats-and-model-binding/media-formatters" target="_blank">Media Formatters</a> post on ASP.NET web site which does the same thing. I was like "Man, come on!" and I noticed that formatter meant to be used with a specific object, so I figured there is still a validity in my implementation.</p>
<p>Here is the drill:</p>
<p>First of all we need to create a class which will be derived from <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter(v=vs.108).aspx" target="_blank">MediaTypeFormatter</a> abstract class. Here is the class with its constructors:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">public</span> <span style="color: blue;">class</span> CSVMediaTypeFormatter : MediaTypeFormatter {

    <span style="color: blue;">public</span> CSVMediaTypeFormatter() {

        SupportedMediaTypes.Add(<span style="color: blue;">new</span> MediaTypeHeaderValue(<span style="color: #a31515;">"text/csv"</span>));
    }
    
    <span style="color: blue;">public</span> CSVMediaTypeFormatter(
        MediaTypeMapping mediaTypeMapping) : <span style="color: blue;">this</span>() {

        MediaTypeMappings.Add(mediaTypeMapping);
    }
    
    <span style="color: blue;">public</span> CSVMediaTypeFormatter(
        IEnumerable&lt;MediaTypeMapping&gt; mediaTypeMappings) : <span style="color: blue;">this</span>() {

        <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> mediaTypeMapping <span style="color: blue;">in</span> mediaTypeMappings) {
            MediaTypeMappings.Add(mediaTypeMapping);
        }
    }
}</pre>
</div>
</div>
<p>Above, no matter which constructor you use, we always add <strong>text/csv</strong> media type to be supported for this formatter. We also allow custom <strong>MediaTypeMappings </strong>to be injected.</p>
<p>Now, we need to override two methods: <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter.canwritetype(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter.canwritetype(v=vs.108).aspx" target="_blank">MediaTypeFormatter.CanWriteType</a> and <a title="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter.onwritetostreamasync(v=vs.108).aspx" href="http://msdn.microsoft.com/en-us/library/system.net.http.formatting.mediatypeformatter.onwritetostreamasync(v=vs.108).aspx" target="_blank">MediaTypeFormatter.OnWriteToStreamAsync</a>.</p>
<p>First of all, here is the CanWriteType method implementation. What this method needs to do is to determine if the type of the object is supported with this formatter or not in order to write it.</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">override</span> <span style="color: blue;">bool</span> CanWriteType(Type type) {

    <span style="color: blue;">if</span> (type == <span style="color: blue;">null</span>)
        <span style="color: blue;">throw</span> <span style="color: blue;">new</span> ArgumentNullException(<span style="color: #a31515;">"type"</span>);

    <span style="color: blue;">return</span> isTypeOfIEnumerable(type);
}

<span style="color: blue;">private</span> <span style="color: blue;">bool</span> isTypeOfIEnumerable(Type type) {

    <span style="color: blue;">foreach</span> (Type interfaceType <span style="color: blue;">in</span> type.GetInterfaces()) {

        <span style="color: blue;">if</span> (interfaceType == <span style="color: blue;">typeof</span>(IEnumerable))
            <span style="color: blue;">return</span> <span style="color: blue;">true</span>;
    }

    <span style="color: blue;">return</span> <span style="color: blue;">false</span>;
}</pre>
</div>
</div>
<p>What this does here is to check if the object has implemented the IEnumerable interface. If so, then it is cool with that and can format the object. If not, it will return false and framework will ignore this formatter for that particular request.</p>
<p>And finally, here is the actual implementation. We need to do some work with reflection here in order to get the property names and values out of the value parameter which is a type of object:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre><span style="color: blue;">protected</span> <span style="color: blue;">override</span> Task OnWriteToStreamAsync(
    Type type,
    <span style="color: blue;">object</span> value,
    Stream stream,
    HttpContentHeaders contentHeaders,
    FormatterContext formatterContext,
    TransportContext transportContext) {

    writeStream(type, value, stream, contentHeaders);
    <span style="color: blue;">var</span> tcs = <span style="color: blue;">new</span> TaskCompletionSource&lt;<span style="color: blue;">int</span>&gt;();
    tcs.SetResult(0);
    <span style="color: blue;">return</span> tcs.Task;
}

<span style="color: blue;">private</span> <span style="color: blue;">void</span> writeStream(Type type, <span style="color: blue;">object</span> value, Stream stream, HttpContentHeaders contentHeaders) {

    <span style="color: green;">//NOTE: We have check the type inside CanWriteType method</span>
    <span style="color: green;">//If request comes this far, the type is IEnumerable. We are safe.</span>

    Type itemType = type.GetGenericArguments()[0];

    StringWriter _stringWriter = <span style="color: blue;">new</span> StringWriter();

    _stringWriter.WriteLine(
        <span style="color: blue;">string</span>.Join&lt;<span style="color: blue;">string</span>&gt;(
            <span style="color: #a31515;">","</span>, itemType.GetProperties().Select(x =&gt; x.Name )
        )
    );

    <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> obj <span style="color: blue;">in</span> (IEnumerable&lt;<span style="color: blue;">object</span>&gt;)value) {

        <span style="color: blue;">var</span> vals = obj.GetType().GetProperties().Select(
            pi =&gt; <span style="color: blue;">new</span> { 
                Value = pi.GetValue(obj, <span style="color: blue;">null</span>)
            }
        );

        <span style="color: blue;">string</span> _valueLine = <span style="color: blue;">string</span>.Empty;

        <span style="color: blue;">foreach</span> (<span style="color: blue;">var</span> val <span style="color: blue;">in</span> vals) {

            <span style="color: blue;">if</span> (val.Value != <span style="color: blue;">null</span>) {

                <span style="color: blue;">var</span> _val = val.Value.ToString();

                <span style="color: green;">//Check if the value contans a comma and place it in quotes if so</span>
                <span style="color: blue;">if</span> (_val.Contains(<span style="color: #a31515;">","</span>))
                    _val = <span style="color: blue;">string</span>.Concat(<span style="color: #a31515;">"\""</span>, _val, <span style="color: #a31515;">"\""</span>);

                <span style="color: green;">//Replace any \r or \n special characters from a new line with a space</span>
                <span style="color: blue;">if</span> (_val.Contains(<span style="color: #a31515;">"\r"</span>))
                    _val = _val.Replace(<span style="color: #a31515;">"\r"</span>, <span style="color: #a31515;">" "</span>);
                <span style="color: blue;">if</span> (_val.Contains(<span style="color: #a31515;">"\n"</span>))
                    _val = _val.Replace(<span style="color: #a31515;">"\n"</span>, <span style="color: #a31515;">" "</span>);

                _valueLine = <span style="color: blue;">string</span>.Concat(_valueLine, _val, <span style="color: #a31515;">","</span>);

            } <span style="color: blue;">else</span> {

                _valueLine = <span style="color: blue;">string</span>.Concat(<span style="color: blue;">string</span>.Empty, <span style="color: #a31515;">","</span>);
            }
        }

        _stringWriter.WriteLine(_valueLine.TrimEnd(<span style="color: #a31515;">','</span>));
    }

    <span style="color: blue;">var</span> streamWriter = <span style="color: blue;">new</span> StreamWriter(stream);
        streamWriter.Write(_stringWriter.ToString());
}</pre>
</div>
</div>
<p>We are partially done. Now, we need to make use out of this. I registered this formatter into the pipeline with the following code inside Global.asax <strong>Application_Start </strong>method:</p>
<div class="code-wrapper border-shadow-1">
<div style="background-color: white; color: black;">
<pre>GlobalConfiguration.Configuration.Formatters.Add(
    <span style="color: blue;">new</span> CSVMediaTypeFormatter(
        <span style="color: blue;">new</span>  QueryStringMapping(<span style="color: #a31515;">"format"</span>, <span style="color: #a31515;">"csv"</span>, <span style="color: #a31515;">"text/csv"</span>)
    )
);</pre>
</div>
</div>
<p>On my sample application, when you navigate to <strong>/api/cars?format=csv</strong>, it will get you a CSV file but without an extension. Go ahead and add the <strong>csv</strong> extension. Then, open it with Excel and you should see something similar to below:</p>
<p><a href="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/8a20fb9e71f1_9BCB/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://www.tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/8a20fb9e71f1_9BCB/image_thumb.png" width="644" height="415" /></a></p>
<p>This implementation is also on my ASP.NET Web API package (<a title="http://nuget.org/packages/TugberkUg.Web.Http" href="http://nuget.org/packages/TugberkUg.Web.Http" target="_blank">TugberkUg.Web.Http</a>) and you can get it via <a title="http://nuget.org" href="http://nuget.org" target="_blank">Nuget</a>:</p>
<div class="nuget-badge">
<p><code>PM&gt; Install-Package TugberkUg.Web.Http -Pre </code></p>
</div>
<p>This package contains other stuff related to ASP.NET Web API. You can check out the source code on <a href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http">https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/TugberkUg.Web.Http</a>.</p>
<p>The sample I used here is also on <a href="http://github.com/">GitHub</a>:<a href="https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/samples/CSVMediaTypeFormatterSample">https://github.com/tugberkugurlu/ASPNETWebAPISamples/tree/master/TugberkUg.Web.Http/src/samples/CSVMediaTypeFormatterSample</a></p>
<p>There are some caveats, though. If your class has nested custom types, then this one does not support that. You will see that, type of the class will be printed under the particular column.</p>