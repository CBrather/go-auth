package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	scopes := strings.Split(c.Scope, " ")

	for _, scope := range scopes {
		if scope == expectedScope {
			return true
		}
	}

	return false
}

func EnsureValidToken() func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse(os.Getenv("AUTH_DOMAIN"))
	if err != nil {
		zap.L().Fatal("Failed to parse the issuer url: %v", zap.Error(err))
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)

	if err != nil {
		zap.L().Fatal("Failed to set up the jwt validator", zap.Error(err))
	}

	errorHandler := func(res http.ResponseWriter, req *http.Request, err error) {
		zap.L().Error("An error occurred during jwt validation", zap.Error(err))

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

func RequireScope(requiredScope string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

		claims := token.CustomClaims.(*CustomClaims)
		if !claims.HasScope(requiredScope) {
			zap.L().Info(fmt.Sprintf("Returning a 403, as actor %s is missing a required scope: %s", token.RegisteredClaims.ID, requiredScope))

			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
