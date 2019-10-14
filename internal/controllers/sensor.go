package controllers

import (
	"context"
	// "log"

	// tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/growlog/rpc/protos"
)

func (s *ThingServer) SetSensor(ctx context.Context, in *pb.SetSensorRequest) (*pb.SetSensorResponse, error) {
	return &pb.SetSensorResponse{
		Message: "Sensor was created",
		Status: true,
	}, nil
}
