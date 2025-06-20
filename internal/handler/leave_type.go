package handler

import (
	"net/http"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/pkg/databases"

	"github.com/gin-gonic/gin"
)

func (cfg *Base) GetLeaveTypes(c *gin.Context) {
	var leaveTypes []domain.LeaveType

	if err := databases.DB.Find(&leaveTypes).Error; err != nil {
		c.JSON(http.StatusNotFound, domain.Response{StatusCode: http.StatusNotFound, Message: "NOT_FOUND"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: leaveTypes})
}
