package api

import (
	"net/http"

	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/config"
	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/inventory"
	"github.com/gin-gonic/gin"
)

func Initialize() {
	r := gin.Default()
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
	if host.Hostname == "" {
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
	name, error := inventory.GetTarget(host.Hostname)
	if error == nil && name.Hostname != "" {
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
	}
	err = inventory.UpdateTarget(host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while updating Host: \n"+err.Error())
	}
	c.JSON(http.StatusAccepted, "OK")
}

func deleteHost(c *gin.Context) {
	target := c.Param("id")
	err, found := inventory.RemoveTarget(target)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while updating Host: \n"+err.Error())
	}
	if found == false {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while updating Host: No such Host.")
	}
	c.JSON(http.StatusAccepted, "OK")
}
