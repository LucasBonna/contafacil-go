package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/handlers"
)

func UserRouters(r *gin.Engine, core *app.CoreDependencies) {
	userHandlers := handlers.NewUserHandlers(core)

	users := r.Group("/users")
	{
		users.GET("/", userHandlers.ListAllUsers())
		users.GET("/:id", userHandlers.GetUser())
		users.POST("/", userHandlers.CreateUser())
		users.PATCH("/:id", userHandlers.UpdateUser())
		users.DELETE("/:id", userHandlers.DeleteUser())
	}
}
