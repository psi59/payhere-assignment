package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	regexpPhoneNumber = regexp.MustCompile(`01\d{8,9}`)
	regexpNumeric     = regexp.MustCompile(`\d`)
	regexpLowerCase   = regexp.MustCompile(`[a-z]`)
	regexpUpperCase   = regexp.MustCompile(`[A-Z]`)
)

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Func:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func ValidateStruct(i any) error {
	validate := validator.New()
	if err := validate.Struct(i); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func ValidatePassword(pwd string) error {
	if len(pwd) == 0 {
		return fmt.Errorf("invalid password: empty password")
	}
	if !regexpNumeric.MatchString(pwd) {
		return fmt.Errorf("invalid password: numeric not contains")
	}
	if !regexpLowerCase.MatchString(pwd) {
		return fmt.Errorf("invalid password: lower case not contains")
	}
	if !regexpUpperCase.MatchString(pwd) {
		return fmt.Errorf("invalid password: upper case not contains")
	}
	if !strings.ContainsAny(pwd, "!@#$%&*+-_=?:;,.|(){}<> ") {
		return fmt.Errorf("invalid password: symbol not contains")
	}

	return nil
}

func ValidatePhoneNumber(s string) error {
	if !regexpPhoneNumber.MatchString(s) {
		return fmt.Errorf("invalid phoneNumber: %q", s)
	}

	return nil
}
