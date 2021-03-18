---
id: 01ETQNQC7C3C481X3TT9W2TBDV
title: Redis Cluster Local Setup Through Docker Compose With Sharding and Replication Enabled
abstract: One of the most effective ways to explore Redis cluster is by having it running locally, and playing with it there in your own safe-to-fail environment. However, getting that setup is not too straight forward. In this post, I want to walk you through on what it takes to get Redis Cluster up and running locally through Docker Compose with sharding and replication enabled, while also exposing the metrics of the running Redis nodes through Grafana.
created_at: 2020-12-29 16:35:00.0000000 +0000 UTC
format: md
tags:
- Redis
- Databases
- Docker
slugs:
- redis-cluster-local-setup-through-docker-compose-with-sharding-and-replication-enabled
---

Just more than a week ago, I wrote about the fundamentals of [Redis Cluster](http://localhost.tugberkugurlu.com/archive/redis-cluster-benefits-of-sharding-and-how-it-works), but and I used a pre-existing local setup there. I have since collected that setup together, and put it on GitHub: [tugberkugurlu/redis-cluster](https://github.com/tugberkugurlu/redis-cluster). It gives you an easy way to get a Redis Cluster up and running locally through Docker Compose, with both sharding and replication enabled. This setup also exposes Redis metrics into a Grafana dashboard so that you can see how your setup works through various out-of-the box metrics within a pre-built dashboard.

This setup is useful for various reasons:

 - It's the cheapest and most efficient way to explore Redis Cluster. For instance, once you have this up and running, you can now wire up your application to this, and test your client setup and its interactions.
 - You can also build on top of this setup, and test some advanced scenarios locally, such as how Redis Cluster behaves when we add/remove a node (a.k.a. [resharding](https://redis.io/topics/cluster-spec#cluster-live-reconfiguration)).
 - You can use this setup on your CI environment to run integrations tests against it. It's often the case that when we mock out the data access layer, we end up seeing bugs on production that we haven't seen before. Even if mocking was tempting for me a while a go (as you can see in [one of my earlier posts](https://www.tugberkugurlu.com/archive/how-to-work-with-generic-repositories-on-asp-net-mvc-and-unit-testing-them-by-mocking)), I tend to test directly against the data storage system as much as possible now.

That said, there are some caveats with this, and it's best to mention those upfront:

 - First and foremost, this setup is put together for local development purposes. So, please refrain from using this on production unless you want to set yourself for failure 😀 .
 - As you might guess, there are already other great resources that tries to do the similar. For instance, [bitnami/bitnami-docker-redis-cluster](https://github.com/bitnami/bitnami-docker-redis-cluster) is a good one, but it seems like that particular setup is also optimized for running the setup on production, and also giving an option to run it on Kubernetes which is somewhat unnecessarily complex for my needs.
 - This is only the first iteration partially taken from a setup we have used at [Deliveroo](https://careers.deliveroo.co.uk/) in a few projects (ahem, [we are hiring](https://careers.deliveroo.co.uk/?country=any&remote=&remote=true&team=engineering-team#filter-careers), great team! Speak to me if you are interested). So, it's likely that there are rough edges and improvement areas. Shameless plug alert: I will also be talking about [one of the use cases of Redis Cluster at Deliveroo in my upcoming talk at NDC London](https://ndc-london.com/agenda/redis-cluster-for-write-intensive-workloads-0xcp/0e9ytrmxsf1) in case you are interested.

All that housekeeping is out of the way, let's dive right into it, and see the details of what it takes to set this up locally.

## TL;DR; I Just Want to Get a Local Setup Running

In that case you just want to get a local setup up and running and not interested in the details, I also got your covered:

```bash
git clone https://github.com/tugberkugurlu/redis-cluster
cd redis-cluster
docker-compose up
```

Now you have a Redis Cluster up and running, If you actually navigate to `http://localhost:3000`, you should also be able to see the metrics on the Grafana dashboard.

One way to make use of this cluster within the context of your own application is to keep this setup entirely separate, and connect your own `docker-compose` setup with the containers of this setup through the defined Docker network by defining [the pre-existing network](https://docs.docker.com/compose/networking/#use-a-pre-existing-network) on your own setup. For instance, this Redis cluster setup works within the defined `redis_cluster_network` network. Note that your app’s network is given a name based on the “project name”, which is based on the name of the directory it lives in. So, if you keep the directory name same as the repository name, this means the network name will be `redis-cluster_redis_cluster_network`. Inside the [./examples/basic-client](https://github.com/tugberkugurlu/redis-cluster/blob/master/examples/basic-client) example, you can see how this might work:

```yml
version: "3"
services:
  redis_client_1:
    image: redis:5
    command: redis-cli -c -h redis_1 cluster nodes
    networks:
      - redis-cluster_redis_cluster_network
networks:
  redis-cluster_redis_cluster_network:
    external: true
```

## Initializing the Cluster

### Enabling Sharding

### Enabling Replication

## Making Redis Cluster Work Under Docker

## Exposing Redis Metrics to Grafana

## Bring All This Together With Docker Compose

## Final Thoughts

## Resources

 - [How To Change The Default Docker Subnet](https://support.zenoss.com/hc/en-us/articles/203582809-How-to-Change-the-Default-Docker-Subnet)
 - [Communication between multiple docker-compose projects](https://stackoverflow.com/a/38089080/463785)
    - Also see "[Use a pre-existing network](https://docs.docker.com/compose/networking/#use-a-pre-existing-network)"
 - [Initializing Grafana with preconfigured dashboards](https://ops.tips/blog/initialize-grafana-with-preconfigured-dashboards/)
 - [Provisioning Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/)
 - [bitnami/redis-cluster](https://hub.docker.com/r/bitnami/redis-cluster)