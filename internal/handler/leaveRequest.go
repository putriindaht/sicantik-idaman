package handler

import (
	"fmt"
	"net/http"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/internal/middleware"
	"sicantik-idaman/pkg/databases"
	"sicantik-idaman/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (cfg *Base) CreateLeaveRequest(c *gin.Context) {
	var (
		req domain.ReqLeaveRequest
	)

	auth := middleware.GetAuth(c)
	fmt.Println(auth.UserID)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: http.StatusBadRequest, Message: "INVALID_VALIDATION"})
		return
	}

	leave := domain.LeaveRequest{
		UserID:      auth.UserID, // <- FK will now pass
		LeaveTypeID: req.LeaveTypeId,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Reason:      req.Reason,
		NotifyTeam:  req.NotifyTeam,
	}

	if err := databases.DB.Create(&leave).Error; err != nil {
		logger.Log.Error("insert leave_request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_INSERT_FAILED"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: leave})
}

func (cfg *Base) GetMyLeaveRequests(c *gin.Context) {
	fmt.Println("Hello from GetMyLeaveRequest")
	auth := middleware.GetAuth(c)
	status := c.Query("status")

	fmt.Println(auth)

	db := databases.DB.
		Table("leave_requests AS lr").
		Select("lr.*").
		Where("lr.is_deleted = ?", false).
		Preload("LeaveType").
		Order("lr.created_at DESC")

	// Filter role
	switch auth.Role {
	case "director", "hr":
	case "manager":
		db = db.Joins("JOIN users AS u ON u.id = lr.user_id").
			Where("u.team_id = ? OR lr.user_id = ?", auth.TeamID, auth.UserID)

	default: // employee
		db = db.Where("lr.user_id = ?", auth.UserID)
	}

	// query filter by status
	if status != "" {
		db = db.Where("lr.status = ?", status)
	}

	var leaves []domain.LeaveRequest
	if err := db.Find(&leaves).Error; err != nil {
		logger.Log.Error("query leave_request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			domain.Response{StatusCode: 500, Message: "DB_QUERY_FAILED"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: leaves})

}

func (cfg *Base) GetApprovedLeaves(c *gin.Context) {
	location := time.FixedZone("WIB", 7*3600) // UTC+7
	now := time.Now().In(location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)

	// Query string: ?start=2025-06-22&end=2025-07-10 (opsional)
	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	switch {
	// default today till one month
	case startStr == "" && endStr == "":
		start = today
		end = today.AddDate(0, 0, 30) // 30 days

	// error if just send one parameter
	case startStr == "" || endStr == "":
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "START_AND_END_REQUIRED"})
		return

	// if 2 parameter sent
	default:
		start, err = time.ParseInLocation("2006-01-02", startStr, location)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "INVALID_START_DATE"})
			return
		}
		end, err = time.ParseInLocation("2006-01-02", endStr, location)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "INVALID_END_DATE"})
			return
		}
		// Normalize day
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, location)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, location)
	}

	// Max range validation to get data
	const maxRange = 31 * 24 * time.Hour
	if end.Before(start) {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "END_BEFORE_START"})
		return
	}
	if end.Sub(start) > maxRange {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: 400, Message: "DATE_RANGE_TOO_LONG"})
		return
	}

	// Query
	var leaves []domain.LeaveRequest
	if err := databases.DB.
		Where("status = ? AND start_date BETWEEN ? AND ?", "approved", start, end).
		Order("start_date DESC").
		Preload("User").
		Find(&leaves).Error; err != nil {

		logger.Log.Error("query leave_request (approved) failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_QUERY_FAILED"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: leaves})
}

func (cfg *Base) UpdateLeaveRequest(c *gin.Context) {
	fmt.Println("UpdateLeaveRequest")
	auth := middleware.GetAuth(c)
	id := c.Param("id")

	var req domain.ReqUpdateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{StatusCode: http.StatusBadRequest, Message: "INVALID_VALIDATION"})
		return
	}

	var leave domain.LeaveRequest
	if err := databases.DB.
		Where("id = ? AND user_id = ?", id, auth.UserID).
		First(&leave).Error; err != nil {

		c.JSON(http.StatusNotFound, domain.Response{StatusCode: 404, Message: "LEAVE_NOT_FOUND"})
		return
	}

	if leave.Status == "approved" {
		c.JSON(http.StatusForbidden, domain.Response{StatusCode: 403, Message: "LEAVE_ALREADY_APPROVED"})
		return
	}

	updates := make(map[string]interface{})

	if req.LeaveTypeId != nil {
		updates["leave_type_id"] = *req.LeaveTypeId
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}
	if req.Reason != nil {
		updates["reason"] = *req.Reason
	}
	if req.NotifyTeam != nil {
		updates["notify_team"] = *req.NotifyTeam
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, domain.Response{
			StatusCode: http.StatusBadRequest, Message: "NO_FIELD_TO_UPDATE",
		})
		return
	}

	if err := databases.DB.Model(&leave).Updates(updates).Error; err != nil {
		logger.Log.Error("update leave_request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_INSERT_FAILED"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: leave})

}

