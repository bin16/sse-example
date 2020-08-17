package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/bin16/sse-demo/channel"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	})
	g.GET("/chat/:id", func(c *gin.Context) {
		u := uuid.Must(uuid.NewRandom())

		c.SetCookie("secret", u.String(), 666666666, "/", "", false, true)
		c.File("static/chat.html")
	})
	g.GET("/chat/:id/messages", func(c *gin.Context) {
		id := c.Param("id")
		sid, _ := c.Cookie("secret")
		ch := channel.Subscribe(id, sid)
		defer channel.UnSubscribe(id, sid)
		c.Stream(func(w io.Writer) bool {
			select {
			case m := <-ch:
				c.SSEvent("message", m)
			}
			return true
		})
		fmt.Println(id, sid, "messages-bottom")
	})
	g.POST("/chat/:id/messages/:text", func(c *gin.Context) {
		id := c.Param("id")
		message := c.Param("text")
		channel.Post(id, message)
		c.String(http.StatusOK, "ok")
	})
	g.Run(":2333")
}
