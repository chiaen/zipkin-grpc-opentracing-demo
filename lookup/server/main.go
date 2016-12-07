/*
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"log"
	"net"
	"time"

	pb "github.com/chiaen/zipkin-grpc-opentracing-demo/lookup"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

const (
	port = ":50053"
)

type server struct{}

func (s *server) ElementarySchool(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		log.Printf("got Parenet ctx in SayHello")
	}
	time.Sleep(time.Second * 1)
	return &pb.SearchReply{Memory: " I know you before!"}, nil
}

func (s *server) HighSchool(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		log.Printf("got Parenet ctx in SayHello")
	}
	time.Sleep(time.Second * 1)
	return &pb.SearchReply{Memory: " I know you before!"}, nil
}

func (s *server) College(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		log.Printf("got Parenet ctx in SayHello")
	}
	time.Sleep(time.Second * 1)
	return &pb.SearchReply{Memory: " I know you before!"}, nil
}

func (s *server) Work(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		log.Printf("got Parenet ctx in SayHello")
	}
	time.Sleep(time.Second * 2)
	return &pb.SearchReply{Memory: " I know you before!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Jaeger tracer can be initialized with a transport that will
	// report tracing Spans to a Zipkin backend
	transport, err := zipkin.NewHTTPTransport(
		"http://localhost:9411/api/v1/spans",
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		log.Fatalf("Cannot initialize HTTP transport: %v", err)
	}
	// create Jaeger tracer
	tracer, closer := jaeger.NewTracer(
		"lookup",
		jaeger.NewConstSampler(true), // sample all traces
		jaeger.NewRemoteReporter(transport, nil),
	)
	defer closer.Close()

	s := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	pb.RegisterLookUpServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
