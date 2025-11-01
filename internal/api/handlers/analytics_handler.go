package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/services"
)

type AnalyticsHandler struct {
	leaveService      *services.LeaveService
	attendanceService *services.AttendanceService
}

func NewAnalyticsHandler(
	leaveService *services.LeaveService,
	attendanceService *services.AttendanceService,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		leaveService:      leaveService,
		attendanceService: attendanceService,
	}
}

func (h *AnalyticsHandler) GetAnalyticsSummary(c *gin.Context) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0) // Default: last month

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

	leaveStats, err := h.leaveService.GetLeaveStats(startDate, endDate)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	lowAttendance, err := h.attendanceService.GetLowAttendanceStudents(75.0)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	summary := gin.H{
		"period": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"leave_statistics":        leaveStats,
		"low_attendance_students": lowAttendance,
	}

	core.SuccessResponse(c, http.StatusOK, "Analytics summary retrieved successfully", summary)
}

func (h *AnalyticsHandler) GetLeaveTypeBreakdown(c *gin.Context) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)

	stats, err := h.leaveService.GetLeaveStats(startDate, endDate)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Leave breakdown retrieved successfully", stats)
}
