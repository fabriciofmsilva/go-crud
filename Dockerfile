# FROM golang

# ADD . /go/src/wscrapper

# RUN go get golang.org/x/net/html
# RUN go get golang.org/go-redis/redis

# RUN go install wscrapper

# ENTRYPOINT /go/bin/wscrapper

FROM golang:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]

EXPOSE 8080
