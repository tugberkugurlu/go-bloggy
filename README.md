# go-bloggy
Markdown-driven version of tugberkugurlu/Bloggy, implemented in Go.

## Facts

 - This codebase structured is influenced by [golang-standards/project-layout](https://github.com/golang-standards/project-layout), 
 even if it doesn't adhere to it completely.

## Build Docker Image

```
docker build -t my-golang-app -f docker-web.dockerfile .
docker run -it --rm -p 9000:8080 --name my-running-app my-golang-app
curl http://localhost:9000/monkeys
```

You can connect to the running container like below to inspect, where `web` is the name of the container:

```
docker exec -it web bash
```