package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func NewSha256(data, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}
