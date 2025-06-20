package routes

import (
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/internal/handler"
	"sicantik-idaman/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, hlp domain.Helper) {
	handler := handler.NewHandler(&hlp)

	r.GET("/", handler.TestApi)

	api := r.Group("/api/v1")
	{
		api.POST("/login", handler.Login)
	}

	middleware := middleware.NewMiddleware(&hlp)

	apiLeaveTypes := r.Group("/api/v1/leaves/types", middleware.Auth())
	{
		apiLeaveTypes.GET("", handler.GetLeaveTypes)
	}

}
