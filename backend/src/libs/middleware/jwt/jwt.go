package jwt

import (
	"bank-service/src/libs/dto"
	"bank-service/src/libs/env"
	"bank-service/src/libs/errors"
	httpUtils "bank-service/src/libs/http"
	"bank-service/src/libs/logger"
	"bank-service/src/utils/constant"
	"bank-service/src/utils/helpers"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

//IJWTMiddleware interface for jwt auth middleware, useful for mocking too
type IJWTMiddleware interface {
	Handler(next http.Handler) http.Handler
}

type jwtMiddleware struct {
	validUserTypes []string
}

//NewAuth0Middleware is a constructor for middleware struct
func NewJWTMiddleware(validUserTypes []string) IJWTMiddleware {
	return &jwtMiddleware{
		validUserTypes: validUserTypes,
	}
}

/*
JWTMiddleware receive the next handler (protected endpoint). Checks if JWT was sent in Authorization
header, verifies with the secret key.
Finally injects jwt values in context
*/
func (j *jwtMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx, err := j.getData(r)
		if err != nil {
			httpUtils.MakeErrorResponse(w, errors.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, rCtx)
	})
}

func (j *jwtMiddleware) getData(r *http.Request) (*http.Request, error) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		logger.GetInstance().Warning("Authorization not found")
		return r, errors.ErrFieldValidation("Authorization", "required", "")
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	signingKey := []byte(env.JWTSecretKey)
	myClaims := &claims{}
	_, err := jwt.ParseWithClaims(tokenString, myClaims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		logger.GetInstance().WithError(err).Warn("Parse JWT error")
		return r, err
	}

	if !helpers.StringInSlice(myClaims.UserType, j.validUserTypes) {
		return nil, errors.ErrUnauthorized
	}

	jwtContext := &dto.JWTContext{
		UserID:   myClaims.UserID,
		UserType: myClaims.UserType,
	}
	rCtx := r.WithContext(context.WithValue(r.Context(), helpers.ContextKey(constant.JWTContext), jwtContext))
	return rCtx, nil
}

//claims is a DTO for JWT claims
type claims struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}
