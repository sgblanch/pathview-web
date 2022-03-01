package server

import (
	"log"
	"net/http"

	"github.com/sgblanch/pathview-web/internal/model"

	"github.com/gin-gonic/gin"
)

type OrganismController struct {
	model model.Organism
}

func (p OrganismController) Search(c *gin.Context) {
	var (
		organisms []model.Organism
		err       error
	)

	if q := c.Query("q"); q == "" {
		organisms, err = p.model.Default()
		if err != nil {
			log.Printf("unable to retrieve organisms: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to retrieve organisms"})
			c.Abort()
			return
		}
	} else {
		organisms, err = p.model.Search(q)
		if err != nil {
			log.Printf("unable to retrieve organisms: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to retrieve organisms"})
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, organisms)
}
