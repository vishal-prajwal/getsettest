package validate

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	panPattern = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

func init() {

	// Instantiate the validator for use.
	validate = validator.New()

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Check validates the provided model against it's declared tags.
func Check(val interface{}) error {
	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Error(),
			}
			fields = append(fields, field)
		}
		return NewFieldError(fields)
	}

	return nil
}

func GetValidDOB(val string) (string, error) {
	dob := strings.ReplaceAll(val, "-", "/")
	dobDate, err := time.Parse("02/01/2006", dob)
	if err != nil {
		return "", err
	}
	diff := time.Now().Year() - dobDate.Year()
	if diff >= 100 {
		return "", fmt.Errorf("yob is invalid")
	}
	return dob, nil
}

/*
If the pan is not personal then we should return an error.
Valid Pans are whose 4th character is 'P' and length is 10.
*/
func ValidatePanNumber(panNumber string) error {
	panNumber = strings.ToUpper(panNumber)
	if !panPattern.MatchString(panNumber) {
		return errors.New("PAN number is invalid")
	}
	if panNumber[3] != 'P' {
		return errors.New("PAN number should be of Individual Type")
	}
	return nil
}
