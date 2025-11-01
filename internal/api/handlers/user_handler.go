package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.service.GetAll(page, pageSize)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	pagination := core.CreatePaginationResponse(page, pageSize, total, users)
	core.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", pagination)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		core.ErrorResponse(c, http.StatusNotFound, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}
