package models

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     Role   `json:"role" binding:"required"`
	Dept     string `json:"dept"`
	Hostel   string `json:"hostel"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ApplyLeaveRequest struct {
	LeaveType string `json:"leave_type" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type ApproveLeaveRequest struct {
	Status  string  `json:"status" binding:"required,oneof=approved rejected"`
	Remarks *string `json:"remarks"`
}

type MarkAttendanceRequest struct {
	StudentID uint   `json:"student_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	Present   bool   `json:"present"`
}
