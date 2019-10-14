/**
 *  THIS IS A SAMPLE FILE WHICH WORKS AND IS FOR TESTING PURPOSES.
 *  - https://github.com/kelseyhightower/grpc-hello-service
 *  - https://davidsbond.github.io/2019/06/14/creating-grpc-interceptors-in-go.html
 *
 */

package internal // github.com/growlog/things-server/internal

import (
	"context"
	"log"
	// "net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/**
 *  The purpose of this interceptor is to make sure we extract the
 *  token string from the authorization metadata and save it to our context
 *  before any function gets called in our server.
 */
func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Extract our metadata from our incoming context.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

    // Confirm an authorization metadata was included else error.
	authorization := md["authorization"]
	if len(authorization) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

    // Extract our token.
	tokenWithBearer := strings.TrimPrefix(authorization[0], "Bearer ")
	token := strings.Replace(tokenWithBearer, "bearer ", "", -1)

	// Attach our `token` string to our context so when any function runs,
	// it will have already our token string.
	ctx = context.WithValue(ctx, "Token", token)
	m, err := handler(ctx, req)
	if err != nil {
		log.Fatalf("RPC failed to serve: %v", err)
	}
	return m, err
}
