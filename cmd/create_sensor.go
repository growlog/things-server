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
    rootCmd.AddCommand(createSensorCmd)
}

var createSensorCmd = &cobra.Command{
    Use:   "create_sensor [FIELDS]",
    Short: "Create sensor.",
    Long:  `Command used to create a GrowLog sensor.`,
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 3 {
          return errors.New("requires the following fields: tenantId, thingId, typeId")
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        // Get our user arguments.
        tenantIdString := args[0]
        thingIdString := args[1]
        typeIdString := args[2]

        // Minor modifications.
        tenantId, _ := strconv.ParseInt(tenantIdString, 10, 64)
        thingId, _ := strconv.ParseInt(thingIdString, 10, 64)
        typeId, _ := strconv.ParseInt(typeIdString, 10, 64)

        // Load up our `environment variables` from our operating system.
        dbHost := os.Getenv("GROWLOG_THING_DB_HOST")
        dbPort := os.Getenv("GROWLOG_THING_DB_PORT")
        dbUser := os.Getenv("GROWLOG_THING_DB_USER")
        dbPassword := os.Getenv("GROWLOG_THING_DB_PASSWORD")
        dbName := os.Getenv("GROWLOG_THING_DB_NAME")

        // Initialize and connect our database layer for the command.
        dal := models.InitDataAccessLayer(dbHost, dbPort, dbUser, dbPassword, dbName)

        thing, err := dal.CreateSensor(tenantId, thingId, typeId)
        if err != nil {
            fmt.Println("Failed creating sensor!")
            fmt.Println(err)
        } else {
            fmt.Println("Sensor created with ID #", thing.Id)
        }
    },
}
