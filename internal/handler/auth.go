package handler

import (
	"net/http"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/pkg/databases"
	"sicantik-idaman/pkg/logger"
	"sicantik-idaman/pkg/token"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *Base) Login(c *gin.Context) {
	var req domain.ReqLogin
	var user domain.User

	log := logger.Log

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: http.StatusBadRequest, Message: "INVALID_VALIDATION"})
		return
	}

	if err := databases.DB.Where("email = ? ", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, domain.Response{StatusCode: http.StatusNotFound, Message: "NOT_FOUND"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Info("compare password not match")
		c.JSON(http.StatusNotFound, domain.Response{StatusCode: http.StatusNotFound, Message: "NOT_FOUND"})
		return
	}

	token, err := token.GenerateToken(cfg.Helper.JwtSecret, &domain.JwtClaims{UserId: user.ID, Name: user.Name, TeamId: *user.TeamID, Role: string(user.Role)})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: gin.H{"token": token}})
}
