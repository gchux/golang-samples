#!/bin/bash

ID_TOKEN="$(gcloud auth print-identity-token | tr -d '\n')"
ENDPOINT='https://helloworld-grpc-web-proxy-114063878166.us-west4.run.app/helloworld.Greeter/SayHello'

curl --raw -D - -iv \
  --http1.1 -XOPTIONS -H "Authorization: Bearer ${ID_TOKEN}" \
  -H 'Access-Control-Request-Method: POST' \
  -H 'Access-Control-Request-Headers: authorization,x-serverless-authorization,x-grpc-web,content-type,accept,user-agent,x-request-id,x-client-id'  \
  -H 'Origin: https//localhost' "${ENDPOINT}" \
  --next \
  --http1.1 -XPOST -H "Authorization: Bearer ${ID_TOKEN}" \
  -H 'Accept: application/grpc-web-text' \
  -H 'Content-Type: application/grpc-web-text' \
  -H 'X-Grpc-Web: 1' "${ENDPOINT}" \
  --data 'AAAAAAYKBHRlc3Q='
