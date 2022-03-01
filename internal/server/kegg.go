package server

import (
	"github.com/gin-gonic/gin"
)

func (p *Router) KeggGroup(g *gin.RouterGroup) {
	organism := new(OrganismController)
	g.GET("organism", organism.Search)

	pathway := new(PathwayController)
	g.GET("pathway", pathway.Search)
}
