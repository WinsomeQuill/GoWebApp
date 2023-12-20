package controllers

import (
	"GoWebApp/models"
	"GoWebApp/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOrder
// @Summary Create order.
// @Description Create order.
// @Tags order
// @Accept application/json
// @Produce json
// @Param User body models.UserDto  true  "Create Order JSON"
// @Success 200
// @Router /order [post]
func CreateOrder(c *gin.Context, pgConnect postgres.PgConnect) {
	reqBody := new(models.UserDto)
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	countInCart, err := pgConnect.GetItemCountInCart(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Unknown error!",
		})
		return
	}

	if countInCart == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Your cart is empty!",
		})
		return
	}

	err = pgConnect.InsertOrder(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// GetOrders
// @Summary Get orders.
// @Description Get orders.
// @Tags order
// @Accept application/json
// @Produce json
// @Param User body models.UserDto  true  "Get Orders JSON"
// @Success 200 {object} []models.UserOrder "ok"
// @Router /orders [post]
func GetOrders(c *gin.Context, pgConnect postgres.PgConnect) {
	reqBody := new(models.UserDto)
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	orders, err := pgConnect.GetOrders(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": orders,
	})
}

// OrderChangeStatus
// @Summary Changed status to order.
// @Description Changed status to order.
// @Tags order
// @Accept application/json
// @Produce json
// @Param OrderStatus body models.UpdateOrderStatusDto  true  "Order Change Status JSON"
// @Success 200
// @Router /order-status [post]
func OrderChangeStatus(c *gin.Context, pgConnect postgres.PgConnect) {
	reqBody := new(models.UpdateOrderStatusDto)
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	err = pgConnect.UpdateOrderStatus(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
