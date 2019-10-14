/**
 *  RemoteAccountClient.go
 *  The purpose of this code is to provide a background service to our
 *  application to be able to make commands with the `GrowLog Account`
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

type RemoteAccountClient struct {
	ticker *time.Ticker
	remoteAccountCon *grpc.ClientConn
	remoteAccount pb.AccountClient
}

type RemoteAccountClaims struct {
	UserId int64
	ThingId int64
	ExpiresAt *timestamp.Timestamp
}

// Function will construct the GrowLog remote account application.
func InitRemoteAccountClient(remoteAccountAddress string) (*RemoteAccountClient) {
	// Set up a direct connection to the `mikapod-remoteAccount` server.
	remoteAccountCon, err := grpc.Dial(
        remoteAccountAddress,
        grpc.WithInsecure(),
		grpc.WithTimeout(10*time.Second),
		//grpc.WithUnaryInterceptor(unaryInterceptor), // Ex. Added `UnaryInterceptor`.
    )
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Set up our protocol buffer interface.
	remoteAccount := pb.NewAccountClient(remoteAccountCon)

	return &RemoteAccountClient{
		ticker: nil,
		remoteAccountCon: remoteAccountCon,
		remoteAccount: remoteAccount,
	}
}

func (c *RemoteAccountClient) VerifyAccessToken(token string) (*RemoteAccountClaims) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.remoteAccount.VerifyAccessToken(ctx, &pb.VerifyAccessTokenRequest{
		Token: token,
	})
	if err != nil {
		log.Fatalf("could not verify token from remoteAccount, reason: %v", err)
	}

    // CASE 1 OF 2: FAILURE
	if r.Status == false {
		return nil
	}

	// CASE 2 OF 2: SUCCESS
	return &RemoteAccountClaims{
		UserId: r.UserId,
		ThingId: r.ThingId,
		ExpiresAt: r.ExpiresAt,
	}
}

func (app *RemoteAccountClient) Shutdown()  {
	// app.timer.Stop()
    app.remoteAccountCon.Close()
}
