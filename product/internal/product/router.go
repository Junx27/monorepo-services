package product

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
	r.group.GET("/product/:id", r.handler.GetProductHandler)
	r.group.PUT("/product/:id", r.handler.UpdateProduct)
	r.group.POST("/product/:id/reduceStock", r.handler.ReduceStock)
	r.group.POST("/product/:id/increaseStock", r.handler.IncreaseStock)
}
