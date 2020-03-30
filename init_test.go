package mysqlx

import (
	"testing"
	"time"

	"github.com/powerman/getenv"
	"github.com/powerman/gotest/testinit"
	_ "github.com/smartystreets/goconvey/convey" // get nice diff in web UI
)

func TestMain(m *testing.M) { testinit.Main(m) }

var (
	testTimeFactor = getenv.Float("GO_TEST_TIME_FACTOR", 1.0)
	testSecond     = time.Duration(float64(time.Second) * testTimeFactor)
)
