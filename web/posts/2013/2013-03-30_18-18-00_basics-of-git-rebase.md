---
id: 18535220-0cd4-4be0-807f-ccaebe01a1aa
title: Basics of Git Rebase
abstract: Git Rebase is only one of the powerful features of git and it allows you
  to have a clean history in a highly branching workflow.
created_at: 2013-03-30 18:18:00 +0000 UTC
tags:
- Git
- Tips
slugs:
- basics-of-git-rebase
---

<p>You may wonder why the title starts with "Basics". The answer is simple: I know only the basics of git rebase :) It's only one of the powerful features of git and it allows you to have a clean history in a highly branching workflow. "Rebase" is quite powerful as mentioned and what I'm about to show you is only one of the reasons why to use rebase. I highly recommend Keith <a href="http://vimeo.com/43659036">Dahlby's NDC talk</a> which he took some time to show the rebase feature.</p>
<p>Let's see the easiest sample where rebase comes handy. We have the following history where we have two branches: master and feature-1.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_thumb.png" width="644" height="390" /></a></p>
<p>Typically, what you would do here is to merge the feature-1 branch onto master which is fairly reasonable and it works. However, it creates you a unnecessary commit + a ridiculous graph which would be a mess if you think of hundreds of branches:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_3.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_thumb_3.png" width="644" height="345" /></a></p>
<p>What you can do with rebase is to patch the feature-1 branch onto master. Later then, you can merge from there. The following command is what you need to run:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_4.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_thumb_4.png" width="644" height="217" /></a></p>
<p>After running the rebase command, we can run "gitk &ndash;all" to see the graph:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_5.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_thumb_5.png" width="644" height="345" /></a></p>
<p>It's now nice clean history. Notice that the master is still pointing where it was. It's because we haven't merge the feature-1 branch yet. Let's checkout to master branch and run "git merge feature-1" to merge feature-1 branch onto master branch:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image3e6a86da-a24a-497e-a74a-afc55251e34e.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_thumb_6.png" width="644" height="296" /></a></p>
<p>Nicely done! Open up the gitk one more time and see the clean history:</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_6.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_thumb_7.png" width="644" height="349" /></a></p>
<p>After we remove the feature-1 branch by running "git branch &ndash;D feature-1", we won't have any trace from feature-1 branch which is absolutely OK as feature branches are just the implementation details, that's all.</p>
<p><a href="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image.png"><img title="image" style="background-image: none; padding-top: 0px; padding-left: 0px; display: inline; padding-right: 0px; border: 0px;" border="0" alt="image" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/image_thumb_8.png" width="644" height="349" /></a></p>
<h3>Rebase can hurt</h3>
<p>With git rebase, at the very basic level, you are messing with the history which can be dangerous depending on the case. On the other hand, when you have a collision, it's not a picnic to solve those collisions with interactive rebase without a deep firsthand knowledge but it's worth looking into even if it seems hard at the first glance <img class="wlEmoticon wlEmoticon-smile" style="border-style: none;" alt="Smile" src="https://www.tugberkugurlu.com/Content/images/Uploadedbyauthors/wlw/Basics-of-Git-Rebase_11F73/wlEmoticon-smile.png" /></p>