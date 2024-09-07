package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	EmailRX = regexp.MustCompile(`[A-Za-z0-9\\._%+\\-]+@[A-Za-z0-9\\.\\-]+\\.[A-Za-z]{2,}`)
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(errors),
	}
}

func (f *Form) Required(feilds ...string) {
	for _, feild := range feilds {
		value := f.Get(feild)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(feild, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This feild is to long (max is %d)", d))
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d)", d))
	}
}
func (f *Form) EmailCheck(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		f.Errors.Add(field, "This email field is invalid")
	}
}
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(field) {
		f.Errors.Add(field, "This field is invalid")
	}

}
