---
id: 01ETB9J6GCNSSECM81HHKFY8W0
title: Configure Free Wildcard SSL Certificate on AWS Application Load Balancer (ALB) Through Terraform
abstract: Last week, I have moved all my personal compute and storage from Azure to AWS, and started managing it through terraform. While doing so, I discovered that you can actually have SSL for your web application for free when using AWS Application Load Balancer. Setting it up was a bit tedious, and I wanted to share that experience here.
created_at: 2020-12-24 21:22:00.0000000 +0000 UTC
format: md
tags:
- AWS
- HTTP
- Security
slugs:
- configure-free-wildcard-ssl-certificate-on-aws-application-load-balancer-through-terraform
---

Last week, I have moved all my personal compute and storage from [Azure](https://azure.microsoft.com/) to [AWS](https://aws.amazon.com/). I took this opportunity as an excuse to also start to manage all that infrastructure through [Terraform](https://www.terraform.io/). Why AWS though? I had had chance to use AWS before for a few projects that I worked on, but since [I joined Deliveroo 2 years ago](https://twitter.com/tourismgeek/status/1091003681003237376), I have been using AWS everyday. So, it's the least friction for me when it comes to working with a cloud provider. That said, this migration itself still has been a really great learning experience for me, and it emphasized it more for me that AWS is million miles ahead in their journey when it comes to developer experience. Things just work, especially when it comes to gluing things together (we will see in an example of that in this post). When they don't, it's also very obvious why which makes it easy to diagnose what's going wrong (although, [it's probably because of IAM](https://nodramadevops.com/2019/11/why-is-aws-iam-so-hard/) for like 99.9% of the cases).

During this migration, I have also discovered that you can actually configure SSL on your own domain free of any additional charges through [AWS Certificate Manager](https://aws.amazon.com/certificate-manager/) (ACM) if you are already using [AWS Application Load Balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) (ALB). This was a valuable find for me, as I needed to enable HTTPS for this blog which I have been procrastinating to get it done, like forever. However, when I think about it, it wasn't only the additional payment I had to make for the SSL certificate that was making me delay getting one. It was probably the cost of maintenance which was the biggest chore that I didn't really want (e.g. certificate renewals and all that). 

ALB and ACM integration addresses these both issues, by providing a way to configure SSL as well as keeping to renewed for you free of any additional charges. To be fair, there is probably also a way to automate this all on Azure, but I have been also away from that world for over 2 years now, and I didn't have the mental capacity to sort it out. Anyway, enough with the excuses, and let's see how to make this all sorted through Terraform.

> Note that I am going to skip what AWS ALB is, how it works, and how to configure it to start directing traffic to your resources (e.g. ECS services, Lambda, EC2 instances, etc.). However, it's worth checking out [the ALB documentation](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) before this post if you don't have a good grasp of its concepts.

## Creating the Certificate Through AWS Certificate Manager

For the purpose of serving the content of this blog through HTTPS, I wanted to create an SSL certificate for `www.tugberkugurlu.com`. However, I also wanted to have the option to serve other content under subdomains. That led me to look into whether I can actually create a wildcard certificate, and this in fact turned out to be possible. As started in [the ACM characteristics docs](https://docs.aws.amazon.com/acm/latest/userguide/acm-certificate.html), ACM allows you to use an asterisk (`*`) in the domain name to create an ACM certificate containing a wildcard name that can protect several sites in the same domain.

With that information, the next step was to see how Terraform would allow me to create a wilcard certificate. Terraform AWS module already has a resource for us to create the certificate, which is called [`aws_acm_certificate`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/acm_certificate):

```
resource "aws_acm_certificate" "tugberkugurlu_com" {
  domain_name               = "tugberkugurlu.com"
  subject_alternative_names = ["*.tugberkugurlu.com"]
  validation_method         = "DNS"
}
```

Let's take a look what each of these things mean:

 - `domain_name`: Fully qualified domain name (FQDN), that you want to secure with an ACM certificate.
 - `subject_alternative_names`: Additional FQDNs to be included in the Subject Alternative Name extension of the ACM certificate. Here, we can use an asterisk (`*`) to create a wildcard certificate that protects several sites in the same domain. However, note that the asterisk (`*`) can protect only one subdomain level when you request a wildcard certificate. For example, `*.tugberkugurlu.com` can protect `foo.tugberkugurlu.com` and `bar.tugberkugurlu.com`, but it cannot protect `foo.bar.tugberkugurlu.com`. Another thing to note here is that `*.tugberkugurlu.com` protects only the subdomains of `tugberkugurlu.com`, it does not protect [the domain apex](https://help.easyredir.com/en/articles/453072-what-is-a-domain-apex) (i.e. `tugberkugurlu.com` in our case here). That's why I am providing that through the `domain_name`.
 - `validation_method`: ACM needs to validate that you actually own this domain before it can issue a public certificate. This validation can be performed through either `EMAIL` or `DNS`. I am going with `DNS` here for several reasons:
   - DNS validation is required to be eligible to renew your certificates automatically through [the managed certificate renewal](https://docs.aws.amazon.com/acm/latest/userguide/managed-renewal.html).
   - DNS validation also allows us to complete the whole process through Terraform if you use [Route53](https://aws.amazon.com/route53/) as your domain's DNS host, and managing its state through Terraform as well.

This is all we have to do to request a creation of the certificate, and like I mentioned, the certificate renewal is going to be performed automatically for us since we are creating the initial certificate with DNS validation. See [this documentation around ACM certificate renewal](Renewal for Domains Validated by DNS
) to find more about how the automatic renewal works. It's also worth mentioning that there is further configuration you can provide. So, I suggest you to check out [the Terraform documentation for [the ACM certificate resource](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/acm_certificate) to learn more about those options in case you end up needing them.

## Domain Name Validation Through Route53 DNS Configuration

Before I applied the changes, I also wanted to make sure that the domain validation side of the story is also sorted out. I use Rouet53 as my DNS service, and this made it so much easier to perform the validation. 

> To be fair, even without this, it shouldn't be too much of a hassle as DNS validation is just going to be a one-off process regardless of the approach. So, it's still low friction even you need to perform this manually.

The main data point I needed to hook into for this was [`domain_validation_options`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/acm_certificate#domain_validation_options) attribute exported through `aws_acm_certificate` resource for my certificate which I declared above. Quoting from the documentation directly, this attribute gives you the domain validation objects which can be used to complete certificate validation. Note that this can have more than one value. So, wee need to keep that in mind when we are using this. 

This is great as we can use this value to create a Route53 record through the [`aws_route53_record`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) resource. This object exports a few attributes for us:

  - `domain_name`: The domain name to be validated.
  - `resource_record_name`: The name of the DNS record to create to validate the certificate
  - `resource_record_type`: The type of DNS record to create, e.g. `A`, `CNAME`, etc.
  - `resource_record_value`: The value the DNS record needs to have.

The only issue was to figure out how to iterate over `domain_validation_options` array, and create a `aws_route53_record` resource for each. Luckly, Terraform has a way to make this work through [`for_each`](https://www.terraform.io/docs/configuration/meta-arguments/for_each.html) meta-argument, which allows us to create an instance for each item in that map or set.

Here is how my `aws_route53_record` resource declaration looked like: 

```
resource "aws_route53_record" "tugberkugurlu_com_acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.tugberkugurlu_com.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id = aws_route53_zone.tugberkugurlu_com.zone_id
  name    = each.value.name
  type    = each.value.type
  ttl     = 60
  records = [
    each.value.record,
  ]

  allow_overwrite = true
} 
```

`zone_id` here refers to the `zone_id` attribute from an [`aws_route53_zone`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone) resource which I already had declared for this domain.

Another interesting bit here is the `allow_override` argument which is used to allow creation of this record in Terraform to overwrite an existing record, if any. It turns out that [`domain_validation_options` can result in duplicate DNS records](https://stackoverflow.com/a/59745029/463785), and [this argument seems to be added just for this purpose](https://github.com/hashicorp/terraform-provider-aws/issues/13653#issuecomment-640237762).

This should on its own be enough to get the validation performed. However, one caveat here is that the validation will happen asynchronously. Therefore, your certificate might be usable right away (a.k.a. the good old eventual consistency). Terraform already has a solution for this, too through the [`aws_acm_certificate_validation`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/acm_certificate_validation) resource. This resource implements a part of the validation workflow and represents a successful validation of an ACM certificate by waiting for validation to complete. Note that this doesn't represent a real-world entity in AWS. So, changing or deleting this resource on its own has no immediate effect.

As we already have the `aws_acm_certificate` and `aws_route53_record`(s) for the validation, we can easy declare a `aws_acm_certificate_validation` resource:

```
resource "aws_acm_certificate_validation" "tugberkugurlu_com" {
  certificate_arn         = aws_acm_certificate.tugberkugurlu_com.arn
  validation_record_fqdns = [for record in aws_route53_record.tugberkugurlu_com_acm_validation : record.fqdn]
} 
```

As there can be multiple `aws_route53_record.tugberkugurlu_com_acm_validation` resources, we make use of the Terraform [`for` expression](https://www.terraform.io/docs/configuration/expressions/for.html) to assign `validation_record_fqdns` argument.

We need everything we need for the certificate and its validation. Once I executed the [`terraform apply`](https://www.terraform.io/docs/commands/apply.html) command, the certificate was created and it was all ready to use.

![](https://tugberkugurlu-blog.s3.us-east-2.amazonaws.com/post-images/01ETDHHY9QJ765MSWM7MTENMF6-Screenshot-2020-12-25-at-18.00.50-resized-3.png)

## Wiring It up with Application Load Balancer

### Redirecting HTTP Traffic Through an ALB Rule

## Resources

 - [ACM FAQ](https://aws.amazon.com/certificate-manager/faqs)
 - [ACM Certificate Characteristics](https://docs.aws.amazon.com/acm/latest/userguide/acm-certificate.html)
 - [AWS ACM RequestCertificate API Reference](https://docs.aws.amazon.com/acm/latest/APIReference/API_RequestCertificate.html)
 - [Managed Renewal for ACM Certificates](https://docs.aws.amazon.com/acm/latest/userguide/managed-renewal.html)