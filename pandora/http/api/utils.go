package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/interactive-solutions/govalidator"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Validate(ctx context.Context, w http.ResponseWriter, r *http.Request, data interface{}) bool {
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// log.WithError(err).
		// 	Error("Failed to read request body")

		WriteInternalServerErrorResponse(w, "failed to read request body")
		return false
	}

	if err = json.Unmarshal(bodyContent, data); err != nil {
		// log.
		// 	WithError(err).
		// 	Error("Failed to parse request body")

		WriteBadRequestResponse(w, "Failed to parse request body: %s", err.Error())
		return false
	}

	valid, errorMap, err := govalidator.ValidateStruct(ctx, data)
	if err != nil {
		// log.
		// 	WithError(err).
		// 	Error("Failed to validate request data")

		WriteInternalServerErrorResponse(w, "Failed to validate request")
		return false
	}

	if !valid {
		// logger := logrus.WithFields(logrus.Fields{
		// 	"validation": reflect.TypeOf(data),
		// 	"errorMap":   errorMap,
		// })

		// if user, ok := r.Context().Value("User").(*internal.User); ok {
		// 	logger = logrus.WithField("user_id", user.Id)
		// }

		//logger.Info("form validation failure")

		WriteValidationErrorResponse(w, errorMap, nil)
		return false
	}

	return true
}

func ValidateGraphqlInput(ctx context.Context, input interface{}) error {
	valid, errorMap, err := govalidator.ValidateStruct(ctx, input)

	if err != nil {
		return err
	}

	if !valid {
		graphql.AddError(ctx, &gqlerror.Error{
			Message:    "Validation errors",
			Extensions: errorMap,
		})

		return errors.New("Validation error")
	}

	return nil
}
