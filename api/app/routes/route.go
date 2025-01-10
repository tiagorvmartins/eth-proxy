package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/tiagorvmartins/eth-proxy/api/app/controllers"
	"github.com/tiagorvmartins/eth-proxy/api/app/middlewares"
)

// Function to setup routers and router groups
func SetupRouters(app *gin.Engine) {
	app.POST("/:token", middlewares.Token(), controllers.BasePath)
	app.GET("/:token/*eth", middlewares.Token(), controllers.BasePath)
}
