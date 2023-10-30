# go-bloggy

![Go](https://github.com/tugberkugurlu/go-bloggy/workflows/Go/badge.svg?branch=master)

Trying out stale bot...

Markdown-driven version of [tugberkugurlu/tugberk-web](https://github.com/tugberkugurlu/tugberk-web), implemented in [Go](https://go.dev/). 

## Facts

 - This codebase structure is influenced by [golang-standards/project-layout](https://github.com/golang-standards/project-layout), 
 even if it doesn't adhere to it completely.

## Build Docker Image

```
docker build -t my-golang-app -f docker-web.dockerfile .
docker run -it --rm -p 9000:8080 --name my-running-app my-golang-app
curl http://localhost:9000
```

You can connect to the running container like below to inspect, where `web` is the name of the container:

```
docker exec -it web bash
```

## Thanks

Huge thanks to below open source projects which helped me get this site up and running more quickly 🙇🏼‍♂️

 - [Semantic-Org/Semantic-UI](https://github.com/Semantic-Org/Semantic-UI): Semantic is a UI framework designed for 
 theming.
 - [gorilla/mux](https://github.com/gorilla/mux): Package gorilla/mux implements a request router and dispatcher for 
 matching incoming requests to their respective handler.
 - [denisenkom/go-mssqldb](github.com/denisenkom/go-mssqldb): A pure Go MSSQL driver for Go's database/sql package. I 
 used this to migrate the data from my old blog which used SQL Server for data storage.
 - [go-yaml/yaml](https://github.com/go-yaml/yaml): YAML support for the Go language. I used this to parse the YAML front 
 matter block inside the markdown file of each post.
 - [github.com/gosimple/slug](https://github.com/gosimple/slug): URL-friendly slugify with multiple languages support. I 
 used this to generate tag URLs.
 - [gorilla/feeds](https://github.com/gorilla/feeds): feeds is a web feed generator library for generating RSS, Atom and 
 JSON feeds from Go applications. I used this to generate my RSS feed.

Also, big thanks to people who contributed to below content which helped me implement this using Go, and Semantic UI:

 - [Serving Static Sites with Go](https://www.alexedwards.net/blog/serving-static-sites-with-go)
 - [Docker for Go Development with Hot Reload](https://levelup.gitconnected.com/docker-for-go-development-a27141f36ba9)
 - [Using Nested Templates in Go for Efficient Web Development](https://levelup.gitconnected.com/using-go-templates-for-effective-web-development-f7df10b0e4a0)
 - [Golang parse HTML, extract all content with <body> </body> tags](https://stackoverflow.com/questions/30109061/golang-parse-html-extract-all-content-with-body-body-tags)
 - [How to compare the length of a list in html/template in golang?](https://stackoverflow.com/questions/35967109/how-to-compare-the-length-of-a-list-in-html-template-in-golang)
 - [How to use base template file for golang html/template?](https://stackoverflow.com/questions/36617949/how-to-use-base-template-file-for-golang-html-template)
 - [How to implement case insensitive URL matching using gorilla mux](https://stackoverflow.com/questions/53593618/how-to-implement-case-insensitive-url-matching-using-gorilla-mux)
 - [gorilla/mux#Static Files](https://github.com/gorilla/mux/tree/75dcda0896e109a2a22c9315bca3bb21b87b2ba5#static-files)
 - [Golang Templates Cheatsheet](https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet)
 - [Go: Is there a modulus I can use inside a template](https://stackoverflow.com/a/36369436/463785)
