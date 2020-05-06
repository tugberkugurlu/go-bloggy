---
title: Is There Anybody Out There Who Doesn't Have Gravatar !
abstract: 'Sometimes we get sick to put our avatar pic. to every web site that we
  have registered to view our avatar. In this point gravatar.com helps us. '
created_at: 2010-11-19 22:14:47 +0000 UTC
tags:
- ASP.Net
- IT Stuff
- NuGet
- Time Saviour
- Tips
slugs:
- is-there-anybody-out-there-who-doesn-t-have-gravatar
---

<p>Sometimes we get sick to put our avatar pic. to every web site that we have registered to view our avatar. In this point&nbsp;<a target="_self" href="http://gravatar.com/">gravatar.com</a>&nbsp;helps us.&nbsp;<br /> <br /> Gravatar works very basicly and securely. Firstly, you need to go to&nbsp;<a target="_self" href="http://gravatar.com/">gravatar.com</a> and sing up. After uploading your avatar and configure it with your e-mail address, your are basically done.<br /> <br /> When you are writing a comment to a blog post, you are mostly required to enter your e-mail address and that e-mail address is enough to put your global avatar along with your comment. Of course, this is valid if the web site accepts gravatar pictures.&nbsp;<br /> <br /> Do not panic whether your e-mail address will be on picture's url. It won't. The picture url will similarly looks like that;<br /> <br /> <a target="_self" href="http://www.gravatar.com/avatar.php?gravatar_id=17698e3ad0e0dc70853cddda166bc573&amp;size=80&amp;default=identicon">http://www.gravatar.com/avatar.php?gravatar_id=17698e3ad0e0dc70853cddda166bc573&amp;size=80&amp;default=identicon</a><br /> <br /> You could easlily take advantages of gravatar and use it on your website. You could read <a target="_self" href="http://gravatar.com/site/implement/">implementation guidance</a> for more information.<br /> <br /> If you are using ASP.Net MVC and you are familiar with NuGet, the implementation will be so easy.<br /> <br /> <span style="font-weight: bold;">Implementation of Gravatar on ASP.Net MVC</span></p>
<ol>
<li>Write this command on Package Management Console : PM&gt; install-package Microsoft.MVC.Helpers</li>
<li>After the instalation, you will see Microsoft.MVC.Helpers.Gravatar.GetHTML(...) method. That method will generate the gravatar automatically.</li>
</ol>
<p>Hope this helps !</p>