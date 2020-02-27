package passwords

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"unicode"

	"github.com/pkg/errors"
)

var (
	errHash   = errors.New("Generate Hash passwords")
	errLen    = errors.New("Password wrong length, minimum 8 chars, maximum 160")
	errLower  = errors.New("password must contain a lower letter")
	errUpper  = errors.New("password must contain a upper letter")
	errNumber = errors.New("password must contain a number")
)

const (
	passwordMinLen = 8
	passwordMaxLen = 160
)

func Verify(s string) error {
	var number, upper, lower bool
	cnt := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		}
		cnt++
	}
	if cnt > passwordMaxLen || cnt < passwordMinLen {
		return errLen
	}
	var errs bytes.Buffer
	if !number {
		errs.WriteString(errNumber.Error())
	}
	if !upper {
		if errs.Len() > 0 {
			errs.WriteString(", ")
		}
		errs.WriteString(errUpper.Error())
	}
	if !lower {
		if errs.Len() > 0 {
			errs.WriteString(", ")
		}
		errs.WriteString(errLower.Error())
	}
	if errs.Len() > 0 {
		return errors.New(errs.String())
	}
	return nil
}

func ValidateMD5(hash, password string) (bool, error) {
	if hash == "" {
		return false, errHash
	}
	h := md5.New()
	if _, err := h.Write([]byte(password)); err != nil {
		return false, errHash
	}
	if fmt.Sprintf("%x", h.Sum(nil)) != hash {
		return false, nil
	}
	return true, nil
}

func ToMD5(password string) (string, error) {
	h := md5.New()
	if _, err := h.Write([]byte(password)); err != nil {
		return "", errHash
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
