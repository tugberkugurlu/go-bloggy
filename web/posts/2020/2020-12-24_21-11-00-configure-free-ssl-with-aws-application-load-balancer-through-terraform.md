---
id: 01ETB9J6GCNSSECM81HHKFY8W0
title: Configure Free SSL Certificate on AWS Application Load Balancer Through Terraform
abstract: Last week, I have moved all my personal compute and storage from Azure to AWS, and started managing it through terraform. While doing so, I discovered that you can actually have SSL for your web application for free when using AWS Application Load Balancer. Setting it up was a bit tedious, and I wanted to share that experience here.
created_at: 2020-12-24 21:22:00.0000000 +0000 UTC
format: md
tags:
- AWS
- HTTP
slugs:
- configure-free-ssl-certificate-on-aws-application-load-balancer-through-terraform
---

Last week, I have moved all my personal compute and storage from [Azure](https://azure.microsoft.com/) to [AWS](https://aws.amazon.com/). I took this opportunity as an excuse to also start to manage all that infrastructure through Terraform. Why AWS though? I had had chance to use AWS before for a few projects that I worked on, but since [I joined Deliveroo 2 years ago](https://twitter.com/tourismgeek/status/1091003681003237376), I have been using AWS everyday. So, it's the least friction for me when it comes to working with a cloud provider. That said, this migration itself still has been a really great learning experience for me, and it emphasized it more for me that AWS is million miles ahead in their journey when it comes to developer experience. Things just work, especially when it comes to gluing things together (we will see in an example of that in this post). When they don't, it's also very obvious why which makes it easy to diagnose what's going wrong (although, [it's probably because of IAM](https://nodramadevops.com/2019/11/why-is-aws-iam-so-hard/) for like 99.9% of the cases).

During this migration, I have also discovered that you can actually configure SSL on your own domain free of any additional charges through [AWS Certificate Manager](https://aws.amazon.com/certificate-manager/) (ACM) if you are already using [AWS Application Load Balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) (ALB). This was a valuable find for me, as I needed to enable HTTPS for this blog which I have been procrastinating to get it done, like forever. However, when I think about it, it wasn't only the additional payment I had to make for the SSL certificate that was making me delay getting one. It was probably the cost of maintenance which was the biggest chore that I didn't really want (e.g. certificate renewals and all that). 

ALB and ACM integration addresses these both issues, by providing a way to configure SSL as well as keeping to renewed for you free of any additional charges. To be fair, there is probably also a way to automate this all on Azure, but I have been also away from that world for over 2 years now, and I didn't have the mental capacity to sort it out. Anyway, enough with the excuses, and let's see how to make this all sorted through Terraform.

> Note that I am going to skip what AWS ALB is, how it works, and how to configure it to start directing traffic to your resources (e.g. ECS services, Lambda, EC2 instances, etc.). However, it's worth checking out [the ALB documentation](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) before this post if you don't have a good grasp of its concepts.

## Creating the Certificate Through AWS Certificate Manager

## Domain Name Validation Through Route53 DNS Configuration

## Wiring It up with Application Load Balancer

### Redirecting HTTP Traffic Through an ALB Rule