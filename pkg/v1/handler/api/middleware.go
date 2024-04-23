package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey           = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header not specified")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %v", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		payload, err := s.tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error to verify token")))
			return
		}
		ctx.Set(authorizationPayloadKey, payload.Email)
		ctx.Next()

	}
}
