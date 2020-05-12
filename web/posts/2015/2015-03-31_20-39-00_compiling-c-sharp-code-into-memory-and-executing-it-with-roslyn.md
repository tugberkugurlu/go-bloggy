---
id: fd385d9d-3eab-4d3e-87f0-4510b1334e86
title: Compiling C# Code Into Memory and Executing It with Roslyn
abstract: Let me show you how to compile a piece of C# code into memory and execute
  it with Roslyn. It is super easy if you believe it or not :)
created_at: 2015-03-31 20:39:00 +0000 UTC
tags:
- .net
- C#
- Roslyn
slugs:
- compiling-c-sharp-code-into-memory-and-executing-it-with-roslyn
- compile-a-piece-of-c-sharp-code-into-memory-and-execute-it-with-roslyn
- compiling-a-piece-of-c-sharp-code-into-memory-and-execute-it-with-roslyn
- compiling-a-piece-of-c-sharp-code-into-memory-and-executing-it-with-roslyn
---

<p>For the last couple of days, I have been looking into how to get Razor view engine running outside <a href="https://github.com/aspnet/Mvc">ASP.NET 5 MVC</a>. It was fairly straight forward but there are a few bits and pieces that you need to stitch together which can be challenging. I will get Razor part in a later post and in this post, I would like to show how to compile a piece of C# code into memory and execute it with Roslyn, which was one of the parts of getting Razor to work outside ASP.NET MVC.</p> <p>First thing is to install C# code analysis library into you project though NuGet. In other words, installing Roslyn :)</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>Install<span style="color: gray">-</span>Package Microsoft.CodeAnalysis.CSharp <span style="color: gray">-</span>pre</pre></div></div>
<p>This will pull down bunch of stuff like <a href="http://www.nuget.org/packages/Microsoft.CodeAnalysis.Analyzers">Microsoft.CodeAnalysis.Analyzers</a>, <a href="http://www.nuget.org/packages/System.Collections.Immutable">System.Collections.Immutable</a>, etc. as its dependencies which is OK. In order to compile the code, we want to first create a <a href="http://source.roslyn.codeplex.com/#Microsoft.CodeAnalysis/Syntax/SyntaxTree.cs,8649488200d5b57a">SyntaxTree</a> instance. We can do this pretty easily by parsing the code block using the <a href="http://source.roslyn.codeplex.com/#Microsoft.CodeAnalysis.CSharp/Syntax/CSharpSyntaxTree.cs,bade6a931ef27795">CSharpSyntaxTree.ParseText</a> static method.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>SyntaxTree syntaxTree = CSharpSyntaxTree.ParseText(<span style="color: #a31515">@"
    using System;

    namespace RoslynCompileSample
    {
        public class Writer
        {
            public void Write(string message)
            {
                Console.WriteLine(message);
            }
        }
    }"</span>);</pre></div></div>
<p>The next step is to create a <a href="http://source.roslyn.codeplex.com/#Microsoft.CodeAnalysis/Compilation/Compilation.cs,ec43f5a2c70b26f1">Compilation</a> object. If you wonder, the compilation object is an immutable representation of a single invocation of the compiler (code comments to the rescue). It is the actual bit which carries the information about syntax trees, reference assemblies and other important stuff which you would usually give as information to the compiler. We can create an instance of a Compilation object through another static method: <a href="http://source.roslyn.codeplex.com/#Microsoft.CodeAnalysis.CSharp/Compilation/CSharpCompilation.cs,cb0be8b9d3027ce8">CSharpCompilation.Create</a>.</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">string</span> assemblyName = Path.GetRandomFileName();
MetadataReference[] references = <span style="color: blue">new</span> MetadataReference[]
{
    MetadataReference.CreateFromFile(<span style="color: blue">typeof</span>(<span style="color: blue">object</span>).Assembly.Location),
    MetadataReference.CreateFromFile(<span style="color: blue">typeof</span>(Enumerable).Assembly.Location)
};

CSharpCompilation compilation = CSharpCompilation.Create(
    assemblyName,
    syntaxTrees: <span style="color: blue">new</span>[] { syntaxTree },
    references: references,
    options: <span style="color: blue">new</span> CSharpCompilationOptions(OutputKind.DynamicallyLinkedLibrary));</pre></div></div>
<p>Hard part is now done. The final bit is actually running the compilation and getting the output (in our case, it is a dynamically linked library). To run the actual compilation, we will use the <a href="http://source.roslyn.codeplex.com/#Microsoft.CodeAnalysis/Compilation/Compilation.cs,9f62285c857030a3">Emit</a> method on the Compilation object. There are a few overloads of this method but we will use the one where we can pass a <a href="https://msdn.microsoft.com/en-us/library/system.io.stream%28v=vs.110%29.aspx">Stream</a> object in and make the Emit method write the assembly bytes into it. Emit method will give us an instance of an <a href="http://source.roslyn.codeplex.com/#Microsoft.CodeAnalysis/Compilation/EmitResult.cs,19d1f60577d83c3c">EmitResult</a> object and we can pull the status of the compilation, warnings, failures, etc. from it. Here is the actual code:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: blue">using</span> (<span style="color: blue">var</span> ms = <span style="color: blue">new</span> MemoryStream())
{
    EmitResult result = compilation.Emit(ms);

    <span style="color: blue">if</span> (!result.Success)
    {
        IEnumerable&lt;Diagnostic&gt; failures = result.Diagnostics.Where(diagnostic =&gt; 
            diagnostic.IsWarningAsError || 
            diagnostic.Severity == DiagnosticSeverity.Error);

        <span style="color: blue">foreach</span> (Diagnostic diagnostic <span style="color: blue">in</span> failures)
        {
            Console.Error.WriteLine(<span style="color: #a31515">"{0}: {1}"</span>, diagnostic.Id, diagnostic.GetMessage());
        }
    }
    <span style="color: blue">else</span>
    {
        ms.Seek(0, SeekOrigin.Begin);
        Assembly assembly = Assembly.Load(ms.ToArray());
    }
}</pre></div></div>
<p>As mentioned before, here, we are getting the EmitResult out as a result and looking for its status. If it’s not a success, we get the errors out and output them. If it’s a success, we load the bytes into an <a href="https://msdn.microsoft.com/en-us/library/system.reflection.assembly%28v=vs.110%29.aspx">Assembly</a> object. The Assembly object you have here is no different the ones that you are used to. From this point on, it’s all up to your ninja reflection skills in order to execute the compiled code. For the purpose of this demo, it was as easy as the below code:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>Type type = assembly.GetType(<span style="color: #a31515">"RoslynCompileSample.Writer"</span>);
<span style="color: blue">object</span> obj = Activator.CreateInstance(type);
type.InvokeMember(<span style="color: #a31515">"Write"</span>,
    BindingFlags.Default | BindingFlags.InvokeMethod,
    <span style="color: blue">null</span>,
    obj,
    <span style="color: blue">new</span> <span style="color: blue">object</span>[] { <span style="color: #a31515">"Hello World"</span> });</pre></div></div>
<p>This was in a console application and after running the whole thing, I got the expected result:</p>
<p><a href="https://tugberkugurlu.blob.core.windows.net/bloggyimages/d03de59f-8806-4942-83b3-d2b9bb93a2e2.png"><img title="image" style="border-top: 0px; border-right: 0px; background-image: none; border-bottom: 0px; padding-top: 0px; padding-left: 0px; border-left: 0px; display: inline; padding-right: 0px" border="0" alt="image" src="https://tugberkugurlu.blob.core.windows.net/bloggyimages/a275313d-31ae-4048-b810-3e4fdf9ea334.png" width="644" height="135"></a></p>
<p>Pretty sweet and easy! <a href="https://github.com/tugberkugurlu/DotNetSamples/tree/0883fb2e8c723420663e2d60140ce7591c7b311a/csharp/RoslynCompileSample">This sample is up on GitHub</a> if you are interested.</p>  