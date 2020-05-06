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

## Thanks

Big thanks to people who contributed to below content which helped me implement this using Go, and Semantic UI:

 - [Serving Static Sites with Go](https://www.alexedwards.net/blog/serving-static-sites-with-go)
 - [Docker for Go Development with Hot Reload](https://levelup.gitconnected.com/docker-for-go-development-a27141f36ba9)
 - [Using Nested Templates in Go for Efficient Web Development](https://levelup.gitconnected.com/using-go-templates-for-effective-web-development-f7df10b0e4a0)
 - [Golang parse HTML, extract all content with <body> </body> tags](https://stackoverflow.com/questions/30109061/golang-parse-html-extract-all-content-with-body-body-tags)
 - [How to compare the length of a list in html/template in golang?](https://stackoverflow.com/questions/35967109/how-to-compare-the-length-of-a-list-in-html-template-in-golang)