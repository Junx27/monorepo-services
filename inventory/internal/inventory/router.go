package inventory

import "github.com/gin-gonic/gin"

type Router struct {
	handler Handler
	group   gin.RouterGroup
}

func NewRouter(handler Handler, group gin.RouterGroup) Router {
	return Router{
		handler: handler,
		group:   group,
	}
}

func (r *Router) Register() {
	r.group.GET("/inventory/:id", r.handler.GetInventory)
	r.group.POST("/inventory", r.handler.CreateInventory)
	r.group.PUT("/inventory/:id", r.handler.UpdateInventory)
}
