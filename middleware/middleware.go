package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/auremsinistram/go-errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

var (
	ErrAuthHeaderMissing = errors.New("authorization header is missing")
	ErrAuthFormatInvalid = errors.New("authorization format is invalid")
)

type (
	ClaimsHandlerFunc func(c *echo.Context, claims jwt.Claims)
	ErrorHandlerFunc  func(c *echo.Context, err error) error
)

type requestIDKey struct{}

func Auth(
	cookieName string,
	keyHandler jwt.Keyfunc,
	claimsHandler ClaimsHandlerFunc,
	errorHandler ErrorHandlerFunc,
	options ...jwt.ParserOption,
) echo.MiddlewareFunc {
	parser := jwt.NewParser(options...)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) (err error) {
			defer func() {
				if err != nil {
					err = errorHandler(c, err)
				} else {
					err = next(c)
				}
			}()

			token, err := cookieToken(c, cookieName)
			if err != nil {
				if !errors.Is(err, http.ErrNoCookie) {
					return errors.Wrap(err, "middleware - Auth - #1")
				}

				token, err = headerToken(c)
				if err != nil {
					return errors.Wrap(err, "middleware - Auth - #2")
				}
			}

			if keyHandler != nil {
				jwtToken, err := parser.Parse(token, keyHandler)
				if err != nil {
					return errors.Wrap(err, "middleware - Auth - #3")
				}

				if claimsHandler != nil {
					claimsHandler(c, jwtToken.Claims)
				}
			}

			return nil
		}
	}
}

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			requestID := req.Header.Get(echo.HeaderXRequestID)

			if requestID == "" {
				requestID = uuid.New().String()
			}

			ctx := context.WithValue(req.Context(), requestIDKey{}, requestID)

			c.SetRequest(req.WithContext(ctx))
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)

			return next(c)
		}
	}
}

func StaticKey(secret string) func(t *jwt.Token) (any, error) {
	return func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}

func cookieToken(c *echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", errors.Wrap(err, "middleware - cookieToken - #1")
	}

	return cookie.Value, nil
}

func headerToken(c *echo.Context) (string, error) {
	authorization := c.Request().Header.Get("Authorization")

	if authorization == "" {
		return "", ErrAuthHeaderMissing
	}

	parts := strings.SplitN(authorization, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrAuthFormatInvalid
	}

	return parts[1], nil
}
