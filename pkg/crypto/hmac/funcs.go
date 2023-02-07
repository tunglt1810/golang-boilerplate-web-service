package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
)

func SHA256(key, data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func SHA512(key, data []byte) ([]byte, error) {
	h := hmac.New(sha512.New, key)
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Sign calcs signature and returns a header value for adding to request.
func Sign(params map[string]string, secretKey []byte) (encodedSignature, headerStr string, err error) {
	signingString := ""
	headers := make([]string, 0, len(params))
	paramKeys := make([]string, 0, len(params))
	for k := range params {
		paramKeys = append(paramKeys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(paramKeys)))
	for _, k := range paramKeys {
		signingString += fmt.Sprintf("%s:%s\n", k, params[k])
		headers = append(headers, k)
	}

	if signingString == "" {
		return "", "", errors.New("signing string is empty")
	}
	if len(signingString) > 0 {
		signingString = signingString[:len(signingString)-1]
	}
	signature, err := SHA256(secretKey, []byte(signingString))
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(signature), strings.Join(headers, " "), nil
}
