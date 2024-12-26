package middleware

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"goland-boilerplate-web-service/config"
	"goland-boilerplate-web-service/pkg/crypto/hmac"
	_errs "goland-boilerplate-web-service/pkg/errors"

	"github.com/labstack/echo/v4"
)

const (
	HeaderDate       = "Date"
	nSignatureParams = 5
)

var (
	authorizationRegex = regexp.MustCompile(`\s*SM-(.+)\s+accessKey=\"(.+)\",\s*headers=\"(.+)\",\s*signature=\"(.+)\"`)
)

// HmacAuthentication to validate hmac signature.
func HmacAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			authHeader := req.Header.Get(echo.HeaderAuthorization)
			headerDate := req.Header.Get(HeaderDate)
			requestLine := fmt.Sprintf("%s %s", strings.ToLower(req.Method), req.RequestURI)
			matches := authorizationRegex.FindAllStringSubmatch(authHeader, 1)
			if len(matches) == 0 {
				return errors.New("cannot retrieve signature params")
			}
			match := matches[0]
			if len(match) != nSignatureParams {
				return errors.New("cannot retrieve signature params")
			}

			signatureStr := match[4]
			params := make(map[string]string)
			params["date"] = headerDate
			params["(request-line)"] = requestLine

			hmac.InitInstance(&config.Config.HMACInternal)
			ok, err := hmac.Instance.ValidateSignature(signatureStr, headerDate, params, []byte(config.Config.HMACInternal.Secret))
			if err != nil {
				return _errs.ErrUnauthorized.WithDetails(err.Error())
			}
			if !ok {
				return _errs.ErrUnauthorized
			}

			return next(c)
		}
	}
}
