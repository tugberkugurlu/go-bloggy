---
id: 01F17XQ4E37MX14WARR80KVGS6
title: "C++, Getting Started with the Basics: Working with Dependencies and Linker"
abstract: I am learning C++, and what better way to make the learning stick more stronger than blogging about my journey and experience, especially thinking that the barrier to entry is quite high and there is too much to learn. So, the reason that this post exists is a bit selfish, but I am hoping it will be helpful to some other folks who are going through the same struggles as I am. In this post, I will go over details of what it takes to work with dependencies in C++ and how we can manage them in a simple C++ project.
created_at: 2021-03-21 22:01:00.0000000 +0000 UTC
format: md
tags:
- C++
slugs:
- cpp-getting-started-with-the-basics-working-with-dependencies-and-linker
---

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

> My understanding around the namespaces in C++ hasn't been formed fully yet. When I thought I grasped it and it seemed to be similar to how it works in C#, I came across an interesting behavior which I documented in [this Stackoverflow question](https://stackoverflow.com/questions/66739009). I would suggest for you to check that out first before basing any assumptions on this.

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

## Specifying Header File Directories to Compiler

```
apt-get download libgflags-dev
```

```
ls
Dockerfile  a.out  build.sh  libgflags-dev_2.2.2-1build1_amd64.deb  main.cpp
```

```
mkdir tmp
dpkg-deb -R libgflags-dev_2.2.2-1build1_amd64.deb tmp
```

```
cd tmp
ls -R  
.:
DEBIAN	usr

./DEBIAN:
control  md5sums

./usr:
include  lib  share

./usr/include:
gflags

./usr/include/gflags:
gflags.h  gflags_completions.h	gflags_declare.h  gflags_gflags.h

./usr/lib:
x86_64-linux-gnu

./usr/lib/x86_64-linux-gnu:
cmake  libgflags.a  libgflags.so  libgflags_nothreads.a  libgflags_nothreads.so  pkgconfig

./usr/lib/x86_64-linux-gnu/cmake:
gflags

./usr/lib/x86_64-linux-gnu/cmake/gflags:
gflags-config-version.cmake		  gflags-nonamespace-targets.cmake
gflags-config.cmake			  gflags-targets-release.cmake
gflags-nonamespace-targets-release.cmake  gflags-targets.cmake

./usr/lib/x86_64-linux-gnu/pkgconfig:
gflags.pc

./usr/share:
doc

./usr/share/doc:
libgflags-dev

./usr/share/doc/libgflags-dev:
changelog.Debian.gz  copyright
```

## Specifying Library Location to Linker

## Resources

These are the resources I benefited from while writing this post. It's only fair I give these some credit. They might not entirely beneficial to you though:

 - [Where Does GCC Look to Find its Header Files?](https://commandlinefanatic.com/cgi-bin/showarticle.cgi?article=art026)
 - [How C++ Works: Understanding Compilation](https://www.toptal.com/c-plus-plus/c-plus-plus-understanding-compilation)
 - [What is the purpose of a .cmake file?](https://stackoverflow.com/questions/46456498/what-is-the-purpose-of-a-cmake-file)
 - [Easily unpack DEB, edit postinst, and repack DEB](https://unix.stackexchange.com/questions/138188/easily-unpack-deb-edit-postinst-and-repack-deb)
 - [Is there an apt command to download a deb file from the repositories to the current directory?](https://askubuntu.com/questions/30482/is-there-an-apt-command-to-download-a-deb-file-from-the-repositories-to-the-curr)
 - [Name lookup](https://en.cppreference.com/w/cpp/language/lookup)
 - [Phases of translation](https://en.cppreference.com/w/cpp/language/translation_phases)
 - [How does the compilation/linking process work?](https://stackoverflow.com/questions/6264249/how-does-the-compilation-linking-process-work)