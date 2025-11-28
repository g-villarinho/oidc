package validation

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const (
	passwordLenght = 8
)

func registerCustomRules(v *validator.Validate) error {
	if err := v.RegisterValidation("strong_password", validateStrongPassword); err != nil {
		return err
	}

	if err := v.RegisterValidation("br_phone", validateBrazilianPhone); err != nil {
		return err
	}

	return nil
}

func registerCustomTranslations(v *validator.Validate, trans ut.Translator) error {
	if err := v.RegisterTranslation("strong_password", trans, func(ut ut.Translator) error {
		return ut.Add("strong_password", "Password must be at least 8 characters long and contain uppercase, lowercase, and numeric characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("strong_password")
		return t
	}); err != nil {
		return err
	}

	if err := v.RegisterTranslation("br_phone", trans, func(ut ut.Translator) error {
		return ut.Add("br_phone", "Invalid Brazilian phone number format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("br_phone")
		return t
	}); err != nil {
		return err
	}

	return nil
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < passwordLenght {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

func validateBrazilianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	phoneClean := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	patterns := []string{
		`^55\d{10,11}$`,
		`^\d{10,11}$`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, phoneClean); matched {
			return true
		}
	}

	return false
}
