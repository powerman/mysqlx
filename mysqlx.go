// Package mysqlx provide helpers for use with Go MySQL driver github.com/go-sql-driver/mysql.
package mysqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var reIdentUnsafe = regexp.MustCompile(`[^0-9a-zA-Z$_]`) //nolint:gochecknoglobals // Regexp.

// EnsureTempDB will drop/create new temporary db with suffix in db name
// and return config with temporary db name together with cleanup func
// which will close and drop temporary db.
//
// Suffix will be sanitized to contain only allowed symbols and joined with
// cfg.DBName using "_" as separator.
//
// Default value for suffix is caller's package import path.
func EnsureTempDB(log mysql.Logger, suffix string, cfg *mysql.Config) (tempDBCfg *mysql.Config, cleanup func(), err error) {
	if suffix == "" {
		pc, _, _, _ := runtime.Caller(1)
		suffix = runtime.FuncForPC(pc).Name()
		suffix = suffix[:strings.LastIndex(suffix, ".")]
	}

	cfg = cfg.Clone()
	prefix := cfg.DBName
	cfg.DBName = ""
	db, err := sql.Open("mysql", cfg.FormatDSN())
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

	cfg.DBName = dbName(prefix, reIdentUnsafe.ReplaceAllString(suffix, "_"))
	sqlDropDB := fmt.Sprintf("DROP DATABASE `%s`", cfg.DBName)     // XXX No escaping.
	sqlCreateDB := fmt.Sprintf("CREATE DATABASE `%s`", cfg.DBName) // XXX No escaping.
	if cfg.Collation != "" {
		sqlCreateDB = fmt.Sprintf("%s COLLATE %s", sqlCreateDB, cfg.Collation)
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
	return cfg, cleanup, nil
}

func dbName(prefix, suffix string) string {
	const maxDBNameLen = 63 // https://dev.mysql.com/doc/refman/5.7/en/identifier-length.html
	if len(prefix) > maxDBNameLen/2 {
		prefix = prefix[:maxDBNameLen/2]
	}
	pos := len(prefix) + 1 + len(suffix) - maxDBNameLen
	if pos < 0 {
		pos = 0
	}
	suffix = suffix[pos:]
	return fmt.Sprintf("%s_%s", prefix, suffix)
}
