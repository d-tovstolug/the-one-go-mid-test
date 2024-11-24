package main

import (
	"github.com/d.tovstoluh/the-one-go-mid-test/rest-api/controller"
	_ "github.com/d.tovstoluh/the-one-go-mid-test/rest-api/docs"
	"github.com/d.tovstoluh/the-one-go-mid-test/rest-api/storage"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title		Task API
// @version		1.0
// @description	Task CRUD API

// @schemes	http
// @host	localhost:8080

func main() {
	taskStorage := storage.NewMemoryTaskStorage()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	controller.NewTaskController(taskStorage).Register(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
