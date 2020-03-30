// +build integration

package mysqlx

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/gotest/testinit"
)

const testDBSuffix = "github.com/powerman/mysqlx"

func init() { testinit.Setup(2, setupIntegration) }

var dsn string

func setupIntegration() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	dbCfg, err := mysql.ParseDSN(os.Getenv("GO_TEST_DSN"))
	if err != nil {
		testinit.Fatal("failed to parse $GO_TEST_DSN: ", err)
	}
	dbCfg.Timeout = 3 * testSecond

	dbCfg, cleanup, err := EnsureTempDB(logger, testDBSuffix, *dbCfg)
	if err != nil {
		testinit.Fatal(err)
	}
	testinit.Teardown(cleanup)

	dsn = dbCfg.FormatDSN()
}
