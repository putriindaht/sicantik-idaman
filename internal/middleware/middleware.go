package middleware

import (
	"fmt"
	"net/http"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Middleware struct {
	Helper *domain.Helper
}

func NewMiddleware(hlp *domain.Helper) Middleware {
	return Middleware{Helper: hlp}
}

func (cfg *Middleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		log := logger.Log

		tokenStr := ctx.GetHeader("Authorization")
		tokenStr = strings.TrimSpace(strings.TrimPrefix(tokenStr, "Bearer "))

		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, domain.Response{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			ctx.Abort()
			return
		}

		claims := &domain.JwtClaims{}

		fmt.Println("token str : ", tokenStr)

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {

			return []byte(cfg.Helper.JwtSecret), nil

		})

		fmt.Println(token, "token")

		if err != nil || !token.Valid {
			log.Error("error validation token", zap.String("err", err.Error()))
			ctx.JSON(http.StatusUnauthorized, domain.Response{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			ctx.Abort()
			return
		}

		log.Info("UserId", zap.Any("userId", claims.UserId))
		ctx.Set("UserId", claims.UserId)
		ctx.Next()

	}
}
