package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"io"

	"github.com/gin-gonic/gin"
)

func PrintFullReques() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Читаем тело запроса
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading request body")
			return
		}

		// Восстанавливаем тело запроса для дальнейшей обработки
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Печатаем полученный XML в консоль
		fmt.Println("Received XML:")
		fmt.Println(string(body))

		c.Next()
	}

}