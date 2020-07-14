# Exploring using nats

# Question
- Can I use it, to decouple the webservice (learnalist) from the blog generation (gohugo)?
- Can I share the same signature to not use nats if I dont want to.
    - I can embed the service, which is an option I do with hugo.
    - Maybe I shouldn't embed at all, as hugo is a command line tool, as stan can be too.
    - I like the flag option "render-content-via-events=true".
- Will decoupling speed up hugo rendering?
- Will all this decoupling make learnalist hugely over engineered? :P

# Run server
```sh
docker run -p 4222:4222 -p 8222:8222 nats:alpine3.10
```

## Subscribe
```sh
go run main.go subscribe
```

## Publish
```sh
go run main.go publish
```

# Run streaming server

```sh
docker run -p 4222:4222 -p 8222:8222 nats-streaming:alpine3.12 --max_age 10s
```

## Publish
```sh
go run main.go publishStreaming
```

## Subscribe
```sh
go run main.go subscriberStreaming test
```

```sh
go run main.go subscriberStreaming test1
```


# Reference
- [Example of using nats to do things on the fly]()https://github.com/mycodesmells/golang-examples/blob/master/nats/pubsub/blog-generator/main.go
- [Example of using stan](https://github.com/mycodesmells/golang-examples/blob/master/nats/streaming/watcher/main.go)
- [Streaming nat](https://github.com/nats-io/stan.go)
- [Sample k8s file for streaming](https://docs.nats.io/nats-on-kubernetes/minimal-setup)
- [Limit setting via config](https://docs.nats.io/nats-server/configuration#limits)

