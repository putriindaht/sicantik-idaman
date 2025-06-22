package handler

import (
	"net/http"
	"sicantik-idaman/internal/domain"

	"github.com/gin-gonic/gin"
)

func (cfg *Base) TestApi(c *gin.Context) {
	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "Success", Data: gin.H{"test": "Hello from sicantik-idaman"}})
}
