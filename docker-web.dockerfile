FROM golang:1.15.6

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

WORKDIR /go/src/app/cmd/web
EXPOSE 80
CMD ["web"]