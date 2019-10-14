package cmd

import (
    "errors"
    "fmt"
    "os"
    "strconv"

    "github.com/spf13/cobra"

    "github.com/growlog/things-server/internal/models"
)

func init() {
    rootCmd.AddCommand(createThingCmd)
}

var createThingCmd = &cobra.Command{
    Use:   "create_thing [FIELDS]",
    Short: "Create thing.",
    Long:  `Command used to create a GrowLog thing.`,
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 3 {
          return errors.New("requires the following fields: userId, name, tenant id")
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        // Get our user arguments.
        userIdString := args[0]
        name := args[1]
        tenantIdString := args[2]

        // Minor modifications.
        userId, _ := strconv.ParseInt(userIdString, 10, 64)
        tenantId, _ := strconv.ParseInt(tenantIdString, 10, 64)

        // Load up our `environment variables` from our operating system.
        dbHost := os.Getenv("GROWLOG_THING_DB_HOST")
        dbPort := os.Getenv("GROWLOG_THING_DB_PORT")
        dbUser := os.Getenv("GROWLOG_THING_DB_USER")
        dbPassword := os.Getenv("GROWLOG_THING_DB_PASSWORD")
        dbName := os.Getenv("GROWLOG_THING_DB_NAME")

        // Initialize and connect our database layer for the command.
        dal := models.InitDataAccessLayer(dbHost, dbPort, dbUser, dbPassword, dbName)

        thing, err := dal.CreateThing(userId, name, tenantId)
        if err != nil {
            fmt.Println("Failed creating thing!")
            fmt.Println(err)
        } else {
            fmt.Println("Thing created with ID #", thing.Id)
        }
    },
}
