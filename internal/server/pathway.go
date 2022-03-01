package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sgblanch/pathview-web/internal/model"

	"github.com/gin-gonic/gin"
)

type PathwayController struct {
	model model.Pathway
}

func (p PathwayController) Search(c *gin.Context) {
	var (
		pathways []model.Pathway
		organism int
		err      error
	)

	if v := c.Query("o"); v == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no organism specified"})
		return
	} else {
		organism, err = strconv.Atoi(strings.TrimPrefix(strings.ToUpper(v), "T"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "unable to parse organism"})
			return
		}
	}

	if q := c.Query("q"); q == "" {
		pathways, err = p.model.Default(organism)
		if err != nil {
			log.Printf("unable to retrieve pathway: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "unable to retrieve pathways"})
			return
		}
	} else {
		pathways, err = p.model.Search(q, organism)
		if err != nil {
			log.Printf("unable to retrieve pathway: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "unable to retrieve pathways"})
			return
		}
	}

	c.JSON(http.StatusOK, pathways)
}
