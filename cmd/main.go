package main

import (

	"github.com/gin-gonic/gin"
)

// type Templates struct {
// 	templates *template.Template
// }
//
// func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }
//
// func newTemplate() *Templates {
// 	return &Templates{
// 		templates: template.Must(template.ParseGlob("views/*.html")),
// 	}
// }

// func CORSMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//         c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//         c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
//         c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
//
//         if c.Request.Method == "OPTIONS" {
//             c.AbortWithStatus(204)
//             return
//         }
//
//         c.Next()
//     }
// }



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

// e := echo.New()
// e.Use(middleware.Logger())
// // e.Use(middleware.Recover())
//
// e.Static("/public", "public")
//
// hub := newHub()
// go hub.run()
//
//
// // e.GET("/game", func(c echo.Context) error {
// // 	return c.Render(200, "index.html", nil)
// // })
//
// e.GET("/ws", func(c echo.Context) error {
// 	fmt.Println("ws")
// 	serveWs(hub, "testPlayer", c.Response(), c.Request())
// 	return nil
// })
//
// e.File("/game", "views/game.html")
// // e.GET("/", func(c echo.Context) error {
// // })
//
//
// e.Logger.Fatal(e.Start(":3030"))
