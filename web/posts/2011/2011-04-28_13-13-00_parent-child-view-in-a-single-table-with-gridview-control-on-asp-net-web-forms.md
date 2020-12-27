---
id: f42cbd1f-0850-443e-886b-e2bfcb5ee880
title: Parent / Child View In a Single Table With GridView Control On ASP.NET Web
  Forms
abstract: This awesome blog post will demonstrate how to create a complete, sub-grouped
  product list in a single grid. Get ready for the awesomeness...
created_at: 2011-04-28 13:13:00 +0000 UTC
tags:
- .NET
- ASP.Net
- C#
- MS SQL
slugs:
- parent-child-view-in-a-single-table-with-gridview-control-on-asp-net-web-forms
---

<p><a href="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/2e5515b18bdc_E211/parent-child-gridview.png"><img height="150" width="150" src="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/2e5515b18bdc_E211/parent-child-gridview_thumb.png" align="right" alt="parent-child-gridview" border="0" title="parent-child-gridview" style="background-image: none; margin: 0px 0px 10px 30px; padding-left: 0px; padding-right: 0px; display: inline; float: right; padding-top: 0px; border: 0px;" /></a>Sometimes we want to create a parent/child report that shows all the records from the child table, organized by parent. I have been thinking about this thing in my head lately but haven&rsquo;t been sure for long time about how to get it done with web forms.</p>
<p>This implementation is actually so easy with ASP.NET MVC framework.The model on the view has its own mappings to relational tables on the database (assuming that we are using an ORM such as Entity Framework) and a foreach loop will do the trick. So, how is this thing done with ASP.NET web forms? Let&rsquo;s demonstrate a sample.</p>
<p>I have created a new ASP.NET Web Application under .Net Framework 4 and also my database under this project. <em>(I&rsquo;m using Visual Studio 2010 as IDE but feel free to use Visual Web Developer Express 2010)</em> After that I created my <strong>ADO.NET Entity Data Model</strong> with database first approach. Our model should look like as following;</p>
<p><a href="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/2e5515b18bdc_E211/image.png"><img height="408" width="644" src="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/2e5515b18bdc_E211/image_thumb.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>After creating our database structure and ORM <em>(and filling our database tables with some fake values for demonstration purpose)</em>, now we&rsquo;ll be playing with ASP.NET GridView control on a web form page.</p>
<p>The fundamental idea is to create a GridView control for the parent table <em>(this is Category class in our case)</em> that contains an embedded GridView for each row. There child GridView controls are added into the parent GridView using a <strong>TemplateField</strong>.</p>
<p>But the hard part is that <strong>you cannot bind the child GirdView controls at the same time that you bind the parent GirdView </strong>because the parent GirdView rows haven&rsquo;t been created yet. So, we need to wait for <strong>GirdView.DataBound</strong> event to fire in the parent view before binding the child GridView controls.</p>
<p>In our example, the parent grid view defines two columns and they are both the TemplateField type. The first column combines the category name and category description as you can see below;</p>
<pre class="brush: xhtml; toolbar: false">          &lt;asp:TemplateField HeaderText="Category"&gt;

                &lt;ItemStyle VerticalAlign="Top" Width="20%" /&gt;
                &lt;ItemTemplate&gt;
                
                    &lt;br /&gt;
                    &lt;b&gt;&lt;%#Eval("CategoryName")%&gt;&lt;/b&gt;
                    &lt;br /&gt;&lt;br /&gt;
                    &lt;%#Eval("CategoryDescription")%&gt;
                    &lt;br /&gt;

                &lt;/ItemTemplate&gt;

            &lt;/asp:TemplateField&gt;</pre>
<p>The second column contains an embedded GridView of products, with two bound columns as you can see below;</p>
<pre class="brush: xhtml; toolbar: false">            &lt;asp:TemplateField HeaderText="Products"&gt;
            
                &lt;ItemStyle VerticalAlign="Top" Width="80%" /&gt;
                &lt;ItemTemplate&gt;
                
                    &lt;asp:GridView ID="productsGrid" runat="server" AutoGenerateColumns="false"&gt;
                        &lt;Columns&gt;
                        
                            &lt;asp:BoundField DataField="ProductName" HeaderText="Product Name" /&gt;
                            &lt;asp:BoundField DataField="Price" HeaderText="Unit Price" DataFormatString="{0:C}" /&gt;

                        &lt;/Columns&gt;
                    &lt;/asp:GridView&gt;

                &lt;/ItemTemplate&gt;

            &lt;/asp:TemplateField&gt;</pre>
