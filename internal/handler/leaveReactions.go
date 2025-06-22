package handler

import (
	"fmt"
	"net/http"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/internal/middleware"
	"sicantik-idaman/pkg/databases"
	"sicantik-idaman/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (cfg *Base) CreateLeaveReaction(c *gin.Context) {
	fmt.Println(("AddLeaveRequestReqction"))
	var (
		req domain.ReqCreateLeaveReaction
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "INVALID_VALIDATION"})
		return
	}

	auth := middleware.GetAuth(c)

	// check approved leave request
	var leave domain.LeaveRequest
	if err := databases.DB.
		Where("id = ? AND status = 'approved' AND is_deleted = false", req.LeaveRequestID).
		First(&leave).Error; err != nil {

		c.JSON(http.StatusNotFound, domain.Response{StatusCode: 404, Message: "LEAVE_NOT_FOUND_OR_NOT_APPROVED"})
		return
	}

	var exists int64
	databases.DB.
		Model(&domain.LeaveReaction{}).
		Where("user_id = ? AND leave_request_id = ? AND is_deleted = false",
			auth.UserID, req.LeaveRequestID).
		Count(&exists)
	if exists > 0 {
		c.JSON(http.StatusConflict, domain.Response{
			StatusCode: http.StatusConflict, Message: "ALREADY_REACTED",
		})
		return
	}

	// create reaction
	reaction := domain.LeaveReaction{
		UserID:         auth.UserID,
		LeaveRequestID: req.LeaveRequestID,
		Reaction:       req.Reaction,
	}

	if err := databases.DB.Create(&reaction).Error; err != nil {
		logger.Log.Error("insert leave_reaction failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_INSERT_FAILED"})
		return
	}

	c.JSON(http.StatusCreated, domain.Response{StatusCode: 201, Message: "SUCCESS", Data: reaction})
}

func (cfg *Base) GetLeaveReactions(c *gin.Context) {
	leaveID := c.Param("id")

	var reactions []domain.LeaveReaction
	if err := databases.DB.
		Where("leave_request_id = ? AND is_deleted = false", leaveID).
		Order("created_at ASC").
		Find(&reactions).Error; err != nil {

		logger.Log.Error("query leave_reaction failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_QUERY_FAILED"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: 200, Message: "SUCCESS", Data: reactions})
}

func (cfg *Base) UpdateLeaveReaction(c *gin.Context) {
	idParam := c.Param("id")
	reactionID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "INVALID_ID"})
		return
	}

	var req domain.ReqUpdateLeaveReaction
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "INVALID_VALIDATION"})
		return
	}
	if req.Reaction == nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "NO_FIELD_TO_UPDATE"})
		return
	}

	auth := middleware.GetAuth(c)

	var reaction domain.LeaveReaction
	if err := databases.DB.First(&reaction, "id = ? AND is_deleted = false", reactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, domain.Response{StatusCode: 404, Message: "REACTION_NOT_FOUND"})
		return
	}

	// Hanya pemilik reaction yg boleh edit
	if reaction.UserID != auth.UserID {
		c.JSON(http.StatusForbidden, domain.Response{StatusCode: 403, Message: "NOT_AUTHORIZED"})
		return
	}

	if err := databases.DB.
		Model(&reaction).
		Update("reaction", *req.Reaction).Error; err != nil {

		logger.Log.Error("update leave_reaction failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_UPDATE_FAILED"})
		return
	}

	reaction.Reaction = *req.Reaction
	c.JSON(http.StatusOK, domain.Response{StatusCode: 200, Message: "SUCCESS", Data: reaction})
}

func (cfg *Base) DeleteLeaveReaction(c *gin.Context) {
	idParam := c.Param("id")
	reactionID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "INVALID_ID"})
		return
	}

	auth := middleware.GetAuth(c)

	var reaction domain.LeaveReaction
	if err := databases.DB.First(&reaction, "id = ?", reactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, domain.Response{StatusCode: 404, Message: "REACTION_NOT_FOUND"})
		return
	}

	// just owner can be delete reaction
	if reaction.UserID != auth.UserID {
		c.JSON(http.StatusForbidden, domain.Response{StatusCode: 403, Message: "NOT_AUTHORIZED"})
		return
	}

	if reaction.IsDeleted {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 200, Message: "SUCCESS"})
		return
	}

	if err := databases.DB.
		Model(&reaction).
		Update("is_deleted", true).Error; err != nil {

		logger.Log.Error("delete leave_reaction failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_DELETE_FAILED"})
		return
	}

	reaction.IsDeleted = true
	c.JSON(http.StatusOK, domain.Response{StatusCode: 200, Message: "SUCCESS", Data: reaction})
}
