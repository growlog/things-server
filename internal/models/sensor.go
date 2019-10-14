package models  // github.com/growlog/mikapod-iam/internal/models


import (
    "database/sql"
    "fmt"
    "time"

    // "github.com/growlog/things-server/internal/utils"
)


type Sensor struct {
    Id                int64          `db:"id"`
    TenantId          int64          `db:"tenant_id"`
    ThingId           int64          `db:"thing_id"`
    TypeId            uint8          `db:"type_id"`
    CreatedAt         int64          `db:"created_at"`
}

/*
 * typeId
 * (1) Temperature
 */

// Special Thanks:
// * https://jmoiron.github.io/sqlx/
// * http://wysocki.in/golang-sqlx/

/**
 *  Function will create the `sensors` table in the database.
 */
func (dal *DataAccessLayer) CreateSensorTable(dropExistingTable bool) {
    if dropExistingTable {
        drop_stmt := "DROP TABLE sensors;"
        results, err := dal.db.Exec(drop_stmt)
        if err != nil {
            fmt.Println("Sensor Model:", results, err)
        }
    }

    // Special thanks:
    // * http://www.postgresqltutorial.com/postgresql-create-table/
    // * https://www.postgresql.org/docs/9.5/datatype.html

    stmt := `CREATE TABLE sensors (
        id bigserial PRIMARY KEY,
        tenant_id bigint NOT NULL,
        thing_id bigint NOT NULL,
        type_id INT NOT NULL,
        created_at BIGINT NOT NULL
    );`
    results, err := dal.db.Exec(stmt)
    if err != nil {
        fmt.Println("Sensor table dropped with error", results, err)
    } else {
        fmt.Println("Sensor table dropped and re-created")
    }
    return
}


/**
 *  Function will return the `sensor` struct if it exists in the database or
 *  return an error.
 */
func (dal *DataAccessLayer) GetSensorByTenantIdAndCreatedAt(tenantId int64, createdAt int64) (*Sensor, error) {
    sensor := Sensor{} // The struct which will be populated from the database.

    // DEVELOPERS NOTE:
    // (1) Lookup the sensor based on the email.
    // (2) PostgreSQL uses an enumerated $1, $2, etc bindvar syntax
    err := dal.db.Get(&sensor, "SELECT * FROM sensors WHERE tenant_id = $1 AND created_at = $2", tenantId, createdAt)

    // Handling non existing item
    if err == sql.ErrNoRows {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    return &sensor, nil
}


/**
 *  Function will create a sensor, if validation passess, and reutrns the `sensor`
 *  struct else returns the error.
 */
func (dal *DataAccessLayer) CreateSensor(tenantId int64, thingId int64, typeId int64) (*Sensor, error) {
    // Step 1: Generate SQL statement for creating a new `user` in `postgres`.
    statement := `INSERT INTO sensors (tenant_id, thing_id, type_id, created_at) VALUES ($1, $2, $3, $4)`

    now := time.Now()        // Current local time. (ex: 2009-11-10 23:00:00 +0000 UTC m=+0.000000000)
    createdAt := now.Unix()  // Number of seconds since January 1, 1970 UTC. (ex: 1257894000)

    // Step 2: Execute our SQL statement and either return our new user or
    //         our error.
    _, err := dal.db.Exec(statement, tenantId, thingId, typeId, createdAt)
    if err != nil {
        return nil, err
    }

    // Step 3:
    return dal.GetSensorByTenantIdAndCreatedAt(tenantId, createdAt)
}
