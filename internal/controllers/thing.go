package controllers

import (
	"context"
	// "log"

	// tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/growlog/rpc/protos"
)

func (s *ThingServer) SetThing(ctx context.Context, in *pb.SetThingRequest) (*pb.SetThingResponse, error) {
	return &pb.SetThingResponse{
		Message: "Thing was created",
		Status: true,
	}, nil
}
