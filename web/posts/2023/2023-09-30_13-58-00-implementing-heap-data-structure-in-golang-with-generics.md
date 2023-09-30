---
id: 01HBK5YS7GWXXKX7XT0RDE4908
title: Implementing Heap Data Structure in Go (Golang) with Generics
abstract: 
created_at: 2023-09-30 13:58:00.0000000 +0000 UTC
format: md
tags:
- Go
- Golang
- Data Structures
- Algorithms
slugs:
- implementing-heap-data-structure-in-go-golang-with-generics
---

I previously wrote about [Heap data structure implementation in Go](https://www.tugberkugurlu.com/archive/usage-of-the-heap-data-structure-in-go-golang-with-examples), but that was at the pre-generics era. While 
Go provides a package called [container/heap](https://golang.org/pkg/container/heap/) which has heap operations for any type that implements [heap.Interface](https://golang.org/pkg/container/heap/#Interface), working with this structure was a bit 
cumbersome since this involved creating separate boilerplate code for each type that you wanted this to work with.

https://github.com/zyedidia/generic