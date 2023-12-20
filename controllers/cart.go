package controllers

import (
	"GoWebApp/models"
	"GoWebApp/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddItemToCart
// @Summary Add Item To Cart.
// @Description Add Item To Cart.
// @Tags cart
// @Accept application/json
// @Produce json
// @Param cart body models.ItemToCartUserDto  true  "Cart JSON"
// @Success 200
// @Router /cart [post]
func AddItemToCart(c *gin.Context, pgConnect postgres.PgConnect) {
	reqBody := new(models.ItemToCartUserDto)
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	if reqBody.Item.Count <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	countTotal, err := pgConnect.GetItemCount(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Unknown error!",
		})
		return
	}

	if countTotal < reqBody.Item.Count {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Sorry, dont have that quantity! Available: %d", countTotal),
		})
		return
	}

	if countTotal == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Current item is out!",
		})
		return
	}

	err = pgConnect.InsertItemToCartUser(reqBody)
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

// RemoveItemFromCart
// @Summary Remove Item From Cart.
// @Description Remove Item From Cart.
// @Tags cart
// @Accept application/json
// @Produce json
// @Param cart body models.ItemToCartUserDto  true  "Cart JSON"
// @Success 200
// @Router /cart [delete]
func RemoveItemFromCart(c *gin.Context, pgConnect postgres.PgConnect) {
	reqBody := new(models.ItemToCartUserDto)
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	countInCart, err := pgConnect.GetItemCountInCart(&reqBody.User)
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

	if countInCart < reqBody.Item.Count {
		reqBody.Item.Count = countInCart
	}

	err = pgConnect.RemoveItemFromCartUser(reqBody)
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

// GetCart
// @Summary Get Items From Cart.
// @Description Get Items From Cart.
// @Tags cart
// @Accept application/json
// @Produce json
// @Param cart body models.UserDto  true  "Cart JSON"
// @Success 200 {object} models.UserCart "ok"
// @Router /cart [get]
func GetCart(c *gin.Context, pgConnect postgres.PgConnect) {
	reqBody := new(models.UserDto)
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	cart, err := pgConnect.GetCart(reqBody)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Invalid request!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": cart,
	})
}
