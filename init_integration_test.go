//go:build integration
// +build integration

package mysqlx_test

import (
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/getenv"
	"github.com/powerman/gotest/testinit"
)

func init() { testinit.Setup(2, setupIntegration) }

var (
	testTimeFactor = getenv.Float("GO_TEST_TIME_FACTOR", 1.0)
	testSecond     = time.Duration(float64(time.Second) * testTimeFactor)
	dsn            string
)

func setupIntegration() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	dbCfg, err := mysql.ParseDSN(os.Getenv("GO_TEST_DSN"))
	if err != nil {
		testinit.Fatal("failed to parse $GO_TEST_DSN: ", err)
	}
	dbCfg.Timeout = 3 * testSecond

	dbCfg, cleanup, err := EnsureTempDB(logger, "", dbCfg)
	if err != nil {
		testinit.Fatal(err)
	}
	testinit.Teardown(cleanup)

	dsn = dbCfg.FormatDSN()
}
