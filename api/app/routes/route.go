package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/tiagorvmartins/eth-proxy/api/app/controllers"
)

// Function to setup routers and router groups
func SetupRouters(app *gin.Engine) {
	app.POST("/", controllers.Example)
}
