#!/bin/bash

curlie --raw --http1.1 -ivL -XPOST \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -H 'User-Agent: grpc-web-javascript/0.1' \
  -H 'Accept: application/grpc-web-text' \
  -H 'Content-Type: application/grpc-web-text' \
  -H 'X-Grpc-Web: 1' \
  'https://helloworld-grpc-web-proxy-114063878166.us-west4.run.app/helloworld.Greeter/SayHello' \
  --data 'AAAAAAYKBHRlc3Q=' -D -
