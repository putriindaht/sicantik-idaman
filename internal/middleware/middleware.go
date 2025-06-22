package middleware

import (
	"net/http"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Middleware struct {
	Helper *domain.Helper
}

type AuthInfo struct {
	UserID uuid.UUID
	TeamID uuid.UUID
	Name   string
	Role   string
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

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {

			return []byte(cfg.Helper.JwtSecret), nil

		})

		if err != nil || !token.Valid {
			log.Error("error validation token", zap.String("err", err.Error()))
			ctx.JSON(http.StatusUnauthorized, domain.Response{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("auth", AuthInfo{claims.UserId, claims.TeamId, claims.Name, claims.Role})

		ctx.Next()

	}
}

func GetAuth(c *gin.Context) AuthInfo {
	return c.MustGet("auth").(AuthInfo)
}
