package utils

import (
	"crypto/md5"
	"fmt"
)

func GenerateKey(keys ...string) string {
	var key string
	for _, value := range keys {
		key += value
	}
	secKey := []byte(key)
	return fmt.Sprintf("%x", md5.Sum(secKey))
}

type ErrorStruct struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}
