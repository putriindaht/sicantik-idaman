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

	apiLeaveRequest := r.Group("/api/v1/leaves/requests", middleware.Auth())
	{
		apiLeaveRequest.POST("", handler.CreateLeaveRequest)
		apiLeaveRequest.GET("/me", handler.GetMyLeaveRequests)
		apiLeaveRequest.GET("/approved", handler.GetApprovedLeaves)
		apiLeaveRequest.PATCH("/:id", handler.UpdateStatusLeaveRequest) // update status for approver
		apiLeaveRequest.PUT("/:id", handler.UpdateLeaveRequest)
		apiLeaveRequest.DELETE("/:id", handler.DeleteLeaveRequest)
		apiLeaveRequest.GET("/:id/reactions", handler.GetLeaveReactions)
	}

	apiLeaveBalanceRequest := r.Group("/api/v1/leaves/balances", middleware.Auth())
	{
		apiLeaveBalanceRequest.GET("/me", handler.GetMyLeaveBalance)
	}

	apiLeaveReaction := r.Group("/api/v1/leaves/reactions", middleware.Auth())
	{
		apiLeaveReaction.POST("", handler.CreateLeaveReaction)
		apiLeaveReaction.PATCH("/:id", handler.UpdateLeaveReaction)
		apiLeaveReaction.DELETE("/:id", handler.DeleteLeaveReaction)
	}

}
