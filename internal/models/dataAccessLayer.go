package models // github.com/growlog/mikapod-iam/internal/models

import (
    "fmt"
    "log"

    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)


type DataAccessLayer struct {
	db *sqlx.DB
}


/**
 *  Function initializes our connection with the `postgres` database for this
 *  web-application and saves the db connection instance in a global variable.
 */
func InitDataAccessLayer(dbHost, dbPort, dbUser, dbPassword, dbName string) (*DataAccessLayer) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+ "password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
    dbInstance, err := sqlx.Connect("postgres", psqlInfo)
    if err != nil {
       log.Fatalln(err)
    }
    err = dbInstance.Ping()
    if err != nil {
        panic(err)
    }
    fmt.Println("Database successfully connected!")
    return &DataAccessLayer{
        db: dbInstance,
    }
}

func (dal *DataAccessLayer) Shutdown() {
    dal.db.Close()
}
