package order

import (
	"net/http"
	"strconv"

	sharedErr "project-structure/pkg/shared/errors"
	"project-structure/pkg/shared/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	validate *validator.Validate
	logger   logger.Logger
	service  *Service
}

func NewHandler(validate *validator.Validate, logger logger.Logger, service *Service) *Handler {
	return &Handler{validate: validate, logger: logger, service: service}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	order, err := h.service.CreateOrder(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create order", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to create order")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Order created successfully", "data": order})
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.service.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Order not found", "error", err)
		appErr := sharedErr.NewNotFound("Order not found")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Order fetched successfully", "data": order})
}

func (h *Handler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	order, err := h.service.UpdateOrder(c.Request.Context(), id, req)
	if err != nil {
		h.logger.Error("Failed to update order", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to update order")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Order updated successfully", "data": order})
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteOrder(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete order", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to delete order")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Order deleted successfully", "data": nil})
}

func (h *Handler) ListOrders(c *gin.Context) {
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "10")

	skip, err := strconv.ParseInt(skipStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid skip parameter", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid skip parameter")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid limit parameter", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid limit parameter")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	orders, err := h.service.ListOrders(c.Request.Context(), skip, limit)
	if err != nil {
		h.logger.Error("Failed to list orders", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to list orders")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Orders fetched successfully", "data": orders})
}