func (cfg *Base) DeleteLeaveRequest(c *gin.Context) {
	fmt.Println("DeleteLeaveRequest")
	auth := middleware.GetAuth(c)
	id := c.Param("id")

	// check leave record
	var leave domain.LeaveRequest
	if err := databases.DB.
		Where("id = ? AND user_id = ?", id, auth.UserID).
		First(&leave).Error; err != nil {

		c.JSON(http.StatusNotFound, domain.Response{StatusCode: 404, Message: "LEAVE_NOT_FOUND"})
		return
	}

	// cannot delete record if its approved or rejected
	if leave.Status == "approved" || leave.Status == "rejected" {
		c.JSON(http.StatusForbidden, domain.Response{StatusCode: 403, Message: "LEAVE_ALREADY_APPROVED_OR_REJECTED"})
		return
	}

	// return success if the record already deleted
	if leave.IsDeleted {
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: http.StatusOK, Message: "SUCCESS",
		})
		return
	}

	leave.IsDeleted = true

	// soft delete
	if err := databases.DB.
		Model(&leave).
		Select("is_deleted").
		Update("is_deleted", true).Error; err != nil {
		logger.Log.Error("delete leave_request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{StatusCode: 500, Message: "DB_DELETE_FAILED"})
		return
	}

	c.JSON(http.StatusOK, domain.Response{StatusCode: http.StatusOK, Message: "SUCCESS", Data: leave})
}

func (cfg *Base) UpdateStatusLeaveRequest(c *gin.Context) {
	fmt.Println("UpdateStatusLeaveRequest")

	auth := middleware.GetAuth(c)
	id := c.Param("id")

	var (
		req domain.ReqUpdateLeaveStatus
	)

	// check request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "INVALID_BODY",
		})
		return
	}

	// check role
	if auth.Role != "manager" && auth.Role != "director" {
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: http.StatusForbidden,
			Message:    "ROLE_NOT_ALLOWED",
		})
		return
	}

	// check leave record
	var leave domain.LeaveRequest
	if err := databases.DB.
		Preload("User").
		Where("id = ?", id).
		First(&leave).Error; err != nil {

		c.JSON(http.StatusNotFound, domain.Response{
			StatusCode: http.StatusNotFound,
			Message:    "LEAVE_NOT_FOUND",
		})
		return
	}

	// validate leave status
	if leave.Status == "approved" || leave.Status == "rejected" {
		c.JSON(http.StatusBadRequest, domain.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "ALREADY_FINALIZED",
		})
		return
	}

	// check team
	if auth.Role == "manager" && *leave.User.TeamID != auth.TeamID {
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: http.StatusForbidden,
			Message:    "NOT_YOUR_TEAM",
		})
		return
	}

	// REQUEST APPROVED
	if req.Approve {
		// count the leave dates
		days := int(leave.EndDate.Sub(leave.StartDate).Hours()/24) + 1
		if days <= 0 {
			c.JSON(http.StatusBadRequest, domain.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "INVALID_DATE_RANGE",
			})
			return
		}

		// transaction begin
		tx := databases.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// check balance
		var balance domain.LeaveBalance
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND leave_type_id = ? AND year = ?",
				leave.UserID, leave.LeaveTypeID, leave.StartDate.Year()).
			First(&balance).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusInternalServerError, domain.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "BALANCE_NOT_FOUND",
			})
			return
		}

		availableDays := balance.TotalDays - balance.UsedDays
		if availableDays < days {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, domain.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "BALANCE_NOT_ENOUGH",
			})
			return
		}

		// Update balace
		if err := tx.Model(&balance).
			Update("used_days", gorm.Expr("used_days + ?", days)).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusInternalServerError, domain.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "UPDATE_BALANCE_FAILED",
			})
			return
		}

		// update request
		now := time.Now()
		if err := tx.Model(&leave).Updates(map[string]interface{}{
			"status":         "approved",
			"approved_by_id": auth.UserID,
			"approved_at":    now,
		}).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusInternalServerError, domain.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "UPDATE_LEAVE_FAILED",
			})
			return
		}

		leave.Status = "approved"
		leave.ApprovedByID = &auth.UserID
		leave.ApprovedAt = &now

		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, domain.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "TX_COMMIT_FAILED",
			})
			return
		}

		c.JSON(http.StatusOK, domain.Response{
			StatusCode: http.StatusOK,
			Message:    "SUCCESS",
			Data:       leave,
		})
		return
	}

	// REQUEST REJECTED
	updateMap := map[string]interface{}{
		"status":         "rejected",
		"approved_by_id": auth.UserID,
		"approved_at":    time.Now(),
		"rejected_note":  nil,
	}

	if req.RejectedNote != "" {
		updateMap["rejected_note"] = req.RejectedNote
	}

	// query
	if err := databases.DB.Model(&leave).Updates(updateMap).Error; err != nil {
		logger.Log.Error("reject leave failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "DB_UPDATE_FAILED",
		})
		return
	}

	leave.Status = "rejected"
	if req.RejectedNote != "" {
		leave.RejectedNote = &req.RejectedNote
	}

	c.JSON(http.StatusOK, domain.Response{
		StatusCode: http.StatusOK,
		Message:    "SUCCESS",
		Data:       leave,
	})
}
