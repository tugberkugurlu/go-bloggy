---
id: ea41a396-edc3-4bf9-9e73-058c1cba5763
title: Sending E-mail to All List of Membership Users with ASP.Net Using Built-in
  Membership API
abstract: Most of the Asp.Net developers are using Membership class of Asp.Net and
  in this blog post we will see how to send e-mail to all of the membership users
  at once...
created_at: 2010-11-26 09:23:28 +0000 UTC
tags:
- .NET
- ASP.Net
- C#
- eCommerce
- MS SQL
slugs:
- sending-e-mail-to-all-list-of-membership-users-with-asp-net-using-built-in-membership-api
---

<p>ASP.Net has a very powerfull built-in membership api and it has all the things we need;</p>
<ul>
<li>User Management</li>
<li>CRUD (Create, Read, Update &amp; Delete)&nbsp;of Membership Model</li>
<li>User approvement</li>
<li>Changing password</li>
<li>Recovering password</li>
</ul>
<p>And the other things. Also, it is safe to say so many developer use this magical thing.&nbsp;</p>
<p>The registered users are so important to our company and we want to let them know we live and our new things once in a while. As a result, newsletters' main purpose is that.</p>
<p>But how can we send e-mail to all of the registered users inside our magical Membership SQL Database? There are some ways but here we will see the simple and most useful way. We will be talking to database of course in order to get the e-mail addresses but we won't be directly talking to membership database. ASP.Net has already done that for us.&nbsp;</p>
<p>As <a target="_blank" title="ScottHa's Blog" href="http://www.hanselman.com/blog/">Scott Hanselman</a>&nbsp;says, talk is cheap, show me the code&nbsp;<img border="0" title="Smile" alt="Smile" src="../Scripts/tiny_mce/plugins/emotions/img/smiley-smile.gif" />&nbsp;Let's start with lesson - 101 : File &gt; New &gt; Project and select ASP.Net Web Application.</p>
<p><img src="https://www.tugberkugurlu.com/Content/Images/BlogUploadedPics/file-new-project-vs-2010-pro.png" alt="file-new-project-vs-2010-pro.png" title="file-new-project-vs-2010-pro.png" width="300" style="border: 0px initial initial;" /><img src="https://www.tugberkugurlu.com/Content/Images/BlogUploadedPics/new-web-application-cs-vs-2010-pro.PNG" alt="new-web-application-cs-vs-2010-pro.PNG" title="new-web-application-cs-vs-2010-pro.PNG" width="300" style="border: 0px initial initial;" /></p>
<p>After creating a new project, I will be doing this example on default page but feel free to do it anywhere you like to. Maybe you could create a class library project and that would be nicer given that you could use it in any other projects. So we will be applying the DRY (Don't repeat yourself) rule.</p>
<p>Inside&nbsp;Default.aspx.cs file, I created the below private class;</p>
<p>&nbsp;</p>
<pre class="brush: c-sharp">        private class Users {

            public string UserName { get; set; }
            public string EmailAddress { get; set; }
        
        } </pre>
<p>&nbsp;</p>
<p>Why did I create that? We will see why in a sec. Our main aim here is to send e-mail to our registered users and in order to do that we need all of the users, right? So, we need to talk to database for that. As mentioned before, we won't be using T-Sql for that. In fact, we won't even use linq queries.</p>
<p>Let's write the code of our private IList of <em>Users</em> (this <em>Users </em>class is the class that we just&nbsp;created above)</p>
<p>&nbsp;</p>
<p>&nbsp;</p>
<pre class="brush: c-sharp">        private IList<users> MembershipUserEmailAddresses() {

            IList<users> addresses = new List<users>();

            foreach (var item in Membership.GetAllUsers()) {

                MembershipUser user = (MembershipUser)item;

                addresses.Add(new Users { EmailAddress = user.Email.ToString(), UserName = user.UserName.ToString() });

            }

            return addresses;
        
        } 
        //The below code has been created by tinyMCE by mistake
        //</users></users></users>
</pre>
<p>&nbsp;</p>
<p>Now we have our data. Pretty easy. Ok, now is the time to write the code of our e-mail send method. Before doing that, make sure that you already configured your web.config file for your host settings. It is pretty straight forward.</p>
<p>&nbsp;</p>
<pre class="brush: xhtml">  <system.net>
    <mailsettings>
      <smtp from="info@example.com">
        <network host="mail.example.com" password="passwordOfyourEmailAddress" port="587" username="info@example.com">
      </network></smtp>
    </mailsettings>
  </system.net> </pre>
<p>After this little cınfiguration, it is time to write the e-mail send method. So, we have multiple users in hand. We have some options here;</p>
<ol>
<li>We could send the e-mail to every each user&nbsp;separately</li>
<li>We could send the e-mail to every each user at once by adding them as Bcc</li>
</ol>
<p>I will demonstrate the 2<sup>nd </sup>option here. Here is the acctual send method;</p>
<p>&nbsp;</p>
<pre class="brush: c-sharp">        private void Send(string Sender, string SenderDisplayName, IList<users> Recivers) {

            MailMessage MyMessage = new MailMessage();
            SmtpClient MySmtpClient = new SmtpClient();

            //Optional
            MyMessage.IsBodyHtml = true;


            //Required Fields
            MyMessage.From = new MailAddress(Sender, SenderDisplayName);
            MyMessage.Subject = "subject of the message";
            MyMessage.Body = "My Message... Here is my message";

            for (int i = 0; i &lt; Recivers.Count(); i++) {

                MyMessage.Bcc.Add(new MailAddress(Recivers[i].EmailAddress,Recivers[i].UserName));

            }

            MySmtpClient.Send(MyMessage);
        
        } 

        //The below code has been created by tinyMCE by mistake
        //</users>

</pre>
<p>&nbsp;</p>
<p>As you see, it is not so complicated. Of course, it could be much more pretty but I figured that we are demonstrating the function here, so I made a quick and simple one.</p>
<p>&nbsp;</p>
<p>What do we do now? Finally, it has come to fruition. I will fire it up on the page_load event but you could do it wherever you want;</p>
<p>&nbsp;</p>
<pre class="brush: c-sharp; highlight: [3];">        protected void Page_Load(object sender, EventArgs e) {

            Send("info@example.com", "My Name", MembershipUserEmailAddresses());

        } </pre>
<p>&nbsp;</p>
<p>That's it. It is very simple and&nbsp;straight forward.</p>
<p>You could do this by connecting to sql directly if you need any further fillitering but it is the way it works and I guess you get the idea.&nbsp;</p>
<p>Enjoy your coding&nbsp;<img border="0" title="Smile" alt="Smile" src="../Scripts/tiny_mce/plugins/emotions/img/smiley-smile.gif" /></p>
<p><strong>The whole code</strong></p>
<pre class="brush: c-sharp; collapse: true">using System;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Web.Security;
using System.Web.UI;
using System.Web.UI.WebControls;
using System.Net.Mail;

namespace MembershipUsersSendEmail
{
    public partial class _Default : System.Web.UI.Page
    {
        protected void Page_Load(object sender, EventArgs e) {

            Send("info@example.com", "My Name", MembershipUserEmailAddresses());

        }

        private class Users {

            public string UserName { get; set; }
            public string EmailAddress { get; set; }
        
        }

        private IList<users> MembershipUserEmailAddresses() {

            IList<users> addresses = new List<users>();

            foreach (var item in Membership.GetAllUsers()) {

                MembershipUser user = (MembershipUser)item;

                addresses.Add(new Users { EmailAddress = user.Email.ToString(), UserName = user.UserName.ToString() });

            }

            return addresses;
        
        }

        private void Send(string Sender, string SenderDisplayName, IList<users> Recivers) {

            MailMessage MyMessage = new MailMessage();
            SmtpClient MySmtpClient = new SmtpClient();

            //Optional
            MyMessage.IsBodyHtml = true;


            //Required Fields
            MyMessage.From = new MailAddress(Sender, SenderDisplayName);
            MyMessage.Subject = "subject of the message";
            MyMessage.Body = "My Message... Here is my message";

            for (int i = 0; i &lt; Recivers.Count(); i++) {

                MyMessage.Bcc.Add(new MailAddress(Recivers[i].EmailAddress,Recivers[i].UserName));

            }

            MySmtpClient.Send(MyMessage);
        
        }
    }
} 

//The below code has been created by tinyMCE by mistake
//</users></users></users></users>

</pre>