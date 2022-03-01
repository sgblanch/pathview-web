package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sgblanch/pathview-web/internal/model"
)

type PathviewController struct{}

func (p *Router) pathviewGroup(g *gin.RouterGroup) {
	pathview := &PathviewController{}

	defaults := gin.H{
		"gene": gin.H{
			"name":  "gene",
			"low":   "#00ff00",
			"mid":   "#d3d3d3",
			"high":  "#ff0000",
			"bins":  10,
			"limit": []int{-1, 1},
		},
		"compound": gin.H{
			"name":  "compound",
			"low":   "#0000ff",
			"mid":   "#d3d3d3",
			"high":  "#ffff00",
			"bins":  10,
			"limit": []int{-1, 1},
		},
	}

	g.GET(".", func(c *gin.Context) {
		c.HTML(200, "pathview.go.html", defaults)
	})

	// g.POST("submit")
	g.POST("upload", pathview.upload)
	// g.POST("upload/:uuid", handler)
}

func (p *PathviewController) upload(c *gin.Context) {
	var (
		upload  *model.FileUploadRequest
		session = sessions.Default(c)
	)

	err := c.ShouldBindWith(&upload, binding.JSON)
	if err != nil {
		log.Print(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "unable to parse request"})
		return
	}

	upload.Session = session.ID()

	if v := session.Get("user"); v != nil {
		user := v.(model.User)
		upload.Owner = sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		}
	}

	if !upload.Owner.Valid && upload.Session == "" {
		log.Print("[pathview upload] user and session empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unable to associate upload to user or session"})
	}

	err = upload.StoreRequest()
	if err != nil {
		log.Print(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	err = session.Save()
	if err != nil {
		log.Printf("save session: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", upload.Token))
	c.Header("Location", fmt.Sprintf("%s/%s", c.FullPath(), upload.ID))
	c.JSON(http.StatusAccepted, gin.H{
		"message":  "upload file to requested location",
		"location": fmt.Sprintf("%s/%s", c.FullPath(), upload.ID),
	})
}
