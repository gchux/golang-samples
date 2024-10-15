/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// src: https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
// ref: https://cloud.google.com/run/docs/triggering/grpc

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	defaultName = "world"
)

var (
	addr    = flag.String("addr", "localhost:8080", "the address to connect to")
	name    = flag.String("name", defaultName, "Name to greet")
	secure  = flag.Bool("secure", true, "gRPC secure transport")
	token   = flag.String("token", "", "bearer token to be set in Authorization header")
	timeout = flag.Uint("timeout", 30, "rpc deadline in seconds")
	id      = flag.String("id", "", "RPC ID added as metadata")
)

func addTransportCredentials(secure bool, opts []grpc.DialOption) []grpc.DialOption {
	if secure {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("x509 error: %v", err)
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		return append(opts, grpc.WithTransportCredentials(cred))
	}
	return append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func saveRequestProto(requestProto *pb.HelloRequest) (int, error) {
	out, _ := proto.Marshal(requestProto)
	file, err := os.OpenFile("/tmp/hello-request.pb.bin",
		os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalf("failed to write rpc proto request: %v", err)
	}
	defer func() {
		file.Sync()
		file.Close()
	}()
	return file.Write(out)
}

func withIdentityToken(ctx context.Context, idToken *string) context.Context {
	if *idToken != "" {
		return grpcMetadata.AppendToOutgoingContext(ctx,
			"Authorization", fmt.Sprintf("Bearer %s", *idToken))
	}
	return ctx
}

func withID(ctx context.Context, rpcID *string) context.Context {
	id := *rpcID
	if id == "" {
		id = uuid.New().String()
	}
	return grpcMetadata.AppendToOutgoingContext(ctx, "x-rpc-id", id)
}

func main() {
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithAuthority(*addr)}

	opts = addTransportCredentials(*secure, opts)

	conn, err := grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	rpcClient := pb.NewGreeterClient(conn)

	helloRequest := &pb.HelloRequest{Name: *name}

	var sizeOfRequestProto int = 0
	if sizeOfRequestProto, err = saveRequestProto(helloRequest); err != nil {
		log.Fatalf("failed to write request proto: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	md := grpcMetadata.Pairs("timestamp", time.Now().Format(time.RFC3339Nano))
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx = withID(ctx, id)
	ctx = withIdentityToken(ctx, token)

	md, _ = grpcMetadata.FromOutgoingContext(ctx)

	log.Printf("%s[%+v] | meta: %+v | size: %d\n",
		helloRequest.ProtoReflect().Descriptor().FullName(),
		protojson.Format(helloRequest), md, sizeOfRequestProto)

	var header, trailer metadata.MD

	startOfRPC := time.Now()
	// see: https://github.com/grpc/grpc-go/blob/v1.67.1/examples/features/metadata/client/main.go#L51
	rpcResponse, err := rpcClient.SayHello(ctx, helloRequest,
		grpc.Header(&header), grpc.Trailer(&trailer))

	log.Printf("latency: %+v\n", time.Since(startOfRPC))
	log.Printf("headers: %+v\n", header)
	log.Printf("trailers: %+v\n", trailer)

	if err != nil {
		log.Fatalf("rpc failed: %v\n", err)
	}

	log.Printf("response: [%+v]\n",
		rpcResponse.ProtoReflect().Descriptor().FullName(),
		protojson.Format(rpcResponse))
}
