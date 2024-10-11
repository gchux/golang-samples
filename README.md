# Go Hello World gRPC server

source code: https://github.com/grpc/grpc-go/tree/master/examples/helloworld

## BUILD SERVER

```sh
docker buildx build --no-cache -f Dockerfile.server -t grpc-hello-world-server:latest .
```

## BUILD CLIENT

```sh
docker buildx build --no-cache -f Dockerfile.client -t grpc-hello-world-client:latest .
```
