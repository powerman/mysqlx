//go:build integration
// +build integration

package mysqlx_test

import (
	"database/sql"
	"testing"

	"github.com/powerman/check"
)

func TestEnsureTempDB(tt *testing.T) {
	t := check.Must(tt)
	t.Contains(dsn, "_mysqlx?")
	db, err := sql.Open("mysql", dsn)
	t.Nil(err)
	_, err = db.Exec("CREATE TABLE a (id INT)")
	t.Nil(err)
}
