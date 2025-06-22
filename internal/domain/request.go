package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReqLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReqLeaveRequest struct {
	LeaveTypeId uuid.UUID `json:"leaveTypeId"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Reason      string    `json:"reason"`
	NotifyTeam  bool      `json:"notifyTeam"`
}

type ReqUpdateLeaveRequest struct {
	LeaveTypeId *uuid.UUID `json:"leaveTypeId,omitempty"`
	StartDate   *time.Time `json:"startDate,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Reason      *string    `json:"reason,omitempty"`
	NotifyTeam  *bool      `json:"notifyTeam,omitempty"`
}

type ReqUpdateLeaveStatus struct {
	Approve      bool   `json:"approve"`
	RejectedNote string `json:"rejectedNote,omitempty"`
}

type ReqCreateLeaveReaction struct {
	LeaveRequestID uuid.UUID `json:"leaveRequestId" binding:"required"`
	Reaction       string    `json:"reaction"          binding:"required"`
}

type ReqUpdateLeaveReaction struct {
	Reaction *string `json:"reaction,omitempty"`
}
