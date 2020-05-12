---
id: 9b407915-6870-44da-a542-6acf87af57dd
title: 'ASP.NET Web Forms : Calling Web Service Page Methods Using JQuery'
abstract: In this blog post we will see how to consume a web page methods using JQuery
  on ASP.NET Web Forms and use ASP.NET page methods as services. You will find some
  cool stuff about other things as well :)
created_at: 2011-05-15 20:58:00 +0000 UTC
tags:
- .net
- ASP.Net
- C#
- JQuery
slugs:
- asp-net-web-forms---calling-web-service-page-methods-using-jquery
---

<p>As it was officially announced last year, Microsoft has been contributing code to <a title="http://jquery.com/" href="http://jquery.com/" target="_blank">JQuery</a> project for over a year. Lots of developers are now using JQuery in their Web Forms project as well.</p>
<p>In web forms, there is a narrowly known feature that you can put your tiny Web Services on your code behind file and call them via client side code easily. Let&rsquo;s look at how easy that.</p>
<p>We have following code inside our Default.aspx page&rsquo;s code behind file;</p>
<pre class="brush: c-sharp; toolbar: false">        [WebMethod]
        public static string CallMe() {

            return "You called me on " + DateTime.Now.ToString();

        }</pre>
<p>Compile your code and run your application with <strong>CTRL + F5</strong>. After that, open up the fiddler. We will make a HTTP POST request to our method with following code;</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb.png" width="644" height="404" /></a></p>
<p>Did you notice the URL? It gets the method name after the '/'. Amazing! After we execute this request, here is what we got;</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_3.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb_3.png" width="644" height="404" /></a></p>
<p>Pretty straight forward and simple. Now let&rsquo;s see how JQuery fits in here and as <a title="http://hanselman.com/" href="http://hanselman.com/" target="_blank">Scott Hanselman</a> always says let&rsquo;s see 'how Lego pieces fit in together'</p>
<p><strong>Getting a Little Work Done by Default</strong></p>
<p>Most of the heavy lifting is done by the following code here. I put this code inside a separate .js file named <strong>JQuery.PageMethod.Call.Helper.js</strong> but feel free to use it inside the your html markup.</p>
<pre class="brush: javascript; toolbar: false">//-----------------------------------------------------------------------------+
// jQuery call AJAX Page Method                                                |
//-----------------------------------------------------------------------------+
function PageMethod(fn, paramArray, successFn, errorFn) {
    var pagePath = window.location.pathname;
    //-------------------------------------------------------------------------+
    // Create list of parameters in the form:                                  |
    // {"paramName1":"paramValue1","paramName2":"paramValue2"}                 |
    //-------------------------------------------------------------------------+
    var paramList = '';
    if (paramArray.length &gt; 0) {
        for (var i = 0; i &lt; paramArray.length; i += 2) {
            if (paramList.length &gt; 0) paramList += ',';
            paramList += '"' + paramArray[i] + '":"' + paramArray[i + 1] + '"';
        }
    }
    paramList = '{' + paramList + '}';
    //Call the page method
    $.ajax({
        type: "POST",
        url: pagePath + "/" + fn,
        contentType: "application/json; charset=utf-8",
        data: paramList,
        dataType: "json",
        success: successFn,
        error: errorFn
    });
}</pre>
<blockquote>
<p>As I mentioned above, I have put this helper code inside a separate .js file. Although this is not a big file, I still would like to minify it so that I can get the minimum size file. I would rather minify that file with <a title="http://ajaxmin.codeplex.com/" href="http://ajaxmin.codeplex.com/" target="_blank">Ajax Minifier</a> via command line but there are a lot of good compression tools. My file was 1.364 bayt before the minify action and it became 291 bayt after that. I have 77 % saving here. Not bad!</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_4.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb_4.png" width="644" height="223" /></a></p>
</blockquote>
<p>When we look at this JQuery code here, we will see that nothing fancy going on. At the line 5, we are grabbing the current page path (which will be Default.aspx in our case). Between line 10 and 17, we are getting the paramArray parameter that we have declared and creating a JSON our of that. We are doing that to pass parameters to our method if our page method has parameters defined. Between line 19 and 26, we are making an HTTP Post request to the service which is a Page Method in our case.</p>
<p>Let&rsquo;s talk less and look at the code. Here is our entire code behind file;</p>
<pre class="brush: c-sharp; toolbar: false">using System;
using System.Web.Services;

namespace JQueryPageMethodCall1 {

    public partial class _Default : System.Web.UI.Page {

        protected void Page_Load(object sender, EventArgs e) {

        }

        [WebMethod]
        public static string CallMe() {

            return "You called me on " + DateTime.Now.ToString();
        }

        [WebMethod]
        public static string GetMeAGUID(string name, string surname, string age) {

            var poo = int.Parse(age);

            return string.Format(
                "Hey, {0} {1}. How is it goin over there? u are {2} years old and here is a Guid for you : {3}", 
                name, surname, poo.ToString(), Guid.NewGuid()
                );
        }

    }
}</pre>
<p>We have two web methods here. One is with parameters defined and the other is parameterless. We will consume those with the help of JQuery. Here is how our web form page looks like;</p>
<pre class="brush: xhtml; toolbar: false">&lt;%@ Page Title="Home Page" Language="C#" MasterPageFile="~/Site.master" AutoEventWireup="true" ClientIDMode="Static"
    CodeBehind="Default.aspx.cs" Inherits="JQueryPageMethodCall1._Default" %&gt;

