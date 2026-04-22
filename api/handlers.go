package api

import (
	"net/http"
	"time"

	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/config"
	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/inventory"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Initialize() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/hosts/:id", getHost)       // Get a host
	r.POST("/hosts", createHost)       // Create a host
	r.PUT("/hosts/:id", updateHost)    // Update a host
	r.DELETE("/hosts/:id", deleteHost) // Delete a host
	r.Run(":8010")
}

func getHost(c *gin.Context) {
	target := c.Param("id")
	host, err := inventory.GetTarget(target)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while retrieving target Host: \n"+err.Error())
		return
	}
	if len(host.Hostname) == 0 {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while retrieving target Host: No such Host.")
		return
	}
	c.JSON(http.StatusOK, host)
}

func createHost(c *gin.Context) {
	var host config.Host
	err := c.ShouldBindBodyWithJSON(&host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while parsing JSON Body: \n"+err.Error())
		return
	}
	name, error := inventory.GetTarget(host.Hostname[0])
	if error == nil && len(name.Hostname) > 0 {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while adding Host: Host already exists")
		return
	}
	err = inventory.AddTarget(host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while adding Host: \n"+err.Error())
		return
	}
	c.JSON(http.StatusAccepted, "OK")
}

func updateHost(c *gin.Context) {
	var host config.Host
	err := c.ShouldBindBodyWithJSON(&host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while parsing JSON Body: \n"+err.Error())
		return
	}
	err = inventory.UpdateTarget(host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while updating Host: \n"+err.Error())
		return
	}
	c.JSON(http.StatusAccepted, "OK")
}

func deleteHost(c *gin.Context) {
	target := c.Param("id")
	err, found := inventory.RemoveTarget(target)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while updating Host: \n"+err.Error())
		return
	}
	if found == false {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while updating Host: No such Host.")
		return
	}
	c.JSON(http.StatusAccepted, "OK")
}
