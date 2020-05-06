---
title: Should I await on Task.FromResult Method Calls?
abstract: Task class has a static method called FromResult which returns an already
  completed (at the RanToCompletion status) Task object. I have seen a few developers
  "await"ing on Task.FromResult method call and this clearly indicates that there
  is a misunderstanding here. I'm hoping to clear the air a bit with this post.
created_at: 2014-02-24 21:14:00 +0000 UTC
tags:
- async
- C#
- TPL
slugs:
- should-i-await-on-task-fromresult-method-calls
---