&lt;asp:Content ID="HeaderContent" runat="server" ContentPlaceHolderID="HeadContent"&gt;
&lt;/asp:Content&gt;

&lt;asp:Content ID="ScriptContent" runat="server" ContentPlaceHolderID="ScriptPlaceHolder"&gt;
    &lt;script type="text/javascript" src="http://ajax.aspnetcdn.com/ajax/jQuery/jquery-1.5.2.min.js"&gt;&lt;/script&gt;
    &lt;script type="text/javascript" src="Scripts/Helpers/JQuery.PageMethod.Call.Helper.min.js"&gt;&lt;/script&gt;
    &lt;script type="text/javascript"&gt;

        $(document).ready(function () {

            var succeededAjaxFn = function (result) {

                $('#result').hide();
                $('&lt;p&gt;' + result.d + '&lt;/p&gt;').css({ background: 'green', padding: '10px', color: 'white' }).appendTo('#result');
                $('#result').fadeIn('slow');
            }

            var failedAjaxFn = function (result) {

                $('#result').hide();
                $('&lt;p&gt;Failed : ' + result.d + '&lt;/p&gt;').css({ background: 'red', padding: '10px', color: 'white' }).appendTo('#result');
                $('#result').fadeIn('slow');

            }

            $('#btnGetDateTime').click(function () {
                PageMethod("CallMe", [], succeededAjaxFn, failedAjaxFn);
            });

            $('#btnGetGUID').click(function () {
                var parameters = ["name", $('#txtName').val(), "surname", $('#txtSurname').val(), "age", $('#txtAge').val()];
                PageMethod("GetMeAGUID", parameters, succeededAjaxFn, failedAjaxFn);
            });
        });

    &lt;/script&gt;
&lt;/asp:Content&gt;

&lt;asp:Content ID="BodyContent" runat="server" ContentPlaceHolderID="MainContent"&gt;
    &lt;h2&gt;
        Get me a guid !
    &lt;/h2&gt;

    &lt;div&gt;
        Name : &lt;asp:TextBox runat="server" ID="txtName" /&gt;
    &lt;/div&gt;
    
    &lt;div&gt;
        Surname : &lt;asp:TextBox runat="server" ID="txtSurname" /&gt;
    &lt;/div&gt;

    &lt;div&gt;
        Age : &lt;asp:TextBox runat="server" ID="txtAge" /&gt;
    &lt;/div&gt;

    &lt;a href="#" id="btnGetGUID"&gt;Get Me A GUID&lt;/a&gt;
    &lt;a href="#" id="btnGetDateTime"&gt;GetCurrentDateTime&lt;/a&gt;

    &lt;div id="result" style="margin:20px 0 20px 0;" /&gt;

&lt;/asp:Content&gt;</pre>
<p>What we have done here is to define succeededAjaxFn function for success event and failedAjaxFn function for error event. Also we are calling the PageMethod function, that we have already defined and wrap it up in a separate file, in every a element click.</p>
<p>This is how our page looks like;</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_5.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb_5.png" width="644" height="418" /></a></p>
<p>Let&rsquo;s go through some scenarios. We&rsquo;ll click the <strong>GetCurrentDateTime </strong>link and see what happens;</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_6.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb_6.png" width="644" height="418" /></a></p>
<p>Boom! Worked. It gets what server gave to it. Now, we will fill in the fields and click the <strong>Get Me A GUID</strong> link;</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_7.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb_7.png" width="644" height="418" /></a></p>
<p>Rock on! This one also worked and we are not making any postbacks here. It is all JQuery. Let&rsquo;s send the request with empty <strong>Age</strong> value;</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_8.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb_8.png" width="644" height="418" /></a></p>
<p>It failed which means our JQuery code is working just fine. (It failed because I try to parse the age value into an integer. I supplied an empty age parameter so the method throws an exception)</p>
<p>Now, I will put a breakpoint on my <strong>GetMeAGUID</strong> method and debug my application. When I run the app, fill in the fields and click the <strong>Get Me A GUID</strong> link, it gets me to my server side code with JQuery code;</p>
<p><a href="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_9.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/image_thumb_9.png" width="644" height="425" /></a></p>
<p>You see, all the values we have provided are sitting right there. It is that easy. Of course, the main purpose of this demo is not so much fitting in the real world needs but you get the idea about what you could accomplish. Feel free to download the code from below and play around with it.</p>
<p>Be well and write good code <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/content/images/uploadedbyauthors/wlw/8668e28cfa76_148E4/wlEmoticon-smile.png" /></p>
<p><iframe title="Preview" scrolling="no" marginheight="0" marginwidth="0" frameborder="0" style="width: 98px; height: 115px; padding: 0; background-color: #fcfcfc;" src="http://cid-0ee89cb310fe3603.office.live.com/embedicon.aspx/Programming/JQueryPageMethodCall1.rar"></iframe></p>