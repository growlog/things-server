package internal // github.com/growlog/things-server/internal

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/growlog/rpc/protos"

    "github.com/growlog/things-server/internal/services"
	"github.com/growlog/things-server/internal/controllers"
	"github.com/growlog/things-server/internal/models"
)

type ThingsServer struct {
	thingAddress string
	dal *models.DataAccessLayer
	remoteAccount *services.RemoteAccountClient
	grpcServer *grpc.Server
}

// Function will construct the GrowLog remote account application.
func InitThingsServer(dbHost, dbPort, dbUser, dbPassword, dbName, thingAddress string, remoteAccountAddress string) (*ThingsServer) {

	// Initialize and connect our database layer for the entire application.
    dal := models.InitDataAccessLayer(dbHost, dbPort, dbUser, dbPassword, dbName)

    // Create our app's models if they haven't been created previously.
    dal.CreateThingTable(false)
	dal.CreateSensorTable(false)
	dal.CreateTimeSeriesDatumTable(false)

    // Initialize our RemoteAccount client.
	remoteAccount := services.InitRemoteAccountClient(remoteAccountAddress)

	// Create our application instance.
 	return &ThingsServer{
		thingAddress: thingAddress,
		dal: dal,
		remoteAccount: remoteAccount,
		grpcServer: nil,
	}
}

// Function will consume the main runtime loop and run the business logic
// of the thing application.
func (app *ThingsServer) RunMainRuntimeLoop() {
	// Open a TCP server to the specified localhost and environment variable
    // specified port number.
    lis, err := net.Listen("tcp", app.thingAddress)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Initialize our gRPC server using our TCP server.
    grpcServer := grpc.NewServer( grpc.UnaryInterceptor(unaryInterceptor), ) // EX: Added `UnaryInterceptor`.

    // Save reference to our application state.
    app.grpcServer = grpcServer

    // For debugging purposes only.
    log.Printf("gRPC server is running.")

    // Block the main runtime loop for accepting and processing gRPC requests.
    pb.RegisterThingServer(grpcServer, &controllers.ThingServer{
        // DEVELOPERS NOTE:
        // We want to attach to every gRPC call the following variables...
        DAL: app.dal,
		RemoteAccount: app.remoteAccount,
    })
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// Function will tell the application to stop the main runtime loop when
// the process has been finished.
func (app *ThingsServer) StopMainRuntimeLoop() {
	// Finish any RPC communication taking place at the moment before
    // shutting down the gRPC server.
    app.grpcServer.GracefulStop()

	// Shutdown our connection with our database.
	app.dal.Shutdown()
}
