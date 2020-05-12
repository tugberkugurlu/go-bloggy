---
id: f5063f89-4b83-49dc-bfc2-143a3db0259e
title: SQL Injection vs. Lethal Injection / Protection Against SQL Injection
abstract: SQL Injection and Lethal Injection... They are both dangerous and they can
  be easily fatal. But how? What is SQL Injection and how it can effect my project?
  The answers are in this blog post.
created_at: 2011-01-29 18:37:00 +0000 UTC
tags:
- .net
- ASP.Net
- MS SQL
- Security
- SQL Injection
slugs:
- sql-injection-vs-lethal-injection-protection-against-sql-injection
---

<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/lethal-injection.jpg" target="_blank"><img style="background-image: none; margin: 0px 15px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="lethal-injection" border="0" alt="lethal-injection" align="left" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/lethal-injection_thumb.jpg" width="244" height="163" /></a>Writing a software, web application code is a real deal. It requires a good quality of talent on programing languages, appropriate analectic approach and most of all, a good prescience on your project. The things I have mentioned are very important and basic features of a qualified programmer.</p>
<p>I am not a student of a computer related course and I haven&rsquo;t been but I support that educational background on computer science makes a difference on the quality of programmer. But the diploma or the course certificate is not enough. Little mistakes could be unforgivable in programming world and your diploma or certificate cannot get those mistakes back or cover them.</p>
<p>As for our topic, SQL injection is one of the most important topic on programming security. I have seen couple of developer&rsquo;s <strong><em>&ldquo;handy&rdquo;</em> </strong>work for last several months and I decided to write this blog post and I would like say all of the developers, with <strong>no offense;</strong></p>
<blockquote>
<p><strong>Please, </strong>if you are creating a project with database structure, for the sake of programming, <strong>be aware of the SQL injection and its effects</strong>. It is not a shame that you haven&rsquo;t heard of that term. What the shame is to write lines of codes creating the proper connection with your database without considering the effects of <strong>SQL injection !</strong></p>
<p><strong>NO OFEENSE !</strong></p>
</blockquote>
<p><strong>What is SQL Injection?</strong></p>
<p>Well, some of you might want to know what the SQL injection is. I won&rsquo;t explore the world from scratch, so here is the clear explanation that I quoted from <a title="http://en.wikipedia.org" href="http://en.wikipedia.org" target="_blank">Wikipedia</a>;</p>
<blockquote>
<p><em><strong>SQL injection</strong> is a code injection technique that <strong>exploits a security vulnerability occurring in the database layer of an application</strong>. The vulnerability is present when user input is either incorrectly filtered for string literal escape characters embedded in SQL statements or user input is not strongly typed and thereby unexpectedly executed. It is an instance of a more general class of vulnerabilities that can occur whenever one programming or scripting language is embedded inside another. <strong>SQL injection attacks</strong> are also known as <strong>SQL insertion attacks</strong>.</em></p>
</blockquote>
<p><em>&nbsp;</em></p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection.jpg"><img style="background-image: none; margin: 0px 15px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="sql-injection" border="0" alt="sql-injection" align="left" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection_thumb.jpg" width="239" height="244" /></a></p>
<p>So the definition supposed to clear the question marks but it might not. Let&rsquo;s demonstrate.</p>
<p>Imagine that you have a web application running on the web and it aims to provide an user interface to your customers to view their account details.</p>
<p>The demo application is pretty small and we will only create 2 pages with one database table. Page 1 will be the wrong scenario&nbsp;and 2nd one will be the right.</p>
<p>In this application, we will see how the end user can easily display the sensetive data you migh have in your database.</p>
<p><em>"I would like to say this, in a nutshell, nobody (I mean a programmer who knows what he/she is doing) developed a kind of application for that kind of purpose but to demonstrate the topic, I have done something like that. The project is not supposed to be a real world example."</em></p>
<p>Our database structure looks like this;</p>
<p>&nbsp;</p>
<p>&nbsp;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-server-explorer-view-for-sql-database-structure.png"><img style="background-image: none; margin: 0px 15px 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; float: left; padding-top: 0px; border-width: 0px;" title="sql-injection-demo-project-server-explorer-view-for-sql-database-structure" border="0" alt="sql-injection-demo-project-server-explorer-view-for-sql-database-structure" align="left" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-server-explorer-view-for-sql-database-structure_thumb.png" width="169" height="244" /></a></p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-data-view-for-sql-database-structure.png"><img style="background-image: none; margin: 0px 0px 10px 10px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="sql-injection-demo-project-data-view-for-sql-database-structure" border="0" alt="sql-injection-demo-project-data-view-for-sql-database-structure" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-data-view-for-sql-database-structure_thumb.png" width="581" height="242" /></a></p>
<p>I won&rsquo;t dive into details, I will post the project code so your could download and dig it letter.</p>
<p><strong>SQL Injectable Page</strong>&nbsp;</p>
<p>I have used <em>GridView</em> to list the data and here is what the user page looks like;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-open-page-view.png"><img style="background-image: none; margin: 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="sql-injection-demo-project-sql-injection-open-page-view" border="0" alt="sql-injection-demo-project-sql-injection-open-page-view" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-open-page-view_thumb.png" width="644" height="188" /></a></p>
<p>The code has been use to provide the data is as exactly below;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [10]">        protected void butn_click(object sender, EventArgs e) {

            GridView1.DataSource = DataProvider(txt1);
            GridView1.DataBind();
        }

        private static DataSet DataProvider(TextBox mytext) {

            string connectionString = WebConfigurationManager.ConnectionStrings["SampleConnectionString"].ConnectionString;
            string sql = "SELECT * FROM Customers WHERE ([TCKimlikNo] = '" + mytext.Text + "')";
            SqlConnection con = new SqlConnection(connectionString);
            SqlCommand cmd = new SqlCommand(sql, con);
            cmd.CommandType = CommandType.Text;
            SqlDataAdapter MyAdapter = new SqlDataAdapter(cmd);
            DataSet ds = new DataSet("MyDs");
            MyAdapter.Fill(ds, "MyDs");

            return ds;
        }</pre>
