package database

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func GoquDateTrunc(text string, timestamp exp.Expression) exp.SQLFunctionExpression {
	return goqu.Func("date_trunc", text, timestamp)
}
