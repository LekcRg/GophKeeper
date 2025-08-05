package valid

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func IsContainsHTTP(value any) error {
	str, ok := value.(string)
	if !ok {
		return errs.ErrValueIsNotString
	}

	if strings.Contains(str, "http://") || strings.Contains(str, "https://") {
		return nil
	}

	return errs.ErrMustContainHTTP
}

func MapString(m map[string]string, rules []*validation.KeyRules) error {
	return validation.Validate(m, validation.Map(rules...))
}
