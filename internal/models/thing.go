package models  // github.com/growlog/mikapod-iam/internal/models


import (
    "database/sql"
    "fmt"
    "time"

    // "github.com/growlog/things-server/internal/utils"
)


type Thing struct {
    Id                int64          `db:"id"`
    TenantId          int64          `db:"tenant_id"`
    UserId            uint8          `db:"user_id"`
    Name              sql.NullString `db:"name"`
    CreatedAt         int64          `db:"created_at"`
}


// Special Thanks:
// * https://jmoiron.github.io/sqlx/
// * http://wysocki.in/golang-sqlx/

/**
 *  Function will create the `things` table in the database.
 */
func (dal *DataAccessLayer) CreateThingTable(dropExistingTable bool) {
    if dropExistingTable {
        drop_stmt := "DROP TABLE things;"
        results, err := dal.db.Exec(drop_stmt)
        if err != nil {
            fmt.Println("Thing table dropped with error:", results, err)
        } else {
            fmt.Println("Thing table dropped and re-created")
        }
    }

    // Special thanks:
    // * http://www.postgresqltutorial.com/postgresql-create-table/
    // * https://www.postgresql.org/docs/9.5/datatype.html

    stmt := `CREATE TABLE things (
        id bigserial PRIMARY KEY,
        tenant_id bigint NOT NULL,
        user_id BIGINT NOT NULL,
        name VARCHAR (50) NOT NULL,
        created_at BIGINT NOT NULL
    );`
    results, err := dal.db.Exec(stmt)
    if err != nil {
        fmt.Println("Thing Model", results, err)
    }
    return
}


/**
 *  Function will return the `thing` struct if it exists in the database or
 *  return an error.
 */
func (dal *DataAccessLayer) GetThingByTenantIdAndCreatedAt(tenantId int64, createdAt int64) (*Thing, error) {
    thing := Thing{} // The struct which will be populated from the database.

    // DEVELOPERS NOTE:
    // (1) Lookup the thing based on the id.
    // (2) PostgreSQL uses an enumerated $1, $2, etc bindvar syntax
    err := dal.db.Get(&thing, "SELECT * FROM things WHERE created_at = $1 AND tenant_id = $2", createdAt, tenantId)

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
func (dal *DataAccessLayer) CreateThing(userId int64, name string, tenantId int64) (*Thing, error) {
    // Step 1: Generate SQL statement for creating a new `user` in `postgres`.
    statement := `INSERT INTO things (user_id, name, tenant_id, created_at) VALUES ($1, $2, $3, $4)`

    now := time.Now()        // Current local time. (ex: 2009-11-10 23:00:00 +0000 UTC m=+0.000000000)
    createdAt := now.Unix()  // Number of seconds since January 1, 1970 UTC. (ex: 1257894000)

    // Step 2: Execute our SQL statement and either return our new user or
    //         our error.
    _, err := dal.db.Exec(statement, userId, name, tenantId, createdAt)
    if err != nil {
        return nil, err
    }

    // Step 3:
    return dal.GetThingByTenantIdAndCreatedAt(tenantId, createdAt)
}
