package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"goland-boilerplate-web-service/pkg/crypto/hmac"
	"goland-boilerplate-web-service/pkg/errors"

	"github.com/rs/zerolog/log"
)

type InternalClientConfig struct {
	URL       string        `mapstructure:"url"`
	SecretKey string        `mapstructure:"secret_key"`
	AccessKey string        `mapstructure:"access_key"`
	Timeout   time.Duration `mapstructure:"timeout"`
}
type InternalClient struct {
	client *http.Client
	cfg    *InternalClientConfig
}

func (p InternalClient) doRequestWithHMacAuth(req *http.Request, any interface{}) error {
	date := time.Now().UTC().Format(http.TimeFormat)
	requestLine := fmt.Sprintf("%s %s", strings.ToLower(req.Method), req.URL.RequestURI())
	signatureBase64, headers, err := hmac.Sign(map[string]string{
		"date":           date,
		"(request-line)": requestLine,
	}, []byte(p.cfg.SecretKey))
	if err != nil {
		return err
	}

	authHeader := `SM-HMAC-SHA256 accessKey="` + p.cfg.AccessKey +
		`", headers="` + headers + `", signature="` + signatureBase64 + `"`

	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Date", date)
	req.Header.Add("Content-Type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		log.Err(err).Caller().Msgf("[DEBUG] error when send request %v", err)
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		errorBody := errors.HTTPError{}
		err1 := json.Unmarshal(body, &errorBody)
		if err1 != nil {
			return err1
		}
		errorBody.HTTPStatus = res.StatusCode
		return &errorBody
	}

	err = json.Unmarshal(body, any)
	if err != nil {
		return err
	}

	return nil
}

func (p InternalClient) doRequest(req *http.Request, any interface{}) error {
	req.Header.Add("Content-Type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		errorBody := errors.HTTPError{}
		if err := json.Unmarshal(body, &errorBody); err != nil {
			return err
		}
		errorBody.HTTPStatus = res.StatusCode
		return &errorBody
	}

	if err := json.Unmarshal(body, any); err != nil {
		return err
	}

	return nil
}
