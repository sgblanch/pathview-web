package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

func ParseKey(k string, size int) ([]byte, error) {
	if k == "" {
		return nil, fmt.Errorf("key is empty")
	} else if !strings.HasPrefix(k, "base64:") {
		return nil, fmt.Errorf("key does not start with 'base64:'")
	}

	decoded, err := base64.StdEncoding.DecodeString(k[7:])
	if err != nil {
		return nil, err
	}

	if len(decoded) != size {
		return nil, fmt.Errorf("key size (%d) does not match expeted size (%d)", len(decoded), size)
	}

	return decoded, nil
}

func RandomToken(n int) (*string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	s := base64.StdEncoding.EncodeToString(b)

	return &s, nil
}
