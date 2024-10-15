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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	defaultName = "world"
)

var (
	addr   = flag.String("addr", "localhost:8080", "the address to connect to")
	name   = flag.String("name", defaultName, "Name to greet")
	secure = flag.Bool("secure", false, "use TLS")
	token  = flag.String("token", "", "Identity Token")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption

	if *addr != "" {
		opts = append(opts, grpc.WithAuthority(*addr))
	}

	fmt.Println(*addr)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if *secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Note: On the Windows platform, use of x509.SystemCertPool() requires
		// go version 1.18 or higher.
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			os.Exit(1)
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	// Set up a connection to the server.
	// conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*token)

	helloRequest := &pb.HelloRequest{Name: *name}

	out, _ := proto.Marshal(helloRequest)
	os.WriteFile("/tmp/hello-request.bin", out, 0644)

	var header, trailer metadata.MD
	// Contact the server and print out its response.
	r, err := c.SayHello(ctx, helloRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s\n", r.GetMessage())
	log.Printf("headers: %+v\n", header)
	log.Printf("trailers: %+v\n", trailer)
}
