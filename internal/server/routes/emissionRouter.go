package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/handlers"
)

func EmissionRouter(r *gin.Engine, deps *app.Dependencies) {
	emissionHandlers := handlers.NewEmissionHandlers(&deps.Core, &deps.External, &deps.Internal)
	emission := r.Group("/emission")
	{
		gnre := emission.Group("/gnre")
		{
			gnre.GET("", emissionHandlers.HandlerListEmissions())
			gnre.POST("", emissionHandlers.IssueGNREHandler())
		}
	}
}
