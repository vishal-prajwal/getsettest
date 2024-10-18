package newrelic

import (
	"fmt"

	"github.com/kataras/iris/v12/context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

const transactionContextKey = "newrelic.transaction"

func MiddleLayer(app *newrelic.Application) context.Handler {

	// Wrap accepts an existing newrelic Application and returns its Iris middleware.

	// Note that the Context's response writer's underline writer will be upgraded to the newrelic's one.

	// See `GetTransaction` to retrieve the transaction created.

	return func(ctx *context.Context) {

		txn := app.StartTransaction(fmt.Sprintf("%s %s", ctx.Path(), ctx.Method()))

		defer txn.End()

		ctx.Values().Set(transactionContextKey, txn)

		ctx.Next()
	}
}

func GetNewrelicTxn(ctx *context.Context) (*newrelic.Transaction, error) {

	val := ctx.Values().Get(transactionContextKey)

	txnObject, ok := val.(*newrelic.Transaction)
	if !ok {
		return nil, fmt.Errorf("unabel to load txn object")
	}
	return txnObject, nil
}
