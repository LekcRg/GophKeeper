package valid

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func isContainsHTTP(value any) error {
	str, ok := value.(string)
	if !ok {
		return errs.ErrValueIsNotString
	}

	if strings.Contains(str, "http://") || strings.Contains(str, "https://") {
		return nil
	}

	return errs.ErrMustContainHTTP
}

func ValidAddr(addr string) error {
	return validation.Validate(
		addr,
		validation.Required,
		is.URL,
		validation.By(isContainsHTTP),
	)
}
