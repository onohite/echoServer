package utils

import (
	"crypto/md5"
	"fmt"
)

const (
	DataLayout = "02.01.2006"
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
