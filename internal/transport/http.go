package transport

import (
	"WST_lab1_server_new1/internal/database/postgres"
	"WST_lab1_server_new1/internal/handlers"
	"WST_lab1_server_new1/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init(httpserver *gin.Engine, storage *postgres.Storage) {
	//Восстановление после паники
	httpserver.Use(gin.Recovery())
	//Логгирование
	httpserver.Use(gin.Logger())
	//debug печать полного запроса
	httpserver.Use(middleware.PrintFullReques())

	handler := &handlers.StorageHandler{Storage: storage}
	//Подключение к БД
	httpserver.POST("/soap", handler.SOAPHandler)
}
