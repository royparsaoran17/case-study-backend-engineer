// Package validator
package validator

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"

	"auth-se/internal/consts"
	"auth-se/pkg/util"
)

const (
	AlphaNumericDash     = `^[0-9a-zA-Z\-]+$`
	AlphaNumeric         = `^[0-9a-zA-Z]+$`
	Numeric              = `^[0-9]+$`
	AlphaNumericSpace    = `^[0-9a-zA-Z\s]+$`
	Alpha                = `^[a-zA-Z]+$`
	AlphaSpace           = `^[a-zA-Z\s]+$`
	AlphaDashSpace       = `^[a-zA-Z\-\s]+$`
	IndonesianPeopleName = `^[a-zA-Z\'’.,\s]+$`
	RtRw                 = `^\d{1,3}\/\d{1,3}$`
	SubDistrict          = `^[0-9a-zA-Z\-\s\(\)]+$`
	Address              = `^[A-Za-z0-9'\.\-\s\,/#_()\[\]]+$`
	Pob                  = `^[A-Za-z'\.\-\s\,/#_()\[\]]+$`
	Email                = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

// ValidateAlphaNumericDash for reference transaction id
func validateAlphaNumericDash(v string) bool {
	pattern := `^[0-9a-zA-Z\-]+$`

	rgx, err := regexp.Compile(pattern)

	if err != nil {
		return false
	}

	return rgx.MatchString(v)
}

func Regex(pattern string) func(v string) bool {
	return func(v string) bool {
		if len(v) == 0 {
			return true
		}
		rgx, err := regexp.Compile(pattern)

		if err != nil {
			return false
		}

		return rgx.MatchString(v)
	}
}

func validDOB(v string) bool {
	var (
		f  = []string{`2006-01-02`, `02-01-2006`}
		tm time.Time
	)

	if len(v) == 0 {
		return true
	}

	for i := 0; i < len(f); i++ {
		t, err := time.Parse(f[i], v)

		if err != nil && i == 0 {
			continue
		}

		if err != nil {
			return false
		}

		tm = t
		break
	}

	if tm.After(time.Now().AddDate(-15, 0, 0)) {
		return false
	}

	return true
}

func isValidEmail(whitelist []string) func(string) bool {
	return func(v string) bool {

		if len(v) == 0 {
			return true
		}

		domain := util.EmailDomain(v)

		ok := Regex(Email)(v)
		if !ok {
			return false
		}

		pref := util.SubstringBefore(v, "@")

		if domain == "gmail.com" {
			if strings.Count(pref, ".") > 2 {
				return false
			}
		}

		ok = Regex(`^[a-z0-9\._]+$`)(pref)
		if !ok {
			return false
		}

		return util.InArray(domain, whitelist)
	}
}

func validDateTime(v string) bool {
	if len(v) == 0 {
		return true
	}

	_, err := time.Parse(consts.LayoutDateTimeFormat, v)
	if err != nil {
		return false
	}

	return true
}

func validIn(in []string) func(string) bool {
	return func(v string) bool {
		return util.InArray(v, in)
	}
}

func ValidAlphaNumericDash() validation.StringRule {
	return validation.NewStringRuleWithError(
		validateAlphaNumericDash,
		validation.NewError("validation_is_alphanumeric", "must contain alpha, digits and dash only"))
}

func ValidDOB() validation.StringRule {
	return validation.NewStringRuleWithError(
		validDOB,
		validation.NewError("validation_is_dob", "must be valid date of bird YYYY-MM-DD or DD-MM-YYYY"))
}

func ValidRegex(pattern string) validation.StringRule {
	return validation.NewStringRuleWithError(
		Regex(pattern),
		validation.NewError("validation_regex", fmt.Sprintf("must be valid regex %s", util.SubstringAfter(pattern, "^"))))
}

func ValidIn(in []string) validation.StringRule {
	return validation.NewStringRuleWithError(
		validIn(in),
		validation.NewError("validation_is_in", fmt.Sprintf("must be valid in: %s", util.StringJoin(in, ",", ""))))
}

func ValidDateTime() validation.StringRule {
	return validation.NewStringRuleWithError(
		validDateTime,
		validation.NewError("validation_is_date_time", fmt.Sprintf("must be valid date fromat: YYYY-MM-DD H:m:s")))
}

func ValidateEmailExact(whitelist []string) validation.StringRule {
	return validation.NewStringRuleWithError(
		isValidEmail(whitelist),
		validation.NewError("validation_email", "must be valid email address"),
	)
}
