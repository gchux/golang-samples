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

## USE IT

### START THE SERVER

```sh
docker run -it --rm -p 8080:8080 grpc-hello-server:latest
```

### USE THE CLIENT

```sh
docker run -it --rm grpc-hello-client:latest -addr 'host.domain.com:443' -name test
```

## DEPLOY SERVER TO Cloud Run

### PUSH THE SERVER IMAGE

```sh
docker tag grpc-hello-server:latest ${LOCATION}-pkg.dev/${PROJECT_ID}/${REPO_NAME}/${IMAGE_NAME}:${IMAGE_VERSION}
docker push ${LOCATION}-pkg.dev/${PROJECT_ID}/${REPO_NAME}/${IMAGE_NAME}:${IMAGE_VERSION}
```

### DEPLOY THE IMAGE

- https://cloud.google.com/run/docs/deploying#service
- set `--use-http2`: http://cloud/sdk/gcloud/reference/run/deploy#--[no-]use-http2

### USE THE CLIENT

```sh
docker run -it --rm grpc-hello-client:latest \
  -addr="${SERVICE_NAME}-${PROJECT_NUMBER}.${REGION}.run.app:443" -name=test \
  -token=$(gcloud auth print-identity-token | tr -d '\n') -secure=true -timeout=30 -id=test-rpc
```
