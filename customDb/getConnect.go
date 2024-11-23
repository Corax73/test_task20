package customDb

import (
	"database/sql"
	"songLibrary/customLog"
	"songLibrary/utils"

	_ "github.com/lib/pq"
)

// GetConnect receives data for the database from the environment file, if successful, returns the connection from the database.
func GetConnect() *sql.DB {
	settings := utils.GetConfFromEnvFile()
	dsnStr := GetDsnString(settings)
	if dsnStr != "" {
		db, err := sql.Open("postgres", dsnStr)
		if err == nil {
			return db
		} else {
			customLog.Logging(err)
		}
	}
	return nil
}

// CloseConnect closes the connection based on the passed DB instance, returns errors when closing.
func CloseConnect(db *sql.DB) {
	err := db.Close()
	if err != nil {
		customLog.Logging(err)
	}
}
