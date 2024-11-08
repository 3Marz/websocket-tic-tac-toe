package main

import (

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")
	r.Static("/public", "public")

	lobby := newLobby()
	go lobby.run()

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		serveWs(lobby, "X", c.Writer, c.Request)
	})

	r.Run(":5050")
}
