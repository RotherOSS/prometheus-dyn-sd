package api

import (
	"net/http"
	"time"

	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/config"
	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/inventory"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Conf *inventory.Config
}

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

	server := &Server{Conf: &inventory.Config{Filepath: config.Filepath}}

	r.GET("/hosts/:id", server.getHost)       // Get a host
	r.POST("/hosts", server.createHost)       // Create a host
	r.PUT("/hosts/:id", server.updateHost)    // Update a host
	r.DELETE("/hosts/:id", server.deleteHost) // Delete a host
	r.Run(":8010")
}

func (server *Server) getHost(c *gin.Context) {
	target := c.Param("id")
	host, err := server.Conf.GetTarget(target)
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

func (server *Server) createHost(c *gin.Context) {
	var host config.Host
	err := c.ShouldBindBodyWithJSON(&host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while parsing JSON Body: \n"+err.Error())
		return
	}
	name, error := server.Conf.GetTarget(host.Hostname[0])
	if error == nil && len(name.Hostname) > 0 {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while adding Host: Host already exists")
		return
	}
	err = server.Conf.AddTarget(host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while adding Host: \n"+err.Error())
		return
	}
	c.JSON(http.StatusAccepted, "OK")
}

func (server *Server) updateHost(c *gin.Context) {
	var host config.Host
	err := c.ShouldBindBodyWithJSON(&host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while parsing JSON Body: \n"+err.Error())
		return
	}
	err = server.Conf.UpdateTarget(host)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Request.Method+" Request failed while updating Host: \n"+err.Error())
		return
	}
	c.JSON(http.StatusAccepted, "OK")
}

func (server *Server) deleteHost(c *gin.Context) {
	target := c.Param("id")
	err, found := server.Conf.RemoveTarget(target)
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
