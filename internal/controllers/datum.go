package controllers

import (
	"context"
	// "log"

	// tspb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/growlog/rpc/protos"
)


func (s *ThingServer) SetTimeSeriesData(ctx context.Context, in *pb.SetTimeSeriesDataRequest) (*pb.SetTimeSeriesDataResponse, error) {
	// STEP 1: Verify token.
	tokenString := ctx.Value("Token").(string)
	claims := s.RemoteAccount.VerifyAccessToken(tokenString)
	if claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token - either expired or incorrect")
	}

	// STEP 2: Confirm the token is valid for the request.


    // // Iterate through all the data that was submitted by the embedded device.
	// for _, v := range in.Data {
	// 	// s.DAL.CreateTimeSeriesDatum(v.TenantId, v.SensorId, v.Value, v.Timestamp.Seconds)
	// 	log.Printf("SAVED: %v", v) // For debugging purposes only.
	// }

	return &pb.SetTimeSeriesDataResponse{
		Message: "Data was successfully uploaded",
		Status: true,
	}, nil
}
