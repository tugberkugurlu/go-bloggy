---
id: f5e6d03e-5bc9-4e15-aa2c-ae6e3502b797
title: How To Use C# (C Sharp) Switch Case, Switch Case Samples
abstract: After you read this article, you will be able to use the 'Switch Case' Function
  on your C# Projects. This function becomes so handy with DropDownList & RadioButtonList
  !
created_at: 2010-04-16 08:51:47 +0000 UTC
tags:
- .net
- ASP.Net
- C#
slugs:
- how-to-use-c-sharp-c-sharp-switch-case-switch-case-samples
---

<p>C# Switch Case function is the same with Select Case function in Visual Basic. It is usually used along with DropDownLists or RadioButtonLists. We wiil demostrate a&nbsp;function here with s&nbsp;dropdownlist.<br /> <br /> Switch Case function enables us to run what we need to run in a case that we want. After the demonstration, you will exactly see what it means !<br /> <br /> <span style="color: #009900;">Code Sample<br /> </span></p>
<hr />
<p>&nbsp;</p>
<p><span style="font-weight: bold;">1-</span> Firstly, put the below code into the body section of your page;<br /> <br /> <span color="#0000ff"><span color="#0000ff"><span color="#0000FF" style="color: #0000ff;">&nbsp;</span> </span></span></p>
<pre class="brush: xhtml"><asp:dropdownlist id="DropDownList1" autopostback="true" runat="server" onselectedindexchanged="DropDownList1_SelectedIndexChanged">

<asp:listitem text="1" value="1">
<asp:listitem text="2" value="2">
<asp:listitem text="3" value="3">

</asp:listitem></asp:listitem></asp:listitem></asp:dropdownlist>

<br /><br />

<asp:label id="Label1" runat="server"></asp:label> </pre>
<p><br /> <br /><span style="font-weight: bold;">2-</span> And then, As you can see, we appointed <span style="color: #0000ff;" color="#0000ff"><span style="font-style: italic; font-size: 7pt;">'DropDownList1_SelectedIndexChanged'</span></span><span style="font-size: 8pt;">&nbsp;</span>to run on DropDownList's SelectedIndexChanged. So on our code page we will use the below code;<br /> <br /> <br /> <span color="#0000FF" style="color: #0000ff;"> </span></p>
<pre class="brush: c-sharp">protected void DropDownList1_SelectedIndexChanged(object sender, EventArgs e) {
          switch(DropDownList1.SelectedValue) {

                  case "1":
                          Label1.Text = "DropDownList Value = 1";
                          break;
                  case "2":
                          Label1.Text = "DropDownList Value = 2";
                          break;
                  default:
                          Label1.Text = "None Of Them";
                          break;
          }

  } </pre>
<p><br /> <br /> <br /> That's It ! Of course, we can give so many examples on switch case function but this is pretty much it is ! <br /> <br /> On our example when we select the item which has the value '1' on our dropdownlist, Label1's Text will be <span style="color: #a31515;">"DropDownList Value = 1" </span>and on value '2' it will be <span style="color: #a31515;">"DropDownList Value = 2" <br /> <br /> </span>As you can see we added the below code at the end of our switch case function;<br /> <br /> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <span style="color: blue;">default</span>:<br /> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Label1.Text = <span style="color: #a31515;">"None Of Them"</span>;<br /> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <span style="color: blue;">break</span>;<br /> <br /> It reffers that when an item which&nbsp;is not indicated as a case is selected, the fuction will run this value !<br /> <br /> I hope it gives you an idea on C# (C Sharp) Switch Case Function !</p>
<p>&nbsp;</p>
<p><strong>Update on&nbsp;2010-12-01&nbsp;12:53:57</strong></p>
<p>In order to use switch case with Enum type in C#, follow the following instructions;</p>
<p>&nbsp;</p>
<pre class="brush: c-sharp">    public enum MyEnum {

        enum1 = 1,
        enum2 = 2,
        enum3 = 3

    }

        public static string EnumSwitchCaseTry(MyEnum myenum) {

            switch (myenum) 
            {
                case MyEnum.enum1:
                    //Do whatever you need to do

                case MyEnum.enum2:
                    //Do whatever you need to do

                case MyEnum.enum3:
                    //Do whatever you need to do

                default:
                    break;
            }

            //.....

            //.....

        } </pre>
<p>&nbsp;</p>