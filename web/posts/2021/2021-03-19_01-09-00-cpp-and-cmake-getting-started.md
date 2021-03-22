---
id: 01F140FYA2CBEG70R4YX6XXNXG
title: "C++, Getting Started with the Basics: Hello World and the Build Pipeline"
abstract: I am learning C++, and what better way to make the learning stronger than blogging about my journey and pinning down my experience. I especially think this will be really beneficial when it comes to C++, since the barrier to entry is quite high and there is too much to learn. So, the reason that this post exists is a bit selfish, but I am hoping it will be helpful to some other folks who are going through the same struggles as I am. In this post, I will go over details of the "Hello World" experience for C++, while also going a bit beyond and understanding how the build pipeline works!
created_at: 2021-03-19 01:09:00.0000000 +0000 UTC
format: md
tags:
- C++
slugs:
- cpp-getting-started-with-the-basics-hello-world-and-build-pipeline
---

## Intro

I am probably the least qualified person to be writing a blog post about [C++](https://en.wikipedia.org/wiki/C%2B%2B). So, please approach this post with some caution. 
But, why am I writing it then? Well, I am currently learning C++, and it has been an unusual experience compared to my other programming-language-learning journeys. 
At this stage, my assumption is that the main difference that causes me to struggle comes from the fact that:

 - C++ doesn't have a universally agreed way to define how projects should be structured and built (well, sort of)
 - C++ also doesn't have one defined way to manage your dependencies (well, sort of)
 - Finally, C++ is not a programming language that will spit out an executable which will perform garbage collection for you out of the box (well, it really doesn't have this)

All of these are quite unusual characteristics for me when learning a programming language. I admit, I have been spoiled! The other thing to admit is that I have used C++ before in projects, but on all those occasions, the projects and build structures were already set up, and I only needed to maintain the codebase, and do occasionally changes in them. Also, these were projects which didn't run under a significant scale to generate those usual high-scale issues. So, nothing really required me to deeply understand how the C++ and its toolchain worked.

However, I knew that a few obstacle wouldn't wear me down, and I had to find a way to regain my perseverance! So, I thought what better way to make the learning stronger than blogging about my journey and pinning down my experience. And, here we are! You now know that the reason this post exists is a bit selfish, but I am hoping it will be helpful to some other folks who are going through the same while also acknowledging that everyone's mental model is different. So, [YMMV](https://www.urbandictionary.com/define.php?term=ymmv)!

If you are still here, let me tell you what this post is all about! I will be going through the "Hello World" experience for C++, while also taking the explanation a bit beyond and understanding how the build pipeline works by attempting to dive deep into the bowels of the compiler (but I also now I am probably only scratching the surface)!

## Hello World

I am learning a programming language. So, I should be concerned about the syntax, right? Well, usually but I really don't care about that at this stage, at least too much. At the moment, I am making an assumption that I will be able to get a handle of the syntax gradually as I start solving actual problems. What I really wanted to focus is the toolchain experience, and how I can glue things together. 

So, I focused on getting a bare minimum program written and understanding what goes in that as much as possible. As you can guess, it's the good-old "Hello World" program. Here is the code for that which is saved within the `main.cpp` file.

```c++
#include <iostream>

int main() {
    std::cout << "Hello World" << std::endl;
}
```

Nothing fancy here. However, I was able to learn so many things just from this small program!

 - First line is the [#include directive](https://docs.microsoft.com/en-us/cpp/preprocessor/hash-include-directive-c-cpp), which allows us to define the dependency between our source file and the "include file" (a.k.a. header file), which are the files that contains the constant and macro definitions, declarations of external variables and complex data types. These files don't contain the actual implementation. I know that you have more questions now, but hopefully I will be able to touch more into this in the upcoming posts.
 - C++ has a [Standard Library](https://en.wikipedia.org/wiki/C%2B%2B_Standard_Library), which contains a collection of classes and functions, which are written in the core language and part of the C++ ISO Standard itself. 
 - [`iostream`](https://en.wikipedia.org/wiki/Input/output_(C%2B%2B)#Input/output_streams) is part of the standard library, and provides C++ input and output fundamentals. `std::` prefix represents the standard library.
 - `cout` is coming from the standard library, and represents the standard output stream.
 - You can write into an output stream using the "[insertion operator](https://docs.microsoft.com/en-us/cpp/standard-library/using-insertion-operators-and-controlling-format)" (i.e. `<<`), which sends bytes to that output stream object.

## The Build Pipeline

The question now is how we can run this. C++ is a static typed language which requires a compilation step (well, it's more involved than compilation but hopefully we will get there, hang in there!). So, we need to compile the source code into an executable which we can run. This is where things got a bit more interesting for me because there isn't one compiler that you can use for C++. There are at least three of them (possibly more?): `g++`, `gcc`, `clang++` (well, _I guess_ you can also count [`c++` which is actually a symbolic link](https://stackoverflow.com/a/1712803/463785)). I honestly don't know at this stage enough about these to be able to tell you about the difference between them. So, for the purposes of not getting stuck, I am going to go with `g++` for now, and there is not a particular reason to why I chose it other than the fact that most examples I have come across so far have been using `g++` 🤷🏻‍♂️.

![](https://tugberkugurlu-blog.s3.us-east-2.amazonaws.com/post-images/01F140WSJCR5KE63NNXD48R488-Screenshot-2021-03-19-at-01.13.44.png)

> Here is something more fun! [`g++` on my Mac actually ends up calling `clang`](https://stackoverflow.com/a/19535525/463785). Go figure 🤷🏻‍♂️. If you are someone who understands the diff between C++ compilers, please direct me to a resource which would make me understand what is really going on here. I have given up on this for now 😕

OK, compiler choice is sorted out, kind of. Let's compile this now. Here is the simplest compilation we can run which will spit out an executable called `hello-world`:

```bash
g++ ./main.cpp -o hello-world
```

If we add `-v` flag, we will actually be able to see in more details what's going on underneath:

```
➜ g++ -v ./main.cpp -o hello-world
Apple LLVM version 10.0.1 (clang-1001.0.46.4)
Target: x86_64-apple-darwin18.7.0
Thread model: posix
InstalledDir: /Library/Developer/CommandLineTools/usr/bin
 "/Library/Developer/CommandLineTools/usr/bin/clang" -cc1 -triple x86_64-apple-macosx10.14.0 -Wdeprecated-objc-isa-usage -Werror=deprecated-objc-isa-usage -emit-obj -mrelax-all -disable-free -disable-llvm-verifier -discard-value-names -main-file-name main.cpp -mrelocation-model pic -pic-level 2 -mthread-model posix -mdisable-fp-elim -fno-strict-return -masm-verbose -munwind-tables -target-sdk-version=10.14 -target-cpu penryn -dwarf-column-info -debugger-tuning=lldb -target-linker-version 450.3 -v -resource-dir /Library/Developer/CommandLineTools/usr/lib/clang/10.0.1 -isysroot /Library/Developer/CommandLineTools/SDKs/MacOSX10.14.sdk -I/usr/local/include -stdlib=libc++ -Wno-atomic-implicit-seq-cst -Wno-framework-include-private-from-public -Wno-atimport-in-framework-header -Wno-quoted-include-in-framework-header -fdeprecated-macro -fdebug-compilation-dir /Users/tugberkugurlu/go/src/github.com/tugberkugurlu/cmake-getting-started/0-hello-world -ferror-limit 19 -fmessage-length 95 -stack-protector 1 -fblocks -fencode-extended-block-signature -fregister-global-dtors-with-atexit -fobjc-runtime=macosx-10.14.0 -fcxx-exceptions -fexceptions -fmax-type-align=16 -fdiagnostics-show-option -fcolor-diagnostics -o /var/folders/l4/2c_f_d8973z3g7lmkb9xcw9h0000gn/T/main-d3ea6f.o -x c++ ./main.cpp
clang -cc1 version 10.0.1 (clang-1001.0.46.4) default target x86_64-apple-darwin18.7.0
ignoring nonexistent directory "/Library/Developer/CommandLineTools/SDKs/MacOSX10.14.sdk/usr/include/c++/v1"
ignoring nonexistent directory "/Library/Developer/CommandLineTools/SDKs/MacOSX10.14.sdk/usr/local/include"
ignoring nonexistent directory "/Library/Developer/CommandLineTools/SDKs/MacOSX10.14.sdk/Library/Frameworks"
#include "..." search starts here:
#include <...> search starts here:
 /usr/local/include
 /Library/Developer/CommandLineTools/usr/include/c++/v1
 /Library/Developer/CommandLineTools/usr/lib/clang/10.0.1/include
 /Library/Developer/CommandLineTools/usr/include
 /Library/Developer/CommandLineTools/SDKs/MacOSX10.14.sdk/usr/include
 /Library/Developer/CommandLineTools/SDKs/MacOSX10.14.sdk/System/Library/Frameworks (framework directory)
End of search list.
 "/Library/Developer/CommandLineTools/usr/bin/ld" -demangle -lto_library /Library/Developer/CommandLineTools/usr/lib/libLTO.dylib -no_deduplicate -dynamic -arch x86_64 -macosx_version_min 10.14.0 -syslibroot /Library/Developer/CommandLineTools/SDKs/MacOSX10.14.sdk -o hello-world /var/folders/l4/2c_f_d8973z3g7lmkb9xcw9h0000gn/T/main-d3ea6f.o -L. -L/Users/tugberkugurlu/.tensorflow-1.11.0/lib -L/usr/local/lib -lc++ -lSystem /Library/Developer/CommandLineTools/usr/lib/clang/10.0.1/lib/darwin/libclang_rt.osx.a
```

Lots of useful things to unpack from this output, which really helped me understand how the compiler is behaving! I cannot say that I currently understand all of it, but let me try to explain what I have been able to extract from this so far:

 - `clang` compiler is called to compile the `main.cpp` file, and spit out the file called `main-d3ea6f.o`, [which contains the compiled object code](https://stackoverflow.com/a/2186252/463785). That file is being put under the temporary `/var/folders/l4/2c_f_d8973z3g7lmkb9xcw9h0000gn/T` folder, so that the compiler can refer back to it later.
 - The compilation will happen for the target `x86_64-apple-darwin18.7.0`. I am assuming this is used as the default because I am performing this compilation on my Mac, and I haven't specified a target for the compiler.
 - While compilation is happening, the compiler is looking under several directories which includes `/usr/local/include` and a few others right after hitting the `#include` directive. These directories are known as include directories, and these are where the header files are being looked for. In my case, `iostream` header file is located under `/Library/Developer/CommandLineTools/usr/include/c++/v1`.
 - Once the compilation is performed, [`ld`](https://linux.die.net/man/1/ld) is invoked. `ld` is the linker (see `man ld`), which combines several object files and libraries, resolves references and produces an output file. If you look at the output, you can see that `main-d3ea6f.o` file which contains our compiled object code is passed into `ld` as one of its arguments. There are also a few folders passed in as an argument here, and one of them is `/Users/tugberkugurlu/.tensorflow-1.11.0/lib`, which is a bit strange. The reason that's there is because it's set as one of the paths through the `LIBRARY_PATH` environment variable for me. [Colon-separated list of directories through this environment variable is used by the linker when searching for special linker files](https://gcc.gnu.org/onlinedocs/gcc/Environment-Variables.html).
 - You can also see the output of the `ld` is specified as `-o hello-world`, which is the name that we have given to `g++` compiler.

> I am most likely glossing over a lot of details here. I have found [this article](https://www.toptal.com/c-plus-plus/c-plus-plus-understanding-compilation) to be a very informative when it comes to explaining the C++ build pipeline, by breaking it down to three steps called preprocessing, compilation and linking. So, please check it out for more thorough explanation.

After the compilation and linking, we end up with an executable file called `hello-world`, and we can execute it to see our super complex output:

```
➜ ./hello-world
Hello World
```

## What's Next?

This was the very basic example. I am sure you can relate to the fact that almost none of the real world problems will be solved with this simple implementation. The next part for me will be to look at how we can work with a multi-file project as well as being able to take an external open source library as a dependency.

I am new on my C++ journey. So, if you see anything wrong and things that I can benefit from, please do leave a comment on this post 🙂

## Resources 

 - [How C++ Works: Understanding Compilation](https://www.toptal.com/c-plus-plus/c-plus-plus-understanding-compilation)
 - [What is the difference between g++ and gcc?](https://stackoverflow.com/questions/172587/what-is-the-difference-between-g-and-gcc)
 - [GCC vs. Clang/LLVM: An In-Depth Comparison of C/C++ Compilers](https://alibabatech.medium.com/gcc-vs-clang-llvm-an-in-depth-comparison-of-c-c-compilers-899ede2be378)
 - [Where Does GCC Look to Find its Header Files](https://commandlinefanatic.com/cgi-bin/showarticle.cgi?article=art026)