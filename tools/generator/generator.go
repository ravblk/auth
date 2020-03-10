package generator

import (
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	LenToken = 32
)

var errLen = errors.New("length generation")

func generator(cnt int) (string, error) {
	b := make([]byte, cnt)
	n, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	if n != len(b) {
		return "", errLen
	}
	return fmt.Sprintf("%x", b), nil
}

func Token() (string, error) {
	return generator(LenToken)
}
