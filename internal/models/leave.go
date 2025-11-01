package models

import "time"

type LeaveType string
type LeaveStatus string

const (
	LeaveTypeMedical   LeaveType = "Medical"
	LeaveTypePersonal  LeaveType = "Personal"
	LeaveTypeEmergency LeaveType = "Emergency"
	LeaveTypeAcademic  LeaveType = "Academic"

	LeaveStatusPending  LeaveStatus = "pending"
	LeaveStatusApproved LeaveStatus = "approved"
	LeaveStatusRejected LeaveStatus = "rejected"
)

type LeaveRequest struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	StudentID  uint        `gorm:"index;not null" json:"student_id" binding:"required"`
	Student    User        `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	LeaveType  LeaveType   `gorm:"type:varchar(50);not null" json:"leave_type" binding:"required"`
	Reason     string      `gorm:"type:text;not null" json:"reason" binding:"required"`
	StartDate  time.Time   `gorm:"not null" json:"start_date" binding:"required"`
	EndDate    time.Time   `gorm:"not null" json:"end_date" binding:"required"`
	Status     LeaveStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ApprovedBy *uint       `gorm:"index" json:"approved_by,omitempty"`
	Approver   *User       `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Remarks    *string     `gorm:"type:text" json:"remarks,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

//checks if leave dates are valid
func (l *LeaveRequest) Validate() error {
	if l.EndDate.Before(l.StartDate) {
		return ErrInvalidDateRange
	}
	if l.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return ErrPastDate
	}
	return nil
}
