package validations

import (
	"fmt"
	"github.com/brienze1/notes-api/internal/domain/exceptions"
	"github.com/go-playground/validator/v10"
)

func getError(tag string, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("'%s' is required", field)
	case "min":
		return fmt.Sprintf("'%s' should be greater in length", field)
	case "uuid":
		return fmt.Sprintf("'%s' should be an uuid", field)
	}

	return fmt.Sprintf("unmapped error for field %s", field)
}

func getField(field string) string {
	return map[string]string{
		"Id":        "id",
		"UserId":    "user_id",
		"Name":      "name",
		"Title":     "title",
		"CreatedAt": "created_at",
		"UpdatedAt": "updated_at",
		"Content":   "content",
	}[field]
}

func Validate[T any](request *T) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(request)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		messages = append(messages, getError(err.ActualTag(), getField(err.Field())))
	}

	return exceptions.NewValidationError(messages...)
}
