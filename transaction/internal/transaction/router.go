package transaction

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
	r.group.GET("/transaction/:id", r.handler.GetTransaction)
	r.group.POST("/transaction", r.handler.CreateTransaction)
	r.group.PUT("/transaction/:id", r.handler.UpdateTransaction)
}
