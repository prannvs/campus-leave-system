package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prannvs/campus-leave-system/internal/api/middleware"
	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/models"
	"github.com/prannvs/campus-leave-system/internal/services"
)

type LeaveHandler struct {
	service *services.LeaveService
}

func NewLeaveHandler(service *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{service: service}
}

func (h *LeaveHandler) ApplyLeave(c *gin.Context) {
	var req models.ApplyLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		core.ErrorResponse(c, http.StatusUnauthorized, err, nil)
		return
	}

	leave, err := h.service.ApplyLeave(userID, req)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusCreated, "Leave request submitted successfully", leave)
}

func (h *LeaveHandler) GetMyLeaves(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		core.ErrorResponse(c, http.StatusUnauthorized, err, nil)
		return
	}

	leaves, err := h.service.GetMyLeaves(userID)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Leaves retrieved successfully", leaves)
}

func (h *LeaveHandler) GetPendingLeaves(c *gin.Context) {
	leaves, err := h.service.GetPendingLeaves()
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Pending leaves retrieved successfully", leaves)
}

func (h *LeaveHandler) ApproveLeave(c *gin.Context) {
	leaveID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	var req models.ApproveLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	approverID, err := middleware.GetUserID(c)
	if err != nil {
		core.ErrorResponse(c, http.StatusUnauthorized, err, nil)
		return
	}

	status := models.LeaveStatus(req.Status)
	err = h.service.ApproveLeave(uint(leaveID), approverID, status, req.Remarks)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Leave "+req.Status+" successfully", nil)
}

func (h *LeaveHandler) DeleteLeave(c *gin.Context) {
	leaveID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	if err := h.service.Delete(uint(leaveID)); err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Leave deleted successfully", nil)
}
