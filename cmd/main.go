package main

import (
	"WST_lab1_server_new1/config"
	"WST_lab1_server_new1/internal/database/postgres"
	"WST_lab1_server_new1/internal/transport"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	storage, err := postgres.Init()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}

	router := gin.Default()

	transport.Init(router, storage)

	router.StaticFile("/favicon.ico", "./favicon.ico")
	err = router.Run("127.0.0.1:8094")
	if err != nil {
		fmt.Println(err)
	}

}
