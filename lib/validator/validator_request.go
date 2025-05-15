package validatorLib

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(obj interface{}) error {
	var errorMessages []string
	validateErr := validate.Struct(obj)
	if validateErr != nil {
		for _, valErr := range validateErr.(validator.ValidationErrors) {
			switch valErr.Tag() {
			case "email":
				errorMessages = append(errorMessages, "Field "+valErr.Field()+" must be email")
			case "required":
				errorMessages = append(errorMessages, "Field "+valErr.Field()+" is required")
			case "min":
				if valErr.Field() == "Password" {
					errorMessages = append(errorMessages, "Field "+valErr.Field()+" must be at least "+valErr.Param()+" characters")
				}
			case "eqfield":
				errorMessages = append(errorMessages, "Field "+valErr.Field()+" must be equal to "+valErr.Param())
			default:
				errorMessages = append(errorMessages, "Field "+valErr.Field()+" is invalid")
			}
		}
		return errors.New("Validation failed: " + strings.Join(errorMessages, ", "))
	}
	return nil
}
