# grpc-boilerplate [![Go Report Card](https://goreportcard.com/badge/github.com/alexferl/grpc-boilerplate)](https://goreportcard.com/report/github.com/alexferl/grpc-boilerplate)

A small boilerplate gRPC app based on the Go [helloworld](https://github.com/grpc/grpc-go/tree/master/examples/helloworld)
with [12-factor](https://12factor.net/) and following golang-standards' [project-layout](https://github.com/golang-standards/project-layout).

### Building & Running locally
```shell script
go build ./cmd/grpcd && ./grpcd
```

### Testing
With [gRPCurl](https://github.com/fullstorydev/grpcurl):

```shell script
grpcurl -plaintext -d '{"name": "world"}' localhost:50051 helloworld.Greeter.SayHello

```

### Usage
```shell script
Usage of ./grpcd:
      --app-name string          The name of the application. (default "app")
      --bind-address ip          The IP address to listen at. (default 127.0.0.1)
      --bind-port uint           The port to listen at. (default 50051)
      --env-name string          The environment of the application. Used to load the right config file. (default "local")
      --env-var-prefix string    Used to prefix environment variables. (default "app")
      --graceful-timeout uint    Timeout for graceful shutdown. (default 30)
      --health-bind-address ip   The IP address to listen at. (default 127.0.0.1)
      --health-bind-port uint    The port to listen at. (default 50052)
      --log-level string         The granularity of log outputs. Valid log levels: 'panic', 'fatal', 'error', 'warn', 'info', 'debug' and 'trace'. (default "info")
      --log-output string        The output to write to. 'stdout' means log to stdout, 'stderr' means log to stderr. (default "stdout")
      --log-writer string        The log writer. Valid writers are: 'console' and 'json'. (default "json")
```

### Building Docker image
```shell script
docker build -t grpcd .
```
