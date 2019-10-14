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
    rootCmd.AddCommand(createFloatDatumCmd)
}

var createFloatDatumCmd = &cobra.Command{
    Use:   "create_datum [FIELDS]",
    Short: "Create time-series (float) datum.",
    Long:  `Command used to create a single time series (float) datum for the particular sensor.`,
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 4 {
          return errors.New("requires the following fields: tenantId, sensorId, value, timestamp")
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        // Get our user arguments.
        tenantIdString := args[0]
        sensorIdString := args[1]
        valueString := args[2]
        timestampString := args[3]

        // Minor modifications.
        tenantId, _ := strconv.ParseInt(tenantIdString, 10, 64)
        sensorId, _ := strconv.ParseInt(sensorIdString, 10, 64)
        value, _ := strconv.ParseFloat(valueString, 64) // create `float64` primitive
        timestamp, _ := strconv.ParseInt(timestampString, 10, 64)

        // Load up our `environment variables` from our operating system.
        dbHost := os.Getenv("GROWLOG_THING_DB_HOST")
        dbPort := os.Getenv("GROWLOG_THING_DB_PORT")
        dbUser := os.Getenv("GROWLOG_THING_DB_USER")
        dbPassword := os.Getenv("GROWLOG_THING_DB_PASSWORD")
        dbName := os.Getenv("GROWLOG_THING_DB_NAME")

        // Initialize and connect our database layer for the command.
        dal := models.InitDataAccessLayer(dbHost, dbPort, dbUser, dbPassword, dbName)

        dal.CreateThingTable(false)
        dal.CreateSensorTable(false)
        dal.CreateTimeSeriesDatumTable(false)

        datum, err := dal.CreateTimeSeriesDatum(tenantId, sensorId, value, timestamp)
        if err != nil {
            fmt.Println("Failed creating time-series (float) datum!")
            fmt.Println(err)
        } else {
            fmt.Println("Time-series (float) datum created with ID #", datum.Id)
        }
    },
}
