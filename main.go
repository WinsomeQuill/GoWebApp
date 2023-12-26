package main

import (
	"GoWebApp/config"
	"GoWebApp/controllers"
	_ "GoWebApp/docs"
	"GoWebApp/logger"
	"GoWebApp/postgres"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
)

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	file, err := logger.OpenLogFile("./logger.log")
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	config := config.NewConfig()
	pgConnect := postgres.NewPostgresPool(&config)

	err = pgConnect.Migration()
	if err != nil {
		log.Println(err)
		return
	}

	r := gin.Default()
	r.Use(AccessControl)

	r.POST("/cart", func(c *gin.Context) {
		controllers.AddItemToCart(c, pgConnect)
	})

	r.DELETE("/cart", func(c *gin.Context) {
		controllers.RemoveItemFromCart(c, pgConnect)
	})

	r.GET("/cart", func(c *gin.Context) {
		controllers.GetCart(c, pgConnect)
	})

	r.POST("/order", func(c *gin.Context) {
		controllers.CreateOrder(c, pgConnect)
	})

	r.GET("/orders", func(c *gin.Context) {
		controllers.GetOrders(c, pgConnect)
	})

	r.POST("/order-status", func(c *gin.Context) {
		controllers.OrderChangeStatus(c, pgConnect)
	})

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run(":8080")
}

func AccessControl(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
}
