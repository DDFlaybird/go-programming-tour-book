package blogserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"code": "2000", "messgae": "Pong !", "msg": nil})
	})

	r.Run(":8081")
}
