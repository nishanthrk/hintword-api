package validator

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// Validate validates the input struct
func Validate(payload interface{}) []*fiber.Error {
	err := validate.Struct(payload)

	if err != nil {
		// Empty errors slice to store the errors
		var errorList []*fiber.Error
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			errorList = append(
				errorList,
				&fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: fmt.Sprintf("Invalid value for %v: %v", fieldName, err.ActualTag()),
				},
			)
		}
		return errorList
	}

	return nil
}

// ParseBody is helper function for parsing the body.
// Is any error occurs it will panic.
// It's just a helper function to avoid writing if condition again n again.
func ParseBody(c *fiber.Ctx, body interface{}) []*fiber.Error {

	if err := c.BodyParser(body); err != nil {
		var errorList []*fiber.Error
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: fmt.Sprintf("%v given value is not matching", err.Error()),
			},
		)

		return errorList
	}

	return nil
}

// ParseBodyAndValidate is helper function for parsing the body.
// Is any error occurs it will panic.
// It's just a helper function to avoid writing if condition again n again.
func ParseBodyAndValidate(c *fiber.Ctx, body interface{}) []*fiber.Error {

	// First We Parse
	if err := ParseBody(c, body); err != nil {
		return err
	}
	// Then We Validate
	return Validate(body)
}

// TODO CUSTOM VALIDATION RULES =========================
// Password validation rule: required,min=6,max=100
var _ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
	l := len(fl.Field().String())

	return l >= 6 && l < 100
})

func init() {
	validate = validator.New()
	// Register custom validation function for FileHeader
	_ = validate.RegisterValidation("fileUpload", func(fl validator.FieldLevel) bool {
		field := fl.Field().Interface()
		fmt.Printf("field : %v", field)
		if fileHeader, ok := field.(multipart.FileHeader); ok {
			// Implement your file upload validation logic here
			// Example: Check file extension
			allowedExtensions := []string{".jpg", ".jpeg", ".png"}
			extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
			for _, allowedExt := range allowedExtensions {
				if allowedExt == extension {
					return true
				}
			}
			return false
		}
		return false
	})
}
