package order

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	orders := router.Group("/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("", h.ListOrders)
		orders.GET(":id", h.GetOrder)
		orders.PUT(":id", h.UpdateOrder)
		orders.DELETE(":id", h.DeleteOrder)
	}
}
