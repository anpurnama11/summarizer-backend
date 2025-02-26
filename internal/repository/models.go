package repository

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type History struct {
	ID        int       `validate:"-"`
	URL       string    `validate:"required,url"`
	Title     *string   `validate:"omitempty,min=1"`
	Content   string    `validate:"required"`
	Summary   string    `validate:"required"`
	StyleID   *int      `validate:"required"`
	Language  *string   `validate:"omitempty,iso639_1"`
	CreatedAt time.Time `validate:"-"`
	Style     *Style    `validate:"-"`
}

func (h *History) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("iso639_1", validateISO639_1)

	return validate.Struct(h)
}

// validateISO639_1 checks if the language code follows ISO 639-1 format
func validateISO639_1(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true // Allow empty string (omitempty)
	}

	code := fl.Field().String()
	// ISO 639-1 codes are exactly 2 characters long
	if len(code) != 2 {
		return false
	}

	// ISO 639-1 codes must be lowercase letters
	for _, c := range code {
		if c < 'a' || c > 'z' {
			return false
		}
	}

	return true
}

type Style struct {
	ID             int       `validate:"required"`
	Name           string    `validate:"required,min=1"`
	Description    *string   `validate:"omitempty,min=1"`
	PromptTemplate string    `validate:"required,min=1"`
	CreatedAt      time.Time `validate:"required"`
}

func (s *Style) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
