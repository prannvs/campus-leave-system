package models

import "errors"

var (
	ErrInvalidDateRange   = errors.New("end date cannot be before start date")
	ErrPastDate           = errors.New("start date cannot be in the past")
	ErrOverlappingLeave   = errors.New("leave request overlaps with existing leave")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrLeaveNotFound      = errors.New("leave request not found")
	ErrInvalidRole        = errors.New("invalid role for this operation")
	ErrAttendanceExists   = errors.New("attendance already marked for this date")
)
