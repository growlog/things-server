/**
 *  IAMClient.go
 *  The purpose of this code is to provide a background service to our
 *  application to be able to make commands with the `GrowLog IAM`
 *  web service (which is responsible for our authentication and authorization).
 */

package services


import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/growlog/rpc/protos"
)

type IAMClient struct {
	ticker *time.Ticker
	iamCon *grpc.ClientConn
	iam pb.AccountClient
}

type IAMClaims struct {
	UserId int64
	ThingId int64
	ExpiresAt *timestamp.Timestamp
}

// Function will construct the Mikapod Remote application.
func InitIAMClient(iamAddress string) (*IAMClient) {
	// Set up a direct connection to the `mikapod-iam` server.
	iamCon, err := grpc.Dial(
        iamAddress,
        grpc.WithInsecure(),
		grpc.WithTimeout(10*time.Second),
		//grpc.WithUnaryInterceptor(unaryInterceptor), // Ex. Added `UnaryInterceptor`.
    )
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Set up our protocol buffer interface.
	iam := pb.NewAccountClient(iamCon)

	return &IAMClient{
		ticker: nil,
		iamCon: iamCon,
		iam: iam,
	}
}

func (c *IAMClient) VerifyAccessToken(token string) (*IAMClaims) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.iam.VerifyAccessToken(ctx, &pb.VerifyAccessTokenRequest{
		Token: token,
	})
	if err != nil {
		log.Fatalf("could not verify token from iam, reason: %v", err)
	}

    // CASE 1 OF 2: FAILURE
	if r.Status == false {
		return nil
	}

	// CASE 2 OF 2: SUCCESS
	return &IAMClaims{
		UserId: r.UserId,
		ThingId: r.ThingId,
		ExpiresAt: r.ExpiresAt,
	}
}

func (app *IAMClient) Shutdown()  {
	// app.timer.Stop()
    app.iamCon.Close()
}
