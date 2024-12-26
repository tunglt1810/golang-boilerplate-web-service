package hmac

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"
)

type HMAC struct {
	cfg *Config
}

// New returns a new HMAC.
func New(cfg *Config) *HMAC {
	return &HMAC{
		cfg: cfg,
	}
}

func (h *HMAC) validateClockSkew(dateStr string) (bool, error) {
	date, err := http.ParseTime(dateStr)
	if err != nil {
		fmt.Println("ParseTime error: ", err)
		return false, err
	}

	now := time.Now().UTC()
	if math.Abs(float64(now.Sub(date))) > float64(time.Duration(h.cfg.AllowedClockSkew)*time.Second) {
		return false, nil
	}

	return true, nil
}

func (h *HMAC) ValidateSignature(signatureStr, dateStr string, params map[string]string, secret []byte) (bool, error) {
	clockSkewValid, err := h.validateClockSkew(dateStr)
	if err != nil {
		return false, err
	}

	if !clockSkewValid {
		return false, errors.New("signature cannot be verify")
	}

	signingString := ""
	paramKeys := make([]string, 0, len(params))
	for k := range params {
		paramKeys = append(paramKeys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(paramKeys)))
	for _, k := range paramKeys {
		signingString += fmt.Sprintf("%s:%s\n", k, params[k])
	}

	if signingString != "" {
		signingString = signingString[:len(signingString)-1]
	}

	signature, err := SHA256(secret, []byte(signingString))
	if err != nil {
		return false, err
	}
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	return signatureStr == signatureBase64, nil
}
