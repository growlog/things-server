package models  // github.com/growlog/mikapod-iam/internal/models


import (
    "database/sql"
    "fmt"

    // "github.com/growlog/things-server/internal/utils"
)


type TimeSeriesDatum struct {
    Id                int64          `db:"id"`
    TenantId          int64          `db:"tenant_id"`
    SensorId          int64          `db:"sensor_id"`
    Value             float64        `db:"value"`
    Timestamp         int64          `db:"timestamp"`
}


// Special Thanks:
// * https://jmoiron.github.io/sqlx/
// * http://wysocki.in/golang-sqlx/

/**
 *  Function will create the `data` table in the database.
 */
func (dal *DataAccessLayer) CreateTimeSeriesDatumTable(dropExistingTable bool) {
    if dropExistingTable {
        drop_stmt := "DROP TABLE data;"
        results, err := dal.db.Exec(drop_stmt)
        if err != nil {
            fmt.Println("TimeSeriesDatum table dropped with error:", results, err)
        } else {
            fmt.Println("TimeSeriesDatum table dropped and re-created")
        }
    }

    // Special thanks:
    // * http://www.postgresqltutorial.com/postgresql-create-table/
    // * https://www.postgresql.org/docs/9.5/datatype.html

    stmt := `CREATE TABLE data (
        id bigserial PRIMARY KEY,
        tenant_id bigint NOT NULL,
        sensor_id BIGINT NOT NULL,
        value FLOAT NULL,
        timestamp BIGINT NOT NULL,
        unique (tenant_id, sensor_id, timestamp)
    );`
    results, err := dal.db.Exec(stmt)
    if err != nil {
        fmt.Println("TimeSeriesDatum Model", results, err)
    }
    return
}


/**
 *  Function will return the `thing` struct if it exists in the database or
 *  return an error.
 */
func (dal *DataAccessLayer) GetTimeSeriesDatumByTenantIdAndCreatedAt(tenantId int64, timestamp int64) (*TimeSeriesDatum, error) {
    thing := TimeSeriesDatum{} // The struct which will be populated from the database.

    // DEVELOPERS NOTE:
    // (1) Lookup the thing based on the id.
    // (2) PostgreSQL uses an enumerated $1, $2, etc bindvar syntax
    err := dal.db.Get(&thing, "SELECT * FROM data WHERE timestamp = $1 AND tenant_id = $2", timestamp, tenantId)

    // Handling non existing item
    if err == sql.ErrNoRows {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    return &thing, nil
}


/**
 *  Function will create a thing, if validation passess, and reutrns the `thing`
 *  struct else returns the error.
 */
func (dal *DataAccessLayer) CreateTimeSeriesDatum(tenantId int64, sensorId int64, value float64, timestamp int64) (*TimeSeriesDatum, error) {
    // Step 1: Generate SQL statement for creating a new `user` in `postgres`.
    statement := `INSERT INTO data (tenant_id, sensor_id, value, timestamp) VALUES ($1, $2, $3, $4)`

    // Step 2: Execute our SQL statement and either return our new user or
    //         our error.
    _, err := dal.db.Exec(statement, tenantId, sensorId, value, timestamp)
    if err != nil {
        return nil, err
    }

    // Step 3:
    return dal.GetTimeSeriesDatumByTenantIdAndCreatedAt(tenantId, timestamp)
}
