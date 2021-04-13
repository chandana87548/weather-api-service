package persistence

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type JDBCPersistence interface {
	GetByLocation(location string, currentAccessTime string) string
	SetLastAccess(location string, lastAccess string, temperature string)
}

type JDBCRepository struct {
	DBClient *sql.DB
}

func (r *JDBCRepository) GetByLocation(location string, currentAccessTime string) string {
	rows, _ := r.DBClient.Query(fmt.Sprintf("SELECT temperature, last_accessed FROM WEATHER where location='%s'", location))
	var temperature string
	var lastAccessed string
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			_ = rows.Scan(&temperature, &lastAccessed)
		}
		if lastAccessed >= currentAccessTime {
			return temperature
		}
		return ""
	} else {
		return ""
	}

}

func (r *JDBCRepository) SetLastAccess(location string, lastAccess string, temperature string) {
	insertSt := fmt.Sprintf("INSERT INTO WEATHER (location,last_accessed,temperature ) VALUES ('%s', '%s', '%s') ON"+
		" CONFLICT(location) DO UPDATE SET last_accessed='%s', temperature='%s'", location, lastAccess, temperature, lastAccess, temperature)
	statement, _ := r.DBClient.Prepare(insertSt)
	_, err := statement.Exec()
	if err != nil {
		log.Error("Failed to save to repository, error ", err)
	}
}
