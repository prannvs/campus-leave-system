package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prannvs/campus-leave-system/internal/api/middleware"
	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/models"
	"github.com/prannvs/campus-leave-system/internal/services"
)

type AttendanceHandler struct {
	service *services.AttendanceService
}

func NewAttendanceHandler(service *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	var req models.MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, "Invalid date format")
		return
	}

	markerID, err := middleware.GetUserID(c)
	if err != nil {
		core.ErrorResponse(c, http.StatusUnauthorized, err, nil)
		return
	}

	err = h.service.MarkAttendance(req.StudentID, date, req.Present, markerID)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusCreated, "Attendance marked successfully", nil)
}

func (h *AttendanceHandler) GetAttendanceStats(c *gin.Context) {
	var studentID uint
	studentIDParam := c.Query("student_id")

	if studentIDParam != "" {
		id, err := strconv.ParseUint(studentIDParam, 10, 32)
		if err != nil {
			core.ErrorResponse(c, http.StatusBadRequest, err, "Invalid student ID")
			return
		}
		studentID = uint(id)
	} else {
		id, err := middleware.GetUserID(c)
		if err != nil {
			core.ErrorResponse(c, http.StatusUnauthorized, err, nil)
			return
		}
		studentID = id
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0) // Default: last 30 days

	if startStr := c.Query("start_date"); startStr != "" {
		parsed, err := time.Parse("2006-01-02", startStr)
		if err == nil {
			startDate = parsed
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		parsed, err := time.Parse("2006-01-02", endStr)
		if err == nil {
			endDate = parsed
		}
	}

	stats, err := h.service.GetStats(studentID, startDate, endDate)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Attendance stats retrieved successfully", stats)
}
func (h *AttendanceHandler) GetLowAttendanceStudents(c *gin.Context) {
	threshold := 75.0
	if thresholdStr := c.Query("threshold"); thresholdStr != "" {
		if t, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
			threshold = t
		}
	}

	students, err := h.service.GetLowAttendanceStudents(threshold)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Low attendance students retrieved successfully", students)
}
