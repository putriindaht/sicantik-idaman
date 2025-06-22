package handler

import (
	"fmt"
	"net/http"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/internal/middleware"
	"sicantik-idaman/pkg/databases"
	"sicantik-idaman/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (cfg *Base) GetMyLeaveBalance(c *gin.Context) {
	fmt.Println("Hello from GetMyLeaveRequest")
	auth := middleware.GetAuth(c)

	fmt.Println(auth)

	var balances []domain.LeaveBalance
	if err := databases.DB.
		Where("user_id = ?", auth.UserID).
		Order("created_at DESC").
		Preload("LeaveType").
		Find(&balances).Error; err != nil {

		logger.Log.Error("query leave_balance (me) failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_GET_FAILED"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: balances})

}
