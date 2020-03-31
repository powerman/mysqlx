// Package mysqlx provide helpers for use with Go MySQL driver github.com/go-sql-driver/mysql.
package mysqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-sql-driver/mysql"
)

var reIdentUnsafe = regexp.MustCompile(`[^0-9a-zA-Z$_]`) //nolint:gochecknoglobals // Regexp.

// EnsureTempDB will drop/create new temporary db with suffix in db name
// and return config with temporary db name together with cleanup func
// which will close and drop temporary db.
//
// Recommended value for suffix is your package's import path (it'll be
// sanitized to contain only allowed symbols).
func EnsureTempDB(log mysql.Logger, suffix string, dbCfg mysql.Config) (tempDBCfg *mysql.Config, cleanup func(), err error) {
	prefix := dbCfg.DBName
	dbCfg.DBName = ""
	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, nil, err
	}
	closeDB := func() {
		if err := db.Close(); err != nil {
			log.Print("failed to close db: ", err)
		}
	}
	defer func() {
		if err != nil {
			closeDB()
		}
	}()

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	dbCfg.DBName = fmt.Sprintf("%s_%s", prefix, reIdentUnsafe.ReplaceAllString(suffix, "_"))
	sqlDropDB := fmt.Sprintf("DROP DATABASE %s", dbCfg.DBName)     // XXX No escaping.
	sqlCreateDB := fmt.Sprintf("CREATE DATABASE %s", dbCfg.DBName) // XXX No escaping.
	if dbCfg.Collation != "" {
		sqlCreateDB = fmt.Sprintf("%s COLLATE %s", sqlCreateDB, dbCfg.Collation)
	}
	if _, err = db.Exec(sqlDropDB); err != nil {
		if err2 := new(mysql.MySQLError); !(errors.As(err, &err2) && err2.Number == 1008) {
			return nil, nil, fmt.Errorf("failed to drop temporary db: %w", err)
		}
	}
	if _, err := db.Exec(sqlCreateDB); err != nil {
		return nil, nil, fmt.Errorf("failed to create temporary db: %w", err)
	}

	cleanup = func() {
		if _, err := db.Exec(sqlDropDB); err != nil {
			log.Print("failed to drop temporary db: ", err)
		}
		closeDB()
	}
	return &dbCfg, cleanup, nil
}
