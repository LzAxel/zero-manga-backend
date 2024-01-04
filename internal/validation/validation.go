package validation

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	PasswordRules = []validation.Rule{
		validation.Length(6, 32),
		validation.Match(regexp.MustCompile(`^[^\"'()+,-./:;<=>?\[\]_{|}~]+$`)),
	}

	UsernameRules = []validation.Rule{
		validation.Length(4, 32),
		validation.Match(regexp.MustCompile(`^[^\"'()+,./:;<=>?\[\]{|}~]+$`)),
	}

	EmailRules = []validation.Rule{
		validation.Required,
		is.Email,
		validation.Length(4, 40),
	}
)
