package helper

import (
	"fmt"
	"strings"

	"goland-boilerplate-web-service/pkg/crypto/hmac"
)

func GenerateHMACAccessToken(method, requestPath, date, secret, accessKey string) (string, error) {
	requestLine := fmt.Sprintf("%s %s", strings.ToLower(method), requestPath)
	signatureBase64, headers, err := hmac.Sign(map[string]string{
		"date":           date,
		"(request-line)": requestLine,
	}, []byte(secret))
	if err != nil {
		return "", err
	}

	authHeader := fmt.Sprintf(`SM-HMAC-SHA256 accessKey=%q,headers=%q,signature=%q`, accessKey, headers, signatureBase64)

	return authHeader, nil
}