<p>You probably realized that markup for the second GirdView does not set the DataSourceID property. That's because the data source for each of these grids will be supplied programmatically as the parent grid is being bound to its data source.</p>
<p>Now we need to create two data sources, one for retrieving the list of categories and the other for retrieving all products in a specified category. As we have our model as ADO.NET Entity Data Model, we will use <a target="_blank" href="http://msdn.microsoft.com/en-us/library/system.web.ui.webcontrols.entitydatasource.aspx" title="http://msdn.microsoft.com/en-us/library/system.web.ui.webcontrols.entitydatasource.aspx">EntityDataSoruce</a> to communicate with the database. The following code for first data source which will fill the parent GirdView;</p>
<pre class="brush: xhtml; toolbar: false">    &lt;asp:EntityDataSource ID="EntityDataSource1" runat="server" 
        ConnectionString="name=ProductsEntities" 
        DefaultContainerName="ProductsEntities" EnableFlattening="False" 
        EntitySetName="Categories"&gt;
    &lt;/asp:EntityDataSource&gt;</pre>
<p>Now, you need to bind the first grid directly to the data source and your markup for the grid view beginning tag should look like this;</p>
<pre class="brush: xhtml; toolbar: false; highlight: [2]">    &lt;asp:GridView ID="categoryGrid" AutoGenerateColumns="false" DataKeyNames="CategoryID"
        DataSourceID="EntityDataSource1" 
        onrowdatabound="categoryGrid_RowDataBound" runat="server" Width="100%"&gt;</pre>
<p>And here we are on the tricky part; binding the child GirdView controls. First, we need a second EntityDataSource. The second data source contains the query that&rsquo;s called multiple times to fill the child GridView. Each time, it retrieves the products that are in a different category. The CategoryID is supplied as a parameter;</p>
<pre class="brush: xhtml; toolbar: false">    &lt;asp:EntityDataSource ID="EntityDataSource2" runat="server" 
        ConnectionString="name=ProductsEntities" 
        DefaultContainerName="ProductsEntities" EnableFlattening="False" Where="it.CategoryID = @categoryid"
        EntitySetName="Products"&gt;
        &lt;WhereParameters&gt;
            &lt;asp:Parameter Name="categoryid" Type="Int32" /&gt;
        &lt;/WhereParameters&gt;
    &lt;/asp:EntityDataSource&gt;</pre>
<p>To bind the child GridView controls, you need to react to the GridView.RowDataBound event, which fires every time a row is generated and bound to the parent GridView. At this point, you can retrieve the child GridView control from the second column and bind it to the product information by programmatically. To ensure that you show only the products in the current category, you must also retrieve the CategoryID field for the current item and pass it as a parameter. Here&rsquo;s the code you need;</p>
<pre class="brush: c-sharp; toolbar: false">        protected void categoryGrid_RowDataBound(object sender, GridViewRowEventArgs e) {

            if (e.Row.RowType == DataControlRowType.DataRow) {

                //get the GridView control in the second column
                GridView gridChild = (GridView)e.Row.Cells[1].Controls[1];

                //set the categoryid parameter so you get the products in the current category only

                string categoryID = categoryGrid.DataKeys[e.Row.DataItemIndex].Value.ToString();
                EntityDataSource2.WhereParameters[0].DefaultValue = categoryID;

                //Bind the grid
                gridChild.DataSource = EntityDataSource2;
                gridChild.DataBind();

            }

        }</pre>
<p>Let&rsquo;s fire up our project and see what happens;</p>
<p><a href="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/2e5515b18bdc_E211/image_3.png"><img height="396" width="644" src="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/2e5515b18bdc_E211/image_thumb_3.png" alt="image" border="0" title="image" style="background-image: none; padding-left: 0px; padding-right: 0px; display: inline; padding-top: 0px; border: 0px;" /></a></p>
<p>We totally nailed it <em>(I added some style to make it look a little bit better)</em>. Perfect.</p>
<p>I hope that you found it useful and it helped <img src="https://www.tugberkugurlu.com/content/images/uploadedbyauthors/wlw/2e5515b18bdc_E211/wlEmoticon-smile.png" alt="Smile" class="wlEmoticon wlEmoticon-smile" style="border-style: none;" /></p>
<p><iframe title="Preview" scrolling="no" marginheight="0" marginwidth="0" frameborder="0" style="width: 98px; height: 115px; padding: 0; background-color: #fcfcfc;" src="http://cid-0ee89cb310fe3603.office.live.com/embedicon.aspx/Programming/SubgroupedProducts.rar"></iframe></p>