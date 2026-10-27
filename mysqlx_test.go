package mysqlx_test

import (
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/mysqlx"
)

func TestQuoteIdent(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	tests := []struct {
		ident string
		want  string
	}{
		{"", "``"},
		{"db_name", "`db_name`"},
		{"a`b", "`a``b`"},
		{"``", "``````"},
		{`a"b'c\d`, "`a\"b'c\\d`"},
	}
	for _, tc := range tests {
		t.Equal(mysqlx.QuoteIdent(tc.ident), tc.want, tc.ident)
	}
}
