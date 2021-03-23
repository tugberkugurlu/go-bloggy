---
id: 01F17XQ4E37MX14WARR80KVGS6
title: "C++, Getting Started with the Basics: Working with Dependencies and Linker"
abstract: I am learning C++, and what better way to make the learning stick more stronger than blogging about my journey and experience, especially thinking that the barrier to entry is quite high and there is too much to learn. So, the reason that this post exists is a bit selfish, but I am hoping it will be helpful to some other folks who are going through the same struggles as I am. In this post, I will go over details of what it takes to work with dependencies in C++ and how the compilation and linking process works.
created_at: 2021-03-21 22:01:00.0000000 +0000 UTC
format: md
tags:
- C++
slugs:
- cpp-getting-started-with-the-basics-working-with-dependencies-and-linker
---

## Intro

As I mentioned in [my previous post about C++](https://www.tugberkugurlu.com/archive/cpp-getting-started-with-the-basics-hello-world-and-build-pipeline), I am learning [C++](https://www.tugberkugurlu.com/tags/cpp). It has been a bumpy ride so far, and C++ is certainly not an easy to pick up programming language! So, I thought what better way to make the learning stronger than blogging about my journey and pinning down my experience. You now know that the reason this post exists is a bit selfish, but I am hoping it will be helpful to some other folks who are going through the same while also acknowledging that everyone's mental model is different. So, [YMMV](https://www.urbandictionary.com/define.php?term=ymmv).

In this post, I want to share my experience of incorporating a 3rd party dependency into my own program, and understanding what goes under the hood during the compilation and linking phase of the build process.

## Referencing a Dependency in Code

The way you should be specifying a dependency is in your code is through the [`#include` directive](https://docs.microsoft.com/en-us/cpp/preprocessor/hash-include-directive-c-cpp?view=msvc-160) by specifying the header file that you want to take a dependency on. As we probably know by now that the header file doesn't actually contain the implementation, but only declares the contract between the library and the consumer. We will shortly touch on how we will be able to tie the header file with a dependency.

As we can see inside the [`gflags` documentation](https://gflags.github.io/gflags/), the header file we want to work with is called `gflags/gflags.h`. Now that immediately raised some questions for me, and I am sure it will for you if you happen to be a newbie in C++ World like me. The biggest one of all is where `gflgas/` folder is relative to. That will become more clear when it comes to the building part. So, for now, let's assume it's magic™️.

As we learned about how to take a dependency on this library within the code, here is how our sample program looks like:

```cpp
#include <iostream>
#include <gflags/gflags.h>

DEFINE_string(name, "Tugberk", "Name of the person to greet");

int main(int argc, char *argv[]) {
    gflags::ParseCommandLineFlags(&argc, &argv, true);
    std::cout << "Hello " << FLAGS_name << std::endl;
}
```

Nothing fancy, and you can see the [`gflags` documentation](https://gflags.github.io/gflags/) about the specifics of our usage here. The purpose of this post is not to explain that, and the only reason that we are using `gflags` here to demonstrate how to take a dependency on an external library.

However, one thing that's worth noting is the usage of `gflags::` before the `ParseCommandLineFlags` function call. That `gflags` that's being referred here is the namespace, which we are betting that it will be declared within the `gflags.h` header file. `gflags::ParseCommandLineFlags` is the fully-qualified reference to the function we want to invoke. 

Alternatively, we could have imported the entire `gflags` namespace, and be able to call `ParseCommandLineFlags` directly w/o namespace declaration like below, which would mean that you can use anything under that namespace directly:

```cpp
#include <iostream>
#include <gflags/gflags.h>

using namespace gflags;

DEFINE_string(name, "Tugberk", "Name of the person to greet");

int main(int argc, char *argv[]) {
    ParseCommandLineFlags(&argc, &argv, true);
    std::cout << "Hello " << FLAGS_name << std::endl;
}
```

Based on my understanding, there is nothing wrong with this in terms of performance of the program or the compiler (I could be wrong, don't quote me on this). However, this will likely increase your changes of having a name collisions, and also it will make it a bit hard to read the code (i.e. it's not immediately clear where `ParseCommandLineFlags` is coming from).

One other alternative is to just declare a using for the type you want to use:

```cpp
#include <iostream>
#include <gflags/gflags.h>

using gflags::ParseCommandLineFlags;

DEFINE_string(name, "Tugberk", "Name of the person to greet");

int main(int argc, char *argv[]) {
    ParseCommandLineFlags(&argc, &argv, true);
    std::cout << "Hello " << FLAGS_name << std::endl;
}
```

Although this still suffers from the same problems I listed above to a certain extent, this is a bit better especially when you are planning to use the defined type a few times within the same file.

Final thing I want to note within this code is the use of `DEFINE_string`. It's also defined within the same header file. However, that's a [Macro](https://docs.microsoft.com/en-us/cpp/preprocessor/macros-c-cpp?view=msvc-160) and [it doesn't seem to be tied to a namespace](https://github.com/gflags/gflags/blob/827c769e5fc98e0f2a34c47cef953cc6328abced/src/gflags.h.in#L595-L620). I don't have much info about Macros at this stage, but wanted to touch on the rationale of why it's being used in this way.

## Building with a Library Dependency

We have our implementation which should give us a command like program where we can call `hello-world --name Bob` and that would print out `Hello Bob` for us. To be able to demonstrate different build variations, I am going to run the build within a [Docker](https://www.tugberkugurlu.com/tags/docker) container. Configuration for this is going to be very simple. The code we have seen above will be inside the `main.cpp` file. Also to start with, we will also have a `build.sh` file with the following content:

```bash
#!/bin/bash

g++ -v ./main.cpp -o hello-world
```

`-v` is here to give verbose output from the compiler which will be handy when it comes to understanding what goes under the hood. The `Dockerfile` content will be as following:

```dockerfile
FROM ubuntu

RUN apt-get update && apt-get -y install build-essential

WORKDIR /opt/
RUN mkdir app
WORKDIR /opt/app

COPY ./ ./

RUN ./build.sh
CMD ["./hello-world", "--name=Bob"]
```

When I run `docker build .` with this setup, I'm getting an error:

```
...
...
Step 7/7 : RUN ./build.sh
 ---> Running in 39ce491a452e
Using built-in specs.
COLLECT_GCC=g++
COLLECT_LTO_WRAPPER=/usr/lib/gcc/x86_64-linux-gnu/9/lto-wrapper
OFFLOAD_TARGET_NAMES=nvptx-none:hsa
OFFLOAD_TARGET_DEFAULT=1
Target: x86_64-linux-gnu
Configured with: ../src/configure -v --with-pkgversion='Ubuntu 9.3.0-17ubuntu1~20.04' --with-bugurl=file:///usr/share/doc/gcc-9/README.Bugs --enable-languages=c,ada,c++,go,brig,d,fortran,objc,obj-c++,gm2 --prefix=/usr --with-gcc-major-version-only --program-suffix=-9 --program-prefix=x86_64-linux-gnu- --enable-shared --enable-linker-build-id --libexecdir=/usr/lib --without-included-gettext --enable-threads=posix --libdir=/usr/lib --enable-nls --enable-clocale=gnu --enable-libstdcxx-debug --enable-libstdcxx-time=yes --with-default-libstdcxx-abi=new --enable-gnu-unique-object --disable-vtable-verify --enable-plugin --enable-default-pie --with-system-zlib --with-target-system-zlib=auto --enable-objc-gc=auto --enable-multiarch --disable-werror --with-arch-32=i686 --with-abi=m64 --with-multilib-list=m32,m64,mx32 --enable-multilib --with-tune=generic --enable-offload-targets=nvptx-none=/build/gcc-9-HskZEa/gcc-9-9.3.0/debian/tmp-nvptx/usr,hsa --without-cuda-driver --enable-checking=release --build=x86_64-linux-gnu --host=x86_64-linux-gnu --target=x86_64-linux-gnu
Thread model: posix
gcc version 9.3.0 (Ubuntu 9.3.0-17ubuntu1~20.04) 
COLLECT_GCC_OPTIONS='-v' '-o' 'hello-world' '-shared-libgcc' '-mtune=generic' '-march=x86-64'
 /usr/lib/gcc/x86_64-linux-gnu/9/cc1plus -quiet -v -imultiarch x86_64-linux-gnu -D_GNU_SOURCE ./main.cpp -quiet -dumpbase main.cpp -mtune=generic -march=x86-64 -auxbase main -version -fasynchronous-unwind-tables -fstack-protector-strong -Wformat -Wformat-security -fstack-clash-protection -fcf-protection -o /tmp/ccebxWeM.s
GNU C++14 (Ubuntu 9.3.0-17ubuntu1~20.04) version 9.3.0 (x86_64-linux-gnu)
	compiled by GNU C version 9.3.0, GMP version 6.2.0, MPFR version 4.0.2, MPC version 1.1.0, isl version isl-0.22.1-GMP

GGC heuristics: --param ggc-min-expand=100 --param ggc-min-heapsize=131072
ignoring duplicate directory "/usr/include/x86_64-linux-gnu/c++/9"
ignoring nonexistent directory "/usr/local/include/x86_64-linux-gnu"
ignoring nonexistent directory "/usr/lib/gcc/x86_64-linux-gnu/9/include-fixed"
ignoring nonexistent directory "/usr/lib/gcc/x86_64-linux-gnu/9/../../../../x86_64-linux-gnu/include"
#include "..." search starts here:
#include <...> search starts here:
 /usr/include/c++/9
 /usr/include/x86_64-linux-gnu/c++/9
 /usr/include/c++/9/backward
 /usr/lib/gcc/x86_64-linux-gnu/9/include
 /usr/local/include
 /usr/include/x86_64-linux-gnu
 /usr/include
End of search list.
GNU C++14 (Ubuntu 9.3.0-17ubuntu1~20.04) version 9.3.0 (x86_64-linux-gnu)
	compiled by GNU C version 9.3.0, GMP version 6.2.0, MPFR version 4.0.2, MPC version 1.1.0, isl version isl-0.22.1-GMP

GGC heuristics: --param ggc-min-expand=100 --param ggc-min-heapsize=131072
Compiler executable checksum: 466f818abe2f30ba03783f22bd12d815
./main.cpp:2:10: fatal error: gflags/gflags.h: No such file or directory
    2 | #include <gflags/gflags.h>
      |          ^~~~~~~~~~~~~~~~~
compilation terminated.
The command '/bin/sh -c ./build.sh' returned a non-zero code: 1
```

There are a few important things to call out here:

 - As you may remember from [the previous C++ post](https://www.tugberkugurlu.com/archive/cpp-getting-started-with-the-basics-hello-world-and-build-pipeline), the compiler is looking under several directories which includes `/usr/local/include` and a few others right after hitting the `#include` directives during its preprocessing stage.
 - We can see that compilation is failing with the following error: `gflags/gflags.h: No such file or directory`. That's giving us an indication that the header file with the path of `gflags/gflags.h` wasn't found in any of the include directories which the compiler was searching under.

This is an expected error at this stage, because gflags is a 3rd party library, this is a fresh box and we didn't install that library. 

## Back to Basics: Compilation of a C++ Program

Let's pause a bit and learn some fundamentals. I kept mentioning compilation, as it's a black box where you give it some input and get an output, compiled object back. Most of the time, this type of thinking will get us where we want to be. However, my aim here is to understand what's going on under the hood a bit more. When I have done that for C++ build process, I have found out that the build step is broken down into three independent steps:

 - **Preprocessing**: This stage handles the preprocessor directives, like `#include` and `#define`. After the processing of these directives, the preprocessor produces a single output.
 - **Compilation**: The compilation step is performed on each output of the preprocessor, and this is the step where the C++ code is converted into assembly code. This step also involves the assembler to turn the assembly code into machine code, then producing an actual binary file (a.k.a. object file). The bit that's super interesting at this stage is that **these object files can refer to symbols that are not defined**, and this is how the header files are being compiled at this stage without any specific implementation.
 - **Linking**: This is the final step within our build process, and this steps is handled through the linker which produces the final output for our program from the object files the compiler produced. This output can be either a library or an executable. It links all the object files by replacing the references to undefined symbols with the correct addresses, and **if the definitions exist in libraries other than the standard one, the linker needs to be informed about these specificity**, which is relevant to what I am trying to achieve in this post (more on this to come later).

You can check out [this incredible Stackoverflow answer](https://stackoverflow.com/a/6264256/463785) on this topic, which explains compilation steps of a C++ program more in-depth, and I copied most of what I mentioned in this section from there.

This 

## Preprocessing the Headers

Let's install this according to [the installation guidelines of this library](https://github.com/gflags/gflags/blob/827c769e5fc98e0f2a34c47cef953cc6328abced/INSTALL.md#installing-a-binary-distribution-package), and rerun the compilation.

```patch
diff --git a/1-dependency/Dockerfile b/1-dependency/Dockerfile
index fbaeba8..58215ea 100644
--- a/1-dependency/Dockerfile
+++ b/1-dependency/Dockerfile
@@ -1,6 +1,7 @@
 FROM ubuntu
 
 RUN apt-get update && apt-get -y install build-essential
+RUN apt-get -y install libgflags-dev
 
 WORKDIR /opt/
 RUN mkdir app
```

If I run `docker build .` command again, it still gives me an error but this time error is different:

```
COLLECT_GCC_OPTIONS='-v' '-o' 'hello-world' '-shared-libgcc' '-mtune=generic' '-march=x86-64'
 /usr/lib/gcc/x86_64-linux-gnu/9/collect2 -plugin /usr/lib/gcc/x86_64-linux-gnu/9/liblto_plugin.so -plugin-opt=/usr/lib/gcc/x86_64-linux-gnu/9/lto-wrapper -plugin-opt=-fresolution=/tmp/ccjLVDaH.res -plugin-opt=-pass-through=-lgcc_s -plugin-opt=-pass-through=-lgcc -plugin-opt=-pass-through=-lc -plugin-opt=-pass-through=-lgcc_s -plugin-opt=-pass-through=-lgcc --build-id --eh-frame-hdr -m elf_x86_64 --hash-style=gnu --as-needed -dynamic-linker /lib64/ld-linux-x86-64.so.2 -pie -z now -z relro -o hello-world /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/Scrt1.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crti.o /usr/lib/gcc/x86_64-linux-gnu/9/crtbeginS.o -L/usr/lib/gcc/x86_64-linux-gnu/9 -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../../lib -L/lib/x86_64-linux-gnu -L/lib/../lib -L/usr/lib/x86_64-linux-gnu -L/usr/lib/../lib -L/usr/lib/gcc/x86_64-linux-gnu/9/../../.. /tmp/ccg11F4K.o -lstdc++ -lm -lgcc_s -lgcc -lc -lgcc_s -lgcc /usr/lib/gcc/x86_64-linux-gnu/9/crtendS.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crtn.o
/usr/bin/ld: /tmp/ccg11F4K.o: in function `main':
main.cpp:(.text+0x27): undefined reference to `google::ParseCommandLineFlags(int*, char***, bool)'
/usr/bin/ld: /tmp/ccg11F4K.o: in function `__static_initialization_and_destruction_0(int, int)':
main.cpp:(.text+0x12e): undefined reference to `google::FlagRegisterer::FlagRegisterer<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > >(char const*, char const*, char const*, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >*, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >*)'
collect2: error: ld returned 1 exit status
The command '/bin/sh -c ./build.sh' returned a non-zero code: 1
```

We are still not quite there yet. However, as a software engineer, you know that this is a great feeling! You made some progress, and the changes that you have just made had some impact to move you forward 🙂

What's happened here is that the compiler was able to find the header file to be able to preprocess the `#include` directives. However, where did it find it? We can try to look for `gflags.h` file inside the container and see where it's located:

```
# find / -iname gflags.h
/usr/include/gflags/gflags.h
```

This makes more sense now as `/usr/include` is one of the directories where the compiler is looking for to find the header files.

# Linking

The error we have received this time seems to be coming from [`ld`, the linker](https://ftp.gnu.org/old-gnu/Manuals/ld-2.9.1/html_node/ld_toc.html), and it seems to be indicating that there are undefined references to several objects and functions under `google` namespace.

```
/usr/bin/ld: /tmp/ccg11F4K.o: in function `main':
main.cpp:(.text+0x27): undefined reference to `google::ParseCommandLineFlags(int*, char***, bool)'
```

> It's worth noting where this `google::` namespace comes from. This library seems to be exposed under two namespaces: `gflags` and `google`. All the documentation is referring to `gflags`. However, it seems like any usage under that namespace eventually seems to be redirected to `google` namespace. It took a while for me to understand why and how, but I documented the investigation in [this Stackoverflow question](https://stackoverflow.com/questions/66739009). I would suggest for you to check that out first before basing any assumptions on the namespace usage.

This is also expected, as haven't told the compiler yet what library dependency we want to link to, a.k.a archive, or static library. For static library files, the filenames always start with `lib`, and end with `.a` (archive, static library) on Unix/Linux (see [this post](https://www.bogotobogo.com/cplusplus/libraries.php) for reference). We can use the `-l` command line option of the `g++` compiler, which would eventually pass this to `ld` to [add the archive file to the list of files to link](https://ftp.gnu.org/old-gnu/Manuals/ld-2.9.1/html_node/ld_3.html). This option may be used any number of times. `ld` will search its path-list for occurrences of `lib{archive}.a` for every `{archive}` specified.

With this in mind, we should be able to complete our compilation journey by passing `-lgflags` option to `g++` compiler:

> The error output above might be confusing you since it seems like `/usr/lib/gcc/x86_64-linux-gnu/9/collect2` is invoked directly, not `ld`. Quick search suggests to me that [`collect2`](https://gcc.gnu.org/onlinedocs/gccint/Collect2.html) is eventually calls `ld` but I am not sure at this stage why and how the compiler located `collect2` at the first place, and decided to call it instead of calling `ld` directly. For simplicity, I will ignore `collect2` for the rest of the post, and only mention `ld`.

```bash
#!/bin/bash

g++ -v ./main.cpp -lgflags -o hello-world
```

Now, let's run `docker build .` with this setup:

```
...
...
GGC heuristics: --param ggc-min-expand=100 --param ggc-min-heapsize=131072
ignoring duplicate directory "/usr/include/x86_64-linux-gnu/c++/9"
ignoring nonexistent directory "/usr/local/include/x86_64-linux-gnu"
ignoring nonexistent directory "/usr/lib/gcc/x86_64-linux-gnu/9/include-fixed"
ignoring nonexistent directory "/usr/lib/gcc/x86_64-linux-gnu/9/../../../../x86_64-linux-gnu/include"
#include "..." search starts here:
#include <...> search starts here:
 /usr/include/c++/9
 /usr/include/x86_64-linux-gnu/c++/9
 /usr/include/c++/9/backward
 /usr/lib/gcc/x86_64-linux-gnu/9/include
 /usr/local/include
 /usr/include/x86_64-linux-gnu
 /usr/include
End of search list.
GNU C++14 (Ubuntu 9.3.0-17ubuntu1~20.04) version 9.3.0 (x86_64-linux-gnu)
	compiled by GNU C version 9.3.0, GMP version 6.2.0, MPFR version 4.0.2, MPC version 1.1.0, isl version isl-0.22.1-GMP

GGC heuristics: --param ggc-min-expand=100 --param ggc-min-heapsize=131072
Compiler executable checksum: 466f818abe2f30ba03783f22bd12d815
COLLECT_GCC_OPTIONS='-v' '-o' 'hello-world' '-shared-libgcc' '-mtune=generic' '-march=x86-64'
 as -v --64 -o /tmp/ccZZScyH.o /tmp/ccvjxfiH.s
GNU assembler version 2.34 (x86_64-linux-gnu) using BFD version (GNU Binutils for Ubuntu) 2.34
COMPILER_PATH=/usr/lib/gcc/x86_64-linux-gnu/9/:/usr/lib/gcc/x86_64-linux-gnu/9/:/usr/lib/gcc/x86_64-linux-gnu/:/usr/lib/gcc/x86_64-linux-gnu/9/:/usr/lib/gcc/x86_64-linux-gnu/
LIBRARY_PATH=/usr/lib/gcc/x86_64-linux-gnu/9/:/usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/:/usr/lib/gcc/x86_64-linux-gnu/9/../../../../lib/:/lib/x86_64-linux-gnu/:/lib/../lib/:/usr/lib/x86_64-linux-gnu/:/usr/lib/../lib/:/usr/lib/gcc/x86_64-linux-gnu/9/../../../:/lib/:/usr/lib/
COLLECT_GCC_OPTIONS='-v' '-o' 'hello-world' '-shared-libgcc' '-mtune=generic' '-march=x86-64'
 /usr/lib/gcc/x86_64-linux-gnu/9/collect2 -plugin /usr/lib/gcc/x86_64-linux-gnu/9/liblto_plugin.so -plugin-opt=/usr/lib/gcc/x86_64-linux-gnu/9/lto-wrapper -plugin-opt=-fresolution=/tmp/cc10fZHH.res -plugin-opt=-pass-through=-lgcc_s -plugin-opt=-pass-through=-lgcc -plugin-opt=-pass-through=-lc -plugin-opt=-pass-through=-lgcc_s -plugin-opt=-pass-through=-lgcc --build-id --eh-frame-hdr -m elf_x86_64 --hash-style=gnu --as-needed -dynamic-linker /lib64/ld-linux-x86-64.so.2 -pie -z now -z relro -o hello-world /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/Scrt1.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crti.o /usr/lib/gcc/x86_64-linux-gnu/9/crtbeginS.o -L/usr/lib/gcc/x86_64-linux-gnu/9 -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../../lib -L/lib/x86_64-linux-gnu -L/lib/../lib -L/usr/lib/x86_64-linux-gnu -L/usr/lib/../lib -L/usr/lib/gcc/x86_64-linux-gnu/9/../../.. /tmp/ccZZScyH.o -lgflags -lstdc++ -lm -lgcc_s -lgcc -lc -lgcc_s -lgcc /usr/lib/gcc/x86_64-linux-gnu/9/crtendS.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crtn.o
COLLECT_GCC_OPTIONS='-v' '-o' 'hello-world' '-shared-libgcc' '-mtune=generic' '-march=x86-64'
Removing intermediate container ce5a3c257fe2
 ---> 455abaa9d2d9
Step 9/9 : CMD ["./hello-world", "--name=Bob"]
 ---> Running in 2b17e00b3210
Removing intermediate container 2b17e00b3210
 ---> cc8ae20c8aa8
Successfully built cc8ae20c8aa8
```

Build passed! If we look at the compiler output from this, we should be able to see that `-lgflags` option is passed to the linker:

![](https://tugberkugurlu-blog.s3.us-east-2.amazonaws.com/post-images/01F1G3FVJPT67X1TW1MWXBGWPJ-Screenshot-2021-03-21-at-22.51.31.png)

Based on the information we have about the linker and with the `-lgflags` option being passed to it now, we know that the linker is looking for `libgflags.a` static library file to use as part of the linking process. Where did it find it though, and how did it knew to look there at the first place? Let's look for that file within the container:

```
➜ docker run -it cc8ae20c8aa8 /bin/sh
# find / -iname libgflags.a
/usr/lib/x86_64-linux-gnu/libgflags.a
```

That seems to be existing under `/usr/lib/x86_64-linux-gnu` folder. This is [the folder where architecture specific libraries live](https://unix.stackexchange.com/a/43214) under ubuntu. If we also look at what's being passed to the linker through the `-L` command line option, which adds a path to the list of paths that `ld` will search for archive libraries and `ld` control scripts, we will see that `/usr/lib/x86_64-linux-gnu` is already bing passed.

![](https://tugberkugurlu-blog.s3.us-east-2.amazonaws.com/post-images/01F1G46H1V0G9KJCGZN2NVDB0W-Screenshot-2021-03-21-at-22.50.58.png)


Nice, the C++ build process is now making more sense for me 🙂  

Just to make sure things are working as expect, I will run the container I have just built.

```
➜ docker run cc8ae20c8aa8                  
Hello Bob
➜ docker run cc8ae20c8aa8 ./hello-world --name=Alice
Hello Alice
```

It works as expected 🎉

## Resources

These are the resources I benefited from while writing this post. It's only fair I give these some credit. They might not entirely beneficial to you though:

 - [Where Does GCC Look to Find its Header Files?](https://commandlinefanatic.com/cgi-bin/showarticle.cgi?article=art026)
 - [How C++ Works: Understanding Compilation](https://www.toptal.com/c-plus-plus/c-plus-plus-understanding-compilation)
 - [What is the purpose of a .cmake file?](https://stackoverflow.com/questions/46456498/what-is-the-purpose-of-a-cmake-file)
 - [Easily unpack DEB, edit postinst, and repack DEB](https://unix.stackexchange.com/questions/138188/easily-unpack-deb-edit-postinst-and-repack-deb)
 - [Is there an apt command to download a deb file from the repositories to the current directory?](https://askubuntu.com/questions/30482/is-there-an-apt-command-to-download-a-deb-file-from-the-repositories-to-the-curr)
 - [Name lookup](https://en.cppreference.com/w/cpp/language/lookup)
 - [Phases of translation](https://en.cppreference.com/w/cpp/language/translation_phases)
 - [How does the compilation/linking process work?](https://stackoverflow.com/questions/6264249)
 - [Using fully qualified names in C++](https://stackoverflow.com/questions/17281214/using-fully-qualified-names-in-c)
 - [0.5 — Introduction to the compiler, linker, and libraries](https://www.learncpp.com/cpp-tutorial/introduction-to-the-compiler-linker-and-libraries/)
 - [The GNU linker](https://ftp.gnu.org/old-gnu/Manuals/ld-2.9.1/html_node/ld_toc.html)
 - [ld: Command Line Options](https://ftp.gnu.org/old-gnu/Manuals/ld-2.9.1/html_node/ld_3.html)