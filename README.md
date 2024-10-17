# Go Hello World gRPC server & client

sources:

- https://github.com/grpc/grpc-go/tree/master/examples/helloworld
- https://github.com/grpc/grpc-web/tree/master/net/grpc/gateway/examples/helloworld

## BUILD SERVER

```sh
docker buildx build --no-cache -f Dockerfile.server -t helloworld-grpc-serverr:latest .
```

## BUILD CLI CLIENT

```sh
docker buildx build --no-cache -f Dockerfile.client -t helloworld-grpc-client:latest .
```

## BUILD WEB CLIENT

```sh
docker buildx build --no-cache -f Dockerfile.web -t helloworld-grpc-web:latest .
```

## BUILD gRPC WEB PROXY

```sh
docker buildx build --no-cache -f Dockerfile.envoy -t helloworld-grpc-web-proxy:latest .
```

The **gRPC WEB proxy** proxy requires the following environment variables:

- `GPRC_AUDIENCE`: [Cloud Run custom audience](https://cloud.google.com/run/docs/configuring/custom-audiences); usually the FQDN pointing Cloud Load Balancing.
- `GRPC_HOST`: FQDN ( `*.run.app` ) of the Cloud Run hosted gRPC server; i/e: `helloworld-grpc-server-${PROJECT_NUMBER}.${GCP_REGION}.run.app`.
- `GRPC_PORT`: TCP port where the gRPC server accepts connections; in Cloud Run this must be set to `443`.
- `GRPC_HOST_IP`: IPv4 assigned to the Cloud Load Balancing frontend.

## BUILD using [Taskfile](https://github.com/go-task/task)

> [!TIP]
> This is the simplest approach to build ALL images

```sh
task -vf docker-build
```

### CREATE x509 CERTIFICATES

```sh
task -vf certs-gen
```

## USE gRPC SERVER LOCALLY

### START THE gRPC SERVER

```sh
docker run -it --rm --network=host grpc-hello-server:latest
```

### USE THE gRPC CLI CLIENT

```sh
docker run -it --rm --network=host grpc-hello-client:latest \
  -addr='127.0.0.1:8080' -name=test -secure=false -id=local-test
```

## DEPLOY SERVER TO Cloud Run

### PUSH THE gRPC SERVER IMAGE

```sh
docker tag helloworld-grpc-server:latest ${LOCATION}-pkg.dev/${PROJECT_ID}/${REPO_NAME}/helloworld-grpc-server:latest
docker push ${LOCATION}-pkg.dev/${PROJECT_ID}/${REPO_NAME}/helloworld-grpc-server:latest
```

### PUSH THE gRPC WEB Proxy IMAGE

```sh
docker tag helloworld-grpc-web-proxy:latest ${LOCATION}-pkg.dev/${PROJECT_ID}/${REPO_NAME}/helloworld-grpc-web-proxy:latest
docker push ${LOCATION}-pkg.dev/${PROJECT_ID}/${REPO_NAME}/helloworld-grpc-web-proxy:latest
```

### DEPLOY THE IMAGES

The deployment includes 2 Cloud Run services:

- 1st one to host the gRPC server
- 2nd one to host the gRPC WEB proxy

See Cloud Run service deployment docs: https://cloud.google.com/run/docs/deploying#service

- (_optional_) set `--use-http2`: http://cloud/sdk/gcloud/reference/run/deploy#--[no-]use-http2

When using the gRPC WEB proxy, the gRPC server must be served via [Cloud Load Balancing and Serverless NEG](https://cloud.google.com/load-balancing/docs/negs/serverless-neg-concepts)

## USE gRPC CLIENT

### gRPC CLI CLIENT

```sh
docker run -it --rm grpc-hello-client:latest \
  -addr="${SERVICE_NAME}-${PROJECT_NUMBER}.${REGION}.run.app:443" -name=test \
  -token=$(gcloud auth print-identity-token | tr -d '\n') -secure=true -timeout=30 -id=remote-test
```

### gRPC WEB PROXY

```sh
gidcurl --raw --http1.1 -iv -XPOST \
  -H 'User-Agent: grpc-web-javascript/0.1' \
  -H 'Accept: application/grpc-web-text' \
  -H 'Content-Type: application/grpc-web-text' \
  -H 'X-Grpc-Web: 1' \
  --data-binary @hello-request.pb.bin \
  "https://helloworld-grpc-web-proxy-${PROJECT_NUMBER}.${GCP_REGION}.run.app/helloworld.Greeter/SayHello' \
  -D -
```

### [grpcurl](https://github.com/fullstorydev/grpcurl)

```sh
grpcurl -vv -proto helloworld.proto -d '{"name":"test"}' \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  "${SERVICE_NAME}-${PROJECT_NUMBER}.${REGION}.run.app:443" helloworld.Greeter/SayHello
```
