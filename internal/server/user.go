package server

import "github.com/gin-gonic/gin"

type UserController struct{}

func (u UserController) Search(c *gin.Context) {

}

func UserGroup(g *gin.RouterGroup) {
	user := new(UserController)
	g.GET("/:id", user.Search)
}
