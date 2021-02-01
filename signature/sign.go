package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func SignPayload(payload string, secret string) (signed string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	var sig = hmac.New(sha256.New, decoded)
	_, err = sig.Write([]byte(payload))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig.Sum(nil)), err
}