<p><em>&ldquo;DataProvider()&rdquo; </em>static method connects to the database and executes some SQL against a SQL Server database that returns the number of rows where the user data supplied by the user matches a row in the database. If the result is one matching row, that row will be displayed as you can see;</p>
<p>&nbsp;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-open-page-view-no-harm.png"><img style="background-image: none; margin: 10px 0px; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="sql-injection-demo-project-sql-injection-open-page-view-no-harm" border="0" alt="sql-injection-demo-project-sql-injection-open-page-view-no-harm" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-open-page-view-no-harm_thumb.png" width="682" height="161" /></a></p>
<p>Let&rsquo;s put a break point on the 10th line and hit it;</p>
<p>&nbsp;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-breakpoint-debug.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="sql-injection-demo-project-sql-injection-breakpoint-debug" border="0" alt="sql-injection-demo-project-sql-injection-breakpoint-debug" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-breakpoint-debug_thumb.png" width="681" height="272" /></a></p>
<p>&nbsp;</p>
<p>The value supplied above for TCKimlikNo is 34265128731. As we can see in the image, the code works perfectly fine and the value is on the place that we wanted. Now let&rsquo;s do some evil things;</p>
<p>&nbsp;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-breakpoint-debug-evil.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="sql-injection-demo-project-sql-injection-breakpoint-debug-evil" border="0" alt="sql-injection-demo-project-sql-injection-breakpoint-debug-evil" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/sql-injection-demo-project-sql-injection-breakpoint-debug-evil_thumb.png" width="678" height="341" /></a></p>
<p>&nbsp;</p>
<p>Now the query explains itself pretty clearly. The <strong>evil user </strong>put this;</p>
<p><em>hi&rsquo; or &lsquo;1&rsquo; = &lsquo;1</em></p>
<p>And the logic fits. Method will return all the rows inside the database table. Look at the result;</p>
<p>&nbsp;</p>
<p><a href="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/image.png"><img style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border-width: 0px;" title="image" border="0" alt="image" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/image_thumb.png" width="679" height="233" /></a></p>
<p>&nbsp;</p>
<p>Boom, you have been hacked ! This is the SQL Injection my friends. This thing is easy to apply and the worse part, this mistake is being made often.</p>
<p>Here is a quote from Mike&rsquo;s blog;</p>
<blockquote>
<p><em>This is SQL Injection. Basically, additional SQL syntax has been injected into the statement to change its behavior. The single quotes are string delimiters as far as T-SQL is concerned, and if you simply allow users to enter these without managing them, you are asking for potential trouble. </em></p>
</blockquote>
<p><strong>What is the Prevention?</strong></p>
<p>Easy ! Just do not create the world from scratch.</p>
<p>If you are a ASP.Net user, use parameters instead of hand made code. Review the following code and compare it with the previous one;</p>
<pre class="brush: c-sharp; toolbar: false; highlight: [2,3,4]">            string connectionString = WebConfigurationManager.ConnectionStrings["SampleConnectionString"].ConnectionString;
            string sql = "SELECT * FROM Customers WHERE ([TCKimlikNo] = @IDParameter)";
            SqlConnection con = new SqlConnection(connectionString);
            SqlCommand cmd = new SqlCommand(sql, con);
            cmd.CommandType = CommandType.Text;
            cmd.Parameters.Add("@IDParameter", SqlDbType.VarChar);
            cmd.Parameters["@IDParameter"].Value = mytext.Text;

            SqlDataAdapter MyAdapter = new SqlDataAdapter(cmd);
            DataSet ds = new DataSet("MyDs");
            MyAdapter.Fill(ds, "MyDs");

            return ds;</pre>
<p>Done ! You will be good to go <img style="border-style: none;" class="wlEmoticon wlEmoticon-smile" alt="Smile" src="http://tugberkugurlu.com/Content/Images/UploadedByAuthors/wlw/SQL-Injection-vs.-Lethal-Injection_D85/wlEmoticon-smile.png" /> But what happened there? Here is a good quote from Mike&rsquo;s blog again;</p>
<blockquote>
<p><strong>Parameter Queries</strong> <br /> <br />Parameters in queries are placeholders for values that are supplied to a SQL query at runtime, in very much the same way as parameters act as placeholders for values supplied to a C# method at runtime. And, just as C# parameters ensure type safety, SQL parameters do a similar thing. If you attempt to pass in a value that cannot be implicitly converted to a numeric where the database field expects one, exceptions are thrown</p>
</blockquote>
<p>Paramaters will protect your data if your building your project this way but another safe way is to use <strong>LINQ to SQL</strong> and <strong>Entity Framework</strong> to protect your project against SQL Injection.</p>
<p><iframe src="http://cid-0ee89cb310fe3603.office.live.com/embedicon.aspx/Programming/SQLInjection.rar" style="width: 98px; height: 115px; padding: 0; background-color: #fcfcfc;" frameborder="0" marginwidth="0" marginheight="0" scrolling="no" title="Preview"></iframe></p>